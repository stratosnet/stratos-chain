package keeper

import (
	"encoding/hex"
	// this line is used by starport scaffolding # 1
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// query file hash
	QueryFileHash = "uploaded_file"
	QueryPrepay   = "prepay"
)

// NewQuerier creates a new querier for sds clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryFileHash:
			return queryFileHash(ctx, req, k)
		case QueryPrepay:
			return queryPrepay(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown sds query endpoint "+req.String()+hex.EncodeToString(req.Data))
		}
	}
}

// queryFileHash fetch an file's hash for the supplied height.
func queryFileHash(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	fileHash, err := k.GetFileHash(ctx, req.Data)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return fileHash, nil
}

// queryFileHash fetch an file's hash for the supplied height.
func queryPrepay(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	balance, err := k.GetPrepayBytes(ctx, req.Data)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return balance, nil
}
