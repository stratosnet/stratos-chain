package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// GenesisState - all register state that must be provided at genesis
type GenesisState struct {
	Params                     Params                  `json:"params" yaml:"params"`
	LastResourceNodeTotalStake sdk.Int                 `json:"last_resource_node_total_stake" yaml:"last_resource_node_total_stake"`
	LastResourceNodeStakes     []LastResourceNodeStake `json:"last_resource_node_stakes" yaml:"last_resource_node_stakes"`
	ResourceNodes              ResourceNodes           `json:"resource_nodes" yaml:"resource_nodes"`
	LastIndexingNodeTotalStake sdk.Int                 `json:"last_indexing_node_total_stake" yaml:"last_indexing_node_total_stake"`
	LastIndexingNodeStakes     []LastIndexingNodeStake `json:"last_indexing_node_stakes" yaml:"last_indexing_node_stakes"`
	IndexingNodes              IndexingNodes           `json:"indexing_nodes" yaml:"indexing_nodes"`
}

// LastResourceNodeStake required for resource node set update logic
type LastResourceNodeStake struct {
	Address sdk.AccAddress
	Stake   sdk.Int
}

// LastIndexingNodeStake required for indexing node set update logic
type LastIndexingNodeStake struct {
	Address sdk.AccAddress
	Stake   sdk.Int
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params,
	lastResourceNodeTotalStake sdk.Int, lastResourceNodeStakes []LastResourceNodeStake, resourceNodes []ResourceNode,
	lastIndexingNodeTotalStake sdk.Int, lastIndexingNodeStakes []LastIndexingNodeStake, indexingNodes []IndexingNode,
) GenesisState {
	return GenesisState{
		Params:                     params,
		LastResourceNodeTotalStake: lastResourceNodeTotalStake,
		LastResourceNodeStakes:     lastResourceNodeStakes,
		ResourceNodes:              resourceNodes,
		LastIndexingNodeTotalStake: lastIndexingNodeTotalStake,
		LastIndexingNodeStakes:     lastIndexingNodeStakes,
		IndexingNodes:              indexingNodes,
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
