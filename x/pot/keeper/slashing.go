package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

func (k Keeper) SlashingResourceNode(ctx sdk.Context, p2pAddr sdk.AccAddress, walletAddr sdk.AccAddress, ozAmt sdk.Int, suspend bool) (amt sdk.Int, err error) {

	node, ok := k.RegisterKeeper.GetResourceNode(ctx, p2pAddr)
	if !ok {
		return sdk.ZeroInt(), regtypes.ErrNoResourceNodeFound
	}
	//if suspend == node.Suspend {
	//	return types.ErrNodeStatusSuspend
	//}
	node.Suspend = suspend

	//slashing amt is equivalent to reward traffic calculation
	_, slash := k.getTrafficReward(ctx, []types.SingleWalletVolume{{
		WalletAddress: node.GetOwnerAddr(),
		Volume:        ozAmt,
	}})

	oldSlashing := k.GetSlashing(ctx, p2pAddr)

	// only slashing the reward token for now.
	newSlashing := oldSlashing.Add(slash.TruncateInt())

	// TODO: (add to reward distribution?) deduct from immature reward (would affect immatureToMature)

	k.RegisterKeeper.SetResourceNode(ctx, node)
	k.RegisterKeeper.SetLastResourceNodeStake(ctx, node.GetNetworkAddr(), node.Tokens)

	k.SetSlashing(ctx, p2pAddr, newSlashing)
	return slash.TruncateInt(), nil
}
