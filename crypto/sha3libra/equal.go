package sha3libra

import "bytes"

func Equal(h1, h2 HashValue) bool {
	if len(h1) != HashSize || len(h2) != HashSize {
		return false
	}
	return bytes.Compare(h1, h2) == 0
}
