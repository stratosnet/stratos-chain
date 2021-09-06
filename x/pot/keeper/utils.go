package keeper

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// QueryPotRewardsParams Params for query 'custom/pot/rewards'
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

// QueryPotRewardsByepochParams Params for query 'custom/pot/rewards'
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

type NodeRewardsInfo struct {
	NodeAddress sdk.AccAddress
	//FoundationAccount   sdk.AccAddress
	Epoch sdk.Int
	//LastMaturedEpoch    sdk.Int
	//TotalUnissuedPrepay sdk.Int
	//TotalMinedTokens    sdk.Coin
	//MinedTokens         sdk.Coin
	IndividualRewards   sdk.Coin
	MatureTotalReward   sdk.Coin
	ImmatureTotalReward sdk.Coin
}

//
//// QueryVolumeReportParams Params for query 'custom/pot/report'
//type QueryVolumeReportParams struct {
//	Epoch sdk.Int
//}
//
//// NewQueryVolumeReportParams creates a new instance of QueryVolumeReportParams
//func NewQueryVolumeReportParams(epoch int64) QueryVolumeReportParams {
//	return QueryVolumeReportParams{
//		Epoch: sdk.NewInt(epoch),
//	}
//}

// NewNodeRewardsInfo creates a new instance of NodeRewardsInfo
func NewNodeRewardsInfo(
	NodeAddress sdk.AccAddress,
	//foundationAccount sdk.AccAddress,
	//epoch,
	//lastMaturedEpoch,
	//totalUnissuedPrepay,
	//totalMinedTokens,
	//minedTokens,
	individualRewards,
	matureTotal,
	immatureTotal sdk.Int,
) NodeRewardsInfo {
	denomName := "ustos"
	return NodeRewardsInfo{
		NodeAddress: NodeAddress,
		//FoundationAccount:   foundationAccount,
		//Epoch: epoch,
		//LastMaturedEpoch:    lastMaturedEpoch,
		//TotalUnissuedPrepay: totalUnissuedPrepay,
		//TotalMinedTokens:    sdk.NewCoin(denomName, totalMinedTokens),
		//MinedTokens:         sdk.NewCoin(denomName, minedTokens),
		IndividualRewards:   sdk.NewCoin(denomName, individualRewards),
		MatureTotalReward:   sdk.NewCoin(denomName, matureTotal),
		ImmatureTotalReward: sdk.NewCoin(denomName, immatureTotal),
	}
}

func (k Keeper) GetNodesRewards(ctx sdk.Context, params QueryPotRewardsParams) (res []NodeRewardsInfo) {

	rewardAddrList := k.GetRewardAddressPool(ctx)

	for _, n := range rewardAddrList {
		// match NodeAddr (if supplied)
		//if !params.NodeAddr.Equals(sdk.AccAddress{}) {
		//	if !n.Equals(params.NodeAddr) {
		//		continue
		//	}
		//} else {
		//	continue
		//}
		ctx.Logger().Info("n", "n", n)
		if !(n.Equals(params.NodeAddr)) {
			continue
		}
		ctx.Logger().Info("equal", "equal", true)

		//foundationAccount := k.GetFoundationAccount(ctx)
		//totalMinedTokens := k.GetTotalMinedTokens(ctx)
		//minedTokens := k.GetMinedTokens(ctx, params.Epoch)

		individualRewards := k.GetIndividualReward(ctx, n, params.Epoch)
		matureTotal := k.GetMatureTotalReward(ctx, n)
		immatureTotal := k.GetImmatureTotalReward(ctx, n)
		//lastMaturedEpoch := k.getLastMaturedEpoch(ctx)
		//totalUnissuedPrepay := k.GetTotalUnissuedPrepay(ctx)

		individualResult := NewNodeRewardsInfo(
			n,
			//foundationAccount,
			//params.Epoch,
			//lastMaturedEpoch,
			//totalUnissuedPrepay,
			//totalMinedTokens,
			//minedTokens,
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

func (k Keeper) GetPotRewardsByEpoch(ctx sdk.Context, params QueryPotRewardsByepochParams) (res []types.Reward) {
	resourceNodesAddr := k.RegisterKeeper.GetAllResourceNodes(ctx)
	indexingNodesAddr := k.RegisterKeeper.GetAllIndexingNodes(ctx)
	filteredNodesAddrStr := make([]string, 0, len(resourceNodesAddr)+len(indexingNodesAddr))
	//filteredIndexingNodes := make([]registerTypes.IndexingNode, 0, len(indexingNodes))

	for _, n := range resourceNodesAddr {
		// match OwnerAddr (if supplied)
		if !params.OwnerAddr.Empty() {
			if !n.OwnerAddress.Equals(params.OwnerAddr) {
				continue
			}
		}
		filteredNodesAddrStr = append(filteredNodesAddrStr, sdk.AccAddress(n.PubKey.Address()).String())
	}
	for _, n := range indexingNodesAddr {
		// match OwnerAddr (if supplied)
		if !params.OwnerAddr.Empty() {
			if !n.OwnerAddress.Equals(params.OwnerAddr) {
				continue
			}
		}
		filteredNodesAddrStr = append(filteredNodesAddrStr, sdk.AccAddress(n.PubKey.Address()).String())
	}

	epochRewards := k.GetEpochReward(ctx, params.Epoch)
	ctx.Logger().Info("epochRewards", "epochRewards", epochRewards)
	for _, v := range epochRewards {
		if stringInSlice(v.NodeAddress.String(), filteredNodesAddrStr) {
			newNodeReward := types.NewReward(v.NodeAddress, v.RewardFromMiningPool, v.RewardFromTrafficPool)
			res = append(res, newNodeReward)
		}

	}
	ctx.Logger().Info("res", "res", res)
	start, end := client.Paginate(len(res), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		return nil
	} else {
		res = res[start:end]
		return res
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
