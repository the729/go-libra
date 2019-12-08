package accumulator

import (
	"encoding/binary"
	"fmt"
	"hash"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/the729/go-libra/crypto/sha3libra"
)

func getTestHash(idx int) HashValue {
	v := make([]byte, sha3libra.HashSize)
	binary.LittleEndian.PutUint64(v, uint64(idx))
	return v
}

// buildTestTree build a full merkle tree with numNode nodes of test hashes.
// Returns a 2-d slice of hashes. 1st dim is tree level, from leaf to root.
// 2nd dim is tree hashes.
func buildTestTree(numNode int, hasher hash.Hash) [][]HashValue {
	tree := make([][]HashValue, 0, 0)

	level := make([]HashValue, 0, numNode)
	// generate leaf nodes
	for i := 0; i < numNode; i++ {
		level = append(level, getTestHash(i))
	}
	tree = append(tree, level)

	for numNode > 1 {
		lastLevel := level
		level = make([]HashValue, 0, (numNode+1)/2)
		for i := 0; i < numNode; i += 2 {
			hasher.Reset()
			hasher.Write(lastLevel[i])
			if i+1 < numNode {
				hasher.Write(lastLevel[i+1])
			} else {
				hasher.Write(sha3libra.AccumulatorPlaceholderHash)
			}
			level = append(level, hasher.Sum([]byte{}))
		}
		tree = append(tree, level)
		numNode = len(level)
	}

	return tree
}

func testAppendOne(t *testing.T, numNode int, subtreeIdxes [][]int) {
	acc := Accumulator{
		Hasher: sha3libra.NewTransactionAccumulator(),
	}
	tree := buildTestTree(numNode, acc.Hasher)

	for i := 0; i < numNode; i++ {
		err := acc.AppendOne(getTestHash(i))
		assert.NoError(t, err)
	}

	expSubtrees := make([]HashValue, 0, 0)
	for _, subtreeIdx := range subtreeIdxes {
		expSubtrees = append(expSubtrees, tree[subtreeIdx[0]][subtreeIdx[1]])
	}

	assert.Equal(t, expSubtrees, acc.FrozenSubtreeRoots, "subtrees should match.")
}

func TestAppendOne(t *testing.T) {
	t.Run("1", func(t *testing.T) { testAppendOne(t, 1, [][]int{{0, 0}}) })
	t.Run("2", func(t *testing.T) { testAppendOne(t, 2, [][]int{{1, 0}}) })
	t.Run("3", func(t *testing.T) { testAppendOne(t, 3, [][]int{{1, 0}, {0, 2}}) })
	t.Run("5", func(t *testing.T) { testAppendOne(t, 5, [][]int{{2, 0}, {0, 4}}) })
	t.Run("7", func(t *testing.T) { testAppendOne(t, 7, [][]int{{2, 0}, {1, 2}, {0, 6}}) })
	t.Run("9", func(t *testing.T) { testAppendOne(t, 9, [][]int{{3, 0}, {0, 8}}) })
	t.Run("10", func(t *testing.T) { testAppendOne(t, 10, [][]int{{3, 0}, {1, 4}}) })
	t.Run("16", func(t *testing.T) { testAppendOne(t, 16, [][]int{{4, 0}}) })
}

func testRootHash(t *testing.T, numNode int) {
	acc := Accumulator{
		Hasher: sha3libra.NewTransactionAccumulator(),
	}
	tree := buildTestTree(numNode, acc.Hasher)

	for i := 0; i < numNode; i++ {
		err := acc.AppendOne(getTestHash(i))
		assert.NoError(t, err)
	}
	root, err := acc.RootHash()
	assert.NoError(t, err)
	expRoot := tree[len(tree)-1][0]
	assert.Equal(t, expRoot, root, "root should match.")
}

func TestHashRoot(t *testing.T) {
	for i := 1; i < 100; i++ {
		name := fmt.Sprintf("%d nodes", i)
		t.Run(name, func(t *testing.T) { testRootHash(t, i) })
	}
}

func testAppendSubtrees(t *testing.T, numNode, numNewNodes int, newSubtreeIdxes [][]int) {
	acc := Accumulator{
		Hasher: sha3libra.NewTransactionAccumulator(),
	}
	tree := buildTestTree(numNode+numNewNodes, acc.Hasher)

	for i := 0; i < numNode; i++ {
		err := acc.AppendOne(getTestHash(i))
		assert.NoError(t, err)
	}

	newSubtrees := make([]HashValue, 0, 0)
	for _, subtreeIdx := range newSubtreeIdxes {
		newSubtrees = append(newSubtrees, tree[subtreeIdx[0]][subtreeIdx[1]])
	}
	err := acc.AppendSubtrees(newSubtrees, uint64(numNewNodes))
	assert.NoError(t, err)

	root, err := acc.RootHash()
	assert.NoError(t, err)
	expRoot := tree[len(tree)-1][0]
	assert.Equal(t, expRoot, root, "root should match.")
}

func TestAppendSubtrees(t *testing.T) {
	t.Run("1+0", func(t *testing.T) { testAppendSubtrees(t, 1, 0, [][]int{}) })
	t.Run("1+1", func(t *testing.T) { testAppendSubtrees(t, 1, 1, [][]int{{0, 1}}) })
	t.Run("1+2", func(t *testing.T) { testAppendSubtrees(t, 1, 2, [][]int{{0, 1}, {0, 2}}) })
	t.Run("1+3", func(t *testing.T) { testAppendSubtrees(t, 1, 3, [][]int{{0, 1}, {1, 1}}) })
	t.Run("1+4", func(t *testing.T) { testAppendSubtrees(t, 1, 4, [][]int{{0, 1}, {1, 1}, {0, 4}}) })
	t.Run("6+1", func(t *testing.T) { testAppendSubtrees(t, 6, 1, [][]int{{0, 6}}) })
	t.Run("6+2", func(t *testing.T) { testAppendSubtrees(t, 6, 2, [][]int{{1, 3}}) })
	t.Run("6+3", func(t *testing.T) { testAppendSubtrees(t, 6, 3, [][]int{{1, 3}, {0, 8}}) })
	t.Run("6+6", func(t *testing.T) { testAppendSubtrees(t, 6, 6, [][]int{{1, 3}, {2, 2}}) })
	t.Run("6+10", func(t *testing.T) { testAppendSubtrees(t, 6, 10, [][]int{{1, 3}, {3, 1}}) })
}
