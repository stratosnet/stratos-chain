package types

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all register state that must be provided at genesis
type GenesisState struct {
	Params                 Params                  `json:"params" yaml:"params"`
	LastResourceNodeStakes []LastResourceNodeStake `json:"last_resource_node_stakes" yaml:"last_resource_node_stakes"`
	ResourceNodes          ResourceNodes           `json:"resource_nodes" yaml:"resource_nodes"`
	LastIndexingNodeStakes []LastIndexingNodeStake `json:"last_indexing_node_stakes" yaml:"last_indexing_node_stakes"`
	IndexingNodes          IndexingNodes           `json:"indexing_nodes" yaml:"indexing_nodes"`
}

// LastResourceNodeStake required for resource node set update logic
type LastResourceNodeStake struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Stake   sdk.Int        `json:"stake" yaml:"stake"`
}

// LastIndexingNodeStake required for indexing node set update logic
type LastIndexingNodeStake struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Stake   sdk.Int        `json:"stake" yaml:"stake"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params,
	lastResourceNodeStakes []LastResourceNodeStake, resourceNodes ResourceNodes,
	lastIndexingNodeStakes []LastIndexingNodeStake, indexingNodes IndexingNodes,
) GenesisState {
	return GenesisState{
		Params:                 params,
		LastResourceNodeStakes: lastResourceNodeStakes,
		ResourceNodes:          resourceNodes,
		LastIndexingNodeStakes: lastIndexingNodeStakes,
		IndexingNodes:          indexingNodes,
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

// ValidateGenesis validates the register genesis parameters
func ValidateGenesis(data GenesisState) error {
	// TODO: Create a sanity check to make sure the state conforms to the modules needs
	return nil
}
