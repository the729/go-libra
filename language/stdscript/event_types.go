package stdscript

import (
	"github.com/the729/go-libra/types"
	"github.com/the729/lcs"
)

// SentPaymentEvent is a p2p sent payment event
type SentPaymentEvent struct {
	Amount uint64
	Payee  types.AccountAddress
}

// ReceivedPaymentEvent is a p2p received payment event
type ReceivedPaymentEvent struct {
	Amount uint64
	Payer  types.AccountAddress
}

// UnmarshalBinary unmarshals raw bytes into this struct.
func (ev *SentPaymentEvent) UnmarshalBinary(data []byte) error {
	return lcs.Unmarshal(data, ev)
}

// UnmarshalBinary unmarshals raw bytes into this struct.
func (ev *ReceivedPaymentEvent) UnmarshalBinary(data []byte) error {
	return lcs.Unmarshal(data, ev)
}
