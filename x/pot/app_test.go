package pot_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/stratosnet/stratos-chain/app"
	stratos "github.com/stratosnet/stratos-chain/types"
	potKeeper "github.com/stratosnet/stratos-chain/x/pot/keeper"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	registerKeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
)

const (
	chainID     = "testchain_1-1"
	stos2wei    = stratos.StosToWei
	rewardDenom = stratos.Utros

	stopFlagOutOfTotalMiningReward = true
	stopFlagSpecificMinedReward    = false
	stopFlagSpecificEpoch          = true
)

var (
	paramSpecificMinedReward = sdk.NewCoins(stratos.NewCoinInt64(160000000000))
	paramSpecificEpoch       = sdk.NewInt(10)

	resNodeSlashingNOZAmt1            = sdk.NewInt(1000000000000000000)
	resNodeSlashingEffectiveTokenAmt1 = sdk.NewInt(1000000000000000000)

	resourceNodeVolume1 = sdk.NewInt(500000)
	resourceNodeVolume2 = sdk.NewInt(300000)
	resourceNodeVolume3 = sdk.NewInt(200000)

	depositForSendingTx, _    = sdk.NewIntFromString("100000000000000000000000000000")
	totalUnissuedPrepayVal, _ = sdk.NewIntFromString("1000000000000")
	totalUnissuedPrepay       = stratos.NewCoin(totalUnissuedPrepayVal)

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

	resNodePubKey1       = secp256k1.GenPrivKey().PubKey()
	resNodeAddr1         = sdk.AccAddress(resNodePubKey1.Address())
	resNodeNetworkId1    = stratos.SdsAddress(resNodePubKey1.Address())
	resNodeInitialStake1 = sdk.NewInt(3 * stos2wei)

	resNodePubKey2       = secp256k1.GenPrivKey().PubKey()
	resNodeAddr2         = sdk.AccAddress(resNodePubKey2.Address())
	resNodeNetworkId2    = stratos.SdsAddress(resNodePubKey2.Address())
	resNodeInitialStake2 = sdk.NewInt(3 * stos2wei)

	resNodePubKey3       = secp256k1.GenPrivKey().PubKey()
	resNodeAddr3         = sdk.AccAddress(resNodePubKey3.Address())
	resNodeNetworkId3    = stratos.SdsAddress(resNodePubKey3.Address())
	resNodeInitialStake3 = sdk.NewInt(3 * stos2wei)

	resNodePubKey4       = secp256k1.GenPrivKey().PubKey()
	resNodeAddr4         = sdk.AccAddress(resNodePubKey4.Address())
	resNodeNetworkId4    = stratos.SdsAddress(resNodePubKey4.Address())
	resNodeInitialStake4 = sdk.NewInt(3 * stos2wei)

	resNodePubKey5       = secp256k1.GenPrivKey().PubKey()
	resNodeAddr5         = sdk.AccAddress(resNodePubKey5.Address())
	resNodeNetworkId5    = stratos.SdsAddress(resNodePubKey5.Address())
	resNodeInitialStake5 = sdk.NewInt(3 * stos2wei)

	idxNodePrivKey1      = secp256k1.GenPrivKey()
	idxNodePubKey1       = idxNodePrivKey1.PubKey()
	idxNodeAddr1         = sdk.AccAddress(idxNodePubKey1.Address())
	idxNodeNetworkId1    = stratos.SdsAddress(idxNodePubKey1.Address())
	idxNodeInitialStake1 = sdk.NewInt(5 * stos2wei)

	idxNodePubKey2       = secp256k1.GenPrivKey().PubKey()
	idxNodeAddr2         = sdk.AccAddress(idxNodePubKey2.Address())
	idxNodeNetworkId2    = stratos.SdsAddress(idxNodePubKey2.Address())
	idxNodeInitialStake2 = sdk.NewInt(5 * stos2wei)

	idxNodePubKey3       = secp256k1.GenPrivKey().PubKey()
	idxNodeAddr3         = sdk.AccAddress(idxNodePubKey3.Address())
	idxNodeNetworkId3    = stratos.SdsAddress(idxNodePubKey3.Address())
	idxNodeInitialStake3 = sdk.NewInt(5 * stos2wei)

	valOpPrivKey1 = secp256k1.GenPrivKey()
	valOpPubKey1  = valOpPrivKey1.PubKey()
	valOpValAddr1 = sdk.ValAddress(valOpPubKey1.Address())
	valOpAccAddr1 = sdk.AccAddress(valOpPubKey1.Address())

	valConsPrivKey1 = ed25519.GenPrivKey()
	valConsPubk1    = valConsPrivKey1.PubKey()
	valInitialStake = sdk.NewInt(5 * stos2wei)
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
	slashingMsg := types.NewMsgSlashingResourceNode(reporters, reportOwner, resNodeNetworkId1, resOwner1, resNodeSlashingNOZAmt1, true, resNodeSlashingEffectiveTokenAmt1)
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

	stApp := app.SetupWithGenesisNodeSet(t, false, valSet, metaNodes, resourceNodes, accs, totalUnissuedPrepay, chainID, balances...)

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
	ctx = stApp.BaseApp.NewContext(false, header)

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

	/********************** loop sending volume report **********************/
	var i int64
	var slashingAmtSetup sdk.Int
	i = 0
	slashingAmtSetup = sdk.ZeroInt()
	for {

		/********************* test slashing msg when i==2 *********************/
		if i == 2 {
			t.Log("********************************* Deliver Slashing Tx START ********************************************")
			slashingMsg := setupSlashingMsg()
			/********************* deliver tx *********************/

			idxOwnerAcc1 := accountKeeper.GetAccount(ctx, idxOwner1)
			ownerAccNum := idxOwnerAcc1.GetAccountNumber()
			ownerAccSeq := idxOwnerAcc1.GetSequence()

			_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{slashingMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, idxOwnerPrivKey1)
			require.NoError(t, err)
			/********************* commit & check result *********************/
			header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
			stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
			ctx = stApp.BaseApp.NewContext(true, header)

			slashingAmtSetup = registerKeeper.GetSlashing(ctx, resOwner1)

			totalConsumedNoz := resNodeSlashingNOZAmt1.ToDec()

			slashingAmtCheck := potKeeper.GetTrafficReward(ctx, totalConsumedNoz)
			t.Log("slashingAmtSetup = " + slashingAmtSetup.String())
			require.Equal(t, slashingAmtSetup, slashingAmtCheck.TruncateInt())

			t.Log("********************************* Deliver Slashing Tx END ********************************************")
		}

		t.Log("*****************************************************************************")
		t.Log("*")
		t.Log("*                    height = ", header.GetHeight())
		t.Log("*")
		t.Log("*****************************************************************************")
		/********************* prepare tx data *********************/
		volumeReportMsg := setupMsgVolumeReport(i + 1)

		lastTotalMinedToken := potKeeper.GetTotalMinedTokens(ctx)
		t.Log("last committed TotalMinedTokens = " + lastTotalMinedToken.String())
		epoch, ok := sdk.NewIntFromString(volumeReportMsg.Epoch.String())
		require.Equal(t, ok, true)

		if isNeedStop(ctx, potKeeper, epoch, lastTotalMinedToken) {
			break
		}

		totalConsumedNoz := potKeeper.GetTotalConsumedNoz(volumeReportMsg.WalletVolumes).ToDec()

		/********************* print info *********************/
		t.Log("epoch " + volumeReportMsg.Epoch.String())
		S := registerKeeper.GetInitialGenesisStakeTotal(ctx).ToDec()
		Pt := registerKeeper.GetTotalUnissuedPrepay(ctx).Amount.ToDec()
		Y := totalConsumedNoz
		Lt := registerKeeper.GetRemainingOzoneLimit(ctx).ToDec()
		R := S.Add(Pt).Mul(Y).Quo(Lt.Add(Y))
		//t.Log("R = (S + Pt) * Y / (Lt + Y)")
		t.Log("S=" + S.String() + "\nPt=" + Pt.String() + "\nY=" + Y.String() + "\nLt=" + Lt.String() + "\nR=" + R.String() + "\n")

		t.Log("---------------------------")
		potKeeper.InitVariable(ctx)
		distributeGoal := types.InitDistributeGoal()
		distributeGoal, err := potKeeper.CalcTrafficRewardInTotal(ctx, distributeGoal, totalConsumedNoz)
		require.NoError(t, err)

		distributeGoal, err = potKeeper.CalcMiningRewardInTotal(ctx, distributeGoal) //for main net
		require.NoError(t, err)
		t.Log(distributeGoal.String())

		t.Log("---------------------------")
		t.Log("distribute detail:")
		rewardDetailMap := make(map[string]types.Reward)
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
		lastFoundationAccBalance := bankKeeper.GetAllBalances(ctx, foundationAccountAddr)
		lastUnissuedPrepay := registerKeeper.GetTotalUnissuedPrepay(ctx)
		lastCommunityPool := sdk.NewCoins(sdk.NewCoin(potKeeper.BondDenom(ctx), potKeeper.DistrKeeper.GetFeePool(ctx).CommunityPool.AmountOf(potKeeper.BondDenom(ctx)).TruncateInt()))
		lastMatureTotalOfResNode1 := potKeeper.GetMatureTotalReward(ctx, resOwner1)

		/********************* deliver tx *********************/
		idxOwnerAcc1 := accountKeeper.GetAccount(ctx, idxOwner1)
		ownerAccNum := idxOwnerAcc1.GetAccountNumber()
		ownerAccSeq := idxOwnerAcc1.GetSequence()

		feePoolAccAddr := accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		require.NotNil(t, feePoolAccAddr)
		feeCollectorToFeePoolAtBeginBlock := bankKeeper.GetBalance(ctx, feePoolAccAddr, potKeeper.BondDenom(ctx))

		t.Log("--------------------------- deliver volumeReportMsg")
		_, _, err = app.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{volumeReportMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, idxOwnerPrivKey1)
		require.NoError(t, err)

		/********************* commit & check result *********************/
		// reward distribution start at height = height + 1 where volume report tx executed
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		stApp.EndBlock(abci.RequestEndBlock{Height: header.Height})
		stApp.Commit()

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

		i++
	}
}

// return : coins - slashing
func deductSlashingAmt(ctx sdk.Context, coins sdk.Coins, slashing sdk.Coin) (ret sdk.Coins) {
	slashingDenom := slashing.Denom
	rewardToken := sdk.NewCoin(slashingDenom, coins.AmountOf(slashingDenom))
	if rewardToken.IsGTE(slashing) {
		ret = coins.Sub(sdk.NewCoins(slashing))
	} else {
		ret = coins.Sub(sdk.NewCoins(rewardToken))
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
	initialSlashingAmt sdk.Int,
	feeCollectorToFeePoolAtBeginBlock sdk.Coin) {

	// print individual reward
	individualRewardTotal := sdk.Coins{}
	newMatureEpoch := currentEpoch.Add(sdk.NewInt(k.MatureEpoch(ctx)))
	k.IteratorIndividualReward(ctx, newMatureEpoch, func(walletAddress sdk.AccAddress, individualReward types.Reward) (stop bool) {
		individualRewardTotal = individualRewardTotal.Add(individualReward.RewardFromTrafficPool...).Add(individualReward.RewardFromMiningPool...)
		t.Log("individualReward of [" + walletAddress.String() + "] = " + individualReward.String())
		return false
	})
	t.Log("---------------------------")
	k.IteratorMatureTotal(ctx, func(walletAddress sdk.AccAddress, matureTotal sdk.Coins) (stop bool) {
		t.Log("MatureTotal of [" + walletAddress.String() + "] = " + matureTotal.String())
		return false
	})
	t.Log("---------------------------")

	feeCollectorAccAddr := accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	require.NotNil(t, feeCollectorAccAddr)
	foundationAccountAddr := accountKeeper.GetModuleAddress(types.FoundationAccount)
	newFoundationAccBalance := bankKeeper.GetAllBalances(ctx, foundationAccountAddr)
	newUnissuedPrepay := sdk.NewCoins(registerKeeper.GetTotalUnissuedPrepay(ctx))
	newCommunityPool := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), k.DistrKeeper.GetFeePool(ctx).CommunityPool.AmountOf(k.BondDenom(ctx)).TruncateInt()))

	t.Log("resource node 1 initial slashingAmt        = " + initialSlashingAmt.String())
	currentSlashingAmt := registerKeeper.GetSlashing(ctx, resOwner1)
	t.Log("resource node 1 currentSlashingAmt         = " + currentSlashingAmt.String())
	slashingDeducted := sdk.NewCoin(k.RewardDenom(ctx), initialSlashingAmt.Sub(currentSlashingAmt))
	t.Log("resource node 1 slashing deducted          = " + slashingDeducted.String())
	matureTotal := k.GetMatureTotalReward(ctx, resOwner1)
	immatureTotal := k.GetImmatureTotalReward(ctx, resOwner1)
	t.Log("resource node 1 matureTotal                = " + matureTotal.String())
	t.Log("resource node 1 immatureTotal              = " + immatureTotal.String())

	// distribution module will send all tokens from "fee_collector" to "distribution" account in the BeginBlocker() method
	feeCollectorValChange := bankKeeper.GetAllBalances(ctx, feeCollectorAccAddr)
	t.Log("reward for validator send to fee_collector = " + feeCollectorValChange.String())
	communityTaxChange := newCommunityPool.Sub(lastCommunityPool).Sub(sdk.NewCoins(feeCollectorToFeePoolAtBeginBlock))
	t.Log("community tax change in community_pool     = " + communityTaxChange.String())
	t.Log("community_pool amount of wei               = " + newCommunityPool.String())

	rewardSrcChange := lastFoundationAccBalance.
		Sub(newFoundationAccBalance).
		Add(lastUnissuedPrepay).
		Sub(newUnissuedPrepay)
	t.Log("rewardSrcChange                            = " + rewardSrcChange.String())

	rewardDestChange := feeCollectorValChange.
		Add(individualRewardTotal...).
		Add(communityTaxChange...)

	t.Log("rewardDestChange                           = " + rewardDestChange.String())

	require.Equal(t, rewardSrcChange, rewardDestChange)

	t.Log("************************ slashing test***********************************")
	t.Log("slashing change                            = " + slashingDeducted.String())

	upcomingMaturedIndividual := sdk.Coins{}
	individualReward, found := k.GetIndividualReward(ctx, resOwner1, currentEpoch)
	if found {
		tmp := individualReward.RewardFromTrafficPool.Add(individualReward.RewardFromMiningPool...)
		upcomingMaturedIndividual = deductSlashingAmt(ctx, tmp, slashingDeducted)
	}
	t.Log("upcomingMaturedIndividual                  = " + upcomingMaturedIndividual.String())

	// get mature total changes
	newMatureTotalOfResNode1 := k.GetMatureTotalReward(ctx, resOwner1)
	matureTotalOfResNode1Change, _ := newMatureTotalOfResNode1.SafeSub(lastMatureTotalOfResNode1)
	if matureTotalOfResNode1Change == nil || matureTotalOfResNode1Change.IsAnyNegative() {
		matureTotalOfResNode1Change = sdk.Coins{}
	}
	t.Log("matureTotalOfResNode1Change                = " + matureTotalOfResNode1Change.String())
	require.Equal(t, matureTotalOfResNode1Change.String(), upcomingMaturedIndividual.String())

	totalRewardPoolAddr := accountKeeper.GetModuleAddress(types.TotalRewardPool)
	totalRewardPoolBalance := bankKeeper.GetAllBalances(ctx, totalRewardPoolAddr)
	t.Log("totalRewardPoolBalance                     = " + totalRewardPoolBalance.String())
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
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialStake1.Add(depositForSendingTx))},
		},
		{
			Address: resOwner2.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialStake2)},
		},
		{
			Address: resOwner3.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialStake3)},
		},
		{
			Address: resOwner4.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialStake4)},
		},
		{
			Address: resOwner5.String(),
			Coins:   sdk.Coins{stratos.NewCoin(resNodeInitialStake5)},
		},
		{
			Address: idxOwner1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(idxNodeInitialStake1)},
		},
		{
			Address: idxOwner2.String(),
			Coins:   sdk.Coins{stratos.NewCoin(idxNodeInitialStake2)},
		},
		{
			Address: idxOwner3.String(),
			Coins:   sdk.Coins{stratos.NewCoin(idxNodeInitialStake3)},
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

func setupAllResourceNodes() []registertypes.ResourceNode {

	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	nodeType := registertypes.STORAGE
	resourceNode1, _ := registertypes.NewResourceNode(resNodeNetworkId1, resNodePubKey1, resOwner1, registertypes.NewDescription("sds://resourceNode1", "", "", "", ""), nodeType, time)
	resourceNode2, _ := registertypes.NewResourceNode(resNodeNetworkId2, resNodePubKey2, resOwner2, registertypes.NewDescription("sds://resourceNode2", "", "", "", ""), nodeType, time)
	resourceNode3, _ := registertypes.NewResourceNode(resNodeNetworkId3, resNodePubKey3, resOwner3, registertypes.NewDescription("sds://resourceNode3", "", "", "", ""), nodeType, time)
	resourceNode4, _ := registertypes.NewResourceNode(resNodeNetworkId4, resNodePubKey4, resOwner4, registertypes.NewDescription("sds://resourceNode4", "", "", "", ""), nodeType, time)
	resourceNode5, _ := registertypes.NewResourceNode(resNodeNetworkId5, resNodePubKey5, resOwner5, registertypes.NewDescription("sds://resourceNode5", "", "", "", ""), nodeType, time)

	resourceNode1 = resourceNode1.AddToken(resNodeInitialStake1)
	resourceNode2 = resourceNode2.AddToken(resNodeInitialStake2)
	resourceNode3 = resourceNode3.AddToken(resNodeInitialStake3)
	resourceNode4 = resourceNode4.AddToken(resNodeInitialStake4)
	resourceNode5 = resourceNode5.AddToken(resNodeInitialStake5)

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

func setupAllMetaNodes() []registertypes.MetaNode {
	var indexingNodes []registertypes.MetaNode

	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	indexingNode1, _ := registertypes.NewMetaNode(stratos.SdsAddress(idxNodeAddr1), idxNodePubKey1, idxOwner1, registertypes.NewDescription("sds://indexingNode1", "", "", "", ""), time)
	indexingNode2, _ := registertypes.NewMetaNode(stratos.SdsAddress(idxNodeAddr2), idxNodePubKey2, idxOwner2, registertypes.NewDescription("sds://indexingNode2", "", "", "", ""), time)
	indexingNode3, _ := registertypes.NewMetaNode(stratos.SdsAddress(idxNodeAddr3), idxNodePubKey3, idxOwner3, registertypes.NewDescription("sds://indexingNode3", "", "", "", ""), time)

	indexingNode1 = indexingNode1.AddToken(idxNodeInitialStake1)
	indexingNode2 = indexingNode2.AddToken(idxNodeInitialStake2)
	indexingNode3 = indexingNode3.AddToken(idxNodeInitialStake3)

	indexingNode1.Status = stakingtypes.Bonded
	indexingNode2.Status = stakingtypes.Bonded
	indexingNode3.Status = stakingtypes.Bonded

	indexingNodes = append(indexingNodes, indexingNode1)
	indexingNodes = append(indexingNodes, indexingNode2)
	indexingNodes = append(indexingNodes, indexingNode3)

	return indexingNodes

}
