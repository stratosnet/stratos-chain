package keeper

import (
	"fmt"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// Keeper of the pot store
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              codec.Codec
	paramSpace       paramstypes.Subspace
	feeCollectorName string // name of the FeeCollector ModuleAccount
	BankKeeper       types.BankKeeper
	AccountKeeper    types.AccountKeeper
	StakingKeeper    types.StakingKeeper
	RegisterKeeper   types.RegisterKeeper
	DistrKeeper      types.DistrKeeper
}

// NewKeeper creates a pot keeper
func NewKeeper(cdc codec.Codec, key sdk.StoreKey, paramSpace paramstypes.Subspace, feeCollectorName string,
	bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper, stakingKeeper types.StakingKeeper,
	registerKeeper types.RegisterKeeper, distrKeeper types.DistrKeeper,
) Keeper {
	keeper := Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace.WithKeyTable(types.ParamKeyTable()),
		feeCollectorName: feeCollectorName,
		BankKeeper:       bankKeeper,
		AccountKeeper:    accountKeeper,
		StakingKeeper:    stakingKeeper,
		RegisterKeeper:   registerKeeper,
		DistrKeeper:      distrKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) VolumeReport(ctx sdk.Context, walletVolumes types.WalletVolumes, reporter stratos.SdsAddress,
	epoch sdk.Int, reportReference string, txHash string) (err error) {

	// tx fail when unhandled epoch > 0
	if k.GetUnhandledEpoch(ctx).GT(sdk.ZeroInt()) {
		return types.ErrRewardDistributionNotComplete
	}

	//record volume report
	reportRecord := types.NewReportRecord(reporter, reportReference, txHash)
	k.SetVolumeReport(ctx, epoch, reportRecord)

	// save for reward distribution in the EndBlock at height + 1
	k.SetIsReadyToDistributeReward(ctx, false)
	k.SetUnhandledEpoch(ctx, epoch)
	k.SetUnhandledReport(ctx, walletVolumes)

	return nil
}

func (k Keeper) IsSPNode(ctx sdk.Context, p2pAddr stratos.SdsAddress) (found bool) {
	_, found = k.RegisterKeeper.GetMetaNode(ctx, p2pAddr)
	return found
}

func (k Keeper) FoundationDeposit(ctx sdk.Context, amount sdk.Coins, from sdk.AccAddress) (err error) {
	err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, from, types.FoundationAccount, amount)
	if err != nil {
		return err
	}
	return nil
}
