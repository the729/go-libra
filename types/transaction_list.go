package types

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types/proof"
)

type TransactionListWithProof struct {
	Transactions []*SubmittedTransaction
	Proof        *proof.AccumulatorRange
}

type ProvenTransactionList struct {
	proven       bool
	transactions []*ProvenTransaction
	ledgerInfo   *ProvenLedgerInfo
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

	if pb.FirstTransactionVersion == nil && len(pb.Transactions) > 0 {
		return errors.New("missing first txn version")
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
		item := &SubmittedTransaction{
			SignedTransaction: txn,
			Info:              info,
			Version:           pb.FirstTransactionVersion.Value + uint64(idx),
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

func (ptl *ProvenTransactionList) GetTransactions() []*ProvenTransaction {
	if !ptl.proven {
		panic("not valid proven transaction list")
	}
	out := make([]*ProvenTransaction, len(ptl.transactions))
	copy(out, ptl.transactions)
	return out
}

func (ptl *ProvenTransactionList) GetLedgerInfo() *ProvenLedgerInfo {
	if !ptl.proven {
		panic("not valid proven transaction list")
	}
	return ptl.ledgerInfo
}
