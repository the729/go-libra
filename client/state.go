package client

import (
	"encoding/hex"
	"fmt"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/types"
	"github.com/the729/go-libra/types/proof/accumulator"
)

// HashValue is a wrap of []byte. It implementes encoding.TextMarshaler and
// encoding.TextUnmarshaler.
type HashValue []byte

// UnmarshalText implements encoding.TextUnmarshaler for HashValue
func (h *HashValue) UnmarshalText(txt []byte) error {
	data, err := hex.DecodeString(string(txt))
	if err != nil {
		return fmt.Errorf("hash value decode error: %v", err)
	}
	*h = data
	return nil
}

// MarshalText implements encoding.TextMarshaler for HashValue
func (h HashValue) MarshalText() (text []byte, err error) {
	return []byte(hex.EncodeToString(h)), nil
}

// State represents the state of a client.
type State struct {
	Waypoint     string                 `toml:"waypoint" json:"waypoint"`
	VSScheme     string                 `toml:"validator_set_scheme" json:"validator_set_scheme,omitempty"`
	ValidatorSet []*types.ValidatorInfo `toml:"validator_set" json:"validator_set,omitempty"`
	Epoch        uint64                 `toml:"epoch" json:"epoch"`
	KnownVersion uint64                 `toml:"known_version" json:"known_version"`
	Subtrees     []HashValue            `toml:"subtrees" json:"subtrees"`
}

// GetState returns the current state of a client.
func (c *Client) GetState() *State {
	c.accMu.RLock()
	defer c.accMu.RUnlock()

	cs := &State{}
	cs.Waypoint = c.lastWaypoint
	if vv, ok := c.verifier.(*types.ValidatorVerifier); ok {
		var vs *types.ValidatorSet
		vs, cs.Epoch = vv.ToValidatorSet()
		cs.ValidatorSet = vs.Payload
		switch vs.Scheme.(type) {
		case types.SchemeED25519:
			cs.VSScheme = "ed25519"
		}
	}
	cs.KnownVersion = c.acc.NumLeaves - 1
	cs.Subtrees = cloneSubtrees2(c.acc.FrozenSubtreeRoots)

	return cs
}

// SetState restores a client to a given state.
func (c *Client) SetState(cs *State) error {
	var verifier types.LedgerInfoVerifier
	if cs.ValidatorSet != nil {
		vv := &types.ValidatorVerifier{}
		vs := &types.ValidatorSet{
			Payload: cs.ValidatorSet,
		}
		switch cs.VSScheme {
		case "ed25519":
			vs.Scheme = types.SchemeED25519{}
		}
		if err := vv.FromValidatorSet(vs, cs.Epoch); err != nil {
			return fmt.Errorf("restore validator set error: %v", err)
		}
		verifier = vv
	} else if cs.Waypoint == "insecure" {
		verifier = &types.ValidatorVerifier{}
		println("Warning: INSECURE! No waypoint or validator set specified.")
	} else {
		wp := &types.Waypoint{}
		if err := wp.UnmarshalText([]byte(cs.Waypoint)); err != nil {
			return fmt.Errorf("restore waypoint error: %v", err)
		}
		verifier = wp
	}

	acc := &accumulator.Accumulator{
		Hasher:    sha3libra.NewTransactionAccumulator(),
		NumLeaves: cs.KnownVersion + 1,
	}
	if cs.Subtrees == nil {
		if wp, ok := verifier.(*types.Waypoint); ok {
			acc.NumLeaves = wp.Version + 1
		}
	} else {
		acc.FrozenSubtreeRoots = cloneSubtrees1(cs.Subtrees)
		if _, err := acc.RootHash(); err != nil {
			return fmt.Errorf("known accumulator invalid: %v", err)
		}
	}

	c.accMu.Lock()
	defer c.accMu.Unlock()
	c.acc = acc
	c.verifier = verifier
	c.lastWaypoint = cs.Waypoint
	return nil
}

func cloneSubtrees(in [][]byte) [][]byte {
	if in == nil {
		return nil
	}
	out := make([][]byte, 0, len(in))
	for _, h := range in {
		h1 := make([]byte, len(h))
		copy(h1, h)
		out = append(out, h1)
	}
	return out
}

func cloneSubtrees1(in []HashValue) [][]byte {
	if in == nil {
		return nil
	}
	out := make([][]byte, 0, len(in))
	for _, h := range in {
		h1 := make([]byte, len(h))
		copy(h1, h)
		out = append(out, h1)
	}
	return out
}

func cloneSubtrees2(in [][]byte) []HashValue {
	if in == nil {
		return nil
	}
	out := make([]HashValue, 0, len(in))
	for _, h := range in {
		h1 := make([]byte, len(h))
		copy(h1, h)
		out = append(out, HashValue(h1))
	}
	return out
}
