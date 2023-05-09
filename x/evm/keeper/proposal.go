package keeper

import (
	"bytes"
	"fmt"
	"math"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"

	keestatedb "github.com/stratosnet/stratos-chain/core/statedb"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/evm/statedb"
	"github.com/stratosnet/stratos-chain/x/evm/tracers"
	"github.com/stratosnet/stratos-chain/x/evm/types"
	"github.com/stratosnet/stratos-chain/x/evm/vm"
)

var (
	emptyCodeHash = crypto.Keccak256Hash(nil)
)

type ProposalCounsil struct {
	keeper         *Keeper
	ctx            sdk.Context
	stateDB        *statedb.StateDB
	evm            *vm.EVM
	consensusOwner common.Address
	proxyOwner     common.Address
	verifier       *vm.GenesisContractVerifier
}

func NewProposalCounsil(k Keeper, ctx sdk.Context) (*ProposalCounsil, error) {
	params := k.GetParams(ctx)

	pc := &ProposalCounsil{
		keeper:         &k,
		ctx:            ctx,
		consensusOwner: common.HexToAddress(params.ProxyProposalParams.ConsensusAddress),
		proxyOwner:     common.HexToAddress(params.ProxyProposalParams.ProxyOwnerAddress),
		verifier:       k.verifier,
	}
	cfg, err := k.EVMConfig(ctx)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to load evm config")
	}

	blockCtx := vm.BlockContext{
		CanTransfer: vm.CanTransfer,
		Transfer:    vm.Transfer,
		GetHash:     k.GetHashFn(ctx),
		Coinbase:    cfg.CoinBase,
		GasLimit:    stratos.BlockGasLimit(ctx),
		BlockNumber: big.NewInt(ctx.BlockHeight()),
		Time:        big.NewInt(ctx.BlockHeader().Time.Unix()),
		Difficulty:  big.NewInt(0), // unused. Only required in PoW context
		BaseFee:     cfg.BaseFee,
	}

	txCtx := vm.TxContext{
		Origin:   pc.consensusOwner,
		GasPrice: big.NewInt(0),
	}
	tracer := tracers.NewNoOpTracer()
	vmConfig := k.VMConfig(ctx, cfg, tracer)

	txConfig := statedb.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash()))
	pc.stateDB = statedb.New(ctx, pc.keeper, txConfig)
	kstatedb := keestatedb.New(ctx)
	pc.evm = vm.NewEVM(blockCtx, txCtx, pc.stateDB, kstatedb, cfg.ChainConfig, vmConfig, pc.verifier)

	return pc, nil
}

func (pc *ProposalCounsil) finalize() error {
	if err := pc.stateDB.Commit(); err != nil {
		return err
	}
	return nil
}

func (pc *ProposalCounsil) call(sender, contractAddress common.Address, data []byte, value *big.Int) error {
	// required
	txCtx := vm.TxContext{
		Origin:   sender,
		GasPrice: big.NewInt(0),
	}
	pc.evm.Reset(txCtx, pc.stateDB)

	nonce := pc.stateDB.GetNonce(sender)
	// we do not care about gas during consil execution
	gas := uint64(math.MaxUint64)
	// for safety
	if value == nil {
		value = big.NewInt(0)
	}

	{
		pc.stateDB.SetNonce(sender, nonce+1)
	}

	if _, _, vmErr := pc.evm.Call(vm.AccountRef(sender), contractAddress, data, gas, value); vmErr != nil {
		return vmErr
	}
	return nil
}

func (pc *ProposalCounsil) create(sender, contractAddress common.Address, data []byte, value *big.Int) (*common.Address, error) {
	// required
	txCtx := vm.TxContext{
		Origin:   sender,
		GasPrice: big.NewInt(0),
	}
	pc.evm.Reset(txCtx, pc.stateDB)

	nonce := pc.stateDB.GetNonce(sender)
	// we do not care about gas during consil execution
	gas := uint64(math.MaxUint64)

	// for safety
	if value == nil {
		value = big.NewInt(0)
	}

	interpreter := vm.NewEVMInterpreter(pc.evm, pc.evm.Config)

	accRef := vm.AccountRef(sender)

	{
		pc.stateDB.SetNonce(accRef.Address(), nonce+1)
	}

	contractHash := pc.evm.StateDB.GetCodeHash(contractAddress)
	if pc.evm.StateDB.GetNonce(contractAddress) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		return nil, vm.ErrContractAddressCollision
	}

	snapshot := pc.evm.StateDB.Snapshot()

	pc.evm.StateDB.CreateAccount(contractAddress)
	pc.evm.StateDB.SetNonce(contractAddress, 1)
	pc.evm.Context.Transfer(pc.evm.StateDB, sender, contractAddress, value)

	contract := vm.NewContract(accRef, vm.AccountRef(contractAddress), value, gas)
	contract.SetCallCode(&contractAddress, common.Hash{}, data)

	ret, err := interpreter.Run(contract, nil, false)
	if err != nil {
		return nil, err
	}

	// Check whether the max code size has been exceeded, assign err if the case.
	if err == nil && len(ret) > params.MaxCodeSize {
		err = vm.ErrMaxCodeSizeExceeded
	}

	// Reject code starting with 0xEF if EIP-3541 is enabled.
	if err == nil && len(ret) >= 1 && ret[0] == 0xEF {
		err = vm.ErrInvalidCode
	}

	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil {
		createDataGas := uint64(len(ret)) * params.CreateDataGas
		if contract.UseGas(createDataGas) {
			pc.evm.StateDB.SetCode(contractAddress, ret)
		} else {
			err = vm.ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil {
		pc.evm.StateDB.RevertToSnapshot(snapshot)
		if err != vm.ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}

	return &contractAddress, nil
}

func (pc *ProposalCounsil) ApplyGenesisState(height uint64) error {
	contracts := pc.verifier.GetContracts(height)

	if len(contracts) == 0 {
		return nil
	}

	for _, contract := range contracts {
		implAddr := crypto.CreateAddress(pc.consensusOwner, pc.stateDB.GetNonce(pc.consensusOwner))
		proxyAddr := common.HexToAddress(contract.GetAddress())

		bin, err := hexutil.Decode(contract.GetBin())
		if err != nil {
			return err
		}

		data, err := hexutil.Decode(contract.GetInit())
		if err != nil {
			return err
		}

		implCode := pc.stateDB.GetCode(implAddr)
		if implCode == nil {
			if _, err := pc.create(pc.proxyOwner, implAddr, bin, nil); err != nil {
				return sdkerrors.Wrapf(err, "failed to get or create address on '%s'", implAddr)
			}
		}

		value := sdk.NewInt(0)
		c := types.NewUpdateImplmentationProposal(
			proxyAddr,
			implAddr,
			data,
			&value,
		)

		if err := c.ValidateBasic(); err != nil {
			return err
		}

		if err = pc.updateProxyImplementation(c.(*types.UpdateImplmentationProposal), false); err != nil {
			return err
		}
	}
	if err := pc.finalize(); err != nil {
		return err
	}
	return nil
}

func (pc *ProposalCounsil) UpdateProxyImplementation(p *types.UpdateImplmentationProposal) error {
	return pc.updateProxyImplementation(p, true)
}

// updateProxyImplementation execute provided function to a proxy contract for impl upgrade
func (pc *ProposalCounsil) updateProxyImplementation(p *types.UpdateImplmentationProposal, commit bool) error {
	if !pc.verifier.IsTrustedAddress(p.ProxyAddress) {
		return fmt.Errorf("proxy '%s' has not been verified", p.ProxyAddress)
	}

	proxyAddress := common.HexToAddress(p.ProxyAddress)
	implAddress := common.HexToAddress(p.ImplementationAddress)

	implCode := pc.stateDB.GetCode(implAddress)
	if implCode == nil {
		return fmt.Errorf("implementation '%s' not found", implAddress)
	}

	if bytes.Equal(implCode, emptyCodeHash[:]) {
		return fmt.Errorf("implementation '%s' is EOA", implAddress)
	}

	proxyCode := pc.stateDB.GetCode(proxyAddress)
	if proxyCode != nil {
		upgradeData, err := types.EncodeContractFunc(
			types.TransparentUpgradableProxyABI,
			"upgradeToAndCall",
			implAddress,
			p.Data,
		)
		if err != nil {
			return err
		}

		if err := pc.call(pc.consensusOwner, proxyAddress, upgradeData, p.Amount.BigInt()); err != nil {
			return err
		}
	} else {
		proxyConstructorData, err := types.EncodeContractFunc(
			types.TransparentUpgradableProxyABI,
			"",
			implAddress,
			pc.consensusOwner,
			p.Data,
		)
		if err != nil {
			return err
		}

		proxyConstructorData = append(common.FromHex(types.TransparentUpgradableProxyBin), proxyConstructorData...)

		if _, err := pc.create(pc.proxyOwner, proxyAddress, proxyConstructorData, p.Amount.BigInt()); err != nil {
			return err
		}
	}

	if commit {
		if err := pc.finalize(); err != nil {
			return err
		}
	}

	return nil
}
