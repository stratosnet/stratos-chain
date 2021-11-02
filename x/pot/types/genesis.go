package types

type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
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
	return nil
}
