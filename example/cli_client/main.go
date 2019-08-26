package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	defaultServer    = "ac.testnet.libra.org:8000"
	trustedPeersFile = "../trusted_peers.config.toml"
	walletFile       = "wallet.toml"
)

var ServerAddr, TrustedPeersFile, WalletFile string

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "server",
			Value:       defaultServer,
			Usage:       "use Libra server `HOST:PORT`",
			Destination: &ServerAddr,
		},
		cli.StringFlag{
			Name:        "peers, s",
			Value:       trustedPeersFile,
			Usage:       "load trusted peers from `FILE`",
			Destination: &TrustedPeersFile,
		},
		cli.StringFlag{
			Name:        "wallet, w",
			Value:       walletFile,
			Usage:       "load or store account private keys in `FILE`",
			Destination: &WalletFile,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "query",
			Aliases: []string{"q"},
			Subcommands: []cli.Command{
				{
					Name:    "ledger_info",
					Usage:   "known_version",
					Aliases: []string{"l"},
					Action:  cmdQueryLedgerInfo,
				},
				{
					Name:    "account_state",
					Usage:   "address_prefix",
					Aliases: []string{"as"},
					Action:  cmdQueryAccountState,
				},
				{
					Name:    "transaction_range",
					Usage:   "start limit",
					Aliases: []string{"tr"},
					Action:  cmdQueryTransactionRange,
				},
				{
					Name:    "transaction_by_seq",
					Usage:   "address_prefix sequence",
					Aliases: []string{"ts"},
					Action:  cmdQueryTransactionByAccountSeq,
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
					Action:  cmdCreateAccounts,
				},
				{
					Name:    "list",
					Aliases: []string{"ls"},
					Action:  cmdListAccounts,
				},
				{
					Name:    "mint",
					Usage:   "address_prefix amount",
					Aliases: []string{"m"},
					Action:  cmdMint,
				},
			},
		},
		{
			Name:    "transfer",
			Usage:   "sender_address_prefix receiver_address_prefix amount [max_gas_amount [gas_unit_price_micro [expiration_seconds]]]",
			Aliases: []string{"t"},
			Action:  cmdTransfer,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
