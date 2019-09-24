package types

import (
	"encoding/hex"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types/validator"
	"github.com/the729/lcs"
)

// LedgerInfo is a information struct of a version (height) of the ledger.
type LedgerInfo struct {
	Version                    uint64
	TransactionAccumulatorHash []byte
	ConsensusDataHash          []byte
	ConsensusBlockID           []byte
	EpochNum                   uint64
	TimestampUsec              uint64
	NextValidatorSet           ValidatorSet `lcs:"optional"`
}

// LedgerInfoWithSignatures is a ledger info with signature from trusted
// validators.
type LedgerInfoWithSignatures struct {
	*LedgerInfo
	Sigs map[string]sha3libra.HashValue
}

// ProvenLedgerInfo is a ledger info proven to be history state of the ledger.
type ProvenLedgerInfo struct {
	proven     bool
	ledgerInfo LedgerInfo
}

// FromProto parses a protobuf struct into this struct.
func (l *LedgerInfo) FromProto(pb *pbtypes.LedgerInfo) error {
	l.Version = pb.Version
	l.TransactionAccumulatorHash = pb.TransactionAccumulatorHash
	l.ConsensusDataHash = pb.ConsensusDataHash
	l.ConsensusBlockID = pb.ConsensusBlockId
	l.EpochNum = pb.EpochNum
	l.TimestampUsec = pb.TimestampUsecs
	if pb.NextValidatorSet != nil {
		if err := l.NextValidatorSet.FromProto(pb.NextValidatorSet); err != nil {
			return err
		}
	}
	return nil
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (l *LedgerInfo) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewLedgerInfo()
	if err := lcs.NewEncoder(hasher).Encode(l); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
}

// FromProto parses a protobuf struct into this struct.
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

// Verify the ledger info with a consensus verifier and output a ProvenLedgerInfo.
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

// GetVersion returns the height of this ledger info.
func (pl *ProvenLedgerInfo) GetVersion() uint64 {
	if !pl.proven {
		panic("not valid proven ledger info")
	}
	return pl.ledgerInfo.Version
}

// GetTransactionAccumulatorHash returns the root hash of the transaction Merkle Tree accumulator.
func (pl *ProvenLedgerInfo) GetTransactionAccumulatorHash() []byte {
	if !pl.proven {
		panic("not valid proven ledger info")
	}
	return cloneBytes(pl.ledgerInfo.TransactionAccumulatorHash)
}

// GetEpochNum returns the epoch number.
func (pl *ProvenLedgerInfo) GetEpochNum() uint64 {
	if !pl.proven {
		panic("not valid proven ledger info")
	}
	return pl.ledgerInfo.EpochNum
}

// GetTimestampUsec returns the timestamp of this version, in microseconds.
func (pl *ProvenLedgerInfo) GetTimestampUsec() uint64 {
	if !pl.proven {
		panic("not valid proven ledger info")
	}
	return pl.ledgerInfo.TimestampUsec
}
