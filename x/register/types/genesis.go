package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params *Params,
	resourceNodes *ResourceNodes,
	indexingNodes *IndexingNodes,
	initialUOzonePrice sdk.Dec,
	totalUnissuedPrepay sdk.Int,
	slashingInfo []*Slashing,
) GenesisState {
	return GenesisState{
		Params:              params,
		ResourceNodes:       resourceNodes,
		IndexingNodes:       indexingNodes,
		InitialUozPrice:     initialUOzonePrice,
		TotalUnissuedPrepay: totalUnissuedPrepay,
		Slashing:            slashingInfo,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:              DefaultParams(),
		ResourceNodes:       &ResourceNodes{},
		IndexingNodes:       &IndexingNodes{},
		InitialUozPrice:     DefaultUozPrice,
		TotalUnissuedPrepay: DefaultTotalUnissuedPrepay,
		Slashing:            make([]*Slashing, 0),
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

// ValidateGenesis validates the register genesis parameters
func ValidateGenesis(data GenesisState) error {
	if err := data.GetParams().Validate(); err != nil {
		return err
	}
	if err := data.GetResourceNodes().Validate(); err != nil {
		return err
	}
	if err := data.GetIndexingNodes().Validate(); err != nil {
		return err
	}

	if (data.InitialUozPrice).LTE(sdk.ZeroDec()) {
		return ErrInitialUOzonePrice
	}

	if data.TotalUnissuedPrepay.LT(sdk.ZeroInt()) {
		return ErrInitialUOzonePrice
	}
	return nil
}

func (v GenesisIndexingNode) ToIndexingNode() (IndexingNode, error) {
	ownerAddress, err := sdk.AccAddressFromBech32(v.OwnerAddress)
	if err != nil {
		return IndexingNode{}, ErrInvalidOwnerAddr
	}

	netAddr, err := stratos.SdsAddressFromBech32(v.GetNetworkAddress())
	if err != nil {
		return IndexingNode{}, ErrInvalidNetworkAddr
	}

	return IndexingNode{
		NetworkAddress: netAddr.String(),
		Pubkey:         v.GetPubkey(),
		Suspend:        v.GetSuspend(),
		Status:         v.GetStatus(),
		Tokens:         v.Tokens,
		OwnerAddress:   ownerAddress.String(),
		Description:    v.GetDescription(),
	}, nil
}

func NewSlashing(walletAddress sdk.AccAddress, value sdk.Int) *Slashing {
	return &Slashing{
		WalletAddress: walletAddress.String(),
		Value:         value.Int64(),
	}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g GenesisState) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range g.IndexingNodes.IndexingNodes {
		if err := g.IndexingNodes.IndexingNodes[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	//for i := range g.ResourceNodes.ResourceNodes {
	//	if err := g.ResourceNodes.ResourceNodes[i].UnpackInterfaces(c); err != nil {
	//		return err
	//	}

	//}
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v GenesisIndexingNode) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(v.Pubkey, &pk)
}
