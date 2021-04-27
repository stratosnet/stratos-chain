package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type ResourceNode struct {
	NetworkAddress string         `json:"network_address" yaml:"network_address"` // network address of the resource node
	PubKey         crypto.PubKey  `json:"pubkey" yaml:"pubkey"`                   // the public key of the resource node; bech encoded in JSON
	Suspend        bool           `json:"suspend" yaml:"suspend"`                 // has the resource node been suspended from bonded status?
	Status         sdk.BondStatus `json:"status" yaml:"status"`                   // resource node bond status (bonded/unbonding/unbonded)
	Tokens         sdk.Int        `json:"tokens" yaml:"tokens"`                   // delegated tokens
	OwnerAddress   sdk.AccAddress `json:"owner_address" yaml:"owner_address"`     // owner address of the resource node
	Description    Description    `json:"description" yaml:"description"`         // description terms for the resource node
}

// NewResourceNode - initialize a new resource node
func NewResourceNode(networkAddr string, pubKey crypto.PubKey, ownerAddr sdk.AccAddress, description Description) ResourceNode {
	return ResourceNode{
		NetworkAddress: networkAddr,
		PubKey:         pubKey,
		Suspend:        false,
		Status:         sdk.Unbonded,
		Tokens:         sdk.ZeroInt(),
		OwnerAddress:   ownerAddr,
		Description:    description,
	}
}

// MustMarshalResourceNode returns the resourceNode bytes. Panics if fails
func MustMarshalResourceNode(cdc *codec.Codec, resourceNode ResourceNode) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(resourceNode)
}

// MustUnmarshalResourceNode unmarshal a resourceNode from a store value. Panics if fails
func MustUnmarshalResourceNode(cdc *codec.Codec, value []byte) ResourceNode {
	resourceNode, err := UnmarshalResourceNode(cdc, value)
	if err != nil {
		panic(err)
	}
	return resourceNode
}

// UnmarshalResourceNode unmarshal a resourceNode from a store value
func UnmarshalResourceNode(cdc *codec.Codec, value []byte) (resourceNode ResourceNode, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &resourceNode)
	return resourceNode, err
}

// AddToken adds tokens to a resource node
func (v ResourceNode) AddToken(amount sdk.Int) ResourceNode {
	v.Tokens = v.Tokens.Add(amount)
	return v
}

// RemoveToken removes tokens from a resource node
func (v ResourceNode) RemoveToken(tokens sdk.Int) ResourceNode {
	if tokens.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", tokens))
	}
	if v.Tokens.LT(tokens) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.Tokens, tokens))
	}
	v.Tokens = v.Tokens.Sub(tokens)
	return v
}

func (v ResourceNode) IsSuspended() bool            { return v.Suspend }
func (v ResourceNode) GetMoniker() string           { return v.Description.Moniker }
func (v ResourceNode) GetStatus() sdk.BondStatus    { return v.Status }
func (v ResourceNode) GetNetworkAddr() string       { return v.NetworkAddress }
func (v ResourceNode) GetPubKey() crypto.PubKey     { return v.PubKey }
func (v ResourceNode) GetAddr() sdk.AccAddress      { return sdk.AccAddress(v.PubKey.Address()) }
func (v ResourceNode) GetTokens() sdk.Int           { return v.Tokens }
func (v ResourceNode) GetOwnerAddr() sdk.AccAddress { return v.OwnerAddress }
