package keeper

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/evm/statedb"
	"github.com/stratosnet/stratos-chain/x/evm/types"
)

var emptyCodeHash = crypto.Keccak256Hash(nil)

type ProposalCounsil struct {
	keeper     *Keeper
	ctx        sdk.Context
	stateDB    *statedb.StateDB
	evm        *vm.EVM
	sender     common.Address
	proxyAdmin common.Address
}

func NewProposalCounsil(k Keeper, ctx sdk.Context) (*ProposalCounsil, error) {
	params := k.GetParams(ctx)

	pc := &ProposalCounsil{
		keeper:     &k,
		ctx:        ctx,
		sender:     common.HexToAddress(params.ProxyProposalParams.ConsensusDeployerAddress),
		proxyAdmin: common.HexToAddress(params.ProxyProposalParams.ProxyAdminAddress),
	}
	cfg, err := k.EVMConfig(ctx)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to load evm config")
	}

	defer pc.prepare()

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

func (pc *ProposalCounsil) prepare() {
	nonce := pc.keeper.GetNonce(pc.ctx, pc.sender)
	pc.stateDB.SetNonce(pc.sender, nonce)
}

func (pc *ProposalCounsil) finalize() error {
	if err := pc.stateDB.Commit(); err != nil {
		return err
	}
	return nil
}

func (pc *ProposalCounsil) create(contractAddress common.Address, data []byte, value *big.Int) (*common.Address, error) {
	nonce := pc.stateDB.GetNonce(pc.sender)
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
	fmt.Printf("Ret of exec: %s\n", ret)
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

// ExecuteProxyFunc execute provided function to a proxy contract
func (pc *ProposalCounsil) UpdateProxyImplementation(p *types.UpdateImplmentationProposal) error {
	// TODO: PROXY: Add validation for this address
	proxyAddress := common.HexToAddress(p.ProxyAddress)
	implAddress := common.HexToAddress(p.ImplementationAddress)

	proxyAdminConstructorData := common.FromHex(types.ProxyAdminBin)
	_, err := pc.getOrCreateContract(pc.proxyAdmin, proxyAdminConstructorData)
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

	_, _, vmErr := pc.evm.Call(vm.AccountRef(pc.sender), pc.proxyAdmin, upgradeData, gas, value)
	if vmErr != nil {
		return vmErr
	}

	if err := pc.finalize(); err != nil {
		return err
	}

	return nil
}

// TODO: PROXY: Remove this, just for testing purposes before test writing
func (pc *ProposalCounsil) TestDeployERC20Mock() (*common.Address, error) {
	erc20Data := common.FromHex("0x60806040523480156200001157600080fd5b506040518060400160405280600981526020016845524332304d6f636b60b81b815250604051806040016040528060048152602001634532304d60e01b81525081600390816200006291906200011f565b5060046200007182826200011f565b505050620001eb565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620000a557607f821691505b602082108103620000c657634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200011a57600081815260208120601f850160051c81016020861015620000f55750805b601f850160051c820191505b81811015620001165782815560010162000101565b5050505b505050565b81516001600160401b038111156200013b576200013b6200007a565b62000153816200014c845462000090565b84620000cc565b602080601f8311600181146200018b5760008415620001725750858301515b600019600386901b1c1916600185901b17855562000116565b600085815260208120601f198616915b82811015620001bc578886015182559484019460019091019084016200019b565b5085821015620001db5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b610aac80620001fb6000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c806340c10f191161008c5780639dc29fac116100665780639dc29fac146101a2578063a457c2d7146101b5578063a9059cbb146101c8578063dd62ed3e146101db57600080fd5b806340c10f191461015c57806370a082311461017157806395d89b411461019a57600080fd5b806306fdde03146100d4578063095ea7b3146100f257806318160ddd1461011557806323b872dd14610127578063313ce5671461013a5780633950935114610149575b600080fd5b6100dc6101ee565b6040516100e991906108ea565b60405180910390f35b61010561010036600461095b565b610280565b60405190151581526020016100e9565b6002545b6040519081526020016100e9565b610105610135366004610985565b610298565b604051601281526020016100e9565b61010561015736600461095b565b6102bc565b61016f61016a36600461095b565b6102de565b005b61011961017f3660046109c1565b6001600160a01b031660009081526020819052604090205490565b6100dc6102ec565b61016f6101b036600461095b565b6102fb565b6101056101c336600461095b565b610305565b6101056101d636600461095b565b610385565b6101196101e93660046109e3565b610393565b6060600380546101fd90610a16565b80601f016020809104026020016040519081016040528092919081815260200182805461022990610a16565b80156102765780601f1061024b57610100808354040283529160200191610276565b820191906000526020600020905b81548152906001019060200180831161025957829003601f168201915b5050505050905090565b60003361028e8185856103be565b5060019392505050565b6000336102a68582856104e3565b6102b185858561055d565b506001949350505050565b60003361028e8185856102cf8383610393565b6102d99190610a50565b6103be565b6102e88282610701565b5050565b6060600480546101fd90610a16565b6102e882826107c0565b600033816103138286610393565b9050838110156103785760405162461bcd60e51b815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f77604482015264207a65726f60d81b60648201526084015b60405180910390fd5b6102b182868684036103be565b60003361028e81858561055d565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6001600160a01b0383166104205760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161036f565b6001600160a01b0382166104815760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161036f565b6001600160a01b0383811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b60006104ef8484610393565b90506000198114610557578181101561054a5760405162461bcd60e51b815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000604482015260640161036f565b61055784848484036103be565b50505050565b6001600160a01b0383166105c15760405162461bcd60e51b815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f206164604482015264647265737360d81b606482015260840161036f565b6001600160a01b0382166106235760405162461bcd60e51b815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201526265737360e81b606482015260840161036f565b6001600160a01b0383166000908152602081905260409020548181101561069b5760405162461bcd60e51b815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e7420657863656564732062604482015265616c616e636560d01b606482015260840161036f565b6001600160a01b03848116600081815260208181526040808320878703905593871680835291849020805487019055925185815290927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a3610557565b6001600160a01b0382166107575760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640161036f565b80600260008282546107699190610a50565b90915550506001600160a01b038216600081815260208181526040808320805486019055518481527fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35050565b6001600160a01b0382166108205760405162461bcd60e51b815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f206164647265736044820152607360f81b606482015260840161036f565b6001600160a01b038216600090815260208190526040902054818110156108945760405162461bcd60e51b815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e604482015261636560f01b606482015260840161036f565b6001600160a01b0383166000818152602081815260408083208686039055600280548790039055518581529192917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91016104d6565b600060208083528351808285015260005b81811015610917578581018301518582016040015282016108fb565b81811115610929576000604083870101525b50601f01601f1916929092016040019392505050565b80356001600160a01b038116811461095657600080fd5b919050565b6000806040838503121561096e57600080fd5b6109778361093f565b946020939093013593505050565b60008060006060848603121561099a57600080fd5b6109a38461093f565b92506109b16020850161093f565b9150604084013590509250925092565b6000602082840312156109d357600080fd5b6109dc8261093f565b9392505050565b600080604083850312156109f657600080fd5b6109ff8361093f565b9150610a0d6020840161093f565b90509250929050565b600181811c90821680610a2a57607f821691505b602082108103610a4a57634e487b7160e01b600052602260045260246000fd5b50919050565b60008219821115610a7157634e487b7160e01b600052601160045260246000fd5b50019056fea2646970667358221220753ab2d79409cf4bcc5a25417d310be474125fdabed362fdbdec3e0025d8281c64736f6c634300080f0033")

	addr, _ := hex.DecodeString("0000000000000000000000000000000000000000")
	addr[len(addr)-1] = byte('t')
	addr[len(addr)-2] = byte('s')
	addr[len(addr)-3] = byte('e')
	addr[len(addr)-4] = byte('t')
	fmt.Println("qwe addr", common.Bytes2Hex(addr))

	contractAddress, err := pc.getOrCreateContract(common.BytesToAddress(addr), erc20Data)
	if err != nil {
		return nil, err
	}

	if err := pc.finalize(); err != nil {
		return nil, err
	}

	return contractAddress, nil
}
