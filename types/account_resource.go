package types

// AccountResource is the Libra coin resource of an account.
type AccountResource struct {
	AuthenticationKey              []byte
	DelegatedKeyRotationCapability bool
	DelegatedWithdrawalCapability  bool
	ReceivedEvents                 *EventHandle
	SentEvents                     *EventHandle
	SequenceNumber                 uint64
	EventGenerator                 uint64
}

// ProvenAccountResource is the Libra coin resource of an account which is proven
// to be included in the ledger.
type ProvenAccountResource struct {
	proven          bool
	accountResource *AccountResource
	addr            AccountAddress
	ledgerInfo      *ProvenLedgerInfo
}

type BalanceResource struct {
	Coin uint64
}

// Clone deep clones this struct.
func (r *AccountResource) Clone() *AccountResource {
	out := &AccountResource{}
	out.AuthenticationKey = cloneBytes(r.AuthenticationKey)
	out.DelegatedWithdrawalCapability = r.DelegatedWithdrawalCapability
	out.ReceivedEvents = r.ReceivedEvents.Clone()
	out.SentEvents = r.SentEvents.Clone()
	out.SequenceNumber = r.SequenceNumber
	out.EventGenerator = r.EventGenerator
	return out
}

// GetLedgerInfo returns the ledger info.
func (pr *ProvenAccountResource) GetLedgerInfo() *ProvenLedgerInfo {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return pr.ledgerInfo
}

// GetSequenceNumber returns sequence number of the account.
func (pr *ProvenAccountResource) GetSequenceNumber() uint64 {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return pr.accountResource.SequenceNumber
}

// GetEventGenerator returns event generator of the account.
func (pr *ProvenAccountResource) GetEventGenerator() uint64 {
	if !pr.proven {
		panic("not valid proven account resource")
	}
	return pr.accountResource.EventGenerator
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
	return pr.addr
}
