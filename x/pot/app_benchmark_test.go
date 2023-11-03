package pot_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"

	sdkmath "cosmossdk.io/math"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratostestutil "github.com/stratosnet/stratos-chain/testutil"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
	sdstypes "github.com/stratosnet/stratos-chain/x/sds/types"
)

const (
	testchainID   = "testchain"
	funderKeyName = "funder"
	reNodeCount   = 1e5
)

var (
	keysMap  = make(map[string]KeyInfo, 0) // map[bech32 P2PAddr]KeyInfo
	keysList = make([]KeyInfo, 0)
	accounts = make([]authtypes.GenesisAccount, 0)
	balances = make([]banktypes.Balance, 0)

	accInitBalance        = sdkmath.NewInt(100).MulRaw(stratos.StosToWei)
	initFoundationDeposit = sdk.NewCoins(sdk.NewCoin(stratos.Wei, sdkmath.NewInt(4e7).MulRaw(stratos.StosToWei)))
	nodeInitDeposit       = sdkmath.NewInt(1 * stratos.StosToWei)
	prepayAmt             = sdk.NewCoins(stratos.NewCoin(sdkmath.NewInt(20).MulRaw(stratos.StosToWei)))
	valP2PAddrBech32      string
)

type KeyInfo struct {
	ed25519PrivKey   *ed25519.PrivKey
	secp256k1PrivKey *secp256k1.PrivKey
}

func (k KeyInfo) P2PPubKey() cryptotypes.PubKey {
	return k.ed25519PrivKey.PubKey()
}

func (k KeyInfo) P2PAddressBech32() string {
	addr := k.ed25519PrivKey.PubKey().Address()
	return stratos.SdsAddress(addr).String()
}

func (k KeyInfo) P2PAddress() stratos.SdsAddress {
	addr := k.ed25519PrivKey.PubKey().Address()
	return stratos.SdsAddress(addr)
}

func (k KeyInfo) OwnerAddress() sdk.AccAddress {
	addr := k.secp256k1PrivKey.PubKey().Address()
	return sdk.AccAddress(addr)
}

func (k KeyInfo) SignKey() *secp256k1.PrivKey {
	return k.secp256k1PrivKey
}

func NewKeyInfo() KeyInfo {
	return KeyInfo{
		ed25519PrivKey:   ed25519.GenPrivKey(),
		secp256k1PrivKey: secp256k1.GenPrivKey(),
	}
}

func TestVolumeReportBenchmark(t *testing.T) {
	/********************* initialize mock app *********************/
	setupKeysAndAccBalance(reNodeCount)
	metaNodes, resourceNodes := setupNodesBenchmark()

	// create validator set with single validator
	consPubKey, err := cryptocodec.ToTmPubKeyInterface(keysList[0].P2PPubKey())
	validator := tmtypes.NewValidator(consPubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	stApp := stratostestutil.SetupWithGenesisNodeSet(t, valSet, metaNodes, resourceNodes, accounts, testchainID, false, balances...)
	accountKeeper := stApp.GetAccountKeeper()

	/********************* foundation account deposit *********************/
	header := tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: testchainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := stApp.BaseApp.NewContext(true, header)

	foundationDepositMsg := types.NewMsgFoundationDeposit(initFoundationDeposit, keysMap[funderKeyName].OwnerAddress())
	txGen := stratostestutil.MakeTestEncodingConfig().TxConfig

	senderAcc := accountKeeper.GetAccount(ctx, keysMap[funderKeyName].OwnerAddress())
	accNum := senderAcc.GetAccountNumber()
	accSeq := senderAcc.GetSequence()
	_, _, err = stratostestutil.SignCheckDeliverWithFee(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, testchainID, []uint64{accNum}, []uint64{accSeq}, true, true, keysMap[funderKeyName].secp256k1PrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(types.FoundationAccount)
	stratostestutil.CheckBalance(t, stApp, foundationAccountAddr, initFoundationDeposit)

	/********************* prepay *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: testchainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(false, header)
	prepayMsg := sdstypes.NewMsgPrepay(resourceNodes[0].OwnerAddress, resourceNodes[0].OwnerAddress, prepayAmt)
	senderAcc = accountKeeper.GetAccount(ctx, keysMap[resourceNodes[0].NetworkAddress].OwnerAddress())

	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()
	_, _, err = stratostestutil.SignCheckDeliverWithFee(t, txGen, stApp.BaseApp, header, []sdk.Msg{prepayMsg}, testchainID, []uint64{accNum}, []uint64{accSeq}, true, true, keysMap[resourceNodes[0].NetworkAddress].secp256k1PrivKey)
	require.NoError(t, err)

	/********************** commit **********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: testchainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	/********************* prepare tx data *********************/
	volumeReportMsg := setupMsgVolumeReportBenchmark(t, sdkmath.NewInt(1), metaNodes, resourceNodes)

	/********************* deliver tx *********************/
	idxOwnerAcc1 := accountKeeper.GetAccount(ctx, keysMap[metaNodes[0].NetworkAddress].OwnerAddress())
	ownerAccNum := idxOwnerAcc1.GetAccountNumber()
	ownerAccSeq := idxOwnerAcc1.GetSequence()

	feePoolAccAddr := accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	require.NotNil(t, feePoolAccAddr)

	t.Log("--------------------------- deliver volumeReportMsg")
	gInfo, _, err := stratostestutil.SignCheckDeliverWithFee(t, txGen, stApp.BaseApp, header, []sdk.Msg{volumeReportMsg}, testchainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, keysMap[metaNodes[0].NetworkAddress].secp256k1PrivKey)
	require.NoError(t, err)
	fmt.Println("##### volume nodes count:", len(volumeReportMsg.WalletVolumes))
	fmt.Println("##### gInfo:", gInfo.String())
}

func setupKeysAndAccBalance(resNodeCnt int) {
	// 1 validator, 3 meta nodes
	totalNodesCnt := resNodeCnt + 4

	for i := 0; i < totalNodesCnt; i++ {
		keyInfo := NewKeyInfo()
		account := &authtypes.BaseAccount{
			Address: keyInfo.OwnerAddress().String(),
		}
		balance := banktypes.Balance{
			Address: keyInfo.OwnerAddress().String(),
			Coins:   sdk.Coins{stratos.NewCoin(accInitBalance)},
		}

		keysMap[keyInfo.P2PAddressBech32()] = keyInfo
		keysList = append(keysList, keyInfo)
		accounts = append(accounts, account)
		balances = append(balances, balance)

		// 1st key is validator key
		if i == 0 {
			valP2PAddrBech32 = keyInfo.P2PAddressBech32()
		}

	}

	funderKey := NewKeyInfo()
	funderAcc := &authtypes.BaseAccount{Address: funderKey.OwnerAddress().String()}
	feeAmt := sdkmath.NewInt(50).MulRaw(stratos.StosToWei)
	funderBalance := banktypes.Balance{
		Address: funderKey.OwnerAddress().String(),
		Coins:   initFoundationDeposit.Add(sdk.NewCoin(stratos.Wei, feeAmt)),
	}
	keysMap[funderKeyName] = funderKey
	accounts = append(accounts, funderAcc)
	balances = append(balances, funderBalance)
}

func setupNodesBenchmark() (metaNodes []registertypes.MetaNode, resourceNodes []registertypes.ResourceNode) {
	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	nodeType := registertypes.STORAGE

	for idx, keyInfo := range keysList {

		if 0 < idx && idx < 4 {
			// 1~3 keys are metaNode keys
			metaNode, _ := registertypes.NewMetaNode(
				keyInfo.P2PAddress(),
				keyInfo.P2PPubKey(),
				keyInfo.OwnerAddress(),
				keyInfo.OwnerAddress(),
				registertypes.NewDescription(keyInfo.P2PAddressBech32(), "", "", "", ""),
				time,
			)
			metaNode = metaNode.AddToken(nodeInitDeposit)
			metaNode.Status = stakingtypes.Bonded
			metaNode.Suspend = false

			metaNodes = append(metaNodes, metaNode)
		} else {
			resourceNode, _ := registertypes.NewResourceNode(
				keyInfo.P2PAddress(),
				keyInfo.P2PPubKey(),
				keyInfo.OwnerAddress(),
				registertypes.NewDescription(keyInfo.P2PAddressBech32(), "", "", "", ""),
				nodeType,
				time,
			)
			resourceNode = resourceNode.AddToken(nodeInitDeposit)
			resourceNode.EffectiveTokens = nodeInitDeposit
			resourceNode.Status = stakingtypes.Bonded
			resourceNode.Suspend = false

			resourceNodes = append(resourceNodes, resourceNode)
		}
	}

	return
}

// initialize data of volume report
func setupMsgVolumeReportBenchmark(t *testing.T, epoch sdkmath.Int, metaNodes []registertypes.MetaNode, resourceNodes []registertypes.ResourceNode) *types.MsgVolumeReport {
	rsNodeVolume := sdkmath.NewInt(50000)

	nodesVolume := make([]types.SingleWalletVolume, 0)
	for _, rsNode := range resourceNodes {
		ownerAddr := keysMap[rsNode.NetworkAddress].OwnerAddress()
		volume := types.NewSingleWalletVolume(ownerAddr, rsNodeVolume)
		nodesVolume = append(nodesVolume, volume)
	}

	reporterNode := metaNodes[0]
	reporterKey := keysMap[reporterNode.NetworkAddress]

	reporter := reporterKey.P2PAddress()
	reportReference := "report for epoch " + epoch.String()
	reporterOwner := reporterKey.OwnerAddress()

	volumeReportMsg := types.NewMsgVolumeReport(nodesVolume, reporter, epoch, reportReference, reporterOwner)
	volumeReportMsg, err := stratostestutil.SignVolumeReport(
		volumeReportMsg,
		keysMap[metaNodes[0].NetworkAddress].ed25519PrivKey.Bytes(),
		keysMap[metaNodes[1].NetworkAddress].ed25519PrivKey.Bytes(),
		keysMap[metaNodes[2].NetworkAddress].ed25519PrivKey.Bytes(),
	)
	require.NoError(t, err)

	return volumeReportMsg
}
