package stdscript

import (
	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/types"
)

// SentPaymentEvent is a p2p sent payment event
type SentPaymentEvent struct {
	Payee  types.AccountAddress
	Amount uint64
}

// ReceivedPaymentEvent is a p2p received payment event
type ReceivedPaymentEvent struct {
	Payer  types.AccountAddress
	Amount uint64
}

// UnmarshalBinary unmarshals raw bytes into this struct.
func (ev *SentPaymentEvent) UnmarshalBinary(data []byte) error {
	ev.Amount = serialization.SimpleDeserializer.Uint64(data)
	data = data[8:]
	addr, err := serialization.SimpleDeserializer.ByteSlice(data)
	if err != nil {
		return err
	}
	ev.Payee = addr
	return nil
}

// UnmarshalBinary unmarshals raw bytes into this struct.
func (ev *ReceivedPaymentEvent) UnmarshalBinary(data []byte) error {
	ev.Amount = serialization.SimpleDeserializer.Uint64(data)
	data = data[8:]
	addr, err := serialization.SimpleDeserializer.ByteSlice(data)
	if err != nil {
		return err
	}
	ev.Payer = addr
	return nil
}
