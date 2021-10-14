package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params          Params  `json:"params" yaml:"params"`
	InitialUozPrice sdk.Dec `json:"initial_uoz_price" yaml:"initial_uoz_price"` //initial price of uoz
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, initialUOzonePrice sdk.Dec) GenesisState {
	return GenesisState{
		Params:          params,
		InitialUozPrice: initialUOzonePrice,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:          DefaultParams(),
		InitialUozPrice: DefaultUozPrice,
	}
}

// ValidateGenesis validates the pot genesis parameters
func ValidateGenesis(data GenesisState) error {
	if data.InitialUozPrice.LTE(sdk.ZeroDec()) {
		return ErrInitialUOzonePrice
	}
	return nil
}
