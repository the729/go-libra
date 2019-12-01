package proof

import (
	"errors"
	"hash"

	"github.com/the729/go-libra/common/bitmap"
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
)

// Accumulator is a proof of a single element's existance in a Merkle tree accumulator.
type Accumulator struct {
	Hasher   hash.Hash
	Siblings []sha3libra.HashValue
}

// FromProto parses a protobuf struct into this struct, and fills all placeholder
// siblings with placeholder hash.
func (a *Accumulator) FromProto(pb *pbtypes.AccumulatorProof) error {
	a.Siblings = siblingsWithPlaceholder(pb.Siblings, sha3libra.AccumulatorPlaceholderHash)
	return nil
}

// siblingsWithPlaceholder fills in placeholders and forms a dense list of hash siblings
func siblingsWithPlaceholder(pbSiblings [][]byte, placeholder []byte) []sha3libra.HashValue {
	siblings := make([]sha3libra.HashValue, 0, len(pbSiblings))
	for _, sibling := range pbSiblings {
		if len(sibling) == 0 {
			siblings = append(siblings, placeholder)
		} else {
			siblings = append(siblings, sibling)
		}
	}
	return siblings
}

// Verify an element exists in a Merkle tree accumulator.
//
// Arguments:
//  - elemIndex: index of the element.
//  - elemHash: hash of the element.
//  - expectedRootHash: expected root hash of the Merkle tree accumulator.
func (a *Accumulator) Verify(elemIndex uint64, elemHash, expectedRootHash sha3libra.HashValue) error {
	if a.Hasher == nil {
		return errors.New("nil hasher")
	}

	bm := bitmap.NewFromUint64(elemIndex)
	if bm.Cap() < len(a.Siblings) {
		return errors.New("merkle tree proof has too many siblings")
	}

	// log.Printf("target hash: %s", hex.EncodeToString(expectedRootHash))
	hash := elemHash
	// log.Printf("initial hash: %s", hex.EncodeToString(hash))
	hasher := a.Hasher
	for i := bm.BitsRev(); i.Next(); {
		idx, b := i.Bit()
		if idx >= len(a.Siblings) {
			break
		}
		hasher.Reset()
		if b {
			hasher.Write(a.Siblings[idx])
			hasher.Write(hash)
		} else {
			hasher.Write(hash)
			hasher.Write(a.Siblings[idx])
		}
		hash = hasher.Sum([]byte{})
		// log.Printf("new hash: %s", hex.EncodeToString(hash))
	}
	if !sha3libra.Equal(hash, expectedRootHash) {
		return errors.New("root hashes do not match")
	}
	return nil
}
