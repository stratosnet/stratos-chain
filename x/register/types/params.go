package types

import (
	"errors"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/stratosnet/stratos-chain/types"
)

var _ paramtypes.ParamSet = &Params{}

// Default parameter namespace
const (
	DefaultBondDenom              = types.USTOS
	DefaultRewardDenom            = "utros"
	DefaultMaxEntries             = uint32(16)
	DefaultResourceNodeRegEnabled = true
)

// Parameter store keys
var (
	KeyBondDenom               = []byte("BondDenom")
	KeyRewardDenom             = []byte("RewardDenom")
	KeyUnbondingThreasholdTime = []byte("UnbondingThreasholdTime")
	KeyUnbondingCompletionTime = []byte("UnbondingCompletionTime")
	KeyMaxEntries              = []byte("MaxEntries")
	KeyResourceNodeRegEnabled  = []byte("ResourceNodeRegEnabled")

	DefaultUnbondingThreasholdTime = 180 * 24 * time.Hour           // threashold for unbonding - by default 180 days
	DefaultUnbondingCompletionTime = 14 * 24 * time.Hour            // lead time to complete unbonding - by default 14 days
	DefaultUozPrice                = sdk.NewDecWithPrec(1000000, 9) // 0.001 ustos -> 1 uoz
	DefaultTotalUnissuedPrepay     = sdk.NewInt(0)
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(bondDenom, rewardDenom string, threashold, completion time.Duration,
	maxEntries uint32, resourceNodeRegEnabled bool) Params {

	return Params{
		BondDenom:               bondDenom,
		RewardDenom:             rewardDenom,
		UnbondingThreasholdTime: threashold,
		UnbondingCompletionTime: completion,
		MaxEntries:              maxEntries,
		ResourceNodeRegEnabled:  resourceNodeRegEnabled,
	}
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
		paramtypes.NewParamSetPair(KeyRewardDenom, &p.RewardDenom, validateRewardDenom),
		paramtypes.NewParamSetPair(KeyUnbondingThreasholdTime, &p.UnbondingThreasholdTime, validateUnbondingThreasholdTime),
		paramtypes.NewParamSetPair(KeyUnbondingCompletionTime, &p.UnbondingCompletionTime, validateUnbondingCompletionTime),
		paramtypes.NewParamSetPair(KeyMaxEntries, &p.MaxEntries, validateMaxEntries),
		paramtypes.NewParamSetPair(KeyResourceNodeRegEnabled, &p.ResourceNodeRegEnabled, validateResourceNodeRegEnabled),
	}
}

func (p Params) Validate() error {
	if err := validateBondDenom(p.BondDenom); err != nil {
		return err
	}
	if err := validateRewardDenom(p.RewardDenom); err != nil {
		return err
	}
	if err := validateUnbondingThreasholdTime(p.UnbondingThreasholdTime); err != nil {
		return err
	}
	if err := validateUnbondingCompletionTime(p.UnbondingCompletionTime); err != nil {
		return err
	}
	if err := validateMaxEntries(p.MaxEntries); err != nil {
		return err
	}
	if err := validateResourceNodeRegEnabled(p.ResourceNodeRegEnabled); err != nil {
		return err
	}
	return nil
}

// DefaultParams defines the parameters for this module
func DefaultParams() *Params {
	p := NewParams(
		DefaultBondDenom,
		DefaultRewardDenom,
		DefaultUnbondingThreasholdTime,
		DefaultUnbondingCompletionTime,
		DefaultMaxEntries,
		DefaultResourceNodeRegEnabled,
	)
	return &p
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

func validateUnbondingThreasholdTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("unbonding threashold time must be positive: %d", v)
	}

	return nil
}

func validateUnbondingCompletionTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("unbonding completion time must be positive: %d", v)
	}

	return nil
}

func validateMaxEntries(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("max entries must be positive: %d", v)
	}

	return nil
}

func validateResourceNodeRegEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
