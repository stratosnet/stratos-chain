package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// SetParams sets the params on the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)

	return nil
}

// GetParams returns the params from the store
func (k Keeper) GetParams(ctx sdk.Context) (p types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return p
	}
	k.cdc.MustUnmarshal(bz, &p)
	return p
}

// BondDenom - Bondable coin denomination
func (k Keeper) BondDenom(ctx sdk.Context) (res string) {
	params := k.GetParams(ctx)
	res = params.GetBondDenom()
	return
}

func (k Keeper) RewardDenom(ctx sdk.Context) (res string) {
	params := k.GetParams(ctx)
	res = params.GetRewardDenom()
	return
}

func (k Keeper) MatureEpoch(ctx sdk.Context) (res int64) {
	params := k.GetParams(ctx)
	res = params.GetMatureEpoch()
	return
}

func (k Keeper) MiningRewardParams(ctx sdk.Context) (res []types.MiningRewardParam) {
	params := k.GetParams(ctx)
	res = params.GetMiningRewardParams()
	return
}

func (k Keeper) GetMiningRewardParamByMinedToken(ctx sdk.Context, minedToken sdk.Coin) (types.MiningRewardParam, error) {
	miningRewardParams := k.MiningRewardParams(ctx)
	for _, param := range miningRewardParams {
		if minedToken.IsGTE(param.TotalMinedValveStart) && minedToken.IsLT(param.TotalMinedValveEnd) {
			return param, nil
		}
	}
	return miningRewardParams[len(miningRewardParams)-1], types.ErrOutOfIssuance
}

func (k Keeper) GetTotalMining(ctx sdk.Context) sdk.Coin {
	miningRewardParams := k.MiningRewardParams(ctx)
	return miningRewardParams[len(miningRewardParams)-1].TotalMinedValveEnd
}

func (k Keeper) GetCommunityTax(ctx sdk.Context) (res sdkmath.LegacyDec) {
	params := k.GetParams(ctx)
	res = params.CommunityTax
	return
}

func (k Keeper) InitialTotalSupply(ctx sdk.Context) (res sdk.Coin) {
	params := k.GetParams(ctx)
	res = params.GetInitialTotalSupply()
	return
}
