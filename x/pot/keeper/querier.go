package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	abci "github.com/tendermint/tendermint/abci/types"
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
		case QueryPotRewardsByEpoch:
			return queryPotRewardsByEpoch(ctx, req, k)
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

	reportRecord := k.GetVolumeReport(ctx, sdk.NewInt(epoch))
	if reportRecord.TxHash == "" {
		bz := []byte(fmt.Sprintf("no volume report at epoch: %d", epoch))
		return bz, nil
	}
	bz, err := codec.MarshalJSONIndent(k.cdc, reportRecord)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

// queryPotRewardsByEpoch fetches total rewards and owner individual rewards from traffic and mining.
func queryPotRewardsByEpoch(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params QueryPotRewardsByEpochParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	potEpochRewards := k.getPotRewardsByEpoch(ctx, params)
	if len(potEpochRewards) < 1 {
		bz, _ := codec.MarshalJSONIndent(k.cdc, fmt.Sprintf("no Pot rewards information at epoch: %s", params.Epoch.String()))
		return bz, nil
	}
	bz, err := codec.MarshalJSONIndent(k.cdc, potEpochRewards)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func (k Keeper) getPotRewardsByEpoch(ctx sdk.Context, params QueryPotRewardsByEpochParams) (res []types.Reward) {
	res = k.getRewardsResult(ctx, params, params.NodeVolumes)
	start, end := client.Paginate(len(res), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		return nil
	} else {
		res = res[start:end]
		return res
	}
}

func (k Keeper) getRewardsResult(ctx sdk.Context, params QueryPotRewardsByEpochParams, nodesVolume []types.SingleNodeVolume) (res []types.Reward) {
	rewardDetailMap := k.tempCalculateNodePotRewards(ctx, nodesVolume)

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
	return
}

func (k Keeper) tempCalculateNodePotRewards(ctx sdk.Context, nodesVolume []types.SingleNodeVolume) map[string]types.Reward {
	distributeGoal := types.InitDistributeGoal()
	rewardDetailMap := make(map[string]types.Reward) //key: node address

	_, distributeGoal, err := k.CalcTrafficRewardInTotal(ctx, nodesVolume, distributeGoal)
	if err != nil {
		return nil
	}
	distributeGoal, err = k.CalcMiningRewardInTotal(ctx, distributeGoal)
	if err != nil && err != types.ErrOutOfIssuance {
		return nil
	}

	distributeGoalBalance := distributeGoal
	rewardDetailMap, distributeGoalBalance = k.CalcRewardForResourceNode(ctx, nodesVolume, distributeGoalBalance, rewardDetailMap)
	rewardDetailMap, distributeGoalBalance = k.CalcRewardForIndexingNode(ctx, distributeGoalBalance, rewardDetailMap)
	return rewardDetailMap
}

func queryPotRewardsWithOwnerHeight(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params QueryPotRewardsWithOwnerHeightParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	ownerRewards := k.getOwnerRewards(ctx, params)
	if len(ownerRewards) == 0 {
		ownerRewards = nil
	}

	record := OwnerRewardsRecord{ownerRewards}
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

func (k Keeper) getOwnerRewards(ctx sdk.Context, params QueryPotRewardsWithOwnerHeightParams) (res []NodeRewardsInfo) {
	nodeOwnerMap := k.getNodeOwnerMap(ctx)
	var r NodeRewardsInfo
	for nodeAccStr, OwnerAcc := range nodeOwnerMap {
		if OwnerAcc.Equals(params.OwnerAddr) {
			nodeAcc, err := sdk.AccAddressFromBech32(nodeAccStr)
			if err != nil {
				return nil
			}
			r.ImmatureTotalReward = sdk.NewCoin(k.BondDenom(ctx), k.GetImmatureTotalReward(ctx, nodeAcc))
			r.MatureTotalReward = sdk.NewCoin(k.BondDenom(ctx), k.GetMatureTotalReward(ctx, nodeAcc))
			r.NodeAddress = nodeAcc
			res = append(res, r)
		}

	}

	start, end := client.Paginate(len(res), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		return nil
	} else {
		res = res[start:end]
		return
	}
}
