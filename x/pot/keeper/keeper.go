package keeper

import (
	"errors"
	"fmt"
	"math"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
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

	//record volume report
	reportRecord := types.NewReportRecord(reporter, reportReference, txHash)
	k.SetVolumeReport(ctx, epoch, reportRecord)

	err = k.DistributePotReward(ctx, walletVolumes.GetVolumes(), epoch)
	if err != nil {
		return err
	}

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

func (k Keeper) SafeMintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	denom := k.BondDenom(ctx)
	communityPollBalance := k.distrKeeper.GetFeePool(ctx).CommunityPool
	availableCoinsCommunity := communityPollBalance.AmountOf(denom)
	if amt.AmountOf(denom).ToDec().GT(availableCoinsCommunity) {
		return errors.New("minting not completed because total supply cap is hit")
	}
	return k.bankKeeper.MintCoins(ctx, moduleName, amt)
}

func (k Keeper) safeBurnCoinsFromCommunityPool(ctx sdk.Context, coins sdk.Coins) error {
<<<<<<< HEAD
	communityPoolBalance := k.distrKeeper.GetFeePool(ctx).CommunityPool
	//ctx.Logger().Info("------communityPoolBalance is " + communityPoolBalance.String())
	if communityPoolBalance.AmountOf(k.BondDenom(ctx)).GTE(coins.AmountOf(k.BondDenom(ctx)).ToDec()) {
=======
	communityPollBalance := k.distrKeeper.GetFeePool(ctx).CommunityPool
	//ctx.Logger().Info("------communityPollBalance is " + communityPollBalance.String())
	if communityPollBalance.AmountOf(k.BondDenom(ctx)).GTE(coins.AmountOf(k.BondDenom(ctx)).ToDec()) {
>>>>>>> cherry-pick fix/qb1866: fix token supply
		k.bankKeeper.BurnCoins(ctx, distrtypes.ModuleName, coins)
		return nil
	}
	return errors.New("burning not completed as a result of insufficient balance in community pool")
}

// Restore total supply to 100M stos
func (k Keeper) RestoreTotalSupply(ctx sdk.Context) (minted, burned sdk.Coins) {
	InitialTotalSupply := k.InitialTotalSupply(ctx).Amount
	currentTotalSupply := k.bankKeeper.GetSupply(ctx, k.BondDenom(ctx)).Amount
	//ctx.Logger().Info("------currentTotalSupply is " + currentTotalSupply.String())
	if InitialTotalSupply.Equal(currentTotalSupply) {
		//ctx.Logger().Info("------no need to restore")
		return sdk.Coins{}, sdk.Coins{}
	}
	supplyDiff := currentTotalSupply.Sub(InitialTotalSupply)
	if supplyDiff.GT(sdk.ZeroInt()) {
		// burn surplus if currentTotalSupply > InitialTotalSupply
		amtToBurn := supplyDiff
		coinToBurn := sdk.NewCoin(k.BondDenom(ctx), amtToBurn)
		coinsToBurn := sdk.NewCoins(coinToBurn)
		err := k.safeBurnCoinsFromCommunityPool(ctx, coinsToBurn)
		if err != nil {
			return sdk.Coins{}, sdk.Coins{}
		}
		return sdk.Coins{}, coinsToBurn
	}
	// mint slack
	amtToMint := supplyDiff.Abs()
	coinToMint := sdk.NewCoin(k.BondDenom(ctx), amtToMint)
	coinsToMint := sdk.NewCoins(coinToMint)
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coinsToMint) // do not use SafeMintCoins here
	if err != nil {
		ctx.Logger().Error("Restore total supply failed:", err.Error())
		return sdk.Coins{}, sdk.Coins{}
	}
	// send new mint coins to community pool
	senderAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	err = k.distrKeeper.FundCommunityPool(ctx, coinsToMint, senderAddr)
	if err != nil {
		ctx.Logger().Error("Restore total supply failed:", err.Error())
		return sdk.Coins{}, sdk.Coins{}
	}
	return coinsToMint, sdk.Coins{}
}

func (k Keeper) GetSupply(ctx sdk.Context) (totalSupply sdk.Coin) {
	return k.bankKeeper.GetSupply(ctx, k.BondDenom(ctx))
}

func (k Keeper) GetCirculationSupply(ctx sdk.Context) (circulationSupply sdk.Coins) {
	// total supply  - validator deposit - resource node deposit -  mining pool - prepay
	totalSupply := k.bankKeeper.GetSupply(ctx, k.BondDenom(ctx))

	validatorBondedPoolAcc := k.accountKeeper.GetModuleAddress(stakingtypes.BondedPoolName)
	validatorStaking := k.bankKeeper.GetBalance(ctx, validatorBondedPoolAcc, k.BondDenom(ctx))

	resourceNodeBondedPoolAcc := k.accountKeeper.GetModuleAddress(registertypes.ResourceNodeBondedPool)
	resourceNodeDeposit := k.bankKeeper.GetBalance(ctx, resourceNodeBondedPoolAcc, k.BondDenom(ctx))

	metaNodeBondedPoolAcc := k.accountKeeper.GetModuleAddress(registertypes.MetaNodeNotBondedPool)
	metaNodeDeposit := k.bankKeeper.GetBalance(ctx, metaNodeBondedPoolAcc, k.BondDenom(ctx))

	miningPoolAcc := k.accountKeeper.GetModuleAddress(types.FoundationAccount)
	miningPool := k.bankKeeper.GetBalance(ctx, miningPoolAcc, k.BondDenom(ctx))

	unissuedPrepayAcc := k.accountKeeper.GetModuleAddress(registertypes.TotalUnissuedPrepay)
	unissuedPrepay := k.bankKeeper.GetBalance(ctx, unissuedPrepayAcc, k.BondDenom(ctx))

	circulationSupplyStos := totalSupply.
		Sub(validatorStaking).
		Sub(resourceNodeDeposit).
		Sub(metaNodeDeposit).
		Sub(miningPool).
		Sub(unissuedPrepay)

	circulationSupply = sdk.NewCoins(circulationSupplyStos)

	return
}
