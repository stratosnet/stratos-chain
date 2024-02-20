package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

type RegisterKeeper interface {
	GetCurrNozPriceParams(ctx sdk.Context) (St, Pt, Lt sdkmath.Int)
	GetEffectiveTotalDeposit(ctx sdk.Context) (deposit sdkmath.Int)
	GetMetaNodeBitMapIndex(ctx sdk.Context, networkAddr stratos.SdsAddress) (index int, err error)
	GetRemainingOzoneLimit(ctx sdk.Context) (value sdkmath.Int)
	GetTotalUnissuedPrepay(ctx sdk.Context) (totalUnissuedPrepay sdk.Coin)
	SetRemainingOzoneLimit(ctx sdk.Context, value sdkmath.Int)
	OwnMetaNode(ctx sdk.Context, ownerAddr sdk.AccAddress, p2pAddr stratos.SdsAddress) bool
	CalculatePurchaseAmount(ctx sdk.Context, amount sdkmath.Int) (sdkmath.Int, sdkmath.Int, error)
}

type PotKeeper interface {
	GetCurrentNozPrice(St, Pt, Lt sdkmath.Int) (currentNozPrice sdkmath.LegacyDec)
	NozSupply(ctx sdk.Context) (remaining, total sdkmath.Int)
}
