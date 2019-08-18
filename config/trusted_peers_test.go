package config

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	lcrypto "github.com/the729/go-libra/crypto"
)

func mustDecodeHex(str string) []byte {
	b, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return b
}

func TestLoadTrustedPeersConfigFromFile(t *testing.T) {
	peers, err := LoadTrustedPeersFromFile("peers_test.config.toml")
	assert.NoError(t, err)
	assert.Equal(t, &TrustedPeersConfig{
		Peers: map[string]TrustedPeer{
			"9102bd7b1ad7e8f31023c500371cc7d2971758b450cfa89c003efb3ab192a4b8": {
				NetworkSigningPubkey:  lcrypto.PublicKey(mustDecodeHex("5f5ecda9576edd942ed22aa4735939092161445177cd456fd087c7bc1d6de403")),
				NetworkIdentityPubkey: lcrypto.PublicKey(mustDecodeHex("b5eb9a2e5814c66df6c01a1dc94252a4ae6733e93a58187c5eb48d1f53be0b28")),
				ConsensusPubkey:       lcrypto.PublicKey(mustDecodeHex("576e91b04632683a11c3be3dc47a19f9f0a31ae947211f59c5fe02dfa2d07d68")),
			},
			"dfb9c683d1788857e961160f28d4c9c79b23f042c80f770f37f0f93ee5fa6a96": {
				NetworkSigningPubkey:  lcrypto.PublicKey(mustDecodeHex("246ca919a3b39c95110e3bee891136ab087a9b3b9e84fa90cbf8f19c8abe62e3")),
				NetworkIdentityPubkey: lcrypto.PublicKey(mustDecodeHex("8aa297d686dd2444de86ea3a68353d74af74b9659990d06ccaf4344e2b629b33")),
				ConsensusPubkey:       lcrypto.PublicKey(mustDecodeHex("3ca1400fb865befa8a21c58e90fc636ef2f84993a8396cb0e10008f876a00afd")),
			},
		},
	}, peers)
}
