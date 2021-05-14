package keeper

import (
	"container/list"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
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
