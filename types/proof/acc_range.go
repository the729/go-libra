package proof

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/common/bitmap"
	"github.com/the729/go-libra/crypto/sha3libra"
)

type AccumulatorRange struct {
	First *Accumulator
	Last  *Accumulator
}

func (r *AccumulatorRange) Verify(firstIndex uint64, hashes []sha3libra.HashValue, expectedRootHash sha3libra.HashValue) error {
	if len(hashes) == 0 {
		if r.First == nil && r.Last == nil {
			return nil
		}
		return errors.New("empty range to verify, expecting nil first and last proofs")
	}
	if r.First == nil {
		return errors.New("nil first proof")
	}
	firstProof := r.First
	lastProof := r.Last
	if len(hashes) == 1 {
		if lastProof != nil {
			// If last is not nil, it should proof the same object as the first
			if err := lastProof.Verify(firstIndex, hashes[0], expectedRootHash); err != nil {
				return err
			}
		}
		lastProof = firstProof
	}
	if len(firstProof.siblings) != len(lastProof.siblings) {
		return errors.New("mismatch first proof and last proof sibling counts")
	}

	lastIndex := firstIndex + uint64(len(hashes)) - 1

	firstBitmap := bitmap.NewFromUint64(firstIndex)
	lastBitmap := bitmap.NewFromUint64(lastIndex)

	hasher := sha3libra.NewTransactionAccumulator()
	for firstIter, lastIter := firstBitmap.BitsRev(), lastBitmap.BitsRev(); firstIter.Next() && lastIter.Next(); {
		fIdx, fBit := firstIter.Bit()
		lIdx, lBit := lastIter.Bit()
		fSibling := firstProof.siblings[len(firstProof.siblings)-fIdx-1]
		lSibling := lastProof.siblings[len(lastProof.siblings)-lIdx-1]

		// log.Printf("hashes:")
		// for _, h := range hashes {
		// 	log.Printf("    %s", hex.EncodeToString(h))
		// }

		if fBit {
			// prepend to hashes
			hashes = append(hashes, nil)
			copy(hashes[1:], hashes)
			hashes[0] = fSibling
		}
		if !lBit {
			hashes = append(hashes, lSibling)
		}
		if !fBit {
			if !sha3libra.Equal(hashes[1], fSibling) {
				return errors.New("first-side sibling hash mismatch")
			}
		}
		if lBit {
			if !sha3libra.Equal(hashes[len(hashes)-2], lSibling) {
				return errors.New("last-side sibling hash mismatch")
			}
		}
		// in-place update new hashes from pairs of hashes
		for i := 0; i < len(hashes)/2; i++ {
			hasher.Reset()
			hasher.Write(hashes[i*2])
			hasher.Write(hashes[i*2+1])
			hashes[i] = hasher.Sum(hashes[i][:0])
		}
		hashes = hashes[:len(hashes)/2]

		if len(firstProof.siblings)-fIdx-1 == 0 {
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
