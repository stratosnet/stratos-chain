package types

import (
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
	LastResourceNodeStakeKey      = []byte{0x11} // prefix for each key to a resource node index, for bonded resource nodes
	LastResourceNodeTotalStakeKey = []byte{0x12} // prefix for the total bonded tokens of resource nodes
	LastIndexingNodeStakeKey      = []byte{0x13} // prefix for each key to a indexing node index, for bonded indexing nodes
	LastIndexingNodeTotalStakeKey = []byte{0x14} // prefix for the total bonded tokens of indexing nodes

	ResourceNodeKey = []byte{0x21} // prefix for each key to a resource node
	IndexingNodeKey = []byte{0x22} // prefix for each key to a indexing node
)

// GetLastResourceNodeStakeKey get the bonded resource node index key for an address
func GetLastResourceNodeStakeKey(nodeAddr sdk.AccAddress) []byte {
	return append(LastResourceNodeStakeKey, nodeAddr...)
}

// GetResourceNodeKey gets the key for the resourceNode with address
// VALUE: ResourceNode
func GetResourceNodeKey(nodeAddr sdk.AccAddress) []byte {
	return append(ResourceNodeKey, nodeAddr.Bytes()...)
}

// GetLastIndexingNodeStakeKey get the bonded indexing node index key for an address
func GetLastIndexingNodeStakeKey(nodeAddr sdk.AccAddress) []byte {
	return append(LastIndexingNodeStakeKey, nodeAddr...)
}

// GetIndexingNodeKey gets the key for the indexingNode with address
// VALUE: ResourceNode
func GetIndexingNodeKey(nodeAddr sdk.AccAddress) []byte {
	return append(IndexingNodeKey, nodeAddr.Bytes()...)
}
