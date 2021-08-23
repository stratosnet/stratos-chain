package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

const (
	Oz2uoz = 1000000000
)

// QueryPotRewardsParams Params for query 'custom/register/resource-nodes'
type QueryPotRewardsParams struct {
	Page     int
	Limit    int
	NodeAddr sdk.AccAddress
}

// NewQueryPotRewardsParams creates a new instance of QueryNodesParams
func NewQueryPotRewardsParams(page, limit int, ownerAddr sdk.AccAddress) QueryPotRewardsParams {
	return QueryPotRewardsParams{
		Page:     page,
		Limit:    limit,
		NodeAddr: ownerAddr,
	}
}

func (k Keeper) GetResourceNodesRewards(ctx sdk.Context, params QueryPotRewardsParams) ([]types.Reward, error) {

	rewardAddrList := k.GetRewardAddressPool(ctx)

	for _, n := range rewardAddrList {
		// match OwnerAddr (if supplied)
		if len(params.NodeAddr) > 0 {
			if !n.Equals(params.NodeAddr) {
				continue
			}
		}

	}

	// use registerKeeper.NewQueryNodesParams to find all filtered nodes
	//newQueryResNodeParams := registerKeeper.NewQueryNodesParams(params.Page, params.Limit, "", "", params.OwnerAddr)

	//resNodes := k.RegisterKeeper.GetResourceNodesFiltered(ctx, newQueryResNodeParams)
	//resNodeRewards := make([]registerTypes.ResourceNode, 0, len(resNodes))

	//err := k.DistributePotReward(ctx, trafficList, epoch2017)
	//if err != nil {
	//	return []types.Reward{}, sdkerrors.ErrUnknownRequest
	//}
	//rewardAddrList := k.GetRewardAddressPool(ctx)
	//k.Logger("address pool: ")
	//for i := 0; i < len(rewardAddrList); i++ {
	//	fmt.Println(rewardAddrList[i].String() + ", ")
	//}
	//
	//
	//rewardDetailMap := make(map[string]types.Reward)
	//
	//for i, n := range resNodes {
	//	distributeGoal := types.InitDistributeGoal()
	//
	//	//build traffic list
	//	trafficList = append(trafficList, types.NewSingleNodeVolume(n.OwnerAddress, sdk.NewInt(ResourceNodeVolume[i])))
	//	rewardDetailMap, distributeGoal = k.CalcRewardForResourceNode(ctx, trafficList, distributeGoal, rewardDetailMap)
	//	totalReward := k.getTrafficReward(ctx, trafficList)
	//
	//	idvRwdResNode1Ep1 := k.GetIndividualReward(ctx, addrRes1, epoch1)
	//	matureTotalResNode1 := k.GetMatureTotalReward(ctx, addrRes1)
	//	immatureTotalResNode1 := k.GetImmatureTotalReward(ctx, addrRes1)
	//	fmt.Println("resourceNode1: address = " + addrRes1.String() + ", individual = " + idvRwdResNode1Ep1.String() + ",\tmatureTotal = " + matureTotalResNode1.String() + ",\timmatureTotal = " + immatureTotalResNode1.String())
	//	require.Equal(t, idvRwdResNode1Ep1, sdk.NewInt(131089476265))

	//	filteredNodes = append(filteredNodes, n)
	//}
	//
	//start, end := client.Paginate(len(filteredNodes), params.Page, params.Limit, QueryDefaultLimit)
	//if start < 0 || end < 0 {
	//	filteredNodes = []types.ResourceNode{}
	//} else {
	//	filteredNodes = filteredNodes[start:end]
	//return resNodeRewards
	return nil, sdkerrors.Error{}
}
