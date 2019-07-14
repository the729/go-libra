package types

import (
	"io"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
)

type TransactionInfo struct {
	signedTransactionHash []byte
	stateRootHash         []byte
	eventRootHash         []byte
	gasUsed               uint64
}

func (t *TransactionInfo) FromProto(pb *pbtypes.TransactionInfo) error {
	t.signedTransactionHash = pb.SignedTransactionHash
	t.stateRootHash = pb.StateRootHash
	t.eventRootHash = pb.EventRootHash
	t.gasUsed = pb.GasUsed
	return nil
}

func (t *TransactionInfo) SerializeTo(w io.Writer) error {
	w.Write(t.signedTransactionHash)
	w.Write(t.stateRootHash)
	w.Write(t.eventRootHash)
	if err := serialization.SimpleSerializer.Write(w, t.gasUsed); err != nil {
		return err
	}
	return nil
}

func (t *TransactionInfo) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewTransactionInfo()
	if err := t.SerializeTo(hasher); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
}
