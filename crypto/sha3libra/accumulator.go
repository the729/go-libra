package sha3libra

import (
	"hash"
	"sync"
)

type accumulatorItem struct {
	hash HashValue
	// level of this item. leaf is 0, counts upwards the root
	level int
}

type accumulator struct {
	sync.RWMutex
	state  []*accumulatorItem
	hasher hash.Hash
}

// NewAccumulator creates a new Merkle Tree Accumulator
// Each `Write()` accumulates to it, and `Sum()` returns the current root hash
func NewAccumulator(merkleNodeHasher hash.Hash) hash.Hash {
	return &accumulator{hasher: merkleNodeHasher}
}

func (a *accumulator) Write(h HashValue) (n int, err error) {
	a.Lock()
	defer a.Unlock()

	a.state = append(a.state, &accumulatorItem{hash: h})

	for len(a.state) > 1 && a.state[len(a.state)-2].level == a.state[len(a.state)-1].level {
		a.mergeItemPair(a.state[len(a.state)-2], a.state[len(a.state)-1], a.state[len(a.state)-2])
		a.state = a.state[:len(a.state)-1]
	}

	// log.Printf("new state after write: %s", spew.Sdump(a.state))
	return len(h), nil
}

func (a *accumulator) mergeItemPair(left, right, out *accumulatorItem) {
	// we do not check left and right are at same level
	// caller must ensure matched level
	a.hasher.Reset()
	a.hasher.Write(left.hash)
	a.hasher.Write(right.hash)
	hash := a.hasher.Sum([]byte{})
	level := left.level + 1
	out.hash, out.level = hash, level
}

func (a *accumulator) Sum(b []byte) HashValue {
	a.RLock()
	defer a.RUnlock()

	if len(a.state) == 0 {
		return append(b, AccumulatorPlaceholderHash...)
	}
	placeHolderItem := &accumulatorItem{hash: AccumulatorPlaceholderHash}

	// initialize: working item set to copy of last state item
	item := &accumulatorItem{
		hash:  a.state[len(a.state)-1].hash,
		level: a.state[len(a.state)-1].level,
	}
	for i := len(a.state) - 2; i >= 0; i-- {
		sibling := a.state[i]
		for item.level < sibling.level {
			a.mergeItemPair(item, placeHolderItem, item)
		}
		a.mergeItemPair(sibling, item, item)
	}

	return append(b, item.hash...)
}

func (a *accumulator) Reset() {
	a.state = nil
}

func (a *accumulator) Size() int {
	return HashSize
}

func (a *accumulator) BlockSize() int {
	return HashSize
}
