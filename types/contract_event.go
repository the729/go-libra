package types

import (
	"io"

	"github.com/the729/go-libra/crypto/sha3libra"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/generated/pbtypes"
)

// ContractEvent is a output event of transaction
type ContractEvent struct {
	AccessPath     *AccessPath
	SequenceNumber uint64
	Data           []byte
}

// EventList is a list of events
type EventList []*ContractEvent

// FromProto parses a protobuf struct into this struct.
func (e *ContractEvent) FromProto(pb *pbtypes.Event) error {
	if pb == nil {
		return ErrNilInput
	}
	e.AccessPath = &AccessPath{}
	if err := e.AccessPath.FromProto(pb.AccessPath); err != nil {
		return err
	}
	e.SequenceNumber = pb.SequenceNumber
	e.Data = pb.EventData

	return nil
}

// SerializeTo serializes this struct into a io.Writer.
func (e *ContractEvent) SerializeTo(w io.Writer) error {
	if err := e.AccessPath.SerializeTo(w); err != nil {
		return err
	}
	if err := serialization.SimpleSerializer.Write(w, e.SequenceNumber); err != nil {
		return err
	}
	if err := serialization.SimpleSerializer.Write(w, e.Data); err != nil {
		return err
	}
	return nil
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (e *ContractEvent) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewContractEvent()
	if err := e.SerializeTo(hasher); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
}

// Clone deep clones this struct.
func (e *ContractEvent) Clone() *ContractEvent {
	out := &ContractEvent{}
	out.AccessPath = e.AccessPath.Clone()
	out.SequenceNumber = e.SequenceNumber
	out.Data = cloneBytes(e.Data)
	return out
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (el EventList) Hash() sha3libra.HashValue {
	nodeHasher := sha3libra.NewEventAccumulator()
	hasher := sha3libra.NewAccumulator(nodeHasher)
	for _, e := range el {
		hasher.Write(e.Hash())
	}
	return hasher.Sum([]byte{})
}

// Clone deep clones this struct.
func (el EventList) Clone() EventList {
	if el == nil {
		return nil
	}
	out := make([]*ContractEvent, 0, len(el))
	for _, e := range el {
		out = append(out, e.Clone())
	}
	return out
}
