package types

import "github.com/the729/lcs"

// Transaction is an abstraction of user transaction and system transaction
// such as WriteSet and BlockMetaData
type Transaction struct {
	Transaction isTransaction `lcs:"enum=transaction"`
}

type isTransaction interface {
	isTransaction()
}

type WriteSet []*WriteOpWithPath

func (*SignedTransaction) isTransaction() {}
func (WriteSet) isTransaction()           {}
func (*BlockMetaData) isTransaction()     {}

// EnumTypes defines enum variants for lcs
func (*Transaction) EnumTypes() []lcs.EnumVariant {
	return []lcs.EnumVariant{
		{
			Name:     "transaction",
			Value:    0, // UserTransaction
			Template: (*SignedTransaction)(nil),
		},
		{
			Name:     "transaction",
			Value:    1, // WriteSet
			Template: WriteSet(nil),
		},
		{
			Name:     "transaction",
			Value:    2, // BlockMetaData
			Template: (*BlockMetaData)(nil),
		},
	}
}
