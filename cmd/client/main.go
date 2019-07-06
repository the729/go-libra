package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/the729/go-libra/client"
)

const (
	defaultServer    = "ac.testnet.libra.org:8000"
	trustedPeersFile = "trusted_peers.config.toml"
	walletFile       = "wallet.toml"
)

func main() {
	c := &client.Client{}

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "server",
			Value:       defaultServer,
			Usage:       "use Libra server `HOST:PORT`",
			Destination: &c.ServerAddr,
		},
		cli.StringFlag{
			Name:        "peers, s",
			Value:       trustedPeersFile,
			Usage:       "load trusted peers from `FILE`",
			Destination: &c.TrustedPeerFile,
		},
		cli.StringFlag{
			Name:        "wallet, w",
			Value:       walletFile,
			Usage:       "load or store account private keys in `FILE`",
			Destination: &c.WalletFile,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "query",
			Aliases: []string{"q"},
			Subcommands: []cli.Command{
				{
					Name:    "account_state",
					Usage:   "address_prefix",
					Aliases: []string{"as"},
					Action:  c.CmdQueryAccountState,
				},
			},
		},
		{
			Name:    "account",
			Aliases: []string{"a"},
			Subcommands: []cli.Command{
				{
					Name:    "create",
					Usage:   "number_of_accounts",
					Aliases: []string{"c"},
					Action:  c.CmdCreateAccounts,
				},
				{
					Name:    "list",
					Aliases: []string{"ls"},
					Action:  c.CmdListAccounts,
				},
				{
					Name:    "mint",
					Usage:   "address_prefix amount",
					Aliases: []string{"m"},
					Action:  c.CmdMint,
				},
			},
		},
		{
			Name:    "transfer",
			Usage:   "sender_address_prefix receiver_address_prefix amount",
			Aliases: []string{"t"},
			Action:  c.CmdTransfer,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
