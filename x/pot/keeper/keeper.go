package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// Keeper of the pot store
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	paramSpace       params.Subspace
	feeCollectorName string // name of the FeeCollector ModuleAccount
	BankKeeper       bank.Keeper
	SupplyKeeper     supply.Keeper
	AccountKeeper    auth.AccountKeeper
	StakingKeeper    staking.Keeper
	RegisterKeeper   register.Keeper
}

// NewKeeper creates a pot keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, feeCollectorName string,
	bankKeeper bank.Keeper, supplyKeeper supply.Keeper, accountKeeper auth.AccountKeeper, stakingKeeper staking.Keeper,
	registerKeeper register.Keeper,
) Keeper {
	keeper := Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace.WithKeyTable(types.ParamKeyTable()),
		feeCollectorName: feeCollectorName,
		BankKeeper:       bankKeeper,
		SupplyKeeper:     supplyKeeper,
		AccountKeeper:    accountKeeper,
		StakingKeeper:    stakingKeeper,
		RegisterKeeper:   registerKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetVolumeReport(ctx sdk.Context, epoch sdk.Int) (res types.VolumeReportRecord) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VolumeReportStoreKey(epoch))
	if bz == nil {
		return types.VolumeReportRecord{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return res
}

func (k Keeper) SetVolumeReport(ctx sdk.Context, epoch sdk.Int, reportRecord types.VolumeReportRecord) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.VolumeReportStoreKey(epoch)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(reportRecord)
	store.Set(storeKey, bz)
}

func (k Keeper) IsSPNode(ctx sdk.Context, p2pAddr sdk.AccAddress) (found bool) {
	_, found = k.RegisterKeeper.GetIndexingNode(ctx, p2pAddr)
	return found
}

func (k Keeper) FoundationDeposit(ctx sdk.Context, amount sdk.Coin, from sdk.AccAddress) (err error) {
	_, err = k.BankKeeper.SubtractCoins(ctx, from, sdk.NewCoins(amount))
	if err != nil {
		return err
	}

	foundationAccountAddr := k.SupplyKeeper.GetModuleAddress(types.FoundationAccount)
	_, err = k.BankKeeper.AddCoins(ctx, foundationAccountAddr, sdk.NewCoins(amount))
	if err != nil {
		return err
	}

	return nil

}
