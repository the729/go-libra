package proof

import (
	"errors"

	"github.com/the729/go-libra/common/bitmap"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
)

type LeafNode struct {
	Key       sha3libra.HashValue
	ValueHash sha3libra.HashValue
}

type SparseMerkle struct {
	leaf     *LeafNode
	siblings []sha3libra.HashValue
}

type InternalNode struct {
	Left  sha3libra.HashValue
	Right sha3libra.HashValue
}

func (n *LeafNode) Hash() sha3libra.HashValue {
	if n == nil {
		return sha3libra.SparseMerklePlaceholderHash
	}
	hasher := sha3libra.NewSparseMerkleLeaf()
	hasher.Write(n.Key)
	hasher.Write(n.ValueHash)
	return hasher.Sum([]byte{})
}

func (n *InternalNode) Hash() sha3libra.HashValue {
	if n == nil {
		return sha3libra.SparseMerklePlaceholderHash
	}
	hasher := sha3libra.NewSparseMerkleLeaf()
	hasher.Write(n.Left)
	hasher.Write(n.Right)
	return hasher.Sum([]byte{})
}

func (m *SparseMerkle) FromProto(pb *pbtypes.SparseMerkleProof) error {
	m.leaf = nil
	if pb.Leaf != nil {
		if len(pb.Leaf) != 2*sha3libra.HashSize {
			return errors.New("leaf wrong length")
		}
		m.leaf = &LeafNode{
			Key:       pb.Leaf[0:sha3libra.HashSize],
			ValueHash: pb.Leaf[sha3libra.HashSize:],
		}
	}

	bm := bitmap.NewFromByteSlice(pb.Bitmap)
	if bm.TrailingZeros() >= 8 {
		return errors.New("bitmap last byte should not be 0")
	}
	if len(pb.NonDefaultSiblings) != bm.OnesCount() {
		return errors.New("wrong count of non-default siblings")
	}

	// Iterate from the MSB of the first byte to the rightmost 1-bit in the bitmap. If a bit is
	// set, the corresponding sibling is non-default and we take the sibling from
	// proto_siblings. Otherwise the sibling on this position is default.

	// Since we already checked len(NonDefaultSiblings), we can safely stop loop when we consumed
	// all NonDefaultSiblings.
	siblings := make([]sha3libra.HashValue, 0, bm.Cap())
	for i, j := bm.Bits(), 0; i.Next() && j < bm.OnesCount(); {
		if _, b := i.Bit(); b {
			siblings = append(siblings, pb.NonDefaultSiblings[j])
			j++
		} else {
			siblings = append(siblings, sha3libra.SparseMerklePlaceholderHash)
		}
	}
	m.siblings = siblings
	return nil
}

func (m *SparseMerkle) VerifyInclusion(elem *LeafNode, expectedRootHash sha3libra.HashValue) error {
	if m.leaf == nil {
		return errors.New("leaf is empty, cannot prove inclusion")
	}
	if !sha3libra.Equal(elem.Key, m.leaf.Key) || !sha3libra.Equal(elem.ValueHash, m.leaf.ValueHash) {
		return errors.New("mismatch element and leaf")
	}
	return m.verify(elem.Key, expectedRootHash)
}

func (m *SparseMerkle) VerifyNonInclusion(elemKey, expectedRootHash sha3libra.HashValue) error {
	if m.leaf != nil {
		if sha3libra.Equal(elemKey, m.leaf.Key) {
			return errors.New("key exists in proof")
		}
		commonBits := 0
		for i, j := bitmap.NewFromByteSlice(elemKey).Bits(), bitmap.NewFromByteSlice(m.leaf.Key).Bits(); i.Next() && j.Next(); {
			_, b1 := i.Bit()
			_, b2 := j.Bit()
			if b1 != b2 {
				break
			}
			commonBits++
		}
		if commonBits < len(m.siblings) {
			return errors.New("key would not have ended up in the subtree where the provided key in proof is the only existing key, if it existed")
		}
	}
	return m.verify(elemKey, expectedRootHash)
}

func (m *SparseMerkle) verify(elemKey, expectedRootHash sha3libra.HashValue) error {
	bm := bitmap.NewFromByteSlice(elemKey)
	if bm.Cap() != sha3libra.HashSize*8 {
		return errors.New("wrong element key size")
	}
	if bm.Cap() < len(m.siblings) {
		return errors.New("merkle tree proof has too many siblings")
	}

	// log.Printf("target hash: %s", hex.EncodeToString(expectedRootHash))
	hash := m.leaf.Hash()
	// log.Printf("initial hash: %s", hex.EncodeToString(hash))
	for i := bm.BitsRev(); i.Next(); {
		idx, b := i.Bit()
		if idx < bm.Cap()-len(m.siblings) {
			// skip bits after len(siblings)
			continue
		}
		hasher := sha3libra.NewSparseMerkleInternal()
		if b {
			// log.Printf("hash: %s with left sibling %s", hex.EncodeToString(hash), hex.EncodeToString(m.siblings[bm.Cap()-idx-1]))
			hasher.Write(m.siblings[bm.Cap()-idx-1])
			hasher.Write(hash)
		} else {
			// log.Printf("hash: %s with right sibling %s", hex.EncodeToString(hash), hex.EncodeToString(m.siblings[bm.Cap()-idx-1]))
			hasher.Write(hash)
			hasher.Write(m.siblings[bm.Cap()-idx-1])
		}
		hash = hasher.Sum([]byte{})
	}
	if !sha3libra.Equal(hash, expectedRootHash) {
		return errors.New("root hashes do not match")
	}
	return nil
}
