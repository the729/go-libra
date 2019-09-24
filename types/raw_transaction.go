package types

import (
	"github.com/the729/lcs"
)

// RawTransaction is a raw transaction struct.
type RawTransaction struct {
	// Sender address.
	Sender AccountAddress

	// SequenceNumber of this transaction corresponding to sender's account.
	SequenceNumber uint64

	// Payload is the transaction script to execute.
	Payload TransactionPayload `lcs:"enum:TransactionPayload"`

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

type TransactionArgument interface {
	isTransactionArgument()
	Clone() TransactionArgument
}
type TxnArgU64 uint64
type TxnArgAddress AccountAddress
type TxnArgString string
type TxnArgBytes []byte

func (TxnArgU64) isTransactionArgument()           {}
func (TxnArgAddress) isTransactionArgument()       {}
func (TxnArgString) isTransactionArgument()        {}
func (TxnArgBytes) isTransactionArgument()         {}
func (v TxnArgU64) Clone() TransactionArgument     { return v }
func (v TxnArgAddress) Clone() TransactionArgument { return TxnArgAddress(cloneBytes(v)) }
func (v TxnArgString) Clone() TransactionArgument  { return v }
func (v TxnArgBytes) Clone() TransactionArgument   { return TxnArgBytes(cloneBytes(v)) }

var txnArgEnumDef = []lcs.EnumVariant{
	{
		Name:     "TransactionArgument",
		Value:    0,
		Template: TxnArgU64(0),
	},
	{
		Name:     "TransactionArgument",
		Value:    1,
		Template: TxnArgAddress(nil),
	},
	{
		Name:     "TransactionArgument",
		Value:    2,
		Template: TxnArgString(""),
	},
	{
		Name:     "TransactionArgument",
		Value:    3,
		Template: TxnArgBytes(nil),
	},
}

type WriteOpWithPath struct {
	AccessPath *AccessPath
	WriteOp    WriteOp `lcs:"enum:WriteOp"`
}

type WriteOp interface {
	isWriteOp()
	Clone() WriteOp
}
type WriteOpValue []byte
type WriteOpDeletion struct{}

func (WriteOpValue) isWriteOp()          {}
func (WriteOpDeletion) isWriteOp()       {}
func (v WriteOpValue) Clone() WriteOp    { return WriteOpValue(cloneBytes(v)) }
func (v WriteOpDeletion) Clone() WriteOp { return v }

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
func (v *WriteOpWithPath) Clone() *WriteOpWithPath {
	return &WriteOpWithPath{AccessPath: v.AccessPath.Clone(), WriteOp: v.WriteOp.Clone()}
}

type TransactionPayload interface {
	isTransactionPayload()
	Clone() TransactionPayload
}

type TxnPayloadProgram struct {
	Code    []byte
	Args    []TransactionArgument `lcs:"enum:TransactionArgument"`
	Modules [][]byte
}

func (*TxnPayloadProgram) EnumTypes() []lcs.EnumVariant { return txnArgEnumDef }
func (v *TxnPayloadProgram) Clone() TransactionPayload {
	c := cloneBytes(v.Code)
	args := make([]TransactionArgument, 0, len(v.Args))
	for _, arg := range v.Args {
		args = append(args, arg.Clone())
	}
	mods := make([][]byte, 0, len(v.Modules))
	for _, mod := range v.Modules {
		mods = append(mods, cloneBytes(mod))
	}
	return &TxnPayloadProgram{Code: c, Args: args, Modules: mods}
}

type TxnPayloadWriteSet []*WriteOpWithPath

func (v TxnPayloadWriteSet) Clone() TransactionPayload {
	n := make([]*WriteOpWithPath, 0, len(v))
	for _, wop := range v {
		n = append(n, wop.Clone())
	}
	return TxnPayloadWriteSet(n)
}

type TxnPayloadScript struct {
	Code []byte
	Args []TransactionArgument `lcs:"enum:TransactionArgument"`
}

func (*TxnPayloadScript) EnumTypes() []lcs.EnumVariant { return txnArgEnumDef }
func (v *TxnPayloadScript) Clone() TransactionPayload {
	c := cloneBytes(v.Code)
	args := make([]TransactionArgument, 0, len(v.Args))
	for _, arg := range v.Args {
		args = append(args, arg.Clone())
	}
	return &TxnPayloadScript{Code: c, Args: args}
}

type TxnPayloadModule []byte

func (v TxnPayloadModule) Clone() TransactionPayload { return TxnPayloadModule(cloneBytes(v)) }

func (*TxnPayloadProgram) isTransactionPayload() {}
func (TxnPayloadWriteSet) isTransactionPayload() {}
func (*TxnPayloadScript) isTransactionPayload()  {}
func (TxnPayloadModule) isTransactionPayload()   {}

func (*RawTransaction) EnumTypes() []lcs.EnumVariant {
	return []lcs.EnumVariant{
		{
			Name:     "TransactionPayload",
			Value:    0,
			Template: (*TxnPayloadProgram)(nil),
		},
		{
			Name:     "TransactionPayload",
			Value:    1,
			Template: TxnPayloadWriteSet(nil),
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

func (rt *RawTransaction) Clone() *RawTransaction {
	return &RawTransaction{
		Sender:         cloneBytes(rt.Sender),
		SequenceNumber: rt.SequenceNumber,
		Payload:        rt.Payload.Clone(),
		MaxGasAmount:   rt.MaxGasAmount,
		GasUnitPrice:   rt.GasUnitPrice,
		ExpirationTime: rt.ExpirationTime,
	}
}
