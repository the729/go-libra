A Libra js client library with crypto verifications, for NodeJS and browsers.

Compatible with libra testnet updated on 4/8/2020.

# Usage

In order to work with browsers, `gopherjs-libra` uses gRPC-Web which is not directly compatible with gRPC. A proxy is needed to forward gRPC-Web requests to gRPC backends. You can setup an Envoy proxy (https://grpc.io/docs/tutorials/basic/web/), or use my demo proxy shown in the examples. 

To setup an Envoy proxy using Docker:
```bash
docker run -d --name envoy-libra -p 38080:8080 wutianji/envoy-libra
```

## Node.js

```bash
npm install gopherjs-libra
```

```js
const { libra } = require("gopherjs-libra");

const defaultServer = "http://hk2.wutj.info:38080",
    waypoint = "insecure";

var client = libra.client(defaultServer, waypoint)
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
    var client = libra.client(defaultServer, "insecure");
</script>
```

# Examples

For NodeJS, there are several examples included in [`example/nodejs`](https://github.com/the729/go-libra/tree/master/example/nodejs) folder.

For browsers, there is a [pure front-end Libra blockchain explorer](http://pg.wutj.info/web_client/), based on gopherjs-libra.

## Generating Libra Accounts

This library does not handle account generation based on mnemonic. You will have to manage your public-private key pairs with other wallet libraries. 

For example, to simply generate a new account:
```js
const { sign } = require("tweetnacl");
const { libra } = require("gopherjs-libra");

keyPair = sign.keyPair();
address = libra.pubkeyToAddress(keyPair.publicKey);
authKey = libra.pubkeyToAuthKey(keyPair.publicKey);
// address is a libra account address (Uint8Array, 16-byte)
// authKey is the full auth key (Uint8Array, 32-byte)
// Use authkey[0:16] if you need authKeyPrefix
// keyPair.secretKey is the corresponding private key (Uint8Array, 64-byte)
```

# API Reference

### `.client(server, initState)`

Create a client using specified server and state or trusted waypoint. 

Returns a Libra Client instance. 

Arguments:
 - `server` (string): gRPC-Web server URL (http://host:port).
 - `initState` (string or object): initial state of the client, or a trusted waypoint string.

If `initState` is a string, it should be one of the following cases:
 - "insecure": the client will trust whatever the ledger has. This option is useful for the testnet,
which gets reset every now and then.
 - a waypoint string in the format of `version:hash`. It will provide a root of trust. The client 
only trust the ledger if it matches the waypoint. Due to an [issue](https://community.libra.org/t/how-can-a-client-detects-a-fork/1243/5) that an arbitrary waypoint does provides enough information for 
consistency check, for now, it's required to provide a waypoint with version=0 as an `initState`.

`initState` can also be an object, with which you can restore a client's previous state (i.e validator set and known version). In this case, you should use whatever `client.getState()` returns. It has the following keys:
 - `epoch` (integer): epoch number of the `validator_set`
 - `validator_set` (array of object): a list of trusted validators at a certain `epoch`. Each validator has the following keys:
   - `addr` (string): hex string of the validator address.
   - `c` (string): hex string of the validator's consensus public key.
   - `power` (integer): the validator's voting power.
 - `knwon_version` (integer): known version for consistency check.
 - `subtrees` (array of string): each string is a hex string of a subtree hash value, used for consistency check.

If you don't need strict blockchain consistency check, you are recommended to use a version 0 waypoint to init a client. 

Otherwise, you should periodically export the client state with `client.getState()` and save the output to local store or file. And reload the state when re-init a new client.

### `.resourcePath(addr, module, name, accesses ...)`

Build a resource path and returns a `Uint8Array`. Note that if you need a trailing '/', put an extra empty string in `accesses`. For example, the following call build the path `0x0.LibraAccount.T/sent_events_count/`
```js
resourcePath(
    new Uint8Array(32),         // addr
    "LibraAccount", "T",        // module and struct name
    "sent_events_count", ""     // access path
)
```
which is equivalent to `accountSentEventPath()`.

Arguments:
 - addr (Uint8Array): 32-byte address.
 - module (string): module name.
 - name (string): type name.
 - accesses (strings): 0 or more access paths. 

### `.accountResourcePath()`

Returns a `Uint8Array`: the raw path to the Libra account resource, which is `0x01+hash(0x0.LibraAccount.T)`.

### `.balanceResourcePath()`

Returns a `Uint8Array`: the raw path to the Libra balance resource, which is `0x01+hash(0x0.LibraAccount.Balance)`.

### `.accountSentEventPath()`

Returns a `Uint8Array`: the raw path to Libra coin sent events, which is `0x01+hash(0x0.LibraAccount.T)/sent_events_count/`.

### `.accountReceivedEventPath()`

Returns a `Uint8Array`: the raw path to Libra coin received events, which is `0x01+hash(0x0.LibraAccount.T)/received_events_count/`.

### `.pubkeyToAddress(publicKey)`

Arguments:
 - publicKey (Uint8Array): 32-byte ed25519 public key.

Returns last 16 bytes of SHA3 hash of input public key, which is used as Libra account address.

### `.pubkeyToAuthKey(publicKey)`

Arguments:
 - publicKey (Uint8Array): 32-byte ed25519 public key.

Returns the full SHA3 hash of input public key, which is used as initial Libra account auth key.

## Object: `Client`

Client represents a Libra client.

### `Client.queryLedgerInfo()`

Return a promise that resolves to a `provenLedgerInfo` object.

### `Client.queryAccountState(address)`

Argument:
 - `address` (Uint8Array): raw address bytes. 

Returns a promise that resolves to a `provenAccountState` object.

### `Client.queryAccountSequenceNumber(address)`

Argument:
 - `address` (Uint8Array): raw address bytes. 

Returns a promise that resolves to the sequence number (integer).

### `Client.pollSequenceUntil(address, seq, expire)`

Polls an account until its sequence number is greater or equal to the given seq.

Arguments:
 - `address` (Uint8Array): raw address bytes. 
 - `seq` (integer): expected sequence number.
 - `expire` (integer): expiration unix timestamp in seconds. The polling fails until the ledger timestamp is greater than `expire`.

Returns a promise that resolves when the expected sequence number is reached.

### `Client.submitP2PTransaction(p2pTxn)`

Submit a p2p Libra coin payment transaction.

Arguments:
 - `p2pTxn` (Object): a p2p transaction object, with following keys
   - `senderAddr` (Uint8Array): sender address
   - `recvAddr` (Uint8Array): receiver address
   - `recvAuthKeyPrefix` (Uint8Array): receiver auth key prefix (first 16 bytes of auth key)
   - `senderPrivateKey` (Uint8Array): sender ed25519 secret key (64 bytes)
   - `senderSeq` (integer): current sender account sequence number
   - `amountMicro` (integer): amount to transfer in micro libra
   - `maxGasAmount` (integer): max gas amount in micro libra
   - `gasUnitPrice` (integer): micro libra per gas
   - `expirationTimestamp` (integer): transaction expiration unix timestamp

Returns a promise that resolves to the expected sequence number of this transaction. Use `pollSequenceUntil` afterward to make sure the transaction is included in the ledger.

### `Client.submitRawTransaction(rawTxn)`

Submit a raw user transaction. There are 2 types of user transactions: script and module. 
Script transaction has a piece of `code` with arguments(`args`), which will be executed on the libra Move VM.
Module transaction contains a code `module`. 

A transaction payload should have either `code` or `module`, not both. 
You will need a Move compiler to generate binary codes for both type of transactions. 

Transaction arguments have 4 types: bool, uint64, bytes, address.
The arguments specified in the payload will be converted to these 4 types according to the following rules.

1. JS bool -> Move bool
2. JS number -> Move uint64
   
   Floating point numbers are truncated, and because of the limitation of float64, integers greater or equal to 1<<53 will be truncated too. If you need large integers, use explicit type definition.

3. JS Uint8Array -> Move Address or raw bytes, based on length

   If the length is 32, it will become an address, otherwise a raw byte array. If you need a raw byte array whose length is exactly 32, use explicit type definition. 

4. Explicit type definition: JS Object with keys:
   - type (string): 'uint64' or 'bytes'
   - value (Uint8Array): if type is 'uint64', value should have exactly 8 bytes.

Arguments:
 - `rawTxn` (Object): a raw transaction object, with following keys
   - `senderAddr` (Uint8Array): sender address
   - `senderPrivateKey` (Uint8Array): sender ed25519 secret key (64 bytes)
   - `senderSeq` (integer): current sender account sequence number
   - `payload` (Object): transaction payload
     - `code` (Uint8Array): binary MoveVM code, for script transaction
     - `args` (Array of object): arguments for the script
     - `module` (Uint8Array): binary MoveVM module, for module transaction
   - `maxGasAmount` (integer): max gas amount in micro libra
   - `gasUnitPrice` (integer): micro libra per gas
   - `expirationTimestamp` (integer): transaction expiration unix timestamp

Returns a promise that resolves to the expected sequence number of this transaction. Use `pollSequenceUntil` afterward to make sure the transaction is included in the ledger.

### `Client.queryTransactionByAccountSeq(address, seq, withEvents)`

Arguments:
 - `address` (Uint8Array): raw address bytes. 
 - `seq` (integer): sequence number to query.
 - `withEvents` (bool): whether to includes events in the returned value.

Returns a promise that resolves to a `provenTransaction` object.

### `Client.queryTransactionRange(start, limit, withEvents)`

Arguments:
 - `start` (integer): first transaction to return.
 - `limit` (integer): max number of transactions to return.
 - `withEvents` (bool): whether to includes events in the returned value.

Returns a promise that resolves to a `provenTransactionList` objects.

### `Client.queryEventsByAccessPath(address, path, start, ascending, limit)`

Arguments:
 - `address` (Uint8Array): raw address bytes. 
 - `path` (Uint8Array): `accountSentEventPath()` or `accountReceivedEventPath()`.
 - `start` (integer): the index of the first event.
 - `limit` (integer): max number of events to return.
 - `ascending` (bool): whether return events in ascending order.

Returns a promise that resolves to a list of `provenEvent` objects.

### `Client.getState()`

Returns an object with current validator set and known version subtrees. The returned object can be used to init a new client with [`.client()`](#clientserver-initstate) or restore a client with `.setState()`. 

See [`.client()`](#clientserver-initstate) for detailed description of the returned object.

### `Client.setState(state)`

Restore client state, i.e. validator set and known version subtrees. See [`.client()`](#clientserver-initstate) for detailed description of the state object.

### `Client.getLatestWaypoint()`

Returns the latest waypoint the client has encountered, in the format of "version:hash" string. 

It is useful when you start the client with "insecure" waypoint, and decided to only trust the current chain from now on. You can either save the whole client state with `getState()`, or you can export a waypoint with this function. 

Note that only a version 0 waypoint can be used to init a client alone. See [`.client()`](#clientserver-initstate) for detailed description.

## Object: `provenLedgerInfo`

Represents a proven state of the ledger at some version and some time. It is proven by quorum of trusted validators. It is the source of trust for all the following structures.

It has a list of getters (with return type), whose names are self-descriptive:
 - `.getVersion()` (integer)
 - `.getTimestampUsec()` (integer)
 - `.getTransactionAccumulatorHash()` (Uint8Array)
 - `.getEpochNum()` (integer)

## Object: `provenAccountState`

Represents an account state that is proven to be included in the ledger at a certain version (or proven to be NOT included if `isNil()`).

It has a list of getters (with return type), whose names are self-descriptive:
 - `.getLedgerInfo()` (`provenLedgerInfo`)
 - `.getVersion()` (integer)
 - `.getAccountBlob()` (`provenAccountBlob`)
 - `.isNil()` (bool): Returns `true` only if the address is proven to be NOT included in the ledger.

## Object: `provenAccountBlob`

Represents an account blob (key-value map) that is proven to be at a certain version.

It has a list of getters (with return type), whose names are self-descriptive:
 - `.getLedgerInfo()` (`provenLedgerInfo`)
 - `.getAddress()` (Uint8Array)
 - `.getResource(path)` (Uint8Array): Returns a binary resource content on the given access path. Use `resourcePath()` to build a path.
 - `.getLibraResources()` (object): Returns the Libra account resources, i.e. `accountResource` and `balanceResource`.

## Object: `provenTransaction`

Represents a transaction proven to be included in the ledger.

It has a list of getters (with return type), whose names are self-descriptive:
 - `.getLedgerInfo()` (`provenLedgerInfo`)
 - `.getVersion()` (integer)
 - `.getMajorStatus()` (integer)
 - `.getGasUsed()` (integer)
 - `.getWithEvents()` (bool)
 - `.getEvents()` (array of `provenEvent`)
 - `.getBlockMetadata()` (null or object)
 - `.getSignedTxn()` (null or object)

## Object: `provenTransactionList`

Represents a list of transactions proven to be included in the ledger.

It has a list of getters (with return type), whose names are self-descriptive:
 - `.getLedgerInfo()` (`provenLedgerInfo`)
 - `.getTransactions()` (array of `provenTransaction`)

## Object: `provenEvent`

Represents an event emitted with a transaction.

It has a list of getters (with return type), whose names are self-descriptive:
 - `.getLedgerInfo()` (`provenLedgerInfo`)
 - `.getTransactionVersion()` (integer)
 - `.getEventIndex()` (integer): index within the transaction
 - `.getEvent()` (object)
