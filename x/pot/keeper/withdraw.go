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
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.TotalRewardPool, targetAddress, amount)
	if err != nil {
		return err
	}
	// mature = 1, slashing= 5  ===> withdraw =0 , slashing=4
	matureRewardBalance := matureReward.Sub(amount)
	k.SetMatureTotalReward(ctx, walletAddress, matureRewardBalance)
	return nil
}
