package keeper

import (
	"container/list"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the register store
type Keeper struct {
	storeKey              sdk.StoreKey
	cdc                   *codec.Codec
	paramSpace            params.Subspace
	accountKeeper         auth.AccountKeeper
	bankKeeper            bank.Keeper
	resourceNodeCache     map[string]cachedResourceNode
	resourceNodeCacheList *list.List
	indexingNodeCache     map[string]cachedIndexingNode
	indexingNodeCacheList *list.List
}

// NewKeeper creates a register keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper) Keeper {

	keeper := Keeper{
		storeKey:              key,
		cdc:                   cdc,
		paramSpace:            paramSpace.WithKeyTable(types.ParamKeyTable()),
		accountKeeper:         accountKeeper,
		bankKeeper:            bankKeeper,
		resourceNodeCache:     make(map[string]cachedResourceNode, resourceNodeCacheSize),
		resourceNodeCacheList: list.New(),
		indexingNodeCache:     make(map[string]cachedIndexingNode, indexingNodeCacheSize),
		indexingNodeCacheList: list.New(),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetInitialGenesisStakeTotal(ctx sdk.Context, stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(stake)
	store.Set(types.InitialGenesisStakeTotalKey, b)
}

func (k Keeper) GetInitialGenesisStakeTotal(ctx sdk.Context) (stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.InitialGenesisStakeTotalKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &stake)
	return
}

func (k Keeper) SetRemainingOzoneLimit(ctx sdk.Context, value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set(types.UpperBoundOfTotalOzoneKey, b)
}

func (k Keeper) GetRemainingOzoneLimit(ctx sdk.Context) (value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.UpperBoundOfTotalOzoneKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &value)
	return
}

func (k Keeper) increaseOzoneLimitByAddStake(ctx sdk.Context, stake sdk.Int) (ozoneLimitChange sdk.Int) {
	initialGenesisDeposit := k.GetInitialGenesisStakeTotal(ctx).ToDec() //ustos
	if initialGenesisDeposit.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("initialGenesisDeposit is zero, increase ozone limit failed")
		return
	}
	currentLimit := k.GetRemainingOzoneLimit(ctx).ToDec() //uoz
	limitToAdd := currentLimit.Mul(stake.ToDec()).Quo(initialGenesisDeposit)
	newLimit := currentLimit.Add(limitToAdd).TruncateInt()
	k.SetRemainingOzoneLimit(ctx, newLimit)
	return limitToAdd.TruncateInt()
}

func (k Keeper) decreaseOzoneLimitBySubtractStake(ctx sdk.Context, stake sdk.Int) (ozoneLimitChange sdk.Int) {
	initialGenesisDeposit := k.GetInitialGenesisStakeTotal(ctx).ToDec() //ustos
	if initialGenesisDeposit.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("initialGenesisDeposit is zero, decrease ozone limit failed")
		return
	}
	currentLimit := k.GetRemainingOzoneLimit(ctx).ToDec() //uoz
	limitToSub := currentLimit.Mul(stake.ToDec()).Quo(initialGenesisDeposit)
	newLimit := currentLimit.Sub(limitToSub).TruncateInt()
	k.SetRemainingOzoneLimit(ctx, newLimit)
	return limitToSub.TruncateInt()
}

// GetResourceNetworksIterator gets an iterator over all network addresses
func (k Keeper) GetResourceNetworksIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.ResourceNodeKey)
}

// GetIndexingNetworksIterator gets an iterator over all network addresses
func (k Keeper) GetIndexingNetworksIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.IndexingNodeKey)
}

func (k Keeper) GetNetworks(ctx sdk.Context, keeper Keeper) (res []byte) {
	var networkList []string
	iterator := keeper.GetResourceNetworksIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		resourceNode := types.MustUnmarshalResourceNode(k.cdc, iterator.Value())
		networkList = append(networkList, resourceNode.NetworkID)
	}
	iter := keeper.GetIndexingNetworksIterator(ctx)
	for ; iter.Valid(); iter.Next() {
		indexingNode := types.MustUnmarshalResourceNode(k.cdc, iter.Value())
		networkList = append(networkList, indexingNode.NetworkID)
	}
	r := removeDuplicateValues(networkList)
	return r
}

func removeDuplicateValues(stringSlice []string) (res []byte) {
	keys := make(map[string]bool)
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			res = append(res, types.ModuleCdc.MustMarshalJSON(entry)...)
			res = append(res, ';')
		}
	}
	return res[:len(res)-1]
}
