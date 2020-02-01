package main

import (
	"context"
	"encoding/binary"
	"errors"
	"reflect"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/miratronix/jopher"

	"github.com/the729/go-libra/client"
	"github.com/the729/go-libra/language/stdscript"
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
		"resourcePath":             types.ResourcePath,
		"accountResourcePath":      types.AccountResourcePath,
		"accountSentEventPath":     types.AccountSentEventPath,
		"accountReceivedEventPath": types.AccountReceivedEventPath,
		"pubkeyToAddress":          client.PubkeyMustToAddress,
		"inferProgramName":         stdscript.InferProgramName,
	})
}

type jsClient struct {
	*js.Object

	queryLedgerInfo              func(...interface{}) *js.Object                                     `js:"queryLedgerInfo"`
	queryAccountState            func(types.AccountAddress) *js.Object                               `js:"queryAccountState"`
	queryAccountSequenceNumber   func(types.AccountAddress) *js.Object                               `js:"queryAccountSequenceNumber"`
	submitP2PTransaction         func(*js.Object) *js.Object                                         `js:"submitP2PTransaction"`
	submitRawTransaction         func(*js.Object) *js.Object                                         `js:"submitRawTransaction"`
	pollSequenceUntil            func(types.AccountAddress, uint64, int64) *js.Object                `js:"pollSequenceUntil"`
	queryTransactionByAccountSeq func(types.AccountAddress, uint64, bool) *js.Object                 `js:"queryTransactionByAccountSeq"`
	queryTransactionRange        func(uint64, uint64, bool) *js.Object                               `js:"queryTransactionRange"`
	queryEventsByAccessPath      func(types.AccountAddress, []byte, uint64, bool, uint64) *js.Object `js:"queryEventsByAccessPath"`
}

func newClient(server, waypoint string) *js.Object {
	c, err := client.New(server, waypoint)
	if err != nil {
		panic(err)
	}
	jc := jsClient{Object: js.Global.Get("Object").New()}
	jc.queryLedgerInfo = jopher.Promisify(func() (*js.Object, error) {
		r, err := c.QueryLedgerInfo(context.TODO())
		return wrapProvenLedgerInfo(r), err
	})

	promiseQueryAccountState := jopher.Promisify(func(addr types.AccountAddress) (*js.Object, error) {
		r, err := c.QueryAccountState(context.TODO(), addr)
		return wrapProvenAccountState(r), err
	})
	jc.queryAccountState = func(addr types.AccountAddress) *js.Object {
		return promiseQueryAccountState(addr)
	}

	promiseQueryAccountSequenceNumber := jopher.Promisify(func(addr types.AccountAddress) (uint64, error) {
		return c.QueryAccountSequenceNumber(context.TODO(), addr)
	})
	jc.queryAccountSequenceNumber = func(addr types.AccountAddress) *js.Object {
		return promiseQueryAccountSequenceNumber(addr)
	}

	promiseSubmitP2PTransaction := jopher.Promisify(func(txn *js.Object) (uint64, error) {
		type jsP2PTxn struct {
			*js.Object
			SenderAddr          [32]byte `js:"senderAddr"`
			SenderPriKey        []byte   `js:"senderPrivateKey"`
			RecvAddr            [32]byte `js:"recvAddr"`
			SenderSeq           uint64   `js:"senderSeq"`
			AmountMicro         uint64   `js:"amountMicro"`
			MaxGasAmount        uint64   `js:"maxGasAmount"`
			GasUnitPrice        uint64   `js:"gasUnitPrice"`
			ExpirationTimestamp int64    `js:"expirationTimestamp"`
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

	promiseSubmitRawTransaction := jopher.Promisify(func(txn *js.Object) (uint64, error) {
		type jsRawTxn struct {
			*js.Object
			SenderAddr   [32]byte `js:"senderAddr"`
			SenderPriKey []byte   `js:"senderPrivateKey"`
			SenderSeq    uint64   `js:"senderSeq"`
			Payload      *struct {
				*js.Object
				Code   []byte        `js:"code"`
				Args   []interface{} `js:"args"`
				Module []byte        `js:"module"`
			} `js:"payload"`
			MaxGasAmount        uint64 `js:"maxGasAmount"`
			GasUnitPrice        uint64 `js:"gasUnitPrice"`
			ExpirationTimestamp uint64 `js:"expirationTimestamp"`
		}
		jstxn := &jsRawTxn{Object: txn}
		if jstxn.Payload.Object == js.Undefined {
			return 0, errors.New("payload cannot be nil")
		}
		rawTxn := &types.RawTransaction{
			Sender:         jstxn.SenderAddr,
			SequenceNumber: jstxn.SenderSeq,
			MaxGasAmount:   jstxn.MaxGasAmount,
			GasUnitPrice:   jstxn.GasUnitPrice,
			ExpirationTime: jstxn.ExpirationTimestamp,
		}
		// Probably because of https://github.com/gopherjs/gopherjs/issues/460,
		// null array in js cannot be converted to nil slice in go, and accessing it
		// will yield errors in js.
		if jstxn.Payload.Get("code") != js.Undefined && jstxn.Payload.Get("module") != js.Undefined {
			return 0, errors.New("module and code cannot be both non-nil")
		}
		if jstxn.Payload.Get("code") != js.Undefined {
			var jsArgs []interface{}
			if jstxn.Payload.Get("args") != js.Undefined {
				jsArgs = jstxn.Payload.Args
			}
			payload := &types.TxnPayloadScript{
				Code: jstxn.Payload.Code,
				Args: make([]types.TransactionArgument, 0, len(jsArgs)),
			}
			for _, arg := range jsArgs {
				var arg1 types.TransactionArgument
				switch v := arg.(type) {
				case bool:
					arg1 = types.TxnArgBool(v)
				case float64:
					arg1 = types.TxnArgU64(v)
				case []uint8:
					if len(v) == 32 {
						a := types.TxnArgAddress{}
						copy(a[:], v)
						arg1 = types.TxnArgAddress(a)
					} else {
						arg1 = types.TxnArgBytes(v)
					}
				case string:
					arg1 = types.TxnArgString(v)
				case map[string]interface{}:
					typ, ok1 := v["type"].(string)
					val, ok2 := v["value"].([]uint8)
					if !ok1 || !ok2 {
						return 0, errors.New("invalid argument type or value")
					}
					switch typ {
					case "uint64":
						if len(val) != 8 {
							return 0, errors.New("invalid length for uint64 argument")
						}
						arg1 = types.TxnArgU64(binary.LittleEndian.Uint64(val))
					case "bytes":
						arg1 = types.TxnArgBytes(val)
					default:
						return 0, errors.New("transaction argument explicit type not supported: " + typ)
					}
				default:
					return 0, errors.New("transaction argument type not supported: " + reflect.TypeOf(v).String())
				}
				payload.Args = append(payload.Args, arg1)
			}
			rawTxn.Payload = payload
		} else {
			rawTxn.Payload = types.TxnPayloadModule(jstxn.Payload.Module)
		}
		return c.SubmitRawTransaction(context.TODO(), rawTxn, jstxn.SenderPriKey)
	})
	jc.submitRawTransaction = func(rawTxn *js.Object) *js.Object {
		return promiseSubmitRawTransaction(rawTxn)
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
