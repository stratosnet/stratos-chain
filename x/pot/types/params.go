package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

const (
	DefaultBondDenom   = stratos.Wei
	DefaultRewardDenom = stratos.Wei
	DefaultMatureEpoch = 2016
)

var (
	DefaultCommunityTax       = sdkmath.LegacyNewDecWithPrec(2, 2) // 2%
	DefaultInitialTotalSupply = sdk.NewCoin(DefaultBondDenom,
		sdkmath.NewInt(1e8).MulRaw(stratos.StosToWei),
	) //100,000,000 stos
)

// NewParams creates a new Params object
func NewParams(bondDenom string, rewardDenom string, matureEpoch int64, miningRewardParams []MiningRewardParam,
	communityTax sdkmath.LegacyDec, initialTotalSupply sdk.Coin) Params {

	return Params{
		BondDenom:          bondDenom,
		RewardDenom:        rewardDenom,
		MatureEpoch:        matureEpoch,
		MiningRewardParams: miningRewardParams,
		CommunityTax:       communityTax,
		InitialTotalSupply: initialTotalSupply,
	}
}

// DefaultParams returns the default distribution parameters
func DefaultParams() Params {
	var miningRewardParams []MiningRewardParam
	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(0)),
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(16819200).MulRaw(stratos.StosToWei)),
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(80).MulRaw(stratos.StosToWei)),
		sdkmath.NewInt(6000), sdkmath.NewInt(2000), sdkmath.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(16819200).MulRaw(stratos.StosToWei)),
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(25228800).MulRaw(stratos.StosToWei)),
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(40).MulRaw(stratos.StosToWei)),
		sdkmath.NewInt(6200), sdkmath.NewInt(1800), sdkmath.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(25228800).MulRaw(stratos.StosToWei)),
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(29433600).MulRaw(stratos.StosToWei)),
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(20).MulRaw(stratos.StosToWei)),
		sdkmath.NewInt(6400), sdkmath.NewInt(1600), sdkmath.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(29433600).MulRaw(stratos.StosToWei)),
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(31536000).MulRaw(stratos.StosToWei)),
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(10).MulRaw(stratos.StosToWei)),
		sdkmath.NewInt(6600), sdkmath.NewInt(1400), sdkmath.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(31536000).MulRaw(stratos.StosToWei)),
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(32587200).MulRaw(stratos.StosToWei)),
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(5).MulRaw(stratos.StosToWei)),
		sdkmath.NewInt(6800), sdkmath.NewInt(1200), sdkmath.NewInt(2000)))

	miningRewardParams = append(miningRewardParams, NewMiningRewardParam(
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(32587200).MulRaw(stratos.StosToWei)),
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(40000000).MulRaw(stratos.StosToWei)),
		sdk.NewCoin(DefaultRewardDenom, sdkmath.NewInt(25).MulRaw(1e17)),
		sdkmath.NewInt(7000), sdkmath.NewInt(1000), sdkmath.NewInt(2000)))

	return NewParams(
		DefaultBondDenom,
		DefaultRewardDenom,
		DefaultMatureEpoch,
		miningRewardParams,
		DefaultCommunityTax,
		DefaultInitialTotalSupply,
	)
}

func (p Params) Validate() error {
	if err := validateBondDenom(p.BondDenom); err != nil {
		return errors.Wrap(ErrInvalidDenom, "failed to validate bond denomination")
	}
	if err := validateRewardDenom(p.RewardDenom); err != nil {
		return errors.Wrap(ErrInvalidDenom, "failed to validate reward denomination")
	}
	if err := validateMatureEpoch(p.MatureEpoch); err != nil {
		return errors.Wrap(ErrMatureEpoch, "failed to validate mature epoch")
	}
	if err := validateMiningRewardParams(p.MiningRewardParams); err != nil {
		return errors.Wrap(ErrMiningRewardParams, "failed to validate mining reward params")
	}
	if err := validateCommunityTax(p.CommunityTax); err != nil {
		return errors.Wrap(ErrCommunityTax, "failed to validate community tax")
	}
	if err := validateInitialTotalSupply(p.InitialTotalSupply); err != nil {
		return errors.Wrap(ErrInitialTotalSupply, err.Error())
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

func validateRewardDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("mining reward denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateMatureEpoch(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("mature epoch must be positive: %d", v)
	}

	return nil
}

func validateMiningRewardParams(i interface{}) error {
	v, ok := i.([]MiningRewardParam)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, p := range v {
		sumOfPercentage := p.BlockChainPercentageInBp.Int64() + p.MetaNodePercentageInBp.Int64() + p.ResourceNodePercentageInBp.Int64()
		if sumOfPercentage != 10000 {
			return fmt.Errorf("sum of block_chain_percentage_in_bp, resource_node_percentage_in_bp, meta_node_percentage_in_bp must be 10000: %v", sumOfPercentage)
		}
	}
	return nil
}

func validateCommunityTax(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("pot community tax must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("pot community tax must be positive: %s", v)
	}
	if v.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("pot community tax too large: %s", v)
	}

	return nil
}

func validateInitialTotalSupply(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("total supply must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("total supply must be positive: %s", v)
	}
	return nil
}
