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

func (ap *AccessPath) FromProto(pb *pbtypes.AccessPath) error {
	if pb == nil {
		return ErrNilInput
	}
	ap.Address = pb.Address
	ap.Path = pb.Path
	return nil
}

func (ap *AccessPath) SerializeTo(w io.Writer) error {
	if err := ap.Address.SerializeTo(w); err != nil {
		return err
	}
	if err := serialization.SimpleSerializer.Write(w, ap.Path); err != nil {
		return err
	}
	return nil
}
