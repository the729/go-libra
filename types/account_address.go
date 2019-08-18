package types

import (
	"encoding/hex"
	"io"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
)

const (
	// AccountAddressLength is the length of an account address, which is 32 bytes.
	AccountAddressLength = sha3libra.HashSize
)

// AccountAddress is an account address.
type AccountAddress []byte

// SerializeTo serializes this struct into a io.Writer.
func (a AccountAddress) SerializeTo(w io.Writer) error {
	if err := serialization.SimpleSerializer.Write(w, []byte(a)); err != nil {
		return err
	}
	return nil
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (a AccountAddress) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewAccountAddress()
	hasher.Write(a)
	return hasher.Sum([]byte{})
}

// UnmarshalText unmarshals the hex representation of an account address.
func (a *AccountAddress) UnmarshalText(txt []byte) error {
	data, err := hex.DecodeString(string(txt))
	if err != nil {
		return err
	}
	*a = data
	return nil
}

// MarshalText marshals the account address into hex representation.
func (a AccountAddress) MarshalText() (text []byte, err error) {
	return []byte(hex.EncodeToString(a)), nil
}
