package types

import (
	"fmt"

	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types/proof"
)

// SignedTransactionWithProof is a submitted transaction with a Merkle Tree accumulator proof
// to prove its inclusion in a version of the ledger.
type SignedTransactionWithProof struct {
	*SubmittedTransaction
	LedgerInfoToTransactionInfoProof *proof.Accumulator
}

// FromProto parses a protobuf struct into this struct.
func (t *SignedTransactionWithProof) FromProto(pb *pbtypes.SignedTransactionWithProof) error {
	var err error
	if pb == nil {
		return ErrNilInput
	}
	t.SubmittedTransaction = &SubmittedTransaction{}
	t.Version = pb.Version

	t.SignedTransaction = &SignedTransaction{}
	err = t.SignedTransaction.FromProto(pb.SignedTransaction)
	if err != nil {
		return err
	}

	t.Events = make([]*ContractEvent, 0, len(pb.Events.Events))
	for _, ev := range pb.Events.Events {
		ev1 := &ContractEvent{}
		if err := ev1.FromProto(ev); err != nil {
			return err
		}
		t.Events = append(t.Events, ev1)
	}

	t.Info = &TransactionInfo{}
	err = t.Info.FromProto(pb.Proof.TransactionInfo)
	if err != nil {
		return err
	}

	t.LedgerInfoToTransactionInfoProof = &proof.Accumulator{}
	err = t.LedgerInfoToTransactionInfoProof.FromProto(pb.Proof.LedgerInfoToTransactionInfoProof)
	if err != nil {
		return err
	}
	return nil
}

// Verify the proof of the transaction, and output a ProvenTransaction if successful.
func (t *SignedTransactionWithProof) Verify(ledgerInfo *ProvenLedgerInfo) (*ProvenTransaction, error) {
	pTxn, err := t.SubmittedTransaction.Verify()
	if err != nil {
		return nil, err
	}

	err = t.LedgerInfoToTransactionInfoProof.Verify(
		t.Version, t.Info.Hash(),
		ledgerInfo.GetTransactionAccumulatorHash(),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot verify transaction info from ledger info: %v", err)
	}

	pTxn.proven = true
	return pTxn, nil
}
