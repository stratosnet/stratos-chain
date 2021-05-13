package types

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"strings"
)

// Default parameter namespace
const (
	DefaultParamSpace = ModuleName
	DefaultBondDenom  = "stos"
)

// Parameter store keys
var (
	KeyBondDenom = []byte("BondDenom")
)

// ParamKeyTable for register module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for register at genesis
type Params struct {
	BondDenom string `json:"bond_denom" yaml:"bond_denom"` // bondable coin denomination
}

// NewParams creates a new Params object
func NewParams(bondDenom string) Params {
	return Params{
		BondDenom: bondDenom,
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`
	// TODO: Return all the params as a string
	`)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(DefaultBondDenom)
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
