package sha3libra

import (
	"bytes"
	"log"
)

func Equal(h1, h2 HashValue) bool {
	if len(h1) != HashSize || len(h2) != HashSize {
		return false
	}
	if bytes.Compare(h1, h2) != 0 {
		log.Printf("mismatch: %x vs %x", h1, h2)
	}
	return bytes.Compare(h1, h2) == 0
}
