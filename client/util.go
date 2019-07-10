package client

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/the729/go-libra/config"
	"github.com/the729/go-libra/generated/pbac"
	"github.com/the729/go-libra/types/validator"
)

type Client struct {
	ServerAddr      string
	TrustedPeerFile string
	WalletFile      string

	conn     *grpc.ClientConn
	ac       pbac.AdmissionControlClient
	verifier validator.Verifier
	accounts map[string]*Account
}

func (c *Client) Connect() error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(c.ServerAddr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}

	acClient := pbac.NewAdmissionControlClient(conn)
	c.conn = conn
	c.ac = acClient
	return nil
}

func (c *Client) Disconnect() {
	c.conn.Close()
}

func (c *Client) LoadTrustedPeers() error {
	peerconf, err := config.LoadTrustedPeersFromFile(c.TrustedPeerFile)
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
