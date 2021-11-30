package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func (k Keeper) setTotalMinedTokens(ctx sdk.Context, totalMinedToken sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(totalMinedToken)
	store.Set(types.TotalMinedTokensKey, b)
}

func (k Keeper) GetTotalMinedTokens(ctx sdk.Context) (totalMinedToken sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.TotalMinedTokensKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &totalMinedToken)
	return
}

func (k Keeper) setMinedTokens(ctx sdk.Context, epoch sdk.Int, minedToken sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(minedToken)
	store.Set(types.GetMinedTokensKey(epoch), b)
}

func (k Keeper) GetMinedTokens(ctx sdk.Context, epoch sdk.Int) (minedToken sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetMinedTokensKey(epoch))
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &minedToken)
	return
}

func (k Keeper) SetTotalUnissuedPrepay(ctx sdk.Context, totalUnissuedPrepay sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(totalUnissuedPrepay)
	store.Set(types.TotalUnissuedPrepayKey, b)
}

func (k Keeper) GetTotalUnissuedPrepay(ctx sdk.Context) (totalUnissuedPrepay sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.TotalUnissuedPrepayKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &totalUnissuedPrepay)
	return
}

func (k Keeper) setRewardAddressPool(ctx sdk.Context, walletAddressList []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(walletAddressList)
	store.Set(types.RewardAddressPoolKey, b)
}

func (k Keeper) GetRewardAddressPool(ctx sdk.Context) (walletAddressList []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.RewardAddressPoolKey)
	if b == nil {
		return nil
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &walletAddressList)
	return
}

func (k Keeper) setLastReportedEpoch(ctx sdk.Context, epoch sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(epoch)
	store.Set(types.LastReportedEpochKey, b)
}

func (k Keeper) GetLastReportedEpoch(ctx sdk.Context) (epoch sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastReportedEpochKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &epoch)
	return
}

func (k Keeper) setIndividualReward(ctx sdk.Context, walletAddress sdk.AccAddress, epoch sdk.Int, value types.Reward) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set(types.GetIndividualRewardKey(walletAddress, epoch), b)
}

func (k Keeper) GetIndividualReward(ctx sdk.Context, walletAddress sdk.AccAddress, epoch sdk.Int) (value types.Reward, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetIndividualRewardKey(walletAddress, epoch))
	if b == nil {
		return value, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &value)
	return value, true
}

func (k Keeper) setMatureTotalReward(ctx sdk.Context, walletAddress sdk.AccAddress, value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set(types.GetMatureTotalRewardKey(walletAddress), b)
}

func (k Keeper) GetMatureTotalReward(ctx sdk.Context, walletAddress sdk.AccAddress) (value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetMatureTotalRewardKey(walletAddress))
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &value)
	return
}

func (k Keeper) setImmatureTotalReward(ctx sdk.Context, walletAddress sdk.AccAddress, value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set(types.GetImmatureTotalRewardKey(walletAddress), b)
}

func (k Keeper) GetImmatureTotalReward(ctx sdk.Context, walletAddress sdk.AccAddress) (value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetImmatureTotalRewardKey(walletAddress))
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &value)
	return
}
