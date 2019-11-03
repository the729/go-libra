package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ed25519"

	"github.com/the729/go-libra/generated/pbac"
	"github.com/the729/go-libra/language/stdscript"
	"github.com/the729/go-libra/types"
)

// NewRawP2PTransaction creates a new serialized raw transaction bytes corresponding to a
// peer-to-peer Libra coin transaction.
func NewRawP2PTransaction(
	senderAddress, receiverAddress types.AccountAddress,
	senderSequenceNumber uint64,
	amount, maxGasAmount, gasUnitPrice uint64,
	expiration time.Time,
) (*types.RawTransaction, error) {
	txn := &types.RawTransaction{
		Sender:         senderAddress,
		SequenceNumber: senderSequenceNumber,
		Payload: &types.TxnPayloadScript{
			Code: stdscript.PeerToPeerTransfer,
			Args: []types.TransactionArgument{
				types.TxnArgAddress{receiverAddress},
				types.TxnArgU64(amount),
			},
		},
		MaxGasAmount:   maxGasAmount,
		GasUnitPrice:   gasUnitPrice,
		ExpirationTime: uint64(expiration.Unix()),
	}

	return txn, nil
}

// SubmitRawTransaction signes and submits a raw transaction.
func (c *Client) SubmitRawTransaction(ctx context.Context, rawTxn *types.RawTransaction, privateKey ed25519.PrivateKey) (uint64, error) {
	signedTxn := types.SignRawTransaction(rawTxn, privateKey)
	pbSignedTxn, _ := signedTxn.ToProto()
	resp, err := c.ac.SubmitTransaction(ctx, &pbac.SubmitTransactionRequest{
		Transaction: pbSignedTxn,
	})
	if err != nil {
		return 0, fmt.Errorf("submit transaction error: %v", err)
	}

	// log.Printf("Result: ")
	// spew.Dump(resp)
	if vmStatus := resp.GetVmStatus(); vmStatus != nil {
		return 0, fmt.Errorf("vm error: %s", vmStatus)
	}
	if mpStatus := resp.GetMempoolStatus(); mpStatus != nil {
		return 0, fmt.Errorf("mempool error: %s", mpStatus)
	}
	if acStatus := resp.GetAcStatus(); acStatus.Code != pbac.AdmissionControlStatusCode_Accepted {
		return 0, fmt.Errorf("ac error: %s", acStatus)
	}

	return rawTxn.SequenceNumber + 1, nil
}

// PollSequenceUntil blocks to repeatedly poll the sequence number of a specific account, until the sequence number
// is greater or equal to specified target sequence number, or the ledger state passes specified expiration time.
func (c *Client) PollSequenceUntil(ctx context.Context, addr types.AccountAddress, targetSeq uint64, expiration time.Time) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		paccount, err := c.QueryAccountState(ctx, addr)
		if err != nil {
			return err
		}
		ledgerInfo := paccount.GetLedgerInfo()
		if !paccount.IsNil() {
			resource, err := c.GetLibraCoinResourceFromAccountBlob(paccount.GetAccountBlob())
			if err != nil {
				return err
			}
			seq := resource.GetSequenceNumber()
			log.Printf("sequence number: %d, ledger version: %d", seq, ledgerInfo.GetVersion())
			if seq >= targetSeq {
				return nil
			}
		}
		if ledgerInfo.GetTimestampUsec() > uint64(expiration.Unix()+1)*1000000 {
			break
		}
	}
	return errors.New("expired")
}
