package types

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/crypto/sha3libra"

	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types/proof"
)

type TransactionListItem struct {
	*SignedTransaction
	Info   *TransactionInfo
	Events EventList
}

type TransactionListWithProof struct {
	Transactions    []*TransactionListItem
	FirstTxnVersion uint64
	Proof           *proof.AccumulatorRange
}

func (tl *TransactionListWithProof) FromProtoResponse(pb *pbtypes.GetTransactionsResponse) error {
	if pb == nil {
		return ErrNilInput
	}
	return tl.FromProto(pb.TxnListWithProof)
}

func (tl *TransactionListWithProof) FromProto(pb *pbtypes.TransactionListWithProof) error {
	if pb == nil {
		return ErrNilInput
	}

	if len(pb.Transactions) != len(pb.Infos) {
		return errors.New("mismatch length: txns and infos")
	}

	var eventsList []*pbtypes.EventsList
	if pb.EventsForVersions != nil {
		if len(pb.EventsForVersions.EventsForVersion) != len(pb.Transactions) {
			return errors.New("mismatch length: txns and events")
		}
		eventsList = pb.EventsForVersions.EventsForVersion
	}

	tl.Transactions = nil
	for idx := range pb.Transactions {
		txn := &SignedTransaction{}
		if err := txn.FromProto(pb.Transactions[idx]); err != nil {
			return err
		}
		info := &TransactionInfo{}
		if err := info.FromProto(pb.Infos[idx]); err != nil {
			return err
		}
		item := &TransactionListItem{
			SignedTransaction: txn,
			Info:              info,
		}

		if eventsList != nil {
			for _, ev := range eventsList[idx].Events {
				ev1 := &ContractEvent{}
				if err := ev1.FromProto(ev); err != nil {
					return err
				}
				item.Events = append(item.Events, ev1)
			}
		}

		tl.Transactions = append(tl.Transactions, item)
	}

	if pb.FirstTransactionVersion != nil {
		tl.FirstTxnVersion = pb.FirstTransactionVersion.Value
	}

	tl.Proof = &proof.AccumulatorRange{}
	if pb.ProofOfFirstTransaction != nil {
		tl.Proof.First = &proof.Accumulator{}
		if err := tl.Proof.First.FromProto(pb.ProofOfFirstTransaction); err != nil {
			return err
		}
	}
	if pb.ProofOfLastTransaction != nil {
		tl.Proof.Last = &proof.Accumulator{}
		if err := tl.Proof.Last.FromProto(pb.ProofOfLastTransaction); err != nil {
			return err
		}
	}
	return nil
}

func (tl *TransactionListWithProof) Verify(ledgerInfo *LedgerInfo) error {
	if len(tl.Transactions) > 0 && tl.FirstTxnVersion+uint64(len(tl.Transactions))-1 > ledgerInfo.Version {
		return errors.New("last transaction version greater than ledger version")
	}

	if tl.Proof == nil {
		return errors.New("nil proof")
	}

	hashes := make([]sha3libra.HashValue, 0)
	// 1. verify signed transactions, and events
	for i, t := range tl.Transactions {
		// according to https://community.libra.org/t/how-to-verify-a-signedtransaction-thoroughly/1214/3,
		// it is unnecessary to verify SignedTransaction itself
		if err := t.SignedTransaction.Verify(); err != nil {
			return fmt.Errorf("txn(%d) signature verification fail: %v", tl.FirstTxnVersion+uint64(i), err)
		}

		// verify SignedTransaction and Events hash from transaction info
		txnHash := t.SignedTransaction.Hash()
		if !sha3libra.Equal(txnHash, t.Info.signedTransactionHash) {
			return fmt.Errorf("signed txn hash mismatch in txn(%d)", tl.FirstTxnVersion+uint64(i))
		}
		eventHash := t.Events.Hash()
		if !sha3libra.Equal(eventHash, t.Info.eventRootHash) {
			return fmt.Errorf("event root hash mismatch in txn(%d)", tl.FirstTxnVersion+uint64(i))
		}

		hashes = append(hashes, t.Info.Hash())
	}

	// 2. verify transaction accumulator
	err := tl.Proof.Verify(tl.FirstTxnVersion, hashes, ledgerInfo.transactionAccumulatorHash)
	if err != nil {
		return fmt.Errorf("accumulator range proof failed: %v", err)
	}
	return nil
}
