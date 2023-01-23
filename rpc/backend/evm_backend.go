package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"

	tmrpctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmjsonrpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/stratosnet/stratos-chain/rpc/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
	tmrpccore "github.com/tendermint/tendermint/rpc/core"
)

var bAttributeKeyEthereumBloom = []byte(evmtypes.AttributeKeyEthereumBloom)

// BlockNumber returns the current block number in abci app state.
// Because abci app state could lag behind from tendermint latest block, it's more stable
// for the client to use the latest block number in abci app state than tendermint rpc.
func (b *Backend) BlockNumber() (hexutil.Uint64, error) {
	res, err := tmrpccore.Block(nil, nil)
	if err != nil {
		return hexutil.Uint64(0), err
	}

	if res.Block == nil {
		return hexutil.Uint64(0), errors.Errorf("block store not loaded")
	}

	return hexutil.Uint64(res.Block.Height), nil
}

// GetBlockByNumber returns the block identified by number.
func (b *Backend) GetBlockByNumber(blockNum types.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	resBlock, err := b.GetTendermintBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	// return if requested block height is greater than the current one
	if resBlock == nil || resBlock.Block == nil {
		return nil, nil
	}

	res, err := b.EthBlockFromTendermint(resBlock.Block, fullTx)
	if err != nil {
		b.logger.Debug("EthBlockFromTendermint failed", "height", blockNum, "error", err.Error())
		return nil, err
	}

	return res, nil
}

// GetBlockByHash returns the block identified by hash.
func (b *Backend) GetBlockByHash(hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	resBlock, err := tmrpccore.BlockByHash(nil, hash.Bytes())
	if err != nil {
		b.logger.Debug("BlockByHash block not found", "hash", hash.Hex(), "error", err.Error())
		return nil, err
	}

	if resBlock == nil || resBlock.Block == nil {
		b.logger.Debug("BlockByHash block not found", "hash", hash.Hex())
		return nil, nil
	}

	return b.EthBlockFromTendermint(resBlock.Block, fullTx)
}

// BlockByNumber returns the block identified by number.
func (b *Backend) BlockByNumber(blockNum types.BlockNumber) (*ethtypes.Block, error) {
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
			return nil, errors.Errorf("incorrect block height: %d", height)
		}
	}

	resBlock, err := tmrpccore.Block(nil, &height)
	if err != nil {
		b.logger.Debug("HeaderByNumber failed", "height", height)
		return nil, err
	}

	if resBlock == nil || resBlock.Block == nil {
		return nil, errors.Errorf("block not found for height %d", height)
	}

	return b.EthBlockFromTm(resBlock.Block)
}

// BlockByHash returns the block identified by hash.
func (b *Backend) BlockByHash(hash common.Hash) (*ethtypes.Block, error) {
	resBlock, err := tmrpccore.BlockByHash(nil, hash.Bytes())
	if err != nil {
		b.logger.Debug("HeaderByHash failed", "hash", hash.Hex())
		return nil, err
	}

	if resBlock == nil || resBlock.Block == nil {
		return nil, errors.Errorf("block not found for hash %s", hash)
	}

	return b.EthBlockFromTm(resBlock.Block)
}

func (b *Backend) EthBlockFromTm(block *tmtypes.Block) (*ethtypes.Block, error) {
	height := block.Height
	bloom, err := b.BlockBloom(&height)
	if err != nil {
		b.logger.Debug("HeaderByNumber BlockBloom failed", "height", height)
	}

	baseFee, err := b.BaseFee()
	if err != nil {
		b.logger.Debug("HeaderByNumber BaseFee failed", "height", height, "error", err.Error())
		return nil, err
	}

	ethHeader := types.EthHeaderFromTendermint(block.Header, bloom, baseFee)

	var txs []*ethtypes.Transaction
	for _, txBz := range block.Txs {
		tx, err := b.clientCtx.TxConfig.TxDecoder()(txBz)
		if err != nil {
			b.logger.Debug("failed to decode transaction in block", "height", height, "error", err.Error())
			continue
		}

		for _, msg := range tx.GetMsgs() {
			ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
			if !ok {
				continue
			}

			tx := ethMsg.AsTransaction()
			txs = append(txs, tx)
		}
	}

	// TODO: add tx receipts
	ethBlock := ethtypes.NewBlock(ethHeader, txs, nil, nil, nil)
	return ethBlock, nil
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
			return nil, errors.Errorf("cannot fetch a negative block height: %d", height)
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
	for _, event := range result.EndBlockEvents {
		if event.Type != evmtypes.EventTypeBlockBloom {
			continue
		}

		for _, attr := range event.Attributes {
			if bytes.Equal(attr.Key, bAttributeKeyEthereumBloom) {
				return ethtypes.BytesToBloom(attr.Value), nil
			}
		}
	}
	return ethtypes.Bloom{}, errors.New("block bloom event is not found")
}

// EthBlockFromTendermint returns a JSON-RPC compatible Ethereum block from a given Tendermint block and its block result.
func (b *Backend) EthBlockFromTendermint(
	block *tmtypes.Block,
	fullTx bool,
) (map[string]interface{}, error) {
	ethRPCTxs := []interface{}{}

	baseFee, err := b.BaseFee()
	if err != nil {
		return nil, err
	}

	resBlockResult, err := tmrpccore.BlockResults(nil, &block.Height)
	if err != nil {
		return nil, err
	}

	txResults := resBlockResult.TxsResults

	for i, tmTx := range block.Txs {
		if !fullTx {
			hash := common.Bytes2Hex(tmTx.Hash())
			ethRPCTxs = append(ethRPCTxs, hash)
			continue
		}

		blockHash := common.BytesToHash(block.Hash())
		blockHeight := uint64(block.Height)
		txIndex := uint64(i)

		rpcTx, err := types.TmTxToEthTx(
			b.clientCtx.TxConfig.TxDecoder(),
			tmTx,
			&blockHash,
			&blockHeight,
			&txIndex,
		)
		if err != nil {
			b.logger.Debug("NewTransactionFromData for receipt failed", "hash", common.Bytes2Hex(tmTx.Hash()), "error", err.Error())
			continue
		}
		ethRPCTxs = append(ethRPCTxs, rpcTx)
	}

	bloom, err := b.BlockBloom(&block.Height)
	if err != nil {
		b.logger.Debug("failed to query BlockBloom", "height", block.Height, "error", err.Error())
	}

	req := &evmtypes.QueryValidatorAccountRequest{
		ConsAddress: sdk.ConsAddress(block.Header.ProposerAddress).String(),
	}

	sdkCtx := b.GetSdkContext(&block.Header)
	res, err := b.GetEVMKeeper().ValidatorAccount(sdk.WrapSDKContext(sdkCtx), req)
	if err != nil {
		b.logger.Debug(
			"failed to query validator operator address",
			"height", block.Height,
			"cons-address", req.ConsAddress,
			"error", err.Error(),
		)
		return nil, err
	}

	addr, err := sdk.AccAddressFromBech32(res.AccountAddress)
	if err != nil {
		return nil, err
	}

	validatorAddr := common.BytesToAddress(addr)

	gasLimit, err := types.BlockMaxGasFromConsensusParams(b.clientCtx, block.Height)
	if err != nil {
		b.logger.Error("failed to query consensus params", "error", err.Error())
	}

	gasUsed := uint64(0)

	for _, txsResult := range txResults {
		// workaround for cosmos-sdk bug. https://github.com/cosmos/cosmos-sdk/issues/10832
		if txsResult.GetCode() == 11 && txsResult.GetLog() == "no block gas left to run tx: out of gas" {
			// block gas limit has exceeded, other txs must have failed with same reason.
			break
		}
		gasUsed += uint64(txsResult.GetGasUsed())
	}

	formattedBlock := types.FormatBlock(
		block.Header, block.Size(),
		gasLimit, new(big.Int).SetUint64(gasUsed),
		ethRPCTxs, bloom, validatorAddr, baseFee,
	)
	return formattedBlock, nil
}

// CurrentHeader returns the latest block header
func (b *Backend) CurrentHeader() *ethtypes.Header {
	header, _ := b.HeaderByNumber(types.EthLatestBlockNumber)
	return header
}

// HeaderByNumber returns the block header identified by height.
func (b *Backend) HeaderByNumber(blockNum types.BlockNumber) (*ethtypes.Header, error) {
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
			return nil, errors.Errorf("incorrect block height: %d", height)
		}
	}

	resBlock, err := tmrpccore.Block(nil, &height)
	if err != nil {
		b.logger.Debug("HeaderByNumber failed")
		return nil, err
	}

	bloom, err := b.BlockBloom(&resBlock.Block.Height)
	if err != nil {
		b.logger.Debug("HeaderByNumber BlockBloom failed", "height", resBlock.Block.Height)
	}

	baseFee, err := b.BaseFee()
	if err != nil {
		b.logger.Debug("HeaderByNumber BaseFee failed", "height", resBlock.Block.Height, "error", err.Error())
		return nil, err
	}

	ethHeader := types.EthHeaderFromTendermint(resBlock.Block.Header, bloom, baseFee)
	return ethHeader, nil
}

// HeaderByHash returns the block header identified by hash.
func (b *Backend) HeaderByHash(blockHash common.Hash) (*ethtypes.Header, error) {
	resBlock, err := tmrpccore.BlockByHash(nil, blockHash.Bytes())
	if err != nil {
		b.logger.Debug("HeaderByHash failed", "hash", blockHash.Hex())
		return nil, err
	}

	if resBlock == nil || resBlock.Block == nil {
		return nil, errors.Errorf("block not found for hash %s", blockHash.Hex())
	}

	bloom, err := b.BlockBloom(&resBlock.Block.Height)
	if err != nil {
		b.logger.Debug("HeaderByHash BlockBloom failed", "height", resBlock.Block.Height)
	}

	baseFee, err := b.BaseFee()
	if err != nil {
		b.logger.Debug("HeaderByHash BaseFee failed", "height", resBlock.Block.Height, "error", err.Error())
		return nil, err
	}

	ethHeader := types.EthHeaderFromTendermint(resBlock.Block.Header, bloom, baseFee)
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

	blockLogs := [][]*ethtypes.Log{}
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
			return nil, errors.Errorf("incorrect block height: %d", height)
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

	res, err := b.GetEVMKeeper().ValidatorAccount(b.ctx, req)
	if err != nil {
		return nil, err
	}

	address, _ := sdk.AccAddressFromBech32(res.AccountAddress)
	return address, nil
}

// GetTransactionByHash returns the Ethereum format transaction identified by Ethereum transaction hash
func (b *Backend) GetTransactionByHash(txHash common.Hash) (*types.RPCTransaction, error) {
	res, err := b.GetTxByHash(txHash)
	// fmt.Printf("debug res TX structure: %+v\n", sdk.NewResponseResultTx(res, nil, "").String())
	if err != nil {
		// TODO: Get chain id value from genesis
		tx, err := types.GetPendingTx(b.clientCtx.TxConfig.TxDecoder(), b.GetMempool(), txHash, b.ChainConfig().ChainID)
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

func (b *Backend) GetTxByHash(hash common.Hash) (*tmrpctypes.ResultTx, error) {
	resTx, err := tmrpccore.Tx(nil, hash.Bytes(), false)
	if err != nil {
		query := fmt.Sprintf("%s.%s='%s'", evmtypes.TypeMsgEthereumTx, evmtypes.AttributeKeyEthereumTxHash, hash.Hex())
		resTxs, err := tmrpccore.TxSearch(new(tmjsonrpctypes.Context), query, false, nil, nil, "")
		if err != nil {
			return nil, err
		}
		if len(resTxs.Txs) == 0 {
			return nil, errors.Errorf("ethereum tx not found for hash %s", hash.Hex())
		}
		return resTxs.Txs[0], nil
	}
	return resTx, nil
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
		return nil, errors.Errorf("ethereum tx not found for block %d index %d", height, index)
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
	sdkCtx := b.GetSdkContext(nil)
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
		return common.Hash{}, errors.New("legacy pre-eip-155 transactions not supported")
	}

	txHash := ethTx.Hash()

	// Broadcast transaction in sync mode (default)
	// NOTE: If error is encountered on the node, the broadcast will not return an error
	syncCtx := b.clientCtx.WithBroadcastMode(flags.BroadcastSync)
	rsp, err := syncCtx.BroadcastTx(txBytes)
	if rsp != nil && rsp.Code != 0 {
		err = sdkerrors.ABCIError(rsp.Codespace, rsp.Code, rsp.RawLog)
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
	sdkCtx := b.GetSdkContext(&resBlock.Block.Header)
	res, err := b.GetEVMKeeper().EstimateGas(sdk.WrapSDKContext(sdkCtx), &req)
	if err != nil {
		return 0, err
	}
	return hexutil.Uint64(res.Gas), nil
}

// GetTransactionCount returns the number of transactions at the given address up to the given block number.
func (b *Backend) GetTransactionCount(address common.Address, blockNum types.BlockNumber) (*hexutil.Uint64, error) {
	// Get nonce (sequence) from account
	from := sdk.AccAddress(address.Bytes())
	accRet := b.clientCtx.AccountRetriever

	err := accRet.EnsureExists(b.clientCtx, from)
	if err != nil {
		// account doesn't exist yet, return 0
		n := hexutil.Uint64(0)
		return &n, nil
	}

	includePending := blockNum == types.EthPendingBlockNumber
	nonce := b.getAccountNonce(address, includePending, blockNum.Int64())
	n := hexutil.Uint64(nonce)
	return &n, nil
}

// RPCGasCap is the global gas cap for eth-call variants.
func (b *Backend) RPCGasCap() uint64 {
	return b.cfg.JSONRPC.GasCap
}

// RPCEVMTimeout is the global evm timeout for eth-call variants.
func (b *Backend) RPCEVMTimeout() time.Duration {
	return b.cfg.JSONRPC.EVMTimeout
}

// RPCGasCap is the global gas cap for eth-call variants.
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
	sdkCtx := b.GetSdkContext(nil)
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
	sdkCtx := b.GetSdkContext(nil)
	params, err := b.GetEVMKeeper().Params(sdk.WrapSDKContext(sdkCtx), &evmtypes.QueryParamsRequest{})
	if err != nil {
		return nil
	}

	return params.Params.ChainConfig.EthereumConfig()
}

// SuggestGasTipCap returns the suggested tip cap
// Although we don't support tx prioritization yet, but we return a positive value to help client to
// mitigate the base fee changes.
func (b *Backend) SuggestGasTipCap() (*big.Int, error) {
	baseFee, err := b.BaseFee()
	if err != nil {
		// london hardfork not enabled or feemarket not enabled
		return big.NewInt(0), nil
	}

	sdkCtx := b.GetSdkContext(nil)
	params, err := b.GetEVMKeeper().Params(sdk.WrapSDKContext(sdkCtx), &evmtypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}
	// calculate the maximum base fee delta in current block, assuming all block gas limit is consumed
	// ```
	// GasTarget = GasLimit / ElasticityMultiplier
	// Delta = BaseFee * (GasUsed - GasTarget) / GasTarget / Denominator
	// ```
	// The delta is at maximum when `GasUsed` is equal to `GasLimit`, which is:
	// ```
	// MaxDelta = BaseFee * (GasLimit - GasLimit / ElasticityMultiplier) / (GasLimit / ElasticityMultiplier) / Denominator
	//          = BaseFee * (ElasticityMultiplier - 1) / Denominator
	// ```
	maxDelta := baseFee.Int64() * (int64(params.Params.FeeMarketParams.ElasticityMultiplier) - 1) / int64(params.Params.FeeMarketParams.BaseFeeChangeDenominator)
	if maxDelta < 0 {
		// impossible if the parameter validation passed.
		maxDelta = 0
	}
	return big.NewInt(maxDelta), nil
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

	sdkCtx := b.GetSdkContext(&resBlock.Block.Header)
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
		err = b.processBlock(tendermintblock, &ethBlock, rewardPercentiles, tendermintBlockResult, &oneFeeHistory)
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
