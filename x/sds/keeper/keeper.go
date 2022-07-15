package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	potKeeper "github.com/stratosnet/stratos-chain/x/pot/keeper"
	registerKeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
	"github.com/stratosnet/stratos-chain/x/sds/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	RewardToUstos = sdk.NewInt(1)
)

// Keeper encodes/decodes files using the go-amino (binary)
// encoding/decoding library.
type Keeper struct {
	key            sdk.StoreKey
	cdc            codec.Codec
	paramSpace     paramtypes.Subspace
	bankKeeper     bankKeeper.Keeper
	RegisterKeeper registerKeeper.Keeper
	PotKeeper      potKeeper.Keeper
}

// NewKeeper returns a new sdk.NewKeeper that uses go-amino to
// (binary) encode and decode concrete sdk.MsgUploadFile.
// nolint
func NewKeeper(
	cdc codec.Codec,
	key sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	bankKeeper bankKeeper.Keeper,
	registerKeeper registerKeeper.Keeper,
	potKeeper potKeeper.Keeper,
) Keeper {
	return Keeper{
		key:            key,
		cdc:            cdc,
		paramSpace:     paramSpace.WithKeyTable(types.ParamKeyTable()),
		bankKeeper:     bankKeeper,
		RegisterKeeper: registerKeeper,
		PotKeeper:      potKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetFileInfoBytesByFileHash Returns the hash of file
func (k Keeper) GetFileInfoBytesByFileHash(ctx sdk.Context, key []byte) ([]byte, error) {
	store := ctx.KVStore(k.key)
	bz := store.Get(types.FileStoreKey(key))
	if bz == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "FileHash %s does not exist", hex.EncodeToString(types.FileStoreKey(key))[2:])
	}
	return bz, nil
}

// SetFileHash Sets sender-fileHash KV pair
func (k Keeper) SetFileHash(ctx sdk.Context, fileHash []byte, fileInfo types.FileInfo) {
	store := ctx.KVStore(k.key)
	storeKey := types.FileStoreKey(fileHash)
	bz := types.MustMarshalFileInfo(k.cdc, fileInfo)
	store.Set(storeKey, bz)
}

// [S] is the initial genesis deposit by all Resource Nodes and Meta Nodes at t=0
// The current unissued prepay Volume Pool [Pt] is the total remaining prepay STOS kept by the Stratos Network but not yet issued to Resource Nodes as rewards.
// The remaining total Ozone limit [Lt] is the upper bound of the total Ozone that users can purchase from the Stratos blockchain.
// [X] is the total amount of STOS token prepaid by user at time t
// the total amount of Ozone the user gets = Lt * X / (S + Pt + X)
func (k Keeper) purchaseUozAndSubCoins(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int) sdk.Int {
	St := k.RegisterKeeper.GetEffectiveGenesisStakeTotal(ctx)
	//S := k.RegisterKeeper.GetInitialGenesisStakeTotal(ctx)
	Pt := k.RegisterKeeper.GetTotalUnissuedPrepay(ctx).Amount
	Lt := k.RegisterKeeper.GetRemainingOzoneLimit(ctx)

	purchased := Lt.ToDec().
		Mul(amount.ToDec()).
		Quo((St.
			Add(Pt).
			Add(amount)).ToDec()).
		TruncateInt()

	// send coins to total unissued prepay pool
	err := k.RegisterKeeper.SendCoinsFromAccount2TotalUnissuedPrepayPool(ctx, from, sdk.NewCoin(k.BondDenom(ctx), amount))
	if err != nil {
		return sdk.ZeroInt()
	}
	//newTotalUnissuedPrepay := Pt.Add(amount)
	//k.RegisterKeeper.SetTotalUnissuedPrepay(ctx, sdk.NewCoin(k.BondDenom(ctx), newTotalUnissuedPrepay))

	// update remaining uoz limit
	newRemainingOzoneLimit := Lt.Sub(purchased)
	k.RegisterKeeper.SetRemainingOzoneLimit(ctx, newRemainingOzoneLimit)

	return purchased
}

func (k Keeper) simulatePurchaseUoz(ctx sdk.Context, amount sdk.Int) sdk.Int {
	St := k.RegisterKeeper.GetEffectiveGenesisStakeTotal(ctx)
	//S := k.RegisterKeeper.GetInitialGenesisStakeTotal(ctx)
	Pt := k.RegisterKeeper.GetTotalUnissuedPrepay(ctx).Amount
	Lt := k.RegisterKeeper.GetRemainingOzoneLimit(ctx)
	purchased := Lt.ToDec().
		Mul(amount.ToDec()).
		Quo((St.
			Add(Pt).
			Add(amount)).ToDec()).
		TruncateInt()
	return purchased
}

// Prepay transfers coins from bank to sds (volumn) pool
func (k Keeper) Prepay(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) (sdk.Int, error) {
	// src - hasCoins?
	//if !k.bankKeeper.HasCoins(ctx, sender, coins) {
	//	return sdk.ZeroInt(), sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "No valid coins to be deducted from acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
	//}
	for _, coin := range coins {
		hasCoin := k.bankKeeper.HasBalance(ctx, sender, coin)
		if !hasCoin {
			return sdk.ZeroInt(), sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "No valid coins to be deducted from acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
		}
	}

	prepay := coins.AmountOf(k.BondDenom(ctx))
	purchased := k.purchaseUozAndSubCoins(ctx, sender, prepay)

	err := k.doPrepay(ctx, sender, coins)
	if err != nil {
		return sdk.ZeroInt(), sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "Failed prepay from acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
	}
	//// sub coins from sender's wallet
	//err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, regtypes.TotalUnissuedPrepayName, coins)
	//if err != nil {
	//	return sdk.ZeroInt(), err
	//}
	//
	//_, err = k.bankKeeper.SubtractCoins(ctx, sender, coins)
	//if err != nil {
	//	return sdk.ZeroInt(), err
	//}

	return purchased, nil
}

// HasPrepay Returns bool indicating if the sender did prepay before
func (k Keeper) hasPrepay(ctx sdk.Context, sender sdk.AccAddress) bool {
	store := ctx.KVStore(k.key)
	return store.Has(types.PrepayBalanceKey(sender))
}

// GetPrepay Returns the existing prepay coins
func (k Keeper) GetPrepay(ctx sdk.Context, sender sdk.AccAddress) (sdk.Int, error) {
	store := ctx.KVStore(k.key)
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

// GetPrepayBytes returns bytearr of the existing prepay coins
func (k Keeper) GetPrepayBytes(ctx sdk.Context, sender sdk.AccAddress) ([]byte, error) {
	store := ctx.KVStore(k.key)
	storeValue := store.Get(types.PrepayBalanceKey(sender))
	if storeValue == nil {
		return []byte{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "No prepayment info for acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
	}
	return storeValue, nil
}

// SetPrepay Sets init coins
func (k Keeper) SetPrepay(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) error {
	store := ctx.KVStore(k.key)
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
func (k Keeper) appendPrepay(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) error {
	store := ctx.KVStore(k.key)
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

func (k Keeper) doPrepay(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) error {
	// tar - has key?
	if k.hasPrepay(ctx, sender) {
		// has key - append coins
		return k.appendPrepay(ctx, sender, coins)
	}
	// doesn't have key - create new
	return k.SetPrepay(ctx, sender, coins)
}

// IterateFileUpload Iterate over all uploaded files.
// Iteration for all uploaded files
func (k Keeper) IterateFileUpload(ctx sdk.Context, handler func(string, types.FileInfo) (stop bool)) {
	store := ctx.KVStore(k.key)
	iter := sdk.KVStorePrefixIterator(store, types.FileStoreKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		fileHash := string(iter.Key()[len(types.FileStoreKeyPrefix):])
		var fileInfo types.FileInfo
		k.cdc.MustUnmarshal(iter.Value(), &fileInfo)
		if handler(fileHash, fileInfo) {
			break
		}
	}
}

// IteratePrepay Iterate over all prepay KVs.
// Iteration for all prepay KVs
func (k Keeper) IteratePrepay(ctx sdk.Context, handler func(sdk.AccAddress, sdk.Int) (stop bool)) {
	store := ctx.KVStore(k.key)
	iter := sdk.KVStorePrefixIterator(store, types.PrepayBalancePrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		senderAddr := sdk.AccAddress(iter.Key()[len(types.PrepayBalancePrefix):])
		var amt sdk.Int
		err := amt.UnmarshalJSON(iter.Value())
		if err != nil {
			panic("invalid prepay amount")
		}
		if handler(senderAddr, amt) {
			break
		}
	}
}
