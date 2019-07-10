package client

import (
	"context"
	"crypto"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/ed25519"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbac"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/language/stdscript"
	"github.com/urfave/cli"
)

func (c *Client) NewRawTransaction(
	sender, receiver *Account,
	amount, maxGasAmount, gasUnitPrice uint64,
	expiration time.Time,
) ([]byte, error) {
	ammountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(ammountBytes, amount)

	txn := &pbtypes.RawTransaction{
		SenderAccount:  sender.Address,
		SequenceNumber: sender.SequenceNumber,
		Payload: &pbtypes.RawTransaction_Program{
			Program: &pbtypes.Program{
				Code: stdscript.PeerToPeerTransfer,
				Arguments: []*pbtypes.TransactionArgument{
					{
						Type: pbtypes.TransactionArgument_ADDRESS,
						Data: receiver.Address,
					},
					{
						Type: pbtypes.TransactionArgument_U64,
						Data: ammountBytes,
					},
				},
				Modules: nil,
			},
		},
		MaxGasAmount:   maxGasAmount,
		GasUnitPrice:   gasUnitPrice,
		ExpirationTime: uint64(expiration.Unix()),
	}

	j, _ := json.MarshalIndent(txn, "", "    ")
	log.Printf("Raw txn: %s", string(j))

	raw, err := proto.Marshal(txn)
	return raw, err
}

func (c *Client) SignRawTransaction(rawTxnBytes []byte, signer ed25519.PrivateKey) *pbtypes.SignedTransaction {
	hasher := sha3libra.NewRawTransaction()
	hasher.Write(rawTxnBytes)
	txnHash := hasher.Sum([]byte{})
	senderPubKey := signer.Public().(ed25519.PublicKey)
	sig, _ := signer.Sign(rand.Reader, txnHash, crypto.Hash(0))

	return &pbtypes.SignedTransaction{
		RawTxnBytes:     rawTxnBytes,
		SenderPublicKey: senderPubKey,
		SenderSignature: sig,
	}
}

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
	rawTxn, err := c.NewRawTransaction(sender, receiver, amountMicro, maxGasAmount, gasUnitPrice, expiration)
	if err != nil {
		return fmt.Errorf("cannot create raw transaction: %v", err)
	}

	signedTxn := c.SignRawTransaction(rawTxn, ed25519.PrivateKey(sender.PrivateKey))
	resp, err := c.SubmitTransactionRequest(signedTxn)
	if err != nil {
		return fmt.Errorf("submit transaction error: %v", err)
	}

	log.Printf("Result: ")
	spew.Dump(resp)
	if resp.GetVmStatus() != nil || resp.GetAcStatus() != pbac.AdmissionControlStatus_Accepted {
		log.Printf("Transaction failed. ")
		return nil
	}

	log.Printf("Waiting until transaction is included in ledger...")
	for range time.Tick(1 * time.Second) {
		seq, ledgerInfo, _ := c.GetAccountSequenceNumber(sender.Address)
		log.Printf("sequence number of sender: %d", seq)
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
