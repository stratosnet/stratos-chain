package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	// this line is used by starport scaffolding # 1
	abci "github.com/tendermint/tendermint/abci/types"
)

const QueryVolumeReportHash = "volume_report"
const QueryNodeVolume = "node_volume"

// NewQuerier creates a new querier for pot clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		// this line is used by starport scaffolding # 2
		case QueryVolumeReportHash:
			return queryVolumeReportHash(ctx, req, k)
		case QueryNodeVolume:
			return querySingleNodeVolume(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown pot query endpoint")
		}
	}
}

// queryVolumeReportHash fetches an hash of report volume for the supplied height.
func queryVolumeReportHash(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	volumeReportHash, err := k.GetVolumeReportHash(ctx, req.Data)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return volumeReportHash, nil
}

// queryVolumeReportHash fetches an hash of report volume for the supplied height.
func querySingleNodeVolume(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	singleVolumeReport, err := k.GetSingleNodeVolume(ctx, req.Data)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return singleVolumeReport, nil
}
