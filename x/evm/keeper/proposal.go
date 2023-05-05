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
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/evm/statedb"
	"github.com/stratosnet/stratos-chain/x/evm/types"
)

var (
	emptyCodeHash = crypto.Keccak256Hash(nil)
)

type ProposalCounsil struct {
	keeper     *Keeper
	ctx        sdk.Context
	stateDB    *statedb.StateDB
	evm        *vm.EVM
	sender     common.Address
	proxyAdmin common.Address
	verifier   *ProposalVerifier
}

func NewProposalCounsil(k Keeper, ctx sdk.Context) (*ProposalCounsil, error) {
	params := k.GetParams(ctx)

	pc := &ProposalCounsil{
		keeper:     &k,
		ctx:        ctx,
		sender:     common.HexToAddress(params.ProxyProposalParams.ConsensusAddress),
		proxyAdmin: common.HexToAddress(params.ProxyProposalParams.ProxyAdminAddress),
		verifier:   NewProposalVerifier(),
	}
	cfg, err := k.EVMConfig(ctx)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to load evm config")
	}

	defer pc.prepare(params)

	blockCtx := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     k.GetHashFn(ctx),
		Coinbase:    cfg.CoinBase,
		GasLimit:    stratos.BlockGasLimit(ctx),
		BlockNumber: big.NewInt(ctx.BlockHeight()),
		Time:        big.NewInt(ctx.BlockHeader().Time.Unix()),
		Difficulty:  big.NewInt(0), // unused. Only required in PoW context
		BaseFee:     cfg.BaseFee,
	}

	txCtx := vm.TxContext{
		Origin:   pc.sender,
		GasPrice: big.NewInt(0),
	}
	tracer := types.NewNoOpTracer()
	vmConfig := k.VMConfig(ctx, cfg, tracer)

	txConfig := statedb.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash()))
	pc.stateDB = statedb.New(ctx, pc.keeper, txConfig)
	pc.evm = vm.NewEVM(blockCtx, txCtx, pc.stateDB, cfg.ChainConfig, vmConfig)

	return pc, nil
}

func (pc *ProposalCounsil) prepare(params types.Params) {
	nonce := pc.keeper.GetNonce(pc.ctx, pc.sender)
	pc.stateDB.SetNonce(pc.sender, nonce)

	pc.verifier.ApplyParamsState(params)
}

func (pc *ProposalCounsil) finalize() error {
	if err := pc.stateDB.Commit(); err != nil {
		return err
	}
	return nil
}

func (pc *ProposalCounsil) create(contractAddress common.Address, data []byte, value *big.Int) (*common.Address, error) {
	nonce := pc.stateDB.GetNonce(pc.sender)
	// we do not care about gas during consil execution
	gas := uint64(math.MaxUint64)

	// for safety
	if value == nil {
		value = big.NewInt(0)
	}

	interpreter := vm.NewEVMInterpreter(pc.evm, pc.evm.Config)

	sender := vm.AccountRef(pc.sender)

	{
		pc.stateDB.SetNonce(sender.Address(), nonce+1)
	}

	contractHash := pc.evm.StateDB.GetCodeHash(contractAddress)
	if pc.evm.StateDB.GetNonce(contractAddress) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		return nil, vm.ErrContractAddressCollision
	}

	snapshot := pc.evm.StateDB.Snapshot()

	pc.evm.StateDB.CreateAccount(contractAddress)
	pc.evm.StateDB.SetNonce(contractAddress, 1)
	pc.evm.Context.Transfer(pc.evm.StateDB, pc.sender, contractAddress, value)

	contract := vm.NewContract(sender, vm.AccountRef(contractAddress), value, gas)
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

func (pc *ProposalCounsil) getOrCreateContract(contractAddress common.Address, data []byte) (*common.Address, error) {
	code := pc.stateDB.GetCode(contractAddress)
	if code != nil {
		return &contractAddress, nil
	}

	if _, err := pc.create(contractAddress, data, big.NewInt(0)); err != nil {
		return nil, err
	}
	return &contractAddress, nil
}

func (pc *ProposalCounsil) ApplyGenesisState(height uint64) error {
	states := pc.verifier.GetStates(height)

	for _, state := range states {
		addr := crypto.CreateAddress(pc.sender, pc.stateDB.GetNonce(pc.sender))

		proxyAddr := common.HexToAddress(state.Address)
		bin, err := hexutil.Decode(state.Bin)
		if err != nil {
			return err
		}

		implAddr, err := pc.getOrCreateContract(addr, bin)
		if err != nil {
			return sdkerrors.Wrapf(err, "failed to get or create address on '%s'", addr)
		}

		amount := sdk.NewInt(0)
		data, err := hexutil.Decode(state.Init)
		if err != nil {
			return err
		}

		c := types.NewUpdateImplmentationProposal(
			proxyAddr,
			*implAddr,
			data,
			&amount,
		)

		err = pc.updateProxyImplementation(c.(*types.UpdateImplmentationProposal), false)
		if err != nil {
			return err
		}
	}
	if len(states) > 0 {
		if err := pc.finalize(); err != nil {
			return err
		}
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

	proxyAdminConstructorData, err := hexutil.Decode(types.ProxyAdminBin)
	if err != nil {
		return err
	}

	_, err = pc.getOrCreateContract(pc.proxyAdmin, proxyAdminConstructorData)
	if err != nil {
		return err
	}

	proxyConstructorData, err := types.EncodeContractFunc(
		types.TransparentUpgradableProxyABI,
		"",
		implAddress,
		pc.proxyAdmin,
		[]byte{},
	)
	if err != nil {
		return err
	}

	proxyConstructorData = append(common.FromHex(types.TransparentUpgradableProxyBin), proxyConstructorData...)

	_, err = pc.getOrCreateContract(proxyAddress, proxyConstructorData)
	if err != nil {
		return err
	}

	upgradeData, err := types.EncodeContractFunc(
		types.ProxyAdminABI,
		"upgradeAndCall",
		proxyAddress,
		implAddress,
		p.Data,
	)
	if err != nil {
		return err
	}

	nonce := pc.stateDB.GetNonce(pc.sender)
	gas := uint64(math.MaxUint64)
	value := p.Amount.BigInt()

	{
		pc.stateDB.SetNonce(pc.sender, nonce+1)
	}

	if _, _, vmErr := pc.evm.Call(vm.AccountRef(pc.sender), pc.proxyAdmin, upgradeData, gas, value); vmErr != nil {
		return vmErr
	}

	if commit {
		if err := pc.finalize(); err != nil {
			return err
		}
	}

	return nil
}
