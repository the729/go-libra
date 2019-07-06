package config

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestLoadTrustedPeersConfigFromFile(t *testing.T) {
	peers, err := LoadTrustedPeersFromFile("trusted_peers.config.toml")
	assert.NoError(t, err)
	spew.Dump(peers)
}
