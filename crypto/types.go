package crypto

import (
	"encoding/hex"

	"golang.org/x/crypto/ed25519"
)

type PublicKey ed25519.PublicKey
type PrivateKey ed25519.PrivateKey

func (k *PublicKey) UnmarshalText(txt []byte) error {
	data, err := hex.DecodeString(string(txt))
	if err != nil {
		return ErrInvalidText
	}
	if len(data) != ed25519.PublicKeySize {
		return ErrWrongSize
	}
	*k = data
	return nil
}

func (k PublicKey) MarshalText() (text []byte, err error) {
	return []byte(hex.EncodeToString(k)), nil
}

func (k *PrivateKey) UnmarshalText(txt []byte) error {
	data, err := hex.DecodeString(string(txt))
	if err != nil {
		return ErrInvalidText
	}
	if len(data) != ed25519.PrivateKeySize {
		return ErrWrongSize
	}
	*k = data
	return nil
}

func (k PrivateKey) MarshalText() (text []byte, err error) {
	return []byte(hex.EncodeToString(k)), nil
}
