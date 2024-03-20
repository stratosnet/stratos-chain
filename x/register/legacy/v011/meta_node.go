package v011

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

// MustUnmarshalMetaNode unmarshal an meta node from a store value. Panics if fails
func MustUnmarshalMetaNode(cdc codec.Codec, value []byte) MetaNode {
	metaNode, err := UnmarshalMetaNode(cdc, value)
	if err != nil {
		panic(err)
	}
	return metaNode
}

// UnmarshalMetaNode unmarshal an meta node from a store value
func UnmarshalMetaNode(cdc codec.Codec, value []byte) (metaNode MetaNode, err error) {
	err = cdc.Unmarshal(value, &metaNode)
	return metaNode, err
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v MetaNode) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(v.Pubkey, &pk)
}
