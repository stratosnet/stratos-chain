package keeper

import (
	"encoding/hex"
	// this line is used by starport scaffolding # 1
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	QueryFileHash       = "uploaded_file"
	QueryPrepay         = "prepay"
	QuerySimulatePrepay = "simulate_prepay"
)

// NewQuerier creates a new querier for sds clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryFileHash:
			return queryFileHash(ctx, req, k)
		case QueryPrepay:
			return queryPrepay(ctx, req, k)
		case QuerySimulatePrepay:
			return querySimulatePrepay(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown sds query endpoint "+req.String()+hex.EncodeToString(req.Data))
		}
	}
}

// queryFileHash fetch an file's hash for the supplied height.
func queryFileHash(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	fileHash, err := k.GetFileInfoBytesByFileHash(ctx, req.Data)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return fileHash, nil
}

// queryPrepay fetch prepaid balance of an account.
func queryPrepay(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	balance, err := k.GetPrepayBytes(ctx, req.Data)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return balance, nil
}

// querySimulatePrepay fetch amt of uoz with a simulated prepay of X ustos.
func querySimulatePrepay(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var amtToPrepay sdk.Int
	err := amtToPrepay.UnmarshalJSON(req.Data)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	uozAmt := k.simulatePurchaseUoz(ctx, amtToPrepay)
	uozAmtByte, _ := uozAmt.MarshalJSON()
	return uozAmtByte, nil
}
