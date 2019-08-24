# Query transaction range

The query parameters are hard-coded in this example.

## Usage and example output

```
$ go run query_txn_range.go
2019/08/24 18:04:30 Txn #2397:
2019/08/24 18:04:30     Raw txn:
2019/08/24 18:04:30         Sender account: c16f24963efcc1c493353bd29e03ff53488fc7ba099ca42e13bbc3ed2d4e53c7
2019/08/24 18:04:30         Sender seq #0
2019/08/24 18:04:30         Program: 4c49425241564d0a010007014a00000004000000034e000000060000000c...
2019/08/24 18:04:30             (program name: peer_to_peer_transfer)
2019/08/24 18:04:30         Arg 0: fb6fcaa8730ffa4c2310fa04c14a7a6665bc1036a4ba70e5a97f4235667f09b7
2019/08/24 18:04:30         Arg 1: 80460a0c00000000
2019/08/24 18:04:30     Gas used: 0
2019/08/24 18:04:30     Events: (2 total)
2019/08/24 18:04:30       #0:
2019/08/24 18:04:30         Seq #0
2019/08/24 18:04:30         Addr: c16f24963efcc1c493353bd29e03ff53488fc7ba099ca42e13bbc3ed2d4e53c7
2019/08/24 18:04:30         Path: 0x0.LibraAccount.T/sent_events_count/
2019/08/24 18:04:30             (Event is: sent payment)
2019/08/24 18:04:30         Data:
2019/08/24 18:04:30             Payee: fb6fcaa8730ffa4c2310fa04c14a7a6665bc1036a4ba70e5a97f4235667f09b7
2019/08/24 18:04:30             Amount: 202
2019/08/24 18:04:30       #1:
2019/08/24 18:04:30         Seq #1
2019/08/24 18:04:30         Addr: fb6fcaa8730ffa4c2310fa04c14a7a6665bc1036a4ba70e5a97f4235667f09b7
2019/08/24 18:04:30         Path: 0x0.LibraAccount.T/received_events_count/
2019/08/24 18:04:30             (Event is: received payment)
2019/08/24 18:04:30         Data:
2019/08/24 18:04:30             Payer: c16f24963efcc1c493353bd29e03ff53488fc7ba099ca42e13bbc3ed2d4e53c7
2019/08/24 18:04:30             Amount: 202
2019/08/24 18:04:30 Txn #2398:
2019/08/24 18:04:30     Raw txn:
2019/08/24 18:04:30         Sender account: 000000000000000000000000000000000000000000000000000000000a550c18
2019/08/24 18:04:30         Sender seq #1874
2019/08/24 18:04:30         Program: 4c49425241564d0a010007014a000000060000000350000000060000000c...
2019/08/24 18:04:30             (program name: mint)
2019/08/24 18:04:30         Arg 0: 854563c50d20788fb6c11fac1010b553d722edb0c02f87c2edbdd3923726d13f
2019/08/24 18:04:30         Arg 1: 002d310100000000
2019/08/24 18:04:30     Gas used: 0
2019/08/24 18:04:30     Events: (2 total)
2019/08/24 18:04:30       #0:
2019/08/24 18:04:30         Seq #1857
2019/08/24 18:04:30         Addr: 000000000000000000000000000000000000000000000000000000000a550c18
2019/08/24 18:04:30         Path: 0x0.LibraAccount.T/sent_events_count/
2019/08/24 18:04:30             (Event is: sent payment)
2019/08/24 18:04:30         Data:
2019/08/24 18:04:30             Payee: 854563c50d20788fb6c11fac1010b553d722edb0c02f87c2edbdd3923726d13f
2019/08/24 18:04:30             Amount: 20
2019/08/24 18:04:30       #1:
2019/08/24 18:04:30         Seq #0
2019/08/24 18:04:30         Addr: 854563c50d20788fb6c11fac1010b553d722edb0c02f87c2edbdd3923726d13f
2019/08/24 18:04:30         Path: 0x0.LibraAccount.T/received_events_count/
2019/08/24 18:04:30             (Event is: received payment)
2019/08/24 18:04:30         Data:
2019/08/24 18:04:30             Payer: 000000000000000000000000000000000000000000000000000000000a550c18
2019/08/24 18:04:30             Amount: 20
```
