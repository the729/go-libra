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
        console.log("Proven at ledger version / time: ", r.getLedgerInfo().getVersion(), r.getLedgerInfo().getTimestampUsec() / 1000000)
    })
    .catch(e => {
        console.log("Error: ", e)
    })
