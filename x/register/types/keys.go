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
	// MetaNodeBondedPoolName stores the total balance of bonded Meta nodes
	MetaNodeBondedPoolName = "meta_node_bonded_pool"
	// MetaNodeNotBondedPoolName stores the total balance of not bonded meta nodes
	MetaNodeNotBondedPoolName = "meta_node_not_bonded_pool"
	// TotalUnissuedPrepayName stores the balance of total unissued prepay
	TotalUnissuedPrepayName = "total_unissued_prepay"
	// TotalSlashedPoolName stores the balance of total unissued prepay
	TotalSlashedPoolName = "total_slashed_pool"
)

var (
	ResourceNodeKey               = []byte{0x01} // prefix for each key to a resource node
	MetaNodeKey                   = []byte{0x02} // prefix for each key to a meta node
	MetaNodeRegistrationVotesKey  = []byte{0x03} // prefix for each key to the vote for meta node registration
	UpperBoundOfTotalOzoneKey     = []byte{0x04}
	SlashingPrefix                = []byte{0x05}
	InitialGenesisStakeTotalKey   = []byte{0x06} // key of initial genesis deposit by all resource nodes and meta nodes at t=0
	InitialUOzonePriceKey         = []byte{0x07} // key of initial uoz price at t=0
	MetaNodeCntKey                = []byte{0x08} // the number of all meta nodes
	ResourceNodeCntKey            = []byte{0x09} // the number of all resource nodes
	EffectiveGenesisStakeTotalKey = []byte{0x10} // key of effective(ongoing) genesis deposit by all resource nodes and meta nodes at time t

	UBDNodeKey      = []byte{0x11} // prefix for each key to an unbonding node
	UBDNodeQueueKey = []byte{0x12} // prefix for the timestamps in unbonding node queue

)

// GetResourceNodeKey gets the key for the resourceNode with address
// VALUE: ResourceNode
func GetResourceNodeKey(nodeAddr stratos.SdsAddress) []byte {
	return append(ResourceNodeKey, nodeAddr.Bytes()...)
}

// GetMetaNodeKey gets the key for the metaNode with address
// VALUE: ResourceNode
func GetMetaNodeKey(nodeAddr stratos.SdsAddress) []byte {
	return append(MetaNodeKey, nodeAddr.Bytes()...)
}

// GetMetaNodeRegistrationVotesKey get the key for the vote for Meta node registration
func GetMetaNodeRegistrationVotesKey(nodeAddr stratos.SdsAddress) []byte {
	return append(MetaNodeRegistrationVotesKey, nodeAddr.Bytes()...)
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
