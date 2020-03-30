package types

import (
	"github.com/the729/go-libra/crypto"
	"github.com/the729/go-libra/generated/pbtypes"
)

// ValidatorInfo is the set of public keys and info of a validator
type ValidatorInfo struct {
	AccountAddress        AccountAddress   `toml:"addr" json:"addr"`
	ConsensusPubkey       crypto.PublicKey `toml:"c" json:"c"`
	ConsensusVotingPower  uint64           `toml:"power" json:"power"`
	NetworkSigningPubkey  crypto.PublicKey `toml:"ns" json:"ns"`
	NetworkIdentityPubkey crypto.PublicKey `toml:"ni" json:"ni"`
}

// ValidatorSet is a set of validators
type ValidatorSet []*ValidatorInfo

// FromProto parses a protobuf struct into this struct.
func (vi *ValidatorInfo) FromProto(pb *pbtypes.ValidatorInfo) error {
	copy(vi.AccountAddress[:], pb.AccountAddress)
	vi.ConsensusPubkey = pb.ConsensusPublicKey
	vi.ConsensusVotingPower = pb.ConsensusVotingPower
	vi.NetworkSigningPubkey = pb.NetworkSigningPublicKey
	vi.NetworkIdentityPubkey = pb.NetworkIdentityPublicKey
	return nil
}

// Clone deep clones this struct.
func (vi *ValidatorInfo) Clone() *ValidatorInfo {
	out := &ValidatorInfo{
		AccountAddress:        vi.AccountAddress,
		ConsensusPubkey:       cloneBytes(vi.ConsensusPubkey),
		ConsensusVotingPower:  vi.ConsensusVotingPower,
		NetworkSigningPubkey:  cloneBytes(vi.NetworkSigningPubkey),
		NetworkIdentityPubkey: cloneBytes(vi.NetworkIdentityPubkey),
	}
	return out
}

// FromProto parses a protobuf struct into this struct.
func (vs *ValidatorSet) FromProto(pb *pbtypes.ValidatorSet) error {
	*vs = nil
	for _, v := range pb.ValidatorInfo {
		v1 := &ValidatorInfo{}
		if err := v1.FromProto(v); err != nil {
			return err
		}
		*vs = append(*vs, v1)
	}
	return nil
}
