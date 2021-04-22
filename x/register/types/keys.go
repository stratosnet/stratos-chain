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
	ResourceNodeByAddrKey        = []byte{0x22} // prefix for each key to a resource node index, by pubkey
	ResourceNodesByPowerIndexKey = []byte{0x23} // prefix for each key to a resource node index, sorted by power

	IndexingNodeKey              = []byte{0x31} // prefix for each key to a indexing node
	IndexingNodeByAddrKey        = []byte{0x32} // prefix for each key to a indexing node index, by pubkey
	IndexingNodesByPowerIndexKey = []byte{0x33} // prefix for each key to a indexing node, sorted by power
)

// GetNodeKey gets the key for the resourceNode/indexingNode with address
// VALUE: staking/node
func GetNodeKey(nodeType NodeType, nodeAddr sdk.ValAddress) []byte {
	switch nodeType {
	case NodeTypeResource:
		return append(ResourceNodeKey, nodeAddr.Bytes()...)
	case NodeTypeIndexing:
		return append(IndexingNodeKey, nodeAddr.Bytes()...)
	default:
		return nil
	}
}

// gets the key for the resource node with pubkey
// VALUE: resource node operator address ([]byte)
func GetResourceNodeByAddrKey(addr sdk.ConsAddress) []byte {
	return append(ResourceNodeByAddrKey, addr.Bytes()...)
}

// GetResourceNodesByPowerIndexKey get the resource node by power index.
// Power index is the key used in the power-store, and represents the relative
// power ranking of the resource node.
// VALUE: resource node operator address ([]byte)
func GetResourceNodesByPowerIndexKey(resourceNode ResourceNode) []byte {
	// NOTE the address doesn't need to be stored because counter bytes must always be different
	return getResourceNodePowerRank(resourceNode)
}

// get the power ranking of a resource node
// NOTE the larger values are of higher value
func getResourceNodePowerRank(resourceNode ResourceNode) []byte {

	resourcePower := TokensToPower(resourceNode.Tokens)
	resourcePowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(resourcePowerBytes, uint64(resourcePower))

	powerBytes := resourcePowerBytes
	powerBytesLen := len(powerBytes) // 8

	// key is of format prefix || powerbytes || addrBytes
	key := make([]byte, 1+powerBytesLen+sdk.AddrLen)

	key[0] = ResourceNodesByPowerIndexKey[0]
	copy(key[1:powerBytesLen+1], powerBytes)
	operAddrInvr := sdk.CopyBytes(resourceNode.OperatorAddress)
	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}
	copy(key[powerBytesLen+1:], operAddrInvr)

	return key
}

// gets the key for the indexing node with pubkey
// VALUE: indexing node operator address ([]byte)
func GetIndexingNodeByAddrKey(addr sdk.ConsAddress) []byte {
	return append(IndexingNodeByAddrKey, addr.Bytes()...)
}

// GetResourceNodesByPowerIndexKey get the indexing node by power index.
// Power index is the key used in the power-store, and represents the relative
// power ranking of the indexing node.
// VALUE: indexing node operator address ([]byte)
func GetIndexingNodesByPowerIndexKey(indexingNode IndexingNode) []byte {
	// NOTE the address doesn't need to be stored because counter bytes must always be different
	return getIndexingNodePowerRank(indexingNode)
}

// get the power ranking of a indexing node
// NOTE the larger values are of higher value
func getIndexingNodePowerRank(indexingNode IndexingNode) []byte {

	indexingPower := TokensToPower(indexingNode.Tokens)
	indexingPowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexingPowerBytes, uint64(indexingPower))

	powerBytes := indexingPowerBytes
	powerBytesLen := len(powerBytes) // 8

	// key is of format prefix || powerbytes || addrBytes
	key := make([]byte, 1+powerBytesLen+sdk.AddrLen)

	key[0] = IndexingNodesByPowerIndexKey[0]
	copy(key[1:powerBytesLen+1], powerBytes)
	operAddrInvr := sdk.CopyBytes(indexingNode.OperatorAddress)
	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}
	copy(key[powerBytesLen+1:], operAddrInvr)

	return key
}
