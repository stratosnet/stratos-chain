package keeper

import (
	"encoding/hex"
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stratosnet/stratos-chain/x/sds/types"

	// this line is used by starport scaffolding # 1
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a new querier for sds clients.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryUploadedFile:
			return queryUploadedFileByHash(ctx, req, k, legacyQuerierCdc)
		case types.QuerySimulatePrepay:
			return querySimulatePrepay(ctx, req, k, legacyQuerierCdc)
		case types.QueryCurrUozPrice:
			return queryCurrUozPrice(ctx, req, k, legacyQuerierCdc)
		case types.QueryUozSupply:
			return queryUozSupply(ctx, req, k, legacyQuerierCdc)
		case types.QueryParams:
			return getSdsParams(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown sds query endpoint "+req.String()+hex.EncodeToString(req.Data))
		}
	}
}

func getSdsParams(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryFileHash fetch a file's hash for the supplied height.
func queryUploadedFileByHash(ctx sdk.Context, req abci.RequestQuery, k Keeper, _ *codec.LegacyAmino) ([]byte, error) {
	fileInfo, err := k.GetFileInfoBytesByFileHash(ctx, req.Data)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return fileInfo, nil
}

// querySimulatePrepay fetch amt of uoz with a simulated prepay of X ustos.
func querySimulatePrepay(ctx sdk.Context, req abci.RequestQuery, k Keeper, _ *codec.LegacyAmino) ([]byte, error) {
	var amtToPrepay sdk.Int
	err := amtToPrepay.UnmarshalJSON(req.Data)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	uozAmt := k.simulatePurchaseUoz(ctx, amtToPrepay)
	uozAmtByte, _ := uozAmt.MarshalJSON()
	return uozAmtByte, nil
}

// queryCurrUozPrice fetch current uoz price.
func queryCurrUozPrice(ctx sdk.Context, _ abci.RequestQuery, k Keeper, _ *codec.LegacyAmino) ([]byte, error) {
	uozPrice := k.RegisterKeeper.CurrUozPrice(ctx)
	uozPriceByte, _ := uozPrice.MarshalJSON()
	return uozPriceByte, nil
}

// queryUozSupply fetch remaining/total uoz supply.
func queryUozSupply(ctx sdk.Context, _ abci.RequestQuery, k Keeper, _ *codec.LegacyAmino) ([]byte, error) {
	type Supply struct {
		Remaining sdk.Int
		Total     sdk.Int
	}
	var uozSupply Supply
	uozSupply.Remaining, uozSupply.Total = k.RegisterKeeper.UozSupply(ctx)
	uozSupplyByte, _ := json.Marshal(uozSupply)
	return uozSupplyByte, nil
}
