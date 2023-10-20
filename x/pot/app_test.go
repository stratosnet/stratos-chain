package pot_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/stratosnet/stratos-chain/app"
	"github.com/stratosnet/stratos-chain/crypto"
	"github.com/stratosnet/stratos-chain/crypto/bls"
	stratostestutil "github.com/stratosnet/stratos-chain/testutil/stratos"
	stratos "github.com/stratosnet/stratos-chain/types"
	potKeeper "github.com/stratosnet/stratos-chain/x/pot/keeper"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	registerKeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
	sdstypes "github.com/stratosnet/stratos-chain/x/sds/types"
)

const (
	chainID     = "testchain_1-1"
	rewardDenom = stratos.Utros

	stopFlagOutOfTotalMiningReward = true
	stopFlagSpecificMinedReward    = false
	stopFlagSpecificEpoch          = true
)

var (
	paramSpecificMinedReward = sdk.NewCoins(stratos.NewCoinInt64(160000000000))
	paramSpecificEpoch       = sdkmath.NewInt(10)

	resNodeSlashingNOZAmt1            = sdkmath.NewInt(1000000000000000000)
	resNodeSlashingEffectiveTokenAmt1 = sdkmath.NewInt(1000000000000000000)

	resourceNodeVolume1 = sdkmath.NewInt(50000)
	resourceNodeVolume2 = sdkmath.NewInt(30000)
	resourceNodeVolume3 = sdkmath.NewInt(20000)

	prepayAmount = sdk.NewCoins(stratos.NewCoin(sdkmath.NewInt(20).Mul(sdkmath.NewInt(stratos.StosToWei))))

	foundationDepositorPrivKey = secp256k1.GenPrivKey()
	foundationDepositorAccAddr = sdk.AccAddress(foundationDepositorPrivKey.PubKey().Address())
	foundationDeposit          = sdk.NewCoins(sdk.NewCoin(rewardDenom, sdkmath.NewInt(40000000000000000)))

	nodeInitialDeposit = sdkmath.NewInt(1 * stratos.StosToWei)
	initBalance        = sdkmath.NewInt(100).Mul(sdkmath.NewInt(stratos.StosToWei))

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

// initialize data of volume report
func setupMsgVolumeReport(t *testing.T, newEpoch int64) *types.MsgVolumeReport {
	volume1 := types.NewSingleWalletVolume(resOwner1, resourceNodeVolume1)
	volume2 := types.NewSingleWalletVolume(resOwner2, resourceNodeVolume2)
	volume3 := types.NewSingleWalletVolume(resOwner3, resourceNodeVolume3)

	nodesVolume := []types.SingleWalletVolume{volume1, volume2, volume3}
	reporter := metaNodeP2PAddr1
	epoch := sdkmath.NewInt(newEpoch)
	reportReference := "report for epoch " + epoch.String()
	reporterOwner := metaOwner1

	signature := types.BLSSignatureInfo{}
	volumeReportMsg := types.NewMsgVolumeReport(nodesVolume, reporter, epoch, reportReference, reporterOwner, signature)

	signBytes := volumeReportMsg.GetSignBytes()
	signBytesHash := crypto.Keccak256(signBytes)

	// set blsSignature
	blsPrivKey1, blsPubKey1, err := bls.NewKeyPairFromBytes(metaNodeP2PPrivKey1.Bytes())
	require.NoError(t, err)
	blsPrivKey2, blsPubKey2, err := bls.NewKeyPairFromBytes(metaNodeP2PPrivKey2.Bytes())
	require.NoError(t, err)
	blsPrivKey3, blsPubKey3, err := bls.NewKeyPairFromBytes(metaNodeP2PPrivKey3.Bytes())
	require.NoError(t, err)

	blsSignature1, err := bls.Sign(signBytesHash, blsPrivKey1)
	require.NoError(t, err)
	blsSignature2, err := bls.Sign(signBytesHash, blsPrivKey2)
	require.NoError(t, err)
	blsSignature3, err := bls.Sign(signBytesHash, blsPrivKey3)
	require.NoError(t, err)
	finalBlsSignature, err := bls.AggregateSignatures(blsSignature1, blsSignature2, blsSignature3)
	require.NoError(t, err)

	pubKeys := make([][]byte, 0)
	pubKeys = append(pubKeys, blsPubKey1, blsPubKey2, blsPubKey3)

	signature = types.NewBLSSignatureInfo(pubKeys, finalBlsSignature, signBytesHash)

	volumeReportMsg.BLSSignature = signature

	return volumeReportMsg
}

func setupSlashingMsg() *types.MsgSlashingResourceNode {
	reporters := make([]stratos.SdsAddress, 0)
	reporters = append(reporters, metaNodeP2PAddr1)
	reportOwner := make([]sdk.AccAddress, 0)
	reportOwner = append(reportOwner, metaOwner1)
	slashingMsg := types.NewMsgSlashingResourceNode(reporters, reportOwner, resNodeP2PAddr1, resOwner1, resNodeSlashingNOZAmt1, true)
	return slashingMsg
}

// Test case termination conditions
// modify stop flag & variable could make the test case stop when reach a specific condition
func isNeedStop(ctx sdk.Context, k potKeeper.Keeper, epoch sdkmath.Int, minedToken sdk.Coin) bool {

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

	stApp := app.SetupWithGenesisNodeSet(t, false, valSet, metaNodes, resourceNodes, accs, chainID, balances...)
	accountKeeper := stApp.GetAccountKeeper()
	bankKeeper := stApp.GetBankKeeper()
	registerKeeper := stApp.GetRegisterKeeper()
	potKeeper := stApp.GetPotKeeper()
	distrKeeper := stApp.GetDistrKeeper()

	/********************* foundation account deposit *********************/
	header := tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := stApp.BaseApp.NewContext(true, header)

	foundationDepositMsg := types.NewMsgFoundationDeposit(foundationDeposit, foundationDepositorAccAddr)
	txGen := app.MakeTestEncodingConfig().TxConfig

	senderAcc := accountKeeper.GetAccount(ctx, foundationDepositorAccAddr)
	accNum := senderAcc.GetAccountNumber()
	accSeq := senderAcc.GetSequence()
	_, _, err := stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, foundationDepositorPrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(types.FoundationAccount)
	app.CheckBalance(t, stApp, foundationAccountAddr, foundationDeposit)

	/********************* create validator with 50% commission *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(false, header)

	commission := stakingtypes.NewCommissionRates(
		sdkmath.LegacyNewDecWithPrec(5, 1),
		sdkmath.LegacyNewDecWithPrec(5, 1),
		sdkmath.LegacyNewDec(0),
	)
	description := stakingtypes.NewDescription("foo_moniker", chainID, "", "", "")
	createValidatorMsg, err := stakingtypes.NewMsgCreateValidator(valOpValAddr1, valConsPubk1, stratos.NewCoin(nodeInitialDeposit), description, commission, sdkmath.OneInt())

	senderAcc = accountKeeper.GetAccount(ctx, valOpAccAddr1)
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()
	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{createValidatorMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, valOpPrivKey1)
	require.NoError(t, err)

	/********************* prepay *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(false, header)
	prepayMsg := sdstypes.NewMsgPrepay(resOwner1.String(), resOwner1.String(), prepayAmount)
	senderAcc = accountKeeper.GetAccount(ctx, resOwner1)
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()
	_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{prepayMsg}, chainID, []uint64{accNum}, []uint64{accSeq}, true, true, resOwnerPrivKey1)
	require.NoError(t, err)

	/********************** commit **********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	validator := checkValidator(t, stApp, valOpValAddr1, true)
	require.Equal(t, stakingtypes.Bonded, validator.Status)
	require.True(sdkmath.IntEq(t, nodeInitialDeposit, validator.BondedTokens()))

	/********************** loop sending volume report **********************/
	var i int64
	var slashingAmtSetup sdkmath.Int
	i = 0
	slashingAmtSetup = sdkmath.ZeroInt()
	for {

		/********************* test slashing msg when i==2 *********************/
		if i == 2 {
			t.Log("********************************* Deliver Slashing Tx START ********************************************")

			totalConsumedNoz := resNodeSlashingNOZAmt1.ToLegacyDec()
			slashingAmtCheck := potKeeper.GetTrafficReward(ctx, totalConsumedNoz)

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

			slashingAmtSetup = registerKeeper.GetSlashing(ctx, resOwner1)

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
		volumeReportMsg := setupMsgVolumeReport(t, i+1)

		lastTotalMinedToken := potKeeper.GetTotalMinedTokens(ctx)
		t.Log("last committed TotalMinedTokens = " + lastTotalMinedToken.String())
		epoch, ok := sdkmath.NewIntFromString(volumeReportMsg.Epoch.String())
		require.Equal(t, ok, true)

		if isNeedStop(ctx, potKeeper, epoch, lastTotalMinedToken) {
			break
		}

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

		t.Log("indexing_wallet1:  address = " + metaOwner1.String())
		t.Log("              miningReward = " + rewardDetailMap[metaOwner1.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[metaOwner1.String()].RewardFromTrafficPool.String())

		t.Log("indexing_wallet2:  address = " + metaOwner2.String())
		t.Log("              miningReward = " + rewardDetailMap[metaOwner2.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[metaOwner2.String()].RewardFromTrafficPool.String())

		t.Log("indexing_wallet3:  address = " + metaOwner3.String())
		t.Log("              miningReward = " + rewardDetailMap[metaOwner3.String()].RewardFromMiningPool.String())
		t.Log("             trafficReward = " + rewardDetailMap[metaOwner3.String()].RewardFromTrafficPool.String())
		t.Log("---------------------------")

		/********************* record data before delivering tx  *********************/
		lastFoundationAccBalance := bankKeeper.GetAllBalances(ctx, foundationAccountAddr)
		lastUnissuedPrepay := registerKeeper.GetTotalUnissuedPrepay(ctx)
		lastCommunityPool := sdk.NewCoins(sdk.NewCoin(potKeeper.BondDenom(ctx), distrKeeper.GetFeePool(ctx).CommunityPool.AmountOf(potKeeper.BondDenom(ctx)).TruncateInt()))
		lastMatureTotalOfResNode1 := potKeeper.GetMatureTotalReward(ctx, resOwner1)
		lastIndividualRewardOfResNode1, individualRewardOfResNode1Found := potKeeper.GetIndividualReward(ctx, resOwner1, epoch)
		/********************* deliver tx *********************/
		idxOwnerAcc1 := accountKeeper.GetAccount(ctx, metaOwner1)
		ownerAccNum := idxOwnerAcc1.GetAccountNumber()
		ownerAccSeq := idxOwnerAcc1.GetSequence()

		feePoolAccAddr := accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		require.NotNil(t, feePoolAccAddr)

		t.Log("--------------------------- deliver volumeReportMsg")
		_, _, err = stratostestutil.SignCheckDeliver(t, txGen, stApp.BaseApp, header, []sdk.Msg{volumeReportMsg}, chainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, metaOwnerPrivKey1)
		require.NoError(t, err)

		/********************* commit & check result *********************/
		feeCollectorToFeePoolAtBeginBlock := bankKeeper.GetAllBalances(ctx, feePoolAccAddr)

		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		stApp.EndBlock(abci.RequestEndBlock{Height: header.Height})
		stApp.Commit()

		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)
		stApp.Commit()
		header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: chainID}
		stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx = stApp.BaseApp.NewContext(true, header)

		epoch, ok = sdkmath.NewIntFromString(volumeReportMsg.Epoch.String())
		require.Equal(t, ok, true)

		checkResult(t, ctx, potKeeper,
			accountKeeper,
			bankKeeper,
			registerKeeper,
			distrKeeper,
			epoch,
			lastFoundationAccBalance,
			lastUnissuedPrepay,
			lastCommunityPool,
			lastMatureTotalOfResNode1,
			slashingAmtSetup,
			feeCollectorToFeePoolAtBeginBlock,
			lastIndividualRewardOfResNode1,
			individualRewardOfResNode1Found,
		)

		i++
	}
}

// return : coins - slashing
func deductSlashingAmt(ctx sdk.Context, coins sdk.Coins, slashing sdk.Coin) (ret sdk.Coins) {
	slashingDenom := slashing.Denom
	rewardToken := sdk.NewCoin(slashingDenom, coins.AmountOf(slashingDenom))
	if rewardToken.IsGTE(slashing) {
		ret = coins.Sub(sdk.NewCoins(slashing)...)
	} else {
		ret = coins.Sub(sdk.NewCoins(rewardToken)...)
	}
	return ret
}

// for main net
func checkResult(t *testing.T, ctx sdk.Context,
	k potKeeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankKeeper.Keeper,
	registerKeeper registerKeeper.Keeper,
	distrKeeper distrkeeper.Keeper,
	currentEpoch sdkmath.Int,
	lastFoundationAccBalance sdk.Coins,
	lastUnissuedPrepay sdk.Coin,
	lastCommunityPool sdk.Coins,
	lastMatureTotalOfResNode1 sdk.Coins,
	initialSlashingAmt sdkmath.Int,
	feeCollectorToFeePoolAtBeginBlock sdk.Coins,
	individualRewardOfResNode1 types.Reward,
	individualRewardOfResNode1Found bool,
) {

	// print individual reward
	individualRewardTotal := sdk.Coins{}
	newMatureEpoch := currentEpoch.Add(sdkmath.NewInt(k.MatureEpoch(ctx)))
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
	newCommunityPool := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), distrKeeper.GetFeePool(ctx).CommunityPool.AmountOf(k.BondDenom(ctx)).TruncateInt()))

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
	t.Log("reward for validator send to fee_collector = " + feeCollectorToFeePoolAtBeginBlock.String())
	stakeRewardFromFeeCollectorToFeePool := sdk.NewCoin(k.BondDenom(ctx), feeCollectorToFeePoolAtBeginBlock.AmountOf(k.BondDenom(ctx)))
	communityTaxChange := newCommunityPool.Sub(lastCommunityPool...).Sub(sdk.NewCoins(stakeRewardFromFeeCollectorToFeePool)...)
	t.Log("community tax change in community_pool     = " + communityTaxChange.String())
	t.Log("community_pool amount of wei               = " + newCommunityPool.String())

	rewardSrcChange := lastFoundationAccBalance.
		Sub(newFoundationAccBalance...).
		Add(lastUnissuedPrepay).
		Sub(newUnissuedPrepay...)

	t.Log("rewardSrcChange                            = " + rewardSrcChange.String())

	rewardDestChange := feeCollectorToFeePoolAtBeginBlock.
		Add(individualRewardTotal...).
		Add(communityTaxChange...)

	t.Log("rewardDestChange                           = " + rewardDestChange.String())

	require.Equal(t, rewardSrcChange, rewardDestChange)

	t.Log("************************ slashing test***********************************")
	t.Log("slashing change                            = " + slashingDeducted.String())

	upcomingMaturedIndividual := sdk.Coins{}

	if individualRewardOfResNode1Found {
		tmp := individualRewardOfResNode1.RewardFromTrafficPool.Add(individualRewardOfResNode1.RewardFromMiningPool...)
		upcomingMaturedIndividual = deductSlashingAmt(ctx, tmp, slashingDeducted)
	}
	t.Log("upcomingMaturedIndividual                  = " + upcomingMaturedIndividual.String())

	// get mature total changes
	newMatureTotalOfResNode1 := k.GetMatureTotalReward(ctx, resOwner1)
	matureTotalOfResNode1Change, _ := newMatureTotalOfResNode1.SafeSub(lastMatureTotalOfResNode1...)
	if matureTotalOfResNode1Change == nil || matureTotalOfResNode1Change.IsAnyNegative() {
		matureTotalOfResNode1Change = sdk.Coins{}
	}
	t.Log("matureTotalOfResNode1Change                = " + matureTotalOfResNode1Change.String())
	require.Equal(t, matureTotalOfResNode1Change.String(), upcomingMaturedIndividual.String())

	totalRewardPoolAddr := accountKeeper.GetModuleAddress(types.TotalRewardPool)
	totalRewardPoolBalance := bankKeeper.GetAllBalances(ctx, totalRewardPoolAddr)
	t.Log("totalRewardPoolBalance                     = " + totalRewardPoolBalance.String())
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
	//************************** setup indexing nodes owners' accounts **************************
	idxOwnerAcc1 := &authtypes.BaseAccount{Address: metaOwner1.String()}
	idxOwnerAcc2 := &authtypes.BaseAccount{Address: metaOwner2.String()}
	idxOwnerAcc3 := &authtypes.BaseAccount{Address: metaOwner3.String()}
	//************************** setup validator delegators' accounts **************************
	valOwnerAcc1 := &authtypes.BaseAccount{Address: valOpAccAddr1.String()}
	////************************** setup indexing nodes' accounts **************************
	//idxNodeAcc1 := &authtypes.BaseAccount{Address: idxNodeAddr1.String()}
	foundationDepositorAcc := &authtypes.BaseAccount{Address: foundationDepositorAccAddr.String()}

	accs := []authtypes.GenesisAccount{
		resOwnerAcc1, resOwnerAcc2, resOwnerAcc3, resOwnerAcc4, resOwnerAcc5,
		idxOwnerAcc1, idxOwnerAcc2, idxOwnerAcc3,
		valOwnerAcc1,
		foundationDepositorAcc,
		//idxNodeAcc1,
	}

	feeAmt, _ := sdkmath.NewIntFromString("50000000000000000000")

	balances := []banktypes.Balance{
		{
			Address: resOwner1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(initBalance)},
		},
		{
			Address: resOwner2.String(),
			Coins:   sdk.Coins{stratos.NewCoin(initBalance)},
		},
		{
			Address: resOwner3.String(),
			Coins:   sdk.Coins{stratos.NewCoin(initBalance)},
		},
		{
			Address: resOwner4.String(),
			Coins:   sdk.Coins{stratos.NewCoin(initBalance)},
		},
		{
			Address: resOwner5.String(),
			Coins:   sdk.Coins{stratos.NewCoin(initBalance)},
		},
		{
			Address: metaOwner1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(initBalance)},
		},
		{
			Address: metaOwner2.String(),
			Coins:   sdk.Coins{stratos.NewCoin(initBalance)},
		},
		{
			Address: metaOwner3.String(),
			Coins:   sdk.Coins{stratos.NewCoin(initBalance)},
		},
		{
			Address: valOpAccAddr1.String(),
			Coins:   sdk.Coins{stratos.NewCoin(initBalance)},
		},
		//{
		//	Address: idxNodeAddr1.String(),
		//	Coins:   sdk.Coins{stratos.NewCoin(sdkmath.ZeroInt())},
		//},
		{
			Address: foundationDepositorAccAddr.String(),
			Coins:   foundationDeposit.Add(sdk.NewCoin(stratos.Wei, feeAmt)),
		},
	}
	return accs, balances
}

func setupAllResourceNodes() []registertypes.ResourceNode {

	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	nodeType := registertypes.STORAGE
	resourceNode1, _ := registertypes.NewResourceNode(resNodeP2PAddr1, resNodeP2PPubKey1, resOwner1, registertypes.NewDescription("resourceNode1", "", "", "", ""), nodeType, time)
	resourceNode2, _ := registertypes.NewResourceNode(resNodeP2PAddr2, resNodeP2PPubKey2, resOwner2, registertypes.NewDescription("resourceNode2", "", "", "", ""), nodeType, time)
	resourceNode3, _ := registertypes.NewResourceNode(resNodeP2PAddr3, resNodeP2PPubKey3, resOwner3, registertypes.NewDescription("resourceNode3", "", "", "", ""), nodeType, time)
	resourceNode4, _ := registertypes.NewResourceNode(resNodeP2PAddr4, resNodeP2PPubKey4, resOwner4, registertypes.NewDescription("resourceNode4", "", "", "", ""), nodeType, time)
	resourceNode5, _ := registertypes.NewResourceNode(resNodeP2PAddr5, resNodeP2PPubKey5, resOwner5, registertypes.NewDescription("resourceNode5", "", "", "", ""), nodeType, time)

	resourceNode1 = resourceNode1.AddToken(nodeInitialDeposit)
	resourceNode2 = resourceNode2.AddToken(nodeInitialDeposit)
	resourceNode3 = resourceNode3.AddToken(nodeInitialDeposit)
	resourceNode4 = resourceNode4.AddToken(nodeInitialDeposit)
	resourceNode5 = resourceNode5.AddToken(nodeInitialDeposit)

	resourceNode1.EffectiveTokens = nodeInitialDeposit
	resourceNode2.EffectiveTokens = nodeInitialDeposit
	resourceNode3.EffectiveTokens = nodeInitialDeposit
	resourceNode4.EffectiveTokens = nodeInitialDeposit
	resourceNode5.EffectiveTokens = nodeInitialDeposit

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

func setupAllMetaNodes() []registertypes.MetaNode {
	var indexingNodes []registertypes.MetaNode

	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	indexingNode1, _ := registertypes.NewMetaNode(metaNodeP2PAddr1, metaNodeP2PPubKey1, metaOwner1, registertypes.NewDescription("indexingNode1", "", "", "", ""), time)
	indexingNode2, _ := registertypes.NewMetaNode(metaNodeP2PAddr2, metaNodeP2PPubKey2, metaOwner2, registertypes.NewDescription("indexingNode2", "", "", "", ""), time)
	indexingNode3, _ := registertypes.NewMetaNode(metaNodeP2PAddr3, metaNodeP2PPubKey3, metaOwner3, registertypes.NewDescription("indexingNode3", "", "", "", ""), time)

	indexingNode1 = indexingNode1.AddToken(nodeInitialDeposit)
	indexingNode2 = indexingNode2.AddToken(nodeInitialDeposit)
	indexingNode3 = indexingNode3.AddToken(nodeInitialDeposit)

	indexingNode1.Status = stakingtypes.Bonded
	indexingNode2.Status = stakingtypes.Bonded
	indexingNode3.Status = stakingtypes.Bonded

	indexingNode1.Suspend = false
	indexingNode2.Suspend = false
	indexingNode3.Suspend = false

	indexingNodes = append(indexingNodes, indexingNode1)
	indexingNodes = append(indexingNodes, indexingNode2)
	indexingNodes = append(indexingNodes, indexingNode3)

	return indexingNodes

}
