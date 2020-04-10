// +build !js
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: validator_info.proto

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

// Protobuf definition for the Rust struct ValidatorInfo
type ValidatorInfo struct {
	// Validator account address
	AccountAddress []byte `protobuf:"bytes,1,opt,name=account_address,json=accountAddress,proto3" json:"account_address,omitempty"`
	// Consensus public key
	ConsensusPublicKey []byte `protobuf:"bytes,2,opt,name=consensus_public_key,json=consensusPublicKey,proto3" json:"consensus_public_key,omitempty"`
	// Validator voting power for consensus
	ConsensusVotingPower uint64 `protobuf:"varint,3,opt,name=consensus_voting_power,json=consensusVotingPower,proto3" json:"consensus_voting_power,omitempty"`
	// Network signing publick key
	NetworkSigningPublicKey []byte `protobuf:"bytes,4,opt,name=network_signing_public_key,json=networkSigningPublicKey,proto3" json:"network_signing_public_key,omitempty"`
	/// Network identity publick key
	NetworkIdentityPublicKey []byte   `protobuf:"bytes,5,opt,name=network_identity_public_key,json=networkIdentityPublicKey,proto3" json:"network_identity_public_key,omitempty"`
	XXX_NoUnkeyedLiteral     struct{} `json:"-"`
	XXX_unrecognized         []byte   `json:"-"`
	XXX_sizecache            int32    `json:"-"`
}

func (m *ValidatorInfo) Reset()         { *m = ValidatorInfo{} }
func (m *ValidatorInfo) String() string { return proto.CompactTextString(m) }
func (*ValidatorInfo) ProtoMessage()    {}
func (*ValidatorInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_58482d830cf3ea80, []int{0}
}

func (m *ValidatorInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidatorInfo.Unmarshal(m, b)
}
func (m *ValidatorInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidatorInfo.Marshal(b, m, deterministic)
}
func (m *ValidatorInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidatorInfo.Merge(m, src)
}
func (m *ValidatorInfo) XXX_Size() int {
	return xxx_messageInfo_ValidatorInfo.Size(m)
}
func (m *ValidatorInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidatorInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ValidatorInfo proto.InternalMessageInfo

func (m *ValidatorInfo) GetAccountAddress() []byte {
	if m != nil {
		return m.AccountAddress
	}
	return nil
}

func (m *ValidatorInfo) GetConsensusPublicKey() []byte {
	if m != nil {
		return m.ConsensusPublicKey
	}
	return nil
}

func (m *ValidatorInfo) GetConsensusVotingPower() uint64 {
	if m != nil {
		return m.ConsensusVotingPower
	}
	return 0
}

func (m *ValidatorInfo) GetNetworkSigningPublicKey() []byte {
	if m != nil {
		return m.NetworkSigningPublicKey
	}
	return nil
}

func (m *ValidatorInfo) GetNetworkIdentityPublicKey() []byte {
	if m != nil {
		return m.NetworkIdentityPublicKey
	}
	return nil
}

func init() {
	proto.RegisterType((*ValidatorInfo)(nil), "types.ValidatorInfo")
}

func init() {
	proto.RegisterFile("validator_info.proto", fileDescriptor_58482d830cf3ea80)
}

var fileDescriptor_58482d830cf3ea80 = []byte{
	// 264 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0x4f, 0x4b, 0xc3, 0x40,
	0x10, 0x47, 0x49, 0x6d, 0x3d, 0x2c, 0xfe, 0x81, 0xa5, 0x68, 0xd0, 0x4b, 0xf1, 0x62, 0x0f, 0x9a,
	0x88, 0x0a, 0x22, 0xe2, 0x41, 0x6f, 0xc5, 0x4b, 0xa9, 0xd0, 0x83, 0x97, 0x65, 0x93, 0x4c, 0xd3,
	0xa5, 0x71, 0x66, 0xd9, 0x9d, 0xb4, 0xe4, 0xbb, 0xf8, 0x61, 0xc5, 0x6d, 0xda, 0x78, 0xdd, 0xf7,
	0xde, 0xfe, 0x60, 0xc4, 0x70, 0xad, 0x2b, 0x53, 0x68, 0x26, 0xa7, 0x0c, 0x2e, 0x28, 0xb1, 0x8e,
	0x98, 0xe4, 0x80, 0x1b, 0x0b, 0xfe, 0xea, 0xa7, 0x27, 0x8e, 0xe7, 0x3b, 0x3e, 0xc1, 0x05, 0xc9,
	0x6b, 0x71, 0xaa, 0xf3, 0x9c, 0x6a, 0x64, 0xa5, 0x8b, 0xc2, 0x81, 0xf7, 0x71, 0x34, 0x8a, 0xc6,
	0x47, 0xb3, 0x93, 0xf6, 0xf9, 0x6d, 0xfb, 0x2a, 0xef, 0xc4, 0x30, 0x27, 0xf4, 0x80, 0xbe, 0xf6,
	0xca, 0xd6, 0x59, 0x65, 0x72, 0xb5, 0x82, 0x26, 0xee, 0x05, 0x5b, 0xee, 0xd9, 0x34, 0xa0, 0x0f,
	0x68, 0xe4, 0xa3, 0x38, 0xeb, 0x8a, 0x35, 0xb1, 0xc1, 0x52, 0x59, 0xda, 0x80, 0x8b, 0x0f, 0x46,
	0xd1, 0xb8, 0x3f, 0xeb, 0xfe, 0x9b, 0x07, 0x38, 0xfd, 0x63, 0xf2, 0x45, 0x5c, 0x20, 0xf0, 0x86,
	0xdc, 0x4a, 0x79, 0x53, 0x62, 0x88, 0xba, 0xb5, 0x7e, 0x58, 0x3b, 0x6f, 0x8d, 0xcf, 0xad, 0xd0,
	0x4d, 0xbe, 0x8a, 0xcb, 0x5d, 0x6c, 0x0a, 0x40, 0x36, 0xdc, 0xfc, 0xaf, 0x07, 0xa1, 0x8e, 0x5b,
	0x65, 0xd2, 0x1a, 0xfb, 0xfc, 0x3d, 0xf9, 0xba, 0x29, 0x0d, 0x2f, 0xeb, 0x2c, 0xc9, 0xe9, 0x3b,
	0xe5, 0x25, 0x3c, 0xdd, 0x3f, 0xa7, 0x25, 0xdd, 0x56, 0x26, 0x73, 0x3a, 0x2d, 0x01, 0xc1, 0x69,
	0x86, 0x22, 0xb5, 0x59, 0x38, 0x67, 0x76, 0x18, 0x8e, 0xfb, 0xf0, 0x1b, 0x00, 0x00, 0xff, 0xff,
	0xea, 0xa2, 0x74, 0x6c, 0x74, 0x01, 0x00, 0x00,
}
