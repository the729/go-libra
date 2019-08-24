# Query transaction by account sequence

The query parameters are hard-coded in this example.

## Usage and example output

```
$ go run query_txn_by_seq.go
2019/08/24 18:07:29 Txn #166:
2019/08/24 18:07:29     Raw txn:
2019/08/24 18:07:29         Sender account: 18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a
2019/08/24 18:07:29         Sender seq #0
2019/08/24 18:07:29         Program: 4c49425241564d0a010007014a00000004000000034e000000060000000c...
2019/08/24 18:07:29             (program name: peer_to_peer_transfer)
2019/08/24 18:07:29         Arg 0: 34d9fb3daedb9cf71c2d9f024efe058e502f848bd2be96f795c79ed0bb56e247
2019/08/24 18:07:29         Arg 1: 8096980000000000
2019/08/24 18:07:29     Gas used: 0
2019/08/24 18:07:29     Events: (2 total)
2019/08/24 18:07:29       #0:
2019/08/24 18:07:29         Seq #0
2019/08/24 18:07:29         Addr: 18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a
2019/08/24 18:07:29         Path: 0x0.LibraAccount.T/sent_events_count/
2019/08/24 18:07:29             (Event is: sent payment)
2019/08/24 18:07:29         Data:
2019/08/24 18:07:29             Payee: 34d9fb3daedb9cf71c2d9f024efe058e502f848bd2be96f795c79ed0bb56e247
2019/08/24 18:07:29             Amount: 10
2019/08/24 18:07:29       #1:
2019/08/24 18:07:29         Seq #0
2019/08/24 18:07:29         Addr: 34d9fb3daedb9cf71c2d9f024efe058e502f848bd2be96f795c79ed0bb56e247
2019/08/24 18:07:29         Path: 0x0.LibraAccount.T/received_events_count/
2019/08/24 18:07:29             (Event is: received payment)
2019/08/24 18:07:29         Data:
2019/08/24 18:07:29             Payer: 18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a
2019/08/24 18:07:29             Amount: 10
```
