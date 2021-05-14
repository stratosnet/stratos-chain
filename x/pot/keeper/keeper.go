package keeper

import (
	"encoding/hex"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// Keeper of the pot store
type Keeper struct {
	BankKeeper bank.Keeper
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	// paramspace types.ParamSubspace
}

// NewKeeper creates a pot keeper
func NewKeeper(bankKeeper bank.Keeper, cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		BankKeeper: bankKeeper,
		storeKey:   key,
		cdc:        cdc,
		// paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
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
