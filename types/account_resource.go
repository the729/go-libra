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

// BalanceResource is LBR balance resource
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
