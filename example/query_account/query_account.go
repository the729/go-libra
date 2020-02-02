package main

import (
	"context"
	"encoding/hex"
	"log"

	"github.com/the729/go-libra/client"
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

	provenState, err := c.QueryAccountState(context.TODO(), addr)
	if err != nil {
		log.Fatal(err)
	}

	if provenState.IsNil() {
		log.Printf("Account %s does not exist at version %d.", addrStr, provenState.GetVersion())
		return
	}

	resource, err := provenState.GetAccountBlob().GetLibraAccountResource()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Balance (microLibra): %d", resource.GetBalance())
	log.Printf("Sequence Number: %d", resource.GetSequenceNumber())
	log.Printf("SentEventsCount: %d", resource.GetSentEvents().Count)
	log.Printf("    Key: %x", resource.GetSentEvents().Key)
	log.Printf("ReceivedEventsCount: %d", resource.GetReceivedEvents().Count)
	log.Printf("    Key: %x", resource.GetReceivedEvents().Key)
	log.Printf("DelegatedWithdrawalCapability: %v", resource.GetDelegatedWithdrawalCapability())
	log.Printf("Authentication key: %v", hex.EncodeToString(resource.GetAuthenticationKey()))
	log.Printf("Event generator: %v", resource.GetEventGenerator())
}
