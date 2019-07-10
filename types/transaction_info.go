package types

import (
	"encoding/binary"

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

func (t *TransactionInfo) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewTransactionInfo()
	hasher.Write(t.signedTransactionHash)
	hasher.Write(t.stateRootHash)
	hasher.Write(t.eventRootHash)
	binary.Write(hasher, binary.LittleEndian, t.gasUsed)
	return hasher.Sum([]byte{})
}
