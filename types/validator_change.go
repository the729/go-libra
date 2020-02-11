package types

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/generated/pbtypes"
)

// ValidatorChangeProof is a vector of LedgerInfo with contiguous increasing
// epoch numbers to prove a sequence of validator changes from the first
// LedgerInfo's epoch.
type ValidatorChangeProof struct {
	LedgerInfoWithSigs []*LedgerInfoWithSignatures
	More               bool
}

type ProvenValidatorChange struct {
	proven         bool
	lastLedgerInfo *LedgerInfo
	genesisHash    []byte
}

// FromProto parses a protobuf struct into this struct.
func (vcp *ValidatorChangeProof) FromProto(pb *pbtypes.ValidatorChangeProof) error {
	var lis []*LedgerInfoWithSignatures
	for _, pbli := range pb.LedgerInfoWithSigs {
		li := &LedgerInfoWithSignatures{}
		if err := li.FromProto(pbli); err != nil {
			return err
		}
		lis = append(lis, li)
	}
	vcp.LedgerInfoWithSigs = lis
	vcp.More = pb.More
	return nil
}

// Verify the ValidatorChangeProof, which is a series of LedgerInfo.
func (vcp *ValidatorChangeProof) Verify(v LedgerInfoVerifier) (*ProvenValidatorChange, error) {
	if len(vcp.LedgerInfoWithSigs) == 0 {
		return nil, errors.New("empty validator change")
	}
	var genesisHash []byte
	for _, li := range vcp.LedgerInfoWithSigs {
		if err := v.Verify(li); err != nil {
			return nil, fmt.Errorf("some ledger info failed to verify: %v", err)
		}
		if li.Version == 0 {
			genesisHash = li.TransactionAccumulatorHash
		}
		if li.NextValidatorSet == nil {
			return nil, errors.New("ledger info doesn't carry validator set")
		}
		vv := &ValidatorVerifier{}
		if err := vv.FromValidatorSet(li.NextValidatorSet, li.Epoch+1); err != nil {
			return nil, fmt.Errorf("init new validator error: %v", err)
		}
		v = vv
	}
	return &ProvenValidatorChange{
		proven:         true,
		lastLedgerInfo: vcp.LedgerInfoWithSigs[len(vcp.LedgerInfoWithSigs)-1].LedgerInfo.Clone(),
		genesisHash:    cloneBytes(genesisHash),
	}, nil
}

// GetLastLedgerInfo returns the last (and latest) ProvenLedgerInfo.
func (pvc *ProvenValidatorChange) GetLastLedgerInfo() *ProvenLedgerInfo {
	if !pvc.proven {
		panic("not valid proven validator change")
	}
	return &ProvenLedgerInfo{
		proven:     true,
		ledgerInfo: pvc.lastLedgerInfo,
	}
}

// GetGenesisHash returns the genesis hash (if extracted from version 0)
func (pvc *ProvenValidatorChange) GetGenesisHash() []byte {
	if !pvc.proven {
		panic("not valid proven validator change")
	}
	return cloneBytes(pvc.genesisHash)
}
