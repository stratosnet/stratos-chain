package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func (k Keeper) SetTotalMinedTokens(ctx sdk.Context, totalMinedToken sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	b := types.ModuleCdc.MustMarshalLengthPrefixed(totalMinedToken)
	store.Set(types.TotalMinedTokensKey, b)
}

func (k Keeper) GetTotalMinedTokens(ctx sdk.Context) (totalMinedToken sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.TotalMinedTokensKey)
	if b == nil {
		return sdk.NewCoin(k.RewardDenom(ctx), sdk.ZeroInt())
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(b, &totalMinedToken)
	return
}

func (k Keeper) setMinedTokens(ctx sdk.Context, epoch sdk.Int, minedToken sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	b := types.ModuleCdc.MustMarshalLengthPrefixed(minedToken)
	store.Set(types.GetMinedTokensKey(epoch), b)
}

func (k Keeper) GetMinedTokens(ctx sdk.Context, epoch sdk.Int) (minedToken sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetMinedTokensKey(epoch))
	if b == nil {
		return sdk.NewCoin(k.RewardDenom(ctx), sdk.ZeroInt())
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(b, &minedToken)
	return
}

func (k Keeper) SetLastReportedEpoch(ctx sdk.Context, epoch sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := types.ModuleCdc.MustMarshalLengthPrefixed(epoch)
	store.Set(types.LastReportedEpochKey, b)
}

func (k Keeper) GetLastReportedEpoch(ctx sdk.Context) (epoch sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastReportedEpochKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(b, &epoch)
	return
}

func (k Keeper) SetIndividualReward(ctx sdk.Context, walletAddress sdk.AccAddress, epoch sdk.Int, value types.Reward) {
	store := ctx.KVStore(k.storeKey)
	b := types.ModuleCdc.MustMarshalLengthPrefixed(value)
	store.Set(types.GetIndividualRewardKey(walletAddress, epoch), b)
}

func (k Keeper) GetIndividualReward(ctx sdk.Context, walletAddress sdk.AccAddress, epoch sdk.Int) (value types.Reward, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetIndividualRewardKey(walletAddress, epoch))
	if b == nil {
		return value, false
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(b, &value)
	return value, true
}

func (k Keeper) SetMatureTotalReward(ctx sdk.Context, walletAddress sdk.AccAddress, value sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	b := types.ModuleCdc.MustMarshalLengthPrefixed(value)
	store.Set(types.GetMatureTotalRewardKey(walletAddress), b)
}

func (k Keeper) GetMatureTotalReward(ctx sdk.Context, walletAddress sdk.AccAddress) (value sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetMatureTotalRewardKey(walletAddress))
	if b == nil {
		return sdk.Coins{}
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(b, &value)
	return
}

func (k Keeper) SetImmatureTotalReward(ctx sdk.Context, walletAddress sdk.AccAddress, value sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	b := types.ModuleCdc.MustMarshalLengthPrefixed(value)
	store.Set(types.GetImmatureTotalRewardKey(walletAddress), b)
}

func (k Keeper) GetImmatureTotalReward(ctx sdk.Context, walletAddress sdk.AccAddress) (value sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetImmatureTotalRewardKey(walletAddress))
	if b == nil {
		return sdk.Coins{}
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(b, &value)
	return
}

func (k Keeper) GetVolumeReport(ctx sdk.Context, epoch sdk.Int) (res types.VolumeReportRecord) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VolumeReportStoreKey(epoch))
	if bz == nil {
		return types.VolumeReportRecord{}
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(bz, &res)
	return res
}

func (k Keeper) SetVolumeReport(ctx sdk.Context, epoch sdk.Int, reportRecord types.VolumeReportRecord) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.VolumeReportStoreKey(epoch)
	bz := types.ModuleCdc.MustMarshalLengthPrefixed(reportRecord)
	store.Set(storeKey, bz)
}
