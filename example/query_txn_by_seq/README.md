# Query transaction by account sequence

The query parameters are hard-coded in this example.

## Usage and example output

```
$ go run query_txn_by_seq.go
2019/12/15 21:55:10 Txn #8207475:
2019/12/15 21:55:10     Raw txn: hash=44c06c0e9d80ab86aa333d6eea95359040343f5f5cedb3b313dd19cae64fb585
2019/12/15 21:55:10         Sender account: 18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a
2019/12/15 21:55:10         Sender seq #0
2019/12/15 21:55:10         Payload is Script ...
2019/12/15 21:55:10         Program: 4c49425241564d0a010007014a00000004000000034e000000060000000d...
2019/12/15 21:55:10             (program name: peer_to_peer_transfer)
2019/12/15 21:55:10         Arg 0: addr (34d9fb3daedb9cf71c2d9f024efe058e502f848bd2be96f795c79ed0bb56e247)
2019/12/15 21:55:10         Arg 1: u64  (10000000)
2019/12/15 21:55:10     Max gas amount (gas units): 140000
2019/12/15 21:55:10     Gas unit price (microLibra/unit): 0
2019/12/15 21:55:10     Expiration timestamp: 1576414055
2019/12/15 21:55:10     Gas used (microLibra): 0
2019/12/15 21:55:10     Major status: 4001 - EXECUTED
2019/12/15 21:55:10     Events: (2 total)
2019/12/15 21:55:10       #0:
2019/12/15 21:55:10         Key: 010000000000000018b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a
2019/12/15 21:55:10         Seq #0
2019/12/15 21:55:10         Raw event: 809698000000000034d9fb3daedb9cf71c2d9f024efe058e502f848bd2be ...
2019/12/15 21:55:10             Amount (microLibra): 10000000
2019/12/15 21:55:10             Opponent address: 34d9fb3daedb9cf71c2d9f024efe058e502f848bd2be96f795c79ed0bb56e247
2019/12/15 21:55:10       #1:
2019/12/15 21:55:10         Key: 000000000000000034d9fb3daedb9cf71c2d9f024efe058e502f848bd2be96f795c79ed0bb56e247
2019/12/15 21:55:10         Seq #0
2019/12/15 21:55:10         Raw event: 809698000000000018b553473df736e5e363e7214bd624735ca66ac22a70 ...
2019/12/15 21:55:10             Amount (microLibra): 10000000
2019/12/15 21:55:10             Opponent address: 18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a
```
