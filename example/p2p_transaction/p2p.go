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
	defaultServer    = "ac.testnet.libra.org:8000"
	trustedPeersFile = "../consensus_peers.config.toml"
)

func main() {
	c, err := client.New(defaultServer, trustedPeersFile)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// We need private key fo the sender to sign the transaction
	priKeyBytes, _ := hex.DecodeString("657cd8ed5e434cc4f874d6822889f637957f0145c67e2b055c9954c936670a61e57ea705e00e3ecaf417b4285cd0a69b1d79406914581456c1ce278b81a48674")
	priKey := ed25519.PrivateKey(priKeyBytes)

	// Transaction parameters
	senderAddr := client.MustToAddress("18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a")
	recvAddr := client.MustToAddress("e89a0d93fcf1ca4423328c1bddebe6c02da666808993c8a888ff7a8bad19ffd5")
	amountMicro := uint64(2 * 1000000)
	maxGasAmount := uint64(140000)
	gasUnitPrice := uint64(0)
	expiration := time.Now().Add(1 * time.Minute)

	log.Printf("Get current account sequence of sender...")
	seq, err := c.QueryAccountSequenceNumber(context.TODO(), senderAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("... is %d", seq)

	rawTxn, err := client.NewRawP2PTransaction(
		senderAddr, recvAddr, seq,
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
