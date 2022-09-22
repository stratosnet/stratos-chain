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
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the register store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.Codec
	// module specific parameter space that can be configured through governance
	paramSpace            paramtypes.Subspace
	accountKeeper         types.AccountKeeper
	bankKeeper            types.BankKeeper
	distrKeeper           types.DistrKeeper
	hooks                 types.RegisterHooks
	resourceNodeCache     map[string]cachedResourceNode
	resourceNodeCacheList *list.List
	metaNodeCache         map[string]cachedMetaNode
	metaNodeCacheList     *list.List
}

// NewKeeper creates a register keeper
func NewKeeper(cdc codec.Codec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, distrKeeper types.DistrKeeper) Keeper {

	keeper := Keeper{
		storeKey:              key,
		cdc:                   cdc,
		paramSpace:            paramSpace.WithKeyTable(types.ParamKeyTable()),
		accountKeeper:         accountKeeper,
		bankKeeper:            bankKeeper,
		distrKeeper:           distrKeeper,
		hooks:                 nil,
		resourceNodeCache:     make(map[string]cachedResourceNode, resourceNodeCacheSize),
		resourceNodeCacheList: list.New(),
		metaNodeCache:         make(map[string]cachedMetaNode, metaNodeCacheSize),
		metaNodeCacheList:     list.New(),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetHooks Set the register hooks
func (k Keeper) SetHooks(sh types.RegisterHooks) Keeper {
	if k.hooks != nil {
		panic("cannot set register hooks twice")
	}
	k.hooks = sh
	return k
}

func (k Keeper) SetInitialUOzonePrice(ctx sdk.Context, price sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := types.ModuleCdc.MustMarshalLengthPrefixed(price)
	store.Set(types.InitialUOzonePriceKey, b)
}

func (k Keeper) GetInitialUOzonePrice(ctx sdk.Context) (price sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.InitialUOzonePriceKey)
	if b == nil {
		panic("Stored initial uOzone price should not have been nil")
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(b, &price)
	return
}

func (k Keeper) SendCoinsFromAccount2TotalUnissuedPrepayPool(ctx sdk.Context, fromWallet sdk.AccAddress, coinToSend sdk.Coin) error {
	fromAcc := k.accountKeeper.GetAccount(ctx, fromWallet)
	if fromAcc == nil {
		return types.ErrUnknownAccountAddress
	}
	hasCoin := k.bankKeeper.HasBalance(ctx, fromWallet, coinToSend)
	if !hasCoin {
		return types.ErrInsufficientBalance
	}
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, fromWallet, types.TotalUnissuedPrepay, sdk.NewCoins(coinToSend))
}

func (k Keeper) GetTotalUnissuedPrepay(ctx sdk.Context) (totalUnissuedPrepay sdk.Coin) {
	totalUnissuedPrepayAccAddr := k.accountKeeper.GetModuleAddress(regtypes.TotalUnissuedPrepay)
	if totalUnissuedPrepayAccAddr == nil {
		ctx.Logger().Error("account address for total unissued prepay does not exist.")
		return sdk.Coin{
			Denom:  types.DefaultBondDenom,
			Amount: sdk.ZeroInt(),
		}
	}
	return k.bankKeeper.GetBalance(ctx, totalUnissuedPrepayAccAddr, k.BondDenom(ctx))
}

func (k Keeper) SetInitialGenesisStakeTotal(ctx sdk.Context, stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := types.ModuleCdc.MustMarshalLengthPrefixed(stake)
	store.Set(types.InitialGenesisStakeTotalKey, b)
}

func (k Keeper) GetInitialGenesisStakeTotal(ctx sdk.Context) (stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.InitialGenesisStakeTotalKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(b, &stake)
	return
}

func (k Keeper) SetRemainingOzoneLimit(ctx sdk.Context, value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := types.ModuleCdc.MustMarshalLengthPrefixed(value)
	store.Set(types.UpperBoundOfTotalOzoneKey, b)
}

func (k Keeper) GetRemainingOzoneLimit(ctx sdk.Context) (value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.UpperBoundOfTotalOzoneKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(b, &value)
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

// GetUnbondingNode return a unbonding UnbondingMetaNode
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

// HasMaxUnbondingMetaNodeEntries - check if unbonding MetaNode has maximum number of entries
func (k Keeper) HasMaxUnbondingNodeEntries(ctx sdk.Context, networkAddr stratos.SdsAddress) bool {
	ubd, found := k.GetUnbondingNode(ctx, networkAddr)
	if !found {
		return false
	}
	return len(ubd.Entries) >= int(k.MaxEntries(ctx))
}

// SetUnbondingNode sets the unbonding MetaNode
func (k Keeper) SetUnbondingNode(ctx sdk.Context, ubd types.UnbondingNode) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalUnbondingNode(k.cdc, ubd)
	networkAddr, err := stratos.SdsAddressFromBech32(ubd.GetNetworkAddr())
	if err != nil {
		return
	}
	key := types.GetUBDNodeKey(networkAddr)
	store.Set(key, bz)
}

// RemoveUnbondingNode removes the unbonding MetaNode object
func (k Keeper) RemoveUnbondingNode(ctx sdk.Context, ubd types.UnbondingNode) {
	store := ctx.KVStore(k.storeKey)
	networkAddr, err := stratos.SdsAddressFromBech32(ubd.GetNetworkAddr())
	if err != nil {
		return
	}
	key := types.GetUBDNodeKey(networkAddr)
	store.Delete(key)
}

// SetUnbondingMetaNodeEntry adds an entry to the unbonding MetaNode at
// the given addresses. It creates the unbonding MetaNode if it does not exist
func (k Keeper) SetUnbondingNodeEntry(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool,
	creationHeight int64, minTime time.Time, balance sdk.Int) types.UnbondingNode {

	ubd, found := k.GetUnbondingNode(ctx, networkAddr)
	if found {
		ubd.AddEntry(creationHeight, minTime, balance)
	} else {
		ubd = types.NewUnbondingNode(networkAddr, isMetaNode, creationHeight, minTime, balance)
	}
	k.SetUnbondingNode(ctx, ubd)
	return ubd
}

// unbonding delegation queue timeslice operations

// GetUnbondingNodeQueueTimeSlice gets a specific unbonding queue timeslice. A timeslice is a slice of DVPairs
// corresponding to unbonding delegations that expire at a certain time.
func (k Keeper) GetUnbondingNodeQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (networkAddrs []stratos.SdsAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetUBDTimeKey(timestamp))
	if bz == nil {
		return []stratos.SdsAddress{}
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(bz, &networkAddrs)
	return networkAddrs
}

// SetUnbondingNodeQueueTimeSlice sets a specific unbonding queue timeslice.
func (k Keeper) SetUnbondingNodeQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []stratos.SdsAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := types.ModuleCdc.MustMarshalLengthPrefixed(keys)
	store.Set(types.GetUBDTimeKey(timestamp), bz)
}

// InsertUnbondingNodeQueue inserts an unbonding delegation to the appropriate timeslice in the unbonding queue
func (k Keeper) InsertUnbondingNodeQueue(ctx sdk.Context, ubd types.UnbondingNode,
	completionTime time.Time) {

	timeSlice := k.GetUnbondingNodeQueueTimeSlice(ctx, completionTime)
	networkAddr, err := stratos.SdsAddressFromBech32(ubd.GetNetworkAddr())
	if err != nil {
		return
	}
	if len(timeSlice) == 0 {
		k.SetUnbondingNodeQueueTimeSlice(ctx, completionTime, []stratos.SdsAddress{networkAddr})
	} else {
		timeSlice = append(timeSlice, networkAddr)
		k.SetUnbondingNodeQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// UnbondingNodeQueueIterator returns all the unbonding queue timeslices from time 0 until endTime
func (k Keeper) UnbondingNodeQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UBDNodeQueueKey,
		sdk.InclusiveEndBytes(types.GetUBDTimeKey(endTime)))
}

// DequeueAllMatureUBDQueue returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue
// Iteration for dequeuing  all mature unbonding queue
func (k Keeper) DequeueAllMatureUBDQueue(ctx sdk.Context,
	currTime time.Time) (matureUnbonds []stratos.SdsAddress) {

	store := ctx.KVStore(k.storeKey)
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	unbondingTimesliceIterator := k.UnbondingNodeQueueIterator(ctx, ctx.BlockHeader().Time)
	defer unbondingTimesliceIterator.Close()

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := []stratos.SdsAddress{}
		value := unbondingTimesliceIterator.Value()
		types.ModuleCdc.MustUnmarshalLengthPrefixed(value, &timeslice)
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
				amt := sdk.NewCoin(bondDenom, *entry.Balance)
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

	return balances, ubd.IsMetaNode, nil
}

func (k Keeper) SubtractUBDNodeStake(ctx sdk.Context, ubd types.UnbondingNode, tokenToSub sdk.Coin) error {
	// case of meta node
	networkAddr, err := stratos.SdsAddressFromBech32(ubd.GetNetworkAddr())
	if err != nil {
		return err
	}
	if ubd.IsMetaNode {
		metaNode, found := k.GetMetaNode(ctx, networkAddr)
		if !found {
			return types.ErrNoMetaNodeFound
		}
		return k.SubtractMetaNodeStake(ctx, metaNode, tokenToSub)
	}
	// case of resource node
	resourceNode, found := k.GetResourceNode(ctx, networkAddr)
	if !found {
		return types.ErrNoMetaNodeFound
	}
	return k.SubtractResourceNodeStake(ctx, resourceNode, tokenToSub)
}

func (k Keeper) UnbondResourceNode(
	ctx sdk.Context, resourceNode types.ResourceNode, amt sdk.Int,
) (ozoneLimitChange sdk.Int, unbondingMatureTime time.Time, err error) {
	params := k.GetParams(ctx)
	ctx.Logger().Info("Params of register module: " + params.String())

	// transfer the node tokens to the not bonded pool
	networkAddr, err := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddress())
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

		// decrease resource node count
		v := k.GetBondedResourceNodeCnt(ctx)
		count := v.Sub(sdk.NewInt(1))
		k.SetBondedResourceNodeCnt(ctx, count)
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

func (k Keeper) UnbondMetaNode(
	ctx sdk.Context, metaNode types.MetaNode, amt sdk.Int,
) (ozoneLimitChange sdk.Int, unbondingMatureTime time.Time, err error) {

	networkAddr, err := stratos.SdsAddressFromBech32(metaNode.GetNetworkAddress())
	if err != nil {
		return sdk.ZeroInt(), time.Now(), errors.New("invalid network address")
	}
	ownerAddr, err := sdk.AccAddressFromBech32(metaNode.GetOwnerAddress())
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

	unbondingMatureTime = calcUnbondingMatureTime(ctx, metaNode.Status, metaNode.CreationTime, k.UnbondingThreasholdTime(ctx), k.UnbondingCompletionTime(ctx))

	bondDenom := k.GetParams(ctx).BondDenom
	coin := sdk.NewCoin(bondDenom, amt)
	if metaNode.GetStatus() == stakingtypes.Bonded {
		// transfer the node tokens to the not bonded pool
		k.bondedToUnbonding(ctx, metaNode, true, coin)
		// adjust ozone limit
		ozoneLimitChange = k.decreaseOzoneLimitBySubtractStake(ctx, amt)
	}
	// change node status to unbonding if unbonding all tokens
	if amt.Equal(metaNode.Tokens) {
		metaNode.Status = stakingtypes.Unbonding
		// decrease meta node count
		v := k.GetBondedMetaNodeCnt(ctx)
		count := v.Sub(sdk.NewInt(1))
		k.SetBondedMetaNodeCnt(ctx, count)
		// set meta node
		k.SetMetaNode(ctx, metaNode)
	}

	// Set the unbonding mature time and completion height appropriately
	unbondingNode := k.SetUnbondingNodeEntry(ctx, networkAddr, true, ctx.BlockHeight(), unbondingMatureTime, amt)
	// Add to unbonding node queue
	k.InsertUnbondingNodeQueue(ctx, unbondingNode, unbondingMatureTime)
	ctx.Logger().Info("Unbonding meta node " + unbondingNode.String() + "\n after mature time" + unbondingMatureTime.String())
	return ozoneLimitChange, unbondingMatureTime, nil
}

// GetAllUnbondingNodesTotalBalance Iteration for getting the total balance of all unbonding nodes
func (k Keeper) GetAllUnbondingNodesTotalBalance(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UBDNodeKey)
	defer iterator.Close()

	var ubdTotal = sdk.ZeroInt()
	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalUnbondingNode(k.cdc, iterator.Value())
		for _, entry := range node.Entries {
			ubdTotal = ubdTotal.Add(*entry.Balance)
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
		balance = balance.Add(*entry.Balance)
	}
	return balance
}

// CurrUozPrice calcs current uoz price
func (k Keeper) CurrUozPrice(ctx sdk.Context) sdk.Dec {
	S := k.GetInitialGenesisStakeTotal(ctx)
	Pt := k.GetTotalUnissuedPrepay(ctx).Amount
	Lt := k.GetRemainingOzoneLimit(ctx)
	currUozPrice := (S.Add(Pt)).ToDec().
		Quo(Lt.ToDec())
	return currUozPrice
}

// UozSupply calc remaining/total supply for uoz
func (k Keeper) UozSupply(ctx sdk.Context) (remaining, total sdk.Int) {
	remaining = k.GetRemainingOzoneLimit(ctx) // Lt
	S := k.GetInitialGenesisStakeTotal(ctx)
	Pt := k.GetTotalUnissuedPrepay(ctx).Amount
	// total supply = Lt * ( 1 + Pt / S )
	total = (Pt.ToDec().Quo(S.ToDec()).TruncateInt().Add(sdk.NewInt(1))).Mul(remaining)
	return remaining, total
}

func (k Keeper) SetBondedResourceNodeCnt(ctx sdk.Context, count sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := types.ModuleCdc.MustMarshalLengthPrefixed(count)
	//ctx.Logger().Info("New ResourceNode count = " + count.String())
	store.Set(types.ResourceNodeCntKey, b)
}

func (k Keeper) SetBondedMetaNodeCnt(ctx sdk.Context, count sdk.Int) {
	if count.LT(sdk.ZeroInt()) {
		//ctx.Logger().Info("count < 0, count = " + count.String())
	}
	store := ctx.KVStore(k.storeKey)
	b := types.ModuleCdc.MustMarshalLengthPrefixed(count)
	//ctx.Logger().Info("New MetaNode count = " + count.String())
	store.Set(types.MetaNodeCntKey, b)
}

func (k Keeper) GetBondedResourceNodeCnt(ctx sdk.Context) (balance sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.ResourceNodeCntKey)
	if value == nil {
		return sdk.ZeroInt()
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(value, &balance)
	//ctx.Logger().Info("ResourceNode count = " + balance.String())
	return
}

func (k Keeper) GetBondedMetaNodeCnt(ctx sdk.Context) (balance sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.MetaNodeCntKey)
	if value == nil {
		return sdk.ZeroInt()
	}
	types.ModuleCdc.MustUnmarshalLengthPrefixed(value, &balance)
	//ctx.Logger().Info("MetaNode count = " + balance.String())
	return
}

func (k Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
}
