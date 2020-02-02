package main

import (
	"context"
	"log"

	"github.com/the729/go-libra/client"
	"github.com/the729/go-libra/example/utils"
)

const (
	defaultServer = "ac.testnet.libra.org:8000"
	waypoint      = "0:59d76f0cb85470b269e754b54edac5d88619f21234d69e02ce85cea6f003089a"
)

func main() {
	c, err := client.New(defaultServer, waypoint)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	provenTxnList, err := c.QueryTransactionRange(context.TODO(), 8207475, 2, true)
	if err != nil {
		log.Fatal(err)
	}

	for _, provenTxn := range provenTxnList.GetTransactions() {
		utils.PrintTxn(provenTxn)
	}
}
