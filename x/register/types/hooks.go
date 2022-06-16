package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

// MultiRegisterHooks combines multiple register hooks, all hook functions are run in array sequence
type MultiRegisterHooks []RegisterHooks

func NewMultiRegisterHooks(hooks ...RegisterHooks) MultiRegisterHooks {
	return hooks
}

// nolint
func (h MultiRegisterHooks) AfterNodeCreated(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool) {
	for i := range h {
		h[i].AfterNodeCreated(ctx, networkAddr, isMetaNode)
	}
}
func (h MultiRegisterHooks) BeforeNodeModified(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool) {
	for i := range h {
		h[i].BeforeNodeModified(ctx, networkAddr, isMetaNode)
	}
}
func (h MultiRegisterHooks) AfterNodeRemoved(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool) {
	for i := range h {
		h[i].AfterNodeRemoved(ctx, networkAddr, isMetaNode)
	}
}
func (h MultiRegisterHooks) AfterNodeBonded(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool) {
	for i := range h {
		h[i].AfterNodeBonded(ctx, networkAddr, isMetaNode)
	}
}
func (h MultiRegisterHooks) AfterNodeBeginUnbonding(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool) {
	for i := range h {
		h[i].AfterNodeBeginUnbonding(ctx, networkAddr, isMetaNode)
	}
}
