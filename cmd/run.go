package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sieniven/zkevm-nubit/config"
	"github.com/sieniven/zkevm-nubit/dataavailability"
	"github.com/sieniven/zkevm-nubit/dataavailability/nubit"
	"github.com/sieniven/zkevm-nubit/etherman"
	"github.com/sieniven/zkevm-nubit/ethtxmanager"
	"github.com/sieniven/zkevm-nubit/log"
	"github.com/sieniven/zkevm-nubit/sequencesender"
	"github.com/urfave/cli/v2"
)

func start(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		return err
	}
	setupLog(c.Log)

	// Initialize ether manager instance
	etherMan, err := newEtherman(*c)
	if err != nil {
		panic(err)
	}

	// Create new data avaiability backend and set sequencer flag to false
	isSequencer := false
	_, pk, err := etherMan.LoadAuthFromKeyStore(c.SequenceSender.DAPermitApiPrivateKey.Path, c.SequenceSender.DAPermitApiPrivateKey.Password)
	if err != nil {
		return err
	}

	log.Infof("from pk %s", crypto.PubkeyToAddress(pk.PublicKey))
	daBackend, err := nubit.NewNubitDABackend(&c.DataAvailability, pk)
	if err != nil {
		return err
	}

	da, err := dataavailability.New(isSequencer, daBackend)
	etherMan.SetDataProvider(da)

	// Initialize eth tx manager instance
	etm := ethtxmanager.New(c.EthTxManager, etherMan)

	// Initialize mock sequence sender
	seqSender := createMockSequenceSender(*c, etm, etherMan)

	// Start mock sequence sender
	go seqSender.Start(cliCtx.Context)

	// Start send sequence flag handler
	reader := bufio.NewReader(os.Stdin)
	seqSender.SendSequenceHandle(cliCtx.Context, reader)
	return nil
}

// createMockSequenceSender is the mock function for PolygonCDK node that
// creates a new instance of the mock sequence sender for the mock node.
func createMockSequenceSender(cfg config.Config, etm *ethtxmanager.Client, etherMan *etherman.Client) *sequencesender.SequenceSender {
	// Create new data avaiability manager
	da, err := newDataAvailability(cfg, etherMan)
	if err != nil {
		panic(err)
	}

	etherMan.SetDataProvider(da)
	_, privKey, err := etherMan.LoadAuthFromKeyStore(cfg.Key.Path, cfg.Key.Password)
	if err != nil {
		panic(err)
	}
	if cfg.SequenceSender.SenderAddress.Cmp(common.Address{}) == 0 {
		panic(errors.New("sequence sender address not found"))
	}
	if privKey == nil { //nolint:staticcheck
		panic(errors.New("private key not found"))
	}
	fmt.Printf("from pk %s, from sender %s\n", crypto.PubkeyToAddress(privKey.PublicKey), cfg.SequenceSender.SenderAddress.String()) //nolint:staticcheck

	// Initialize new sequence sender instance
	seqSender, err := sequencesender.New(cfg.SequenceSender, etherMan, etm)
	if err != nil {
		panic(err)
	}
	seqSender.SetDataProvider(da)
	return seqSender
}

func newEtherman(c config.Config) (*etherman.Client, error) {
	return etherman.NewClient(c.Etherman, c.L1Config)
}

func newDataAvailability(c config.Config, etherMan *etherman.Client) (*dataavailability.DataAvailability, error) {
	isSequencer := false
	_, pk, err := etherMan.LoadAuthFromKeyStore(c.SequenceSender.DAPermitApiPrivateKey.Path, c.SequenceSender.DAPermitApiPrivateKey.Password)
	if err != nil {
		return nil, err
	}

	log.Infof("from pk %s", crypto.PubkeyToAddress(pk.PublicKey))
	daBackend, err := nubit.NewNubitDABackend(&c.DataAvailability, pk)
	if err != nil {
		return nil, err
	}

	return dataavailability.New(isSequencer, daBackend)
}

func setupLog(c log.Config) {
	log.Init(c)
}
