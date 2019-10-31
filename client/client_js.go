// +build js

/*
Package client implements a gRPC client to Libra RPC service.

Features include:
  - Query ledger information
  - Query account state
  - Query transactions by range
  - Query single transaction by account and sequence number
  - Sign and submit raw transactions

All queries are cryptographically verified to proof their inclusion and integrity in the blockchain.
*/
package client

import (
	"fmt"

	"github.com/the729/go-libra/config"
	"github.com/the729/go-libra/generated/pbac"
	"github.com/the729/go-libra/types/validator"
)

// Client is a Libra client.
// It has a gRPC client to a Libra RPC server, with public keys to trusted peers.
type Client struct {
	ac       pbac.AdmissionControlClient
	verifier validator.Verifier
}

// New creates a new Libra Client.
func New(ServerAddr, TrustedPeerData string) (*Client, error) {
	c := &Client{}
	if err := c.loadTrustedPeers(TrustedPeerData); err != nil {
		return nil, err
	}
	c.ac = pbac.NewAdmissionControlClient(ServerAddr)
	return c, nil
}

// Close the client (do nothing).
func (c *Client) Close() {}

func (c *Client) loadTrustedPeers(tomlData string) error {
	peerconf, err := config.LoadTrustedPeers(tomlData)
	if err != nil {
		return fmt.Errorf("load conf err: %v", err)
	}
	verifier, err := validator.NewConsensusVerifier(peerconf)
	if err != nil {
		return fmt.Errorf("new verifier err: %v", err)
	}
	c.verifier = verifier
	return nil
}
