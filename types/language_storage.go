package types

import (
	"io"

	"github.com/the729/go-libra/crypto/sha3libra"

	"github.com/the729/go-libra/common/canonical_serialization"
)

type StructTag struct {
	Address    AccountAddress
	Module     string
	Name       string
	typeParams []*StructTag
}

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

func (t *StructTag) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewAccessPath()
	t.SerializeTo(hasher)
	return hasher.Sum([]byte{})
}
