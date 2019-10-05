package main

import (
	"context"
	"encoding/hex"
	"log"
	"strconv"
	"time"

	"github.com/urfave/cli"

	"github.com/the729/go-libra/client"
)

func cmdTransfer(ctx *cli.Context) error {
	c, err := client.New(ServerAddr, TrustedPeersFile)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	wallet, err := LoadAccounts(WalletFile)
	if err != nil {
		log.Fatal(err)
	}

	sender, err := wallet.GetAccount(ctx.Args().Get(0))
	if err != nil {
		return err
	}
	receiver, err := wallet.GetAccount(ctx.Args().Get(1))
	if err != nil {
		return err
	}

	amount, err := strconv.Atoi(ctx.Args().Get(2))
	if err != nil {
		return err
	}
	amountMicro := uint64(amount) * 1000000

	log.Printf("Going to transfer %d microLibra from %s to %s", amountMicro, hex.EncodeToString(sender.Address), hex.EncodeToString(receiver.Address))

	maxGasAmount := uint64(140000) // must > 70996?
	gasUnitPrice := uint64(0)
	expiration := time.Now().Add(1 * time.Minute)

	if ctx.Args().Get(3) != "" {
		if maxGasAmount, err = strconv.ParseUint(ctx.Args().Get(3), 10, 64); err != nil {
			return err
		}
	}
	if ctx.Args().Get(4) != "" {
		if gasUnitPrice, err = strconv.ParseUint(ctx.Args().Get(4), 10, 64); err != nil {
			return err
		}
	}
	if ctx.Args().Get(5) != "" {
		expSec, err := strconv.Atoi(ctx.Args().Get(5))
		if err != nil {
			return err
		}
		expiration = time.Now().Add(time.Duration(expSec) * time.Second)
	}

	log.Printf("Max gas: %d, Gas price: %d, Expiration: %v", maxGasAmount, gasUnitPrice, expiration)

	log.Printf("Get current account sequence of sender...")
	seq, err := c.GetAccountSequenceNumber(context.Background(), sender.Address)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("... is %d", seq)

	rawTxn, err := client.NewRawP2PTransaction(
		sender.Address, receiver.Address, seq,
		amountMicro, maxGasAmount, gasUnitPrice, expiration,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Submit transaction...")
	err = c.SubmitRawTransaction(context.Background(), rawTxn, sender.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Waiting until transaction is included in ledger...")
	err = c.PollSequenceUntil(context.Background(), sender.Address, seq+1, expiration)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("done.")
	return nil
}
