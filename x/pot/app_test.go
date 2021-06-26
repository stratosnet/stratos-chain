package pot

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

const (
	stopFlagOutOfTotalMiningReward = true
	stopFlagSpecificMinedReward    = true
	stopFlagSpecificEpoch          = true
)

var (
	paramSpecificMinedReward = sdk.NewInt(160000000000)
	paramSpecificEpoch       = sdk.NewInt(100)
)

// initialize data of volume report
func setupMsgVolumeReport(newEpoch int64) types.MsgVolumeReport {
	volume1 := types.NewSingleNodeVolume(addrRes1, resourceNodeVolume1)
	volume2 := types.NewSingleNodeVolume(addrRes2, resourceNodeVolume2)
	volume3 := types.NewSingleNodeVolume(addrRes3, resourceNodeVolume3)

	nodesVolume := []types.SingleNodeVolume{volume1, volume2, volume3}
	reporter := addrIdx1
	epoch := sdk.NewInt(newEpoch)
	reportReference := "report for epoch " + epoch.String()

	volumeReportMsg := types.NewMsgVolumeReport(nodesVolume, reporter, epoch, reportReference)

	return volumeReportMsg
}

// Test case termination conditions
// modify stop flag & variable could make the test case stop when reach a specific condition
func isNeedStop(ctx sdk.Context, k Keeper, epoch sdk.Int, minedToken sdk.Int) bool {

	if stopFlagOutOfTotalMiningReward && minedToken.GTE(foundationDeposit.AmountOf("ustos")) {
		return true
	}
	if stopFlagSpecificMinedReward && minedToken.GTE(paramSpecificMinedReward) {
		return true
	}
	if stopFlagSpecificEpoch && epoch.GTE(paramSpecificEpoch) {
		return true
	}
	return false
}

func TestPotVolumeReportMsgs(t *testing.T) {

	/********************* initialize mock app *********************/
	SetConfig()
	//mApp, k, accountKeeper, bankKeeper, stakingKeeper, registerKeeper := getMockApp(t)
	mApp, k, stakingKeeper, bankKeeper, supplyKeeper := getMockApp(t)
	accs := setupAccounts(mApp)
	mock.SetGenesis(mApp, accs)
	mock.CheckBalance(t, mApp, foundationAccAddr, foundationDeposit)

	/********************* create validator with 50% commission *********************/
	commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	description := staking.NewDescription("foo_moniker", "", "", "", "")
	createValidatorMsg := staking.NewMsgCreateValidator(valOpValAddr1, valConsPubk1, sdk.NewCoin("ustos", valInitialStake), description, commission, sdk.OneInt())

	header := abci.Header{Height: mApp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, header, []sdk.Msg{createValidatorMsg}, []uint64{8}, []uint64{0}, true, true, valOpPrivKey1)
	mock.CheckBalance(t, mApp, valOpAccAddr1, nil)

	header = abci.Header{Height: mApp.LastBlockHeight() + 1}
	mApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	stakingKeeper.ApplyAndReturnValidatorSetUpdates(mApp.BaseApp.NewContext(true, header))
	validator := checkValidator(t, mApp, stakingKeeper, valOpValAddr1, true)

	require.Equal(t, valOpValAddr1, validator.OperatorAddress)
	require.Equal(t, sdk.Bonded, validator.Status)
	require.True(sdk.IntEq(t, valInitialStake, validator.BondedTokens()))
	header = abci.Header{Height: mApp.LastBlockHeight() + 1}
	ctx := mApp.BaseApp.NewContext(true, header)
	/*
		the sequence of the account list is related to the value of parameter "accNums" of mock.SignCheckDeliver() method
		accs := []authexported.Account{
			resOwnerAcc1, resOwnerAcc2, resOwnerAcc3, resOwnerAcc4, resOwnerAcc5,
			idxOwnerAcc1, idxOwnerAcc2, idxOwnerAcc3, valOwnerAcc1,
			resNodeAcc1, resNodeAcc2, resNodeAcc3, resNodeAcc4, resNodeAcc5,
			idxNodeAcc1, idxNodeAcc2, idxNodeAcc3,
		}
	*/
	var i int64
	i = 0
	for {
		ctx.Logger().Info("*****************************************************************************")
		/********************* prepare tx data *********************/
		volumeReportMsg := setupMsgVolumeReport(i + 1)

		lastTotalMinedToken := k.GetTotalMinedTokens(ctx)
		ctx.Logger().Info("last committed mined token = " + lastTotalMinedToken.String())
		if isNeedStop(ctx, k, volumeReportMsg.Epoch, lastTotalMinedToken) {
			break
		}

		/********************* print info *********************/
		ctx.Logger().Info("epoch " + volumeReportMsg.Epoch.String())
		S := k.RegisterKeeper.GetInitialGenesisStakeTotal(ctx).ToDec()
		Pt := k.GetTotalUnissuedPrepay(ctx).ToDec()
		Y := k.GetTotalConsumedOzone(volumeReportMsg.NodesVolume).ToDec()
		Lt := k.RegisterKeeper.GetRemainingOzoneLimit(ctx).ToDec()
		R := S.Add(Pt).Mul(Y).Quo(Lt.Add(Y))
		//ctx.Logger().Info("R = (S + Pt) * Y / (Lt + Y)")
		ctx.Logger().Info("S=" + S.String() + "\nPt=" + Pt.String() + "\nY=" + Y.String() + "\nLt=" + Lt.String() + "\nR=" + R.String() + "\n")

		ctx.Logger().Info("---------------------------")
		distributeGoal := types.InitDistributeGoal()
		distributeGoal, err := k.CalcTrafficRewardInTotal(ctx, volumeReportMsg.NodesVolume, distributeGoal)
		require.NoError(t, err)
		distributeGoal, err = k.CalcMiningRewardInTotal(ctx, distributeGoal)
		require.NoError(t, err)
		ctx.Logger().Info(distributeGoal.String())

		ctx.Logger().Info("---------------------------")
		distributeGoalBalance := distributeGoal
		rewardDetailMap := make(map[string]types.Reward)
		rewardDetailMap, distributeGoalBalance = k.CalcRewardForResourceNode(ctx, volumeReportMsg.NodesVolume, distributeGoalBalance, rewardDetailMap)
		rewardDetailMap, distributeGoalBalance = k.CalcRewardForIndexingNode(ctx, distributeGoalBalance, rewardDetailMap)
		ctx.Logger().Info("resourceNode1:  address = " + addrRes1.String())
		ctx.Logger().Info("           miningReward = " + rewardDetailMap[addrRes1.String()].RewardFromMiningPool.String())
		ctx.Logger().Info("          trafficReward = " + rewardDetailMap[addrRes1.String()].RewardFromTrafficPool.String())

		ctx.Logger().Info("resourceNode2:  address = " + addrRes2.String())
		ctx.Logger().Info("           miningReward = " + rewardDetailMap[addrRes2.String()].RewardFromMiningPool.String())
		ctx.Logger().Info("          trafficReward = " + rewardDetailMap[addrRes2.String()].RewardFromTrafficPool.String())

		ctx.Logger().Info("resourceNode3:  address = " + addrRes3.String())
		ctx.Logger().Info("           miningReward = " + rewardDetailMap[addrRes3.String()].RewardFromMiningPool.String())
		ctx.Logger().Info("          trafficReward = " + rewardDetailMap[addrRes3.String()].RewardFromTrafficPool.String())

		ctx.Logger().Info("resourceNode4:  address = " + addrRes4.String())
		ctx.Logger().Info("           miningReward = " + rewardDetailMap[addrRes4.String()].RewardFromMiningPool.String())
		ctx.Logger().Info("          trafficReward = " + rewardDetailMap[addrRes4.String()].RewardFromTrafficPool.String())

		ctx.Logger().Info("resourceNode5:  address = " + addrRes5.String())
		ctx.Logger().Info("           miningReward = " + rewardDetailMap[addrRes5.String()].RewardFromMiningPool.String())
		ctx.Logger().Info("          trafficReward = " + rewardDetailMap[addrRes5.String()].RewardFromTrafficPool.String())

		ctx.Logger().Info("indexingNode1:  address = " + addrIdx1.String())
		ctx.Logger().Info("           miningReward = " + rewardDetailMap[addrIdx1.String()].RewardFromMiningPool.String())
		ctx.Logger().Info("          trafficReward = " + rewardDetailMap[addrIdx1.String()].RewardFromTrafficPool.String())

		ctx.Logger().Info("indexingNode2:  address = " + addrIdx2.String())
		ctx.Logger().Info("           miningReward = " + rewardDetailMap[addrIdx2.String()].RewardFromMiningPool.String())
		ctx.Logger().Info("          trafficReward = " + rewardDetailMap[addrIdx2.String()].RewardFromTrafficPool.String())

		ctx.Logger().Info("indexingNode3:  address = " + addrIdx3.String())
		ctx.Logger().Info("           miningReward = " + rewardDetailMap[addrIdx3.String()].RewardFromMiningPool.String())
		ctx.Logger().Info("          trafficReward = " + rewardDetailMap[addrIdx3.String()].RewardFromTrafficPool.String())
		ctx.Logger().Info("---------------------------")

		/********************* record data before delivering tx  *********************/
		feePoolAccAddr := supplyKeeper.GetModuleAddress(k.FeeCollectorName)
		lastFoundationAccBalance := bankKeeper.GetCoins(ctx, foundationAccAddr).AmountOf("ustos")
		lastFeePool := bankKeeper.GetCoins(ctx, feePoolAccAddr).AmountOf("ustos")
		lastUnissuedPrepay := k.GetTotalUnissuedPrepay(ctx)

		/********************* deliver tx *********************/
		SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, header, []sdk.Msg{volumeReportMsg}, []uint64{14}, []uint64{uint64(i)}, true, true, privKeyIdx1)

		/********************* commit & check result *********************/
		header = abci.Header{Height: mApp.LastBlockHeight() + 1}
		ctx = mApp.BaseApp.NewContext(true, header)
		checkResult(t, ctx, k, volumeReportMsg.Epoch, lastFoundationAccBalance, lastUnissuedPrepay, lastFeePool)

		i++
	}

}

func checkResult(t *testing.T, ctx sdk.Context, k Keeper, currentEpoch sdk.Int,
	lastFoundationAccBalance sdk.Int, lastUnissuedPrepay sdk.Int, lastFeePool sdk.Int) {

	individualRewardTotal := sdk.ZeroInt()
	newMatureEpoch := currentEpoch.Add(sdk.NewInt(k.MatureEpoch(ctx)))
	rewardAddrList := k.GetRewardAddressPool(ctx)
	for _, addr := range rewardAddrList {
		individualReward := k.GetIndividualReward(ctx, addr, newMatureEpoch)
		individualRewardTotal = individualRewardTotal.Add(individualReward)

		ctx.Logger().Info("individualReward of [" + addr.String() + "] = " + individualReward.String())
	}

	feePoolAccAddr := k.SupplyKeeper.GetModuleAddress(k.FeeCollectorName)
	newFoundationAccBalance := k.BankKeeper.GetCoins(ctx, foundationAccAddr).AmountOf("ustos")
	newUnissuedPrepay := k.GetTotalUnissuedPrepay(ctx)

	rewardSrcChange := lastFoundationAccBalance.
		Sub(newFoundationAccBalance).
		Add(lastUnissuedPrepay).
		Sub(newUnissuedPrepay)

	newFeePool := k.BankKeeper.GetCoins(ctx, feePoolAccAddr).AmountOf("ustos")

	feePoolValChange := newFeePool.Sub(lastFeePool)
	ctx.Logger().Info("reward send to validator fee pool                               = " + feePoolValChange.String())

	rewardDestChange := feePoolValChange.Add(individualRewardTotal)

	require.Equal(t, rewardSrcChange, rewardDestChange)

}

func checkValidator(t *testing.T, mApp *mock.App, stakingKeeper staking.Keeper,
	addr sdk.ValAddress, expFound bool) staking.Validator {

	ctxCheck := mApp.BaseApp.NewContext(true, abci.Header{})
	validator, found := stakingKeeper.GetValidator(ctxCheck, addr)

	require.Equal(t, expFound, found)
	return validator
}

func getMockApp(t *testing.T) (*mock.App, Keeper, staking.Keeper, bank.Keeper, supply.Keeper) {
	mApp := mock.NewApp()

	RegisterCodec(mApp.Cdc)
	supply.RegisterCodec(mApp.Cdc)
	staking.RegisterCodec(mApp.Cdc)
	register.RegisterCodec(mApp.Cdc)

	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keyRegister := sdk.NewKVStoreKey(register.StoreKey)
	keyPot := sdk.NewKVStoreKey(StoreKey)

	feeCollector := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollector.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true

	bankKeeper := bank.NewBaseKeeper(mApp.AccountKeeper, mApp.ParamsKeeper.Subspace(bank.DefaultParamspace), blacklistedAddrs)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     {"fee_collector"},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(mApp.Cdc, keySupply, mApp.AccountKeeper, bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(mApp.Cdc, keyStaking, supplyKeeper, mApp.ParamsKeeper.Subspace(staking.DefaultParamspace))
	registerKeeper := register.NewKeeper(mApp.Cdc, keyRegister, mApp.AccountKeeper, bankKeeper, mApp.ParamsKeeper.Subspace(register.DefaultParamSpace))

	keeper := NewKeeper(mApp.Cdc, keyPot, mApp.ParamsKeeper.Subspace(DefaultParamSpace), auth.FeeCollectorName, bankKeeper, supplyKeeper, mApp.AccountKeeper, stakingKeeper, registerKeeper)

	mApp.Router().AddRoute(staking.RouterKey, staking.NewHandler(stakingKeeper))
	mApp.Router().AddRoute(RouterKey, NewHandler(keeper))
	mApp.SetEndBlocker(getEndBlocker(keeper))
	mApp.SetInitChainer(getInitChainer(mApp, keeper, mApp.AccountKeeper, supplyKeeper,
		[]supplyexported.ModuleAccountI{feeCollector, notBondedPool, bondPool}, stakingKeeper, registerKeeper))

	err := mApp.CompleteSetup(keyStaking, keySupply, keyRegister, keyPot)
	require.NoError(t, err)

	return mApp, keeper, stakingKeeper, bankKeeper, supplyKeeper
}

// getInitChainer initializes the chainer of the mock app and sets the genesis
// state. It returns an empty ResponseInitChain.
func getInitChainer(mapp *mock.App, keeper Keeper, accountKeeper auth.AccountKeeper, supplyKeeper supply.Keeper,
	blacklistedAddrs []supplyexported.ModuleAccountI, stakingKeeper staking.Keeper, registerKeeper register.Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}

		mapp.InitChainer(ctx, req)

		lastResourceNodeTotalStake := initialStakeRes1.Add(initialStakeRes2).Add(initialStakeRes3).Add(initialStakeRes4).Add(initialStakeRes5)
		lastIndexingNodeTotalStake := initialStakeIdx1.Add(initialStakeIdx2).Add(initialStakeIdx3)

		var lastResourceNodeStakes []register.LastResourceNodeStake
		lastResourceNodeStakes = append(lastResourceNodeStakes, register.LastResourceNodeStake{Address: addrRes1, Stake: initialStakeRes1})
		lastResourceNodeStakes = append(lastResourceNodeStakes, register.LastResourceNodeStake{Address: addrRes2, Stake: initialStakeRes2})
		lastResourceNodeStakes = append(lastResourceNodeStakes, register.LastResourceNodeStake{Address: addrRes3, Stake: initialStakeRes3})
		lastResourceNodeStakes = append(lastResourceNodeStakes, register.LastResourceNodeStake{Address: addrRes4, Stake: initialStakeRes4})
		lastResourceNodeStakes = append(lastResourceNodeStakes, register.LastResourceNodeStake{Address: addrRes5, Stake: initialStakeRes5})

		var lastIndexingNodeStakes []register.LastIndexingNodeStake
		lastIndexingNodeStakes = append(lastIndexingNodeStakes, register.LastIndexingNodeStake{Address: addrIdx1, Stake: initialStakeIdx1})
		lastIndexingNodeStakes = append(lastIndexingNodeStakes, register.LastIndexingNodeStake{Address: addrIdx2, Stake: initialStakeIdx2})
		lastIndexingNodeStakes = append(lastIndexingNodeStakes, register.LastIndexingNodeStake{Address: addrIdx3, Stake: initialStakeIdx3})

		resourceNodes := setupAllResourceNodes()
		indexingNodes := setupAllIndexingNodes()

		registerGenesis := register.NewGenesisState(register.DefaultParams(), lastResourceNodeTotalStake, lastResourceNodeStakes, resourceNodes,
			lastIndexingNodeTotalStake, lastIndexingNodeStakes, indexingNodes)

		register.InitGenesis(ctx, registerKeeper, registerGenesis)

		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}

		stakingGenesis := staking.NewGenesisState(staking.NewParams(staking.DefaultUnbondingTime, staking.DefaultMaxValidators, staking.DefaultMaxEntries, 0, "ustos"), nil, nil)

		totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000)))
		supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}

		validators := staking.InitGenesis(ctx, stakingKeeper, accountKeeper, supplyKeeper, stakingGenesis)

		//preset
		registerKeeper.SetRemainingOzoneLimit(ctx, remainingOzoneLimit)
		keeper.SetTotalUnissuedPrepay(ctx, totalUnissuedPrepay)

		//pot genesis data load
		InitGenesis(ctx, keeper, NewGenesisState(types.DefaultParams(), foundationAccAddr, initialOzonePrice))

		return abci.ResponseInitChain{
			Validators: validators,
		}
	}

}

// getEndBlocker returns a staking endblocker.
func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		validatorUpdates := keeper.StakingKeeper.BlockValidatorUpdates(ctx)

		return abci.ResponseEndBlock{
			ValidatorUpdates: validatorUpdates,
		}
	}
	return nil
}
