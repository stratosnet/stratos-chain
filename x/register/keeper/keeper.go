package keeper

import (
	"container/list"
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the register store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.BinaryCodec
	// module specific parameter space that can be configured through governance
	paramSpace            paramtypes.Subspace
	accountKeeper         types.AccountKeeper
	bankKeeper            types.BankKeeper
	hooks                 types.RegisterHooks
	resourceNodeCache     map[string]cachedResourceNode
	resourceNodeCacheList *list.List
	indexingNodeCache     map[string]cachedIndexingNode
	indexingNodeCacheList *list.List
}

// NewKeeper creates a register keeper
func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper) Keeper {

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

// SetHooks Set the register hooks
func (k *Keeper) SetHooks(sh types.RegisterHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set register hooks twice")
	}
	k.hooks = sh
	return k
}

func (k Keeper) SetInitialUOzonePrice(ctx sdk.Context, price sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := amino.MustMarshalBinaryLengthPrefixed(price)
	store.Set(types.InitialUOzonePriceKey, b)
}

func (k Keeper) GetInitialUOzonePrice(ctx sdk.Context) (price sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.InitialUOzonePriceKey)
	if b == nil {
		panic("Stored initial uOzone price should not have been nil")
	}
	amino.MustUnmarshalBinaryLengthPrefixed(b, &price)
	return
}

func (k Keeper) SetTotalUnissuedPrepay(ctx sdk.Context, totalUnissuedPrepay sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	b := amino.MustMarshalBinaryLengthPrefixed(totalUnissuedPrepay)
	store.Set(types.TotalUnissuedPrepayKey, b)
}

func (k Keeper) GetTotalUnissuedPrepay(ctx sdk.Context) (totalUnissuedPrepay sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.TotalUnissuedPrepayKey)
	if b == nil {
		return sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())
	}
	amino.MustUnmarshalBinaryLengthPrefixed(b, &totalUnissuedPrepay)
	return
}

func (k Keeper) SetInitialGenesisStakeTotal(ctx sdk.Context, stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := amino.MustMarshalBinaryLengthPrefixed(stake)
	store.Set(types.InitialGenesisStakeTotalKey, b)
}

func (k Keeper) GetInitialGenesisStakeTotal(ctx sdk.Context) (stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.InitialGenesisStakeTotalKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	amino.MustUnmarshalBinaryLengthPrefixed(b, &stake)
	return
}

func (k Keeper) SetRemainingOzoneLimit(ctx sdk.Context, value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := amino.MustMarshalBinaryLengthPrefixed(value)
	store.Set(types.UpperBoundOfTotalOzoneKey, b)
}

func (k Keeper) GetRemainingOzoneLimit(ctx sdk.Context) (value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.UpperBoundOfTotalOzoneKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	amino.MustUnmarshalBinaryLengthPrefixed(b, &value)
	return
}

func (k Keeper) increaseOzoneLimitByAddStake(ctx sdk.Context, stake sdk.Int) (ozoneLimitChange sdk.Int) {
	initialGenesisDeposit := k.GetInitialGenesisStakeTotal(ctx).ToDec() //ustos
	if initialGenesisDeposit.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("initialGenesisDeposit is zero, increase ozone limit failed")
		return sdk.ZeroInt()
	}
	initialUozonePrice := k.GetInitialUOzonePrice(ctx)
	if initialUozonePrice.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("initialUozonePrice is zero, increase ozone limit failed")
		return sdk.ZeroInt()
	}
	initialOzoneLimit := initialGenesisDeposit.Quo(initialUozonePrice)
	//ctx.Logger().Debug("----- initialOzoneLimit is " + initialOzoneLimit.String() + " uoz", )
	currentLimit := k.GetRemainingOzoneLimit(ctx).ToDec() //uoz
	//ctx.Logger().Info("----- currentLimit is " + currentLimit.String() + " uoz")
	limitToAdd := initialOzoneLimit.Mul(stake.ToDec()).Quo(initialGenesisDeposit)
	//ctx.Logger().Info("----- limitToAdd is " + limitToAdd.String() + " uoz")
	newLimit := currentLimit.Add(limitToAdd).TruncateInt()
	//ctx.Logger().Info("----- newLimit is " + newLimit.String() + " uoz")
	k.SetRemainingOzoneLimit(ctx, newLimit)
	return limitToAdd.TruncateInt()
}

func (k Keeper) decreaseOzoneLimitBySubtractStake(ctx sdk.Context, stake sdk.Int) (ozoneLimitChange sdk.Int) {
	initialGenesisDeposit := k.GetInitialGenesisStakeTotal(ctx).ToDec() //ustos
	if initialGenesisDeposit.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("initialGenesisDeposit is zero, decrease ozone limit failed")
		return sdk.ZeroInt()
	}
	initialUozonePrice := k.GetInitialUOzonePrice(ctx)
	if initialUozonePrice.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("initialUozonePrice is zero, increase ozone limit failed")
		return sdk.ZeroInt()
	}
	initialOzoneLimit := initialGenesisDeposit.Quo(initialUozonePrice)
	currentLimit := k.GetRemainingOzoneLimit(ctx).ToDec() //uoz
	limitToSub := initialOzoneLimit.Mul(stake.ToDec()).Quo(initialGenesisDeposit)
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
	var networkList []stratos.SdsAddress
	iterator := keeper.GetResourceNetworksIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		resourceNode := types.MustUnmarshalResourceNode(k.cdc, iterator.Value())
		networkAddr, err := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddr())
		if err != nil {
			continue
		}
		networkList = append(networkList, networkAddr)
	}
	iter := keeper.GetIndexingNetworksIterator(ctx)
	for ; iter.Valid(); iter.Next() {
		indexingNode := types.MustUnmarshalResourceNode(k.cdc, iter.Value())
		networkAddr, err := stratos.SdsAddressFromBech32(indexingNode.GetNetworkAddr())
		if err != nil {
			continue
		}
		networkList = append(networkList, networkAddr)
	}
	r := removeDuplicateValues(networkList)
	return r
}

func removeDuplicateValues(stringSlice []stratos.SdsAddress) (res []byte) {
	keys := make(map[string]bool)
	for _, entry := range stringSlice {
		if _, value := keys[entry.String()]; !value {
			keys[entry.String()] = true
			res = append(res, types.ModuleCdc.MustMarshalJSON(entry)...)
			res = append(res, ';')
		}
	}
	return res[:len(res)-1]
}

// return a given amount of all the UnbondingIndexingNodes
func (k Keeper) GetUnbondingNodes(ctx sdk.Context, networkAddr stratos.SdsAddress,
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
	networkAddr stratos.SdsAddress) (ubd types.UnbondingNode, found bool) {

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
func (k Keeper) HasMaxUnbondingNodeEntries(ctx sdk.Context, networkAddr stratos.SdsAddress) bool {
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
func (k Keeper) SetUnbondingNodeEntry(ctx sdk.Context, networkAddr stratos.SdsAddress, isIndexingNode bool,
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
func (k Keeper) GetUnbondingNodeQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (networkAddrs []stratos.SdsAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetUBDTimeKey(timestamp))
	if bz == nil {
		return []stratos.SdsAddress{}
	}
	amino.MustUnmarshalBinaryLengthPrefixed(bz, &networkAddrs)
	return networkAddrs
}

// Sets a specific unbonding queue timeslice.
func (k Keeper) SetUnbondingNodeQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []stratos.SdsAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := amino.MustMarshalBinaryLengthPrefixed(keys)
	store.Set(types.GetUBDTimeKey(timestamp), bz)
}

// Insert an unbonding delegation to the appropriate timeslice in the unbonding queue
func (k Keeper) InsertUnbondingNodeQueue(ctx sdk.Context, ubd types.UnbondingNode,
	completionTime time.Time) {

	timeSlice := k.GetUnbondingNodeQueueTimeSlice(ctx, completionTime)
	networkAddr := ubd.NetworkAddr
	if len(timeSlice) == 0 {
		k.SetUnbondingNodeQueueTimeSlice(ctx, completionTime, []stratos.SdsAddress{networkAddr})
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
	currTime time.Time) (matureUnbonds []stratos.SdsAddress) {

	store := ctx.KVStore(k.storeKey)
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	unbondingTimesliceIterator := k.UnbondingNodeQueueIterator(ctx, ctx.BlockHeader().Time)
	defer unbondingTimesliceIterator.Close()

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := []stratos.SdsAddress{}
		value := unbondingTimesliceIterator.Value()
		amino.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)
		matureUnbonds = append(matureUnbonds, timeslice...)
		store.Delete(unbondingTimesliceIterator.Key())
	}
	ctx.Logger().Debug(fmt.Sprintf("DequeueAllMatureUBDQueue, %d matured unbonding nodes detected", len(matureUnbonds)))
	return matureUnbonds
}

// CompleteUnbondingWithAmount completes the unbonding of all mature entries in
// the retrieved unbonding delegation object and returns the total unbonding
// balance or an error upon failure.
func (k Keeper) CompleteUnbondingWithAmount(ctx sdk.Context, networkAddr stratos.SdsAddress) (sdk.Coins, bool, error) {
	ubd, found := k.GetUnbondingNode(ctx, networkAddr)
	if !found {
		ctx.Logger().Info(fmt.Sprintf("NetworAddr: %s not found while completing UnbondingWithAmount", networkAddr))
		return nil, false, types.ErrNoUnbondingNode
	}

	bondDenom := k.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time
	ctx.Logger().Debug(fmt.Sprintf("Completing UnbondingWithAmount, networAddr: %s", networkAddr))
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
					return nil, false, err
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

	return balances, ubd.IsIndexingNode, nil
}

// CompleteUnbonding performs the same logic as CompleteUnbondingWithAmount except
// it does not return the total unbonding amount.
func (k Keeper) CompleteUnbonding(ctx sdk.Context, networkAddr stratos.SdsAddress) error {
	_, _, err := k.CompleteUnbondingWithAmount(ctx, networkAddr)
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

func (k Keeper) UnbondResourceNode(
	ctx sdk.Context, resourceNode types.ResourceNode, amt sdk.Int,
) (ozoneLimitChange sdk.Int, unbondingMatureTime time.Time, err error) {
	params := k.GetParams(ctx)
	ctx.Logger().Info("Params of register module: " + params.String())

	// transfer the node tokens to the not bonded pool
	networkAddr, err := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddr())
	if err != nil {
		return sdk.ZeroInt(), time.Now(), errors.New("invalid network address")
	}
	ownerAddr, err := sdk.AccAddressFromBech32(resourceNode.GetOwnerAddress())
	if err != nil {
		return sdk.ZeroInt(), time.Now(), errors.New("invalid wallet address")
	}
	ownerAcc := k.accountKeeper.GetAccount(ctx, ownerAddr)
	if ownerAcc == nil {
		return sdk.ZeroInt(), time.Time{}, types.ErrNoOwnerAccountFound
	}

	if k.HasMaxUnbondingNodeEntries(ctx, networkAddr) {
		return sdk.ZeroInt(), time.Time{}, types.ErrMaxUnbondingNodeEntries
	}
	unbondingMatureTime = calcUnbondingMatureTime(ctx, resourceNode.Status, resourceNode.CreationTime, k.UnbondingThreasholdTime(ctx), k.UnbondingCompletionTime(ctx))

	bondDenom := k.GetParams(ctx).BondDenom
	coin := sdk.NewCoin(bondDenom, amt)
	if resourceNode.GetStatus() == stakingtypes.Bonded {
		// transfer the node tokens to the not bonded pool
		k.bondedToUnbonding(ctx, resourceNode, false, coin)
		// adjust ozone limit
		ozoneLimitChange = k.decreaseOzoneLimitBySubtractStake(ctx, amt)
	}

	// change node status to unbonding if unbonding all tokens
	if amt.Equal(resourceNode.Tokens) {
		resourceNode.Status = stakingtypes.Unbonding
		k.SetResourceNode(ctx, resourceNode)
	}

	// set the unbonding mature time and completion height appropriately
	ctx.Logger().Info(fmt.Sprintf("Calculating mature time: creationTime[%s], threasholdTime[%s], completionTime[%s], matureTime[%s]",
		resourceNode.CreationTime, k.UnbondingThreasholdTime(ctx), k.UnbondingCompletionTime(ctx), unbondingMatureTime,
	))
	unbondingNode := k.SetUnbondingNodeEntry(ctx, networkAddr, false, ctx.BlockHeight(), unbondingMatureTime, amt)
	// Add to unbonding node queue
	k.InsertUnbondingNodeQueue(ctx, unbondingNode, unbondingMatureTime)
	ctx.Logger().Info("Unbonding resource node " + unbondingNode.String() + "\n after mature time" + unbondingMatureTime.String())

	return ozoneLimitChange, unbondingMatureTime, nil
}

func (k Keeper) UnbondIndexingNode(
	ctx sdk.Context, indexingNode types.IndexingNode, amt sdk.Int,
) (ozoneLimitChange sdk.Int, unbondingMatureTime time.Time, err error) {

	networkAddr, err := stratos.SdsAddressFromBech32(indexingNode.GetNetworkAddr())
	if err != nil {
		return sdk.ZeroInt(), time.Now(), errors.New("invalid network address")
	}
	ownerAddr, err := sdk.AccAddressFromBech32(indexingNode.GetOwnerAddress())
	if err != nil {
		return sdk.ZeroInt(), time.Now(), errors.New("invalid wallet address")
	}

	ownerAcc := k.accountKeeper.GetAccount(ctx, ownerAddr)
	if ownerAcc == nil {
		return sdk.ZeroInt(), time.Time{}, types.ErrNoOwnerAccountFound
	}

	if k.HasMaxUnbondingNodeEntries(ctx, networkAddr) {
		return sdk.ZeroInt(), time.Time{}, types.ErrMaxUnbondingNodeEntries
	}

	unbondingMatureTime = calcUnbondingMatureTime(ctx, indexingNode.Status, indexingNode.CreationTime, k.UnbondingThreasholdTime(ctx), k.UnbondingCompletionTime(ctx))

	bondDenom := k.GetParams(ctx).BondDenom
	coin := sdk.NewCoin(bondDenom, amt)
	if indexingNode.GetStatus() == stakingtypes.Bonded {
		// transfer the node tokens to the not bonded pool
		k.bondedToUnbonding(ctx, indexingNode, true, coin)
		// adjust ozone limit
		ozoneLimitChange = k.decreaseOzoneLimitBySubtractStake(ctx, amt)
	}
	// change node status to unbonding if unbonding all tokens
	if amt.Equal(indexingNode.Tokens) {
		indexingNode.Status = stakingtypes.Unbonding
		k.SetIndexingNode(ctx, indexingNode)
	}

	// Set the unbonding mature time and completion height appropriately
	unbondingNode := k.SetUnbondingNodeEntry(ctx, networkAddr, true, ctx.BlockHeight(), unbondingMatureTime, amt)
	// Add to unbonding node queue
	k.InsertUnbondingNodeQueue(ctx, unbondingNode, unbondingMatureTime)
	ctx.Logger().Info("Unbonding indexing node " + unbondingNode.String() + "\n after mature time" + unbondingMatureTime.String())
	return ozoneLimitChange, unbondingMatureTime, nil
}

// GetAllUnbondingNodes get the set of all ubd nodes with no limits, used during genesis dump
func (k Keeper) GetAllUnbondingNodes(ctx sdk.Context) (unbondingNodes []types.UnbondingNode) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UBDNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalUnbondingNode(k.cdc, iterator.Value())
		unbondingNodes = append(unbondingNodes, node)
	}
	return unbondingNodes
}

func (k Keeper) GetAllUnbondingNodesTotalBalance(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UBDNodeKey)
	defer iterator.Close()

	var ubdTotal = sdk.ZeroInt()
	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalUnbondingNode(k.cdc, iterator.Value())
		for _, entry := range node.Entries {
			ubdTotal = ubdTotal.Add(entry.Balance)
		}
	}
	return ubdTotal
}

// GetUnbondingNodeBalance returns an unbonding balance and an UnbondingNode
func (k Keeper) GetUnbondingNodeBalance(ctx sdk.Context,
	networkAddr stratos.SdsAddress) sdk.Int {

	balance := sdk.ZeroInt()

	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDNodeKey(networkAddr)
	value := store.Get(key)
	if value == nil {
		return balance
	}

	ubd := types.MustUnmarshalUnbondingNode(k.cdc, value)
	for _, entry := range ubd.Entries {
		balance = balance.Add(entry.Balance)
	}
	return balance
}

// calc current uoz price
func (k Keeper) CurrUozPrice(ctx sdk.Context) sdk.Dec {
	S := k.GetInitialGenesisStakeTotal(ctx)
	Pt := k.GetTotalUnissuedPrepay(ctx).Amount
	Lt := k.GetRemainingOzoneLimit(ctx)
	currUozPrice := (S.Add(Pt)).ToDec().
		Quo(Lt.ToDec())
	return currUozPrice
}

// calc remaining/total supply for uoz
func (k Keeper) UozSupply(ctx sdk.Context) (remaining, total sdk.Int) {
	remaining = k.GetRemainingOzoneLimit(ctx) // Lt
	S := k.GetInitialGenesisStakeTotal(ctx)
	Pt := k.GetTotalUnissuedPrepay(ctx).Amount
	// total supply = Lt * ( 1 + Pt / S )
	total = (Pt.ToDec().Quo(S.ToDec()).TruncateInt().Add(sdk.NewInt(1))).Mul(remaining)
	return remaining, total
}
