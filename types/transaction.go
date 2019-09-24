package types

import (
	"crypto"
	"crypto/rand"
	"errors"
	"fmt"

	"golang.org/x/crypto/ed25519"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/lcs"
)

// SignedTransaction is a signed transaction, which consists of a raw transaction
// and the signature and public key.
type SignedTransaction struct {
	// RawTxn is the raw transaction.
	RawTxn *RawTransaction

	// PublicKey is the public key of the sender.
	PublicKey []byte

	// Signature is the signature.
	Signature []byte
}

// SubmittedTransaction is a signed transaction with execution outputs.
// It is not guaranteed to be included in the ledger.
type SubmittedTransaction struct {
	RawSignedTxn []byte

	// Info is the transaction info.
	Info *TransactionInfo

	// Events is a list of output events.
	Events EventList

	// Version is height of this transaction in the ledger.
	Version uint64
}

// ProvenTransaction is a transaction which has been proven to be included in the ledger.
type ProvenTransaction struct {
	proven      bool
	withEvents  bool
	signedTxn   *SignedTransaction
	events      EventList
	version     uint64
	gasUsed     uint64
	majorStatus VMStatusCode
}

// // FromProto parses a protobuf struct into this struct.
// func (t *SignedTransaction) FromProto(pb *pbtypes.SignedTransaction) error {
// 	if pb == nil {
// 		return ErrNilInput
// 	}
// 	return lcs.Unmarshal(pb.SignedTxn, t)
// }

// ToProto builds a protobuf struct from this struct.
func (t *SignedTransaction) ToProto() (*pbtypes.SignedTransaction, error) {
	b, err := lcs.Marshal(t)
	if err != nil {
		return nil, err
	}
	return &pbtypes.SignedTransaction{SignedTxn: b}, nil
}

// // Hash ouptuts the hash of this struct, using the appropriate hash function.
// func (t *SignedTransaction) Hash() sha3libra.HashValue {
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

// Verify the submitted transaction, and output a ProvenTransaction which is NOT fully proven yet.
//
// Your should not need to call this function.
// To fully prove a submitted transaction, you will need to verify a SignedTransactionWithProof or
// a TransactionListWithProof.
func (st *SubmittedTransaction) Verify() (*ProvenTransaction, error) {
	// according to https://community.libra.org/t/how-to-verify-a-signedtransaction-thoroughly/1214/3,
	// it is unnecessary to verify SignedTransaction itself

	// verify SignedTransaction and Events hash from transaction info
	hasher := sha3libra.NewSignedTransaction()
	if _, err := hasher.Write(st.RawSignedTxn); err != nil {
		panic(err)
	}
	txnHash := hasher.Sum([]byte{})
	if !sha3libra.Equal(txnHash, st.Info.SignedTransactionHash) {
		return nil, fmt.Errorf("signed txn hash mismatch in txn(%d)", st.Version)
	}
	eventHash := st.Events.Hash()
	withEvents := true
	if !sha3libra.Equal(eventHash, st.Info.EventRootHash) {
		if st.Events != nil {
			return nil, fmt.Errorf("event root hash mismatch in txn(%d)", st.Version)
		}
		// if event hash does not match, and events is nil, must be without events
		withEvents = false
	}

	decodedTxn := &SignedTransaction{}
	if err := lcs.Unmarshal(st.RawSignedTxn, decodedTxn); err != nil {
		return nil, fmt.Errorf("lcs unmarshal signedtxn error: %v", err)
	}
	return &ProvenTransaction{
		// this verification alone does not prove ledger inclusion
		proven:      false,
		withEvents:  withEvents,
		signedTxn:   decodedTxn,
		events:      st.Events.Clone(),
		version:     st.Version,
		gasUsed:     st.Info.GasUsed,
		majorStatus: st.Info.MajorStatus,
	}, nil
}

// GetSignedTxn returns a copy of the underlying signed transaction.
func (pt *ProvenTransaction) GetSignedTxn() *SignedTransaction {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	return pt.signedTxn.Clone()
}

// GetWithEvents returns whether this proven transaction has output events included.
func (pt *ProvenTransaction) GetWithEvents() bool {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	return pt.withEvents
}

// GetEvents returns a copy of the underlying events list.
//
// Nil output does not necessarily mean empty output event list. It could be this proven
// transaction does not have output events list included. Call GetWithEvents() to find out.
func (pt *ProvenTransaction) GetEvents() []*ContractEvent {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	return pt.events.Clone()
}

// GetVersion returns the height of this transaction.
func (pt *ProvenTransaction) GetVersion() uint64 {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	return pt.version
}

// GetGasUsed returns the gas used to process this transaction.
func (pt *ProvenTransaction) GetGasUsed() uint64 {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	return pt.gasUsed
}

// GetMajorStatus returns the major VM status returned from this transaction.
func (pt *ProvenTransaction) GetMajorStatus() VMStatusCode {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	return pt.majorStatus
}

// SignRawTransaction signes a raw transaction with a private key.
func SignRawTransaction(rawTxn *RawTransaction, signer ed25519.PrivateKey) *SignedTransaction {
	hasher := sha3libra.NewRawTransaction()
	if err := lcs.NewEncoder(hasher).Encode(rawTxn); err != nil {
		panic(err)
	}
	txnHash := hasher.Sum([]byte{})
	senderPubKey := signer.Public().(ed25519.PublicKey)
	sig, _ := signer.Sign(rand.Reader, txnHash, crypto.Hash(0))

	return &SignedTransaction{
		RawTxn:    rawTxn,
		PublicKey: senderPubKey,
		Signature: sig,
	}
}
