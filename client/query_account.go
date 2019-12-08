package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types"
)

// QueryAccountState queries account state from RPC server by account address, and does necessary
// crypto verifications.
func (c *Client) QueryAccountState(ctx context.Context, addr types.AccountAddress) (*types.ProvenAccountState, error) {
	c.accMu.RLock()
	frozenSubtreeRoots := cloneSubtrees(c.acc.FrozenSubtreeRoots)
	numLeaves := c.acc.NumLeaves
	c.accMu.RUnlock()

	resp, err := c.ac.UpdateToLatestLedger(ctx, &pbtypes.UpdateToLatestLedgerRequest{
		ClientKnownVersion: numLeaves - 1,
		RequestedItems: []*pbtypes.RequestItem{
			&pbtypes.RequestItem{
				RequestedItems: &pbtypes.RequestItem_GetAccountStateRequest{
					GetAccountStateRequest: &pbtypes.GetAccountStateRequest{
						Address: addr,
					},
				},
			},
		},
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

	account := &types.AccountStateWithProof{}
	err = account.FromProtoResponse(resp.ResponseItems[0].GetGetAccountStateResponse())
	if err != nil {
		return nil, fmt.Errorf("account state with proof from proto failed: %v", err)
	}

	paccount, err := account.Verify(addr, pli)
	if err != nil {
		return nil, fmt.Errorf("account state with proof verification failed: %v", err)
	}

	return paccount, nil
}

// GetLibraCoinResourceFromAccountBlob decodes the resource of Libra coin from a proven account blob.
func (c *Client) GetLibraCoinResourceFromAccountBlob(blob *types.ProvenAccountBlob) (*types.ProvenAccountResource, error) {
	res, err := blob.GetResource(types.AccountResourcePath())
	if err != nil {
		return nil, fmt.Errorf("get resource failed: %v", err)
	}
	return res, nil
}

// QueryAccountSequenceNumber queries sequence number of an account from RPC server, and does necessary
// crypto verifications.
func (c *Client) QueryAccountSequenceNumber(ctx context.Context, addr types.AccountAddress) (uint64, error) {
	paccount, err := c.QueryAccountState(ctx, addr)
	if err != nil {
		return 0, err
	}
	if paccount.IsNil() {
		return 0, errors.New("sender account not present in ledger")
	}
	resource, err := c.GetLibraCoinResourceFromAccountBlob(paccount.GetAccountBlob())
	if err != nil {
		return 0, err
	}
	return resource.GetSequenceNumber(), nil
}
