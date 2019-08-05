package client

import (
	"context"
	"fmt"
	"time"

	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types"
)

// QueryLedgerInfo queries ledger info from RPC server, and does necessary crypto verifications.
func (c *Client) QueryLedgerInfo(knownVersion uint64) (*types.ProvenLedgerInfo, error) {
	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.ac.UpdateToLatestLedger(ctx1, &pbtypes.UpdateToLatestLedgerRequest{
		ClientKnownVersion: knownVersion,
	})
	if err != nil {
		return nil, err
	}

	// respj, _ := json.MarshalIndent(resp, "", "    ")
	// log.Println(string(respj))

	li := &types.LedgerInfoWithSignatures{}
	li.FromProto(resp.LedgerInfoWithSigs)
	pli, err := li.Verify(c.verifier)
	if err != nil {
		return nil, fmt.Errorf("ledger info verification failed: %v", err)
	}

	return pli, nil
}
