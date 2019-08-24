package types

import (
	"io"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
)

const (
	CodeTag     = 0
	ResourceTag = 1
)

// AccessPathTag is an interface that can be used to construct the root of an access path
type AccessPathTag interface {
	Hash() sha3libra.HashValue
	TypePrefix() byte
}

// StructTag is a tag to form a resource path.
//
// StructTag implements AccessPathTag interface
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

// Hash outputs the hash of this struct, using the appropriate hash function.
func (t *StructTag) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewAccessPath()
	t.SerializeTo(hasher)
	return hasher.Sum([]byte{})
}

// TypePrefix returns type byte of this tag, which is '0x01'
func (t *StructTag) TypePrefix() byte { return ResourceTag }

// RawTag is a tag with raw hash values. It implements AccessPathTag interface.
type RawTag struct {
	HashVal sha3libra.HashValue
	TypeVal byte
}

// Hash returns HashVal
func (t *RawTag) Hash() sha3libra.HashValue { return t.HashVal }

// TypePrefix returns TypeVal
func (t *RawTag) TypePrefix() byte { return t.TypeVal }
