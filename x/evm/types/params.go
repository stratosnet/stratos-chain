package types

import (
	"fmt"
	"math/big"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"

	stratos "github.com/stratosnet/stratos-chain/types"
)

var _ paramtypes.ParamSet = &Params{}

const (
	DefaultEVMDenom = stratos.Wei
)

// Parameter keys
var (
	ParamStoreKeyEVMDenom     = []byte("EVMDenom")
	ParamStoreKeyEnableCreate = []byte("EnableCreate")
	ParamStoreKeyEnableCall   = []byte("EnableCall")
	ParamStoreKeyExtraEIPs    = []byte("EnableExtraEIPs")
	ParamStoreKeyChainConfig  = []byte("ChainConfig")

	// AvailableExtraEIPs define the list of all EIPs that can be enabled by the
	// EVM interpreter. These EIPs are applied in order and can override the
	// instruction sets from the latest hard fork enabled by the ChainConfig. For
	// more info check:
	// https://github.com/ethereum/go-ethereum/blob/master/core/vm/interpreter.go#L97
	AvailableExtraEIPs = []int64{1344, 1884, 2200, 2929, 3198, 3529}

	// fee market
	ParamStoreKeyNoBaseFee                = []byte("NoBaseFee")
	ParamStoreKeyBaseFeeChangeDenominator = []byte("BaseFeeChangeDenominator")
	ParamStoreKeyElasticityMultiplier     = []byte("ElasticityMultiplier")
	ParamStoreKeyBaseFee                  = []byte("BaseFee")
	ParamStoreKeyEnableHeight             = []byte("EnableHeight")

	// proposal proxy
	ParamStoreKeyConsensusAddress  = []byte("ProposalConsensusAddress")
	ParamStoreKeyProxyAdminAddress = []byte("ProxyAdminAddress")
	ParamStoreKeyContracts         = []byte("ProposalContracts")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(evmDenom string, enableCreate, enableCall bool, config ChainConfig, feeMarketParams FeeMarketParams, proxyProposalParams ProxyProposalParams, extraEIPs ...int64) Params {
	return Params{
		EvmDenom:            evmDenom,
		EnableCreate:        enableCreate,
		EnableCall:          enableCall,
		ExtraEIPs:           extraEIPs,
		ChainConfig:         config,
		FeeMarketParams:     feeMarketParams,
		ProxyProposalParams: proxyProposalParams,
	}
}

// DefaultParams returns default evm parameters
// ExtraEIPs is empty to prevent overriding the latest hard fork instruction set
func DefaultParams() Params {
	return Params{
		EvmDenom:            DefaultEVMDenom,
		EnableCreate:        true,
		EnableCall:          true,
		ChainConfig:         DefaultChainConfig(),
		ExtraEIPs:           nil,
		FeeMarketParams:     DefaultFeeMarketParams(),
		ProxyProposalParams: DefaultProxyProposalParams(),
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyEVMDenom, &p.EvmDenom, validateEVMDenom),
		paramtypes.NewParamSetPair(ParamStoreKeyEnableCreate, &p.EnableCreate, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyEnableCall, &p.EnableCall, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyExtraEIPs, &p.ExtraEIPs, validateEIPs),
		paramtypes.NewParamSetPair(ParamStoreKeyChainConfig, &p.ChainConfig, validateChainConfig),
		//fee market
		paramtypes.NewParamSetPair(ParamStoreKeyNoBaseFee, &p.FeeMarketParams.NoBaseFee, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyBaseFeeChangeDenominator, &p.FeeMarketParams.BaseFeeChangeDenominator, validateBaseFeeChangeDenominator),
		paramtypes.NewParamSetPair(ParamStoreKeyElasticityMultiplier, &p.FeeMarketParams.ElasticityMultiplier, validateElasticityMultiplier),
		paramtypes.NewParamSetPair(ParamStoreKeyBaseFee, &p.FeeMarketParams.BaseFee, validateBaseFee),
		paramtypes.NewParamSetPair(ParamStoreKeyEnableHeight, &p.FeeMarketParams.EnableHeight, validateEnableHeight),
		//proposal proxy
		paramtypes.NewParamSetPair(ParamStoreKeyConsensusAddress, &p.ProxyProposalParams.ConsensusAddress, validateAddress),
		paramtypes.NewParamSetPair(ParamStoreKeyProxyAdminAddress, &p.ProxyProposalParams.ProxyAdminAddress, validateAddress),
		paramtypes.NewParamSetPair(ParamStoreKeyContracts, &p.ProxyProposalParams.Contracts, validateContracts),
	}
}

// Validate performs basic validation on evm parameters.
func (p Params) Validate() error {
	if err := sdk.ValidateDenom(p.EvmDenom); err != nil {
		return err
	}

	if err := validateEIPs(p.ExtraEIPs); err != nil {
		return err
	}

	if err := p.FeeMarketParams.Validate(); err != nil {
		return err
	}

	return p.ChainConfig.Validate()
}

// EIPs returns the ExtraEips as a int slice
func (p Params) EIPs() []int {
	eips := make([]int, len(p.ExtraEIPs))
	for i, eip := range p.ExtraEIPs {
		eips[i] = int(eip)
	}
	return eips
}

func validateEVMDenom(i interface{}) error {
	denom, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter EVM denom type: %T", i)
	}

	return sdk.ValidateDenom(denom)
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateEIPs(i interface{}) error {
	eips, ok := i.([]int64)
	if !ok {
		return fmt.Errorf("invalid EIP slice type: %T", i)
	}

	for _, eip := range eips {
		if !vm.ValidEip(int(eip)) {
			return fmt.Errorf("EIP %d is not activateable, valid EIPS are: %s", eip, vm.ActivateableEips())
		}
	}

	return nil
}

func validateChainConfig(i interface{}) error {
	cfg, ok := i.(ChainConfig)
	if !ok {
		return fmt.Errorf("invalid chain config type: %T", i)
	}

	return cfg.Validate()
}

// IsLondon returns if london hardfork is enabled.
func IsLondon(ethConfig *params.ChainConfig, height int64) bool {
	return ethConfig.IsLondon(big.NewInt(height))
}

// creates a new FeeMarketParams instance
func NewFeeMarketParams(noBaseFee bool, baseFeeChangeDenom, elasticityMultiplier uint32, baseFee uint64, enableHeight int64) FeeMarketParams {
	return FeeMarketParams{
		NoBaseFee:                noBaseFee,
		BaseFeeChangeDenominator: baseFeeChangeDenom,
		ElasticityMultiplier:     elasticityMultiplier,
		BaseFee:                  sdk.NewIntFromUint64(baseFee),
		EnableHeight:             enableHeight,
	}
}

// DefaultParams returns default evm parameters
func DefaultFeeMarketParams() FeeMarketParams {
	return NewFeeMarketParams(
		false,
		params.BaseFeeChangeDenominator,
		params.ElasticityMultiplier,
		params.InitialBaseFee,
		0,
	)
}

// Validate performs basic validation on fee market parameters.
func (p FeeMarketParams) Validate() error {
	if p.BaseFeeChangeDenominator == 0 {
		return fmt.Errorf("base fee change denominator cannot be 0")
	}

	if p.BaseFee.IsNegative() {
		return fmt.Errorf("initial base fee cannot be negative: %s", p.BaseFee)
	}

	if p.EnableHeight < 0 {
		return fmt.Errorf("enable height cannot be negative: %d", p.EnableHeight)
	}

	return nil
}

func (p *FeeMarketParams) IsBaseFeeEnabled(height int64) bool {
	return !p.NoBaseFee && height >= p.EnableHeight
}

func validateBaseFeeChangeDenominator(i interface{}) error {
	value, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value == 0 {
		return fmt.Errorf("base fee change denominator cannot be 0")
	}

	return nil
}

func validateElasticityMultiplier(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateBaseFee(i interface{}) error {
	value, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value.IsNegative() {
		return fmt.Errorf("base fee cannot be negative")
	}

	return nil
}

func validateEnableHeight(i interface{}) error {
	value, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value < 0 {
		return fmt.Errorf("enable height cannot be negative: %d", value)
	}

	return nil
}

// creates a new ProxyProposalParams instance
func NewProxyProposalParams(consensusAddress, proxyAdminAddress common.Address) ProxyProposalParams {
	return ProxyProposalParams{
		ConsensusAddress:  consensusAddress.Hex(),
		ProxyAdminAddress: proxyAdminAddress.Hex(),
		Contracts:         make(map[string]*ProxyContractInitState),
	}
}

// ProxyProposalParams returns default proxy parameters
func DefaultProxyProposalParams() ProxyProposalParams {
	ppp := NewProxyProposalParams(
		common.HexToAddress("0x1000000000000000000000000000000000000000"),
		common.HexToAddress("0x1000000000000000000000000000000000000001"),
	)
	// TODO: Maybe create some pretty method?
	ppp.Contracts["prepay"] = &ProxyContractInitState{
		Height:  1,
		Address: "0x1000000000000000000000000000000000010101",
		Bin:     "0x608060405234801561001057600080fd5b506114a9806100206000396000f3fe60806040526004361061008a5760003560e01c8063a846c2fd11610059578063a846c2fd1461010b578063bf22c45714610134578063e50701f414610171578063f2fde38b1461019a578063ffa1ad74146101c357610091565b8063040e21f114610096578063715018a6146100b25780638129fc1c146100c95780638da5cb5b146100e057610091565b3661009157005b600080fd5b6100b060048036038101906100ab9190610d71565b6101ee565b005b3480156100be57600080fd5b506100c7610460565b005b3480156100d557600080fd5b506100de610474565b005b3480156100ec57600080fd5b506100f56105ba565b6040516101029190610dad565b60405180910390f35b34801561011757600080fd5b50610132600480360381019061012d9190610dfe565b6105e4565b005b34801561014057600080fd5b5061015b60048036038101906101569190610dfe565b610824565b6040516101689190610f28565b60405180910390f35b34801561017d57600080fd5b5061019860048036038101906101939190610dfe565b61094c565b005b3480156101a657600080fd5b506101c160048036038101906101bc9190610d71565b6109d4565b005b3480156101cf57600080fd5b506101d8610a57565b6040516101e59190610f5f565b60405180910390f35b60003403610231576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161022890610fd7565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036102a0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161029790611043565b60405180910390fd5b60006102ac6066610a5c565b90506102b86066610a6a565b60006040518060a00160405280838152602001600060028111156102df576102de610e3a565b5b81526020013373ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff16815260200134815250606560008481526020019081526020016000206000820151816000015560208201518160010160006101000a81548160ff0219169083600281111561036857610367610e3a565b5b021790555060408201518160010160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060608201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550608082015181600301559050503373ffffffffffffffffffffffffffffffffffffffff16827f1235076423ee44e5e8c0dae33d57026663a202b06bdc96acad474dd499fcb55085843460405161045393929190611081565b60405180910390a3505050565b610468610a80565b6104726000610afe565b565b60008060019054906101000a900460ff161590508080156104a55750600160008054906101000a900460ff1660ff16105b806104d257506104b430610bc4565b1580156104d15750600160008054906101000a900460ff1660ff16145b5b610511576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105089061112a565b60405180910390fd5b60016000806101000a81548160ff021916908360ff160217905550801561054e576001600060016101000a81548160ff0219169083151502179055505b610556610be7565b61055e610c38565b80156105b75760008060016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249860016040516105ae919061118f565b60405180910390a15b50565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000606560008381526020019081526020016000206040518060a0016040529081600082015481526020016001820160009054906101000a900460ff16600281111561063357610632610e3a565b5b600281111561064557610644610e3a565b5b81526020016001820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600382015481525050905060028081111561071457610713610e3a565b5b8160200151600281111561072b5761072a610e3a565b5b1461076b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610762906111f6565b60405180910390fd5b6000816040015173ffffffffffffffffffffffffffffffffffffffff16826080015160405161079990611247565b60006040518083038185875af1925050503d80600081146107d6576040519150601f19603f3d011682016040523d82523d6000602084013e6107db565b606091505b505090508061081f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610816906112a8565b60405180910390fd5b505050565b61082c610ca1565b606560008381526020019081526020016000206040518060a0016040529081600082015481526020016001820160009054906101000a900460ff16600281111561087957610878610e3a565b5b600281111561088b5761088a610e3a565b5b81526020016001820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016003820154815250509050919050565b610954610a80565b600060029050806065600084815260200190815260200160002060010160006101000a81548160ff0219169083600281111561099357610992610e3a565b5b0217905550817fb17f281799264da36ce0121b71ef89948a9e0f0e7ee16a7a0eb2e0efec3565f0826040516109c891906112c8565b60405180910390a25050565b6109dc610a80565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610a4b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a4290611355565b60405180910390fd5b610a5481610afe565b50565b600081565b600081600001549050919050565b6001816000016000828254019250508190555050565b610a88610c99565b73ffffffffffffffffffffffffffffffffffffffff16610aa66105ba565b73ffffffffffffffffffffffffffffffffffffffff1614610afc576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610af3906113c1565b60405180910390fd5b565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081603360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b600060019054906101000a900460ff16610c36576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c2d90611453565b60405180910390fd5b565b600060019054906101000a900460ff16610c87576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c7e90611453565b60405180910390fd5b610c97610c92610c99565b610afe565b565b600033905090565b6040518060a001604052806000815260200160006002811115610cc757610cc6610e3a565b5b8152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600081525090565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610d3e82610d13565b9050919050565b610d4e81610d33565b8114610d5957600080fd5b50565b600081359050610d6b81610d45565b92915050565b600060208284031215610d8757610d86610d0e565b5b6000610d9584828501610d5c565b91505092915050565b610da781610d33565b82525050565b6000602082019050610dc26000830184610d9e565b92915050565b6000819050919050565b610ddb81610dc8565b8114610de657600080fd5b50565b600081359050610df881610dd2565b92915050565b600060208284031215610e1457610e13610d0e565b5b6000610e2284828501610de9565b91505092915050565b610e3481610dc8565b82525050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60038110610e7a57610e79610e3a565b5b50565b6000819050610e8b82610e69565b919050565b6000610e9b82610e7d565b9050919050565b610eab81610e90565b82525050565b610eba81610d33565b82525050565b60a082016000820151610ed66000850182610e2b565b506020820151610ee96020850182610ea2565b506040820151610efc6040850182610eb1565b506060820151610f0f6060850182610eb1565b506080820151610f226080850182610e2b565b50505050565b600060a082019050610f3d6000830184610ec0565b92915050565b600060ff82169050919050565b610f5981610f43565b82525050565b6000602082019050610f746000830184610f50565b92915050565b600082825260208201905092915050565b7f503a205a45524f5f414d4f554e54000000000000000000000000000000000000600082015250565b6000610fc1600e83610f7a565b9150610fcc82610f8b565b602082019050919050565b60006020820190508181036000830152610ff081610fb4565b9050919050565b7f503a205a45524f5f414444524553530000000000000000000000000000000000600082015250565b600061102d600f83610f7a565b915061103882610ff7565b602082019050919050565b6000602082019050818103600083015261105c81611020565b9050919050565b61106c81610e90565b82525050565b61107b81610dc8565b82525050565b60006060820190506110966000830186610d9e565b6110a36020830185611063565b6110b06040830184611072565b949350505050565b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b6000611114602e83610f7a565b915061111f826110b8565b604082019050919050565b6000602082019050818103600083015261114381611107565b9050919050565b6000819050919050565b6000819050919050565b600061117961117461116f8461114a565b611154565b610f43565b9050919050565b6111898161115e565b82525050565b60006020820190506111a46000830184611180565b92915050565b7f503a204f4e4c595f4641494c45445f414c4c4f57454400000000000000000000600082015250565b60006111e0601683610f7a565b91506111eb826111aa565b602082019050919050565b6000602082019050818103600083015261120f816111d3565b9050919050565b600081905092915050565b50565b6000611231600083611216565b915061123c82611221565b600082019050919050565b600061125282611224565b9150819050919050565b7f503a2053454e445f524556455254000000000000000000000000000000000000600082015250565b6000611292600e83610f7a565b915061129d8261125c565b602082019050919050565b600060208201905081810360008301526112c181611285565b9050919050565b60006020820190506112dd6000830184611063565b92915050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b600061133f602683610f7a565b915061134a826112e3565b604082019050919050565b6000602082019050818103600083015261136e81611332565b9050919050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b60006113ab602083610f7a565b91506113b682611375565b602082019050919050565b600060208201905081810360008301526113da8161139e565b9050919050565b7f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960008201527f6e697469616c697a696e67000000000000000000000000000000000000000000602082015250565b600061143d602b83610f7a565b9150611448826113e1565b604082019050919050565b6000602082019050818103600083015261146c81611430565b905091905056fea2646970667358221220606f457f98698781eb6281c0d93466c73143ef1c2cedcd26bdaad90f785cb73864736f6c63430008120033",
		Init:    "0x8129fc1c",
	}
	return ppp
}

func validateAddress(i interface{}) error {
	value, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !common.IsHexAddress(value) {
		return fmt.Errorf("address '%s' is not a hex", value)
	}
	return nil
}

func validateContracts(i interface{}) error {
	contracts, ok := i.(map[string]*ProxyContractInitState)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	keys := make([]string, 0)
	for k := range contracts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		state := contracts[k]
		if err := validateAddress(state.Address); err != nil {
			return err
		}
		if state.Height < 1 {
			return fmt.Errorf("invalid height '%d'", state.Height)
		}
		if _, err := hexutil.Decode(state.Bin); err != nil {
			return err
		}
		if _, err := hexutil.Decode(state.Init); len(state.Init) != 0 && err != nil {
			return err
		}
	}

	return nil
}
