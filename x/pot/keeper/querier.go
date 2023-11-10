package keeper

import (
	"fmt"
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/stratosnet/stratos-chain/client"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// NewQuerier creates a new querier for pot clients.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryVolumeReport:
			return queryVolumeReport(ctx, req, k, legacyQuerierCdc)
		case types.QueryIndividualRewardsByReportEpoch:
			return queryIndividualRewardsByReportEpoch(ctx, req, k, legacyQuerierCdc)
		case types.QueryRewardsByWalletAddr:
			return queryRewardsByWalletAddress(ctx, req, k, legacyQuerierCdc)
		case types.QuerySlashingByWalletAddr:
			return querySlashingByWalletAddress(ctx, req, k, legacyQuerierCdc)
		case types.QueryPotParams:
			return getPotParams(ctx, req, k, legacyQuerierCdc)
		case types.QueryTotalMinedToken:
			return getTotalMinedToken(ctx, req, k, legacyQuerierCdc)
		case types.QueryCirculationSupply:
			return getCirculationSupply(ctx, req, k, legacyQuerierCdc)
		case types.QueryTotalRewardByEpoch:
			return getTotalRewardByEpoch(ctx, req, k, legacyQuerierCdc)
		case types.QueryMetrics:
			return getMetrics(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown pot query endpoint")
		}
	}
}

func getPotParams(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryVolumeReport fetches a hash of report volume for the supplied epoch.
func queryVolumeReport(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	epoch, err := strconv.ParseInt(string(req.Data), 10, 64)
	if err != nil {
		return []byte{}, err
	}

	reportRecord := k.GetVolumeReport(ctx, sdk.NewInt(epoch))
	if reportRecord.TxHash == "" {
		e := sdkerrors.Wrapf(types.ErrCannotFindReport,
			fmt.Sprintf("no volume report found at epoch %d. Current epoch is %s",
				epoch, k.GetLastDistributedEpoch(ctx).String()))
		return []byte{}, e
	}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, reportRecord)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

// queryPotRewardsByReportEpoch fetches individual rewards from traffic and mining pool.
func queryIndividualRewardsByReportEpoch(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryIndividualRewardsByReportEpochParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	potEpochRewards := k.getIndividualRewardsByReportEpoch(ctx, params)
	if len(potEpochRewards) < 1 {
		e := sdkerrors.Wrapf(types.ErrCannotFindReward, fmt.Sprintf("no Pot rewards information at epoch %s", params.Epoch.String()))
		return []byte{}, e
	}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, potEpochRewards)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func (k Keeper) getIndividualRewardsByReportEpoch(ctx sdk.Context, params types.QueryIndividualRewardsByReportEpochParams) (res []types.Reward) {
	matureEpoch := params.Epoch.Add(sdk.NewInt(k.MatureEpoch(ctx)))

	start, end := client.Paginate(params.Page, params.Limit, types.QueryDefaultLimit)
	if start < 0 || end < 0 {
		return nil
	}

	index := 0
	k.IteratorIndividualReward(ctx, matureEpoch, func(walletAddress sdk.AccAddress, individualReward types.Reward) (stop bool) {
		if !((individualReward.RewardFromMiningPool.Empty() || individualReward.RewardFromMiningPool.IsZero()) &&
			(individualReward.RewardFromTrafficPool.Empty() || individualReward.RewardFromTrafficPool.IsZero())) {
			if index >= end {
				return true
			}
			if index >= start && index < end {
				res = append(res, individualReward)
			}
			index++
		}
		return false
	})

	return res
}

// When param "Epoch" not exists: Returns wallet_address, mature_total & immature_total
// When Param "Epoch" exists: Returns walletAddress & individual_reward at that epoch
func queryRewardsByWalletAddress(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryRewardsByWalletAddrParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	// param epoch exists
	if !params.Epoch.Equal(sdk.ZeroInt()) {
		matureEpoch := params.Epoch.Add(sdk.NewInt(k.MatureEpoch(ctx)))
		rewardsByByWalletAddressAndEpoch, found := k.GetIndividualReward(ctx, params.WalletAddr, matureEpoch)
		if !found {
			e := sdkerrors.Wrapf(types.ErrCannotFindReward, fmt.Sprintf("no Pot rewards information at epoch %s", params.Epoch.String()))
			return []byte{}, e
		}

		bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, rewardsByByWalletAddressAndEpoch)
		if err != nil {
			return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil
	}

	// param epoch not exists
	immatureTotalReward := k.GetImmatureTotalReward(ctx, params.WalletAddr)
	matureTotalReward := k.GetMatureTotalReward(ctx, params.WalletAddr)
	reward := types.NewRewardInfo(params.WalletAddr, matureTotalReward, immatureTotalReward)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, reward)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func querySlashingByWalletAddress(ctx sdk.Context, req abci.RequestQuery, k Keeper, _ *codec.LegacyAmino) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(string(req.Data))
	if err != nil {
		return []byte(sdk.ZeroInt().String()), types.ErrUnknownAccountAddress
	}

	return []byte(k.registerKeeper.GetSlashing(ctx, addr).String()), nil
}

func getTotalMinedToken(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	totalMinedToken := k.GetTotalMinedTokens(ctx)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, totalMinedToken)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func getCirculationSupply(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	circulationSupply := k.GetCirculationSupply(ctx)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, circulationSupply)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func getTotalRewardByEpoch(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryTotalRewardByEpochParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	totalReward := k.GetTotalReward(ctx, params.Epoch)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, totalReward)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func getMetrics(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	metrics := k.GetMetrics(ctx)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, metrics)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
