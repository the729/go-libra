package types

import (
	"fmt"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/lcs"
)

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

// SubmittedTransaction is an abstract transaction (user txn, writeset, or block metadata)
// with execution outputs. It is not guaranteed to be included in the ledger.
type SubmittedTransaction struct {
	// RawTxn is raw (bytes) abstract transaction (user txn, writeset, or block metadata).
	RawTxn []byte

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
	txn         isTransaction
	txnHash     HashValue
	events      EventList
	version     uint64
	gasUsed     uint64
	majorStatus VMStatusCode
	ledgerInfo  *ProvenLedgerInfo
}

// Verify the submitted transaction, and output a ProvenTransaction which is NOT fully proven yet.
//
// Your should not need to call this function.
// To fully prove a submitted transaction, you will need to verify a SignedTransactionWithProof or
// a TransactionListWithProof.
func (st *SubmittedTransaction) Verify() (*ProvenTransaction, error) {
	// according to https://community.libra.org/t/how-to-verify-a-signedtransaction-thoroughly/1214/3,
	// it is unnecessary to verify SignedTransaction itself

	// verify Events hash from transaction info
	eventHash := st.Events.Hash()
	withEvents := true
	if !sha3libra.Equal(eventHash, st.Info.EventRootHash) {
		if st.Events != nil {
			return nil, fmt.Errorf("event root hash mismatch in txn(%d)", st.Version)
		}
		// if event hash does not match, and events is nil, must be without events
		withEvents = false
	}

	// verify Transaction hash from transaction info
	hasher := sha3libra.NewTransaction()
	if _, err := hasher.Write(st.RawTxn); err != nil {
		panic(err)
	}
	txnHash := hasher.Sum([]byte{})
	if !sha3libra.Equal(txnHash, st.Info.TransactionHash) {
		return nil, fmt.Errorf("txn hash mismatch in txn(%d)", st.Version)
	}

	decodedTxn := &Transaction{}
	if err := lcs.Unmarshal(st.RawTxn, decodedTxn); err != nil {
		return nil, fmt.Errorf("lcs unmarshal signedtxn error: %v", err)
	}

	return &ProvenTransaction{
		// this verification alone does not prove ledger inclusion
		proven:      false,
		withEvents:  withEvents,
		txn:         decodedTxn.Transaction,
		txnHash:     st.Info.Hash(),
		events:      st.Events.Clone(),
		version:     st.Version,
		gasUsed:     st.Info.GasUsed,
		majorStatus: st.Info.MajorStatus,
	}, nil
}

// GetLedgerInfo returns the ledger info.
func (pt *ProvenTransaction) GetLedgerInfo() *ProvenLedgerInfo {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	return pt.ledgerInfo
}

// GetSignedTxn returns a copy of the underlying signed user transaction.
// It returns nil if the transaction is not a user transaction.
func (pt *ProvenTransaction) GetSignedTxn() *SignedTransaction {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	signedUserTxn, ok := pt.txn.(*SignedTransaction)
	if !ok {
		// This transaction is not a signed user transaction
		return nil
	}
	return signedUserTxn.Clone()
}

// GetBlockMetadata returns a copy of the underlying block metadata.
// It returns nil if the transaction is not a block metadata.
func (pt *ProvenTransaction) GetBlockMetadata() *BlockMetaData {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	blockMetadata, ok := pt.txn.(*BlockMetaData)
	if !ok {
		// This transaction is not a signed user transaction
		return nil
	}
	return blockMetadata.Clone()
}

// GetHash returns a copy of the transaction info hash
func (pt *ProvenTransaction) GetHash() HashValue {
	if !pt.proven {
		panic("not valid proven transaction")
	}
	return cloneBytes(pt.txnHash)
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
