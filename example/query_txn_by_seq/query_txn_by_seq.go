package main

import (
	"context"
	"log"

	"github.com/the729/go-libra/client"
	"github.com/the729/go-libra/example/utils"
)

const (
	defaultServer = "ac.testnet.libra.org:8000"
	waypoint      = "0:4d4d0feaa9378069f8fcee71980e142273837e108702d8d7f93a8419e2736f3f"
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
