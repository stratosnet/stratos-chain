package keeper

import (
	"encoding/hex"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

// GetVolumeReport returns the hash of volume report
//func (k Keeper) GetVolumeReport(ctx sdk.Context, reporter sdk.AccAddress) ([]byte, error) {
//	store := ctx.KVStore(k.storeKey)
//	bz := store.Get(types.VolumeReportStoreKey(reporter))
//	if bz == nil {
//		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress,
//			"key %s does not exist", hex.EncodeToString(types.VolumeReportStoreKey(reporter)))
//	}
//	return bz, nil
//}

func (k Keeper) GetVolumeReport(ctx sdk.Context, epoch sdk.Int) ([]byte, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VolumeReportStoreKey(epoch))
	if bz == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress,
			"key %s does not exist", hex.EncodeToString(types.VolumeReportStoreKey(epoch)))
	}
	return bz, nil
}

//func (k Keeper) SetVolumeReport(ctx sdk.Context, reporter sdk.AccAddress, reportReference string) {
//	store := ctx.KVStore(k.storeKey)
//	storeKey := types.VolumeReportStoreKey(reporter)
//	store.Set(storeKey, []byte(reportReference))
//}

func (k Keeper) SetVolumeReport(ctx sdk.Context, epoch sdk.Int, reportReference string) {
	store := ctx.KVStore(k.storeKey)
	//storeKey := types.VolumeReportStoreKey(reporter)
	storeKey := types.VolumeReportStoreKey(epoch)
	store.Set(storeKey, []byte(reportReference))
}

func (k Keeper) DeleteVolumeReport(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
}

func (k Keeper) IsSPNode(ctx sdk.Context, addr sdk.AccAddress) (found bool) {
	_, found = k.RegisterKeeper.GetIndexingNode(ctx, addr)
	return found
}
