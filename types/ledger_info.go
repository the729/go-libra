package types

import (
	"encoding/hex"
	"io"

	"github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types/validator"
)

type LedgerInfo struct {
	Version                    uint64
	transactionAccumulatorHash []byte
	consensusDataHash          []byte
	consensusBlockID           []byte
	epochNum                   uint64
	TimestampUsec              uint64
}

type LedgerInfoWithSignatures struct {
	LedgerInfo *LedgerInfo
	sigs       map[string]sha3libra.HashValue
}

func (l *LedgerInfo) FromProto(pb *pbtypes.LedgerInfo) error {
	l.Version = pb.Version
	l.transactionAccumulatorHash = pb.TransactionAccumulatorHash
	l.consensusDataHash = pb.ConsensusDataHash
	l.consensusBlockID = pb.ConsensusBlockId
	l.epochNum = pb.EpochNum
	l.TimestampUsec = pb.TimestampUsecs
	return nil
}

func (l *LedgerInfo) SerializeTo(w io.Writer) error {
	serialization.SimpleSerializer.Write(w, l.Version)
	w.Write(l.transactionAccumulatorHash)
	w.Write(l.consensusDataHash)
	w.Write(l.consensusBlockID)
	serialization.SimpleSerializer.Write(w, l.epochNum)
	serialization.SimpleSerializer.Write(w, l.TimestampUsec)
	return nil
}

func (l *LedgerInfo) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewLedgerInfo()
	l.SerializeTo(hasher)
	hash := hasher.Sum([]byte{})

	return hash
}

func (l *LedgerInfoWithSignatures) FromProto(pb *pbtypes.LedgerInfoWithSignatures) error {
	l.LedgerInfo = &LedgerInfo{}
	l.LedgerInfo.FromProto(pb.LedgerInfo)

	sigs := make(map[string]sha3libra.HashValue)
	for _, s := range pb.Signatures {
		sigs[hex.EncodeToString(s.ValidatorId)] = s.Signature
	}
	l.sigs = sigs
	return nil
}

func (l *LedgerInfoWithSignatures) Verify(v validator.Verifier) error {
	return v.Verify(l.LedgerInfo.Hash(), l.sigs)
}
