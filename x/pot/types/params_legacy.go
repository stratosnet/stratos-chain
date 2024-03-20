package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyBondDenom          = []byte("BondDenom")
	KeyRewardDenom        = []byte("RewardDenom")
	KeyMatureEpoch        = []byte("MatureEpoch")
	KeyMiningRewardParams = []byte("MiningRewardParams")
	KeyCommunityTax       = []byte("CommunityTax")
	KeyInitialTotalSupply = []byte("InitialTotalSupply")
)

// ParamKeyTable for pot module
// Deprecated: now params can be accessed on key `0x20` on the pot store.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs - Implements params.ParamSet
// Deprecated.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
		paramtypes.NewParamSetPair(KeyRewardDenom, &p.RewardDenom, validateRewardDenom),
		paramtypes.NewParamSetPair(KeyMatureEpoch, &p.MatureEpoch, validateMatureEpoch),
		paramtypes.NewParamSetPair(KeyMiningRewardParams, &p.MiningRewardParams, validateMiningRewardParams),
		paramtypes.NewParamSetPair(KeyCommunityTax, &p.CommunityTax, validateCommunityTax),
		paramtypes.NewParamSetPair(KeyInitialTotalSupply, &p.InitialTotalSupply, validateInitialTotalSupply),
	}
}
