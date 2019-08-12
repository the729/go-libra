package types

import (
	"io"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
)

// TransactionInfo is a information struct of a submitted transaction.
type TransactionInfo struct {
	// SignedTransactionHash is the hash of this transaction.
	SignedTransactionHash []byte

	// StateRootHash is the root hash of a Merkle Tree accumulator built form states of all
	// existing accounts in the ledger.
	// It represents the whole state of the ledger after execution of this transaction.
	StateRootHash []byte

	// EventRootHash is the root hash of a Merkle Tree accumulator built from all output events.
	EventRootHash []byte

	// GasUsed is the actual gas used to process this transaction, in microLibra.
	GasUsed uint64
}

// FromProto parses a protobuf struct into this struct.
func (t *TransactionInfo) FromProto(pb *pbtypes.TransactionInfo) error {
	t.SignedTransactionHash = pb.SignedTransactionHash
	t.StateRootHash = pb.StateRootHash
	t.EventRootHash = pb.EventRootHash
	t.GasUsed = pb.GasUsed
	return nil
}

// SerializeTo serializes this struct into a io.Writer.
func (t *TransactionInfo) SerializeTo(w io.Writer) error {
	w.Write(t.SignedTransactionHash)
	w.Write(t.StateRootHash)
	w.Write(t.EventRootHash)
	if err := serialization.SimpleSerializer.Write(w, t.GasUsed); err != nil {
		return err
	}
	return nil
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (t *TransactionInfo) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewTransactionInfo()
	if err := t.SerializeTo(hasher); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
}
