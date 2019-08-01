package serialization

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var SimpleSerializer = simpleSerializer{binary.LittleEndian}

type simpleSerializer struct {
	byteOrder binary.ByteOrder
}

func (s simpleSerializer) WriteByteSlice(w io.Writer, v []byte) error {
	if err := binary.Write(w, s.byteOrder, uint32(len(v))); err != nil {
		return err
	}
	n, err := io.Copy(w, bytes.NewReader(v))
	if err != nil {
		return err
	}
	if int(n) != len(v) {
		return errors.New("wrong size")
	}
	return nil
}

func (s simpleSerializer) Write(w io.Writer, v interface{}) error {
	if vv, ok := v.([]byte); ok {
		return s.WriteByteSlice(w, vv)
	} else if vv, ok := v.(bool); ok {
		v8 := uint8(0)
		if vv {
			v8 = uint8(1)
		}
		return binary.Write(w, s.byteOrder, v8)
	}
	return binary.Write(w, s.byteOrder, v)
}
