package types

import (
	"crypto"
	"crypto/rand"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/lcs"
	"golang.org/x/crypto/ed25519"
)

// RawTransaction is a raw user transaction.
//
// TODO: according to a comment in libra codebase(types/src/transaction/mod.rs Line#1065),
// should be renamed to RawUserTransaction.
type RawTransaction struct {
	// Sender address.
	Sender AccountAddress

	// SequenceNumber of this transaction corresponding to sender's account.
	SequenceNumber uint64

	// Payload is the transaction script to execute.
	Payload TransactionPayload `lcs:"enum=TransactionPayload"`

	// Maximal total gas specified by wallet to spend for this transaction.
	MaxGasAmount uint64

	// Maximal price can be paid per gas.
	GasUnitPrice uint64

	// Expiration time for this transaction.  If storage is queried and
	// the time returned is greater than or equal to this time and this
	// transaction has not been included, you can be certain that it will
	// never be included.
	// A transaction that doesn't expire is represented by a very large value like
	// u64::max_value().
	ExpirationTime uint64
}

// TransactionArgument is the enum type of TransactionArgument
type TransactionArgument interface {
	isTransactionArgument()
	Clone() TransactionArgument
}

// TxnArgU64 is uint64 transaction argument
type TxnArgU64 uint64

// TxnArgAddress is transaction argument of account address type
type TxnArgAddress AccountAddress

// TxnArgBytes is byte array transaction argument
type TxnArgBytes []byte

// TxnArgBool is boolean transaction argument
type TxnArgBool bool

func (TxnArgU64) isTransactionArgument()     {}
func (TxnArgAddress) isTransactionArgument() {}
func (TxnArgBytes) isTransactionArgument()   {}
func (TxnArgBool) isTransactionArgument()    {}

// Clone the argument
func (v TxnArgU64) Clone() TransactionArgument { return v }

// Clone the argument
func (v TxnArgAddress) Clone() TransactionArgument { return v }

// Clone the argument
func (v TxnArgBytes) Clone() TransactionArgument { return TxnArgBytes(cloneBytes(v)) }

// Clone the argument
func (v TxnArgBool) Clone() TransactionArgument { return v }

var txnArgEnumDef = []lcs.EnumVariant{
	{
		Name:     "TransactionArgument",
		Value:    0,
		Template: TxnArgU64(0),
	},
	{
		Name:     "TransactionArgument",
		Value:    1,
		Template: TxnArgAddress{},
	},
	{
		Name:     "TransactionArgument",
		Value:    2,
		Template: TxnArgBytes(nil),
	},
	{
		Name:     "TransactionArgument",
		Value:    3,
		Template: TxnArgBool(false),
	},
}

// WriteOpWithPath is write op with access path
type WriteOpWithPath struct {
	AccessPath *AccessPath
	WriteOp    WriteOp `lcs:"enum=WriteOp"`
}

// WriteOp is an enum type of either value or deletion
type WriteOp interface {
	isWriteOp()
	Clone() WriteOp
}

// WriteOpValue is a variant of WriteOp
type WriteOpValue []byte

// WriteOpDeletion is a variant of WriteOp
type WriteOpDeletion struct{}

func (WriteOpValue) isWriteOp()    {}
func (WriteOpDeletion) isWriteOp() {}

// Clone the WriteOp
func (v WriteOpValue) Clone() WriteOp { return WriteOpValue(cloneBytes(v)) }

// Clone the WriteOp
func (v WriteOpDeletion) Clone() WriteOp { return v }

// EnumTypes defines enum variants for lcs
func (*WriteOpWithPath) EnumTypes() []lcs.EnumVariant {
	return []lcs.EnumVariant{
		{
			Name:     "WriteOp",
			Value:    0,
			Template: WriteOpDeletion(struct{}{}),
		},
		{
			Name:     "WriteOp",
			Value:    1,
			Template: WriteOpValue(nil),
		},
	}
}

// Clone the WriteOpWithPath
func (v *WriteOpWithPath) Clone() *WriteOpWithPath {
	return &WriteOpWithPath{AccessPath: v.AccessPath.Clone(), WriteOp: v.WriteOp.Clone()}
}

// TransactionPayload is the enum type of transaction payload
type TransactionPayload interface {
	isTransactionPayload()
	Clone() TransactionPayload
}

// TxnPayloadWriteSet is variant of TransactionPayload
type TxnPayloadWriteSet struct {
	WriteSet []*WriteOpWithPath
	Events   []*ContractEvent
}

// Clone the transaction payload
func (v *TxnPayloadWriteSet) Clone() TransactionPayload {
	ws := make([]*WriteOpWithPath, 0, len(v.WriteSet))
	for _, wop := range v.WriteSet {
		ws = append(ws, wop.Clone())
	}
	ev := make([]*ContractEvent, 0, len(v.Events))
	for _, ev1 := range v.Events {
		ev = append(ev, ev1.Clone())
	}
	return &TxnPayloadWriteSet{ws, ev}
}

// TxnPayloadScript is variant of TransactionPayload
type TxnPayloadScript struct {
	Code []byte
	Args []TransactionArgument `lcs:"enum=TransactionArgument"`
}

// EnumTypes defines enum variants for lcs
func (*TxnPayloadScript) EnumTypes() []lcs.EnumVariant { return txnArgEnumDef }

// Clone the transaction payload
func (v *TxnPayloadScript) Clone() TransactionPayload {
	c := cloneBytes(v.Code)
	args := make([]TransactionArgument, 0, len(v.Args))
	for _, arg := range v.Args {
		args = append(args, arg.Clone())
	}
	return &TxnPayloadScript{Code: c, Args: args}
}

// TxnPayloadModule is variant of TransactionPayload
type TxnPayloadModule []byte

// Clone the transaction payload
func (v TxnPayloadModule) Clone() TransactionPayload { return TxnPayloadModule(cloneBytes(v)) }

func (*TxnPayloadWriteSet) isTransactionPayload() {}
func (*TxnPayloadScript) isTransactionPayload()   {}
func (TxnPayloadModule) isTransactionPayload()    {}

// EnumTypes defines enum variants for lcs
func (*RawTransaction) EnumTypes() []lcs.EnumVariant {
	return []lcs.EnumVariant{
		{
			Name:     "TransactionPayload",
			Value:    1,
			Template: (*TxnPayloadWriteSet)(nil),
		},
		{
			Name:     "TransactionPayload",
			Value:    2,
			Template: (*TxnPayloadScript)(nil),
		},
		{
			Name:     "TransactionPayload",
			Value:    3,
			Template: TxnPayloadModule(nil),
		},
	}
}

// Clone the raw transaction
func (rt *RawTransaction) Clone() *RawTransaction {
	return &RawTransaction{
		Sender:         rt.Sender,
		SequenceNumber: rt.SequenceNumber,
		Payload:        rt.Payload.Clone(),
		MaxGasAmount:   rt.MaxGasAmount,
		GasUnitPrice:   rt.GasUnitPrice,
		ExpirationTime: rt.ExpirationTime,
	}
}

// Sign the raw transaction with a private key. It panics when the private key is
// invalid.
func (rt *RawTransaction) Sign(signer ed25519.PrivateKey) *SignedTransaction {
	hasher := sha3libra.NewRawTransaction()
	if err := lcs.NewEncoder(hasher).Encode(rt); err != nil {
		panic(err)
	}
	txnHash := hasher.Sum([]byte{})
	senderPubKey := signer.Public().(ed25519.PublicKey)
	sig, err := signer.Sign(rand.Reader, txnHash, crypto.Hash(0))
	if err != nil {
		panic(err)
	}

	return &SignedTransaction{
		RawTxn:    rt,
		PublicKey: senderPubKey,
		Signature: sig,
	}
}
