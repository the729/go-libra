package main

import (
	"context"
	"encoding/hex"
	"log"
	"time"

	"github.com/the729/go-libra/client"
	"golang.org/x/crypto/ed25519"
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

	// We need private key fo the sender to sign the transaction
	priKeyBytes, _ := hex.DecodeString("996911072ee011ffa44a1325e0da593ff3b9374e255115f223cbdffb6bfa0bcfba60d1f8edd6923f59cf9125d3ac80e389afa4e2b8d0e4f1183a30a0270fde71")
	priKey := ed25519.PrivateKey(priKeyBytes)

	// Transaction parameters
	senderAddr := client.MustToAddress("42f5745128c05452a0c68272de8042b1")
	recvAddr := client.MustToAddress("5817cd6e6e84c110c43efca22df54172")
	recvAuthKeyPrefix, _ := hex.DecodeString("26c7bfaa8e0f32206f35bf6d44b43c9c")
	amountMicro := uint64(2 * 1000000)
	maxGasAmount := uint64(500000)
	gasUnitPrice := uint64(0)
	expiration := time.Now().Add(1 * time.Minute)

	log.Printf("Get current account sequence of sender...")
	seq, err := c.QueryAccountSequenceNumber(context.TODO(), senderAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("... is %d", seq)

	rawTxn, err := client.NewRawP2PTransaction(
		senderAddr, recvAddr, recvAuthKeyPrefix, seq,
		amountMicro, maxGasAmount, gasUnitPrice, expiration,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Submit transaction...")
	expectedSeq, err := c.SubmitRawTransaction(context.TODO(), rawTxn, priKey)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Waiting until transaction is included in ledger...")
	err = c.PollSequenceUntil(context.TODO(), senderAddr, expectedSeq, expiration)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("done.")
}
