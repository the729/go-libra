package types

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/lcs"
)

// Waypoint keeps information about the LedgerInfo on a given reconfiguration.
// A trusted waypoint verifies the LedgerInfo for a certain epoch change.
type Waypoint struct {
	Version uint64
	Value   HashValue
}

type ledger2WaypointConverter struct {
	Epoch            uint64
	RootHash         []byte
	Version          uint64
	TimestampUsec    uint64
	NextValidatorSet *ValidatorSet
}

func (l2wp *ledger2WaypointConverter) Hash() HashValue {
	hasher := sha3libra.NewWaypointLedgerInfo()
	if err := lcs.NewEncoder(hasher).Encode(l2wp); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
}

// FromLedgerInfo builds a waypoint from a ledger info.
// Ledger info should be on the boundary of epochs.
func (wp *Waypoint) FromLedgerInfo(li *LedgerInfo) *Waypoint {
	// Here we do a shallow copy, because ledger2WaypointConverter is local.
	c := ledger2WaypointConverter{
		Epoch:            li.Epoch,
		RootHash:         li.TransactionAccumulatorHash,
		Version:          li.Version,
		TimestampUsec:    li.TimestampUsec,
		NextValidatorSet: li.NextValidatorSet,
	}
	wp.Value = c.Hash()
	wp.Version = c.Version
	return wp
}

// FromProvenLedgerInfo builds a waypoint from a proven ledger info.
func (wp *Waypoint) FromProvenLedgerInfo(pli *ProvenLedgerInfo) *Waypoint {
	return wp.FromLedgerInfo(pli.ledgerInfo)
}

// MarshalText outputs a text representation of this waypoint, in
// the format of version:ledger_info_hash
func (wp *Waypoint) MarshalText() (text []byte, err error) {
	text = strconv.AppendUint([]byte{}, wp.Version, 10)
	text = append(text, ':')
	h := make([]byte, len(wp.Value)*2)
	hex.Encode(h, wp.Value)
	text = append(text, h...)
	return text, nil
}

// UnmarshalText parses a text representation of waypoint into
// the receiver struct. The format is version:ledger_info_hash
func (wp *Waypoint) UnmarshalText(text []byte) error {
	sep := bytes.IndexByte(text, ':')
	if sep < 0 {
		return errors.New("waypoint invalid format")
	}
	v, err := strconv.ParseUint(string(text[0:sep]), 10, 64)
	if err != nil {
		return fmt.Errorf("waypoint invalid version: %v", err)
	}
	h, err := hex.DecodeString(string(text[sep+1:]))
	if err != nil {
		return fmt.Errorf("waypoint invalid hash value: %v", err)
	}
	wp.Version, wp.Value = v, h
	return nil
}

// Verify whether the given ledger info matches this waypoint.
func (wp *Waypoint) Verify(li *LedgerInfoWithSignatures) error {
	li0 := li.Value.(*LedgerInfoWithSignaturesV0)
	if wp.Version != li0.Version {
		return errors.New("waypoint version mismatch")
	}
	wp1 := &Waypoint{}
	wp1.FromLedgerInfo(li0.LedgerInfo)
	if !sha3libra.Equal(wp.Value, wp1.Value) {
		return errors.New("waypoint hash value mismatch")
	}
	return nil
}

// EpochChangeVerificationRequired always returns true.
func (wp *Waypoint) EpochChangeVerificationRequired(_ uint64) bool {
	return true
}
