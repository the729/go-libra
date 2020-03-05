// Use following require statement when using npm.
// const { libra } = require("gopherjs-libra");
const { libra } = require("../../gopherjs-libra/gopherjs-libra.js");

const fromHexString = hexString =>
    new Uint8Array(hexString.match(/.{1,2}/g).map(byte => parseInt(byte, 16)));

const toHexString = bytes =>
    bytes.reduce((str, byte) => str + byte.toString(16).padStart(2, '0'), '');

const defaultServer = "http://hk2.wutj.info:38080",
    waypoint = "0:bf7e1eef81af68cc6b4801c3739da6029c778a72e67118a8adf0dd759f188908";

var addrStr = "18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a"
var addr = fromHexString(addrStr)

var client = libra.client(defaultServer, waypoint)
client.queryEventsByAccessPath(addr, libra.accountSentEventPath(), 0, true, 10)
    .then(r => {
        r.map(ev => {
            console.log("Txn #", ev.getTransactionVersion(), " event #", ev.getEventIndex())
            var evBody = ev.getEvent()
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
