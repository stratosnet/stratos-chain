package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// GetParams returns the total set of pot parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the pot parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// BondDenom - Bondable coin denomination
func (k Keeper) BondDenom(ctx sdk.Context) (res string) {
	k.paramSpace.Get(ctx, types.KeyBondDenom, &res)
	return
}

func (k Keeper) MiningRewardParams(ctx sdk.Context) (res []types.MiningRewardParam) {
	k.paramSpace.Get(ctx, types.KeyMiningRewardParams, &res)
	return
}

func (k Keeper) GetMiningRewardParamByMinedToken(ctx sdk.Context, minedToken sdk.Int) (types.MiningRewardParam, error) {
	miningRewardParams := k.MiningRewardParams(ctx)
	for _, param := range miningRewardParams {
		if minedToken.GTE(param.TotalMinedValveStart) && minedToken.LT(param.TotalMinedValveEnd) {
			return param, nil
		}
	}
	return miningRewardParams[0], types.ErrOutOfIssuance
}
