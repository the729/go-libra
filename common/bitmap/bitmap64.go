package bitmap

import (
	"math/bits"
)

type bitmap64 uint64

type bitmap64iter struct {
	b   bitmap64
	pos int
}

// NewFromUint64 build a new Bitmap from uint64
func NewFromUint64(b uint64) Bitmap {
	return bitmap64(b)
}

func (b bitmap64) Cap() int {
	return 64
}

func (b bitmap64) LeadingZeros() int {
	return bits.LeadingZeros64(uint64(b))
}

func (b bitmap64) TrailingZeros() int {
	return bits.TrailingZeros64(uint64(b))
}

func (b bitmap64) OnesCount() int {
	return bits.OnesCount64(uint64(b))
}

func (b bitmap64) Bits() Iter {
	return &bitmap64iter{b, -1}
}

func (b bitmap64) BitsRev() Iter {
	return &bitmap64iter{bitmap64(bits.Reverse64(uint64(b))), -1}
}

func (i *bitmap64iter) Bit() (idx int, bit bool) {
	idx = i.pos
	if idx >= 0 && idx < i.b.Cap() {
		bit = (uint64(i.b) & (uint64(1) << uint(i.b.Cap()-1-idx))) > 0
	}
	return
}

func (i *bitmap64iter) Next() bool {
	if i.pos >= i.b.Cap()-1 {
		return false
	}
	i.pos++
	return true
}
