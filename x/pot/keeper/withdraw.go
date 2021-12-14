package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) Withdraw(ctx sdk.Context, amount sdk.Coins, walletAddress sdk.AccAddress, targetAddress sdk.AccAddress) error {
	matureReward := k.GetMatureTotalReward(ctx, walletAddress)
	matureRewardBalance := matureReward.Sub(amount)
	_, err := k.BankKeeper.AddCoins(ctx, targetAddress, amount)
	if err != nil {
		return err
	}
	k.setMatureTotalReward(ctx, walletAddress, matureRewardBalance)
	return nil
}
