package sds_test

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

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
	chainID           = "testchain_1-1"
	stos2wei          = stratos.StosToWei
	StosToWeiSuffix   = "000000000000000000" // 1 Stos = 1e18 wei
	rewardDenom       = stratos.Utros
	depositNozRateStr = "100000"
)

var (
	depositNozRateInt, _ = sdk.NewIntFromString(depositNozRateStr)

	paramSpecificMinedReward = sdk.NewCoins(stratos.NewCoinInt64(160000000000))
	paramSpecificEpoch       = sdk.NewInt(10)

	resNodeSlashingNOZAmt1            = sdk.NewInt(100000000000)
	resNodeSlashingEffectiveTokenAmt1 = sdk.NewInt(1000000000000000000)

	resourceNodeVolume1 = sdk.NewInt(537500000000)
	resourceNodeVolume2 = sdk.NewInt(200000000000)
	resourceNodeVolume3 = sdk.NewInt(200000000000)

	depositForSendingTx, _ = sdk.NewIntFromString("100000000000000000000000000000")
	totalUnissuedPrepayVal = sdk.ZeroInt()
	totalUnissuedPrepay    = stratos.NewCoin(totalUnissuedPrepayVal)

	foundationDepositorPrivKey = secp256k1.GenPrivKey()
	foundationDepositorAccAddr = sdk.AccAddress(foundationDepositorPrivKey.PubKey().Address())
	foundationDeposit          = sdk.NewCoins(sdk.NewCoin(rewardDenom, sdk.NewInt(40000000000000000)))

	resOwnerPrivKey1 = secp256k1.GenPrivKey()
	resOwnerPrivKey2 = secp256k1.GenPrivKey()
	resOwnerPrivKey3 = secp256k1.GenPrivKey()
	resOwnerPrivKey4 = secp256k1.GenPrivKey()
	resOwnerPrivKey5 = secp256k1.GenPrivKey()
	idxOwnerPrivKey1 = secp256k1.GenPrivKey()
	idxOwnerPrivKey2 = secp256k1.GenPrivKey()
	idxOwnerPrivKey3 = secp256k1.GenPrivKey()

	resOwner1 = sdk.AccAddress(resOwnerPrivKey1.PubKey().Address())
	resOwner2 = sdk.AccAddress(resOwnerPrivKey2.PubKey().Address())
	resOwner3 = sdk.AccAddress(resOwnerPrivKey3.PubKey().Address())
	resOwner4 = sdk.AccAddress(resOwnerPrivKey4.PubKey().Address())
	resOwner5 = sdk.AccAddress(resOwnerPrivKey5.PubKey().Address())
	idxOwner1 = sdk.AccAddress(idxOwnerPrivKey1.PubKey().Address())
	idxOwner2 = sdk.AccAddress(idxOwnerPrivKey2.PubKey().Address())
	idxOwner3 = sdk.AccAddress(idxOwnerPrivKey3.PubKey().Address())

	resNodeInitialDepositForMultipleNodes, _ = sdk.NewIntFromString("3" + StosToWeiSuffix)
	//resNodeInitialDepositForMultipleNodes = sdk.NewInt(3 * stos2wei)

	resNodePubKey1            = ed25519.GenPrivKey().PubKey()
	resNodeAddr1              = sdk.AccAddress(resNodePubKey1.Address())
	resNodeNetworkId1         = stratos.SdsAddress(resNodePubKey1.Address())
	resNodeInitialDeposit1, _ = sdk.NewIntFromString("3" + StosToWeiSuffix)
	//resNodeInitialDeposit1 = sdk.NewInt(3 * stos2wei)

	resNodePubKey2            = ed25519.GenPrivKey().PubKey()
	resNodeAddr2              = sdk.AccAddress(resNodePubKey2.Address())
	resNodeNetworkId2         = stratos.SdsAddress(resNodePubKey2.Address())
	resNodeInitialDeposit2, _ = sdk.NewIntFromString("3" + StosToWeiSuffix)
	//resNodeInitialDeposit2 = sdk.NewInt(3 * stos2wei)

	resNodePubKey3            = ed25519.GenPrivKey().PubKey()
	resNodeAddr3              = sdk.AccAddress(resNodePubKey3.Address())
	resNodeNetworkId3         = stratos.SdsAddress(resNodePubKey3.Address())
	resNodeInitialDeposit3, _ = sdk.NewIntFromString("3" + StosToWeiSuffix)
	//resNodeInitialDeposit3 = sdk.NewInt(3 * stos2wei)

	resNodePubKey4         = ed25519.GenPrivKey().PubKey()
	resNodeAddr4           = sdk.AccAddress(resNodePubKey4.Address())
	resNodeNetworkId4      = stratos.SdsAddress(resNodePubKey4.Address())
	resNodeInitialDeposit4 = sdk.NewInt(3 * stos2wei)

	resNodePubKey5         = ed25519.GenPrivKey().PubKey()
	resNodeAddr5           = sdk.AccAddress(resNodePubKey5.Address())
	resNodeNetworkId5      = stratos.SdsAddress(resNodePubKey5.Address())
	resNodeInitialDeposit5 = sdk.NewInt(3 * stos2wei)

	idxNodePrivKey1           = ed25519.GenPrivKey()
	idxNodePubKey1            = idxNodePrivKey1.PubKey()
	idxNodeAddr1              = sdk.AccAddress(idxNodePubKey1.Address())
	idxNodeNetworkId1         = stratos.SdsAddress(idxNodePubKey1.Address())
	idxNodeInitialDeposit1, _ = sdk.NewIntFromString("5" + StosToWeiSuffix)
	//idxNodeInitialDeposit1 = sdk.NewInt(5 * stos2wei)

	idxNodePubKey2            = ed25519.GenPrivKey().PubKey()
	idxNodeAddr2              = sdk.AccAddress(idxNodePubKey2.Address())
	idxNodeNetworkId2         = stratos.SdsAddress(idxNodePubKey2.Address())
	idxNodeInitialDeposit2, _ = sdk.NewIntFromString("5" + StosToWeiSuffix)
	//idxNodeInitialDeposit2 = sdk.NewInt(5 * stos2wei)

	idxNodePubKey3            = ed25519.GenPrivKey().PubKey()
	idxNodeAddr3              = sdk.AccAddress(idxNodePubKey3.Address())
	idxNodeNetworkId3         = stratos.SdsAddress(idxNodePubKey3.Address())
	idxNodeInitialDeposit3, _ = sdk.NewIntFromString("5" + StosToWeiSuffix)
	//idxNodeInitialDeposit3 = sdk.NewInt(5 * stos2wei)

	valOpPrivKey1 = secp256k1.GenPrivKey()
	valOpPubKey1  = valOpPrivKey1.PubKey()
	valOpValAddr1 = sdk.ValAddress(valOpPubKey1.Address())
	valOpAccAddr1 = sdk.AccAddress(valOpPubKey1.Address())

	valConsPrivKey1    = ed25519.GenPrivKey()
	valConsPubk1       = valConsPrivKey1.PubKey()
	valInitialStake, _ = sdk.NewIntFromString("15" + StosToWeiSuffix)
)

type NozPriceFactors struct {
	NOzonePrice           sdk.Dec
	InitialTotalDeposit   sdk.Int
	EffectiveTotalDeposit sdk.Int
	TotalUnissuedPrepay   sdk.Int
	DepositAndPrepay      sdk.Int
	OzoneLimit            sdk.Int
	NozSupply             sdk.Int
}

func TestPriceCurve(t *testing.T) {

	NUM_TESTS := 100

	initFactorsBefore := &NozPriceFactors{
		NOzonePrice:           nozPrice,
		InitialTotalDeposit:   initialTotalDepositStore,
		EffectiveTotalDeposit: initialTotalDepositStore,
		TotalUnissuedPrepay:   totalUnissuedPrepayStore,
		DepositAndPrepay:      initialTotalDepositStore.Add(totalUnissuedPrepayStore),
		OzoneLimit:            initialTotalDepositStore.ToDec().Quo(nozPrice).TruncateInt(),
		NozSupply:             initialTotalDepositStore.ToDec().Quo(depositNozRateInt.ToDec()).TruncateInt(),
	}

	initFactorsBefore, _, _ = simulatePriceChange(t, &PriceChangeEvent{
		depositDelta:        sdk.ZeroInt(),
		unissuedPrepayDelta: sdk.ZeroInt(),
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
		depositDeltaChange, _ := sdk.NewIntFromString(strconv.Itoa(depositChangePerm[i]) + StosToWeiSuffix)
		unissuedPrepayDeltaChange, _ := sdk.NewIntFromString(strconv.Itoa(prepayChangePerm[i]) + StosToWeiSuffix)
		change := &PriceChangeEvent{
			depositDelta:        depositDeltaChange.Mul(sdk.NewInt(int64(tempDepositSign))),
			unissuedPrepayDelta: unissuedPrepayDeltaChange.Mul(sdk.NewInt(int64(tempPrepaySign))),
		}
		t.Logf("\ndepositDeltaOri: %d, unissuedPrepayDeltaOri: %d\n", depositChangePerm[i], prepayChangePerm[i])
		t.Logf("\ndepositDelta: %v, unissuedPrepayDelta: %v\n", change.depositDelta.String(), change.unissuedPrepayDelta.String())
		initFactorsBefore, _, _ = simulatePriceChange(t, change, initFactorsBefore)
	}
}

func TestOzPriceChange(t *testing.T) {
	/********************* initialize mock app *********************/
	//mApp, k, stakingKeeper, bankKeeper, supplyKeeper, registerKeeper := getMockApp(t)
	accs, balances := setupAccounts()
	//stApp := app.SetupWithGenesisAccounts(accs, chainID, balances...)
	validators := make([]*tmtypes.Validator, 0)
	valSet := tmtypes.NewValidatorSet(validators)
	metaNodes := setupAllMetaNodes()
	//resourceNodes := setupAllResourceNodes()
	resourceNodes := make([]registertypes.ResourceNode, 0)

	stApp := app.SetupWithGenesisNodeSet(t, true, valSet, metaNodes, resourceNodes, accs, chainID, balances...)

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
	txGen := app.MakeTestEncodingConfig().TxConfig

	foundationDepositorAcc := accountKeeper.GetAccount(ctx, foundationDepositorAccAddr)
	accNum := foundationDepositorAcc.GetAccountNumber()
	accSeq := foundationDepositorAcc.GetSequence()
	_, _, err := app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, foundationDepositorPrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(pottypes.FoundationAccount)
	app.CheckBalance(t, stApp, foundationAccountAddr, foundationDeposit)

	/********************* create validator with 50% commission *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	commission := stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	description := stakingtypes.NewDescription("foo_moniker", chainID, "", "", "")
	createValidatorMsg, err := stakingtypes.NewMsgCreateValidator(valOpValAddr1, valConsPubk1, stratos.NewCoin(valInitialStake), description, commission, sdk.OneInt())

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

	resOwnerAcc := accountKeeper.GetAccount(ctx, resOwner1)
	ownerAccNum := resOwnerAcc.GetAccountNumber()
	ownerAccSeq := resOwnerAcc.GetSequence()

	_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{prepayMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, resOwnerPrivKey1)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	nozPriceFactorsSeq1, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq0)
	require.True(t, nozPricePercentage.GT(sdk.ZeroDec()), "noz price should increase after PREPAY")
	require.True(t, ozoneLimitPercentage.LT(sdk.ZeroDec()), "OzLimit should decrease after PREPAY")
	t.Log("********************************* Deliver Prepay Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver CreateResourceNode Tx START ********************************************")
	createResourceNodeMsg := setupMsgCreateResourceNode1()
	/********************* deliver tx *********************/

	resOwnerAcc = accountKeeper.GetAccount(ctx, resOwner1)
	ownerAccNum = resOwnerAcc.GetAccountNumber()
	ownerAccSeq = resOwnerAcc.GetSequence()

	_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, resOwnerPrivKey1)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	nozPriceFactorsSeq2, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq1)
	require.True(t, nozPricePercentage.Equal(sdk.ZeroDec()), "noz price shouldn't change after CreateResourceNode")
	require.True(t, ozoneLimitPercentage.Equal(sdk.ZeroDec()), "OzLimit shouldn't change  after CreateResourceNode")
	t.Log("********************************* Deliver CreateResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver UnsuspendResourceNode Tx (Slashing) START ********************************************")
	unsuspendMsg := setupUnsuspendMsg()
	/********************* deliver tx *********************/

	idxOwnerAcc := accountKeeper.GetAccount(ctx, idxOwner1)
	ownerAccNum = idxOwnerAcc.GetAccountNumber()
	ownerAccSeq = idxOwnerAcc.GetSequence()

	_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{unsuspendMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, idxOwnerPrivKey1)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	slashingAmtSetup := registerKeeper.GetSlashing(ctx, resOwner1)

	totalConsumedNoz := sdk.ZeroInt().ToDec()

	slashingAmtCheck := potKeeper.GetTrafficReward(ctx, totalConsumedNoz)
	t.Log("slashingAmtSetup=" + slashingAmtSetup.String())
	require.Equal(t, slashingAmtSetup, slashingAmtCheck.TruncateInt())

	nozPriceFactorsSeq3, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq2)
	require.True(t, nozPricePercentage.LT(sdk.ZeroDec()), "noz price should decrease after UnsuspendResourceNode")
	require.True(t, ozoneLimitPercentage.GT(sdk.ZeroDec()), "OzLimit should increase after UnsuspendResourceNode")
	t.Log("********************************* Deliver UnsuspendResourceNode Tx (Slashing) END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver SuspendResourceNode Tx (Slashing) START ********************************************")
	slashingMsg := setupSlashingMsg()
	/********************* deliver tx *********************/

	idxOwnerAcc = accountKeeper.GetAccount(ctx, idxOwner1)
	ownerAccNum = idxOwnerAcc.GetAccountNumber()
	ownerAccSeq = idxOwnerAcc.GetSequence()

	_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{slashingMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, idxOwnerPrivKey1)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	//slashingAmtSetup = registerKeeper.GetSlashing(ctx, resOwner1)
	//
	//totalConsumedNoz = resNodeSlashingNOZAmt1.ToDec()
	//
	//slashingAmtCheck = potKeeper.GetTrafficReward(ctx, totalConsumedNoz)
	//t.Log("slashingAmtSetup=" + slashingAmtSetup.String())
	//require.Equal(t, slashingAmtSetup, slashingAmtCheck.TruncateInt())

	nozPriceFactorsSeq4, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq3)
	require.True(t, nozPricePercentage.GT(sdk.ZeroDec()), "noz price should increase after SlashResourceNode")
	require.True(t, ozoneLimitPercentage.LT(sdk.ZeroDec()), "OzLimit should decrease after SlashResourceNode")

	_, nozPricePercentage42, ozoneLimitPercentage42 := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq2)
	require.True(t, nozPricePercentage42.Equal(sdk.ZeroDec()), "noz price after SlashResourceNode should be same with the price when node hadn't been activated")
	require.True(t, ozoneLimitPercentage42.Equal(sdk.ZeroDec()), "OzLimit after SlashResourceNode should be same with the ozLimit when node hadn't been activated")
	t.Log("********************************* Deliver SuspendResourceNode Tx (Slashing) END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver VolumeReport Tx START ********************************************")
	/********************* prepare tx data *********************/
	volumeReportMsg := setupMsgVolumeReport(1)

	lastTotalMinedToken := potKeeper.GetTotalMinedTokens(ctx)
	t.Log("last committed TotalMinedTokens = " + lastTotalMinedToken.String())
	epoch, ok := sdk.NewIntFromString(volumeReportMsg.Epoch.String())
	require.Equal(t, ok, true)

	totalConsumedNoz = potKeeper.GetTotalConsumedNoz(volumeReportMsg.WalletVolumes).ToDec()
	remaining, total := potKeeper.NozSupply(ctx)
	require.True(t, potKeeper.GetTotalConsumedNoz(volumeReportMsg.WalletVolumes).Add(remaining).LTE(total), "remaining+consumed Noz exceeds total Noz supply")

	/********************* print info *********************/
	t.Log("epoch " + volumeReportMsg.Epoch.String())
	StDec := registerKeeper.GetEffectiveTotalDeposit(ctx).ToDec()
	PtDec := registerKeeper.GetTotalUnissuedPrepay(ctx).Amount.ToDec()
	Y := totalConsumedNoz
	LtDec := registerKeeper.GetRemainingOzoneLimit(ctx).ToDec()
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

	t.Log("indexing_wallet1:  address = " + idxOwner1.String())
	t.Log("              miningReward = " + rewardDetailMap[idxOwner1.String()].RewardFromMiningPool.String())
	t.Log("             trafficReward = " + rewardDetailMap[idxOwner1.String()].RewardFromTrafficPool.String())

	t.Log("indexing_wallet2:  address = " + idxOwner2.String())
	t.Log("              miningReward = " + rewardDetailMap[idxOwner2.String()].RewardFromMiningPool.String())
	t.Log("             trafficReward = " + rewardDetailMap[idxOwner2.String()].RewardFromTrafficPool.String())

	t.Log("indexing_wallet3:  address = " + idxOwner3.String())
	t.Log("              miningReward = " + rewardDetailMap[idxOwner3.String()].RewardFromMiningPool.String())
	t.Log("             trafficReward = " + rewardDetailMap[idxOwner3.String()].RewardFromTrafficPool.String())
	t.Log("---------------------------")

	/********************* record data before delivering tx  *********************/
	lastFoundationAccBalance := bankKeeper.GetAllBalances(ctx, foundationAccountAddr)
	lastUnissuedPrepay := registerKeeper.GetTotalUnissuedPrepay(ctx)
	lastCommunityPool := sdk.NewCoins(sdk.NewCoin(potKeeper.BondDenom(ctx), distrKeeper.GetFeePool(ctx).CommunityPool.AmountOf(potKeeper.BondDenom(ctx)).TruncateInt()))
	lastMatureTotalOfResNode1 := potKeeper.GetMatureTotalReward(ctx, resOwner1)

	/********************* deliver tx *********************/
	idxOwnerAcc = accountKeeper.GetAccount(ctx, idxOwner1)
	ownerAccNum = idxOwnerAcc.GetAccountNumber()
	ownerAccSeq = idxOwnerAcc.GetSequence()

	feePoolAccAddr := accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	require.NotNil(t, feePoolAccAddr)
	feeCollectorToFeePoolAtBeginBlock := bankKeeper.GetBalance(ctx, feePoolAccAddr, potKeeper.BondDenom(ctx))

	_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{volumeReportMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, idxOwnerPrivKey1)
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	stApp.EndBlock(abci.RequestEndBlock{Height: header.Height})
	stApp.Commit()
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
		lastCommunityPool,
		lastMatureTotalOfResNode1,
		slashingAmtSetup,
		feeCollectorToFeePoolAtBeginBlock,
	)

	nozPriceFactorsSeq5, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq4)
	require.True(t, nozPricePercentage.LT(sdk.ZeroDec()), "noz price should decrease after VolumeReport")
	require.True(t, ozoneLimitPercentage.GT(sdk.ZeroDec()), "OzLimit shouldn't change after VolumeReport")
	t.Log("********************************* Deliver VolumeReport Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver CreateResourceNode Tx START ********************************************")
	createResourceNodeMsg2 := setupMsgCreateResourceNode2()
	/********************* deliver tx *********************/

	resOwnerAcc = accountKeeper.GetAccount(ctx, resOwner2)
	ownerAccNum = resOwnerAcc.GetAccountNumber()
	ownerAccSeq = resOwnerAcc.GetSequence()

	_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg2}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, resOwnerPrivKey2)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	nozPriceFactorsSeq6, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq5)
	require.True(t, nozPricePercentage.Equal(sdk.ZeroDec()), "noz price shouldn't change after CreateResourceNode")
	require.True(t, ozoneLimitPercentage.Equal(sdk.ZeroDec()), "OzLimit shouldn't change  after CreateResourceNode")
	t.Log("********************************* Deliver CreateResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver CreateResourceNode Tx START ********************************************")
	createResourceNodeMsg3 := setupMsgCreateResourceNode3()
	/********************* deliver tx *********************/

	resOwnerAcc = accountKeeper.GetAccount(ctx, resOwner3)
	ownerAccNum = resOwnerAcc.GetAccountNumber()
	ownerAccSeq = resOwnerAcc.GetSequence()

	_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg3}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, resOwnerPrivKey3)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	nozPriceFactorsSeq7, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq6)
	require.True(t, nozPricePercentage.Equal(sdk.ZeroDec()), "noz price shouldn't change after CreateResourceNode")
	require.True(t, ozoneLimitPercentage.Equal(sdk.ZeroDec()), "OzLimit shouldn't change  after CreateResourceNode")
	t.Log("********************************* Deliver CreateResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver CreateResourceNode Tx START ********************************************")
	createResourceNodeMsg4 := setupMsgCreateResourceNode4()
	/********************* deliver tx *********************/

	resOwnerAcc = accountKeeper.GetAccount(ctx, resOwner4)
	ownerAccNum = resOwnerAcc.GetAccountNumber()
	ownerAccSeq = resOwnerAcc.GetSequence()

	_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg4}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, resOwnerPrivKey4)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	nozPriceFactorsSeq8, nozPricePercentage, ozoneLimitPercentage := printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq7)
	require.True(t, nozPricePercentage.Equal(sdk.ZeroDec()), "noz price shouldn't change after CreateResourceNode")
	require.True(t, ozoneLimitPercentage.Equal(sdk.ZeroDec()), "OzLimit shouldn't change  after CreateResourceNode")
	t.Log("********************************* Deliver CreateResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

	t.Log("********************************* Deliver CreateResourceNode Tx START ********************************************")
	createResourceNodeMsg5 := setupMsgCreateResourceNode5()
	/********************* deliver tx *********************/

	resOwnerAcc = accountKeeper.GetAccount(ctx, resOwner5)
	ownerAccNum = resOwnerAcc.GetAccountNumber()
	ownerAccSeq = resOwnerAcc.GetSequence()

	_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg5}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, resOwnerPrivKey5)
	require.NoError(t, err)
	/********************* commit & check result *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	_, nozPricePercentage, ozoneLimitPercentage = printCurrNozPrice(t, ctx, potKeeper, registerKeeper, nozPriceFactorsSeq8)
	require.True(t, nozPricePercentage.Equal(sdk.ZeroDec()), "noz price shouldn't change after CreateResourceNode")
	require.True(t, ozoneLimitPercentage.Equal(sdk.ZeroDec()), "OzLimit shouldn't change  after CreateResourceNode")
	t.Log("********************************* Deliver CreateResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

}

// initialize data of volume report
func setupMsgVolumeReport(newEpoch int64) *pottypes.MsgVolumeReport {
	volume1 := pottypes.NewSingleWalletVolume(resOwner1, resourceNodeVolume1)
	volume2 := pottypes.NewSingleWalletVolume(resOwner2, resourceNodeVolume2)
	volume3 := pottypes.NewSingleWalletVolume(resOwner3, resourceNodeVolume3)

	nodesVolume := []pottypes.SingleWalletVolume{volume1, volume2, volume3}
	reporter := idxNodeNetworkId1
	epoch := sdk.NewInt(newEpoch)
	reportReference := "report for epoch " + epoch.String()
	reporterOwner := idxOwner1

	pubKeys := make([][]byte, 1)
	for i := range pubKeys {
		pubKeys[i] = make([]byte, 1)
	}

	signature := pottypes.NewBLSSignatureInfo(pubKeys, []byte("signature"), []byte("txData"))

	volumeReportMsg := pottypes.NewMsgVolumeReport(nodesVolume, reporter, epoch, reportReference, reporterOwner, signature, sdk.NewInt(0))

	return volumeReportMsg
}

func setupSlashingMsg() *pottypes.MsgSlashingResourceNode {
	reporters := make([]stratos.SdsAddress, 0)
	reporters = append(reporters, idxNodeNetworkId1)
	reportOwner := make([]sdk.AccAddress, 0)
	reportOwner = append(reportOwner, idxOwner1)
	slashingMsg := pottypes.NewMsgSlashingResourceNode(reporters, reportOwner, resNodeNetworkId1, resOwner1, resNodeSlashingNOZAmt1, true)
	return slashingMsg
}

func setupSuspendMsgByIndex(i int, resNodeNetworkId stratos.SdsAddress, resNodePubKey cryptotypes.PubKey, resOwner sdk.AccAddress) *pottypes.MsgSlashingResourceNode {
	reporters := make([]stratos.SdsAddress, 0)
	reporters = append(reporters, idxNodeNetworkId1)
	reportOwner := make([]sdk.AccAddress, 0)
	reportOwner = append(reportOwner, idxOwner1)
	slashingMsg := pottypes.NewMsgSlashingResourceNode(reporters, reportOwner, resNodeNetworkId, resOwner, resNodeSlashingNOZAmt1, true)
	return slashingMsg
}

func setupUnsuspendMsg() *pottypes.MsgSlashingResourceNode {
	reporters := make([]stratos.SdsAddress, 0)
	reporters = append(reporters, idxNodeNetworkId1)
	reportOwner := make([]sdk.AccAddress, 0)
	reportOwner = append(reportOwner, idxOwner1)
	slashingMsg := pottypes.NewMsgSlashingResourceNode(reporters, reportOwner, resNodeNetworkId1, resOwner1, sdk.ZeroInt(), false)
	return slashingMsg
}
func setupPrepayMsg() *sdstypes.MsgPrepay {
	sender := resOwner1
	amount, _ := sdk.NewIntFromString("1" + StosToWeiSuffix)
	coin := sdk.NewCoin(stratos.Wei, amount)
	prepayMsg := sdstypes.NewMsgPrepay(sender.String(), sender.String(), sdk.NewCoins(coin))
	return prepayMsg
}

func setupPrepayMsgWithResOwner(resOwner sdk.AccAddress) *sdstypes.MsgPrepay {
	sender := resOwner
	amount, _ := sdk.NewIntFromString("3" + StosToWeiSuffix)
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
	createResourceNodeMsg, _ := registertypes.NewMsgCreateResourceNode(resNodeNetworkId, resNodePubKey, sdk.NewCoin(stratos.Wei, resNodeInitialDepositForMultipleNodes), resOwner, registertypes.NewDescription("sds://resourceNode"+strconv.Itoa(i+1), "", "", "", ""), nodeType)
	return createResourceNodeMsg
}

func setupUnsuspendMsgByIndex(i int, resNodeNetworkId stratos.SdsAddress, resNodePubKey cryptotypes.PubKey, resOwner sdk.AccAddress) *pottypes.MsgSlashingResourceNode {
	reporters := make([]stratos.SdsAddress, 0)
	reporters = append(reporters, idxNodeNetworkId1)
	reportOwner := make([]sdk.AccAddress, 0)
	reportOwner = append(reportOwner, idxOwner1)
	slashingMsg := pottypes.NewMsgSlashingResourceNode(reporters, reportOwner, resNodeNetworkId, resOwner, sdk.ZeroInt(), false)
	return slashingMsg
}

func setupMsgCreateResourceNode1() *registertypes.MsgCreateResourceNode {
	nodeType := uint32(registertypes.STORAGE)
	createResourceNodeMsg, _ := registertypes.NewMsgCreateResourceNode(resNodeNetworkId1, resNodePubKey1, sdk.NewCoin(stratos.Wei, resNodeInitialDeposit1), resOwner1, registertypes.NewDescription("sds://resourceNode1", "", "", "", ""), nodeType)
	return createResourceNodeMsg
}
func setupMsgCreateResourceNode2() *registertypes.MsgCreateResourceNode {
	nodeType := uint32(registertypes.STORAGE)
	createResourceNodeMsg, _ := registertypes.NewMsgCreateResourceNode(resNodeNetworkId2, resNodePubKey2, sdk.NewCoin(stratos.Wei, resNodeInitialDeposit2), resOwner2, registertypes.NewDescription("sds://resourceNode2", "", "", "", ""), nodeType)
	return createResourceNodeMsg
}
func setupMsgCreateResourceNode3() *registertypes.MsgCreateResourceNode {
	nodeType := uint32(registertypes.STORAGE)
	createResourceNodeMsg, _ := registertypes.NewMsgCreateResourceNode(resNodeNetworkId3, resNodePubKey3, sdk.NewCoin(stratos.Wei, resNodeInitialDeposit3), resOwner3, registertypes.NewDescription("sds://resourceNode3", "", "", "", ""), nodeType)
	return createResourceNodeMsg
}

func setupMsgCreateResourceNode4() *registertypes.MsgCreateResourceNode {
	nodeType := uint32(registertypes.STORAGE)
	createResourceNodeMsg, _ := registertypes.NewMsgCreateResourceNode(resNodeNetworkId4, resNodePubKey4, sdk.NewCoin(stratos.Wei, resNodeInitialDeposit4), resOwner4, registertypes.NewDescription("sds://resourceNode4", "", "", "", ""), nodeType)
	return createResourceNodeMsg
}

func setupMsgCreateResourceNode5() *registertypes.MsgCreateResourceNode {
	nodeType := uint32(registertypes.STORAGE)
	createResourceNodeMsg, _ := registertypes.NewMsgCreateResourceNode(resNodeNetworkId5, resNodePubKey5, sdk.NewCoin(stratos.Wei, resNodeInitialDeposit5), resOwner5, registertypes.NewDescription("sds://resourceNode5", "", "", "", ""), nodeType)
	return createResourceNodeMsg
}

func printCurrNozPrice(t *testing.T, ctx sdk.Context, potKeeper potKeeper.Keeper, registerKeeper registerKeeper.Keeper, nozPriceFactorsBefore NozPriceFactors) (NozPriceFactors, sdk.Dec, sdk.Dec) {
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

	nozPricePercentage := nozPriceDelta.Quo(nozPriceFactorsBefore.NOzonePrice).MulInt(sdk.NewInt(100))
	//initialTotalDepositPercentage := initialTotalDepositDelta.Quo(nozPriceFactorsBefore.InitialTotalDeposit)
	//effectiveTotalDepositPercentage := effectiveTotalDepositDelta.Quo(nozPriceFactorsBefore.EffectiveTotalDeposit)
	//totalUnissuedPrepayPercentage := totalUnissuedPrepayDelta.Quo(nozPriceFactorsBefore.TotalUnissuedPrepay)
	//depositAndPrepayPercentage := depositAndPrepayDelta.Quo(nozPriceFactorsBefore.DepositAndPrepay)
	ozoneLimitPercentage := ozoneLimitDelta.ToDec().Quo(nozPriceFactorsBefore.OzoneLimit.ToDec()).MulInt(sdk.NewInt(100))

	t.Log("===>>>>>>>>>>>>>>     Current noz Price    ===>>>>>>>>>>>>>>")
	t.Log("NOzonePrice: 									" + nozPriceFactorsAfter.NOzonePrice.String() + "(delta: " + nozPriceDelta.String() + ", " + nozPricePercentage.String()[:5] + "%)")
	t.Log("InitialTotalDeposit: 							" + nozPriceFactorsAfter.InitialTotalDeposit.String() + "(delta: " + initialTotalDepositDelta.String() + ")")
	t.Log("EffectiveTotalDeposit: 							" + nozPriceFactorsAfter.EffectiveTotalDeposit.String() + "(delta: " + effectiveTotalDepositDelta.String() + ")")
	t.Log("TotalUnissuedPrepay: 							" + nozPriceFactorsAfter.TotalUnissuedPrepay.String() + "(delta: " + totalUnissuedPrepayDelta.String() + ")")
	t.Log("InitialTotalDeposit+TotalUnissuedPrepay:			" + nozPriceFactorsAfter.DepositAndPrepay.String() + "(delta: " + depositAndPrepayDelta.String() + ")")
	t.Log("OzoneLimit: 									" + nozPriceFactorsAfter.OzoneLimit.String() + "(delta: " + ozoneLimitDelta.String() + ", " + ozoneLimitPercentage.String()[:5] + "%)")
	t.Log("NozSupply: 									    " + nozPriceFactorsAfter.NozSupply.String() + "(delta: " + nozSupplyDelta.String() + ")")

	return nozPriceFactorsAfter, nozPricePercentage, ozoneLimitPercentage
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

// for main net
func checkResult(t *testing.T, ctx sdk.Context,
	k potKeeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankKeeper.Keeper,
	registerKeeper registerKeeper.Keeper,
	currentEpoch sdk.Int,
	lastFoundationAccBalance sdk.Coins,
	lastUnissuedPrepay sdk.Coin,
	lastCommunityPool sdk.Coins,
	lastMatureTotalOfResNode1 sdk.Coins,
	slashingAmtSetup sdk.Int,
	feeCollectorToFeePoolAtBeginBlock sdk.Coin) {

	currentSlashing := registerKeeper.GetSlashing(ctx, resNodeAddr2)
	t.Log("currentSlashing					= " + currentSlashing.String())

	individualRewardTotal := sdk.Coins{}
	newMatureEpoch := currentEpoch.Add(sdk.NewInt(k.MatureEpoch(ctx)))

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
		Sub(newFoundationAccBalance).
		Add(lastUnissuedPrepay).
		Sub(newUnissuedPrepay)
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
	matureTotalOfResNode1Change, _ := newMatureTotalOfResNode1.SafeSub(lastMatureTotalOfResNode1)
	if matureTotalOfResNode1Change == nil || matureTotalOfResNode1Change.IsAnyNegative() {
		matureTotalOfResNode1Change = sdk.Coins{}
	}
	t.Log("matureTotalOfResNode1Change		= " + matureTotalOfResNode1Change.String())
	require.Equal(t, matureTotalOfResNode1Change.String(), upcomingMaturedIndividual.String())
}

func checkValidator(t *testing.T, app *app.NewApp, addr sdk.ValAddress, expFound bool) stakingtypes.Validator {
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
	//************************** setup indexing nodes owners' accounts **************************
	idxOwnerAcc1 := &authtypes.BaseAccount{Address: idxOwner1.String()}
	idxOwnerAcc2 := &authtypes.BaseAccount{Address: idxOwner2.String()}
	idxOwnerAcc3 := &authtypes.BaseAccount{Address: idxOwner3.String()}
	//************************** setup validator delegators' accounts **************************
	valOwnerAcc1 := &authtypes.BaseAccount{Address: valOpAccAddr1.String()}
	//************************** setup indexing nodes' accounts **************************
	idxNodeAcc1 := &authtypes.BaseAccount{Address: idxNodeAddr1.String()}
	foundationDepositorAcc := &authtypes.BaseAccount{Address: foundationDepositorAccAddr.String()}

	accs := []authtypes.GenesisAccount{
		resOwnerAcc1, resOwnerAcc2, resOwnerAcc3, resOwnerAcc4, resOwnerAcc5,
		idxOwnerAcc1, idxOwnerAcc2, idxOwnerAcc3,
		valOwnerAcc1,
		foundationDepositorAcc,
		idxNodeAcc1,
	}

	balances := []banktypes.Balance{
		{
			Address: resOwner1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialDeposit1.Add(depositForSendingTx))},
		},
		{
			Address: resOwner2.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialDeposit2)},
		},
		{
			Address: resOwner3.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialDeposit3)},
		},
		{
			Address: resOwner4.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialDeposit4)},
		},
		{
			Address: resOwner5.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialDeposit5)},
		},
		{
			Address: idxOwner1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(idxNodeInitialDeposit1)},
		},
		{
			Address: idxOwner2.String(),
			Coins:   sdk.Coins{stratos.NewCoin(idxNodeInitialDeposit2)},
		},
		{
			Address: idxOwner3.String(),
			Coins:   sdk.Coins{stratos.NewCoin(idxNodeInitialDeposit3)},
		},
		{
			Address: valOpAccAddr1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(valInitialStake)},
		},
		{
			Address: idxNodeAddr1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(sdk.ZeroInt())},
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
	//resOwnerAcc1 := &authtypes.BaseAccount{Address: resOwner1.String()}
	//resOwnerAcc2 := &authtypes.BaseAccount{Address: resOwner2.String()}
	//resOwnerAcc3 := &authtypes.BaseAccount{Address: resOwner3.String()}
	//resOwnerAcc4 := &authtypes.BaseAccount{Address: resOwner4.String()}
	//resOwnerAcc5 := &authtypes.BaseAccount{Address: resOwner5.String()}
	//************************** setup indexing nodes owners' accounts **************************
	idxOwnerAcc1 := &authtypes.BaseAccount{Address: idxOwner1.String()}
	idxOwnerAcc2 := &authtypes.BaseAccount{Address: idxOwner2.String()}
	idxOwnerAcc3 := &authtypes.BaseAccount{Address: idxOwner3.String()}
	//************************** setup validator delegators' accounts **************************
	valOwnerAcc1 := &authtypes.BaseAccount{Address: valOpAccAddr1.String()}
	//************************** setup indexing nodes' accounts **************************
	idxNodeAcc1 := &authtypes.BaseAccount{Address: idxNodeAddr1.String()}
	foundationDepositorAcc := &authtypes.BaseAccount{Address: foundationDepositorAccAddr.String()}

	accs := []authtypes.GenesisAccount{
		//resOwnerAcc1, resOwnerAcc2, resOwnerAcc3, resOwnerAcc4, resOwnerAcc5,
		idxOwnerAcc1, idxOwnerAcc2, idxOwnerAcc3,
		valOwnerAcc1,
		foundationDepositorAcc,
		idxNodeAcc1,
	}

	balances := []banktypes.Balance{
		{
			Address: idxOwner1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(idxNodeInitialDeposit1)},
		},
		{
			Address: idxOwner2.String(),
			Coins:   sdk.Coins{stratos.NewCoin(idxNodeInitialDeposit2)},
		},
		{
			Address: idxOwner3.String(),
			Coins:   sdk.Coins{stratos.NewCoin(idxNodeInitialDeposit3)},
		},
		{
			Address: valOpAccAddr1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(valInitialStake)},
		},
		{
			Address: idxNodeAddr1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(sdk.ZeroInt())},
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
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialDepositForMultipleNodes.Add(depositForSendingTx))},
		})
	}

	return accs, balances
}

func setupAllResourceNodes() []registertypes.ResourceNode {

	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	nodeType := registertypes.STORAGE
	resourceNode1, _ := registertypes.NewResourceNode(resNodeNetworkId1, resNodePubKey1, resOwner1, registertypes.NewDescription("sds://resourceNode1", "", "", "", ""), nodeType, time)
	resourceNode2, _ := registertypes.NewResourceNode(resNodeNetworkId2, resNodePubKey2, resOwner2, registertypes.NewDescription("sds://resourceNode2", "", "", "", ""), nodeType, time)
	resourceNode3, _ := registertypes.NewResourceNode(resNodeNetworkId3, resNodePubKey3, resOwner3, registertypes.NewDescription("sds://resourceNode3", "", "", "", ""), nodeType, time)
	resourceNode4, _ := registertypes.NewResourceNode(resNodeNetworkId4, resNodePubKey4, resOwner4, registertypes.NewDescription("sds://resourceNode4", "", "", "", ""), nodeType, time)
	resourceNode5, _ := registertypes.NewResourceNode(resNodeNetworkId5, resNodePubKey5, resOwner5, registertypes.NewDescription("sds://resourceNode5", "", "", "", ""), nodeType, time)

	resourceNode1 = resourceNode1.AddToken(resNodeInitialDeposit1)
	resourceNode2 = resourceNode2.AddToken(resNodeInitialDeposit2)
	resourceNode3 = resourceNode3.AddToken(resNodeInitialDeposit3)
	resourceNode4 = resourceNode4.AddToken(resNodeInitialDeposit4)
	resourceNode5 = resourceNode5.AddToken(resNodeInitialDeposit5)

	resourceNode1.Status = stakingtypes.Bonded
	resourceNode2.Status = stakingtypes.Bonded
	resourceNode3.Status = stakingtypes.Bonded
	resourceNode4.Status = stakingtypes.Bonded
	resourceNode5.Status = stakingtypes.Bonded

	var resourceNodes []registertypes.ResourceNode
	resourceNodes = append(resourceNodes, resourceNode1)
	resourceNodes = append(resourceNodes, resourceNode2)
	resourceNodes = append(resourceNodes, resourceNode3)
	resourceNodes = append(resourceNodes, resourceNode4)
	resourceNodes = append(resourceNodes, resourceNode5)
	return resourceNodes
}

func setupMultipleResourceNodes(resOwnerPrivKeys []*secp256k1.PrivKey, resNodePubKeys []cryptotypes.PubKey, resOwners []sdk.AccAddress, resNodeNetworkIds []stratos.SdsAddress) []registertypes.ResourceNode {
	if len(resOwnerPrivKeys) != len(resNodePubKeys) ||
		len(resNodePubKeys) != len(resOwners) ||
		len(resOwners) != len(resNodeNetworkIds) {
		return nil
	}

	numOfNodes := len(resOwnerPrivKeys)
	resourceNodes := make([]registertypes.ResourceNode, 0, numOfNodes)

	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	nodeType := registertypes.STORAGE

	for i, _ := range resOwnerPrivKeys {
		resourceNodeTmp, _ := registertypes.NewResourceNode(resNodeNetworkIds[i], resNodePubKeys[i], resOwners[i], registertypes.NewDescription("sds://resourceNode"+strconv.Itoa(i+1), "", "", "", ""), nodeType, time)
		resourceNodeTmp = resourceNodeTmp.AddToken(resNodeInitialDepositForMultipleNodes)
		resourceNodeTmp.Status = stakingtypes.Bonded
		resourceNodes = append(resourceNodes, resourceNodeTmp)
	}

	return resourceNodes
}

func setupAllMetaNodes() []registertypes.MetaNode {
	var indexingNodes []registertypes.MetaNode

	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	indexingNode1, _ := registertypes.NewMetaNode(stratos.SdsAddress(idxNodeAddr1), idxNodePubKey1, idxOwner1, registertypes.NewDescription("sds://indexingNode1", "", "", "", ""), time)
	indexingNode2, _ := registertypes.NewMetaNode(stratos.SdsAddress(idxNodeAddr2), idxNodePubKey2, idxOwner2, registertypes.NewDescription("sds://indexingNode2", "", "", "", ""), time)
	indexingNode3, _ := registertypes.NewMetaNode(stratos.SdsAddress(idxNodeAddr3), idxNodePubKey3, idxOwner3, registertypes.NewDescription("sds://indexingNode3", "", "", "", ""), time)

	indexingNode1.Suspend = false
	indexingNode2.Suspend = false
	indexingNode3.Suspend = false

	indexingNode1 = indexingNode1.AddToken(idxNodeInitialDeposit1)
	indexingNode2 = indexingNode2.AddToken(idxNodeInitialDeposit2)
	indexingNode3 = indexingNode3.AddToken(idxNodeInitialDeposit3)

	indexingNode1.Status = stakingtypes.Bonded
	indexingNode2.Status = stakingtypes.Bonded
	indexingNode3.Status = stakingtypes.Bonded

	indexingNodes = append(indexingNodes, indexingNode1)
	indexingNodes = append(indexingNodes, indexingNode2)
	indexingNodes = append(indexingNodes, indexingNode3)

	return indexingNodes
}

var (
	initialTotalDepositStore   = sdk.NewInt(1500000000000)
	effectiveTotalDepositStore = sdk.NewInt(1500000000000)
	remainOzoneLimitStore      = sdk.NewInt(1500000000000000)
	totalUnissuedPrepayStore   = sdk.ZeroInt()
	nozPrice                   = sdk.NewDecWithPrec(1000000, 9)
	//priceChangeChan = make(chan PriceChangeEvent, 0)
)

type PriceChangeEvent struct {
	depositDelta        sdk.Int
	unissuedPrepayDelta sdk.Int
}

func simulatePriceChange(t *testing.T, priceChangeEvent *PriceChangeEvent, nozPriceFactorsBefore *NozPriceFactors) (*NozPriceFactors, sdk.Dec, sdk.Dec) {
	nozPriceFactorsAfter := &NozPriceFactors{}
	nozPriceFactorsAfter.InitialTotalDeposit = nozPriceFactorsBefore.InitialTotalDeposit
	nozPriceFactorsAfter.TotalUnissuedPrepay = nozPriceFactorsBefore.TotalUnissuedPrepay.Add(priceChangeEvent.unissuedPrepayDelta)
	nozPriceFactorsAfter.DepositAndPrepay = nozPriceFactorsAfter.InitialTotalDeposit.Add(nozPriceFactorsAfter.TotalUnissuedPrepay)
	nozPriceFactorsAfter.EffectiveTotalDeposit = nozPriceFactorsBefore.EffectiveTotalDeposit.Add(priceChangeEvent.depositDelta)
	deltaNozLimit := sdk.ZeroInt()
	nozPriceFactorsAfter.NozSupply = nozPriceFactorsBefore.NozSupply
	if !priceChangeEvent.depositDelta.Equal(sdk.ZeroInt()) {
		ozoneLimitChangeByDeposit := priceChangeEvent.depositDelta.ToDec().Quo(depositNozRateInt.ToDec()).TruncateInt()
		//ozoneLimitChangeByDeposit := nozPriceFactorsBefore.OzoneLimit.ToDec().Quo(nozPriceFactorsBefore.InitialTotalDeposit.ToDec()).Mul(priceChangeEvent.depositDelta.ToDec()).TruncateInt()
		deltaNozLimit = deltaNozLimit.Add(ozoneLimitChangeByDeposit)
		nozPriceFactorsAfter.NozSupply = nozPriceFactorsBefore.NozSupply.Add(ozoneLimitChangeByDeposit)
	}
	if !priceChangeEvent.unissuedPrepayDelta.Equal(sdk.ZeroInt()) {
		ozoneLimitChangeByPrepay := nozPriceFactorsBefore.OzoneLimit.ToDec().
			Mul(priceChangeEvent.unissuedPrepayDelta.ToDec()).
			Quo(nozPriceFactorsBefore.EffectiveTotalDeposit.Add(nozPriceFactorsBefore.TotalUnissuedPrepay).Add(priceChangeEvent.unissuedPrepayDelta).ToDec()).
			TruncateInt()
		//Sub(nozPriceFactorsBefore.OzoneLimit)
		if priceChangeEvent.unissuedPrepayDelta.GT(sdk.ZeroInt()) {
			// positive value of prepay leads to limit decrease
			deltaNozLimit = deltaNozLimit.Sub(ozoneLimitChangeByPrepay)
		} else {
			// nagative value of prepay (reward distribution) leads to limit increase
			deltaNozLimit = deltaNozLimit.Add(ozoneLimitChangeByPrepay)
		}
	}

	nozPriceFactorsAfter.OzoneLimit = nozPriceFactorsBefore.OzoneLimit.Add(deltaNozLimit)

	nozPriceFactorsAfter.NOzonePrice = nozPriceFactorsAfter.DepositAndPrepay.ToDec().Quo(nozPriceFactorsAfter.OzoneLimit.ToDec())
	nozPriceFactorsAfter.EffectiveTotalDeposit = nozPriceFactorsBefore.EffectiveTotalDeposit.Add(priceChangeEvent.depositDelta)

	nozPriceDelta := nozPriceFactorsAfter.NOzonePrice.Sub(nozPriceFactorsBefore.NOzonePrice)
	initialTotalDepositDelta := nozPriceFactorsAfter.InitialTotalDeposit.Sub(nozPriceFactorsBefore.InitialTotalDeposit)
	effectiveTotalDepositDelta := nozPriceFactorsAfter.EffectiveTotalDeposit.Sub(nozPriceFactorsBefore.EffectiveTotalDeposit)
	totalUnissuedPrepayDelta := nozPriceFactorsAfter.TotalUnissuedPrepay.Sub(nozPriceFactorsBefore.TotalUnissuedPrepay)
	depositAndPrepayDelta := nozPriceFactorsAfter.DepositAndPrepay.Sub(nozPriceFactorsBefore.DepositAndPrepay)
	ozoneLimitDelta := nozPriceFactorsAfter.OzoneLimit.Sub(nozPriceFactorsBefore.OzoneLimit)
	nozSupplyDelta := nozPriceFactorsAfter.NozSupply.Sub(nozPriceFactorsBefore.NozSupply)

	nozPricePercentage := nozPriceDelta.Quo(nozPriceFactorsBefore.NOzonePrice).MulInt(sdk.NewInt(100))
	//initialTotalDepositPercentage := initialTotalDepositDelta.Quo(nozPriceFactorsBefore.InitialTotalDeposit)
	//effectiveTotalDepositPercentage := effectiveTotalDepositDelta.Quo(nozPriceFactorsBefore.EffectiveTotalDeposit)
	//totalUnissuedPrepayPercentage := totalUnissuedPrepayDelta.Quo(nozPriceFactorsBefore.TotalUnissuedPrepay)
	//depositAndPrepayPercentage := depositAndPrepayDelta.Quo(nozPriceFactorsBefore.DepositAndPrepay)
	ozoneLimitPercentage := ozoneLimitDelta.ToDec().Quo(nozPriceFactorsBefore.OzoneLimit.ToDec()).MulInt(sdk.NewInt(100))

	t.Log("===>>>>>>>>>>>>>>     Current noz Price    ===>>>>>>>>>>>>>>")
	t.Log("NOzonePrice: 									" + nozPriceFactorsAfter.NOzonePrice.String() + "(delta: " + nozPriceDelta.String() + ", " + nozPricePercentage.String()[:5] + "%)")
	t.Log("InitialTotalDeposit: 							" + nozPriceFactorsAfter.InitialTotalDeposit.String() + "(delta: " + initialTotalDepositDelta.String() + ")")
	t.Log("EffectiveTotalDeposit: 							" + nozPriceFactorsAfter.EffectiveTotalDeposit.String() + "(delta: " + effectiveTotalDepositDelta.String() + ")")
	t.Log("TotalUnissuedPrepay: 							" + nozPriceFactorsAfter.TotalUnissuedPrepay.String() + "(delta: " + totalUnissuedPrepayDelta.String() + ")")
	t.Log("InitialTotalDeposit+TotalUnissuedPrepay:			" + nozPriceFactorsAfter.DepositAndPrepay.String() + "(delta: " + depositAndPrepayDelta.String() + ")")
	t.Log("OzoneLimit: 									" + nozPriceFactorsAfter.OzoneLimit.String() + "(delta: " + ozoneLimitDelta.String() + ", " + ozoneLimitPercentage.String()[:5] + "%)")
	t.Log("NozSupply: 				     					" + nozPriceFactorsAfter.NozSupply.String() + "(delta: " + nozSupplyDelta.String() + ")")

	return nozPriceFactorsAfter, nozPricePercentage, ozoneLimitPercentage
}

func TestOzPriceChangePrepay(t *testing.T) {
	NUM_OF_SAMPLE := 100
	dataToExcel := make([]NozPriceFactors, 0, NUM_OF_SAMPLE)
	/********************* initialize mock app *********************/
	//mApp, k, stakingKeeper, bankKeeper, supplyKeeper, registerKeeper := getMockApp(t)
	accs, balances := setupAccounts()
	//stApp := app.SetupWithGenesisAccounts(accs, chainID, balances...)
	validators := make([]*tmtypes.Validator, 0)
	valSet := tmtypes.NewValidatorSet(validators)
	metaNodes := setupAllMetaNodes()
	//resourceNodes := setupAllResourceNodes()
	resourceNodes := make([]registertypes.ResourceNode, 0)

	stApp := app.SetupWithGenesisNodeSet(t, true, valSet, metaNodes, resourceNodes, accs, chainID, balances...)

	accountKeeper := stApp.GetAccountKeeper()
	//bankKeeper := stApp.GetBankKeeper()
	registerKeeper := stApp.GetRegisterKeeper()
	potKeeper := stApp.GetPotKeeper()

	/********************* foundation account deposit *********************/
	header := tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := stApp.BaseApp.NewContext(true, header)

	foundationDepositMsg := pottypes.NewMsgFoundationDeposit(foundationDeposit, foundationDepositorAccAddr)
	txGen := app.MakeTestEncodingConfig().TxConfig

	foundationDepositorAcc := accountKeeper.GetAccount(ctx, foundationDepositorAccAddr)
	accNum := foundationDepositorAcc.GetAccountNumber()
	accSeq := foundationDepositorAcc.GetSequence()
	_, _, err := app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, foundationDepositorPrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(pottypes.FoundationAccount)
	app.CheckBalance(t, stApp, foundationAccountAddr, foundationDeposit)

	/********************* create validator with 50% commission *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	commission := stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	description := stakingtypes.NewDescription("foo_moniker", chainID, "", "", "")
	createValidatorMsg, err := stakingtypes.NewMsgCreateValidator(valOpValAddr1, valConsPubk1, stratos.NewCoin(valInitialStake), description, commission, sdk.OneInt())

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

		resOwnerAcc := accountKeeper.GetAccount(ctx, resOwner1)
		ownerAccNum := resOwnerAcc.GetAccountNumber()
		ownerAccSeq := resOwnerAcc.GetSequence()

		_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{prepayMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, resOwnerPrivKey1)
		require.NoError(t, err)
		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)
		priceAfter, nozPricePercentage, ozoneLimitPercentage = printCurrNozPrice(t, ctx, potKeeper, registerKeeper, priceBefore)
		dataToExcel = append(dataToExcel, priceAfter)
		require.True(t, nozPricePercentage.GT(sdk.ZeroDec()), "noz price should increase after PREPAY")
		require.True(t, ozoneLimitPercentage.LT(sdk.ZeroDec()), "OzLimit should not change after PREPAY")
		t.Log("********************************* Deliver Prepay Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

		priceBefore = priceAfter
	}
	exportToCSV(t, dataToExcel)
}

func TestOzPriceChangeVolumeReport(t *testing.T) {
	NUM_OF_SAMPLE := 100
	dataToExcel := make([]NozPriceFactors, 0, NUM_OF_SAMPLE)
	/********************* initialize mock app *********************/
	//mApp, k, stakingKeeper, bankKeeper, supplyKeeper, registerKeeper := getMockApp(t)
	accs, balances := setupAccounts()
	//stApp := app.SetupWithGenesisAccounts(accs, chainID, balances...)
	validators := make([]*tmtypes.Validator, 0)
	valSet := tmtypes.NewValidatorSet(validators)
	metaNodes := setupAllMetaNodes()
	//resourceNodes := setupAllResourceNodes()
	resourceNodes := make([]registertypes.ResourceNode, 0)

	stApp := app.SetupWithGenesisNodeSet(t, true, valSet, metaNodes, resourceNodes, accs, chainID, balances...)

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
	txGen := app.MakeTestEncodingConfig().TxConfig

	foundationDepositorAcc := accountKeeper.GetAccount(ctx, foundationDepositorAccAddr)
	accNum := foundationDepositorAcc.GetAccountNumber()
	accSeq := foundationDepositorAcc.GetSequence()
	_, _, err := app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, foundationDepositorPrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(pottypes.FoundationAccount)
	app.CheckBalance(t, stApp, foundationAccountAddr, foundationDeposit)

	/********************* create validator with 50% commission *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	commission := stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	description := stakingtypes.NewDescription("foo_moniker", chainID, "", "", "")
	createValidatorMsg, err := stakingtypes.NewMsgCreateValidator(valOpValAddr1, valConsPubk1, stratos.NewCoin(valInitialStake), description, commission, sdk.OneInt())

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

		resOwnerAcc := accountKeeper.GetAccount(ctx, resOwner1)
		ownerAccNum := resOwnerAcc.GetAccountNumber()
		ownerAccSeq := resOwnerAcc.GetSequence()

		_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{prepayMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, resOwnerPrivKey1)
		require.NoError(t, err)
		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)
		priceAfter, nozPricePercentage, ozoneLimitPercentage = printCurrNozPrice(t, ctx, potKeeper, registerKeeper, priceBefore)
		dataToExcel = append(dataToExcel, priceAfter)
		require.True(t, nozPricePercentage.GT(sdk.ZeroDec()), "noz price should increase after PREPAY")
		require.True(t, ozoneLimitPercentage.LT(sdk.ZeroDec()), "OzLimit should decrease after PREPAY")
		t.Log("********************************* Deliver Prepay Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

		priceBefore = priceAfter
	}

	for i := 0; i < NUM_OF_SAMPLE; i++ {
		t.Log("********************************* Deliver VolumeReport Tx START ********************************************")
		/********************* prepare tx data *********************/
		volumeReportMsg := setupMsgVolumeReport(int64(i + 1))

		lastTotalMinedToken := potKeeper.GetTotalMinedTokens(ctx)
		t.Log("last committed TotalMinedTokens = " + lastTotalMinedToken.String())
		_, ok := sdk.NewIntFromString(volumeReportMsg.Epoch.String())
		require.Equal(t, ok, true)
		totalConsumedNoz := potKeeper.GetTotalConsumedNoz(volumeReportMsg.WalletVolumes).ToDec()

		/********************* print info *********************/
		t.Log("epoch " + volumeReportMsg.Epoch.String())
		S := registerKeeper.GetInitialGenesisDepositTotal(ctx).ToDec()
		Pt := registerKeeper.GetTotalUnissuedPrepay(ctx).Amount.ToDec()
		Y := totalConsumedNoz
		Lt := registerKeeper.GetRemainingOzoneLimit(ctx).ToDec()
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

		t.Log("indexing_wallet1:  address = " + idxOwner1.String())
		t.Log("              miningReward = " + rewardDetailMap[idxOwner1.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[idxOwner1.String()].RewardFromTrafficPool.String())

		t.Log("indexing_wallet2:  address = " + idxOwner2.String())
		t.Log("              miningReward = " + rewardDetailMap[idxOwner2.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[idxOwner2.String()].RewardFromTrafficPool.String())

		t.Log("indexing_wallet3:  address = " + idxOwner3.String())
		t.Log("              miningReward = " + rewardDetailMap[idxOwner3.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[idxOwner3.String()].RewardFromTrafficPool.String())
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

		resOwnerAcc := accountKeeper.GetAccount(ctx, idxOwner1)
		ownerAccNum := resOwnerAcc.GetAccountNumber()
		ownerAccSeq := resOwnerAcc.GetSequence()

		feePoolAccAddr := accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		require.NotNil(t, feePoolAccAddr)
		_ = bankKeeper.GetBalance(ctx, feePoolAccAddr, potKeeper.BondDenom(ctx))
		//feeCollectorToFeePoolAtBeginBlock := bankKeeper.GetBalance(ctx, feePoolAccAddr, potKeeper.BondDenom(ctx))

		_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{volumeReportMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, idxOwnerPrivKey1)
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

	resOwners := make([]sdk.AccAddress, 0, NUM_OF_SAMPLE)
	resOwnerPrivKeys := make([]*secp256k1.PrivKey, 0, NUM_OF_SAMPLE)
	resOwnerPubkeys := make([]cryptotypes.PubKey, 0, NUM_OF_SAMPLE)
	resNodeNetworkIds := make([]stratos.SdsAddress, 0, NUM_OF_SAMPLE)

	for i := 0; i < NUM_OF_SAMPLE; i++ {
		resOwnerPrivKeyTmp := secp256k1.GenPrivKey()
		resOwnerPrivKeys = append(resOwnerPrivKeys, resOwnerPrivKeyTmp)
		resOwnerPubKeyTmp := resOwnerPrivKeyTmp.PubKey()
		resOwnerPubkeys = append(resOwnerPubkeys, resOwnerPrivKeyTmp.PubKey())
		resNodeAddrTmp := sdk.AccAddress(resOwnerPubKeyTmp.Address())
		resOwners = append(resOwners, resNodeAddrTmp)
		resNodeNetworkIds = append(resNodeNetworkIds, stratos.SdsAddress(resNodeAddrTmp))
	}

	/********************* initialize mock app *********************/
	//mApp, k, stakingKeeper, bankKeeper, supplyKeeper, registerKeeper := getMockApp(t)
	accs, balances := setupAccountsMultipleResNodes(resOwners)
	//stApp := app.SetupWithGenesisAccounts(accs, chainID, balances...)
	validators := make([]*tmtypes.Validator, 0)
	valSet := tmtypes.NewValidatorSet(validators)
	metaNodes := setupAllMetaNodes()
	//resourceNodes := setupAllResourceNodes()
	resourceNodes := make([]registertypes.ResourceNode, 0)

	stApp := app.SetupWithGenesisNodeSet(t, true, valSet, metaNodes, resourceNodes, accs, chainID, balances...)

	accountKeeper := stApp.GetAccountKeeper()
	//bankKeeper := stApp.GetBankKeeper()
	registerKeeper := stApp.GetRegisterKeeper()
	potKeeper := stApp.GetPotKeeper()

	/********************* foundation account deposit *********************/
	header := tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := stApp.BaseApp.NewContext(true, header)

	foundationDepositMsg := pottypes.NewMsgFoundationDeposit(foundationDeposit, foundationDepositorAccAddr)
	txGen := app.MakeTestEncodingConfig().TxConfig

	foundationDepositorAcc := accountKeeper.GetAccount(ctx, foundationDepositorAccAddr)
	accNum := foundationDepositorAcc.GetAccountNumber()
	accSeq := foundationDepositorAcc.GetSequence()
	_, _, err := app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, foundationDepositorPrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(pottypes.FoundationAccount)
	app.CheckBalance(t, stApp, foundationAccountAddr, foundationDeposit)

	/********************* create validator with 50% commission *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	commission := stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	description := stakingtypes.NewDescription("foo_moniker", chainID, "", "", "")
	createValidatorMsg, err := stakingtypes.NewMsgCreateValidator(valOpValAddr1, valConsPubk1, stratos.NewCoin(valInitialStake), description, commission, sdk.OneInt())

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

	resOwnerAcc := accountKeeper.GetAccount(ctx, resOwners[0])
	ownerAccNum := resOwnerAcc.GetAccountNumber()
	ownerAccSeq := resOwnerAcc.GetSequence()

	_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{prepayMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, resOwnerPrivKeys[0])
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
		createResourceNodeMsg := setupMsgCreateResourceNode(i, resNodeNetworkIds[i], resOwnerPubkeys[i], resOwners[i])
		/********************* deliver tx *********************/

		resOwnerAcc := accountKeeper.GetAccount(ctx, resOwners[i])
		ownerAccNum := resOwnerAcc.GetAccountNumber()
		ownerAccSeq := resOwnerAcc.GetSequence()

		_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createResourceNodeMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, resOwnerPrivKeys[i])
		require.NoError(t, err)
		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)

		unsuspendMsg := setupUnsuspendMsgByIndex(i, resNodeNetworkIds[i], resOwnerPubkeys[i], resOwners[i])
		/********************* deliver tx *********************/

		idxOwnerAcc := accountKeeper.GetAccount(ctx, idxOwner1)
		ownerAccNum = idxOwnerAcc.GetAccountNumber()
		ownerAccSeq = idxOwnerAcc.GetSequence()

		_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{unsuspendMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, idxOwnerPrivKey1)
		require.NoError(t, err)
		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)

		priceAfter, nozPricePercentage, ozoneLimitPercentage = printCurrNozPrice(t, ctx, potKeeper, registerKeeper, priceBefore)
		require.True(t, nozPricePercentage.LT(sdk.ZeroDec()), "noz price should decrease after CreateResourceNode")
		require.True(t, ozoneLimitPercentage.GT(sdk.ZeroDec()), "OzLimit should increase after CreateResourceNode")
		t.Log("********************************* Deliver Create and unsuspend ResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

		dataToExcel = append(dataToExcel, priceAfter)
		priceBefore = priceAfter
	}

	for i := 0; i < NUM_OF_SAMPLE; i++ {
		slashingMsg := setupSuspendMsgByIndex(i, resNodeNetworkIds[i], resOwnerPubkeys[i], resOwners[i])
		/********************* deliver tx *********************/

		idxOwnerAcc := accountKeeper.GetAccount(ctx, idxOwner1)
		ownerAccNum = idxOwnerAcc.GetAccountNumber()
		ownerAccSeq = idxOwnerAcc.GetSequence()

		_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{slashingMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, idxOwnerPrivKey1)
		require.NoError(t, err)
		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)

		//createResourceNodeMsg := setupMsgRemoveResourceNode(i, resNodeNetworkIds[i], resOwners[i])
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
		require.True(t, nozPricePercentage.GT(sdk.ZeroDec()), "noz price should increase after RemoveResourceNode")
		require.True(t, ozoneLimitPercentage.LT(sdk.ZeroDec()), "OzLimit should decrease after RemoveResourceNode")
		t.Log("********************************* Deliver CreateResourceNode Tx END ********************************************\n\n...\n[NEXT TEST CASE]")

		dataToExcel = append(dataToExcel, priceAfter)
		priceBefore = priceAfter
	}

	exportToCSV(t, dataToExcel)

}

func TestOzPriceChangeRemoveMultipleResourceNodeAfterGenesis(t *testing.T) {
	NUM_OF_SAMPLE := 100
	dataToExcel := make([]NozPriceFactors, 0, NUM_OF_SAMPLE)

	resOwners := make([]sdk.AccAddress, 0, NUM_OF_SAMPLE)
	resOwnerPrivKeys := make([]*secp256k1.PrivKey, 0, NUM_OF_SAMPLE)
	resOwnerPubkeys := make([]cryptotypes.PubKey, 0, NUM_OF_SAMPLE)
	resNodeNetworkIds := make([]stratos.SdsAddress, 0, NUM_OF_SAMPLE)

	for i := 0; i < NUM_OF_SAMPLE; i++ {
		resOwnerPrivKeyTmp := secp256k1.GenPrivKey()
		resOwnerPrivKeys = append(resOwnerPrivKeys, resOwnerPrivKeyTmp)
		resOwnerPubKeyTmp := resOwnerPrivKeyTmp.PubKey()
		resOwnerPubkeys = append(resOwnerPubkeys, resOwnerPrivKeyTmp.PubKey())
		resNodeAddrTmp := sdk.AccAddress(resOwnerPubKeyTmp.Address())
		resOwners = append(resOwners, resNodeAddrTmp)
		resNodeNetworkIds = append(resNodeNetworkIds, stratos.SdsAddress(resNodeAddrTmp))
	}

	/********************* initialize mock app *********************/
	//mApp, k, stakingKeeper, bankKeeper, supplyKeeper, registerKeeper := getMockApp(t)
	accs, balances := setupAccountsMultipleResNodes(resOwners)
	//stApp := app.SetupWithGenesisAccounts(accs, chainID, balances...)
	validators := make([]*tmtypes.Validator, 0)
	valSet := tmtypes.NewValidatorSet(validators)
	metaNodes := setupAllMetaNodes()
	resourceNodes := setupMultipleResourceNodes(resOwnerPrivKeys, resOwnerPubkeys, resOwners, resNodeNetworkIds)

	stApp := app.SetupWithGenesisNodeSet(t, true, valSet, metaNodes, resourceNodes, accs, chainID, balances...)

	accountKeeper := stApp.GetAccountKeeper()
	//bankKeeper := stApp.GetBankKeeper()
	registerKeeper := stApp.GetRegisterKeeper()
	potKeeper := stApp.GetPotKeeper()

	/********************* foundation account deposit *********************/
	header := tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := stApp.BaseApp.NewContext(true, header)

	foundationDepositMsg := pottypes.NewMsgFoundationDeposit(foundationDeposit, foundationDepositorAccAddr)
	txGen := app.MakeTestEncodingConfig().TxConfig

	foundationDepositorAcc := accountKeeper.GetAccount(ctx, foundationDepositorAccAddr)
	accNum := foundationDepositorAcc.GetAccountNumber()
	accSeq := foundationDepositorAcc.GetSequence()
	_, _, err := app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, foundationDepositorPrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(pottypes.FoundationAccount)
	app.CheckBalance(t, stApp, foundationAccountAddr, foundationDeposit)

	/********************* create validator with 50% commission *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	commission := stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	description := stakingtypes.NewDescription("foo_moniker", chainID, "", "", "")
	createValidatorMsg, err := stakingtypes.NewMsgCreateValidator(valOpValAddr1, valConsPubk1, stratos.NewCoin(valInitialStake), description, commission, sdk.OneInt())

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

	for i := 0; i < NUM_OF_SAMPLE; i++ {
		unsuspendMsg := setupUnsuspendMsgByIndex(i, resNodeNetworkIds[i], resOwnerPubkeys[i], resOwners[i])
		/********************* deliver tx *********************/

		idxOwnerAcc := accountKeeper.GetAccount(ctx, idxOwner1)
		ownerAccNum := idxOwnerAcc.GetAccountNumber()
		ownerAccSeq := idxOwnerAcc.GetSequence()

		_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{unsuspendMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, idxOwnerPrivKey1)
		require.NoError(t, err)
		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)

	}

	// start testing
	t.Log("********************************* Deliver RemoveResourceNode Tx START ********************************************")

	prepayMsg := setupPrepayMsgWithResOwner(resOwners[0])
	/********************* deliver tx *********************/

	resOwnerAcc := accountKeeper.GetAccount(ctx, resOwners[0])
	ownerAccNum := resOwnerAcc.GetAccountNumber()
	ownerAccSeq := resOwnerAcc.GetSequence()

	_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{prepayMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, resOwnerPrivKeys[0])
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
		slashingMsg := setupSuspendMsgByIndex(i, resNodeNetworkIds[i], resOwnerPubkeys[i], resOwners[i])
		/********************* deliver tx *********************/

		idxOwnerAcc := accountKeeper.GetAccount(ctx, idxOwner1)
		ownerAccNum = idxOwnerAcc.GetAccountNumber()
		ownerAccSeq = idxOwnerAcc.GetSequence()

		_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{slashingMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, idxOwnerPrivKey1)
		require.NoError(t, err)
		/********************* commit & check result *********************/
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)
		//createResourceNodeMsg := setupMsgRemoveResourceNode(i, resNodeNetworkIds[i], resOwners[i])
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
		require.True(t, nozPricePercentage.GT(sdk.ZeroDec()), "noz price should increase after RemoveResourceNode")
		require.True(t, ozoneLimitPercentage.LT(sdk.ZeroDec()), "OzLimit should decrease after RemoveResourceNode")
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
