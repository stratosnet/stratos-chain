package sds

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/mock"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	chainID             = ""
	StratosBech32Prefix = "st"
	DefaultDenom        = "ustos"
	stos2ustos          = 1000000000
)

var (
	testFileHashHex = "c03661732294feb49caf6dc16c7cbb2534986d73"

	AccountPubKeyPrefix    = StratosBech32Prefix + "pub"
	ValidatorAddressPrefix = StratosBech32Prefix + "valoper"
	ValidatorPubKeyPrefix  = StratosBech32Prefix + "valoperpub"
	ConsNodeAddressPrefix  = StratosBech32Prefix + "valcons"
	ConsNodePubKeyPrefix   = StratosBech32Prefix + "valconspub"
	SdsNodeP2PKeyPrefix    = StratosBech32Prefix + "sdsp2p"

	resourceNodeVolume1 = sdk.NewInt(500000000000)
	resourceNodeVolume2 = sdk.NewInt(300000000000)
	resourceNodeVolume3 = sdk.NewInt(200000000000)
	prepayAmt           = sdk.NewInt(2 * stos2ustos)

	depositForSendingTx, _             = sdk.NewIntFromString("100000000000000000000000000000")
	initialUOzonePrice                 = sdk.NewDecWithPrec(10000000, 9) // 0.001 ustos -> 1 uoz
	totalUnissuedPrepayVal, _          = sdk.NewIntFromString("100000000000000000")
	totalUnissuedPrepay                = sdk.NewCoin("ustos", totalUnissuedPrepayVal)
	remainingOzoneLimit, _             = sdk.NewIntFromString("500000000000000000000")
	totalUnissuedPrepayTestPurchase, _ = sdk.NewIntFromString("0")
	remainingOzoneLimitTestPurchase, _ = sdk.NewIntFromString("100000000000")
	initialUOzonePriceTestPurchase     = sdk.NewDecWithPrec(1000000, 9) // 0.001 ustos -> 1 uoz

	foundationDeposit = sdk.NewInt(40000000000000000)

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

	pubKeyRes1                   = secp256k1.GenPrivKey().PubKey()
	addrRes1                     = sdk.AccAddress(pubKeyRes1.Address())
	initialStakeRes1             = sdk.NewInt(3 * stos2ustos)
	initialStakeRes1TestPurchase = sdk.NewInt(100000000000)

	pubKeyRes2       = secp256k1.GenPrivKey().PubKey()
	addrRes2         = sdk.AccAddress(pubKeyRes2.Address())
	initialStakeRes2 = sdk.NewInt(3 * stos2ustos)

	pubKeyRes3       = secp256k1.GenPrivKey().PubKey()
	addrRes3         = sdk.AccAddress(pubKeyRes3.Address())
	initialStakeRes3 = sdk.NewInt(3 * stos2ustos)

	pubKeyRes4       = secp256k1.GenPrivKey().PubKey()
	addrRes4         = sdk.AccAddress(pubKeyRes4.Address())
	initialStakeRes4 = sdk.NewInt(3 * stos2ustos)

	pubKeyRes5       = secp256k1.GenPrivKey().PubKey()
	addrRes5         = sdk.AccAddress(pubKeyRes5.Address())
	initialStakeRes5 = sdk.NewInt(3 * stos2ustos)

	privKeyIdx1                  = secp256k1.GenPrivKey()
	pubKeyIdx1                   = privKeyIdx1.PubKey()
	addrIdx1                     = sdk.AccAddress(pubKeyIdx1.Address())
	initialStakeIdx1             = sdk.NewInt(5 * stos2ustos)
	initialStakeIdx1TestPurchase = sdk.NewInt(100 * stos2ustos)

	pubKeyIdx2       = secp256k1.GenPrivKey().PubKey()
	addrIdx2         = sdk.AccAddress(pubKeyIdx2.Address())
	initialStakeIdx2 = sdk.NewInt(5 * stos2ustos)

	pubKeyIdx3       = secp256k1.GenPrivKey().PubKey()
	addrIdx3         = sdk.AccAddress(pubKeyIdx3.Address())
	initialStakeIdx3 = sdk.NewInt(5 * stos2ustos)

	valOpPrivKey1 = secp256k1.GenPrivKey()
	valOpPubKey1  = valOpPrivKey1.PubKey()
	valOpValAddr1 = sdk.ValAddress(valOpPubKey1.Address())
	valOpAccAddr1 = sdk.AccAddress(valOpPubKey1.Address())

	valConsPrivKey1 = secp256k1.GenPrivKey()
	valConsPubk1    = valConsPrivKey1.PubKey()
	valInitialStake = sdk.NewInt(15 * stos2ustos)

	// accs for sds module
	sdsAccPrivKey1      = secp256k1.GenPrivKey()
	sdsAccPubKey1       = sdsAccPrivKey1.PubKey()
	sdsAccAddr1         = sdk.AccAddress(sdsAccPubKey1.Address())
	sdsAccBal1          = sdk.NewInt(100 * stos2ustos)
	initialStakeSdsIdx1 = sdk.NewInt(5 * stos2ustos)

	sdsAccPrivKey2 = secp256k1.GenPrivKey()
	sdsAccPubKey2  = sdsAccPrivKey2.PubKey()
	sdsAccAddr2    = sdk.AccAddress(sdsAccPubKey2.Address())
	sdsAccBal2     = sdk.NewInt(100 * stos2ustos)

	sdsAccPrivKey3 = secp256k1.GenPrivKey()
	sdsAccPubKey3  = sdsAccPrivKey3.PubKey()
	sdsAccAddr3    = sdk.AccAddress(sdsAccPubKey3.Address())
	sdsAccBal3     = sdk.NewInt(100 * stos2ustos)

	// sp node used in sds module
	spNodePrivKeyIdx1      = secp256k1.GenPrivKey()
	spNodePubKeyIdx1       = spNodePrivKeyIdx1.PubKey()
	spNodeAddrIdx1         = sdk.AccAddress(spNodePubKeyIdx1.Address())
	spNodeInitialStakeIdx1 = sdk.NewInt(5 * stos2ustos)
)

func SetConfig() {
	config := stratos.GetConfig()
	config.SetBech32PrefixForAccount(StratosBech32Prefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.SetBech32PrefixForSdsNodeP2P(SdsNodeP2PKeyPrefix)
}

func setupAccounts(mApp *mock.App) []authexported.Account {

	str, _ := stratos.Bech32ifyPubKey(stratos.Bech32PubKeyTypeSdsP2PPub, sdsAccPubKey1)
	fmt.Println("sdsAccPubKey1=" + str)

	//************************** setup resource nodes owners' accounts **************************
	resOwnerAcc1 := &auth.BaseAccount{
		Address: resOwner1,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, sdk.NewInt(10000000000000000))},
	}
	resOwnerAcc2 := &auth.BaseAccount{
		Address: resOwner2,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, sdk.NewInt(10000000000000000))},
	}
	resOwnerAcc3 := &auth.BaseAccount{
		Address: resOwner3,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, sdk.NewInt(10000000000000000))},
	}
	resOwnerAcc4 := &auth.BaseAccount{
		Address: resOwner4,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, sdk.NewInt(10000000000000000))},
	}
	resOwnerAcc5 := &auth.BaseAccount{
		Address: resOwner5,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, sdk.NewInt(10000000000000000))},
	}

	//************************** setup indexing nodes owners' accounts **************************
	idxOwnerAcc1 := &auth.BaseAccount{
		Address: idxOwner1,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, sdk.ZeroInt())},
	}
	idxOwnerAcc2 := &auth.BaseAccount{
		Address: idxOwner2,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, sdk.ZeroInt())},
	}
	idxOwnerAcc3 := &auth.BaseAccount{
		Address: idxOwner3,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, sdk.ZeroInt())},
	}

	//************************** setup validator delegators' accounts **************************
	valOwnerAcc1 := &auth.BaseAccount{
		Address: valOpAccAddr1,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, valInitialStake)},
	}

	//************************** setup resource nodes' accounts **************************
	resNodeAcc1 := &auth.BaseAccount{
		Address: addrRes1,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, initialStakeRes1)},
	}
	resNodeAcc2 := &auth.BaseAccount{
		Address: addrRes2,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, initialStakeRes2)},
	}
	resNodeAcc3 := &auth.BaseAccount{
		Address: addrRes3,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, initialStakeRes3)},
	}
	resNodeAcc4 := &auth.BaseAccount{
		Address: addrRes4,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, initialStakeRes4)},
	}
	resNodeAcc5 := &auth.BaseAccount{
		Address: addrRes5,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, initialStakeRes5)},
	}

	//************************** setup indexing nodes' accounts **************************
	idxNodeAcc1 := &auth.BaseAccount{
		Address: addrIdx1,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, initialStakeIdx1.Add(depositForSendingTx))},
	}
	idxNodeAcc2 := &auth.BaseAccount{
		Address: addrIdx2,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, initialStakeIdx2)},
	}
	idxNodeAcc3 := &auth.BaseAccount{
		Address: addrIdx3,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, initialStakeIdx3)},
	}
	spNodeIdxNodeAcc1 := &auth.BaseAccount{
		Address: spNodeAddrIdx1,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, spNodeInitialStakeIdx1)},
	}

	//************************** setup sds module's accounts **************************
	sdsAcc1 := &auth.BaseAccount{ // sp node owner
		Address: sdsAccAddr1,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, sdsAccBal1)},
	}
	sdsAcc2 := &auth.BaseAccount{
		Address: sdsAccAddr2,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, sdsAccBal2)},
	}
	sdsAcc3 := &auth.BaseAccount{
		Address: sdsAccAddr3,
		Coins:   sdk.Coins{sdk.NewCoin(DefaultDenom, sdsAccBal3)},
	}

	// the sequence of the account list is related to the value of parameter "accNums" of mock.SignCheckDeliver() method
	accs := []authexported.Account{
		resOwnerAcc1, resOwnerAcc2, resOwnerAcc3, resOwnerAcc4, resOwnerAcc5,
		idxOwnerAcc1, idxOwnerAcc2, idxOwnerAcc3,
		valOwnerAcc1,
		resNodeAcc1, resNodeAcc2, resNodeAcc3, resNodeAcc4, resNodeAcc5,
		idxNodeAcc1, idxNodeAcc2, idxNodeAcc3, spNodeIdxNodeAcc1,
		sdsAcc1, sdsAcc2, sdsAcc3,
	}

	ctx1 := mApp.BaseApp.NewContext(true, abci.Header{})
	ctx1.Logger().Info("resNodeAcc1 -> " + resNodeAcc1.String())
	ctx1.Logger().Info("resNodeAcc2 -> " + resNodeAcc2.String())
	ctx1.Logger().Info("resNodeAcc3 -> " + resNodeAcc3.String())
	ctx1.Logger().Info("resNodeAcc4 -> " + resNodeAcc4.String())
	ctx1.Logger().Info("resNodeAcc5 -> " + resNodeAcc5.String())
	ctx1.Logger().Info("idxNodeAcc1 -> " + idxNodeAcc1.String())
	ctx1.Logger().Info("idxNodeAcc2 -> " + idxNodeAcc2.String())
	ctx1.Logger().Info("idxNodeAcc3 -> " + idxNodeAcc3.String())
	ctx1.Logger().Info("spNodeIdxNodeAcc1 -> " + spNodeIdxNodeAcc1.String())
	//ctx1.Logger().Info("foundationAcc -> " + foundationAcc.String())
	ctx1.Logger().Info("sdsAcc1 -> " + sdsAcc1.String())
	ctx1.Logger().Info("sdsAcc2 -> " + sdsAcc2.String())
	ctx1.Logger().Info("sdsAcc3 -> " + sdsAcc3.String())

	return accs
}

func setupAllResourceNodes() []register.ResourceNode {
	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	resourceNode1 := register.NewResourceNode("sds://resourceNode1", pubKeyRes1, resOwner1, register.NewDescription("sds://resourceNode1", "", "", "", ""), "4", time)
	resourceNode2 := register.NewResourceNode("sds://resourceNode2", pubKeyRes2, resOwner2, register.NewDescription("sds://resourceNode2", "", "", "", ""), "4", time)
	resourceNode3 := register.NewResourceNode("sds://resourceNode3", pubKeyRes3, resOwner3, register.NewDescription("sds://resourceNode3", "", "", "", ""), "4", time)
	resourceNode4 := register.NewResourceNode("sds://resourceNode4", pubKeyRes4, resOwner4, register.NewDescription("sds://resourceNode4", "", "", "", ""), "4", time)
	resourceNode5 := register.NewResourceNode("sds://resourceNode5", pubKeyRes5, resOwner5, register.NewDescription("sds://resourceNode5", "", "", "", ""), "4", time)

	resourceNode1 = resourceNode1.AddToken(initialStakeRes1)
	resourceNode2 = resourceNode2.AddToken(initialStakeRes2)
	resourceNode3 = resourceNode3.AddToken(initialStakeRes3)
	resourceNode4 = resourceNode4.AddToken(initialStakeRes4)
	resourceNode5 = resourceNode5.AddToken(initialStakeRes5)

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

	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	indexingNode1 := register.NewIndexingNode("sds://indexingNode1", pubKeyIdx1, idxOwner1, register.NewDescription("sds://indexingNode1", "", "", "", ""), time)
	indexingNode2 := register.NewIndexingNode("sds://indexingNode2", pubKeyIdx2, idxOwner2, register.NewDescription("sds://indexingNode2", "", "", "", ""), time)
	indexingNode3 := register.NewIndexingNode("sds://indexingNode3", pubKeyIdx3, idxOwner3, register.NewDescription("sds://indexingNode3", "", "", "", ""), time)
	spNodeIndexingNode1 := register.NewIndexingNode("sds://sdsIndexingNode1", spNodePubKeyIdx1, sdsAccAddr1, register.NewDescription("sds://sdsIndexingNode1", "", "", "", ""), time)

	indexingNode1 = indexingNode1.AddToken(initialStakeIdx1)
	indexingNode2 = indexingNode2.AddToken(initialStakeIdx2)
	indexingNode3 = indexingNode3.AddToken(initialStakeIdx3)
	spNodeIndexingNode1 = spNodeIndexingNode1.AddToken(spNodeInitialStakeIdx1)

	indexingNode1.Status = sdk.Bonded
	indexingNode2.Status = sdk.Bonded
	indexingNode3.Status = sdk.Bonded
	spNodeIndexingNode1.Status = sdk.Bonded

	indexingNodes = append(indexingNodes, indexingNode1)
	indexingNodes = append(indexingNodes, indexingNode2)
	indexingNodes = append(indexingNodes, indexingNode3)
	indexingNodes = append(indexingNodes, spNodeIndexingNode1)

	return indexingNodes

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
