package types

import (
	"encoding/hex"

	"github.com/the729/go-libra/crypto/sha3libra"
)

const (
	// AccountAddressLength is the length of an account address, which is 16 bytes.
	AccountAddressLength = 16
)

// AccountAddress is an account address.
type AccountAddress [AccountAddressLength]byte

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (a AccountAddress) Hash() HashValue {
	hasher := sha3libra.NewAccountAddress()
	hasher.Write(a[:])
	return hasher.Sum([]byte{})
}

// UnmarshalText unmarshals the hex representation of an account address.
func (a *AccountAddress) UnmarshalText(txt []byte) error {
	data, err := hex.DecodeString(string(txt))
	if err != nil {
		return err
	}
	copy(a[:], data)
	return nil
}

// MarshalText marshals the account address into hex representation.
func (a AccountAddress) MarshalText() (text []byte, err error) {
	return []byte(hex.EncodeToString(a[:])), nil
}
