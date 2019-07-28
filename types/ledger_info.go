package types

import (
	"encoding/hex"
	"io"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types/validator"
)

type LedgerInfo struct {
	Version                    uint64
	TransactionAccumulatorHash []byte
	ConsensusDataHash          []byte
	ConsensusBlockID           []byte
	EpochNum                   uint64
	TimestampUsec              uint64
}

type LedgerInfoWithSignatures struct {
	*LedgerInfo
	Sigs map[string]sha3libra.HashValue
}

type ProvenLedgerInfo struct {
	proven     bool
	ledgerInfo LedgerInfo
}

func (l *LedgerInfo) FromProto(pb *pbtypes.LedgerInfo) error {
	l.Version = pb.Version
	l.TransactionAccumulatorHash = pb.TransactionAccumulatorHash
	l.ConsensusDataHash = pb.ConsensusDataHash
	l.ConsensusBlockID = pb.ConsensusBlockId
	l.EpochNum = pb.EpochNum
	l.TimestampUsec = pb.TimestampUsecs
	return nil
}

func (l *LedgerInfo) SerializeTo(w io.Writer) error {
	serialization.SimpleSerializer.Write(w, l.Version)
	w.Write(l.TransactionAccumulatorHash)
	w.Write(l.ConsensusDataHash)
	w.Write(l.ConsensusBlockID)
	serialization.SimpleSerializer.Write(w, l.EpochNum)
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
	l.Sigs = sigs
	return nil
}

func (l *LedgerInfoWithSignatures) Verify(v validator.Verifier) (*ProvenLedgerInfo, error) {
	if err := v.Verify(l.LedgerInfo.Hash(), l.Sigs); err != nil {
		return nil, err
	}
	return &ProvenLedgerInfo{
		proven: true,
		ledgerInfo: LedgerInfo{
			Version:                    l.LedgerInfo.Version,
			TransactionAccumulatorHash: cloneBytes(l.LedgerInfo.TransactionAccumulatorHash),
			ConsensusDataHash:          cloneBytes(l.LedgerInfo.ConsensusDataHash),
			ConsensusBlockID:           cloneBytes(l.LedgerInfo.ConsensusBlockID),
			EpochNum:                   l.LedgerInfo.EpochNum,
			TimestampUsec:              l.LedgerInfo.TimestampUsec,
		},
	}, nil
}

func (pl *ProvenLedgerInfo) GetVersion() uint64 {
	if !pl.proven {
		panic("not valid proven ledger info")
	}
	return pl.ledgerInfo.Version
}

func (pl *ProvenLedgerInfo) GetTransactionAccumulatorHash() []byte {
	if !pl.proven {
		panic("not valid proven ledger info")
	}
	return cloneBytes(pl.ledgerInfo.TransactionAccumulatorHash)
}

func (pl *ProvenLedgerInfo) GetEpochNum() uint64 {
	if !pl.proven {
		panic("not valid proven ledger info")
	}
	return pl.ledgerInfo.EpochNum
}

func (pl *ProvenLedgerInfo) GetTimestampUsec() uint64 {
	if !pl.proven {
		panic("not valid proven ledger info")
	}
	return pl.ledgerInfo.TimestampUsec
}
