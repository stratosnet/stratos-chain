package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	stos2ustos = 1000000000
	oz2uoz     = 1000000000

	resourceNodeVolume1 = 500 * oz2uoz
	resourceNodeVolume2 = 300 * oz2uoz
	resourceNodeVolume3 = 200 * oz2uoz
	totalVolume         = resourceNodeVolume1 + resourceNodeVolume2 + resourceNodeVolume3
)

var (
	foundationDeposit = sdk.NewCoins(sdk.NewCoin("ustos", sdk.NewInt(40000000*stos2ustos)))

	resOwner1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	resOwner2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	resOwner3 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	resOwner4 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	resOwner5 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	idxOwner1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	idxOwner2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	idxOwner3 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	pubKeyRes1       = ed25519.GenPrivKey().PubKey()
	addrRes1         = sdk.AccAddress(pubKeyRes1.Address())
	initialStakeRes1 = sdk.NewCoin("ustos", sdk.NewInt(3*stos2ustos))

	pubKeyRes2       = ed25519.GenPrivKey().PubKey()
	addrRes2         = sdk.AccAddress(pubKeyRes2.Address())
	initialStakeRes2 = sdk.NewCoin("ustos", sdk.NewInt(3*stos2ustos))

	pubKeyRes3       = ed25519.GenPrivKey().PubKey()
	addrRes3         = sdk.AccAddress(pubKeyRes3.Address())
	initialStakeRes3 = sdk.NewCoin("ustos", sdk.NewInt(3*stos2ustos))

	pubKeyRes4       = ed25519.GenPrivKey().PubKey()
	addrRes4         = sdk.AccAddress(pubKeyRes4.Address())
	initialStakeRes4 = sdk.NewCoin("ustos", sdk.NewInt(3*stos2ustos))

	pubKeyRes5       = ed25519.GenPrivKey().PubKey()
	addrRes5         = sdk.AccAddress(pubKeyRes5.Address())
	initialStakeRes5 = sdk.NewCoin("ustos", sdk.NewInt(3*stos2ustos))

	pubKeyIdx1       = ed25519.GenPrivKey().PubKey()
	addrIdx1         = sdk.AccAddress(pubKeyIdx1.Address())
	initialStakeIdx1 = sdk.NewCoin("ustos", sdk.NewInt(5*stos2ustos))

	pubKeyIdx2       = ed25519.GenPrivKey().PubKey()
	addrIdx2         = sdk.AccAddress(pubKeyIdx2.Address())
	initialStakeIdx2 = sdk.NewCoin("ustos", sdk.NewInt(5*stos2ustos))

	pubKeyIdx3       = ed25519.GenPrivKey().PubKey()
	addrIdx3         = sdk.AccAddress(pubKeyIdx3.Address())
	initialStakeIdx3 = sdk.NewCoin("ustos", sdk.NewInt(5*stos2ustos))

	valOpPk1        = ed25519.GenPrivKey().PubKey()
	valOpAddr1      = sdk.ValAddress(valOpPk1.Address())
	valAccAddr1     = sdk.AccAddress(valOpPk1.Address())
	valConsPk1      = ed25519.GenPrivKey().PubKey()
	valInitialStake = sdk.NewCoin("ustos", sdk.NewInt(15*stos2ustos))

	totalUnissuedPrePay = sdk.NewInt(5000 * stos2ustos)
	remainingOzoneLimit = sdk.NewInt(5000 * oz2uoz)
	//initialUOzonePrice = sdk.NewDec(1000000)
	initialUOzonePrice = sdk.NewDecWithPrec(1000000, 9) // 0.001 ustos -> 1 uoz

	epoch1    = sdk.NewInt(1)
	epoch2017 = epoch1.Add(sdk.NewInt(2016))
	epoch4033 = epoch2017.Add(sdk.NewInt(2016))
)

func Test(t *testing.T) {

	//prepare keepers
	ctx, accountKeeper, bankKeeper, k, stakingKeeper, _, supplyKeeper, registerKeeper := CreateTestInput(t, false)

	// create validator with 50% commission
	stakingHandler := staking.NewHandler(stakingKeeper)
	createAccount(t, ctx, accountKeeper, bankKeeper, valAccAddr1, sdk.NewCoins(valInitialStake))
	commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	msgVal := staking.NewMsgCreateValidator(valOpAddr1, valConsPk1, valInitialStake, staking.Description{}, commission, sdk.OneInt())
	res, err := stakingHandler(ctx, msgVal)
	require.NoError(t, err)
	require.NotNil(t, res)
	stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)

	/************************************************** pot reward distribution part start ************************************************************/
	//initial genesis stake total value
	initialGenesisStakeTotal := initialStakeRes1.Add(initialStakeRes2).Add(initialStakeRes3).Add(initialStakeRes4).Add(initialStakeRes5).
		Add(initialStakeIdx1).Add(initialStakeIdx2).Add(initialStakeIdx3)
	registerKeeper.SetInitialGenesisStakeTotal(ctx, initialGenesisStakeTotal.Amount)

	//PrePay
	k.SetTotalUnissuedPrepay(ctx, totalUnissuedPrePay)
	//remaining ozone limit
	registerKeeper.SetRemainingOzoneLimit(ctx, remainingOzoneLimit)
	//initial uoz price
	registerKeeper.SetInitialUOzonePrice(ctx, initialUOzonePrice)

	//pot genesis data load
	foundationAccountAddr := supplyKeeper.GetModuleAddress(types.FoundationAccount)
	err = bankKeeper.SetCoins(ctx, foundationAccountAddr, foundationDeposit)
	require.NoError(t, err)

	//initialize owner accounts
	createAccount(t, ctx, accountKeeper, bankKeeper, resOwner1, sdk.NewCoins(initialStakeRes1))
	createAccount(t, ctx, accountKeeper, bankKeeper, resOwner2, sdk.NewCoins(initialStakeRes2))
	createAccount(t, ctx, accountKeeper, bankKeeper, resOwner3, sdk.NewCoins(initialStakeRes3))
	createAccount(t, ctx, accountKeeper, bankKeeper, resOwner4, sdk.NewCoins(initialStakeRes4))
	createAccount(t, ctx, accountKeeper, bankKeeper, resOwner5, sdk.NewCoins(initialStakeRes5))
	createAccount(t, ctx, accountKeeper, bankKeeper, idxOwner1, sdk.NewCoins(initialStakeIdx1))
	createAccount(t, ctx, accountKeeper, bankKeeper, idxOwner2, sdk.NewCoins(initialStakeIdx2))
	createAccount(t, ctx, accountKeeper, bankKeeper, idxOwner3, sdk.NewCoins(initialStakeIdx3))
	//initialize sds node register msg
	msgRes1 := register.NewMsgCreateResourceNode("sds://resourceNode1", pubKeyRes1, initialStakeRes1, resOwner1, register.NewDescription("sds://resourceNode1", "", "", "", ""), "4")
	msgRes2 := register.NewMsgCreateResourceNode("sds://resourceNode2", pubKeyRes2, initialStakeRes2, resOwner2, register.NewDescription("sds://resourceNode2", "", "", "", ""), "4")
	msgRes3 := register.NewMsgCreateResourceNode("sds://resourceNode3", pubKeyRes3, initialStakeRes3, resOwner3, register.NewDescription("sds://resourceNode3", "", "", "", ""), "4")
	msgRes4 := register.NewMsgCreateResourceNode("sds://resourceNode4", pubKeyRes4, initialStakeRes4, resOwner4, register.NewDescription("sds://resourceNode4", "", "", "", ""), "4")
	msgRes5 := register.NewMsgCreateResourceNode("sds://resourceNode5", pubKeyRes5, initialStakeRes5, resOwner5, register.NewDescription("sds://resourceNode5", "", "", "", ""), "4")
	msgIdx1 := register.NewMsgCreateIndexingNode("sds://indexingNode1", pubKeyIdx1, initialStakeIdx1, idxOwner1, register.NewDescription("sds://indexingNode1", "", "", "", ""))
	msgIdx2 := register.NewMsgCreateIndexingNode("sds://indexingNode2", pubKeyIdx2, initialStakeIdx2, idxOwner2, register.NewDescription("sds://indexingNode2", "", "", "", ""))
	msgIdx3 := register.NewMsgCreateIndexingNode("sds://indexingNode3", pubKeyIdx3, initialStakeIdx3, idxOwner3, register.NewDescription("sds://indexingNode3", "", "", "", ""))

	//register sds nodes
	registerHandler := register.NewHandler(registerKeeper)
	res, err = registerHandler(ctx, msgRes1)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = registerHandler(ctx, msgRes2)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = registerHandler(ctx, msgRes3)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = registerHandler(ctx, msgRes4)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = registerHandler(ctx, msgRes5)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = registerHandler(ctx, msgIdx1)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = registerHandler(ctx, msgIdx2)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = registerHandler(ctx, msgIdx3)
	require.NoError(t, err)
	require.NotNil(t, res)

	// set the status of indexing nodes to bonded
	idxUnBondedPool := k.RegisterKeeper.GetIndexingNodeNotBondedToken(ctx)
	k.RegisterKeeper.SetIndexingNodeBondedToken(ctx, idxUnBondedPool)
	k.RegisterKeeper.SetIndexingNodeNotBondedToken(ctx, sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt()))
	idxNode1, _ := k.RegisterKeeper.GetIndexingNode(ctx, addrIdx1)
	idxNode2, _ := k.RegisterKeeper.GetIndexingNode(ctx, addrIdx2)
	idxNode3, _ := k.RegisterKeeper.GetIndexingNode(ctx, addrIdx3)
	idxNode1.Status = sdk.Bonded
	idxNode2.Status = sdk.Bonded
	idxNode3.Status = sdk.Bonded
	k.RegisterKeeper.SetIndexingNode(ctx, idxNode1)
	k.RegisterKeeper.SetIndexingNode(ctx, idxNode2)
	k.RegisterKeeper.SetIndexingNode(ctx, idxNode3)

	//build traffic list
	var trafficList []types.SingleWalletVolume
	trafficList = append(trafficList, types.NewSingleWalletVolume(resOwner1, sdk.NewInt(resourceNodeVolume1)))
	trafficList = append(trafficList, types.NewSingleWalletVolume(resOwner2, sdk.NewInt(resourceNodeVolume2)))
	trafficList = append(trafficList, types.NewSingleWalletVolume(resOwner3, sdk.NewInt(resourceNodeVolume3)))

	//check prepared data
	S := k.RegisterKeeper.GetInitialGenesisStakeTotal(ctx).ToDec()
	fmt.Println("S=" + S.String())
	Pt := k.GetTotalUnissuedPrepay(ctx).ToDec()
	fmt.Println("Pt=" + Pt.String())
	Y := k.GetTotalConsumedUoz(trafficList).ToDec()
	fmt.Println("Y=" + Y.String())
	Lt := k.RegisterKeeper.GetRemainingOzoneLimit(ctx).ToDec()
	fmt.Println("Lt=" + Lt.String())
	R := S.Add(Pt).Mul(Y).Quo(Lt.Add(Y))
	fmt.Println("R=" + R.String())

	fmt.Println("***************************************************************************************")

	testBlockChainRewardFromTrafficPool(t, ctx, k, bankKeeper, trafficList)
	testMetaNodeRewardFromTrafficPool(t, ctx, k, bankKeeper, trafficList)
	testTrafficRewardFromTrafficPool(t, ctx, k, bankKeeper, trafficList)

	testBlockChainRewardFromMiningPool(t, ctx, k, bankKeeper, trafficList)
	testMetaNodeRewardFromMiningPool(t, ctx, k, bankKeeper, trafficList)
	testTrafficRewardFromMiningPool(t, ctx, k, bankKeeper, trafficList)

	testFullDistributeProcessAtEpoch1(t, ctx, k, trafficList)
	testFullDistributeProcessAtEpoch2017(t, ctx, k, trafficList)
	testWithdraw(t, ctx, k, bankKeeper)

}

func testWithdraw(t *testing.T, ctx sdk.Context, k Keeper, bankKeeper bank.Keeper) {
	AccountBalanceBefore := bankKeeper.GetCoins(ctx, resOwner1)

	err := k.Withdraw(ctx, sdk.NewCoin("ustos", sdk.NewInt(68846296294)), resOwner2, resOwner2)
	require.Error(t, err, types.ErrNotTheOwner)

	err = k.Withdraw(ctx, sdk.NewCoin("ustos", sdk.NewInt(68846296295)), resOwner1, resOwner1)
	require.Error(t, err, types.ErrInsufficientMatureTotal)

	err = k.Withdraw(ctx, sdk.NewCoin("ustos", sdk.NewInt(68846296294)), resOwner1, resOwner1)
	require.NoError(t, err)

	AccountBalanceAfter := bankKeeper.GetCoins(ctx, resOwner1)
	require.Equal(t, AccountBalanceAfter.Sub(AccountBalanceBefore).AmountOf("ustos"), sdk.NewInt(68846296294))

	matureTotalResNode1 := k.GetMatureTotalReward(ctx, resOwner1)
	require.Equal(t, matureTotalResNode1, sdk.ZeroInt())
}

func testFullDistributeProcessAtEpoch2017(t *testing.T, ctx sdk.Context, k Keeper, trafficList []types.SingleWalletVolume) {
	_, err := k.DistributePotReward(ctx, trafficList, epoch2017)
	require.NoError(t, err)
	fmt.Println("Distribution result at Epoch2017: ")
	rewardAddrList := k.GetRewardAddressPool(ctx)
	fmt.Println("address pool: ")
	for i := 0; i < len(rewardAddrList); i++ {
		fmt.Println(rewardAddrList[i].String() + ", ")
	}
	fmt.Println("----------------------------------------------------------------------------------")

	idvRwdResNode1Ep1, _ := k.GetIndividualReward(ctx, resOwner1, epoch4033)
	individualTotalReward := idvRwdResNode1Ep1.RewardFromMiningPool.Add(idvRwdResNode1Ep1.RewardFromTrafficPool)
	matureTotalResNode1 := k.GetMatureTotalReward(ctx, resOwner1)
	immatureTotalResNode1 := k.GetImmatureTotalReward(ctx, resOwner1)
	fmt.Println("resource_wallet1: address = " + resOwner1.String() + ", individual = " + individualTotalReward.String() + ",\tmatureTotal = " + matureTotalResNode1.String() + ",\timmatureTotal = " + immatureTotalResNode1.String())
	require.Equal(t, individualTotalReward, sdk.NewInt(67630195471))
	require.Equal(t, matureTotalResNode1, sdk.NewInt(68846296294))
	require.Equal(t, immatureTotalResNode1, sdk.NewInt(67630195471))

	idvRwdResNode2Ep1, _ := k.GetIndividualReward(ctx, resOwner2, epoch4033)
	individualTotalReward = idvRwdResNode2Ep1.RewardFromMiningPool.Add(idvRwdResNode2Ep1.RewardFromTrafficPool)
	matureTotalResNode2 := k.GetMatureTotalReward(ctx, resOwner2)
	immatureTotalResNode2 := k.GetImmatureTotalReward(ctx, resOwner2)
	require.Equal(t, individualTotalReward, sdk.NewInt(41729269545))
	require.Equal(t, matureTotalResNode2, sdk.NewInt(42479629627))
	require.Equal(t, immatureTotalResNode2, sdk.NewInt(41729269545))
	fmt.Println("resource_wallet2: address = " + resOwner2.String() + ", individual = " + individualTotalReward.String() + ",\tmatureTotal = " + matureTotalResNode2.String() + ",\timmatureTotal = " + immatureTotalResNode2.String())

	idvRwdResNode3Ep1, _ := k.GetIndividualReward(ctx, resOwner3, epoch4033)
	individualTotalReward = idvRwdResNode3Ep1.RewardFromMiningPool.Add(idvRwdResNode3Ep1.RewardFromTrafficPool)
	matureTotalResNode3 := k.GetMatureTotalReward(ctx, resOwner3)
	immatureTotalResNode3 := k.GetImmatureTotalReward(ctx, resOwner3)
	require.Equal(t, individualTotalReward, sdk.NewInt(28778806582))
	require.Equal(t, matureTotalResNode3, sdk.NewInt(29296296294))
	require.Equal(t, immatureTotalResNode3, sdk.NewInt(28778806582))
	fmt.Println("resource_wallet3: address = " + resOwner3.String() + ", individual = " + individualTotalReward.String() + ",\tmatureTotal = " + matureTotalResNode3.String() + ",\timmatureTotal = " + immatureTotalResNode3.String())

	idvRwdResNode4Ep1, _ := k.GetIndividualReward(ctx, resOwner4, epoch4033)
	individualTotalReward = idvRwdResNode4Ep1.RewardFromMiningPool.Add(idvRwdResNode4Ep1.RewardFromTrafficPool)
	matureTotalResNode4 := k.GetMatureTotalReward(ctx, resOwner4)
	immatureTotalResNode4 := k.GetImmatureTotalReward(ctx, resOwner4)
	require.Equal(t, individualTotalReward, sdk.NewInt(2877880657))
	require.Equal(t, matureTotalResNode4, sdk.NewInt(2929629628))
	require.Equal(t, immatureTotalResNode4, sdk.NewInt(2877880657))
	fmt.Println("resource_wallet4: address = " + resOwner4.String() + ", individual = " + individualTotalReward.String() + ", \tmatureTotal = " + matureTotalResNode4.String() + ",\timmatureTotal = " + immatureTotalResNode4.String())

	idvRwdResNode5Ep1, _ := k.GetIndividualReward(ctx, resOwner5, epoch4033)
	individualTotalReward = idvRwdResNode5Ep1.RewardFromMiningPool.Add(idvRwdResNode5Ep1.RewardFromTrafficPool)
	matureTotalResNode5 := k.GetMatureTotalReward(ctx, resOwner5)
	immatureTotalResNode5 := k.GetImmatureTotalReward(ctx, resOwner5)
	require.Equal(t, individualTotalReward, sdk.NewInt(2877880657))
	require.Equal(t, matureTotalResNode5, sdk.NewInt(2929629628))
	require.Equal(t, immatureTotalResNode5, sdk.NewInt(2877880657))
	fmt.Println("resource_wallet5: address = " + resOwner5.String() + ", individual = " + individualTotalReward.String() + ", \tmatureTotal = " + matureTotalResNode5.String() + ",\timmatureTotal = " + immatureTotalResNode5.String())

	idvRwdIdxNode1Ep1, _ := k.GetIndividualReward(ctx, idxOwner1, epoch4033)
	individualTotalReward = idvRwdIdxNode1Ep1.RewardFromMiningPool.Add(idvRwdIdxNode1Ep1.RewardFromTrafficPool)
	matureTotalIdxNode1 := k.GetMatureTotalReward(ctx, idxOwner1)
	immatureTotalIdxNode1 := k.GetImmatureTotalReward(ctx, idxOwner1)
	require.Equal(t, individualTotalReward, sdk.NewInt(19185871053))
	require.Equal(t, matureTotalIdxNode1, sdk.NewInt(19530864195))
	require.Equal(t, immatureTotalIdxNode1, sdk.NewInt(19185871053))
	fmt.Println("indexing_wallet1: address = " + idxOwner1.String() + ", individual = " + individualTotalReward.String() + ",\tmatureTotal = " + matureTotalIdxNode1.String() + ",\timmatureTotal = " + immatureTotalIdxNode1.String())

	idvRwdIdxNode2Ep1, _ := k.GetIndividualReward(ctx, idxOwner2, epoch4033)
	individualTotalReward = idvRwdIdxNode2Ep1.RewardFromMiningPool.Add(idvRwdIdxNode2Ep1.RewardFromTrafficPool)
	matureTotalIdxNode2 := k.GetMatureTotalReward(ctx, idxOwner2)
	immatureTotalIdxNode2 := k.GetImmatureTotalReward(ctx, idxOwner2)
	require.Equal(t, individualTotalReward, sdk.NewInt(19185871053))
	require.Equal(t, matureTotalIdxNode2, sdk.NewInt(19530864195))
	require.Equal(t, immatureTotalIdxNode2, sdk.NewInt(19185871053))
	fmt.Println("indexing_wallet2: address = " + idxOwner2.String() + ", individual = " + individualTotalReward.String() + ",\tmatureTotal = " + matureTotalIdxNode2.String() + ",\timmatureTotal = " + immatureTotalIdxNode2.String())

	idvRwdIdxNode3Ep1, _ := k.GetIndividualReward(ctx, idxOwner3, epoch4033)
	individualTotalReward = idvRwdIdxNode3Ep1.RewardFromMiningPool.Add(idvRwdIdxNode3Ep1.RewardFromTrafficPool)
	matureTotalIdxNode3 := k.GetMatureTotalReward(ctx, idxOwner3)
	immatureTotalIdxNode3 := k.GetImmatureTotalReward(ctx, idxOwner3)
	require.Equal(t, individualTotalReward, sdk.NewInt(19185871053))
	require.Equal(t, matureTotalIdxNode3, sdk.NewInt(19530864195))
	require.Equal(t, immatureTotalIdxNode3, sdk.NewInt(19185871053))
	fmt.Println("indexing_wallet3: address = " + idxOwner3.String() + ", individual = " + individualTotalReward.String() + ",\tmatureTotal = " + matureTotalIdxNode3.String() + ",\timmatureTotal = " + immatureTotalIdxNode3.String())
	fmt.Println("***************************************************************************************")
}

func testFullDistributeProcessAtEpoch1(t *testing.T, ctx sdk.Context, k Keeper, trafficList []types.SingleWalletVolume) {
	//PrePay
	k.SetTotalUnissuedPrepay(ctx, totalUnissuedPrePay)

	_, err := k.DistributePotReward(ctx, trafficList, epoch1)
	require.NoError(t, err)

	fmt.Println("Distribution result at Epoch1: ")
	rewardAddrList := k.GetRewardAddressPool(ctx)
	fmt.Println("address pool: ")
	for i := 0; i < len(rewardAddrList); i++ {
		fmt.Println(rewardAddrList[i].String() + ", ")
	}
	fmt.Println("----------------------------------------------------------------------------------")

	idvRwdResNode1Ep1, _ := k.GetIndividualReward(ctx, resOwner1, epoch2017)
	individualTotalReward := idvRwdResNode1Ep1.RewardFromMiningPool.Add(idvRwdResNode1Ep1.RewardFromTrafficPool)
	matureTotalResNode1 := k.GetMatureTotalReward(ctx, resOwner1)
	immatureTotalResNode1 := k.GetImmatureTotalReward(ctx, resOwner1)
	fmt.Println("resource_wallet1: address = " + resOwner1.String() + ", individual = " + individualTotalReward.String() + ",\tmatureTotal = " + matureTotalResNode1.String() + ",\timmatureTotal = " + immatureTotalResNode1.String())
	require.Equal(t, individualTotalReward, sdk.NewInt(68846296294))
	require.Equal(t, matureTotalResNode1, sdk.ZeroInt())
	require.Equal(t, immatureTotalResNode1, sdk.NewInt(68846296294))
	fmt.Println("resource_wallet1: address = " + resOwner1.String() + ", individual = " + individualTotalReward.String() + ",\tmatureTotal = " + matureTotalResNode1.String() + ",\timmatureTotal = " + immatureTotalResNode1.String())

	idvRwdResNode2Ep1, _ := k.GetIndividualReward(ctx, resOwner2, epoch2017)
	individualTotalReward = idvRwdResNode2Ep1.RewardFromMiningPool.Add(idvRwdResNode2Ep1.RewardFromTrafficPool)
	matureTotalResNode2 := k.GetMatureTotalReward(ctx, resOwner2)
	immatureTotalResNode2 := k.GetImmatureTotalReward(ctx, resOwner2)
	require.Equal(t, individualTotalReward, sdk.NewInt(42479629627))
	require.Equal(t, matureTotalResNode2, sdk.ZeroInt())
	require.Equal(t, immatureTotalResNode2, sdk.NewInt(42479629627))
	fmt.Println("resource_wallet2: address = " + resOwner2.String() + ", individual = " + individualTotalReward.String() + ",\tmatureTotal = " + matureTotalResNode2.String() + ",\timmatureTotal = " + immatureTotalResNode2.String())

	idvRwdResNode3Ep1, _ := k.GetIndividualReward(ctx, resOwner3, epoch2017)
	individualTotalReward = idvRwdResNode3Ep1.RewardFromMiningPool.Add(idvRwdResNode3Ep1.RewardFromTrafficPool)
	matureTotalResNode3 := k.GetMatureTotalReward(ctx, resOwner3)
	immatureTotalResNode3 := k.GetImmatureTotalReward(ctx, resOwner3)
	require.Equal(t, individualTotalReward, sdk.NewInt(29296296294))
	require.Equal(t, matureTotalResNode3, sdk.ZeroInt())
	require.Equal(t, immatureTotalResNode3, sdk.NewInt(29296296294))
	fmt.Println("resource_wallet3: address = " + resOwner3.String() + ", individual = " + individualTotalReward.String() + ",\tmatureTotal = " + matureTotalResNode3.String() + ",\timmatureTotal = " + immatureTotalResNode3.String())

	idvRwdResNode4Ep1, _ := k.GetIndividualReward(ctx, resOwner4, epoch2017)
	individualTotalReward = idvRwdResNode4Ep1.RewardFromMiningPool.Add(idvRwdResNode4Ep1.RewardFromTrafficPool)
	matureTotalResNode4 := k.GetMatureTotalReward(ctx, resOwner4)
	immatureTotalResNode4 := k.GetImmatureTotalReward(ctx, resOwner4)
	require.Equal(t, individualTotalReward, sdk.NewInt(2929629628))
	require.Equal(t, matureTotalResNode4, sdk.ZeroInt())
	require.Equal(t, immatureTotalResNode4, sdk.NewInt(2929629628))
	fmt.Println("resource_wallet4: address = " + resOwner4.String() + ", individual = " + individualTotalReward.String() + ", \tmatureTotal = " + matureTotalResNode4.String() + ",\timmatureTotal = " + immatureTotalResNode4.String())

	idvRwdResNode5Ep1, _ := k.GetIndividualReward(ctx, resOwner5, epoch2017)
	individualTotalReward = idvRwdResNode5Ep1.RewardFromMiningPool.Add(idvRwdResNode5Ep1.RewardFromTrafficPool)
	matureTotalResNode5 := k.GetMatureTotalReward(ctx, resOwner5)
	immatureTotalResNode5 := k.GetImmatureTotalReward(ctx, resOwner5)
	require.Equal(t, individualTotalReward, sdk.NewInt(2929629628))
	require.Equal(t, matureTotalResNode5, sdk.ZeroInt())
	require.Equal(t, immatureTotalResNode5, sdk.NewInt(2929629628))
	fmt.Println("resource_wallet5: address = " + resOwner5.String() + ", individual = " + individualTotalReward.String() + ", \tmatureTotal = " + matureTotalResNode5.String() + ",\timmatureTotal = " + immatureTotalResNode5.String())

	idvRwdIdxNode1Ep1, _ := k.GetIndividualReward(ctx, idxOwner1, epoch2017)
	individualTotalReward = idvRwdIdxNode1Ep1.RewardFromMiningPool.Add(idvRwdIdxNode1Ep1.RewardFromTrafficPool)
	matureTotalIdxNode1 := k.GetMatureTotalReward(ctx, idxOwner1)
	immatureTotalIdxNode1 := k.GetImmatureTotalReward(ctx, idxOwner1)
	require.Equal(t, individualTotalReward, sdk.NewInt(19530864195))
	require.Equal(t, matureTotalIdxNode1, sdk.ZeroInt())
	require.Equal(t, immatureTotalIdxNode1, sdk.NewInt(19530864195))
	fmt.Println("indexing_wallet1: address = " + idxOwner1.String() + ", individual = " + individualTotalReward.String() + ",\tmatureTotal = " + matureTotalIdxNode1.String() + ",\timmatureTotal = " + immatureTotalIdxNode1.String())

	idvRwdIdxNode2Ep1, _ := k.GetIndividualReward(ctx, idxOwner2, epoch2017)
	individualTotalReward = idvRwdIdxNode2Ep1.RewardFromMiningPool.Add(idvRwdIdxNode2Ep1.RewardFromTrafficPool)
	matureTotalIdxNode2 := k.GetMatureTotalReward(ctx, idxOwner2)
	immatureTotalIdxNode2 := k.GetImmatureTotalReward(ctx, idxOwner2)
	require.Equal(t, individualTotalReward, sdk.NewInt(19530864195))
	require.Equal(t, matureTotalIdxNode2, sdk.ZeroInt())
	require.Equal(t, immatureTotalIdxNode2, sdk.NewInt(19530864195))
	fmt.Println("indexing_wallet2: address = " + idxOwner2.String() + ", individual = " + individualTotalReward.String() + ",\tmatureTotal = " + matureTotalIdxNode2.String() + ",\timmatureTotal = " + immatureTotalIdxNode2.String())

	idvRwdIdxNode3Ep1, _ := k.GetIndividualReward(ctx, idxOwner3, epoch2017)
	individualTotalReward = idvRwdIdxNode3Ep1.RewardFromMiningPool.Add(idvRwdIdxNode3Ep1.RewardFromTrafficPool)
	matureTotalIdxNode3 := k.GetMatureTotalReward(ctx, idxOwner3)
	immatureTotalIdxNode3 := k.GetImmatureTotalReward(ctx, idxOwner3)
	require.Equal(t, individualTotalReward, sdk.NewInt(19530864195))
	require.Equal(t, matureTotalIdxNode3, sdk.ZeroInt())
	require.Equal(t, immatureTotalIdxNode3, sdk.NewInt(19530864195))
	fmt.Println("indexing_wallet3: address = " + idxOwner3.String() + ", individual = " + individualTotalReward.String() + ",\tmatureTotal = " + matureTotalIdxNode3.String() + ",\timmatureTotal = " + immatureTotalIdxNode3.String())
	fmt.Println("***************************************************************************************")
}

// 20% of traffic reward distribute to all validators/delegators by shares of stake
func testBlockChainRewardFromTrafficPool(t *testing.T, ctx sdk.Context, k Keeper, bankKeeper bank.Keeper, trafficList []types.SingleWalletVolume) {
	distributeGoal := types.InitDistributeGoal()
	rewardDetailMap := make(map[string]types.Reward)

	//1, calc traffic reward in total
	_, distributeGoal, err := k.CalcTrafficRewardInTotal(ctx, trafficList, distributeGoal)
	require.NoError(t, err)

	// stake reward split by the amount of delegation/deposit
	// total delegation of validator/resource node/indexing node is 15stos
	require.Equal(t, distributeGoal.BlockChainRewardToValidatorFromTrafficPool, distributeGoal.BlockChainRewardToIndexingNodeFromTrafficPool)
	require.Equal(t, distributeGoal.BlockChainRewardToIndexingNodeFromTrafficPool, distributeGoal.BlockChainRewardToResourceNodeFromTrafficPool)

	//Only keep blockchain reward to test
	distributeGoal.MetaNodeRewardToIndexingNodeFromTrafficPool = sdk.ZeroInt()
	distributeGoal.TrafficRewardToResourceNodeFromTrafficPool = sdk.ZeroInt()
	fmt.Println("testBlockChainRewardFromTrafficPool: \n" + distributeGoal.String())

	//Get excepted reward before calculation method changed the value of distributeGoal
	exceptedValRwd := distributeGoal.BlockChainRewardToValidatorFromTrafficPool
	exceptedResNodeRwd := distributeGoal.BlockChainRewardToResourceNodeFromTrafficPool.ToDec().Quo(sdk.NewDec(5)).TruncateInt()
	exceptedIdxNodeRwd := distributeGoal.BlockChainRewardToIndexingNodeFromTrafficPool.ToDec().Quo(sdk.NewDec(3)).TruncateInt()
	feePoolBefore := getFeePoolBalance(t, ctx, k, bankKeeper)

	/********************************* after calculation method, value of distributeGoal object will change ******************************************/
	//3, calc reward for resource node
	rewardDetailMap, distributeGoal = k.CalcRewardForResourceNode(ctx, trafficList, distributeGoal, rewardDetailMap)
	//4, calc reward from indexing node
	rewardDetailMap, distributeGoal = k.CalcRewardForIndexingNode(ctx, distributeGoal, rewardDetailMap)
	//5, deduct reward from provider account
	err = k.deductRewardFromRewardProviderAccount(ctx, distributeGoal, epoch1)
	require.NoError(t, err)
	//6, distribute skate reward to fee pool for validators
	distributeGoal, err = k.distributeValidatorRewardToFeePool(ctx, distributeGoal)
	require.NoError(t, err)

	feePoolAfter := getFeePoolBalance(t, ctx, k, bankKeeper)

	require.Equal(t, feePoolBefore.Add(sdk.NewCoin(k.BondDenom(ctx), exceptedValRwd)), feePoolAfter)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner1.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner2.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner3.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner4.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner5.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner1.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner2.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner3.String()].RewardFromTrafficPool)

	fmt.Println("reward to fee pool： " + feePoolAfter.Sub(feePoolBefore).String())
	fmt.Println("resource_wallet1： address = " + resOwner1.String() + ", reward = " + rewardDetailMap[resOwner1.String()].RewardFromTrafficPool.String())
	fmt.Println("resource_wallet2： address = " + resOwner2.String() + ", reward = " + rewardDetailMap[resOwner2.String()].RewardFromTrafficPool.String())
	fmt.Println("resource_wallet3： address = " + resOwner3.String() + ", reward = " + rewardDetailMap[resOwner3.String()].RewardFromTrafficPool.String())
	fmt.Println("resource_wallet4： address = " + resOwner4.String() + ", reward = " + rewardDetailMap[resOwner4.String()].RewardFromTrafficPool.String())
	fmt.Println("resource_wallet5： address = " + resOwner5.String() + ", reward = " + rewardDetailMap[resOwner5.String()].RewardFromTrafficPool.String())
	fmt.Println("indexing_wallet1： address = " + idxOwner1.String() + ", reward = " + rewardDetailMap[idxOwner1.String()].RewardFromTrafficPool.String())
	fmt.Println("indexing_wallet2： address = " + idxOwner2.String() + ", reward = " + rewardDetailMap[idxOwner2.String()].RewardFromTrafficPool.String())
	fmt.Println("indexing_wallet3： address = " + idxOwner3.String() + ", reward = " + rewardDetailMap[idxOwner3.String()].RewardFromTrafficPool.String())
	fmt.Println("***************************************************************************************")
}

// 20% of traffic reward equally distribute to all indexing nodes
func testMetaNodeRewardFromTrafficPool(t *testing.T, ctx sdk.Context, k Keeper, bankKeeper bank.Keeper, trafficList []types.SingleWalletVolume) {
	distributeGoal := types.InitDistributeGoal()
	rewardDetailMap := make(map[string]types.Reward)

	_, totalReward := k.getTrafficReward(ctx, trafficList)

	//1, calc traffic reward in total
	_, distributeGoal, err := k.CalcTrafficRewardInTotal(ctx, trafficList, distributeGoal)
	require.NoError(t, err)

	//Only keep meta node reward to test
	distributeGoal.BlockChainRewardToValidatorFromTrafficPool = sdk.ZeroInt()
	distributeGoal.BlockChainRewardToIndexingNodeFromTrafficPool = sdk.ZeroInt()
	distributeGoal.BlockChainRewardToResourceNodeFromTrafficPool = sdk.ZeroInt()
	distributeGoal.TrafficRewardToResourceNodeFromTrafficPool = sdk.ZeroInt()
	fmt.Println("testMetaNodeRewardFromTrafficPool: \n" + distributeGoal.String())

	//20% of traffic reward to meta nodes
	exceptedTotalRewardToMetaNodes := totalReward.Mul(sdk.NewDecWithPrec(20, 2)).TruncateInt()
	require.Equal(t, exceptedTotalRewardToMetaNodes, distributeGoal.MetaNodeRewardToIndexingNodeFromTrafficPool)

	//indexing node 1,2,3 have the same share of the meta node reward
	exceptedIdxNodeRwd := distributeGoal.MetaNodeRewardToIndexingNodeFromTrafficPool.ToDec().Quo(sdk.NewDec(3)).TruncateInt()
	exceptedResNodeRwd := sdk.ZeroInt()
	feePoolBefore := getFeePoolBalance(t, ctx, k, bankKeeper)

	/********************************* after calculation method, value of distributeGoal object will change ******************************************/
	//3, calc reward for resource node
	rewardDetailMap, distributeGoal = k.CalcRewardForResourceNode(ctx, trafficList, distributeGoal, rewardDetailMap)
	//4, calc reward from indexing node
	rewardDetailMap, distributeGoal = k.CalcRewardForIndexingNode(ctx, distributeGoal, rewardDetailMap)
	//5, deduct reward from provider account
	err = k.deductRewardFromRewardProviderAccount(ctx, distributeGoal, epoch1)
	require.NoError(t, err)
	//6, distribute skate reward to fee pool for validators
	distributeGoal, err = k.distributeValidatorRewardToFeePool(ctx, distributeGoal)
	require.NoError(t, err)

	feePoolAfter := getFeePoolBalance(t, ctx, k, bankKeeper)

	require.Equal(t, feePoolBefore, feePoolAfter)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner1.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner2.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner3.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner4.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner5.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner1.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner2.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner3.String()].RewardFromTrafficPool)

	fmt.Println("reward to fee pool： " + feePoolAfter.Sub(feePoolBefore).String())
	fmt.Println("resource_wallet1： address = " + resOwner1.String() + ", reward = " + rewardDetailMap[resOwner1.String()].RewardFromTrafficPool.String())
	fmt.Println("resource_wallet2： address = " + resOwner2.String() + ", reward = " + rewardDetailMap[resOwner2.String()].RewardFromTrafficPool.String())
	fmt.Println("resource_wallet3： address = " + resOwner3.String() + ", reward = " + rewardDetailMap[resOwner3.String()].RewardFromTrafficPool.String())
	fmt.Println("resource_wallet4： address = " + resOwner4.String() + ", reward = " + rewardDetailMap[resOwner4.String()].RewardFromTrafficPool.String())
	fmt.Println("resource_wallet5： address = " + resOwner5.String() + ", reward = " + rewardDetailMap[resOwner5.String()].RewardFromTrafficPool.String())
	fmt.Println("indexing_wallet1： address = " + idxOwner1.String() + ", reward = " + rewardDetailMap[idxOwner1.String()].RewardFromTrafficPool.String())
	fmt.Println("indexing_wallet2： address = " + idxOwner2.String() + ", reward = " + rewardDetailMap[idxOwner2.String()].RewardFromTrafficPool.String())
	fmt.Println("indexing_wallet3： address = " + idxOwner3.String() + ", reward = " + rewardDetailMap[idxOwner3.String()].RewardFromTrafficPool.String())
	fmt.Println("***************************************************************************************")
}

// 60% of traffic reward distribute to resource nodes by traffic
func testTrafficRewardFromTrafficPool(t *testing.T, ctx sdk.Context, k Keeper, bankKeeper bank.Keeper, trafficList []types.SingleWalletVolume) {
	distributeGoal := types.InitDistributeGoal()
	rewardDetailMap := make(map[string]types.Reward)

	_, totalReward := k.getTrafficReward(ctx, trafficList)

	//1, calc traffic reward in total
	_, distributeGoal, err := k.CalcTrafficRewardInTotal(ctx, trafficList, distributeGoal)
	require.NoError(t, err)

	//Only keep traffic reward to test
	distributeGoal.BlockChainRewardToValidatorFromTrafficPool = sdk.ZeroInt()
	distributeGoal.BlockChainRewardToIndexingNodeFromTrafficPool = sdk.ZeroInt()
	distributeGoal.BlockChainRewardToResourceNodeFromTrafficPool = sdk.ZeroInt()
	distributeGoal.MetaNodeRewardToIndexingNodeFromTrafficPool = sdk.ZeroInt()
	fmt.Println("testTrafficRewardFromTrafficPool: \n" + distributeGoal.String())

	//60% of traffic reward to resource nodes
	exceptedTotalRewardToResNodes := totalReward.Mul(sdk.NewDecWithPrec(60, 2)).TruncateInt()
	require.Equal(t, exceptedTotalRewardToResNodes, distributeGoal.TrafficRewardToResourceNodeFromTrafficPool)

	//resource node 1,2,3 are in the volume report, so they have stake reward AND traffic reward in this epoch
	exceptedResNode1Rwd := distributeGoal.TrafficRewardToResourceNodeFromTrafficPool.ToDec().Mul(sdk.NewDec(resourceNodeVolume1)).Quo(sdk.NewDec(totalVolume)).TruncateInt()
	exceptedResNode2Rwd := distributeGoal.TrafficRewardToResourceNodeFromTrafficPool.ToDec().Mul(sdk.NewDec(resourceNodeVolume2)).Quo(sdk.NewDec(totalVolume)).TruncateInt()
	exceptedResNode3Rwd := distributeGoal.TrafficRewardToResourceNodeFromTrafficPool.ToDec().Mul(sdk.NewDec(resourceNodeVolume3)).Quo(sdk.NewDec(totalVolume)).TruncateInt()
	//resource node 4&5 are not in the volume report, so they only have stake reward in this epoch
	exceptedResNode4Rwd := sdk.ZeroInt()
	exceptedResNode5Rwd := sdk.ZeroInt()
	//indexing node 1,2,3 only have stake reward and meta node reward in this epoch
	exceptedIdxNodeRwd := sdk.ZeroInt()
	feePoolBefore := getFeePoolBalance(t, ctx, k, bankKeeper)

	/********************************* after calculation method, value of distributeGoal object will change ******************************************/
	//3, calc reward for resource node
	rewardDetailMap, distributeGoal = k.CalcRewardForResourceNode(ctx, trafficList, distributeGoal, rewardDetailMap)
	//4, calc reward from indexing node
	rewardDetailMap, distributeGoal = k.CalcRewardForIndexingNode(ctx, distributeGoal, rewardDetailMap)
	//5, deduct reward from provider account
	err = k.deductRewardFromRewardProviderAccount(ctx, distributeGoal, epoch1)
	require.NoError(t, err)
	//6, distribute skate reward to fee pool for validators
	distributeGoal, err = k.distributeValidatorRewardToFeePool(ctx, distributeGoal)
	require.NoError(t, err)
	feePoolAfter := getFeePoolBalance(t, ctx, k, bankKeeper)

	require.Equal(t, feePoolBefore, feePoolAfter)
	require.Equal(t, exceptedResNode1Rwd, rewardDetailMap[resOwner1.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedResNode2Rwd, rewardDetailMap[resOwner2.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedResNode3Rwd, rewardDetailMap[resOwner3.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedResNode4Rwd, rewardDetailMap[resOwner4.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedResNode5Rwd, rewardDetailMap[resOwner5.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner1.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner2.String()].RewardFromTrafficPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner3.String()].RewardFromTrafficPool)

	fmt.Println("reward to fee pool： " + feePoolAfter.Sub(feePoolBefore).String())
	fmt.Println("resource_wallet1： address = " + resOwner1.String() + ", reward = " + rewardDetailMap[resOwner1.String()].RewardFromTrafficPool.String())
	fmt.Println("resource_wallet2： address = " + resOwner2.String() + ", reward = " + rewardDetailMap[resOwner2.String()].RewardFromTrafficPool.String())
	fmt.Println("resource_wallet3： address = " + resOwner3.String() + ", reward = " + rewardDetailMap[resOwner3.String()].RewardFromTrafficPool.String())
	fmt.Println("resource_wallet4： address = " + resOwner4.String() + ", reward = " + rewardDetailMap[resOwner4.String()].RewardFromTrafficPool.String())
	fmt.Println("resource_wallet5： address = " + resOwner5.String() + ", reward = " + rewardDetailMap[resOwner5.String()].RewardFromTrafficPool.String())
	fmt.Println("indexing_wallet1： address = " + idxOwner1.String() + ", reward = " + rewardDetailMap[idxOwner1.String()].RewardFromTrafficPool.String())
	fmt.Println("indexing_wallet2： address = " + idxOwner2.String() + ", reward = " + rewardDetailMap[idxOwner2.String()].RewardFromTrafficPool.String())
	fmt.Println("indexing_wallet3： address = " + idxOwner3.String() + ", reward = " + rewardDetailMap[idxOwner3.String()].RewardFromTrafficPool.String())
	fmt.Println("***************************************************************************************")
}

// 20% of mining reward distribute to all validators/delegators by shares of stake
func testBlockChainRewardFromMiningPool(t *testing.T, ctx sdk.Context, k Keeper, bankKeeper bank.Keeper, trafficList []types.SingleWalletVolume) {
	distributeGoal := types.InitDistributeGoal()
	rewardDetailMap := make(map[string]types.Reward)

	//2, calc mining reward in total
	distributeGoal, err := k.CalcMiningRewardInTotal(ctx, distributeGoal)
	require.NoError(t, err)

	totalMiningReward := distributeGoal.BlockChainRewardToValidatorFromMiningPool.Add(distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool).
		Add(distributeGoal.BlockChainRewardToResourceNodeFromMiningPool)

	// since validators, indexing nodes and resource nodes have the same total stake in this test case,
	// total block chain reward from mining pool needs to be divisible by 3
	exceptedTotalReward := sdk.NewDec(80000000000).Mul(sdk.NewDecWithPrec(20, 2)).Quo(sdk.NewDec(3)).TruncateDec().Mul(sdk.NewDec(3)).TruncateInt()
	require.Equal(t, exceptedTotalReward, totalMiningReward)
	// stake reward split by the amount of delegation/deposit
	// total delegation of validator/resource node/indexing node is 15stos
	require.Equal(t, distributeGoal.BlockChainRewardToValidatorFromMiningPool, distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool)
	require.Equal(t, distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool, distributeGoal.BlockChainRewardToResourceNodeFromMiningPool)

	//Only keep blockchain reward to test
	distributeGoal.MetaNodeRewardToIndexingNodeFromMiningPool = sdk.ZeroInt()
	distributeGoal.TrafficRewardToResourceNodeFromMiningPool = sdk.ZeroInt()
	fmt.Println("testBlockChainRewardFromMiningPool: \n" + distributeGoal.String())

	//Get excepted reward before calculation method changed the value of distributeGoal
	exceptedValRwd := distributeGoal.BlockChainRewardToValidatorFromMiningPool
	exceptedResNodeRwd := distributeGoal.BlockChainRewardToResourceNodeFromMiningPool.ToDec().Quo(sdk.NewDec(5)).TruncateInt()
	exceptedIdxNodeRwd := distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool.ToDec().Quo(sdk.NewDec(3)).TruncateInt()
	feePoolBefore := getFeePoolBalance(t, ctx, k, bankKeeper)

	/********************************* after calculation method, value of distributeGoal object will change ******************************************/
	//3, calc reward for resource node
	rewardDetailMap, distributeGoal = k.CalcRewardForResourceNode(ctx, trafficList, distributeGoal, rewardDetailMap)
	//4, calc reward from indexing node
	rewardDetailMap, distributeGoal = k.CalcRewardForIndexingNode(ctx, distributeGoal, rewardDetailMap)
	//5, deduct reward from provider account
	err = k.deductRewardFromRewardProviderAccount(ctx, distributeGoal, epoch1)
	require.NoError(t, err)
	//6, distribute skate reward to fee pool for validators
	distributeGoal, err = k.distributeValidatorRewardToFeePool(ctx, distributeGoal)
	require.NoError(t, err)

	feePoolAfter := getFeePoolBalance(t, ctx, k, bankKeeper)

	require.Equal(t, feePoolBefore.Add(sdk.NewCoin(k.BondDenom(ctx), exceptedValRwd)), feePoolAfter)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner1.String()].RewardFromMiningPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner2.String()].RewardFromMiningPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner3.String()].RewardFromMiningPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner4.String()].RewardFromMiningPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner5.String()].RewardFromMiningPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner1.String()].RewardFromMiningPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner2.String()].RewardFromMiningPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner3.String()].RewardFromMiningPool)

	fmt.Println("reward to fee pool： " + feePoolAfter.Sub(feePoolBefore).String())
	fmt.Println("resource_wallet1： address = " + resOwner1.String() + ", reward = " + rewardDetailMap[resOwner1.String()].RewardFromMiningPool.String())
	fmt.Println("resource_wallet2： address = " + resOwner2.String() + ", reward = " + rewardDetailMap[resOwner2.String()].RewardFromMiningPool.String())
	fmt.Println("resource_wallet3： address = " + resOwner3.String() + ", reward = " + rewardDetailMap[resOwner3.String()].RewardFromMiningPool.String())
	fmt.Println("resource_wallet4： address = " + resOwner4.String() + ", reward = " + rewardDetailMap[resOwner4.String()].RewardFromMiningPool.String())
	fmt.Println("resource_wallet5： address = " + resOwner5.String() + ", reward = " + rewardDetailMap[resOwner5.String()].RewardFromMiningPool.String())
	fmt.Println("indexing_wallet1： address = " + idxOwner1.String() + ", reward = " + rewardDetailMap[idxOwner1.String()].RewardFromMiningPool.String())
	fmt.Println("indexing_wallet2： address = " + idxOwner2.String() + ", reward = " + rewardDetailMap[idxOwner2.String()].RewardFromMiningPool.String())
	fmt.Println("indexing_wallet3： address = " + idxOwner3.String() + ", reward = " + rewardDetailMap[idxOwner3.String()].RewardFromMiningPool.String())
	fmt.Println("***************************************************************************************")
}

// 20% of mining reward equally distribute to all indexing nodes
func testMetaNodeRewardFromMiningPool(t *testing.T, ctx sdk.Context, k Keeper, bankKeeper bank.Keeper, trafficList []types.SingleWalletVolume) {
	distributeGoal := types.InitDistributeGoal()
	rewardDetailMap := make(map[string]types.Reward)

	totalReward := sdk.NewDec(80000000000)

	//2, calc mining reward in total
	distributeGoal, err := k.CalcMiningRewardInTotal(ctx, distributeGoal)
	require.NoError(t, err)

	//Only keep meta node reward to test
	distributeGoal.BlockChainRewardToValidatorFromMiningPool = sdk.ZeroInt()
	distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool = sdk.ZeroInt()
	distributeGoal.BlockChainRewardToResourceNodeFromMiningPool = sdk.ZeroInt()
	distributeGoal.TrafficRewardToResourceNodeFromMiningPool = sdk.ZeroInt()
	fmt.Println("testMetaNodeRewardFromMiningPool: \n" + distributeGoal.String())

	//20% of mining reward to meta nodes
	exceptedTotalRewardToMetaNodes := totalReward.Mul(sdk.NewDecWithPrec(20, 2)).TruncateInt()
	require.Equal(t, exceptedTotalRewardToMetaNodes, distributeGoal.MetaNodeRewardToIndexingNodeFromMiningPool)

	//indexing node 1,2,3 have the same share of the meta node reward
	exceptedIdxNodeRwd := distributeGoal.MetaNodeRewardToIndexingNodeFromMiningPool.ToDec().Quo(sdk.NewDec(3)).TruncateInt()
	exceptedResNodeRwd := sdk.ZeroInt()
	feePoolBefore := getFeePoolBalance(t, ctx, k, bankKeeper)

	/********************************* after calculation method, value of distributeGoal object will change ******************************************/
	//3, calc reward for resource node
	rewardDetailMap, distributeGoal = k.CalcRewardForResourceNode(ctx, trafficList, distributeGoal, rewardDetailMap)
	//4, calc reward from indexing node
	rewardDetailMap, distributeGoal = k.CalcRewardForIndexingNode(ctx, distributeGoal, rewardDetailMap)
	//5, deduct reward from provider account
	err = k.deductRewardFromRewardProviderAccount(ctx, distributeGoal, epoch1)
	require.NoError(t, err)
	//6, distribute skate reward to fee pool for validators
	distributeGoal, err = k.distributeValidatorRewardToFeePool(ctx, distributeGoal)
	require.NoError(t, err)

	feePoolAfter := getFeePoolBalance(t, ctx, k, bankKeeper)

	require.Equal(t, feePoolBefore, feePoolAfter)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner1.String()].RewardFromMiningPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner2.String()].RewardFromMiningPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner3.String()].RewardFromMiningPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner4.String()].RewardFromMiningPool)
	require.Equal(t, exceptedResNodeRwd, rewardDetailMap[resOwner5.String()].RewardFromMiningPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner1.String()].RewardFromMiningPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner2.String()].RewardFromMiningPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner3.String()].RewardFromMiningPool)

	fmt.Println("reward to fee pool： " + feePoolAfter.Sub(feePoolBefore).String())
	fmt.Println("resource_wallet1： address = " + resOwner1.String() + ", reward = " + rewardDetailMap[resOwner1.String()].RewardFromMiningPool.String())
	fmt.Println("resource_wallet2： address = " + resOwner2.String() + ", reward = " + rewardDetailMap[resOwner2.String()].RewardFromMiningPool.String())
	fmt.Println("resource_wallet3： address = " + resOwner3.String() + ", reward = " + rewardDetailMap[resOwner3.String()].RewardFromMiningPool.String())
	fmt.Println("resource_wallet4： address = " + resOwner4.String() + ", reward = " + rewardDetailMap[resOwner4.String()].RewardFromMiningPool.String())
	fmt.Println("resource_wallet5： address = " + resOwner5.String() + ", reward = " + rewardDetailMap[resOwner5.String()].RewardFromMiningPool.String())
	fmt.Println("indexing_wallet1： address = " + idxOwner1.String() + ", reward = " + rewardDetailMap[idxOwner1.String()].RewardFromMiningPool.String())
	fmt.Println("indexing_wallet2： address = " + idxOwner2.String() + ", reward = " + rewardDetailMap[idxOwner2.String()].RewardFromMiningPool.String())
	fmt.Println("indexing_wallet3： address = " + idxOwner3.String() + ", reward = " + rewardDetailMap[idxOwner3.String()].RewardFromMiningPool.String())
	fmt.Println("***************************************************************************************")
}

// 60% of mining reward distribute to resource nodes by traffic
func testTrafficRewardFromMiningPool(t *testing.T, ctx sdk.Context, k Keeper, bankKeeper bank.Keeper, trafficList []types.SingleWalletVolume) {
	distributeGoal := types.InitDistributeGoal()
	rewardDetailMap := make(map[string]types.Reward)

	totalReward := sdk.NewDec(80000000000)

	//2, calc mining reward in total
	distributeGoal, err := k.CalcMiningRewardInTotal(ctx, distributeGoal)
	require.NoError(t, err)

	//Only keep traffic reward to test
	distributeGoal.BlockChainRewardToValidatorFromMiningPool = sdk.ZeroInt()
	distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool = sdk.ZeroInt()
	distributeGoal.BlockChainRewardToResourceNodeFromMiningPool = sdk.ZeroInt()
	distributeGoal.MetaNodeRewardToIndexingNodeFromMiningPool = sdk.ZeroInt()
	fmt.Println("testTrafficRewardFromMiningPool: \n" + distributeGoal.String())

	//60% of mining reward to resource nodes
	exceptedTotalRewardToResNodes := totalReward.Mul(sdk.NewDecWithPrec(60, 2)).TruncateInt()
	require.Equal(t, exceptedTotalRewardToResNodes, distributeGoal.TrafficRewardToResourceNodeFromMiningPool)

	//resource node 1,2,3 are in the volume report, so they have stake reward AND traffic reward in this epoch
	exceptedResNode1Rwd := distributeGoal.TrafficRewardToResourceNodeFromMiningPool.ToDec().Mul(sdk.NewDec(resourceNodeVolume1)).Quo(sdk.NewDec(totalVolume)).TruncateInt()
	exceptedResNode2Rwd := distributeGoal.TrafficRewardToResourceNodeFromMiningPool.ToDec().Mul(sdk.NewDec(resourceNodeVolume2)).Quo(sdk.NewDec(totalVolume)).TruncateInt()
	exceptedResNode3Rwd := distributeGoal.TrafficRewardToResourceNodeFromMiningPool.ToDec().Mul(sdk.NewDec(resourceNodeVolume3)).Quo(sdk.NewDec(totalVolume)).TruncateInt()
	//resource node 4&5 are not in the volume report, so they only have stake reward in this epoch
	exceptedResNode4Rwd := sdk.ZeroInt()
	exceptedResNode5Rwd := sdk.ZeroInt()
	//indexing node 1,2,3 only have stake reward and meta node reward in this epoch
	exceptedIdxNodeRwd := sdk.ZeroInt()
	feePoolBefore := getFeePoolBalance(t, ctx, k, bankKeeper)

	/********************************* after calculation method, value of distributeGoal object will change ******************************************/
	//3, calc reward for resource node
	rewardDetailMap, distributeGoal = k.CalcRewardForResourceNode(ctx, trafficList, distributeGoal, rewardDetailMap)
	//4, calc reward from indexing node
	rewardDetailMap, distributeGoal = k.CalcRewardForIndexingNode(ctx, distributeGoal, rewardDetailMap)
	//5, deduct reward from provider account
	err = k.deductRewardFromRewardProviderAccount(ctx, distributeGoal, epoch1)
	require.NoError(t, err)
	//6, distribute skate reward to fee pool for validators
	distributeGoal, err = k.distributeValidatorRewardToFeePool(ctx, distributeGoal)
	require.NoError(t, err)
	feePoolAfter := getFeePoolBalance(t, ctx, k, bankKeeper)

	require.Equal(t, feePoolBefore, feePoolAfter)
	require.Equal(t, exceptedResNode1Rwd, rewardDetailMap[resOwner1.String()].RewardFromMiningPool)
	require.Equal(t, exceptedResNode2Rwd, rewardDetailMap[resOwner2.String()].RewardFromMiningPool)
	require.Equal(t, exceptedResNode3Rwd, rewardDetailMap[resOwner3.String()].RewardFromMiningPool)
	require.Equal(t, exceptedResNode4Rwd, rewardDetailMap[resOwner4.String()].RewardFromMiningPool)
	require.Equal(t, exceptedResNode5Rwd, rewardDetailMap[resOwner5.String()].RewardFromMiningPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner1.String()].RewardFromMiningPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner2.String()].RewardFromMiningPool)
	require.Equal(t, exceptedIdxNodeRwd, rewardDetailMap[idxOwner3.String()].RewardFromMiningPool)

	fmt.Println("reward to fee pool： " + feePoolAfter.Sub(feePoolBefore).String())
	fmt.Println("resourceNode1： address = " + resOwner1.String() + ", reward = " + rewardDetailMap[resOwner1.String()].RewardFromMiningPool.String())
	fmt.Println("resourceNode2： address = " + resOwner2.String() + ", reward = " + rewardDetailMap[resOwner2.String()].RewardFromMiningPool.String())
	fmt.Println("resourceNode3： address = " + resOwner3.String() + ", reward = " + rewardDetailMap[resOwner3.String()].RewardFromMiningPool.String())
	fmt.Println("resourceNode4： address = " + resOwner4.String() + ", reward = " + rewardDetailMap[resOwner4.String()].RewardFromMiningPool.String())
	fmt.Println("resourceNode5： address = " + resOwner5.String() + ", reward = " + rewardDetailMap[resOwner5.String()].RewardFromMiningPool.String())
	fmt.Println("indexingNode1： address = " + idxOwner1.String() + ", reward = " + rewardDetailMap[idxOwner1.String()].RewardFromMiningPool.String())
	fmt.Println("indexingNode2： address = " + idxOwner2.String() + ", reward = " + rewardDetailMap[idxOwner2.String()].RewardFromMiningPool.String())
	fmt.Println("indexingNode3： address = " + idxOwner3.String() + ", reward = " + rewardDetailMap[idxOwner3.String()].RewardFromMiningPool.String())
	fmt.Println("***************************************************************************************")
}

func createAccount(t *testing.T, ctx sdk.Context, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, acc sdk.AccAddress, coins sdk.Coins) {
	account := accountKeeper.GetAccount(ctx, acc)
	if account == nil {
		account = accountKeeper.NewAccountWithAddress(ctx, acc)
		//fmt.Printf("create account: " + account.String() + "\n")
	}
	coins, err := bankKeeper.AddCoins(ctx, acc, coins)
	require.NoError(t, err)
}

func getFeePoolBalance(t *testing.T, ctx sdk.Context, k Keeper, bankKeeper bank.Keeper) sdk.Coins {
	feePoolAccAddr := k.SupplyKeeper.GetModuleAddress(k.feeCollectorName)
	require.NotNil(t, feePoolAccAddr)
	coins := bankKeeper.GetCoins(ctx, feePoolAccAddr)
	return coins
}
