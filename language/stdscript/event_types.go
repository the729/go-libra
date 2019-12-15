package stdscript

import (
	"github.com/the729/go-libra/types"
	"github.com/the729/lcs"
)

// PaymentEvent is a standard p2p sent or received payment event
type PaymentEvent struct {
	Amount  uint64
	Address types.AccountAddress
}

// UnmarshalBinary unmarshals raw bytes into this struct.
func (ev *PaymentEvent) UnmarshalBinary(data []byte) error {
	return lcs.Unmarshal(data, ev)
}
