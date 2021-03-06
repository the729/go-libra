// +build !js
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mempool_status.proto

package pbtypes

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// The statuses and errors produced by the Mempool during transaction insertion
type MempoolStatus struct {
	// e.g. assertion violation
	Code                 uint64   `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message              string   `protobuf:"bytes,5,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MempoolStatus) Reset()         { *m = MempoolStatus{} }
func (m *MempoolStatus) String() string { return proto.CompactTextString(m) }
func (*MempoolStatus) ProtoMessage()    {}
func (*MempoolStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_cad4a86f8a5465be, []int{0}
}

func (m *MempoolStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MempoolStatus.Unmarshal(m, b)
}
func (m *MempoolStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MempoolStatus.Marshal(b, m, deterministic)
}
func (m *MempoolStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MempoolStatus.Merge(m, src)
}
func (m *MempoolStatus) XXX_Size() int {
	return xxx_messageInfo_MempoolStatus.Size(m)
}
func (m *MempoolStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_MempoolStatus.DiscardUnknown(m)
}

var xxx_messageInfo_MempoolStatus proto.InternalMessageInfo

func (m *MempoolStatus) GetCode() uint64 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *MempoolStatus) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*MempoolStatus)(nil), "types.MempoolStatus")
}

func init() {
	proto.RegisterFile("mempool_status.proto", fileDescriptor_cad4a86f8a5465be)
}

var fileDescriptor_cad4a86f8a5465be = []byte{
	// 150 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc9, 0x4d, 0xcd, 0x2d,
	0xc8, 0xcf, 0xcf, 0x89, 0x2f, 0x2e, 0x49, 0x2c, 0x29, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x62, 0x2d, 0xa9, 0x2c, 0x48, 0x2d, 0x56, 0xb2, 0xe5, 0xe2, 0xf5, 0x85, 0x48, 0x07, 0x83,
	0x65, 0x85, 0x84, 0xb8, 0x58, 0x92, 0xf3, 0x53, 0x52, 0x25, 0x18, 0x15, 0x18, 0x35, 0x58, 0x82,
	0xc0, 0x6c, 0x21, 0x09, 0x2e, 0xf6, 0xdc, 0xd4, 0xe2, 0xe2, 0xc4, 0xf4, 0x54, 0x09, 0x56, 0x05,
	0x46, 0x0d, 0xce, 0x20, 0x18, 0xd7, 0x49, 0x2f, 0x4a, 0x27, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49,
	0x2f, 0x39, 0x3f, 0x57, 0xbf, 0x24, 0x23, 0xd5, 0xdc, 0xc8, 0x52, 0x3f, 0x3d, 0x5f, 0x37, 0x27,
	0x33, 0xa9, 0x28, 0x51, 0x3f, 0x3d, 0x35, 0x2f, 0xb5, 0x28, 0xb1, 0x24, 0x35, 0x45, 0xbf, 0x20,
	0x09, 0x6c, 0x5d, 0x12, 0x1b, 0xd8, 0x72, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xec, 0xc4,
	0x0a, 0x75, 0x94, 0x00, 0x00, 0x00,
}
