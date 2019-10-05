package main

import (
	"context"
	"log"

	"github.com/the729/go-libra/client"
	"github.com/the729/go-libra/example/utils"
)

const (
	defaultServer    = "ac.testnet.libra.org:8000"
	trustedPeersFile = "../trusted_peers.config.toml"
)

func main() {
	c, err := client.New(defaultServer, trustedPeersFile)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	provenTxnList, err := c.QueryTransactionRange(context.TODO(), 2397, 2, true)
	if err != nil {
		log.Fatal(err)
	}

	for _, provenTxn := range provenTxnList.GetTransactions() {
		utils.PrintTxn(provenTxn)
	}
}
