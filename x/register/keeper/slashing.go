package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/register/types"
)

// deduct slashing amount from coins, return the coins that after deduction
func (k Keeper) DeductSlashing(ctx sdk.Context, walletAddress sdk.AccAddress, coins sdk.Coins) sdk.Coins {
	slashing := k.GetSlashing(ctx, walletAddress)
	if slashing.LTE(sdk.ZeroInt()) || coins.Empty() || coins.IsZero() {
		return coins
	}

	for _, coin := range coins {
		if coin.Amount.GTE(slashing) {
			coin = coin.Sub(sdk.NewCoin(coin.Denom, slashing))
			slashing = sdk.ZeroInt()
			break
		} else {
			coin = sdk.NewCoin(coin.Denom, sdk.ZeroInt())
			slashing = slashing.Sub(coin.Amount)
		}
	}
	k.SetSlashing(ctx, walletAddress, slashing)
	return coins
}

func (k Keeper) IteratorSlashingInfo(ctx sdk.Context, handler func(walletAddress sdk.AccAddress, slashing sdk.Int) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.SlashingPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		walletAddress := sdk.AccAddress(iter.Key()[len(types.SlashingPrefix):])
		var slashing sdk.Int
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &slashing)
		if handler(walletAddress, slashing) {
			break
		}
	}
}

func (k Keeper) SetSlashing(ctx sdk.Context, walletAddress sdk.AccAddress, slashing sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.GetSlashingKey(walletAddress)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(slashing)
	store.Set(storeKey, bz)
}

func (k Keeper) GetSlashing(ctx sdk.Context, walletAddress sdk.AccAddress) (res sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSlashingKey(walletAddress))
	if bz == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}
