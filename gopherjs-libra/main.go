package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/miratronix/jopher"

	"github.com/the729/go-libra/client"
	"github.com/the729/go-libra/language/stdscript"
	"github.com/the729/go-libra/types"
)

var (
	jsTypeOf *js.Object
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
		"balanceResourcePath":      types.BalanceResourcePath,
		"accountSentEventPath":     types.AccountSentEventPath,
		"accountReceivedEventPath": types.AccountReceivedEventPath,
		"pubkeyToAddress":          client.PubkeyMustToAddress,
		"pubkeyToAuthKey":          client.PubkeyMustToAuthKey,
		"inferProgramName":         stdscript.InferProgramName,
	})
	jsTypeOf = js.Global.Call("eval", `(function(x){return typeof(x);})`)
}

func newClient(server string, state *js.Object) *js.Object {
	t := jsTypeOf.Invoke(state).String()
	var c *client.Client
	var err error
	switch t {
	case "string":
		c, err = client.New(server, state.String())
		jopher.ThrowOnError(err)
	case "object":
		cs, err := unwrapClientState(state)
		jopher.ThrowOnError(err)
		c, err = client.NewFromState(server, cs)
		jopher.ThrowOnError(err)
	default:
		panic("state must be string or object")
	}
	return wrapClientObject(c)
}
