package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// GenesisState - all register state that must be provided at genesis
type GenesisState struct {
	Params                     Params                  `json:"params" yaml:"params"`
	LastResourceNodeTotalPower sdk.Int                 `json:"last_resource_node_total_power" yaml:"last_resource_node_total_power"`
	LastResourceNodePowers     []LastResourceNodePower `json:"last_resource_node_powers" yaml:"last_resource_node_powers"`
	ResourceNodes              ResourceNodes           `json:"resource_nodes" yaml:"resource_nodes"`
	LastIndexingNodeTotalPower sdk.Int                 `json:"last_indexing_node_total_power" yaml:"last_indexing_node_total_power"`
	LastIndexingNodePowers     []LastIndexingNodePower `json:"last_indexing_node_powers" yaml:"last_indexing_node_powers"`
	IndexingNodes              IndexingNodes           `json:"indexing_nodes" yaml:"indexing_nodes"`
	Exported                   bool                    `json:"exported" yaml:"exported"`
}

// LastResourceNodePower required for resource node set update logic
type LastResourceNodePower struct {
	Address sdk.AccAddress
	Power   int64
}

// LastIndexingNodePower required for indexing node set update logic
type LastIndexingNodePower struct {
	Address sdk.AccAddress
	Power   int64
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, resourceNodes []ResourceNode, indexingNodes []IndexingNode) GenesisState {
	return GenesisState{
		Params:        params,
		ResourceNodes: resourceNodes,
		IndexingNodes: indexingNodes,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the register genesis parameters
func ValidateGenesis(data GenesisState) error {
	// TODO: Create a sanity check to make sure the state conforms to the modules needs
	return nil
}
