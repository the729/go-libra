package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ed25519"

	"github.com/the729/go-libra/generated/pbac"
	"github.com/the729/go-libra/types"
)

// SubmitRawTransaction signes and submits a raw transaction.
func (c *Client) SubmitRawTransaction(rawTxn []byte, privateKey ed25519.PrivateKey) error {
	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	signedTxn := types.SignRawTransaction(rawTxn, privateKey)
	pbSignedTxn, _ := signedTxn.ToProto()
	resp, err := c.ac.SubmitTransaction(ctx1, &pbac.SubmitTransactionRequest{
		SignedTxn: pbSignedTxn,
	})
	if err != nil {
		return fmt.Errorf("submit transaction error: %v", err)
	}

	// log.Printf("Result: ")
	// spew.Dump(resp)
	if vmStatus := resp.GetVmStatus(); vmStatus != nil {
		return fmt.Errorf("vm error: %s", vmStatus)
	}
	if mpStatus := resp.GetMempoolStatus(); mpStatus != nil {
		return fmt.Errorf("mempool error: %s", mpStatus)
	}
	if acStatus := resp.GetAcStatus(); acStatus.Code != pbac.AdmissionControlStatusCode_Accepted {
		return fmt.Errorf("ac error: %s", acStatus)
	}

	return nil
}

// PollSequenceUntil blocks to repeatedly poll the sequence number of a specific account, until the sequence number
// is greater or equal to specified target sequence number, or the ledger state passes specified expiration time.
func (c *Client) PollSequenceUntil(addr types.AccountAddress, targetSeq uint64, expiration time.Time) error {
	for range time.Tick(1 * time.Second) {
		paccount, err := c.QueryAccountState(addr)
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
