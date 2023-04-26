package keeper

import (
	"bytes"
	"fmt"
	"math"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// Keeper of the pot store
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              codec.Codec
	paramSpace       paramstypes.Subspace
	feeCollectorName string // name of the FeeCollector ModuleAccount
	bankKeeper       types.BankKeeper
	accountKeeper    types.AccountKeeper
	stakingKeeper    types.StakingKeeper
	registerKeeper   types.RegisterKeeper
	distrKeeper      types.DistrKeeper
}

// NewKeeper creates a pot keeper
func NewKeeper(cdc codec.Codec, key sdk.StoreKey, paramSpace paramstypes.Subspace, feeCollectorName string,
	bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper, stakingKeeper types.StakingKeeper,
	registerKeeper types.RegisterKeeper, distrKeeper types.DistrKeeper,
) Keeper {
	keeper := Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace.WithKeyTable(types.ParamKeyTable()),
		feeCollectorName: feeCollectorName,
		bankKeeper:       bankKeeper,
		accountKeeper:    accountKeeper,
		stakingKeeper:    stakingKeeper,
		registerKeeper:   registerKeeper,
		distrKeeper:      distrKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) VolumeReport(ctx sdk.Context, walletVolumes types.WalletVolumes, reporter stratos.SdsAddress,
	epoch sdk.Int, reportReference string, txHash string) (err error) {

	// tx fail when unhandled epoch > 0
	if k.GetUnhandledEpoch(ctx).GT(sdk.ZeroInt()) {
		return types.ErrRewardDistributionNotComplete
	}

	//record volume report
	reportRecord := types.NewReportRecord(reporter, reportReference, txHash)
	k.SetVolumeReport(ctx, epoch, reportRecord)

	// save for reward distribution in the EndBlock at height + 1
	k.SetIsReadyToDistributeReward(ctx, false)
	k.SetUnhandledEpoch(ctx, epoch)
	k.SetUnhandledReport(ctx, walletVolumes)

	return nil
}

func (k Keeper) HasReachedThreshold(ctx sdk.Context, pubKeys [][]byte) bool {
	totalMetaNodes := k.registerKeeper.GetBondedMetaNodeCnt(ctx).Int64()
	signedMetaNodes := len(pubKeys)

	threshold := int(math.Max(1, math.Floor(float64(totalMetaNodes)*2/3)))

	return signedMetaNodes >= threshold
}

func (k Keeper) FoundationDeposit(ctx sdk.Context, amount sdk.Coins, from sdk.AccAddress) (err error) {
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.FoundationAccount, amount)
	if err != nil {
		return err
	}
	return nil
}

// Restore total supply to 100M stos
func (k Keeper) RestoreTotalSupply(ctx sdk.Context) (minter sdk.AccAddress, mintCoins sdk.Coins) {
	// reset total supply to 100M stos
	events := ctx.EventManager().Events()
	attrKeyAmtBytes := []byte(sdk.AttributeKeyAmount)

	totalBurnedCoins := sdk.Coins{}
	for _, event := range events {
		if event.Type == banktypes.EventTypeCoinBurn {
			attributes := event.Attributes
			for _, attr := range attributes {
				if bytes.Equal(attr.Key, attrKeyAmtBytes) {
					amount, err := sdk.ParseCoinsNormalized(string(attr.Value))
					if err != nil {
						ctx.Logger().Error("An error occurred while parsing burned amount. ", "ErrMsg", err.Error())
						break
					}
					totalBurnedCoins = totalBurnedCoins.Add(amount...)
				}
			}
		}
	}

	MaxTotalSupply := sdk.NewInt(types.DefaultTotalSupply)
	currentTotalSupply := k.bankKeeper.GetSupply(ctx, k.BondDenom(ctx)).Amount
	totalBurned := totalBurnedCoins.AmountOf(k.BondDenom(ctx))

	if totalBurned.Add(currentTotalSupply).GT(MaxTotalSupply) {
		mintCoins = sdk.NewCoins(
			sdk.NewCoin(k.BondDenom(ctx), MaxTotalSupply.Sub(currentTotalSupply)),
		)
	} else {
		mintCoins = totalBurnedCoins
	}

	// mint coins
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins)
	if err != nil {
		ctx.Logger().Error("Restore total supply failed:", err.Error())
		return sdk.AccAddress{}, sdk.Coins{}
	}

	// send new mint coins to community pool
	senderAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	err = k.distrKeeper.FundCommunityPool(ctx, mintCoins, senderAddr)
	if err != nil {
		ctx.Logger().Error("Restore total supply failed:", err.Error())
		return sdk.AccAddress{}, sdk.Coins{}
	}

	return
}
