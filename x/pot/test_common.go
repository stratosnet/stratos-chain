package pot

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

const (
	chainID              = ""
	AccountAddressPrefix = "st"

	stos2ustos = 1000000000
)

var (
	AccountPubKeyPrefix    = AccountAddressPrefix + "pub"
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	ValidatorPubKeyPrefix  = AccountAddressPrefix + "valoperpub"
	ConsNodeAddressPrefix  = AccountAddressPrefix + "valcons"
	ConsNodePubKeyPrefix   = AccountAddressPrefix + "valconspub"

	resourceNodeVolume1 = sdk.NewInt(500000000000)
	resourceNodeVolume2 = sdk.NewInt(300000000000)
	resourceNodeVolume3 = sdk.NewInt(200000000000)

	depositForSendingTx, _ = sdk.NewIntFromString("100000000000000000000000000000")
	totalUnissuedPrepay, _ = sdk.NewIntFromString("100000000000000000")
	remainingOzoneLimit, _ = sdk.NewIntFromString("500000000000000000000")
	initialOzonePrice      = sdk.NewInt(10000000000)
	foundationAccAddr      = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	foundationDeposit      = sdk.NewCoins(sdk.NewCoin("ustos", sdk.NewInt(40000000000000000)))

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
	resNodeInitialStake1 = sdk.NewInt(3 * stos2ustos)

	resNodePubKey2       = secp256k1.GenPrivKey().PubKey()
	resNodeAddr2         = sdk.AccAddress(resNodePubKey2.Address())
	resNodeInitialStake2 = sdk.NewInt(3 * stos2ustos)

	resNodePubKey3       = secp256k1.GenPrivKey().PubKey()
	resNodeAddr3         = sdk.AccAddress(resNodePubKey3.Address())
	resNodeInitialStake3 = sdk.NewInt(3 * stos2ustos)

	resNodePubKey4       = secp256k1.GenPrivKey().PubKey()
	resNodeAddr4         = sdk.AccAddress(resNodePubKey4.Address())
	resNodeInitialStake4 = sdk.NewInt(3 * stos2ustos)

	resNodePubKey5       = secp256k1.GenPrivKey().PubKey()
	resNodeAddr5         = sdk.AccAddress(resNodePubKey5.Address())
	resNodeInitialStake5 = sdk.NewInt(3 * stos2ustos)

	idxNodePrivKey1      = secp256k1.GenPrivKey()
	idxNodePubKey1       = idxNodePrivKey1.PubKey()
	idxNodeAddr1         = sdk.AccAddress(idxNodePubKey1.Address())
	idxNodeInitialStake1 = sdk.NewInt(5 * stos2ustos)

	idxNodePubKey2       = secp256k1.GenPrivKey().PubKey()
	idxNodeAddr2         = sdk.AccAddress(idxNodePubKey2.Address())
	idxNodeInitialStake2 = sdk.NewInt(5 * stos2ustos)

	idxNodePubKey3       = secp256k1.GenPrivKey().PubKey()
	idxNodeAddr3         = sdk.AccAddress(idxNodePubKey3.Address())
	idxNodeInitialStake3 = sdk.NewInt(5 * stos2ustos)

	valOpPrivKey1 = secp256k1.GenPrivKey()
	valOpPubKey1  = valOpPrivKey1.PubKey()
	valOpValAddr1 = sdk.ValAddress(valOpPubKey1.Address())
	valOpAccAddr1 = sdk.AccAddress(valOpPubKey1.Address())

	valConsPrivKey1 = secp256k1.GenPrivKey()
	valConsPubk1    = valConsPrivKey1.PubKey()
	valInitialStake = sdk.NewInt(15 * stos2ustos)
)

func SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}

func setupAccounts(mApp *mock.App) []authexported.Account {

	//************************** setup resource nodes owners' accounts **************************
	resOwnerAcc1 := &auth.BaseAccount{
		Address: resOwner1,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", resNodeInitialStake1.Add(depositForSendingTx))},
	}
	resOwnerAcc2 := &auth.BaseAccount{
		Address: resOwner2,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", resNodeInitialStake2)},
	}
	resOwnerAcc3 := &auth.BaseAccount{
		Address: resOwner3,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", resNodeInitialStake3)},
	}
	resOwnerAcc4 := &auth.BaseAccount{
		Address: resOwner4,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", resNodeInitialStake4)},
	}
	resOwnerAcc5 := &auth.BaseAccount{
		Address: resOwner5,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", resNodeInitialStake5)},
	}

	//************************** setup indexing nodes owners' accounts **************************
	idxOwnerAcc1 := &auth.BaseAccount{
		Address: idxOwner1,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", idxNodeInitialStake1)},
	}
	idxOwnerAcc2 := &auth.BaseAccount{
		Address: idxOwner2,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", idxNodeInitialStake2)},
	}
	idxOwnerAcc3 := &auth.BaseAccount{
		Address: idxOwner3,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", idxNodeInitialStake3)},
	}

	//************************** setup validator delegators' accounts **************************
	valOwnerAcc1 := &auth.BaseAccount{
		Address: valOpAccAddr1,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", valInitialStake)},
	}

	//************************** setup indexing nodes' accounts **************************
	idxNodeAcc1 := &auth.BaseAccount{
		Address: idxNodeAddr1,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", sdk.ZeroInt())},
	}

	foundationAcc := &auth.BaseAccount{
		Address: foundationAccAddr,
		Coins:   foundationDeposit,
	}

	accs := []authexported.Account{
		resOwnerAcc1, resOwnerAcc2, resOwnerAcc3, resOwnerAcc4, resOwnerAcc5,
		idxOwnerAcc1, idxOwnerAcc2, idxOwnerAcc3,
		valOwnerAcc1,
		idxNodeAcc1,
		foundationAcc,
	}

	ctx1 := mApp.BaseApp.NewContext(true, abci.Header{})
	ctx1.Logger().Info("idxNodeAcc1 -> " + idxNodeAcc1.String())
	ctx1.Logger().Info("foundationAcc -> " + foundationAcc.String())

	return accs
}

func setupAllResourceNodes() []register.ResourceNode {
	resourceNode1 := register.NewResourceNode("sds://resourceNode1", resNodePubKey1, resOwner1, register.NewDescription("sds://resourceNode1", "", "", "", ""), "4")
	resourceNode2 := register.NewResourceNode("sds://resourceNode2", resNodePubKey2, resOwner2, register.NewDescription("sds://resourceNode2", "", "", "", ""), "4")
	resourceNode3 := register.NewResourceNode("sds://resourceNode3", resNodePubKey3, resOwner3, register.NewDescription("sds://resourceNode3", "", "", "", ""), "4")
	resourceNode4 := register.NewResourceNode("sds://resourceNode4", resNodePubKey4, resOwner4, register.NewDescription("sds://resourceNode4", "", "", "", ""), "4")
	resourceNode5 := register.NewResourceNode("sds://resourceNode5", resNodePubKey5, resOwner5, register.NewDescription("sds://resourceNode5", "", "", "", ""), "4")

	resourceNode1 = resourceNode1.AddToken(resNodeInitialStake1)
	resourceNode2 = resourceNode2.AddToken(resNodeInitialStake2)
	resourceNode3 = resourceNode3.AddToken(resNodeInitialStake3)
	resourceNode4 = resourceNode4.AddToken(resNodeInitialStake4)
	resourceNode5 = resourceNode5.AddToken(resNodeInitialStake5)

	resourceNode1.Status = sdk.Bonded
	resourceNode2.Status = sdk.Bonded
	resourceNode3.Status = sdk.Bonded
	resourceNode4.Status = sdk.Bonded
	resourceNode5.Status = sdk.Bonded

	var resourceNodes []register.ResourceNode
	resourceNodes = append(resourceNodes, resourceNode1)
	resourceNodes = append(resourceNodes, resourceNode2)
	resourceNodes = append(resourceNodes, resourceNode3)
	resourceNodes = append(resourceNodes, resourceNode4)
	resourceNodes = append(resourceNodes, resourceNode5)
	return resourceNodes
}

func setupAllIndexingNodes() []register.IndexingNode {
	var indexingNodes []register.IndexingNode
	indexingNode1 := register.NewIndexingNode("sds://indexingNode1", idxNodePubKey1, idxOwner1, register.NewDescription("sds://indexingNode1", "", "", "", ""))
	indexingNode2 := register.NewIndexingNode("sds://indexingNode2", idxNodePubKey2, idxOwner2, register.NewDescription("sds://indexingNode2", "", "", "", ""))
	indexingNode3 := register.NewIndexingNode("sds://indexingNode3", idxNodePubKey3, idxOwner3, register.NewDescription("sds://indexingNode3", "", "", "", ""))

	indexingNode1 = indexingNode1.AddToken(idxNodeInitialStake1)
	indexingNode2 = indexingNode2.AddToken(idxNodeInitialStake2)
	indexingNode3 = indexingNode3.AddToken(idxNodeInitialStake3)

	indexingNode1.Status = sdk.Bonded
	indexingNode2.Status = sdk.Bonded
	indexingNode3.Status = sdk.Bonded

	indexingNodes = append(indexingNodes, indexingNode1)
	indexingNodes = append(indexingNodes, indexingNode2)
	indexingNodes = append(indexingNodes, indexingNode3)

	return indexingNodes

}

// SignCheckDeliver checks a generated signed transaction and simulates a
// block commitment with the given transaction. A test assertion is made using
// the parameter 'expPass' against the result. A corresponding result is
// returned.
func SignCheckDeliver(
	t *testing.T, cdc *codec.Codec, app *baseapp.BaseApp, header abci.Header, msgs []sdk.Msg,
	accNums, seq []uint64, expSimPass, expPass bool, priv ...crypto.PrivKey,
) (sdk.GasInfo, *sdk.Result, error) {

	tx := GenTx(msgs, accNums, seq, priv...)

	txBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)
	require.Nil(t, err)

	// Must simulate now as CheckTx doesn't run Msgs anymore
	_, res, err := app.Simulate(txBytes, tx)

	if expSimPass {
		require.NoError(t, err)
		require.NotNil(t, res)
	} else {
		require.Error(t, err)
		require.Nil(t, res)
	}

	// Simulate a sending a transaction and committing a block
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	gInfo, res, err := app.Deliver(tx)

	if expPass {
		require.NoError(t, err)
		require.NotNil(t, res)
	} else {
		require.Error(t, err)
		require.Nil(t, res)
	}

	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()

	return gInfo, res, err
}

// GenTx generates a signed mock transaction.
func GenTx(msgs []sdk.Msg, accnums []uint64, seq []uint64, priv ...crypto.PrivKey) auth.StdTx {
	// Make the transaction free
	fee := auth.StdFee{
		Amount: sdk.NewCoins(sdk.NewInt64Coin("foocoin", 0)),
		Gas:    300000,
	}

	sigs := make([]auth.StdSignature, len(priv))
	memo := "testmemotestmemo"

	for i, p := range priv {
		sig, err := p.Sign(auth.StdSignBytes(chainID, accnums[i], seq[i], fee, msgs, memo))
		if err != nil {
			panic(err)
		}

		sigs[i] = auth.StdSignature{
			PubKey:    p.PubKey(),
			Signature: sig,
		}
	}

	return auth.NewStdTx(msgs, fee, sigs, memo)
}
