package types

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types/proof"
)

type SignedTransactionWithProof struct {
	*SignedTransaction
	Version uint64
	Events  EventList
	Proof   *SignedTransactionProof
}

type SignedTransactionProof struct {
	ledgerInfoToTransactionInfoProof *proof.Accumulator
	transactionInfo                  *TransactionInfo
}

func (t *SignedTransactionWithProof) FromProto(pb *pbtypes.SignedTransactionWithProof) error {
	var err error
	if pb == nil {
		return ErrNilInput
	}
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

	t.Proof = &SignedTransactionProof{}
	err = t.Proof.FromProto(pb.Proof)
	if err != nil {
		return err
	}
	return nil
}

func (tp *SignedTransactionProof) FromProto(pb *pbtypes.SignedTransactionProof) error {
	var err error
	if pb == nil {
		return ErrNilInput
	}

	tp.ledgerInfoToTransactionInfoProof = &proof.Accumulator{}
	err = tp.ledgerInfoToTransactionInfoProof.FromProto(pb.LedgerInfoToTransactionInfoProof)
	if err != nil {
		return err
	}
	tp.transactionInfo = &TransactionInfo{}
	err = tp.transactionInfo.FromProto(pb.TransactionInfo)
	if err != nil {
		return err
	}
	return nil
}

func (t *SignedTransactionWithProof) Verify(ledgerInfo *LedgerInfo) error {
	// according to https://community.libra.org/t/how-to-verify-a-signedtransaction-thoroughly/1214/3,
	// it is unnecessary to verify SignedTransaction itself

	txnHash := t.SignedTransaction.Hash()
	if !sha3libra.Equal(txnHash, t.Proof.transactionInfo.SignedTransactionHash) {
		return errors.New("signed txn hash mismatch")
	}
	eventHash := t.Events.Hash()
	if !sha3libra.Equal(eventHash, t.Proof.transactionInfo.EventRootHash) {
		return errors.New("event hash mismatch")
	}

	err := t.Proof.ledgerInfoToTransactionInfoProof.Verify(
		t.Version, t.Proof.transactionInfo.Hash(),
		ledgerInfo.TransactionAccumulatorHash,
	)
	if err != nil {
		return fmt.Errorf("cannot verify transaction info from ledger info: %v", err)
	}

	return nil
}
