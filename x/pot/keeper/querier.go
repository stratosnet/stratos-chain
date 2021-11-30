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
	QueryVolumeReport            = "query_volume_report"
	QueryPotRewardsByReportEpoch = "query_pot_rewards_by_report_epoch"
	QueryPotRewardsByWalletAddr  = "query_pot_rewards_by_wallet_address"
	QueryDefaultLimit            = 100
)

// NewQuerier creates a new querier for pot clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryVolumeReport:
			return queryVolumeReport(ctx, req, k)
		case QueryPotRewardsByReportEpoch:
			return queryPotRewardsByReportEpoch(ctx, req, k)
		case QueryPotRewardsByWalletAddr:
			return queryPotRewardsByWalletAddress(ctx, req, k)
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

// queryPotRewardsByReportEpoch fetches total rewards and owner individual rewards from traffic and mining.
func queryPotRewardsByReportEpoch(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryPotRewardsByReportEpochParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	potEpochRewards := k.getPotRewardsByReportEpoch(ctx, params)
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

func (k Keeper) getPotRewardsByReportEpoch(ctx sdk.Context, params types.QueryPotRewardsByReportEpochParams) (res []types.Reward) {
	matureEpoch := params.Epoch.Add(sdk.NewInt(k.MatureEpoch(ctx)))

	if !params.WalletAddress.Empty() {
		reward, found := k.GetIndividualReward(ctx, params.WalletAddress, matureEpoch)
		if found {
			res = append(res, reward)
		}
	} else {
		rewardAddressPool := k.GetRewardAddressPool(ctx)
		for _, walletAddress := range rewardAddressPool {
			reward, found := k.GetIndividualReward(ctx, walletAddress, matureEpoch)
			if found {
				res = append(res, reward)
			}
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

func queryPotRewardsByWalletAddress(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryPotRewardsByWalletAddrParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	immatureTotalReward := sdk.NewCoin(k.BondDenom(ctx), k.GetImmatureTotalReward(ctx, params.WalletAddr))
	matureTotalReward := sdk.NewCoin(k.BondDenom(ctx), k.GetMatureTotalReward(ctx, params.WalletAddr))
	reward := types.NewPotRewardInfo(params.WalletAddr, immatureTotalReward, matureTotalReward)

	bz, err := codec.MarshalJSONIndent(k.cdc, reward)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
