package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// Implements RegisterHooks interface
var _ types.RegisterHooks = Keeper{}

// AfterNodeCreated - call hook if registered
func (k Keeper) AfterNodeCreated(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool) {
	if k.hooks != nil {
		k.hooks.AfterNodeCreated(ctx, networkAddr, isIndexingNode)
	}
}

// BeforeNodeModified - call hook if registered
func (k Keeper) BeforeNodeModified(ctx sdk.Context, network sdk.AccAddress, isIndexingNode bool) {
	if k.hooks != nil {
		k.hooks.BeforeNodeModified(ctx, network, isIndexingNode)
	}
}

// AfterNodeRemoved - call hook if registered
func (k Keeper) AfterNodeRemoved(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool) {
	if k.hooks != nil {
		k.hooks.AfterNodeRemoved(ctx, networkAddr, isIndexingNode)
	}
}

// AfterNodeBonded - call hook if registered
func (k Keeper) AfterNodeBonded(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool) {
	if k.hooks != nil {
		k.hooks.AfterNodeBonded(ctx, networkAddr, isIndexingNode)
	}
}

// AfterNodeBeginUnbonding - call hook if registered
func (k Keeper) AfterNodeBeginUnbonding(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool) {
	if k.hooks != nil {
		k.hooks.AfterNodeBeginUnbonding(ctx, networkAddr, isIndexingNode)
	}
}
