package main

import (
	"github.com/gopherjs/gopherjs/js"

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

func newClient(server, waypoint string) *js.Object {
	c, err := client.New(server, waypoint)
	if err != nil {
		panic(err)
	}
	return wrapClientObject(c)
}
