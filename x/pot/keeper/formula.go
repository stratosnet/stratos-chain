package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// [S] is initial genesis deposit by all resource nodes and meta nodes at t=0
// The current unissued prepay Volume Pool [pt] is the total remaining prepay wei kept by Stratos Network but not issued to Resource Node as rewards. At time t=0,  pt=0
// total consumed Ozone is [Y]
// The remaining total Ozone limit [lt] is the upper bound of total Ozone that users can purchase from Stratos blockchain.
// the total generated traffic rewards as [R]
// R = (S + Pt) * Y / (Lt + Y)
func (k Keeper) GetTrafficReward(ctx sdk.Context, totalConsumedNoz sdk.Dec) (result sdk.Dec) {
	St := k.registerKeeper.GetEffectiveTotalDeposit(ctx).ToDec()
	if St.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("effective genesis deposit by all resource nodes and meta nodes is 0")
	}
	Pt := k.registerKeeper.GetTotalUnissuedPrepay(ctx).Amount.ToDec()
	if Pt.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("total remaining prepay not issued is 0")
	}
	Y := totalConsumedNoz
	if Y.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("total consumed noz is 0")
	}
	Lt := k.registerKeeper.GetRemainingOzoneLimit(ctx).ToDec()
	if Lt.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("remaining total noz limit is 0")
	}
	R := St.Add(Pt).Mul(Y).Quo(Lt.Add(Y))
	if R.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("traffic reward to distribute is 0")
	}
	return R
}

func (k Keeper) GetPrepayAmount(Lt, tokenAmount, St, Pt sdk.Int) (purchaseNoz, remainingNoz sdk.Int, err error) {
	purchaseNoz = Lt.ToDec().
		Mul(tokenAmount.ToDec()).
		Quo((St.
			Add(Pt).
			Add(tokenAmount)).ToDec()).
		TruncateInt()

	if purchaseNoz.GT(Lt) {
		return sdk.ZeroInt(), sdk.ZeroInt(), fmt.Errorf("not enough remaining ozone limit to complete prepay")
	}
	remainingNoz = Lt.Sub(purchaseNoz)

	return purchaseNoz, remainingNoz, nil
}

func (k Keeper) GetCurrentNozPrice(St, Pt, Lt sdk.Int) (currentNozPrice sdk.Dec) {
	currentNozPrice = (St.Add(Pt)).ToDec().
		Quo(Lt.ToDec())
	return
}

// NozSupply calc remaining/total supply for noz
func (k Keeper) NozSupply(ctx sdk.Context) (remaining, total sdk.Int) {
	remaining = k.registerKeeper.GetRemainingOzoneLimit(ctx) // Lt
	depositNozRate := k.registerKeeper.GetDepositNozRate(ctx)
	St := k.registerKeeper.GetEffectiveTotalDeposit(ctx)
	total = St.ToDec().Quo(depositNozRate).TruncateInt()
	return remaining, total
}
