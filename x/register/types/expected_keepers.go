package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	stratos "github.com/stratosnet/stratos-chain/types"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI // only used for simulation
	GetModuleAddress(name string) sdk.AccAddress
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) error
}

// RegisterHooks event hooks for registered node object (noalias)
type RegisterHooks interface {
	AfterNodeCreated(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool)        // Must be called when a node is created
	BeforeNodeModified(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool)      // Must be called when a node's state changes
	AfterNodeRemoved(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool)        // Must be called when a node is deleted
	AfterNodeBonded(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool)         // Must be called when a node is bonded
	AfterNodeBeginUnbonding(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool) // Must be called when a node begins unbonding
}

type DistrKeeper interface {
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
}
