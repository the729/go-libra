package validator

import (
	"errors"

	"github.com/the729/go-libra/config"
	"golang.org/x/crypto/ed25519"
)

var (
	VerifyErrUnknownAuthor     = errors.New("unknown author")
	VerifyErrInvalidSignature  = errors.New("invalid signature")
	VerifyErrTooFewSignatures  = errors.New("too few signatures")
	VerifyErrTooManySignatures = errors.New("too many signatures")
)

type Verifier interface {
	Verify(hash []byte, sigs map[string][]byte) error
}

type verifier struct {
	publicKeyMap map[string]ed25519.PublicKey
	quorum       int
}

func NewConsensusVerifier(conf *config.TrustedPeersConfig) (Verifier, error) {
	v := &verifier{
		publicKeyMap: make(map[string]ed25519.PublicKey),
	}
	for accountAddr, peer := range conf.Peers {
		if peer.ConsensusPubkey != nil {
			v.publicKeyMap[accountAddr] = ed25519.PublicKey(peer.ConsensusPubkey)
		}
	}
	// Total 3f + 1 validators, 2f + 1 correct signatures are required.
	// If < 4 validators, all validators have to agree.
	v.quorum = len(v.publicKeyMap)*2/3 + 1
	return v, nil
}

func (v *verifier) verifySingle(author string, hash, sig []byte) error {
	pubk, ok := v.publicKeyMap[author]
	if !ok {
		return VerifyErrUnknownAuthor
	}
	ok = ed25519.Verify(pubk, hash, sig)
	if !ok {
		return VerifyErrInvalidSignature
	}
	return nil
}

func (v *verifier) Verify(hash []byte, sigs map[string][]byte) error {
	if len(sigs) < v.quorum {
		return VerifyErrTooFewSignatures
	}
	if len(sigs) > len(v.publicKeyMap) {
		return VerifyErrTooManySignatures
	}
	for author, sig := range sigs {
		err := v.verifySingle(author, hash, sig)
		if err != nil {
			return err
		}
	}
	return nil
}
