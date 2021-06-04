package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

const (
	foundationAccAddr = "st1qr9set2jaayzjjpm9tw4f3n6f5zfu3hef8wtaw"
	foundationDeposit = 400000000000000

	ownerAddr1 = "st1xgw0rrdgjdun404kwm8kl88x2pn5ku7q0ahysn"
	ownerAddr2 = "st1lc9sg3wq7guvkqv2d8vd2wycvj9fsspxq6qtg3"
	ownerAddr3 = "st15lm2e9h79j4d2zhyyf99j40uuy5vr404vurd0e"
	ownerAddr4 = "st1sysfc0hrjt63zqdywtzu7wu367uc23mcq7cz99"
	ownerAddr5 = "st1ax999nnc6z456axup4kq4wlxkyjwa363f8mglu"
	ownerAddr6 = "st1usjxfuzr0t70wny7rgtlr57ys3hsgsx0t3fl80"

	resNodePubKey1       = "stpub1addwnpepqwx64naamqf460ulupmyfnvy7arvvavr55as6q4h5my9tlcyx6fdu7q9mas"
	resNodeAddr1         = "st1vclaf9cdq4cd3fl4cnz4puhl7wp7hduagjaf82"
	resNodeInitialStake1 = 1000000000

	resNodePubKey2       = "stpub1addwnpepqgk5fz0hwwxsp4eed4yfqfywd46fccp84v0pc30lvvzlkgfhfezvsg3cla0"
	resNodeAddr2         = "st16qqwwt8uw3xj8qrcj3qm6r52zm45jmzat3mlrg"
	resNodeInitialStake2 = 1000000000

	resNodePubKey3       = "stpub1addwnpepqwhzesvudsfz93q22s7n2fecsj578z3ddl0dlf579dr6ffdma55mstuqr4z"
	resNodeAddr3         = "st1dta9wqhhmlpjen30dtcfjvsnf4qjwu6l572nus"
	resNodeInitialStake3 = 1000000000

	idxNodePubKey1       = "stpub1addwnpepq0t9el66pkwr5rspd0daq44m7a755u8sn0vrynjhqma3chpu6nw3kxellss"
	idxNodeAddr1         = "st1srwkzl2z5934ph69dzwqrcd7xytr8p3q3rru3r"
	idxNodeInitialStake1 = 1000000000

	idxNodePubKey2       = "stpub1addwnpepqf8p0s5m2dpqq49nf65nzm54fagr3j6zxzl6u507r0lfkf82uhpng7wt7uz"
	idxNodeAddr2         = "st1az73m086m9knql092w4ezpqeer9zs4ec050wnd"
	idxNodeInitialStake2 = 1000000000

	idxNodePubKey3       = "stpub1addwnpepq26rc94qwawkjckuh7new6uswpwnl527hca66xsmex55g9cnmluhuemnc8l"
	idxNodeAddr3         = "st1m80hu9d3qd0jmw5gef5ntgzcz9cps6jysye9nn"
	idxNodeInitialStake3 = 1000000000

	valInitialStake = 1000000000

	resourceNodeVolume1 = 10000000
	resourceNodeVolume2 = 20000000
	resourceNodeVolume3 = 30000000
	epoch1              = 1

	totalUnissuedPrePay = 10000000000
	remainingOzoneLimit = 10000000000
)

var (
	valOpPk1    = ed25519.GenPrivKey().PubKey()
	valOpAddr1  = sdk.ValAddress(valOpPk1.Address())
	valAccAddr1 = sdk.AccAddress(valOpPk1.Address())

	valConsPk1 = ed25519.GenPrivKey().PubKey()
)

func Test(t *testing.T) {

	//prepare keepers
	ctx, accountKeeper, bankKeeper, k, stakingKeeper, _, _, registerKeeper := CreateTestInput(t, false)
	bondDenom := k.BondDenom(ctx)

	// create validator with 50% commission
	stakingHandler := staking.NewHandler(stakingKeeper)
	createAccount(t, ctx, accountKeeper, bankKeeper, valAccAddr1, sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.NewInt(valInitialStake))))
	commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	msgVal := staking.NewMsgCreateValidator(
		valOpAddr1, valConsPk1,
		sdk.NewCoin(bondDenom, sdk.NewInt(valInitialStake)), staking.Description{}, commission, sdk.OneInt(),
	)
	res, err := stakingHandler(ctx, msgVal)
	require.NoError(t, err)
	require.NotNil(t, res)

	//initial genesis stake total value
	initialGenesisStakeTotal := sdk.NewInt(resNodeInitialStake1).Add(sdk.NewInt(resNodeInitialStake2)).Add(sdk.NewInt(resNodeInitialStake3)).
		Add(sdk.NewInt(idxNodeInitialStake1)).Add(sdk.NewInt(idxNodeInitialStake2)).Add(sdk.NewInt(idxNodeInitialStake3))
	registerKeeper.SetInitialGenesisStakeTotal(ctx, initialGenesisStakeTotal)

	//createAccount(t, ctx, accountKeeper, bankKeeper, nodeAddr, sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.ZeroInt())))
	foundationAcc, err := sdk.AccAddressFromBech32(foundationAccAddr)
	require.NoError(t, err)
	createAccount(t, ctx, accountKeeper, bankKeeper, foundationAcc, sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.ZeroInt())))
	k.SetFoundationAccount(ctx, foundationAcc)
	bankKeeper.AddCoins(ctx, foundationAcc, sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), sdk.NewInt(foundationDeposit))))

	registerHandler := register.NewHandler(registerKeeper)

	msgRes1, resAddr1 := getResNodeRegisterMsg(t, ctx, accountKeeper, bankKeeper, ownerAddr1, resNodePubKey1, resNodeAddr1, "sds://resourceNode1", bondDenom, resNodeInitialStake1)
	msgRes2, resAddr2 := getResNodeRegisterMsg(t, ctx, accountKeeper, bankKeeper, ownerAddr2, resNodePubKey2, resNodeAddr2, "sds://resourceNode2", bondDenom, resNodeInitialStake2)
	msgRes3, resAddr3 := getResNodeRegisterMsg(t, ctx, accountKeeper, bankKeeper, ownerAddr3, resNodePubKey3, resNodeAddr3, "sds://resourceNode3", bondDenom, resNodeInitialStake3)
	msgIdx1, idxAddr1 := getIdxNodeRegisterMsg(t, ctx, accountKeeper, bankKeeper, ownerAddr4, idxNodePubKey1, idxNodeAddr1, "sds://indexingNode1", bondDenom, idxNodeInitialStake1)
	msgIdx2, idxAddr2 := getIdxNodeRegisterMsg(t, ctx, accountKeeper, bankKeeper, ownerAddr5, idxNodePubKey2, idxNodeAddr2, "sds://indexingNode2", bondDenom, idxNodeInitialStake2)
	msgIdx3, idxAddr3 := getIdxNodeRegisterMsg(t, ctx, accountKeeper, bankKeeper, ownerAddr6, idxNodePubKey3, idxNodeAddr3, "sds://indexingNode3", bondDenom, idxNodeInitialStake3)

	//register resource node1
	res, err = registerHandler(ctx, msgRes1)
	require.NoError(t, err)
	require.NotNil(t, res)
	//register resource node2
	res, err = registerHandler(ctx, msgRes2)
	require.NoError(t, err)
	require.NotNil(t, res)
	//register resource node3
	res, err = registerHandler(ctx, msgRes3)
	require.NoError(t, err)
	require.NotNil(t, res)
	//register indexing node1
	res, err = registerHandler(ctx, msgIdx1)
	require.NoError(t, err)
	require.NotNil(t, res)
	//register indexing node2
	res, err = registerHandler(ctx, msgIdx2)
	require.NoError(t, err)
	require.NotNil(t, res)
	//register indexing node3
	res, err = registerHandler(ctx, msgIdx3)
	require.NoError(t, err)
	require.NotNil(t, res)

	//build traffic list
	var trafficList []types.SingleNodeVolume
	trafficList = append(trafficList, types.NewSingleNodeVolume(resAddr1, sdk.NewInt(resourceNodeVolume1)))
	trafficList = append(trafficList, types.NewSingleNodeVolume(resAddr2, sdk.NewInt(resourceNodeVolume2)))
	trafficList = append(trafficList, types.NewSingleNodeVolume(resAddr3, sdk.NewInt(resourceNodeVolume3)))

	//PrePay
	k.setTotalUnissuedPrepay(ctx, sdk.NewInt(totalUnissuedPrePay))
	//remaining ozone limit
	registerKeeper.SetRemainingOzoneLimit(ctx, sdk.NewInt(remainingOzoneLimit))

	//check prepared data
	S := k.registerKeeper.GetInitialGenesisStakeTotal(ctx).ToDec()
	fmt.Println("S=" + S.String())
	Pt := k.getTotalUnissuedPrepay(ctx).ToDec()
	fmt.Println("Pt=" + Pt.String())
	Y := k.getTotalConsumedOzone(trafficList).ToDec()
	fmt.Println("Y=" + Y.String())
	Lt := k.registerKeeper.GetRemainingOzoneLimit(ctx).ToDec()
	fmt.Println("Lt=" + Lt.String())
	R := S.Add(Pt).Mul(Y).Quo(Lt.Add(Y))
	fmt.Println("R=" + R.String())

	k.DistributePotReward(ctx, trafficList, sdk.NewInt(epoch1))

	idvRwdResNode1Ep1 := k.getIndividualReward(ctx, resAddr1, epoch1)
	fmt.Println("idvRwdResNode1Ep1=" + idvRwdResNode1Ep1.String())
	matureTotalResNode1 := k.getMatureTotalReward(ctx, resAddr1)
	fmt.Println("matureTotalResNode1=" + matureTotalResNode1.String())
	immatureTotalResNode1 := k.getImmatureTotalReward(ctx, resAddr1)
	fmt.Println("immatureTotalResNode1=" + immatureTotalResNode1.String())

	idvRwdResNode2Ep1 := k.getIndividualReward(ctx, resAddr2, epoch1)
	fmt.Println("idvRwdResNode2Ep1=" + idvRwdResNode2Ep1.String())
	matureTotalResNode2 := k.getMatureTotalReward(ctx, resAddr2)
	fmt.Println("matureTotalResNode2=" + matureTotalResNode2.String())
	immatureTotalResNode2 := k.getImmatureTotalReward(ctx, resAddr2)
	fmt.Println("immatureTotalResNode2=" + immatureTotalResNode2.String())

	idvRwdResNode3Ep1 := k.getIndividualReward(ctx, resAddr3, epoch1)
	fmt.Println("idvRwdResNode3Ep1=" + idvRwdResNode3Ep1.String())
	matureTotalResNode3 := k.getMatureTotalReward(ctx, resAddr3)
	fmt.Println("matureTotalResNode3=" + matureTotalResNode3.String())
	immatureTotalResNode3 := k.getImmatureTotalReward(ctx, resAddr3)
	fmt.Println("immatureTotalResNode3=" + immatureTotalResNode3.String())

	idvRwdIdxNode1Ep1 := k.getIndividualReward(ctx, idxAddr1, epoch1)
	fmt.Println("idvRwdIdxNode1Ep1=" + idvRwdIdxNode1Ep1.String())
	matureTotalIdxNode1 := k.getMatureTotalReward(ctx, idxAddr1)
	fmt.Println("matureTotalIdxNode1=" + matureTotalIdxNode1.String())
	immatureTotalIdxNode1 := k.getImmatureTotalReward(ctx, idxAddr1)
	fmt.Println("immatureTotalIdxNode1=" + immatureTotalIdxNode1.String())

	idvRwdIdxNode2Ep1 := k.getIndividualReward(ctx, idxAddr2, epoch1)
	fmt.Println("idvRwdIdxNode2Ep1=" + idvRwdIdxNode2Ep1.String())
	matureTotalIdxNode2 := k.getMatureTotalReward(ctx, idxAddr2)
	fmt.Println("matureTotalIdxNode2=" + matureTotalIdxNode2.String())
	immatureTotalIdxNode2 := k.getImmatureTotalReward(ctx, idxAddr2)
	fmt.Println("immatureTotalIdxNode2=" + immatureTotalIdxNode2.String())

	idvRwdIdxNode3Ep1 := k.getIndividualReward(ctx, idxAddr3, epoch1)
	fmt.Println("idvRwdIdxNode3Ep1=" + idvRwdIdxNode3Ep1.String())
	matureTotalIdxNode3 := k.getMatureTotalReward(ctx, idxAddr3)
	fmt.Println("matureTotalIdxNode3=" + matureTotalIdxNode3.String())
	immatureTotalIdxNode3 := k.getImmatureTotalReward(ctx, idxAddr3)
	fmt.Println("immatureTotalIdxNode3=" + immatureTotalIdxNode3.String())
}

func getResNodeRegisterMsg(t *testing.T, ctx sdk.Context, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper,
	ownerAddrStr string, pubKeyStr string, nodeAddrStr string, networkId string, bondDenom string, stake int64) (
	msg register.MsgCreateResourceNode, nodeAddr sdk.AccAddress) {

	nodeAddr, err := sdk.AccAddressFromBech32(nodeAddrStr)
	require.NoError(t, err)
	owner, err := sdk.AccAddressFromBech32(ownerAddrStr)
	require.NoError(t, err)

	//create owner account
	createAccount(t, ctx, accountKeeper, bankKeeper, owner, sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.NewInt(stake))))
	//create node account
	createAccount(t, ctx, accountKeeper, bankKeeper, nodeAddr, sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.ZeroInt())))

	pubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, pubKeyStr)
	require.NoError(t, err)
	msg = register.NewMsgCreateResourceNode(networkId, pubKey, sdk.NewCoin(bondDenom, sdk.NewInt(stake)), owner,
		register.NewDescription(networkId, "", "", "", ""), "4")

	return msg, nodeAddr
}

func getIdxNodeRegisterMsg(t *testing.T, ctx sdk.Context, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper,
	ownerAddrStr string, pubKeyStr string, nodeAddrStr string, networkId string, bondDenom string, stake int64) (
	msg register.MsgCreateIndexingNode, nodeAddr sdk.AccAddress) {

	nodeAddr, err := sdk.AccAddressFromBech32(nodeAddrStr)
	require.NoError(t, err)
	owner, err := sdk.AccAddressFromBech32(ownerAddrStr)
	require.NoError(t, err)

	//create owner account
	createAccount(t, ctx, accountKeeper, bankKeeper, owner, sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.NewInt(stake))))
	//create node account
	createAccount(t, ctx, accountKeeper, bankKeeper, nodeAddr, sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.ZeroInt())))

	pubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, pubKeyStr)
	require.NoError(t, err)
	msg = register.NewMsgCreateIndexingNode(networkId, pubKey, sdk.NewCoin(bondDenom, sdk.NewInt(stake)), owner,
		register.NewDescription(networkId, "", "", "", ""))

	return msg, nodeAddr
}

func createAccount(t *testing.T, ctx sdk.Context, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, acc sdk.AccAddress, coins sdk.Coins) {
	account := accountKeeper.GetAccount(ctx, acc)
	if account == nil {
		account = accountKeeper.NewAccountWithAddress(ctx, acc)
		fmt.Printf("create account: " + account.String() + "\n")
	}
	coins, err := bankKeeper.AddCoins(ctx, acc, coins)
	require.NoError(t, err)
}
