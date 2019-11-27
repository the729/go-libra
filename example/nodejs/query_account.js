const { libra } = require("gopherjs-libra");

const fromHexString = hexString =>
    new Uint8Array(hexString.match(/.{1,2}/g).map(byte => parseInt(byte, 16)));

const toHexString = bytes =>
    bytes.reduce((str, byte) => str + byte.toString(16).padStart(2, '0'), '');

const defaultServer = "http://hk2.wutj.info:38080";

var addrStr = "18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a"
var addr = fromHexString(addrStr)

var client = libra.client(defaultServer, libra.trustedPeersFile)
client.queryAccountState(addr)
    .then(r => {
        if (r.isNil()) {
            throw "Account " + addrStr + " does not exist at version " + r.getVersion() + "."
        }
        return r.getAccountBlob()
    })
    .then(r => r.getResource(libra.accountResourcePath()))
    .then(r => {
        console.log("Balance (microLibra): %d", r.getBalance())
        console.log("Sequence Number: ", r.getSequenceNumber())
        console.log("Sent Events: ", r.getSentEvents())
        console.log("Received Events: ", r.getReceivedEvents())
        console.log("DelegatedWithdrawalCapability: ", r.getDelegatedWithdrawalCapability())
        console.log("Proven at ledger version / time: ", r.getLedgerInfo().getVersion(), r.getLedgerInfo().getTimestampUsec() / 1000000)
    })
    .catch(e => {
        console.log("Error: ", e)
    })
