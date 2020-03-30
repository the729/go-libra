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

	verifier := c.verifier
	lastWaypoint := ""
	if verifier.EpochChangeVerificationRequired(li.Epoch) {
		vcp := &types.ValidatorChangeProof{}
		if err := vcp.FromProto(resp.ValidatorChangeProof); err != nil {
			return nil, fmt.Errorf("validator change proof invalid: %v", err)
		}
		epochChange, err := vcp.Verify(verifier)
		if err != nil {
			return nil, fmt.Errorf("validator change proof verification error: %v", err)
		}
		if genesisHash := epochChange.GetGenesisHash(); genesisHash != nil {
			// this is the genesis block, update accumulator
			numLeaves = 1
			frozenSubtreeRoots = [][]byte{genesisHash}
		}
		pli := epochChange.GetLastLedgerInfo()
		v, err := pli.ToVerifier()
		if err != nil {
			return nil, err
		}

		// log.Printf("Ledger info verifier updated, epoch = %d, version = %d", pli.GetEpochNum()+1, pli.GetVersion())
		// spew.Dump(v)

		verifier = v
		lastWaypointB, _ := (&types.Waypoint{}).FromProvenLedgerInfo(pli).MarshalText()
		lastWaypoint = string(lastWaypointB)
	}
	pli, err := li.Verify(verifier)
	if err != nil {
		return nil, fmt.Errorf("ledger info verification failed: %v", err)
	}
	if frozenSubtreeRoots != nil {
		numLeaves, frozenSubtreeRoots, err = pli.VerifyConsistency(
			numLeaves,
			frozenSubtreeRoots,
			resp.GetLedgerConsistencyProof().GetSubtrees(),
		)
		if err != nil {
			return nil, fmt.Errorf("ledger not consistent with known version: %v", err)
		}
	} else {
		numLeaves = pli.GetVersion() + 1
	}

	c.accMu.Lock()
	if numLeaves > c.acc.NumLeaves {
		c.acc.FrozenSubtreeRoots, c.acc.NumLeaves = frozenSubtreeRoots, numLeaves
	}
	if lastWaypoint != "" {
		c.verifier = verifier
		c.lastWaypoint = lastWaypoint
	}
	c.accMu.Unlock()

	return pli, nil
}
