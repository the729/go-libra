package types

import (
	"io"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/generated/pbtypes"
)

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
