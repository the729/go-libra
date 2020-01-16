A Libra js client library with crypto verifications, for Nodejs and browsers.

# Usage

In order to work with browsers, `gopherjs-libra` uses gRPC-Web which is not directly compatible with gRPC. A proxy is needed to forward gRPC-Web requests to gRPC backends. You can setup an Envoy proxy (https://grpc.io/docs/tutorials/basic/web/), or use my demo proxy shown in the examples. 

## Node.js

```bash
npm install gopherjs-libra
```

```js
const { libra } = require("gopherjs-libra");

const defaultServer = "http://hk2.wutj.info:38080";

var client = libra.client(defaultServer, libra.trustedPeersFile)
client.queryTransactionRange(100, 2, true)
    .then(r => {
        r.getTransactions().map(txn => {
            console.log("Txn #", txn.getVersion())
            console.log("    Gas used (microLibra): ", txn.getGasUsed())
            console.log("    Major status: ", txn.getMajorStatus())
        })
    })
    .catch(e => {
        console.log("Error: ", e)
    })
```

## Browser

```html
<script src="gopherjs-libra.js"></script>
<script>
    var defaultServer = "http://hk2.wutj.info:38080";
    var client = libra.client(defaultServer, libra.trustedPeersFile);
</script>
```

# Examples

Several examples are included in [`example/nodejs`](https://github.com/the729/go-libra/tree/master/example/nodejs) folder.

There is a [pure frontend Libra blockchain explorer](http://pg.wutj.info/web_client/), based on gopherjs-libra.

## Generating Libra Accounts

This library does not handle account generation based on mnemonic. You will have to manage your public-private key pairs with other wallet libraries. 

For example, to simply generate a new account:
```js
const { sign } = require("tweetnacl");
const { libra } = require("gopherjs-libra");

keyPair = sign.keyPair();
address = libra.pubkeyToAddress(keyPair.publicKey);
// address is a libra account address (Uint8Array, 32-byte)
// keyPair.secretKey is the corresponding private key (Uint8Array, 64-byte)
```

# API Reference

## .client(server, trustedPeers)

Create a client using specified server and trusted peers. 

### Arguments
 - server (string): gRPC-Web server URL.
 - trustedPeers (string): a TOML formated string containing configurations of trusted peers. You can use `libra.trustedPeersFile`.

Returns a Libra Client instance. 

## .trustedPeersFile

A constant string. TOML formated string of the default trusted peers of the libra testnet.

## .accountResourcePath()

Returns a `Uint8Array`: the raw path to the Libra account resource, which is `0x01+hash(0x0.LibraAccount.T)`.

## .accountSentEventPath()

Returns a `Uint8Array`: the raw path to Libra coin sent events, which is `0x01+hash(0x0.LibraAccount.T)/sent_events_count/`.

## .accountReceivedEventPath()

Returns a `Uint8Array`: the raw path to Libra coin received events, which is `0x01+hash(0x0.LibraAccount.T)/received_events_count/`.

## .pubkeyToAddress(publicKey)

### Arguments
 - publicKey (Uint8Array): 32-byte ed25519 public key.

Returns SHA3 hash of input public key, which is used as Libra account address.

## Object: Client

Client represents a Libra client.

### Client.queryLedgerInfo()

Return a promise that resolves to a `provenLedgerInfo` object.

### Client.queryAccountState(address)

Argument:
 - address (Uint8Array): raw address bytes. 

Returns a promise that resolves to a `provenAccountState` object.

### Client.queryAccountSequenceNumber(address)

Argument:
 - address (Uint8Array): raw address bytes. 

Returns a promise that resolves to the sequence number (integer).

### Client.pollSequenceUntil(address, seq, expire)

Polls an account until its sequence number is greater or equal to the given seq.

Arguments:
 - address (Uint8Array): raw address bytes. 
 - seq (integer): expected sequence number.
 - expire (integer): expiration unix timestamp in seconds. The polling fails until the ledger timestamp is greater than `expire`.

Returns a promise that resolves when the expected sequence number is reached.

### Client.submitP2PTransaction(rawTxn)

Arguments:
 - rawTxn (Object): the raw transaction object, with following keys
   - senderAddr (Uint8Array): sender address
   - recvAddr (Uint8Array): receiver address
   - senderPrivateKey (Uint8Array): sender ed25519 secret key (64 bytes)
   - senderSeq (integer): current sender account sequence number
   - amountMicro (integer): amount to transfer in micro libra
   - maxGasAmount (integer): max gas amount in micro libra
   - gasUnitPrice (integer): micro libra per gas
   - expirationTimestamp (integer): transaction expiration unix timestamp

Returns a promise that resolves to the expected sequence number of this transaction. Use `pollSequenceUntil` afterward to make sure the transaction is included in the ledger.

### Client.submitContractTransaction(rawTxn)

Arguments:
 - rawTxn (Object): the raw transaction object, with following keys
   - senderAddr (Uint8Array): sender address
   - senderPrivateKey (Uint8Array): sender ed25519 secret key (64 bytes)
   - senderSeq (integer): current sender account sequence number
   - maxGasAmount (integer): max gas amount in micro libra
   - gasUnitPrice (integer): micro libra per gas
   - expirationTimestamp (integer): transaction expiration unix timestamp
   - payload (types.TransactionPayload): the code and/or args to publish/execute

Returns a promise that resolves to the expected sequence number of this transaction. Use `pollSequenceUntil` afterward to make sure the transaction is included in the ledger.

### Client.queryTransactionByAccountSeq(address, seq, withEvents)

Arguments:
 - address (Uint8Array): raw address bytes. 
 - seq (integer): sequence number to query.
 - withEvents (bool): whether to includes events in the returned value.

Returns a promise that resolves to a `provenTransaction` object.

### Client.queryTransactionRange(start, limit, withEvents)

Arguments:
 - start (integer): first transaction to return.
 - limit (integer): max number of transactions to return.
 - withEvents (bool): whether to includes events in the returned value.

Returns a promise that resolves to a `provenTransactionList` objects.

### Client.queryEventsByAccessPath(address, path, start, ascending, limit)

Arguments:
 - address (Uint8Array): raw address bytes. 
 - path (Uint8Array): `accountSentEventPath()` or `accountReceivedEventPath()`.
 - start (integer): the index of the first event.
 - limit (integer): max number of events to return.
 - ascending (bool): whether return events in ascending order.

Returns a promise that resolves to a list of `provenEvent` objects.

## Object: provenLedgerInfo

A `provenLedgerInfo` represents a proven state of the ledger at some version and some time. It is proven by quorum of trusted validators. It is the source of trust for all the following structures.

The member functions are self-descriptive. 

### .getVersion()
### .getTimestampUsec()
### .getTransactionAccumulatorHash()
### .getEpochNum()

## Object: provenAccountState

A `provenAccountState` represents an account state that is proven to be included in the ledger at a certain version (or proven included if isNil()).

### .getLedgerInfo()

Returns the `provenLedgerInfo` which proofs this object.

### .getVersion()

Returns the ledger version.

### .getAccountBlob()

Returns libra account blob `provenAccountBlob`.

### .isNil()

Returns `true` if the address is proven not included in the ledger.

## Object: provenAccountBlob

### .getLedgerInfo()

Returns the `provenLedgerInfo` which proofs this object.

### .getAddress()

Returns address (Uint8Array).

### getResource(path)

Returns `provenAccountResource` on the given path. Use `accountResourcePath()` as the path.

## Object: provenAccountResource

### .getLedgerInfo()
### .getAddress()
### .getBalance()
### .getSequenceNumber()
### .getSentEvents()
### .getReceivedEvents()
### .getDelegatedWithdrawalCapability()

## Object: provenTransaction

### .getLedgerInfo()
### .getVersion()
### .getMajorStatus()
### .getGasUsed()
### .getWithEvents()
### .getEvents()
### .getSignedTxn()

## Object: provenTransactionList

### .getLedgerInfo()
### .getTransactions()

## Object: provenEvent

### .getLedgerInfo()
### .getTransactionVersion()
### .getEventIndex()
### .getEvent()
