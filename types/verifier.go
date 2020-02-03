package types

// LedgerInfoVerifier can verify a LedgerInfoWithSignatures.
type LedgerInfoVerifier interface {
	// Verify the given LedgerInfoWithSignatures
	Verify(*LedgerInfoWithSignatures) error

	// Returns whether this verifier maybe outdated at a given epoch
	EpochChangeVerificationRequired(epoch uint64) bool
}
