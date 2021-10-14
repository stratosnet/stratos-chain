package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	"sort"
)

type QueryPotRewardsParams struct {
	Page     int
	Limit    int
	NodeAddr sdk.AccAddress
	Epoch    sdk.Int
}

func sortDetailMapToSlice(rewardDetailMap map[string]types.Reward) (rewardDetailList []types.Reward) {
	keys := make([]string, 0, len(rewardDetailMap))
	for key, _ := range rewardDetailMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		reward := rewardDetailMap[key]
		rewardDetailList = append(rewardDetailList, reward)
	}
	return rewardDetailList
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

type QueryPotRewardsByEpochParams struct {
	Page        int
	Limit       int
	OwnerAddr   sdk.AccAddress
	Epoch       sdk.Int
	NodeVolumes []types.SingleNodeVolume
}

// NewQueryPotRewardsByEpochParams creates a new instance of QueryPotRewardsParams
func NewQueryPotRewardsByEpochParams(page, limit int, ownerAddr sdk.AccAddress, epoch sdk.Int, nodeVolumes []types.SingleNodeVolume) QueryPotRewardsByEpochParams {
	return QueryPotRewardsByEpochParams{
		Page:        page,
		Limit:       limit,
		OwnerAddr:   ownerAddr,
		Epoch:       epoch,
		NodeVolumes: nodeVolumes,
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
	NodeDetails []NodeRewardsInfo
}

type NodeRewardsRecord struct {
	Record NodeRewardsInfo
}

// NewNodeRewardsRecord creates a new instance of NodeRewardsRecord
func NewNodeRewardsRecord(
	record NodeRewardsInfo,

) NodeRewardsRecord {
	return NodeRewardsRecord{
		Record: record,
	}
}
