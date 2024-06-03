package types

import (
	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryTypeAll      = 0
	QueryTypeSP       = 1
	QueryTypePP       = 2
	QueryDefaultLimit = 100
)

// NewQueryDepositTotalInfo creates a new instance of QueryDepositTotalResponse
func NewQueryDepositTotalInfo(bondDenom string, ResourceNodeTotalDeposit, MetaNodeTotalDeposit, totalBondedDeposit,
	totalUnbondedDeposit, totalUnbondingDeposit sdkmath.Int) *QueryDepositTotalResponse {

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
	unBondingDeposit sdkmath.Int,
	unBondedDeposit sdkmath.Int,
	bondedDeposit sdkmath.Int,

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
	unBondingDeposit sdkmath.Int,
	unBondedDeposit sdkmath.Int,
	bondedDeposit sdkmath.Int,
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

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v DepositInfo) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(v.Pubkey, &pk)
}
