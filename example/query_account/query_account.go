package main

import (
	"context"
	"encoding/hex"
	"log"

	"github.com/the729/go-libra/client"
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

	provenState, err := c.QueryAccountState(context.TODO(), addr)
	if err != nil {
		log.Fatal(err)
	}

	if provenState.IsNil() {
		log.Printf("Account %s does not exist at version %d.", addrStr, provenState.GetVersion())
		return
	}

	ar, br, err := provenState.GetAccountBlob().GetLibraResources()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Balance: %d", br.Coin)
	log.Printf("Sequence Number: %d", ar.SequenceNumber)
	log.Printf("SentEventsCount: %d", ar.SentEvents.Count)
	log.Printf("    Key: %x", ar.SentEvents.Key)
	log.Printf("ReceivedEventsCount: %d", ar.ReceivedEvents.Count)
	log.Printf("    Key: %x", ar.ReceivedEvents.Key)
	log.Printf("DelegatedWithdrawalCapability: %v", ar.DelegatedWithdrawalCapability)
	log.Printf("Authentication key: %v", hex.EncodeToString(ar.AuthenticationKey))
	log.Printf("Event generator: %v", ar.EventGenerator)
}
