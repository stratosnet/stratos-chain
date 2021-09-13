package keeper_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/stratosnet/stratos-chain/app"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/evm/keeper"
	"github.com/stratosnet/stratos-chain/x/evm/types"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	abci "github.com/tendermint/tendermint/abci/types"
)

const addrHex = "0x756F45E3FA69347A9A973A725E3C98bC4db0b4c1"
const hex = "0x0d87a3a5f73140f46aac1bf419263e4e94e87c292f25007700ab7f2060e2af68"

var (
	hash = ethcmn.FromHex(hex)
)

type KeeperTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	querier sdk.Querier
	app     *app.NewApp
	address ethcmn.Address
}

func (suite *KeeperTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "stratos-3", Time: time.Now().UTC()})
	suite.querier = keeper.NewQuerier(*suite.app.GetEvmKeeper())
	suite.address = ethcmn.HexToAddress(addrHex)

	balance := sdk.NewCoins(sdk.NewCoin("ustos", sdk.ZeroInt()))
	acc := &stratos.StAccount{
		BaseAccount: auth.NewBaseAccount(sdk.AccAddress(suite.address.Bytes()), balance, nil, 0, 0),
		CodeHash:    ethcrypto.Keccak256(nil),
	}

	suite.app.GetAccountKeeper().SetAccount(suite.ctx, acc)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestTransactionLogs() {
	ethHash := ethcmn.BytesToHash(hash)
	log := &ethtypes.Log{
		Address:     suite.address,
		Data:        []byte("log"),
		BlockNumber: 10,
	}
	log2 := &ethtypes.Log{
		Address:     suite.address,
		Data:        []byte("log2"),
		BlockNumber: 11,
	}
	expLogs := []*ethtypes.Log{log}

	err := suite.app.GetEvmKeeper().SetLogs(suite.ctx, ethHash, expLogs)
	suite.Require().NoError(err)

	logs, err := suite.app.GetEvmKeeper().GetLogs(suite.ctx, ethHash)
	suite.Require().NoError(err)
	suite.Require().Equal(expLogs, logs)

	expLogs = []*ethtypes.Log{log2, log}

	// add another log under the zero hash
	suite.app.GetEvmKeeper().AddLog(suite.ctx, log2)
	logs = suite.app.GetEvmKeeper().AllLogs(suite.ctx)
	suite.Require().Equal(expLogs, logs)

	// add another log under the zero hash
	log3 := &ethtypes.Log{
		Address:     suite.address,
		Data:        []byte("log3"),
		BlockNumber: 10,
	}
	suite.app.GetEvmKeeper().AddLog(suite.ctx, log3)

	txLogs := suite.app.GetEvmKeeper().GetAllTxLogs(suite.ctx)
	suite.Require().Equal(2, len(txLogs))

	suite.Require().Equal(ethcmn.Hash{}.String(), txLogs[0].Hash)
	suite.Require().Equal([]*ethtypes.Log{log2, log3}, txLogs[0].Logs)

	suite.Require().Equal(ethHash.String(), txLogs[1].Hash)
	suite.Require().Equal([]*ethtypes.Log{log}, txLogs[1].Logs)
}

func (suite *KeeperTestSuite) TestDBStorage() {
	// Perform state transitions
	suite.app.GetEvmKeeper().CreateAccount(suite.ctx, suite.address)
	suite.app.GetEvmKeeper().SetBalance(suite.ctx, suite.address, big.NewInt(5))
	suite.app.GetEvmKeeper().SetNonce(suite.ctx, suite.address, 4)
	suite.app.GetEvmKeeper().SetState(suite.ctx, suite.address, ethcmn.HexToHash("0x2"), ethcmn.HexToHash("0x3"))
	suite.app.GetEvmKeeper().SetCode(suite.ctx, suite.address, []byte{0x1})

	// Test block hash mapping functionality
	suite.app.GetEvmKeeper().SetBlockHash(suite.ctx, hash, 7)
	height, found := suite.app.GetEvmKeeper().GetBlockHash(suite.ctx, hash)
	suite.Require().True(found)
	suite.Require().Equal(int64(7), height)

	suite.app.GetEvmKeeper().SetBlockHash(suite.ctx, []byte{0x43, 0x32}, 8)

	// Test block height mapping functionality
	testBloom := ethtypes.BytesToBloom([]byte{0x1, 0x3})
	suite.app.GetEvmKeeper().SetBlockBloom(suite.ctx, 4, testBloom)

	// Get those state transitions
	suite.Require().Equal(suite.app.GetEvmKeeper().GetBalance(suite.ctx, suite.address).Cmp(big.NewInt(5)), 0)
	suite.Require().Equal(suite.app.GetEvmKeeper().GetNonce(suite.ctx, suite.address), uint64(4))
	suite.Require().Equal(suite.app.GetEvmKeeper().GetState(suite.ctx, suite.address, ethcmn.HexToHash("0x2")), ethcmn.HexToHash("0x3"))
	suite.Require().Equal(suite.app.GetEvmKeeper().GetCode(suite.ctx, suite.address), []byte{0x1})

	height, found = suite.app.GetEvmKeeper().GetBlockHash(suite.ctx, hash)
	suite.Require().True(found)
	suite.Require().Equal(height, int64(7))
	height, found = suite.app.GetEvmKeeper().GetBlockHash(suite.ctx, []byte{0x43, 0x32})
	suite.Require().True(found)
	suite.Require().Equal(height, int64(8))

	bloom, found := suite.app.GetEvmKeeper().GetBlockBloom(suite.ctx, 4)
	suite.Require().True(found)
	suite.Require().Equal(bloom, testBloom)

	// commit stateDB
	_, err := suite.app.GetEvmKeeper().Commit(suite.ctx, false)
	suite.Require().NoError(err, "failed to commit StateDB")

	// simulate BaseApp EndBlocker commitment
	suite.app.Commit()
}

func (suite *KeeperTestSuite) TestChainConfig() {
	config, found := suite.app.GetEvmKeeper().GetChainConfig(suite.ctx)
	suite.Require().True(found)
	suite.Require().Equal(types.DefaultChainConfig(), config)

	config.EIP150Block = sdk.NewInt(100)
	suite.app.GetEvmKeeper().SetChainConfig(suite.ctx, config)
	newConfig, found := suite.app.GetEvmKeeper().GetChainConfig(suite.ctx)
	suite.Require().True(found)
	suite.Require().Equal(config, newConfig)
}
