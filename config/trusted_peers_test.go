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
				ConsensusPubkey: lcrypto.PublicKey(mustDecodeHex("576e91b04632683a11c3be3dc47a19f9f0a31ae947211f59c5fe02dfa2d07d68")),
			},
			"dfb9c683d1788857e961160f28d4c9c79b23f042c80f770f37f0f93ee5fa6a96": {
				ConsensusPubkey: lcrypto.PublicKey(mustDecodeHex("3ca1400fb865befa8a21c58e90fc636ef2f84993a8396cb0e10008f876a00afd")),
			},
		},
	}, peers)
}
