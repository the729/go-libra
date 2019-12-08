package accumulator

import (
	"errors"
	"hash"
	"math/bits"

	"github.com/the729/go-libra/crypto/sha3libra"
)

const (
	// MaxAccumulatorProofDepth is the max accumulator depth
	MaxAccumulatorProofDepth uint = 63
	// MaxAccumulatorLeaves is the max number of leaves in an accumulator
	MaxAccumulatorLeaves uint64 = 1 << MaxAccumulatorProofDepth
)

// HashValue is equivalent to sha3libra.HashValue, which is []byte
type HashValue = sha3libra.HashValue

// Accumulator is the state of a Merkle tree accumulator.
type Accumulator struct {
	Hasher hash.Hash

	// FrozenSubtreeRoots are hashes of all full subtrees, from left to right.
	FrozenSubtreeRoots []HashValue

	// Total number of leaves
	NumLeaves uint64
}

// RootHash computes root hash of current accumulator
func (a *Accumulator) RootHash() (HashValue, error) {
	if len(a.FrozenSubtreeRoots) == 0 {
		return sha3libra.AccumulatorPlaceholderHash, nil
	}
	if len(a.FrozenSubtreeRoots) == 1 {
		return a.FrozenSubtreeRoots[0], nil
	}
	currHash := make([]byte, sha3libra.HashSize)
	copy(currHash, sha3libra.AccumulatorPlaceholderHash)
	bitmap := a.NumLeaves >> uint(bits.TrailingZeros64(a.NumLeaves))
	subtrees := a.FrozenSubtreeRoots
	for bitmap > 0 {
		a.Hasher.Reset()
		if bitmap&1 != 0 {
			if len(subtrees) == 0 {
				return nil, errors.New("invalid accumulator: too few subtrees")
			}
			a.Hasher.Write(subtrees[len(subtrees)-1])
			a.Hasher.Write(currHash)
			currHash = a.Hasher.Sum(currHash[:0])
			subtrees = subtrees[:len(subtrees)-1]
		} else {
			a.Hasher.Write(currHash)
			a.Hasher.Write(sha3libra.AccumulatorPlaceholderHash)
			currHash = a.Hasher.Sum(currHash[:0])
		}
		bitmap >>= 1
	}
	if len(subtrees) != 0 {
		return nil, errors.New("invalid accumulator: too many subtrees")
	}
	return currHash, nil
}

// AppendOne appends a leaf to the accumulator
func (a *Accumulator) AppendOne(leafHash HashValue) error {
	a.FrozenSubtreeRoots = append(a.FrozenSubtreeRoots, leafHash)

	numTrailingOnes := bits.TrailingZeros64(^a.NumLeaves)
	for i := 0; i < numTrailingOnes; i++ {
		if len(a.FrozenSubtreeRoots) < 2 {
			return errors.New("invalid accumulator")
		}
		a.Hasher.Reset()
		a.Hasher.Write(a.FrozenSubtreeRoots[len(a.FrozenSubtreeRoots)-2])
		a.Hasher.Write(a.FrozenSubtreeRoots[len(a.FrozenSubtreeRoots)-1])
		newHash := a.Hasher.Sum([]byte{})
		a.FrozenSubtreeRoots[len(a.FrozenSubtreeRoots)-2] = newHash
		a.FrozenSubtreeRoots = a.FrozenSubtreeRoots[:len(a.FrozenSubtreeRoots)-1]
	}
	a.NumLeaves++

	return nil
}

// AppendSubtrees appends a list of new leaves to the existing accumulator. The new leaves
// are represented as subtrees.
func (a *Accumulator) AppendSubtrees(subtrees []HashValue, numNewLeaves uint64) error {
	if numNewLeaves > MaxAccumulatorLeaves-a.NumLeaves {
		return errors.New("too many new leaves")
	}
	if len(subtrees) == 0 {
		if numNewLeaves == 0 {
			return nil
		}
		return errors.New("too few subtrees")
	}

	currFrozen := a.FrozenSubtreeRoots
	currNumLeaves := a.NumLeaves
	remainingNewLeaves := numNewLeaves

	for {
		rightmostFrozenSubtreeSize := uint64(1) << uint(bits.TrailingZeros64(currNumLeaves))
		if currNumLeaves == 0 || rightmostFrozenSubtreeSize > remainingNewLeaves {
			break
		}
		if len(subtrees) == 0 {
			return errors.New("too few subtrees")
		}
		currHash := make([]byte, sha3libra.HashSize)
		copy(currHash, subtrees[0])
		subtrees = subtrees[1:]
		for mask := rightmostFrozenSubtreeSize; (currNumLeaves & mask) != 0; mask <<= 1 {
			a.Hasher.Reset()
			a.Hasher.Write(currFrozen[len(currFrozen)-1])
			a.Hasher.Write(currHash)
			currHash = a.Hasher.Sum(currHash[:0])
			currFrozen = currFrozen[:len(currFrozen)-1]
		}
		currFrozen = append(currFrozen, currHash)
		currNumLeaves += rightmostFrozenSubtreeSize
		remainingNewLeaves -= rightmostFrozenSubtreeSize
	}
	currFrozen = append(currFrozen, subtrees...)
	currNumLeaves += remainingNewLeaves

	a.FrozenSubtreeRoots = currFrozen
	a.NumLeaves = currNumLeaves

	return nil
}
