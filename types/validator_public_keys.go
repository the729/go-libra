package types

import (
	"golang.org/x/crypto/ed25519"

	"github.com/the729/go-libra/generated/pbtypes"
)

type ValidatorPublicKeys struct {
	AccountAddress        AccountAddress
	ConsensusPubkey       ed25519.PublicKey
	NetworkSigningPubkey  ed25519.PublicKey
	NetworkIdentityPubkey ed25519.PublicKey
}

type ValidatorSet []*ValidatorPublicKeys

func (vk *ValidatorPublicKeys) FromProto(pb *pbtypes.ValidatorPublicKeys) error {
	vk.AccountAddress = pb.AccountAddress
	vk.ConsensusPubkey = pb.ConsensusPublicKey
	vk.NetworkSigningPubkey = pb.NetworkSigningPublicKey
	vk.NetworkIdentityPubkey = pb.NetworkIdentityPublicKey
	return nil
}

func (vs ValidatorSet) FromProto(pb *pbtypes.ValidatorSet) error {
	for _, v := range pb.ValidatorPublicKeys {
		v1 := &ValidatorPublicKeys{}
		if err := v1.FromProto(v); err != nil {
			return err
		}
		vs = append(vs, v1)
	}
	return nil
}
