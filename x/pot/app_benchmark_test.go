package pot_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/stratosnet/stratos-chain/app"
	"github.com/stratosnet/stratos-chain/crypto"
	"github.com/stratosnet/stratos-chain/crypto/bls"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
	sdstypes "github.com/stratosnet/stratos-chain/x/sds/types"
)

const (
	testchainID = "testchain"

	reNodeCount = 1e5
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

var (
	keysMap  = make(map[string]KeyInfo, 0) // map[bech32 P2PAddr]KeyInfo
	keysList = make([]KeyInfo, 0)
	accounts = make([]authtypes.GenesisAccount, 0)
	balances = make([]banktypes.Balance, 0)

	accInitBalance        = sdk.NewInt(100).Mul(sdk.NewInt(stratos.StosToWei))
	initFoundationDeposit = sdk.NewCoins(sdk.NewCoin(stratos.Wei, sdk.NewInt(40000000000000000).MulRaw(stratos.GweiToWei)))

	nodeInitDeposit  = sdk.NewInt(1 * stratos.StosToWei)
	prepayAmt        = sdk.NewCoins(stratos.NewCoin(sdk.NewInt(20).Mul(sdk.NewInt(stratos.StosToWei))))
	valP2PAddrBech32 string
)

func TestVolumeReportBenchmark(t *testing.T) {
	/********************* initialize mock app *********************/
	setupKeysAndAccBalance(reNodeCount)
	createValidatorMsg, metaNodes, resourceNodes := setupNodesBenchmark()

	validators := make([]*tmtypes.Validator, 0)
	valSet := tmtypes.NewValidatorSet(validators)

	//fmt.Println("##### accounts: ", accounts)
	//fmt.Println("!!!!! balances: ", balances)

	stApp := app.SetupWithGenesisNodeSet(t, false, valSet, metaNodes, resourceNodes, accounts, testchainID, balances...)
	accountKeeper := stApp.GetAccountKeeper()

	/********************* foundation account deposit *********************/
	header := tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: testchainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := stApp.BaseApp.NewContext(true, header)

	foundationDepositMsg := types.NewMsgFoundationDeposit(initFoundationDeposit, keysMap["foundationDepositorKey"].OwnerAddress())
	txGen := app.MakeTestEncodingConfig().TxConfig

	senderAcc := accountKeeper.GetAccount(ctx, keysMap["foundationDepositorKey"].OwnerAddress())
	accNum := senderAcc.GetAccountNumber()
	accSeq := senderAcc.GetSequence()
	_, _, err := app.SignCheckDeliverWithFee(t, txGen, stApp.BaseApp, header, []sdk.Msg{foundationDepositMsg}, testchainID, []uint64{accNum}, []uint64{accSeq}, true, true, keysMap["foundationDepositorKey"].secp256k1PrivKey)
	require.NoError(t, err)
	foundationAccountAddr := accountKeeper.GetModuleAddress(types.FoundationAccount)
	app.CheckBalance(t, stApp, foundationAccountAddr, initFoundationDeposit)

	/********************* create validator with 50% commission *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: testchainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(false, header)

	senderAcc = accountKeeper.GetAccount(ctx, keysMap[valP2PAddrBech32].OwnerAddress())
	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()
	_, _, err = app.SignCheckDeliverWithFee(t, txGen, stApp.BaseApp, header, []sdk.Msg{createValidatorMsg}, testchainID, []uint64{accNum}, []uint64{accSeq}, true, true, keysMap[valP2PAddrBech32].secp256k1PrivKey)
	require.NoError(t, err)

	/********************* prepay *********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: testchainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(false, header)
	prepayMsg := sdstypes.NewMsgPrepay(resourceNodes[0].OwnerAddress, resourceNodes[0].OwnerAddress, prepayAmt)
	senderAcc = accountKeeper.GetAccount(ctx, keysMap[resourceNodes[0].NetworkAddress].OwnerAddress())

	accNum = senderAcc.GetAccountNumber()
	accSeq = senderAcc.GetSequence()
	_, _, err = app.SignCheckDeliverWithFee(t, txGen, stApp.BaseApp, header, []sdk.Msg{prepayMsg}, testchainID, []uint64{accNum}, []uint64{accSeq}, true, true, keysMap[resourceNodes[0].NetworkAddress].secp256k1PrivKey)
	require.NoError(t, err)

	/********************** commit **********************/
	header = tmproto.Header{Height: stApp.LastBlockHeight() + 1, ChainID: testchainID}
	stApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx = stApp.BaseApp.NewContext(true, header)

	/********************* prepare tx data *********************/
	volumeReportMsg := setupMsgVolumeReportBenchmark(t, sdk.NewInt(1), metaNodes, resourceNodes)

	/********************* deliver tx *********************/
	idxOwnerAcc1 := accountKeeper.GetAccount(ctx, keysMap[metaNodes[0].NetworkAddress].OwnerAddress())
	ownerAccNum := idxOwnerAcc1.GetAccountNumber()
	ownerAccSeq := idxOwnerAcc1.GetSequence()

	feePoolAccAddr := accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	require.NotNil(t, feePoolAccAddr)

	t.Log("--------------------------- deliver volumeReportMsg")
	gInfo, _, err := app.SignCheckDeliverWithFee(t, txGen, stApp.BaseApp, header, []sdk.Msg{volumeReportMsg}, testchainID, []uint64{ownerAccNum}, []uint64{ownerAccSeq}, true, true, keysMap[metaNodes[0].NetworkAddress].secp256k1PrivKey)
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

	foundationDepositorKey := NewKeyInfo()
	foundationDepositorAcc := &authtypes.BaseAccount{Address: foundationDepositorKey.OwnerAddress().String()}
	feeAmt, _ := sdk.NewIntFromString("50000000000000000000")
	foundationDepositorBalance := banktypes.Balance{
		Address: foundationDepositorKey.OwnerAddress().String(),
		Coins:   append(initFoundationDeposit, sdk.NewCoin(stratos.Wei, feeAmt)),
	}
	keysMap["foundationDepositorKey"] = foundationDepositorKey
	accounts = append(accounts, foundationDepositorAcc)
	balances = append(balances, foundationDepositorBalance)
}

func setupNodesBenchmark() (createValidatorMsg *stakingtypes.MsgCreateValidator, metaNodes []registertypes.MetaNode, resourceNodes []registertypes.ResourceNode) {
	time, _ := time.Parse(time.RubyDate, "Fri Sep 24 10:37:13 -0400 2021")
	nodeType := registertypes.STORAGE

	for idx, keyInfo := range keysList {

		if idx == 0 {
			// first key is validator key
			commission := stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
			description := stakingtypes.NewDescription("foo_moniker", testchainID, "", "", "")
			createValidatorMsg, _ = stakingtypes.NewMsgCreateValidator(
				sdk.ValAddress(keyInfo.OwnerAddress()),
				keyInfo.P2PPubKey(),
				stratos.NewCoin(nodeInitDeposit),
				description,
				commission,
				sdk.OneInt(),
			)
		} else if idx < 4 {
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
func setupMsgVolumeReportBenchmark(t *testing.T, epoch sdk.Int, metaNodes []registertypes.MetaNode, resourceNodes []registertypes.ResourceNode) *types.MsgVolumeReport {
	rsNodeVolume := sdk.NewInt(50000)

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

	signature := types.BLSSignatureInfo{}
	volumeReportMsg := types.NewMsgVolumeReport(nodesVolume, reporter, epoch, reportReference, reporterOwner, signature)

	signBytes := volumeReportMsg.GetBLSSignBytes()
	signBytesHash := crypto.Keccak256(signBytes)

	// set blsSignature
	blsPrivKey1, blsPubKey1, err := bls.NewKeyPairFromBytes(keysMap[metaNodes[0].NetworkAddress].ed25519PrivKey.Bytes())
	require.NoError(t, err)
	blsPrivKey2, blsPubKey2, err := bls.NewKeyPairFromBytes(keysMap[metaNodes[1].NetworkAddress].ed25519PrivKey.Bytes())
	require.NoError(t, err)
	blsPrivKey3, blsPubKey3, err := bls.NewKeyPairFromBytes(keysMap[metaNodes[2].NetworkAddress].ed25519PrivKey.Bytes())
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
