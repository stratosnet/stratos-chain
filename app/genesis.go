package app

import (
	"cosmossdk.io/simapp"
	"github.com/cosmos/cosmos-sdk/codec"
)

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) simapp.GenesisState {
	return ModuleBasics.DefaultGenesis(cdc)
}
