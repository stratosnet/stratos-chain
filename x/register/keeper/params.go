package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// GetParams returns the total set of register parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the register parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// BondDenom - Bondable coin denomination
func (k Keeper) BondDenom(ctx sdk.Context) (res string) {
	k.paramSpace.Get(ctx, types.KeyBondDenom, &res)
	return
}

// MaxEntries - Maximum number of simultaneous unbonding
// delegations or redelegations (per pair/trio)
func (k Keeper) MaxEntries(ctx sdk.Context) (res uint32) {
	k.paramSpace.Get(ctx, types.KeyMaxEntries, &res)
	return
}

// UnbondingThreasholdTime
func (k Keeper) UnbondingThreasholdTime(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.KeyUnbondingThreasholdTime, &res)
	return
}

// UnbondingCompletionTime
func (k Keeper) UnbondingCompletionTime(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.KeyUnbondingCompletionTime, &res)
	return
}

// ResourceNodeRegEnabled
func (k Keeper) ResourceNodeRegEnabled(ctx sdk.Context) (res bool) {
	k.paramSpace.Get(ctx, types.KeyResourceNodeRegEnabled, &res)
	return
}

func (k Keeper) ResourceNodeMinDeposit(ctx sdk.Context) (res sdk.Coin) {
	k.paramSpace.Get(ctx, types.KeyResourceNodeMinDeposit, &res)
	return
}

// VotingPeriod
func (k Keeper) VotingPeriod(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.KeyVotingPeriod, &res)
	return
}
