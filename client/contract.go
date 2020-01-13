package client

import (
	"github.com/the729/go-libra/types"
	"time"
)

// NewRawCustomTransaction creates a new serialized raw transaction bytes corresponding to a
// custom transaction.
func NewRawCustomTransaction(
	senderAddress types.AccountAddress,
	senderSequenceNumber uint64,
	maxGasAmount, gasUnitPrice uint64,
	expiration time.Time,
	payload types.TransactionPayload,
) (*types.RawTransaction, error) {
	txn := &types.RawTransaction{
		Sender:         senderAddress,
		SequenceNumber: senderSequenceNumber,
		Payload:        payload,
		MaxGasAmount:   maxGasAmount,
		GasUnitPrice:   gasUnitPrice,
		ExpirationTime: uint64(expiration.Unix()),
	}

	return txn, nil
}
