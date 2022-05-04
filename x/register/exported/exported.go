package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/tendermint/tendermint/crypto"
)

// ResourceNodeI expected resource node functions
type ResourceNodeI interface {
	IsSuspended() bool                  // whether the node is jailed
	GetMoniker() string                 // moniker of the node
	GetStatus() stakingtypes.BondStatus // status of the node
	GetPubKey() crypto.PubKey           // pubkey of the node
	GetNetworkAddr() stratos.SdsAddress // network address of the node
	GetTokens() sdk.Int                 // staking tokens of the node
	GetOwnerAddr() sdk.AccAddress       // owner address of the node
	GetNodeType() string                // node type
}

// IndexingNodeI expected indexing node functions
type IndexingNodeI interface {
	IsSuspended() bool                  // whether the node is jailed
	GetMoniker() string                 // moniker of the node
	GetStatus() stakingtypes.BondStatus // status of the node
	GetPubKey() crypto.PubKey           // pubkey of the node
	GetNetworkAddr() stratos.SdsAddress // network address of the node
	GetTokens() sdk.Int                 // staking tokens of the node
	GetOwnerAddr() sdk.AccAddress       // owner address of the node
}
