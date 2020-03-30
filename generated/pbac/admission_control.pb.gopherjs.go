// +build js
// Code generated by protoc-gen-gopherjs. DO NOT EDIT.
// source: admission_control.proto

/*
	Package pbac is a generated protocol buffer package.

	It is generated from these files:
		admission_control.proto

	It has these top-level messages:
		AdmissionControlMsg
		SubmitTransactionRequest
		AdmissionControlStatus
		SubmitTransactionResponse
*/
package pbac

import jspb "github.com/johanbrandhorst/protobuf/jspb"
import types10 "github.com/the729/go-libra/generated/pbtypes"
import types11 "github.com/the729/go-libra/generated/pbtypes"
import types8 "github.com/the729/go-libra/generated/pbtypes"
import types12 "github.com/the729/go-libra/generated/pbtypes"

import (
	context "context"

	grpcweb "github.com/johanbrandhorst/protobuf/grpcweb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the jspb package it is being compiled against.
const _ = jspb.JspbPackageIsVersion2

// Additional statuses that are possible from admission control in addition
// to VM statuses.
type AdmissionControlStatusCode int

const (
	// Validator accepted the transaction.
	AdmissionControlStatusCode_Accepted AdmissionControlStatusCode = 0
	// The sender is blacklisted.
	AdmissionControlStatusCode_Blacklisted AdmissionControlStatusCode = 1
	// The transaction is rejected, e.g. due to incorrect signature.
	AdmissionControlStatusCode_Rejected AdmissionControlStatusCode = 2
)

var AdmissionControlStatusCode_name = map[int]string{
	0: "Accepted",
	1: "Blacklisted",
	2: "Rejected",
}
var AdmissionControlStatusCode_value = map[string]int{
	"Accepted":    0,
	"Blacklisted": 1,
	"Rejected":    2,
}

func (x AdmissionControlStatusCode) String() string {
	return AdmissionControlStatusCode_name[int(x)]
}

// The request for submitting a transaction to an upstream validator or full
// node.
type AdmissionControlMsg struct {
	// Types that are valid to be assigned to Message:
	//	*AdmissionControlMsg_SubmitTransactionRequest
	//	*AdmissionControlMsg_SubmitTransactionResponse
	Message isAdmissionControlMsg_Message
}

// isAdmissionControlMsg_Message is used to distinguish types assignable to Message
type isAdmissionControlMsg_Message interface{ isAdmissionControlMsg_Message() }

// AdmissionControlMsg_SubmitTransactionRequest is assignable to Message
type AdmissionControlMsg_SubmitTransactionRequest struct {
	SubmitTransactionRequest *SubmitTransactionRequest
}

// AdmissionControlMsg_SubmitTransactionResponse is assignable to Message
type AdmissionControlMsg_SubmitTransactionResponse struct {
	SubmitTransactionResponse *SubmitTransactionResponse
}

func (*AdmissionControlMsg_SubmitTransactionRequest) isAdmissionControlMsg_Message()  {}
func (*AdmissionControlMsg_SubmitTransactionResponse) isAdmissionControlMsg_Message() {}

// GetMessage gets the Message of the AdmissionControlMsg.
func (m *AdmissionControlMsg) GetMessage() (x isAdmissionControlMsg_Message) {
	if m == nil {
		return x
	}
	return m.Message
}

// GetSubmitTransactionRequest gets the SubmitTransactionRequest of the AdmissionControlMsg.
func (m *AdmissionControlMsg) GetSubmitTransactionRequest() (x *SubmitTransactionRequest) {
	if v, ok := m.GetMessage().(*AdmissionControlMsg_SubmitTransactionRequest); ok {
		return v.SubmitTransactionRequest
	}
	return x
}

// GetSubmitTransactionResponse gets the SubmitTransactionResponse of the AdmissionControlMsg.
func (m *AdmissionControlMsg) GetSubmitTransactionResponse() (x *SubmitTransactionResponse) {
	if v, ok := m.GetMessage().(*AdmissionControlMsg_SubmitTransactionResponse); ok {
		return v.SubmitTransactionResponse
	}
	return x
}

// MarshalToWriter marshals AdmissionControlMsg to the provided writer.
func (m *AdmissionControlMsg) MarshalToWriter(writer jspb.Writer) {
	if m == nil {
		return
	}

	switch t := m.Message.(type) {
	case *AdmissionControlMsg_SubmitTransactionRequest:
		if t.SubmitTransactionRequest != nil {
			writer.WriteMessage(1, func() {
				t.SubmitTransactionRequest.MarshalToWriter(writer)
			})
		}
	case *AdmissionControlMsg_SubmitTransactionResponse:
		if t.SubmitTransactionResponse != nil {
			writer.WriteMessage(2, func() {
				t.SubmitTransactionResponse.MarshalToWriter(writer)
			})
		}
	}

	return
}

// Marshal marshals AdmissionControlMsg to a slice of bytes.
func (m *AdmissionControlMsg) Marshal() []byte {
	writer := jspb.NewWriter()
	m.MarshalToWriter(writer)
	return writer.GetResult()
}

// UnmarshalFromReader unmarshals a AdmissionControlMsg from the provided reader.
func (m *AdmissionControlMsg) UnmarshalFromReader(reader jspb.Reader) *AdmissionControlMsg {
	for reader.Next() {
		if m == nil {
			m = &AdmissionControlMsg{}
		}

		switch reader.GetFieldNumber() {
		case 1:
			reader.ReadMessage(func() {
				m.Message = &AdmissionControlMsg_SubmitTransactionRequest{
					SubmitTransactionRequest: new(SubmitTransactionRequest).UnmarshalFromReader(reader),
				}
			})
		case 2:
			reader.ReadMessage(func() {
				m.Message = &AdmissionControlMsg_SubmitTransactionResponse{
					SubmitTransactionResponse: new(SubmitTransactionResponse).UnmarshalFromReader(reader),
				}
			})
		default:
			reader.SkipField()
		}
	}

	return m
}

// Unmarshal unmarshals a AdmissionControlMsg from a slice of bytes.
func (m *AdmissionControlMsg) Unmarshal(rawBytes []byte) (*AdmissionControlMsg, error) {
	reader := jspb.NewReader(rawBytes)

	m = m.UnmarshalFromReader(reader)

	if err := reader.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

// -----------------------------------------------------------------------------
// ---------------- Submit transaction
// -----------------------------------------------------------------------------
// The request for transaction submission.
type SubmitTransactionRequest struct {
	// Transaction submitted by user.
	Transaction *types8.SignedTransaction
}

// GetTransaction gets the Transaction of the SubmitTransactionRequest.
func (m *SubmitTransactionRequest) GetTransaction() (x *types8.SignedTransaction) {
	if m == nil {
		return x
	}
	return m.Transaction
}

// MarshalToWriter marshals SubmitTransactionRequest to the provided writer.
func (m *SubmitTransactionRequest) MarshalToWriter(writer jspb.Writer) {
	if m == nil {
		return
	}

	if m.Transaction != nil {
		writer.WriteMessage(1, func() {
			m.Transaction.MarshalToWriter(writer)
		})
	}

	return
}

// Marshal marshals SubmitTransactionRequest to a slice of bytes.
func (m *SubmitTransactionRequest) Marshal() []byte {
	writer := jspb.NewWriter()
	m.MarshalToWriter(writer)
	return writer.GetResult()
}

// UnmarshalFromReader unmarshals a SubmitTransactionRequest from the provided reader.
func (m *SubmitTransactionRequest) UnmarshalFromReader(reader jspb.Reader) *SubmitTransactionRequest {
	for reader.Next() {
		if m == nil {
			m = &SubmitTransactionRequest{}
		}

		switch reader.GetFieldNumber() {
		case 1:
			reader.ReadMessage(func() {
				m.Transaction = m.Transaction.UnmarshalFromReader(reader)
			})
		default:
			reader.SkipField()
		}
	}

	return m
}

// Unmarshal unmarshals a SubmitTransactionRequest from a slice of bytes.
func (m *SubmitTransactionRequest) Unmarshal(rawBytes []byte) (*SubmitTransactionRequest, error) {
	reader := jspb.NewReader(rawBytes)

	m = m.UnmarshalFromReader(reader)

	if err := reader.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

// AC response status containing code and optionally an error message.
type AdmissionControlStatus struct {
	Code    AdmissionControlStatusCode
	Message string
}

// GetCode gets the Code of the AdmissionControlStatus.
func (m *AdmissionControlStatus) GetCode() (x AdmissionControlStatusCode) {
	if m == nil {
		return x
	}
	return m.Code
}

// GetMessage gets the Message of the AdmissionControlStatus.
func (m *AdmissionControlStatus) GetMessage() (x string) {
	if m == nil {
		return x
	}
	return m.Message
}

// MarshalToWriter marshals AdmissionControlStatus to the provided writer.
func (m *AdmissionControlStatus) MarshalToWriter(writer jspb.Writer) {
	if m == nil {
		return
	}

	if int(m.Code) != 0 {
		writer.WriteEnum(1, int(m.Code))
	}

	if len(m.Message) > 0 {
		writer.WriteString(2, m.Message)
	}

	return
}

// Marshal marshals AdmissionControlStatus to a slice of bytes.
func (m *AdmissionControlStatus) Marshal() []byte {
	writer := jspb.NewWriter()
	m.MarshalToWriter(writer)
	return writer.GetResult()
}

// UnmarshalFromReader unmarshals a AdmissionControlStatus from the provided reader.
func (m *AdmissionControlStatus) UnmarshalFromReader(reader jspb.Reader) *AdmissionControlStatus {
	for reader.Next() {
		if m == nil {
			m = &AdmissionControlStatus{}
		}

		switch reader.GetFieldNumber() {
		case 1:
			m.Code = AdmissionControlStatusCode(reader.ReadEnum())
		case 2:
			m.Message = reader.ReadString()
		default:
			reader.SkipField()
		}
	}

	return m
}

// Unmarshal unmarshals a AdmissionControlStatus from a slice of bytes.
func (m *AdmissionControlStatus) Unmarshal(rawBytes []byte) (*AdmissionControlStatus, error) {
	reader := jspb.NewReader(rawBytes)

	m = m.UnmarshalFromReader(reader)

	if err := reader.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

// The response for transaction submission.
//
// How does a client know if their transaction was included?
// A response from the transaction submission only means that the transaction
// was successfully added to mempool, but not that it is guaranteed to be
// included in the chain.  Each transaction should include an expiration time in
// the signed transaction.  Let's call this T0.  As a client, I submit my
// transaction to a validator. I now need to poll for the transaction I
// submitted.  I can use the query that takes my account and sequence number. If
// I receive back that the transaction is completed, I will verify the proofs to
// ensure that this is the transaction I expected.  If I receive a response that
// my transaction is not yet completed, I must check the latest timestamp in the
// ledgerInfo that I receive back from the query.  If this time is greater than
// T0, I can be certain that my transaction will never be included.  If this
// time is less than T0, I need to continue polling.
type SubmitTransactionResponse struct {
	// The status of a transaction submission can either be a VM status, or
	// some other admission control/mempool specific status e.g. Blacklisted.
	//
	// Types that are valid to be assigned to Status:
	//	*SubmitTransactionResponse_VmStatus
	//	*SubmitTransactionResponse_AcStatus
	//	*SubmitTransactionResponse_MempoolStatus
	Status isSubmitTransactionResponse_Status
	// Public key(id) of the validator that processed this transaction
	ValidatorId []byte
}

// isSubmitTransactionResponse_Status is used to distinguish types assignable to Status
type isSubmitTransactionResponse_Status interface{ isSubmitTransactionResponse_Status() }

// SubmitTransactionResponse_VmStatus is assignable to Status
type SubmitTransactionResponse_VmStatus struct {
	VmStatus *types12.VMStatus
}

// SubmitTransactionResponse_AcStatus is assignable to Status
type SubmitTransactionResponse_AcStatus struct {
	AcStatus *AdmissionControlStatus
}

// SubmitTransactionResponse_MempoolStatus is assignable to Status
type SubmitTransactionResponse_MempoolStatus struct {
	MempoolStatus *types11.MempoolStatus
}

func (*SubmitTransactionResponse_VmStatus) isSubmitTransactionResponse_Status()      {}
func (*SubmitTransactionResponse_AcStatus) isSubmitTransactionResponse_Status()      {}
func (*SubmitTransactionResponse_MempoolStatus) isSubmitTransactionResponse_Status() {}

// GetStatus gets the Status of the SubmitTransactionResponse.
func (m *SubmitTransactionResponse) GetStatus() (x isSubmitTransactionResponse_Status) {
	if m == nil {
		return x
	}
	return m.Status
}

// GetVmStatus gets the VmStatus of the SubmitTransactionResponse.
func (m *SubmitTransactionResponse) GetVmStatus() (x *types12.VMStatus) {
	if v, ok := m.GetStatus().(*SubmitTransactionResponse_VmStatus); ok {
		return v.VmStatus
	}
	return x
}

// GetAcStatus gets the AcStatus of the SubmitTransactionResponse.
func (m *SubmitTransactionResponse) GetAcStatus() (x *AdmissionControlStatus) {
	if v, ok := m.GetStatus().(*SubmitTransactionResponse_AcStatus); ok {
		return v.AcStatus
	}
	return x
}

// GetMempoolStatus gets the MempoolStatus of the SubmitTransactionResponse.
func (m *SubmitTransactionResponse) GetMempoolStatus() (x *types11.MempoolStatus) {
	if v, ok := m.GetStatus().(*SubmitTransactionResponse_MempoolStatus); ok {
		return v.MempoolStatus
	}
	return x
}

// GetValidatorId gets the ValidatorId of the SubmitTransactionResponse.
func (m *SubmitTransactionResponse) GetValidatorId() (x []byte) {
	if m == nil {
		return x
	}
	return m.ValidatorId
}

// MarshalToWriter marshals SubmitTransactionResponse to the provided writer.
func (m *SubmitTransactionResponse) MarshalToWriter(writer jspb.Writer) {
	if m == nil {
		return
	}

	switch t := m.Status.(type) {
	case *SubmitTransactionResponse_VmStatus:
		if t.VmStatus != nil {
			writer.WriteMessage(1, func() {
				t.VmStatus.MarshalToWriter(writer)
			})
		}
	case *SubmitTransactionResponse_AcStatus:
		if t.AcStatus != nil {
			writer.WriteMessage(2, func() {
				t.AcStatus.MarshalToWriter(writer)
			})
		}
	case *SubmitTransactionResponse_MempoolStatus:
		if t.MempoolStatus != nil {
			writer.WriteMessage(3, func() {
				t.MempoolStatus.MarshalToWriter(writer)
			})
		}
	}

	if len(m.ValidatorId) > 0 {
		writer.WriteBytes(4, m.ValidatorId)
	}

	return
}

// Marshal marshals SubmitTransactionResponse to a slice of bytes.
func (m *SubmitTransactionResponse) Marshal() []byte {
	writer := jspb.NewWriter()
	m.MarshalToWriter(writer)
	return writer.GetResult()
}

// UnmarshalFromReader unmarshals a SubmitTransactionResponse from the provided reader.
func (m *SubmitTransactionResponse) UnmarshalFromReader(reader jspb.Reader) *SubmitTransactionResponse {
	for reader.Next() {
		if m == nil {
			m = &SubmitTransactionResponse{}
		}

		switch reader.GetFieldNumber() {
		case 1:
			reader.ReadMessage(func() {
				m.Status = &SubmitTransactionResponse_VmStatus{
					VmStatus: new(types12.VMStatus).UnmarshalFromReader(reader),
				}
			})
		case 2:
			reader.ReadMessage(func() {
				m.Status = &SubmitTransactionResponse_AcStatus{
					AcStatus: new(AdmissionControlStatus).UnmarshalFromReader(reader),
				}
			})
		case 3:
			reader.ReadMessage(func() {
				m.Status = &SubmitTransactionResponse_MempoolStatus{
					MempoolStatus: new(types11.MempoolStatus).UnmarshalFromReader(reader),
				}
			})
		case 4:
			m.ValidatorId = reader.ReadBytes()
		default:
			reader.SkipField()
		}
	}

	return m
}

// Unmarshal unmarshals a SubmitTransactionResponse from a slice of bytes.
func (m *SubmitTransactionResponse) Unmarshal(rawBytes []byte) (*SubmitTransactionResponse, error) {
	reader := jspb.NewReader(rawBytes)

	m = m.UnmarshalFromReader(reader)

	if err := reader.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpcweb.Client

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpcweb package it is being compiled against.
const _ = grpcweb.GrpcWebPackageIsVersion3

// Client API for AdmissionControl service

// -----------------------------------------------------------------------------
// ---------------- Service definition
// -----------------------------------------------------------------------------
type AdmissionControlClient interface {
	// Public API to submit transaction to a validator.
	SubmitTransaction(ctx context.Context, in *SubmitTransactionRequest, opts ...grpcweb.CallOption) (*SubmitTransactionResponse, error)
	// This API is used to update the client to the latest ledger version and
	// optionally also request 1..n other pieces of data.  This allows for batch
	// queries.  All queries return proofs that a client should check to validate
	// the data. Note that if a client only wishes to update to the latest
	// LedgerInfo and receive the proof of this latest version, they can simply
	// omit the requested_items (or pass an empty list)
	UpdateToLatestLedger(ctx context.Context, in *types10.UpdateToLatestLedgerRequest, opts ...grpcweb.CallOption) (*types10.UpdateToLatestLedgerResponse, error)
}

type admissionControlClient struct {
	client *grpcweb.Client
}

// NewAdmissionControlClient creates a new gRPC-Web client.
func NewAdmissionControlClient(hostname string, opts ...grpcweb.DialOption) AdmissionControlClient {
	return &admissionControlClient{
		client: grpcweb.NewClient(hostname, "admission_control.AdmissionControl", opts...),
	}
}

func (c *admissionControlClient) SubmitTransaction(ctx context.Context, in *SubmitTransactionRequest, opts ...grpcweb.CallOption) (*SubmitTransactionResponse, error) {
	resp, err := c.client.RPCCall(ctx, "SubmitTransaction", in.Marshal(), opts...)
	if err != nil {
		return nil, err
	}

	return new(SubmitTransactionResponse).Unmarshal(resp)
}

func (c *admissionControlClient) UpdateToLatestLedger(ctx context.Context, in *types10.UpdateToLatestLedgerRequest, opts ...grpcweb.CallOption) (*types10.UpdateToLatestLedgerResponse, error) {
	resp, err := c.client.RPCCall(ctx, "UpdateToLatestLedger", in.Marshal(), opts...)
	if err != nil {
		return nil, err
	}

	return new(types10.UpdateToLatestLedgerResponse).Unmarshal(resp)
}
