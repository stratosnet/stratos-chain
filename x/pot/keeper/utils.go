package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type QueryPotRewardsParams struct {
	Page     int
	Limit    int
	NodeAddr sdk.AccAddress
	Epoch    sdk.Int
}

// NewQueryPotRewardsParams creates a new instance of QueryPotRewardsParams
func NewQueryPotRewardsParams(page, limit int, nodeAddr sdk.AccAddress, epoch sdk.Int) QueryPotRewardsParams {
	return QueryPotRewardsParams{
		Page:     page,
		Limit:    limit,
		NodeAddr: nodeAddr,
		Epoch:    epoch,
	}
}

type QueryPotRewardsByepochParams struct {
	Page      int
	Limit     int
	OwnerAddr sdk.AccAddress
	Epoch     sdk.Int
}

// NewQueryPotRewardsByepochParams creates a new instance of QueryPotRewardsParams
func NewQueryPotRewardsByepochParams(page, limit int, ownerAddr sdk.AccAddress, epoch sdk.Int) QueryPotRewardsByepochParams {
	return QueryPotRewardsByepochParams{
		Page:      page,
		Limit:     limit,
		OwnerAddr: ownerAddr,
		Epoch:     epoch,
	}
}

type QueryPotRewardsByOwnerParams struct {
	Page      int
	Limit     int
	OwnerAddr sdk.AccAddress
	Height    int64
}

// NewQueryPotRewardsByOwnerParams creates a new instance of QueryPotRewardsParams
func NewQueryPotRewardsByOwnerParams(page, limit int, ownerAddr sdk.AccAddress, height int64) QueryPotRewardsByOwnerParams {
	return QueryPotRewardsByOwnerParams{
		Page:      page,
		Limit:     limit,
		OwnerAddr: ownerAddr,
		Height:    height,
	}
}

type QueryPotRewardsWithOwnerHeightParams struct {
	Page      int
	Limit     int
	OwnerAddr sdk.AccAddress
	Height    int64
}

func NewQueryPotRewardsWithOwnerHeightParams(page, limit int, ownerAddr sdk.AccAddress, height int64) QueryPotRewardsWithOwnerHeightParams {
	return QueryPotRewardsWithOwnerHeightParams{
		Page:      page,
		Limit:     limit,
		OwnerAddr: ownerAddr,
		Height:    height,
	}
}

type NodeRewardsInfo struct {
	NodeAddress         sdk.AccAddress
	MatureTotalReward   sdk.Coin
	ImmatureTotalReward sdk.Coin
}

// NewNodeRewardsInfo creates a new instance of NodeRewardsInfo
func NewNodeRewardsInfo(
	nodeAddress sdk.AccAddress,
	matureTotal,
	immatureTotal sdk.Int,
) NodeRewardsInfo {
	denomName := "ustos"
	return NodeRewardsInfo{
		NodeAddress:         nodeAddress,
		MatureTotalReward:   sdk.NewCoin(denomName, matureTotal),
		ImmatureTotalReward: sdk.NewCoin(denomName, immatureTotal),
	}
}

type OwnerRewardsRecord struct {
	PotRewardsRecordHeight int64
	PotRewardsRecordEpoch  sdk.Int
	NodeDetails            []NodeRewardsInfo
}

type NodeRewardsRecord struct {
	PotRewardsRecordHeight int64
	PotRewardsRecordEpoch  sdk.Int
	Record                 NodeRewardsInfo
}

// NewNodeRewardsRecord creates a new instance of NodeRewardsRecord
func NewNodeRewardsRecord(
	potRewardsRecordHeight int64,
	potRewardsRecordEpoch sdk.Int,
	record NodeRewardsInfo,

) NodeRewardsRecord {
	return NodeRewardsRecord{
		PotRewardsRecordHeight: potRewardsRecordHeight,
		PotRewardsRecordEpoch:  potRewardsRecordEpoch,
		Record:                 record,
	}
}
