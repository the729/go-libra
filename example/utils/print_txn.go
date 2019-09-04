package utils

import (
	"encoding/binary"
	"encoding/hex"
	"log"

	"github.com/the729/go-libra/generated/pbtypes"
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
	switch rawTxn.GetPayload().(type) {
	case *pbtypes.RawTransaction_Program:
		log.Printf("        Payload is Program ...")
	case *pbtypes.RawTransaction_WriteSet:
		log.Printf("        Payload is WriteSet ...")
		return
	case *pbtypes.RawTransaction_Script:
		log.Printf("        Payload is Script ...")
		return
	case *pbtypes.RawTransaction_Module:
		log.Printf("        Payload is Module ...")
		return
	}
	log.Printf("        Program: %v...", hex.EncodeToString(rawTxn.GetProgram().Code[:30]))
	progName := stdscript.InferProgramName(rawTxn.GetProgram().Code)
	log.Printf("            (program name: %s)", progName)
	switch progName {
	case "peer_to_peer_transfer", "mint":
		log.Printf("        Arg 0 (receiver address): %v", hex.EncodeToString(rawTxn.GetProgram().Arguments[0].Data))
		log.Printf("        Arg 1 (amount microLibra): %v", binary.LittleEndian.Uint64(rawTxn.GetProgram().Arguments[1].Data))
	default:
		for i, d := range rawTxn.GetProgram().Arguments {
			log.Printf("        Arg %d: %v", i, hex.EncodeToString(d.Data))
		}
	}
	log.Printf("    Max gas amount (gas units): %v", rawTxn.GetMaxGasAmount())
	log.Printf("    Gas unit price (microLibra/unit): %v", rawTxn.GetGasUnitPrice())
	log.Printf("    Expiration timestamp: %v", rawTxn.GetExpirationTime())
	log.Printf("    Gas used (microLibra): %v", txn.GetGasUsed())
	if txn.GetWithEvents() {
		log.Printf("    Events: (%d total)", len(txn.GetEvents()))
		for idx, ev := range txn.GetEvents() {
			log.Printf("      #%d:", idx)
			log.Printf("        Seq #%d", ev.SequenceNumber)
			log.Printf("        Key: %v", hex.EncodeToString(ev.Key))
			// if dp, err := ev.AccessPath.DecodePath(); err != nil {
			// 	log.Printf("        Raw path: %v (can not decode: %v)", hex.EncodeToString(ev.AccessPath.Path), err)
			// 	log.Printf("        Raw data: %v", hex.EncodeToString(ev.Data))
			// } else {
			// 	if tagName := types.InferPathTagName(dp.Tag); tagName == "0x0.LibraAccount.T" {
			// 		log.Printf("        Path: %v/%v", tagName, strings.Join(dp.Accesses, "/"))
			// 		if dp.IsEqual(types.AccountSentEventPath()) {
			// 			log.Printf("            (Event is: sent payment)")
			// 			log.Printf("        Data:")
			// 			evBody := &stdscript.SentPaymentEvent{}
			// 			if err := evBody.UnmarshalBinary(ev.Data); err != nil {
			// 				log.Printf("        Raw data: %v (cannot decode: %v)", hex.EncodeToString(ev.Data), err)
			// 			} else {
			// 				log.Printf("            Payee: %v", hex.EncodeToString(evBody.Payee))
			// 				log.Printf("            Amount: %v", float64(evBody.Amount)/1000000.0)
			// 			}
			// 		} else if dp.IsEqual(types.AccountReceivedEventPath()) {
			// 			log.Printf("            (Event is: received payment)")
			// 			log.Printf("        Data:")
			// 			evBody := &stdscript.ReceivedPaymentEvent{}
			// 			if err := evBody.UnmarshalBinary(ev.Data); err != nil {
			// 				log.Printf("        Raw data: %v (cannot decode: %v)", hex.EncodeToString(ev.Data), err)
			// 			} else {
			// 				log.Printf("            Payer: %v", hex.EncodeToString(evBody.Payer))
			// 				log.Printf("            Amount: %v", float64(evBody.Amount)/1000000.0)
			// 			}
			// 		}
			// 	} else {
			// 		log.Printf("        Path: %v+%v/%v", dp.Tag.TypePrefix(), hex.EncodeToString(dp.Tag.Hash()), strings.Join(dp.Accesses, "/"))
			// 	}
			// }
		}
	} else {
		log.Printf("    Events not present")
	}
}
