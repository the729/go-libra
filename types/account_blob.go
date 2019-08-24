package types

import (
	"errors"
	"fmt"

	serialization "github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
)

// RawAccountBlob is the raw blob of an account.
type RawAccountBlob []byte

// AccountBlob is the blob of an account.
//
// It consists of the raw blob, and the decoded map of resources.
type AccountBlob struct {
	Raw []byte
	Map map[string][]byte
}

// ProvenAccountBlob is and account blob proven to be included in the ledger.
type ProvenAccountBlob struct {
	proven      bool
	accountBlob AccountBlob
	addr        AccountAddress
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
	return RawAccountBlob(b.Raw).Hash()
}

// ParseToMap parses the raw blob into a map of resources.
func (b *AccountBlob) ParseToMap() error {
	data := b.Raw
	l := int(serialization.SimpleDeserializer.Uint32(data))
	data = data[4:]
	m := make(map[string][]byte)
	for i := 0; i < l; i++ {
		key, err := serialization.SimpleDeserializer.ByteSlice(data)
		if err != nil {
			return errors.New("error deserizaing key")
		}
		data = data[len(key)+4:]
		val, err := serialization.SimpleDeserializer.ByteSlice(data)
		if err != nil {
			return errors.New("error deserizaing val")
		}
		data = data[len(val)+4:]
		m[string(key)] = val
	}
	b.Map = m
	return nil
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
	err := r.UnmarshalBinary(val)
	if err != nil {
		return nil, fmt.Errorf("unmarshal resource error: %v", err)
	}
	return r, nil
}

// GetRawBlob returns a copy of raw blob
func (pb *ProvenAccountBlob) GetRawBlob() []byte {
	if !pb.proven {
		panic("not valid proven account blob")
	}
	return cloneBytes(pb.accountBlob.Raw)
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
	return &ProvenAccountResource{
		proven: true,
		accountResource: AccountResource{
			Balance:             ar.Balance,
			SequenceNumber:      ar.SequenceNumber,
			AuthenticationKey:   cloneBytes(ar.AuthenticationKey),
			SentEventsCount:     ar.SentEventsCount,
			ReceivedEventsCount: ar.ReceivedEventsCount,
		},
		addr: cloneBytes(pb.addr),
	}, nil
}

// GetAddress returns a copy of account address.
func (pb *ProvenAccountBlob) GetAddress() AccountAddress {
	if !pb.proven {
		panic("not valid proven account blob")
	}
	return AccountAddress(cloneBytes(pb.addr))
}
