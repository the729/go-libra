# go-libra

[![Go Report Card](https://goreportcard.com/badge/github.com/the729/go-libra)](https://goreportcard.com/report/github.com/the729/go-libra)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/15abfbbb81354b7fae9656baa6204002)](https://www.codacy.com/app/the729/go-libra?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=the729/go-libra&amp;utm_campaign=Badge_Grade)

A golang client library for [Libra blockchain](https://github.com/libra/libra). 

It has all cryptographic verification algorithms, including validator-signature-based consensus verification, ledger history accumulator proof, and account state sparse Merkle tree proof, etc. 

## Features

Compatible with testnet 2019/09/03 (commit hash [05364e6da2](https://github.com/libra/libra/commit/05364e6da23a58948c1243c09023df4c9be74783)).

- Data models with all necessary cryptographic verification algorithms
  - Ledger state: signature-based consensus verification
  - Transaction info: ledger history accumulator proof
  - Transaction list: ledger history accumulator proof on a range of transactions
  - Transaction signature: ed25519 signature
  - Account state: sparse Merkle tree proof
  - Events: event list hash based on Merkle tree accumulator
- RPC functions
  - Query account states
  - Make P2P transaction, and wait for ledger inclusion
  - Query transactions by ledger version
  - Query account transaction by sequence number

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
	"log"

	"github.com/the729/go-libra/client"
)

const (
	defaultServer    = "ac.testnet.libra.org:8000"
	trustedPeersFile = "../trusted_peers.config.toml"
)

func main() {
	c, err := client.New(defaultServer, trustedPeersFile)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	addrStr := "18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a"
	// Parse hex string into binary address
	addr := client.MustToAddress(addrStr)

	// provenState is cryptographically proven state of account
	provenState, err := c.QueryAccountState(addr)
	if err != nil {
		log.Fatal(err)
	}

	if provenState.IsNil() {
		log.Printf("Account %s does not exist at version %d.", addrStr, provenState.GetVersion())
		return
	}

	provenResource, err := c.GetLibraCoinResourceFromAccountBlob(provenState.GetAccountBlob())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Balance (microLibra): %d", provenResource.GetBalance())
	log.Printf("Sequence Number: %d", provenResource.GetSequenceNumber())
}
```

### Make peer-to-peer transaction

```golang
// Get current account sequence of sender
seq, err := c.GetAccountSequenceNumber(senderAddr)
if err != nil {
	log.Fatal(err)
}

// Build a raw transaction
rawTxn, err := types.NewRawP2PTransaction(
	senderAddr, recvAddr, seq,
	amountMicro, maxGasAmount, gasUnitPrice, expiration,
)
if err != nil {
	log.Fatal(err)
}

// Sign and submit transaction
err = c.SubmitRawTransaction(rawTxn, priKey)
if err != nil {
	log.Fatal(err)
}

// Wait until transaction is included in ledger, or timeout
err = c.PollSequenceUntil(senderAddr, seq+1, expiration)
if err != nil {
	log.Fatal(err)
}
```

### Other examples

Several examples are included in `example` folder.
- [cli_client](example/cli_client): A fully functional Libra CLI client
- [query_account](example/query_account): Query specific account states
- [query_txn_range](example/query_txn_range): Query a range of transactions
- [query_txn_by_seq](example/query_txn_by_seq): Query a transaction by specific account and sequence number
- [p2p_transaction](example/p2p_transaction): Make P2P transaction

## Contributions are welcome
