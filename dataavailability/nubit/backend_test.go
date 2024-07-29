package nubit

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollkit/go-da/proxy"
	"github.com/sieniven/zkevm-nubit/log"
	"github.com/stretchr/testify/require"
)

func TestOffchainPipeline(t *testing.T) {
	cfg := Config{
		NubitRpcURL: "http://127.0.0.1:26658",
		NubitAuthKey: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.DAMv0s7915Ahx-kDFSzDT1ATz4Q9WwktWcHmjp7_99Q"
		NubitNamespace: "xlayer",
		NubitMaxBatchesSize: "102400",
		NubitGetProofMaxRetry: "10",
		NubitGetProofWaitPeriod: "5s",
	}
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	backend, err := NewNubitDABackend(&cfg, pk)
	require.NoError(t, err)

	// Generate mock string batch data
	stringData := "hihihihihihihihihihihihihihihihihihi"
	data := []byte(stringData)

	// Generate mock string sequence
	mockBatches := [][]byte{}
	for i := 0; i < 10; i++ {
		mockBatches = append(mockBatches, data)
	}

	msg, err := backend.PostSequence(context.Background(), mockBatches)
	fmt.Println("DA msg: ", msg)
	require.NoError(t, err)
	time.Sleep(600 * time.Millisecond)

	blobData, err := TryDecodeFromDataAvailabilityMessage(msg)
	require.NoError(t, err)
	require.NotNil(blobData.BlobID)
	require.NotNil(blobData.Signature)
	require.NotZero(len(blobData.BlobID))
	require.NotZero(len(blobData.Signature))
	fmt.Println("Decoding DA msg successful")

	// Retrieve sequence with provider
	returnData, err := backend.GetSequence(context.Background(), []common.Hash{}, msg)

	// Validate retrieved data
	require.NoError(t, err)
	require.Equal(t, 10, len(returnData))
	for _, batchData := range returnData {
		assert.Equal(t, stringData, string(batchData))
	}
}

func TestOffchainPipelineWithRandomData(t *testing.T) {
	cfg := Config{
		NubitRpcURL: "http://127.0.0.1:26658",
		NubitAuthKey: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.DAMv0s7915Ahx-kDFSzDT1ATz4Q9WwktWcHmjp7_99Q"
		NubitNamespace: "xlayer",
		NubitMaxBatchesSize: "102400",
		NubitGetProofMaxRetry: "10",
		NubitGetProofWaitPeriod: "5s",
	}
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	backend, err := NewNubitDABackend(&cfg, pk)
	require.NoError(t, err)

	// Define Different DataSizes
	dataSize := []int{100000, 200000, 1000, 80, 30000}

	// Disperse Blob with different DataSizes
	rand.Seed(time.Now().UnixNano())                         //nolint:gosec,staticcheck
	data := make([]byte, dataSize[rand.Intn(len(dataSize))]) //nolint:gosec,staticcheck
	_, err := rand.Read(data)                                //nolint:gosec,staticcheck
	assert.NoError(t, err)

	// Generate mock string sequence
	mockBatches := [][]byte{}
	for i := 0; i < 10; i++ {
		mockBatches = append(mockBatches, data)
	}

	msg, err := backend.PostSequence(context.Background(), mockBatches)
	fmt.Println("DA msg: ", msg)
	require.NoError(t, err)
	time.Sleep(600 * time.Millisecond)

	blobData, err := TryDecodeFromDataAvailabilityMessage(msg)
	require.NoError(t, err)
	require.NotNil(blobData.BlobID)
	require.NotNil(blobData.Signature)
	require.NotZero(len(blobData.BlobID))
	require.NotZero(len(blobData.Signature))
	fmt.Println("Decoding DA msg successful")

	// Retrieve sequence with provider
	returnData, err := backend.GetSequence(context.Background(), []common.Hash{}, msg)

	// Validate retrieved data
	require.NoError(t, err)
	require.Equal(t, 10, len(returnData))
	for idx, batchData := range batchesData {
		assert.Equal(t, mockBatches[idx], batchData)
	}
}
 
func NewMockNubitDABackend(url string, authKey string, pk *ecdsa.PrivateKey) (*NubitDABackend, error) {
	cn, err := proxy.NewClient(url, authKey)
	if err != nil || cn == nil {
		return nil, err
	}

	name, err := hex.DecodeString("xlayer")
	if err != nil {
		return nil, err
	}

	log.Infof("Nubit Namespace: %s ", string(name))
	return &NubitDABackend{
		namespace:  name,
		client:     cn,
		privKey:    pk,
		commitTime: time.Now(),
	}, nil
}
