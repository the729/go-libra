package types

import (
	"github.com/the729/go-libra/crypto"
	"github.com/the729/go-libra/generated/pbtypes"
)

// ValidatorPublicKeys is the set of public keys of a validator
type ValidatorPublicKeys struct {
	AccountAddress        AccountAddress   `toml:"addr" json:"addr"`
	ConsensusPubkey       crypto.PublicKey `toml:"c" json:"c"`
	ConsensusVotingPower  uint64           `toml:"power" json:"power"`
	NetworkSigningPubkey  crypto.PublicKey `toml:"ns" json:"ns"`
	NetworkIdentityPubkey crypto.PublicKey `toml:"ni" json:"ni"`
}

// ValidatorSet is a set of validators
type ValidatorSet []*ValidatorPublicKeys

// FromProto parses a protobuf struct into this struct.
func (vk *ValidatorPublicKeys) FromProto(pb *pbtypes.ValidatorPublicKeys) error {
	copy(vk.AccountAddress[:], pb.AccountAddress)
	vk.ConsensusPubkey = pb.ConsensusPublicKey
	vk.ConsensusVotingPower = pb.ConsensusVotingPower
	vk.NetworkSigningPubkey = pb.NetworkSigningPublicKey
	vk.NetworkIdentityPubkey = pb.NetworkIdentityPublicKey
	return nil
}

// Clone deep clones this struct.
func (vk *ValidatorPublicKeys) Clone() *ValidatorPublicKeys {
	out := &ValidatorPublicKeys{
		AccountAddress:        vk.AccountAddress,
		ConsensusPubkey:       cloneBytes(vk.ConsensusPubkey),
		ConsensusVotingPower:  vk.ConsensusVotingPower,
		NetworkSigningPubkey:  cloneBytes(vk.NetworkSigningPubkey),
		NetworkIdentityPubkey: cloneBytes(vk.NetworkIdentityPubkey),
	}
	return out
}

// FromProto parses a protobuf struct into this struct.
func (vs *ValidatorSet) FromProto(pb *pbtypes.ValidatorSet) error {
	*vs = nil
	for _, v := range pb.ValidatorPublicKeys {
		v1 := &ValidatorPublicKeys{}
		if err := v1.FromProto(v); err != nil {
			return err
		}
		*vs = append(*vs, v1)
	}
	return nil
}
