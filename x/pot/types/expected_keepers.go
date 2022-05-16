package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
	//authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	//"github.com/cosmos/cosmos-sdk/x/params"
)

// ParamSubSpace defines the expected Subspace interface
type ParamSubSpace interface {
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

type RegisterKeeper interface {
	GetIndexingNode(ctx sdk.Context, p2pAddress stratos.SdsAddress) (indexingNode types.IndexingNode, found bool)
	SetIndexingNode(ctx sdk.Context, indexingNode types.IndexingNode)

	GetResourceNode(ctx sdk.Context, p2pAddress stratos.SdsAddress) (resourceNode types.ResourceNode, found bool)
	SetResourceNode(ctx sdk.Context, resourceNode types.ResourceNode)

	GetSlashing(ctx sdk.Context, walletAddress sdk.AccAddress) (res sdk.Int)
	SetSlashing(ctx sdk.Context, walletAddress sdk.AccAddress, slashing sdk.Int)
	DeductSlashing(ctx sdk.Context, walletAddress sdk.AccAddress, coins sdk.Coins) sdk.Coins

	GetRemainingOzoneLimit(ctx sdk.Context) (value sdk.Int)
	SetRemainingOzoneLimit(ctx sdk.Context, value sdk.Int)
	GetTotalUnissuedPrepay(ctx sdk.Context) (totalUnissuedPrepay sdk.Coin)
	SetTotalUnissuedPrepay(ctx sdk.Context, totalUnissuedPrepay sdk.Coin)

	GetResourceNodeBondedToken(ctx sdk.Context) (token sdk.Coin)
	SetResourceNodeBondedToken(ctx sdk.Context, token sdk.Coin)
	GetIndexingNodeBondedToken(ctx sdk.Context) (token sdk.Coin)
	SetIndexingNodeBondedToken(ctx sdk.Context, token sdk.Coin)

	GetInitialGenesisStakeTotal(ctx sdk.Context) (stake sdk.Int)
	SetInitialGenesisStakeTotal(ctx sdk.Context, stake sdk.Int)

	GetAllResourceNodes(ctx sdk.Context) (resourceNodes *types.ResourceNodes)
	GetAllIndexingNodes(ctx sdk.Context) (indexingNodes *types.IndexingNodes)
}

type StakingKeeper interface {
	TotalBondedTokens(ctx sdk.Context) sdk.Int
}
