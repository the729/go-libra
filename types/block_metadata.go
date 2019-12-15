package types

type BlockMetaData struct {
	ID                 HashValue
	TimestampUSec      uint64
	PreviousBlockVotes map[AccountAddress][]byte
	Proposer           AccountAddress
}
