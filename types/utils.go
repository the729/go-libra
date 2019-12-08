package types

import "github.com/the729/go-libra/crypto/sha3libra"

// HashValue is equivalent to sha3libra.HashValue, which is []byte
type HashValue = sha3libra.HashValue

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	c := make([]byte, len(b))
	copy(c, b)
	return c
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
