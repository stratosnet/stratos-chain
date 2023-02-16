package eth

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/log"
	tmrpctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/stratosnet/stratos-chain/crypto/hd"
	"github.com/stratosnet/stratos-chain/ethereum/eip712"
	"github.com/stratosnet/stratos-chain/rpc/backend"
	"github.com/stratosnet/stratos-chain/rpc/types"
	rpctypes "github.com/stratosnet/stratos-chain/rpc/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/evm/pool"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
	tmrpccore "github.com/tendermint/tendermint/rpc/core"
)

// PublicAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicAPI struct {
	ctx          context.Context
	clientCtx    client.Context
	chainIDEpoch *big.Int
	logger       log.Logger
	backend      backend.BackendI
	nonceLock    *rpctypes.AddrLocker
	signer       ethtypes.Signer
}

// NewPublicAPI creates an instance of the public ETH Web3 API.
func NewPublicAPI(
	logger log.Logger,
	clientCtx client.Context,
	backend backend.BackendI,
	nonceLock *rpctypes.AddrLocker,
) *PublicAPI {
	algos, _ := clientCtx.Keyring.SupportedAlgorithms()

	if !algos.Contains(hd.EthSecp256k1) {
		kr, err := keyring.New(
			sdk.KeyringServiceName(),
			viper.GetString(flags.FlagKeyringBackend),
			clientCtx.KeyringDir,
			clientCtx.Input,
			hd.EthSecp256k1Option(),
		)
		if err != nil {
			panic(err)
		}

		clientCtx = clientCtx.WithKeyring(kr)
	}

	// The signer used by the API should always be the 'latest' known one because we expect
	// signers to be backwards-compatible with old transactions.
	cfg := backend.ChainConfig()
	if cfg == nil {
		cfg = evmtypes.DefaultChainConfig().EthereumConfig()
	}

	signer := ethtypes.LatestSigner(cfg)

	api := &PublicAPI{
		ctx:          context.Background(),
		clientCtx:    clientCtx,
		chainIDEpoch: cfg.ChainID,
		logger:       logger.With("client", "json-rpc"),
		backend:      backend,
		nonceLock:    nonceLock,
		signer:       signer,
	}

	return api
}

// ClientCtx returns client context
func (e *PublicAPI) ClientCtx() client.Context {
	return e.clientCtx
}

func (e *PublicAPI) Ctx() context.Context {
	return e.ctx
}

// ProtocolVersion returns the supported Ethereum protocol version.
func (e *PublicAPI) ProtocolVersion() hexutil.Uint {
	e.logger.Debug("eth_protocolVersion")
	return hexutil.Uint(stratos.ProtocolVersion)
}

// ChainId is the EIP-155 replay-protection chain id for the current ethereum chain config.
func (e *PublicAPI) ChainId() (*hexutil.Big, error) { // nolint
	e.logger.Debug("eth_chainId")
	ctx := e.backend.GetEVMContext().GetSdkContext()
	params := e.backend.GetEVMKeeper().GetParams(ctx)
	return (*hexutil.Big)(params.ChainConfig.ChainID.BigInt()), nil
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronize from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (e *PublicAPI) Syncing() (interface{}, error) {
	e.logger.Debug("eth_syncing")

	if !e.backend.GetConsensusReactor().WaitSync() {
		return false, nil
	}

	status, err := tmrpccore.Status(nil)
	if err != nil {
		return false, err
	}

	if !status.SyncInfo.CatchingUp {
		return false, nil
	}

	return map[string]interface{}{
		"startingBlock": hexutil.Uint64(status.SyncInfo.EarliestBlockHeight),
		"currentBlock":  hexutil.Uint64(status.SyncInfo.LatestBlockHeight),
		"highestBlock":  nil, // NA
		"pulledStates":  nil, // NA
		"knownStates":   nil, // NA
	}, nil
}

// Coinbase is the address that staking rewards will be send to (alias for Etherbase).
func (e *PublicAPI) Coinbase() (string, error) {
	e.logger.Debug("eth_coinbase")

	coinbase, err := e.backend.GetCoinbase()
	if err != nil {
		return "", err
	}
	ethAddr := common.BytesToAddress(coinbase.Bytes())
	return ethAddr.Hex(), nil
}

// Mining returns whether or not this node is currently mining. Always false.
func (e *PublicAPI) Mining() bool {
	e.logger.Debug("eth_mining")
	return false
}

// Hashrate returns the current node's hashrate. Always 0.
func (e *PublicAPI) Hashrate() hexutil.Uint64 {
	e.logger.Debug("eth_hashrate")
	return 0
}

// GasPrice returns the current gas price based on stratos's gas price oracle.
func (e *PublicAPI) GasPrice() (*hexutil.Big, error) {
	e.logger.Debug("eth_gasPrice")
	var (
		result *big.Int
		err    error
	)
	baseFee, err := e.backend.BaseFee()
	if err != nil {
		return nil, err
	}
	if baseFee != nil {
		result, err = e.backend.SuggestGasTipCap()
		if err != nil {
			return nil, err
		}
		result = result.Add(result, baseFee)
	} else {
		result = big.NewInt(e.backend.RPCMinGasPrice())
	}

	return (*hexutil.Big)(result), nil
}

// MaxPriorityFeePerGas returns a suggestion for a gas tip cap for dynamic fee transactions.
func (e *PublicAPI) MaxPriorityFeePerGas() (*hexutil.Big, error) {
	e.logger.Debug("eth_maxPriorityFeePerGas")
	tipcap, err := e.backend.SuggestGasTipCap()
	if err != nil {
		return nil, err
	}
	return (*hexutil.Big)(tipcap), nil
}

func (e *PublicAPI) FeeHistory(blockCount rpc.DecimalOrHex, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*rpctypes.FeeHistoryResult, error) {
	e.logger.Debug("eth_feeHistory")
	return e.backend.FeeHistory(blockCount, lastBlock, rewardPercentiles)
}

// Accounts returns the list of accounts available to this node.
func (e *PublicAPI) Accounts() ([]common.Address, error) {
	e.logger.Debug("eth_accounts")

	addresses := make([]common.Address, 0) // return [] instead of nil if empty

	infos, err := e.clientCtx.Keyring.List()
	if err != nil {
		return addresses, err
	}

	for _, info := range infos {
		addressBytes := info.GetPubKey().Address().Bytes()
		addresses = append(addresses, common.BytesToAddress(addressBytes))
	}

	return addresses, nil
}

// BlockNumber returns the current block number.
func (e *PublicAPI) BlockNumber() (hexutil.Uint64, error) {
	e.logger.Debug("eth_blockNumber")
	return e.backend.BlockNumber()
}

// GetBalance returns the provided account's balance up to the provided block number in wei.
func (e *PublicAPI) GetBalance(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (*hexutil.Big, error) {
	e.logger.Debug("eth_getBalance", "address", address.String(), "block number or hash", blockNrOrHash)

	blockNum, err := e.getBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	resBlock, err := e.backend.GetTendermintBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	// return if requested block height is greater than the current one or chain not synced
	if resBlock == nil || resBlock.Block == nil {
		return nil, nil
	}

	sdkCtx, err := e.backend.GetEVMContext().GetSdkContextWithHeader(&resBlock.Block.Header)
	if err != nil {
		return nil, err
	}
	balance := e.backend.GetEVMKeeper().GetBalance(sdkCtx, address)

	return (*hexutil.Big)(balance), nil
}

// GetStorageAt returns the contract storage at the given address, block number, and key.
func (e *PublicAPI) GetStorageAt(address common.Address, key string, blockNrOrHash rpctypes.BlockNumberOrHash) (hexutil.Bytes, error) {
	e.logger.Debug("eth_getStorageAt", "address", address.Hex(), "key", key, "block number or hash", blockNrOrHash)

	blockNum, err := e.getBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	resBlock, err := e.backend.GetTendermintBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	// return if requested block height is greater than the current one or chain not synced
	if resBlock == nil || resBlock.Block == nil {
		return nil, nil
	}

	sdkCtx, err := e.backend.GetEVMContext().GetSdkContextWithHeader(&resBlock.Block.Header)
	if err != nil {
		return nil, err
	}
	state := e.backend.GetEVMKeeper().GetState(sdkCtx, address, common.HexToHash(key))
	return state.Bytes(), nil
}

// GetTransactionCount returns the number of transactions at the given address up to the given block number.
func (e *PublicAPI) GetTransactionCount(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (hexutil.Uint64, error) {
	e.logger.Debug("eth_getTransactionCount", "address", address.Hex(), "block number or hash", blockNrOrHash)
	blockNum, err := e.getBlockNumber(blockNrOrHash)
	if err != nil {
		return 0, err
	}
	return e.backend.GetTransactionCount(address, blockNum)
}

// GetBlockTransactionCountByHash returns the number of transactions in the block identified by hash.
func (e *PublicAPI) GetBlockTransactionCountByHash(hash common.Hash) *hexutil.Uint {
	e.logger.Debug("eth_getBlockTransactionCountByHash", "hash", hash.Hex())

	block, err := tmrpccore.BlockByHash(nil, hash.Bytes())
	if err != nil {
		e.logger.Debug("block not found", "hash", hash.Hex(), "error", err.Error())
		return nil
	}

	if block.Block == nil {
		e.logger.Debug("block not found", "hash", hash.Hex())
		return nil
	}

	n := hexutil.Uint(len(block.Block.Txs))
	return &n
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block identified by number.
func (e *PublicAPI) GetBlockTransactionCountByNumber(blockNum rpctypes.BlockNumber) *hexutil.Uint {
	e.logger.Debug("eth_getBlockTransactionCountByNumber", "height", blockNum.Int64())
	block, err := tmrpccore.Block(nil, blockNum.TmHeight())
	if err != nil {
		e.logger.Debug("block not found", "height", blockNum.Int64(), "error", err.Error())
		return nil
	}

	if block.Block == nil {
		e.logger.Debug("block not found", "height", blockNum.Int64())
		return nil
	}

	n := hexutil.Uint(len(block.Block.Txs))
	return &n
}

// GetUncleCountByBlockHash returns the number of uncles in the block identified by hash. Always zero.
func (e *PublicAPI) GetUncleCountByBlockHash(hash common.Hash) hexutil.Uint {
	return 0
}

// GetUncleCountByBlockNumber returns the number of uncles in the block identified by number. Always zero.
func (e *PublicAPI) GetUncleCountByBlockNumber(blockNum rpctypes.BlockNumber) hexutil.Uint {
	return 0
}

// GetCode returns the contract code at the given address and block number.
func (e *PublicAPI) GetCode(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (hexutil.Bytes, error) {
	e.logger.Debug("eth_getCode", "address", address.Hex(), "block number or hash", blockNrOrHash)

	blockNum, err := e.getBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	req := &evmtypes.QueryCodeRequest{
		Address: address.String(),
	}

	resBlock, err := e.backend.GetTendermintBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	// return if requested block height is greater than the current one or chain not synced
	if resBlock == nil || resBlock.Block == nil {
		return nil, nil
	}

	sdkCtx, err := e.backend.GetEVMContext().GetSdkContextWithHeader(&resBlock.Block.Header)
	if err != nil {
		return nil, err
	}
	res, err := e.backend.GetEVMKeeper().Code(sdk.WrapSDKContext(sdkCtx), req)
	if err != nil {
		return nil, err
	}

	return res.Code, nil
}

// GetTransactionLogs returns the logs given a transaction hash.
func (e *PublicAPI) GetTransactionLogs(txHash common.Hash) ([]*ethtypes.Log, error) {
	e.logger.Debug("eth_getTransactionLogs", "hash", txHash)

	hexTx := txHash.Hex()
	res, err := pool.GetTmTxByHash(txHash)
	if err != nil {
		e.logger.Debug("tx not found", "hash", hexTx, "error", err.Error())
		return nil, nil
	}

	msgIndex, _ := rpctypes.FindTxAttributes(res.TxResult.Events, hexTx)
	if msgIndex < 0 {
		return nil, fmt.Errorf("ethereum tx not found in msgs: %s", hexTx)
	}
	// parse tx logs from events
	return backend.TxLogsFromEvents(res.TxResult.Events, msgIndex)
}

// Sign signs the provided data using the private key of address via Geth's signature standard.
func (e *PublicAPI) Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	e.logger.Debug("eth_sign", "address", address.Hex(), "data", common.Bytes2Hex(data))

	from := sdk.AccAddress(address.Bytes())

	_, err := e.clientCtx.Keyring.KeyByAddress(from)
	if err != nil {
		e.logger.Error("failed to find key in keyring", "address", address.String())
		return nil, fmt.Errorf("%s; %s", keystore.ErrNoMatch, err.Error())
	}

	// Sign the requested hash with the wallet
	signature, _, err := e.clientCtx.Keyring.SignByAddress(from, data)
	if err != nil {
		e.logger.Error("keyring.SignByAddress failed", "address", address.Hex())
		return nil, err
	}

	signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return signature, nil
}

// SignTypedData signs EIP-712 conformant typed data
func (e *PublicAPI) SignTypedData(address common.Address, typedData apitypes.TypedData) (hexutil.Bytes, error) {
	e.logger.Debug("eth_signTypedData", "address", address.Hex(), "data", typedData)
	from := sdk.AccAddress(address.Bytes())

	_, err := e.clientCtx.Keyring.KeyByAddress(from)
	if err != nil {
		e.logger.Error("failed to find key in keyring", "address", address.String())
		return nil, fmt.Errorf("%s; %s", keystore.ErrNoMatch, err.Error())
	}

	sigHash, err := eip712.ComputeTypedDataHash(typedData)
	if err != nil {
		return nil, err
	}

	// Sign the requested hash with the wallet
	signature, _, err := e.clientCtx.Keyring.SignByAddress(from, sigHash)
	if err != nil {
		e.logger.Error("keyring.SignByAddress failed", "address", address.Hex())
		return nil, err
	}

	signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return signature, nil
}

// SendTransaction sends an Ethereum transaction.
func (e *PublicAPI) SendTransaction(args evmtypes.TransactionArgs) (common.Hash, error) {
	e.logger.Debug("eth_sendTransaction", "args", args.String())
	return e.backend.SendTransaction(args)
}

// FillTransaction fills the defaults (nonce, gas, gasPrice or 1559 fields)
// on a given unsigned transaction, and returns it to the caller for further
// processing (signing + broadcast).
func (e *PublicAPI) FillTransaction(args evmtypes.TransactionArgs) (*rpctypes.SignTransactionResult, error) {
	// Set some sanity defaults and terminate on failure
	args, err := e.backend.SetTxDefaults(args)
	if err != nil {
		return nil, err
	}

	// Assemble the transaction and obtain rlp
	tx := args.ToTransaction().AsTransaction()

	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return &rpctypes.SignTransactionResult{
		Raw: data,
		Tx:  tx,
	}, nil
}

// SendRawTransaction send a raw Ethereum transaction.
func (e *PublicAPI) SendRawTransaction(data hexutil.Bytes) (common.Hash, error) {
	e.logger.Debug("eth_sendRawTransaction", "length", len(data))

	// RLP decode raw transaction bytes
	tx := &ethtypes.Transaction{}
	if err := tx.UnmarshalBinary(data); err != nil {
		e.logger.Error("transaction decoding failed", "error", err.Error())
		return common.Hash{}, err
	}

	_, err := e.backend.GetTxPool().Add(tx)
	if err != nil {
		e.logger.Error("failed to add eth tx into tx evm pool", "error", err.Error())
		return common.Hash{}, err
	}

	return tx.Hash(), nil
}

// checkTxFee is an internal function used to check whether the fee of
// the given transaction is _reasonable_(under the cap).
func checkTxFee(gasPrice *big.Int, gas uint64, cap float64) error {
	// Short circuit if there is no cap for transaction fee at all.
	if cap == 0 {
		return nil
	}
	totalfee := new(big.Float).SetInt(new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gas)))
	// 1 stos in 10^18 wei
	oneToken := new(big.Float).SetInt(big.NewInt(stratos.StosToWei))
	// quo = rounded(x/y)
	feeEth := new(big.Float).Quo(totalfee, oneToken)
	// no need to check error from parsing
	feeFloat, _ := feeEth.Float64()
	if feeFloat > cap {
		return fmt.Errorf("tx fee (%.2f ether) exceeds the configured cap (%.2f ether)", feeFloat, cap)
	}
	return nil
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func (e *PublicAPI) Resend(ctx context.Context, args evmtypes.TransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	e.logger.Debug("eth_resend", "args", args.String())
	if args.Nonce == nil {
		return common.Hash{}, fmt.Errorf("missing transaction nonce in transaction spec")
	}

	args, err := e.backend.SetTxDefaults(args)
	if err != nil {
		return common.Hash{}, err
	}

	matchTx := args.ToTransaction().AsTransaction()

	// Before replacing the old transaction, ensure the _new_ transaction fee is reasonable.
	price := matchTx.GasPrice()
	if gasPrice != nil {
		price = gasPrice.ToInt()
	}
	gas := matchTx.Gas()
	if gasLimit != nil {
		gas = uint64(*gasLimit)
	}
	if err := checkTxFee(price, gas, e.backend.RPCTxFeeCap()); err != nil {
		return common.Hash{}, err
	}

	pending, err := e.backend.PendingTransactions()
	if err != nil {
		return common.Hash{}, err
	}

	for _, tx := range pending {
		p, err := evmtypes.UnwrapEthereumMsg(tx, common.Hash{})
		if err != nil {
			// not valid ethereum tx
			continue
		}

		pTx := p.AsTransaction()

		wantSigHash := e.signer.Hash(matchTx)
		pFrom, err := ethtypes.Sender(e.signer, pTx)
		if err != nil {
			continue
		}

		if pFrom == *args.From && e.signer.Hash(pTx) == wantSigHash {
			// Match. Re-sign and send the transaction.
			if gasPrice != nil && (*big.Int)(gasPrice).Sign() != 0 {
				args.GasPrice = gasPrice
			}
			if gasLimit != nil && *gasLimit != 0 {
				args.Gas = gasLimit
			}

			return e.backend.SendTransaction(args) // TODO: this calls SetTxDefaults again, refactor to avoid calling it twice
		}
	}

	return common.Hash{}, fmt.Errorf("transaction %#x not found", matchTx.Hash())
}

// Call performs a raw contract call.
func (e *PublicAPI) Call(args evmtypes.TransactionArgs, blockNrOrHash rpctypes.BlockNumberOrHash, _ *rpctypes.StateOverride) (hexutil.Bytes, error) {
	e.logger.Debug("eth_call", "args", args.String(), "block number or hash", blockNrOrHash)

	blockNum, err := e.getBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}
	data, err := e.doCall(args, blockNum)
	if err != nil {
		return []byte{}, err
	}

	return (hexutil.Bytes)(data.Ret), nil
}

// DoCall performs a simulated call operation through the evmtypes. It returns the
// estimated gas used on the operation or an error if fails.
func (e *PublicAPI) doCall(
	args evmtypes.TransactionArgs, blockNr rpctypes.BlockNumber,
) (*evmtypes.MsgEthereumTxResponse, error) {
	bz, err := json.Marshal(&args)
	if err != nil {
		return nil, err
	}

	req := evmtypes.EthCallRequest{
		Args:   bz,
		GasCap: e.backend.RPCGasCap(),
	}

	resBlock, err := e.backend.GetTendermintBlockByNumber(blockNr)
	if err != nil {
		return nil, err
	}

	// return if requested block height is greater than the current one or chain not synced
	if resBlock == nil || resBlock.Block == nil {
		return nil, nil
	}

	sdkCtx, err := e.backend.GetEVMContext().GetSdkContextWithHeader(&resBlock.Block.Header)
	if err != nil {
		return nil, err
	}

	// it will return an empty context and the gRPC query will use
	// the latest block height for querying.
	timeout := e.backend.RPCEVMTimeout()

	ctx := sdk.WrapSDKContext(sdkCtx)

	// Setup context so it may be canceled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	// Make sure the context is canceled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	res, err := e.backend.GetEVMKeeper().EthCall(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		if res.VmError != vm.ErrExecutionReverted.Error() {
			return nil, status.Error(codes.Internal, res.VmError)
		}
		return nil, evmtypes.NewExecErrorWithReason(res.Ret)
	}

	return res, nil
}

// EstimateGas returns an estimate of gas usage for the given smart contract call.
func (e *PublicAPI) EstimateGas(args evmtypes.TransactionArgs, blockNrOptional *rpctypes.BlockNumber) (hexutil.Uint64, error) {
	e.logger.Debug("eth_estimateGas")
	return e.backend.EstimateGas(args, blockNrOptional)
}

// GetBlockByHash returns the block identified by hash.
func (e *PublicAPI) GetBlockByHash(hash common.Hash, fullTx bool) (*types.Block, error) {
	e.logger.Debug("eth_getBlockByHash", "hash", hash.Hex(), "full", fullTx)
	return e.backend.GetBlockByHash(hash, fullTx)
}

// GetBlockByNumber returns the block identified by number.
func (e *PublicAPI) GetBlockByNumber(ethBlockNum rpctypes.BlockNumber, fullTx bool) (*types.Block, error) {
	e.logger.Debug("eth_getBlockByNumber", "number", ethBlockNum, "full", fullTx)
	return e.backend.GetBlockByNumber(ethBlockNum, fullTx)
}

// GetTransactionByHash returns the transaction identified by hash.
func (e *PublicAPI) GetTransactionByHash(hash common.Hash) (*rpctypes.Transaction, error) {
	e.logger.Debug("eth_getTransactionByHash", "hash", hash.Hex())
	return e.backend.GetTransactionByHash(hash)
}

// getTransactionByBlockAndIndex is the common code shared by `GetTransactionByBlockNumberAndIndex` and `GetTransactionByBlockHashAndIndex`.
func (e *PublicAPI) getTransactionByBlockAndIndex(block *tmrpctypes.ResultBlock, idx hexutil.Uint) (*rpctypes.Transaction, error) {
	// return if index out of bounds
	if uint64(idx) >= uint64(len(block.Block.Txs)) {
		return nil, nil
	}

	tx := block.Block.Txs[idx]

	blockHash := common.BytesToHash(block.Block.Hash())
	blockHeight := uint64(block.Block.Height)
	txIndex := uint64(idx)

	return rpctypes.TmTxToEthTx(
		e.clientCtx.TxConfig.TxDecoder(),
		tx,
		&blockHash,
		&blockHeight,
		&txIndex,
	)
}

// GetTransactionByBlockHashAndIndex returns the transaction identified by hash and index.
func (e *PublicAPI) GetTransactionByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) (*rpctypes.Transaction, error) {
	e.logger.Debug("eth_getTransactionByBlockHashAndIndex", "hash", hash.Hex(), "index", idx)

	block, err := tmrpccore.BlockByHash(nil, hash.Bytes())
	if err != nil {
		e.logger.Debug("block not found", "hash", hash.Hex(), "error", err.Error())
		return nil, nil
	}

	if block.Block == nil {
		e.logger.Debug("block not found", "hash", hash.Hex())
		return nil, nil
	}

	return e.getTransactionByBlockAndIndex(block, idx)
}

// GetTransactionByBlockNumberAndIndex returns the transaction identified by number and index.
func (e *PublicAPI) GetTransactionByBlockNumberAndIndex(blockNum rpctypes.BlockNumber, idx hexutil.Uint) (*rpctypes.Transaction, error) {
	e.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "number", blockNum, "index", idx)

	block, err := tmrpccore.Block(nil, blockNum.TmHeight())
	if err != nil {
		e.logger.Debug("block not found", "height", blockNum.Int64(), "error", err.Error())
		return nil, nil
	}

	if block.Block == nil {
		e.logger.Debug("block not found", "height", blockNum.Int64())
		return nil, nil
	}

	return e.getTransactionByBlockAndIndex(block, idx)
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (e *PublicAPI) GetTransactionReceipt(hash common.Hash) (*rpctypes.TransactionReceipt, error) {
	e.logger.Debug("eth_getTransactionReceipt", "hash", hash)
	res, err := pool.GetTmTxByHash(hash)
	if err != nil {
		return nil, nil
	}

	block := e.backend.GetBlockStore().LoadBlock(res.Height)
	if block == nil {
		e.logger.Debug("eth_getTransactionReceipt", "hash", hash, "block not found")
		return nil, nil
	}

	blockResults, err := tmrpccore.BlockResults(nil, &block.Height)
	if err != nil {
		e.logger.Debug("eth_getTransactionReceipt", "hash", hash, "block not found")
		return nil, nil
	}

	blockHash := common.BytesToHash(block.Hash())
	blockHeight := uint64(res.Height)
	txIndex := uint64(res.Index)

	rpcTx, err := rpctypes.TmTxToEthTx(
		e.clientCtx.TxConfig.TxDecoder(),
		res.Tx,
		&blockHash,
		&blockHeight,
		&txIndex,
	)
	if err != nil {
		return nil, err
	}

	cumulativeGasUsed := uint64(res.TxResult.GasUsed)
	if *rpcTx.TransactionIndex != 0 {
		cumulativeGasUsed += rpctypes.GetBlockCumulativeGas(blockResults, int(*rpcTx.TransactionIndex))
	}

	_, attrs := rpctypes.FindTxAttributes(res.TxResult.Events, hash.Hex())

	var (
		contractAddress *common.Address
		bloom           = ethtypes.BytesToBloom(make([]byte, 6))
		logs            = make([]*ethtypes.Log, 0)
	)
	// Set status codes based on tx result
	status := ethtypes.ReceiptStatusSuccessful
	if res.TxResult.GetCode() == 1 {
		status = ethtypes.ReceiptStatusFailed
	} else {
		// Get the transaction result from the log
		_, found := attrs[evmtypes.AttributeKeyEthereumTxFailed]
		if found {
			status = ethtypes.ReceiptStatusFailed
		}

		if status == ethtypes.ReceiptStatusSuccessful {
			// parse tx logs from events
			logs, err = backend.TxLogsFromEvents(res.TxResult.Events, 0)
			if err != nil {
				e.logger.Debug("logs not found", "hash", hash, "error", err.Error())
			}
			if logs == nil {
				logs = make([]*ethtypes.Log, 0)
			}
			if rpcTx.To == nil {
				// TODO: Rewrite on more optimal way in order to get a contract address
				tx, err := e.clientCtx.TxConfig.TxDecoder()(res.Tx)
				if err != nil {
					e.logger.Debug("decoding failed", "error", err.Error())
					return nil, fmt.Errorf("failed to decode tx: %w", err)
				}

				// the `msgIndex` is inferred from tx events, should be within the bound.
				msg := tx.GetMsgs()[0]
				ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
				if !ok {
					e.logger.Debug(fmt.Sprintf("invalid tx type: %T", msg))
					return nil, fmt.Errorf("invalid tx type: %T", msg)
				}

				txData, err := evmtypes.UnpackTxData(ethMsg.Data)
				if err != nil {
					e.logger.Error("failed to unpack tx data", "error", err.Error())
					return nil, err
				}
				contractAddress = new(common.Address)
				*contractAddress = crypto.CreateAddress(rpcTx.From, txData.GetNonce())
			}
			bloom = ethtypes.BytesToBloom(ethtypes.LogsBloom(logs))
		}
	}

	receipt := &rpctypes.TransactionReceipt{
		Status:            hexutil.Uint64(status),
		CumulativeGasUsed: hexutil.Uint64(cumulativeGasUsed),
		LogsBloom:         bloom,
		Logs:              logs,
		TransactionHash:   rpcTx.Hash,
		ContractAddress:   contractAddress,
		GasUsed:           hexutil.Uint64(res.TxResult.GasUsed),
		BlockHash:         *rpcTx.BlockHash,
		BlockNumber:       *rpcTx.BlockNumber,
		TransactionIndex:  *rpcTx.TransactionIndex,
		From:              rpcTx.From,
		To:                rpcTx.To,
	}

	return receipt, nil
}

// GetPendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (e *PublicAPI) GetPendingTransactions() ([]*rpctypes.Transaction, error) {
	e.logger.Debug("eth_getPendingTransactions")

	txs := e.backend.GetMempool().ReapMaxTxs(100)

	result := make([]*rpctypes.Transaction, 0, len(txs))
	for _, tx := range txs {
		rpctx, err := rpctypes.TmTxToEthTx(e.clientCtx.TxConfig.TxDecoder(), tx, nil, nil, nil)
		if err != nil {
			return nil, err
		}

		result = append(result, rpctx)
	}

	return result, nil
}

// GetUncleByBlockHashAndIndex returns the uncle identified by hash and index. Always returns nil.
func (e *PublicAPI) GetUncleByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) map[string]interface{} {
	return nil
}

// GetUncleByBlockNumberAndIndex returns the uncle identified by number and index. Always returns nil.
func (e *PublicAPI) GetUncleByBlockNumberAndIndex(number, idx hexutil.Uint) map[string]interface{} {
	return nil
}

// GetProof returns an account object with proof and any storage proofs
func (e *PublicAPI) GetProof(address common.Address, storageKeys []string, blockNrOrHash rpctypes.BlockNumberOrHash) (*rpctypes.AccountResult, error) {
	e.logger.Debug("eth_getProof", "address", address.Hex(), "keys", storageKeys, "block number or hash", blockNrOrHash)

	blockNum, err := e.getBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	height := blockNum.Int64()

	// if the height is equal to zero, meaning the query condition of the block is either "pending" or "latest"
	if height == 0 {
		bn, err := e.backend.BlockNumber()
		if err != nil {
			return nil, err
		}

		if bn > math.MaxInt64 {
			return nil, fmt.Errorf("not able to query block number greater than MaxInt64")
		}

		height = int64(bn)
	}

	// query storage proofs
	storageProofs := make([]rpctypes.StorageResult, len(storageKeys))

	resBlock, err := e.backend.GetTendermintBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	// return if requested block height is greater than the current one or chain not synced
	if resBlock == nil || resBlock.Block == nil {
		return nil, nil
	}

	sdkCtx, err := e.backend.GetEVMContext().GetSdkContextWithHeader(&resBlock.Block.Header)
	if err != nil {
		return nil, err
	}

	for i, key := range storageKeys {
		hexKey := common.HexToHash(key)
		valueBz, proof, err := types.GetProof(height, evmtypes.StoreKey, evmtypes.StateKey(address, hexKey.Bytes()))
		if err != nil {
			return nil, err
		}

		// check for proof
		var proofStr string
		if proof != nil {
			proofStr = proof.String()
		}

		storageProofs[i] = rpctypes.StorageResult{
			Key:   key,
			Value: (*hexutil.Big)(new(big.Int).SetBytes(valueBz)),
			Proof: []string{proofStr},
		}
	}

	// query EVM account
	acc := e.backend.GetEVMKeeper().GetAccountOrEmpty(sdkCtx, address)

	// query account proofs
	accountKey := authtypes.AddressStoreKey(sdk.AccAddress(address.Bytes()))
	_, proof, err := types.GetProof(height, authtypes.StoreKey, accountKey)
	if err != nil {
		return nil, err
	}

	// check for proof
	var accProofStr string
	if proof != nil {
		accProofStr = proof.String()
	}

	return &rpctypes.AccountResult{
		Address:      address,
		AccountProof: []string{accProofStr},
		Balance:      (*hexutil.Big)(acc.Balance),
		CodeHash:     common.BytesToHash(acc.CodeHash),
		Nonce:        hexutil.Uint64(acc.Nonce),
		StorageHash:  common.Hash{}, // NOTE: stratos doesn't have a storage hash. TODO: implement?
		StorageProof: storageProofs,
	}, nil
}

// getBlockNumber returns the BlockNumber from BlockNumberOrHash
func (e *PublicAPI) getBlockNumber(blockNrOrHash rpctypes.BlockNumberOrHash) (rpctypes.BlockNumber, error) {
	switch {
	case blockNrOrHash.BlockHash == nil && blockNrOrHash.BlockNumber == nil:
		return rpctypes.EthEarliestBlockNumber, fmt.Errorf("types BlockHash and BlockNumber cannot be both nil")
	case blockNrOrHash.BlockHash != nil:
		blockHeader, err := e.backend.HeaderByHash(*blockNrOrHash.BlockHash)
		if err != nil {
			return rpctypes.EthEarliestBlockNumber, err
		}
		return rpctypes.NewBlockNumber(blockHeader.Number), nil
	case blockNrOrHash.BlockNumber != nil:
		return *blockNrOrHash.BlockNumber, nil
	default:
		return rpctypes.EthEarliestBlockNumber, nil
	}
}
