package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	defaultDenom  = "ustos"
	QueryType_All = 0
	QueryType_SP  = 1
	QueryType_PP  = 2
)

// QueryNodesParams Params for query 'custom/register/resource-nodes'
type QueryNodesParams struct {
	Page      int
	Limit     int
	NetworkID string
	Moniker   string
	OwnerAddr sdk.AccAddress
}

// NewQueryNodesParams creates a new instance of QueryNodesParams
func NewQueryNodesParams(page, limit int, networkID, moniker string, ownerAddr sdk.AccAddress) QueryNodesParams {
	return QueryNodesParams{
		Page:      page,
		Limit:     limit,
		NetworkID: networkID,
		Moniker:   moniker,
		OwnerAddr: ownerAddr,
	}
}

type QueryNodeStakingParams struct {
	AccAddr   sdk.AccAddress
	QueryType int64 //0:All(Default) 1: indexingNode; 2: ResourceNode
}

// NewQueryNodeStakingParams creates a new instance of QueryNodesParams
func NewQueryNodeStakingParams(nodeAddr sdk.AccAddress, queryType int64) QueryNodeStakingParams {
	return QueryNodeStakingParams{
		AccAddr:   nodeAddr,
		QueryType: queryType,
	}
}

// NodesStakingInfo Params for query 'custom/register/staking'
type NodesStakingInfo struct {
	TotalStakeOfResourceNodes sdk.Coin
	TotalStakeOfIndexingNodes sdk.Coin
	TotalBondedStake          sdk.Coin
	TotalUnbondedStake        sdk.Coin
	TotalUnbondingStake       sdk.Coin
}

// NewQueryNodesStakingInfo creates a new instance of NodesStakingInfo
func NewQueryNodesStakingInfo(
	totalStakeOfResourceNodes,
	totalStakeOfIndexingNodes,
	totalBondedStake,
	totalUnbondedStake,
	totalUnbondingStake sdk.Int,
) NodesStakingInfo {
	return NodesStakingInfo{
		TotalStakeOfResourceNodes: sdk.NewCoin(defaultDenom, totalStakeOfResourceNodes),
		TotalStakeOfIndexingNodes: sdk.NewCoin(defaultDenom, totalStakeOfIndexingNodes),
		TotalBondedStake:          sdk.NewCoin(defaultDenom, totalBondedStake),
		TotalUnbondedStake:        sdk.NewCoin(defaultDenom, totalUnbondedStake),
		TotalUnbondingStake:       sdk.NewCoin(defaultDenom, totalUnbondingStake),
	}
}

type StakingInfo struct {
	NetworkID      string         `json:"network_id"`
	PubKey         crypto.PubKey  `json:"pub_key"`
	Suspend        bool           `json:"suspend"`
	Status         sdk.BondStatus `json:"status"`
	Tokens         sdk.Int        `json:"tokens"`
	OwnerAddress   sdk.AccAddress `json:"owner_address"`
	Description    Description    `json:"description"`
	NodeType       string         `json:"node_type"`
	CreationTime   time.Time      `json:"creation_time"`
	BondedStake    sdk.Coin       `json:"bonded_stake"`
	UnBondingStake sdk.Coin       `json:"un_bonding_stake"`
	UnBondedStake  sdk.Coin       `json:"un_bonded_stake"`
}

// NewStakingInfoByResourceNodeAddr creates a new instance of StakingInfoByNodeAddr
func NewStakingInfoByResourceNodeAddr(
	resourceNode ResourceNode,
	unBondingStake sdk.Int,
	unBondedStake sdk.Int,
	bondedStake sdk.Int,

) StakingInfo {
	return StakingInfo{
		NetworkID:      resourceNode.NetworkID,
		PubKey:         resourceNode.PubKey,
		Suspend:        resourceNode.Suspend,
		Status:         resourceNode.Status,
		Tokens:         resourceNode.Tokens,
		OwnerAddress:   resourceNode.OwnerAddress,
		Description:    resourceNode.Description,
		NodeType:       resourceNode.NodeType,
		CreationTime:   resourceNode.CreationTime,
		UnBondingStake: sdk.NewCoin(defaultDenom, unBondingStake),
		UnBondedStake:  sdk.NewCoin(defaultDenom, unBondedStake),
		BondedStake:    sdk.NewCoin(defaultDenom, bondedStake),
	}
}

// NewStakingInfoByIndexingNodeAddr creates a new instance of StakingInfoByNodeAddr
func NewStakingInfoByIndexingNodeAddr(
	indexingNode IndexingNode,
	unBondingStake sdk.Int,
	unBondedStake sdk.Int,
	bondedStake sdk.Int,
) StakingInfo {
	return StakingInfo{
		NetworkID:      indexingNode.NetworkID,
		PubKey:         indexingNode.PubKey,
		Suspend:        indexingNode.Suspend,
		Status:         indexingNode.Status,
		Tokens:         indexingNode.Tokens,
		OwnerAddress:   indexingNode.OwnerAddress,
		Description:    indexingNode.Description,
		NodeType:       "SP",
		CreationTime:   indexingNode.CreationTime,
		UnBondingStake: sdk.NewCoin(defaultDenom, unBondingStake),
		UnBondedStake:  sdk.NewCoin(defaultDenom, unBondedStake),
		BondedStake:    sdk.NewCoin(defaultDenom, bondedStake),
	}
}
