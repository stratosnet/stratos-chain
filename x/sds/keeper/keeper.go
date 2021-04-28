package keeper

import (
	"encoding/hex"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/sds/types"
)

// Keeper encodes/decodes files using the go-amino (binary)
// encoding/decoding library.
type Keeper struct {
	CoinKeeper bank.Keeper
	key        sdk.StoreKey
	cdc        *codec.Codec
}

// NewKeeper returns a new sdk.NewKeeper that uses go-amino to
// (binary) encode and decode concrete sdk.MsgUploadFile.
// nolint
func NewKeeper(
	coinKeeper bank.Keeper,
	cdc *codec.Codec,
	key sdk.StoreKey,
) Keeper {
	return Keeper{
		CoinKeeper: coinKeeper,
		key:        key,
		cdc:        cdc,
	}
}

// Logger returns a module-specific logger.
func (fk Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetFileHash Returns the hash of file
func (fk Keeper) GetFileHash(ctx sdk.Context, key []byte) ([]byte, error) {
	store := ctx.KVStore(fk.key)
	bz := store.Get(types.FileStoreKey(key))
	if bz == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "key %s does not exist", hex.EncodeToString(types.FileStoreKey(key)))
	}
	return bz, nil
}

// SetFileHash Sets sender-fileHash KV pair
func (fk Keeper) SetFileHash(ctx sdk.Context, sender []byte, fileHash []byte) {
	store := ctx.KVStore(fk.key)
	storeKey := types.FileStoreKey(sender)
	store.Set(storeKey, fileHash)
}
