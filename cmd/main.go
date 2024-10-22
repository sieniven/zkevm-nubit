package main

import (
	"fmt"
	"os"

	"github.com/sieniven/zkevm-nubit/config"
	"github.com/urfave/cli/v2"
)

const appName = "node"

const (
	// SEQUENCER is the sequencer component identifier
	SEQUENCER = "sequencer"
	// ETHTXMANAGER is the service that managers the tx sent to the L1
	ETHTXMANAGER = "eth-tx-manager"
	// SEQUENCE_SENDER is the sequence sender component identifier
	SEQUENCE_SENDER = "sequence-sender"
)

const (
	// NODE_CONFIGFILE name to identify the node config-file
	NODE_CONFIGFILE = "zkevm-node"
)

var configFileFlag = cli.StringFlag{
	Name:     config.FlagCfg,
	Aliases:  []string{"c"},
	Usage:    "Configuration `FILE`",
	Required: true,
}

var networkJsonFlag = cli.BoolFlag{
	Name:     config.FlagNetwork,
	Aliases:  []string{"net"},
	Usage:    "Use JSON network configuration",
	Required: false,
	Value:    false,
}

var toFlag = cli.StringFlag{
	Name:     config.FlagTo,
	Aliases:  []string{"toaddress"},
	Usage:    "receiving address to send test eth",
	Required: false,
	Value:    "",
}

var passwordFlag = cli.StringFlag{
	Name:     config.FlagPassword,
	Aliases:  []string{"p"},
	Usage:    "set keystore password",
	Required: false,
	Value:    "password",
}

func main() {
	app := cli.NewApp()
	app.Name = appName
	flags := []cli.Flag{
		&configFileFlag,
		&networkJsonFlag,
		&toFlag,
		&passwordFlag,
	}
	app.Commands = []*cli.Command{
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   "Run the mock zkevm-node",
			Action:  start,
			Flags:   flags,
		},
		{
			Name:    "create-keystore",
			Aliases: []string{},
			Usage:   "Create new eth keystore",
			Action:  createKeystore,
			Flags:   flags,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
