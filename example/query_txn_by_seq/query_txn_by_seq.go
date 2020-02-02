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

	addrStr := "18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a"
	addr := client.MustToAddress(addrStr)

	provenTxn, err := c.QueryTransactionByAccountSeq(context.TODO(), addr, 0, true)
	if err != nil {
		log.Fatal(err)
	}

	utils.PrintTxn(provenTxn)
}
