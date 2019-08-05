package client

import (
	"encoding/hex"

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
