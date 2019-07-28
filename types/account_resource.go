package types

import serialization "github.com/the729/go-libra/common/canonical_serialization"

type AccountResource struct {
	Balance             uint64
	SequenceNumber      uint64
	AuthenticationKey   []byte
	SentEventsCount     uint64
	ReceivedEventsCount uint64
}

type ProvenAccountResource struct {
	proven          bool
	accountResource AccountResource
	addr            AccountAddress
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

func (pr *ProvenAccountResource) GetBalance() uint64 {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return pr.accountResource.Balance
}

func (pr *ProvenAccountResource) GetSequenceNumber() uint64 {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return pr.accountResource.SequenceNumber
}

func (pr *ProvenAccountResource) GetAuthenticationKey() []byte {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return cloneBytes(pr.accountResource.AuthenticationKey)
}

func (pr *ProvenAccountResource) GetSentEventsCount() uint64 {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return pr.accountResource.SentEventsCount
}

func (pr *ProvenAccountResource) GetReceivedEventsCount() uint64 {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return pr.accountResource.ReceivedEventsCount
}

func (pr *ProvenAccountResource) GetAddress() AccountAddress {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return AccountAddress(cloneBytes(pr.addr))
}
