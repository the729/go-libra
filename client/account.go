package client

import (
	"encoding/hex"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"

	"github.com/the729/go-libra/types"
)

// MustToAddress converts hex string represent of an address into types.AccountAddress.
// Input string should be a hex string with exactly 64 hex digits.
func MustToAddress(str string) types.AccountAddress {
	addr, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	if len(addr) != types.AccountAddressLength {
		panic("wrong address length")
	}
	return types.AccountAddress(addr)
}

// PubkeyMustToAddress converts an ed25519 public key (32 bytes) into types.AccountAddress.
func PubkeyMustToAddress(pubkey []byte) types.AccountAddress {
	if len(pubkey) != ed25519.PublicKeySize {
		panic("wrong pubkey length")
	}
	hasher := sha3.New256()
	hasher.Write(pubkey)
	return hasher.Sum([]byte{})
}
