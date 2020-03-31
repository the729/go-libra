package types

type BlockMetaData struct {
	ID                 HashValue
	Round              uint64
	TimestampUSec      uint64
	PreviousBlockVotes []AccountAddress
	Proposer           AccountAddress
}

// Clone deep clones this struct.
func (bm *BlockMetaData) Clone() *BlockMetaData {
	out := &BlockMetaData{
		ID:                 cloneBytes(bm.ID),
		TimestampUSec:      bm.TimestampUSec,
		PreviousBlockVotes: make([]AccountAddress, 0, len(bm.PreviousBlockVotes)),
		Proposer:           bm.Proposer,
	}
	for _, v := range bm.PreviousBlockVotes {
		out.PreviousBlockVotes = append(out.PreviousBlockVotes, v)
	}
	return out
}
