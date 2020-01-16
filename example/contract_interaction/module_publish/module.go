package main

import (
	"context"
	"encoding/hex"
	"github.com/the729/go-libra/types"
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
	maxGasAmount := uint64(140000)
	gasUnitPrice := uint64(0)
	expiration := time.Now().Add(1 * time.Minute)

	// Contract parameters
	code := []byte{76, 73, 66, 82, 65, 86, 77, 10, 1, 0, 8, 1, 83, 0, 0, 0, 4, 0, 0, 0, 2, 87, 0, 0, 0, 4, 0, 0, 0, 3, 91, 0, 0, 0, 3, 0, 0, 0, 13, 94, 0, 0, 0, 10, 0, 0, 0, 14, 104, 0, 0, 0, 5, 0, 0, 0, 5, 109, 0, 0, 0, 24, 0, 0, 0, 4, 133, 0, 0, 0, 64, 0, 0, 0, 11, 197, 0, 0, 0, 10, 0, 0, 0, 0, 0, 1, 1, 1, 2, 1, 0, 0, 3, 0, 2, 1, 7, 0, 0, 1, 7, 0, 0, 0, 3, 1, 7, 0, 0, 8, 77, 121, 77, 111, 100, 117, 108, 101, 9, 76, 105, 98, 114, 97, 67, 111, 105, 110, 1, 84, 2, 105, 100, 234, 131, 65, 108, 43, 73, 164, 206, 220, 184, 161, 170, 140, 26, 16, 133, 232, 205, 32, 81, 132, 88, 130, 252, 229, 76, 227, 54, 246, 174, 11, 235, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 2, 0, 12, 0, 2}

	log.Printf("Get current account sequence of sender...")
	seq, err := c.QueryAccountSequenceNumber(context.TODO(), senderAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("... is %d", seq)

	rawTxn, err := client.NewRawCustomModuleTransaction(
		senderAddr, seq,
		maxGasAmount, gasUnitPrice, expiration, code,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Submit transaction...")
	expectedSeq, err := c.SubmitRawTransaction(context.TODO(), rawTxn, priKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(expectedSeq)

	log.Printf("Waiting until transaction is included in ledger...")
	err = c.PollSequenceUntil(context.TODO(), senderAddr, expectedSeq, expiration)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("done.")
}
