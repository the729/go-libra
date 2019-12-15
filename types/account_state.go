package types

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/go-libra/generated/pbtypes"
	"github.com/the729/go-libra/types/proof"
)

// AccountState is an account state.
type AccountState struct {
	Version uint64
	RawBlob RawAccountBlob
}

// AccountStateWithProof is an account state with proof.
type AccountStateWithProof struct {
	*AccountState
	Proof *AccountStateProof
}

// AccountStateProof is a chain of proof that a certain account state is included
// in the ledger, or the account does not exist.
type AccountStateProof struct {
	// LedgerInfoToTransactionInfoProof is a Merkle Tree accumulator to prove that TransactionInfo
	// is included in the ledger.
	LedgerInfoToTransactionInfoProof *proof.Accumulator

	// TransactionInfo is the info of the transaction that leads to this version of the ledger.
	*TransactionInfo

	// TransactionInfoToAccountProof is a Sparse Merkle Tree proof that the account state is part of
	// the whole ledger state.
	TransactionInfoToAccountProof *proof.SparseMerkle
}

// ProvenAccountState is an account state proven to be equal to the state of the ledger.
//
// Either the account is included in the ledger, or it does not exist.
type ProvenAccountState struct {
	proven       bool
	accountState AccountState
	addr         AccountAddress
	ledgerInfo   *ProvenLedgerInfo
}

// FromProtoResponse parses a protobuf struct into this struct.
func (a *AccountStateWithProof) FromProtoResponse(pb *pbtypes.GetAccountStateResponse) error {
	if pb == nil {
		return ErrNilInput
	}
	return a.FromProto(pb.AccountStateWithProof)
}

// FromProto parses a protobuf struct into this struct.
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

// FromProto parses a protobuf struct into this struct.
func (ap *AccountStateProof) FromProto(pb *pbtypes.AccountStateProof) error {
	var err error
	if pb == nil {
		return ErrNilInput
	}

	ap.LedgerInfoToTransactionInfoProof = &proof.Accumulator{Hasher: sha3libra.NewTransactionAccumulator()}
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

// Verify the proof of the account state, and output a ProvenAccountState if successful.
func (a *AccountStateWithProof) Verify(addr AccountAddress, provenLedgerInfo *ProvenLedgerInfo) (*ProvenAccountState, error) {
	addrHash := addr.Hash()
	blobHash := a.RawBlob.Hash()

	var err error
	if blobHash == nil {
		err = a.Proof.TransactionInfoToAccountProof.VerifyNonInclusion(
			addrHash,
			a.Proof.TransactionInfo.StateRootHash,
		)
	} else {
		err = a.Proof.TransactionInfoToAccountProof.VerifyInclusion(
			&proof.LeafNode{Key: addrHash, ValueHash: blobHash},
			a.Proof.TransactionInfo.StateRootHash,
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
		addr:       addr,
		ledgerInfo: provenLedgerInfo,
	}, nil
}

// GetLedgerInfo returns the ledger info.
func (pas *ProvenAccountState) GetLedgerInfo() *ProvenLedgerInfo {
	if !pas.proven {
		panic("not valid proven account state")
	}
	return pas.ledgerInfo
}

// GetVersion returns the version.
func (pas *ProvenAccountState) GetVersion() uint64 {
	if !pas.proven {
		panic("not valid proven account state")
	}
	return pas.accountState.Version
}

// GetAddress returns a copy of the address of the account.
func (pas *ProvenAccountState) GetAddress() AccountAddress {
	if !pas.proven {
		panic("not valid proven account state")
	}
	return pas.addr
}

// GetAccountBlob returns a copy of the account blob, as a proven struct.
//
// GetAccountBlob returns nil if the account does not exist. You can call IsNil()
// to check whether the account exists.
func (pas *ProvenAccountState) GetAccountBlob() *ProvenAccountBlob {
	if !pas.proven {
		panic("not valid proven account state")
	}
	if pas.IsNil() {
		return nil
	}
	pab := &ProvenAccountBlob{
		proven:     true,
		addr:       pas.addr,
		ledgerInfo: pas.ledgerInfo,
	}
	if err := pab.accountBlob.ParseToMap(cloneBytes(pas.accountState.RawBlob)); err != nil {
		panic(err)
	}
	return pab
}

// IsNil returns whether this account is null (not exists, or not created yet).
func (pas *ProvenAccountState) IsNil() bool {
	if !pas.proven {
		panic("not valid proven account state")
	}
	return len(pas.accountState.RawBlob) == 0
}
