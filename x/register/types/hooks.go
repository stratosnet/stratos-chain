package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// combine multiple register hooks, all hook functions are run in array sequence
type MultiRegisterHooks []RegisterHooks

func NewMultiRegisterHooks(hooks ...RegisterHooks) MultiRegisterHooks {
	return hooks
}

// nolint
func (h MultiRegisterHooks) AfterNodeCreated(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool) {
	for i := range h {
		h[i].AfterNodeCreated(ctx, networkAddr, isIndexingNode)
	}
}
func (h MultiRegisterHooks) BeforeNodeModified(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool) {
	for i := range h {
		h[i].BeforeNodeModified(ctx, networkAddr, isIndexingNode)
	}
}
func (h MultiRegisterHooks) AfterNodeRemoved(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool) {
	for i := range h {
		h[i].AfterNodeRemoved(ctx, networkAddr, isIndexingNode)
	}
}
func (h MultiRegisterHooks) AfterNodeBonded(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool) {
	for i := range h {
		h[i].AfterNodeBonded(ctx, networkAddr, isIndexingNode)
	}
}
func (h MultiRegisterHooks) AfterNodeBeginUnbonding(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool) {
	for i := range h {
		h[i].AfterNodeBeginUnbonding(ctx, networkAddr, isIndexingNode)
	}
}
