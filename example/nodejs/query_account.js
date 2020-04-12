// Use following require statement when using npm.
// const { libra } = require("gopherjs-libra");
const { libra } = require("../../gopherjs-libra/gopherjs-libra.js");

const fromHexString = hexString =>
    new Uint8Array(hexString.match(/.{1,2}/g).map(byte => parseInt(byte, 16)));

const toHexString = bytes =>
    bytes.reduce((str, byte) => str + byte.toString(16).padStart(2, '0'), '');

const defaultServer = "http://hk2.wutj.info:38080",
    waypoint = "0:4d4d0feaa9378069f8fcee71980e142273837e108702d8d7f93a8419e2736f3f";

var addrStr = "42f5745128c05452a0c68272de8042b1"
var addr = fromHexString(addrStr)

var client = libra.client(defaultServer, waypoint)
client.queryAccountState(addr)
    .then(r => {
        if (r.isNil()) {
            throw "Account " + addrStr + " does not exist at version " + r.getVersion() + "."
        }
        console.log("Proven at ledger version / time: ", r.getLedgerInfo().getVersion(), r.getLedgerInfo().getTimestampUsec() / 1000000)
        return r.getAccountBlob()
    })
    .then(r => r.getLibraResources())
    .then(r => {
        console.log("Address: ", addrStr)
        console.log("Balance (microLibra): %d", r.balanceResource.Coin)
        console.log("Sequence Number: ", r.accountResource.SequenceNumber)
        console.log("Sent Events: ", r.accountResource.SentEvents)
        console.log("Received Events: ", r.accountResource.ReceivedEvents)
        console.log("Delegated withdrawal capability: ", r.accountResource.DelegatedWithdrawalCapability)
        console.log("Event generator: ", r.accountResource.EventGenerator)
    })
    .catch(e => {
        console.log("Error: ", e)
    })
