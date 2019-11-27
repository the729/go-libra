package types

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/crypto/sha3libra"
	"github.com/the729/lcs"
)

// RawAccountBlob is the raw blob of an account.
type RawAccountBlob []byte

// AccountBlob is the blob of an account. It is a map of resources.
type AccountBlob struct {
	Map map[string][]byte
}

// ProvenAccountBlob is and account blob proven to be included in the ledger.
type ProvenAccountBlob struct {
	proven      bool
	accountBlob AccountBlob
	addr        AccountAddress
	ledgerInfo  *ProvenLedgerInfo
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (b RawAccountBlob) Hash() sha3libra.HashValue {
	if b == nil {
		return nil
	}
	hasher := sha3libra.NewAccountStateBlob()
	hasher.Write(b)
	return hasher.Sum([]byte{})
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (b *AccountBlob) Hash() sha3libra.HashValue {
	if b == nil {
		return nil
	}
	raw, err := lcs.Marshal(b)
	if err != nil {
		panic(err)
	}
	return RawAccountBlob(raw).Hash()
}

// ParseToMap parses the raw blob into a map of resources.
func (b *AccountBlob) ParseToMap(raw RawAccountBlob) error {
	return lcs.Unmarshal(raw, b)
}

// GetResource gets a resource from the account blob by its path.
//
// The account blob should be already parsed into map of resources. To get Libra coin account resource,
// use AccountResourcePath() to generate the path.
func (b *AccountBlob) GetResource(path []byte) (*AccountResource, error) {
	key := string(path)
	val, ok := b.Map[key]
	if !ok {
		return nil, errors.New("resource not found")
	}
	r := &AccountResource{}
	if err := lcs.Unmarshal(val, r); err != nil {
		return nil, fmt.Errorf("unmarshal resource error: %v", err)
	}
	return r, nil
}

// GetLedgerInfo returns the ledger info.
func (pb *ProvenAccountBlob) GetLedgerInfo() *ProvenLedgerInfo {
	if !pb.proven {
		panic("not valid proven account blob")
	}
	return pb.ledgerInfo
}

// GetResource gets a resource from a proven account blob by its path.
//
// To get Libra coin account resource, use AccountResourcePath() to generate the path.
func (pb *ProvenAccountBlob) GetResource(path []byte) (*ProvenAccountResource, error) {
	if !pb.proven {
		panic("not valid proven account blob")
	}
	ar, err := pb.accountBlob.GetResource(path)
	if err != nil {
		return nil, err
	}
	par := &ProvenAccountResource{
		proven:     true,
		addr:       cloneBytes(pb.addr),
		ledgerInfo: pb.ledgerInfo,
	}
	par.accountResource = *(ar.Clone())
	return par, nil
}

// GetAddress returns a copy of account address.
func (pb *ProvenAccountBlob) GetAddress() AccountAddress {
	if !pb.proven {
		panic("not valid proven account blob")
	}
	return AccountAddress(cloneBytes(pb.addr))
}
