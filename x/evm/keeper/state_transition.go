package keeper

import (
	"math"
	"math/big"
	"sync"

	tmtypes "github.com/cometbft/cometbft/types"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/evm/statedb"
	"github.com/stratosnet/stratos-chain/x/evm/tracers"
	"github.com/stratosnet/stratos-chain/x/evm/types"
	"github.com/stratosnet/stratos-chain/x/evm/vm"
)

// GasToRefund calculates the amount of gas the state machine should refund to the sender. It is
// capped by the refund quotient value.
// Note: do not pass 0 to refundQuotient
func GasToRefund(availableRefund, gasConsumed, refundQuotient uint64) uint64 {
	// Apply refund counter
	refund := gasConsumed / refundQuotient
	if refund > availableRefund {
		return availableRefund
	}
	return refund
}

// EVMConfig creates the EVMConfig based on current state
func (k *Keeper) EVMConfig(ctx sdk.Context) (*types.EVMConfig, error) {
	evmParams := k.GetParams(ctx)
	ethCfg := evmParams.ChainConfig.EthereumConfig()

	coinbase, err := k.GetCoinbaseAddress(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain coinbase address")
	}

	baseFee := k.GetBaseFee(ctx, ethCfg)
	return &types.EVMConfig{
		Params:      evmParams,
		ChainConfig: ethCfg,
		CoinBase:    coinbase,
		BaseFee:     baseFee,
	}, nil
}

// TxConfig loads `TxConfig` from current transient storage
func (k *Keeper) TxConfig(ctx sdk.Context, txHash common.Hash) statedb.TxConfig {
	return statedb.NewTxConfig(
		common.BytesToHash(ctx.HeaderHash()), // BlockHash
		txHash,                               // TxHash
		uint(k.GetTxIndexTransient(ctx)),     // TxIndex
		uint(k.GetLogSizeTransient(ctx)),     // LogIndex
	)
}

// NewEVM generates a go-ethereum VM from the provided Message fields and the chain parameters
// (ChainConfig and module Params). It additionally sets the validator operator address as the
// coinbase address to make it available for the COINBASE opcode, even though there is no
// beneficiary of the coinbase transaction (since we're not mining).
func (k *Keeper) NewEVM(
	ctx sdk.Context,
	msg core.Message,
	cfg *types.EVMConfig,
	tracer vm.EVMLogger,
	stateDB vm.StateDB,
) *vm.EVM {
	blockCtx := vm.BlockContext{
		CanTransfer:        vm.CanTransfer,
		Transfer:           vm.Transfer,
		GetHash:            k.GetHashFn(ctx),
		RunSdkMsg:          k.RunSdkMsgFn(ctx, false),
		ParseProtoFromData: k.ParseProtoFromDataFn(ctx),
		Coinbase:           cfg.CoinBase,
		GasLimit:           stratos.BlockGasLimit(ctx),
		BlockNumber:        big.NewInt(ctx.BlockHeight()),
		Time:               big.NewInt(ctx.BlockHeader().Time.Unix()),
		Difficulty:         big.NewInt(0), // unused. Only required in PoW context
		BaseFee:            cfg.BaseFee,
	}

	txCtx := vm.NewEVMTxContext(msg)
	if tracer == nil {
		tracer = k.Tracer(ctx, msg, cfg.ChainConfig)
	}
	vmConfig := k.VMConfig(ctx, cfg, tracer)
	return vm.NewEVM(blockCtx, txCtx, stateDB, cfg.ChainConfig, vmConfig, k.verifier)
}

// VMConfig creates an EVM configuration from the debug setting and the extra EIPs enabled on the
// module parameters. The config generated uses the default JumpTable from the EVM.
func (k Keeper) VMConfig(ctx sdk.Context, cfg *types.EVMConfig, tracer vm.EVMLogger) vm.Config {
	noBaseFee := true
	if types.IsLondon(cfg.ChainConfig, ctx.BlockHeight()) {
		noBaseFee = k.GetParams(ctx).FeeMarketParams.NoBaseFee
	}

	var debug bool

	if _, ok := tracer.(tracers.NoOpTracer); !ok {
		debug = true
		noBaseFee = true
	}

	return vm.Config{
		Debug:     debug,
		Tracer:    tracer,
		NoBaseFee: noBaseFee,
		ExtraEips: cfg.Params.EIPs(),
	}
}

// GetHashFn implements vm.GetHashFunc for stratos. It handles 3 cases:
//  1. The requested height matches the current height from context (and thus same epoch number)
//  2. The requested height is from a previous height from the same chain epoch
//  3. The requested height is from a height greater than the latest one
func (k Keeper) GetHashFn(ctx sdk.Context) vm.GetHashFunc {
	cache := make(map[int64]common.Hash)
	var rw sync.Mutex

	return func(height uint64) common.Hash {
		h, err := stratos.SafeInt64(height)
		if err != nil {
			k.Logger(ctx).Error("failed to cast height to int64", "error", err)
			return common.Hash{}
		}

		switch {
		case ctx.BlockHeight() == h:
			// Case 1: The requested height matches the one from the context, so we can retrieve the header
			// hash directly from the context.
			// Note: The headerHash is only set at begin block, it will be nil in case of a query context
			headerHash := ctx.HeaderHash()
			return common.BytesToHash(headerHash)

		case ctx.BlockHeight() > h:
			// Case 2: if the chain is not the current height we need to retrieve the hash from the store for the
			// current chain epoch. This only applies if the current height is greater than the requested height.

			// NOTE: In case of concurrency
			rw.Lock()
			defer rw.Unlock()

			if hash, ok := cache[h]; ok {
				return hash
			}

			histInfo, found := k.stakingKeeper.GetHistoricalInfo(ctx, h)
			if !found {
				k.Logger(ctx).Debug("historical info not found", "height", h)
				return common.Hash{}
			}

			header, err := tmtypes.HeaderFromProto(&histInfo.Header)
			if err != nil {
				k.Logger(ctx).Error("failed to cast tendermint header from proto", "error", err)
				return common.Hash{}
			}

			hash := common.BytesToHash(header.Hash())
			cache[h] = hash
			return hash
		default:
			// Case 3: heights greater than the current one returns an empty hash.
			return common.Hash{}
		}
	}
}

// GetSdkMsg return prepared msg from data
func (k Keeper) GetSdkMsg(from sdk.AccAddress, data []byte) (*types.MsgCosmosData, error) {
	msg := types.NewMsgCosmosData(k.cdc, from)
	if err := msg.Parse(data); err != nil {
		return nil, err
	}
	return msg, nil
}

// ParseProtoFromDataFn parse data and get proto info
func (k *Keeper) ParseProtoFromDataFn(ctx sdk.Context) vm.ParseProtoFromDataFunc {
	return func(data []byte, gas uint64) ([]byte, uint64, error) {
		gasBefore := ctx.GasMeter().GasConsumed()
		getGas := func() uint64 {
			return gas - types.Max(gasBefore-ctx.GasMeter().GasConsumed(), 2500)
		}
		any, err := types.TxDataToAny(data)
		if err != nil {
			return []byte{}, getGas(), err
		}

		cMsg, ok := any.GetCachedValue().(sdk.Msg)
		if !ok {
			return []byte{}, getGas(), err
		}

		if len(cMsg.GetSigners()) == 0 {
			return []byte{}, getGas(), errors.Wrapf(sdkerrors.ErrorInvalidSigner, "signer not found")
		}

		signer := cMsg.GetSigners()[0]

		addressTy, _ := abi.NewType("address", "", nil)
		bytesTy, _ := abi.NewType("bytes", "", nil)

		arguments := abi.Arguments{
			{
				Type: addressTy,
			},
			{
				Type: bytesTy,
			},
		}
		res, err := arguments.Pack(
			common.BytesToAddress(signer.Bytes()),
			[]byte(any.TypeUrl),
		)
		if err != nil {
			return []byte{}, getGas(), err
		}

		return res, getGas(), nil
	}
}

// RunSdkMsg execute cosmos msg from payload (NOTE: simulate=true always in smart contracts!!!)
func (k *Keeper) RunSdkMsgFn(ctx sdk.Context, simulate bool) vm.RunSdkMsgFunc {
	return func(from common.Address, data []byte, gas uint64) ([]byte, uint64, error) {
		if k.msgServiceRouter == nil {
			return nil, 0, errors.Wrapf(sdkerrors.ErrUnknownRequest, "service router not set")
		}

		cMsg, err := k.GetSdkMsg(from.Bytes(), data)
		if err != nil {
			return nil, 0, errors.Wrap(err, "failed to get cosmos msg")
		}

		// NOTE: in simulation we should check this in order avoid tx sends from estimations
		if simulate {
			if err := cMsg.ValidateBasic(); err != nil {
				return nil, 0, err
			}
		}

		var ret []byte
		gasBefore := ctx.GasMeter().GasConsumed()
		{
			msg := cMsg.GetMsgs()[0]
			handler := k.msgServiceRouter.Handler(msg)
			if handler == nil {
				return nil, 0, errors.Wrapf(sdkerrors.ErrUnknownRequest, "can't route message %+v", msg)
			}
			// ADR 031 request type routing
			msgResult, err := handler(ctx, msg)
			// NOTE: Error should be returned because we do not know at which state error occured
			if err != nil {
				return nil, 0, errors.Wrapf(err, "failed to execute message")
			}
			// TODO: Maybe pack them also as ethereum events?
			ctx.EventManager().EmitEvents(msgResult.GetEvents())

			if len(msgResult.MsgResponses) > 0 {
				msgResponse := msgResult.MsgResponses[0]
				if msgResponse == nil {
					return nil, 0, sdkerrors.ErrLogic.Wrapf("got nil Msg response for msg %s", sdk.MsgTypeURL(msg))
				}
				ret = msgResponse.GetValue()
			}
		}
		gasUsed := types.Max(ctx.GasMeter().GasConsumed()-gasBefore, params.TxGas)
		if gas < gasUsed {
			return nil, 0, errors.Wrap(types.ErrGasOverflow, "apply message")
		}
		return ret, gasUsed, nil
	}
}

// ApplyTransaction runs and attempts to perform a state transition with the given transaction (i.e. Message), that will
// only be persisted (committed) to the underlying KVStore if the transaction does not fail.
//
// # Gas tracking
//
// Ethereum consumes gas according to the EVM opcodes instead of general reads and writes to store. Because of this, the
// state transition needs to ignore the SDK gas consumption mechanism defined by the GasKVStore and instead consume the
// amount of gas used by the VM execution. The amount of gas used is tracked by the EVM and returned in the execution
// result.
//
// Prior to the execution, the starting tx gas meter is saved and replaced with an infinite gas meter in a new context
// in order to ignore the SDK gas consumption config values (read, write, has, delete).
// After the execution, the gas used from the message execution will be added to the starting gas consumed, taking into
// consideration the amount of gas returned. Finally, the context is updated with the EVM gas consumed value prior to
// returning.
//
// For relevant discussion see: https://github.com/cosmos/cosmos-sdk/discussions/9072
func (k *Keeper) ApplyTransaction(ctx sdk.Context, tx *ethtypes.Transaction) (*types.MsgEthereumTxResponse, error) {
	var (
		bloom        *big.Int
		bloomReceipt ethtypes.Bloom
	)

	cfg, err := k.EVMConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load evm config")
	}
	txConfig := k.TxConfig(ctx, tx.Hash())

	// get the signer according to the chain rules from the config and block height
	signer := ethtypes.MakeSigner(cfg.ChainConfig, big.NewInt(ctx.BlockHeight()))
	msg, err := tx.AsMessage(signer, cfg.BaseFee)
	if err != nil {
		return nil, errors.Wrap(err, "failed to return ethereum transaction as core message")
	}

	// snapshot to contain the tx processing and post-processing in same scope
	var commit func()
	tmpCtx := ctx
	if k.hooks != nil {
		// Create a cache context to revert state when tx hooks fails,
		// the cache context is only committed when both tx and hooks executed successfully.
		// Didn't use `Snapshot` because the context stack has exponential complexity on certain operations,
		// thus restricted to be used only inside `ApplyMessage`.
		tmpCtx, commit = ctx.CacheContext()
	}

	res, err := k.ApplyAutoMessageWithConfig(tmpCtx, msg, nil, true, cfg, txConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply ethereum core message")
	}

	logs := types.LogsToEthereum(res.Logs)

	// Compute block bloom filter
	if len(logs) > 0 {
		bloom = k.GetBlockBloomTransient(ctx)
		bloom.Or(bloom, big.NewInt(0).SetBytes(ethtypes.LogsBloom(logs)))
		bloomReceipt = ethtypes.BytesToBloom(bloom.Bytes())
	}

	cumulativeGasUsed := res.GasUsed
	if ctx.BlockGasMeter() != nil {
		limit := ctx.BlockGasMeter().Limit()
		consumed := ctx.BlockGasMeter().GasConsumed()
		cumulativeGasUsed = uint64(math.Min(float64(cumulativeGasUsed+consumed), float64(limit)))
	}

	var contractAddr common.Address
	if msg.To() == nil {
		contractAddr = crypto.CreateAddress(msg.From(), msg.Nonce())
	}

	receipt := &ethtypes.Receipt{
		Type:              tx.Type(),
		PostState:         nil, // TODO: intermediate state root
		CumulativeGasUsed: cumulativeGasUsed,
		Bloom:             bloomReceipt,
		Logs:              logs,
		TxHash:            txConfig.TxHash,
		ContractAddress:   contractAddr,
		GasUsed:           res.GasUsed,
		BlockHash:         txConfig.BlockHash,
		BlockNumber:       big.NewInt(ctx.BlockHeight()),
		TransactionIndex:  txConfig.TxIndex,
	}

	if !res.Failed() {
		receipt.Status = ethtypes.ReceiptStatusSuccessful
		// Only call hooks if tx executed successfully.
		if err = k.PostTxProcessing(tmpCtx, msg, receipt); err != nil {
			// If hooks return error, revert the whole tx.
			res.VmError = types.ErrPostTxProcessing.Error()
			k.Logger(ctx).Error("tx post processing failed", "error", err)
		} else if commit != nil {
			// PostTxProcessing is successful, commit the tmpCtx
			commit()
			ctx.EventManager().EmitEvents(tmpCtx.EventManager().Events())
		}
	}

	// refund gas in order to match the Ethereum gas consumption instead of the default SDK one.
	if err = k.RefundGas(ctx, msg, msg.Gas()-res.GasUsed, cfg.Params.EvmDenom); err != nil {
		return nil, errors.Wrapf(err, "failed to refund gas leftover gas to sender %s", msg.From())
	}

	if len(receipt.Logs) > 0 {
		// Update transient block bloom filter
		k.SetBlockBloomTransient(ctx, receipt.Bloom.Big())
		k.SetLogSizeTransient(ctx, uint64(txConfig.LogIndex)+uint64(len(receipt.Logs)))
	}

	k.SetTxIndexTransient(ctx, uint64(txConfig.TxIndex)+1)

	totalGasUsed, err := k.AddTransientGasUsed(ctx, res.GasUsed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add transient gas used")
	}

	// reset the gas meter for current cosmos transaction
	k.ResetGasMeterAndConsumeGas(ctx, totalGasUsed)
	return res, nil
}

// ApplyAutoMessageWithConfig automatic handle where msg should be proceed
func (k *Keeper) ApplyAutoMessageWithConfig(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool, cfg *types.EVMConfig, txConfig statedb.TxConfig) (*types.MsgEthereumTxResponse, error) {
	var (
		res *types.MsgEthereumTxResponse
		err error
	)

	if !types.IsCosmosHandler(msg.To()) {
		// 1) EVM Executioner
		res, err = k.ApplyMessageWithConfig(ctx, msg, tracer, commit, cfg, txConfig)
	} else {
		// 2) Cosmos msg router executioner
		res, err = k.ApplyCosmosMessageWithConfig(ctx, msg, tracer, commit, cfg, txConfig)
	}
	return res, err
}

// ApplyCosmosMessageWithConfig runs and attempts to perform a state transition with the given transaction for cosmos msd (i.e. Message), that will
// only be persisted (committed) to the underlying KVStore if the transaction does not fail.
// It will be executed only if all checks are fine. In case tx aborted, it should not been even commited because it does not use statedb
func (k *Keeper) ApplyCosmosMessageWithConfig(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool, cfg *types.EVMConfig, txConfig statedb.TxConfig) (*types.MsgEthereumTxResponse, error) {
	ret, gasUsed, err := k.RunSdkMsgFn(ctx, !commit)(msg.From(), msg.Data(), msg.Gas())
	if err != nil {
		return nil, err
	}

	response := &types.MsgEthereumTxResponse{
		Hash:    txConfig.TxHash.Hex(),
		Logs:    types.NewLogsFromEth(nil),
		GasUsed: gasUsed,
		Ret:     ret,
	}

	return response, nil
}

// ApplyMessageWithConfig computes the new state by applying the given message against the existing state.
// If the message fails, the VM execution error with the reason will be returned to the client
// and the transaction won't be committed to the store.
//
// # Reverted state
//
// The snapshot and rollback are supported by the `statedb.StateDB`.
//
// # Different Callers
//
// It's called in three scenarios:
// 1. `ApplyTransaction`, in the transaction processing flow.
// 2. `EthCall/EthEstimateGas` grpc query handler.
// 3. Called by other native modules directly.
//
// # Prechecks and Preprocessing
//
// All relevant state transition prechecks for the MsgEthereumTx are performed on the AnteHandler,
// prior to running the transaction against the state. The prechecks run are the following:
//
// 1. the nonce of the message caller is correct
// 2. caller has enough balance to cover transaction fee(gaslimit * gasprice)
// 3. the amount of gas required is available in the block
// 4. the purchased gas is enough to cover intrinsic usage
// 5. there is no overflow when calculating intrinsic gas
// 6. caller has enough balance to cover asset transfer for **topmost** call
//
// The preprocessing steps performed by the AnteHandler are:
//
// 1. set up the initial access list (iff fork > Berlin)
//
// # Tracer parameter
//
// It should be a `vm.Tracer` object or nil, if pass `nil`, it'll create a default one based on keeper options.
//
// # Commit parameter
//
// If commit is true, the `StateDB` will be committed, otherwise discarded.
func (k *Keeper) ApplyMessageWithConfig(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool, cfg *types.EVMConfig, txConfig statedb.TxConfig) (*types.MsgEthereumTxResponse, error) {
	var (
		ret   []byte // return bytes from evm execution
		vmErr error  // vm errors do not effect consensus and are therefore not assigned to err
	)

	// return error if contract creation or call are disabled through governance
	if !cfg.Params.EnableCreate && msg.To() == nil {
		return nil, errors.Wrap(types.ErrCreateDisabled, "failed to create new contract")
	} else if !cfg.Params.EnableCall && msg.To() != nil {
		return nil, errors.Wrap(types.ErrCallDisabled, "failed to call contract")
	}

	stateDB := statedb.New(ctx, k, txConfig)
	evm := k.NewEVM(ctx, msg, cfg, tracer, stateDB)

	defer func() { evm.Restore() }()

	sender := vm.AccountRef(msg.From())
	contractCreation := msg.To() == nil
	isLondon := cfg.ChainConfig.IsLondon(evm.Context.BlockNumber)

	intrinsicGas, err := k.GetEthIntrinsicGas(ctx, msg, cfg.ChainConfig, contractCreation)
	if err != nil {
		// should have already been checked on Ante Handler
		return nil, errors.Wrap(err, "intrinsic gas failed")
	}
	// Should check again even if it is checked on Ante Handler, because eth_call don't go through Ante Handler.
	if msg.Gas() < intrinsicGas {
		// eth_estimateGas will check for this exact error
		return nil, errors.Wrap(core.ErrIntrinsicGas, "apply message")
	}
	leftoverGas := msg.Gas() - intrinsicGas

	// access list preparation is moved from ante handler to here, because it's needed when `ApplyMessage` is called
	// under contexts where ante handlers are not run, for example `eth_call` and `eth_estimateGas`.
	if rules := cfg.ChainConfig.Rules(big.NewInt(ctx.BlockHeight()), cfg.ChainConfig.MergeNetsplitBlock != nil); rules.IsBerlin {
		evm.StateDB.PrepareAccessList(msg.From(), msg.To(), vm.ActivePrecompiles(rules), msg.AccessList())
	}

	// NOTE: In order to achieve this, nonce should be checked in ante handler and increased, otherwise
	// it could make potential nonce override with double spend or contract rewrite
	// take over the nonce management from evm:
	// reset sender's nonce to msg.Nonce() before calling evm on msg nonce
	// as nonce already increased in db
	evm.StateDB.SetNonce(sender.Address(), msg.Nonce())

	if contractCreation {
		// no need to increase nonce here as contract as during contract creation:
		// - tx.origin nonce increase automatically
		// - if IsEIP158 enabled, contract nonce will be set as 1
		ret, _, leftoverGas, vmErr = evm.Create(sender, msg.Data(), leftoverGas, msg.Value())
	} else {
		// should be increased before call on nonce from msg, so we make sure nonce remaining same as on init tx
		evm.StateDB.SetNonce(sender.Address(), msg.Nonce()+1)
		ret, leftoverGas, vmErr = evm.Call(sender, *msg.To(), msg.Data(), leftoverGas, msg.Value())
	}

	// NEW: Kill tx for cosmos msg to abort everything
	if commit && (vmErr != nil || evm.Cancelled()) && evm.Killed() {
		return nil, errors.Wrap(vmErr, "cosmos failure")
	}

	refundQuotient := params.RefundQuotient

	// After EIP-3529: refunds are capped to gasUsed / 5
	if isLondon {
		refundQuotient = params.RefundQuotientEIP3529
	}

	// calculate gas refund
	if msg.Gas() < leftoverGas {
		return nil, errors.Wrap(types.ErrGasOverflow, "apply message")
	}
	gasUsed := msg.Gas() - leftoverGas
	refund := GasToRefund(evm.StateDB.GetRefund(), gasUsed, refundQuotient)
	if refund > gasUsed {
		return nil, errors.Wrap(types.ErrGasOverflow, "apply message")
	}
	gasUsed -= refund

	// EVM execution error needs to be available for the JSON-RPC client
	var vmError string
	if vmErr != nil {
		vmError = vmErr.Error()
	}

	// The dirty states in `StateDB` is either committed or discarded after return
	if commit {
		if err := evm.StateDB.Commit(); err != nil {
			return nil, errors.Wrap(err, "failed to commit stateDB")
		}
	}

	k.ApplyEvents(ctx, vmErr != nil)

	return &types.MsgEthereumTxResponse{
		GasUsed: gasUsed,
		VmError: vmError,
		Ret:     ret,
		Logs:    types.NewLogsFromEth(evm.StateDB.Logs()),
		Hash:    txConfig.TxHash.Hex(),
	}, nil
}

// ApplyMessage calls ApplyMessageWithConfig with default EVMConfig
func (k *Keeper) ApplyMessage(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*types.MsgEthereumTxResponse, error) {
	cfg, err := k.EVMConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load evm config")
	}
	txConfig := statedb.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash()))
	return k.ApplyMessageWithConfig(ctx, msg, tracer, commit, cfg, txConfig)
}

// GetEthIntrinsicGas returns the intrinsic gas cost for the transaction
func (k *Keeper) GetEthIntrinsicGas(ctx sdk.Context, msg core.Message, cfg *params.ChainConfig, isContractCreation bool) (uint64, error) {
	height := big.NewInt(ctx.BlockHeight())
	homestead := cfg.IsHomestead(height)
	istanbul := cfg.IsIstanbul(height)

	return core.IntrinsicGas(msg.Data(), msg.AccessList(), isContractCreation, homestead, istanbul)
}

// RefundGas transfers the leftover gas to the sender of the message, caped to half of the total gas
// consumed in the transaction. Additionally, the function sets the total gas consumed to the value
// returned by the EVM execution, thus ignoring the previous intrinsic gas consumed during in the
// AnteHandler.
func (k *Keeper) RefundGas(ctx sdk.Context, msg core.Message, leftoverGas uint64, denom string) error {
	// Return EVM tokens for remaining gas, exchanged at the original rate.
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(leftoverGas), msg.GasPrice())

	switch remaining.Sign() {
	case -1:
		// negative refund errors
		return errors.Wrapf(types.ErrInvalidRefund, "refunded amount value cannot be negative %d", remaining.Int64())
	case 1:
		// positive amount refund
		refundedCoins := sdk.Coins{sdk.NewCoin(denom, sdkmath.NewIntFromBigInt(remaining))}

		// refund to sender from the fee collector module account, which is the escrow account in charge of collecting tx fees

		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, msg.From().Bytes(), refundedCoins)
		if err != nil {
			err = errors.Wrapf(sdkerrors.ErrInsufficientFunds, "fee collector account failed to refund fees: %s", err.Error())
			return errors.Wrapf(err, "failed to refund %d leftover gas (%s)", leftoverGas, refundedCoins.String())
		}
	default:
		// no refund, consume gas and update the tx gas meter
	}

	return nil
}

// ResetGasMeterAndConsumeGas reset first the gas meter consumed value to zero and set it back to the new value
// 'gasUsed'
func (k *Keeper) ResetGasMeterAndConsumeGas(ctx sdk.Context, gasUsed uint64) {
	// reset the gas count
	ctx.GasMeter().RefundGas(ctx.GasMeter().GasConsumed(), "reset the gas count")
	ctx.GasMeter().ConsumeGas(gasUsed, "apply evm transaction")
}

// GetCoinbaseAddress returns the block proposer's validator operator address.
func (k Keeper) GetCoinbaseAddress(ctx sdk.Context) (common.Address, error) {
	consAddr := sdk.ConsAddress(ctx.BlockHeader().ProposerAddress)
	validator, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		return common.Address{}, errors.Wrapf(
			stakingtypes.ErrNoValidatorFound,
			"failed to retrieve validator from block proposer address %s",
			consAddr.String(),
		)
	}

	coinbase := common.BytesToAddress(validator.GetOperator())
	return coinbase, nil
}
