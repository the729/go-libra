package client

import (
	"fmt"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/types/proof/accumulator"
)

func (c *Client) SetKnownVersion(knownVersion uint64, subtrees [][]byte) error {
	acc := &accumulator.Accumulator{
		Hasher:             sha3libra.NewTransactionAccumulator(),
		FrozenSubtreeRoots: cloneSubtrees(subtrees),
		NumLeaves:          knownVersion + 1,
	}
	_, err := acc.RootHash()
	if err != nil {
		return fmt.Errorf("known accumulator invalid: %s", err)
	}
	c.accMu.Lock()
	defer c.accMu.Unlock()
	c.acc = acc
	return nil
}

func (c *Client) GetKnownVersion() (knownVersion uint64, subtrees [][]byte) {
	c.accMu.RLock()
	defer c.accMu.RUnlock()
	return c.acc.NumLeaves - 1, cloneSubtrees(c.acc.FrozenSubtreeRoots)
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
