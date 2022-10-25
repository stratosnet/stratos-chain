package register

//
//import (
//	"time"
//
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/cosmos/cosmos-sdk/x/auth"
//	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
//	"github.com/cosmos/cosmos-sdk/x/mock"
//	stratos "github.com/stratosnet/stratos-chain/types"
//	"github.com/tendermint/tendermint/crypto/ed25519"
//	"github.com/tendermint/tendermint/crypto/secp256k1"
//)
//
//var (
//	resOwnerPrivKey1 = secp256k1.GenPrivKey()
//	resOwnerPrivKey2 = secp256k1.GenPrivKey()
//	//resOwnerPrivKey3 = ed25519.GenPrivKey()
//	resOwnerPrivKey3 = secp256k1.GenPrivKey()
//	idxOwnerPrivKey1 = secp256k1.GenPrivKey()
//	idxOwnerPrivKey2 = secp256k1.GenPrivKey()
//	idxOwnerPrivKey3 = secp256k1.GenPrivKey()
//
//	resOwnerAddr1 = sdk.AccAddress(resOwnerPrivKey1.PubKey().Address())
//	resOwnerAddr2 = sdk.AccAddress(resOwnerPrivKey2.PubKey().Address())
//	resOwnerAddr3 = sdk.AccAddress(resOwnerPrivKey3.PubKey().Address())
//	idxOwnerAddr1 = sdk.AccAddress(idxOwnerPrivKey1.PubKey().Address())
//	idxOwnerAddr2 = sdk.AccAddress(idxOwnerPrivKey2.PubKey().Address())
//	idxOwnerAddr3 = sdk.AccAddress(idxOwnerPrivKey3.PubKey().Address())
//
//	resOwnerInitBalance = sdk.NewInt(1000000000000)
//	idxOwnerInitBalance = sdk.NewInt(1000000000000)
//
//	resNodePrivKey1 = secp256k1.GenPrivKey()
//	resNodePrivKey2 = secp256k1.GenPrivKey()
//	resNodePrivKey3 = ed25519.GenPrivKey()
//	//resNodePrivKey3 = secp256k1.GenPrivKey()
//	idxNodePrivKey1 = secp256k1.GenPrivKey()
//	idxNodePrivKey2 = secp256k1.GenPrivKey()
//	idxNodePrivKey3 = secp256k1.GenPrivKey()
//
//	resNodePubKey1 = resNodePrivKey1.PubKey()
//	resNodePubKey2 = resNodePrivKey2.PubKey()
//	resNodePubKey3 = resNodePrivKey3.PubKey()
//	idxNodePubKey1 = idxNodePrivKey1.PubKey()
//	idxNodePubKey2 = idxNodePrivKey2.PubKey()
//	idxNodePubKey3 = idxNodePrivKey3.PubKey()
//
//	resNodeAddr1 = sdk.AccAddress(resNodePubKey1.Address())
//	resNodeAddr2 = sdk.AccAddress(resNodePubKey2.Address())
//	resNodeAddr3 = sdk.AccAddress(resNodePubKey3.Address())
//	idxNodeAddr1 = sdk.AccAddress(idxNodePubKey1.Address())
//	idxNodeAddr2 = sdk.AccAddress(idxNodePubKey2.Address())
//	idxNodeAddr3 = sdk.AccAddress(idxNodePubKey3.Address())
//
//	resNodeNetworkId1 = stratos.SdsAddress(resNodePubKey1.Address())
//	resNodeNetworkId2 = stratos.SdsAddress(resNodePubKey2.Address())
//	resNodeNetworkId3 = stratos.SdsAddress(resNodePubKey3.Address())
//	idxNodeNetworkId1 = stratos.SdsAddress(idxNodePubKey1.Address())
//	idxNodeNetworkId2 = stratos.SdsAddress(idxNodePubKey2.Address())
//	idxNodeNetworkId3 = stratos.SdsAddress(idxNodePubKey3.Address())
//
//	resNodeInitStake   = sdk.NewInt(10000000000)
//	idxNodeInitStake   = sdk.NewInt(10000000000)
//	initialNOzonePrice = sdk.NewDec(1000000) // 0.001 gwei -> 1 noz
//)
//
//func setupAllResourceNodes() []ResourceNode {
//	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
//	resourceNode1 := NewResourceNode(resNodeNetworkId1, resNodePubKey1, resOwnerAddr1, NewDescription("sds://resourceNode1", "", "", "", ""), 4, time)
//	resourceNode1 = resourceNode1.AddToken(resNodeInitStake)
//	resourceNode1.Status = sdk.Bonded
//
//	resourceNode3 := NewResourceNode(resNodeNetworkId3, resNodePubKey3, resOwnerAddr3, NewDescription("sds://resourceNode3", "", "", "", ""), 4, time)
//	resourceNode3 = resourceNode3.AddToken(resNodeInitStake)
//	resourceNode3.Status = sdk.Bonded
//
//	var resourceNodes []ResourceNode
//	resourceNodes = append(resourceNodes, resourceNode1, resourceNode3)
//	return resourceNodes
//}
//
//func setupAllIndexingNodes() []IndexingNode {
//	var indexingNodes []IndexingNode
//	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
//	indexingNode1 := NewIndexingNode(stratos.SdsAddress(idxNodeAddr1), idxNodePubKey1, idxOwnerAddr1, NewDescription("sds://indexingNode1", "", "", "", ""), time)
//	indexingNode2 := NewIndexingNode(stratos.SdsAddress(idxNodeAddr2), idxNodePubKey2, idxOwnerAddr2, NewDescription("sds://indexingNode2", "", "", "", ""), time)
//
//	indexingNode1 = indexingNode1.AddToken(idxNodeInitStake)
//	indexingNode2 = indexingNode2.AddToken(idxNodeInitStake)
//
//	indexingNode1.Status = sdk.Bonded
//	indexingNode2.Status = sdk.Unbonded
//
//	indexingNodes = append(indexingNodes, indexingNode1)
//	indexingNodes = append(indexingNodes, indexingNode2)
//
//	return indexingNodes
//
//}
//
//func setupAccounts(mApp *mock.App) []authexported.Account {
//	//************************** setup resource nodes owners' accounts **************************
//	resOwnerAcc1 := &auth.BaseAccount{
//		Address: resOwnerAddr1,
//		Coins:   sdk.Coins{sdk.NewCoin("wei", resOwnerInitBalance)},
//	}
//	resOwnerAcc2 := &auth.BaseAccount{
//		Address: resOwnerAddr2,
//		Coins:   sdk.Coins{sdk.NewCoin("wei", resOwnerInitBalance)},
//	}
//
//	resOwnerAcc3 := &auth.BaseAccount{
//		Address: resOwnerAddr3,
//		Coins:   sdk.Coins{sdk.NewCoin("wei", resOwnerInitBalance)},
//	}
//
//	idxOwnerAcc1 := &auth.BaseAccount{
//		Address: idxOwnerAddr1,
//		Coins:   sdk.Coins{sdk.NewCoin("wei", idxOwnerInitBalance)},
//	}
//	idxOwnerAcc2 := &auth.BaseAccount{
//		Address: idxOwnerAddr2,
//		Coins:   sdk.Coins{sdk.NewCoin("wei", idxOwnerInitBalance)},
//	}
//	idxOwnerAcc3 := &auth.BaseAccount{
//		Address: idxOwnerAddr3,
//		Coins:   sdk.Coins{sdk.NewCoin("wei", idxOwnerInitBalance)},
//	}
//
//	resNodeAcc2 := &auth.BaseAccount{
//		Address: resNodeAddr2,
//		Coins:   sdk.Coins{sdk.NewCoin("wei", sdk.ZeroInt())},
//	}
//
//	resNodeAcc3 := &auth.BaseAccount{
//		Address: resNodeAddr3,
//		Coins:   sdk.Coins{sdk.NewCoin("wei", sdk.ZeroInt())},
//	}
//
//	idxNodeAcc1 := &auth.BaseAccount{
//		Address: idxNodeAddr1,
//		Coins:   sdk.Coins{sdk.NewCoin("wei", sdk.ZeroInt())},
//	}
//
//	idxNodeAcc3 := &auth.BaseAccount{
//		Address: idxNodeAddr3,
//		Coins:   sdk.Coins{sdk.NewCoin("wei", sdk.ZeroInt())},
//	}
//
//	accs := []authexported.Account{
//		resOwnerAcc1, resOwnerAcc2, resOwnerAcc3, idxOwnerAcc1, idxOwnerAcc2, idxOwnerAcc3, resNodeAcc2, resNodeAcc3, idxNodeAcc1, idxNodeAcc3,
//	}
//
//	return accs
//}
