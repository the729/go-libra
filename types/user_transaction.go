package types

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ed25519"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/lcs"
)

// SignedTransaction is a signed user transaction, which consists of a raw transaction
// and the signature and public key.
//
// TODO: according to a comment in libra codebase(types/src/transaction/mod.rs Line#1065),
// should be renamed to SignedUserTransaction.
type SignedTransaction struct {
	// RawTxn is the raw transaction.
	RawTxn *RawTransaction

	// PublicKey is the public key of the sender.
	PublicKey []byte

	// Signature is the signature.
	Signature []byte
}

// ToProto builds a protobuf struct from this struct.
func (t *SignedTransaction) ToProto() (*pbtypes.SignedTransaction, error) {
	b, err := lcs.Marshal(t)
	if err != nil {
		return nil, err
	}
	return &pbtypes.SignedTransaction{TxnBytes: b}, nil
}

// // Hash ouptuts the hash of this struct, using the appropriate hash function.
// func (t *SignedTransaction) Hash() HashValue {
// 	hasher := sha3libra.NewSignedTransaction()
// 	if err := lcs.NewEncoder(hasher).Encode(t); err != nil {
// 		panic(err)
// 	}
// 	return hasher.Sum([]byte{})
// }

// Clone deep clones this struct.
func (t *SignedTransaction) Clone() *SignedTransaction {
	out := &SignedTransaction{}
	out.RawTxn = t.RawTxn.Clone()
	out.PublicKey = cloneBytes(t.PublicKey)
	out.Signature = cloneBytes(t.Signature)
	return out
}

// VerifySignature verifies the signature of this transaction.
// Correct signature does NOT prove inclusion in the ledger.
func (t *SignedTransaction) VerifySignature() error {
	// 1. decode raw transaction, compare sender account and sender public key
	// Account address sometimes is different from hash(publickey), e.g.
	// libracore account 0x0.

	// rawTxn, err := t.UnmarshalRawTransaction()
	// if err != nil {
	// 	return errors.New("invalid raw transaction")
	// }
	// addrHasher := sha3.New256()
	// addrHasher.Write(t.SenderPublicKey)
	// gotAddr := addrHasher.Sum([]byte{})
	// if !sha3libra.Equal(gotAddr, rawTxn.SenderAccount) {
	// 	return errors.New("transaction sender does not match signer")
	// }

	// 2. verify signature
	txnHasher := sha3libra.NewRawTransaction()
	if err := lcs.NewEncoder(txnHasher).Encode(t.RawTxn); err != nil {
		return fmt.Errorf("marshal raw txn error: %v", err)
	}
	txnHash := txnHasher.Sum([]byte{})

	k := ed25519.PublicKey(t.PublicKey)
	if !ed25519.Verify(k, txnHash, t.Signature) {
		return errors.New("signature verification fail")
	}

	return nil
}
