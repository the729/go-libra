package client

import (
	"context"
	"errors"
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

func (c *Client) CmdQueryTransactionByAccountSeq(ctx *cli.Context) error {
	c.Connect()
	defer c.Disconnect()
	c.LoadTrustedPeers()
	c.LoadAccounts()

	account, err := c.GetAccount(ctx.Args().Get(0))
	if err != nil {
		return err
	}

	sequence, err := strconv.Atoi(ctx.Args().Get(1))
	if err != nil {
		return err
	}

	resp, err := c.ac.UpdateToLatestLedger(context.Background(), &pbtypes.UpdateToLatestLedgerRequest{
		RequestedItems: []*pbtypes.RequestItem{
			{
				RequestedItems: &pbtypes.RequestItem_GetAccountTransactionBySequenceNumberRequest{
					GetAccountTransactionBySequenceNumberRequest: &pbtypes.GetAccountTransactionBySequenceNumberRequest{
						Account:        account.Address,
						SequenceNumber: uint64(sequence),
						FetchEvents:    true,
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

	resp1 := resp.ResponseItems[0].GetGetAccountTransactionBySequenceNumberResponse()
	if resp1 == nil {
		return errors.New("nil response")
	}

	if resp1.SignedTransactionWithProof != nil {
		txn := &types.SignedTransactionWithProof{}
		if err = txn.FromProto(resp1.SignedTransactionWithProof); err != nil {
			return err
		}

		if err = txn.Verify(li.LedgerInfo); err != nil {
			return fmt.Errorf("transaction verify failed: %v", err)
		}

		log.Printf("Version: %d", txn.Version)
		log.Printf("Transaction detail:")
		spew.Dump(txn.SignedTransaction)
		log.Printf("Events:")
		spew.Dump(txn.Events)
	} else {
		state := &types.AccountStateWithProof{}
		err = state.FromProto(resp1.ProofOfCurrentSequenceNumber)
		if err != nil {
			return fmt.Errorf("account state with proof from proto failed: %v", err)
		}

		err = state.Verify(account.Address, li.LedgerInfo)
		if err != nil {
			return fmt.Errorf("account state with proof verify failed: %v", err)
		}

		if state.Blob != nil {
			err = state.Blob.ParseToMap()
			if err != nil {
				return fmt.Errorf("account blob cannot parse to map: %v", err)
			}

			resource, err := c.GetLibraCoinResourceFromAccountBlob(state.Blob)
			if err != nil {
				return err
			}
			log.Printf("Latest sequence number of the queried account is: %d", resource.SequenceNumber)
		} else {
			log.Printf("Queried account is not present")
		}
	}

	return nil
}
