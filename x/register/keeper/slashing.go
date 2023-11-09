package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DeductSlashing deduct slashing amount from coins, return the coins that after deduction
func (k Keeper) DeductSlashing(ctx sdk.Context, walletAddress sdk.AccAddress, coins sdk.Coins, slashingDenom string) (remaining, deducted sdk.Coins) {
	slashing := k.GetSlashing(ctx, walletAddress)
	remaining = coins
	deducted = sdk.Coins{}

	if slashing.LTE(sdkmath.ZeroInt()) || coins.Empty() || coins.AmountOf(slashingDenom).IsZero() {
		return
	}

	coinAmt := coins.AmountOf(slashingDenom)
	if coinAmt.GTE(slashing) {
		deducted = sdk.NewCoins(sdk.NewCoin(slashingDenom, slashing))
		coinAmt = coinAmt.Sub(slashing)
		remaining = remaining.Sub(deducted...)
		slashing = sdkmath.ZeroInt()
	} else {
		deducted = sdk.NewCoins(sdk.NewCoin(slashingDenom, coinAmt))
		slashing = slashing.Sub(coinAmt)
		remaining = remaining.Sub(deducted...)
		coinAmt = sdkmath.ZeroInt()
	}

	k.SetSlashing(ctx, walletAddress, slashing)
	return
}
