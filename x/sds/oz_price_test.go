package sds_test

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"

	sdkmath "cosmossdk.io/math"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/stratosnet/stratos-chain/app"
	stratostestutil "github.com/stratosnet/stratos-chain/testutil"
	stratos "github.com/stratosnet/stratos-chain/types"
	potKeeper "github.com/stratosnet/stratos-chain/x/pot/keeper"
	pottypes "github.com/stratosnet/stratos-chain/x/pot/types"
	registerKeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
	sdstypes "github.com/stratosnet/stratos-chain/x/sds/types"
)

/**
Test scenarios:
1. intt chain

*/

const (
	chainID = "testchain_1-1"
)

var (
	depositNozRateInt = sdkmath.NewInt(1e5)

	paramSpecificMinedReward = sdk.NewCoins(stratos.NewCoinInt64(160000000000))
	paramSpecificEpoch       = sdkmath.NewInt(10)

	resNodeSlashingNOZAmt1            = sdkmath.NewInt(1e11)
	resNodeSlashingEffectiveTokenAmt1 = sdkmath.NewInt(1).MulRaw(stratos.StosToWei)

	resourceNodeVolume1 = sdkmath.NewInt(537500000000)
	resourceNodeVolume2 = sdkmath.NewInt(200000000000)
	resourceNodeVolume3 = sdkmath.NewInt(200000000000)

	prepayAmount = sdk.NewCoins(stratos.NewCoin(sdkmath.NewInt(20).MulRaw(stratos.StosToWei)))

	depositForSendingTx    = sdkmath.NewInt(1000).MulRaw(stratos.StosToWei)
	totalUnissuedPrepayVal = sdkmath.ZeroInt()
	totalUnissuedPrepay    = stratos.NewCoin(totalUnissuedPrepayVal)

	foundationDepositorPrivKey = secp256k1.GenPrivKey()
	foundationDepositorAccAddr = sdk.AccAddress(foundationDepositorPrivKey.PubKey().Address())
	foundationDeposit          = sdk.NewCoins(sdk.NewCoin(stratos.Wei, sdkmath.NewInt(4e7).MulRaw(stratos.StosToWei)))

	valInitialStake        = sdkmath.NewInt(15).MulRaw(stratos.StosToWei)
	resNodeInitialDeposit  = sdkmath.NewInt(3).MulRaw(stratos.StosToWei)
	metaNodeInitialDeposit = sdkmath.NewInt(1000).MulRaw(stratos.StosToWei)

	// wallet private keys
	resOwnerPrivKey1  = secp256k1.GenPrivKey()
	resOwnerPrivKey2  = secp256k1.GenPrivKey()
	resOwnerPrivKey3  = secp256k1.GenPrivKey()
	resOwnerPrivKey4  = secp256k1.GenPrivKey()
	resOwnerPrivKey5  = secp256k1.GenPrivKey()
	metaOwnerPrivKey1 = secp256k1.GenPrivKey()
	metaOwnerPrivKey2 = secp256k1.GenPrivKey()
	metaOwnerPrivKey3 = secp256k1.GenPrivKey()

	// wallet addresses
	resOwner1  = sdk.AccAddress(resOwnerPrivKey1.PubKey().Address())
	resOwner2  = sdk.AccAddress(resOwnerPrivKey2.PubKey().Address())
	resOwner3  = sdk.AccAddress(resOwnerPrivKey3.PubKey().Address())
	resOwner4  = sdk.AccAddress(resOwnerPrivKey4.PubKey().Address())
	resOwner5  = sdk.AccAddress(resOwnerPrivKey5.PubKey().Address())
	metaOwner1 = sdk.AccAddress(metaOwnerPrivKey1.PubKey().Address())
	metaOwner2 = sdk.AccAddress(metaOwnerPrivKey2.PubKey().Address())
	metaOwner3 = sdk.AccAddress(metaOwnerPrivKey3.PubKey().Address())

	// P2P public key of resource nodes
	resNodeP2PPubKey1 = ed25519.GenPrivKey().PubKey()
	resNodeP2PPubKey2 = ed25519.GenPrivKey().PubKey()
	resNodeP2PPubKey3 = ed25519.GenPrivKey().PubKey()
	resNodeP2PPubKey4 = ed25519.GenPrivKey().PubKey()
	resNodeP2PPubKey5 = ed25519.GenPrivKey().PubKey()
	// P2P address of resource nodes
	resNodeP2PAddr1 = stratos.SdsAddress(resNodeP2PPubKey1.Address())
	resNodeP2PAddr2 = stratos.SdsAddress(resNodeP2PPubKey2.Address())
	resNodeP2PAddr3 = stratos.SdsAddress(resNodeP2PPubKey3.Address())
	resNodeP2PAddr4 = stratos.SdsAddress(resNodeP2PPubKey4.Address())
	resNodeP2PAddr5 = stratos.SdsAddress(resNodeP2PPubKey5.Address())

	// P2P private key of meta nodes
	metaNodeP2PPrivKey1 = ed25519.GenPrivKey()
	metaNodeP2PPrivKey2 = ed25519.GenPrivKey()
	metaNodeP2PPrivKey3 = ed25519.GenPrivKey()
	// P2P public key of meta nodes
	metaNodeP2PPubKey1 = metaNodeP2PPrivKey1.PubKey()
	metaNodeP2PPubKey2 = metaNodeP2PPrivKey2.PubKey()
	metaNodeP2PPubKey3 = metaNodeP2PPrivKey3.PubKey()
	// P2P address of meta nodes
	metaNodeP2PAddr1 = stratos.SdsAddress(metaNodeP2PPubKey1.Address())
	metaNodeP2PAddr2 = stratos.SdsAddress(metaNodeP2PPubKey2.Address())
	metaNodeP2PAddr3 = stratos.SdsAddress(metaNodeP2PPubKey3.Address())

	valOpPrivKey1 = secp256k1.GenPrivKey()
	valOpPubKey1  = valOpPrivKey1.PubKey()
	valOpValAddr1 = sdk.ValAddress(valOpPubKey1.Address())
	valOpAccAddr1 = sdk.AccAddress(valOpPubKey1.Address())

	valConsPrivKey1 = ed25519.GenPrivKey()
	valConsPubk1    = valConsPrivKey1.PubKey()
)

type NozPriceFactors struct {
	NOzonePrice           sdkmath.LegacyDec
	InitialTotalDeposit   sdkmath.Int
	EffectiveTotalDeposit sdkmath.Int
	TotalUnissuedPrepay   sdkmath.Int
	DepositAndPrepay      sdkmath.Int
	OzoneLimit            sdkmath.Int
	NozSupply             sdkmath.Int
}

func TestPriceCurve(t *testing.T) {

	NUM_TESTS := 100

	initFactorsBefore := &NozPriceFactors{
		NOzonePrice:           nozPrice,
		InitialTotalDeposit:   initialTotalDepositStore,
		EffectiveTotalDeposit: initialTotalDepositStore,
		TotalUnissuedPrepay:   totalUnissuedPrepayStore,
		DepositAndPrepay:      initialTotalDepositStore.Add(totalUnissuedPrepayStore),
		OzoneLimit:            initialTotalDepositStore.ToLegacyDec().Quo(nozPrice).TruncateInt(),
		NozSupply:             initialTotalDepositStore.ToLegacyDec().Quo(depositNozRateInt.ToLegacyDec()).TruncateInt(),
	}

	initFactorsBefore, _, _ = simulatePriceChange(t, &PriceChangeEvent{
		depositDelta:        sdkmath.ZeroInt(),
		unissuedPrepayDelta: sdkmath.ZeroInt(),
	}, initFactorsBefore)

	depositChangePerm := rand.Perm(NUM_TESTS)
	prepayChangePerm := rand.Perm(NUM_TESTS)

	for i := 0; i < NUM_TESTS; i++ {
		tempDepositSign := 1
		if i > 50 && rand.Intn(5) >= 3 {
			tempDepositSign = -1
		}
		tempPrepaySign := 1
		if i > 50 && rand.Intn(5) >= 3 {
			tempPrepaySign = -1
		}
		depositDeltaChange := sdkmath.NewInt(int64(depositChangePerm[i])).MulRaw(stratos.StosToWei)
		unissuedPrepayDeltaChange := sdkmath.NewInt(int64(prepayChangePerm[i])).MulRaw(stratos.StosToWei)
		change := &PriceChangeEvent{
			depositDelta:        depositDeltaChange.MulRaw(int64(tempDepositSign)),
			unissuedPrepayDelta: unissuedPrepayDeltaChange.MulRaw(int64(tempPrepaySign)),
		}
		t.Logf("\ndepositDeltaOri: %d, unissuedPrepayDeltaOri: %d\n", depositChangePerm[i], prepayChangePerm[i])
		t.Logf("\ndepositDelta: %v, unissuedPrepayDelta: %v\n", change.depositDelta.String(), change.unissuedPrepayDelta.String())
		initFactorsBefore, _, _ = simulatePriceChange(t, change, initFactorsBefore)
	}
}

func TestOzPriceChange(t *testing.T) {
	/********************* initialize mock app *********************/
	accs, balances := setupAccounts()

	// create validator set with single validator
	consPubKey, err := cryptocodec.ToTmPubKeyInterface(valConsPubk1)
	validator := tmtypes.NewValidator(consPubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	metaNodes := setupAllMetaNodes()
	//resourceNodes := setupAllResourceNodes()
	resourceNodes := make([]registertypes.ResourceNode, 0)

	stApp := stratostestutil.SetupWithGenesisNodeSet(t, valSet, metaNodes, resourceNodes, accs, chainID, false, balances...)

	accountKeeper := stApp.GetAccountKeeper()
	bankKeeper := stApp.GetBankKeeper()
	registerKeeper := stApp.GetRegisterKeeper()
	potKeeper := stApp.GetPotKeeper()
	distrKeeper := stApp.GetDistrKeeper()

	/********************* foundation account deposit *********************/
	header := tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := stApp.BaseApp.NewContext(true, header)

	foundationDepositMsg := pottypes.NewMsgFoundationDeposit(foundationDeposit, foundationDepositorAccAddr)
	txGen := stratostestutil.MakeTestEncodingConfig().TxConfig

	senderAcc := accountKeeper.GetAccount(ctx, foundationDepositorAccAddr)
	accNum := senderAcc.GetAccountNumber()
	accSeq := senderAcc.GetSequence()
	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, foundationDepositorPrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(pottypes.FoundationAccount)
	stratostestutil.CheckBalance(t, stApp, foundationAccountAddr, foundationDeposit)

	/********************** commit **********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	_, nozSupply := potKeeper.NozSupply(ctx)
	St, Pt, Lt := registerKeeper.GetCurrNozPriceParams(ctx)
	nozPriceFactorsSeq0, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, NozPriceFactors{
		NOzonePrice:           potKeeper.GetCurrentNozPrice(St, Pt, Lt),
		InitialTotalDeposit:   registerKeeper.GetInitialGenesisDepositTotal(ctx),
		EffectiveTotalDeposit: registerKeeper.GetEffectiveTotalDeposit(ctx),
		TotalUnissuedPrepay:   registerKeeper.GetTotalUnissuedPrepay(ctx).Amount,
		DepositAndPrepay:      registerKeeper.GetInitialGenesisDepositTotal(ctx).Add(registerKeeper.GetTotalUnissuedPrepay(ctx).Amount),
		OzoneLimit:            registerKeeper.GetRemainingOzoneLimit(ctx),
		NozSupply:             nozSupply,
	})

	// start testing
	t.Log("\n********************************* Deliver Prepay Tx START ********************************************")
	prepayMsg := setupPrepayMsg()
	/********************* deliver tx *********************/

	senderAcc = accountKeeper.GetAccount(ctx, resOwner1)
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()

	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{prepayMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, resOwnerPrivKey1)
	require.NoError(t, err)
	t.Log("********************************* Deliver Prepay Tx END ********************************************\n\n...\n[NEXT TEST CASE]")
	/********************* new height & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	nozPriceFactorsSeq1, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq0)
	require.True(t, nozPricePercentage.GT(sdkmath.LegacyZeroDec()), "noz price should increase after PREPAY")
	require.True(t, ozoneLimitPercentage.LT(sdkmath.LegacyZeroDec()), "OzLimit should decrease after PREPAY")

	t.Log("********************************* Deliver CreateResourceNode Tx START ********************************************")
	createResourceNodeMsg := setupMsgCreateResourceNode1()
	/********************* deliver tx *********************/

	senderAcc = accountKeeper.GetAccount(ctx, resOwner1)
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()

	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, resOwnerPrivKey1)
	require.NoError(t, err)
	t.Log("********************************* Deliver CreateResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")
	/********************* new height & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	nozPriceFactorsSeq2, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq1)
	require.True(t, nozPricePercentage.Equal(sdkmath.LegacyZeroDec()), "noz price shouldn't change after CreateResourceNode")
	require.True(t, ozoneLimitPercentage.Equal(sdkmath.LegacyZeroDec()), "OzLimit shouldn't change  after CreateResourceNode")

	t.Log("********************************* Deliver UpdateEffectiveDeposit Tx START ********************************************")
	updateEffectiveDepositMsg := setupMsgUpdateEffectiveDeposit1()
	/********************* deliver tx *********************/

	senderAcc = accountKeeper.GetAccount(ctx, metaOwner1)
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()

	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{updateEffectiveDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, metaOwnerPrivKey1)
	require.NoError(t, err)
	t.Log("********************************* Deliver UpdateEffectiveDeposit Tx END ********************************************\n\n...\n[NEXT TEST CASE]")
	/********************* new height & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	t.Log("********************************* Deliver UnsuspendResourceNode Tx (Slashing) START ********************************************")
	unsuspendMsg := setupUnsuspendMsg()

	/********************* deliver tx *********************/
	senderAcc = accountKeeper.GetAccount(ctx, metaOwner1)
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()
	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{unsuspendMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, metaOwnerPrivKey1)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	slashingAmtSetup := registerKeeper.GetSlashing(ctx, resOwner1)

	totalConsumedNoz := sdkmath.ZeroInt().ToLegacyDec()

	slashingAmtCheck := potKeeper.GetTrafficReward(ctx, totalConsumedNoz)
	t.Log("slashingAmtSetup=" + slashingAmtSetup.String())
	require.Equal(t, slashingAmtSetup, slashingAmtCheck.TruncateInt())

	nozPriceFactorsSeq3, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq2)
	require.True(t, nozPricePercentage.LT(sdkmath.LegacyZeroDec()), "noz price should decrease after UnsuspendResourceNode")
	require.True(t, ozoneLimitPercentage.GT(sdkmath.LegacyZeroDec()), "OzLimit should increase after UnsuspendResourceNode")
	t.Log("********************************* Deliver UnsuspendResourceNode Tx (Slashing) END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver SuspendResourceNode Tx (Slashing) START ********************************************")
	slashingMsg := setupSlashingMsg()
	/********************* deliver tx *********************/

	senderAcc = accountKeeper.GetAccount(ctx, metaOwner1)
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()

	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{slashingMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, metaOwnerPrivKey1)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	//slashingAmtSetup = registerKeeper.GetSlashing(ctx, resOwner1)
	//
	//totalConsumedNoz = resNodeSlashingNOZAmt1.ToLegacyDec()
	//
	//slashingAmtCheck = potKeeper.GetTrafficReward(ctx, totalConsumedNoz)
	//t.Log("slashingAmtSetup=" + slashingAmtSetup.String())
	//require.Equal(t, slashingAmtSetup, slashingAmtCheck.TruncateInt())

	nozPriceFactorsSeq4, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq3)
	require.True(t, nozPricePercentage.GT(sdkmath.LegacyZeroDec()), "noz price should increase after SlashResourceNode")
	require.True(t, ozoneLimitPercentage.LT(sdkmath.LegacyZeroDec()), "OzLimit should decrease after SlashResourceNode")

	_, nozPricePercentage42, ozoneLimitPercentage42 := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq2)
	require.True(t, nozPricePercentage42.Equal(sdkmath.LegacyZeroDec()), "noz price after SlashResourceNode should be same with the price when node hadn't been activated")
	require.True(t, ozoneLimitPercentage42.Equal(sdkmath.LegacyZeroDec()), "OzLimit after SlashResourceNode should be same with the ozLimit when node hadn't been activated")
	t.Log("********************************* Deliver SuspendResourceNode Tx (Slashing) END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver VolumeReport Tx START ********************************************")
	/********************* prepare tx data *********************/
	volumeReportMsg := setupMsgVolumeReport(t, 1)

	lastTotalMinedToken := potKeeper.GetTotalMinedTokens(ctx)
	t.Log("last committed TotalMinedTokens = " + lastTotalMinedToken.String())
	epoch, ok := sdkmath.NewIntFromString(volumeReportMsg.Epoch.String())
	require.Equal(t, ok, true)

	totalConsumedNoz = potKeeper.GetTotalConsumedNoz(volumeReportMsg.WalletVolumes).ToLegacyDec()
	remaining, total := potKeeper.NozSupply(ctx)
	require.True(t, potKeeper.GetTotalConsumedNoz(volumeReportMsg.WalletVolumes).Add(remaining).LTE(total), "remaining+consumed Noz exceeds total Noz supply")

	/********************* print info *********************/
	t.Log("epoch " + volumeReportMsg.Epoch.String())
	StDec := registerKeeper.GetEffectiveTotalDeposit(ctx).ToLegacyDec()
	PtDec := registerKeeper.GetTotalUnissuedPrepay(ctx).Amount.ToLegacyDec()
	Y := totalConsumedNoz
	LtDec := registerKeeper.GetRemainingOzoneLimit(ctx).ToLegacyDec()
	R := StDec.Add(PtDec).Mul(Y).Quo(LtDec.Add(Y))
	//t.Log("R = (S + Pt) * Y / (Lt + Y)")
	t.Log("St=" + StDec.String() + "\nPt=" + PtDec.String() + "\nY=" + Y.String() + "\nLt=" + LtDec.String() + "\nR=" + R.String() + "\n")

	t.Log("---------------------------")
	potKeeper.InitVariable(ctx)
	distributeGoal := pottypes.InitDistributeGoal()
	distributeGoal, err = potKeeper.CalcTrafficRewardInTotal(ctx, distributeGoal, totalConsumedNoz)
	require.NoError(t, err)

	distributeGoal, err = potKeeper.CalcMiningRewardInTotal(ctx, distributeGoal) //for main net
	require.NoError(t, err)
	t.Log(distributeGoal.String())

	t.Log("---------------------------")
	t.Log("distribute detail:")
	distributeGoalBalance := distributeGoal
	rewardDetailMap := make(map[string]pottypes.Reward)
	rewardDetailMap = potKeeper.CalcRewardForResourceNode(ctx, totalConsumedNoz, volumeReportMsg.WalletVolumes, distributeGoalBalance, rewardDetailMap)
	rewardDetailMap = potKeeper.CalcRewardForMetaNode(ctx, distributeGoalBalance, rewardDetailMap)

	t.Log("resource_wallet1:  address = " + resOwner1.String())
	t.Log("              miningReward = " + rewardDetailMap[resOwner1.String()].RewardFromMiningPool.String())
	t.Log("             trafficReward = " + rewardDetailMap[resOwner1.String()].RewardFromTrafficPool.String())

	t.Log("resource_wallet2:  address = " + resOwner2.String())
	t.Log("              miningReward = " + rewardDetailMap[resOwner2.String()].RewardFromMiningPool.String())
	t.Log("             trafficReward = " + rewardDetailMap[resOwner2.String()].RewardFromTrafficPool.String())

	t.Log("resource_wallet3:  address = " + resOwner3.String())
	t.Log("              miningReward = " + rewardDetailMap[resOwner3.String()].RewardFromMiningPool.String())
	t.Log("             trafficReward = " + rewardDetailMap[resOwner3.String()].RewardFromTrafficPool.String())

	t.Log("resource_wallet4:  address = " + resOwner4.String())
	t.Log("              miningReward = " + rewardDetailMap[resOwner4.String()].RewardFromMiningPool.String())
	t.Log("             trafficReward = " + rewardDetailMap[resOwner4.String()].RewardFromTrafficPool.String())

	t.Log("resource_wallet5:  address = " + resOwner5.String())
	t.Log("              miningReward = " + rewardDetailMap[resOwner5.String()].RewardFromMiningPool.String())
	t.Log("             trafficReward = " + rewardDetailMap[resOwner5.String()].RewardFromTrafficPool.String())

	t.Log("meta_wallet1:      address = " + metaOwner1.String())
	t.Log("              miningReward = " + rewardDetailMap[metaOwner1.String()].RewardFromMiningPool.String())
	t.Log("             trafficReward = " + rewardDetailMap[metaOwner1.String()].RewardFromTrafficPool.String())

	t.Log("meta_wallet2:      address = " + metaOwner2.String())
	t.Log("              miningReward = " + rewardDetailMap[metaOwner2.String()].RewardFromMiningPool.String())
	t.Log("             trafficReward = " + rewardDetailMap[metaOwner2.String()].RewardFromTrafficPool.String())

	t.Log("meta_wallet3:      address = " + metaOwner3.String())
	t.Log("              miningReward = " + rewardDetailMap[metaOwner3.String()].RewardFromMiningPool.String())
	t.Log("             trafficReward = " + rewardDetailMap[metaOwner3.String()].RewardFromTrafficPool.String())
	t.Log("---------------------------")

	/********************* record data before delivering tx  *********************/
	lastFoundationAccBalance := bankKeeper.GetAllBalances(ctx, foundationAccountAddr)
	lastUnissuedPrepay := registerKeeper.GetTotalUnissuedPrepay(ctx)
	lastCommunityPool := sdk.NewCoins(sdk.NewCoin(potKeeper.BondDenom(ctx), distrKeeper.GetFeePool(ctx).CommunityPool.AmountOf(potKeeper.BondDenom(ctx)).TruncateInt()))
	lastMatureTotalOfResNode1 := potKeeper.GetMatureTotalReward(ctx, resOwner1)

	/********************* deliver tx *********************/
	senderAcc = accountKeeper.GetAccount(ctx, metaOwner1)
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()

	feePoolAccAddr := accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	require.NotNil(t, feePoolAccAddr)
	feeCollectorToFeePoolAtBeginBlock := bankKeeper.GetBalance(ctx, feePoolAccAddr, potKeeper.BondDenom(ctx))

	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{volumeReportMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, metaOwnerPrivKey1)
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	stApp.EndBlock(abci.RequestEndBlock{Height: header.Height})
	stApp.Commit()
	require.NoError(t, err)

	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	epoch, ok = sdkmath.NewIntFromString(volumeReportMsg.Epoch.String())
	require.Equal(t, ok, true)

	checkResult(t, ctx, potKeeper,
		accountKeeper,
		bankKeeper,
		registerKeeper,
		epoch,
		lastFoundationAccBalance,
		lastUnissuedPrepay,
		lastCommunityPool,
		lastMatureTotalOfResNode1,
		slashingAmtSetup,
		feeCollectorToFeePoolAtBeginBlock,
	)

	nozPriceFactorsSeq5, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq4)
	require.True(t, nozPricePercentage.LT(sdkmath.LegacyZeroDec()), "noz price should decrease after VolumeReport")
	require.True(t, ozoneLimitPercentage.GT(sdkmath.LegacyZeroDec()), "OzLimit shouldn't change after VolumeReport")
	t.Log("********************************* Deliver VolumeReport Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver CreateResourceNode Tx START ********************************************")
	createResourceNodeMsg2 := setupMsgCreateResourceNode2()
	/********************* deliver tx *********************/

	senderAcc = accountKeeper.GetAccount(ctx, resOwner2)
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()

	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg2}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, resOwnerPrivKey2)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	nozPriceFactorsSeq6, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq5)
	require.True(t, nozPricePercentage.Equal(sdkmath.LegacyZeroDec()), "noz price shouldn't change after CreateResourceNode")
	require.True(t, ozoneLimitPercentage.Equal(sdkmath.LegacyZeroDec()), "OzLimit shouldn't change  after CreateResourceNode")
	t.Log("********************************* Deliver CreateResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver CreateResourceNode Tx START ********************************************")
	createResourceNodeMsg3 := setupMsgCreateResourceNode3()
	/********************* deliver tx *********************/

	senderAcc = accountKeeper.GetAccount(ctx, resOwner3)
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()

	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg3}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, resOwnerPrivKey3)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	nozPriceFactorsSeq7, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq6)
	require.True(t, nozPricePercentage.Equal(sdkmath.LegacyZeroDec()), "noz price shouldn't change after CreateResourceNode")
	require.True(t, ozoneLimitPercentage.Equal(sdkmath.LegacyZeroDec()), "OzLimit shouldn't change  after CreateResourceNode")
	t.Log("********************************* Deliver CreateResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver CreateResourceNode Tx START ********************************************")
	createResourceNodeMsg4 := setupMsgCreateResourceNode4()
	/********************* deliver tx *********************/

	senderAcc = accountKeeper.GetAccount(ctx, resOwner4)
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()

	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg4}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, resOwnerPrivKey4)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	nozPriceFactorsSeq8, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq7)
	require.True(t, nozPricePercentage.Equal(sdkmath.LegacyZeroDec()), "noz price shouldn't change after CreateResourceNode")
	require.True(t, ozoneLimitPercentage.Equal(sdkmath.LegacyZeroDec()), "OzLimit shouldn't change  after CreateResourceNode")
	t.Log("********************************* Deliver CreateResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver CreateResourceNode Tx START ********************************************")
	createResourceNodeMsg5 := setupMsgCreateResourceNode5()
	/********************* deliver tx *********************/

	senderAcc = accountKeeper.GetAccount(ctx, resOwner5)
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()

	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg5}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, resOwnerPrivKey5)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	_, nozPricePercentage, ozoneLimitPercentage = printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq8)
	require.True(t, nozPricePercentage.Equal(sdkmath.LegacyZeroDec()), "noz price shouldn't change after CreateResourceNode")
	require.True(t, ozoneLimitPercentage.Equal(sdkmath.LegacyZeroDec()), "OzLimit shouldn't change  after CreateResourceNode")
	t.Log("********************************* Deliver CreateResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

}

// initialize data of volume report
func setupMsgVolumeReport(t *testing.T, newEpoch int64) *pottypes.MsgVolumeReport {
	volume1 := pottypes.NewSingleWalletVolume(resOwner1, resourceNodeVolume1)
	volume2 := pottypes.NewSingleWalletVolume(resOwner2, resourceNodeVolume2)
	volume3 := pottypes.NewSingleWalletVolume(resOwner3, resourceNodeVolume3)

	nodesVolume := []pottypes.SingleWalletVolume{volume1, volume2, volume3}
	reporter := metaNodeP2PAddr1
	epoch := sdkmath.NewInt(newEpoch)
	reportReference := "report for epoch " + epoch.String()
	reporterOwner := metaOwner1

	volumeReportMsg := pottypes.NewMsgVolumeReport(nodesVolume, reporter, epoch, reportReference, reporterOwner)

	volumeReportMsg, err := stratostestutil.SignVolumeReport(
		volumeReportMsg,
		metaNodeP2PPrivKey1.Bytes(),
		metaNodeP2PPrivKey2.Bytes(),
		metaNodeP2PPrivKey3.Bytes(),
	)
	require.NoError(t, err)

	return volumeReportMsg
}

func setupSlashingMsg() *pottypes.MsgSlashingResourceNode {
	reporters := make([]stratos.SdsAddress, 0)
	reporters = append(reporters, metaNodeP2PAddr1)
	reportOwner := make([]sdk.AccAddress, 0)
	reportOwner = append(reportOwner, metaOwner1)
	slashingMsg := pottypes.NewMsgSlashingResourceNode(reporters, reportOwner, resNodeP2PAddr1, resOwner1, resNodeSlashingNOZAmt1, true)
	return slashingMsg
}

func setupSuspendMsgByIndex(i int, resNodeNetworkId stratos.SdsAddress, resNodePubKey cryptotypes.PubKey, resOwner sdk.AccAddress) *pottypes.MsgSlashingResourceNode {
	reporters := make([]stratos.SdsAddress, 0)
	reporters = append(reporters, metaNodeP2PAddr1)
	reportOwner := make([]sdk.AccAddress, 0)
	reportOwner = append(reportOwner, metaOwner1)
	slashingMsg := pottypes.NewMsgSlashingResourceNode(reporters, reportOwner, resNodeNetworkId, resOwner, resNodeSlashingNOZAmt1, true)
	return slashingMsg
}

func setupUnsuspendMsg() *pottypes.MsgSlashingResourceNode {
	reporters := make([]stratos.SdsAddress, 0)
	reporters = append(reporters, metaNodeP2PAddr1)
	reportOwner := make([]sdk.AccAddress, 0)
	reportOwner = append(reportOwner, metaOwner1)
	slashingMsg := pottypes.NewMsgSlashingResourceNode(reporters, reportOwner, resNodeP2PAddr1, resOwner1, sdkmath.ZeroInt(), false)
	return slashingMsg
}
func setupPrepayMsg() *sdstypes.MsgPrepay {
	sender := resOwner1
	amount := sdkmath.NewInt(1).MulRaw(stratos.StosToWei)
	coin := sdk.NewCoin(stratos.Wei, amount)
	prepayMsg := sdstypes.NewMsgPrepay(sender.String(), sender.String(), sdk.NewCoins(coin))
	return prepayMsg
}

func setupPrepayMsgWithResOwner(resOwner sdk.AccAddress) *sdstypes.MsgPrepay {
	sender := resOwner
	amount := sdkmath.NewInt(3).MulRaw(stratos.StosToWei)
	coin := sdk.NewCoin(stratos.Wei, amount)
	prepayMsg := sdstypes.NewMsgPrepay(sender.String(), sender.String(), sdk.NewCoins(coin))
	return prepayMsg
}

func setupMsgRemoveResourceNode(i int, resNodeNetworkId stratos.SdsAddress, resOwner sdk.AccAddress) *registertypes.MsgRemoveResourceNode {
	removeResourceNodeMsg := registertypes.NewMsgRemoveResourceNode(resNodeNetworkId, resOwner)
	return removeResourceNodeMsg
}
func setupMsgCreateResourceNode(i int, resNodeNetworkId stratos.SdsAddress, resNodePubKey cryptotypes.PubKey, resOwner sdk.AccAddress) *registertypes.MsgCreateResourceNode {
	nodeType := uint32(registertypes.STORAGE)
	createResourceNodeMsg, _ := registertypes.NewMsgCreateResourceNode(resNodeNetworkId, resNodePubKey, sdk.NewCoin(stratos.Wei, resNodeInitialDeposit), resOwner, registertypes.NewDescription("sds://resourceNode"+strconv.Itoa(i+1), "", "", "", ""), nodeType)
	return createResourceNodeMsg
}

func setupUnsuspendMsgByIndex(i int, resNodeNetworkId stratos.SdsAddress, resOwner sdk.AccAddress) *pottypes.MsgSlashingResourceNode {
	reporters := make([]stratos.SdsAddress, 0)
	reporters = append(reporters, metaNodeP2PAddr1)
	reportOwner := make([]sdk.AccAddress, 0)
	reportOwner = append(reportOwner, metaOwner1)
	slashingMsg := pottypes.NewMsgSlashingResourceNode(reporters, reportOwner, resNodeNetworkId, resOwner, sdkmath.ZeroInt(), false)
	return slashingMsg
}

func setupMsgCreateResourceNode1() *registertypes.MsgCreateResourceNode {
	nodeType := uint32(registertypes.STORAGE)
	createResourceNodeMsg, _ := registertypes.NewMsgCreateResourceNode(resNodeP2PAddr1, resNodeP2PPubKey1, sdk.NewCoin(stratos.Wei, resNodeInitialDeposit), resOwner1, registertypes.NewDescription("sds://resourceNode1", "", "", "", ""), nodeType)
	return createResourceNodeMsg
}
func setupMsgCreateResourceNode2() *registertypes.MsgCreateResourceNode {
	nodeType := uint32(registertypes.STORAGE)
	createResourceNodeMsg, _ := registertypes.NewMsgCreateResourceNode(resNodeP2PAddr2, resNodeP2PPubKey2, sdk.NewCoin(stratos.Wei, resNodeInitialDeposit), resOwner2, registertypes.NewDescription("sds://resourceNode2", "", "", "", ""), nodeType)
	return createResourceNodeMsg
}
func setupMsgCreateResourceNode3() *registertypes.MsgCreateResourceNode {
	nodeType := uint32(registertypes.STORAGE)
	createResourceNodeMsg, _ := registertypes.NewMsgCreateResourceNode(resNodeP2PAddr3, resNodeP2PPubKey3, sdk.NewCoin(stratos.Wei, resNodeInitialDeposit), resOwner3, registertypes.NewDescription("sds://resourceNode3", "", "", "", ""), nodeType)
	return createResourceNodeMsg
}

func setupMsgCreateResourceNode4() *registertypes.MsgCreateResourceNode {
	nodeType := uint32(registertypes.STORAGE)
	createResourceNodeMsg, _ := registertypes.NewMsgCreateResourceNode(resNodeP2PAddr4, resNodeP2PPubKey4, sdk.NewCoin(stratos.Wei, resNodeInitialDeposit), resOwner4, registertypes.NewDescription("sds://resourceNode4", "", "", "", ""), nodeType)
	return createResourceNodeMsg
}

func setupMsgCreateResourceNode5() *registertypes.MsgCreateResourceNode {
	nodeType := uint32(registertypes.STORAGE)
	createResourceNodeMsg, _ := registertypes.NewMsgCreateResourceNode(resNodeP2PAddr5, resNodeP2PPubKey5, sdk.NewCoin(stratos.Wei, resNodeInitialDeposit), resOwner5, registertypes.NewDescription("sds://resourceNode5", "", "", "", ""), nodeType)
	return createResourceNodeMsg
}

func setupMsgUpdateEffectiveDeposit1() *registertypes.MsgUpdateEffectiveDeposit {
	reporters := make([]stratos.SdsAddress, 0)
	reporters = append(reporters, metaNodeP2PAddr1)
	reporterOwner := make([]sdk.AccAddress, 0)
	reporterOwner = append(reporterOwner, metaOwner1)
	msg := registertypes.NewMsgUpdateEffectiveDeposit(reporters, reporterOwner, resNodeP2PAddr1, resNodeInitialDeposit)
	return msg
}

func setupMsgUpdateEffectiveDeposit(resNodeP2PAddr stratos.SdsAddress) *registertypes.MsgUpdateEffectiveDeposit {
	reporters := make([]stratos.SdsAddress, 0)
	reporters = append(reporters, metaNodeP2PAddr1)
	reporterOwner := make([]sdk.AccAddress, 0)
	reporterOwner = append(reporterOwner, metaOwner1)
	msg := registertypes.NewMsgUpdateEffectiveDeposit(reporters, reporterOwner, resNodeP2PAddr, resNodeInitialDeposit)
	return msg
}

func printCurrNozPrice(t *testing.T, ctx sdk.Context, potKeeper potKeeper.Keeper, registerKeeper registerKeeper.Keeper, nozPriceFactorsBefore NozPriceFactors) (NozPriceFactors, sdkmath.LegacyDec, sdkmath.LegacyDec) {
	nozPriceFactorsAfter := NozPriceFactors{}
	nozPriceFactorsAfter.InitialTotalDeposit = registerKeeper.GetInitialGenesisDepositTotal(ctx)
	nozPriceFactorsAfter.EffectiveTotalDeposit = registerKeeper.GetEffectiveTotalDeposit(ctx)
	nozPriceFactorsAfter.TotalUnissuedPrepay = registerKeeper.GetTotalUnissuedPrepay(ctx).Amount
	nozPriceFactorsAfter.DepositAndPrepay = nozPriceFactorsAfter.InitialTotalDeposit.Add(nozPriceFactorsAfter.TotalUnissuedPrepay)
	nozPriceFactorsAfter.OzoneLimit = registerKeeper.GetRemainingOzoneLimit(ctx)
	St, Pt, Lt := registerKeeper.GetCurrNozPriceParams(ctx)
	nozPriceFactorsAfter.NOzonePrice = potKeeper.GetCurrentNozPrice(St, Pt, Lt)
	_, nozPriceFactorsAfter.NozSupply = potKeeper.NozSupply(ctx)

	nozPriceDelta := nozPriceFactorsAfter.NOzonePrice.Sub(nozPriceFactorsBefore.NOzonePrice)
	initialTotalDepositDelta := nozPriceFactorsAfter.InitialTotalDeposit.Sub(nozPriceFactorsBefore.InitialTotalDeposit)
	effectiveTotalDepositDelta := nozPriceFactorsAfter.EffectiveTotalDeposit.Sub(nozPriceFactorsBefore.EffectiveTotalDeposit)
	totalUnissuedPrepayDelta := nozPriceFactorsAfter.TotalUnissuedPrepay.Sub(nozPriceFactorsBefore.TotalUnissuedPrepay)
	depositAndPrepayDelta := nozPriceFactorsAfter.DepositAndPrepay.Sub(nozPriceFactorsBefore.DepositAndPrepay)
	ozoneLimitDelta := nozPriceFactorsAfter.OzoneLimit.Sub(nozPriceFactorsBefore.OzoneLimit)
	nozSupplyDelta := nozPriceFactorsAfter.NozSupply.Sub(nozPriceFactorsBefore.NozSupply)

	nozPricePercentage := nozPriceDelta.Quo(nozPriceFactorsBefore.NOzonePrice).MulInt(sdkmath.NewInt(100))
	//initialTotalDepositPercentage := initialTotalDepositDelta.Quo(nozPriceFactorsBefore.InitialTotalDeposit)
	//effectiveTotalDepositPercentage := effectiveTotalDepositDelta.Quo(nozPriceFactorsBefore.EffectiveTotalDeposit)
	//totalUnissuedPrepayPercentage := totalUnissuedPrepayDelta.Quo(nozPriceFactorsBefore.TotalUnissuedPrepay)
	//depositAndPrepayPercentage := depositAndPrepayDelta.Quo(nozPriceFactorsBefore.DepositAndPrepay)
	ozoneLimitPercentage := ozoneLimitDelta.ToLegacyDec().Quo(nozPriceFactorsBefore.OzoneLimit.ToLegacyDec()).MulInt(sdkmath.NewInt(100))

	t.Log("===>>>>>>>>>>>>>>     Current noz Price    ===>>>>>>>>>>>>>>")
	t.Log("NOzonePrice:                                   " + nozPriceFactorsAfter.NOzonePrice.String() + "(delta: " + nozPriceDelta.String() + ", " + nozPricePercentage.String()[:5] + "%)")
	t.Log("InitialTotalDeposit:                           " + nozPriceFactorsAfter.InitialTotalDeposit.String() + "(delta: " + initialTotalDepositDelta.String() + ")")
	t.Log("EffectiveTotalDeposit:                         " + nozPriceFactorsAfter.EffectiveTotalDeposit.String() + "(delta: " + effectiveTotalDepositDelta.String() + ")")
	t.Log("TotalUnissuedPrepay:                           " + nozPriceFactorsAfter.TotalUnissuedPrepay.String() + "(delta: " + totalUnissuedPrepayDelta.String() + ")")
	t.Log("InitialTotalDeposit+TotalUnissuedPrepay:       " + nozPriceFactorsAfter.DepositAndPrepay.String() + "(delta: " + depositAndPrepayDelta.String() + ")")
	t.Log("OzoneLimit:                                    " + nozPriceFactorsAfter.OzoneLimit.String() + "(delta: " + ozoneLimitDelta.String() + ", " + ozoneLimitPercentage.String()[:5] + "%)")
	t.Log("NozSupply:                                     " + nozPriceFactorsAfter.NozSupply.String() + "(delta: " + nozSupplyDelta.String() + ")")

	return nozPriceFactorsAfter, nozPricePercentage, ozoneLimitPercentage
}

// return : coins - slashing
func deductSlashingAmt(ctx sdk.Context, coins sdk.Coins, slashing sdkmath.Int) sdk.Coins {
	ret := sdk.Coins{}
	for _, coin := range coins {
		if coin.Amount.GTE(slashing) {
			coin = coin.Sub(sdk.NewCoin(coin.Denom, slashing))
			ret = ret.Add(coin)
			slashing = sdkmath.ZeroInt()
		} else {
			slashing = slashing.Sub(coin.Amount)
			coin = sdk.NewCoin(coin.Denom, sdkmath.ZeroInt())
			ret = ret.Add(coin)
		}
	}
	return ret
}

// for main net
func checkResult(t *testing.T, ctx sdk.Context,
	k potKeeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankKeeper.Keeper,
	registerKeeper registerKeeper.Keeper,
	currentEpoch sdkmath.Int,
	lastFoundationAccBalance sdk.Coins,
	lastUnissuedPrepay sdk.Coin,
	lastCommunityPool sdk.Coins,
	lastMatureTotalOfResNode1 sdk.Coins,
	slashingAmtSetup sdkmath.Int,
	feeCollectorToFeePoolAtBeginBlock sdk.Coin) {

	currentSlashing := registerKeeper.GetSlashing(ctx, resOwner2)
	t.Log("currentSlashing					= " + currentSlashing.String())

	individualRewardTotal := sdk.Coins{}
	newMatureEpoch := currentEpoch.Add(sdkmath.NewInt(k.MatureEpoch(ctx)))

	k.IteratorIndividualReward(ctx, newMatureEpoch, func(walletAddress sdk.AccAddress, individualReward pottypes.Reward) (stop bool) {
		individualRewardTotal = individualRewardTotal.Add(individualReward.RewardFromTrafficPool...).Add(individualReward.RewardFromMiningPool...)
		t.Log("individualReward of [" + walletAddress.String() + "] = " + individualReward.String())
		return false
	})

	feePoolAccAddr := accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	require.NotNil(t, feePoolAccAddr)
	foundationAccountAddr := accountKeeper.GetModuleAddress(pottypes.FoundationAccount)
	newFoundationAccBalance := bankKeeper.GetAllBalances(ctx, foundationAccountAddr)
	newUnissuedPrepay := sdk.NewCoins(registerKeeper.GetTotalUnissuedPrepay(ctx))

	slashingChange := slashingAmtSetup.Sub(registerKeeper.GetSlashing(ctx, resOwner1))
	t.Log("resource node 1 slashing change		= " + slashingChange.String())
	matureTotal := k.GetMatureTotalReward(ctx, resOwner1)
	immatureTotal := k.GetImmatureTotalReward(ctx, resOwner1)
	t.Log("resource node 1 matureTotal		= " + matureTotal.String())
	t.Log("resource node 1 immatureTotal		= " + immatureTotal.String())

	rewardSrcChange := lastFoundationAccBalance.
		Sub(newFoundationAccBalance...).
		Add(lastUnissuedPrepay).
		Sub(newUnissuedPrepay...)
	t.Log("rewardSrcChange				= " + rewardSrcChange.String())

	// distribution module will send all tokens from "fee_collector" to "distribution" account in the BeginBlocker() method
	feePoolValChange := bankKeeper.GetAllBalances(ctx, feePoolAccAddr)
	t.Log("reward send to validator fee pool	= " + feePoolValChange.String())

	rewardDestChange := feePoolValChange.Add(individualRewardTotal...)
	t.Log("rewardDestChange			= " + rewardDestChange.String())

	//require.Equal(t, rewardSrcChange, rewardDestChange)

	t.Log("************************ slashing test***********************************")
	t.Log("slashing change				= " + slashingChange.String())

	upcomingMaturedIndividual := sdk.Coins{}
	individualReward, found := k.GetIndividualReward(ctx, resOwner1, currentEpoch)
	if found {
		tmp := individualReward.RewardFromTrafficPool.Add(individualReward.RewardFromMiningPool...)
		upcomingMaturedIndividual = deductSlashingAmt(ctx, tmp, slashingChange)
	}
	t.Log("upcomingMaturedIndividual		= " + upcomingMaturedIndividual.String())

	// get mature total changes
	newMatureTotalOfResNode1 := k.GetMatureTotalReward(ctx, resOwner1)
	matureTotalOfResNode1Change, _ := newMatureTotalOfResNode1.SafeSub(lastMatureTotalOfResNode1...)
	if matureTotalOfResNode1Change == nil || matureTotalOfResNode1Change.IsAnyNegative() {
		matureTotalOfResNode1Change = sdk.Coins{}
	}
	t.Log("matureTotalOfResNode1Change		= " + matureTotalOfResNode1Change.String())
	require.Equal(t, matureTotalOfResNode1Change.String(), upcomingMaturedIndividual.String())
}

func checkValidator(t *testing.T, app *app.StratosApp, addr sdk.ValAddress, expFound bool) stakingtypes.Validator {
	ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})
	validator, found := app.GetStakingKeeper().GetValidator(ctxCheck, addr)

	require.Equal(t, expFound, found)
	return validator
}

func TestMain(m *testing.M) {
	config := stratos.GetConfig()
	config.Seal()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func setupAccounts() ([]authtypes.GenesisAccount, []banktypes.Balance) {

	//************************** setup resource nodes owners' accounts **************************
	resOwnerAcc1 := &authtypes.BaseAccount{Address: resOwner1.String()}
	resOwnerAcc2 := &authtypes.BaseAccount{Address: resOwner2.String()}
	resOwnerAcc3 := &authtypes.BaseAccount{Address: resOwner3.String()}
	resOwnerAcc4 := &authtypes.BaseAccount{Address: resOwner4.String()}
	resOwnerAcc5 := &authtypes.BaseAccount{Address: resOwner5.String()}
	//************************** setup meta nodes owners' accounts **************************
	metaOwnerAcc1 := &authtypes.BaseAccount{Address: metaOwner1.String()}
	metaOwnerAcc2 := &authtypes.BaseAccount{Address: metaOwner2.String()}
	metaOwnerAcc3 := &authtypes.BaseAccount{Address: metaOwner3.String()}
	//************************** setup validator delegators' accounts **************************
	valOwnerAcc1 := &authtypes.BaseAccount{Address: valOpAccAddr1.String()}
	//************************** setup meta nodes' accounts **************************
	foundationDepositorAcc := &authtypes.BaseAccount{Address: foundationDepositorAccAddr.String()}

	accs := []authtypes.GenesisAccount{
		resOwnerAcc1, resOwnerAcc2, resOwnerAcc3, resOwnerAcc4, resOwnerAcc5,
		metaOwnerAcc1, metaOwnerAcc2, metaOwnerAcc3,
		valOwnerAcc1,
		foundationDepositorAcc,
	}

	balances := []banktypes.Balance{
		{
			Address: resOwner1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialDeposit.Add(depositForSendingTx))},
		},
		{
			Address: resOwner2.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialDeposit)},
		},
		{
			Address: resOwner3.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialDeposit)},
		},
		{
			Address: resOwner4.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialDeposit)},
		},
		{
			Address: resOwner5.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialDeposit)},
		},
		{
			Address: metaOwner1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(metaNodeInitialDeposit)},
		},
		{
			Address: metaOwner2.String(),
			Coins:   sdk.Coins{stratos.NewCoin(metaNodeInitialDeposit)},
		},
		{
			Address: metaOwner3.String(),
			Coins:   sdk.Coins{stratos.NewCoin(metaNodeInitialDeposit)},
		},
		{
			Address: valOpAccAddr1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(valInitialStake)},
		},
		{
			Address: foundationDepositorAccAddr.String(),
			Coins:   foundationDeposit,
		},
	}
	return accs, balances
}

func setupAccountsMultipleResNodes(resOwners []sdk.AccAddress) ([]authtypes.GenesisAccount, []banktypes.Balance) {

	resOwnerAccs := make([]*authtypes.BaseAccount, 0, len(resOwners))
	//************************** setup resource nodes owners' accounts **************************
	for _, resOwner := range resOwners {
		resOwnerAccs = append(resOwnerAccs, &authtypes.BaseAccount{Address: resOwner.String()})
	}
	//************************** setup meta nodes owners' accounts **************************
	idxOwnerAcc1 := &authtypes.BaseAccount{Address: metaOwner1.String()}
	idxOwnerAcc2 := &authtypes.BaseAccount{Address: metaOwner2.String()}
	idxOwnerAcc3 := &authtypes.BaseAccount{Address: metaOwner3.String()}
	//************************** setup validator delegators' accounts **************************
	valOwnerAcc1 := &authtypes.BaseAccount{Address: valOpAccAddr1.String()}
	//************************** setup meta nodes' accounts **************************
	foundationDepositorAcc := &authtypes.BaseAccount{Address: foundationDepositorAccAddr.String()}

	accs := []authtypes.GenesisAccount{
		idxOwnerAcc1, idxOwnerAcc2, idxOwnerAcc3,
		valOwnerAcc1,
		foundationDepositorAcc,
	}
	for _, resAcc := range resOwnerAccs {
		accs = append(accs, resAcc)
	}

	balances := []banktypes.Balance{
		{
			Address: metaOwner1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(metaNodeInitialDeposit)},
		},
		{
			Address: metaOwner2.String(),
			Coins:   sdk.Coins{stratos.NewCoin(metaNodeInitialDeposit)},
		},
		{
			Address: metaOwner3.String(),
			Coins:   sdk.Coins{stratos.NewCoin(metaNodeInitialDeposit)},
		},
		{
			Address: valOpAccAddr1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(valInitialStake)},
		},
		{
			Address: foundationDepositorAccAddr.String(),
			Coins:   foundationDeposit,
		},
	}

	for _, resOwnerAcc := range resOwnerAccs {
		accs = append(accs, resOwnerAcc)
		balances = append(balances, banktypes.Balance{
			Address: resOwnerAcc.Address,
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialDeposit.Add(depositForSendingTx))},
		})
	}

	return accs, balances
}

func setupAllResourceNodes() []registertypes.ResourceNode {

	createTime, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	nodeType := registertypes.STORAGE
	resourceNode1, _ := registertypes.NewResourceNode(resNodeP2PAddr1, resNodeP2PPubKey1, resOwner1, registertypes.NewDescription("sds://resourceNode1", "", "", "", ""), nodeType, createTime)
	resourceNode2, _ := registertypes.NewResourceNode(resNodeP2PAddr2, resNodeP2PPubKey2, resOwner2, registertypes.NewDescription("sds://resourceNode2", "", "", "", ""), nodeType, createTime)
	resourceNode3, _ := registertypes.NewResourceNode(resNodeP2PAddr3, resNodeP2PPubKey3, resOwner3, registertypes.NewDescription("sds://resourceNode3", "", "", "", ""), nodeType, createTime)
	resourceNode4, _ := registertypes.NewResourceNode(resNodeP2PAddr4, resNodeP2PPubKey4, resOwner4, registertypes.NewDescription("sds://resourceNode4", "", "", "", ""), nodeType, createTime)
	resourceNode5, _ := registertypes.NewResourceNode(resNodeP2PAddr5, resNodeP2PPubKey5, resOwner5, registertypes.NewDescription("sds://resourceNode5", "", "", "", ""), nodeType, createTime)

	resourceNode1 = resourceNode1.AddToken(resNodeInitialDeposit)
	resourceNode2 = resourceNode2.AddToken(resNodeInitialDeposit)
	resourceNode3 = resourceNode3.AddToken(resNodeInitialDeposit)
	resourceNode4 = resourceNode4.AddToken(resNodeInitialDeposit)
	resourceNode5 = resourceNode5.AddToken(resNodeInitialDeposit)

	resourceNode1.EffectiveTokens = resNodeInitialDeposit
	resourceNode2.EffectiveTokens = resNodeInitialDeposit
	resourceNode3.EffectiveTokens = resNodeInitialDeposit
	resourceNode4.EffectiveTokens = resNodeInitialDeposit
	resourceNode5.EffectiveTokens = resNodeInitialDeposit

	resourceNode1.Status = stakingtypes.Bonded
	resourceNode2.Status = stakingtypes.Bonded
	resourceNode3.Status = stakingtypes.Bonded
	resourceNode4.Status = stakingtypes.Bonded
	resourceNode5.Status = stakingtypes.Bonded

	resourceNode1.Suspend = false
	resourceNode2.Suspend = false
	resourceNode3.Suspend = false
	resourceNode4.Suspend = false
	resourceNode5.Suspend = false

	var resourceNodes []registertypes.ResourceNode
	resourceNodes = append(resourceNodes, resourceNode1)
	resourceNodes = append(resourceNodes, resourceNode2)
	resourceNodes = append(resourceNodes, resourceNode3)
	resourceNodes = append(resourceNodes, resourceNode4)
	resourceNodes = append(resourceNodes, resourceNode5)
	return resourceNodes
}

func setupMultipleResourceNodes(resNodeP2PAddresses []stratos.SdsAddress, resNodeP2PPubKeys []cryptotypes.PubKey, resOwners []sdk.AccAddress) []registertypes.ResourceNode {
	if len(resNodeP2PPubKeys) != len(resOwners) ||
		len(resOwners) != len(resNodeP2PAddresses) {
		return nil
	}

	numOfNodes := len(resNodeP2PPubKeys)
	resourceNodes := make([]registertypes.ResourceNode, 0, numOfNodes)

	createTime, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	nodeType := registertypes.STORAGE

	for i, _ := range resNodeP2PPubKeys {
		resourceNodeTmp, _ := registertypes.NewResourceNode(resNodeP2PAddresses[i], resNodeP2PPubKeys[i], resOwners[i], registertypes.NewDescription("sds://resourceNode"+strconv.Itoa(i+1), "", "", "", ""), nodeType, createTime)
		resourceNodeTmp = resourceNodeTmp.AddToken(resNodeInitialDeposit)
		resourceNodeTmp.Status = stakingtypes.Bonded
		resourceNodeTmp.EffectiveTokens = resNodeInitialDeposit
		resourceNodeTmp.Suspend = false
		resourceNodes = append(resourceNodes, resourceNodeTmp)
	}

	return resourceNodes
}

func setupAllMetaNodes() []registertypes.MetaNode {
	var metaNodes []registertypes.MetaNode

	createTime, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	metaNode1, _ := registertypes.NewMetaNode(metaNodeP2PAddr1, metaNodeP2PPubKey1, metaOwner1, metaOwner1, registertypes.NewDescription("sds://metaNode1", "", "", "", ""), createTime)
	metaNode2, _ := registertypes.NewMetaNode(metaNodeP2PAddr2, metaNodeP2PPubKey2, metaOwner2, metaOwner2, registertypes.NewDescription("sds://metaNode2", "", "", "", ""), createTime)
	metaNode3, _ := registertypes.NewMetaNode(metaNodeP2PAddr3, metaNodeP2PPubKey3, metaOwner3, metaOwner3, registertypes.NewDescription("sds://metaNode3", "", "", "", ""), createTime)

	metaNode1 = metaNode1.AddToken(metaNodeInitialDeposit)
	metaNode2 = metaNode2.AddToken(metaNodeInitialDeposit)
	metaNode3 = metaNode3.AddToken(metaNodeInitialDeposit)

	metaNode1.Status = stakingtypes.Bonded
	metaNode2.Status = stakingtypes.Bonded
	metaNode3.Status = stakingtypes.Bonded

	metaNode1.Suspend = false
	metaNode2.Suspend = false
	metaNode3.Suspend = false

	metaNodes = append(metaNodes, metaNode1)
	metaNodes = append(metaNodes, metaNode2)
	metaNodes = append(metaNodes, metaNode3)

	return metaNodes
}

var (
	initialTotalDepositStore   = sdkmath.NewInt(1500000000000)
	effectiveTotalDepositStore = sdkmath.NewInt(1500000000000)
	remainOzoneLimitStore      = sdkmath.NewInt(1500000000000000)
	totalUnissuedPrepayStore   = sdkmath.ZeroInt()
	nozPrice                   = sdkmath.LegacyNewDecWithPrec(1000000, 9)
	//priceChangeChan = make(chan PriceChangeEvent, 0)
)

type PriceChangeEvent struct {
	depositDelta        sdkmath.Int
	unissuedPrepayDelta sdkmath.Int
}

func simulatePriceChange(t *testing.T, priceChangeEvent *PriceChangeEvent, nozPriceFactorsBefore *NozPriceFactors) (*NozPriceFactors, sdkmath.LegacyDec, sdkmath.LegacyDec) {
	nozPriceFactorsAfter := &NozPriceFactors{}
	nozPriceFactorsAfter.InitialTotalDeposit = nozPriceFactorsBefore.InitialTotalDeposit
	nozPriceFactorsAfter.TotalUnissuedPrepay = nozPriceFactorsBefore.TotalUnissuedPrepay.Add(priceChangeEvent.unissuedPrepayDelta)
	nozPriceFactorsAfter.DepositAndPrepay = nozPriceFactorsAfter.InitialTotalDeposit.Add(nozPriceFactorsAfter.TotalUnissuedPrepay)
	nozPriceFactorsAfter.EffectiveTotalDeposit = nozPriceFactorsBefore.EffectiveTotalDeposit.Add(priceChangeEvent.depositDelta)
	deltaNozLimit := sdkmath.ZeroInt()
	nozPriceFactorsAfter.NozSupply = nozPriceFactorsBefore.NozSupply
	if !priceChangeEvent.depositDelta.Equal(sdkmath.ZeroInt()) {
		ozoneLimitChangeByDeposit := priceChangeEvent.depositDelta.ToLegacyDec().Quo(depositNozRateInt.ToLegacyDec()).TruncateInt()
		//ozoneLimitChangeByDeposit := nozPriceFactorsBefore.OzoneLimit.ToLegacyDec().Quo(nozPriceFactorsBefore.InitialTotalDeposit.ToLegacyDec()).Mul(priceChangeEvent.depositDelta.ToDec()).TruncateInt()
		deltaNozLimit = deltaNozLimit.Add(ozoneLimitChangeByDeposit)
		nozPriceFactorsAfter.NozSupply = nozPriceFactorsBefore.NozSupply.Add(ozoneLimitChangeByDeposit)
	}
	if !priceChangeEvent.unissuedPrepayDelta.Equal(sdkmath.ZeroInt()) {
		ozoneLimitChangeByPrepay := nozPriceFactorsBefore.OzoneLimit.ToLegacyDec().
			Mul(priceChangeEvent.unissuedPrepayDelta.ToLegacyDec()).
			Quo(nozPriceFactorsBefore.EffectiveTotalDeposit.Add(nozPriceFactorsBefore.TotalUnissuedPrepay).Add(priceChangeEvent.unissuedPrepayDelta).ToLegacyDec()).
			TruncateInt()
		//Sub(nozPriceFactorsBefore.OzoneLimit)
		if priceChangeEvent.unissuedPrepayDelta.GT(sdkmath.ZeroInt()) {
			// positive value of prepay leads to limit decrease
			deltaNozLimit = deltaNozLimit.Sub(ozoneLimitChangeByPrepay)
		} else {
			// nagative value of prepay (reward distribution) leads to limit increase
			deltaNozLimit = deltaNozLimit.Add(ozoneLimitChangeByPrepay)
		}
	}

	nozPriceFactorsAfter.OzoneLimit = nozPriceFactorsBefore.OzoneLimit.Add(deltaNozLimit)

	nozPriceFactorsAfter.NOzonePrice = nozPriceFactorsAfter.DepositAndPrepay.ToLegacyDec().Quo(nozPriceFactorsAfter.OzoneLimit.ToLegacyDec())
	nozPriceFactorsAfter.EffectiveTotalDeposit = nozPriceFactorsBefore.EffectiveTotalDeposit.Add(priceChangeEvent.depositDelta)

	nozPriceDelta := nozPriceFactorsAfter.NOzonePrice.Sub(nozPriceFactorsBefore.NOzonePrice)
	initialTotalDepositDelta := nozPriceFactorsAfter.InitialTotalDeposit.Sub(nozPriceFactorsBefore.InitialTotalDeposit)
	effectiveTotalDepositDelta := nozPriceFactorsAfter.EffectiveTotalDeposit.Sub(nozPriceFactorsBefore.EffectiveTotalDeposit)
	totalUnissuedPrepayDelta := nozPriceFactorsAfter.TotalUnissuedPrepay.Sub(nozPriceFactorsBefore.TotalUnissuedPrepay)
	depositAndPrepayDelta := nozPriceFactorsAfter.DepositAndPrepay.Sub(nozPriceFactorsBefore.DepositAndPrepay)
	ozoneLimitDelta := nozPriceFactorsAfter.OzoneLimit.Sub(nozPriceFactorsBefore.OzoneLimit)
	nozSupplyDelta := nozPriceFactorsAfter.NozSupply.Sub(nozPriceFactorsBefore.NozSupply)

	nozPricePercentage := nozPriceDelta.Quo(nozPriceFactorsBefore.NOzonePrice).MulInt(sdkmath.NewInt(100))
	//initialTotalDepositPercentage := initialTotalDepositDelta.Quo(nozPriceFactorsBefore.InitialTotalDeposit)
	//effectiveTotalDepositPercentage := effectiveTotalDepositDelta.Quo(nozPriceFactorsBefore.EffectiveTotalDeposit)
	//totalUnissuedPrepayPercentage := totalUnissuedPrepayDelta.Quo(nozPriceFactorsBefore.TotalUnissuedPrepay)
	//depositAndPrepayPercentage := depositAndPrepayDelta.Quo(nozPriceFactorsBefore.DepositAndPrepay)
	ozoneLimitPercentage := ozoneLimitDelta.ToLegacyDec().Quo(nozPriceFactorsBefore.OzoneLimit.ToLegacyDec()).MulInt(sdkmath.NewInt(100))

	t.Log("===>>>>>>>>>>>>>>     Current noz Price    ===>>>>>>>>>>>>>>")
	t.Log("NOzonePrice:                                   " + nozPriceFactorsAfter.NOzonePrice.String() + "(delta: " + nozPriceDelta.String() + ", " + nozPricePercentage.String()[:5] + "%)")
	t.Log("InitialTotalDeposit:                           " + nozPriceFactorsAfter.InitialTotalDeposit.String() + "(delta: " + initialTotalDepositDelta.String() + ")")
	t.Log("EffectiveTotalDeposit:                         " + nozPriceFactorsAfter.EffectiveTotalDeposit.String() + "(delta: " + effectiveTotalDepositDelta.String() + ")")
	t.Log("TotalUnissuedPrepay:                           " + nozPriceFactorsAfter.TotalUnissuedPrepay.String() + "(delta: " + totalUnissuedPrepayDelta.String() + ")")
	t.Log("InitialTotalDeposit+TotalUnissuedPrepay:       " + nozPriceFactorsAfter.DepositAndPrepay.String() + "(delta: " + depositAndPrepayDelta.String() + ")")
	t.Log("OzoneLimit:                                    " + nozPriceFactorsAfter.OzoneLimit.String() + "(delta: " + ozoneLimitDelta.String() + ", " + ozoneLimitPercentage.String()[:5] + "%)")
	t.Log("NozSupply:                                     " + nozPriceFactorsAfter.NozSupply.String() + "(delta: " + nozSupplyDelta.String() + ")")

	return nozPriceFactorsAfter, nozPricePercentage, ozoneLimitPercentage
}

func TestOzPriceChangePrepay(t *testing.T) {
	NUM_OF_SAMPLE := 100
	dataToExcel := make([]NozPriceFactors, 0, NUM_OF_SAMPLE)
	/********************* initialize mock app *********************/
	accs, balances := setupAccounts()

	// create validator set with single validator
	consPubKey, err := cryptocodec.ToTmPubKeyInterface(valConsPubk1)
	validator := tmtypes.NewValidator(consPubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	metaNodes := setupAllMetaNodes()
	//resourceNodes := setupAllResourceNodes()
	resourceNodes := make([]registertypes.ResourceNode, 0)

	stApp := stratostestutil.SetupWithGenesisNodeSet(t, valSet, metaNodes, resourceNodes, accs, chainID, false, balances...)

	accountKeeper := stApp.GetAccountKeeper()
	//bankKeeper := stApp.GetBankKeeper()
	registerKeeper := stApp.GetRegisterKeeper()
	potKeeper := stApp.GetPotKeeper()

	/********************* foundation account deposit *********************/
	header := tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := stApp.BaseApp.NewContext(true, header)

	foundationDepositMsg := pottypes.NewMsgFoundationDeposit(foundationDeposit, foundationDepositorAccAddr)
	txGen := stratostestutil.MakeTestEncodingConfig().TxConfig

	senderAcc := accountKeeper.GetAccount(ctx, foundationDepositorAccAddr)
	accNum := senderAcc.GetAccountNumber()
	accSeq := senderAcc.GetSequence()
	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, foundationDepositorPrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(pottypes.FoundationAccount)
	stratostestutil.CheckBalance(t, stApp, foundationAccountAddr, foundationDeposit)

	/********************** commit **********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	_, nozSupply := potKeeper.NozSupply(ctx)
	St, Pt, Lt := registerKeeper.GetCurrNozPriceParams(ctx)
	nozPriceFactorsSeq0, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, NozPriceFactors{
		NOzonePrice:           potKeeper.GetCurrentNozPrice(St, Pt, Lt),
		InitialTotalDeposit:   registerKeeper.GetInitialGenesisDepositTotal(ctx),
		EffectiveTotalDeposit: registerKeeper.GetEffectiveTotalDeposit(ctx),
		TotalUnissuedPrepay:   registerKeeper.GetTotalUnissuedPrepay(ctx).Amount,
		DepositAndPrepay:      registerKeeper.GetInitialGenesisDepositTotal(ctx).Add(registerKeeper.GetTotalUnissuedPrepay(ctx).Amount),
		OzoneLimit:            registerKeeper.GetRemainingOzoneLimit(ctx),
		NozSupply:             nozSupply,
	})

	// start testing
	t.Log("\n********************************* Deliver Prepay Tx START ********************************************")

	priceBefore := nozPriceFactorsSeq0
	priceAfter := nozPriceFactorsSeq0
	dataToExcel = append(dataToExcel, priceBefore)

	for i := 0; i < NUM_OF_SAMPLE; i++ {
		prepayMsg := setupPrepayMsg()
		/********************* deliver tx *********************/

		senderAcc = accountKeeper.GetAccount(ctx, resOwner1)
		accNum = senderAcc.GetAccountNumber()
		accSeq = senderAcc.GetSequence()
		_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{prepayMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, resOwnerPrivKey1)
		require.NoError(t, err)
		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)
		priceAfter, nozPricePercentage, ozoneLimitPercentage = printCurrNozPrice(t, ctx, potKeeper, registerKeeper, priceBefore)
		dataToExcel = append(dataToExcel, priceAfter)
		require.True(t, nozPricePercentage.GT(sdkmath.LegacyZeroDec()), "noz price should increase after PREPAY")
		require.True(t, ozoneLimitPercentage.LT(sdkmath.LegacyZeroDec()), "OzLimit should not change after PREPAY")
		t.Log("********************************* Deliver Prepay Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

		priceBefore = priceAfter
	}
	exportToCSV(t, dataToExcel)
}

func TestOzPriceChangeVolumeReport(t *testing.T) {
	NUM_OF_SAMPLE := 100
	dataToExcel := make([]NozPriceFactors, 0, NUM_OF_SAMPLE)
	/********************* initialize mock app *********************/
	accs, balances := setupAccounts()

	// create validator set with single validator
	consPubKey, err := cryptocodec.ToTmPubKeyInterface(valConsPubk1)
	validator := tmtypes.NewValidator(consPubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	metaNodes := setupAllMetaNodes()
	//resourceNodes := setupAllResourceNodes()
	resourceNodes := make([]registertypes.ResourceNode, 0)

	stApp := stratostestutil.SetupWithGenesisNodeSet(t, valSet, metaNodes, resourceNodes, accs, chainID, false, balances...)

	accountKeeper := stApp.GetAccountKeeper()
	bankKeeper := stApp.GetBankKeeper()
	registerKeeper := stApp.GetRegisterKeeper()
	potKeeper := stApp.GetPotKeeper()
	distrKeeper := stApp.GetDistrKeeper()

	/********************* foundation account deposit *********************/
	header := tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := stApp.BaseApp.NewContext(true, header)

	foundationDepositMsg := pottypes.NewMsgFoundationDeposit(foundationDeposit, foundationDepositorAccAddr)
	txGen := stratostestutil.MakeTestEncodingConfig().TxConfig

	senderAcc := accountKeeper.GetAccount(ctx, foundationDepositorAccAddr)
	accNum := senderAcc.GetAccountNumber()
	accSeq := senderAcc.GetSequence()

	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, foundationDepositorPrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(pottypes.FoundationAccount)
	stratostestutil.CheckBalance(t, stApp, foundationAccountAddr, foundationDeposit)

	/********************** commit **********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	_, nozSupply := potKeeper.NozSupply(ctx)
	St, Pt, Lt := registerKeeper.GetCurrNozPriceParams(ctx)
	nozPriceFactorsSeq0, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, NozPriceFactors{
		NOzonePrice:           potKeeper.GetCurrentNozPrice(St, Pt, Lt),
		InitialTotalDeposit:   registerKeeper.GetInitialGenesisDepositTotal(ctx),
		EffectiveTotalDeposit: registerKeeper.GetEffectiveTotalDeposit(ctx),
		TotalUnissuedPrepay:   registerKeeper.GetTotalUnissuedPrepay(ctx).Amount,
		DepositAndPrepay:      registerKeeper.GetInitialGenesisDepositTotal(ctx).Add(registerKeeper.GetTotalUnissuedPrepay(ctx).Amount),
		OzoneLimit:            registerKeeper.GetRemainingOzoneLimit(ctx),
		NozSupply:             nozSupply,
	})

	// start testing
	t.Log("\n********************************* Deliver Prepay Tx START ********************************************")

	priceBefore := nozPriceFactorsSeq0
	priceAfter := nozPriceFactorsSeq0
	dataToExcel = append(dataToExcel, priceBefore)

	for i := 0; i < NUM_OF_SAMPLE; i++ {
		prepayMsg := setupPrepayMsg()
		/********************* deliver tx *********************/

		senderAcc = accountKeeper.GetAccount(ctx, resOwner1)
		accNum = senderAcc.GetAccountNumber()
		accSeq = senderAcc.GetSequence()
		_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{prepayMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, resOwnerPrivKey1)
		require.NoError(t, err)
		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)
		priceAfter, nozPricePercentage, ozoneLimitPercentage = printCurrNozPrice(t, ctx, potKeeper, registerKeeper, priceBefore)
		dataToExcel = append(dataToExcel, priceAfter)
		require.True(t, nozPricePercentage.GT(sdkmath.LegacyZeroDec()), "noz price should increase after PREPAY")
		require.True(t, ozoneLimitPercentage.LT(sdkmath.LegacyZeroDec()), "OzLimit should decrease after PREPAY")
		t.Log("********************************* Deliver Prepay Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

		priceBefore = priceAfter
	}

	for i := 0; i < NUM_OF_SAMPLE; i++ {
		t.Log("********************************* Deliver VolumeReport Tx START ********************************************")
		/********************* prepare tx data *********************/
		volumeReportMsg := setupMsgVolumeReport(t, int64(i+1))

		lastTotalMinedToken := potKeeper.GetTotalMinedTokens(ctx)
		t.Log("last committed TotalMinedTokens = " + lastTotalMinedToken.String())
		_, ok := sdkmath.NewIntFromString(volumeReportMsg.Epoch.String())
		require.Equal(t, ok, true)
		totalConsumedNoz := potKeeper.GetTotalConsumedNoz(volumeReportMsg.WalletVolumes).ToLegacyDec()

		/********************* print info *********************/
		t.Log("epoch " + volumeReportMsg.Epoch.String())
		S := registerKeeper.GetInitialGenesisDepositTotal(ctx).ToLegacyDec()
		Pt := registerKeeper.GetTotalUnissuedPrepay(ctx).Amount.ToLegacyDec()
		Y := totalConsumedNoz
		Lt := registerKeeper.GetRemainingOzoneLimit(ctx).ToLegacyDec()
		R := S.Add(Pt).Mul(Y).Quo(Lt.Add(Y))
		//t.Log("R = (S + Pt) * Y / (Lt + Y)")
		t.Log("S=" + S.String() + "\nPt=" + Pt.String() + "\nY=" + Y.String() + "\nLt=" + Lt.String() + "\nR=" + R.String() + "\n")

		t.Log("---------------------------")
		potKeeper.InitVariable(ctx)
		distributeGoal := pottypes.InitDistributeGoal()
		distributeGoal, err := potKeeper.CalcTrafficRewardInTotal(ctx, distributeGoal, totalConsumedNoz)
		require.NoError(t, err)

		distributeGoal, err = potKeeper.CalcMiningRewardInTotal(ctx, distributeGoal) //for main net
		require.NoError(t, err)
		t.Log(distributeGoal.String())

		t.Log("---------------------------")
		t.Log("distribute detail:")
		rewardDetailMap := make(map[string]pottypes.Reward)
		rewardDetailMap = potKeeper.CalcRewardForResourceNode(ctx, totalConsumedNoz, volumeReportMsg.WalletVolumes, distributeGoal, rewardDetailMap)
		rewardDetailMap = potKeeper.CalcRewardForMetaNode(ctx, distributeGoal, rewardDetailMap)

		t.Log("resource_wallet1:  address = " + resOwner1.String())
		t.Log("              miningReward = " + rewardDetailMap[resOwner1.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[resOwner1.String()].RewardFromTrafficPool.String())

		t.Log("resource_wallet2:  address = " + resOwner2.String())
		t.Log("              miningReward = " + rewardDetailMap[resOwner2.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[resOwner2.String()].RewardFromTrafficPool.String())

		t.Log("resource_wallet3:  address = " + resOwner3.String())
		t.Log("              miningReward = " + rewardDetailMap[resOwner3.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[resOwner3.String()].RewardFromTrafficPool.String())

		t.Log("resource_wallet4:  address = " + resOwner4.String())
		t.Log("              miningReward = " + rewardDetailMap[resOwner4.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[resOwner4.String()].RewardFromTrafficPool.String())

		t.Log("resource_wallet5:  address = " + resOwner5.String())
		t.Log("              miningReward = " + rewardDetailMap[resOwner5.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[resOwner5.String()].RewardFromTrafficPool.String())

		t.Log("meta_wallet1:      address = " + metaOwner1.String())
		t.Log("              miningReward = " + rewardDetailMap[metaOwner1.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[metaOwner1.String()].RewardFromTrafficPool.String())

		t.Log("meta_wallet2:      address = " + metaOwner2.String())
		t.Log("              miningReward = " + rewardDetailMap[metaOwner2.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[metaOwner2.String()].RewardFromTrafficPool.String())

		t.Log("meta_wallet3:      address = " + metaOwner3.String())
		t.Log("              miningReward = " + rewardDetailMap[metaOwner3.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[metaOwner3.String()].RewardFromTrafficPool.String())
		t.Log("---------------------------")

		/********************* record data before delivering tx  *********************/
		_ = bankKeeper.GetAllBalances(ctx, foundationAccountAddr)
		_ = registerKeeper.GetTotalUnissuedPrepay(ctx)
		_ = sdk.NewCoins(sdk.NewCoin(potKeeper.BondDenom(ctx), distrKeeper.GetFeePool(ctx).CommunityPool.AmountOf(potKeeper.BondDenom(ctx)).TruncateInt()))
		_ = potKeeper.GetMatureTotalReward(ctx, resOwner1)
		//lastFoundationAccBalance := bankKeeper.GetAllBalances(ctx, foundationAccountAddr)
		//lastUnissuedPrepay := registerKeeper.GetTotalUnissuedPrepay(ctx)
		//lastCommunityPool := sdk.NewCoins(sdk.NewCoin(potKeeper.BondDenom(ctx), potKeeper.DistrKeeper.GetFeePool(ctx).CommunityPool.AmountOf(potKeeper.BondDenom(ctx)).TruncateInt()))
		//lastMatureTotalOfResNode1 := potKeeper.GetMatureTotalReward(ctx, resOwner1)

		feePoolAccAddr := accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		require.NotNil(t, feePoolAccAddr)
		_ = bankKeeper.GetBalance(ctx, feePoolAccAddr, potKeeper.BondDenom(ctx))
		//feeCollectorToFeePoolAtBeginBlock := bankKeeper.GetBalance(ctx, feePoolAccAddr, potKeeper.BondDenom(ctx))

		senderAcc = accountKeeper.GetAccount(ctx, metaOwner1)
		accNum = senderAcc.GetAccountNumber()
		accSeq = senderAcc.GetSequence()

		_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{volumeReportMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, metaOwnerPrivKey1)
		require.NoError(t, err)
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		stApp.EndBlock(abci.RequestEndBlock{Height: header.Height})
		stApp.Commit()

		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)
		priceAfter, nozPricePercentage, ozoneLimitPercentage = printCurrNozPrice(t, ctx, potKeeper, registerKeeper, priceBefore)
		dataToExcel = append(dataToExcel, priceAfter)
		priceBefore = priceAfter
	}
	exportToCSV(t, dataToExcel)
}

func TestOzPriceChangeAddMultipleResourceNodeAndThenRemove(t *testing.T) {
	NUM_OF_SAMPLE := 100
	dataToExcel := make([]NozPriceFactors, 0, NUM_OF_SAMPLE)

	resOwnerPrivKeys := make([]*secp256k1.PrivKey, 0, NUM_OF_SAMPLE)
	resOwners := make([]sdk.AccAddress, 0, NUM_OF_SAMPLE)
	resNodeP2PPubkeys := make([]cryptotypes.PubKey, 0, NUM_OF_SAMPLE)
	resNodeP2PAddresses := make([]stratos.SdsAddress, 0, NUM_OF_SAMPLE)

	for i := 0; i < NUM_OF_SAMPLE; i++ {
		resOwnerPrivKeyTmp := secp256k1.GenPrivKey()
		resOwnerPrivKeys = append(resOwnerPrivKeys, resOwnerPrivKeyTmp)
		resOwnerAddrTmp := sdk.AccAddress(resOwnerPrivKeyTmp.PubKey().Address())
		resOwners = append(resOwners, resOwnerAddrTmp)

		resNodeP2PPubkeyTmp := ed25519.GenPrivKey().PubKey()
		resNodeP2PPubkeys = append(resNodeP2PPubkeys, resNodeP2PPubkeyTmp)
		resNodeP2PAddresses = append(resNodeP2PAddresses, stratos.SdsAddress(resNodeP2PPubkeyTmp.Address()))
	}

	/********************* initialize mock app *********************/
	accs, balances := setupAccountsMultipleResNodes(resOwners)

	// create validator set with single validator
	consPubKey, err := cryptocodec.ToTmPubKeyInterface(valConsPubk1)
	validator := tmtypes.NewValidator(consPubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	metaNodes := setupAllMetaNodes()
	//resourceNodes := setupAllResourceNodes()
	resourceNodes := make([]registertypes.ResourceNode, 0)

	stApp := stratostestutil.SetupWithGenesisNodeSet(t, valSet, metaNodes, resourceNodes, accs, chainID, false, balances...)

	accountKeeper := stApp.GetAccountKeeper()
	//bankKeeper := stApp.GetBankKeeper()
	registerKeeper := stApp.GetRegisterKeeper()
	potKeeper := stApp.GetPotKeeper()

	/********************* foundation account deposit *********************/
	header := tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := stApp.BaseApp.NewContext(true, header)

	foundationDepositMsg := pottypes.NewMsgFoundationDeposit(foundationDeposit, foundationDepositorAccAddr)
	txGen := stratostestutil.MakeTestEncodingConfig().TxConfig

	senderAcc := accountKeeper.GetAccount(ctx, foundationDepositorAccAddr)
	accNum := senderAcc.GetAccountNumber()
	accSeq := senderAcc.GetSequence()

	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, foundationDepositorPrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(pottypes.FoundationAccount)
	stratostestutil.CheckBalance(t, stApp, foundationAccountAddr, foundationDeposit)

	/********************** commit **********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	_, nozSupply := potKeeper.NozSupply(ctx)
	//nozPriceFactorsSeq0, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, registerKeeper, NozPriceFactors{
	//	NOzonePrice:          registerKeeper.GetCurrNozPriceParams(ctx),
	//	InitialTotalDeposit:   registerKeeper.GetInitialGenesisDepositTotal(ctx),
	//	EffectiveTotalDeposit: registerKeeper.GetEffectiveTotalDeposit(ctx),
	//	TotalUnissuedPrepay:  registerKeeper.GetTotalUnissuedPrepay(ctx).Amount,
	//	DepositAndPrepay:       registerKeeper.GetInitialGenesisDepositTotal(ctx).Add(registerKeeper.GetTotalUnissuedPrepay(ctx).Amount),
	//	OzoneLimit:           registerKeeper.GetRemainingOzoneLimit(ctx),
	//	NozSupply:            nozSupply,
	//})

	// start testing
	t.Log("********************************* Deliver Create and unsuspend ResourceNode Tx START ********************************************")

	prepayMsg := setupPrepayMsgWithResOwner(resOwners[0])
	/********************* deliver tx *********************/

	senderAcc = accountKeeper.GetAccount(ctx, resOwners[0])
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()
	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{prepayMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, resOwnerPrivKeys[0])
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	St, Pt, Lt := registerKeeper.GetCurrNozPriceParams(ctx)
	nozPriceFactorsSeq0, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, NozPriceFactors{
		NOzonePrice:           potKeeper.GetCurrentNozPrice(St, Pt, Lt),
		InitialTotalDeposit:   registerKeeper.GetInitialGenesisDepositTotal(ctx),
		EffectiveTotalDeposit: registerKeeper.GetEffectiveTotalDeposit(ctx),
		TotalUnissuedPrepay:   registerKeeper.GetTotalUnissuedPrepay(ctx).Amount,
		DepositAndPrepay:      registerKeeper.GetInitialGenesisDepositTotal(ctx).Add(registerKeeper.GetTotalUnissuedPrepay(ctx).Amount),
		OzoneLimit:            registerKeeper.GetRemainingOzoneLimit(ctx),
		NozSupply:             nozSupply,
	})
	priceBefore := nozPriceFactorsSeq0
	priceAfter := nozPriceFactorsSeq0
	dataToExcel = append(dataToExcel, priceBefore)

	//resOwnerAcc := accountKeeper.GetAccount(ctx, resOwner1)
	//ownerAccNum := resOwnerAcc.GetAccountNumber()
	//ownerAccSeq := resOwnerAcc.GetSequence()

	for i := 0; i < NUM_OF_SAMPLE; i++ {
		createResourceNodeMsg := setupMsgCreateResourceNode(i, resNodeP2PAddresses[i], resNodeP2PPubkeys[i], resOwners[i])
		/********************* deliver tx *********************/

		senderAcc = accountKeeper.GetAccount(ctx, resOwners[i])
		accNum = senderAcc.GetAccountNumber()
		accSeq = senderAcc.GetSequence()
		_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, resOwnerPrivKeys[i])
		require.NoError(t, err)
		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)

		unsuspendMsg := setupUnsuspendMsgByIndex(i, resNodeP2PAddresses[i], resOwners[i])
		/********************* deliver tx *********************/

		senderAcc = accountKeeper.GetAccount(ctx, metaOwner1)
		accNum = senderAcc.GetAccountNumber()
		accSeq = senderAcc.GetSequence()

		_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{unsuspendMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, metaOwnerPrivKey1)
		require.NoError(t, err)
		/********************* new height & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)

		t.Log("********************************* Deliver UpdateEffectiveDeposit Tx START ********************************************")
		updateEffectiveDepositMsg := setupMsgUpdateEffectiveDeposit(resNodeP2PAddresses[i])
		/********************* deliver tx *********************/

		senderAcc = accountKeeper.GetAccount(ctx, metaOwner1)
		accNum = senderAcc.GetAccountNumber()
		accSeq = senderAcc.GetSequence()

		_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{updateEffectiveDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, metaOwnerPrivKey1)
		require.NoError(t, err)
		t.Log("********************************* Deliver UpdateEffectiveDeposit Tx END ********************************************\n\n...\n[NEXT TEST CASE]")
		/********************* new height & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)

		priceAfter, nozPricePercentage, ozoneLimitPercentage = printCurrNozPrice(t, ctx, potKeeper, registerKeeper, priceBefore)
		require.True(t, nozPricePercentage.LT(sdkmath.LegacyZeroDec()), "noz price should decrease after CreateResourceNode")
		require.True(t, ozoneLimitPercentage.GT(sdkmath.LegacyZeroDec()), "OzLimit should increase after CreateResourceNode")
		t.Log("********************************* Deliver Create and unsuspend ResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

		dataToExcel = append(dataToExcel, priceAfter)
		priceBefore = priceAfter
	}

	for i := 0; i < NUM_OF_SAMPLE; i++ {
		slashingMsg := setupSuspendMsgByIndex(i, resNodeP2PAddresses[i], resNodeP2PPubkeys[i], resOwners[i])
		/********************* deliver tx *********************/

		senderAcc = accountKeeper.GetAccount(ctx, metaOwner1)
		accNum = senderAcc.GetAccountNumber()
		accSeq = senderAcc.GetSequence()

		_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{slashingMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, metaOwnerPrivKey1)
		require.NoError(t, err)
		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)

		//createResourceNodeMsg := setupMsgRemoveResourceNode(i, resNodeP2PAddresses[i], resOwners[i])
		///********************* deliver tx *********************/
		//
		//resOwnerAcc := accountKeeper.GetAccount(ctx, resOwners[i])
		//ownerAccNum := resOwnerAcc.GetAccountNumber()
		//ownerAccSeq := resOwnerAcc.GetSequence()
		//
		//_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, resOwnerPrivKeys[i])
		//require.NoError(t, err)
		///********************* commit & check result *********************/
		//header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		//stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		//ctx = stApp.BaseApp.NewContext(true, header)

		priceAfter, nozPricePercentage, ozoneLimitPercentage = printCurrNozPrice(t, ctx, potKeeper, registerKeeper, priceBefore)
		require.True(t, nozPricePercentage.GT(sdkmath.LegacyZeroDec()), "noz price should increase after RemoveResourceNode")
		require.True(t, ozoneLimitPercentage.LT(sdkmath.LegacyZeroDec()), "OzLimit should decrease after RemoveResourceNode")
		t.Log("********************************* Deliver CreateResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

		dataToExcel = append(dataToExcel, priceAfter)
		priceBefore = priceAfter
	}

	exportToCSV(t, dataToExcel)

}

func TestOzPriceChangeRemoveMultipleResourceNodeAfterGenesis(t *testing.T) {
	NUM_OF_SAMPLE := 100
	dataToExcel := make([]NozPriceFactors, 0, NUM_OF_SAMPLE)

	resOwnerPrivKeys := make([]*secp256k1.PrivKey, 0, NUM_OF_SAMPLE)
	resOwners := make([]sdk.AccAddress, 0, NUM_OF_SAMPLE)

	resNodeP2PPubKeys := make([]cryptotypes.PubKey, 0, NUM_OF_SAMPLE)
	resNodeP2PAddresses := make([]stratos.SdsAddress, 0, NUM_OF_SAMPLE)

	for i := 0; i < NUM_OF_SAMPLE; i++ {
		resOwnerPrivKeyTmp := secp256k1.GenPrivKey()
		resOwnerPrivKeys = append(resOwnerPrivKeys, resOwnerPrivKeyTmp)
		resOwnerAddrTmp := sdk.AccAddress(resOwnerPrivKeyTmp.PubKey().Address())
		resOwners = append(resOwners, resOwnerAddrTmp)

		resNodeP2PPubKeyTmp := ed25519.GenPrivKey().PubKey()
		resNodeP2PPubKeys = append(resNodeP2PPubKeys, resNodeP2PPubKeyTmp)
		resNodeP2PAddresses = append(resNodeP2PAddresses, stratos.SdsAddress(resNodeP2PPubKeyTmp.Address()))
	}

	/********************* initialize mock app *********************/
	accs, balances := setupAccountsMultipleResNodes(resOwners)

	// create validator set with single validator
	consPubKey, err := cryptocodec.ToTmPubKeyInterface(valConsPubk1)
	validator := tmtypes.NewValidator(consPubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	metaNodes := setupAllMetaNodes()
	resourceNodes := setupMultipleResourceNodes(resNodeP2PAddresses, resNodeP2PPubKeys, resOwners)

	stApp := stratostestutil.SetupWithGenesisNodeSet(t, valSet, metaNodes, resourceNodes, accs, chainID, false, balances...)

	accountKeeper := stApp.GetAccountKeeper()
	//bankKeeper := stApp.GetBankKeeper()
	registerKeeper := stApp.GetRegisterKeeper()
	potKeeper := stApp.GetPotKeeper()

	/********************* foundation account deposit *********************/
	header := tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := stApp.BaseApp.NewContext(true, header)

	foundationDepositMsg := pottypes.NewMsgFoundationDeposit(foundationDeposit, foundationDepositorAccAddr)
	txGen := stratostestutil.MakeTestEncodingConfig().TxConfig

	senderAcc := accountKeeper.GetAccount(ctx, foundationDepositorAccAddr)
	accNum := senderAcc.GetAccountNumber()
	accSeq := senderAcc.GetSequence()

	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, foundationDepositorPrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(pottypes.FoundationAccount)
	stratostestutil.CheckBalance(t, stApp, foundationAccountAddr, foundationDeposit)

	/********************** commit **********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	_, nozSupply := potKeeper.NozSupply(ctx)

	//for i := 0; i < NUM_OF_SAMPLE; i++ {
	//	unsuspendMsg := setupUnsuspendMsgByIndex(i, resNodeP2PAddresses[i], resOwners[i])
	//	/********************* deliver tx *********************/
	//
	//	senderAcc = accountKeeper.GetAccount(ctx, metaOwner1)
	//	accNum = senderAcc.GetAccountNumber()
	//	accSeq = senderAcc.GetSequence()
	//
	//	fmt.Println("!!!!!!!!!!!!!!!!! ------------- i = ", i)
	//
	//	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{unsuspendMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, metaOwnerPrivKey1)
	//	require.NoError(t, err)
	//	/********************* commit & check result *********************/
	//	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	//	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	//	ctx = stApp.BaseApp.NewContext(true, header)
	//
	//}

	// start testing
	t.Log("********************************* Deliver RemoveResourceNode Tx START ********************************************")

	prepayMsg := setupPrepayMsgWithResOwner(resOwners[0])
	/********************* deliver tx *********************/

	senderAcc = accountKeeper.GetAccount(ctx, resOwners[0])
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()

	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{prepayMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, resOwnerPrivKeys[0])
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)
	St, Pt, Lt := registerKeeper.GetCurrNozPriceParams(ctx)
	nozPriceFactorsSeq0, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, NozPriceFactors{
		NOzonePrice:           potKeeper.GetCurrentNozPrice(St, Pt, Lt),
		InitialTotalDeposit:   registerKeeper.GetInitialGenesisDepositTotal(ctx),
		EffectiveTotalDeposit: registerKeeper.GetEffectiveTotalDeposit(ctx),
		TotalUnissuedPrepay:   registerKeeper.GetTotalUnissuedPrepay(ctx).Amount,
		DepositAndPrepay:      registerKeeper.GetInitialGenesisDepositTotal(ctx).Add(registerKeeper.GetTotalUnissuedPrepay(ctx).Amount),
		OzoneLimit:            registerKeeper.GetRemainingOzoneLimit(ctx),
		NozSupply:             nozSupply,
	})

	priceBefore := nozPriceFactorsSeq0
	priceAfter := nozPriceFactorsSeq0
	dataToExcel = append(dataToExcel, priceBefore)

	for i := 0; i < NUM_OF_SAMPLE; i++ {
		slashingMsg := setupSuspendMsgByIndex(i, resNodeP2PAddresses[i], resNodeP2PPubKeys[i], resOwners[i])
		/********************* deliver tx *********************/

		senderAcc = accountKeeper.GetAccount(ctx, metaOwner1)
		accNum = senderAcc.GetAccountNumber()
		accSeq = senderAcc.GetSequence()
		_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{slashingMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, metaOwnerPrivKey1)
		require.NoError(t, err)
		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)
		//createResourceNodeMsg := setupMsgRemoveResourceNode(i, resNodeP2PAddresses[i], resOwners[i])
		///********************* deliver tx *********************/
		//
		//resOwnerAcc := accountKeeper.GetAccount(ctx, resOwners[i])
		//ownerAccNum := resOwnerAcc.GetAccountNumber()
		//ownerAccSeq := resOwnerAcc.GetSequence()
		//
		//_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, resOwnerPrivKeys[i])
		//require.NoError(t, err)
		///********************* commit & check result *********************/
		//header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		//stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		//ctx = stApp.BaseApp.NewContext(true, header)

		priceAfter, nozPricePercentage, ozoneLimitPercentage = printCurrNozPrice(t, ctx, potKeeper, registerKeeper, priceBefore)
		require.True(t, nozPricePercentage.GT(sdkmath.LegacyZeroDec()), "noz price should increase after RemoveResourceNode")
		require.True(t, ozoneLimitPercentage.LT(sdkmath.LegacyZeroDec()), "OzLimit should decrease after RemoveResourceNode")
		t.Log("********************************* Deliver CreateResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

		dataToExcel = append(dataToExcel, priceAfter)
		priceBefore = priceAfter
	}

	exportToCSV(t, dataToExcel)

}

func exportToCSV(t *testing.T, factors []NozPriceFactors) {
	t.Logf("\n%v, %v, %v, %v, %v, %v, %v", "Index", "InitialTotalDeposit", "TotalUnissuedPrepay", "DepositAndPrepay",
		"NOzonePrice", "RemainingOzoneLimit", "TotalNozSupply")
	for i, factor := range factors {
		t.Logf("\n%v, %v, %v, %v, %v, %v, %v", i+1, factor.InitialTotalDeposit.String(), factor.TotalUnissuedPrepay.String(), factor.DepositAndPrepay.String(),
			factor.NOzonePrice.String(), factor.OzoneLimit.String(), factor.NozSupply.String())
	}
	t.Log("\n")
}
