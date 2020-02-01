package types

type LedgerInfoVerifier interface {
	Verify(*LedgerInfoWithSignatures) error
	EpochChangeVerificationRequired(epoch uint64) bool
}
