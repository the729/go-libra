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

	"google.golang.org/grpc"

	"github.com/the729/go-libra/config"
	"github.com/the729/go-libra/generated/pbac"
	"github.com/the729/go-libra/types/validator"
)

// Client is a Libra client.
// It has a gRPC client to a Libra RPC server, with public keys to trusted peers.
type Client struct {
	conn     *grpc.ClientConn
	ac       pbac.AdmissionControlClient
	verifier validator.Verifier
}

// New creates a new Libra Client.
func New(ServerAddr, TrustedPeerFile string) (*Client, error) {
	c := &Client{}
	if err := c.loadTrustedPeers(TrustedPeerFile); err != nil {
		return nil, err
	}
	if err := c.connect(ServerAddr); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) connect(server string) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("grpc dial error: %v", err)
	}

	acClient := pbac.NewAdmissionControlClient(conn)
	c.conn = conn
	c.ac = acClient
	return nil
}

// Close the client.
func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) loadTrustedPeers(file string) error {
	peerconf, err := config.LoadTrustedPeersFromFile(file)
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
