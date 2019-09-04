package types

import (
	"github.com/the729/lcs"
)

// AccountResource is the Libra coin resource of an account.
type AccountResource struct {
	AuthenticationKey             []byte
	Balance                       uint64
	DelegatedWithdrawalCapability bool
	ReceivedEvents                *EventHandle
	SentEvents                    *EventHandle
	SequenceNumber                uint64
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
	return lcs.Unmarshal(data, r)
}

// Clone deep clones this struct.
func (r *AccountResource) Clone() *AccountResource {
	out := &AccountResource{}
	out.AuthenticationKey = cloneBytes(r.AuthenticationKey)
	out.Balance = r.Balance
	out.DelegatedWithdrawalCapability = r.DelegatedWithdrawalCapability
	out.ReceivedEvents = r.ReceivedEvents.Clone()
	out.SentEvents = r.SentEvents.Clone()
	out.SequenceNumber = r.SequenceNumber
	return out
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

// GetSentEvents returns sent events handle.
func (pr *ProvenAccountResource) GetSentEvents() *EventHandle {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	c := &EventHandle{}
	*c = *pr.accountResource.SentEvents
	return c
}

// GetReceivedEvents returns received events handle.
func (pr *ProvenAccountResource) GetReceivedEvents() *EventHandle {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	c := &EventHandle{}
	*c = *pr.accountResource.ReceivedEvents
	return c
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
