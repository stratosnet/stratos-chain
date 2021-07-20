package register

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	chainID              = ""
	AccountAddressPrefix = "st"
)

var (
	AccountPubKeyPrefix    = AccountAddressPrefix + "pub"
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	ValidatorPubKeyPrefix  = AccountAddressPrefix + "valoperpub"
	ConsNodeAddressPrefix  = AccountAddressPrefix + "valcons"
	ConsNodePubKeyPrefix   = AccountAddressPrefix + "valconspub"

	resOwnerPrivKey1 = secp256k1.GenPrivKey()
	resOwnerPrivKey2 = secp256k1.GenPrivKey()
	idxOwnerPrivKey1 = secp256k1.GenPrivKey()
	idxOwnerPrivKey2 = secp256k1.GenPrivKey()
	idxOwnerPrivKey3 = secp256k1.GenPrivKey()

	resOwnerAddr1 = sdk.AccAddress(resOwnerPrivKey1.PubKey().Address())
	resOwnerAddr2 = sdk.AccAddress(resOwnerPrivKey2.PubKey().Address())
	idxOwnerAddr1 = sdk.AccAddress(idxOwnerPrivKey1.PubKey().Address())
	idxOwnerAddr2 = sdk.AccAddress(idxOwnerPrivKey2.PubKey().Address())
	idxOwnerAddr3 = sdk.AccAddress(idxOwnerPrivKey3.PubKey().Address())

	resOwnerInitBalance = sdk.NewInt(1000000000000)
	idxOwnerInitBalance = sdk.NewInt(1000000000000)

	resNodePrivKey1 = secp256k1.GenPrivKey()
	resNodePrivKey2 = secp256k1.GenPrivKey()
	idxNodePrivKey1 = secp256k1.GenPrivKey()
	idxNodePrivKey2 = secp256k1.GenPrivKey()
	idxNodePrivKey3 = secp256k1.GenPrivKey()

	resNodePubKey1 = resNodePrivKey1.PubKey()
	resNodePubKey2 = resNodePrivKey2.PubKey()
	idxNodePubKey1 = idxNodePrivKey1.PubKey()
	idxNodePubKey2 = idxNodePrivKey2.PubKey()
	idxNodePubKey3 = idxNodePrivKey3.PubKey()

	resNodeAddr1 = sdk.AccAddress(resNodePubKey1.Address())
	resNodeAddr2 = sdk.AccAddress(resNodePubKey2.Address())
	idxNodeAddr1 = sdk.AccAddress(idxNodePubKey1.Address())
	idxNodeAddr2 = sdk.AccAddress(idxNodePubKey2.Address())
	idxNodeAddr3 = sdk.AccAddress(idxNodePubKey3.Address())

	resNodeInitStake = sdk.NewInt(10000000000)
	idxNodeInitStake = sdk.NewInt(10000000000)
)

func setupAllResourceNodes() []ResourceNode {
	resourceNode1 := NewResourceNode("sds://resourceNode1", resNodePubKey1, resOwnerAddr1, NewDescription("sds://resourceNode1", "", "", "", ""), "4")
	resourceNode1 = resourceNode1.AddToken(resNodeInitStake)
	resourceNode1.Status = sdk.Bonded

	var resourceNodes []ResourceNode
	resourceNodes = append(resourceNodes, resourceNode1)
	return resourceNodes
}

func setupAllIndexingNodes() []IndexingNode {
	var indexingNodes []IndexingNode
	indexingNode1 := NewIndexingNode("sds://indexingNode1", idxNodePubKey1, idxOwnerAddr1, NewDescription("sds://indexingNode1", "", "", "", ""))
	indexingNode2 := NewIndexingNode("sds://indexingNode2", idxNodePubKey2, idxOwnerAddr2, NewDescription("sds://indexingNode2", "", "", "", ""))

	indexingNode1 = indexingNode1.AddToken(idxNodeInitStake)
	indexingNode2 = indexingNode2.AddToken(idxNodeInitStake)

	indexingNode1.Status = sdk.Bonded
	indexingNode2.Status = sdk.Unbonded

	indexingNodes = append(indexingNodes, indexingNode1)
	indexingNodes = append(indexingNodes, indexingNode2)

	return indexingNodes

}

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
		Address: resOwnerAddr1,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", resOwnerInitBalance)},
	}
	resOwnerAcc2 := &auth.BaseAccount{
		Address: resOwnerAddr2,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", resOwnerInitBalance)},
	}

	idxOwnerAcc1 := &auth.BaseAccount{
		Address: idxOwnerAddr1,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", idxOwnerInitBalance)},
	}
	idxOwnerAcc2 := &auth.BaseAccount{
		Address: idxOwnerAddr2,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", idxOwnerInitBalance)},
	}
	idxOwnerAcc3 := &auth.BaseAccount{
		Address: idxOwnerAddr3,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", idxOwnerInitBalance)},
	}

	idxNodeAcc1 := &auth.BaseAccount{
		Address: idxNodeAddr1,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", sdk.ZeroInt())},
	}

	accs := []authexported.Account{
		resOwnerAcc1, resOwnerAcc2, idxOwnerAcc1, idxOwnerAcc2, idxOwnerAcc3, idxNodeAcc1,
	}

	return accs
}
