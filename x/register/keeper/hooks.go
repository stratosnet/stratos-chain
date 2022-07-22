package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// Implements RegisterHooks interface
var _ types.RegisterHooks = Keeper{}

// AfterNodeCreated - call hook if registered
func (k Keeper) AfterNodeCreated(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool) {
	if k.hooks != nil {
		k.hooks.AfterNodeCreated(ctx, networkAddr, isMetaNode)
	}
}

// BeforeNodeModified - call hook if registered
func (k Keeper) BeforeNodeModified(ctx sdk.Context, network stratos.SdsAddress, isMetaNode bool) {
	if k.hooks != nil {
		k.hooks.BeforeNodeModified(ctx, network, isMetaNode)
	}
}

// AfterNodeRemoved - call hook if registered
func (k Keeper) AfterNodeRemoved(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool) {
	if k.hooks != nil {
		k.hooks.AfterNodeRemoved(ctx, networkAddr, isMetaNode)
	}
}

// AfterNodeBonded - call hook if registered
func (k Keeper) AfterNodeBonded(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool) {
	if k.hooks != nil {
		k.hooks.AfterNodeBonded(ctx, networkAddr, isMetaNode)
	}
}

// AfterNodeBeginUnbonding - call hook if registered
func (k Keeper) AfterNodeBeginUnbonding(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool) {
	if k.hooks != nil {
		k.hooks.AfterNodeBeginUnbonding(ctx, networkAddr, isMetaNode)
	}
}
