package types

type BlockMetaData struct {
	ID                 HashValue
	TimestampUSec      uint64
	PreviousBlockVotes map[AccountAddress][]byte
	Proposer           AccountAddress
}

// Clone deep clones this struct.
func (bm *BlockMetaData) Clone() *BlockMetaData {
	out := &BlockMetaData{
		ID:                 cloneBytes(bm.ID),
		TimestampUSec:      bm.TimestampUSec,
		PreviousBlockVotes: make(map[AccountAddress][]byte),
		Proposer:           bm.Proposer,
	}
	for k, v := range bm.PreviousBlockVotes {
		out.PreviousBlockVotes[k] = cloneBytes(v)
	}
	return out
}
