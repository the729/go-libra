// Use following require statement when using npm.
// const { libra } = require("gopherjs-libra");
const { libra } = require("../../gopherjs-libra/gopherjs-libra.js");

const fromHexString = hexString =>
    new Uint8Array(hexString.match(/.{1,2}/g).map(byte => parseInt(byte, 16)));

const toHexString = bytes =>
    bytes.reduce((str, byte) => str + byte.toString(16).padStart(2, '0'), '');

const defaultServer = "http://hk2.wutj.info:38080",
    waypoint = "0:4d4d0feaa9378069f8fcee71980e142273837e108702d8d7f93a8419e2736f3f";

var senderAddr = fromHexString("42f5745128c05452a0c68272de8042b1"),
    priKey = fromHexString("996911072ee011ffa44a1325e0da593ff3b9374e255115f223cbdffb6bfa0bcfba60d1f8edd6923f59cf9125d3ac80e389afa4e2b8d0e4f1183a30a0270fde71"),
    recvAddr = fromHexString("5817cd6e6e84c110c43efca22df54172"),
    recvAuthKeyPrefix = fromHexString("26c7bfaa8e0f32206f35bf6d44b43c9c");

var client = libra.client(defaultServer, waypoint)
client.queryAccountSequenceNumber(senderAddr)
    .then(r => {
        var txn = {
            "senderAddr": senderAddr,
            "recvAddr": recvAddr,
            "recvAuthKeyPrefix": recvAuthKeyPrefix,
            "senderPrivateKey": priKey,
            "senderSeq": r,
            "amountMicro": 2 * 1000000,
            "maxGasAmount": 500000,
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
