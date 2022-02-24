package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all sds state that must be provided at genesis
type GenesisState struct {
	Params     Params       `json:"params" yaml:"params"`
	FileUpload []FileUpload `json:"file_upload" yaml:"file_upload"`
	Prepay     []Prepay     `json:"prepay" yaml:"prepay"`
}

// FileUpload required for fileInfo set update logic
type FileUpload struct {
	FileHash string   `json:"file_hash" yaml:"file_hash"`
	FileInfo FileInfo `json:"file_info" yaml:"file_info"`
}

// Prepay required for prepay set update logic
type Prepay struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	Coins  sdk.Coins      `json:"coins" yaml:"coins"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, fileUpload []FileUpload, prepay []Prepay) GenesisState {
	return GenesisState{
		Params:     params,
		FileUpload: fileUpload,
		Prepay:     prepay,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// GetGenesisStateFromAppState returns x/auth GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
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

	if len(data.FileUpload) > 0 {
		for _, upload := range data.FileUpload {
			if len(upload.FileHash) == 0 {
				return ErrEmptyFileHash
			}
			if upload.FileInfo.Height.LT(sdk.ZeroInt()) {
				return ErrInvalidHeight
			}
			if upload.FileInfo.Reporter.Empty() {
				return ErrEmptyReporterAddr
			}
			if upload.FileInfo.Uploader.Empty() {
				return ErrEmptyUploaderAddr
			}
		}
	}

	if len(data.Prepay) > 0 {
		for _, p := range data.Prepay {
			if p.Sender.Empty() {
				return ErrEmptySenderAddr
			}
			if !p.Coins.IsValid() {
				return ErrInvalidHeight
			}
		}
	}
	return nil
}
