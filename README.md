# go-libra

[![Build Status](https://travis-ci.org/the729/go-libra.svg?branch=master)](https://travis-ci.org/the729/go-libra)
[![codecov](https://codecov.io/gh/the729/go-libra/branch/master/graph/badge.svg)](https://codecov.io/gh/the729/go-libra)
[![Go Report Card](https://goreportcard.com/badge/github.com/the729/go-libra)](https://goreportcard.com/report/github.com/the729/go-libra)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/15abfbbb81354b7fae9656baa6204002)](https://www.codacy.com/app/the729/go-libra?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=the729/go-libra&amp;utm_campaign=Badge_Grade)

A golang client library for [Libra blockchain](https://github.com/libra/libra). 

Thanks to [gopherjs](https://gopherjs.github.io/), go-libra is also available as a Javascript client library: [gopherjs-libra](https://www.npmjs.com/package/gopherjs-libra). It works for both NodeJS and browsers. 

It has all cryptographic verification algorithms, including validator-signature-based consensus verification, ledger history accumulator proof, and account state sparse Merkle tree proof, etc. 

## Features

Compatible with testnet 2020/4/8 (commit hash [718ace82](https://github.com/libra/libra/commit/718ace82250e7bd64e08d7d61951bfaa8cee9ea4)).

- Data models with all necessary cryptographic verification algorithms
  - Ledger state: signature-based consensus verification
  - Ledger consistency verification: detects reset or hard-forks
  - Transaction info and event: Merkle tree accumulator proof
  - Transaction list: Merkle tree accumulator proof on a range of transactions
  - Transaction signature: ed25519 signature
  - Account state: sparse Merkle tree proof
  - Events: event list hash based on Merkle tree accumulator
- RPC functions
  - Query account states
  - Make P2P transaction, and wait for ledger inclusion
  - Query transaction list by ledger version
  - Query account transaction by sequence number
  - Query sent or received event list by account

## Installation

```bash
$ # download the code
$ go get -u github.com/the729/go-libra

$ # build example client
$ cd example/cli_client && go build

$ # see example/cli_client/README.md
$ ./cli_client a ls
```

## Usage

Godoc reference to [client package](https://godoc.org/github.com/the729/go-libra/client) and  [types package](https://godoc.org/github.com/the729/go-libra/types).

### Get account balance, cryptographically proven

```golang
package main

import (
	"context"
	"log"

	"github.com/the729/go-libra/client"
)

const (
	defaultServer = "ac.testnet.libra.org:8000"
	waypoint      = "0:4d4d0feaa9378069f8fcee71980e142273837e108702d8d7f93a8419e2736f3f"
)

func main() {
	c, err := client.New(defaultServer, waypoint)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	addrStr := "42f5745128c05452a0c68272de8042b1"
	addr := client.MustToAddress(addrStr)

	// provenState is cryptographically proven state of account
	provenState, err := c.QueryAccountState(context.TODO(), addr)
	if err != nil {
		log.Fatal(err)
	}

	if provenState.IsNil() {
		log.Printf("Account %s does not exist at version %d.", addrStr, provenState.GetVersion())
		return
	}

	ar, br, err := provenState.GetAccountBlob().GetLibraResources()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Balance: %d", br.Coin)
	log.Printf("Sequence Number: %d", ar.SequenceNumber)
	log.Printf("Authentication key: %v", hex.EncodeToString(ar.AuthenticationKey))
}
```

### Make peer-to-peer transaction

```golang
log.Printf("Get current account sequence of sender...")
seq, err := c.QueryAccountSequenceNumber(context.TODO(), senderAddr)
if err != nil {
	log.Fatal(err)
}
log.Printf("... is %d", seq)

rawTxn, err := client.NewRawP2PTransaction(
	senderAddr, recvAddr, recvAuthKeyPrefix, seq,
	amountMicro, maxGasAmount, gasUnitPrice, expiration,
)
if err != nil {
	log.Fatal(err)
}

log.Printf("Submit transaction...")
expectedSeq, err := c.SubmitRawTransaction(context.TODO(), rawTxn, priKey)
if err != nil {
	log.Fatal(err)
}

log.Printf("Waiting until transaction is included in ledger...")
err = c.PollSequenceUntil(context.TODO(), senderAddr, expectedSeq, expiration)
if err != nil {
	log.Fatal(err)
}
log.Printf("done.")
```

### Other examples

Several examples are included in `example` folder.
- [cli_client](example/cli_client): A fully functional Libra CLI client
- [query_account](example/query_account): Query specific account states
- [query_txn_range](example/query_txn_range): Query a range of transactions
- [query_txn_by_seq](example/query_txn_by_seq): Query a transaction by specific account and sequence number
- [p2p_transaction](example/p2p_transaction): Make P2P transaction
- [NodeJS examples](example/nodejs) based on gopherjs-libra
- [web_client](example/web_client): a pure front-end libra explorer based on gopherjs-libra. [See it in action](http://pg.wutj.info/web_client/)

## Contributions are welcome
