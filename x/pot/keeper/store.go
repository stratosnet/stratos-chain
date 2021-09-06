package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func (k Keeper) SetFoundationAccount(ctx sdk.Context, acc sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := k.Cdc.MustMarshalBinaryLengthPrefixed(acc)
	store.Set(types.FoundationAccountKey, b)
}

func (k Keeper) GetFoundationAccount(ctx sdk.Context) (acc sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.FoundationAccountKey)
	if b == nil {
		panic("Stored foundation account should not have been nil")
	}
	k.Cdc.MustUnmarshalBinaryLengthPrefixed(b, &acc)
	return
}

func (k Keeper) SetInitialUOzonePrice(ctx sdk.Context, price sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.Cdc.MustMarshalBinaryLengthPrefixed(price)
	store.Set(types.InitialUOzonePriceKey, b)
}

func (k Keeper) GetInitialUOzonePrice(ctx sdk.Context) (price sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.InitialUOzonePriceKey)
	if b == nil {
		panic("Stored initial uOzone price should not have been nil")
	}
	k.Cdc.MustUnmarshalBinaryLengthPrefixed(b, &price)
	return
}

func (k Keeper) setTotalMinedTokens(ctx sdk.Context, totalMinedToken sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.Cdc.MustMarshalBinaryLengthPrefixed(totalMinedToken)
	store.Set(types.TotalMinedTokensKey, b)
}

func (k Keeper) GetTotalMinedTokens(ctx sdk.Context) (totalMinedToken sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.TotalMinedTokensKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.Cdc.MustUnmarshalBinaryLengthPrefixed(b, &totalMinedToken)
	return
}

func (k Keeper) setMinedTokens(ctx sdk.Context, epoch sdk.Int, minedToken sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.Cdc.MustMarshalBinaryLengthPrefixed(minedToken)
	store.Set(types.GetMinedTokensKey(epoch), b)
}

func (k Keeper) GetMinedTokens(ctx sdk.Context, epoch sdk.Int) (minedToken sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetMinedTokensKey(epoch))
	if b == nil {
		return sdk.ZeroInt()
	}
	k.Cdc.MustUnmarshalBinaryLengthPrefixed(b, &minedToken)
	return
}

func (k Keeper) SetTotalUnissuedPrepay(ctx sdk.Context, totalUnissuedPrepay sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.Cdc.MustMarshalBinaryLengthPrefixed(totalUnissuedPrepay)
	store.Set(types.TotalUnissuedPrepayKey, b)
}

func (k Keeper) GetTotalUnissuedPrepay(ctx sdk.Context) (totalUnissuedPrepay sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.TotalUnissuedPrepayKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.Cdc.MustUnmarshalBinaryLengthPrefixed(b, &totalUnissuedPrepay)
	return
}

func (k Keeper) setRewardAddressPool(ctx sdk.Context, addressList []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := k.Cdc.MustMarshalBinaryLengthPrefixed(addressList)
	store.Set(types.RewardAddressPoolKey, b)
}

func (k Keeper) GetRewardAddressPool(ctx sdk.Context) (addressList []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.RewardAddressPoolKey)
	if b == nil {
		return nil
	}
	k.Cdc.MustUnmarshalBinaryLengthPrefixed(b, &addressList)
	return
}

func (k Keeper) setLastMaturedEpoch(ctx sdk.Context, epoch sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.Cdc.MustMarshalBinaryLengthPrefixed(epoch)
	store.Set(types.LastMaturedEpochKey, b)
}

func (k Keeper) getLastMaturedEpoch(ctx sdk.Context) (epoch sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastMaturedEpochKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.Cdc.MustUnmarshalBinaryLengthPrefixed(b, &epoch)
	return
}

func (k Keeper) setIndividualReward(ctx sdk.Context, acc sdk.AccAddress, epoch sdk.Int, value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.Cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set(types.GetIndividualRewardKey(acc, epoch), b)
}

func (k Keeper) GetIndividualReward(ctx sdk.Context, acc sdk.AccAddress, epoch sdk.Int) (value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetIndividualRewardKey(acc, epoch))
	if b == nil {
		return sdk.ZeroInt()
	}
	k.Cdc.MustUnmarshalBinaryLengthPrefixed(b, &value)
	return
}

func (k Keeper) setMatureTotalReward(ctx sdk.Context, acc sdk.AccAddress, value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.Cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set(types.GetMatureTotalRewardKey(acc), b)
}

func (k Keeper) GetMatureTotalReward(ctx sdk.Context, acc sdk.AccAddress) (value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetMatureTotalRewardKey(acc))
	if b == nil {
		return sdk.ZeroInt()
	}
	k.Cdc.MustUnmarshalBinaryLengthPrefixed(b, &value)
	return
}

func (k Keeper) setImmatureTotalReward(ctx sdk.Context, acc sdk.AccAddress, value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.Cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set(types.GetImmatureTotalRewardKey(acc), b)
}

func (k Keeper) GetImmatureTotalReward(ctx sdk.Context, acc sdk.AccAddress) (value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetImmatureTotalRewardKey(acc))
	if b == nil {
		return sdk.ZeroInt()
	}
	k.Cdc.MustUnmarshalBinaryLengthPrefixed(b, &value)
	return
}

func (k Keeper) setEpochReward(ctx sdk.Context, epoch sdk.Int, value []types.Reward) {
	//ctx.Logger().Info("Enter setEpochReward", "value", value)
	//var res []types.Reward
	//for _, v := range value {
	//	newNodeReward := types.NewReward(v.NodeAddress, v.RewardFromMiningPool, v.RewardFromTrafficPool)
	//	res = append(res, newNodeReward)
	//}
	//ctx.Logger().Info("InsetEpochReward", "res", res)
	store := ctx.KVStore(k.storeKey)
	b := k.Cdc.MustMarshalBinaryLengthPrefixed(value)
	ctx.Logger().Info("In setEpochReward", "MustMarshalBinaryLengthPrefixed", true, "b", b)
	key := types.GetEpochRewardsKey(epoch)
	ctx.Logger().Info("In setEpochReward", "key", key)
	store.Set(key, b)
	ctx.Logger().Info("Leave setEpochReward")
}

func (k Keeper) GetEpochReward(ctx sdk.Context, epoch sdk.Int) (value []types.Reward) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetEpochRewardsKey(epoch))
	if b == nil {
		return nil
	}
	k.Cdc.MustUnmarshalBinaryLengthPrefixed(b, &value)
	return
}
