package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

func (k Keeper) SetInitialGenesisStakeTotal(ctx sdk.Context, stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalLengthPrefixed(&stratos.Int{Value: &stake})
	store.Set(types.InitialGenesisStakeTotalKey, b)
}

func (k Keeper) GetInitialGenesisStakeTotal(ctx sdk.Context) (stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.InitialGenesisStakeTotalKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	value := stratos.Int{}
	k.cdc.MustUnmarshalLengthPrefixed(b, &value)
	stake = *value.Value
	return
}

func (k Keeper) SetRemainingOzoneLimit(ctx sdk.Context, value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalLengthPrefixed(&stratos.Int{Value: &value})
	store.Set(types.UpperBoundOfTotalOzoneKey, b)
}

func (k Keeper) GetRemainingOzoneLimit(ctx sdk.Context) (value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.UpperBoundOfTotalOzoneKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	intVal := stratos.Int{}
	k.cdc.MustUnmarshalLengthPrefixed(b, &intVal)
	value = *intVal.Value
	return
}

func (k Keeper) IsUnbondable(ctx sdk.Context, unbondAmt sdk.Int) bool {
	remaining := k.GetRemainingOzoneLimit(ctx)
	stakeNozRate := k.GetStakeNozRate(ctx)
	return remaining.ToDec().GTE(unbondAmt.ToDec().Quo(stakeNozRate))
}

// SetUnbondingNode sets the unbonding node
func (k Keeper) SetUnbondingNode(ctx sdk.Context, ubd types.UnbondingNode) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&ubd)
	networkAddr, err := stratos.SdsAddressFromBech32(ubd.GetNetworkAddr())
	if err != nil {
		return
	}
	key := types.GetUBDNodeKey(networkAddr)
	store.Set(key, bz)
}

// RemoveUnbondingNode removes the unbonding node object
func (k Keeper) RemoveUnbondingNode(ctx sdk.Context, networkAddr stratos.SdsAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDNodeKey(networkAddr)
	store.Delete(key)
}

// GetUnbondingNode return a unbonding node
func (k Keeper) GetUnbondingNode(ctx sdk.Context, networkAddr stratos.SdsAddress) (ubd types.UnbondingNode, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDNodeKey(networkAddr)
	value := store.Get(key)
	if value == nil {
		return ubd, false
	}
	k.cdc.MustUnmarshalLengthPrefixed(value, &ubd)
	return ubd, true
}

// SetUnbondingNodeQueueTimeSlice sets a specific unbonding queue timeslice.
func (k Keeper) SetUnbondingNodeQueueTimeSlice(ctx sdk.Context, timestamp time.Time, networkAddrs []string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&stratos.SdsAddresses{Addresses: networkAddrs})
	store.Set(types.GetUBDTimeKey(timestamp), bz)
}

// GetUnbondingNodeQueueTimeSlice gets a specific unbonding queue timeslice. A timeslice is a slice of DVPairs
// corresponding to unbonding delegations that expire at a certain time.
func (k Keeper) GetUnbondingNodeQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (networkAddrs []string) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetUBDTimeKey(timestamp))
	if bz == nil {
		return make([]string, 0)
	}

	addrValue := stratos.SdsAddresses{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &addrValue)
	networkAddrs = addrValue.GetAddresses()
	return networkAddrs
}

// UnbondingNodeQueueIterator returns all the unbonding queue timeslices from time 0 until endTime
func (k Keeper) UnbondingNodeQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UBDNodeQueueKey, sdk.InclusiveEndBytes(types.GetUBDTimeKey(endTime)))
}

func (k Keeper) SetBondedResourceNodeCnt(ctx sdk.Context, count sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&stratos.Int{Value: &count})
	store.Set(types.ResourceNodeCntKey, bz)
}

func (k Keeper) GetBondedResourceNodeCnt(ctx sdk.Context) (balance sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ResourceNodeCntKey)
	if bz == nil {
		return sdk.ZeroInt()
	}
	intValue := stratos.Int{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intValue)
	balance = *intValue.Value
	return
}

func (k Keeper) SetBondedMetaNodeCnt(ctx sdk.Context, count sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&stratos.Int{Value: &count})
	store.Set(types.MetaNodeCntKey, bz)
}

func (k Keeper) GetBondedMetaNodeCnt(ctx sdk.Context) (balance sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MetaNodeCntKey)
	if bz == nil {
		return sdk.ZeroInt()
	}
	intValue := stratos.Int{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intValue)
	balance = *intValue.Value
	return
}

func (k Keeper) SetMetaNodeRegistrationVotePool(ctx sdk.Context, votePool types.MetaNodeRegistrationVotePool) {
	nodeAddr := votePool.GetNetworkAddress()
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&votePool)
	node, _ := stratos.SdsAddressFromBech32(nodeAddr)
	store.Set(types.GetMetaNodeRegistrationVotesKey(node), bz)
}
func (k Keeper) GetMetaNodeRegistrationVotePool(ctx sdk.Context, nodeAddr stratos.SdsAddress) (votePool types.MetaNodeRegistrationVotePool, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMetaNodeRegistrationVotesKey(nodeAddr))
	if bz == nil {
		return votePool, false
	}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &votePool)
	return votePool, true
}

func (k Keeper) SetEffectiveTotalStake(ctx sdk.Context, stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&stratos.Int{Value: &stake})
	store.Set(types.EffectiveGenesisStakeTotalKey, bz)
}

func (k Keeper) GetEffectiveTotalStake(ctx sdk.Context) (stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EffectiveGenesisStakeTotalKey)
	if bz == nil {
		return sdk.ZeroInt()
	}
	intValue := stratos.Int{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intValue)
	stake = *intValue.Value
	return
}

func (k Keeper) SetStakeNozRate(ctx sdk.Context, stakeNozRate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&stratos.Dec{Value: &stakeNozRate})
	store.Set(types.StakeNozRateKey, bz)
}

func (k Keeper) GetStakeNozRate(ctx sdk.Context) (stakeNozRate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.StakeNozRateKey)
	if bz == nil {
		panic("Stored stake noz rate should not be nil")
	}
	decValue := stratos.Dec{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &decValue)
	stakeNozRate = *decValue.Value
	return
}

// IteratorSlashingInfo Iteration for each slashing
func (k Keeper) IteratorSlashingInfo(ctx sdk.Context, handler func(walletAddress sdk.AccAddress, slashing sdk.Int) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.SlashingPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		walletAddress := sdk.AccAddress(iter.Key()[len(types.SlashingPrefix):])
		intValue := stratos.Int{}
		k.cdc.MustUnmarshalLengthPrefixed(iter.Value(), &intValue)
		slashing := *intValue.Value
		if handler(walletAddress, slashing) {
			break
		}
	}
}

func (k Keeper) SetSlashing(ctx sdk.Context, walletAddress sdk.AccAddress, slashing sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.GetSlashingKey(walletAddress)
	bz := k.cdc.MustMarshalLengthPrefixed(&stratos.Int{Value: &slashing})
	store.Set(storeKey, bz)
}

func (k Keeper) GetSlashing(ctx sdk.Context, walletAddress sdk.AccAddress) (res sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSlashingKey(walletAddress))
	if bz == nil {
		return sdk.ZeroInt()
	}
	intValue := stratos.Int{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intValue)
	res = *intValue.Value
	return
}
