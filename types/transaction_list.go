package types

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types/proof"
)

// TransactionListWithProof is a consecutive list of transactions with a Merkle Tree accumulator
// range proof to prove inclusion of all transactions in the list.
type TransactionListWithProof struct {
	Transactions []*SubmittedTransaction
	Proof        *proof.AccumulatorRange
}

// ProvenTransactionList is a consecutive list of transactions which has been proven to be included
// in the ledger.
type ProvenTransactionList struct {
	proven       bool
	transactions []*ProvenTransaction
	ledgerInfo   *ProvenLedgerInfo
}

// FromProtoResponse parses a protobuf struct into this struct.
func (tl *TransactionListWithProof) FromProtoResponse(pb *pbtypes.GetTransactionsResponse) error {
	if pb == nil {
		return ErrNilInput
	}
	return tl.FromProto(pb.TxnListWithProof)
}

// FromProto parses a protobuf struct into this struct.
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

	if pb.FirstTransactionVersion == nil && len(pb.Transactions) > 0 {
		return errors.New("missing first txn version")
	}

	tl.Transactions = nil
	for idx := range pb.Transactions {
		info := &TransactionInfo{}
		if err := info.FromProto(pb.Infos[idx]); err != nil {
			return err
		}
		item := &SubmittedTransaction{
			RawSignedTxn: pb.Transactions[idx].SignedTxn,
			Info:         info,
			Version:      pb.FirstTransactionVersion.Value + uint64(idx),
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

// Verify the proof of the transaction list, and output a ProvenTransactionList if successful.
func (tl *TransactionListWithProof) Verify(ledgerInfo *ProvenLedgerInfo) (*ProvenTransactionList, error) {
	var firstVersion uint64

	if len(tl.Transactions) > 0 {
		// verify that submitted txn list contains consecutive txns
		firstVersion = tl.Transactions[0].Version
		for idx, txn := range tl.Transactions {
			if txn.Version != firstVersion+uint64(idx) {
				return nil, errors.New("transaction version not consective")
			}
		}
		if firstVersion+uint64(len(tl.Transactions))-1 > ledgerInfo.GetVersion() {
			return nil, errors.New("last transaction version greater than ledger version")
		}
	}

	if tl.Proof == nil {
		return nil, errors.New("nil proof")
	}

	hashes := make([]sha3libra.HashValue, 0)
	provenTxns := make([]*ProvenTransaction, 0, len(tl.Transactions))
	// 1. verify signed transactions, and events
	for _, t := range tl.Transactions {
		provenTxn, err := t.Verify()
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, t.Info.Hash())
		provenTxn.proven = true
		provenTxns = append(provenTxns, provenTxn)
	}

	// 2. verify transaction accumulator
	err := tl.Proof.Verify(firstVersion, hashes, ledgerInfo.GetTransactionAccumulatorHash())
	if err != nil {
		return nil, fmt.Errorf("accumulator range proof failed: %v", err)
	}

	return &ProvenTransactionList{
		proven:       true,
		transactions: provenTxns,
		ledgerInfo:   ledgerInfo,
	}, nil
}

// GetTransactions returns a copy of the underlying proven transaction list.
func (ptl *ProvenTransactionList) GetTransactions() []*ProvenTransaction {
	if !ptl.proven {
		panic("not valid proven transaction list")
	}
	out := make([]*ProvenTransaction, len(ptl.transactions))
	copy(out, ptl.transactions)
	return out
}

// GetLedgerInfo returns the ledger info which proofs this transaction list.
func (ptl *ProvenTransactionList) GetLedgerInfo() *ProvenLedgerInfo {
	if !ptl.proven {
		panic("not valid proven transaction list")
	}
	return ptl.ledgerInfo
}
