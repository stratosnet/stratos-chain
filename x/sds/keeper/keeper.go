package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	potKeeper "github.com/stratosnet/stratos-chain/x/pot/keeper"
	registerKeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
	"github.com/stratosnet/stratos-chain/x/sds/types"
)

// Keeper encodes/decodes files using the go-amino (binary)
// encoding/decoding library.
type Keeper struct {
	key            sdk.StoreKey
	cdc            codec.Codec
	paramSpace     paramtypes.Subspace
	bankKeeper     bankKeeper.Keeper
	registerKeeper registerKeeper.Keeper
	potKeeper      potKeeper.Keeper
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
		registerKeeper: registerKeeper,
		potKeeper:      potKeeper,
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
func (k Keeper) purchaseNozAndSubCoins(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int) (sdk.Int, error) {
	St := k.registerKeeper.GetEffectiveTotalStake(ctx)
	Pt := k.registerKeeper.GetTotalUnissuedPrepay(ctx).Amount
	Lt := k.registerKeeper.GetRemainingOzoneLimit(ctx)

	purchased := Lt.ToDec().
		Mul(amount.ToDec()).
		Quo((St.
			Add(Pt).
			Add(amount)).ToDec()).
		TruncateInt()

	// send coins to total unissued prepay pool
	prepayAmt := sdk.NewCoin(k.BondDenom(ctx), amount)
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, registertypes.TotalUnissuedPrepay, sdk.NewCoins(prepayAmt))
	if err != nil {
		return sdk.ZeroInt(), err
	}

	// update remaining noz limit
	newRemainingOzoneLimit := Lt.Sub(purchased)
	k.registerKeeper.SetRemainingOzoneLimit(ctx, newRemainingOzoneLimit)

	return purchased, nil
}

func (k Keeper) simulatePurchaseNoz(ctx sdk.Context, coins sdk.Coins) sdk.Int {
	amount := coins.AmountOf(k.BondDenom(ctx))

	St := k.registerKeeper.GetEffectiveTotalStake(ctx)
	Pt := k.registerKeeper.GetTotalUnissuedPrepay(ctx).Amount
	Lt := k.registerKeeper.GetRemainingOzoneLimit(ctx)
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
	validPrepayAmt := coins.AmountOf(k.BondDenom(ctx))
	hasCoin := k.bankKeeper.HasBalance(ctx, sender, sdk.NewCoin(k.BondDenom(ctx), validPrepayAmt))
	if !hasCoin {
		return sdk.ZeroInt(), sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "Insufficient balance in the acc %s", sender.String())
	}

	return k.purchaseNozAndSubCoins(ctx, sender, validPrepayAmt)
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
