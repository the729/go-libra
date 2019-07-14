package types

import (
	"encoding/hex"
	"io"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
)

const (
	AccountAddressLength = sha3libra.HashSize
)

type AccountAddress []byte

func (a AccountAddress) SerializeTo(w io.Writer) error {
	if err := serialization.SimpleSerializer.Write(w, []byte(a)); err != nil {
		return err
	}
	return nil
}

func (a AccountAddress) Hash() sha3libra.HashValue {
	hasher := sha3libra.NewAccountAddress()
	hasher.Write(a)
	return hasher.Sum([]byte{})
}

func (a *AccountAddress) UnmarshalText(txt []byte) error {
	data, err := hex.DecodeString(string(txt))
	if err != nil {
		return ErrInvalidText
	}
	*a = data
	return nil
}

func (a AccountAddress) MarshalText() (text []byte, err error) {
	return []byte(hex.EncodeToString(a)), nil
}
