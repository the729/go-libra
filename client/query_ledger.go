package client

import (
	"context"
	"fmt"

	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types"
)

// QueryLedgerInfo queries ledger info from RPC server, and does necessary crypto verifications.
func (c *Client) QueryLedgerInfo(ctx context.Context) (*types.ProvenLedgerInfo, error) {
	c.accMu.RLock()
	frozenSubtreeRoots := cloneSubtrees(c.acc.FrozenSubtreeRoots)
	numLeaves := c.acc.NumLeaves
	c.accMu.RUnlock()

	resp, err := c.ac.UpdateToLatestLedger(ctx, &pbtypes.UpdateToLatestLedgerRequest{
		ClientKnownVersion: numLeaves - 1,
	})
	if err != nil {
		return nil, err
	}

	// respj, _ := json.MarshalIndent(resp, "", "    ")
	// log.Println(string(respj))

	pli, err := c.verifyLedgerInfoAndConsistency(resp, numLeaves, frozenSubtreeRoots)
	if err != nil {
		return nil, err
	}

	return pli, nil
}

func (c *Client) verifyLedgerInfoAndConsistency(
	resp *pbtypes.UpdateToLatestLedgerResponse,
	numLeaves uint64, frozenSubtreeRoots [][]byte,
) (*types.ProvenLedgerInfo, error) {

	li := &types.LedgerInfoWithSignatures{}
	li.FromProto(resp.LedgerInfoWithSigs)
	pli, err := li.Verify(c.verifier)
	if err != nil {
		return nil, fmt.Errorf("ledger info verification failed: %v", err)
	}
	numLeaves, frozenSubtreeRoots, err = pli.VerifyConsistency(
		numLeaves,
		frozenSubtreeRoots,
		resp.GetLedgerConsistencyProof().GetSubtrees(),
	)
	if err != nil {
		return nil, fmt.Errorf("ledger not consistent with known version: %v", err)
	}
	c.accMu.Lock()
	if numLeaves > c.acc.NumLeaves {
		c.acc.FrozenSubtreeRoots, c.acc.NumLeaves = frozenSubtreeRoots, numLeaves
	}
	c.accMu.Unlock()

	return pli, nil
}
