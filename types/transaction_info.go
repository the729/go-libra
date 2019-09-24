package types

import (
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/lcs"
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

	MajorStatus VMStatusCode
}

// FromProto parses a protobuf struct into this struct.
func (t *TransactionInfo) FromProto(pb *pbtypes.TransactionInfo) error {
	t.SignedTransactionHash = pb.SignedTransactionHash
	t.StateRootHash = pb.StateRootHash
	t.EventRootHash = pb.EventRootHash
	t.GasUsed = pb.GasUsed
	t.MajorStatus = VMStatusCode(pb.MajorStatus)
	return nil
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (t *TransactionInfo) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewTransactionInfo()
	if err := lcs.NewEncoder(hasher).Encode(t); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
}
