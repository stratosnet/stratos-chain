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
	"time"
)

// Keeper of the register store
type Keeper struct {
	storeKey              sdk.StoreKey
	cdc                   *codec.Codec
	paramSpace            params.Subspace
	accountKeeper         auth.AccountKeeper
	bankKeeper            bank.Keeper
	hooks                 types.RegisterHooks
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
		hooks:                 nil,
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

// Set the validator hooks
func (k *Keeper) SetHooks(sh types.RegisterHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set register hooks twice")
	}
	k.hooks = sh
	return k
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

func (k Keeper) increaseOzoneLimitByAddStake(ctx sdk.Context, stake sdk.Int) {
	initialGenesisDeposit := k.GetInitialGenesisStakeTotal(ctx).ToDec() //ustos
	if initialGenesisDeposit.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("initialGenesisDeposit is zero, increase ozone limit failed")
		return
	}
	currentLimit := k.GetRemainingOzoneLimit(ctx).ToDec() //uoz
	limitToAdd := currentLimit.Mul(stake.ToDec()).Quo(initialGenesisDeposit)
	newLimit := currentLimit.Add(limitToAdd).TruncateInt()
	k.SetRemainingOzoneLimit(ctx, newLimit)
}

func (k Keeper) decreaseOzoneLimitBySubtractStake(ctx sdk.Context, stake sdk.Int) {
	initialGenesisDeposit := k.GetInitialGenesisStakeTotal(ctx).ToDec() //ustos
	if initialGenesisDeposit.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("initialGenesisDeposit is zero, decrease ozone limit failed")
		return
	}
	currentLimit := k.GetRemainingOzoneLimit(ctx).ToDec() //uoz
	limitToSub := currentLimit.Mul(stake.ToDec()).Quo(initialGenesisDeposit)
	newLimit := currentLimit.Sub(limitToSub).TruncateInt()
	k.SetRemainingOzoneLimit(ctx, newLimit)
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

//
//// return a given amount of all the UnbondingIndexingNodes
//func (k Keeper) GetUINs(ctx sdk.Context, networkAddr sdk.AccAddress,
//	maxRetrieve uint16) (unbondingIndexingNodes []types.UnbondingIndexingNode) {
//
//	unbondingIndexingNodes = make([]types.UnbondingIndexingNode, maxRetrieve)
//
//	store := ctx.KVStore(k.storeKey)
//	indexingNodePrefixKey := types.GetUINKey(networkAddr)
//	iterator := sdk.KVStorePrefixIterator(store, indexingNodePrefixKey)
//	defer iterator.Close()
//
//	i := 0
//	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
//		unbondingIndexingNode := types.MustUnmarshalUIN(k.cdc, iterator.Value())
//		unbondingIndexingNodes[i] = unbondingIndexingNode
//		i++
//	}
//	return unbondingIndexingNodes[:i] // trim if the array length < maxRetrieve
//}
//
//// return a unbonding UnbondingIndexingNode
//func (k Keeper) GetUIN(ctx sdk.Context,
//	networkAddr sdk.AccAddress) (ubd types.UnbondingIndexingNode, found bool) {
//
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetUINKey(networkAddr)
//	value := store.Get(key)
//	if value == nil {
//		return ubd, false
//	}
//
//	ubd = types.MustUnmarshalUIN(k.cdc, value)
//	return ubd, true
//}
//
//// iterate through all of the unbonding indexingNodes
//func (k Keeper) IterateUINs(ctx sdk.Context, fn func(index int64, ubd types.UnbondingIndexingNode) (stop bool)) {
//	store := ctx.KVStore(k.storeKey)
//	iterator := sdk.KVStorePrefixIterator(store, types.UBDIndexingNodeKey)
//	defer iterator.Close()
//
//	for i := int64(0); iterator.Valid(); iterator.Next() {
//		ubd := types.MustUnmarshalUIN(k.cdc, iterator.Value())
//		if stop := fn(i, ubd); stop {
//			break
//		}
//		i++
//	}
//}
//
//// HasMaxUnbondingIndexingNodeEntries - check if unbonding IndexingNode has maximum number of entries
//func (k Keeper) HasMaxUINEntries(ctx sdk.Context, networkAddr sdk.AccAddress) bool {
//	ubd, found := k.GetUIN(ctx, networkAddr)
//	if !found {
//		return false
//	}
//	return len(ubd.Entries) >= int(k.MaxEntries(ctx))
//}
//
//// set the unbonding IndexingNode
//func (k Keeper) SetUIN(ctx sdk.Context, ubd types.UnbondingIndexingNode) {
//	store := ctx.KVStore(k.storeKey)
//	bz := types.MustMarshalUIN(k.cdc, ubd)
//	key := types.GetUINKey(ubd.GetNetworkAddr())
//	store.Set(key, bz)
//}
//
//// remove the unbonding IndexingNode object
//func (k Keeper) RemoveUIN(ctx sdk.Context, ubd types.UnbondingIndexingNode) {
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetUINKey(ubd.GetNetworkAddr())
//	store.Delete(key)
//}
//
//// SetUnbondingIndexingNodeEntry adds an entry to the unbonding IndexingNode at
//// the given addresses. It creates the unbonding IndexingNode if it does not exist
//func (k Keeper) SetUINEntry(ctx sdk.Context, networkAddr sdk.AccAddress,
//	creationHeight int64, minTime time.Time, balance sdk.Int) types.UnbondingIndexingNode {
//
//	ubd, found := k.GetUIN(ctx, networkAddr)
//	if found {
//		ubd.AddEntry(creationHeight, minTime, balance)
//	} else {
//		ubd = types.NewUnbondingIndexingNode(networkAddr, creationHeight, minTime, balance)
//	}
//	k.SetUIN(ctx, ubd)
//	return ubd
//}
//
//// unbonding delegation queue timeslice operations
//
//// gets a specific unbonding queue timeslice. A timeslice is a slice of DVPairs
//// corresponding to unbonding delegations that expire at a certain time.
//func (k Keeper) GetUINQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (networkAddrs []sdk.AccAddress) {
//	store := ctx.KVStore(k.storeKey)
//	bz := store.Get(types.GetUINTimeKey(timestamp))
//	if bz == nil {
//		return []sdk.AccAddress{}
//	}
//	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &networkAddrs)
//	return networkAddrs
//}
//
//// Sets a specific unbonding queue timeslice.
//func (k Keeper) SetUINQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []sdk.AccAddress) {
//	store := ctx.KVStore(k.storeKey)
//	bz := k.cdc.MustMarshalBinaryLengthPrefixed(keys)
//	store.Set(types.GetUINTimeKey(timestamp), bz)
//}
//
//// Insert an unbonding delegation to the appropriate timeslice in the unbonding queue
//func (k Keeper) InsertUINQueue(ctx sdk.Context, ubd types.UnbondingIndexingNode,
//	completionTime time.Time) {
//
//	timeSlice := k.GetUINQueueTimeSlice(ctx, completionTime)
//	networkAddr := ubd.NetworkAddr
//	if len(timeSlice) == 0 {
//		k.SetUINQueueTimeSlice(ctx, completionTime, []sdk.AccAddress{networkAddr})
//	} else {
//		timeSlice = append(timeSlice, networkAddr)
//		k.SetUINQueueTimeSlice(ctx, completionTime, timeSlice)
//	}
//}
//
//// Returns all the unbonding queue timeslices from time 0 until endTime
//func (k Keeper) UINQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
//	store := ctx.KVStore(k.storeKey)
//	return store.Iterator(types.UBDIndexingNodeQueueKey,
//		sdk.InclusiveEndBytes(types.GetUINTimeKey(endTime)))
//}
//
//// Returns a concatenated list of all the timeslices inclusively previous to
//// currTime, and deletes the timeslices from the queue
//func (k Keeper) DequeueAllMatureUBDQueue(ctx sdk.Context,
//	currTime time.Time) (matureUnbonds []sdk.AccAddress) {
//
//	store := ctx.KVStore(k.storeKey)
//	// gets an iterator for all timeslices from time 0 until the current Blockheader time
//	unbondingTimesliceIterator := k.UINQueueIterator(ctx, ctx.BlockHeader().Time)
//	defer unbondingTimesliceIterator.Close()
//
//	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
//		timeslice := []sdk.AccAddress{}
//		value := unbondingTimesliceIterator.Value()
//		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)
//		matureUnbonds = append(matureUnbonds, timeslice...)
//		store.Delete(unbondingTimesliceIterator.Key())
//	}
//	return matureUnbonds
//}

//==========

// return a given amount of all the UnbondingIndexingNodes
func (k Keeper) GetUnbondingNodes(ctx sdk.Context, networkAddr sdk.AccAddress,
	maxRetrieve uint16) (unbondingIndexingNodes []types.UnbondingNode) {

	unbondingIndexingNodes = make([]types.UnbondingNode, maxRetrieve)

	store := ctx.KVStore(k.storeKey)
	indexingNodePrefixKey := types.GetUBDNodeKey(networkAddr)
	iterator := sdk.KVStorePrefixIterator(store, indexingNodePrefixKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		unbondingIndexingNode := types.MustUnmarshalUnbondingNode(k.cdc, iterator.Value())
		unbondingIndexingNodes[i] = unbondingIndexingNode
		i++
	}
	return unbondingIndexingNodes[:i] // trim if the array length < maxRetrieve
}

// return a unbonding UnbondingIndexingNode
func (k Keeper) GetUnbondingNode(ctx sdk.Context,
	networkAddr sdk.AccAddress) (ubd types.UnbondingNode, found bool) {

	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDNodeKey(networkAddr)
	value := store.Get(key)
	if value == nil {
		return ubd, false
	}

	ubd = types.MustUnmarshalUnbondingNode(k.cdc, value)
	return ubd, true
}

// iterate through all of the unbonding indexingNodes
func (k Keeper) IterateUnbondingNodes(ctx sdk.Context, fn func(index int64, ubd types.UnbondingNode) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UBDNodeKey)
	defer iterator.Close()

	for i := int64(0); iterator.Valid(); iterator.Next() {
		ubd := types.MustUnmarshalUnbondingNode(k.cdc, iterator.Value())
		if stop := fn(i, ubd); stop {
			break
		}
		i++
	}
}

// HasMaxUnbondingIndexingNodeEntries - check if unbonding IndexingNode has maximum number of entries
func (k Keeper) HasMaxUnbondingNodeEntries(ctx sdk.Context, networkAddr sdk.AccAddress) bool {
	ubd, found := k.GetUnbondingNode(ctx, networkAddr)
	if !found {
		return false
	}
	return len(ubd.Entries) >= int(k.MaxEntries(ctx))
}

// set the unbonding IndexingNode
func (k Keeper) SetUnbondingNode(ctx sdk.Context, ubd types.UnbondingNode) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalUnbondingNode(k.cdc, ubd)
	key := types.GetUBDNodeKey(ubd.GetNetworkAddr())
	store.Set(key, bz)
}

// remove the unbonding IndexingNode object
func (k Keeper) RemoveUnbondingNode(ctx sdk.Context, ubd types.UnbondingNode) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDNodeKey(ubd.GetNetworkAddr())
	store.Delete(key)
}

// SetUnbondingIndexingNodeEntry adds an entry to the unbonding IndexingNode at
// the given addresses. It creates the unbonding IndexingNode if it does not exist
func (k Keeper) SetUnbondingNodeEntry(ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool,
	creationHeight int64, minTime time.Time, balance sdk.Int) types.UnbondingNode {

	ubd, found := k.GetUnbondingNode(ctx, networkAddr)
	if found {
		ubd.AddEntry(creationHeight, minTime, balance)
	} else {
		ubd = types.NewUnbondingNode(networkAddr, isIndexingNode, creationHeight, minTime, balance)
	}
	k.SetUnbondingNode(ctx, ubd)
	return ubd
}

// unbonding delegation queue timeslice operations

// gets a specific unbonding queue timeslice. A timeslice is a slice of DVPairs
// corresponding to unbonding delegations that expire at a certain time.
func (k Keeper) GetUnbondingNodeQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (networkAddrs []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetUBDTimeKey(timestamp))
	if bz == nil {
		return []sdk.AccAddress{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &networkAddrs)
	return networkAddrs
}

// Sets a specific unbonding queue timeslice.
func (k Keeper) SetUnbondingNodeQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(keys)
	store.Set(types.GetUBDTimeKey(timestamp), bz)
}

// Insert an unbonding delegation to the appropriate timeslice in the unbonding queue
func (k Keeper) InsertUnbondingNodeQueue(ctx sdk.Context, ubd types.UnbondingNode,
	completionTime time.Time) {

	timeSlice := k.GetUnbondingNodeQueueTimeSlice(ctx, completionTime)
	networkAddr := ubd.NetworkAddr
	if len(timeSlice) == 0 {
		k.SetUnbondingNodeQueueTimeSlice(ctx, completionTime, []sdk.AccAddress{networkAddr})
	} else {
		timeSlice = append(timeSlice, networkAddr)
		k.SetUnbondingNodeQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// Returns all the unbonding queue timeslices from time 0 until endTime
func (k Keeper) UnbondingNodeQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UBDNodeQueueKey,
		sdk.InclusiveEndBytes(types.GetUBDTimeKey(endTime)))
}

// Returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue
func (k Keeper) DequeueAllMatureUBDQueue(ctx sdk.Context,
	currTime time.Time) (matureUnbonds []sdk.AccAddress) {

	store := ctx.KVStore(k.storeKey)
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	unbondingTimesliceIterator := k.UnbondingNodeQueueIterator(ctx, ctx.BlockHeader().Time)
	defer unbondingTimesliceIterator.Close()

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := []sdk.AccAddress{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)
		matureUnbonds = append(matureUnbonds, timeslice...)
		store.Delete(unbondingTimesliceIterator.Key())
	}
	ctx.Logger().Info(fmt.Sprintf("DequeueAllMatureUBDQueue, %d matured unbonding nodes detected", len(matureUnbonds)))
	return matureUnbonds
}

// CompleteUnbondingWithAmount completes the unbonding of all mature entries in
// the retrieved unbonding delegation object and returns the total unbonding
// balance or an error upon failure.
func (k Keeper) CompleteUnbondingWithAmount(ctx sdk.Context, networkAddr sdk.AccAddress) (sdk.Coins, error) {
	ubd, found := k.GetUnbondingNode(ctx, networkAddr)
	if !found {
		ctx.Logger().Info(fmt.Sprintf("NetworAddr: %s not found while completing UnbondingWithAmount", networkAddr))
		return nil, types.ErrNoUnbondingNode
	}

	bondDenom := k.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time
	ctx.Logger().Info(fmt.Sprintf("Completing UnbondingWithAmount, networAddr: %s", networkAddr))
	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(ctxTime) {
			ubd.RemoveEntry(int64(i))
			i--

			// track undelegation only when remaining or truncated shares are non-zero
			if !entry.Balance.IsZero() {
				amt := sdk.NewCoin(bondDenom, entry.Balance)
				err := k.SubtractUBDNodeStake(ctx, ubd, amt)
				if err != nil {
					return nil, err
				}

				balances = balances.Add(amt)
			}
		}
	}

	// set the unbonding node or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveUnbondingNode(ctx, ubd)
	} else {
		k.SetUnbondingNode(ctx, ubd)
	}

	return balances, nil
}

// CompleteUnbonding performs the same logic as CompleteUnbondingWithAmount except
// it does not return the total unbonding amount.
func (k Keeper) CompleteUnbonding(ctx sdk.Context, networkAddr sdk.AccAddress) error {
	_, err := k.CompleteUnbondingWithAmount(ctx, networkAddr)
	return err
}

func (k Keeper) SubtractUBDNodeStake(ctx sdk.Context, ubd types.UnbondingNode, tokenToSub sdk.Coin) error {
	// case of indexing node
	if ubd.IsIndexingNode {
		indexingNode, found := k.GetIndexingNode(ctx, ubd.NetworkAddr)
		if !found {
			return types.ErrNoIndexingNodeFound
		}
		return k.SubtractIndexingNodeStake(ctx, indexingNode, tokenToSub)
	}
	// case of resource node
	resourceNode, found := k.GetResourceNode(ctx, ubd.NetworkAddr)
	if !found {
		return types.ErrNoIndexingNodeFound
	}
	return k.SubtractResourceNodeStake(ctx, resourceNode, tokenToSub)
}

func (k Keeper) DoRemoveResourceNode(
	ctx sdk.Context, resourceNode types.ResourceNode, amt sdk.Int,
) (time.Time, error) {
	// transfer the node tokens to the not bonded pool
	params := k.GetParams(ctx)
	ctx.Logger().Info("Params of register module: " + params.String())

	ownerAcc := k.accountKeeper.GetAccount(ctx, resourceNode.OwnerAddress)
	if ownerAcc == nil {
		return time.Time{}, types.ErrNoOwnerAccountFound
	}

	networkAddr := resourceNode.GetNetworkAddr()
	if k.HasMaxUnbondingNodeEntries(ctx, networkAddr) {
		return time.Time{}, types.ErrMaxUnbondingNodeEntries
	}

	returnAmount, err := k.unbond(ctx, networkAddr, false, amt)
	if err != nil {
		return time.Time{}, err
	}

	if resourceNode.GetStatus() == sdk.Bonded {
		k.bondedToUnbonding(ctx, resourceNode, false)
	}

	// set the unbonding mature time and completion height appropriately
	unbondingMatureTime := calcUnbondingMatureTime(resourceNode.CreationTime, params.UnbondingThreasholdTime, params.UnbondingCompletionTime)
	ctx.Logger().Info(fmt.Sprintf("Calculating mature time: creationTime[%s], threasholdTime[%s], completionTime[%s], matureTime[%s]",
		resourceNode.CreationTime.String(), params.UnbondingThreasholdTime.String(), params.UnbondingCompletionTime.String(), unbondingMatureTime.String(),
	))
	//unbondingNode := types.NewUnbondingNode(resourceNode.GetNetworkAddr(), false, ctx.BlockHeight(), unbondingMatureTime, returnAmount)
	unbondingNode := k.SetUnbondingNodeEntry(ctx, resourceNode.GetNetworkAddr(), false, ctx.BlockHeight(), unbondingMatureTime, returnAmount)
	// Adds to unbonding node queue
	k.InsertUnbondingNodeQueue(ctx, unbondingNode, unbondingMatureTime)
	ctx.Logger().Info("Unbonding resource node " + unbondingNode.String() + "\n after mature time" + unbondingMatureTime.String())

	return unbondingMatureTime, nil
}

func (k Keeper) DoRemoveIndexingNode(
	ctx sdk.Context, indexingNode types.IndexingNode, amt sdk.Int,
) (time.Time, error) {

	ownerAcc := k.accountKeeper.GetAccount(ctx, indexingNode.OwnerAddress)
	if ownerAcc == nil {
		return time.Time{}, types.ErrNoOwnerAccountFound
	}

	networkAddr := indexingNode.GetNetworkAddr()
	if k.HasMaxUnbondingNodeEntries(ctx, networkAddr) {
		return time.Time{}, types.ErrMaxUnbondingNodeEntries
	}

	returnAmount, err := k.unbond(ctx, networkAddr, true, amt)
	if err != nil {
		return time.Time{}, err
	}

	// transfer the node tokens to the not bonded pool
	if indexingNode.GetStatus() == sdk.Bonded {
		k.bondedToUnbonding(ctx, indexingNode, true)
	}

	params := k.GetParams(ctx)
	// set the unbonding mature time and completion height appropriately
	unbondingMatureTime := calcUnbondingMatureTime(indexingNode.CreationTime, params.UnbondingThreasholdTime, params.UnbondingCompletionTime)
	//unbondingNode := types.NewUnbondingNode(indexingNode.GetNetworkAddr(), true, ctx.BlockHeight(), unbondingMatureTime, returnAmount)
	unbondingNode := k.SetUnbondingNodeEntry(ctx, indexingNode.GetNetworkAddr(), true, ctx.BlockHeight(), unbondingMatureTime, returnAmount)
	// Adds to unbonding node queue
	k.InsertUnbondingNodeQueue(ctx, unbondingNode, unbondingMatureTime)
	ctx.Logger().Info("Unbonding indexing node " + unbondingNode.String() + "\n after mature time" + unbondingMatureTime.String())
	return unbondingMatureTime, nil
}

// unbond a particular node and perform associated store operations
func (k Keeper) unbond(
	ctx sdk.Context, networkAddr sdk.AccAddress, isIndexingNode bool, amt sdk.Int,
) (amount sdk.Int, err error) {

	//// check if a node object exists in the store
	//if isIndexingNode {
	//	node2Unbond, found := k.GetIndexingNode(ctx, networkAddr)
	//	if !found {
	//		return amount, types.ErrNoNodeForAddress
	//	}
	//} else {
	//	node2Unbond, found := k.GetResourceNode(ctx, networkAddr)
	//	if !found {
	//		return amount, types.ErrNoNodeForAddress
	//	}
	//}
	//
	//// call the before-delegation-modified hook
	//k.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	//
	//// ensure that we have enough shares to remove
	//if delegation.Shares.LT(shares) {
	//	return amount, sdkerrors.Wrap(types.ErrNotEnoughDelegationShares, delegation.Shares.String())
	//}
	//
	//// get validator
	//validator, found := k.GetValidator(ctx, valAddr)
	//if !found {
	//	return amount, types.ErrNoValidatorFound
	//}
	//
	//// subtract shares from delegation
	//delegation.Shares = delegation.Shares.Sub(shares)
	//
	//isValidatorOperator := delegation.DelegatorAddress.Equals(validator.OperatorAddress)
	//
	//// if the delegation is the operator of the validator and undelegating will decrease the validator's self delegation below their minimum
	//// trigger a jail validator
	//if isValidatorOperator && !validator.Jailed &&
	//	validator.TokensFromShares(delegation.Shares).TruncateInt().LT(validator.MinSelfDelegation) {
	//
	//	k.jailValidator(ctx, validator)
	//	validator = k.mustGetValidator(ctx, validator.OperatorAddress)
	//}
	//
	//// remove the delegation
	//if delegation.Shares.IsZero() {
	//	k.RemoveDelegation(ctx, delegation)
	//} else {
	//	k.SetDelegation(ctx, delegation)
	//	// call the after delegation modification hook
	//	k.AfterDelegationModified(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	//}
	//
	//// remove the shares and coins from the validator
	//// NOTE that the amount is later (in keeper.Delegation) moved between staking module pools
	//validator, amount = k.RemoveValidatorTokensAndShares(ctx, validator, shares)
	//
	//if validator.DelegatorShares.IsZero() && validator.IsUnbonded() {
	//	// if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
	//	k.RemoveValidator(ctx, validator.OperatorAddress)
	//}

	return amount, nil
}
