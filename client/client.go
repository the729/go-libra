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
	"encoding/hex"
	"sync"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbac"
	"github.com/the729/go-libra/types/proof/accumulator"
	"github.com/the729/go-libra/types/validator"
)

// Client is a Libra client.
// It has a gRPC client to a Libra RPC server, with public keys to trusted peers.
type Client struct {
	closeFunc func()
	ac        pbac.AdmissionControlClient
	verifier  validator.Verifier
	acc       *accumulator.Accumulator
	accMu     sync.RWMutex
}

// New creates a new Libra Client.
//
// For normal usage, ServerAddr is in host:port format. TrustedPeer is a TOML file that contains
// the trusted peers.
//
// For use with Javascript, ServerAddr is in http://host:port format. TrustedPeer is a TOML formated
// text of the trusted peers config.
func New(ServerAddr, TrustedPeer string) (*Client, error) {
	c := &Client{}
	if err := c.loadTrustedPeers(TrustedPeer); err != nil {
		return nil, err
	}
	if err := c.connect(ServerAddr); err != nil {
		return nil, err
	}

	genesisHash, _ := hex.DecodeString("b1f2c172f22b8a9e7fd89a64f75b3b10d64431ccbb1f89a00aff5725e4284fb1")
	c.acc = &accumulator.Accumulator{
		Hasher:             sha3libra.NewTransactionAccumulator(),
		FrozenSubtreeRoots: [][]byte{genesisHash},
		NumLeaves:          1,
	}

	return c, nil
}

// Close the client.
func (c *Client) Close() {
	if c.closeFunc != nil {
		c.closeFunc()
	}
}
