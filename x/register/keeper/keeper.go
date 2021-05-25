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
		paramstore:            paramstore.WithKeyTable(types.ParamKeyTable()),
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

// GetLastResourceNodeTotalStake Load the last total stake of resource nodes.
func (k Keeper) GetLastResourceNodeTotalStake(ctx sdk.Context) (stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastResourceNodeTotalStakeKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &stake)
	return
}

// SetLastResourceNodeTotalStake Set the last total stake of resource nodes.
func (k Keeper) SetLastResourceNodeTotalStake(ctx sdk.Context, stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(stake)
	store.Set(types.LastResourceNodeTotalStakeKey, b)
}

// GetLastIndexingNodeTotalStake Load the last total stake of indexing nodes.
func (k Keeper) GetLastIndexingNodeTotalStake(ctx sdk.Context) (stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastIndexingNodeTotalStakeKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &stake)
	return
}

// SetLastIndexingNodeTotalStake Set the last total stake of indexing nodes.
func (k Keeper) SetLastIndexingNodeTotalStake(ctx sdk.Context, stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(stake)
	store.Set(types.LastIndexingNodeTotalStakeKey, b)
}
