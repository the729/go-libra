package serialization

import (
	"encoding/binary"
)

var SimpleDeserializer = simpleDeserializer{binary.LittleEndian}

type simpleDeserializer struct {
	binary.ByteOrder
}

func (d simpleDeserializer) ByteSlice(b []byte) ([]byte, error) {
	if len(b) < 4 {
		return nil, ErrWrongSize
	}
	l := d.Uint32(b[:4])
	if uint32(len(b[4:])) < l {
		return nil, ErrWrongSize
	}
	r := make([]byte, l)
	copy(r, b[4:])
	return r, nil
}
