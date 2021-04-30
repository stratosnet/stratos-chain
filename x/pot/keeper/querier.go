package keeper

import (
	// this line is used by starport scaffolding # 1
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const QueryVolumeReport = "volume_report"

// NewQuerier creates a new querier for pot clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryVolumeReport:
			return queryVolumeReport(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown pot query endpoint")
		}
	}
}

// queryVolumeReport fetches an hash of report volume for the supplied height.
func queryVolumeReport(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	volumeReportHash, err := k.GetVolumeReport(ctx, req.Data)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return volumeReportHash, nil
}
