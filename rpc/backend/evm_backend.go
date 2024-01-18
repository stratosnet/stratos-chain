package backend

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"

	tmrpccore "github.com/cometbft/cometbft/rpc/core"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	tmjsonrpctypes "github.com/cometbft/cometbft/rpc/jsonrpc/types"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/rpc/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
)

// BlockNumber returns the current block number in abci app state.
// Because abci app state could lag behind from tendermint latest block, it's more stable
// for the client to use the latest block number in abci app state than tendermint rpc.
func (b *Backend) BlockNumber() (hexutil.Uint64, error) {
	res, err := tmrpccore.Block(nil, nil)
	if err != nil {
		return hexutil.Uint64(0), err
	}

	if res.Block == nil {
		return hexutil.Uint64(0), fmt.Errorf("block store not loaded")
	}

	return hexutil.Uint64(res.Block.Height), nil
}

func (b *Backend) getBlockFromResultBlock(resBlock *tmrpctypes.ResultBlock, fullTx bool) (*types.Block, error) {
	// return if requested block height is greater than the current one
	if resBlock == nil || resBlock.Block == nil {
		return nil, nil
	}

	res, err := types.EthBlockFromTendermint(b.clientCtx.TxConfig.TxDecoder(), resBlock.Block, fullTx)
	if err != nil {
		b.logger.Debug("EthBlockFromTendermint failed", "height", resBlock.Block.Height, "hash", resBlock.Block.Hash(), "error", err.Error())
		return nil, nil
	}

	// override dynamicly miner address
	sdkCtx, err := b.GetEVMContext().GetSdkContextWithHeader(&resBlock.Block.Header)
	if err != nil {
		b.logger.Debug("GetSdkContextWithHeader context", "height", resBlock.Block.Height, "hash", resBlock.Block.Hash(), "error", err.Error())
		return nil, err
	}

	validator, err := b.GetEVMKeeper().GetCoinbaseAddress(sdkCtx)
	if err != nil {
		b.logger.Debug("GetCoinbaseAddress no validator", "height", resBlock.Block.Height, "hash", resBlock.Block.Hash(), "error", err.Error())
		return nil, err
	}
	res.Miner = validator

	feeResp, err := b.GetEVMKeeper().BaseFee(sdk.WrapSDKContext(sdkCtx), nil)
	if err != nil {
		return nil, err
	}
	res.BaseFee = (*hexutil.Big)(feeResp.BaseFee.BigInt())

	return res, nil
}

// GetBlockByNumber returns the block identified by number.
func (b *Backend) GetBlockByNumber(blockNum types.BlockNumber, fullTx bool) (*types.Block, error) {
	resBlock, err := b.GetTendermintBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}
	return b.getBlockFromResultBlock(resBlock, fullTx)
}

// GetBlockByHash returns the block identified by hash.
func (b *Backend) GetBlockByHash(hash common.Hash, fullTx bool) (*types.Block, error) {
	resBlock, err := tmrpccore.BlockByHash(nil, hash.Bytes())
	if err != nil {
		b.logger.Debug("BlockByHash block not found", "hash", hash.Hex(), "error", err.Error())
		return nil, err
	}
	return b.getBlockFromResultBlock(resBlock, fullTx)
}

// GetTendermintBlockByNumber returns a Tendermint format block by block number
func (b *Backend) GetTendermintBlockByNumber(blockNum types.BlockNumber) (*tmrpctypes.ResultBlock, error) {
	height := blockNum.Int64()
	currentBlockNumber, _ := b.BlockNumber()

	switch blockNum {
	case types.EthLatestBlockNumber:
		if currentBlockNumber > 0 {
			height = int64(currentBlockNumber)
		}
	case types.EthPendingBlockNumber:
		if currentBlockNumber > 0 {
			height = int64(currentBlockNumber)
		}
	case types.EthEarliestBlockNumber:
		height = 1
	default:
		if blockNum < 0 {
			return nil, fmt.Errorf("cannot fetch a negative block height: %d", height)
		}
		if height > int64(currentBlockNumber) {
			return nil, nil
		}
	}

	resBlock, err := tmrpccore.Block(nil, &height)
	if err != nil {
		if resBlock, err = tmrpccore.Block(nil, nil); err != nil {
			b.logger.Debug("tendermint client failed to get latest block", "height", height, "error", err.Error())
			return nil, nil
		}
	}

	if resBlock.Block == nil {
		b.logger.Debug("GetBlockByNumber block not found", "height", height)
		return nil, nil
	}

	return resBlock, nil
}

// GetTendermintBlockByHash returns a Tendermint format block by block number
func (b *Backend) GetTendermintBlockByHash(blockHash common.Hash) (*tmrpctypes.ResultBlock, error) {
	resBlock, err := tmrpccore.BlockByHash(nil, blockHash.Bytes())
	if err != nil {
		b.logger.Debug("tendermint client failed to get block", "blockHash", blockHash.Hex(), "error", err.Error())
	}

	if resBlock == nil || resBlock.Block == nil {
		b.logger.Debug("GetBlockByNumber block not found", "blockHash", blockHash.Hex())
		return nil, nil
	}

	return resBlock, nil
}

// BlockBloom query block bloom filter from block results
func (b *Backend) BlockBloom(height *int64) (ethtypes.Bloom, error) {
	result, err := tmrpccore.BlockResults(nil, height)
	if err != nil {
		return ethtypes.Bloom{}, err
	}
	return types.GetBlockBloom(result)
}

// CurrentHeader returns the latest block header
func (b *Backend) CurrentHeader() *types.Header {
	header, _ := b.HeaderByNumber(types.EthLatestBlockNumber)
	return header
}

// HeaderByNumber returns the block header identified by height.
func (b *Backend) HeaderByNumber(blockNum types.BlockNumber) (*types.Header, error) {
	height := blockNum.Int64()

	switch blockNum {
	case types.EthLatestBlockNumber:
		currentBlockNumber, _ := b.BlockNumber()
		if currentBlockNumber > 0 {
			height = int64(currentBlockNumber)
		}
	case types.EthPendingBlockNumber:
		currentBlockNumber, _ := b.BlockNumber()
		if currentBlockNumber > 0 {
			height = int64(currentBlockNumber)
		}
	case types.EthEarliestBlockNumber:
		height = 1
	default:
		if blockNum < 0 {
			return nil, fmt.Errorf("incorrect block height: %d", height)
		}
	}

	resBlock, err := tmrpccore.Block(nil, &height)
	if err != nil {
		b.logger.Debug("HeaderByNumber failed")
		return nil, err
	}

	ethHeader, err := types.EthHeaderFromTendermint(resBlock.Block.Header)
	if err != nil {
		b.logger.Debug("HeaderByNumber EthHeaderFromTendermint failed", "height", resBlock.Block.Height, "error", err.Error())
		return nil, err
	}

	// override dynamicly miner address
	sdkCtx, err := b.GetEVMContext().GetSdkContextWithHeader(&resBlock.Block.Header)
	if err != nil {
		b.logger.Debug("GetSdkContextWithHeader context", "height", blockNum, "error", err.Error())
		return nil, err
	}

	validator, err := b.evmkeeper.GetCoinbaseAddress(sdkCtx)
	if err != nil {
		b.logger.Debug("EthBlockFromTendermint no validator", "height", blockNum, "error", err.Error())
		return nil, err
	}
	ethHeader.Coinbase = validator

	feeResp, err := b.GetEVMKeeper().BaseFee(sdk.WrapSDKContext(sdkCtx), nil)
	if err != nil {
		return nil, err
	}
	ethHeader.BaseFee = feeResp.BaseFee.BigInt()
	return ethHeader, nil
}

// HeaderByHash returns the block header identified by hash.
func (b *Backend) HeaderByHash(blockHash common.Hash) (*types.Header, error) {
	resBlock, err := tmrpccore.BlockByHash(nil, blockHash.Bytes())
	if err != nil {
		b.logger.Debug("HeaderByHash failed", "hash", blockHash.Hex())
		return nil, err
	}

	if resBlock == nil || resBlock.Block == nil {
		return nil, fmt.Errorf("block not found for hash %s", blockHash.Hex())
	}

	ethHeader, err := types.EthHeaderFromTendermint(resBlock.Block.Header)
	if err != nil {
		b.logger.Debug("HeaderByHash EthHeaderFromTendermint failed", "height", resBlock.Block.Height, "error", err.Error())
		return nil, err
	}

	// override dynamicly miner address
	sdkCtx, err := b.GetEVMContext().GetSdkContextWithHeader(&resBlock.Block.Header)
	if err != nil {
		b.logger.Debug("GetSdkContextWithHeader context", "hash", blockHash.Hex(), "error", err.Error())
		return nil, err
	}

	validator, err := b.evmkeeper.GetCoinbaseAddress(sdkCtx)
	if err != nil {
		b.logger.Debug("EthBlockFromTendermint no validator", "hash", blockHash.Hex(), "error", err.Error())
		return nil, err
	}
	ethHeader.Coinbase = validator

	feeResp, err := b.GetEVMKeeper().BaseFee(sdk.WrapSDKContext(sdkCtx), nil)
	if err != nil {
		return nil, err
	}
	ethHeader.BaseFee = feeResp.BaseFee.BigInt()
	return ethHeader, nil
}

// PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (b *Backend) PendingTransactions() ([]*sdk.Tx, error) {
	txs := b.GetMempool().ReapMaxTxs(100)

	result := make([]*sdk.Tx, 0, len(txs))
	for _, txBz := range txs {
		tx, err := b.clientCtx.TxConfig.TxDecoder()(txBz)
		if err != nil {
			return nil, err
		}
		result = append(result, &tx)
	}

	return result, nil
}

// GetLogsByHeight returns all the logs from all the ethereum transactions in a block.
func (b *Backend) GetLogsByHeight(height *int64) ([][]*ethtypes.Log, error) {
	// NOTE: we query the state in case the tx result logs are not persisted after an upgrade.
	blockRes, err := tmrpccore.BlockResults(nil, height)
	if err != nil {
		return nil, err
	}

	blockLogs := make([][]*ethtypes.Log, 0)
	for _, txResult := range blockRes.TxsResults {
		logs, err := AllTxLogsFromEvents(txResult.Events)
		if err != nil {
			return nil, err
		}

		blockLogs = append(blockLogs, logs...)
	}

	return blockLogs, nil
}

// GetLogs returns all the logs from all the ethereum transactions in a block.
func (b *Backend) GetLogs(hash common.Hash) ([][]*ethtypes.Log, error) {
	block, err := tmrpccore.BlockByHash(nil, hash.Bytes())
	if err != nil {
		return nil, err
	}
	return b.GetLogsByHeight(&block.Block.Header.Height)
}

func (b *Backend) GetLogsByNumber(blockNum types.BlockNumber) ([][]*ethtypes.Log, error) {
	height := blockNum.Int64()

	switch blockNum {
	case types.EthLatestBlockNumber:
		currentBlockNumber, _ := b.BlockNumber()
		if currentBlockNumber > 0 {
			height = int64(currentBlockNumber)
		}
	case types.EthPendingBlockNumber:
		currentBlockNumber, _ := b.BlockNumber()
		if currentBlockNumber > 0 {
			height = int64(currentBlockNumber)
		}
	case types.EthEarliestBlockNumber:
		height = 1
	default:
		if blockNum < 0 {
			return nil, fmt.Errorf("incorrect block height: %d", height)
		}
	}

	return b.GetLogsByHeight(&height)
}

// BloomStatus returns the BloomBitsBlocks and the number of processed sections maintained
// by the chain indexer.
func (b *Backend) BloomStatus() (uint64, uint64) {
	return 4096, 0
}

// GetCoinbase is the address that staking rewards will be send to (alias for Etherbase).
func (b *Backend) GetCoinbase() (sdk.AccAddress, error) {
	status, err := tmrpccore.Status(nil)
	if err != nil {
		return nil, err
	}

	req := &evmtypes.QueryValidatorAccountRequest{
		ConsAddress: sdk.ConsAddress(status.ValidatorInfo.Address).String(),
	}

	ctx := b.GetEVMContext().GetSdkContext()
	res, err := b.GetEVMKeeper().ValidatorAccount(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return nil, err
	}

	address, _ := sdk.AccAddressFromBech32(res.AccountAddress)
	return address, nil
}

// GetTransactionByHash returns the Ethereum format transaction identified by Ethereum transaction hash
func (b *Backend) GetTransactionByHash(txHash common.Hash) (*types.Transaction, error) {
	res, err := evmtypes.GetTmTxByHash(txHash)
	if err != nil {
		// TODO: Get chain id value from genesis
		tx, err := types.GetPendingTx(b.clientCtx.TxConfig.TxDecoder(), b.GetMempool(), txHash)
		if err != nil {
			b.logger.Debug("tx not found", "hash", txHash, "error", err.Error())
			return nil, nil
		}
		return tx, nil
	}

	block := b.GetBlockStore().LoadBlock(res.Height)
	if block == nil {
		b.logger.Debug("eth_getTransactionByHash", "hash", txHash, "block not found")
		return nil, err
	}

	blockHash := common.BytesToHash(block.Hash())
	blockHeight := uint64(res.Height)
	txIndex := uint64(res.Index)

	return types.TmTxToEthTx(
		b.clientCtx.TxConfig.TxDecoder(),
		res.Tx,
		&blockHash,
		&blockHeight,
		&txIndex,
	)
}

// GetTxByTxIndex uses `/tx_query` to find transaction by tx index of valid ethereum txs
func (b *Backend) GetTxByTxIndex(height int64, index uint) (*tmrpctypes.ResultTx, error) {
	query := fmt.Sprintf("tx.height=%d AND %s.%s=%d",
		height, evmtypes.TypeMsgEthereumTx,
		evmtypes.AttributeKeyTxIndex, index,
	)
	resTxs, err := tmrpccore.TxSearch(new(tmjsonrpctypes.Context), query, false, nil, nil, "")
	if err != nil {
		return nil, err
	}
	if len(resTxs.Txs) == 0 {
		return nil, fmt.Errorf("ethereum tx not found for block %d index %d", height, index)
	}
	return resTxs.Txs[0], nil
}

func (b *Backend) SendTransaction(args evmtypes.TransactionArgs) (common.Hash, error) {
	// Look up the wallet containing the requested signer
	_, err := b.clientCtx.Keyring.KeyByAddress(sdk.AccAddress(args.From.Bytes()))
	if err != nil {
		b.logger.Error("failed to find key in keyring", "address", args.From, "error", err.Error())
		return common.Hash{}, fmt.Errorf("%s; %s", keystore.ErrNoMatch, err.Error())
	}

	args, err = b.SetTxDefaults(args)
	if err != nil {
		return common.Hash{}, err
	}

	msg := args.ToTransaction()
	if err := msg.ValidateBasic(); err != nil {
		b.logger.Debug("tx failed basic validation", "error", err.Error())
		return common.Hash{}, err
	}

	bn, err := b.BlockNumber()
	if err != nil {
		b.logger.Debug("failed to fetch latest block number", "error", err.Error())
		return common.Hash{}, err
	}

	signer := ethtypes.MakeSigner(b.ChainConfig(), new(big.Int).SetUint64(uint64(bn)))

	// Sign transaction
	if err := msg.Sign(signer, b.clientCtx.Keyring); err != nil {
		b.logger.Debug("failed to sign tx", "error", err.Error())
		return common.Hash{}, err
	}

	// Query params to use the EVM denomination
	sdkCtx := b.GetEVMContext().GetSdkContext()
	res, err := b.GetEVMKeeper().Params(sdk.WrapSDKContext(sdkCtx), &evmtypes.QueryParamsRequest{})
	if err != nil {
		b.logger.Error("failed to query evm params", "error", err.Error())
		return common.Hash{}, err
	}

	// Assemble transaction from fields
	tx, err := msg.BuildTx(b.clientCtx.TxConfig.NewTxBuilder(), res.Params.EvmDenom)
	if err != nil {
		b.logger.Error("build cosmos tx failed", "error", err.Error())
		return common.Hash{}, err
	}

	// Encode transaction by default Tx encoder
	txEncoder := b.clientCtx.TxConfig.TxEncoder()
	txBytes, err := txEncoder(tx)
	if err != nil {
		b.logger.Error("failed to encode eth tx using default encoder", "error", err.Error())
		return common.Hash{}, err
	}

	ethTx := msg.AsTransaction()
	if !ethTx.Protected() {
		// Ensure only eip155 signed transactions are submitted.
		return common.Hash{}, fmt.Errorf("legacy pre-eip-155 transactions not supported")
	}

	txHash := ethTx.Hash()

	// Broadcast transaction in sync mode (default)
	// NOTE: If error is encountered on the node, the broadcast will not return an error
	syncCtx := b.clientCtx.WithBroadcastMode(flags.BroadcastSync)
	rsp, err := syncCtx.BroadcastTx(txBytes)
	if rsp != nil && rsp.Code != abci.CodeTypeOK {
		err = errors.ABCIError(rsp.Codespace, rsp.Code, rsp.RawLog)
	}
	if err != nil {
		b.logger.Error("failed to broadcast tx", "error", err.Error())
		return txHash, err
	}

	// Return transaction hash
	return txHash, nil
}

// EstimateGas returns an estimate of gas usage for the given smart contract call.
func (b *Backend) EstimateGas(args evmtypes.TransactionArgs, blockNrOptional *types.BlockNumber) (hexutil.Uint64, error) {
	blockNr := types.EthPendingBlockNumber
	if blockNrOptional != nil {
		blockNr = *blockNrOptional
	}

	bz, err := json.Marshal(&args)
	if err != nil {
		return 0, err
	}

	req := evmtypes.EthCallRequest{
		Args:   bz,
		GasCap: b.RPCGasCap(),
	}

	resBlock, err := b.GetTendermintBlockByNumber(blockNr)
	if err != nil {
		return 0, err
	}

	// return if requested block height is greater than the current one or chain not synced
	if resBlock == nil || resBlock.Block == nil {
		return 0, nil
	}

	// it will return an empty context and the sdk.Context will use
	// the latest block height for querying.
	sdkCtx, err := b.GetEVMContext().GetSdkContextWithHeader(&resBlock.Block.Header)
	if err != nil {
		return 0, err
	}
	res, err := b.GetEVMKeeper().EstimateGas(sdk.WrapSDKContext(sdkCtx), &req)
	if err != nil {
		return 0, err
	}
	return hexutil.Uint64(res.Gas), nil
}

// GetTransactionCount returns the number of transactions at the given address up to the given block number.
func (b *Backend) GetTransactionCount(address common.Address, blockNum types.BlockNumber) (hexutil.Uint64, error) {
	// Get nonce (sequence) from account
	from := sdk.AccAddress(address.Bytes())
	accRet := b.clientCtx.AccountRetriever

	err := accRet.EnsureExists(b.clientCtx, from)
	if err != nil {
		return hexutil.Uint64(0), nil
	}

	nonce, err := b.getAccountNonce(address, blockNum)
	if err != nil {
		return hexutil.Uint64(0), err
	}
	return hexutil.Uint64(nonce), nil
}

// RPCGasCap is the global gas cap for eth-call variants.
func (b *Backend) RPCGasCap() uint64 {
	return b.cfg.JSONRPC.GasCap
}

// RPCEVMTimeout is the global evm timeout for eth-call variants.
func (b *Backend) RPCEVMTimeout() time.Duration {
	return b.cfg.JSONRPC.EVMTimeout
}

// RPCTxFeeCap is the global gas cap for eth-call variants.
func (b *Backend) RPCTxFeeCap() float64 {
	return b.cfg.JSONRPC.TxFeeCap
}

// RPCFilterCap is the limit for total number of filters that can be created
func (b *Backend) RPCFilterCap() int32 {
	return b.cfg.JSONRPC.FilterCap
}

// RPCFeeHistoryCap is the limit for total number of blocks that can be fetched
func (b *Backend) RPCFeeHistoryCap() int32 {
	return b.cfg.JSONRPC.FeeHistoryCap
}

// RPCLogsCap defines the max number of results can be returned from single `eth_getLogs` query.
func (b *Backend) RPCLogsCap() int32 {
	return b.cfg.JSONRPC.LogsCap
}

// RPCBlockRangeCap defines the max block range allowed for `eth_getLogs` query.
func (b *Backend) RPCBlockRangeCap() int32 {
	return b.cfg.JSONRPC.BlockRangeCap
}

// RPCMinGasPrice returns the minimum gas price for a transaction obtained from
// the node config. If set value is 0, it will default to 20.

func (b *Backend) RPCMinGasPrice() int64 {
	sdkCtx := b.GetEVMContext().GetSdkContext()
	evmParams, err := b.GetEVMKeeper().Params(sdk.WrapSDKContext(sdkCtx), &evmtypes.QueryParamsRequest{})
	if err != nil {
		return stratos.DefaultGasPrice
	}

	minGasPrice := b.cfg.GetMinGasPrices()
	amt := minGasPrice.AmountOf(evmParams.Params.EvmDenom).TruncateInt64()
	if amt == 0 {
		return stratos.DefaultGasPrice
	}

	return amt
}

// ChainConfig returns the latest ethereum chain configuration
func (b *Backend) ChainConfig() *params.ChainConfig {
	sdkCtx := b.GetEVMContext().GetSdkContext()
	evmParams, err := b.GetEVMKeeper().Params(sdk.WrapSDKContext(sdkCtx), &evmtypes.QueryParamsRequest{})
	if err != nil {
		return nil
	}

	return evmParams.Params.ChainConfig.EthereumConfig()
}

// SuggestGasTipCap returns the suggested tip cap
// Although we don't support tx prioritization yet, but we return a positive value to help client to
// mitigate the base fee changes.
func (b *Backend) SuggestGasTipCap() (*big.Int, error) {
	// baseFee, err := b.BaseFee()
	// if err != nil {
	// 	// london hardfork not enabled or feemarket not enabled
	// 	return big.NewInt(0), nil
	// }

	// sdkCtx := b.GetEVMContext().GetSdkContext()
	// params, err := b.GetEVMKeeper().Params(sdk.WrapSDKContext(sdkCtx), &evmtypes.QueryParamsRequest{})
	// if err != nil {
	// 	return nil, err
	// }
	// // calculate the maximum base fee delta in current block, assuming all block gas limit is consumed
	// // ```
	// // GasTarget = GasLimit / ElasticityMultiplier
	// // Delta = BaseFee * (GasUsed - GasTarget) / GasTarget / Denominator
	// // ```
	// // The delta is at maximum when `GasUsed` is equal to `GasLimit`, which is:
	// // ```
	// // MaxDelta = BaseFee * (GasLimit - GasLimit / ElasticityMultiplier) / (GasLimit / ElasticityMultiplier) / Denominator
	// //          = BaseFee * (ElasticityMultiplier - 1) / Denominator
	// // ```
	// maxDelta := baseFee.Int64() * (int64(params.Params.FeeMarketParams.ElasticityMultiplier) - 1) / int64(params.Params.FeeMarketParams.BaseFeeChangeDenominator)
	// if maxDelta < 0 {
	// 	// impossible if the parameter validation passed.
	// 	maxDelta = 0
	// }
	// return big.NewInt(maxDelta), nil

	// NOTE: Commented as validators do not receive tips
	// but I left a logic in case we want to have this in future
	return big.NewInt(0), nil
}

// BaseFee returns the base fee tracked by the Fee Market module.
// If the base fee is not enabled globally, the query returns nil.
// If the London hard fork is not activated at the current height, the query will
// return nil.
func (b *Backend) BaseFee() (*big.Int, error) {
	resBlock, err := b.GetTendermintBlockByNumber(types.EthLatestBlockNumber)
	if err != nil {
		return nil, err
	}

	// return if requested block height is greater than the current one or chain not synced
	if resBlock == nil || resBlock.Block == nil {
		return nil, nil
	}

	sdkCtx, err := b.GetEVMContext().GetSdkContextWithHeader(&resBlock.Block.Header)
	if err != nil {
		return nil, err
	}
	// return BaseFee if London hard fork is activated and feemarket is enabled
	res, err := b.GetEVMKeeper().BaseFee(sdk.WrapSDKContext(sdkCtx), nil)
	if err != nil {
		return nil, err
	}

	if res.BaseFee == nil {
		return nil, nil
	}

	return res.BaseFee.BigInt(), nil
}

// FeeHistory returns data relevant for fee estimation based on the specified range of blocks.
func (b *Backend) FeeHistory(
	userBlockCount rpc.DecimalOrHex, // number blocks to fetch, maximum is 100
	lastBlock rpc.BlockNumber, // the block to start search , to oldest
	rewardPercentiles []float64, // percentiles to fetch reward
) (*types.FeeHistoryResult, error) {
	blockEnd := int64(lastBlock)

	if blockEnd <= 0 {
		blockNumber, err := b.BlockNumber()
		if err != nil {
			return nil, err
		}
		blockEnd = int64(blockNumber)
	}
	userBlockCountInt := int64(userBlockCount)
	maxBlockCount := int64(b.cfg.JSONRPC.FeeHistoryCap)
	if userBlockCountInt > maxBlockCount {
		return nil, fmt.Errorf("FeeHistory user block count %d higher than %d", userBlockCountInt, maxBlockCount)
	}
	blockStart := blockEnd - userBlockCountInt
	if blockStart < 0 {
		blockStart = 0
	}

	blockCount := blockEnd - blockStart

	oldestBlock := (*hexutil.Big)(big.NewInt(blockStart))

	// prepare space
	reward := make([][]*hexutil.Big, blockCount)
	rewardCount := len(rewardPercentiles)
	for i := 0; i < int(blockCount); i++ {
		reward[i] = make([]*hexutil.Big, rewardCount)
	}
	thisBaseFee := make([]*hexutil.Big, blockCount)
	thisGasUsedRatio := make([]float64, blockCount)

	// rewards should only be calculated if reward percentiles were included
	calculateRewards := rewardCount != 0

	// fetch block
	for blockID := blockStart; blockID < blockEnd; blockID++ {
		index := int32(blockID - blockStart)
		// eth block
		ethBlock, err := b.GetBlockByNumber(types.BlockNumber(blockID), true)
		if ethBlock == nil {
			return nil, err
		}

		// tendermint block
		tendermintblock, err := b.GetTendermintBlockByNumber(types.BlockNumber(blockID))
		if tendermintblock == nil {
			return nil, err
		}

		// tendermint block result
		tendermintBlockResult, err := tmrpccore.BlockResults(nil, &tendermintblock.Block.Height)
		if tendermintBlockResult == nil {
			b.logger.Debug("block result not found", "height", tendermintblock.Block.Height, "error", err.Error())
			return nil, err
		}

		oneFeeHistory := types.OneFeeHistory{}
		err = b.processBlock(tendermintblock, ethBlock, rewardPercentiles, tendermintBlockResult, &oneFeeHistory)
		if err != nil {
			return nil, err
		}

		// copy
		thisBaseFee[index] = (*hexutil.Big)(oneFeeHistory.BaseFee)
		thisGasUsedRatio[index] = oneFeeHistory.GasUsedRatio
		if calculateRewards {
			for j := 0; j < rewardCount; j++ {
				reward[index][j] = (*hexutil.Big)(oneFeeHistory.Reward[j])
				if reward[index][j] == nil {
					reward[index][j] = (*hexutil.Big)(big.NewInt(0))
				}
			}
		}
	}

	feeHistory := types.FeeHistoryResult{
		OldestBlock:  oldestBlock,
		BaseFee:      thisBaseFee,
		GasUsedRatio: thisGasUsedRatio,
	}

	if calculateRewards {
		feeHistory.Reward = reward
	}

	return &feeHistory, nil
}
