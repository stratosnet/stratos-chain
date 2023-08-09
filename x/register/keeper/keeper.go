package keeper

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

// Keeper of the register store
type Keeper struct {
	storeKey                       sdk.StoreKey
	cdc                            codec.Codec
	paramSpace                     paramtypes.Subspace
	accountKeeper                  types.AccountKeeper
	bankKeeper                     types.BankKeeper
	distrKeeper                    types.DistrKeeper
	hooks                          types.RegisterHooks
	metaNodeBitMapIndexCache       map[string]int
	metaNodeBitMapIndexCacheStatus *types.CacheStatus
	cacheMutex                     *sync.Mutex
}

// NewKeeper creates a register keeper
func NewKeeper(cdc codec.Codec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, distrKeeper types.DistrKeeper) Keeper {

	keeper := Keeper{
		storeKey:                       key,
		cdc:                            cdc,
		paramSpace:                     paramSpace.WithKeyTable(types.ParamKeyTable()),
		accountKeeper:                  accountKeeper,
		bankKeeper:                     bankKeeper,
		distrKeeper:                    distrKeeper,
		hooks:                          nil,
		metaNodeBitMapIndexCache:       make(map[string]int),
		metaNodeBitMapIndexCacheStatus: &types.CacheStatus{Status: types.CACHE_DIRTY},
		cacheMutex:                     &sync.Mutex{},
	}
	return keeper
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
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

func (k Keeper) GetTotalUnissuedPrepay(ctx sdk.Context) (totalUnissuedPrepay sdk.Coin) {
	totalUnissuedPrepayAccAddr := k.accountKeeper.GetModuleAddress(regtypes.TotalUnissuedPrepay)
	if totalUnissuedPrepayAccAddr == nil {
		ctx.Logger().Error("account address for total unissued prepay does not exist.")
		return sdk.Coin{
			Denom:  k.BondDenom(ctx),
			Amount: sdk.ZeroInt(),
		}
	}
	return k.bankKeeper.GetBalance(ctx, totalUnissuedPrepayAccAddr, k.BondDenom(ctx))
}

func (k Keeper) IncreaseOzoneLimitByAddDeposit(ctx sdk.Context, deposit sdk.Int) (ozoneLimitChange sdk.Int) {
	// get remainingOzoneLimit before adding deposit
	remainingBefore := k.GetRemainingOzoneLimit(ctx)
	depositNozRate := k.GetDepositNozRate(ctx)

	// update effectiveTotalDeposit
	effectiveTotalDepositBefore := k.GetEffectiveTotalDeposit(ctx)
	effectiveTotalDepositAfter := effectiveTotalDepositBefore.Add(deposit)
	k.SetEffectiveTotalDeposit(ctx, effectiveTotalDepositAfter)

	effectiveGenesisDeposit := effectiveTotalDepositBefore.ToDec() //wei
	if effectiveGenesisDeposit.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("effectiveGenesisDeposit is zero, increase ozone limit failed")
		return sdk.ZeroInt()
	}

	limitToAdd := deposit.ToDec().Quo(depositNozRate)
	k.SetRemainingOzoneLimit(ctx, remainingBefore.ToDec().Add(limitToAdd).TruncateInt())

	//ctx.Logger().Debug("----- IncreaseOzoneLimitByAddDeposit, ",
	//	"effectiveTotalDepositBefore=", effectiveTotalDepositBefore.String(),
	//	"effectiveTotalDepositAfter=", effectiveTotalDepositAfter.String(),
	//	"remainingBefore=", remainingBefore.String(),
	//	"remainingAfter=", k.GetRemainingOzoneLimit(ctx).String(),
	//)
	return limitToAdd.TruncateInt()
}

func (k Keeper) DecreaseOzoneLimitBySubtractDeposit(ctx sdk.Context, deposit sdk.Int) (ozoneLimitChange sdk.Int) {
	// get remainingOzoneLimit before adding deposit
	remainingBefore := k.GetRemainingOzoneLimit(ctx)
	depositNozRate := k.GetDepositNozRate(ctx)

	// update effectiveTotalDeposit
	effectiveTotalDepositBefore := k.GetEffectiveTotalDeposit(ctx)
	effectiveTotalDepositAfter := effectiveTotalDepositBefore.Sub(deposit)
	k.SetEffectiveTotalDeposit(ctx, effectiveTotalDepositAfter)

	effectiveGenesisDeposit := effectiveTotalDepositBefore.ToDec() //wei
	if effectiveGenesisDeposit.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("effectiveGenesisDeposit is zero, increase ozone limit failed")
		return sdk.ZeroInt()
	}
	limitToSub := deposit.ToDec().Quo(depositNozRate)
	k.SetRemainingOzoneLimit(ctx, remainingBefore.ToDec().Sub(limitToSub).TruncateInt())

	//ctx.Logger().Debug("----- DecreaseOzoneLimitBySubtractDeposit, ",
	//	"effectiveTotalDepositBefore=", effectiveTotalDepositBefore.String(),
	//	"effectiveTotalDepositAfter=", effectiveTotalDepositAfter.String(),
	//	"remainingBefore=", remainingBefore.String(),
	//	"remainingAfter=", k.GetRemainingOzoneLimit(ctx).String(),
	//)
	return limitToSub.TruncateInt()
}

// HasMaxUnbondingNodeEntries - check if unbonding node has maximum number of entries
func (k Keeper) HasMaxUnbondingNodeEntries(ctx sdk.Context, networkAddr stratos.SdsAddress) bool {
	ubd, found := k.GetUnbondingNode(ctx, networkAddr)
	if !found {
		return false
	}
	return len(ubd.Entries) >= int(k.MaxEntries(ctx))
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

// InsertUnbondingNodeQueue inserts an unbonding delegation to the appropriate timeslice in the unbonding queue
func (k Keeper) InsertUnbondingNodeQueue(ctx sdk.Context, ubd types.UnbondingNode, completionTime time.Time) {
	timeSlice := k.GetUnbondingNodeQueueTimeSlice(ctx, completionTime)
	networkAddr := ubd.GetNetworkAddr()

	if len(timeSlice) == 0 {
		k.SetUnbondingNodeQueueTimeSlice(ctx, completionTime, []string{networkAddr})
	} else {
		timeSlice = append(timeSlice, networkAddr)
		k.SetUnbondingNodeQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// DequeueAllMatureUBDQueue returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue
// Iteration for dequeuing  all mature unbonding queue
// TODO: Unused parameter: currTime
func (k Keeper) DequeueAllMatureUBDQueue(ctx sdk.Context, currTime time.Time) (matureUnbonds []string) {
	keysToDelete := make([][]byte, 0)
	store := ctx.KVStore(k.storeKey)
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	unbondingTimesliceIterator := k.UnbondingNodeQueueIterator(ctx, ctx.BlockHeader().Time)
	defer unbondingTimesliceIterator.Close()

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeSliceVal := stratos.SdsAddresses{} //[]stratos.SdsAddress{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshalLengthPrefixed(value, &timeSliceVal)
		timeSlice := timeSliceVal.GetAddresses()
		matureUnbonds = append(matureUnbonds, timeSlice...)
		keysToDelete = append(keysToDelete, unbondingTimesliceIterator.Key())
	}
	// safe removal
	for _, key := range keysToDelete {
		store.Delete(key)
	}
	ctx.Logger().Debug(fmt.Sprintf("DequeueAllMatureUBDQueue, %d matured unbonding nodes detected", len(matureUnbonds)))
	return matureUnbonds
}

// CompleteUnbondingWithAmount completes the unbonding of all mature entries in
// the retrieved unbonding delegation object and returns the total unbonding
// balance or an error upon failure.
func (k Keeper) CompleteUnbondingWithAmount(ctx sdk.Context, networkAddrBech32 string) (sdk.Coins, bool, error) {
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrBech32)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("NetworAddr: %s is invalid", networkAddrBech32))
		return nil, false, types.ErrInvalidNetworkAddr
	}

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
				err := k.subtractUBDNodeDeposit(ctx, ubd, amt)
				if err != nil {
					return nil, false, err
				}

				balances = balances.Add(amt)
			}
		}
	}

	// set the unbonding node or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveUnbondingNode(ctx, networkAddr)
	} else {
		k.SetUnbondingNode(ctx, ubd)
	}

	return balances, ubd.IsMetaNode, nil
}

func (k Keeper) subtractUBDNodeDeposit(ctx sdk.Context, ubd types.UnbondingNode, tokenToSub sdk.Coin) error {
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
		return k.SubtractMetaNodeDeposit(ctx, metaNode, tokenToSub)
	}
	// case of resource node
	resourceNode, found := k.GetResourceNode(ctx, networkAddr)
	if !found {
		return types.ErrNoMetaNodeFound
	}
	return k.SubtractResourceNodeDeposit(ctx, resourceNode, tokenToSub)
}

// Unbond all tokens of resource node
func (k Keeper) UnbondResourceNode(ctx sdk.Context, networkAddr stratos.SdsAddress, ownerAddr sdk.AccAddress,
) (depositToRemove sdk.Int, unbondingMatureTime time.Time, err error) {

	resourceNode, found := k.GetResourceNode(ctx, networkAddr)
	if !found {
		return sdk.ZeroInt(), time.Time{}, types.ErrNoResourceNodeFound
	}
	ownerAddrNode, _ := sdk.AccAddressFromBech32(resourceNode.GetOwnerAddress())
	if !ownerAddrNode.Equals(ownerAddr) {
		return sdk.ZeroInt(), time.Time{}, types.ErrInvalidOwnerAddr
	}
	if resourceNode.GetStatus() != stakingtypes.Bonded {
		return sdk.ZeroInt(), time.Time{}, types.ErrInvalidNodeStat
	}
	// suspended node cannot be unbonded (avoid dup deposit decrease with node suspension)
	if resourceNode.GetSuspend() {
		return sdk.ZeroInt(), time.Time{}, types.ErrInvalidSuspensionStatForUnbondNode
	}
	if k.HasMaxUnbondingNodeEntries(ctx, networkAddr) {
		return sdk.ZeroInt(), time.Time{}, types.ErrMaxUnbondingNodeEntries
	}

	// check if node_token - unbonding_token > 0
	unbondingDeposit := k.GetUnbondingNodeBalance(ctx, networkAddr)
	depositToRemove = resourceNode.Tokens.Sub(unbondingDeposit)
	if depositToRemove.LTE(sdk.ZeroInt()) {
		return sdk.ZeroInt(), time.Time{}, types.ErrInsufficientBalance
	}

	unbondingMatureTime = calcUnbondingMatureTime(ctx, resourceNode.Status, resourceNode.CreationTime, k.UnbondingThreasholdTime(ctx), k.UnbondingCompletionTime(ctx))

	// transfer the node tokens to the not bonded pool
	k.bondedToUnbonding(ctx, resourceNode, false, sdk.NewCoin(k.BondDenom(ctx), depositToRemove))
	// change node status to unbonding if unbonding all available tokens
	resourceNode.Status = stakingtypes.Unbonding
	k.SetResourceNode(ctx, resourceNode)
	// decrease resource node count
	v := k.GetBondedResourceNodeCnt(ctx)
	count := v.Sub(sdk.OneInt())
	k.SetBondedResourceNodeCnt(ctx, count)

	// set the unbonding mature time and completion height appropriately
	ctx.Logger().Info(fmt.Sprintf("Calculating mature time: creationTime[%s], threasholdTime[%s], completionTime[%s], matureTime[%s]",
		resourceNode.CreationTime, k.UnbondingThreasholdTime(ctx), k.UnbondingCompletionTime(ctx), unbondingMatureTime,
	))
	unbondingNode := k.SetUnbondingNodeEntry(ctx, networkAddr, false, ctx.BlockHeight(), unbondingMatureTime, depositToRemove)

	// Add to unbonding node queue
	k.InsertUnbondingNodeQueue(ctx, unbondingNode, unbondingMatureTime)
	ctx.Logger().Info("Unbonding resource node " + unbondingNode.String() + "\n after mature time" + unbondingMatureTime.String())

	return depositToRemove, unbondingMatureTime, nil
}

func (k Keeper) UnbondMetaNode(ctx sdk.Context, metaNode types.MetaNode, amt sdk.Int,
) (ozoneLimitChange, availableTokenAmtBefore, availableTokenAmtAfter sdk.Int, unbondingMatureTime time.Time, err error) {
	if metaNode.GetStatus() == stakingtypes.Unbonding {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), time.Time{}, types.ErrUnbondingNode
	}

	networkAddr, err := stratos.SdsAddressFromBech32(metaNode.GetNetworkAddress())
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), time.Time{}, errors.New("invalid network address")
	}
	ownerAddr, err := sdk.AccAddressFromBech32(metaNode.GetOwnerAddress())
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), time.Time{}, errors.New("invalid wallet address")
	}

	ownerAcc := k.accountKeeper.GetAccount(ctx, ownerAddr)
	if ownerAcc == nil {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), time.Time{}, types.ErrNoOwnerAccountFound
	}

	// suspended node cannot be unbonded (avoid dup deposit decrease with node suspension)
	if metaNode.Suspend {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), time.Time{}, types.ErrInvalidSuspensionStatForUnbondNode
	}

	// check if node_token - unbonding_token > amt_to_unbond
	unbondingDeposit := k.GetUnbondingNodeBalance(ctx, networkAddr)
	availableTokenAmtBefore = metaNode.Tokens.Sub(unbondingDeposit)
	if availableTokenAmtBefore.LT(amt) {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), time.Time{}, types.ErrInsufficientBalance
	}

	if k.HasMaxUnbondingNodeEntries(ctx, networkAddr) {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), time.Time{}, types.ErrMaxUnbondingNodeEntries
	}

	unbondingMatureTime = calcUnbondingMatureTime(ctx, metaNode.Status, metaNode.CreationTime, k.UnbondingThreasholdTime(ctx), k.UnbondingCompletionTime(ctx))

	bondDenom := k.GetParams(ctx).BondDenom
	coin := sdk.NewCoin(bondDenom, amt)
	if metaNode.GetStatus() == stakingtypes.Bonded {
		// to prevent remainingOzoneLimit from being negative value
		if !k.IsUnbondable(ctx, amt) {
			return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), time.Time{}, types.ErrInsufficientBalance
		}
		// transfer the node tokens to the not bonded pool
		k.bondedToUnbonding(ctx, metaNode, true, coin)
		availableTokenAmtAfter = availableTokenAmtBefore.Sub(amt)
		// adjust ozone limit
		ozoneLimitChange = k.DecreaseOzoneLimitBySubtractDeposit(ctx, amt)
	}
	// change node status to unbonding if unbonding all available tokens
	if amt.Equal(availableTokenAmtBefore) {
		metaNode.Status = stakingtypes.Unbonding
		// decrease meta node count
		v := k.GetBondedMetaNodeCnt(ctx)
		count := v.Sub(sdk.NewInt(1))
		k.SetBondedMetaNodeCnt(ctx, count)
		// set meta node
		k.SetMetaNode(ctx, metaNode)
		// remove record from vote pool
		if _, found := k.GetMetaNodeRegistrationVotePool(ctx, networkAddr); found {
			ctx.Logger().Info("DeleteMetaNodeRegistrationVotePool of meta node " + networkAddr.String())
			k.DeleteMetaNodeRegistrationVotePool(ctx, networkAddr)
		}
	}

	// Set the unbonding mature time and completion height appropriately
	unbondingNode := k.SetUnbondingNodeEntry(ctx, networkAddr, true, ctx.BlockHeight(), unbondingMatureTime, amt)
	// Add to unbonding node queue
	k.InsertUnbondingNodeQueue(ctx, unbondingNode, unbondingMatureTime)
	ctx.Logger().Info("Unbonding meta node " + unbondingNode.String() + "\n after mature time" + unbondingMatureTime.String())
	return ozoneLimitChange, availableTokenAmtBefore, availableTokenAmtAfter, unbondingMatureTime, nil
}

// GetAllUnbondingNodesTotalBalance Iteration for getting the total balance of all unbonding nodes
func (k Keeper) GetAllUnbondingNodesTotalBalance(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UBDNodeKey)
	defer iterator.Close()

	var ubdTotal = sdk.ZeroInt()
	for ; iterator.Valid(); iterator.Next() {
		node := types.UnbondingNode{}
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &node)
		for _, entry := range node.Entries {
			ubdTotal = ubdTotal.Add(*entry.Balance)
		}
	}
	return ubdTotal
}

// GetUnbondingNodeBalance returns an unbonding balance and an UnbondingNode
func (k Keeper) GetUnbondingNodeBalance(ctx sdk.Context, networkAddr stratos.SdsAddress) sdk.Int {
	balance := sdk.ZeroInt()

	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDNodeKey(networkAddr)
	value := store.Get(key)
	if value == nil {
		return balance
	}
	ubd := types.UnbondingNode{}
	k.cdc.MustUnmarshalLengthPrefixed(value, &ubd)
	for _, entry := range ubd.Entries {
		balance = balance.Add(*entry.Balance)
	}
	return balance
}

// GetCurrNozPriceParams calcs current noz price
func (k Keeper) GetCurrNozPriceParams(ctx sdk.Context) (St, Pt, Lt sdk.Int) {
	St = k.GetEffectiveTotalDeposit(ctx)
	Pt = k.GetTotalUnissuedPrepay(ctx).Amount
	Lt = k.GetRemainingOzoneLimit(ctx)
	return
}
