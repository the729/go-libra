package client

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/urfave/cli"
)

func (c *Client) CmdQueryTransactionRange(ctx *cli.Context) error {
	c.Connect()
	defer c.Disconnect()
	c.LoadTrustedPeers()

	start, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return err
	}
	limit, err := strconv.Atoi(ctx.Args().Get(1))
	if err != nil {
		return err
	}

	resp, err := c.ac.UpdateToLatestLedger(context.Background(), &pbtypes.UpdateToLatestLedgerRequest{
		RequestedItems: []*pbtypes.RequestItem{
			{
				RequestedItems: &pbtypes.RequestItem_GetTransactionsRequest{
					GetTransactionsRequest: &pbtypes.GetTransactionsRequest{
						StartVersion: uint64(start),
						Limit:        uint64(limit),
						FetchEvents:  true,
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(resp, "", "    ")
	log.Printf("resp: %s", string(b))
	// log.Printf("Ledger info: version %d, time %d", ledgerInfo.Version, ledgerInfo.TimestampUsec)
	return nil
}
