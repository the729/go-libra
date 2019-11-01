// +build !js
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: access_path.proto

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

type AccessPath struct {
	Address              []byte   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Path                 []byte   `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccessPath) Reset()         { *m = AccessPath{} }
func (m *AccessPath) String() string { return proto.CompactTextString(m) }
func (*AccessPath) ProtoMessage()    {}
func (*AccessPath) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec5cf8547713a5d2, []int{0}
}

func (m *AccessPath) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccessPath.Unmarshal(m, b)
}
func (m *AccessPath) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccessPath.Marshal(b, m, deterministic)
}
func (m *AccessPath) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccessPath.Merge(m, src)
}
func (m *AccessPath) XXX_Size() int {
	return xxx_messageInfo_AccessPath.Size(m)
}
func (m *AccessPath) XXX_DiscardUnknown() {
	xxx_messageInfo_AccessPath.DiscardUnknown(m)
}

var xxx_messageInfo_AccessPath proto.InternalMessageInfo

func (m *AccessPath) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *AccessPath) GetPath() []byte {
	if m != nil {
		return m.Path
	}
	return nil
}

func init() {
	proto.RegisterType((*AccessPath)(nil), "types.AccessPath")
}

func init() { proto.RegisterFile("access_path.proto", fileDescriptor_ec5cf8547713a5d2) }

var fileDescriptor_ec5cf8547713a5d2 = []byte{
	// 143 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0x4c, 0x4e, 0x4e,
	0x2d, 0x2e, 0x8e, 0x2f, 0x48, 0x2c, 0xc9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d,
	0xa9, 0x2c, 0x48, 0x2d, 0x56, 0xb2, 0xe2, 0xe2, 0x72, 0x04, 0xcb, 0x05, 0x24, 0x96, 0x64, 0x08,
	0x49, 0x70, 0xb1, 0x27, 0xa6, 0xa4, 0x14, 0xa5, 0x16, 0x17, 0x4b, 0x30, 0x2a, 0x30, 0x6a, 0xf0,
	0x04, 0xc1, 0xb8, 0x42, 0x42, 0x5c, 0x2c, 0x20, 0xcd, 0x12, 0x4c, 0x60, 0x61, 0x30, 0xdb, 0x49,
	0x2f, 0x4a, 0x27, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0xbf, 0x24, 0x23,
	0xd5, 0xdc, 0xc8, 0x52, 0x3f, 0x3d, 0x5f, 0x37, 0x27, 0x33, 0xa9, 0x28, 0x51, 0x3f, 0x3d, 0x35,
	0x2f, 0xb5, 0x28, 0xb1, 0x24, 0x35, 0x45, 0xbf, 0x20, 0x09, 0x6c, 0x57, 0x12, 0x1b, 0xd8, 0x66,
	0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x33, 0xc8, 0xf8, 0xf8, 0x8e, 0x00, 0x00, 0x00,
}
