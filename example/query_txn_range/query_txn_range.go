package main

import (
	"context"
	"log"

	"github.com/the729/go-libra/client"
	"github.com/the729/go-libra/example/utils"
)

const (
	defaultServer = "ac.testnet.libra.org:8000"
	waypoint      = "0:bf7e1eef81af68cc6b4801c3739da6029c778a72e67118a8adf0dd759f188908"
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
