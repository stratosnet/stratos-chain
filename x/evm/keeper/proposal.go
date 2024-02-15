package keeper

import (
	"bytes"
	"fmt"
	"math"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/evm/statedb"
	"github.com/stratosnet/stratos-chain/x/evm/tracers"
	"github.com/stratosnet/stratos-chain/x/evm/types"
	"github.com/stratosnet/stratos-chain/x/evm/vm"
)

var (
	emptyCodeHash = crypto.Keccak256Hash(nil)
)

type MigratableContract struct {
	Address common.Address
	Bin     string
	Init    string
}

type ProposalCounsil struct {
	keeper         *Keeper
	ctx            sdk.Context
	stateDB        *statedb.StateDB
	evm            *vm.EVM
	consensusOwner common.Address
	proxyOwner     common.Address
	verifier       *vm.GenesisContractVerifier
}

func NewProposalCounsil(k *Keeper, ctx sdk.Context) (*ProposalCounsil, error) {
	pc := &ProposalCounsil{
		keeper:         k,
		ctx:            ctx,
		consensusOwner: vm.ConsensusAddress,
		proxyOwner:     vm.ProxyOwnerAddress,
		verifier:       k.verifier,
	}
	cfg, err := k.EVMConfig(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to load evm config")
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
	pc.evm = vm.NewEVM(blockCtx, txCtx, pc.stateDB, cfg.ChainConfig, vmConfig, pc.verifier)

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

// Migrate1to2 for initialize contract on v0.12
func (pc *ProposalCounsil) Migrate1to2() error {
	contracts := []*MigratableContract{
		{
			Address: vm.PrepayContractAddress,
			Bin:     "0x608060405234801561001057600080fd5b50610b19806100206000396000f3fe6080604052600436106100595760003560e01c806334fe1d1e14610065578063715018a61461006f5780638129fc1c146100865780638da5cb5b1461009d578063f2fde38b146100c8578063ffa1ad74146100f157610060565b3661006057005b600080fd5b61006d61011c565b005b34801561007b57600080fd5b50610084610241565b005b34801561009257600080fd5b5061009b610255565b005b3480156100a957600080fd5b506100b261039b565b6040516100bf91906106d1565b60405180910390f35b3480156100d457600080fd5b506100ef60048036038101906100ea919061071d565b6103c5565b005b3480156100fd57600080fd5b50610106610448565b6040516101139190610766565b60405180910390f35b600034905060008103610164576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161015b906107de565b60405180910390fd5b61016c61066e565b3373ffffffffffffffffffffffffffffffffffffffff1681600060018110610197576101966107fe565b5b6020020181815250506101a861066e565b6020816020848660f1600019f16101be57600080fd5b823373ffffffffffffffffffffffffffffffffffffffff163073ffffffffffffffffffffffffffffffffffffffff167fa9fdf2e446d7225a2b445bc7c21ca59dcea69b5b23f5c4e6f54f87a5db6cdaee84600060018110610222576102216107fe565b5b60200201516040516102349190610846565b60405180910390a4505050565b61024961044d565b61025360006104cb565b565b60008060019054906101000a900460ff161590508080156102865750600160008054906101000a900460ff1660ff16105b806102b3575061029530610591565b1580156102b25750600160008054906101000a900460ff1660ff16145b5b6102f2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102e9906108d3565b60405180910390fd5b60016000806101000a81548160ff021916908360ff160217905550801561032f576001600060016101000a81548160ff0219169083151502179055505b6103376105b4565b61033f610605565b80156103985760008060016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498600160405161038f9190610938565b60405180910390a15b50565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6103cd61044d565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361043c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610433906109c5565b60405180910390fd5b610445816104cb565b50565b600081565b610455610666565b73ffffffffffffffffffffffffffffffffffffffff1661047361039b565b73ffffffffffffffffffffffffffffffffffffffff16146104c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104c090610a31565b60405180910390fd5b565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081603360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b600060019054906101000a900460ff16610603576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105fa90610ac3565b60405180910390fd5b565b600060019054906101000a900460ff16610654576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161064b90610ac3565b60405180910390fd5b61066461065f610666565b6104cb565b565b600033905090565b6040518060200160405280600190602082028036833780820191505090505090565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006106bb82610690565b9050919050565b6106cb816106b0565b82525050565b60006020820190506106e660008301846106c2565b92915050565b600080fd5b6106fa816106b0565b811461070557600080fd5b50565b600081359050610717816106f1565b92915050565b600060208284031215610733576107326106ec565b5b600061074184828501610708565b91505092915050565b600060ff82169050919050565b6107608161074a565b82525050565b600060208201905061077b6000830184610757565b92915050565b600082825260208201905092915050565b7f503a205a45524f5f414d4f554e54000000000000000000000000000000000000600082015250565b60006107c8600e83610781565b91506107d382610792565b602082019050919050565b600060208201905081810360008301526107f7816107bb565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000819050919050565b6108408161082d565b82525050565b600060208201905061085b6000830184610837565b92915050565b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b60006108bd602e83610781565b91506108c882610861565b604082019050919050565b600060208201905081810360008301526108ec816108b0565b9050919050565b6000819050919050565b6000819050919050565b600061092261091d610918846108f3565b6108fd565b61074a565b9050919050565b61093281610907565b82525050565b600060208201905061094d6000830184610929565b92915050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b60006109af602683610781565b91506109ba82610953565b604082019050919050565b600060208201905081810360008301526109de816109a2565b9050919050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b6000610a1b602083610781565b9150610a26826109e5565b602082019050919050565b60006020820190508181036000830152610a4a81610a0e565b9050919050565b7f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960008201527f6e697469616c697a696e67000000000000000000000000000000000000000000602082015250565b6000610aad602b83610781565b9150610ab882610a51565b604082019050919050565b60006020820190508181036000830152610adc81610aa0565b905091905056fea2646970667358221220af723f0b659da0cdaa0d447e54bb8c0c3af22886c04139bf5fc7ee74f487451c64736f6c63430008120033",
			Init:    "0x8129fc1c",
		},
	}

	for _, contract := range contracts {
		implAddr := crypto.CreateAddress(pc.consensusOwner, pc.stateDB.GetNonce(pc.consensusOwner))
		proxyAddr := contract.Address

		bin, err := hexutil.Decode(contract.Bin)
		if err != nil {
			return err
		}

		data, err := hexutil.Decode(contract.Init)
		if err != nil {
			return err
		}

		implCode := pc.stateDB.GetCode(implAddr)
		if implCode == nil {
			if _, err := pc.create(pc.proxyOwner, implAddr, bin, nil); err != nil {
				return errorsmod.Wrapf(err, "failed to get or create address on '%s'", implAddr)
			}
		}

		value := sdk.NewInt(0)
		c := types.NewUpdateImplmentationProposal(
			proxyAddr.Hex(),
			implAddr.Hex(),
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
