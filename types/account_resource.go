package types

import (
	"github.com/the729/go-libra/common/canonical_serialization"
)

type AccountResource struct {
	Balance             uint64
	SequenceNumber      uint64
	AuthenticationKey   []byte
	SentEventsCount     uint64
	ReceivedEventsCount uint64
}

func (r *AccountResource) UnmarshalBinary(data []byte) error {
	akey, err := serialization.SimpleDeserializer.ByteSlice(data)
	if err != nil {
		return err
	}
	r.AuthenticationKey = akey
	data = data[len(akey)+4:]
	r.Balance = serialization.SimpleDeserializer.Uint64(data)
	data = data[8:]
	r.ReceivedEventsCount = serialization.SimpleDeserializer.Uint64(data)
	data = data[8:]
	r.SentEventsCount = serialization.SimpleDeserializer.Uint64(data)
	data = data[8:]
	r.SequenceNumber = serialization.SimpleDeserializer.Uint64(data)
	return nil
}
