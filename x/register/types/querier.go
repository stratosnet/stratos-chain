package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

const (
	defaultDenom  = "ustos"
	QueryType_All = 0
	QueryType_SP  = 1
	QueryType_PP  = 2
)

// QueryNodesParams Params for query 'custom/register/resource-nodes'
type QueryNodesParams struct {
	Page        int
	Limit       int
	NetworkAddr stratos.SdsAddress
	Moniker     string
	OwnerAddr   sdk.AccAddress
}

// NewQueryNodesParams creates a new instance of QueryNodesParams
func NewQueryNodesParams(page, limit int, networkAddr stratos.SdsAddress, moniker string, ownerAddr sdk.AccAddress) QueryNodesParams {
	return QueryNodesParams{
		Page:        page,
		Limit:       limit,
		NetworkAddr: networkAddr,
		Moniker:     moniker,
		OwnerAddr:   ownerAddr,
	}
}

type QueryNodeStakingParams struct {
	AccAddr   stratos.SdsAddress
	QueryType int64 //0:All(Default) 1: indexingNode; 2: ResourceNode
}

// NewQueryNodeStakingParams creates a new instance of QueryNodesParams
func NewQueryNodeStakingParams(nodeAddr stratos.SdsAddress, queryType int64) QueryNodeStakingParams {
	return QueryNodeStakingParams{
		AccAddr:   nodeAddr,
		QueryType: queryType,
	}
}

// NewQueryNodesStakingInfo creates a new instance of TotalStakesResponse
func NewQueryNodesStakingInfo(ResourceNodeTotalStake, IndexingNodeTotalStake, totalBondedStake, totalUnbondedStake, totalUnbondingStake sdk.Int) *TotalStakesResponse {
	resValue := sdk.NewCoin(defaultDenom, ResourceNodeTotalStake)
	indValue := sdk.NewCoin(defaultDenom, IndexingNodeTotalStake)
	bonedValue := sdk.NewCoin(defaultDenom, totalBondedStake)
	unBondedValue := sdk.NewCoin(defaultDenom, totalUnbondedStake)
	unBondingValue := sdk.NewCoin(defaultDenom, totalUnbondingStake)

	return &TotalStakesResponse{
		ResourceNodesTotalStake: &resValue,
		IndexingNodesTotalStake: &indValue,
		TotalBondedStake:        &bonedValue,
		TotalUnbondedStake:      &unBondedValue,
		TotalUnbondingStake:     &unBondingValue,
	}
}

// NewStakingInfoByResourceNodeAddr creates a new instance of StakingInfoByNodeAddr
func NewStakingInfoByResourceNodeAddr(
	resourceNode ResourceNode,
	unBondingStake sdk.Int,
	unBondedStake sdk.Int,
	bondedStake sdk.Int,

) StakingInfo {
	bonedValue := sdk.NewCoin(defaultDenom, bondedStake)
	unBondedValue := sdk.NewCoin(defaultDenom, unBondedStake)
	unBondingValue := sdk.NewCoin(defaultDenom, unBondingStake)

	return StakingInfo{
		NetworkAddr:    resourceNode.GetNetworkAddr(),
		PubKey:         resourceNode.GetPubKey(),
		Suspend:        resourceNode.GetSuspend(),
		Status:         resourceNode.GetStatus(),
		Tokens:         &resourceNode.Tokens,
		OwnerAddress:   resourceNode.GetOwnerAddress(),
		Description:    resourceNode.GetDescription(),
		NodeType:       resourceNode.GetNodeType(),
		CreationTime:   resourceNode.GetCreationTime(),
		UnBondingStake: &unBondingValue,
		UnBondedStake:  &unBondedValue,
		BondedStake:    &bonedValue,
	}
}

// NewStakingInfoByIndexingNodeAddr creates a new instance of StakingInfoByNodeAddr
func NewStakingInfoByIndexingNodeAddr(
	indexingNode IndexingNode,
	unBondingStake sdk.Int,
	unBondedStake sdk.Int,
	bondedStake sdk.Int,
) StakingInfo {
	bonedValue := sdk.NewCoin(defaultDenom, bondedStake)
	unBondedValue := sdk.NewCoin(defaultDenom, unBondedStake)
	unBondingValue := sdk.NewCoin(defaultDenom, unBondingStake)
	return StakingInfo{
		NetworkAddr:    indexingNode.GetNetworkAddr(),
		PubKey:         indexingNode.GetPubKey(),
		Suspend:        indexingNode.Suspend,
		Status:         indexingNode.Status,
		Tokens:         &indexingNode.Tokens,
		OwnerAddress:   indexingNode.GetOwnerAddress(),
		Description:    indexingNode.Description,
		NodeType:       "metanode",
		CreationTime:   indexingNode.CreationTime,
		UnBondingStake: &unBondingValue,
		UnBondedStake:  &unBondedValue,
		BondedStake:    &bonedValue,
	}
}
