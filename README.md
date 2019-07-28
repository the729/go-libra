# go-libra
This is a golang implementation of the Libra blockchain (https://github.com/libra/libra) client library. 

It has all cryptographic verification algorithms, including validator-signature-based consensus verification, ledger history accumulator proof, and account state sparse Merkle tree proof, etc. 

## Features

- ✓ Connect to testnet AdmissionControl server with gRPC
- ✓ Data models with all necessary cryptographic verification algorithms
  - ✓ Ledger state: signature-based consensus verification
  - ✓ Transaction info: ledger history accumulator proof
  - ✓ Transaction list: ledger history accumulator proof on a range of transactions
  - ✓ Transaction signature: ed25519 signature
  - ✓ Account state: sparse Merkle tree proof
  - ✓ Events: event list hash based on Merkle tree accumulator
- ✓ Query account states
- ✓ Make P2P transaction, and wait for ledger inclusion
- ✓ Query transactions by ledger version
- ✓ Query account transaction by sequence number

## Examples

See `example/` folder. More examples coming soon.

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
	addr := client.MustToAddress(addrStr)

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

## Implementation

I deliberately name the packages, functions, and variables similar to official rust project, with subtle changes to suit golang idioms. It is easier to keep up with the rust project in this way.

After the official rust project is more stable, we can refactor the code to make it tastes more like golang, and to make the packages reusable for other projects.

Contributions are welcome.
