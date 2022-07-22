package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params *Params, fileUploads []*FileUpload) GenesisState {
	return GenesisState{
		Params:      params,
		FileUploads: fileUploads,
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
	if err := data.Params.ValidateBasic(); err != nil {
		return err
	}

	if len(data.FileUploads) > 0 {
		for _, upload := range data.FileUploads {
			if len(upload.FileHash) == 0 {
				return ErrEmptyFileHash
			}
			if upload.FileInfo.Height.LT(sdk.ZeroInt()) {
				return ErrInvalidHeight
			}
			if len(upload.FileInfo.Reporter) == 0 {
				return ErrEmptyReporterAddr
			}
			if len(upload.FileInfo.Uploader) == 0 {
				return ErrEmptyUploaderAddr
			}
		}
	}
	return nil
}
