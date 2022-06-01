package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	stratos "github.com/stratosnet/stratos-chain/types"
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
	// ResourceNodeBondedPoolName stores the total balance of bonded resource nodes
	ResourceNodeBondedPoolName = "resource_node_bonded_pool"
	// ResourceNodeNotBondedPoolName stores the total balance of not bonded resource nodes
	ResourceNodeNotBondedPoolName = "resource_node_not_bonded_pool"
	// IndexingNodeBondedPoolName stores the total balance of bonded indexing nodes
	IndexingNodeBondedPoolName = "indexing_node_bonded_pool"
	// IndexingNodeNotBondedPoolName stores the total balance of not bonded indexing nodes
	IndexingNodeNotBondedPoolName = "indexing_node_not_bonded_pool"
	// TotalUnissuedPrepayName stores the balance of total unissued prepay
	TotalUnissuedPrepayName = "total_unissued_prepay"
	// TotalSlashedPoolName stores the balance of total unissued prepay
	TotalSlashedPoolName = "total_slashed_pool"
)

var (
	ResourceNodeKey                  = []byte{0x01} // prefix for each key to a resource node
	IndexingNodeKey                  = []byte{0x02} // prefix for each key to a indexing node
	IndexingNodeRegistrationVotesKey = []byte{0x03} // prefix for each key to the vote for Indexing node registration
	UpperBoundOfTotalOzoneKey        = []byte{0x04}
	SlashingPrefix                   = []byte{0x05}
	InitialGenesisStakeTotalKey      = []byte{0x06} // key of initial genesis deposit by all resource nodes and meta nodes at t=0
	InitialUOzonePriceKey            = []byte{0x07} // key of initial uoz price at t=0

	UBDNodeKey      = []byte{0x11} // prefix for each key to an unbonding node
	UBDNodeQueueKey = []byte{0x12} // prefix for the timestamps in unbonding node queue
)

// GetResourceNodeKey gets the key for the resourceNode with address
// VALUE: ResourceNode
func GetResourceNodeKey(nodeAddr stratos.SdsAddress) []byte {
	return append(ResourceNodeKey, nodeAddr.Bytes()...)
}

// GetIndexingNodeKey gets the key for the indexingNode with address
// VALUE: ResourceNode
func GetIndexingNodeKey(nodeAddr stratos.SdsAddress) []byte {
	return append(IndexingNodeKey, nodeAddr.Bytes()...)
}

// GetIndexingNodeRegistrationVotesKey get the key for the vote for Indexing node registration
func GetIndexingNodeRegistrationVotesKey(nodeAddr stratos.SdsAddress) []byte {
	return append(IndexingNodeRegistrationVotesKey, nodeAddr.Bytes()...)
}

// GetURNKey gets the key for the unbonding Node with address
func GetUBDNodeKey(nodeAddr stratos.SdsAddress) []byte {
	return append(UBDNodeKey, nodeAddr.Bytes()...)
}

// gets the prefix for all unbonding Node
func GetUBDTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(UBDNodeQueueKey, bz...)
}

func GetSlashingKey(walletAddress sdk.AccAddress) []byte {
	key := append(SlashingPrefix, walletAddress...)
	return key
}
