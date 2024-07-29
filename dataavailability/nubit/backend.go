package nubit

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"time"

	daTypes "github.com/0xPolygon/cdk-data-availability/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rollkit/go-da"
	"github.com/rollkit/go-da/proxy"
)

// NubitDABackend implements the DA integration with Nubit DA layer
type NubitDABackend struct {
	client      da.DA
	config      *Config
	ns          da.Namespace
	privKey     *ecdsa.PrivateKey
	commitTime  time.Time
	batchesData [][]byte
}

// NewNubitDABackend is the factory method to create a new instance of NubitDABackend
func NewNubitDABackend(
	cfg *Config,
	privKey *ecdsa.PrivateKey,
) (*NubitDABackend, error) {
	log.Infof("NubitDABackend config: %#v ", cfg)
	cn, err := proxy.NewClient(cfg.NubitRpcURL, cfg.NubitAuthKey)
	if err != nil {
		return nil, err
	}
	// TODO: Check if name byte array requires zero padding
	name, err := hex.DecodeString(cfg.NubitNamespace)
	if err != nil {
		log.Errorf("error decoding NubitDA namespace config: %+v", err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	log.Infof("NubitDABackend namespace: %s ", string(name))

	return &NubitDABackend{
		config:      cfg,
		privKey:     privKey,
		ns:          name,
		client:      cn,
		commitTime:  time.Now(),
		batchesData: [][]byte{},
	}, nil
}

// Init initializes the NubitDA backend
func (backend *NubitDABackend) Init() error {
	return nil
}

// PostSequence sends the sequence data to the data availability backend, and returns the dataAvailabilityMessage
// as expected by the contract
func (backend *NubitDABackend) PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error) {
	// Add to batches data cache
	backend.batchesData = append(backend.batchesData, batchesData...)
	batchesSize := uint64(len(backend.batchesData))
	if batchesSize < backend.config.NubitMaxBatchesSize {
		log.Infof("Added batches data to NubitDABackend cache, current length: %+v", batchesSize)
		return nil, nil
	}

	// Commit time interval validation
	lastCommitTime := time.Since(backend.commitTime)
	if lastCommitTime < NubitMinCommitTime {
		time.Sleep(NubitMinCommitTime - lastCommitTime)
	}

	// Encode NubitDA blob data
	data := EncodeSequence(backend.batchesData)
	ids, err := backend.client.Submit(ctx, [][]byte{data}, -1, backend.ns)
	// Ensure only a single blob ID returned
	if err != nil || len(ids) != 1 {
		log.Errorf("Submit batch data with NubitDA client failed: %s", err)
		return nil, err
	}
	blobID := ids[0]
	log.Infof("Data submitted to Nubit DA: %d bytes against namespace %v sent with id %#x", batchesSize, backend.ns, blobID)

	// Reset batches data cache and DA commit time
	backend.commitTime = time.Now()
	backend.batchesData = [][]byte{}

	// Get proof
	tries := uint64(0)
	posted := false
	for tries < backend.config.NubitGetProofMaxRetry {
		dataProof, err := backend.client.GetProofs(ctx, [][]byte{blobID}, backend.ns)
		if err != nil {
			log.Infof("Proof not available: %s", err)
		}
		if len(dataProof) == 1 {
			// TODO: add data proof to DA message
			log.Infof("Data proof from Nubit DA received: %+v", dataProof)
			posted = true
			break
		}

		// Retries
		tries += 1
		time.Sleep(backend.config.NubitGetProofWaitPeriod)
	}
	if !posted {
		log.Errorf("Get blob proof on Nubit DA failed: %s", err)
		return nil, err
	}

	// Get abi-encoded data availability message
	sequence := daTypes.Sequence{}
	for _, seq := range batchesData {
		sequence = append(sequence, seq)
	}
	signedSequence, err := sequence.Sign(backend.privKey)
	if err != nil {
		log.Errorf("Failed to sign sequence with pk: %v", err)
		return nil, err
	}
	signature := append(sequence.HashToSign(), signedSequence.Signature...)
	blobData := BlobData{
		BlobID:    blobID,
		Signature: signature,
	}

	return TryEncodeToDataAvailabilityMessage(blobData)
}

// GetSequence gets the sequence data from NubitDA layer
func (backend *NubitDABackend) GetSequence(ctx context.Context, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error) {
	blobData, err := TryDecodeFromDataAvailabilityMessage(dataAvailabilityMessage)
	if err != nil {
		log.Error("Error decoding from da message: ", err)
		return nil, err
	}

	reply, err := backend.client.Get(ctx, [][]byte{blobData.BlobID}, backend.ns)
	if err != nil || len(reply) != 1 {
		log.Error("Error retrieving blob from NubitDA client: ", err)
		return nil, err
	}

	batchesData, _ := DecodeSequence(reply[0])
	return batchesData, nil
}
