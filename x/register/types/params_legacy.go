package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyBondDenom               = []byte("BondDenom")
	KeyUnbondingThreasholdTime = []byte("UnbondingThreasholdTime")
	KeyUnbondingCompletionTime = []byte("UnbondingCompletionTime")
	KeyMaxEntries              = []byte("MaxEntries")
	KeyResourceNodeRegEnabled  = []byte("ResourceNodeRegEnabled")
	KeyResourceNodeMinDeposit  = []byte("ResourceNodeMinDeposit")
	KeyVotingPeriod            = []byte("VotingPeriod")
)

// ParamKeyTable returns the parameter key table.
// Deprecated: now params can be accessed on key `0x20` on the register store.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs - Implements params.ParamSet
// Deprecated.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
		paramtypes.NewParamSetPair(KeyUnbondingThreasholdTime, &p.UnbondingThreasholdTime, validateUnbondingThreasholdTime),
		paramtypes.NewParamSetPair(KeyUnbondingCompletionTime, &p.UnbondingCompletionTime, validateUnbondingCompletionTime),
		paramtypes.NewParamSetPair(KeyMaxEntries, &p.MaxEntries, validateMaxEntries),
		paramtypes.NewParamSetPair(KeyResourceNodeRegEnabled, &p.ResourceNodeRegEnabled, validateResourceNodeRegEnabled),
		paramtypes.NewParamSetPair(KeyResourceNodeMinDeposit, &p.ResourceNodeMinDeposit, validateResourceNodeMinDeposit),
		paramtypes.NewParamSetPair(KeyVotingPeriod, &p.VotingPeriod, validateVotingPeriod),
	}
}
