package main

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"github.com/miratronix/jopher"
	"github.com/the729/go-libra/client"
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
	getAddress              interface{} `js:"getAddress"`
	getResource             interface{} `js:"getResource"`
	getResourcePaths        interface{} `js:"getResourcePaths"`
	getLibraAccountResource interface{} `js:"getLibraAccountResource"`
	getLedgerInfo           interface{} `js:"getLedgerInfo"`
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
	getMajorStatus   interface{} `js:"getMajorStatus"`
	getVersion       interface{} `js:"getVersion"`
	getGasUsed       interface{} `js:"getGasUsed"`
	getWithEvents    interface{} `js:"getWithEvents"`
	getSignedTxn     interface{} `js:"getSignedTxn"`
	getBlockMetadata interface{} `js:"getBlockMetadata"`
	getEvents        interface{} `js:"getEvents"`
	getLedgerInfo    interface{} `js:"getLedgerInfo"`
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

type jsValidator struct {
	*js.Object
	address         string `js:"addr"`
	consensusPubkey string `js:"c"`
	votingPower     uint64 `js:"power"`
}

type jsClientState struct {
	*js.Object
	waypoint     string       `js:"waypoint"`
	epoch        uint64       `js:"epoch"`
	validatorSet []*js.Object `js:"validator_set"`
	knownVersion uint64       `js:"known_version"`
	subtrees     []string     `js:"subtrees"`
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
	j.getLibraAccountResource = jopher.Promisify(func() (*js.Object, error) {
		r, err := g.GetLibraAccountResource()
		if err != nil {
			return nil, err
		}
		return wrapProvenAccountResource(r), nil
	})
	j.getResource = func(path []byte) []byte {
		r, _ := g.GetResource(path)
		return r
	}
	j.getResourcePaths = func() [][]byte {
		return g.GetResourcePaths()
	}
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
	j.getBlockMetadata = g.GetBlockMetadata
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

func wrapClientState(cs *client.ClientState) *js.Object {
	if cs == nil {
		return nil
	}
	j := &jsClientState{Object: js.Global.Get("Object").New()}
	j.waypoint = cs.Waypoint
	j.epoch = cs.Epoch
	j.knownVersion = cs.KnownVersion
	jsvs := make([]*js.Object, 0, len(cs.ValidatorSet))
	for _, v := range cs.ValidatorSet {
		jsv := &jsValidator{Object: js.Global.Get("Object").New()}
		jsv.address = hex.EncodeToString(v.AccountAddress[:])
		jsv.consensusPubkey = hex.EncodeToString(v.ConsensusPubkey)
		jsv.votingPower = v.ConsensusVotingPower

		jsvs = append(jsvs, jsv.Object)
	}
	j.validatorSet = jsvs
	jsSubtrees := make([]string, 0, len(cs.Subtrees))
	for _, t := range cs.Subtrees {
		jsSubtrees = append(jsSubtrees, hex.EncodeToString(t))
	}
	j.subtrees = jsSubtrees
	return j.Object
}

func unwrapClientState(csObj *js.Object) (*client.ClientState, error) {
	j := &jsClientState{Object: csObj}
	cs := &client.ClientState{
		Waypoint:     j.waypoint,
		Epoch:        j.epoch,
		KnownVersion: j.knownVersion,
	}
	if j.Get("validatorSet") != js.Undefined {
		for i, vObj := range j.validatorSet {
			jv := &jsValidator{Object: vObj}
			var err error
			v := &types.ValidatorPublicKeys{}
			v.ConsensusPubkey, err = hex.DecodeString(jv.consensusPubkey)
			if err != nil {
				return nil, fmt.Errorf("unable to decode consensus pubkey #%d: %v", i, err)
			}
			if l, err := hex.Decode(v.AccountAddress[:], []byte(jv.address)); l != 32 || err != nil {
				if err == nil {
					err = errors.New("wrong length")
				}
				return nil, fmt.Errorf("unable to decode address #%d: %v", i, err)
			}
			v.ConsensusVotingPower = jv.votingPower
			cs.ValidatorSet = append(cs.ValidatorSet, v)
		}
	}
	if j.Get("subtrees") != js.Undefined {
		for i, t := range j.subtrees {
			subtree, err := hex.DecodeString(t)
			if err != nil {
				return nil, fmt.Errorf("unable to decode subtree #%d: %v", i, err)
			}
			cs.Subtrees = append(cs.Subtrees, subtree)
		}
	}
	return cs, nil
}
