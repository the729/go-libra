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

func newHasher(salt []byte) hash.Hash {
	saltHasher := sha3.New256()
	saltHasher.Write(salt)
	saltHasher.Write([]byte(libraHashSuffix))
	saltHash := saltHasher.Sum([]byte{})

	hasher := sha3.New256()
	hasher.Write(saltHash)
	return hasher
}

func NewAccessPath() hash.Hash             { return newHasher([]byte("VM_ACCESS_PATH")) }
func NewAccountAddress() hash.Hash         { return newHasher([]byte("AccountAddress")) }
func NewLedgerInfo() hash.Hash             { return newHasher([]byte("LedgerInfo")) }
func NewTransactionAccumulator() hash.Hash { return newHasher([]byte("TransactionAccumulator")) }
func NewEventAccumulator() hash.Hash       { return newHasher([]byte("EventAccumulator")) }
func NewSparseMerkleInternal() hash.Hash   { return newHasher([]byte("SparseMerkleInternal")) }
func NewSparseMerkleLeaf() hash.Hash       { return newHasher([]byte("SparseMerkleLeaf")) }
func NewAccountStateBlob() hash.Hash       { return newHasher([]byte("AccountStateBlob")) }
func NewTransactionInfo() hash.Hash        { return newHasher([]byte("TransactionInfo")) }
func NewRawTransaction() hash.Hash         { return newHasher([]byte("RawTransaction")) }
func NewSignedTransaction() hash.Hash      { return newHasher([]byte("SignedTransaction")) }
func NewBlock() hash.Hash                  { return newHasher([]byte("BlockId")) }
func NewPacemakerTimeout() hash.Hash       { return newHasher([]byte("PacemakerTimeout")) }
func NewTimeoutMsg() hash.Hash             { return newHasher([]byte("TimeoutMsg")) }
func NewVoteMsg() hash.Hash                { return newHasher([]byte("VoteMsg")) }
func NewContractEvent() hash.Hash          { return newHasher([]byte("ContractEvent")) }
func NewDiscoveryMsg() hash.Hash           { return newHasher([]byte("DiscoveryMsg")) }
