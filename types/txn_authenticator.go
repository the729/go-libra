package types

import "github.com/the729/lcs"

type TxnAuthenticator interface {
	Clone() TxnAuthenticator
}

type ED25519Authenticator struct {
	// PublicKey is the public key of the sender.
	PublicKey []byte

	// Signature is the signature.
	Signature []byte
}

// type MultiED25519Authenticator struct {

// }

// Clone the TxnAuthenticator
func (v *ED25519Authenticator) Clone() TxnAuthenticator {
	out := &ED25519Authenticator{}
	out.PublicKey = cloneBytes(v.PublicKey)
	out.Signature = cloneBytes(v.Signature)
	return out
}

var txnAuthenticatorEnumDef = []lcs.EnumVariant{
	{
		Name:     "TxnAuthenticator",
		Value:    0,
		Template: (*ED25519Authenticator)(nil),
	},
}
