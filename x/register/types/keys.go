package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
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
	ResourceNodeNotBondedTokenKey = []byte{0x01}
	ResourceNodeBondedTokenKey    = []byte{0x02}
	IndexingNodeNotBondedTokenKey = []byte{0x03}
	IndexingNodeBondedTokenKey    = []byte{0x04}
	UpperBoundOfTotalOzoneKey     = []byte{0x05}

	LastResourceNodeStakeKey    = []byte{0x11} // prefix for each key to a resource node index, for bonded resource nodes
	LastIndexingNodeStakeKey    = []byte{0x12} // prefix for each key to a indexing node index, for bonded indexing nodes
	InitialGenesisStakeTotalKey = []byte{0x13} // key of initial genesis deposit by all resource nodes and meta nodes at t=0
	InitialUOzonePriceKey       = []byte{0x14} // key of initial uoz price at t=0

	ResourceNodeKey                  = []byte{0x21} // prefix for each key to a resource node
	IndexingNodeKey                  = []byte{0x22} // prefix for each key to a indexing node
	IndexingNodeRegistrationVotesKey = []byte{0x23} // prefix for each key to the vote for Indexing node registration

	UBDNodeKey = []byte{0x31} // prefix for each key to an unbonding node

	UBDNodeQueueKey = []byte{0x41} // prefix for the timestamps in unbonding node queue
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

// GetIndexingNodeRegistrationVotesKey get the key for the vote for Indexing node registration
func GetIndexingNodeRegistrationVotesKey(nodeAddr sdk.AccAddress) []byte {
	return append(IndexingNodeRegistrationVotesKey, nodeAddr.Bytes()...)
}

// GetURNKey gets the key for the unbonding Node with address
func GetUBDNodeKey(nodeAddr sdk.AccAddress) []byte {
	return append(UBDNodeKey, nodeAddr.Bytes()...)
}

// gets the prefix for all unbonding Node
func GetUBDTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(UBDNodeQueueKey, bz...)
}
