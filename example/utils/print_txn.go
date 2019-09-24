package utils

import (
	"encoding/hex"
	"log"

	"github.com/the729/go-libra/language/stdscript"
	"github.com/the729/go-libra/types"
)

// PrintTxn prints a proven transaction, using standard logger
func PrintTxn(txn *types.ProvenTransaction) {
	log.Printf("Txn #%d:", txn.GetVersion())
	rawTxn := txn.GetSignedTxn().RawTxn
	log.Printf("    Raw txn:")
	log.Printf("        Sender account: %v", hex.EncodeToString(rawTxn.Sender))
	log.Printf("        Sender seq #%v", rawTxn.SequenceNumber)
	switch pld := rawTxn.Payload.(type) {
	case *types.TxnPayloadProgram:
		log.Printf("        Payload is Program.")
		return
	case types.TxnPayloadWriteSet:
		log.Printf("        Payload is WriteSet.")
		return
	case *types.TxnPayloadScript:
		log.Printf("        Payload is Script ...")
		log.Printf("        Program: %v...", hex.EncodeToString(pld.Code[:30]))
		progName := stdscript.InferProgramName(pld.Code)
		log.Printf("            (program name: %s)", progName)
		for i, arg := range pld.Args {
			switch arg := arg.(type) {
			case types.TxnArgU64:
				log.Printf("        Arg %d: u64  (%v)", i, arg)
			case types.TxnArgAddress:
				log.Printf("        Arg %d: addr (%v)", i, hex.EncodeToString(arg))
			case types.TxnArgString:
				log.Printf("        Arg %d: str  (%v)", i, arg)
			case types.TxnArgBytes:
				log.Printf("        Arg %d: bytes(%v)", i, hex.EncodeToString(arg))
			}
		}
	case types.TxnPayloadModule:
		log.Printf("        Payload is Module.")
		return
	}
	log.Printf("    Max gas amount (gas units): %v", rawTxn.MaxGasAmount)
	log.Printf("    Gas unit price (microLibra/unit): %v", rawTxn.GasUnitPrice)
	log.Printf("    Expiration timestamp: %v", rawTxn.ExpirationTime)
	log.Printf("    Gas used (microLibra): %v", txn.GetGasUsed())
	log.Printf("    Major status: %d - %s", txn.GetMajorStatus(), txn.GetMajorStatus())
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
