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
	potEpochRewards := k.GetPotRewardsByEpoch(ctx, params)
	bz, err := codec.MarshalJSONIndent(k.cdc, potEpochRewards)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func (k Keeper) GetPotRewardsByEpoch(ctx sdk.Context, params QueryPotRewardsByepochParams) (res []types.Reward) {
	filteredNodesAddr := getFilteredNodesAddrByOwner(ctx, params.OwnerAddr, k)

	epochRewards := k.GetEpochReward(ctx, params.Epoch)
	epochRewardsMap := make(map[string]types.Reward)
	for _, v := range epochRewards {
		epochRewardsMap[v.NodeAddress.String()] = v
	}

	for _, n := range filteredNodesAddr {
		if newNodeReward, found := epochRewardsMap[n.String()]; found {
			res = append(res, newNodeReward)
		}
	}
	start, end := client.Paginate(len(res), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		return nil
	} else {
		res = res[start:end]
		return res
	}
}

func getFilteredNodesAddrByOwner(ctx sdk.Context, ownerAddress sdk.AccAddress, k Keeper) []sdk.AccAddress {
	resourceNodesAddr := k.RegisterKeeper.GetAllResourceNodes(ctx)
	indexingNodesAddr := k.RegisterKeeper.GetAllIndexingNodes(ctx)
	filteredNodesAddr := make([]sdk.AccAddress, 0, len(resourceNodesAddr)+len(indexingNodesAddr))

	for _, n := range resourceNodesAddr {
		// match OwnerAddr (if supplied)
		if ownerAddress.Empty() || n.OwnerAddress.Equals(ownerAddress) {
			filteredNodesAddr = append(filteredNodesAddr, sdk.AccAddress(n.PubKey.Address()))
		}

	}
	for _, n := range indexingNodesAddr {
		// match OwnerAddr (if supplied)
		if ownerAddress.Empty() || n.OwnerAddress.Equals(ownerAddress) {
			filteredNodesAddr = append(filteredNodesAddr, sdk.AccAddress(n.PubKey.Address()))
		}
	}
	return filteredNodesAddr
}

func queryPotRewardsWithOwnerHeight(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params QueryPotRewardsWithOwnerHeightParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	recordHeight, recordEpoch, ownerRewards := k.getOwnerRewards(ctx, params)
	if len(ownerRewards) == 0 {
		ownerRewards = nil
	}

	record := OwnerRewardsRecord{recordHeight, recordEpoch, ownerRewards}
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
