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
	// ResourceNodeBondedPool stores the total balance of bonded resource nodes
	ResourceNodeBondedPool = "resource_node_bonded_pool"
	// ResourceNodeNotBondedPool stores the total balance of not bonded resource nodes
	ResourceNodeNotBondedPool = "resource_node_not_bonded_pool"
	// MetaNodeBondedPool stores the total balance of bonded Meta nodes
	MetaNodeBondedPool = "meta_node_bonded_pool"
	// MetaNodeNotBondedPool stores the total balance of not bonded meta nodes
	MetaNodeNotBondedPool = "meta_node_not_bonded_pool"
	// TotalUnissuedPrepay stores the balance of total unissued prepay
	TotalUnissuedPrepay = "total_unissued_prepay"
)

var (
	ResourceNodeKey                 = []byte{0x01} // prefix for each key to a resource node
	MetaNodeKey                     = []byte{0x02} // prefix for each key to a meta node
	MetaNodeRegistrationVotesKey    = []byte{0x03} // prefix for each key to the vote for meta node registration
	UpperBoundOfTotalOzoneKey       = []byte{0x04}
	SlashingPrefix                  = []byte{0x05}
	InitialGenesisDepositTotalKey   = []byte{0x06} // deprecated. key of initial genesis deposit by all resource nodes and meta nodes at t=0
	DepositNozRateKey               = []byte{0x07} // key of fixed rate of deposit unit over Noz unit
	MetaNodeCntKey                  = []byte{0x08} // the number of bonded meta nodes
	ResourceNodeCntKey              = []byte{0x09} // the number of bonded resource nodes
	EffectiveGenesisDepositTotalKey = []byte{0x10} // key of effective(ongoing) genesis deposit by all resource nodes and meta nodes at time t
	UBDNodeKey                      = []byte{0x11} // prefix for each key to an unbonding node
	UBDNodeQueueKey                 = []byte{0x12} // prefix for the timestamps in unbonding node queue
	KickMetaNodeVotesKey            = []byte{0x13} // prefix for the key to the vote for kicking meta node

	ParamsKey = []byte{0x20}
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

// GetKickMetaNodeVotesKey get the key prefix for kicking meta node vote pool
func GetKickMetaNodeVotesKey(nodeAddr stratos.SdsAddress) []byte {
	return append(KickMetaNodeVotesKey, nodeAddr.Bytes()...)
}

// GetUBDNodeKey gets the key for the unbonding Node with address
func GetUBDNodeKey(nodeAddr stratos.SdsAddress) []byte {
	return append(UBDNodeKey, nodeAddr.Bytes()...)
}

// GetUBDTimeKey gets the prefix for all unbonding Node
func GetUBDTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(UBDNodeQueueKey, bz...)
}

func GetSlashingKey(walletAddress sdk.AccAddress) []byte {
	key := append(SlashingPrefix, walletAddress...)
	return key
}
