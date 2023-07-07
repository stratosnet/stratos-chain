package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pagiquery "github.com/cosmos/cosmos-sdk/types/query"

	stratos "github.com/stratosnet/stratos-chain/types"
)

const (
	QueryType_All     = 0
	QueryType_SP      = 1
	QueryType_PP      = 2
	QueryDefaultLimit = 100
)

// QueryNodesParams Params for query 'custom/register/resource-nodes'
type QueryNodesParams struct {
	PageQuery   pagiquery.PageRequest
	NetworkAddr stratos.SdsAddress
	Moniker     string
	OwnerAddr   sdk.AccAddress
}

// NewQueryNodesParams creates a new instance of QueryNodesParams
func NewQueryNodesParams(networkAddr stratos.SdsAddress, moniker string, ownerAddr sdk.AccAddress, pageQuery pagiquery.PageRequest) QueryNodesParams {
	return QueryNodesParams{
		PageQuery:   pageQuery,
		NetworkAddr: networkAddr,
		Moniker:     moniker,
		OwnerAddr:   ownerAddr,
	}
}

type QueryNodeDepositParams struct {
	AccAddr   stratos.SdsAddress
	QueryType int64 //0:All(Default) 1: MetaNode; 2: ResourceNode
}

// NewQueryNodeDepositParams creates a new instance of QueryNodesParams
func NewQueryNodeDepositParams(nodeAddr stratos.SdsAddress, queryType int64) QueryNodeDepositParams {
	return QueryNodeDepositParams{
		AccAddr:   nodeAddr,
		QueryType: queryType,
	}
}

// NewQueryDepositTotalInfo creates a new instance of QueryDepositTotalResponse
func NewQueryDepositTotalInfo(bondDenom string, ResourceNodeTotalDeposit, MetaNodeTotalDeposit, totalBondedDeposit,
	totalUnbondedDeposit, totalUnbondingDeposit sdk.Int) *QueryDepositTotalResponse {

	resValue := sdk.NewCoin(bondDenom, ResourceNodeTotalDeposit)
	metaValue := sdk.NewCoin(bondDenom, MetaNodeTotalDeposit)
	bonedValue := sdk.NewCoin(bondDenom, totalBondedDeposit)
	unBondedValue := sdk.NewCoin(bondDenom, totalUnbondedDeposit)
	unBondingValue := sdk.NewCoin(bondDenom, totalUnbondingDeposit)

	return &QueryDepositTotalResponse{
		ResourceNodesTotalDeposit: &resValue,
		MetaNodesTotalDeposit:     &metaValue,
		TotalBondedDeposit:        &bonedValue,
		TotalUnbondedDeposit:      &unBondedValue,
		TotalUnbondingDeposit:     &unBondingValue,
	}
}

// NewDepositInfoByResourceNodeAddr creates a new instance of DepositInfoByNodeAddr
func NewDepositInfoByResourceNodeAddr(
	bondDenom string,
	resourceNode ResourceNode,
	unBondingDeposit sdk.Int,
	unBondedDeposit sdk.Int,
	bondedDeposit sdk.Int,

) DepositInfo {
	bonedValue := sdk.NewCoin(bondDenom, bondedDeposit)
	unBondedValue := sdk.NewCoin(bondDenom, unBondedDeposit)
	unBondingValue := sdk.NewCoin(bondDenom, unBondingDeposit)

	return DepositInfo{
		NetworkAddress:   resourceNode.GetNetworkAddress(),
		Pubkey:           resourceNode.GetPubkey(),
		Suspend:          resourceNode.GetSuspend(),
		Status:           resourceNode.GetStatus(),
		Tokens:           resourceNode.Tokens,
		OwnerAddress:     resourceNode.GetOwnerAddress(),
		Description:      resourceNode.GetDescription(),
		NodeType:         resourceNode.GetNodeType(),
		CreationTime:     resourceNode.GetCreationTime(),
		UnBondingDeposit: unBondingValue,
		UnBondedDeposit:  unBondedValue,
		BondedDeposit:    bonedValue,
	}
}

// NewDepositInfoByMetaNodeAddr creates a new instance of DepositInfoByNodeAddr
func NewDepositInfoByMetaNodeAddr(
	bondDenom string,
	metaNode MetaNode,
	unBondingDeposit sdk.Int,
	unBondedDeposit sdk.Int,
	bondedDeposit sdk.Int,
) DepositInfo {
	bonedValue := sdk.NewCoin(bondDenom, bondedDeposit)
	unBondedValue := sdk.NewCoin(bondDenom, unBondedDeposit)
	unBondingValue := sdk.NewCoin(bondDenom, unBondingDeposit)
	return DepositInfo{
		NetworkAddress:   metaNode.GetNetworkAddress(),
		Pubkey:           metaNode.GetPubkey(),
		Suspend:          metaNode.Suspend,
		Status:           metaNode.Status,
		Tokens:           metaNode.Tokens,
		OwnerAddress:     metaNode.GetOwnerAddress(),
		Description:      metaNode.Description,
		NodeType:         uint32(0),
		CreationTime:     metaNode.CreationTime,
		UnBondingDeposit: unBondingValue,
		UnBondedDeposit:  unBondedValue,
		BondedDeposit:    bonedValue,
	}
}

type DepositInfos []DepositInfo

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v DepositInfo) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(v.Pubkey, &pk)
}

func (v DepositInfos) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range v {
		if err := v[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}
