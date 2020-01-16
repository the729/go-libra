package client

import (
	"github.com/the729/go-libra/types"
	"time"
)

// NewRawCustomModuleTransaction creates a new serialized raw transaction bytes corresponding to a
// custom module transaction.
func NewRawCustomModuleTransaction(
	senderAddress types.AccountAddress,
	senderSequenceNumber uint64,
	maxGasAmount, gasUnitPrice uint64,
	expiration time.Time,
	module types.TxnPayloadModule,
) (*types.RawTransaction, error) {
	txn := &types.RawTransaction{
		Sender:         senderAddress,
		SequenceNumber: senderSequenceNumber,
		Payload:        module,
		MaxGasAmount:   maxGasAmount,
		GasUnitPrice:   gasUnitPrice,
		ExpirationTime: uint64(expiration.Unix()),
	}

	return txn, nil
}
