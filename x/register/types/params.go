package types

import (
	"errors"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// Default parameter namespace
const (
	DefaultParamSpace                            = ModuleName
	DefaultBondDenom                             = "ustos"
	DefaultUnbondingThreasholdTime time.Duration = 180 * 24 * time.Hour // threashold for unbonding - by default 180 days
	DefaultUnbondingCompletionTime time.Duration = 14 * 24 * time.Hour  // lead time to complete unbonding - by default 14 days
	DefaultMaxEntries                            = uint16(16)
)

// Parameter store keys
var (
	KeyBondDenom               = []byte("BondDenom")
	KeyUnbondingThreasholdTime = []byte("UnbondingThreasholdTime")
	KeyUnbondingCompletionTime = []byte("UnbondingCompletionTime")
	KeyMaxEntries              = []byte("KeyMaxEntries")

	DefaultUozPrice = sdk.NewDecWithPrec(1000000, 9) // 0.001 ustos -> 1 uoz
)

var _ subspace.ParamSet = &Params{}

// ParamKeyTable for register module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for register at genesis
type Params struct {
	BondDenom               string        `json:"bond_denom" yaml:"bond_denom"`                               // bondable coin denomination
	UnbondingThreasholdTime time.Duration `json:"unbonding_threashold_time" yaml:"unbonding_threashold_time"` // threashold for unbonding - by default 180 days
	UnbondingCompletionTime time.Duration `json:"unbonding_completion_time" yaml:"unbonding_completion_time"` // lead time to complete unbonding - by default 14 days
	MaxEntries              uint16        `json:"max_entries" yaml:"max_entries"`                             // max entries for either unbonding delegation or redelegation (per pair/trio)
}

// NewParams creates a new Params object
func NewParams(bondDenom string, threashold, completion time.Duration, maxEntries uint16) Params {
	return Params{
		BondDenom:               bondDenom,
		UnbondingThreasholdTime: threashold,
		UnbondingCompletionTime: completion,
		MaxEntries:              maxEntries,
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`Register Params:
	  BondDenom:             		%s
	  Unbonding Threashold Time:  	%s
	  Unbonding Completion Time:  	%s
	  Max Entries:        			%d
`,
		p.BondDenom, p.UnbondingThreasholdTime, p.UnbondingCompletionTime, p.MaxEntries,
	)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
		params.NewParamSetPair(KeyUnbondingThreasholdTime, &p.UnbondingThreasholdTime, validateUnbondingThreasholdTime),
		params.NewParamSetPair(KeyUnbondingCompletionTime, &p.UnbondingCompletionTime, validateUnbondingCompletionTime),
		params.NewParamSetPair(KeyMaxEntries, &p.MaxEntries, validateMaxEntries),
	}
}

func (p Params) Validate() error {
	if err := validateBondDenom(p.BondDenom); err != nil {
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
	return nil
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(DefaultBondDenom, DefaultUnbondingThreasholdTime, DefaultUnbondingCompletionTime, DefaultMaxEntries)
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
	v, ok := i.(uint16)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("max entries must be positive: %d", v)
	}

	return nil
}
