package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func (k Keeper) SetFoundationAccount(ctx sdk.Context, acc sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(acc)
	store.Set(types.FoundationAccountKey, b)
}

func (k Keeper) GetFoundationAccount(ctx sdk.Context) (acc sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.FoundationAccountKey)
	if b == nil {
		panic("Stored foundation account should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &acc)
	return
}

func (k Keeper) SetInitialUOzonePrice(ctx sdk.Context, price sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(price)
	store.Set(types.InitialUOzonePriceKey, b)
}

func (k Keeper) GetInitialUOzonePrice(ctx sdk.Context) (price sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.InitialUOzonePriceKey)
	if b == nil {
		panic("Stored initial uOzone price should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &price)
	return
}

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

func (k Keeper) setRewardAddressPool(ctx sdk.Context, addressList []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(addressList)
	store.Set(types.RewardAddressPoolKey, b)
}

func (k Keeper) GetRewardAddressPool(ctx sdk.Context) (addressList []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.RewardAddressPoolKey)
	if b == nil {
		return nil
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &addressList)
	return
}

func (k Keeper) setLastMaturedEpoch(ctx sdk.Context, epoch sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(epoch)
	store.Set(types.LastMaturedEpochKey, b)
}

func (k Keeper) getLastMaturedEpoch(ctx sdk.Context) (epoch sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastMaturedEpochKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &epoch)
	return
}

func (k Keeper) setIndividualReward(ctx sdk.Context, acc sdk.AccAddress, epoch sdk.Int, value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set(types.GetIndividualRewardKey(acc, epoch), b)
}

func (k Keeper) GetIndividualReward(ctx sdk.Context, acc sdk.AccAddress, epoch sdk.Int) (value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetIndividualRewardKey(acc, epoch))
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &value)
	return
}

func (k Keeper) setMatureTotalReward(ctx sdk.Context, acc sdk.AccAddress, value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set(types.GetMatureTotalRewardKey(acc), b)
}

func (k Keeper) GetMatureTotalReward(ctx sdk.Context, acc sdk.AccAddress) (value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetMatureTotalRewardKey(acc))
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &value)
	return
}

func (k Keeper) setImmatureTotalReward(ctx sdk.Context, acc sdk.AccAddress, value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set(types.GetImmatureTotalRewardKey(acc), b)
}

func (k Keeper) GetImmatureTotalReward(ctx sdk.Context, acc sdk.AccAddress) (value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetImmatureTotalRewardKey(acc))
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &value)
	return
}

func (k Keeper) setEpochReward(ctx sdk.Context, epoch sdk.Int, value []types.Reward) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	key := types.GetEpochRewardsKey(epoch)
	store.Set(key, b)
}

func (k Keeper) GetEpochReward(ctx sdk.Context, epoch sdk.Int) (value []types.Reward) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetEpochRewardsKey(epoch))
	if b == nil {
		return nil
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &value)
	return
}

func (k Keeper) setRewardsByEpoch(ctx sdk.Context, rewardDetailMap map[string]types.Reward, epoch sdk.Int) {
	var res []types.Reward
	for _, v := range rewardDetailMap {
		newNodeReward := types.NewReward(v.NodeAddress, v.RewardFromMiningPool, v.RewardFromTrafficPool)
		res = append(res, newNodeReward)
	}
	k.setEpochReward(ctx, epoch, res)
}

func (k Keeper) setPotRewardRecord(ctx sdk.Context, epoch sdk.Int, ownerAddr string, value []NodeRewardsInfo) {
	store := ctx.KVStore(k.storeKey)
	res := OwnerRewardsRecord{ctx.BlockHeight(), epoch, value}
	b := k.cdc.MustMarshalBinaryLengthPrefixed(res)
	key := types.GetPotRewardsRecordKey(ownerAddr)
	//ctx.Logger().Info("setKey", "setKey", string(key), "value", res)
	store.Set(key, b)
}

func getIteratorKey(params QueryPotRewardsWithOwnerHeightParams) (prefix []byte) {
	prefix = types.PotRewardsRecordKeyPrefix
	prefix = append(prefix, []byte("potRewards_owner_")...)
	prefix = append(prefix, []byte(params.OwnerAddr.String())...)
	//prefix = append(prefix, []byte("_height_")...)
	return
}

func (k Keeper) GetPotRewardRecords(ctx sdk.Context, params QueryPotRewardsWithOwnerHeightParams) (int64, sdk.Int, []NodeRewardsInfo) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	store := ctx.KVStore(k.storeKey)
	key := getIteratorKey(params)
	//ctx.Logger().Info("QueryKey", "key", key)

	var record OwnerRewardsRecord
	b := store.Get(key)
	if b == nil {
		return 0, sdk.ZeroInt(), nil
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &record)
	//ctx.Logger().Info("Queryres", "Queryres", record)
	var value []NodeRewardsInfo
	value = append(value, record.NodeDetails...)
	return record.PotRewardsRecordHeight, record.PotRewardsRecordEpoch, value
}

//func (k Keeper) GetPotRewardRecords(ctx sdk.Context, params QueryPotRewardsWithOwnerHeightParams) (ownerTotalMatureRewards, ownerTotalImmatureRewards sdk.Int, value []OwnerRewardsInfo) {
//	ownerTotalMatureRewards = sdk.ZeroInt()
//	ownerTotalImmatureRewards = sdk.ZeroInt()
//
//	// filter params
//	height := params.Height
//	var epoch sdk.Int
//	if params.Epoch.IsNil() {
//		epoch = sdk.ZeroInt()
//
//	} else {
//		epoch = params.Epoch
//	}
//
//	currentHeight := ctx.BlockHeight()
//
//	if height > currentHeight {
//		return sdk.ZeroInt(), sdk.ZeroInt(), nil
//	}
//
//	store := ctx.KVStore(k.storeKey)
//	prefix := getIteratorKey(params)
//
//	if height > 0 {
//
//		prefix = append(prefix, []byte(strconv.FormatInt(height, 10))...)
//		prefix = append(prefix, []byte("_epoch_")...)
//
//		if !epoch.IsZero() && epoch.LTE(k.getLastMaturedEpoch(ctx)) {
//			prefix = append(prefix, []byte(epoch.String())...)
//		}
//	}
//
//	ctx.Logger().Info("QueryPrefix", "prefix", prefix)
//	iter := sdk.KVStorePrefixIterator(store, prefix)
//
//	defer iter.Close()
//	ctx.Logger().Info("iter.Valid()", "iter.Valid()", iter.Valid())
//	if !iter.Valid() {
//		return sdk.ZeroInt(), sdk.ZeroInt(), nil
//	}
//	delimiter := []byte{'_'}
//	var (
//		record []NodeRewardsInfo
//		//value OwnerRewardsInfo
//		recordHeight int64
//		recordEpoch  sdk.Int
//		err          error
//	)
//	for ; iter.Valid(); iter.Next() {
//		key := iter.Key()
//		ctx.Logger().Info("Querykey", "key", key)
//		splitSlice := bytes.Split(key, delimiter)
//		//recordOwnerAddr := splitSlice[2]
//		recordHeight, err = strconv.ParseInt(string(splitSlice[4]), 10, 64)
//		ctx.Logger().Info("recordHeight", "recordHeight", recordHeight)
//		if err != nil {
//			return sdk.ZeroInt(), sdk.ZeroInt(), nil
//		}
//		recordEpochInt64, err := strconv.ParseInt(string(splitSlice[6]), 10, 64)
//		ctx.Logger().Info("recordEpochInt64", "recordEpochInt64", recordEpochInt64)
//		if err != nil {
//			return sdk.ZeroInt(), sdk.ZeroInt(), nil
//		}
//		recordEpoch = sdk.NewInt(recordEpochInt64)
//
//		if !epoch.IsZero() {
//			if recordEpochInt64 == epoch.Int64() {
//				b := store.Get(key)
//				if b == nil {
//					continue
//				}
//				k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &record)
//
//				for _, v := range record {
//					newRecord := OwnerRewardsInfo{recordHeight, recordEpoch, v}
//					ownerTotalMatureRewards = ownerTotalMatureRewards.Add(v.MatureTotalReward.Amount)
//					ownerTotalImmatureRewards = ownerTotalImmatureRewards.Add(v.ImmatureTotalReward.Amount)
//					value = append(value, newRecord)
//				}
//			}
//		} else {
//			b := store.Get(key)
//			if b == nil {
//				continue
//			}
//			k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &record)
//
//			for _, v := range record {
//				newRecord := OwnerRewardsInfo{recordHeight, recordEpoch, v}
//				ownerTotalMatureRewards = ownerTotalMatureRewards.Add(v.MatureTotalReward.Amount)
//				ownerTotalImmatureRewards = ownerTotalImmatureRewards.Add(v.ImmatureTotalReward.Amount)
//				value = append(value, newRecord)
//			}
//		}
//
//	}
//	return
//}

//func (k Keeper) getMatchedPotRewardRecords(
//	store sdk.KVStore,
//	key, recordHeight []byte,
//	params QueryPotRewardsWithOwnerHeightParams,
//) (res []NodeRewardsInfo) {
//
//	if params.Height != 0 {
//		if bytes.Equal(recordHeight, []byte(strconv.FormatInt(params.Height, 10))) {
//			record := k.getRewardsRecord(store, key)
//			if record != nil {
//				res = append(res, record...)
//			}
//			return
//		}
//
//	} else {
//		record := k.getRewardsRecord(store, key)
//		if record != nil {
//			res = append(res, record...)
//		}
//	}
//	return
//}

//func (k Keeper) getRewardsRecord(store sdk.KVStore, key []byte) (record []NodeRewardsInfo) {
//	b := store.Get(key)
//	if b == nil {
//		return nil
//	}
//	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &record)
//	return record
//}
