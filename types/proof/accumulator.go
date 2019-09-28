package proof

import (
	"errors"
	"hash"

	"github.com/the729/go-libra/common/bitmap"
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
)

type Accumulator struct {
	Hasher   hash.Hash
	Siblings []sha3libra.HashValue
}

func (a *Accumulator) FromProto(pb *pbtypes.AccumulatorProof) error {
	bm := bitmap.NewFromUint64(pb.Bitmap)
	if len(pb.NonDefaultSiblings) != bm.OnesCount() {
		return errors.New("wrong count of non-default siblings")
	}

	siblings := make([]sha3libra.HashValue, 0, bm.Cap())
	for i, j, seenOne := bm.Bits(), 0, false; i.Next(); {
		if _, b := i.Bit(); b {
			seenOne = true
			siblings = append(siblings, pb.NonDefaultSiblings[j])
			j++
		} else if seenOne {
			siblings = append(siblings, sha3libra.AccumulatorPlaceholderHash)
		}
	}
	a.Siblings = siblings
	return nil
}

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
		hasher.Reset()
		if b {
			hasher.Write(a.Siblings[len(a.Siblings)-idx-1])
			hasher.Write(hash)
		} else {
			hasher.Write(hash)
			hasher.Write(a.Siblings[len(a.Siblings)-idx-1])
		}
		hash = hasher.Sum([]byte{})
		// log.Printf("new hash: %s", hex.EncodeToString(hash))
		if len(a.Siblings)-idx-1 == 0 {
			break
		}
	}
	if !sha3libra.Equal(hash, expectedRootHash) {
		return errors.New("root hashes do not match")
	}
	return nil
}
