package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/miratronix/jopher"
	"github.com/the729/go-libra/types"
)

type jsProvenLedgerInfo struct {
	*js.Object
	getVersion                    interface{} `js:"getVersion"`
	getTransactionAccumulatorHash interface{} `js:"getTransactionAccumulatorHash"`
	getEpochNum                   interface{} `js:"getEpochNum"`
	getTimestampUsec              interface{} `js:"getTimestampUsec"`
}

type jsProvenAccountState struct {
	*js.Object
	getVersion     interface{} `js:"getVersion"`
	getAccountBlob interface{} `js:"getAccountBlob"`
	isNil          interface{} `js:"isNil"`
	getLedgerInfo  interface{} `js:"getLedgerInfo"`
}

type jsProvenAccountBlob struct {
	*js.Object
	getAddress    interface{} `js:"getAddress"`
	getResource   interface{} `js:"getResource"`
	getLedgerInfo interface{} `js:"getLedgerInfo"`
}

type jsProvenAccountResource struct {
	*js.Object
	getAddress                       interface{} `js:"getAddress"`
	getBalance                       interface{} `js:"getBalance"`
	getSequenceNumber                interface{} `js:"getSequenceNumber"`
	getSentEvents                    interface{} `js:"getSentEvents"`
	getReceivedEvents                interface{} `js:"getReceivedEvents"`
	getDelegatedWithdrawalCapability interface{} `js:"getDelegatedWithdrawalCapability"`
	getEventGenerator                interface{} `js:"getEventGenerator"`
	getLedgerInfo                    interface{} `js:"getLedgerInfo"`
}

type jsProvenTransaction struct {
	*js.Object
	getMajorStatus interface{} `js:"getMajorStatus"`
	getVersion     interface{} `js:"getVersion"`
	getGasUsed     interface{} `js:"getGasUsed"`
	getWithEvents  interface{} `js:"getWithEvents"`
	getSignedTxn   interface{} `js:"getSignedTxn"`
	getEvents      interface{} `js:"getEvents"`
	getLedgerInfo  interface{} `js:"getLedgerInfo"`
}

type jsProvenTransactionList struct {
	*js.Object
	getTransactions interface{} `js:"getTransactions"`
	getLedgerInfo   interface{} `js:"getLedgerInfo"`
}

type jsProvenEvent struct {
	*js.Object
	getTransactionVersion interface{} `js:"getTransactionVersion"`
	getEventIndex         interface{} `js:"getEventIndex"`
	getEvent              interface{} `js:"getEvent"`
	getLedgerInfo         interface{} `js:"getLedgerInfo"`
}

func wrapProvenLedgerInfo(g *types.ProvenLedgerInfo) *js.Object {
	if g == nil {
		return nil
	}
	j := &jsProvenLedgerInfo{Object: js.Global.Get("Object").New()}
	j.getVersion = g.GetVersion
	j.getTransactionAccumulatorHash = g.GetTransactionAccumulatorHash
	j.getEpochNum = g.GetEpochNum
	j.getTimestampUsec = g.GetTimestampUsec
	return j.Object
}

func wrapProvenAccountState(g *types.ProvenAccountState) *js.Object {
	if g == nil {
		return nil
	}
	j := &jsProvenAccountState{Object: js.Global.Get("Object").New()}
	j.getVersion = g.GetVersion
	j.getAccountBlob = func() *js.Object {
		return wrapProvenAccountBlob(g.GetAccountBlob())
	}
	j.isNil = g.IsNil
	j.getLedgerInfo = func() *js.Object {
		return wrapProvenLedgerInfo(g.GetLedgerInfo())
	}
	return j.Object
}

func wrapProvenAccountBlob(g *types.ProvenAccountBlob) *js.Object {
	if g == nil {
		return nil
	}
	j := &jsProvenAccountBlob{Object: js.Global.Get("Object").New()}
	j.getAddress = g.GetAddress
	j.getResource = jopher.Promisify(func(path []byte) (*js.Object, error) {
		r, err := g.GetResource(path)
		if err != nil {
			return nil, err
		}
		return wrapProvenAccountResource(r), nil
	})
	j.getLedgerInfo = func() *js.Object {
		return wrapProvenLedgerInfo(g.GetLedgerInfo())
	}
	return j.Object
}

func wrapProvenAccountResource(g *types.ProvenAccountResource) *js.Object {
	if g == nil {
		return nil
	}
	j := &jsProvenAccountResource{Object: js.Global.Get("Object").New()}
	j.getAddress = g.GetAddress
	j.getBalance = g.GetBalance
	j.getSequenceNumber = g.GetSequenceNumber
	j.getSentEvents = g.GetSentEvents
	j.getReceivedEvents = g.GetReceivedEvents
	j.getDelegatedWithdrawalCapability = g.GetDelegatedWithdrawalCapability
	j.getEventGenerator = g.GetEventGenerator
	j.getLedgerInfo = func() *js.Object {
		return wrapProvenLedgerInfo(g.GetLedgerInfo())
	}
	return j.Object
}

func wrapProvenTransaction(g *types.ProvenTransaction) *js.Object {
	if g == nil {
		return nil
	}
	j := &jsProvenTransaction{Object: js.Global.Get("Object").New()}
	j.getMajorStatus = g.GetMajorStatus
	j.getGasUsed = g.GetGasUsed
	j.getVersion = g.GetVersion
	j.getWithEvents = g.GetWithEvents
	j.getEvents = g.GetEvents
	j.getSignedTxn = g.GetSignedTxn
	j.getLedgerInfo = func() *js.Object {
		return wrapProvenLedgerInfo(g.GetLedgerInfo())
	}
	return j.Object
}

func wrapProvenTransactionList(g *types.ProvenTransactionList) *js.Object {
	if g == nil {
		return nil
	}
	j := &jsProvenTransactionList{Object: js.Global.Get("Object").New()}
	j.getTransactions = func() []*js.Object {
		r := make([]*js.Object, 0, 0)
		for _, txn := range g.GetTransactions() {
			r = append(r, wrapProvenTransaction(txn))
		}
		return r
	}
	j.getLedgerInfo = func() *js.Object {
		return wrapProvenLedgerInfo(g.GetLedgerInfo())
	}
	return j.Object
}

func wrapProvenEvent(g *types.ProvenEvent) *js.Object {
	if g == nil {
		return nil
	}
	j := &jsProvenEvent{Object: js.Global.Get("Object").New()}
	j.getTransactionVersion = g.GetTransactionVersion
	j.getEventIndex = g.GetEventIndex
	j.getEvent = g.GetEvent
	j.getLedgerInfo = func() *js.Object {
		return wrapProvenLedgerInfo(g.GetLedgerInfo())
	}
	return j.Object
}
