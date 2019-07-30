package types

import (
	"crypto"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/ed25519"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/language/stdscript"
)

type SignedTransaction struct {
	RawTxnBytes     []byte
	SenderPublicKey []byte
	SenderSignature []byte
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

func (t *SignedTransaction) Verify() error {
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

func NewRawP2PTransaction(
	senderAddress, receiverAddress AccountAddress,
	senderSequenceNumber uint64,
	amount, maxGasAmount, gasUnitPrice uint64,
	expiration time.Time,
) ([]byte, error) {
	ammountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(ammountBytes, amount)

	txn := &pbtypes.RawTransaction{
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
