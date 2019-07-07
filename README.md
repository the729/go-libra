# go-libra
This is a golang implementation of the Libra blockchain (https://github.com/libra/libra). Currently, only a client is implemented. 

It is not a simple gRPC client. It has all the cryptographic verification algorithms implemented, including validator-signature-based consensus verification, ledger history accumulator proof, and account state sparse Merkle tree proof. 

As Libra itself is in an early stage, go-libra is developed as a proof of concept, only for experimental and learning purposes. 

## Features

- ✓ Connect to testnet AdmissionControl server with gRPC
- ✓ Data models with all necessary cryptographic verification algorithms
  - ✓ Ledger state: signature-based consensus verification
  - ✓ Transaction info: ledger history accumulator proof
  - ✓ Account state: sparse Merkle tree proof
- ✓ Query account states, including balance, sequence number, from the ledger, and do all necessary verifications
- ✓ Mint LibraCoin through the 'official' faucet service
- ✓ Peer to peer transfer LibraCoin, and wait for ledger inclusion
- X Query transactions from the ledger
- X Use mnemonics to manage wallet and private keys
- X Compile Move IR into bytecode

## Installation

1. Install protobuf and gRPC (instructions here: https://grpc.io/docs/quickstart/go/). 

2. Install go-libra:

```bash
$ # download the code
$ go get -u github.com/the729/go-libra

$ # build protobuf, see comments in gen_proto.sh
$ cd $GOPATH/src/github.com/the729/go-libra
$ scripts/gen_proto.sh

$ # build client
$ cd cmd/client && go build
```

This should build the client binary executable in cmd/client/client.

## Usage

The commands are similar to those of the official rust implementation. However, this is not an interactive CLI program, meaning that every time when you execute a new command, a new process is created, do the work, and terminated. 

This guarantees that no state is preserved between commands, except the config files. It is easier to see what must be done in order to finish each command, without any prior knowledge about the ledger state.

Following steps demonstrate how to make a transaction.

### Create 2 new accounts

```
$ ./client a c 2
2019/07/06 16:43:24 generating 2 accounts...
2019/07/06 16:43:25 account: c3838f7d165cc5a6f7b19315378712d1973d507a4e9bf6769ad6aeab5d9e89bf
2019/07/06 16:43:25 account: 7f5e114409a3e780110a0ec4e8e1f5b78948aac9724e2ff2a62c618702ad97ed
```

You can see the two newly generated account addresses. Yours should be different. so use your own addresses to finish this demonstration. 

The private keys of these accounts are saved in a wallet file (default wallet.toml), IN PLAIN TEXT. 

Later on, you can reference the accounts with a prefix of their addresses, just like what you do with docker command. For example, 'c' or 'c383' both references 'c3838f7d165cc5a6f7b19315378712d1973d507a4e9bf6769ad6aeab5d9e89bf'. 

You can also use full addresses not included in the wallet file. 

### Mint 100 coins into account c3...

```
$ ./client a m c3 100
2019/07/06 17:02:49 Please visit the following faucet service:
2019/07/06 17:02:49 http://faucet.testnet.libra.org/?amount=100000000&address=c3838f7d165cc5a6f7b19315378712d1973d507a4e9bf6769ad6aeab5d9e89bf
```

Copy & paste the link into you browser to actually mint the coins. 

### Check account state and balances

```
$ ./client q as c3
2019/07/06 17:04:36 Ledger info: version 485913, time 1562403876325982
2019/07/06 17:04:36 Account version: 485913
2019/07/06 17:04:36 Libra coin resource:
(*types.AccountResource)(0xc42001f880)({
 Balance: (uint64) 100000000,
 SequenceNumber: (uint64) 0,
 AuthenticationKey: ([]uint8) (len=32 cap=32) {
  00000000  c3 83 8f 7d 16 5c c5 a6  f7 b1 93 15 37 87 12 d1  |...}.\......7...|
  00000010  97 3d 50 7a 4e 9b f6 76  9a d6 ae ab 5d 9e 89 bf  |.=PzN..v....]...|
 },
 SentEventsCount: (uint64) 0,
 ReceivedEventsCount: (uint64) 0
})
```

Here you can see the ledger version, and 0x0.LibraAccount.T resource. The balance is 100,000,000 micro libra. 

Now if you check the other account 7f..., you will find it is not present in the ledger yet.

### Transfer 10 coins from c3... to 7f...

```
$ ./client t c3 7f 10
2019/07/06 17:09:07 Going to transfer 10000000 microLibra from c3838f7d165cc5a6f7b19315378712d1973d507a4e9bf6769ad6aeab5d9e89bf to 7f5e114409a3e780110a0ec4e8e1f5b78948aac9724e2ff2a62c618702ad97ed
2019/07/06 17:09:08 Refreshed sequence number of sender: 0
2019/07/06 17:09:08 Raw txn: {
    "sender_account": "w4OPfRZcxab3sZMVN4cS0Zc9UHpOm/Z2mtauq12eib8=",
...
    "max_gas_amount": 10000,
    "expiration_time": 1562404208
}
2019/07/06 17:09:08 Result:
(*ac.SubmitTransactionResponse)(0xc4201ad180)(ac_status:Accepted )
2019/07/06 17:09:08 Waiting until transaction is included in ledger...
2019/07/06 17:09:10 sequence number of sender: 1
```

Now if you check account 7f..., you will find a balance of 10,000,000 micro libra. And the account c3... has 90 left.

Note that the gas price is fixed to 0, and max gas amount 10,000. Actually, if you make a transaction with gas price set to 1, you will find that it takes 2000-3000 gas to execute the transaction. 

## Implementation

I deliberately name the packages, functions, and variables similar to official rust project, with subtle changes to suit golang idioms. It is easier to keep up with the rust project in this way.

After the official rust project is more stable, we can refactor the code to make it tastes more like golang, and to make the packages reusable for other projects.

Contributions are welcome.
