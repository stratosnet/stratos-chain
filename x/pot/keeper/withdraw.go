package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) Withdraw(ctx sdk.Context, amount sdk.Coins, walletAddress sdk.AccAddress, targetAddress sdk.AccAddress) error {
	matureReward := k.GetMatureTotalReward(ctx, walletAddress)
	withdrawCoins := matureReward.Sub(amount)
	_, err := k.BankKeeper.AddCoins(ctx, targetAddress, withdrawCoins)
	if err != nil {
		return err
	}
	newMatureReward := matureReward.Sub(withdrawCoins)
	k.setMatureTotalReward(ctx, walletAddress, newMatureReward)
	return nil
}
