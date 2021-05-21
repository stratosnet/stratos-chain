package keeper

import (
	"encoding/hex"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// Keeper of the pot store
type Keeper struct {
	BankKeeper     bank.Keeper
	storeKey       sdk.StoreKey
	cdc            *codec.Codec
	RegisterKeeper *register.Keeper
	//paramspace types.ParamSubspace
}

// NewKeeper creates a pot keeper
func NewKeeper(bankKeeper bank.Keeper, cdc *codec.Codec, key sdk.StoreKey, registerKeeper *register.Keeper) Keeper {
	keeper := Keeper{
		BankKeeper:     bankKeeper,
		storeKey:       key,
		cdc:            cdc,
		RegisterKeeper: registerKeeper,
		//paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
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

func (k Keeper) SetVolumeReport(ctx sdk.Context, reporter sdk.AccAddress, reportReference string) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.VolumeReportStoreKey(reporter)
	store.Set(storeKey, []byte(reportReference))
}

func (k Keeper) DeleteVolumeReport(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
}

func (k Keeper) IsIndexingNode(ctx sdk.Context, addr sdk.AccAddress) (found bool) {
	//ctx.Logger().Info("Inside IsIndexingNode start")
	//ctx.Logger().Info("ctx in IsIndexingNode:" + string(types.ModuleCdc.MustMarshalJSON(ctx)))
	//ctx.Logger().Info("addr:" + string(types.ModuleCdc.MustMarshalJSON(addr)))

	_, found = k.RegisterKeeper.GetIndexingNode(ctx, addr)
	//ctx.Logger().Info("Inside IsIndexingNode end")
	return found
}
