package client

import (
	"github.com/the729/go-libra/types"
	"time"
)

// NewRawCustomScriptTransaction creates a new serialized raw transaction bytes corresponding to a
// custom script transaction.
func NewRawCustomScriptTransaction(
	senderAddress types.AccountAddress,
	senderSequenceNumber uint64,
	maxGasAmount, gasUnitPrice uint64,
	expiration time.Time,
	code []byte,
	args []types.TransactionArgument,
) (*types.RawTransaction, error) {
	txn := &types.RawTransaction{
		Sender:         senderAddress,
		SequenceNumber: senderSequenceNumber,
		Payload: &types.TxnPayloadScript{
			Code: code,
			Args: args,
		},
		MaxGasAmount:   maxGasAmount,
		GasUnitPrice:   gasUnitPrice,
		ExpirationTime: uint64(expiration.Unix()),
	}

	return txn, nil
}
