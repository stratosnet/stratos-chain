package types

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"strings"
)

// DefaultParamSpace Default parameter namespace
const (
	DefaultParamSpace = ModuleName
	DefaultBondDenom  = "ustos"
)

// Parameter store keys
var (
	KeyBondDenom = []byte("BondDenom")
)

var _ subspace.ParamSet = &Params{}

// Params - used for initializing default parameter for pot at genesis
type Params struct {
	BondDenom string `json:"bond_denom" yaml:"bond_denom"` // bondable coin denomination
}

// ParamKeyTable for pot module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(bondDenom string) Params {
	return Params{
		BondDenom: bondDenom,
	}
}

// DefaultParams returns the default distribution parameters
func DefaultParams() Params {
	return NewParams(DefaultBondDenom)
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`Params:
	BondDenom:			%s`,
		p.BondDenom)
}

// ParamSetPairs - Implements params.ParamSet
func (p Params) ParamSetPairs() subspace.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
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

func (p Params) ValidateBasic() error {
	if err := validateBondDenom(p.BondDenom); err != nil {
		return err
	}
	return nil
}
