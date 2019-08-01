package types

import (
	"io"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
)

type TransactionInfo struct {
	SignedTransactionHash []byte
	StateRootHash         []byte
	EventRootHash         []byte
	GasUsed               uint64
}

func (t *TransactionInfo) FromProto(pb *pbtypes.TransactionInfo) error {
	t.SignedTransactionHash = pb.SignedTransactionHash
	t.StateRootHash = pb.StateRootHash
	t.EventRootHash = pb.EventRootHash
	t.GasUsed = pb.GasUsed
	return nil
}

func (t *TransactionInfo) SerializeTo(w io.Writer) error {
	w.Write(t.SignedTransactionHash)
	w.Write(t.StateRootHash)
	w.Write(t.EventRootHash)
	if err := serialization.SimpleSerializer.Write(w, t.GasUsed); err != nil {
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
