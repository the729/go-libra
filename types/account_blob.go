package types

import (
	"errors"
	"fmt"

	"github.com/the729/go-libra/common/canonical_serialization"
	"github.com/the729/go-libra/crypto/sha3libra"
)

type AccountBlob struct {
	Raw []byte
	Map map[string][]byte
}

func (b *AccountBlob) Hash() sha3libra.HashValue {
	if b == nil {
		return nil
	}
	hasher := sha3libra.NewAccountStateBlob()
	hasher.Write(b.Raw)
	return hasher.Sum([]byte{})
}

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

func (b *AccountBlob) GetResource(tag *StructTag) (*AccountResource, error) {
	tagHash := tag.Hash()
	key := "\x01" + string(tagHash)
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
