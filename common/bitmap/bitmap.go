package bitmap

// Bitmap interface
type Bitmap interface {
	// Cap returns the capacity of the bitmap
	Cap() int
	// LeadingZeros returns number of leading zeros
	LeadingZeros() int
	// TrailingZeros returns number of trailing zeros
	TrailingZeros() int
	// OnesCount returns number of ones in the bitmap
	OnesCount() int
	// Bits returns a iterator on all bits from left to right.
	// The total number of bits will be equal to Cap()
	Bits() Iter
	// Bits returns a iterator on all bits from right to left.
	// The total number of bits will be equal to Cap()
	BitsRev() Iter
}

// Iter is an iterator over the bits in a bitmap
type Iter interface {
	// Bit returns current bit and the index of it.
	// index is always counting from 0, 1, ..., regardless of
	// whether this is a reverse iterator
	// If index < 0, the bit is invalid
	Bit() (int, bool)
	// Next moves the iterator to next bit.
	// returns true if successful.
	Next() bool
}
