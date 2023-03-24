package keeper

import (
	"fmt"

	"github.com/kelindar/bitmap"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	stratos "github.com/stratosnet/stratos-chain/types"
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
	height := sdk.NewInt(ctx.BlockHeight())

	newFileInfo := types.NewFileInfo(height, fileUploadReporters.ToBytes(), uploader.String())

	k.SetFileInfo(ctx, []byte(fileHash), newFileInfo)

	return nil
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
