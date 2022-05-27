package types

import (
	"encoding/json"
	"fmt"

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

func (v GenesisIndexingNode) ToIndexingNode() IndexingNode {

	//fmt.Printf("v.GetPubkey().Value: %v, \r\n", v.GetPubkey().Value)
	//pubkey, ok := v.GetPubkey().GetCachedValue().(cryptotypes.PubKey)
	//
	//if !ok {
	//	fmt.Printf("pubkey: %v, \r\n", pubkey)
	//}

	//stStr, err := stratos.SdsPubKeyFromBech32("stsdspub1zcjduepqzgqd566qdnj4kna050jz505vamjhglxcpdepqctkregdt6snxm6spxdk2l")
	stPubkey, err := stratos.SdsPubKeyFromBech32("stsdspub1zcjduepqzgqd566qdnj4kna050jz505vamjhglxcpdepqctkregdt6snxm6spxdk2l")
	//fmt.Printf("stPubkey: %v\r\n", stPubkey.Bytes())

	//stStr, err := stratos.SdsPubkeyToBech32(pubkey)
	any, err := codectypes.NewAnyWithValue(stPubkey)
	fmt.Printf("any: %v, \r\n", any.Value)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("stStr: %s\r\n", stStr)
	//fmt.Printf("pubkey.String(): %v\r\n", pubkey.Bytes())
	//if err != nil {
	//	panic(err)
	//}

	ownerAddress, err := sdk.AccAddressFromBech32(v.OwnerAddress)
	if err != nil {
		panic(err)
	}

	fmt.Printf("GetNetworkAddress: %s\r\n", v.GetNetworkAddress())
	netAddr, err := stratos.SdsAddressFromBech32(v.GetNetworkAddress())
	fmt.Printf("netAddr: %s\r\n", netAddr)
	if err != nil {
		panic(err)
	}

	return IndexingNode{
		NetworkAddress: netAddr.String(),
		Pubkey:         any,
		Suspend:        v.GetSuspend(),
		Status:         v.GetStatus(),
		Tokens:         v.Tokens,
		OwnerAddress:   ownerAddress.String(),
		Description:    v.GetDescription(),
	}
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
	for i := range g.ResourceNodes.ResourceNodes {
		if err := g.ResourceNodes.ResourceNodes[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g GenesisIndexingNode) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(g.Pubkey, &pk)
}
