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
func (b RawAccountBlob) Hash() HashValue {
	if b == nil {
		return nil
	}
	hasher := sha3libra.NewAccountStateBlob()
	hasher.Write(b)
	return hasher.Sum([]byte{})
}

// Hash ouptuts the hash of this struct, using the appropriate hash function.
func (b *AccountBlob) Hash() HashValue {
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
func (b *AccountBlob) GetResource(path []byte) ([]byte, error) {
	key := string(path)
	val, ok := b.Map[key]
	if !ok {
		return nil, errors.New("resource not found")
	}
	return val, nil
}

// GetLibraAccountResource gets 0x0.LibraAccount.T resource from the account blob.
//
// The account blob should be already parsed into map of resources.
func (b *AccountBlob) GetLibraAccountResource() (*AccountResource, error) {
	val, err := b.GetResource(AccountResourcePath())
	if err != nil {
		return nil, err
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
func (pb *ProvenAccountBlob) GetResource(path []byte) ([]byte, error) {
	if !pb.proven {
		panic("not valid proven account blob")
	}
	ar, err := pb.accountBlob.GetResource(path)
	if err != nil {
		return nil, err
	}
	return cloneBytes(ar), nil
}

// GetResourcePaths gets a list of resource paths from a proven account blob.
func (pb *ProvenAccountBlob) GetResourcePaths() [][]byte {
	if !pb.proven {
		panic("not valid proven account blob")
	}
	paths := make([][]byte, 0, len(pb.accountBlob.Map))
	for p := range pb.accountBlob.Map {
		paths = append(paths, []byte(p))
	}
	return paths
}

// GetLibraAccountResource gets 0x0.LibraAccount.T resource from a proven account blob.
func (pb *ProvenAccountBlob) GetLibraAccountResource() (*ProvenAccountResource, error) {
	if !pb.proven {
		panic("not valid proven account blob")
	}
	ar, err := pb.accountBlob.GetLibraAccountResource()
	if err != nil {
		return nil, err
	}
	par := &ProvenAccountResource{
		proven:          true,
		accountResource: ar.Clone(),
		addr:            pb.addr,
		ledgerInfo:      pb.ledgerInfo,
	}
	return par, nil
}

// GetAddress returns a copy of account address.
func (pb *ProvenAccountBlob) GetAddress() AccountAddress {
	if !pb.proven {
		panic("not valid proven account blob")
	}
	return pb.addr
}
