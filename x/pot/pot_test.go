package pot_test

import (
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"

	//"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	stratos "github.com/stratosnet/stratos-chain/types"
)

const (
	chainID    = "testchain_1-1"
	stos2ustos = 1000000000
)

var (
	resNodeSlashingUOZAmt1 = sdk.NewInt(1000000000000000000)

	resourceNodeVolume1 = sdk.NewInt(500000)
	resourceNodeVolume2 = sdk.NewInt(300000)
	resourceNodeVolume3 = sdk.NewInt(200000)

	depositForSendingTx, _    = sdk.NewIntFromString("100000000000000000000000000000")
	totalUnissuedPrepayVal, _ = sdk.NewIntFromString("1000000000000")
	totalUnissuedPrepay       = sdk.NewCoin("ustos", totalUnissuedPrepayVal)
	initialUOzonePrice        = sdk.NewDecWithPrec(10000000, 9) // 0.001 ustos -> 1 uoz

	foundationDepositorPrivKey = secp256k1.GenPrivKey()
	foundationDepositorAccAddr = sdk.AccAddress(foundationDepositorPrivKey.PubKey().Address())
	foundationDeposit          = sdk.NewCoins(sdk.NewCoin("utros", sdk.NewInt(40000000000000000)))

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
	resNodeInitialStake1 = sdk.NewInt(3 * stos2ustos)

	resNodePubKey2       = secp256k1.GenPrivKey().PubKey()
	resNodeAddr2         = sdk.AccAddress(resNodePubKey2.Address())
	resNodeNetworkId2    = stratos.SdsAddress(resNodePubKey2.Address())
	resNodeInitialStake2 = sdk.NewInt(3 * stos2ustos)

	resNodePubKey3       = secp256k1.GenPrivKey().PubKey()
	resNodeAddr3         = sdk.AccAddress(resNodePubKey3.Address())
	resNodeNetworkId3    = stratos.SdsAddress(resNodePubKey3.Address())
	resNodeInitialStake3 = sdk.NewInt(3 * stos2ustos)

	resNodePubKey4       = secp256k1.GenPrivKey().PubKey()
	resNodeAddr4         = sdk.AccAddress(resNodePubKey4.Address())
	resNodeNetworkId4    = stratos.SdsAddress(resNodePubKey4.Address())
	resNodeInitialStake4 = sdk.NewInt(3 * stos2ustos)

	resNodePubKey5       = secp256k1.GenPrivKey().PubKey()
	resNodeAddr5         = sdk.AccAddress(resNodePubKey5.Address())
	resNodeNetworkId5    = stratos.SdsAddress(resNodePubKey5.Address())
	resNodeInitialStake5 = sdk.NewInt(3 * stos2ustos)

	idxNodePrivKey1      = secp256k1.GenPrivKey()
	idxNodePubKey1       = idxNodePrivKey1.PubKey()
	idxNodeAddr1         = sdk.AccAddress(idxNodePubKey1.Address())
	idxNodeNetworkId1    = stratos.SdsAddress(idxNodePubKey1.Address())
	idxNodeInitialStake1 = sdk.NewInt(5 * stos2ustos)

	idxNodePubKey2       = secp256k1.GenPrivKey().PubKey()
	idxNodeAddr2         = sdk.AccAddress(idxNodePubKey2.Address())
	idxNodeNetworkId2    = stratos.SdsAddress(idxNodePubKey2.Address())
	idxNodeInitialStake2 = sdk.NewInt(5 * stos2ustos)

	idxNodePubKey3       = secp256k1.GenPrivKey().PubKey()
	idxNodeAddr3         = sdk.AccAddress(idxNodePubKey3.Address())
	idxNodeNetworkId3    = stratos.SdsAddress(idxNodePubKey3.Address())
	idxNodeInitialStake3 = sdk.NewInt(5 * stos2ustos)

	valOpPrivKey1 = secp256k1.GenPrivKey()
	valOpPubKey1  = valOpPrivKey1.PubKey()
	valOpValAddr1 = sdk.ValAddress(valOpPubKey1.Address())
	valOpAccAddr1 = sdk.AccAddress(valOpPubKey1.Address())

	valConsPrivKey1 = ed25519.GenPrivKey()
	valConsPubk1    = valConsPrivKey1.PubKey()
	valInitialStake = sdk.NewInt(15 * stos2ustos)
)

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
			Coins:   sdk.Coins{sdk.NewCoin("ustos", resNodeInitialStake1.Add(depositForSendingTx))},
		},
		{
			Address: resOwner2.String(),
			Coins:   sdk.Coins{sdk.NewCoin("ustos", resNodeInitialStake2)},
		},
		{
			Address: resOwner3.String(),
			Coins:   sdk.Coins{sdk.NewCoin("ustos", resNodeInitialStake3)},
		},
		{
			Address: resOwner4.String(),
			Coins:   sdk.Coins{sdk.NewCoin("ustos", resNodeInitialStake4)},
		},
		{
			Address: resOwner5.String(),
			Coins:   sdk.Coins{sdk.NewCoin("ustos", resNodeInitialStake5)},
		},
		{
			Address: idxOwner1.String(),
			Coins:   sdk.Coins{sdk.NewCoin("ustos", idxNodeInitialStake1)},
		},
		{
			Address: idxOwner2.String(),
			Coins:   sdk.Coins{sdk.NewCoin("ustos", idxNodeInitialStake2)},
		},
		{
			Address: idxOwner3.String(),
			Coins:   sdk.Coins{sdk.NewCoin("ustos", idxNodeInitialStake3)},
		},
		{
			Address: valOpAccAddr1.String(),
			Coins:   sdk.Coins{sdk.NewCoin("ustos", valInitialStake)},
		},
		{
			Address: idxNodeAddr1.String(),
			Coins:   sdk.Coins{sdk.NewCoin("ustos", sdk.ZeroInt())},
		},
		{
			Address: foundationDepositorAccAddr.String(),
			Coins:   foundationDeposit,
		},
	}

	//app := simapp.SetupWithGenesisAccounts(accs, balances...)
	//simapp.CheckBalance(t, app, addr1, sdk.Coins{genCoin})
	//simapp.CheckBalance(t, app, addr2, sdk.Coins{genCoin})

	//ctx1 := mApp.BaseApp.NewContext(true, abci.Header{})
	//ctx1.Logger().Info("idxNodeAcc1 -> " + idxNodeAcc1.String())

	return accs, balances
}

func setupAllResourceNodes() []registertypes.ResourceNode {

	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	nodeType := registertypes.STORAGE
	resourceNode1, _ := registertypes.NewResourceNode(resNodeNetworkId1, resNodePubKey1, resOwner1, registertypes.NewDescription("sds://resourceNode1", "", "", "", ""), &nodeType, time)
	resourceNode2, _ := registertypes.NewResourceNode(resNodeNetworkId2, resNodePubKey2, resOwner2, registertypes.NewDescription("sds://resourceNode2", "", "", "", ""), &nodeType, time)
	resourceNode3, _ := registertypes.NewResourceNode(resNodeNetworkId3, resNodePubKey3, resOwner3, registertypes.NewDescription("sds://resourceNode3", "", "", "", ""), &nodeType, time)
	resourceNode4, _ := registertypes.NewResourceNode(resNodeNetworkId4, resNodePubKey4, resOwner4, registertypes.NewDescription("sds://resourceNode4", "", "", "", ""), &nodeType, time)
	resourceNode5, _ := registertypes.NewResourceNode(resNodeNetworkId5, resNodePubKey5, resOwner5, registertypes.NewDescription("sds://resourceNode5", "", "", "", ""), &nodeType, time)

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

// SignCheckDeliver checks a generated signed transaction and simulates a
// block commitment with the given transaction. A test assertion is made using
// the parameter 'expPass' against the result. A corresponding result is
// returned.
//func SignCheckDeliver(
//	t *testing.T, cdc *codec.Codec, app *baseapp.BaseApp, header abci.Header, msgs []sdk.Msg,
//	accNums, seq []uint64, expSimPass, expPass bool, priv ...crypto.PrivKey,
//) (sdk.GasInfo, *sdk.Result, error) {
//
//	tx := GenTx(msgs, accNums, seq, priv...)
//
//	txBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)
//	require.Nil(t, err)
//
//	// Must simulate now as CheckTx doesn't run Msgs anymore
//	_, res, err := app.Simulate(txBytes, tx)
//
//	if expSimPass {
//		require.NoError(t, err)
//		require.NotNil(t, res)
//	} else {
//		require.Error(t, err)
//		require.Nil(t, res)
//	}
//
//	// Simulate a sending a transaction and committing a block
//	app.BeginBlock(abci.RequestBeginBlock{Header: header})
//	gInfo, res, err := app.Deliver(tx)
//
//	if expPass {
//		require.NoError(t, err)
//		require.NotNil(t, res)
//	} else {
//		require.Error(t, err)
//		require.Nil(t, res)
//	}
//
//	app.EndBlock(abci.RequestEndBlock{})
//	app.Commit()
//
//	return gInfo, res, err
//}

// GenTx generates a signed mock transaction.
//func GenTx(msgs []sdk.Msg, accnums []uint64, seq []uint64, priv ...crypto.PrivKey) auth.StdTx {
//	// Make the transaction free
//	fee := auth.StdFee{
//		Amount: sdk.NewCoins(sdk.NewInt64Coin("foocoin", 0)),
//		Gas:    5000000,
//	}
//
//	sigs := make([]auth.StdSignature, len(priv))
//	memo := "testmemotestmemo"
//
//	for i, p := range priv {
//		sig, err := p.Sign(auth.StdSignBytes(chainID, accnums[i], seq[i], fee, msgs, memo))
//		if err != nil {
//			panic(err)
//		}
//
//		sigs[i] = auth.StdSignature{
//			PubKey:    p.PubKey(),
//			Signature: sig,
//		}
//	}
//
//	return auth.NewStdTx(msgs, fee, sigs, memo)
//}
