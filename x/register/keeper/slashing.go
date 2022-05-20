package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// DeductSlashing deduct slashing amount from coins, return the coins that after deduction
func (k Keeper) DeductSlashing(ctx sdk.Context, walletAddress sdk.AccAddress, coins sdk.Coins) sdk.Coins {
	slashing := k.GetSlashing(ctx, walletAddress)
	if slashing.LTE(sdk.ZeroInt()) || coins.Empty() || coins.IsZero() {
		return coins
	}

	ret := sdk.Coins{}
	for _, coin := range coins {
		if coin.Amount.GTE(slashing) {
			coin = coin.Sub(sdk.NewCoin(coin.Denom, slashing))
			ret = ret.Add(coin)
			slashing = sdk.ZeroInt()
		} else {
			slashing = slashing.Sub(coin.Amount)
			coin = sdk.NewCoin(coin.Denom, sdk.ZeroInt())
			ret = ret.Add(coin)
		}
	}
	k.SetSlashing(ctx, walletAddress, slashing)
	return ret
}

func (k Keeper) IteratorSlashingInfo(ctx sdk.Context, handler func(walletAddress sdk.AccAddress, slashing sdk.Int) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.SlashingPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		walletAddress := sdk.AccAddress(iter.Key()[len(types.SlashingPrefix):])
		var slashing sdk.Int
		types.ModuleCdc.MustUnmarshalLengthPrefixed(iter.Value(), &slashing)
		if handler(walletAddress, slashing) {
			break
		}
	}
}

func (k Keeper) SetSlashing(ctx sdk.Context, walletAddress sdk.AccAddress, slashing sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.GetSlashingKey(walletAddress)
	bz := types.ModuleCdc.MustMarshalLengthPrefixed(slashing)
	store.Set(storeKey, bz)
}

func (k Keeper) GetSlashing(ctx sdk.Context, walletAddress sdk.AccAddress) (res sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSlashingKey(walletAddress))
	if bz == nil {
		return sdk.ZeroInt()
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(bz, &res)
	return
}
