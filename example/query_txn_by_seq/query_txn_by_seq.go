package main

import (
	"encoding/hex"
	"log"

	"github.com/the729/go-libra/client"
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

	addrStr := "18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a"
	addr := client.MustToAddress(addrStr)

	provenTxn, err := c.QueryTransactionByAccountSeq(addr, 0, true)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Txn #%d:", provenTxn.GetVersion())
	rawTxn, _ := provenTxn.GetSignedTxn().UnmarshalRawTransaction()
	log.Printf("    Raw txn:")
	log.Printf("        Sender account: %v", hex.EncodeToString(rawTxn.SenderAccount))
	log.Printf("        Sender seq #%v", rawTxn.SequenceNumber)
	log.Printf("        Program: %v...", hex.EncodeToString(rawTxn.GetProgram().Code[:30]))
	log.Printf("        Arg 0: %v", hex.EncodeToString(rawTxn.GetProgram().Arguments[0].Data))
	log.Printf("        Arg 1: %v", hex.EncodeToString(rawTxn.GetProgram().Arguments[1].Data))
	log.Printf("    Gas used: %v", provenTxn.GetGasUsed())
	if provenTxn.GetWithEvents() {
		log.Printf("    Events: (%d total)", len(provenTxn.GetEvents()))
		for idx, ev := range provenTxn.GetEvents() {
			log.Printf("      #%d:", idx)
			log.Printf("        Seq #%d", ev.SequenceNumber)
			log.Printf("        Addr: %v", hex.EncodeToString(ev.AccessPath.Address))
			log.Printf("        Raw path: %v", hex.EncodeToString(ev.AccessPath.Path))
			log.Printf("        Raw data: %v", hex.EncodeToString(ev.Data))
		}
	} else {
		log.Printf("    Events not present")
	}
}