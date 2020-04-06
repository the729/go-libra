package main

import (
	"context"
	"log"

	"github.com/the729/go-libra/client"
	"github.com/the729/go-libra/example/utils"
)

const (
	defaultServer = "ac.testnet.libra.org:8000"
	waypoint      = "0:a69511cc7e6d609efcf03e64098056bc3c96d383e0adcf752464a111b081b808"
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
