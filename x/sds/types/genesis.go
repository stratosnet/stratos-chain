package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, root []byte, previousRoots []MerkleRoot) *GenesisState {
	return &GenesisState{
		Params:              params,
		MerkleRoot:          root,
		PreviousMerkleRoots: previousRoots,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// GetGenesisStateFromAppState returns x/auth GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return genesisState
}

// ValidateGenesis validates the sds genesis parameters
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	for _, root := range data.PreviousMerkleRoots {
		if len(root.Root) == 0 {
			return ErrEmptyMerkleRoot
		}
		if root.Height < 0 {
			return ErrInvalidHeight
		}
	}
	return nil
}
