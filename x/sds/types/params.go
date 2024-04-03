package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stratos "github.com/stratosnet/stratos-chain/types"
)

const (
	DefaultBondDenom = stratos.Wei
)

// NewParams creates a new Params object
func NewParams(bondDenom string) Params {
	return Params{
		BondDenom: bondDenom,
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	p := NewParams(DefaultBondDenom)
	return p
}

func (p Params) Validate() error {
	if err := validateBondDenom(p.BondDenom); err != nil {
		return errors.Wrap(ErrInvalidDenom, "failed to validate bond denomination")
	}
	return nil
}

func validateBondDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("bond denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}
