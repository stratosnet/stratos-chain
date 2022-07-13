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
	QueryPotSlashingByWalletAddr = "query_pot_slashing_by_wallet_address"
	QueryPotParams               = "query_pot_params"
	QueryDefaultLimit            = 100
)

// NewQuerier creates a new querier for pot clients.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryVolumeReport:
			return queryVolumeReport(ctx, req, k, legacyQuerierCdc)
		case QueryPotRewardsByReportEpoch:
			return queryPotRewardsByReportEpoch(ctx, req, k, legacyQuerierCdc)
		case QueryPotRewardsByWalletAddr:
			return queryPotRewardsByWalletAddress(ctx, req, k, legacyQuerierCdc)
		case QueryPotSlashingByWalletAddr:
			return queryPotSlashingByWalletAddress(ctx, req, k, legacyQuerierCdc)
		case QueryPotParams:
			return getPotParams(ctx, req, k, legacyQuerierCdc)
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
				epoch, k.GetLastReportedEpoch(ctx).String()))
		return []byte{}, e
	}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, reportRecord)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

// queryPotRewardsByReportEpoch fetches total rewards and owner individual rewards from traffic and mining.
func queryPotRewardsByReportEpoch(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryPotRewardsByReportEpochParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	potEpochRewards := k.getPotRewardsByReportEpoch(ctx, params)
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

func (k Keeper) getPotRewardsByReportEpoch(ctx sdk.Context, params types.QueryPotRewardsByReportEpochParams) (res []types.Reward) {
	matureEpoch := params.Epoch.Add(sdk.NewInt(k.MatureEpoch(ctx)))

	if !params.WalletAddress.Empty() {
		reward, found := k.GetIndividualReward(ctx, params.WalletAddress, matureEpoch)
		if found {
			res = append(res, reward)
		}
	} else {
		k.IteratorIndividualReward(ctx, matureEpoch, func(walletAddress sdk.AccAddress, individualReward types.Reward) (stop bool) {
			if !((individualReward.RewardFromMiningPool.Empty() || individualReward.RewardFromMiningPool.IsZero()) &&
				(individualReward.RewardFromTrafficPool.Empty() || individualReward.RewardFromTrafficPool.IsZero())) {
				res = append(res, individualReward)
			}
			return false
		})
	}

	start, end := client.Paginate(len(res), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		return nil
	} else {
		res = res[start:end]
		return res
	}
}

func queryPotRewardsByWalletAddress(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryPotRewardsByWalletAddrParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	immatureTotalReward := k.GetImmatureTotalReward(ctx, params.WalletAddr)
	matureTotalReward := k.GetMatureTotalReward(ctx, params.WalletAddr)
	reward := types.NewPotRewardInfo(params.WalletAddr, matureTotalReward, immatureTotalReward)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, reward)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryPotSlashingByWalletAddress(ctx sdk.Context, req abci.RequestQuery, k Keeper, _ *codec.LegacyAmino) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(string(req.Data))
	if err != nil {
		return []byte(sdk.ZeroInt().String()), types.ErrUnknownAccountAddress
	}

	return []byte(k.RegisterKeeper.GetSlashing(ctx, addr).String()), nil
}
