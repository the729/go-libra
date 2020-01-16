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
	code := []byte{76, 73, 66, 82, 65, 86, 77, 10, 1, 0, 8, 1, 83, 0, 0, 0, 8, 0, 0, 0, 2, 91, 0, 0, 0, 4, 0, 0, 0, 3, 95, 0, 0, 0, 12, 0, 0, 0, 13, 107, 0, 0, 0, 31, 0, 0, 0, 14, 138, 0, 0, 0, 8, 0, 0, 0, 5, 146, 0, 0, 0, 78, 0, 0, 0, 4, 224, 0, 0, 0, 64, 0, 0, 0, 8, 32, 1, 0, 0, 24, 0, 0, 0, 0, 0, 1, 1, 1, 2, 0, 3, 2, 5, 1, 0, 0, 4, 0, 1, 6, 1, 3, 7, 2, 1, 8, 3, 2, 0, 1, 2, 0, 2, 1, 7, 0, 0, 1, 2, 0, 2, 1, 7, 0, 0, 1, 7, 0, 0, 0, 2, 0, 2, 4, 7, 0, 0, 0, 3, 2, 2, 7, 0, 0, 3, 0, 6, 60, 83, 69, 76, 70, 62, 12, 76, 105, 98, 114, 97, 65, 99, 99, 111, 117, 110, 116, 9, 76, 105, 98, 114, 97, 67, 111, 105, 110, 8, 77, 121, 77, 111, 100, 117, 108, 101, 4, 109, 97, 105, 110, 1, 84, 20, 119, 105, 116, 104, 100, 114, 97, 119, 95, 102, 114, 111, 109, 95, 115, 101, 110, 100, 101, 114, 2, 105, 100, 7, 100, 101, 112, 111, 115, 105, 116, 48, 10, 22, 130, 119, 45, 13, 237, 76, 255, 173, 135, 90, 62, 132, 104, 125, 111, 150, 46, 136, 86, 138, 138, 73, 221, 192, 110, 59, 123, 7, 218, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 2, 0, 8, 0, 12, 0, 19, 1, 1, 13, 1, 45, 12, 1, 19, 2, 1, 19, 3, 1, 2}
	args := []types.TransactionArgument{
		types.TxnArgU64(10),
	}

	log.Printf("Get current account sequence of sender...")
	seq, err := c.QueryAccountSequenceNumber(context.TODO(), senderAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("... is %d", seq)

	rawTxn, err := client.NewRawCustomScriptTransaction(
		senderAddr, seq,
		maxGasAmount, gasUnitPrice, expiration, code, args,
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
