package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// ResourceNodeI expected resource node functions
type ResourceNodeI interface {
	IsSuspended() bool            // whether the node is jailed
	GetMoniker() string           // moniker of the node
	GetStatus() sdk.BondStatus    // status of the node
	GetNetworkAddr() string       // network address of the node
	GetPubKey() crypto.PubKey     // pubkey of the node
	GetAddr() sdk.AccAddress      // address of the node
	GetTokens() sdk.Int           // staking tokens of the node
	GetOwnerAddr() sdk.AccAddress // owner address of the node
	GetNodeType() []string        // node type
}

// IndexingNodeI expected indexing node functions
type IndexingNodeI interface {
	IsSuspended() bool            // whether the node is jailed
	GetMoniker() string           // moniker of the node
	GetStatus() sdk.BondStatus    // status of the node
	GetNetworkAddr() string       // network address of the node
	GetPubKey() crypto.PubKey     // pubkey of the node
	GetAddr() sdk.AccAddress      // address of the node
	GetTokens() sdk.Int           // staking tokens of the node
	GetOwnerAddr() sdk.AccAddress // owner address of the node
}
