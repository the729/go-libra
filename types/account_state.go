package types

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types/proof"
)

type AccountStateWithProof struct {
	Version uint64
	Blob    *AccountBlob
	Proof   *AccountStateProof
}

type AccountStateProof struct {
	ledgerInfoToTransactionInfoProof *proof.Accumulator
	transactionInfo                  *TransactionInfo
	transactionInfoToAccountProof    *proof.SparseMerkle
}

func (a *AccountStateWithProof) FromProtoResponse(pb *pbtypes.GetAccountStateResponse) error {
	if pb.AccountStateWithProof == nil {
		return errors.New("nil pb.AccountStateWithProof")
	}
	return a.FromProto(pb.AccountStateWithProof)
}

func (a *AccountStateWithProof) FromProto(pb *pbtypes.AccountStateWithProof) error {
	a.Version = pb.Version
	if pb.Blob != nil {
		a.Blob = &AccountBlob{Raw: pb.Blob.Blob}
	}
	a.Proof = &AccountStateProof{}
	return a.Proof.FromProto(pb.Proof)
}

func (ap *AccountStateProof) FromProto(pb *pbtypes.AccountStateProof) error {
	var err error

	ap.ledgerInfoToTransactionInfoProof = &proof.Accumulator{}
	err = ap.ledgerInfoToTransactionInfoProof.FromProto(pb.LedgerInfoToTransactionInfoProof)
	if err != nil {
		return err
	}
	ap.transactionInfo = &TransactionInfo{}
	err = ap.transactionInfo.FromProto(pb.TransactionInfo)
	if err != nil {
		return err
	}
	ap.transactionInfoToAccountProof = &proof.SparseMerkle{}
	err = ap.transactionInfoToAccountProof.FromProto(pb.TransactionInfoToAccountProof)
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountStateWithProof) Verify(addr AccountAddress, ledgerInfo *LedgerInfo) error {
	addrHash := addr.Hash()
	blobHash := a.Blob.Hash()

	var err error
	if blobHash == nil {
		err = a.Proof.transactionInfoToAccountProof.VerifyNonInclusion(
			addrHash,
			a.Proof.transactionInfo.stateRootHash,
		)
	} else {
		err = a.Proof.transactionInfoToAccountProof.VerifyInclusion(
			&proof.LeafNode{addrHash, blobHash},
			a.Proof.transactionInfo.stateRootHash,
		)
	}
	if err != nil {
		return fmt.Errorf("cannot verify account state from transaction info: %v", err)
	}

	if a.Version > ledgerInfo.Version {
		return errors.New("account version > ledger version")
	}

	err = a.Proof.ledgerInfoToTransactionInfoProof.Verify(
		a.Version, a.Proof.transactionInfo.Hash(),
		ledgerInfo.transactionAccumulatorHash,
	)
	if err != nil {
		return fmt.Errorf("cannot verify transaction info from ledger info: %v", err)
	}

	return nil
}
