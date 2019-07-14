package client

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ed25519"

	"github.com/the729/go-libra/generated/pbac"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types"
)

func (c *Client) SubmitTransactionRequest(signedTxn *pbtypes.SignedTransaction) (*pbac.SubmitTransactionResponse, error) {
	req := &pbac.SubmitTransactionRequest{
		SignedTxn: signedTxn,
	}

	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.ac.SubmitTransaction(ctx1, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) CmdTransfer(ctx *cli.Context) error {
	c.Connect()
	defer c.Disconnect()
	c.LoadTrustedPeers()
	c.LoadAccounts()

	sender, err := c.GetAccount(ctx.Args().Get(0))
	if err != nil {
		return err
	}
	receiver, err := c.GetAccount(ctx.Args().Get(1))
	if err != nil {
		return err
	}

	amount, err := strconv.Atoi(ctx.Args().Get(2))
	if err != nil {
		return err
	}
	amountMicro := uint64(amount) * 1000000

	log.Printf("Going to transfer %d microLibra from %s to %s", amountMicro, hex.EncodeToString(sender.Address), hex.EncodeToString(receiver.Address))

	seq, _, err := c.GetAccountSequenceNumber(sender.Address)
	if err != nil {
		return errors.New("sender account not present in ledger. ")
	}
	sender.SequenceNumber = seq
	log.Printf("Refreshed sequence number of sender: %d", sender.SequenceNumber)

	expiration := time.Now().Add(1 * time.Minute)
	maxGasAmount := uint64(10000) // 1000 is too little
	gasUnitPrice := uint64(0)
	rawTxn, err := types.NewRawTransaction(
		sender.Address, receiver.Address, sender.SequenceNumber,
		amountMicro, maxGasAmount, gasUnitPrice, expiration,
	)
	if err != nil {
		return fmt.Errorf("cannot create raw transaction: %v", err)
	}

	signedTxn := types.SignRawTransaction(rawTxn, ed25519.PrivateKey(sender.PrivateKey))
	pbSignedTxn, _ := signedTxn.ToProto()
	resp, err := c.SubmitTransactionRequest(pbSignedTxn)
	if err != nil {
		return fmt.Errorf("submit transaction error: %v", err)
	}

	log.Printf("Result: ")
	spew.Dump(resp)
	if vmStatus := resp.GetVmStatus(); vmStatus != nil {
		log.Printf("VM Error.")
		return nil
	}
	if mpStatus := resp.GetMempoolStatus(); mpStatus != nil {
		log.Printf("Mempool Error: ")
		log.Printf("         Code: %d", mpStatus.Code)
		log.Printf("      Message: %s", mpStatus.Message)
		return nil
	}
	if acStatus := resp.GetAcStatus(); acStatus.Code != pbac.AdmissionControlStatusCode_Accepted {
		log.Printf("AC Error: ")
		log.Printf("    Code: %d", acStatus.Code)
		log.Printf(" Message: %s", acStatus.Message)
		return nil
	}

	log.Printf("Waiting until transaction is included in ledger...")
	for range time.Tick(1 * time.Second) {
		seq, ledgerInfo, _ := c.GetAccountSequenceNumber(sender.Address)
		log.Printf("sequence number of sender: %d (ledger version: %d)", seq, ledgerInfo.Version)
		if seq >= sender.SequenceNumber+1 {
			break
		}
		if ledgerInfo.TimestampUsec > uint64(expiration.Unix()+1)*1000000 {
			log.Printf("transaction expired. ")
			break
		}
	}

	return nil
}
