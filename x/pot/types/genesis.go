package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params             Params         `json:"params" yaml:"params"`
	FoundationAccount  sdk.AccAddress `json:"foundation_account" yaml:"foundation_account"`       //foundation account address
	InitialUOzonePrice sdk.Int        `json:"initial_u_ozone_price" yaml:"initial_u_ozone_price"` //initial price of uOzone
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, foundationAccount sdk.AccAddress, initialUOzonePrice sdk.Int) GenesisState {
	return GenesisState{
		Params:             params,
		FoundationAccount:  foundationAccount,
		InitialUOzonePrice: initialUOzonePrice,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the pot genesis parameters
func ValidateGenesis(data GenesisState) error {
	if data.InitialUOzonePrice.LTE(sdk.ZeroInt()) {
		return ErrInitialUOzonePrice
	}
	return nil
}
