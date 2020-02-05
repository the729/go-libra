// Use following require statement when using npm.
// const { libra } = require("gopherjs-libra");
const { libra } = require("../../gopherjs-libra/gopherjs-libra.js");

const fromHexString = hexString =>
    new Uint8Array(hexString.match(/.{1,2}/g).map(byte => parseInt(byte, 16)));

const toHexString = bytes =>
    bytes.reduce((str, byte) => str + byte.toString(16).padStart(2, '0'), '');

const defaultServer = "http://hk2.wutj.info:38080",
    waypoint = "0:997acd1b112a19eb1d2d3dff78677a0009343727926071c3858aeff2ea3499bf";

var senderAddr = fromHexString("18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a"),
    priKey = fromHexString("657cd8ed5e434cc4f874d6822889f637957f0145c67e2b055c9954c936670a61e57ea705e00e3ecaf417b4285cd0a69b1d79406914581456c1ce278b81a48674"),
    recvAddr = fromHexString("e89a0d93fcf1ca4423328c1bddebe6c02da666808993c8a888ff7a8bad19ffd5");

var client = libra.client(defaultServer, waypoint)
client.queryAccountSequenceNumber(senderAddr)
    .then(r => {
        var txn = {
            "senderAddr": senderAddr,
            "recvAddr": recvAddr,
            "senderPrivateKey": priKey,
            "senderSeq": r,
            "amountMicro": 2 * 1000000,
            "maxGasAmount": 140000,
            "gasUnitPrice": 0,
            "expirationTimestamp": parseInt(Date.now() / 1000) + 60,
        };
        return client.submitP2PTransaction(txn);
    })
    .then(r => {
        console.log("Polling sequence number until ", r)
        return client.pollSequenceUntil(senderAddr, r, parseInt(Date.now() / 1000) + 60)
    })
    .then(r => {
        console.log("done.")
    })
    .catch(e => {
        console.log("Error:", e)
    })
