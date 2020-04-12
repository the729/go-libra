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
client.queryEventsByAccessPath(addr, libra.accountSentEventPath(), 0, true, 10)
    .then(r => {
        r.map(ev => {
            console.log("Txn #", ev.getTransactionVersion(), " event #", ev.getEventIndex())
            var evBody = ev.getEvent().Value
            console.log("    Key: ", toHexString(evBody.Key))
            console.log("    Seq num: ", evBody.SequenceNumber)
            console.log("    Raw data: ", toHexString(evBody.Data))
        })
        if (r.length > 0) {
            console.log("Proven at ledger version / time: ", r[0].getLedgerInfo().getVersion(), r[0].getLedgerInfo().getTimestampUsec() / 1000000)
        }
    })
    .catch(e => {
        console.log("Error: ", e)
    })
