package utils

import (
	"encoding/hex"
	"log"

	"github.com/the729/go-libra/language/stdscript"
	"github.com/the729/go-libra/types"
	"github.com/the729/lcs"
)

// PrintTxn prints a proven transaction, using standard logger
func PrintTxn(txn *types.ProvenTransaction) {
	log.Printf("Txn #%d:", txn.GetVersion())
	if txn.GetSignedTxn() == nil {
		log.Printf("    is not a user transaction.")
		return
	}
	rawTxn := txn.GetSignedTxn().RawTxn
	log.Printf("    Raw txn: hash=%x", txn.GetHash())
	log.Printf("        Sender account: %v", hex.EncodeToString(rawTxn.Sender[:]))
	log.Printf("        Sender seq #%v", rawTxn.SequenceNumber)
	switch pld := rawTxn.Payload.(type) {
	case *types.TxnPayloadWriteSet:
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
				log.Printf("        Arg %d: addr (%v)", i, hex.EncodeToString(arg[:]))
			case types.TxnArgBool:
				log.Printf("        Arg %d: bool (%v)", i, arg)
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
	log.Printf("    Gas specifier: %+v", rawTxn.GasSpecifier)
	log.Printf("    Major status: %d - %s", txn.GetMajorStatus(), txn.GetMajorStatus())
	if txn.GetWithEvents() {
		log.Printf("    Events: (%d total)", len(txn.GetEvents()))
		for idx, ev := range txn.GetEvents() {
			ev0 := ev.Value.(*types.ContractEventV0)
			log.Printf("      #%d:", idx)
			log.Printf("        Key: %v", hex.EncodeToString(ev0.Key))
			log.Printf("        Seq #%d", ev0.SequenceNumber)
			if len(ev0.Data) > 44 {
				log.Printf("        Raw event: %s ...(len=%d)", hex.EncodeToString(ev0.Data[:44]), len(ev0.Data))
			} else {
				log.Printf("        Raw event: %s", hex.EncodeToString(ev0.Data))
			}
			pev := &stdscript.PaymentEvent{}
			if err := lcs.Unmarshal(ev0.Data, pev); err != nil {
				log.Printf("            (Unknown event type)")
			} else {
				log.Printf("            Amount (microLibra): %d", pev.Amount)
				log.Printf("            Opponent address: %s", hex.EncodeToString(pev.Address[:]))
			}
		}
	} else {
		log.Printf("    Events not present")
	}
}
