package types

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types/proof"
	"github.com/the729/go-libra/types/proof/accumulator"
	"github.com/the729/lcs"
)

type EventKey []byte

type EventHandle struct {
	Count uint64
	Key   EventKey
}

type ContractEvent struct {
	Value isContractEvent `lcs:"enum=isContractEvent"`
}

type isContractEvent interface {
	Clone() isContractEvent
}

// ContractEventV0 is a output event of transaction
type ContractEventV0 struct {
	Key            EventKey
	SequenceNumber uint64
	TypeTag        TypeTag
	Data           []byte
}

var contractEventEnumDef = []lcs.EnumVariant{
	{
		Name:     "isContractEvent",
		Value:    0,
		Template: (*ContractEventV0)(nil),
	},
}

// EnumTypes defines enum variants for lcs
func (*ContractEvent) EnumTypes() []lcs.EnumVariant { return contractEventEnumDef }

// EventList is a list of events
type EventList []*ContractEvent

// EventProof is a chain of proof that a event is included in the ledger
type EventProof struct {
	// LedgerInfoToTransactionInfoProof is a Merkle Tree accumulator to prove that TransactionInfo
	// is included in the ledger.
	LedgerInfoToTransactionInfoProof *proof.Accumulator

	// TransactionInfo is the info of the transaction that leads to this version of the ledger.
	*TransactionInfo

	// TransactionInfoToEventProof is an accumulator proof from event root hash in TransactionInfo
	// to actual event.
	TransactionInfoToEventProof *proof.Accumulator
}

// EventWithProof is an event with proof
type EventWithProof struct {
	TransactionVersion uint64
	EventIndex         uint64
	Event              *ContractEvent
	Proof              *EventProof
}

// ProvenEvent is an event proven to be included in the ledger.
type ProvenEvent struct {
	proven     bool
	txnVersion uint64
	eventIndex uint64
	event      *ContractEvent
	ledgerInfo *ProvenLedgerInfo
}

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
	e0 := &ContractEventV0{
		Key:            pb.Key,
		SequenceNumber: pb.SequenceNumber,
		Data:           pb.EventData,
	}
	if err := lcs.Unmarshal(pb.TypeTag, &e0.TypeTag); err != nil {
		return err
	}
	e.Value = e0

	return nil
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (e *ContractEvent) Hash() HashValue {
	hasher := sha3libra.NewContractEvent()
	if err := lcs.NewEncoder(hasher).Encode(e); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
}

// Clone deep clones this struct.
func (e *ContractEventV0) Clone() isContractEvent {
	out := &ContractEventV0{}
	out.Key = cloneBytes(e.Key)
	out.SequenceNumber = e.SequenceNumber
	out.Data = cloneBytes(e.Data)
	// out.TypeTag = e.TypeTag.Clone()
	return out
}

// Clone deep clones this struct.
func (e *ContractEvent) Clone() *ContractEvent {
	out := &ContractEvent{}
	out.Value = e.Value.Clone()
	return out
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (el EventList) Hash() HashValue {
	nodeHasher := sha3libra.NewEventAccumulator()
	acc := accumulator.Accumulator{Hasher: nodeHasher}
	for _, e := range el {
		acc.AppendOne(e.Hash())
	}
	hash, err := acc.RootHash()
	if err != nil {
		panic(err)
	}
	return hash
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

// FromProto parses a protobuf struct into this struct.
func (ep *EventProof) FromProto(pb *pbtypes.EventProof) error {
	var err error
	if pb == nil {
		return ErrNilInput
	}

	ep.LedgerInfoToTransactionInfoProof = &proof.Accumulator{Hasher: sha3libra.NewTransactionAccumulator()}
	err = ep.LedgerInfoToTransactionInfoProof.FromProto(pb.LedgerInfoToTransactionInfoProof)
	if err != nil {
		return err
	}
	ep.TransactionInfo = &TransactionInfo{}
	err = ep.TransactionInfo.FromProto(pb.TransactionInfo)
	if err != nil {
		return err
	}
	ep.TransactionInfoToEventProof = &proof.Accumulator{Hasher: sha3libra.NewEventAccumulator()}
	err = ep.TransactionInfoToEventProof.FromProto(pb.TransactionInfoToEventProof)
	if err != nil {
		return err
	}
	return nil
}

// FromProto parses a protobuf struct into this struct.
func (ep *EventWithProof) FromProto(pb *pbtypes.EventWithProof) error {
	if pb == nil {
		return ErrNilInput
	}
	ep.TransactionVersion = pb.TransactionVersion
	ep.EventIndex = pb.EventIndex
	ep.Event = &ContractEvent{}
	if err := ep.Event.FromProto(pb.Event); err != nil {
		return err
	}
	ep.Proof = &EventProof{}
	if err := ep.Proof.FromProto(pb.Proof); err != nil {
		return err
	}
	return nil
}

// Verify the proof of the event, and output a ProvenEvent if successful.
func (ep *EventWithProof) Verify(provenLedgerInfo *ProvenLedgerInfo) (*ProvenEvent, error) {
	var err error

	eventHash := ep.Event.Hash()
	err = ep.Proof.TransactionInfoToEventProof.Verify(ep.EventIndex, eventHash, ep.Proof.EventRootHash)
	if err != nil {
		return nil, fmt.Errorf("cannot verify event from transaction info: %v", err)
	}

	if ep.TransactionVersion > provenLedgerInfo.GetVersion() {
		return nil, errors.New("event txn version > ledger version")
	}

	err = ep.Proof.LedgerInfoToTransactionInfoProof.Verify(
		ep.TransactionVersion, ep.Proof.TransactionInfo.Hash(),
		provenLedgerInfo.GetTransactionAccumulatorHash(),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot verify transaction info from ledger info: %v", err)
	}

	return &ProvenEvent{
		proven:     true,
		txnVersion: ep.TransactionVersion,
		eventIndex: ep.EventIndex,
		event:      ep.Event.Clone(),
		ledgerInfo: provenLedgerInfo,
	}, nil
}

// GetLedgerInfo returns the ledger info.
func (pe *ProvenEvent) GetLedgerInfo() *ProvenLedgerInfo {
	if !pe.proven {
		panic("not valid proven event")
	}
	return pe.ledgerInfo
}

// GetTransactionVersion returns the transaction version
func (pe *ProvenEvent) GetTransactionVersion() uint64 {
	if !pe.proven {
		panic("not valid proven event")
	}
	return pe.txnVersion
}

// GetEventIndex returns the index of the event in all output events of the transaction
func (pe *ProvenEvent) GetEventIndex() uint64 {
	if !pe.proven {
		panic("not valid proven event")
	}
	return pe.eventIndex
}

// GetEvent returns a copy of the actual event struct
func (pe *ProvenEvent) GetEvent() *ContractEvent {
	if !pe.proven {
		panic("not valid proven event")
	}
	return pe.event.Clone()
}
