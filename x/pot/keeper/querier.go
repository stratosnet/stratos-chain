package keeper

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
)

const (
	QueryVolumeReport      = "volume_report"
	QueryPotRewardsByEpoch = "pot_rewards_by_epoch"
	QueryPotRewardsByOwner = "pot_rewards_by_owner"
	QueryDefaultLimit      = 100
)

// NewQuerier creates a new querier for pot clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryVolumeReport:
			return queryVolumeReport(ctx, req, k)
		//case QueryPotRewards:
		//	return queryPotRewards(ctx, req, k)
		case QueryPotRewardsByEpoch:
			return queryPotRewardsByEpoch(ctx, req, k)
		//case QueryPotRewardsByOwner:
		//	return queryPotRewardsByOwner(ctx, req, k)
		case QueryPotRewardsByOwner:
			return queryPotRewardsWithOwnerHeight(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown pot query endpoint")
		}
	}
}

// queryVolumeReport fetches a hash of report volume for the supplied epoch.
func queryVolumeReport(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	epoch, err := strconv.ParseInt(string(req.Data), 10, 64)
	if err != nil {
		return nil, err
	}

	reportRecord, err := k.GetVolumeReport(ctx, sdk.NewInt(epoch))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	bz, err := codec.MarshalJSONIndent(k.cdc, reportRecord)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

// queryPotRewardsByEpoch fetches total rewards and owner individual rewards from traffic and mining.
func queryPotRewardsByEpoch(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params QueryPotRewardsByepochParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	potEpochRewards := k.getPotRewardsByEpoch(ctx, params)
	bz, err := codec.MarshalJSONIndent(k.cdc, potEpochRewards)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func (k Keeper) getPotRewardsByEpoch(ctx sdk.Context, params QueryPotRewardsByepochParams) (res []types.Reward) {
	// get volume report based on the given epoch
	reportRecord, err := k.GetVolumeReport(ctx, params.Epoch)
	if err != nil {
		return nil
	}

	res = k.getRewardsResult(ctx, params, reportRecord)

	start, end := client.Paginate(len(res), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		return nil
	} else {
		res = res[start:end]
		return res
	}
}

func (k Keeper) getRewardsResult(ctx sdk.Context, params QueryPotRewardsByepochParams, reportRecord types.ReportRecord) (res []types.Reward) {
	rewardDetailMap := k.tempClaculateNodePotRewards(ctx, reportRecord)

	for _, value := range rewardDetailMap {
		if !params.OwnerAddr.Empty() {
			nodeOwnerMap := make(map[string]sdk.AccAddress)

			nodeOwnerMap = k.RegisterKeeper.GetNodeOwnerMapFromIndexingNodes(ctx, nodeOwnerMap)
			if ownerAddr, ok := nodeOwnerMap[value.NodeAddress.String()]; ok {
				if ownerAddr.Equals(params.OwnerAddr) {
					res = append(res, value)
				}
			} else {
				nodeOwnerMap = k.RegisterKeeper.GetNodeOwnerMapFromResourceNodes(ctx, nodeOwnerMap)
				if ownerAddr, ok := nodeOwnerMap[value.NodeAddress.String()]; ok {
					if ownerAddr.Equals(params.OwnerAddr) {
						res = append(res, value)
					}
				}
			}

		} else {
			res = append(res, value)
		}

	}
	return res
}

func (k Keeper) tempClaculateNodePotRewards(ctx sdk.Context, reportRecord types.ReportRecord) map[string]types.Reward {
	distributeGoal := types.InitDistributeGoal()
	rewardDetailMap := make(map[string]types.Reward) //key: node address

	//1, calc traffic reward in total
	_, distributeGoal, err := k.CalcTrafficRewardInTotal(ctx, reportRecord.NodesVolume, distributeGoal)
	if err != nil {
		return nil
	}

	//2, calc mining reward in total
	distributeGoal, err = k.CalcMiningRewardInTotal(ctx, distributeGoal)
	if err != nil && err != types.ErrOutOfIssuance {
		return nil
	}

	distributeGoalBalance := distributeGoal

	//3, calc reward for resource node
	rewardDetailMap, distributeGoalBalance = k.CalcRewardForResourceNode(ctx, reportRecord.NodesVolume, distributeGoalBalance, rewardDetailMap)

	//4, calc reward from indexing node
	rewardDetailMap, distributeGoalBalance = k.CalcRewardForIndexingNode(ctx, distributeGoalBalance, rewardDetailMap)
	return rewardDetailMap
}

//func getFilteredNodesAddrByOwner(ctx sdk.Context, ownerAddress sdk.AccAddress, k Keeper) []sdk.AccAddress {
//	resourceNodesAddr := k.RegisterKeeper.GetAllResourceNodes(ctx)
//	indexingNodesAddr := k.RegisterKeeper.GetAllIndexingNodes(ctx)
//	filteredNodesAddr := make([]sdk.AccAddress, 0, len(resourceNodesAddr)+len(indexingNodesAddr))
//
//	for _, n := range resourceNodesAddr {
//		// match OwnerAddr (if supplied)
//		if ownerAddress.Empty() || n.OwnerAddress.Equals(ownerAddress) {
//			filteredNodesAddr = append(filteredNodesAddr, sdk.AccAddress(n.PubKey.Address()))
//		}
//
//	}
//	for _, n := range indexingNodesAddr {
//		// match OwnerAddr (if supplied)
//		if ownerAddress.Empty() || n.OwnerAddress.Equals(ownerAddress) {
//			filteredNodesAddr = append(filteredNodesAddr, sdk.AccAddress(n.PubKey.Address()))
//		}
//	}
//	return filteredNodesAddr
//}

func queryPotRewardsWithOwnerHeight(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params QueryPotRewardsWithOwnerHeightParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	recordHeight, recordEpoch, ownerRewards := k.getOwnerRewards(ctx, params)
	if len(ownerRewards) == 0 {
		ownerRewards = nil
	}

	record := OwnerRewardsRecord{recordHeight, recordEpoch, ownerRewards}
	if len(record.NodeDetails) < 1 {
		bz, _ := codec.MarshalJSONIndent(k.cdc, "No Pot rewards information at this height")
		return bz, nil
	}
	bz, err := codec.MarshalJSONIndent(k.cdc, record)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func (k Keeper) getOwnerRewards(ctx sdk.Context, params QueryPotRewardsWithOwnerHeightParams) (recordHeight int64, recordEpoch sdk.Int, res []NodeRewardsInfo) {
	recordHeight, recordEpoch, res = k.GetPotRewardRecords(ctx, params)

	start, end := client.Paginate(len(res), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		return 0, sdk.ZeroInt(), nil
	} else {
		res = res[start:end]
		return
	}
}
