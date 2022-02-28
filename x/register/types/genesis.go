package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

// GenesisState - all register state that must be provided at genesis
type GenesisState struct {
	Params                 Params                  `json:"params" yaml:"params"`
	LastResourceNodeStakes []LastResourceNodeStake `json:"last_resource_node_stakes" yaml:"last_resource_node_stakes"`
	ResourceNodes          ResourceNodes           `json:"resource_nodes" yaml:"resource_nodes"`
	LastIndexingNodeStakes []LastIndexingNodeStake `json:"last_indexing_node_stakes" yaml:"last_indexing_node_stakes"`
	IndexingNodes          IndexingNodes           `json:"indexing_nodes" yaml:"indexing_nodes"`
	InitialUozPrice        sdk.Dec                 `json:"initial_uoz_price" yaml:"initial_uoz_price"` //initial price of uoz
	TotalUnissuedPrepay    sdk.Int                 `json:"total_unissued_prepay" yaml:"total_unissued_prepay"`
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
	initialUOzonePrice sdk.Dec, totalUnissuedPrepay sdk.Int,
) GenesisState {
	return GenesisState{
		Params:                 params,
		LastResourceNodeStakes: lastResourceNodeStakes,
		ResourceNodes:          resourceNodes,
		LastIndexingNodeStakes: lastIndexingNodeStakes,
		IndexingNodes:          indexingNodes,
		InitialUozPrice:        initialUOzonePrice,
		TotalUnissuedPrepay:    totalUnissuedPrepay,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:              DefaultParams(),
		InitialUozPrice:     DefaultUozPrice,
		TotalUnissuedPrepay: DefaultTotalUnissuedPrepay,
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
	if err := data.Params.Validate(); err != nil {
		return err
	}
	if err := data.ResourceNodes.Validate(); err != nil {
		return err
	}
	if err := data.IndexingNodes.Validate(); err != nil {
		return err
	}

	if data.LastResourceNodeStakes != nil {
		for _, nodeStake := range data.LastResourceNodeStakes {
			if nodeStake.Address.Empty() {
				return ErrEmptyNetworkAddr
			}
			if nodeStake.Stake.LT(sdk.ZeroInt()) {
				return ErrValueNegative
			}
		}
	}

	if data.LastIndexingNodeStakes != nil {
		for _, nodeStake := range data.LastIndexingNodeStakes {
			if nodeStake.Address.Empty() {
				return ErrEmptyNetworkAddr
			}
			if nodeStake.Stake.LT(sdk.ZeroInt()) {
				return ErrValueNegative
			}
		}
	}
	if data.InitialUozPrice.LTE(sdk.ZeroDec()) {
		return ErrInitialUOzonePrice
	}

	if data.TotalUnissuedPrepay.LT(sdk.ZeroInt()) {
		return ErrInitialUOzonePrice
	}
	return nil
}

type GenesisIndexingNode struct {
	NetworkID    string         `json:"network_id" yaml:"network_id"`       // network address of the indexing node
	PubKey       string         `json:"pubkey" yaml:"pubkey"`               // the consensus public key of the indexing node; bech encoded in JSON
	Suspend      bool           `json:"suspend" yaml:"suspend"`             // has the indexing node been suspended from bonded status?
	Status       sdk.BondStatus `json:"status" yaml:"status"`               // indexing node status (bonded/unbonding/unbonded)
	Tokens       string         `json:"tokens" yaml:"tokens"`               // delegated tokens
	OwnerAddress string         `json:"owner_address" yaml:"owner_address"` // owner address of the indexing node
	Description  Description    `json:"description" yaml:"description"`     // description terms for the indexing node
}

func (v GenesisIndexingNode) ToIndexingNode() IndexingNode {
	pubKey, err := stratos.GetPubKeyFromBech32(stratos.Bech32PubKeyTypeSdsP2PPub, v.PubKey)
	if err != nil {
		panic(err)
	}

	tokens, ok := sdk.NewIntFromString(v.Tokens)
	if !ok {
		panic(ErrInvalidGenesisToken)
	}

	ownerAddress, err := sdk.AccAddressFromBech32(v.OwnerAddress)
	if err != nil {
		panic(err)
	}

	return IndexingNode{
		NetworkID:    v.NetworkID,
		PubKey:       pubKey,
		Suspend:      v.Suspend,
		Status:       v.Status,
		Tokens:       tokens,
		OwnerAddress: ownerAddress,
		Description:  v.Description,
	}
}
