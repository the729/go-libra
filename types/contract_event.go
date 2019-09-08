package types

import (
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/lcs"
)

type EventKey []byte

type EventHandle struct {
	Count uint64
	Key   EventKey
}

// ContractEvent is a output event of transaction
type ContractEvent struct {
	Key            EventKey
	SequenceNumber uint64
	Data           []byte
}

// EventList is a list of events
type EventList []*ContractEvent

// Clone deep clones this struct.
func (eh *EventHandle) Clone() *EventHandle {
	out := &EventHandle{}
	out.Key = cloneBytes(eh.Key)
	out.Count = eh.Count
	return out
}

// FromProto parses a protobuf struct into this struct.
func (e *ContractEvent) FromProto(pb *pbtypes.Event) error {
	if pb == nil {
		return ErrNilInput
	}
	e.Key = pb.Key
	e.SequenceNumber = pb.SequenceNumber
	e.Data = pb.EventData

	return nil
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (e *ContractEvent) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewContractEvent()
	if err := lcs.NewEncoder(hasher).Encode(e); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
}

// Clone deep clones this struct.
func (e *ContractEvent) Clone() *ContractEvent {
	out := &ContractEvent{}
	out.Key = cloneBytes(e.Key)
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
