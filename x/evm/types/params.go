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
	ParamStoreKeyProxyOwnerAddress = []byte("ProxyOwnerAddress")
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
		paramtypes.NewParamSetPair(ParamStoreKeyProxyOwnerAddress, &p.ProxyProposalParams.ProxyOwnerAddress, validateAddress),
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
func NewProxyProposalParams(consensusAddress, proxyOwnerAddress common.Address) ProxyProposalParams {
	return ProxyProposalParams{
		ConsensusAddress:  consensusAddress.Hex(),
		ProxyOwnerAddress: proxyOwnerAddress.Hex(),
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
		Bin:     "0x608060405234801561001057600080fd5b50610b19806100206000396000f3fe6080604052600436106100595760003560e01c806334fe1d1e14610065578063715018a61461006f5780638129fc1c146100865780638da5cb5b1461009d578063f2fde38b146100c8578063ffa1ad74146100f157610060565b3661006057005b600080fd5b61006d61011c565b005b34801561007b57600080fd5b50610084610241565b005b34801561009257600080fd5b5061009b610255565b005b3480156100a957600080fd5b506100b261039b565b6040516100bf91906106d1565b60405180910390f35b3480156100d457600080fd5b506100ef60048036038101906100ea919061071d565b6103c5565b005b3480156100fd57600080fd5b50610106610448565b6040516101139190610766565b60405180910390f35b600034905060008103610164576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161015b906107de565b60405180910390fd5b61016c61066e565b3373ffffffffffffffffffffffffffffffffffffffff1681600060018110610197576101966107fe565b5b6020020181815250506101a861066e565b6020816020848660f1600019f16101be57600080fd5b823373ffffffffffffffffffffffffffffffffffffffff163073ffffffffffffffffffffffffffffffffffffffff167fa9fdf2e446d7225a2b445bc7c21ca59dcea69b5b23f5c4e6f54f87a5db6cdaee84600060018110610222576102216107fe565b5b60200201516040516102349190610846565b60405180910390a4505050565b61024961044d565b61025360006104cb565b565b60008060019054906101000a900460ff161590508080156102865750600160008054906101000a900460ff1660ff16105b806102b3575061029530610591565b1580156102b25750600160008054906101000a900460ff1660ff16145b5b6102f2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102e9906108d3565b60405180910390fd5b60016000806101000a81548160ff021916908360ff160217905550801561032f576001600060016101000a81548160ff0219169083151502179055505b6103376105b4565b61033f610605565b80156103985760008060016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498600160405161038f9190610938565b60405180910390a15b50565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6103cd61044d565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361043c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610433906109c5565b60405180910390fd5b610445816104cb565b50565b600081565b610455610666565b73ffffffffffffffffffffffffffffffffffffffff1661047361039b565b73ffffffffffffffffffffffffffffffffffffffff16146104c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104c090610a31565b60405180910390fd5b565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081603360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b600060019054906101000a900460ff16610603576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105fa90610ac3565b60405180910390fd5b565b600060019054906101000a900460ff16610654576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161064b90610ac3565b60405180910390fd5b61066461065f610666565b6104cb565b565b600033905090565b6040518060200160405280600190602082028036833780820191505090505090565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006106bb82610690565b9050919050565b6106cb816106b0565b82525050565b60006020820190506106e660008301846106c2565b92915050565b600080fd5b6106fa816106b0565b811461070557600080fd5b50565b600081359050610717816106f1565b92915050565b600060208284031215610733576107326106ec565b5b600061074184828501610708565b91505092915050565b600060ff82169050919050565b6107608161074a565b82525050565b600060208201905061077b6000830184610757565b92915050565b600082825260208201905092915050565b7f503a205a45524f5f414d4f554e54000000000000000000000000000000000000600082015250565b60006107c8600e83610781565b91506107d382610792565b602082019050919050565b600060208201905081810360008301526107f7816107bb565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000819050919050565b6108408161082d565b82525050565b600060208201905061085b6000830184610837565b92915050565b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b60006108bd602e83610781565b91506108c882610861565b604082019050919050565b600060208201905081810360008301526108ec816108b0565b9050919050565b6000819050919050565b6000819050919050565b600061092261091d610918846108f3565b6108fd565b61074a565b9050919050565b61093281610907565b82525050565b600060208201905061094d6000830184610929565b92915050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b60006109af602683610781565b91506109ba82610953565b604082019050919050565b600060208201905081810360008301526109de816109a2565b9050919050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b6000610a1b602083610781565b9150610a26826109e5565b602082019050919050565b60006020820190508181036000830152610a4a81610a0e565b9050919050565b7f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960008201527f6e697469616c697a696e67000000000000000000000000000000000000000000602082015250565b6000610aad602b83610781565b9150610ab882610a51565b604082019050919050565b60006020820190508181036000830152610adc81610aa0565b905091905056fea2646970667358221220af723f0b659da0cdaa0d447e54bb8c0c3af22886c04139bf5fc7ee74f487451c64736f6c63430008120033",
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
