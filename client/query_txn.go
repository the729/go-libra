package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types"
)

// QueryTransactionRange queries a list of transactions from RPC server, and does necessary
// crypto verifications.
func (c *Client) QueryTransactionRange(ctx context.Context, start, limit uint64, withEvents bool) (*types.ProvenTransactionList, error) {
	resp, err := c.ac.UpdateToLatestLedger(ctx, &pbtypes.UpdateToLatestLedgerRequest{
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

// QueryTransactionByAccountSeq queries the transaction that is sent from a specific account at a specific sequence number,
// and does necessary crypto verifications.
func (c *Client) QueryTransactionByAccountSeq(ctx context.Context, addr types.AccountAddress, sequence uint64, withEvents bool) (*types.ProvenTransaction, error) {
	resp, err := c.ac.UpdateToLatestLedger(ctx, &pbtypes.UpdateToLatestLedgerRequest{
		RequestedItems: []*pbtypes.RequestItem{
			{
				RequestedItems: &pbtypes.RequestItem_GetAccountTransactionBySequenceNumberRequest{
					GetAccountTransactionBySequenceNumberRequest: &pbtypes.GetAccountTransactionBySequenceNumberRequest{
						Account:        addr,
						SequenceNumber: sequence,
						FetchEvents:    withEvents,
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

	resp1 := resp.ResponseItems[0].GetGetAccountTransactionBySequenceNumberResponse()
	if resp1 == nil {
		return nil, errors.New("nil response")
	}

	if resp1.TransactionWithProof == nil {
		state := &types.AccountStateWithProof{}
		err = state.FromProto(resp1.ProofOfCurrentSequenceNumber)
		if err != nil {
			return nil, fmt.Errorf("account state with proof from proto failed: %v", err)
		}

		pstate, err := state.Verify(addr, pli)
		if err != nil {
			return nil, fmt.Errorf("account state with proof verify failed: %v", err)
		}

		if pstate.IsNil() {
			return nil, errors.New("account not exist")
		}

		pres, err := c.GetLibraCoinResourceFromAccountBlob(pstate.GetAccountBlob())
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("sequence too large, should < %v", pres.GetSequenceNumber())
	}

	txn := &types.TransactionWithProof{}
	if err = txn.FromProto(resp1.TransactionWithProof); err != nil {
		return nil, err
	}

	ptxn, err := txn.Verify(pli)
	if err != nil {
		return nil, fmt.Errorf("transaction verify failed: %v", err)
	}
	return ptxn, nil
}
