// Use following require statement when using npm.
// const { libra } = require("gopherjs-libra");
const { libra } = require("../../gopherjs-libra/gopherjs-libra.js");

const defaultServer = "http://hk2.wutj.info:38080",
    waypoint = "0:4d4d0feaa9378069f8fcee71980e142273837e108702d8d7f93a8419e2736f3f";

var client = libra.client(defaultServer, waypoint)
client.queryTransactionRange(13942242, 2, true)
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
