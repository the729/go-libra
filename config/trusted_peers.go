package config

import (
	"fmt"

	"github.com/BurntSushi/toml"

	lcrypto "github.com/the729/go-libra/crypto"
)

type TrustedPeer struct {
	NetworkSigningPubkey  lcrypto.PublicKey `toml:"ns"`
	NetworkIdentityPubkey lcrypto.PublicKey `toml:"ni"`
	ConsensusPubkey       lcrypto.PublicKey `toml:"c"`
}

type TrustedPeersConfig struct {
	Peers map[string]TrustedPeer
}

func LoadTrustedPeersFromFile(fn string) (*TrustedPeersConfig, error) {
	peersConf := &TrustedPeersConfig{}
	peersConf.Peers = make(map[string]TrustedPeer)

	if _, err := toml.DecodeFile(fn, &peersConf.Peers); err != nil {
		return nil, fmt.Errorf("TrustedPeersConfig toml decode error: %v", err)
	}

	return peersConf, nil
}
