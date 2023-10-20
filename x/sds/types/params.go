package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

const (
	DefaultBondDenom = stratos.Wei
)

var (
	KeyBondDenom = []byte("BondDenom")
)

// ParamKeyTable for sds module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(bondDenom string) Params {
	return Params{
		BondDenom: bondDenom,
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() *Params {
	p := NewParams(DefaultBondDenom)
	return &p
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
	}
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

func (p Params) ValidateBasic() error {
	if err := validateBondDenom(p.BondDenom); err != nil {
		return errors.Wrap(ErrInvalidDenom, "failed to validate bond denomination")
	}
	return nil
}
