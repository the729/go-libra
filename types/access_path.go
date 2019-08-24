package types

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/the729/go-libra/crypto/sha3libra"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/generated/pbtypes"
)

// AccessPath is access path to an event.
type AccessPath struct {
	Address AccountAddress
	Path    []byte
}

// FromProto parses a protobuf struct into this struct.
func (ap *AccessPath) FromProto(pb *pbtypes.AccessPath) error {
	if pb == nil {
		return ErrNilInput
	}
	ap.Address = pb.Address
	ap.Path = pb.Path
	return nil
}

// SerializeTo serializes this struct into a io.Writer.
func (ap *AccessPath) SerializeTo(w io.Writer) error {
	if err := ap.Address.SerializeTo(w); err != nil {
		return err
	}
	if err := serialization.SimpleSerializer.Write(w, ap.Path); err != nil {
		return err
	}
	return nil
}

// Clone deep clones this struct.
func (ap *AccessPath) Clone() *AccessPath {
	out := &AccessPath{}
	out.Address = cloneBytes(ap.Address)
	out.Path = cloneBytes(ap.Path)
	return out
}

// DecodePath decodes raw path byte slice into DecodedPath struct.
func (ap *AccessPath) DecodePath() (*DecodedPath, error) {
	dp := &DecodedPath{}
	if err := dp.UnmarshalBinary(ap.Path); err != nil {
		return nil, fmt.Errorf("decode path error: %v", err)
	}
	return dp, nil
}

// DecodedPath is a decoded path
type DecodedPath struct {
	Tag      AccessPathTag
	Accesses []string
}

// UnmarshalBinary unmarshals raw bytes into this struct.
func (dp *DecodedPath) UnmarshalBinary(data []byte) error {
	if len(data) < 1+sha3libra.HashSize {
		return errors.New("input too short")
	}
	dp.Tag = &RawTag{
		TypeVal: data[0],
		HashVal: data[1 : 1+sha3libra.HashSize],
	}
	data = data[1+sha3libra.HashSize:]
	dp.Accesses = nil
	if len(data) == 0 {
		return nil
	}
	if data[0] != '/' {
		return fmt.Errorf("unexpected char: %v", data[0])
	}
	dp.Accesses = strings.Split(string(data[1:]), "/")
	return nil
}

// SerializeTo serializes this struct into an io.Writer.
func (dp *DecodedPath) SerializeTo(w io.Writer) error {
	if dp.Tag == nil {
		return errors.New("tag is nil")
	}
	if _, err := w.Write(append([]byte{dp.Tag.TypePrefix()}, dp.Tag.Hash()...)); err != nil {
		return err
	}
	for _, a := range dp.Accesses {
		if _, err := w.Write(append([]byte{'/'}, []byte(a)...)); err != nil {
			return err
		}
	}
	return nil
}

// MarshalBinary serializes this struct into a byte slice.
func (dp *DecodedPath) MarshalBinary() (data []byte, err error) {
	var b bytes.Buffer
	if err := dp.SerializeTo(&b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// IsEqual checks whether this DecodedPath is equal to a given raw path
func (dp *DecodedPath) IsEqual(path []byte) bool {
	b, err := dp.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return bytes.Equal(b, path)
}
