package types

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types/proof/accumulator"
	"github.com/the729/lcs"
)

// LedgerInfo is a information struct of a version (height) of the ledger.
type LedgerInfo struct {
	Epoch                      uint64
	Round                      uint64
	ConsensusBlockID           []byte
	TransactionAccumulatorHash []byte
	Version                    uint64
	TimestampUsec              uint64
	NextValidatorSet           *ValidatorSet `lcs:"optional"`
	ConsensusDataHash          []byte
}

type LedgerInfoWithSignatures struct {
	Value isLedgerInfoWithSignatures `lcs:"enum=isLedgerInfoWithSignatures"`
}

type isLedgerInfoWithSignatures interface {
	isLedgerInfoWithSignatures()
}

// LedgerInfoWithSignaturesV0 is a ledger info with signature from trusted
// validators.
type LedgerInfoWithSignaturesV0 struct {
	*LedgerInfo
	Sigs map[AccountAddress]HashValue
}

func (*LedgerInfoWithSignaturesV0) isLedgerInfoWithSignatures() {}

var ledgerInfoWSigsEnumDef = []lcs.EnumVariant{
	{
		Name:     "isLedgerInfoWithSignatures",
		Value:    0,
		Template: (*LedgerInfoWithSignaturesV0)(nil),
	},
}

// EnumTypes defines enum variants for lcs
func (*LedgerInfoWithSignatures) EnumTypes() []lcs.EnumVariant { return ledgerInfoWSigsEnumDef }

// ProvenLedgerInfo is a ledger info proven to be history state of the ledger.
type ProvenLedgerInfo struct {
	proven     bool
	ledgerInfo *LedgerInfo
}

// FromProto parses a protobuf struct into this struct.
func (l *LedgerInfo) FromProto(pb *pbtypes.LedgerInfo) error {
	l.Version = pb.Version
	l.TransactionAccumulatorHash = pb.TransactionAccumulatorHash
	l.ConsensusDataHash = pb.ConsensusDataHash
	l.ConsensusBlockID = pb.ConsensusBlockId
	l.Epoch = pb.Epoch
	l.Round = pb.Round
	l.TimestampUsec = pb.TimestampUsecs
	if pb.NextValidatorSet != nil {
		if err := l.NextValidatorSet.FromProto(pb.NextValidatorSet); err != nil {
			return err
		}
	} else {
		l.NextValidatorSet = nil
	}
	return nil
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (l *LedgerInfo) Hash() HashValue {
	hasher := sha3libra.NewLedgerInfo()
	if err := lcs.NewEncoder(hasher).Encode(l); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
}

// Clone deep clones this struct.
func (l *LedgerInfo) Clone() *LedgerInfo {
	out := &LedgerInfo{
		Epoch:                      l.Epoch,
		Round:                      l.Round,
		ConsensusBlockID:           cloneBytes(l.ConsensusBlockID),
		TransactionAccumulatorHash: cloneBytes(l.TransactionAccumulatorHash),
		Version:                    l.Version,
		TimestampUsec:              l.TimestampUsec,
		NextValidatorSet:           l.NextValidatorSet.Clone(),
		ConsensusDataHash:          cloneBytes(l.ConsensusDataHash),
	}
	return out
}

// FromProto parses a protobuf struct into this struct.
func (l *LedgerInfoWithSignatures) FromProto(pb *pbtypes.LedgerInfoWithSignatures) error {
	return lcs.Unmarshal(pb.Bytes, l)
}

// Verify the ledger info with a consensus verifier and output a ProvenLedgerInfo.
func (l *LedgerInfoWithSignaturesV0) Verify(v LedgerInfoVerifier) (*ProvenLedgerInfo, error) {
	if err := v.Verify(&LedgerInfoWithSignatures{l}); err != nil {
		return nil, err
	}
	return &ProvenLedgerInfo{
		proven:     true,
		ledgerInfo: l.LedgerInfo.Clone(),
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
	return pl.ledgerInfo.Epoch
}

// GetTimestampUsec returns the timestamp of this version, in microseconds.
func (pl *ProvenLedgerInfo) GetTimestampUsec() uint64 {
	if !pl.proven {
		panic("not valid proven ledger info")
	}
	return pl.ledgerInfo.TimestampUsec
}

// ToVerifier builds a ValidatorVerifier using the next validator set in this
// LedgerInfo. Only works when this LedgerInfo is at a boundary of epochs.
func (pl *ProvenLedgerInfo) ToVerifier() (LedgerInfoVerifier, error) {
	if !pl.proven {
		panic("not valid proven ledger info")
	}
	if pl.ledgerInfo.NextValidatorSet == nil {
		return nil, errors.New("empty validator set")
	}
	vv := &ValidatorVerifier{}
	vv.FromValidatorSet(pl.ledgerInfo.NextValidatorSet, pl.ledgerInfo.Epoch+1)
	return vv, nil
}

// VerifyConsistency verifies a new version of ledger is consistent with a known version
// (and the frozen subtrees at that version).
//
// If successful, it outputs the new accumulator states (i.e. numLeaves and subtrees).
func (pl *ProvenLedgerInfo) VerifyConsistency(numLeaves uint64, oldSubtrees, newSubtrees []HashValue) (uint64, []HashValue, error) {
	acc1 := accumulator.Accumulator{
		Hasher:             sha3libra.NewTransactionAccumulator(),
		FrozenSubtreeRoots: cloneSubtrees(oldSubtrees),
		NumLeaves:          numLeaves,
	}
	err := acc1.AppendSubtrees(newSubtrees, pl.ledgerInfo.Version+1-numLeaves)
	if err != nil {
		return 0, nil, fmt.Errorf("append subtree error: %s", err)
	}
	hash, err := acc1.RootHash()
	if err != nil {
		return 0, nil, fmt.Errorf("new accumulator invalid: %s", err)
	}
	if !sha3libra.Equal(hash, pl.ledgerInfo.TransactionAccumulatorHash) {
		return 0, nil, errors.New("hash mismatch, ledger not consistent")
	}
	return acc1.NumLeaves, acc1.FrozenSubtreeRoots, nil
}
