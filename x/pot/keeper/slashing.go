package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

func (k Keeper) SlashingResourceNode(ctx sdk.Context, p2pAddr sdk.AccAddress, walletAddr sdk.AccAddress, amt sdk.Int, suspend bool) (err error) {

	node, ok := k.RegisterKeeper.GetResourceNode(ctx, p2pAddr)
	if !ok {
		return regtypes.ErrNoResourceNodeFound
	}
	if suspend == node.Suspend {
		return types.ErrNodeStatusSuspend
	}
	node.Suspend = suspend

	//slashing amt is equivalent to reward traffic calculation
	_, slash := k.getTrafficReward(ctx, []types.SingleWalletVolume{{
		WalletAddress: node.GetOwnerAddr(),
		Volume:        amt,
	}})

	oldMatureTotal := k.GetMatureTotalReward(ctx, walletAddr)
	oldImmatureTotal := k.GetImmatureTotalReward(ctx, walletAddr)
	oldSlashing := k.GetSlashing(ctx, p2pAddr)

	stakeDenom := k.BondDenom(ctx)
	rewardDenom := k.RewardDenom(ctx)
	// only slashing the reward token for now.
	slashingCoins := sdk.NewCoins(sdk.NewCoin(rewardDenom, slash.TruncateInt())).Add(oldSlashing.SlashingCoins...)

	// deduct from matured reward
	deductFromMatureReward := sdk.ZeroInt()
	if slashingCoins.AmountOf(rewardDenom).GT(oldMatureTotal.AmountOf(rewardDenom)) {
		deductFromMatureReward = oldMatureTotal.AmountOf(rewardDenom)
	} else {
		deductFromMatureReward = slashingCoins.AmountOf(rewardDenom)
	}
	newMatureTotal := oldMatureTotal.Sub(sdk.NewCoins(sdk.NewCoin(rewardDenom, deductFromMatureReward)))
	slashingCoins = slashingCoins.Sub(sdk.NewCoins(sdk.NewCoin(rewardDenom, deductFromMatureReward)))

	// TODO: (add to reward distribution?) deduct from immature reward (would affect immatureToMature)

	oldTotal := oldMatureTotal.Add(oldImmatureTotal...)
	if !oldTotal.IsAllGTE(slashingCoins) {
		// need to deduct from stake
		deductFromStake := sdk.ZeroInt()
		if slashingCoins.AmountOf(rewardDenom).GT(oldTotal.AmountOf(stakeDenom)) {
			deductFromStake = oldTotal.AmountOf(stakeDenom)
		} else {
			deductFromStake = slashingCoins.AmountOf(rewardDenom)
		}
		node = node.SubToken(deductFromStake)
		slashingCoins = slashingCoins.Sub(sdk.NewCoins(sdk.NewCoin(rewardDenom, deductFromStake)))
	}
	k.RegisterKeeper.SetResourceNode(ctx, node)
	k.RegisterKeeper.SetLastResourceNodeStake(ctx, node.GetNetworkAddr(), node.Tokens)

	k.setMatureTotalReward(ctx, walletAddr, newMatureTotal)
	newSlashing := types.NewSlashing(p2pAddr, slashingCoins)
	k.SetSlashing(ctx, p2pAddr, newSlashing)
	return nil
}
