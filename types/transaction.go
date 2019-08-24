package types

import (
	"crypto"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/ed25519"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
)

// RawTransaction is a raw transaction struct.
type RawTransaction = pbtypes.RawTransaction

// SignedTransaction is a signed transaction, which consists of a serialized raw transaction
// and the signature and public key.
type SignedTransaction struct {
	// RawTxnBytes is the serialized raw transaction.
	RawTxnBytes []byte

	// SenderPublicKey is the public key used to sign this transaction.
	SenderPublicKey []byte

	// SenderSignature is the signature.
	SenderSignature []byte
}

// SubmittedTransaction is a signed transaction with execution outputs.
// It is not guaranteed to be included in the ledger.
type SubmittedTransaction struct {
	*SignedTransaction

	// Info is the transaction info.
	Info *TransactionInfo

	// Events is a list of output events.
	Events EventList

	// Version is height of this transaction in the ledger.
	Version uint64
}

// ProvenTransaction is a transaction which has been proven to be included in the ledger.
type ProvenTransaction struct {
	proven     bool
	withEvents bool
	signedTxn  *SignedTransaction
	events     EventList
	version    uint64
	gasUsed    uint64
}

// FromProto parses a protobuf struct into this struct.
func (t *SignedTransaction) FromProto(pb *pbtypes.SignedTransaction) error {
	if pb == nil {
		return ErrNilInput
	}
	t.RawTxnBytes = pb.GetRawTxnBytes()
	t.SenderPublicKey = pb.GetSenderPublicKey()
	t.SenderSignature = pb.GetSenderSignature()
	return nil
}

// ToProto builds a protobuf struct from this struct.
func (t *SignedTransaction) ToProto() (*pbtypes.SignedTransaction, error) {
	return &pbtypes.SignedTransaction{
		RawTxnBytes:     t.RawTxnBytes,
		SenderPublicKey: t.SenderPublicKey,
		SenderSignature: t.SenderSignature,
	}, nil
}

// SerializeTo serializes this struct into a io.Writer.
func (t *SignedTransaction) SerializeTo(w io.Writer) error {
	if err := serialization.SimpleSerializer.Write(w, t.RawTxnBytes); err != nil {
		return err
	}
	if err := serialization.SimpleSerializer.Write(w, t.SenderPublicKey); err != nil {
		return err
	}
	if err := serialization.SimpleSerializer.Write(w, t.SenderSignature); err != nil {
		return err
	}
	return nil
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (t *SignedTransaction) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewSignedTransaction()
	if err := t.SerializeTo(hasher); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
}

// Clone deep clones this struct.
func (t *SignedTransaction) Clone() *SignedTransaction {
	out := &SignedTransaction{}
	out.RawTxnBytes = cloneBytes(t.RawTxnBytes)
	out.SenderPublicKey = cloneBytes(t.SenderPublicKey)
	out.SenderSignature = cloneBytes(t.SenderSignature)
	return out
}

// UnmarshalRawTransaction deserialize the raw transaction bytes.
func (t *SignedTransaction) UnmarshalRawTransaction() (*RawTransaction, error) {
	rt := &RawTransaction{}
	err := proto.Unmarshal(t.RawTxnBytes, rt)
	return rt, err
}

// VerifySignature verifies the signature of this transaction.
// Correct signature does NOT prove inclusion in the ledger.
func (t *SignedTransaction) VerifySignature() error {
	// 1. decode raw transaction, compare sender account and sender public key
	// Account address sometimes is different from hash(publickey), e.g.
	// libracore account 0x0.

	// rawTxn := &pbtypes.RawTransaction{}
	// if err := proto.Unmarshal(t.RawTxnBytes, rawTxn); err != nil {
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
	txnHasher.Write(t.RawTxnBytes)
	txnHash := txnHasher.Sum([]byte{})

	k := ed25519.PublicKey(t.SenderPublicKey)
	if !ed25519.Verify(k, txnHash, t.SenderSignature) {
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
	if err := st.SignedTransaction.VerifySignature(); err != nil {
		return nil, fmt.Errorf("txn(%d) signature verification fail: %v", st.Version, err)
	}

	// verify SignedTransaction and Events hash from transaction info
	txnHash := st.SignedTransaction.Hash()
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

	return &ProvenTransaction{
		// this verification alone does not prove ledger inclusion
		proven:     false,
		withEvents: withEvents,
		signedTxn:  st.SignedTransaction.Clone(),
		events:     st.Events.Clone(),
		version:    st.Version,
		gasUsed:    st.Info.GasUsed,
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

// SignRawTransaction signes a raw transaction with a private key.
func SignRawTransaction(rawTxnBytes []byte, signer ed25519.PrivateKey) *SignedTransaction {
	hasher := sha3libra.NewRawTransaction()
	hasher.Write(rawTxnBytes)
	txnHash := hasher.Sum([]byte{})
	senderPubKey := signer.Public().(ed25519.PublicKey)
	sig, _ := signer.Sign(rand.Reader, txnHash, crypto.Hash(0))

	return &SignedTransaction{
		RawTxnBytes:     rawTxnBytes,
		SenderPublicKey: senderPubKey,
		SenderSignature: sig,
	}
}
