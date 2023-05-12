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
	DefaultBondDenom              = types.Wei
	DefaultMaxEntries             = uint32(16)
	DefaultResourceNodeRegEnabled = true
)

// Parameter store keys
var (
	KeyBondDenom               = []byte("BondDenom")
	KeyUnbondingThreasholdTime = []byte("UnbondingThreasholdTime")
	KeyUnbondingCompletionTime = []byte("UnbondingCompletionTime")
	KeyMaxEntries              = []byte("MaxEntries")
	KeyResourceNodeRegEnabled  = []byte("ResourceNodeRegEnabled")
	KeyResourceNodeMinStake    = []byte("ResourceNodeMinStake")
	KeyVotingValidityPeriod    = []byte("VotingValidityPeriod")

	DefaultUnbondingThreasholdTime = 180 * 24 * time.Hour // threashold for unbonding - by default 180 days
	DefaultUnbondingCompletionTime = 14 * 24 * time.Hour  // lead time to complete unbonding - by default 14 days
	DefaultVotingValidityPeriod    = 7 * 24 * time.Hour   // expiration time of registration voting - by default 7 days
	DefaultStakeNozRate            = sdk.NewDec(1000000)  // 0.001gwei -> 1noz = 1000000wei -> 1noz
	DefaultRemainingNozLimit       = sdk.NewInt(0)
	DefaultResourceNodeMinStake    = sdk.NewCoin(DefaultBondDenom, sdk.NewInt(1e18))
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(bondDenom string, threashold, completion time.Duration, maxEntries uint32,
	resourceNodeRegEnabled bool, resourceNodeMinStake sdk.Coin, votingValidityPeriod time.Duration) Params {

	return Params{
		BondDenom:               bondDenom,
		UnbondingThreasholdTime: threashold,
		UnbondingCompletionTime: completion,
		MaxEntries:              maxEntries,
		ResourceNodeRegEnabled:  resourceNodeRegEnabled,
		ResourceNodeMinStake:    resourceNodeMinStake,
		VotingValidityPeriod:    votingValidityPeriod,
	}
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
		paramtypes.NewParamSetPair(KeyUnbondingThreasholdTime, &p.UnbondingThreasholdTime, validateUnbondingThreasholdTime),
		paramtypes.NewParamSetPair(KeyUnbondingCompletionTime, &p.UnbondingCompletionTime, validateUnbondingCompletionTime),
		paramtypes.NewParamSetPair(KeyMaxEntries, &p.MaxEntries, validateMaxEntries),
		paramtypes.NewParamSetPair(KeyResourceNodeRegEnabled, &p.ResourceNodeRegEnabled, validateResourceNodeRegEnabled),
		paramtypes.NewParamSetPair(KeyResourceNodeMinStake, &p.ResourceNodeMinStake, validateResourceNodeMinStake),
		paramtypes.NewParamSetPair(KeyVotingValidityPeriod, &p.VotingValidityPeriod, validateVotingValidityPeriod),
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
	if err := validateResourceNodeRegEnabled(p.ResourceNodeRegEnabled); err != nil {
		return err
	}
	if err := validateResourceNodeMinStake(p.ResourceNodeMinStake); err != nil {
		return err
	}
	if err := validateVotingValidityPeriod(p.VotingValidityPeriod); err != nil {
		return err
	}
	return nil
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(
		DefaultBondDenom,
		DefaultUnbondingThreasholdTime,
		DefaultUnbondingCompletionTime,
		DefaultMaxEntries,
		DefaultResourceNodeRegEnabled,
		DefaultResourceNodeMinStake,
		DefaultVotingValidityPeriod,
	)
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

func validateResourceNodeMinStake(i interface{}) error {
	_, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateVotingValidityPeriod(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("voting validity period must be positive: %d", v)
	}

	return nil
}
