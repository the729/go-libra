package types

import (
	"io"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
)

// StructTag is a tag to form a resource path.
type StructTag struct {
	Address    AccountAddress
	Module     string
	Name       string
	typeParams []*StructTag
}

// SerializeTo serializes this struct into a io.Writer.
func (t *StructTag) SerializeTo(w io.Writer) error {
	if err := serialization.SimpleSerializer.WriteByteSlice(w, t.Address); err != nil {
		return err
	}
	if err := serialization.SimpleSerializer.WriteByteSlice(w, []byte(t.Module)); err != nil {
		return err
	}
	if err := serialization.SimpleSerializer.WriteByteSlice(w, []byte(t.Name)); err != nil {
		return err
	}

	if err := serialization.SimpleSerializer.Write(w, uint32(len(t.typeParams))); err != nil {
		return err
	}
	for _, v := range t.typeParams {
		if err := v.SerializeTo(w); err != nil {
			return err
		}
	}
	return nil
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (t *StructTag) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewAccessPath()
	t.SerializeTo(hasher)
	return hasher.Sum([]byte{})
}
