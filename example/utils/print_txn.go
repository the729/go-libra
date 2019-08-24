package utils

import (
	"encoding/hex"
	"log"
	"strings"

	"github.com/the729/go-libra/language/stdscript"
	"github.com/the729/go-libra/types"
)

// PrintTxn prints a proven transaction, using standard logger
func PrintTxn(txn *types.ProvenTransaction) {
	log.Printf("Txn #%d:", txn.GetVersion())
	rawTxn, _ := txn.GetSignedTxn().UnmarshalRawTransaction()
	log.Printf("    Raw txn:")
	log.Printf("        Sender account: %v", hex.EncodeToString(rawTxn.SenderAccount))
	log.Printf("        Sender seq #%v", rawTxn.SequenceNumber)
	log.Printf("        Program: %v...", hex.EncodeToString(rawTxn.GetProgram().Code[:30]))
	log.Printf("            (program name: %s)", stdscript.InferProgramName(rawTxn.GetProgram().Code))
	log.Printf("        Arg 0: %v", hex.EncodeToString(rawTxn.GetProgram().Arguments[0].Data))
	log.Printf("        Arg 1: %v", hex.EncodeToString(rawTxn.GetProgram().Arguments[1].Data))
	log.Printf("    Gas used: %v", txn.GetGasUsed())
	if txn.GetWithEvents() {
		log.Printf("    Events: (%d total)", len(txn.GetEvents()))
		for idx, ev := range txn.GetEvents() {
			log.Printf("      #%d:", idx)
			log.Printf("        Seq #%d", ev.SequenceNumber)
			log.Printf("        Addr: %v", hex.EncodeToString(ev.AccessPath.Address))
			if dp, err := ev.AccessPath.DecodePath(); err != nil {
				log.Printf("        Raw path: %v (can not decode: %v)", hex.EncodeToString(ev.AccessPath.Path), err)
				log.Printf("        Raw data: %v", hex.EncodeToString(ev.Data))
			} else {
				if tagName := types.InferPathTagName(dp.Tag); tagName == "0x0.LibraAccount.T" {
					log.Printf("        Path: %v/%v", tagName, strings.Join(dp.Accesses, "/"))
					if dp.IsEqual(types.AccountSentEventPath()) {
						log.Printf("            (Event is: sent payment)")
						log.Printf("        Data:")
						evBody := &stdscript.SentPaymentEvent{}
						if err := evBody.UnmarshalBinary(ev.Data); err != nil {
							log.Printf("        Raw data: %v (cannot decode: %v)", hex.EncodeToString(ev.Data), err)
						} else {
							log.Printf("            Payee: %v", hex.EncodeToString(evBody.Payee))
							log.Printf("            Amount: %v", float64(evBody.Amount)/1000000.0)
						}
					} else if dp.IsEqual(types.AccountReceivedEventPath()) {
						log.Printf("            (Event is: received payment)")
						log.Printf("        Data:")
						evBody := &stdscript.ReceivedPaymentEvent{}
						if err := evBody.UnmarshalBinary(ev.Data); err != nil {
							log.Printf("        Raw data: %v (cannot decode: %v)", hex.EncodeToString(ev.Data), err)
						} else {
							log.Printf("            Payer: %v", hex.EncodeToString(evBody.Payer))
							log.Printf("            Amount: %v", float64(evBody.Amount)/1000000.0)
						}
					}
				} else {
					log.Printf("        Path: %v+%v/%v", dp.Tag.TypePrefix(), hex.EncodeToString(dp.Tag.Hash()), strings.Join(dp.Accesses, "/"))
				}
			}
		}
	} else {
		log.Printf("    Events not present")
	}
}
