package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type IndexingNode struct {
	NetworkAddress string         `json:"network_address" yaml:"network_address"` // network address of the indexing node
	PubKey         crypto.PubKey  `json:"pubkey" yaml:"pubkey"`                   // the consensus public key of the indexing node; bech encoded in JSON
	Suspend        bool           `json:"suspend" yaml:"suspend"`                 // has the indexing node been suspended from bonded status?
	Status         sdk.BondStatus `json:"status" yaml:"status"`                   // indexing node status (bonded/unbonding/unbonded)
	Tokens         sdk.Int        `json:"tokens" yaml:"tokens"`                   // delegated tokens
	OwnerAddress   sdk.AccAddress `json:"owner_address" yaml:"owner_address"`     // owner address of the indexing node
	Description    Description    `json:"description" yaml:"description"`         // description terms for the indexing node
}

// NewIndexingNode - initialize a new indexing node
func NewIndexingNode(networkAddr string, pubKey crypto.PubKey, ownerAddr sdk.AccAddress, description Description) IndexingNode {
	return IndexingNode{
		NetworkAddress: networkAddr,
		PubKey:         pubKey,
		Suspend:        false,
		Status:         sdk.Unbonded,
		Tokens:         sdk.ZeroInt(),
		OwnerAddress:   ownerAddr,
		Description:    description,
	}
}

// MustMarshalIndexingNode returns the indexingNode bytes. Panics if fails
func MustMarshalIndexingNode(cdc *codec.Codec, indexingNode IndexingNode) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(indexingNode)
}

// MustUnmarshalIndexingNode unmarshal a indexing node from a store value. Panics if fails
func MustUnmarshalIndexingNode(cdc *codec.Codec, value []byte) IndexingNode {
	indexingNode, err := UnmarshalIndexingNode(cdc, value)
	if err != nil {
		panic(err)
	}
	return indexingNode
}

// UnmarshalIndexingNode unmarshal a indexing node from a store value
func UnmarshalIndexingNode(cdc *codec.Codec, value []byte) (indexingNode IndexingNode, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &indexingNode)
	return indexingNode, err
}

// AddToken adds tokens to a indexing node
func (v IndexingNode) AddToken(amount sdk.Int) IndexingNode {
	v.Tokens = v.Tokens.Add(amount)
	return v
}

// RemoveToken removes tokens from a indexing node
func (v IndexingNode) RemoveToken(tokens sdk.Int) IndexingNode {
	if tokens.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", tokens))
	}
	if v.Tokens.LT(tokens) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.Tokens, tokens))
	}
	v.Tokens = v.Tokens.Sub(tokens)
	return v
}

func (v IndexingNode) IsSuspended() bool            { return v.Suspend }
func (v IndexingNode) GetMoniker() string           { return v.Description.Moniker }
func (v IndexingNode) GetStatus() sdk.BondStatus    { return v.Status }
func (v IndexingNode) GetNetworkAddr() string       { return v.NetworkAddress }
func (v IndexingNode) GetPubKey() crypto.PubKey     { return v.PubKey }
func (v IndexingNode) GetAddr() sdk.AccAddress      { return sdk.AccAddress(v.PubKey.Address()) }
func (v IndexingNode) GetTokens() sdk.Int           { return v.Tokens }
func (v IndexingNode) GetOwnerAddr() sdk.AccAddress { return v.OwnerAddress }
