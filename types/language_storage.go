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

type isTypeTag interface {
	isTypeTag()
}

// TypeTagBool is bool
type TypeTagBool bool

// TypeTagU64 is uint64
type TypeTagU64 uint64

// TypeTagBytes is byte slice
type TypeTagBytes []byte

// TypeTagAddress is account address type
type TypeTagAddress AccountAddress

// TypeTagStructTag is StructTag
type TypeTagStructTag = StructTag

func (TypeTagBool) isTypeTag()       {}
func (TypeTagU64) isTypeTag()        {}
func (TypeTagBytes) isTypeTag()      {}
func (TypeTagAddress) isTypeTag()    {}
func (*TypeTagStructTag) isTypeTag() {}

var typeTagEnumDef = []lcs.EnumVariant{
	{
		Name:     "TypeTag",
		Value:    0,
		Template: TypeTagBool(false),
	},
	{
		Name:     "TypeTag",
		Value:    1,
		Template: TypeTagU64(0),
	},
	{
		Name:     "TypeTag",
		Value:    2,
		Template: TypeTagBytes(nil),
	},
	{
		Name:     "TypeTag",
		Value:    3,
		Template: TypeTagAddress([32]byte{}),
	},
	{
		Name:     "TypeTag",
		Value:    4,
		Template: (*TypeTagStructTag)(nil),
	},
}

type TypeTag struct {
	TypeTag isTypeTag `lcs:"enum=TypeTag"`
}

// EnumTypes defines enum variants for lcs
func (*TypeTag) EnumTypes() []lcs.EnumVariant { return typeTagEnumDef }

// StructTag is a tag to form a resource path.
//
// StructTag implements AccessPathTag interface
type StructTag struct {
	Address    AccountAddress
	Module     string
	Name       string
	TypeParams []TypeTag
}

// Hash outputs the hash of this struct, using the appropriate hash function.
func (t *StructTag) Hash() HashValue {
	hasher := sha3libra.NewStructTag()
	if err := lcs.NewEncoder(hasher).Encode(t); err != nil {
		panic(err)
	}
	return hasher.Sum([]byte{})
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
