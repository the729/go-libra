package main

import (
	"context"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/miratronix/jopher"
	"github.com/the729/go-libra/client"
	"github.com/the729/go-libra/types"
)

func main() {
	var exports *js.Object
	if js.Module == js.Undefined {
		exports = js.Global
	} else {
		exports = js.Module.Get("exports")
	}
	exports.Set("libra", map[string]interface{}{
		"client":                   newClient,
		"trustedPeersFile":         trustedPeersFile,
		"accountResourcePath":      types.AccountResourcePath,
		"accountSentEventPath":     types.AccountSentEventPath,
		"accountReceivedEventPath": types.AccountReceivedEventPath,
		"pubkeyToAddress":          client.PubkeyMustToAddress,
	})
}

type jsClient struct {
	*js.Object

	queryAccountState            func(...interface{}) *js.Object                                     `js:"queryAccountState"`
	queryAccountSequenceNumber   func(...interface{}) *js.Object                                     `js:"queryAccountSequenceNumber"`
	submitP2PTransaction         func(*js.Object) *js.Object                                         `js:"submitP2PTransaction"`
	pollSequenceUntil            func(types.AccountAddress, uint64, int64) *js.Object                `js:"pollSequenceUntil"`
	queryTransactionByAccountSeq func(types.AccountAddress, uint64, bool) *js.Object                 `js:"queryTransactionByAccountSeq"`
	queryTransactionRange        func(uint64, uint64, bool) *js.Object                               `js:"queryTransactionRange"`
	queryEventsByAccessPath      func(types.AccountAddress, []byte, uint64, bool, uint64) *js.Object `js:"queryEventsByAccessPath"`
}

func newClient(server, trustedPeers string) *js.Object {
	c, err := client.New(server, trustedPeers)
	if err != nil {
		panic(err)
	}
	jc := jsClient{Object: js.Global.Get("Object").New()}
	jc.queryAccountState = jopher.Promisify(func(addr types.AccountAddress) (*js.Object, error) {
		r, err := c.QueryAccountState(context.TODO(), addr)
		return wrapProvenAccountState(r), err
	})
	jc.queryAccountSequenceNumber = jopher.Promisify(func(addr types.AccountAddress) (uint64, error) {
		return c.QueryAccountSequenceNumber(context.TODO(), addr)
	})

	promiseSubmitP2PTransaction := jopher.Promisify(func(txn *js.Object) (uint64, error) {
		type jsP2PTxn struct {
			*js.Object
			SenderAddr          []byte `js:"senderAddr"`
			SenderPriKey        []byte `js:"senderPrivateKey"`
			RecvAddr            []byte `js:"recvAddr"`
			SenderSeq           uint64 `js:"senderSeq"`
			AmountMicro         uint64 `js:"amountMicro"`
			MaxGasAmount        uint64 `js:"maxGasAmount"`
			GasUnitPrice        uint64 `js:"gasUnitPrice"`
			ExpirationTimestamp int64  `js:"expirationTimestamp"`
		}
		jstxn := &jsP2PTxn{Object: txn}
		rawTxn, _ := client.NewRawP2PTransaction(
			jstxn.SenderAddr, jstxn.RecvAddr,
			jstxn.SenderSeq,
			jstxn.AmountMicro, jstxn.MaxGasAmount, jstxn.GasUnitPrice,
			time.Unix(jstxn.ExpirationTimestamp, 0),
		)
		return c.SubmitRawTransaction(context.TODO(), rawTxn, jstxn.SenderPriKey)
	})
	jc.submitP2PTransaction = func(rawTxn *js.Object) *js.Object {
		return promiseSubmitP2PTransaction(rawTxn)
	}

	promisePollSequenceUntil := jopher.Promisify(func(addr types.AccountAddress, seq uint64, expirationTimestamp int64) error {
		return c.PollSequenceUntil(context.TODO(), addr, seq, time.Unix(expirationTimestamp, 0))
	})
	jc.pollSequenceUntil = func(addr types.AccountAddress, seq uint64, expirationTimestamp int64) *js.Object {
		return promisePollSequenceUntil(addr, seq, expirationTimestamp)
	}

	promiseQueryTransactionByAccountSeq := jopher.Promisify(func(addr types.AccountAddress, seq uint64, withEvents bool) (*js.Object, error) {
		txn, err := c.QueryTransactionByAccountSeq(context.TODO(), addr, seq, withEvents)
		return wrapProvenTransaction(txn), err
	})
	jc.queryTransactionByAccountSeq = func(addr types.AccountAddress, seq uint64, withEvents bool) *js.Object {
		return promiseQueryTransactionByAccountSeq(addr, seq, withEvents)
	}

	promiseQueryTransactionRange := jopher.Promisify(func(start, limit uint64, withEvents bool) (*js.Object, error) {
		txnList, err := c.QueryTransactionRange(context.TODO(), start, limit, withEvents)
		return wrapProvenTransactionList(txnList), err
	})
	jc.queryTransactionRange = func(start, limit uint64, withEvents bool) *js.Object {
		return promiseQueryTransactionRange(start, limit, withEvents)
	}

	promiseQueryEventsByAccessPath := jopher.Promisify(func(addr types.AccountAddress, path []byte, start uint64, ascending bool, limit uint64) ([]*js.Object, error) {
		eventList, err := c.QueryEventsByAccessPath(context.TODO(),
			&types.AccessPath{
				Address: addr,
				Path:    path,
			},
			start, ascending, limit)
		jsevList := make([]*js.Object, 0, len(eventList))
		for _, ev := range eventList {
			jsevList = append(jsevList, wrapProvenEvent(ev))
		}
		return jsevList, err
	})
	jc.queryEventsByAccessPath = func(addr types.AccountAddress, path []byte, start uint64, ascending bool, limit uint64) *js.Object {
		return promiseQueryEventsByAccessPath(addr, path, start, ascending, limit)
	}

	return jc.Object
}
