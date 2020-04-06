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

	addrStr := "42f5745128c05452a0c68272de8042b1"
	addr := client.MustToAddress(addrStr)

	provenTxn, err := c.QueryTransactionByAccountSeq(context.TODO(), addr, 0, true)
	if err != nil {
		log.Fatal(err)
	}

	utils.PrintTxn(provenTxn)
}
