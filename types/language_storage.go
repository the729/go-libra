package types

import (
	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/lcs"
)

const (
	CodeTag     = 0
	ResourceTag = 1
)

// AccessPathTag is an interface that can be used to construct the root of an access path
type AccessPathTag interface {
	Hash() HashValue
	TypePrefix() byte
}

type TypeTag interface {
	Clone() TypeTag
}

// TypeTagBool is bool
type TypeTagBool bool

// TypeTagU8 is uint8
type TypeTagU8 uint8

// TypeTagU64 is uint64
type TypeTagU64 uint64

// TypeTagU128 is [16]byte
type TypeTagU128 [16]byte

// TypeTagAddress is account address type
type TypeTagAddress AccountAddress

// TypeTagTypeTags is a vector of TypeTags
type TypeTagTypeTags []TypeTag

// TypeTagStructTag is StructTag
type TypeTagStructTag = StructTag

func (v TypeTagBool) Clone() TypeTag    { return v }
func (v TypeTagU8) Clone() TypeTag      { return v }
func (v TypeTagU64) Clone() TypeTag     { return v }
func (v TypeTagU128) Clone() TypeTag    { return v }
func (v TypeTagAddress) Clone() TypeTag { return v }
func (v TypeTagTypeTags) Clone() TypeTag {
	n := make([]TypeTag, 0, len(v))
	for _, e := range v {
		n = append(n, e.Clone())
	}
	return TypeTagTypeTags(n)
}

var typeTagEnumDef = []lcs.EnumVariant{
	{
		Name:     "TypeTag",
		Value:    0,
		Template: TypeTagBool(false),
	},
	{
		Name:     "TypeTag",
		Value:    1,
		Template: TypeTagU8(0),
	},
	{
		Name:     "TypeTag",
		Value:    2,
		Template: TypeTagU64(0),
	},
	{
		Name:     "TypeTag",
		Value:    3,
		Template: TypeTagU128([16]byte{}),
	},
	{
		Name:     "TypeTag",
		Value:    4,
		Template: TypeTagAddress([AccountAddressLength]byte{}),
	},
	{
		Name:     "TypeTag",
		Value:    5,
		Template: TypeTagTypeTags(nil),
	},
	{
		Name:     "TypeTag",
		Value:    6,
		Template: (*TypeTagStructTag)(nil),
	},
}

type TypeTagWrap struct {
	Value TypeTag `lcs:"enum=TypeTag"`
}

// EnumTypes defines enum variants for lcs
func (*TypeTagWrap) EnumTypes() []lcs.EnumVariant { return typeTagEnumDef }

// StructTag is a tag to form a resource path.
//
// StructTag implements AccessPathTag interface
type StructTag struct {
	Address    AccountAddress
	Module     string
	Name       string
	TypeParams []TypeTag `lcs:"enum=TypeTag"`
}

// EnumTypes defines enum variants for lcs
func (*StructTag) EnumTypes() []lcs.EnumVariant { return typeTagEnumDef }

// Hash outputs the hash of this struct, using the appropriate hash function.
func (t *StructTag) Hash() HashValue {
	hasher := sha3libra.NewStructTag()
	if err := lcs.NewEncoder(hasher).Encode(t); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
}

func (t *StructTag) Clone() TypeTag {
	out := &StructTag{
		Address: t.Address,
		Module:  t.Module,
		Name:    t.Name,
	}
	n := make([]TypeTag, 0, len(t.TypeParams))
	for _, e := range t.TypeParams {
		n = append(n, e.Clone())
	}
	return out
}

// TypePrefix returns type byte of this tag, which is '0x01'
func (t *StructTag) TypePrefix() byte { return ResourceTag }

// RawTag is a tag with raw hash values. It implements AccessPathTag interface.
type RawTag struct {
	HashVal HashValue
	TypeVal byte
}

// Hash returns HashVal
func (t *RawTag) Hash() HashValue { return t.HashVal }

// TypePrefix returns TypeVal
func (t *RawTag) TypePrefix() byte { return t.TypeVal }
