package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTrafficReward
// [S] is initial genesis deposit by all resource nodes and meta nodes at t=0
// The current unissued prepay Volume Pool [pt] is the total remaining prepay wei kept by Stratos Network but not issued to Resource Node as rewards. At time t=0,  pt=0
// total consumed Ozone is [Y]
// The remaining total Ozone limit [lt] is the upper bound of total Ozone that users can purchase from Stratos blockchain.
// the total generated traffic rewards as [R]
// R = (S + Pt) * Y / (Lt + Y)
func (k Keeper) GetTrafficReward(ctx sdk.Context, totalConsumedNoz sdkmath.LegacyDec) (result sdkmath.LegacyDec) {
	St := k.registerKeeper.GetEffectiveTotalDeposit(ctx).ToLegacyDec()
	if St.Equal(sdkmath.LegacyZeroDec()) {
		ctx.Logger().Info("effective genesis deposit by all resource nodes and meta nodes is 0")
	}
	Pt := k.registerKeeper.GetTotalUnissuedPrepay(ctx).Amount.ToLegacyDec()
	if Pt.Equal(sdkmath.LegacyZeroDec()) {
		ctx.Logger().Info("total remaining prepay not issued is 0")
	}
	Y := totalConsumedNoz
	if Y.Equal(sdkmath.LegacyZeroDec()) {
		ctx.Logger().Info("total consumed noz is 0")
	}
	Lt := k.registerKeeper.GetRemainingOzoneLimit(ctx).ToLegacyDec()
	if Lt.Equal(sdkmath.LegacyZeroDec()) {
		ctx.Logger().Info("remaining total noz limit is 0")
	}
	R := St.Add(Pt).Mul(Y).Quo(Lt.Add(Y))
	if R.Equal(sdkmath.LegacyZeroDec()) {
		ctx.Logger().Info("traffic reward to distribute is 0")
	}
	return R
}

func (k Keeper) GetPrepayAmount(Lt, tokenAmount, St, Pt sdkmath.Int) (purchaseNoz, remainingNoz sdkmath.Int, err error) {
	purchaseNoz = Lt.ToLegacyDec().
		Mul(tokenAmount.ToLegacyDec()).
		Quo((St.
			Add(Pt).
			Add(tokenAmount)).ToLegacyDec()).
		TruncateInt()

	if purchaseNoz.GT(Lt) {
		return sdkmath.ZeroInt(), sdkmath.ZeroInt(), fmt.Errorf("not enough remaining ozone limit to complete prepay")
	}
	remainingNoz = Lt.Sub(purchaseNoz)

	return purchaseNoz, remainingNoz, nil
}

func (k Keeper) GetCurrentNozPrice(St, Pt, Lt sdkmath.Int) (currentNozPrice sdkmath.LegacyDec) {
	currentNozPrice = (St.Add(Pt)).ToLegacyDec().
		Quo(Lt.ToLegacyDec())
	return
}

// NozSupply calc remaining/total supply for noz
func (k Keeper) NozSupply(ctx sdk.Context) (remaining, total sdkmath.Int) {
	remaining = k.registerKeeper.GetRemainingOzoneLimit(ctx) // Lt
	depositNozRate := k.registerKeeper.GetDepositNozRate(ctx)
	St := k.registerKeeper.GetEffectiveTotalDeposit(ctx)
	total = St.ToLegacyDec().Quo(depositNozRate).TruncateInt()
	return remaining, total
}
