package types

import (
	"github.com/the729/go-libra/crypto"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/lcs"
)

// VSScheme is the enum type of ValidatorSet Scheme, currently only ED25519
type VSScheme interface {
	Clone() VSScheme
}

type SchemeED25519 struct{}

func (SchemeED25519) Clone() VSScheme { return SchemeED25519{} }

var vsSchemeEnumDef = []lcs.EnumVariant{
	{
		Name:     "ValidatorSetScheme",
		Value:    0,
		Template: SchemeED25519{},
	},
}

type ValidatorSet struct {
	Scheme  VSScheme `lcs:"enum=ValidatorSetScheme"`
	Payload []*ValidatorInfo
}

// EnumTypes defines enum variants for lcs
func (*ValidatorSet) EnumTypes() []lcs.EnumVariant { return vsSchemeEnumDef }

// ValidatorInfo is the set of public keys and info of a validator
type ValidatorInfo struct {
	AccountAddress        AccountAddress   `toml:"addr" json:"addr"`
	ConsensusPubkey       crypto.PublicKey `toml:"c" json:"c"`
	ConsensusVotingPower  uint64           `toml:"power" json:"power"`
	NetworkSigningPubkey  crypto.PublicKey `toml:"ns" json:"ns"`
	NetworkIdentityPubkey crypto.PublicKey `toml:"ni" json:"ni"`
}

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
	err := lcs.Unmarshal(pb.Bytes, vs)
	return err
}

// Clone deep clones this struct.
func (vs *ValidatorSet) Clone() *ValidatorSet {
	if vs == nil {
		return nil
	}
	out := &ValidatorSet{
		Scheme: vs.Scheme.Clone(),
	}
	pl := make([]*ValidatorInfo, 0, len(vs.Payload))
	for _, v := range vs.Payload {
		pl = append(pl, v.Clone())
	}
	out.Payload = pl
	return out
}
