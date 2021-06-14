package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bech32PubKeyType defines a string type alias for a Bech32 public key type.
type Bech32PubKeyType string

const (
	Bech32PubKeyTypesdsPub Bech32PubKeyType = "sdspub"
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
	UpperBoundOfTotalOzoneKey = []byte{0x01}

	LastResourceNodeStakeKey      = []byte{0x11} // prefix for each key to a resource node index, for bonded resource nodes
	LastResourceNodeTotalStakeKey = []byte{0x12} // prefix for the total bonded tokens of resource nodes
	LastIndexingNodeStakeKey      = []byte{0x13} // prefix for each key to a indexing node index, for bonded indexing nodes
	LastIndexingNodeTotalStakeKey = []byte{0x14} // prefix for the total bonded tokens of indexing nodes
	InitialGenesisStakeTotalKey   = []byte{0x15} // key of initial genesis deposit by all resource nodes and meta nodes at t=0

	ResourceNodeKey        = []byte{0x21} // prefix for each key to a resource node
	IndexingNodeKey        = []byte{0x22} // prefix for each key to a indexing node
	SpRegistrationVotesKey = []byte{0x23} // prefix for each key to the vote for SP node registration
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

// GetSpRegistrationVotesKey get the key for the vote for SP node registration
func GetSpRegistrationVotesKey(nodeAddr sdk.AccAddress) []byte {
	return append(SpRegistrationVotesKey, nodeAddr.Bytes()...)
}
