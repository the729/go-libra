// +build !js
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: transaction.proto

package pbtypes

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type TransactionArgument_ArgType int32

const (
	TransactionArgument_U64       TransactionArgument_ArgType = 0
	TransactionArgument_ADDRESS   TransactionArgument_ArgType = 1
	TransactionArgument_STRING    TransactionArgument_ArgType = 2
	TransactionArgument_BYTEARRAY TransactionArgument_ArgType = 3
)

var TransactionArgument_ArgType_name = map[int32]string{
	0: "U64",
	1: "ADDRESS",
	2: "STRING",
	3: "BYTEARRAY",
}

var TransactionArgument_ArgType_value = map[string]int32{
	"U64":       0,
	"ADDRESS":   1,
	"STRING":    2,
	"BYTEARRAY": 3,
}

func (x TransactionArgument_ArgType) String() string {
	return proto.EnumName(TransactionArgument_ArgType_name, int32(x))
}

func (TransactionArgument_ArgType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{0, 0}
}

// An argument to the transaction if the transaction takes arguments
type TransactionArgument struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TransactionArgument) Reset()         { *m = TransactionArgument{} }
func (m *TransactionArgument) String() string { return proto.CompactTextString(m) }
func (*TransactionArgument) ProtoMessage()    {}
func (*TransactionArgument) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{0}
}

func (m *TransactionArgument) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionArgument.Unmarshal(m, b)
}
func (m *TransactionArgument) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionArgument.Marshal(b, m, deterministic)
}
func (m *TransactionArgument) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionArgument.Merge(m, src)
}
func (m *TransactionArgument) XXX_Size() int {
	return xxx_messageInfo_TransactionArgument.Size(m)
}
func (m *TransactionArgument) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionArgument.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionArgument proto.InternalMessageInfo

// A generic structure that represents signed RawTransaction
type SignedTransaction struct {
	// LCS bytes representation of a SignedTransaction.
	TxnBytes             []byte   `protobuf:"bytes,5,opt,name=txn_bytes,json=txnBytes,proto3" json:"txn_bytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignedTransaction) Reset()         { *m = SignedTransaction{} }
func (m *SignedTransaction) String() string { return proto.CompactTextString(m) }
func (*SignedTransaction) ProtoMessage()    {}
func (*SignedTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{1}
}

func (m *SignedTransaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedTransaction.Unmarshal(m, b)
}
func (m *SignedTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedTransaction.Marshal(b, m, deterministic)
}
func (m *SignedTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedTransaction.Merge(m, src)
}
func (m *SignedTransaction) XXX_Size() int {
	return xxx_messageInfo_SignedTransaction.Size(m)
}
func (m *SignedTransaction) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedTransaction.DiscardUnknown(m)
}

var xxx_messageInfo_SignedTransaction proto.InternalMessageInfo

func (m *SignedTransaction) GetTxnBytes() []byte {
	if m != nil {
		return m.TxnBytes
	}
	return nil
}

// A generic structure that represents a transaction, covering all possible
// variants.
type Transaction struct {
	Transaction          []byte   `protobuf:"bytes,1,opt,name=transaction,proto3" json:"transaction,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}
func (*Transaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{2}
}

func (m *Transaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transaction.Unmarshal(m, b)
}
func (m *Transaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transaction.Marshal(b, m, deterministic)
}
func (m *Transaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transaction.Merge(m, src)
}
func (m *Transaction) XXX_Size() int {
	return xxx_messageInfo_Transaction.Size(m)
}
func (m *Transaction) XXX_DiscardUnknown() {
	xxx_messageInfo_Transaction.DiscardUnknown(m)
}

var xxx_messageInfo_Transaction proto.InternalMessageInfo

func (m *Transaction) GetTransaction() []byte {
	if m != nil {
		return m.Transaction
	}
	return nil
}

type TransactionWithProof struct {
	// The version of the returned signed transaction.
	Version uint64 `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	// The transaction itself.
	Transaction *Transaction `protobuf:"bytes,2,opt,name=transaction,proto3" json:"transaction,omitempty"`
	// The proof authenticating the transaction.
	Proof *TransactionProof `protobuf:"bytes,3,opt,name=proof,proto3" json:"proof,omitempty"`
	// The events yielded by executing the transaction, if requested.
	Events               *EventsList `protobuf:"bytes,4,opt,name=events,proto3" json:"events,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *TransactionWithProof) Reset()         { *m = TransactionWithProof{} }
func (m *TransactionWithProof) String() string { return proto.CompactTextString(m) }
func (*TransactionWithProof) ProtoMessage()    {}
func (*TransactionWithProof) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{3}
}

func (m *TransactionWithProof) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionWithProof.Unmarshal(m, b)
}
func (m *TransactionWithProof) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionWithProof.Marshal(b, m, deterministic)
}
func (m *TransactionWithProof) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionWithProof.Merge(m, src)
}
func (m *TransactionWithProof) XXX_Size() int {
	return xxx_messageInfo_TransactionWithProof.Size(m)
}
func (m *TransactionWithProof) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionWithProof.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionWithProof proto.InternalMessageInfo

func (m *TransactionWithProof) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *TransactionWithProof) GetTransaction() *Transaction {
	if m != nil {
		return m.Transaction
	}
	return nil
}

func (m *TransactionWithProof) GetProof() *TransactionProof {
	if m != nil {
		return m.Proof
	}
	return nil
}

func (m *TransactionWithProof) GetEvents() *EventsList {
	if m != nil {
		return m.Events
	}
	return nil
}

// A generic structure that represents a block of transactions originated from a
// particular validator instance.
type SignedTransactionsBlock struct {
	// Set of Signed Transactions
	Transactions []*SignedTransaction `protobuf:"bytes,1,rep,name=transactions,proto3" json:"transactions,omitempty"`
	// Public key of the validator that created this block
	ValidatorPublicKey []byte `protobuf:"bytes,2,opt,name=validator_public_key,json=validatorPublicKey,proto3" json:"validator_public_key,omitempty"`
	// Signature of the validator that created this block
	ValidatorSignature   []byte   `protobuf:"bytes,3,opt,name=validator_signature,json=validatorSignature,proto3" json:"validator_signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignedTransactionsBlock) Reset()         { *m = SignedTransactionsBlock{} }
func (m *SignedTransactionsBlock) String() string { return proto.CompactTextString(m) }
func (*SignedTransactionsBlock) ProtoMessage()    {}
func (*SignedTransactionsBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{4}
}

func (m *SignedTransactionsBlock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedTransactionsBlock.Unmarshal(m, b)
}
func (m *SignedTransactionsBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedTransactionsBlock.Marshal(b, m, deterministic)
}
func (m *SignedTransactionsBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedTransactionsBlock.Merge(m, src)
}
func (m *SignedTransactionsBlock) XXX_Size() int {
	return xxx_messageInfo_SignedTransactionsBlock.Size(m)
}
func (m *SignedTransactionsBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedTransactionsBlock.DiscardUnknown(m)
}

var xxx_messageInfo_SignedTransactionsBlock proto.InternalMessageInfo

func (m *SignedTransactionsBlock) GetTransactions() []*SignedTransaction {
	if m != nil {
		return m.Transactions
	}
	return nil
}

func (m *SignedTransactionsBlock) GetValidatorPublicKey() []byte {
	if m != nil {
		return m.ValidatorPublicKey
	}
	return nil
}

func (m *SignedTransactionsBlock) GetValidatorSignature() []byte {
	if m != nil {
		return m.ValidatorSignature
	}
	return nil
}

// Account state as a whole.
// After execution, updates to accounts are passed in this form to storage for
// persistence.
type AccountState struct {
	// Account address
	Address []byte `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// Account state blob
	Blob                 []byte   `protobuf:"bytes,2,opt,name=blob,proto3" json:"blob,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccountState) Reset()         { *m = AccountState{} }
func (m *AccountState) String() string { return proto.CompactTextString(m) }
func (*AccountState) ProtoMessage()    {}
func (*AccountState) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{5}
}

func (m *AccountState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountState.Unmarshal(m, b)
}
func (m *AccountState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountState.Marshal(b, m, deterministic)
}
func (m *AccountState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountState.Merge(m, src)
}
func (m *AccountState) XXX_Size() int {
	return xxx_messageInfo_AccountState.Size(m)
}
func (m *AccountState) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountState.DiscardUnknown(m)
}

var xxx_messageInfo_AccountState proto.InternalMessageInfo

func (m *AccountState) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *AccountState) GetBlob() []byte {
	if m != nil {
		return m.Blob
	}
	return nil
}

// Transaction struct to commit to storage
type TransactionToCommit struct {
	// The signed transaction which was executed
	Transaction *Transaction `protobuf:"bytes,1,opt,name=transaction,proto3" json:"transaction,omitempty"`
	// State db updates
	AccountStates []*AccountState `protobuf:"bytes,2,rep,name=account_states,json=accountStates,proto3" json:"account_states,omitempty"`
	// Events yielded by the transaction.
	Events []*Event `protobuf:"bytes,3,rep,name=events,proto3" json:"events,omitempty"`
	// The amount of gas used.
	GasUsed uint64 `protobuf:"varint,4,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	// The major status of executing the transaction.
	MajorStatus          uint64   `protobuf:"varint,5,opt,name=major_status,json=majorStatus,proto3" json:"major_status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TransactionToCommit) Reset()         { *m = TransactionToCommit{} }
func (m *TransactionToCommit) String() string { return proto.CompactTextString(m) }
func (*TransactionToCommit) ProtoMessage()    {}
func (*TransactionToCommit) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{6}
}

func (m *TransactionToCommit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionToCommit.Unmarshal(m, b)
}
func (m *TransactionToCommit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionToCommit.Marshal(b, m, deterministic)
}
func (m *TransactionToCommit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionToCommit.Merge(m, src)
}
func (m *TransactionToCommit) XXX_Size() int {
	return xxx_messageInfo_TransactionToCommit.Size(m)
}
func (m *TransactionToCommit) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionToCommit.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionToCommit proto.InternalMessageInfo

func (m *TransactionToCommit) GetTransaction() *Transaction {
	if m != nil {
		return m.Transaction
	}
	return nil
}

func (m *TransactionToCommit) GetAccountStates() []*AccountState {
	if m != nil {
		return m.AccountStates
	}
	return nil
}

func (m *TransactionToCommit) GetEvents() []*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *TransactionToCommit) GetGasUsed() uint64 {
	if m != nil {
		return m.GasUsed
	}
	return 0
}

func (m *TransactionToCommit) GetMajorStatus() uint64 {
	if m != nil {
		return m.MajorStatus
	}
	return 0
}

// A list of consecutive transactions with proof. This is mainly used for state
// synchronization when a validator would request a list of transactions from a
// peer, verify the proof, execute the transactions and persist them. Note that
// the transactions are supposed to belong to the same epoch E, otherwise
// verification will fail.
type TransactionListWithProof struct {
	// The list of transactions.
	Transactions []*Transaction `protobuf:"bytes,1,rep,name=transactions,proto3" json:"transactions,omitempty"`
	// The list of corresponding Event objects (only present if fetch_events was set to true in req)
	EventsForVersions *EventsForVersions `protobuf:"bytes,2,opt,name=events_for_versions,json=eventsForVersions,proto3" json:"events_for_versions,omitempty"`
	// If the list is not empty, the version of the first transaction.
	FirstTransactionVersion *wrappers.UInt64Value `protobuf:"bytes,3,opt,name=first_transaction_version,json=firstTransactionVersion,proto3" json:"first_transaction_version,omitempty"`
	// The proof authenticating the transactions and events.When this is used
	// for state synchronization, the validator who requests the transactions
	// will provide a version in the request and the proofs will be relative to
	// the given version. When this is returned in GetTransactionsResponse, the
	// proofs will be relative to the ledger info returned in
	// UpdateToLatestLedgerResponse.
	Proof                *TransactionListProof `protobuf:"bytes,4,opt,name=proof,proto3" json:"proof,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *TransactionListWithProof) Reset()         { *m = TransactionListWithProof{} }
func (m *TransactionListWithProof) String() string { return proto.CompactTextString(m) }
func (*TransactionListWithProof) ProtoMessage()    {}
func (*TransactionListWithProof) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{7}
}

func (m *TransactionListWithProof) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionListWithProof.Unmarshal(m, b)
}
func (m *TransactionListWithProof) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionListWithProof.Marshal(b, m, deterministic)
}
func (m *TransactionListWithProof) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionListWithProof.Merge(m, src)
}
func (m *TransactionListWithProof) XXX_Size() int {
	return xxx_messageInfo_TransactionListWithProof.Size(m)
}
func (m *TransactionListWithProof) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionListWithProof.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionListWithProof proto.InternalMessageInfo

func (m *TransactionListWithProof) GetTransactions() []*Transaction {
	if m != nil {
		return m.Transactions
	}
	return nil
}

func (m *TransactionListWithProof) GetEventsForVersions() *EventsForVersions {
	if m != nil {
		return m.EventsForVersions
	}
	return nil
}

func (m *TransactionListWithProof) GetFirstTransactionVersion() *wrappers.UInt64Value {
	if m != nil {
		return m.FirstTransactionVersion
	}
	return nil
}

func (m *TransactionListWithProof) GetProof() *TransactionListProof {
	if m != nil {
		return m.Proof
	}
	return nil
}

func init() {
	proto.RegisterEnum("types.TransactionArgument_ArgType", TransactionArgument_ArgType_name, TransactionArgument_ArgType_value)
	proto.RegisterType((*TransactionArgument)(nil), "types.TransactionArgument")
	proto.RegisterType((*SignedTransaction)(nil), "types.SignedTransaction")
	proto.RegisterType((*Transaction)(nil), "types.Transaction")
	proto.RegisterType((*TransactionWithProof)(nil), "types.TransactionWithProof")
	proto.RegisterType((*SignedTransactionsBlock)(nil), "types.SignedTransactionsBlock")
	proto.RegisterType((*AccountState)(nil), "types.AccountState")
	proto.RegisterType((*TransactionToCommit)(nil), "types.TransactionToCommit")
	proto.RegisterType((*TransactionListWithProof)(nil), "types.TransactionListWithProof")
}

func init() { proto.RegisterFile("transaction.proto", fileDescriptor_2cc4e03d2c28c490) }

var fileDescriptor_2cc4e03d2c28c490 = []byte{
	// 677 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0x51, 0x53, 0xd3, 0x40,
	0x10, 0x36, 0x6d, 0xa1, 0xb0, 0x2d, 0x4e, 0x7b, 0x65, 0x24, 0x80, 0xe3, 0xd4, 0x8c, 0x0f, 0x38,
	0x23, 0x09, 0x22, 0x83, 0x23, 0xc3, 0x4b, 0x2b, 0xa8, 0x8c, 0x8e, 0x83, 0xd7, 0x82, 0xe2, 0x4b,
	0xe6, 0x92, 0x5e, 0xd3, 0x48, 0x9b, 0xcb, 0xdc, 0x5d, 0x90, 0xfe, 0x24, 0xff, 0x82, 0x4f, 0xfe,
	0x1f, 0xff, 0x84, 0x93, 0x4b, 0x02, 0x57, 0xca, 0x83, 0x6f, 0xd9, 0xdd, 0xef, 0xdb, 0xec, 0xee,
	0xb7, 0x7b, 0xd0, 0x94, 0x9c, 0x44, 0x82, 0xf8, 0x32, 0x64, 0x91, 0x1d, 0x73, 0x26, 0x19, 0x5a,
	0x90, 0xd3, 0x98, 0x8a, 0x8d, 0x26, 0xf1, 0x7d, 0x2a, 0x84, 0x1b, 0x13, 0x39, 0xca, 0x22, 0x1b,
	0x75, 0x7a, 0x45, 0x23, 0x29, 0x72, 0xab, 0x16, 0x73, 0xc6, 0x86, 0xb9, 0xf1, 0x48, 0xcb, 0xe3,
	0x86, 0xd1, 0x90, 0xe5, 0xfe, 0x27, 0x01, 0x63, 0xc1, 0x98, 0x3a, 0xca, 0xf2, 0x92, 0xa1, 0xf3,
	0x93, 0x93, 0x38, 0xa6, 0x3c, 0x4f, 0x62, 0x7d, 0x81, 0x56, 0xff, 0x96, 0xd9, 0xe1, 0x41, 0x32,
	0xa1, 0x91, 0xb4, 0x0e, 0xa0, 0xda, 0xe1, 0x41, 0x7f, 0x1a, 0x53, 0x54, 0x85, 0xf2, 0xd9, 0xfe,
	0x5e, 0xe3, 0x01, 0xaa, 0x41, 0xb5, 0x73, 0x74, 0x84, 0x8f, 0x7b, 0xbd, 0x86, 0x81, 0x00, 0x16,
	0x7b, 0x7d, 0x7c, 0xf2, 0xf9, 0x7d, 0xa3, 0x84, 0x56, 0x60, 0xb9, 0x7b, 0xd1, 0x3f, 0xee, 0x60,
	0xdc, 0xb9, 0x68, 0x94, 0xad, 0x1d, 0x68, 0xf6, 0xc2, 0x20, 0xa2, 0x03, 0x2d, 0x31, 0xda, 0x84,
	0x65, 0x79, 0x1d, 0xb9, 0xde, 0x54, 0x52, 0x61, 0x2e, 0xb4, 0x8d, 0xad, 0x3a, 0x5e, 0x92, 0xd7,
	0x51, 0x37, 0xb5, 0x2d, 0x07, 0x6a, 0x3a, 0xb6, 0x0d, 0x35, 0xad, 0x1b, 0xd3, 0x50, 0x68, 0xdd,
	0x65, 0xfd, 0x31, 0x60, 0x55, 0x63, 0x7c, 0x0d, 0xe5, 0xe8, 0x34, 0x1d, 0x06, 0x32, 0xa1, 0x7a,
	0x45, 0xb9, 0x28, 0x68, 0x15, 0x5c, 0x98, 0x68, 0x6f, 0x36, 0x69, 0xa9, 0x6d, 0x6c, 0xd5, 0x76,
	0x91, 0xad, 0x66, 0x6d, 0x6b, 0xb9, 0x66, 0x7e, 0x84, 0xb6, 0x61, 0x41, 0x4d, 0xd9, 0x2c, 0x2b,
	0xfc, 0xda, 0x3c, 0x5e, 0xfd, 0x17, 0x67, 0x28, 0xf4, 0x1c, 0x16, 0x33, 0x89, 0xcc, 0x8a, 0xc2,
	0x37, 0x73, 0xfc, 0xb1, 0x72, 0x7e, 0x0a, 0x85, 0xc4, 0x39, 0xc0, 0xfa, 0x6d, 0xc0, 0xda, 0xdc,
	0x98, 0x44, 0x77, 0xcc, 0xfc, 0x4b, 0x74, 0x08, 0x75, 0xad, 0x08, 0x61, 0x1a, 0xed, 0xf2, 0x56,
	0x6d, 0xd7, 0xcc, 0x93, 0xcd, 0xb1, 0xf0, 0x0c, 0x1a, 0xed, 0xc0, 0xea, 0x15, 0x19, 0x87, 0x03,
	0x22, 0x19, 0x77, 0xe3, 0xc4, 0x1b, 0x87, 0xbe, 0x7b, 0x49, 0xa7, 0xaa, 0xe5, 0x3a, 0x46, 0x37,
	0xb1, 0x53, 0x15, 0xfa, 0x48, 0xa7, 0xc8, 0x81, 0xd6, 0x2d, 0x43, 0x84, 0x41, 0x44, 0x64, 0xc2,
	0xa9, 0xea, 0x59, 0x27, 0xf4, 0x8a, 0x88, 0x75, 0x08, 0xf5, 0x8e, 0xef, 0xb3, 0x24, 0x92, 0x3d,
	0x49, 0x24, 0x4d, 0xc7, 0x4e, 0x06, 0x03, 0x4e, 0x85, 0xc8, 0xd5, 0x2a, 0x4c, 0x84, 0xa0, 0xe2,
	0x8d, 0x99, 0x97, 0xff, 0x5c, 0x7d, 0x5b, 0x7f, 0x8d, 0x99, 0xa5, 0xeb, 0xb3, 0xb7, 0x6c, 0x32,
	0x09, 0xe5, 0x5d, 0x89, 0x8c, 0xff, 0x93, 0xe8, 0x00, 0x1e, 0x92, 0xac, 0x16, 0x57, 0xa4, 0xc5,
	0x08, 0xb3, 0xa4, 0xc6, 0xd5, 0xca, 0x89, 0x7a, 0xa1, 0x78, 0x85, 0x68, 0x96, 0x40, 0xcf, 0x6e,
	0xf4, 0x2a, 0x2b, 0x4e, 0x5d, 0xd7, 0xab, 0x90, 0x0a, 0xad, 0xc3, 0x52, 0x40, 0x84, 0x9b, 0x08,
	0x3a, 0x50, 0xba, 0x56, 0x70, 0x35, 0x20, 0xe2, 0x4c, 0xd0, 0x01, 0x7a, 0x0a, 0xf5, 0x09, 0xf9,
	0x91, 0x4e, 0x4d, 0x12, 0x99, 0x64, 0x9b, 0x5d, 0xc1, 0x35, 0xe5, 0xeb, 0x29, 0x97, 0xf5, 0xab,
	0x04, 0xa6, 0x56, 0x7c, 0xba, 0x04, 0xb7, 0xfb, 0xba, 0x7f, 0xaf, 0xd2, 0xf7, 0xf5, 0x3c, 0xab,
	0xf1, 0x07, 0x68, 0x65, 0xc5, 0xb9, 0x43, 0xc6, 0xdd, 0x7c, 0xc7, 0x45, 0xbe, 0xd5, 0xe6, 0xcc,
	0xd6, 0xbd, 0x63, 0xfc, 0x3c, 0x8f, 0xe3, 0x26, 0xbd, 0xeb, 0x42, 0xdf, 0x60, 0x7d, 0x18, 0x72,
	0x21, 0x5d, 0xfd, 0x01, 0x29, 0x6e, 0x28, 0xdb, 0xfa, 0xc7, 0x76, 0xf6, 0x88, 0xd8, 0xc5, 0x23,
	0x62, 0x9f, 0x9d, 0x44, 0x72, 0x7f, 0xef, 0x9c, 0x8c, 0x13, 0x8a, 0xd7, 0x14, 0x5d, 0x2b, 0x35,
	0x4f, 0x8d, 0x5e, 0x16, 0xb7, 0x93, 0xdd, 0xc2, 0xe6, 0x7c, 0x53, 0xe9, 0x2c, 0xf4, 0xfb, 0xe9,
	0xda, 0xdf, 0x5f, 0x04, 0xa1, 0x1c, 0x25, 0x9e, 0xed, 0xb3, 0x89, 0x23, 0x47, 0xf4, 0xf5, 0xee,
	0x1b, 0x27, 0x60, 0xdb, 0xe3, 0xd0, 0xe3, 0xc4, 0x09, 0x68, 0x44, 0x39, 0x91, 0x74, 0xe0, 0xc4,
	0x9e, 0xca, 0xe5, 0x2d, 0xaa, 0x8a, 0x5e, 0xfd, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x01, 0x0a, 0x89,
	0xc4, 0x46, 0x05, 0x00, 0x00,
}
