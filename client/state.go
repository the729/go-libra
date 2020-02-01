package client

import (
	"encoding/hex"
	"fmt"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/types"
	"github.com/the729/go-libra/types/proof/accumulator"
)

type HashValue []byte

func (h *HashValue) UnmarshalText(txt []byte) error {
	data, err := hex.DecodeString(string(txt))
	if err != nil {
		return fmt.Errorf("hash value decode error: %v", err)
	}
	*h = data
	return nil
}

func (h HashValue) MarshalText() (text []byte, err error) {
	return []byte(hex.EncodeToString(h)), nil
}

type ClientState struct {
	Waypoint     string             `toml:"waypoint" json:"waypoint"`
	ValidatorSet types.ValidatorSet `toml:"validator_set" json:"validator_set,omitempty"`
	Epoch        uint64             `toml:"epoch" json:"epoch"`
	KnownVersion uint64             `toml:"known_version" json:"known_version"`
	Subtrees     []HashValue        `toml:"subtrees" json:"subtrees"`
}

func (c *Client) GetState() *ClientState {
	c.accMu.RLock()
	defer c.accMu.RUnlock()

	cs := &ClientState{}
	cs.Waypoint = c.lastWaypoint
	if vv, ok := c.verifier.(*types.ValidatorVerifier); ok {
		cs.ValidatorSet, cs.Epoch = vv.ToValidatorSet()
	}
	cs.KnownVersion = c.acc.NumLeaves - 1
	cs.Subtrees = cloneSubtrees2(c.acc.FrozenSubtreeRoots)

	return cs
}

func (c *Client) SetState(cs *ClientState) error {
	var verifier types.LedgerInfoVerifier
	if len(cs.ValidatorSet) > 0 {
		vv := &types.ValidatorVerifier{}
		if err := vv.FromValidatorSet(cs.ValidatorSet, cs.Epoch); err != nil {
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
		Hasher:             sha3libra.NewTransactionAccumulator(),
		FrozenSubtreeRoots: cloneSubtrees1(cs.Subtrees),
		NumLeaves:          cs.KnownVersion + 1,
	}
	if _, err := acc.RootHash(); err != nil {
		return fmt.Errorf("known accumulator invalid: %v", err)
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
