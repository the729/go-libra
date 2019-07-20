package client

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types"
	"github.com/urfave/cli"
)

func (c *Client) QueryLedgerInfo(knownVersion uint64) (*types.LedgerInfo, error) {
	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.ac.UpdateToLatestLedger(ctx1, &pbtypes.UpdateToLatestLedgerRequest{
		ClientKnownVersion: knownVersion,
	})
	if err != nil {
		return nil, fmt.Errorf("rpc failed: %v", err)
	}

	// respj, _ := json.MarshalIndent(resp, "", "    ")
	// log.Println(string(respj))

	li := &types.LedgerInfoWithSignatures{}
	li.FromProto(resp.LedgerInfoWithSigs)
	err = li.Verify(c.verifier)
	if err != nil {
		return nil, fmt.Errorf("verify failed: %v", err)
	}

	return li.LedgerInfo, nil
}

func (c *Client) CmdQueryLedgerInfo(ctx *cli.Context) error {
	c.Connect()
	defer c.Disconnect()
	c.LoadTrustedPeers()
	c.LoadAccounts()

	knownVersion, err := strconv.ParseUint(ctx.Args().Get(0), 10, 64)
	if err != nil {
		return err
	}
	ledgerInfo, err := c.QueryLedgerInfo(knownVersion)
	if err != nil {
		return err
	}

	log.Printf("Ledger info: version %d, time %d", ledgerInfo.Version, ledgerInfo.TimestampUsec)
	return nil
}
