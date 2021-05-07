package keeper

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stratosnet/stratos-chain/x/sds/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper encodes/decodes files using the go-amino (binary)
// encoding/decoding library.
type Keeper struct {
	BankKeeper bank.Keeper
	key        sdk.StoreKey
	cdc        *codec.Codec
}

// NewKeeper returns a new sdk.NewKeeper that uses go-amino to
// (binary) encode and decode concrete sdk.MsgUploadFile.
// nolint
func NewKeeper(
	bankKeeper bank.Keeper,
	cdc *codec.Codec,
	key sdk.StoreKey,
) Keeper {
	return Keeper{
		BankKeeper: bankKeeper,
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
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "FileHash %s does not exist", hex.EncodeToString(types.FileStoreKey(key)))
	}
	return bz, nil
}

// SetFileHash Sets sender-fileHash KV pair
func (fk Keeper) SetFileHash(ctx sdk.Context, fileHash []byte, height []byte) {
	store := ctx.KVStore(fk.key)
	storeKey := types.FileStoreKey(fileHash)
	store.Set(storeKey, height)
}

// Prepay transfers coins from bank to sds (volumn) pool
func (fk Keeper) Prepay(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) error {
	// src - hasCoins?
	if fk.BankKeeper.HasCoins(ctx, sender, coins) {
		err := fk.doPrepay(ctx, sender, coins)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "Failed prepay from acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
		}
		_, error := fk.BankKeeper.SubtractCoins(ctx, sender, coins)
		return error
	}
	return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "No valid coins to be deducted from acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
}

// HasPrepay Returns bool indicating if the sender did prepay before
func (fk Keeper) hasPrepay(ctx sdk.Context, sender sdk.AccAddress) bool {
	store := ctx.KVStore(fk.key)
	return store.Has(types.PrepayBalanceKey(sender))
}

// GetPrepay Returns the existing prepay coins
func (fk Keeper) GetPrepay(ctx sdk.Context, sender sdk.AccAddress) (sdk.Int, error) {
	store := ctx.KVStore(fk.key)
	storeValue := store.Get(types.PrepayBalanceKey(sender))
	if storeValue == nil {
		return sdk.Int{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "No prepayment info for acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
	}
	var prepaidBalance sdk.Int
	err := prepaidBalance.UnmarshalJSON(storeValue)
	if err != nil {
		return sdk.Int{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "Unmarshal failed for acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
	}
	return prepaidBalance, nil
}

// GetPrepay Returns bytearr of the existing prepay coins
func (fk Keeper) GetPrepayBytes(ctx sdk.Context, sender sdk.AccAddress) ([]byte, error) {
	store := ctx.KVStore(fk.key)
	storeValue := store.Get(types.PrepayBalanceKey(sender))
	if storeValue == nil {
		return []byte{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "No prepayment info for acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
	}
	return storeValue, nil
}

// SetPrepay Sets init coins
func (fk Keeper) setPrepay(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) error {
	store := ctx.KVStore(fk.key)
	storeKey := types.PrepayBalanceKey(sender)
	balance := sdk.NewInt(0)
	for _, coin := range coins {
		balance = balance.Add(coin.Amount)
	}
	storeValue, err := balance.MarshalJSON()
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "Marshalling failed for acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
	}
	store.Set(storeKey, storeValue)
	return nil
}

// AppendPrepay adds more coins to existing coins
func (fk Keeper) appendPrepay(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) error {
	store := ctx.KVStore(fk.key)
	storeKey := types.PrepayBalanceKey(sender)
	storeValue := store.Get(storeKey)
	var prepaidBalance sdk.Int
	err := prepaidBalance.UnmarshalJSON(storeValue)
	if err != nil {
		//prepaidBalance = sdk.NewInt(0)
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "Unmarshal failed for acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
	}
	for _, coin := range coins {
		prepaidBalance = prepaidBalance.Add(coin.Amount)
	}
	newStoreValue, err := prepaidBalance.MarshalJSON()
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "Marshal failed for acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
	}
	store.Set(storeKey, newStoreValue)
	return nil
}

func (fk Keeper) doPrepay(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) error {
	// tar - has key?
	if fk.hasPrepay(ctx, sender) {
		// has key - append coins
		return fk.appendPrepay(ctx, sender, coins)
	}
	// doesn't have key - create new
	return fk.setPrepay(ctx, sender, coins)
}
