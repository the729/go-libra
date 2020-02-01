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
// Returns the last ProvenLedgerInfo.
func (vcp *ValidatorChangeProof) Verify(v LedgerInfoVerifier) (*ProvenLedgerInfo, error) {
	if len(vcp.LedgerInfoWithSigs) == 0 {
		return nil, errors.New("empty validator change")
	}
	for _, li := range vcp.LedgerInfoWithSigs {
		if err := v.Verify(li); err != nil {
			return nil, fmt.Errorf("some ledger info failed to verify: %v", err)
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
	return &ProvenLedgerInfo{
		proven:     true,
		ledgerInfo: vcp.LedgerInfoWithSigs[len(vcp.LedgerInfoWithSigs)-1].LedgerInfo.Clone(),
	}, nil
}
