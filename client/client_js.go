// +build js

package client

import (
	"github.com/the729/go-libra/generated/pbac"
)

func (c *Client) connect(server string) error {
	c.ac = pbac.NewAdmissionControlClient(server)
	return nil
}
