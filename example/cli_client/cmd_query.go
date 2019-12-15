package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"

	"github.com/urfave/cli"

	"github.com/the729/go-libra/client"
	"github.com/the729/go-libra/example/utils"
	"github.com/the729/go-libra/language/stdscript"
	"github.com/the729/go-libra/types"
	"github.com/the729/lcs"
)

func cmdQueryLedgerInfo(ctx *cli.Context) error {
	c, err := client.New(ServerAddr, TrustedPeersFile)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	loadKnownVersion(c, KnownVersionFile)
	defer saveKnownVersion(c, KnownVersionFile)

	ledgerInfo, err := c.QueryLedgerInfo(context.Background())
	if err != nil {
		return err
	}

	log.Printf("Ledger info: version %d, time %d", ledgerInfo.GetVersion(), ledgerInfo.GetTimestampUsec())
	return nil
}

func cmdQueryAccountState(ctx *cli.Context) error {
	c, err := client.New(ServerAddr, TrustedPeersFile)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	loadKnownVersion(c, KnownVersionFile)
	defer saveKnownVersion(c, KnownVersionFile)

	wallet, err := LoadAccounts(WalletFile)
	if err != nil {
		log.Fatal(err)
	}

	addr := ctx.Args().Get(0)
	account, err := wallet.GetAccount(addr)
	if err != nil {
		return err
	}
	accountState, err := c.QueryAccountState(context.Background(), account.Address)
	if err != nil {
		return err
	}

	ledgerInfo := accountState.GetLedgerInfo()
	log.Printf("Ledger info: version %d, time %d", ledgerInfo.GetVersion(), ledgerInfo.GetTimestampUsec())

	if accountState.IsNil() {
		fmt.Println("Account is not present in the ledger.")
	} else {
		log.Printf("Account version: %d", accountState.GetVersion())
		resource, err := c.GetLibraCoinResourceFromAccountBlob(accountState.GetAccountBlob())
		if err != nil {
			log.Printf("Account does not contain libra coin resource, err: %v", err)
			return nil
		}

		addr := resource.GetAddress()
		log.Printf("Address: %v", hex.EncodeToString(addr[:]))
		log.Printf("Balance (microLibra): %d", resource.GetBalance())
		log.Printf("Sequence Number: %d", resource.GetSequenceNumber())
		log.Printf("SentEventsCount: %d", resource.GetSentEvents().Count)
		log.Printf("    Key: %x", resource.GetSentEvents().Key)
		log.Printf("ReceivedEventsCount: %d", resource.GetReceivedEvents().Count)
		log.Printf("    Key: %x", resource.GetReceivedEvents().Key)
		log.Printf("DelegatedWithdrawalCapability: %v", resource.GetDelegatedWithdrawalCapability())
		log.Printf("Authentication key: %v", hex.EncodeToString(resource.GetAuthenticationKey()))
		log.Printf("Event generator: %v", resource.GetEventGenerator())
	}
	return nil
}

func cmdQueryTransactionRange(ctx *cli.Context) error {
	c, err := client.New(ServerAddr, TrustedPeersFile)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	loadKnownVersion(c, KnownVersionFile)
	defer saveKnownVersion(c, KnownVersionFile)

	start, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return err
	}
	limit, err := strconv.Atoi(ctx.Args().Get(1))
	if err != nil {
		return err
	}
	if limit > 100 {
		log.Printf("Limit>100, set to 100.")
		limit = 100
	}
	provenTxnList, err := c.QueryTransactionRange(context.Background(), uint64(start), uint64(limit), true)
	if err != nil {
		log.Fatal(err)
	}

	for _, provenTxn := range provenTxnList.GetTransactions() {
		utils.PrintTxn(provenTxn)
	}
	return nil
}

func cmdQueryTransactionByAccountSeq(ctx *cli.Context) error {
	c, err := client.New(ServerAddr, TrustedPeersFile)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	loadKnownVersion(c, KnownVersionFile)
	defer saveKnownVersion(c, KnownVersionFile)

	wallet, err := LoadAccounts(WalletFile)
	if err != nil {
		log.Fatal(err)
	}

	account, err := wallet.GetAccount(ctx.Args().Get(0))
	if err != nil {
		return err
	}

	sequence, err := strconv.Atoi(ctx.Args().Get(1))
	if err != nil {
		return err
	}

	provenTxn, err := c.QueryTransactionByAccountSeq(context.Background(), account.Address, uint64(sequence), true)
	if err != nil {
		log.Fatal(err)
	}

	utils.PrintTxn(provenTxn)
	return nil
}

func cmdQueryEvents(ctx *cli.Context) error {
	c, err := client.New(ServerAddr, TrustedPeersFile)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	loadKnownVersion(c, KnownVersionFile)
	defer saveKnownVersion(c, KnownVersionFile)

	wallet, err := LoadAccounts(WalletFile)
	if err != nil {
		log.Fatal(err)
	}

	account, err := wallet.GetAccount(ctx.Args().Get(0))
	if err != nil {
		return err
	}

	evType := ctx.Args().Get(1)

	start, err := strconv.Atoi(ctx.Args().Get(2))
	if err != nil {
		return err
	}

	ordering := ctx.Args().Get(3)

	limit, err := strconv.Atoi(ctx.Args().Get(4))
	if err != nil {
		return err
	}

	ap := &types.AccessPath{
		Address: account.Address,
	}
	switch evType {
	case "sent":
		ap.Path = types.AccountSentEventPath()
	case "received":
		ap.Path = types.AccountReceivedEventPath()
	default:
		return fmt.Errorf("unknown event type: %s, should be either sent or received", evType)
	}

	var ascending bool
	switch ordering {
	case "asc":
		ascending = true
	case "desc":
		ascending = false
	default:
		return fmt.Errorf("unknown ordering: %s, should be either asc or desc", evType)
	}
	evs, err := c.QueryEventsByAccessPath(context.Background(), ap, uint64(start), ascending, uint64(limit))
	if err != nil {
		log.Fatal(err)
	}

	for i, ev := range evs {
		log.Printf("#%d: txn #%d event #%d", i, ev.GetTransactionVersion(), ev.GetEventIndex())
		evBody := ev.GetEvent()
		log.Printf("    Key: %s", hex.EncodeToString(evBody.Key))
		log.Printf("    Seq number: %d", evBody.SequenceNumber)
		if len(evBody.Data) > 30 {
			log.Printf("    Raw event: %s ...", hex.EncodeToString(evBody.Data[:30]))
		} else {
			log.Printf("    Raw event: %s", hex.EncodeToString(evBody.Data))
		}

		pev := &stdscript.PaymentEvent{}
		if err := lcs.Unmarshal(evBody.Data, pev); err != nil {
			log.Printf("        (Unknown event type)")
		} else {
			log.Printf("        Amount (microLibra): %d", pev.Amount)
			log.Printf("        Opponent address: %s", hex.EncodeToString(pev.Address[:]))
		}
	}

	return nil
}
