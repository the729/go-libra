package client

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli"

	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types"
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
	if limit > 100 {
		log.Printf("Limit>100, set to 100.")
		limit = 100
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

	// b, err := json.MarshalIndent(resp, "", "    ")
	// log.Printf("resp: %s", string(b))

	li := &types.LedgerInfoWithSignatures{}
	li.FromProto(resp.LedgerInfoWithSigs)
	err = li.Verify(c.verifier)
	if err != nil {
		return fmt.Errorf("verify failed: %v", err)
	}
	log.Printf("Ledger info: version %d, time %d", li.LedgerInfo.Version, li.LedgerInfo.TimestampUsec)

	// 1. 验证txn数量，恰好是请求的数量或不超过ledger version的数量
	txnList := &types.TransactionListWithProof{}
	err = txnList.FromProtoResponse(resp.ResponseItems[0].GetGetTransactionsResponse())
	if err != nil {
		return err
	}

	err = txnList.Verify(li.LedgerInfo)
	if err != nil {
		return fmt.Errorf("transaction list verify failed: %v", err)
	}

	spew.Dump(txnList.Transactions)

	return nil
}
