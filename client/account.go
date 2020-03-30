package client

import (
	"encoding/hex"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"

	"github.com/the729/go-libra/types"
)

// MustToAddress converts hex string represent of an address into types.AccountAddress.
// Input string should be a hex string with exactly 64 hex digits.
func MustToAddress(str string) (out types.AccountAddress) {
	addr, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	if len(addr) != types.AccountAddressLength {
		panic("wrong address length")
	}
	copy(out[:], addr)
	return
}

// PubkeyMustToAddress converts an ed25519 public key (32 bytes) into types.AccountAddress (16 bytes).
func PubkeyMustToAddress(pubkey []byte) (out types.AccountAddress) {
	if len(pubkey) != ed25519.PublicKeySize {
		panic("wrong pubkey length")
	}
	hasher := sha3.New256()
	hasher.Write(pubkey)
	keyHash := hasher.Sum([]byte{})
	copy(out[:], keyHash[hasher.Size()-types.AccountAddressLength:])
	return
}
