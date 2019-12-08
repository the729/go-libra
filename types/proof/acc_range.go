package proof

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/common/bitmap"
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
)

// AccumulatorRange is a proof that a consecutive list of elements exist in a
// Merkle tree accumulator.
type AccumulatorRange struct {
	LeftSiblings  []HashValue // hash siblings of the first element
	RightSiblings []HashValue // hash siblings of the last element
}

// FromProto parses a protobuf struct into this struct, and fills all placeholder
// siblings with placeholder hash.
func (r *AccumulatorRange) FromProto(pb *pbtypes.AccumulatorRangeProof) error {
	r.LeftSiblings = siblingsWithPlaceholder(pb.LeftSiblings, sha3libra.AccumulatorPlaceholderHash)
	r.RightSiblings = siblingsWithPlaceholder(pb.RightSiblings, sha3libra.AccumulatorPlaceholderHash)
	return nil
}

// Verify that a consecutive list of elements exist in a Merkle tree accumulator.
//
// Arguments:
//  - firstIndex: index of the first element.
//  - hashes: hashes of the consecutive list of elements. len(hashes) determines
//    the number of elements.
//  - expectedRootHash: expected root hash of the Merkle tree accumulator.
func (r *AccumulatorRange) Verify(firstIndex uint64, hashes []HashValue, expectedRootHash HashValue) error {
	if len(hashes) == 0 {
		if len(r.LeftSiblings) == 0 && len(r.RightSiblings) == 0 {
			return nil
		}
		return errors.New("empty range to verify, expecting nil first and last proofs")
	}

	lastIndex := firstIndex + uint64(len(hashes)) - 1

	firstBitmap := bitmap.NewFromUint64(firstIndex)
	lastBitmap := bitmap.NewFromUint64(lastIndex)

	leftSiblings := r.LeftSiblings
	rightSiblings := r.RightSiblings

	hasher := sha3libra.NewTransactionAccumulator()
	for firstIter, lastIter := firstBitmap.BitsRev(), lastBitmap.BitsRev(); firstIter.Next() && lastIter.Next(); {
		_, fBit := firstIter.Bit()
		_, lBit := lastIter.Bit()

		if fBit {
			// prepend to hashes
			hashes = append(hashes, nil)
			copy(hashes[1:], hashes)
			hashes[0] = leftSiblings[0]
			leftSiblings = leftSiblings[1:]
		}
		if !lBit {
			hashes = append(hashes, rightSiblings[0])
			rightSiblings = rightSiblings[1:]
		}

		// in-place update new hashes from pairs of hashes
		for i := 0; i < len(hashes)/2; i++ {
			hasher.Reset()
			hasher.Write(hashes[i*2])
			hasher.Write(hashes[i*2+1])
			hashes[i] = hasher.Sum(hashes[i][:0])
		}
		hashes = hashes[:len(hashes)/2]

		if len(leftSiblings) == 0 && len(rightSiblings) == 0 {
			break
		}
	}

	if len(hashes) != 1 {
		return fmt.Errorf("unexpected length: len(hashes)=%d", len(hashes))
	}
	if !sha3libra.Equal(hashes[0], expectedRootHash) {
		return errors.New("root hashes do not match")
	}

	return nil
}
