package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI // only used for simulation
	GetModuleAddress(name string) sdk.AccAddress
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	MintCoins(ctx sdk.Context, moduleName string, amounts sdk.Coins) error
}

type RegisterKeeper interface {
	GetMetaNode(ctx sdk.Context, p2pAddress stratos.SdsAddress) (metaNode types.MetaNode, found bool)
	GetResourceNode(ctx sdk.Context, p2pAddress stratos.SdsAddress) (resourceNode types.ResourceNode, found bool)
	SetResourceNode(ctx sdk.Context, resourceNode types.ResourceNode)

	GetSlashing(ctx sdk.Context, walletAddress sdk.AccAddress) (res sdk.Int)
	SetSlashing(ctx sdk.Context, walletAddress sdk.AccAddress, slashing sdk.Int)
	DeductSlashing(ctx sdk.Context, walletAddress sdk.AccAddress, coins sdk.Coins, slashingDenom string) (remaining, deducted sdk.Coins)

	GetRemainingOzoneLimit(ctx sdk.Context) (value sdk.Int)
	SetRemainingOzoneLimit(ctx sdk.Context, value sdk.Int)
	GetTotalUnissuedPrepay(ctx sdk.Context) (totalUnissuedPrepay sdk.Coin)

	GetResourceNodeBondedToken(ctx sdk.Context) (token sdk.Coin)
	GetMetaNodeBondedToken(ctx sdk.Context) (token sdk.Coin)

	GetEffectiveTotalStake(ctx sdk.Context) (stake sdk.Int)

	GetResourceNodeIterator(ctx sdk.Context) sdk.Iterator
	GetMetaNodeIterator(ctx sdk.Context) sdk.Iterator
	GetBondedMetaNodeCnt(ctx sdk.Context) sdk.Int

	DecreaseOzoneLimitBySubtractStake(ctx sdk.Context, stake sdk.Int) (ozoneLimitChange sdk.Int)
	IncreaseOzoneLimitByAddStake(ctx sdk.Context, stake sdk.Int) (ozoneLimitChange sdk.Int)
	GetUnbondingNodeBalance(ctx sdk.Context, networkAddr stratos.SdsAddress) sdk.Int

	NozSupply(ctx sdk.Context) (remaining, total sdk.Int)
	OwnMetaNode(ctx sdk.Context, ownerAddr sdk.AccAddress, p2pAddr stratos.SdsAddress) (found bool)
}

type StakingKeeper interface {
	TotalBondedTokens(ctx sdk.Context) sdk.Int
}

type DistrKeeper interface {
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
	GetFeePool(ctx sdk.Context) (feePool distrtypes.FeePool)
}
