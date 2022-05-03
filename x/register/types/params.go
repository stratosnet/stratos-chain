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
	DefaultParamSpace = ModuleName
	DefaultBondDenom  = types.USTOS
	DefaultMaxEntries = uint32(16)
)

// Parameter store keys
var (
	KeyBondDenom               = []byte("BondDenom")
	KeyUnbondingThreasholdTime = []byte("UnbondingThreasholdTime")
	KeyUnbondingCompletionTime = []byte("UnbondingCompletionTime")
	KeyMaxEntries              = []byte("KeyMaxEntries")

	DefaultUnbondingThreasholdTime = (180 * 24 * time.Hour).String() // threashold for unbonding - by default 180 days
	DefaultUnbondingCompletionTime = (14 * 24 * time.Hour).String()  // lead time to complete unbonding - by default 14 days
	DefaultUozPrice                = sdk.NewDecWithPrec(1000000, 9)  // 0.001 ustos -> 1 uoz
	DefaultTotalUnissuedPrepay     = sdk.NewInt(0)
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(bondDenom string, threashold, completion string, maxEntries uint32) Params {
	return Params{
		BondDenom:               bondDenom,
		UnbondingThreasholdTime: threashold,
		UnbondingCompletionTime: completion,
		MaxEntries:              maxEntries,
	}
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
		paramtypes.NewParamSetPair(KeyUnbondingThreasholdTime, &p.UnbondingThreasholdTime, validateUnbondingThreasholdTime),
		paramtypes.NewParamSetPair(KeyUnbondingCompletionTime, &p.UnbondingCompletionTime, validateUnbondingCompletionTime),
		paramtypes.NewParamSetPair(KeyMaxEntries, &p.MaxEntries, validateMaxEntries),
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
func DefaultParams() *Params {
	p := NewParams(DefaultBondDenom, DefaultUnbondingThreasholdTime, DefaultUnbondingCompletionTime, DefaultMaxEntries)
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
