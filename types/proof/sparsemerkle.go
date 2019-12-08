package proof

import (
	"errors"

	"github.com/the729/go-libra/common/bitmap"
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
)

// LeafNode of a sparse Merkle tree.
type LeafNode struct {
	Key       HashValue
	ValueHash HashValue
}

// SparseMerkle is a proof that an element exists in a sparse Merkle tree,
// or an element key does not exist in the tree.
type SparseMerkle struct {
	// If Leaf == nil, this struct can only prove non-existance.
	// Otherwise, this struct can prove existance of the Leaf, or non-existance
	// of other elements.
	Leaf *LeafNode

	// Sibling hashes from root to leaf.
	Siblings []HashValue
}

// type InternalNode struct {
// 	Left  HashValue
// 	Right HashValue
// }

// Hash of the struct.
func (n *LeafNode) Hash() HashValue {
	if n == nil {
		return sha3libra.SparseMerklePlaceholderHash
	}
	hasher := sha3libra.NewSparseMerkleLeaf()
	hasher.Write(n.Key)
	hasher.Write(n.ValueHash)
	return hasher.Sum([]byte{})
}

// func (n *InternalNode) Hash() HashValue {
// 	if n == nil {
// 		return sha3libra.SparseMerklePlaceholderHash
// 	}
// 	hasher := sha3libra.NewSparseMerkleLeaf()
// 	hasher.Write(n.Left)
// 	hasher.Write(n.Right)
// 	return hasher.Sum([]byte{})
// }

// FromProto parses a protobuf struct into this struct.
func (m *SparseMerkle) FromProto(pb *pbtypes.SparseMerkleProof) error {
	m.Leaf = nil
	if pb.Leaf != nil {
		if len(pb.Leaf) != 2*sha3libra.HashSize {
			return errors.New("leaf wrong length")
		}
		m.Leaf = &LeafNode{
			Key:       pb.Leaf[0:sha3libra.HashSize],
			ValueHash: pb.Leaf[sha3libra.HashSize:],
		}
	}

	m.Siblings = siblingsWithPlaceholder(pb.Siblings, sha3libra.SparseMerklePlaceholderHash)
	return nil
}

// VerifyInclusion verifies that an element represented as LeafNode exists in the sparse Merkle tree.
func (m *SparseMerkle) VerifyInclusion(elem *LeafNode, expectedRootHash HashValue) error {
	if m.Leaf == nil {
		return errors.New("leaf is empty, cannot prove inclusion")
	}
	if !sha3libra.Equal(elem.Key, m.Leaf.Key) || !sha3libra.Equal(elem.ValueHash, m.Leaf.ValueHash) {
		return errors.New("mismatch element and leaf")
	}
	return m.verify(elem.Key, expectedRootHash)
}

// VerifyNonInclusion verifies that a given element key does not exist in the sparse Merkle tree.
func (m *SparseMerkle) VerifyNonInclusion(elemKey, expectedRootHash HashValue) error {
	if m.Leaf != nil {
		if sha3libra.Equal(elemKey, m.Leaf.Key) {
			return errors.New("key exists in proof")
		}
		commonBits := 0
		for i, j := bitmap.NewFromByteSlice(elemKey).Bits(), bitmap.NewFromByteSlice(m.Leaf.Key).Bits(); i.Next() && j.Next(); {
			_, b1 := i.Bit()
			_, b2 := j.Bit()
			if b1 != b2 {
				break
			}
			commonBits++
		}
		if commonBits < len(m.Siblings) {
			return errors.New("key would not have ended up in the subtree where the provided key in proof is the only existing key, if it existed")
		}
	}
	return m.verify(elemKey, expectedRootHash)
}

func (m *SparseMerkle) verify(elemKey, expectedRootHash HashValue) error {
	bm := bitmap.NewFromByteSlice(elemKey)
	if bm.Cap() != sha3libra.HashSize*8 {
		return errors.New("wrong element key size")
	}
	if bm.Cap() < len(m.Siblings) {
		return errors.New("merkle tree proof has too many siblings")
	}

	siblings := m.Siblings
	// log.Printf("target hash: %s", hex.EncodeToString(expectedRootHash))
	hash := m.Leaf.Hash()
	// log.Printf("initial hash: %s", hex.EncodeToString(hash))
	for i := bm.BitsRev(); i.Next(); {
		idx, b := i.Bit()
		if idx < bm.Cap()-len(m.Siblings) {
			// skip bits after len(siblings)
			continue
		}
		hasher := sha3libra.NewSparseMerkleInternal()
		if b {
			// log.Printf("%d hash: %s with left sibling %s", idx, hex.EncodeToString(hash), hex.EncodeToString(m.siblings[j]))
			hasher.Write(siblings[0])
			hasher.Write(hash)
		} else {
			// log.Printf("%d hash: %s with right sibling %s", idx, hex.EncodeToString(hash), hex.EncodeToString(m.siblings[j]))
			hasher.Write(hash)
			hasher.Write(siblings[0])
		}
		siblings = siblings[1:]
		hash = hasher.Sum([]byte{})
	}
	// log.Printf("final hash: %s", hex.EncodeToString(hash))
	if !sha3libra.Equal(hash, expectedRootHash) {
		return errors.New("root hashes do not match")
	}
	return nil
}
