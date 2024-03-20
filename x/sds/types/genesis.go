package types

import (
	"encoding/json"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, files []GenesisFileInfo) *GenesisState {
	return &GenesisState{
		Params: params,
		Files:  files,
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

	if len(data.GetFiles()) > 0 {
		for _, file := range data.GetFiles() {
			if len(file.FileHash) == 0 {
				return ErrEmptyFileHash
			}
			if file.FileInfo.Height.LT(sdkmath.ZeroInt()) {
				return ErrInvalidHeight
			}
			if len(file.FileInfo.Reporters) == 0 {
				return ErrEmptyReporters
			}
			if len(file.FileInfo.Uploader) == 0 {
				return ErrEmptyUploaderAddr
			}
		}
	}
	return nil
}
