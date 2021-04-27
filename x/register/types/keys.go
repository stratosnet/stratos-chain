package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "register"
	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
	// RouterKey to be used for routing msgs
	RouterKey = ModuleName
	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName
)

var (
	ResourceNodeKey              = []byte{0x21} // prefix for each key to a resource node
	ResourceNodesByPowerIndexKey = []byte{0x22} // prefix for each key to a resource node index, sorted by power

	IndexingNodeKey              = []byte{0x23} // prefix for each key to a indexing node
	IndexingNodesByPowerIndexKey = []byte{0x24} // prefix for each key to a indexing node, sorted by power
)

// GetResourceNodeKey gets the key for the resourceNode with address
// VALUE: ResourceNode
func GetResourceNodeKey(nodeAddr sdk.AccAddress) []byte {
	return append(ResourceNodeKey, nodeAddr.Bytes()...)
}

// GetResourceNodesByPowerIndexKey get the resource node by power index.
// Power index is the key used in the power-store, and represents the relative
// power ranking of the resource node.
// VALUE: resource node address ([]byte)
func GetResourceNodesByPowerIndexKey(resourceNode ResourceNode) []byte {
	resourcePower := TokensToPower(resourceNode.Tokens)
	resourcePowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(resourcePowerBytes, uint64(resourcePower))

	powerBytes := resourcePowerBytes
	powerBytesLen := len(powerBytes) // 8

	// key is of format prefix || powerbytes || addrBytes
	key := make([]byte, 1+powerBytesLen+sdk.AddrLen)

	key[0] = ResourceNodesByPowerIndexKey[0]
	copy(key[1:powerBytesLen+1], powerBytes)
	addrInvr := sdk.CopyBytes(resourceNode.GetAddr())
	for i, b := range addrInvr {
		addrInvr[i] = ^b
	}
	copy(key[powerBytesLen+1:], addrInvr)

	return key
}

// GetIndexingNodeKey gets the key for the indexingNode with address
// VALUE: ResourceNode
func GetIndexingNodeKey(nodeAddr sdk.AccAddress) []byte {
	return append(IndexingNodeKey, nodeAddr.Bytes()...)
}

// GetIndexingNodesByPowerIndexKey get the indexing node by power index.
// Power index is the key used in the power-store, and represents the relative
// power ranking of the indexing node.
// VALUE: indexing node address ([]byte)
func GetIndexingNodesByPowerIndexKey(indexingNode IndexingNode) []byte {
	indexingPower := TokensToPower(indexingNode.Tokens)
	indexingPowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexingPowerBytes, uint64(indexingPower))

	powerBytes := indexingPowerBytes
	powerBytesLen := len(powerBytes) // 8

	// key is of format prefix || powerbytes || addrBytes
	key := make([]byte, 1+powerBytesLen+sdk.AddrLen)

	key[0] = IndexingNodesByPowerIndexKey[0]
	copy(key[1:powerBytesLen+1], powerBytes)
	addrInvr := sdk.CopyBytes(indexingNode.GetAddr())
	for i, b := range addrInvr {
		addrInvr[i] = ^b
	}
	copy(key[powerBytesLen+1:], addrInvr)

	return key
}
