package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

/*
When a module wishes to interact with another module, it is good practice to define what it will use
as an interface so the module cannot use things that are not permitted.
TODO: Create interfaces of what you expect the other keepers to have to be able to use this module.
type BankKeeper interface {
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
*/

// ParamSubspace defines the expected Subspace interface
type ParamSubspace interface {
	WithKeyTable(table paramstypes.KeyTable) paramstypes.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps paramstypes.ParamSet)
	SetParamSet(ctx sdk.Context, ps paramstypes.ParamSet)
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	IterateAccounts(ctx sdk.Context, process func(authtypes.AccountI) (stop bool))
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI // only used for simulation

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI

	// SetModuleAccount TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, authtypes.ModuleAccountI)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	GetSupply(ctx sdk.Context, denom string) sdk.Coin

	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
}

// RegisterHooks event hooks for registered node object (noalias)
type RegisterHooks interface {
	AfterNodeCreated(ctx sdk.Context, networkAddr stratos.SdsAddress, isIndexingNode bool)   // Must be called when a node is created
	BeforeNodeModified(ctx sdk.Context, networkAddr stratos.SdsAddress, isIndexingNode bool) // Must be called when a node's state changes
	AfterNodeRemoved(ctx sdk.Context, networkAddr stratos.SdsAddress, isIndexingNode bool)   // Must be called when a node is deleted

	AfterNodeBonded(ctx sdk.Context, networkAddr stratos.SdsAddress, isIndexingNode bool)         // Must be called when a node is bonded
	AfterNodeBeginUnbonding(ctx sdk.Context, networkAddr stratos.SdsAddress, isIndexingNode bool) // Must be called when a node begins unbonding

	//BeforeNodeCreated(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool)  // Must be called when a node is created
	//BeforeNodeModified(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool) // Must be called when a node's shares are modified
	//BeforeNodeRemoved(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool)  // Must be called when a node is removed
	//AfterNodeModified(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool)
}
