package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func (k Keeper) Withdraw(ctx sdk.Context, amount sdk.Coin, walletAddress sdk.AccAddress, targetAddress sdk.AccAddress) error {
	matureRewardVal := k.GetMatureTotalReward(ctx, walletAddress)
	matureReward := sdk.NewCoin(k.BondDenom(ctx), matureRewardVal)
	if matureReward.IsLT(amount) {
		return types.ErrInsufficientMatureTotal
	}

	_, err := k.BankKeeper.AddCoins(ctx, targetAddress, sdk.NewCoins(amount))
	if err != nil {
		return err
	}
	newMatureReward := matureReward.Sub(amount)

	k.setMatureTotalReward(ctx, walletAddress, newMatureReward.Amount)

	return nil
}
