package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params            Params         `json:"params" yaml:"params"`
	FoundationAccount sdk.AccAddress `json:"foundation_account" yaml:"foundation_account"` //foundation account address
	InitialUozPrice   sdk.Int        `json:"initial_uoz_price" yaml:"initial_uoz_price"`   //initial price of uoz
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, foundationAccount sdk.AccAddress, initialUOzonePrice sdk.Int) GenesisState {
	return GenesisState{
		Params:            params,
		FoundationAccount: foundationAccount,
		InitialUozPrice:   initialUOzonePrice,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:          DefaultParams(),
		InitialUozPrice: sdk.NewInt(10000000000),
	}
}

// ValidateGenesis validates the pot genesis parameters
func ValidateGenesis(data GenesisState) error {
	if data.FoundationAccount == nil {
		return ErrFoundationAccount
	}
	if data.InitialUozPrice.LTE(sdk.ZeroInt()) {
		return ErrInitialUOzonePrice
	}
	return nil
}
