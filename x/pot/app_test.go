package pot_test

import (
	"testing"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stratosnet/stratos-chain/app"
	stratos "github.com/stratosnet/stratos-chain/types"
	potKeeper "github.com/stratosnet/stratos-chain/x/pot/keeper"
	registerKeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

const (
	stopFlagOutOfTotalMiningReward = true
	stopFlagSpecificMinedReward    = false
	stopFlagSpecificEpoch          = true
)

var (
	paramSpecificMinedReward = sdk.NewCoins(sdk.NewCoin("ustos", sdk.NewInt(160000000000)))
	paramSpecificEpoch       = sdk.NewInt(10)
)

// initialize data of volume report
func setupMsgVolumeReport(newEpoch int64) *types.MsgVolumeReport {
	volume1 := types.NewSingleWalletVolume(resOwner1, resourceNodeVolume1)
	volume2 := types.NewSingleWalletVolume(resOwner2, resourceNodeVolume2)
	volume3 := types.NewSingleWalletVolume(resOwner3, resourceNodeVolume3)

	nodesVolume := []*types.SingleWalletVolume{volume1, volume2, volume3}
	reporter := idxNodeNetworkId1
	epoch := sdk.NewInt(newEpoch)
	reportReference := "report for epoch " + epoch.String()
	reporterOwner := idxOwner1

	pubKeys := make([][]byte, 1)
	for i := range pubKeys {
		pubKeys[i] = make([]byte, 1)
	}

	signature := types.NewBLSSignatureInfo(pubKeys, []byte("signature"), []byte("txData"))

	volumeReportMsg := types.NewMsgVolumeReport(nodesVolume, reporter, epoch, reportReference, reporterOwner, signature)

	return volumeReportMsg
}

func setupSlashingMsg() *types.MsgSlashingResourceNode {
	reporters := make([]stratos.SdsAddress, 0)
	reporters = append(reporters, idxNodeNetworkId1)
	reportOwner := make([]sdk.AccAddress, 0)
	reportOwner = append(reportOwner, idxOwner1)
	slashingMsg := types.NewMsgSlashingResourceNode(reporters, reportOwner, resNodeNetworkId1, resOwner1, resNodeSlashingUOZAmt1, true)
	return slashingMsg
}

// Test case termination conditions
// modify stop flag & variable could make the test case stop when reach a specific condition
func isNeedStop(ctx sdk.Context, k potKeeper.Keeper, epoch sdk.Int, minedToken sdk.Coin) bool {

	if stopFlagOutOfTotalMiningReward && (minedToken.Amount.GT(foundationDeposit.AmountOf(k.RewardDenom(ctx))) ||
		minedToken.Amount.GT(foundationDeposit.AmountOf(k.RewardDenom(ctx)))) {
		return true
	}
	if stopFlagSpecificMinedReward && minedToken.Amount.GT(paramSpecificMinedReward.AmountOf(k.BondDenom(ctx))) {
		return true
	}
	if stopFlagSpecificEpoch && epoch.GT(paramSpecificEpoch) {
		return true
	}
	return false
}

func TestPotVolumeReportMsgs(t *testing.T) {
	/********************* initialize mock app *********************/
	//mApp, k, stakingKeeper, bankKeeper, supplyKeeper, registerKeeper := getMockApp(t)
	accs, balances := setupAccounts()
	//stApp := app.SetupWithGenesisAccounts(accs, chainID, balances...)
	validators := make([]*tmtypes.Validator, 0)
	valSet := tmtypes.NewValidatorSet(validators)
	metaNodes := setupAllMetaNodes()
	resourceNodes := setupAllResourceNodes()

	stApp := app.SetupWithGenesisNodeSet(t, valSet, metaNodes, resourceNodes, accs, totalUnissuedPrepay, chainID, balances...)

	accountKeeper := stApp.GetAccountKeeper()
	bankKeeper := stApp.GetBankKeeper()
	registerKeeper := stApp.GetRegisterKeeper()
	potKeeper := stApp.GetPotKeeper()

	/********************* foundation account deposit *********************/
	header := tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := stApp.BaseApp.NewContext(true, header)

	foundationDepositMsg := types.NewMsgFoundationDeposit(foundationDeposit, foundationDepositorAccAddr)
	txGen := app.MakeTestEncodingConfig().TxConfig

	foundationDepositorAcc := accountKeeper.GetAccount(ctx, foundationDepositorAccAddr)
	accNum := foundationDepositorAcc.GetAccountNumber()
	accSeq := foundationDepositorAcc.GetSequence()
	_, _, err := app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, foundationDepositorPrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(types.FoundationAccount)
	app.CheckBalance(t, stApp, foundationAccountAddr, foundationDeposit)

	/********************* create validator with 50% commission *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	commission := stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	description := stakingtypes.NewDescription("foo_moniker", chainID, "", "", "")
	createValidatorMsg, err := stakingtypes.NewMsgCreateValidator(valOpValAddr1, valConsPubk1, sdk.NewCoin("ustos", valInitialStake), description, commission, sdk.OneInt())

	valOpAcc1 := accountKeeper.GetAccount(ctx, valOpAccAddr1)
	accNum = valOpAcc1.GetAccountNumber()
	accSeq = valOpAcc1.GetSequence()
	_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createValidatorMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, valOpPrivKey1)
	require.NoError(t, err)
	app.CheckBalance(t, stApp, valOpAccAddr1, nil)

	/********************** commit **********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	validator := checkValidator(t, stApp, valOpValAddr1, true)
	require.Equal(t, stakingtypes.Bonded, validator.Status)
	require.True(sdk.IntEq(t, valInitialStake, validator.BondedTokens()))

	/********************** loop sending volume report **********************/
	var i int64
	var slashingAmtSetup sdk.Int
	i = 0
	slashingAmtSetup = sdk.ZeroInt()
	for {

		/********************* test slashing msg when i==2 *********************/
		if i == 2 {
			println("********************************* Deliver Slashing Tx START ********************************************")
			slashingMsg := setupSlashingMsg()
			/********************* deliver tx *********************/

			idxOwnerAcc1 := accountKeeper.GetAccount(ctx, idxOwner1)
			ownerAccNum := idxOwnerAcc1.GetAccountNumber()
			ownerAccSeq := idxOwnerAcc1.GetSequence()

			_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{slashingMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, idxOwnerPrivKey1)
			require.NoError(t, err)
			/********************* commit & check result *********************/
			header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
			ctx = stApp.BaseApp.NewContext(true, header)

			slashingAmtSetup = registerKeeper.GetSlashing(ctx, resOwner1)

			totalConsumedUoz := resNodeSlashingUOZAmt1.ToDec()

			slashingAmtCheck := potKeeper.GetTrafficReward(ctx, totalConsumedUoz)
			println("slashingAmtSetup=" + slashingAmtSetup.String())
			require.Equal(t, slashingAmtSetup, slashingAmtCheck.TruncateInt())

			println("********************************* Deliver Slashing Tx END ********************************************")
		}

		println("*****************************************************************************")
		/********************* prepare tx data *********************/
		volumeReportMsg := setupMsgVolumeReport(i + 1)

		lastTotalMinedToken := potKeeper.GetTotalMinedTokens(ctx)
		println("last committed TotalMinedTokens = " + lastTotalMinedToken.String())
		epoch, ok := sdk.NewIntFromString(volumeReportMsg.Epoch.String())
		require.Equal(t, ok, true)

		if isNeedStop(ctx, potKeeper, epoch, lastTotalMinedToken) {
			break
		}

		totalConsumedUoz := potKeeper.GetTotalConsumedUoz(volumeReportMsg.WalletVolumes).ToDec()

		/********************* print info *********************/
		println("epoch " + volumeReportMsg.Epoch.String())
		S := registerKeeper.GetInitialGenesisStakeTotal(ctx).ToDec()
		Pt := registerKeeper.GetTotalUnissuedPrepay(ctx).Amount.ToDec()
		Y := totalConsumedUoz
		Lt := registerKeeper.GetRemainingOzoneLimit(ctx).ToDec()
		R := S.Add(Pt).Mul(Y).Quo(Lt.Add(Y))
		//println("R = (S + Pt) * Y / (Lt + Y)")
		println("S=" + S.String() + "\nPt=" + Pt.String() + "\nY=" + Y.String() + "\nLt=" + Lt.String() + "\nR=" + R.String() + "\n")

		println("---------------------------")
		distributeGoal := types.InitDistributeGoal()
		distributeGoal, err := potKeeper.CalcTrafficRewardInTotal(ctx, distributeGoal, totalConsumedUoz)
		require.NoError(t, err)

		distributeGoal, err = potKeeper.CalcMiningRewardInTotal(ctx, distributeGoal) //for main net
		require.NoError(t, err)
		println(distributeGoal.String())

		println("---------------------------")
		println("distribute detail:")
		distributeGoalBalance := distributeGoal
		rewardDetailMap := make(map[string]types.Reward)
		rewardDetailMap, distributeGoalBalance = potKeeper.CalcRewardForResourceNode(ctx, totalConsumedUoz, volumeReportMsg.WalletVolumes, distributeGoalBalance, rewardDetailMap)
		rewardDetailMap, distributeGoalBalance = potKeeper.CalcRewardForMetaNode(ctx, distributeGoalBalance, rewardDetailMap)

		println("resource_wallet1:  address = " + resOwner1.String())
		println("              miningReward = " + rewardDetailMap[resOwner1.String()].RewardFromMiningPool.String())
		println("             trafficReward = " + rewardDetailMap[resOwner1.String()].RewardFromTrafficPool.String())

		println("resource_wallet2:  address = " + resOwner2.String())
		println("              miningReward = " + rewardDetailMap[resOwner2.String()].RewardFromMiningPool.String())
		println("             trafficReward = " + rewardDetailMap[resOwner2.String()].RewardFromTrafficPool.String())

		println("resource_wallet3:  address = " + resOwner3.String())
		println("              miningReward = " + rewardDetailMap[resOwner3.String()].RewardFromMiningPool.String())
		println("             trafficReward = " + rewardDetailMap[resOwner3.String()].RewardFromTrafficPool.String())

		println("resource_wallet4:  address = " + resOwner4.String())
		println("              miningReward = " + rewardDetailMap[resOwner4.String()].RewardFromMiningPool.String())
		println("             trafficReward = " + rewardDetailMap[resOwner4.String()].RewardFromTrafficPool.String())

		println("resource_wallet5:  address = " + resOwner5.String())
		println("              miningReward = " + rewardDetailMap[resOwner5.String()].RewardFromMiningPool.String())
		println("             trafficReward = " + rewardDetailMap[resOwner5.String()].RewardFromTrafficPool.String())

		println("indexing_wallet1:  address = " + idxOwner1.String())
		println("              miningReward = " + rewardDetailMap[idxOwner1.String()].RewardFromMiningPool.String())
		println("             trafficReward = " + rewardDetailMap[idxOwner1.String()].RewardFromTrafficPool.String())

		println("indexing_wallet2:  address = " + idxOwner2.String())
		println("              miningReward = " + rewardDetailMap[idxOwner2.String()].RewardFromMiningPool.String())
		println("             trafficReward = " + rewardDetailMap[idxOwner2.String()].RewardFromTrafficPool.String())

		println("indexing_wallet3:  address = " + idxOwner3.String())
		println("              miningReward = " + rewardDetailMap[idxOwner3.String()].RewardFromMiningPool.String())
		println("             trafficReward = " + rewardDetailMap[idxOwner3.String()].RewardFromTrafficPool.String())
		println("---------------------------")

		/********************* record data before delivering tx  *********************/
		lastFoundationAccBalance := bankKeeper.GetAllBalances(ctx, foundationAccountAddr)
		lastUnissuedPrepay := registerKeeper.GetTotalUnissuedPrepay(ctx)
		lastMatureTotalOfResNode1 := potKeeper.GetMatureTotalReward(ctx, resOwner1)

		/********************* deliver tx *********************/
		idxOwnerAcc1 := accountKeeper.GetAccount(ctx, idxOwner1)
		ownerAccNum := idxOwnerAcc1.GetAccountNumber()
		ownerAccSeq := idxOwnerAcc1.GetSequence()

		_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{volumeReportMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, idxOwnerPrivKey1)
		require.NoError(t, err)

		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)

		epoch, ok = sdk.NewIntFromString(volumeReportMsg.Epoch.String())
		require.Equal(t, ok, true)

		checkResult(t, ctx, potKeeper,
			accountKeeper,
			bankKeeper,
			registerKeeper,
			epoch,
			lastFoundationAccBalance,
			lastUnissuedPrepay,
			lastMatureTotalOfResNode1,
			slashingAmtSetup,
		)

		i++
	}
}

// return : coins - slashing
func deductSlashingAmt(ctx sdk.Context, coins sdk.Coins, slashing sdk.Int) sdk.Coins {
	ret := sdk.Coins{}
	for _, coin := range coins {
		if coin.Amount.GTE(slashing) {
			coin = coin.Sub(sdk.NewCoin(coin.Denom, slashing))
			ret = ret.Add(coin)
			slashing = sdk.ZeroInt()
		} else {
			slashing = slashing.Sub(coin.Amount)
			coin = sdk.NewCoin(coin.Denom, sdk.ZeroInt())
			ret = ret.Add(coin)
		}
	}
	return ret
}

//for main net
func checkResult(t *testing.T, ctx sdk.Context,
	k potKeeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankKeeper.Keeper,
	registerKeeper registerKeeper.Keeper,
	currentEpoch sdk.Int,
	lastFoundationAccBalance sdk.Coins,
	lastUnissuedPrepay sdk.Coin,
	lastMatureTotalOfResNode1 sdk.Coins,
	slashingAmtSetup sdk.Int) {

	currentSlashing := registerKeeper.GetSlashing(ctx, resNodeAddr2)
	println("currentSlashing					= " + currentSlashing.String())

	individualRewardTotal := sdk.Coins{}
	newMatureEpoch := currentEpoch.Add(sdk.NewInt(k.MatureEpoch(ctx)))

	k.IteratorIndividualReward(ctx, newMatureEpoch, func(walletAddress sdk.AccAddress, individualReward types.Reward) (stop bool) {
		individualRewardTotal = individualRewardTotal.Add(individualReward.RewardFromTrafficPool...).Add(individualReward.RewardFromMiningPool...)
		println("individualReward of [" + walletAddress.String() + "] = " + individualReward.String())
		return false
	})

	feePoolAccAddr := accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	require.NotNil(t, feePoolAccAddr)
	foundationAccountAddr := accountKeeper.GetModuleAddress(types.FoundationAccount)
	newFoundationAccBalance := bankKeeper.GetAllBalances(ctx, foundationAccountAddr)
	newUnissuedPrepay := sdk.NewCoins(registerKeeper.GetTotalUnissuedPrepay(ctx))

	slashingChange := slashingAmtSetup.Sub(registerKeeper.GetSlashing(ctx, resOwner1))
	println("resource node 1 slashing change		= " + slashingChange.String())
	matureTotal := k.GetMatureTotalReward(ctx, resOwner1)
	immatureTotal := k.GetImmatureTotalReward(ctx, resOwner1)
	println("resource node 1 matureTotal		= " + matureTotal.String())
	println("resource node 1 immatureTotal		= " + immatureTotal.String())

	rewardSrcChange := lastFoundationAccBalance.
		Sub(newFoundationAccBalance).
		Add(lastUnissuedPrepay).
		Sub(newUnissuedPrepay)
	println("rewardSrcChange				= " + rewardSrcChange.String())

	// distribution module will send all tokens from "fee_collector" to "distribution" account in the BeginBlocker() method
	feePoolValChange := bankKeeper.GetAllBalances(ctx, feePoolAccAddr)
	println("reward send to validator fee pool	= " + feePoolValChange.String())

	rewardDestChange := feePoolValChange.Add(individualRewardTotal...)
	println("rewardDestChange			= " + rewardDestChange.String())

	require.Equal(t, rewardSrcChange, rewardDestChange)

	println("************************ slashing test***********************************")
	println("slashing change				= " + slashingChange.String())

	upcomingMaturedIndividual := sdk.Coins{}
	individualReward, found := k.GetIndividualReward(ctx, resOwner1, currentEpoch)
	if found {
		tmp := individualReward.RewardFromTrafficPool.Add(individualReward.RewardFromMiningPool...)
		upcomingMaturedIndividual = deductSlashingAmt(ctx, tmp, slashingChange)
	}
	println("upcomingMaturedIndividual		= " + upcomingMaturedIndividual.String())

	// get mature total changes
	newMatureTotalOfResNode1 := k.GetMatureTotalReward(ctx, resOwner1)
	matureTotalOfResNode1Change, _ := newMatureTotalOfResNode1.SafeSub(lastMatureTotalOfResNode1)
	if matureTotalOfResNode1Change == nil || matureTotalOfResNode1Change.IsAnyNegative() {
		matureTotalOfResNode1Change = sdk.Coins{}
	}
	println("matureTotalOfResNode1Change		= " + matureTotalOfResNode1Change.String())
	require.Equal(t, matureTotalOfResNode1Change, upcomingMaturedIndividual)
}

func checkValidator(t *testing.T, app *app.NewApp, addr sdk.ValAddress, expFound bool) stakingtypes.Validator {
	ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})
	validator, found := app.GetStakingKeeper().GetValidator(ctxCheck, addr)

	require.Equal(t, expFound, found)
	return validator
}

//func checkValidator(t *testing.T, mApp *mock.App, stakingKeeper staking.Keeper,
//	addr sdk.ValAddress, expFound bool) staking.Validator {
//
//	ctxCheck := mApp.BaseApp.NewContext(true, abci.Header{})
//	validator, found := stakingKeeper.GetValidator(ctxCheck, addr)
//
//	require.Equal(t, expFound, found)
//	return validator
//}

//func getMockApp(t *testing.T) (*mock.App, Keeper, staking.Keeper, bank.Keeper, supply.Keeper, register.Keeper) {
//	mApp := mock.NewApp()
//
//	RegisterCodec(mApp.Cdc)
//	supply.RegisterCodec(mApp.Cdc)
//	staking.RegisterCodec(mApp.Cdc)
//	register.RegisterCodec(mApp.Cdc)
//
//	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
//	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
//	keyRegister := sdk.NewKVStoreKey(register.StoreKey)
//	keyPot := sdk.NewKVStoreKey(StoreKey)
//
//	feeCollector := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
//	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
//	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)
//	foundationAccount := supply.NewEmptyModuleAccount(types.FoundationAccount)
//
//	blacklistedAddrs := make(map[string]bool)
//	blacklistedAddrs[feeCollector.GetAddress().String()] = true
//	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
//	blacklistedAddrs[bondPool.GetAddress().String()] = true
//	blacklistedAddrs[foundationAccount.GetAddress().String()] = true
//
//	bankKeeper := bank.NewBaseKeeper(mApp.AccountKeeper, mApp.ParamsKeeper.Subspace(bank.DefaultParamspace), blacklistedAddrs)
//	maccPerms := map[string][]string{
//		auth.FeeCollectorName:     nil,
//		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
//		staking.BondedPoolName:    {supply.Burner, supply.Staking},
//		types.FoundationAccount:   nil,
//	}
//	supplyKeeper := supply.NewKeeper(mApp.Cdc, keySupply, mApp.AccountKeeper, bankKeeper, maccPerms)
//	stakingKeeper := staking.NewKeeper(mApp.Cdc, keyStaking, supplyKeeper, mApp.ParamsKeeper.Subspace(staking.DefaultParamspace))
//	registerKeeper := register.NewKeeper(mApp.Cdc, keyRegister, mApp.ParamsKeeper.Subspace(register.DefaultParamSpace), mApp.AccountKeeper, bankKeeper)
//
//	keeper := NewKeeper(mApp.Cdc, keyPot, mApp.ParamsKeeper.Subspace(DefaultParamSpace), auth.FeeCollectorName, bankKeeper, supplyKeeper, mApp.AccountKeeper, stakingKeeper, registerKeeper)
//
//	mApp.Router().AddRoute(staking.RouterKey, staking.NewHandler(stakingKeeper))
//	mApp.Router().AddRoute(RouterKey, NewHandler(keeper))
//	mApp.SetEndBlocker(getEndBlocker(keeper))
//	mApp.SetInitChainer(getInitChainer(mApp, keeper, mApp.AccountKeeper, supplyKeeper,
//		[]supplyexported.ModuleAccountI{feeCollector, notBondedPool, bondPool}, stakingKeeper, registerKeeper))
//
//	err := mApp.CompleteSetup(keyStaking, keySupply, keyRegister, keyPot)
//	require.NoError(t, err)
//
//	return mApp, keeper, stakingKeeper, bankKeeper, supplyKeeper, registerKeeper
//}

// getInitChainer initializes the chainer of the mock app and sets the genesis
// state. It returns an empty ResponseInitChain.
//func getInitChainer(mapp *mock.App, keeper Keeper, accountKeeper auth.AccountKeeper, supplyKeeper supply.Keeper,
//	blacklistedAddrs []supplyexported.ModuleAccountI, stakingKeeper staking.Keeper, registerKeeper register.Keeper) sdk.InitChainer {
//	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
//		// set module accounts
//		for _, macc := range blacklistedAddrs {
//			supplyKeeper.SetModuleAccount(ctx, macc)
//		}
//
//		mapp.InitChainer(ctx, req)
//
//		resourceNodes := setupAllResourceNodes()
//		metaNodes := setupAllMetaNodes()
//
//		registerGenesis := registertypes.NewGenesisState(
//			register.DefaultParams(),
//			resourceNodes,
//			metaNodes,
//			initialUOzonePrice,
//			sdk.ZeroInt(),
//			make([]register.Slashing, 0),
//		)
//
//		register.InitGenesis(ctx, registerKeeper, registerGenesis)
//
//		// set module accounts
//		for _, macc := range blacklistedAddrs {
//			supplyKeeper.SetModuleAccount(ctx, macc)
//		}
//
//		stakingGenesis := staking.NewGenesisState(staking.NewParams(staking.DefaultUnbondingTime, staking.DefaultMaxValidators, staking.DefaultMaxEntries, 0, "ustos"), nil, nil)
//
//		totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000)))
//		supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))
//
//		// set module accounts
//		for _, macc := range blacklistedAddrs {
//			supplyKeeper.SetModuleAccount(ctx, macc)
//		}
//
//		validators := staking.InitGenesis(ctx, stakingKeeper, accountKeeper, supplyKeeper, stakingGenesis)
//
//		//preset
//		keeper.RegisterKeeper.SetTotalUnissuedPrepay(ctx, totalUnissuedPrepay)
//
//		//pot genesis data load
//		InitGenesis(ctx, keeper, NewGenesisState(
//			types.DefaultParams(),
//			sdk.NewCoin(types.DefaultRewardDenom, sdk.ZeroInt()),
//			0,
//			make([]types.ImmatureTotal, 0),
//			make([]types.MatureTotal, 0),
//			make([]types.Reward, 0),
//		))
//
//		return abci.ResponseInitChain{
//			Validators: validators,
//		}
//	}
//
//}

// getEndBlocker returns a staking endblocker.
//func getEndBlocker(keeper Keeper) sdk.EndBlocker {
//	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
//		validatorUpdates := keeper.StakingKeeper.BlockValidatorUpdates(ctx)
//
//		return abci.ResponseEndBlock{
//			ValidatorUpdates: validatorUpdates,
//		}
//	}
//	return nil
//}
