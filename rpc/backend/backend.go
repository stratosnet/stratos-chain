package backend

import (
	"context"
	"fmt"
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

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/stratosnet/stratos-chain/rpc/types"
	"github.com/stratosnet/stratos-chain/server/config"
	evmkeeper "github.com/stratosnet/stratos-chain/x/evm/keeper"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
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
	GetSdkContext() sdk.Context
	GetSdkContextWithHeader(header *tmtypes.Header) (sdk.Context, error)
}

// EVMBackend implements the functionality shared within ethereum namespaces
// as defined by EIP-1474: https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1474.md
// Implemented by Backend.
type EVMBackend interface {
	// General Ethereum API
	RPCGasCap() uint64            // global gas cap for eth_call over rpc: DoS protection
	RPCEVMTimeout() time.Duration // global timeout for eth_call over rpc: DoS protection
	RPCTxFeeCap() float64         // RPCTxFeeCap is the global transaction fee(price * gaslimit) cap for send-transaction variants. The unit is ether.
	RPCFilterCap() int32
	RPCMinGasPrice() int64
	SuggestGasTipCap() (*big.Int, error)
	RPCLogsCap() int32
	RPCBlockRangeCap() int32

	// Blockchain API
	BlockNumber() (hexutil.Uint64, error)
	GetTendermintBlockByNumber(blockNum types.BlockNumber) (*tmrpctypes.ResultBlock, error)
	GetTendermintBlockByHash(blockHash common.Hash) (*tmrpctypes.ResultBlock, error)
	GetBlockByNumber(blockNum types.BlockNumber, fullTx bool) (*types.Block, error)
	GetBlockByHash(hash common.Hash, fullTx bool) (*types.Block, error)
	CurrentHeader() *types.Header
	HeaderByNumber(blockNum types.BlockNumber) (*types.Header, error)
	HeaderByHash(blockHash common.Hash) (*types.Header, error)
	PendingTransactions() ([]*sdk.Tx, error)
	GetTransactionCount(address common.Address, blockNum types.BlockNumber) (hexutil.Uint64, error)
	SendTransaction(args evmtypes.TransactionArgs) (common.Hash, error)
	GetCoinbase() (sdk.AccAddress, error)
	GetTransactionByHash(txHash common.Hash) (*types.Transaction, error)
	GetTxByHash(txHash common.Hash) (*tmrpctypes.ResultTx, error)
	GetTxByTxIndex(height int64, txIndex uint) (*tmrpctypes.ResultTx, error)
	EstimateGas(args evmtypes.TransactionArgs, blockNrOptional *types.BlockNumber) (hexutil.Uint64, error)
	BaseFee() (*big.Int, error)
	GetLogsByNumber(blockNum types.BlockNumber) ([][]*ethtypes.Log, error)
	BlockBloom(height *int64) (ethtypes.Bloom, error)

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
	ms        storetypes.MultiStore
	logger    log.Logger
	cfg       config.Config
}

// NewBackend creates a new Backend instance for cosmos and ethereum namespaces
func NewBackend(ctx *server.Context, tmNode *node.Node, evmkeeper *evmkeeper.Keeper, ms storetypes.MultiStore, logger log.Logger, clientCtx client.Context) *Backend {
	appConf, err := config.GetConfig(ctx.Viper)
	if err != nil {
		panic(err)
	}

	return &Backend{
		ctx:       context.Background(),
		clientCtx: clientCtx,
		tmNode:    tmNode,
		evmkeeper: evmkeeper,
		ms:        ms,
		logger:    logger.With("module", "backend"),
		cfg:       appConf,
	}
}

func (b *Backend) GetEVMKeeper() *evmkeeper.Keeper {
	return b.evmkeeper
}

func (b *Backend) copySdkContext(ms storetypes.MultiStore, header *tmtypes.Header) sdk.Context {
	sdkCtx := sdk.NewContext(ms, tmproto.Header{}, true, b.logger)
	if header != nil {
		return sdkCtx.WithHeaderHash(
			header.Hash(),
		).WithBlockHeader(
			types.FormatTmHeaderToProto(header),
		).WithBlockHeight(
			header.Height,
		).WithProposer(
			sdk.ConsAddress(header.ProposerAddress),
		)
	}
	return sdkCtx
}

func (b *Backend) GetSdkContext() sdk.Context {
	return b.copySdkContext(b.ms.CacheMultiStore(), nil)
}

func (b *Backend) GetSdkContextWithHeader(header *tmtypes.Header) (sdk.Context, error) {
	if header == nil {
		return b.GetSdkContext(), nil
	}
	latestHeight := b.GetBlockStore().Height()
	if latestHeight == 0 {
		return sdk.Context{}, fmt.Errorf("block store not loaded")
	}
	if latestHeight == header.Height {
		return b.copySdkContext(b.ms.CacheMultiStore(), header), nil
	}

	cms, err := b.ms.CacheMultiStoreWithVersion(header.Height)
	if err != nil {
		return sdk.Context{}, err
	}
	return b.copySdkContext(cms, header), nil
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
