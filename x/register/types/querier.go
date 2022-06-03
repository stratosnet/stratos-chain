package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
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
	QueryType int64 //0:All(Default) 1: MetaNode; 2: ResourceNode
}

// NewQueryNodeStakingParams creates a new instance of QueryNodesParams
func NewQueryNodeStakingParams(nodeAddr stratos.SdsAddress, queryType int64) QueryNodeStakingParams {
	return QueryNodeStakingParams{
		AccAddr:   nodeAddr,
		QueryType: queryType,
	}
}

// NewQueryNodesStakingInfo creates a new instance of TotalStakesResponse
func NewQueryNodesStakingInfo(ResourceNodeTotalStake, MetaNodeTotalStake, totalBondedStake, totalUnbondedStake, totalUnbondingStake sdk.Int) *TotalStakesResponse {
	resValue := sdk.NewCoin(defaultDenom, ResourceNodeTotalStake)
	metaValue := sdk.NewCoin(defaultDenom, MetaNodeTotalStake)
	bonedValue := sdk.NewCoin(defaultDenom, totalBondedStake)
	unBondedValue := sdk.NewCoin(defaultDenom, totalUnbondedStake)
	unBondingValue := sdk.NewCoin(defaultDenom, totalUnbondingStake)

	return &TotalStakesResponse{
		ResourceNodesTotalStake: &resValue,
		MetaNodesTotalStake:     &metaValue,
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
		NetworkAddress: resourceNode.GetNetworkAddress(),
		Pubkey:         resourceNode.GetPubkey(),
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

// NewStakingInfoByMetaNodeAddr creates a new instance of StakingInfoByNodeAddr
func NewStakingInfoByMetaNodeAddr(
	metaNode MetaNode,
	unBondingStake sdk.Int,
	unBondedStake sdk.Int,
	bondedStake sdk.Int,
) StakingInfo {
	bonedValue := sdk.NewCoin(defaultDenom, bondedStake)
	unBondedValue := sdk.NewCoin(defaultDenom, unBondedStake)
	unBondingValue := sdk.NewCoin(defaultDenom, unBondingStake)
	return StakingInfo{
		NetworkAddress: metaNode.GetNetworkAddress(),
		Pubkey:         metaNode.GetPubkey(),
		Suspend:        metaNode.Suspend,
		Status:         metaNode.Status,
		Tokens:         &metaNode.Tokens,
		OwnerAddress:   metaNode.GetOwnerAddress(),
		Description:    metaNode.Description,
		NodeType:       "metanode",
		CreationTime:   metaNode.CreationTime,
		UnBondingStake: &unBondingValue,
		UnBondedStake:  &unBondedValue,
		BondedStake:    &bonedValue,
	}
}

type StakingInfos []StakingInfo

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v StakingInfo) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(v.Pubkey, &pk)
}

func (v StakingInfos) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range v {
		if err := v[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}
