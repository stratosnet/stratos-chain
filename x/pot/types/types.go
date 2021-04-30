package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type SingleNodeVolume struct {
	NodeAddress sdk.AccAddress `json:"node_address" yaml:"node_address"`
	Volume      sdk.Int        `json:"node_volume" yaml:"node_volume"`
}

// NewSingleNodeVolume creates a new Msg<Action> instance
func NewSingleNodeVolume(nodeAddress sdk.AccAddress, volume sdk.Int,
) SingleNodeVolume {
	return SingleNodeVolume{
		NodeAddress: nodeAddress,
		Volume:      volume,
	}
}
