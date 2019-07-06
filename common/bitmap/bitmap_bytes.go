package bitmap

import (
	"math/bits"
)

type bitmapBytes []byte

type bytesIter struct {
	data bitmapBytes
	pos  int
	rev  bool
}

func NewFromByteSlice(data []byte) Bitmap {
	return bitmapBytes(data)
}

func (b bitmapBytes) Cap() int {
	return 8 * len(b)
}

func (b bitmapBytes) LeadingZeros() int {
	n := 0
	for _, d := range b {
		n += bits.LeadingZeros8(d)
		if d != 0 {
			break
		}
	}
	return n
}

func (b bitmapBytes) TrailingZeros() int {
	n := 0
	for i := len(b) - 1; i >= 0; i-- {
		d := b[i]
		n += bits.TrailingZeros8(d)
		if d != 0 {
			break
		}
	}
	return n
}

func (b bitmapBytes) OnesCount() int {
	n := 0
	for _, d := range b {
		n += bits.OnesCount8(d)
	}
	return n
}

func (b bitmapBytes) Bits() Iter {
	return &bytesIter{b, -1, false}
}
func (b bitmapBytes) BitsRev() Iter {
	return &bytesIter{b, -1, true}
}

func (iter *bytesIter) Bit() (idx int, bit bool) {
	idx = iter.pos
	if idx >= 0 && idx < iter.data.Cap() {
		var byteIdx, bitOffset int
		if iter.rev {
			byteIdx = len(iter.data) - 1 - idx/8
			bitOffset = idx % 8
		} else {
			byteIdx = idx / 8
			bitOffset = 7 - idx%8
		}
		bit = (iter.data[byteIdx] & (uint8(1) << uint(bitOffset))) > 0
	}
	return
}

func (iter *bytesIter) Next() bool {
	if iter.pos >= iter.data.Cap()-1 {
		return false
	}
	iter.pos++
	return true
}
