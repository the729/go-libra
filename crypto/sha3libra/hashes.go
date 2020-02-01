package sha3libra

import (
	"hash"

	"golang.org/x/crypto/sha3"
)

const (
	libraHashSuffix = "@@$$LIBRA$$@@"

	HashSize = 32
)

type HashValue = []byte

type state struct {
	hash.Hash
	salt []byte
}

func (s *state) Reset() {
	s.Hash.Reset()
	s.Write(s.salt)
}

func newHasher(salt []byte) hash.Hash {
	saltHasher := sha3.New256()
	saltHasher.Write(salt)
	saltHasher.Write([]byte(libraHashSuffix))
	saltHash := saltHasher.Sum([]byte{})

	hasher := &state{
		Hash: sha3.New256(),
		salt: saltHash,
	}
	hasher.Reset()
	return hasher
}

func NewStructTag() hash.Hash {
	return newHasher([]byte("StructTag::libra_types::language_storage"))
}
func NewAccountAddress() hash.Hash {
	return newHasher([]byte("AccountAddress::libra_types::account_address"))
}
func NewLedgerInfo() hash.Hash { return newHasher([]byte("LedgerInfo::libra_types::ledger_info")) }
func NewWaypointLedgerInfo() hash.Hash {
	return newHasher([]byte("Ledger2WaypointConverter::libra_types::waypoint"))
}
func NewTransactionAccumulator() hash.Hash { return newHasher([]byte("TransactionAccumulator")) }
func NewEventAccumulator() hash.Hash       { return newHasher([]byte("EventAccumulator")) }
func NewSparseMerkleInternal() hash.Hash   { return newHasher([]byte("SparseMerkleInternal")) }
func NewSparseMerkleLeaf() hash.Hash {
	return newHasher([]byte("SparseMerkleLeafNode::libra_types::proof"))
}
func NewAccountStateBlob() hash.Hash {
	return newHasher([]byte("AccountStateBlob::libra_types::account_state_blob"))
}
func NewTransactionInfo() hash.Hash {
	return newHasher([]byte("TransactionInfo::libra_types::transaction"))
}
func NewTransaction() hash.Hash { return newHasher([]byte("Transaction::libra_types::transaction")) }
func NewRawTransaction() hash.Hash {
	return newHasher([]byte("RawTransaction::libra_types::transaction"))
}
func NewSignedTransaction() hash.Hash {
	return newHasher([]byte("SignedTransaction::libra_types::transaction"))
}
func NewBlock() hash.Hash            { return newHasher([]byte("BlockId")) }
func NewPacemakerTimeout() hash.Hash { return newHasher([]byte("PacemakerTimeout")) }
func NewTimeoutMsg() hash.Hash       { return newHasher([]byte("TimeoutMsg")) }
func NewVoteMsg() hash.Hash          { return newHasher([]byte("VoteMsg")) }
func NewContractEvent() hash.Hash {
	return newHasher([]byte("ContractEvent::libra_types::contract_event"))
}
func NewDiscoveryMsg() hash.Hash { return newHasher([]byte("DiscoveryMsg")) }
