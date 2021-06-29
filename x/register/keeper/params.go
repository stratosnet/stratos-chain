package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// GetParams returns the total set of register parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the register parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// BondDenom - Bondable coin denomination
func (k Keeper) BondDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyBondDenom, &res)
	return
}
