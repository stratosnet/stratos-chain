package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

func (k Keeper) SlashingResourceNode(ctx sdk.Context, p2pAddr sdk.AccAddress, walletAddr sdk.AccAddress, ozAmt sdk.Int, suspend bool) (err error) {

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
		Volume:        ozAmt,
	}})

	oldSlashing := k.GetSlashing(ctx, p2pAddr)

	// only slashing the reward token for now.
	newSlashing := oldSlashing.Add(slash.TruncateInt())

	// TODO: (add to reward distribution?) deduct from immature reward (would affect immatureToMature)

	k.RegisterKeeper.SetResourceNode(ctx, node)
	k.RegisterKeeper.SetLastResourceNodeStake(ctx, node.GetNetworkAddr(), node.Tokens)

	k.SetSlashing(ctx, p2pAddr, newSlashing)
	return nil
}

func (k Keeper) IteratorSlashingInfo(ctx sdk.Context, handler func(p2pAddress sdk.AccAddress, slashing sdk.Int) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.SlashingPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.AccAddress(iter.Key()[len(types.SlashingPrefix):])
		var slashing sdk.Int
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &slashing)
		if handler(addr, slashing) {
			break
		}
	}
}
