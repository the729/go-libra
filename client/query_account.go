package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types"
)

func (c *Client) QueryAccountState(addr types.AccountAddress) (*types.ProvenAccountState, error) {
	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.ac.UpdateToLatestLedger(ctx1, &pbtypes.UpdateToLatestLedgerRequest{
		ClientKnownVersion: 0,
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
		return nil, fmt.Errorf("rpc failed: %v", err)
	}
	// respj, _ := json.MarshalIndent(resp, "", "    ")
	// log.Println(string(respj))

	li := &types.LedgerInfoWithSignatures{}
	li.FromProto(resp.LedgerInfoWithSigs)
	pli, err := li.Verify(c.verifier)
	if err != nil {
		return nil, fmt.Errorf("verify failed: %v", err)
	}

	account := &types.AccountStateWithProof{}
	err = account.FromProtoResponse(resp.ResponseItems[0].GetGetAccountStateResponse())
	if err != nil {
		return nil, fmt.Errorf("account state with proof from proto failed: %v", err)
	}

	paccount, err := account.Verify(addr, pli)
	if err != nil {
		return nil, fmt.Errorf("account state with proof verify failed: %v", err)
	}

	return paccount, nil
}

func (c *Client) GetLibraCoinResourceFromAccountBlob(blob *types.ProvenAccountBlob) (*types.ProvenAccountResource, error) {
	res, err := blob.GetResource(&types.StructTag{
		Address: make([]byte, 32),
		Module:  "LibraAccount",
		Name:    "T",
	})
	if err != nil {
		return nil, fmt.Errorf("get resource failed: %v", err)
	}
	return res, nil
}

func (c *Client) GetAccountSequenceNumber(addr types.AccountAddress) (uint64, error) {
	paccount, err := c.QueryAccountState(addr)
	if err != nil {
		return 0, err
	}
	if paccount.IsNil() {
		return 0, errors.New("sender account not present in ledger.")
	}
	resource, err := c.GetLibraCoinResourceFromAccountBlob(paccount.GetAccountBlob())
	if err != nil {
		return 0, err
	}
	return resource.GetSequenceNumber(), nil
}
