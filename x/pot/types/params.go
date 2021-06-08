package types

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"strings"
)

// DefaultParamSpace Default parameter namespace
const (
	DefaultParamSpace = ModuleName
	DefaultBondDenom  = "ustos"
)

// Parameter store keys
var (
	KeyBondDenom          = []byte("BondDenom")
	KeyMiningRewardParams = []byte("MiningRewardParams")
)

// Params - used for initializing default parameter for pot at genesis
type Params struct {
	BondDenom          string              `json:"bond_denom" yaml:"bond_denom"` // bondable coin denomination
	MiningRewardParams []MiningRewardParam `json:"mining_reward_params" yaml:"mining_reward_params"`
}

// ParamKeyTable for pot module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(bondDenom string) Params {
	var miningRewardParams []MiningRewardParam
	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewInt(0), sdk.NewInt(16819200000000000), sdk.NewInt(80000000000),
		sdk.NewInt(6000), sdk.NewInt(2000), sdk.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewInt(16819200000000000), sdk.NewInt(25228800000000000), sdk.NewInt(40000000000),
		sdk.NewInt(6200), sdk.NewInt(1800), sdk.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewInt(25228800000000000), sdk.NewInt(29433600000000000), sdk.NewInt(20000000000),
		sdk.NewInt(6400), sdk.NewInt(1600), sdk.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewInt(29433600000000000), sdk.NewInt(31536000000000000), sdk.NewInt(10000000000),
		sdk.NewInt(6600), sdk.NewInt(1400), sdk.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewInt(31536000000000000), sdk.NewInt(32587200000000000), sdk.NewInt(5000000000),
		sdk.NewInt(6800), sdk.NewInt(1200), sdk.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewInt(32587200000000000), sdk.NewInt(40000000000000000), sdk.NewInt(2500000000),
		sdk.NewInt(7000), sdk.NewInt(1000), sdk.NewInt(2000)))

	return Params{
		BondDenom:          bondDenom,
		MiningRewardParams: miningRewardParams,
	}
}

// DefaultParams returns the default distribution parameters
func DefaultParams() Params {
	return NewParams(DefaultBondDenom)
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`Params:
	BondDenom:			%s
  	MiningRewardParams:	%s`,
		p.BondDenom, p.MiningRewardParams)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
		params.NewParamSetPair(KeyMiningRewardParams, &p.MiningRewardParams, validateMiningRewardParams),
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

func validateMiningRewardParams(i interface{}) error {
	return nil
}

func (p Params) ValidateBasic() error {
	if err := validateBondDenom(p.BondDenom); err != nil {
		return err
	}
	return nil
}
