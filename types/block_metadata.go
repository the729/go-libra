package types

import "github.com/the729/go-libra/crypto/sha3libra"

type BlockMetaData struct {
	ID                 sha3libra.HashValue
	TimestampUSec      uint64
	PreviousBlockVotes map[string][]byte
	Proposer           AccountAddress
}
