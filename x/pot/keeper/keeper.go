package keeper

import (
	"encoding/hex"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/stratosnet/stratos-chain/x/sds"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// Keeper of the pot store
type Keeper struct {
	bankKeeper       bank.Keeper
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	paramSpace       params.Subspace
	feeCollectorName string // name of the FeeCollector ModuleAccount
	distrKeeper      distribution.Keeper
	sdsKeeper        sds.Keeper
	supplyKeeper     supply.Keeper
	accountKeeper    auth.AccountKeeper
	stakingKeeper    staking.Keeper
	registerKeeper   register.Keeper
}

// NewKeeper creates a pot keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, feeCollectorName string,
	bankKeeper bank.Keeper, distrKeeper distribution.Keeper, sdsKeeper sds.Keeper, supplyKeeper supply.Keeper,
	accountKeeper auth.AccountKeeper, stakingKeeper staking.Keeper, registerKeeper register.Keeper,
) Keeper {
	keeper := Keeper{
		bankKeeper:     bankKeeper,
		storeKey:       key,
		cdc:            cdc,
		paramSpace:     paramSpace.WithKeyTable(types.ParamKeyTable()),
		feeCollectorName:feeCollectorName,
		distrKeeper:    distrKeeper,
		sdsKeeper:      sdsKeeper,
		supplyKeeper:   supplyKeeper,
		accountKeeper:  accountKeeper,
		stakingKeeper:  stakingKeeper,
		registerKeeper: registerKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetVolumeReport returns the hash of volume report
func (k Keeper) GetVolumeReport(ctx sdk.Context, reporter sdk.AccAddress) ([]byte, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VolumeReportStoreKey(reporter))
	if bz == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress,
			"key %s does not exist", hex.EncodeToString(types.VolumeReportStoreKey(reporter)))
	}
	return bz, nil
}

func (k Keeper) SetVolumeReport(ctx sdk.Context, volumeReport *types.MsgVolumeReport) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.VolumeReportStoreKey(volumeReport.Reporter)
	store.Set(storeKey, []byte(volumeReport.ReportReference))
}

func (k Keeper) DeleteVolumeReport(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
}
