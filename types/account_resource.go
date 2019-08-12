package types

import serialization "github.com/the729/go-libra/common/canonical_serialization"

// AccountResource is the Libra coin resource of an account.
type AccountResource struct {
	Balance                       uint64
	SequenceNumber                uint64
	AuthenticationKey             []byte
	SentEventsCount               uint64
	ReceivedEventsCount           uint64
	DelegatedWithdrawalCapability bool
}

// ProvenAccountResource is the Libra coin resource of an account which is proven
// to be included in the ledger.
type ProvenAccountResource struct {
	proven          bool
	accountResource AccountResource
	addr            AccountAddress
}

// UnmarshalBinary unmarshals raw bytes into this account resource struct.
func (r *AccountResource) UnmarshalBinary(data []byte) error {
	akey, err := serialization.SimpleDeserializer.ByteSlice(data)
	if err != nil {
		return err
	}
	r.AuthenticationKey = akey
	data = data[len(akey)+4:]
	r.Balance = serialization.SimpleDeserializer.Uint64(data)
	data = data[8:]
	r.DelegatedWithdrawalCapability = serialization.SimpleDeserializer.Bool(data)
	data = data[1:]
	r.ReceivedEventsCount = serialization.SimpleDeserializer.Uint64(data)
	data = data[8:]
	r.SentEventsCount = serialization.SimpleDeserializer.Uint64(data)
	data = data[8:]
	r.SequenceNumber = serialization.SimpleDeserializer.Uint64(data)
	return nil
}

// GetBalance returns Libra coin balance in microLibra.
func (pr *ProvenAccountResource) GetBalance() uint64 {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return pr.accountResource.Balance
}

// GetSequenceNumber returns sequence number of the account.
func (pr *ProvenAccountResource) GetSequenceNumber() uint64 {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return pr.accountResource.SequenceNumber
}

// GetAuthenticationKey returns a copy of the hash of public key current in use.
func (pr *ProvenAccountResource) GetAuthenticationKey() []byte {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return cloneBytes(pr.accountResource.AuthenticationKey)
}

// GetSentEventsCount returns count of sent events.
func (pr *ProvenAccountResource) GetSentEventsCount() uint64 {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return pr.accountResource.SentEventsCount
}

// GetReceivedEventsCount returns count of received events.
func (pr *ProvenAccountResource) GetReceivedEventsCount() uint64 {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return pr.accountResource.ReceivedEventsCount
}

// GetDelegatedWithdrawalCapability returns delegated withdrawal capability.
func (pr *ProvenAccountResource) GetDelegatedWithdrawalCapability() bool {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return pr.accountResource.DelegatedWithdrawalCapability
}

// GetAddress returns a copy of the address to which this resource belongs.
func (pr *ProvenAccountResource) GetAddress() AccountAddress {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return AccountAddress(cloneBytes(pr.addr))
}
