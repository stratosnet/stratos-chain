package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func (k Keeper) Withdraw(ctx sdk.Context, amount sdk.Coin, nodeAddress sdk.AccAddress, ownerAddress sdk.AccAddress) error {
	if k.checkOwner(ctx, nodeAddress, ownerAddress) == false {
		return types.ErrNotTheOwner
	}

	matureRewardVal := k.getMatureTotalReward(ctx, nodeAddress)
	matureReward := sdk.NewCoin(k.BondDenom(ctx), matureRewardVal)
	if matureReward.IsLT(amount) {
		return types.ErrInsufficientMatureTotal
	}

	_, err := k.BankKeeper.AddCoins(ctx, ownerAddress, sdk.NewCoins(amount))
	if err != nil {
		return err
	}
	newMatureReward := matureReward.Sub(amount)

	k.setMatureTotalReward(ctx, nodeAddress, newMatureReward.Amount)

	return nil
}

func (k Keeper) checkOwner(ctx sdk.Context, nodeAddress sdk.AccAddress, ownerAddress sdk.Address) (found bool) {
	resourceNode, found := k.RegisterKeeper.GetResourceNode(ctx, nodeAddress)
	if found && resourceNode.OwnerAddress.Equals(ownerAddress) {
		return true
	}

	indexingNode, found := k.RegisterKeeper.GetIndexingNode(ctx, nodeAddress)
	if found && indexingNode.OwnerAddress.Equals(ownerAddress) {
		return true
	}

	return false
}
