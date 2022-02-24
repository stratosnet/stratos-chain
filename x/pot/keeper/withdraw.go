package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func (k Keeper) Withdraw(ctx sdk.Context, amount sdk.Coins, walletAddress sdk.AccAddress, targetAddress sdk.AccAddress) error {
	matureReward := k.GetMatureTotalReward(ctx, walletAddress)
	if !matureReward.IsAllGTE(amount) {
		return types.ErrInsufficientMatureTotal
	}
	_, err := k.BankKeeper.AddCoins(ctx, targetAddress, amount)
	if err != nil {
		return err
	}
	matureRewardBalance := matureReward.Sub(amount)
	k.setMatureTotalReward(ctx, walletAddress, matureRewardBalance)
	return nil
}
