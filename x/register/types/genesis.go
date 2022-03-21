package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

// GenesisState - all register state that must be provided at genesis
type GenesisState struct {
	Params              Params        `json:"params" yaml:"params"`
	ResourceNodes       ResourceNodes `json:"resource_nodes" yaml:"resource_nodes"`
	IndexingNodes       IndexingNodes `json:"indexing_nodes" yaml:"indexing_nodes"`
	InitialUozPrice     sdk.Dec       `json:"initial_uoz_price" yaml:"initial_uoz_price"` //initial price of uoz
	TotalUnissuedPrepay sdk.Int       `json:"total_unissued_prepay" yaml:"total_unissued_prepay"`
	SlashingInfo        []Slashing    `json:"slashing_info" yaml:"slashing_info"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params,
	resourceNodes ResourceNodes,
	indexingNodes IndexingNodes,
	initialUOzonePrice sdk.Dec,
	totalUnissuedPrepay sdk.Int,
	slashingInfo []Slashing,
) GenesisState {
	return GenesisState{
		Params:              params,
		ResourceNodes:       resourceNodes,
		IndexingNodes:       indexingNodes,
		InitialUozPrice:     initialUOzonePrice,
		TotalUnissuedPrepay: totalUnissuedPrepay,
		SlashingInfo:        slashingInfo,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:              DefaultParams(),
		InitialUozPrice:     DefaultUozPrice,
		TotalUnissuedPrepay: DefaultTotalUnissuedPrepay,
		SlashingInfo:        make([]Slashing, 0),
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

	if data.InitialUozPrice.LTE(sdk.ZeroDec()) {
		return ErrInitialUOzonePrice
	}

	if data.TotalUnissuedPrepay.LT(sdk.ZeroInt()) {
		return ErrInitialUOzonePrice
	}
	return nil
}

type GenesisIndexingNode struct {
	NetworkAddr  string         `json:"network_address" yaml:"network_address"` // network address of the indexing node
	PubKey       string         `json:"pubkey" yaml:"pubkey"`                   // the consensus public key of the indexing node; bech encoded in JSON
	Suspend      bool           `json:"suspend" yaml:"suspend"`                 // has the indexing node been suspended from bonded status?
	Status       sdk.BondStatus `json:"status" yaml:"status"`                   // indexing node status (bonded/unbonding/unbonded)
	Tokens       string         `json:"tokens" yaml:"tokens"`                   // delegated tokens
	OwnerAddress string         `json:"owner_address" yaml:"owner_address"`     // owner address of the indexing node
	Description  Description    `json:"description" yaml:"description"`         // description terms for the indexing node
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

	netAddr, err := stratos.SdsAddressFromBech32(v.NetworkAddr)
	if err != nil {
		panic(err)
	}

	return IndexingNode{
		NetworkAddr:  netAddr,
		PubKey:       pubKey,
		Suspend:      v.Suspend,
		Status:       v.Status,
		Tokens:       tokens,
		OwnerAddress: ownerAddress,
		Description:  v.Description,
	}
}

type Slashing struct {
	WalletAddress sdk.AccAddress
	Value         sdk.Int
}

func NewSlashing(walletAddress sdk.AccAddress, value sdk.Int) Slashing {
	return Slashing{
		WalletAddress: walletAddress,
		Value:         value,
	}
}
