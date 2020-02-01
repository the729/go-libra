// +build !js

package client

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/the729/go-libra/generated/pbac"
)

func (c *Client) connect(server string) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("grpc dial error: %v", err)
	}
	c.closeFunc = func() { conn.Close() }
	c.ac = pbac.NewAdmissionControlClient(conn)
	return nil
}
