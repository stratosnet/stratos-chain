package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stratosnet/stratos-chain/x/pot"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/stratosnet/stratos-chain/x/sds/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper encodes/decodes files using the go-amino (binary)
// encoding/decoding library.
type Keeper struct {
	key            sdk.StoreKey
	cdc            *codec.Codec
	BankKeeper     bank.Keeper
	RegisterKeeper register.Keeper
	PotKeeper      pot.Keeper
}

// NewKeeper returns a new sdk.NewKeeper that uses go-amino to
// (binary) encode and decode concrete sdk.MsgUploadFile.
// nolint
func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	bankKeeper bank.Keeper,
	registerKeeper register.Keeper,
	potKeeper pot.Keeper,
) Keeper {
	return Keeper{
		key:            key,
		cdc:            cdc,
		BankKeeper:     bankKeeper,
		RegisterKeeper: registerKeeper,
		PotKeeper:      potKeeper,
	}
}

// Logger returns a module-specific logger.
func (fk Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetFileInfoBytesByFileHash Returns the hash of file
func (fk Keeper) GetFileInfoBytesByFileHash(ctx sdk.Context, key []byte) ([]byte, error) {
	store := ctx.KVStore(fk.key)
	bz := store.Get(types.FileStoreKey(key))
	if bz == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "FileHash %s does not exist", hex.EncodeToString(types.FileStoreKey(key))[2:])
	}
	return bz, nil
}

// SetFileHash Sets sender-fileHash KV pair
func (fk Keeper) SetFileHash(ctx sdk.Context, fileHash []byte, fileInfo types.FileInfo) {
	store := ctx.KVStore(fk.key)
	storeKey := types.FileStoreKey(fileHash)
	bz := types.MustMarshalFileInfo(fk.cdc, fileInfo)
	store.Set(storeKey, bz)
}

// [S] is the initial genesis deposit by all Resource Nodes and Meta Nodes at t=0
// The current unissued prepay Volume Pool [Pt] is the total remaining prepay STOS kept by the Stratos Network but not yet issued to Resource Nodes as rewards.
// The remaining total Ozone limit [Lt] is the upper bound of the total Ozone that users can purchase from the Stratos blockchain.
// [X] is the total amount of STOS token prepaid by user at time t
// the total amount of Ozone the user gets = Lt * X / (S + Pt + X)
func (fk Keeper) purchaseUoz(ctx sdk.Context, amount sdk.Int) sdk.Int {
	S := fk.RegisterKeeper.GetInitialGenesisStakeTotal(ctx)
	Pt := fk.PotKeeper.GetTotalUnissuedPrepay(ctx)
	Lt := fk.RegisterKeeper.GetRemainingOzoneLimit(ctx)

	purchased := Lt.ToDec().
		Mul(amount.ToDec()).
		Quo((S.
			Add(Pt).
			Add(amount)).ToDec()).
		TruncateInt()

	// update total unissued prepay
	newTotalUnissuedPrepay := Pt.Add(amount)
	fk.PotKeeper.SetTotalUnissuedPrepay(ctx, newTotalUnissuedPrepay)

	// update remaining uoz limit
	newRemainingOzoneLimit := Lt.Sub(purchased)
	fk.RegisterKeeper.SetRemainingOzoneLimit(ctx, newRemainingOzoneLimit)

	return purchased
}

func (fk Keeper) simulatePurchaseUoz(ctx sdk.Context, amount sdk.Int) sdk.Int {
	S := fk.RegisterKeeper.GetInitialGenesisStakeTotal(ctx)
	Pt := fk.PotKeeper.GetTotalUnissuedPrepay(ctx)
	Lt := fk.RegisterKeeper.GetRemainingOzoneLimit(ctx)
	purchased := Lt.ToDec().
		Mul(amount.ToDec()).
		Quo((S.
			Add(Pt).
			Add(amount)).ToDec()).
		TruncateInt()
	return purchased
}

// calc current uoz price
func (fk Keeper) currUozPrice(ctx sdk.Context) sdk.Dec {
	S := fk.RegisterKeeper.GetInitialGenesisStakeTotal(ctx)
	Pt := fk.PotKeeper.GetTotalUnissuedPrepay(ctx)
	Lt := fk.RegisterKeeper.GetRemainingOzoneLimit(ctx)
	currUozPrice := (S.Add(Pt)).ToDec().
		Quo(Lt.ToDec())
	return currUozPrice
}

// calc remaining/total supply for uoz
func (fk Keeper) uozSupply(ctx sdk.Context) (remaining, total sdk.Int) {
	remaining = fk.RegisterKeeper.GetRemainingOzoneLimit(ctx)
	// TODO create a dedicated storeKey in pot module for total ozone supply and keep updating it along with related operations (like AddResourceNodeStake)
	// fake return of total
	total, _ = sdk.NewIntFromString("-1")
	return remaining, total
}

// Prepay transfers coins from bank to sds (volumn) pool
func (fk Keeper) Prepay(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) (sdk.Int, error) {
	// src - hasCoins?
	if !fk.BankKeeper.HasCoins(ctx, sender, coins) {
		return sdk.ZeroInt(), sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "No valid coins to be deducted from acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
	}
	adjustedCoins := coins
	for _, coin := range adjustedCoins {
		switch coin.Denom {
		case types.DefaultRewardDenom:
			coin.Amount = coin.Amount.Mul(types.RewardToUstos)
		case types.DefaultBondDenom:
		default:
			return sdk.ZeroInt(), sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "Prepay contains unknown coins %s: ", coin.String())
		}
	}

	err := fk.doPrepay(ctx, sender, adjustedCoins)
	if err != nil {
		return sdk.ZeroInt(), sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "Failed prepay from acc %s", hex.EncodeToString(types.PrepayBalanceKey(sender)))
	}

	_, err = fk.BankKeeper.SubtractCoins(ctx, sender, coins)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	//TODO: move the definition of default denomination to params.go
	prepay := sdk.ZeroInt()
	purchased := sdk.ZeroInt()
	for _, coin := range adjustedCoins {
		prepay = prepay.Add(coins.AmountOf(coin.Denom))
	}
	purchased = purchased.Add(fk.purchaseUoz(ctx, prepay))
	return purchased, nil
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
