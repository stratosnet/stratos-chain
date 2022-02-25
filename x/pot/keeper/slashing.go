package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

func (k Keeper) SlashingResourceNode(ctx sdk.Context, p2pAddr sdk.AccAddress, amt sdk.Int, suspend bool) (err error) {

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

	oldMatureTotal := k.GetMatureTotalReward(ctx, p2pAddr)
	oldImmatureTotal := k.GetImmatureTotalReward(ctx, p2pAddr)
	epoch := k.GetLastReportedEpoch(ctx).Add(sdk.NewInt(1))

	//todo: add previous slashing
	// only slashing the reward token for now.
	slashingCoins := sdk.NewCoins(sdk.NewCoin(k.RewardDenom(ctx), slash.TruncateInt()))

	oldtotal := oldMatureTotal.Add(oldImmatureTotal...)
	if !oldtotal.IsAllGTE(slashingCoins) {
		// need to deduct from stake
		stakeDenom := k.BondDenom(ctx)
		deductFromStake := sdk.NewCoin(stakeDenom, sdk.ZeroInt())
		if tmp := slashingCoins.AmountOf(stakeDenom).Sub(oldtotal.AmountOf(stakeDenom)); tmp.GT(sdk.ZeroInt()) {
			deductFromStake.Add(sdk.NewCoin(stakeDenom, tmp))
			node.Tokens = node.Tokens.Sub(tmp)
		}
		slashingCoins = slashingCoins.Sub(sdk.NewCoins(deductFromStake))
	}
	k.RegisterKeeper.SetResourceNode(ctx, node)
	k.RegisterKeeper.SetLastResourceNodeStake(ctx, node.GetNetworkAddr(), node.Tokens)
	k.SetSlashing(ctx, epoch, p2pAddr, slashingCoins)
	return nil
}
