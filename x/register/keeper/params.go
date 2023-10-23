package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// SetParams sets the params on the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
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

// MaxEntries - Maximum number of simultaneous unbonding
// delegations or redelegations (per pair/trio)
func (k Keeper) MaxEntries(ctx sdk.Context) (res uint32) {
	params := k.GetParams(ctx)
	res = params.GetMaxEntries()
	return
}

func (k Keeper) UnbondingThreasholdTime(ctx sdk.Context) (res time.Duration) {
	params := k.GetParams(ctx)
	res = params.GetUnbondingThreasholdTime()
	return
}

func (k Keeper) UnbondingCompletionTime(ctx sdk.Context) (res time.Duration) {
	params := k.GetParams(ctx)
	res = params.GetUnbondingCompletionTime()
	return
}

func (k Keeper) ResourceNodeRegEnabled(ctx sdk.Context) (res bool) {
	params := k.GetParams(ctx)
	res = params.GetResourceNodeRegEnabled()
	return
}

func (k Keeper) ResourceNodeMinDeposit(ctx sdk.Context) (res sdk.Coin) {
	params := k.GetParams(ctx)
	res = params.GetResourceNodeMinDeposit()
	return
}

func (k Keeper) VotingPeriod(ctx sdk.Context) (res time.Duration) {
	params := k.GetParams(ctx)
	res = params.GetVotingPeriod()
	return
}
