package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	stratos "github.com/stratosnet/stratos-chain/types"
)

// DefaultParamSpace Default parameter namespace
const (
	DefaultBondDenom   = stratos.Wei
	DefaultRewardDenom = stratos.Utros
	DefaultMatureEpoch = 2016
)

// Parameter store keys
var (
	KeyBondDenom          = []byte("BondDenom")
	KeyRewardDenom        = []byte("RewardDenom")
	KeyMatureEpoch        = []byte("MatureEpoch")
	KeyMiningRewardParams = []byte("MiningRewardParams")
	KeyCommunityTax       = []byte("CommunityTax")
	KeyInitialTotalSupply = []byte("InitialTotalSupply")

	DefaultCommunityTax       = sdk.NewDecWithPrec(2, 2) // 2%
	DefaultInitialTotalSupply = sdk.NewCoin(DefaultBondDenom,
		sdk.NewInt(1e8).Mul(sdk.NewInt(stratos.StosToWei)),
	) //100,000,000 stos
)

// ParamKeyTable for pot module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(bondDenom string, rewardDenom string, matureEpoch int64, miningRewardParams []MiningRewardParam,
	communityTax sdk.Dec, initialTotalSupply sdk.Coin) Params {

	return Params{
		BondDenom:          bondDenom,
		RewardDenom:        rewardDenom,
		MatureEpoch:        matureEpoch,
		MiningRewardParams: miningRewardParams,
		CommunityTax:       communityTax,
		InitialTotalSupply: initialTotalSupply,
	}
}

// DefaultParams returns the default distribution parameters
func DefaultParams() Params {
	var miningRewardParams []MiningRewardParam
	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(0)),
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(16819200000000000)),
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(80000000000)),
		sdk.NewInt(6000), sdk.NewInt(2000), sdk.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(16819200000000000)),
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(25228800000000000)),
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(40000000000)),
		sdk.NewInt(6200), sdk.NewInt(1800), sdk.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(25228800000000000)),
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(29433600000000000)),
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(20000000000)),
		sdk.NewInt(6400), sdk.NewInt(1600), sdk.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(29433600000000000)),
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(31536000000000000)),
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(10000000000)),
		sdk.NewInt(6600), sdk.NewInt(1400), sdk.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(31536000000000000)),
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(32587200000000000)),
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(5000000000)),
		sdk.NewInt(6800), sdk.NewInt(1200), sdk.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(32587200000000000)),
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(40000000000000000)),
		sdk.NewCoin(DefaultRewardDenom, sdk.NewInt(2500000000)),
		sdk.NewInt(7000), sdk.NewInt(1000), sdk.NewInt(2000)))

	return NewParams(
		DefaultBondDenom,
		DefaultRewardDenom,
		DefaultMatureEpoch,
		miningRewardParams,
		DefaultCommunityTax,
		DefaultInitialTotalSupply,
	)
}

// HrpString implements the stringer interface for Params
func (p Params) HrpString() string {
	return fmt.Sprintf(`Params:
    BondDenom:          %s
    RewardDenom:        %s
    MatureEpoch:        %d
    MiningRewardParams: %s
    CommunitiyTax:      %v
    InitialTotalSupply: %v`,
		p.BondDenom, p.RewardDenom, p.MatureEpoch, p.MiningRewardParams, p.CommunityTax, p.InitialTotalSupply)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
		paramtypes.NewParamSetPair(KeyRewardDenom, &p.RewardDenom, validateRewardDenom),
		paramtypes.NewParamSetPair(KeyMatureEpoch, &p.MatureEpoch, validateMatureEpoch),
		paramtypes.NewParamSetPair(KeyMiningRewardParams, &p.MiningRewardParams, validateMiningRewardParams),
		paramtypes.NewParamSetPair(KeyCommunityTax, &p.CommunityTax, validateCommunityTax),
		paramtypes.NewParamSetPair(KeyInitialTotalSupply, &p.InitialTotalSupply, validateInitialTotalSupply),
	}
}

func validateBondDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("bond denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateRewardDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mining reward denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateMatureEpoch(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("mature epoch must be positive: %d", v)
	}

	return nil
}

func validateMiningRewardParams(i interface{}) error {
	v, ok := i.([]MiningRewardParam)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, p := range v {
		sumOfPercentage := p.BlockChainPercentageInBp.Int64() + p.MetaNodePercentageInBp.Int64() + p.ResourceNodePercentageInBp.Int64()
		if sumOfPercentage != 10000 {
			return fmt.Errorf("sum of block_chain_percentage_in_bp, resource_node_percentage_in_bp, meta_node_percentage_in_bp must be 10000: %v", sumOfPercentage)
		}
	}
	return nil
}

func validateCommunityTax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("pot community tax must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("pot community tax must be positive: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("pot community tax too large: %s", v)
	}

	return nil
}

func validateInitialTotalSupply(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("total supply must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("total supply must be positive: %s", v)
	}
	return nil
}

func (p Params) ValidateBasic() error {
	if err := validateBondDenom(p.BondDenom); err != nil {
		return sdkerrors.Wrap(ErrInvalidDenom, "failed to validate bond denomination")
	}
	if err := validateRewardDenom(p.RewardDenom); err != nil {
		return sdkerrors.Wrap(ErrInvalidDenom, "failed to validate reward denomination")
	}
	if err := validateMatureEpoch(p.MatureEpoch); err != nil {
		return sdkerrors.Wrap(ErrMatureEpoch, "failed to validate mature epoch")
	}
	if err := validateMiningRewardParams(p.MiningRewardParams); err != nil {
		return sdkerrors.Wrap(ErrMiningRewardParams, "failed to validate mining reward params")
	}
	if err := validateCommunityTax(p.CommunityTax); err != nil {
		return sdkerrors.Wrap(ErrCommunityTax, "failed to validate community tax")
	}
	if err := validateInitialTotalSupply(p.InitialTotalSupply); err != nil {
		return sdkerrors.Wrap(ErrInitialTotalSupply, err.Error())
	}
	return nil
}
