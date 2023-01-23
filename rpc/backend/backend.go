package backend

import (
	"context"
	"math/big"
	"time"

	cs "github.com/tendermint/tendermint/consensus"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/mempool"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	tmrpctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/store"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/stratosnet/stratos-chain/rpc/types"
	"github.com/stratosnet/stratos-chain/server/config"
	evmkeeper "github.com/stratosnet/stratos-chain/x/evm/keeper"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// BackendI implements the Cosmos and EVM backend.
type BackendI interface { // nolint: revive
	CosmosBackend
	EVMBackend
	TMBackend
}

// CosmosBackend implements the functionality shared within cosmos namespaces
// as defined by Wallet Connect V2: https://docs.walletconnect.com/2.0/json-rpc/cosmos.
// Implemented by Backend.
type CosmosBackend interface {
	// TODO: define
	// GetAccounts()
	// SignDirect()
	// SignAmino()
	GetEVMKeeper() *evmkeeper.Keeper
	GetSdkContext(header *tmtypes.Header) sdk.Context
}

// EVMBackend implements the functionality shared within ethereum namespaces
// as defined by EIP-1474: https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1474.md
// Implemented by Backend.
type EVMBackend interface {
	// General Ethereum API
	RPCGasCap() uint64            // global gas cap for eth_call over rpc: DoS protection
	RPCEVMTimeout() time.Duration // global timeout for eth_call over rpc: DoS protection
	RPCTxFeeCap() float64         // RPCTxFeeCap is the global transaction fee(price * gaslimit) cap for send-transaction variants. The unit is ether.

	RPCMinGasPrice() int64
	SuggestGasTipCap() (*big.Int, error)

	// Blockchain API
	BlockNumber() (hexutil.Uint64, error)
	GetTendermintBlockByNumber(blockNum types.BlockNumber) (*tmrpctypes.ResultBlock, error)
	GetTendermintBlockByHash(blockHash common.Hash) (*tmrpctypes.ResultBlock, error)
	GetBlockByNumber(blockNum types.BlockNumber, fullTx bool) (map[string]interface{}, error)
	GetBlockByHash(hash common.Hash, fullTx bool) (map[string]interface{}, error)
	BlockByNumber(blockNum types.BlockNumber) (*ethtypes.Block, error)
	BlockByHash(blockHash common.Hash) (*ethtypes.Block, error)
	CurrentHeader() *ethtypes.Header
	HeaderByNumber(blockNum types.BlockNumber) (*ethtypes.Header, error)
	HeaderByHash(blockHash common.Hash) (*ethtypes.Header, error)
	PendingTransactions() ([]*sdk.Tx, error)
	GetTransactionCount(address common.Address, blockNum types.BlockNumber) (*hexutil.Uint64, error)
	SendTransaction(args evmtypes.TransactionArgs) (common.Hash, error)
	GetCoinbase() (sdk.AccAddress, error)
	GetTransactionByHash(txHash common.Hash) (*types.RPCTransaction, error)
	GetTxByHash(txHash common.Hash) (*tmrpctypes.ResultTx, error)
	GetTxByTxIndex(height int64, txIndex uint) (*tmrpctypes.ResultTx, error)
	EstimateGas(args evmtypes.TransactionArgs, blockNrOptional *types.BlockNumber) (hexutil.Uint64, error)
	BaseFee() (*big.Int, error)

	// Fee API
	FeeHistory(blockCount rpc.DecimalOrHex, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*types.FeeHistoryResult, error)

	// Filter API
	BloomStatus() (uint64, uint64)
	GetLogs(hash common.Hash) ([][]*ethtypes.Log, error)
	GetLogsByHeight(height *int64) ([][]*ethtypes.Log, error)
	ChainConfig() *params.ChainConfig
	SetTxDefaults(args evmtypes.TransactionArgs) (evmtypes.TransactionArgs, error)
}

type TMBackend interface {
	// tendermint helpers
	GetNode() *node.Node
	GetBlockStore() *store.BlockStore
	GetMempool() mempool.Mempool
	GetConsensusReactor() *cs.Reactor
	GetSwitch() *p2p.Switch
}

var _ BackendI = (*Backend)(nil)

// Backend implements the BackendI interface
type Backend struct {
	ctx       context.Context
	clientCtx client.Context
	tmNode    *node.Node // directly tendermint access, new impl
	evmkeeper *evmkeeper.Keeper
	sdkCtx    sdk.Context
	logger    log.Logger
	cfg       config.Config
}

// NewBackend creates a new Backend instance for cosmos and ethereum namespaces
func NewBackend(ctx *server.Context, tmNode *node.Node, evmkeeper *evmkeeper.Keeper, sdkCtx sdk.Context, logger log.Logger, clientCtx client.Context) *Backend {
	appConf, err := config.GetConfig(ctx.Viper)
	if err != nil {
		panic(err)
	}

	return &Backend{
		ctx:       context.Background(),
		clientCtx: clientCtx,
		tmNode:    tmNode,
		evmkeeper: evmkeeper,
		sdkCtx:    sdkCtx,
		logger:    logger.With("module", "backend"),
		cfg:       appConf,
	}
}

func (b *Backend) GetEVMKeeper() *evmkeeper.Keeper {
	return b.evmkeeper
}

func (b *Backend) GetSdkContext(header *tmtypes.Header) sdk.Context {
	sdkCtx := b.sdkCtx
	if header != nil {
		sdkCtx = sdkCtx.WithHeaderHash(header.Hash())
		header := types.FormatTmHeaderToProto(header)
		sdkCtx = sdkCtx.WithBlockHeader(header)
	}
	return sdkCtx
}

func (b *Backend) GetNode() *node.Node {
	return b.tmNode
}

func (b *Backend) GetBlockStore() *store.BlockStore {
	return b.tmNode.BlockStore()
}

func (b *Backend) GetMempool() mempool.Mempool {
	return b.tmNode.Mempool()
}

func (b *Backend) GetConsensusReactor() *cs.Reactor {
	return b.tmNode.ConsensusReactor()
}

func (b *Backend) GetSwitch() *p2p.Switch {
	return b.tmNode.Switch()
}
