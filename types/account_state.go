package types

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types/proof"
)

type AccountState struct {
	Version uint64
	RawBlob RawAccountBlob
}

type AccountStateWithProof struct {
	*AccountState
	Proof *AccountStateProof
}

type AccountStateProof struct {
	*TransactionInfo
	LedgerInfoToTransactionInfoProof *proof.Accumulator
	TransactionInfoToAccountProof    *proof.SparseMerkle
}

type ProvenAccountState struct {
	proven       bool
	accountState AccountState
	addr         AccountAddress
	ledgerInfo   *ProvenLedgerInfo
}

func (a *AccountStateWithProof) FromProtoResponse(pb *pbtypes.GetAccountStateResponse) error {
	if pb == nil {
		return ErrNilInput
	}
	return a.FromProto(pb.AccountStateWithProof)
}

func (a *AccountStateWithProof) FromProto(pb *pbtypes.AccountStateWithProof) error {
	if pb == nil {
		return ErrNilInput
	}
	a.AccountState = &AccountState{}
	a.Version = pb.Version
	if pb.Blob != nil {
		a.RawBlob = RawAccountBlob(pb.Blob.Blob)
	}
	a.Proof = &AccountStateProof{}
	return a.Proof.FromProto(pb.Proof)
}

func (ap *AccountStateProof) FromProto(pb *pbtypes.AccountStateProof) error {
	var err error
	if pb == nil {
		return ErrNilInput
	}

	ap.LedgerInfoToTransactionInfoProof = &proof.Accumulator{}
	err = ap.LedgerInfoToTransactionInfoProof.FromProto(pb.LedgerInfoToTransactionInfoProof)
	if err != nil {
		return err
	}
	ap.TransactionInfo = &TransactionInfo{}
	err = ap.TransactionInfo.FromProto(pb.TransactionInfo)
	if err != nil {
		return err
	}
	ap.TransactionInfoToAccountProof = &proof.SparseMerkle{}
	err = ap.TransactionInfoToAccountProof.FromProto(pb.TransactionInfoToAccountProof)
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountStateWithProof) Verify(addr AccountAddress, provenLedgerInfo *ProvenLedgerInfo) (*ProvenAccountState, error) {
	addrHash := addr.Hash()
	blobHash := a.RawBlob.Hash()

	var err error
	if blobHash == nil {
		err = a.Proof.TransactionInfoToAccountProof.VerifyNonInclusion(
			addrHash,
			a.Proof.TransactionInfo.stateRootHash,
		)
	} else {
		err = a.Proof.TransactionInfoToAccountProof.VerifyInclusion(
			&proof.LeafNode{addrHash, blobHash},
			a.Proof.TransactionInfo.stateRootHash,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot verify account state from transaction info: %v", err)
	}

	if a.Version > provenLedgerInfo.GetVersion() {
		return nil, errors.New("account version > ledger version")
	}

	err = a.Proof.LedgerInfoToTransactionInfoProof.Verify(
		a.Version, a.Proof.TransactionInfo.Hash(),
		provenLedgerInfo.GetTransactionAccumulatorHash(),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot verify transaction info from ledger info: %v", err)
	}

	return &ProvenAccountState{
		proven: true,
		accountState: AccountState{
			Version: a.Version,
			RawBlob: cloneBytes(a.RawBlob),
		},
		addr:       cloneBytes(addr),
		ledgerInfo: provenLedgerInfo,
	}, nil
}

func (pas *ProvenAccountState) GetLedgerInfo() *ProvenLedgerInfo {
	if !pas.proven {
		panic("not valid proven account state")
	}
	return pas.ledgerInfo
}

func (pas *ProvenAccountState) GetVersion() uint64 {
	if !pas.proven {
		panic("not valid proven account state")
	}
	return pas.accountState.Version
}

func (pas *ProvenAccountState) GetAddress() AccountAddress {
	if !pas.proven {
		panic("not valid proven account state")
	}
	return AccountAddress(cloneBytes(pas.addr))
}

func (pas *ProvenAccountState) GetAccountBlob() *ProvenAccountBlob {
	if !pas.proven {
		panic("not valid proven account state")
	}
	pab := &ProvenAccountBlob{
		proven: true,
		addr:   cloneBytes(pas.addr),
	}
	pab.accountBlob.Raw = cloneBytes(pas.accountState.RawBlob)
	pab.accountBlob.ParseToMap()
	return pab
}

func (pas *ProvenAccountState) IsNil() bool {
	if !pas.proven {
		panic("not valid proven account state")
	}
	return pas.accountState.RawBlob == nil
}
