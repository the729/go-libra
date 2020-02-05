// Use following require statement when using npm.
// const { libra } = require("gopherjs-libra");
const { libra } = require("../../gopherjs-libra/gopherjs-libra.js");

const defaultServer = "http://hk2.wutj.info:38080",
    waypoint = "0:997acd1b112a19eb1d2d3dff78677a0009343727926071c3858aeff2ea3499bf";

var client = libra.client(defaultServer, waypoint)
client.queryTransactionRange(8207475, 2, true)
    .then(r => {
        r.getTransactions().map(txn => {
            console.log("Txn #", txn.getVersion())
            console.log("    Gas used (microLibra): ", txn.getGasUsed())
            console.log("    Major status: ", txn.getMajorStatus())
            console.log("    Signed txn: ", txn.getSignedTxn())
        })
        console.log("Proven at ledger version / time: ", r.getLedgerInfo().getVersion(), r.getLedgerInfo().getTimestampUsec() / 1000000)
    })
    .catch(e => {
        console.log("Error: ", e)
    })
