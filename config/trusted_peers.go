package config

import (
	"fmt"

	"github.com/BurntSushi/toml"

	lcrypto "github.com/the729/go-libra/crypto"
)

type TrustedPeer struct {
	NetworkSigningPubkey  lcrypto.PublicKey `toml:"network_signing_pubkey"`
	NetworkIdentityPubkey lcrypto.PublicKey `toml:"network_identity_pubkey"`
	ConsensusPubkey       lcrypto.PublicKey `toml:"consensus_pubkey"`
}

type TrustedPeersConfig struct {
	Peers map[string]TrustedPeer `toml:"peers"`
}

func LoadTrustedPeersFromFile(fn string) (*TrustedPeersConfig, error) {
	peersConf := &TrustedPeersConfig{}

	if _, err := toml.DecodeFile(fn, peersConf); err != nil {
		return nil, fmt.Errorf("TrustedPeersConfig toml decode error: %v", err)
	}

	return peersConf, nil
}
