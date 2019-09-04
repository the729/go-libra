# CLI Client example

## Usage

The commands are similar to those of the official rust implementation. However, this is not an interactive CLI program, meaning that every time when you execute a new command, a new process is created to finish the work, and terminated.

This guarantees that no state is preserved between commands, except the config files. It is easier to see what must be done in order to finish each command, without any prior knowledge about the ledger state.

Following steps demonstrate how to make a transaction.

### Create 2 new accounts

```
$ ./cli_client a c 2
2019/07/06 16:43:24 generating 2 accounts...
2019/07/06 16:43:25 account: 18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a
2019/07/06 16:43:25 account: 34d9fb3daedb9cf71c2d9f024efe058e502f848bd2be96f795c79ed0bb56e247
```

You can see the two newly generated account addresses. (Your accounts should be different, so finish the following demo with your own addresses.)

Be careful that the private keys of these accounts are saved in a wallet file (default wallet.toml), IN PLAIN TEXT. 

Later on, you can reference the accounts with a prefix of their addresses, just like what you do with docker command. For example, '1' or '18b' both references '18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a'. 

You can also use full addresses not included in the wallet file. 

### Mint 100 coins into account c3...

```
$ ./cli_client a m 18b 100
2019/09/04 23:01:18 Going to POST to faucet service: http://faucet.testnet.libra.org/?amount=100000000&address=18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a
2019/09/04 23:01:20 Respone (code=200): 916
```

Copy & paste the link into you browser to actually mint the coins. 

### Check account state and balances

```
$ ./cli_client q as 18b
2019/09/04 22:58:55 Ledger info: version 939, time 1567609135441133
2019/09/04 22:58:55 Account version: 939
2019/09/04 22:58:55 Balance (microLibra): 100000000
2019/09/04 22:58:55 Sequence Number: 0
2019/09/04 22:58:55 SentEventsCount: 0
2019/09/04 22:58:55     Key: f6e51256374f8072e9e56a06b9ddd69e4131e5779dd2d2e4679a110b04b22aa7
2019/09/04 22:58:55 ReceivedEventsCount: 2
2019/09/04 22:58:55     Key: 149e1203240cd9b8ad0013abb0ff3e65f85decca62a65b871bed7cf9855dc6cc
2019/09/04 22:58:55 DelegatedWithdrawalCapability: false
```

Here you can see the ledger version, and `0x0.LibraAccount.T` resource. The balance is 100,000,000 micro libra. 

Now if you check the other account 34d..., you will find it is not present in the ledger yet.

### Transfer 10 coins from 18b... to 34d...

```
$ ./cli_client t 18b 34d 10
2019/09/04 23:03:34 Going to transfer 10000000 microLibra from 18b553473df736e5e363e7214bd624735ca66ac22a7048e3295c9b9b9adfc26a to 34d9fb3daedb9cf71c2d9f024efe058e502f848bd2be96f795c79ed0bb56e247
2019/09/04 23:03:34 Max gas: 140000, Gas price: 0, Expiration: 2019-09-04 23:04:34.8165507 +0800 CST m=+60.006835001
2019/09/04 23:03:34 Get current account sequence of sender...
2019/09/04 23:03:35 ... is 4
2019/09/04 23:03:35 Submit transaction...
2019/09/04 23:03:35 Waiting until transaction is included in ledger...
2019/09/04 23:03:37 sequence number: 5, ledger version: 947
2019/09/04 23:03:37 done.
```

Now if you check account 34d..., you will find a balance of 10,000,000 micro libra. And the account 18b... has 90 left.
