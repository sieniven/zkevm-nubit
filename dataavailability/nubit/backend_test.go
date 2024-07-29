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

func TestBackendSendDataOnChain(t *testing.T) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	nubit, err := NewMockNubitDABackend(
		"http://127.0.0.1:26658",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.DAMv0s7915Ahx-kDFSzDT1ATz4Q9WwktWcHmjp7_99Q",
		pk,
	)
	if err != nil {
		t.Fatal(err)
	}

	testMsg := "test txs"
	mockSequence := [][]byte{}
	for i := 0; i < 10; i++ {
		mockSequence = append(mockSequence, []byte(testMsg))
	}
	daMessage, err := nubit.PostSequence(context.Background(), mockSequence)
	require.NoError(t, err)
	time.Sleep(600 * time.Millisecond)

	returnData, err := nubit.GetSequence(context.Background(), []common.Hash{}, daMessage)
	require.NoError(t, err)
	require.Equal(t, len(mockSequence), len(returnData))

	for _, data := range returnData {
		require.Equal(t, []byte(testMsg), data)
	}
}

func NewMockNubitDABackend(url string, authKey string, pk *ecdsa.PrivateKey) (*NubitDABackend, error) {
	cn, err := proxy.NewClient(url, authKey)
	if err != nil || cn == nil {
		return nil, err
	}

	name, err := hex.DecodeString("00000000000000000000000000000000000000000000706f6c79676f6e")
	if err != nil {
		return nil, err
	}

	log.Infof("Nubit Namespace: %s ", string(name))
	return &NubitDABackend{
		ns:         name,
		client:     cn,
		privKey:    pk,
		commitTime: time.Now(),
	}, nil
}
