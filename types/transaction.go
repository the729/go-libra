package types

import (
	"crypto"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/ed25519"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/language/stdscript"
)

type RawTransaction = pbtypes.RawTransaction

type SignedTransaction struct {
	RawTxnBytes     []byte
	SenderPublicKey []byte
	SenderSignature []byte
}

type SubmittedTransaction struct {
	*SignedTransaction
	Info    *TransactionInfo
	Events  EventList
	Version uint64
}

type ProvenTransaction struct {
	proven     bool
	withEvents bool
	signedTxn  *SignedTransaction
	events     EventList
	version    uint64
	gasUsed    uint64
}

func (t *SignedTransaction) FromProto(pb *pbtypes.SignedTransaction) error {
	if pb == nil {
		return ErrNilInput
	}
	t.RawTxnBytes = pb.GetRawTxnBytes()
	t.SenderPublicKey = pb.GetSenderPublicKey()
	t.SenderSignature = pb.GetSenderSignature()
	return nil
}

func (t *SignedTransaction) ToProto() (*pbtypes.SignedTransaction, error) {
	return &pbtypes.SignedTransaction{
		RawTxnBytes:     t.RawTxnBytes,
		SenderPublicKey: t.SenderPublicKey,
		SenderSignature: t.SenderSignature,
	}, nil
}

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

func (t *SignedTransaction) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewSignedTransaction()
	if err := t.SerializeTo(hasher); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
}

func (t *SignedTransaction) Clone() *SignedTransaction {
	out := &SignedTransaction{}
	out.RawTxnBytes = cloneBytes(t.RawTxnBytes)
	out.SenderPublicKey = cloneBytes(t.SenderPublicKey)
	out.SenderSignature = cloneBytes(t.SenderSignature)
	return out
}

func (t *SignedTransaction) UnmarshalRawTransaction() (*RawTransaction, error) {
	rt := &RawTransaction{}
	err := proto.Unmarshal(t.RawTxnBytes, rt)
	return rt, err
}

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
		proven:     true,
		withEvents: withEvents,
		signedTxn:  st.SignedTransaction.Clone(),
		events:     st.Events.Clone(),
		version:    st.Version,
		gasUsed:    st.Info.GasUsed,
	}, nil
}

func (pt *ProvenTransaction) GetSignedTxn() *SignedTransaction {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	return pt.signedTxn.Clone()
}

func (pt *ProvenTransaction) GetWithEvents() bool {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	return pt.withEvents
}

func (pt *ProvenTransaction) GetEvents() []*ContractEvent {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	return pt.events.Clone()
}

func (pt *ProvenTransaction) GetVersion() uint64 {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	return pt.version
}

func (pt *ProvenTransaction) GetGasUsed() uint64 {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	return pt.gasUsed
}

func NewRawP2PTransaction(
	senderAddress, receiverAddress AccountAddress,
	senderSequenceNumber uint64,
	amount, maxGasAmount, gasUnitPrice uint64,
	expiration time.Time,
) ([]byte, error) {
	ammountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(ammountBytes, amount)

	txn := &RawTransaction{
		SenderAccount:  senderAddress,
		SequenceNumber: senderSequenceNumber,
		Payload: &pbtypes.RawTransaction_Program{
			Program: &pbtypes.Program{
				Code: stdscript.PeerToPeerTransfer,
				Arguments: []*pbtypes.TransactionArgument{
					{
						Type: pbtypes.TransactionArgument_ADDRESS,
						Data: receiverAddress,
					},
					{
						Type: pbtypes.TransactionArgument_U64,
						Data: ammountBytes,
					},
				},
				Modules: nil,
			},
		},
		MaxGasAmount:   maxGasAmount,
		GasUnitPrice:   gasUnitPrice,
		ExpirationTime: uint64(expiration.Unix()),
	}

	// j, _ := json.MarshalIndent(txn, "", "    ")
	// log.Printf("Raw txn: %s", string(j))

	raw, err := proto.Marshal(txn)
	return raw, err
}

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
