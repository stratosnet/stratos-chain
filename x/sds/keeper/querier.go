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
		case types.QueryCurrNozPrice:
			return queryCurrNozPrice(ctx, req, k, legacyQuerierCdc)
		case types.QueryNozSupply:
			return queryNozSupply(ctx, req, k, legacyQuerierCdc)
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

// querySimulatePrepay fetch amt of noz with a simulated prepay of X wei.
func querySimulatePrepay(ctx sdk.Context, req abci.RequestQuery, k Keeper, _ *codec.LegacyAmino) ([]byte, error) {
	var amtToPrepay sdk.Int
	err := amtToPrepay.UnmarshalJSON(req.Data)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	// temporary solution, avoid to modify Rest api. After upgrade to cosmos sdk v0.46.x, legacy Rest API will be removed
	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), amtToPrepay))
	nozAmt := k.simulatePurchaseNoz(ctx, coins)
	nozAmtByte, _ := nozAmt.MarshalJSON()
	return nozAmtByte, nil
}

// queryCurrNozPrice fetch current noz price.
func queryCurrNozPrice(ctx sdk.Context, _ abci.RequestQuery, k Keeper, _ *codec.LegacyAmino) ([]byte, error) {
	nozPrice := k.registerKeeper.CurrNozPrice(ctx)
	nozPriceByte, _ := nozPrice.MarshalJSON()
	return nozPriceByte, nil
}

// queryNozSupply fetch remaining/total noz supply.
func queryNozSupply(ctx sdk.Context, _ abci.RequestQuery, k Keeper, _ *codec.LegacyAmino) ([]byte, error) {
	type Supply struct {
		Remaining sdk.Int
		Total     sdk.Int
	}
	var nozSupply Supply
	nozSupply.Remaining, nozSupply.Total = k.registerKeeper.NozSupply(ctx)
	nozSupplyByte, _ := json.Marshal(nozSupply)
	return nozSupplyByte, nil
}
