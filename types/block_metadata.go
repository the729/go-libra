package types

type BlockMetaData struct {
	ID                 HashValue
	TimestampUSec      uint64
	PreviousBlockVotes map[string][]byte
	Proposer           AccountAddress `lcs:"len=32"`
}
