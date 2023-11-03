package types

import (
	"errors"
	"fmt"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stratos "github.com/stratosnet/stratos-chain/types"
)

// Default parameter namespace
const (
	DefaultBondDenom              = stratos.Wei
	DefaultMaxEntries             = uint32(16)
	DefaultResourceNodeRegEnabled = true
)

// Parameter store keys
var (
	DefaultUnbondingThreasholdTime = 180 * 24 * time.Hour      // threashold for unbonding - by default 180 days
	DefaultUnbondingCompletionTime = 14 * 24 * time.Hour       // lead time to complete unbonding - by default 14 days
	DefaultVotingPeriod            = 7 * 24 * time.Hour        // expiration time of registration voting - by default 7 days
	DefaultDepositNozRate          = sdkmath.LegacyNewDec(1e6) // 0.001gwei -> 1noz = 1000000wei -> 1noz
	DefaultRemainingNozLimit       = sdkmath.ZeroInt()
	DefaultResourceNodeMinDeposit  = sdk.NewCoin(DefaultBondDenom, sdkmath.NewInt(stratos.StosToWei))
)

// NewParams creates a new Params object
func NewParams(bondDenom string, threashold, completion time.Duration, maxEntries uint32,
	resourceNodeRegEnabled bool, resourceNodeMinDeposit sdk.Coin, votingPeriod time.Duration) Params {

	return Params{
		BondDenom:               bondDenom,
		UnbondingThreasholdTime: threashold,
		UnbondingCompletionTime: completion,
		MaxEntries:              maxEntries,
		ResourceNodeRegEnabled:  resourceNodeRegEnabled,
		ResourceNodeMinDeposit:  resourceNodeMinDeposit,
		VotingPeriod:            votingPeriod,
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(
		DefaultBondDenom,
		DefaultUnbondingThreasholdTime,
		DefaultUnbondingCompletionTime,
		DefaultMaxEntries,
		DefaultResourceNodeRegEnabled,
		DefaultResourceNodeMinDeposit,
		DefaultVotingPeriod,
	)
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
	if err := validateResourceNodeRegEnabled(p.ResourceNodeRegEnabled); err != nil {
		return err
	}
	if err := validateResourceNodeMinDeposit(p.ResourceNodeMinDeposit); err != nil {
		return err
	}
	if err := validateVotingPeriod(p.VotingPeriod); err != nil {
		return err
	}
	return nil
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

func validateResourceNodeMinDeposit(i interface{}) error {
	_, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateVotingPeriod(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("voting period must be positive: %d", v)
	}

	return nil
}
