package keeper

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// QueryPotRewardsParams Params for query 'custom/pot/potRewards'
type QueryPotRewardsParams struct {
	Page     int
	Limit    int
	NodeAddr sdk.AccAddress
	Epoch    sdk.Int
}

// NewQueryPotRewardsParams creates a new instance of QueryNodesParams
func NewQueryPotRewardsParams(page, limit int, nodeAddr sdk.AccAddress, epoch int64) QueryPotRewardsParams {
	return QueryPotRewardsParams{
		Page:     page,
		Limit:    limit,
		NodeAddr: nodeAddr,
		Epoch:    sdk.NewInt(epoch),
	}
}

type NodeRewardsInfo struct {
	NodeAddr          sdk.AccAddress
	FoundationAccount sdk.AccAddress
	Epoch             sdk.Int
	TotalMinedTokens  sdk.Coin
	MinedTokens       sdk.Coin
	IndividualRewards sdk.Coin
	MatureTotal       sdk.Coin
	ImmatureTotal     sdk.Coin
}

// NewNodeRewardsInfo creates a new instance of NodeRewardsInfo
func NewNodeRewardsInfo(nodeAddr, foundationAccount sdk.AccAddress, epoch, totalMinedTokens, minedTokens, individualRewards, matureTotal, immatureTotal sdk.Int) NodeRewardsInfo {
	denomName := "ustos"
	return NodeRewardsInfo{
		NodeAddr:          nodeAddr,
		FoundationAccount: foundationAccount,
		Epoch:             epoch,
		TotalMinedTokens:  sdk.NewCoin(denomName, totalMinedTokens),
		MinedTokens:       sdk.NewCoin(denomName, minedTokens),
		IndividualRewards: sdk.NewCoin(denomName, individualRewards),
		MatureTotal:       sdk.NewCoin(denomName, matureTotal),
		ImmatureTotal:     sdk.NewCoin(denomName, immatureTotal),
	}
}

func (k Keeper) GetResourceNodesRewards(ctx sdk.Context, params QueryPotRewardsParams) (res []NodeRewardsInfo) {

	rewardAddrList := k.GetRewardAddressPool(ctx)

	for _, n := range rewardAddrList {
		// match NodeAddr (if supplied)
		if !params.NodeAddr.Equals(sdk.AccAddress{}) {
			if !n.Equals(params.NodeAddr) {
				continue
			}
		}

		foundationAccount := k.GetFoundationAccount(ctx)
		totalMinedTokens := k.GetTotalMinedTokens(ctx)
		minedTokens := k.GetMinedTokens(ctx, params.Epoch)
		individualRewards := k.GetIndividualReward(ctx, n, params.Epoch)
		matureTotal := k.GetMatureTotalReward(ctx, n)
		immatureTotal := k.GetImmatureTotalReward(ctx, n)

		individualResult := NewNodeRewardsInfo(
			params.NodeAddr,
			foundationAccount,
			params.Epoch,
			totalMinedTokens,
			minedTokens,
			individualRewards,
			matureTotal,
			immatureTotal,
		)

		res = append(res, individualResult)
	}

	start, end := client.Paginate(len(res), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		return []NodeRewardsInfo{}
	} else {
		res = res[start:end]
		return res
	}
}
