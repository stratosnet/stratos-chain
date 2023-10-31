package keeper

import (
	"github.com/kelindar/bitmap"

	"github.com/cometbft/cometbft/libs/log"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stratos "github.com/stratosnet/stratos-chain/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
	"github.com/stratosnet/stratos-chain/x/sds/types"
)

// Keeper encodes/decodes files using the go-amino (binary)
// encoding/decoding library.
type Keeper struct {
	storeKey       storetypes.StoreKey
	cdc            codec.Codec
	bankKeeper     types.BankKeeper
	registerKeeper types.RegisterKeeper
	potKeeper      types.PotKeeper

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper returns a new sdk.NewKeeper that uses go-amino to
// (binary) encode and decode concrete sdk.MsgUploadFile.
// nolint
func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
	registerKeeper types.RegisterKeeper,
	potKeeper types.PotKeeper,
	authority string,
) Keeper {
	return Keeper{
		storeKey:       storeKey,
		cdc:            cdc,
		bankKeeper:     bankKeeper,
		registerKeeper: registerKeeper,
		potKeeper:      potKeeper,
		authority:      authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) FileUpload(ctx sdk.Context, fileHash string, reporter stratos.SdsAddress, reporterOwner, uploader sdk.AccAddress) (err error) {
	if !(k.registerKeeper.OwnMetaNode(ctx, reporterOwner, reporter)) {
		return types.ErrReporterAddressOrOwner
	}

	var fileUploadReporters bitmap.Bitmap
	// query exist fileInfo which sent by other meta node
	fileInfo, found := k.GetFileInfoByFileHash(ctx, []byte(fileHash))
	if !found {
		fileUploadReporters = bitmap.Bitmap{}
	} else {
		fileUploadReporters = bitmap.FromBytes(fileInfo.GetReporters())
	}
	reporterIndex, err := k.registerKeeper.GetMetaNodeBitMapIndex(ctx, reporter)
	fileUploadReporters.Set(uint32(reporterIndex))
	height := sdkmath.NewInt(ctx.BlockHeight())

	newFileInfo := types.NewFileInfo(height, fileUploadReporters.ToBytes(), uploader.String())

	k.SetFileInfo(ctx, []byte(fileHash), newFileInfo)

	return nil
}

// [S] is the initial genesis deposit by all Resource Nodes and Meta Nodes at t=0
// The current unissued prepay Volume Pool [Pt] is the total remaining prepay STOS kept by the Stratos Network but not yet issued to Resource Nodes as rewards.
// The remaining total Ozone limit [Lt] is the upper bound of the total Ozone that users can purchase from the Stratos blockchain.
// [X] is the total amount of STOS token prepaid by user at time t
// the total amount of Ozone the user gets = Lt * X / (S + Pt + X)
func (k Keeper) purchaseNozAndSubCoins(ctx sdk.Context, from sdk.AccAddress, amount sdkmath.Int) (sdkmath.Int, error) {
	St := k.registerKeeper.GetEffectiveTotalDeposit(ctx)
	Pt := k.registerKeeper.GetTotalUnissuedPrepay(ctx).Amount
	Lt := k.registerKeeper.GetRemainingOzoneLimit(ctx)

	purchased := Lt.ToLegacyDec().
		Mul(amount.ToLegacyDec()).
		Quo((St.
			Add(Pt).
			Add(amount)).ToLegacyDec()).
		TruncateInt()

	if purchased.GT(Lt) {
		return sdkmath.ZeroInt(), types.ErrOzoneLimitNotEnough
	}
	// send coins to total unissued prepay pool
	prepayAmt := sdk.NewCoin(k.BondDenom(ctx), amount)
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, registertypes.TotalUnissuedPrepay, sdk.NewCoins(prepayAmt))
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	// update remaining noz limit
	newRemainingOzoneLimit := Lt.Sub(purchased)
	k.registerKeeper.SetRemainingOzoneLimit(ctx, newRemainingOzoneLimit)

	return purchased, nil
}

func (k Keeper) simulatePurchaseNoz(ctx sdk.Context, coins sdk.Coins) sdkmath.Int {
	amount := coins.AmountOf(k.BondDenom(ctx))

	St := k.registerKeeper.GetEffectiveTotalDeposit(ctx)
	Pt := k.registerKeeper.GetTotalUnissuedPrepay(ctx).Amount
	Lt := k.registerKeeper.GetRemainingOzoneLimit(ctx)
	purchased := Lt.ToLegacyDec().
		Mul(amount.ToLegacyDec()).
		Quo((St.
			Add(Pt).
			Add(amount)).ToLegacyDec()).
		TruncateInt()
	return purchased
}

// Prepay transfers coins from bank to sds (volume) pool
func (k Keeper) Prepay(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) (sdkmath.Int, error) {
	validPrepayAmt := coins.AmountOf(k.BondDenom(ctx))
	hasCoin := k.bankKeeper.HasBalance(ctx, sender, sdk.NewCoin(k.BondDenom(ctx), validPrepayAmt))
	if !hasCoin {
		return sdkmath.ZeroInt(), errors.Wrapf(sdkerrors.ErrInsufficientFunds, "Insufficient balance in the acc %s", sender.String())
	}

	return k.purchaseNozAndSubCoins(ctx, sender, validPrepayAmt)
}
