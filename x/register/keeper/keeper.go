package keeper

import (
	"fmt"
	"time"

	"github.com/cometbft/cometbft/libs/log"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

// Keeper of the register store
type Keeper struct {
	storeKey      storetypes.StoreKey
	cdc           codec.Codec
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	distrKeeper   types.DistrKeeper
	hooks         types.RegisterHooks

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper creates a register keeper
func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrKeeper types.DistrKeeper,
	authority string,
) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrKeeper:   distrKeeper,
		authority:     authority,
	}
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
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
			Amount: sdkmath.ZeroInt(),
		}
	}
	return k.bankKeeper.GetBalance(ctx, totalUnissuedPrepayAccAddr, k.BondDenom(ctx))
}

func (k Keeper) IncreaseOzoneLimitByAddDeposit(ctx sdk.Context, deposit sdkmath.Int) (ozoneLimitChange sdkmath.Int) {
	// get remainingOzoneLimit before adding deposit
	remainingBefore := k.GetRemainingOzoneLimit(ctx)
	depositNozRate := k.GetDepositNozRate(ctx)

	// update effectiveTotalDeposit
	effectiveTotalDepositBefore := k.GetEffectiveTotalDeposit(ctx)
	effectiveTotalDepositAfter := effectiveTotalDepositBefore.Add(deposit)
	k.SetEffectiveTotalDeposit(ctx, effectiveTotalDepositAfter)

	effectiveGenesisDeposit := effectiveTotalDepositBefore.ToLegacyDec() //wei
	if effectiveGenesisDeposit.Equal(sdkmath.LegacyZeroDec()) {
		ctx.Logger().Info("effectiveGenesisDeposit is zero, increase ozone limit failed")
		return sdkmath.ZeroInt()
	}

	limitToAdd := deposit.ToLegacyDec().Quo(depositNozRate)
	k.SetRemainingOzoneLimit(ctx, remainingBefore.ToLegacyDec().Add(limitToAdd).TruncateInt())

	//ctx.Logger().Debug("----- IncreaseOzoneLimitByAddDeposit, ",
	//	"effectiveTotalDepositBefore=", effectiveTotalDepositBefore.String(),
	//	"effectiveTotalDepositAfter=", effectiveTotalDepositAfter.String(),
	//	"remainingBefore=", remainingBefore.String(),
	//	"remainingAfter=", k.GetRemainingOzoneLimit(ctx).String(),
	//)
	return limitToAdd.TruncateInt()
}

func (k Keeper) DecreaseOzoneLimitBySubtractDeposit(ctx sdk.Context, deposit sdkmath.Int) (ozoneLimitChange sdkmath.Int) {
	// get remainingOzoneLimit before adding deposit
	remainingBefore := k.GetRemainingOzoneLimit(ctx)
	depositNozRate := k.GetDepositNozRate(ctx)

	// update effectiveTotalDeposit
	effectiveTotalDepositBefore := k.GetEffectiveTotalDeposit(ctx)
	effectiveTotalDepositAfter := effectiveTotalDepositBefore.Sub(deposit)
	k.SetEffectiveTotalDeposit(ctx, effectiveTotalDepositAfter)

	effectiveGenesisDeposit := effectiveTotalDepositBefore.ToLegacyDec() //wei
	if effectiveGenesisDeposit.Equal(sdkmath.LegacyZeroDec()) {
		ctx.Logger().Info("effectiveGenesisDeposit is zero, increase ozone limit failed")
		return sdkmath.ZeroInt()
	}
	limitToSub := deposit.ToLegacyDec().Quo(depositNozRate)
	k.SetRemainingOzoneLimit(ctx, remainingBefore.ToLegacyDec().Sub(limitToSub).TruncateInt())

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

// SetUnbondingNodeEntry adds an entry to the unbonding MetaNode at
// the given addresses. It creates the unbonding MetaNode if it does not exist
func (k Keeper) SetUnbondingNodeEntry(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool,
	creationHeight int64, minTime time.Time, balance sdkmath.Int) types.UnbondingNode {

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

// GetAllUnbondingNodesTotalBalance Iteration for getting the total balance of all unbonding nodes
func (k Keeper) GetAllUnbondingNodesTotalBalance(ctx sdk.Context) sdkmath.Int {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UBDNodeKey)
	defer iterator.Close()

	var ubdTotal = sdkmath.ZeroInt()
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
func (k Keeper) GetUnbondingNodeBalance(ctx sdk.Context, networkAddr stratos.SdsAddress) sdkmath.Int {
	balance := sdkmath.ZeroInt()

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
func (k Keeper) GetCurrNozPriceParams(ctx sdk.Context) (St, Pt, Lt sdkmath.Int) {
	St = k.GetEffectiveTotalDeposit(ctx)
	Pt = k.GetTotalUnissuedPrepay(ctx).Amount
	Lt = k.GetRemainingOzoneLimit(ctx)
	return
}
