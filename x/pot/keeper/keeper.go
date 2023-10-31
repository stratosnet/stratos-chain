package keeper

import (
	"context"
	"encoding/hex"
	"math"

	"github.com/cometbft/cometbft/libs/log"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
)

// Keeper of the pot store
type Keeper struct {
	storeKey       storetypes.StoreKey
	cdc            codec.Codec
	accountKeeper  types.AccountKeeper
	bankKeeper     types.BankKeeper
	distrKeeper    types.DistrKeeper
	registerKeeper types.RegisterKeeper
	stakingKeeper  types.StakingKeeper

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper creates a pot keeper
func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrKeeper types.DistrKeeper,
	registerKeeper types.RegisterKeeper,
	stakingKeeper types.StakingKeeper,
	authority string,
) Keeper {
	keeper := Keeper{
		cdc:            cdc,
		storeKey:       key,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		distrKeeper:    distrKeeper,
		registerKeeper: registerKeeper,
		stakingKeeper:  stakingKeeper,
		authority:      authority,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) VolumeReport(ctx sdk.Context, walletVolumes types.WalletVolumes, reporter stratos.SdsAddress,
	epoch sdkmath.Int, reportReference string, txHash string) (err error) {

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
	if amt.AmountOf(denom).ToLegacyDec().GT(availableCoinsCommunity) {
		return types.ErrTotalSupplyCapHit
	}
	return k.bankKeeper.MintCoins(ctx, moduleName, amt)
}

func (k Keeper) safeBurnCoinsFromCommunityPool(ctx sdk.Context, coins sdk.Coins) error {
	communityPoolBalance := k.distrKeeper.GetFeePool(ctx).CommunityPool
	//ctx.Logger().Info("------communityPoolBalance is " + communityPoolBalance.String())
	if communityPoolBalance.AmountOf(k.BondDenom(ctx)).GTE(coins.AmountOf(k.BondDenom(ctx)).ToLegacyDec()) {
		k.bankKeeper.BurnCoins(ctx, distrtypes.ModuleName, coins)
		return nil
	}
	return types.ErrInsufficientCommunityPool
}

// RestoreTotalSupply Restore total supply to 100M stos
func (k Keeper) RestoreTotalSupply(ctx sdk.Context) (minted, burned sdk.Coins) {
	InitialTotalSupply := k.InitialTotalSupply(ctx).Amount
	currentTotalSupply := k.bankKeeper.GetSupply(ctx, k.BondDenom(ctx)).Amount
	//ctx.Logger().Info("------currentTotalSupply is " + currentTotalSupply.String())
	if InitialTotalSupply.Equal(currentTotalSupply) {
		//ctx.Logger().Info("------no need to restore")
		return sdk.Coins{}, sdk.Coins{}
	}
	supplyDiff := currentTotalSupply.Sub(InitialTotalSupply)
	if supplyDiff.GT(sdkmath.ZeroInt()) {
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

	totalMining := k.GetTotalMining(ctx)
	totalMinedTokens := k.GetTotalMinedTokens(ctx)

	unissuedPrepayAcc := k.accountKeeper.GetModuleAddress(registertypes.TotalUnissuedPrepay)
	unissuedPrepay := k.bankKeeper.GetBalance(ctx, unissuedPrepayAcc, k.BondDenom(ctx))

	circulationSupplyStos := totalSupply.
		Sub(validatorStaking).
		Sub(resourceNodeDeposit).
		Sub(metaNodeDeposit).
		Sub(totalMining.Sub(totalMinedTokens)).
		Sub(unissuedPrepay)

	circulationSupply = sdk.NewCoins(circulationSupplyStos)

	return
}

func (k Keeper) GetTotalReward(ctx sdk.Context, epoch sdkmath.Int) (totalReward types.TotalReward) {
	volumeReport := k.GetVolumeReport(ctx, epoch)

	if volumeReport == (types.VolumeReportRecord{}) {
		return types.TotalReward{}
	}
	hash, err := hex.DecodeString(volumeReport.TxHash)
	if err != nil {
		return types.TotalReward{}
	}

	clientCtx := client.Context{}.WithViper("")
	clientCtx, err = config.ReadFromClientConfig(clientCtx)
	if err != nil {
		return types.TotalReward{}
	}

	node, err := clientCtx.GetNode()
	if err != nil {
		return types.TotalReward{}
	}

	resTx, err := node.Tx(context.Background(), hash, true)
	if err != nil {
		return types.TotalReward{}
	}

	senderAddr := k.accountKeeper.GetModuleAddress(registertypes.TotalUnissuedPrepay)
	if senderAddr == nil {

		panic(errors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", registertypes.TotalUnissuedPrepay))
	}

	trafficReward := sdk.NewCoin(k.BondDenom(ctx), sdkmath.ZeroInt())
	txEvents := resTx.TxResult.GetEvents()
	for _, event := range txEvents {
		if event.Type == "coin_received" {
			attributes := event.GetAttributes()
			for _, attr := range attributes {
				if attr.GetKey() == "amount" {
					received, err := sdk.ParseCoinNormalized(attr.GetValue())
					if err != nil {
						continue
					}
					trafficReward = trafficReward.Add(received)
				}
			}
		}
	}
	miningReward := sdk.NewCoin(types.DefaultRewardDenom, sdkmath.NewInt(80).MulRaw(stratos.StosToWei))
	trafficReward = trafficReward.Sub(miningReward)
	totalReward = types.TotalReward{
		MiningReward:  sdk.NewCoins(miningReward),
		TrafficReward: sdk.NewCoins(trafficReward),
	}
	return
}

func (k Keeper) GetMetrics(ctx sdk.Context) types.Metrics {
	totalSupply := k.bankKeeper.GetSupply(ctx, k.BondDenom(ctx))
	totalMining := k.GetTotalMining(ctx)
	totalMinedTokens := k.GetTotalMinedTokens(ctx)

	totalBondedDepositOfResourceNodes := k.registerKeeper.GetResourceNodeBondedToken(ctx).Amount
	totalUnbondedDepositOfResourceNodes := k.registerKeeper.GetResourceNodeNotBondedToken(ctx).Amount
	totalResourceNodesDeposit := totalBondedDepositOfResourceNodes.Add(totalUnbondedDepositOfResourceNodes)

	denom := k.BondDenom(ctx)

	validatorBondedPoolAcc := k.accountKeeper.GetModuleAddress(stakingtypes.BondedPoolName)
	boundedDelegation := k.bankKeeper.GetBalance(ctx, validatorBondedPoolAcc, denom)

	validatorUnbondedPoolAcc := k.accountKeeper.GetModuleAddress(stakingtypes.NotBondedPoolName)
	unbondedDelegation := k.bankKeeper.GetBalance(ctx, validatorUnbondedPoolAcc, denom)

	unbondingDelegation := sdk.NewCoin(denom, sdkmath.ZeroInt())
	// NOTE: Uncomment to get all unboundings, not tested on performance
	// k.stakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd stakingtypes.UnbondingDelegation) (stop bool) {
	// 	for _, entry := range ubd.Entries {
	// 		unbondingDelegation = unbondingDelegation.Add(entry.Balance)
	// 	}
	// 	return false
	// })

	circulationSupply := k.GetCirculationSupply(ctx)

	totalMiningReward := sdkmath.ZeroInt()
	chainMiningReward := sdkmath.ZeroInt()
	resourceMiningReward := sdkmath.ZeroInt()
	metaMiningReward := sdkmath.ZeroInt()
	if miningParam, err := k.GetMiningRewardParamByMinedToken(ctx, totalMinedTokens); err == nil {
		totalMiningReward = miningParam.MiningReward.Amount
		chainMiningReward = totalMiningReward.Mul(miningParam.BlockChainPercentageInBp).Quo(sdkmath.NewInt(10000))
		resourceMiningReward = totalMiningReward.Mul(miningParam.ResourceNodePercentageInBp).Quo(sdkmath.NewInt(10000))
		metaMiningReward = totalMiningReward.Mul(miningParam.MetaNodePercentageInBp).Quo(sdkmath.NewInt(10000))
	}

	return types.Metrics{
		TotalSupply:               totalSupply.Amount,
		TotalMiningSupply:         totalMining.Amount,
		TotalMinedTokens:          totalMinedTokens.Amount,
		TotalResourceNodesDeposit: totalResourceNodesDeposit,
		TotalBondedDelegation:     boundedDelegation.Amount,
		TotalUnbondedDelegation:   unbondedDelegation.Amount,
		TotalUnbondingDelegation:  unbondingDelegation.Amount,
		CirculationSupply:         circulationSupply.AmountOf(denom),
		TotalMiningReward:         totalMiningReward,
		ChainMiningReward:         chainMiningReward,
		ResourceMiningReward:      resourceMiningReward,
		MetaMiningReward:          metaMiningReward,
	}
}
