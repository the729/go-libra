package proof

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/common/bitmap"
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
)

type AccumulatorRange struct {
	LeftSiblings  []sha3libra.HashValue
	RightSiblings []sha3libra.HashValue
}

func (r *AccumulatorRange) FromProto(pb *pbtypes.AccumulatorRangeProof) error {
	r.LeftSiblings = siblingsWithPlaceholder(pb.LeftSiblings)
	r.RightSiblings = siblingsWithPlaceholder(pb.RightSiblings)
	return nil
}

func (r *AccumulatorRange) Verify(firstIndex uint64, hashes []sha3libra.HashValue, expectedRootHash sha3libra.HashValue) error {
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
