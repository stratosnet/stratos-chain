package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
)

const (
	QueryVolumeReport      = "volume_report"
	QueryPotRewards        = "pot_rewards"
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
		case QueryPotRewards:
			return queryPotRewards(ctx, req, k)
		case QueryPotRewardsByEpoch:
			return queryPotRewardsByEpoch(ctx, req, k)
		case QueryPotRewardsByOwner:
			return queryPotRewardsByOwner(ctx, req, k)
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
	bz, err := codec.MarshalJSONIndent(k.Cdc, reportRecord)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

// queryPotRewards fetches total rewards and owner individual rewards from traffic and mining.
func queryPotRewards(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params QueryPotRewardsParams
	err := k.Cdc.UnmarshalJSON(req.Data, &params)
	ctx.Logger().Info("params", "params", params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	nodeRewards := k.GetNodesRewards(ctx, params)
	if len(nodeRewards) == 0 {
		nodeRewards = []NodeRewardsInfo{}
	}

	bz, err := codec.MarshalJSONIndent(k.Cdc, nodeRewards)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// queryPotRewardsByEpoch fetches total rewards and owner individual rewards from traffic and mining.
func queryPotRewardsByEpoch(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params QueryPotRewardsByepochParams
	err := k.Cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	potEpochRewards := k.GetPotRewardsByEpoch(ctx, params)

	bz, err := codec.MarshalJSONIndent(k.Cdc, potEpochRewards)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryPotRewardsByOwner(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params QueryPotRewardsByOwnerParams
	err := k.Cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	ownerRewards := k.GetNodesRewardsByOwner(ctx, params)
	if len(ownerRewards) == 0 {
		ownerRewards = nil
	}

	bz, err := codec.MarshalJSONIndent(k.Cdc, ownerRewards)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
