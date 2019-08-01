package client

import (
	"context"
	"fmt"
	"time"

	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types"
)

func (c *Client) QueryTransactionRange(start, limit uint64, withEvents bool) (*types.ProvenTransactionList, error) {
	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.ac.UpdateToLatestLedger(ctx1, &pbtypes.UpdateToLatestLedgerRequest{
		RequestedItems: []*pbtypes.RequestItem{
			{
				RequestedItems: &pbtypes.RequestItem_GetTransactionsRequest{
					GetTransactionsRequest: &pbtypes.GetTransactionsRequest{
						StartVersion: start,
						Limit:        limit,
						FetchEvents:  withEvents,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// b, err := json.MarshalIndent(resp, "", "    ")
	// log.Printf("resp: %s", string(b))

	li := &types.LedgerInfoWithSignatures{}
	li.FromProto(resp.LedgerInfoWithSigs)
	pli, err := li.Verify(c.verifier)
	if err != nil {
		return nil, fmt.Errorf("ledger info verification failed: %v", err)
	}
	// log.Printf("Ledger info: version %d, time %d", li.LedgerInfo.Version, li.LedgerInfo.TimestampUsec)

	txnList := &types.TransactionListWithProof{}
	err = txnList.FromProtoResponse(resp.ResponseItems[0].GetGetTransactionsResponse())
	if err != nil {
		return nil, err
	}

	ptl, err := txnList.Verify(pli)
	if err != nil {
		return nil, fmt.Errorf("transaction list verification failed: %v", err)
	}

	// spew.Dump(ptl)

	return ptl, nil
}
