package sha3libra

var (
	AccumulatorPlaceholderHash  HashValue
	SparseMerklePlaceholderHash HashValue
	PreGenesisBlockID           HashValue
	GenesisBlockID              HashValue
)

func init() {
	AccumulatorPlaceholderHash = make([]byte, HashSize)
	SparseMerklePlaceholderHash = make([]byte, HashSize)
	PreGenesisBlockID = make([]byte, HashSize)
	GenesisBlockID = make([]byte, HashSize)

	copy(AccumulatorPlaceholderHash, []byte("ACCUMULATOR_PLACEHOLDER_HASH"))
	copy(SparseMerklePlaceholderHash, []byte("SPARSE_MERKLE_PLACEHOLDER_HASH"))
	copy(PreGenesisBlockID, []byte("PRE_GENESIS_BLOCK_ID"))
	copy(GenesisBlockID, []byte("GENESIS_BLOCK_ID"))
}
