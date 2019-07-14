package types

import (
	"io"

	"github.com/the729/go-libra/crypto/sha3libra"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/generated/pbtypes"
)

type ContractEvent struct {
	AccessPath     *AccessPath
	SequenceNumber uint64
	Data           []byte
}

type EventList []*ContractEvent

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

func (e *ContractEvent) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewContractEvent()
	if err := e.SerializeTo(hasher); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
}

func (el EventList) Hash() sha3libra.HashValue {
	nodeHasher := sha3libra.NewEventAccumulator()
	hasher := sha3libra.NewAccumulator(nodeHasher)
	for _, e := range el {
		hasher.Write(e.Hash())
	}
	return hasher.Sum([]byte{})
}
