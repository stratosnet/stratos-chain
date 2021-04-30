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

// GetVolumeReportHash returns the hash of volume report
func (k Keeper) GetVolumeReportHash(ctx sdk.Context, reporter sdk.AccAddress) ([]byte, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VolumeReportStoreKey(reporter))
	if bz == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress,
			"key %s does not exist", hex.EncodeToString(types.VolumeReportStoreKey(reporter)))
	}
	return bz, nil
}

func (k Keeper) SetVolumeReportHash(ctx sdk.Context, reporter sdk.AccAddress, reportReferenceHash string) {
	store := ctx.KVStore(k.storeKey)
	sKey := types.VolumeReportStoreKey(reporter)
	store.Set(sKey, []byte(reportReferenceHash))
	ctx.Logger().Info("SetVolumeReportHash Key " + "------" + string(types.VolumeReportStoreKeyPrefix) + reporter.String())
	ctx.Logger().Info("SetVolumeReportHash Value" + reportReferenceHash)
}

func (k Keeper) DeleteVolumeReportHash(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
}

func (k Keeper) GetSingleNodeVolume(ctx sdk.Context, nodeAddress []byte) ([]byte, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SingleNodeVolumeStoreKey(nodeAddress))
	if bz == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress,
			"key %s does not exist", hex.EncodeToString(types.VolumeReportStoreKey(nodeAddress)))
	}
	//r, _ := strconv.ParseInt(string(bz), 10, 64)

	//res, err =
	//if err != nil {
	//	return nil, fmt.Errorf("the value of key %s cannot be converted to int64",
	//		types.SingleNodeVolumeStoreKey(nodeAddress))
	//}
	return bz, nil
}

func (k Keeper) SetSingleNodeVolume(ctx sdk.Context, singleNodeVolume types.SingleNodeVolume) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.VolumeReportStoreKey(singleNodeVolume.NodeAddress)
	bz := []byte(singleNodeVolume.Volume.String())
	store.Set(storeKey, bz)
	ctx.Logger().Info("node_volume: ", fmt.Sprintf("%s---%s", hex.EncodeToString(storeKey), bz))
}
