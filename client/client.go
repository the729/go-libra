/*
Package client implements a gRPC client to Libra RPC service.

Features include:
  - Query ledger information
  - Query account state
  - Query transactions by range
  - Query single transaction by account and sequence number
  - Sign and submit raw transactions

All queries are cryptographically verified to proof their inclusion and integrity in the blockchain.

The client can also keep track of the consistency of the ledger. This function will detect hard forks or
block chain reset.

When a client is newly constructed, it knows only the genesis block (version 0) of the ledger. The hash of
the genesis block is hardcoded. After each query to the ledger, the client updates its knowledge about the latest
version and the Merkle tree accumulator.

You should extract the known-version state of a client instance before destroying it, by calling GetKnownVersion(),
and saving the result somewhere. Later, when a new client instance is constructed, you should use SetKnownVersion()
to restore the known-version state.
*/
package client

import (
	"fmt"
	"sync"

	"github.com/the729/go-libra/generated/pbac"
	"github.com/the729/go-libra/types"
	"github.com/the729/go-libra/types/proof/accumulator"
)

// Client is a Libra client.
// It has a gRPC client to a Libra RPC server, with public keys to trusted peers.
type Client struct {
	closeFunc    func()
	ac           pbac.AdmissionControlClient
	verifier     types.LedgerInfoVerifier
	acc          *accumulator.Accumulator
	accMu        sync.RWMutex
	lastWaypoint string
}

// New creates a new Libra Client.
//
// For normal usage, ServerAddr is in host:port format. TrustedPeer is a TOML file that contains
// the trusted peers.
//
// For use with Javascript, ServerAddr is in http://host:port format. TrustedPeer is a TOML formated
// text of the trusted peers config.
func New(ServerAddr, Waypoint string) (*Client, error) {
	return NewFromState(ServerAddr, &ClientState{Waypoint: Waypoint})
}

func NewFromState(ServerAddr string, state *ClientState) (*Client, error) {
	c := &Client{}
	if err := c.SetState(state); err != nil {
		return nil, fmt.Errorf("invalid state: %v", err)
	}
	if err := c.connect(ServerAddr); err != nil {
		return nil, err
	}
	return c, nil
}

// Close the client.
func (c *Client) Close() {
	if c.closeFunc != nil {
		c.closeFunc()
	}
}

// GetLatestWaypoint returns the latest waypoint string that this client has
// reached.
func (c *Client) GetLatestWaypoint() string {
	return c.lastWaypoint
}
