package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

const (
	DefaultParamspace = types.ModuleName
)

// ParamKeyTable ParamTable for register module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// BondDenom - Bondable coin denomination
func (k Keeper) BondDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyBondDenom, &res)
	return
}

// GetParams returns the total set of register parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the register parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
