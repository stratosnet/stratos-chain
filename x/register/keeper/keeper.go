package keeper

import (
	"container/list"
	"encoding/json"
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
	accountKeeper         auth.AccountKeeper
	bankKeeper            bank.Keeper
	paramstore            params.Subspace
	resourceNodeCache     map[string]cachedResourceNode
	resourceNodeCacheList *list.List
	indexingNodeCache     map[string]cachedIndexingNode
	indexingNodeCacheList *list.List
}

// NewKeeper creates a register keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, paramstore params.Subspace) Keeper {
	keeper := Keeper{
		storeKey:              key,
		cdc:                   cdc,
		accountKeeper:         accountKeeper,
		bankKeeper:            bankKeeper,
		paramstore:            paramstore.WithKeyTable(ParamKeyTable()),
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

// GetLastResourceNodeTotalPower Load the last total power of resource nodes.
func (k Keeper) GetLastResourceNodeTotalPower(ctx sdk.Context) (power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastResourceNodeTotalPowerKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &power)
	return
}

// SetLastResourceNodeTotalPower Set the last total power of resource nodes.
func (k Keeper) SetLastResourceNodeTotalPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(power)
	store.Set(types.LastResourceNodeTotalPowerKey, b)
}

// GetLastIndexingNodeTotalPower Load the last total power of indexing nodes.
func (k Keeper) GetLastIndexingNodeTotalPower(ctx sdk.Context) (power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastIndexingNodeTotalPowerKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &power)
	return
}

// SetLastIndexingNodeTotalPower Set the last total power of indexing nodes.
func (k Keeper) SetLastIndexingNodeTotalPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(power)
	store.Set(types.LastIndexingNodeTotalPowerKey, b)
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

func (k Keeper) GetNetworks(ctx sdk.Context, keeper Keeper) (res []byte, err error) {
	var networkList []string
	iterator := keeper.GetResourceNetworksIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		resourceNode := types.MustUnmarshalResourceNode(k.cdc, iterator.Value())
		networkList = append(networkList, resourceNode.NetworkAddress)
	}
	iter := keeper.GetIndexingNetworksIterator(ctx)
	for ; iter.Valid(); iter.Next() {
		indexingNode := types.MustUnmarshalResourceNode(k.cdc, iter.Value())
		networkList = append(networkList, indexingNode.NetworkAddress)
	}

	r := removeDuplicateValues(networkList)
	ctx.Logger().Info("r: ", r)
	bz, err2 := json.Marshal(r)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}
	ctx.Logger().Info("bz: ", bz)
	return bz, nil
}

func removeDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	var res []string

	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			res = append(res, entry)
		}
	}
	return res
}
