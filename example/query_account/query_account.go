package main

import (
	"context"
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

	provenState, err := c.QueryAccountState(context.TODO(), addr)
	if err != nil {
		log.Fatal(err)
	}

	if provenState.IsNil() {
		log.Printf("Account %s does not exist at version %d.", addrStr, provenState.GetVersion())
		return
	}

	provenResource, err := c.GetLibraCoinResourceFromAccountBlob(provenState.GetAccountBlob())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Balance (microLibra): %d", provenResource.GetBalance())
	log.Printf("Sequence Number: %d", provenResource.GetSequenceNumber())
	log.Printf("SentEventsCount: %d", provenResource.GetSentEvents().Count)
	log.Printf("    Key: %x", provenResource.GetSentEvents().Key)
	log.Printf("ReceivedEventsCount: %d", provenResource.GetReceivedEvents().Count)
	log.Printf("    Key: %x", provenResource.GetReceivedEvents().Key)
	log.Printf("DelegatedWithdrawalCapability: %v", provenResource.GetDelegatedWithdrawalCapability())
}
