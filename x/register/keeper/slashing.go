package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// DeductSlashing deduct slashing amount from coins, return the coins that after deduction
func (k Keeper) DeductSlashing(ctx sdk.Context, walletAddress sdk.AccAddress, coins sdk.Coins) (remaining, deducted sdk.Coins) {
	slashing := k.GetSlashing(ctx, walletAddress)
	remaining = sdk.Coins{}
	deducted = sdk.Coins{}
	if slashing.LTE(sdk.ZeroInt()) || coins.Empty() || coins.IsZero() {
		return coins, deducted
	}
	fmt.Println("!!!!!!!!!!!!!! DeductSlashing(): walletAddress = " + walletAddress.String() + ", coins = " + coins.String())
	for _, coin := range coins {
		if coin.Amount.GTE(slashing) {
			coin = coin.Sub(sdk.NewCoin(coin.Denom, slashing))
			remaining = remaining.Add(coin)
			deducted = deducted.Add(sdk.NewCoin(coin.Denom, slashing))
			slashing = sdk.ZeroInt()
		} else {
			slashing = slashing.Sub(coin.Amount)
			deducted = deducted.Add(coin)
			coin = sdk.NewCoin(coin.Denom, sdk.ZeroInt())
			remaining = remaining.Add(coin)
		}
	}
	k.SetSlashing(ctx, walletAddress, slashing)
	fmt.Println("!!!!!!!!!!!!!! DeductSlashing(): remaining = " + remaining.String() + ", deducted = " + deducted.String())
	return remaining, deducted
}

// Iteration for each slashing
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
