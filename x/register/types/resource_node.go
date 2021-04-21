package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type ResourceNode struct {
	OperatorAddress sdk.ValAddress `json:"operator_address" yaml:"operator_address"` // address of the resource node's operator; bech encoded in JSON
	PubKey          crypto.PubKey  `json:"pubkey" yaml:"pubkey"`                     // the public key of the resource node; bech encoded in JSON
	Jailed          bool           `json:"jailed" yaml:"jailed"`                     // has the resource node been jailed from bonded status?
	Status          sdk.BondStatus `json:"status" yaml:"status"`                     // resource node status (bonded/unbonding/unbonded)
	Tokens          sdk.Int        `json:"tokens" yaml:"tokens"`                     // delegated tokens (incl. self-delegation)
	Shares          sdk.Dec        `json:"shares" yaml:"shares"`                     // total shares issued to a resource node
	Description     Description    `json:"description" yaml:"description"`           // description terms for the resource node
}

// NewValidator - initialize a new validator
func NewResourceNode(operator sdk.ValAddress, pubKey crypto.PubKey, description Description) ResourceNode {
	return ResourceNode{
		OperatorAddress: operator,
		PubKey:          pubKey,
		Jailed:          false,
		Status:          sdk.Unbonded,
		Tokens:          sdk.ZeroInt(),
		Description:     description,
	}
}

// MustMarshalResourceNode returns the resourceNode bytes. Panics if fails
func MustMarshalResourceNode(cdc *codec.Codec, resourceNode ResourceNode) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(resourceNode)
}

// unmarshal a redelegation from a store value
func MustUnmarshalResourceNode(cdc *codec.Codec, value []byte) ResourceNode {
	resourceNode, err := UnmarshalResourceNode(cdc, value)
	if err != nil {
		panic(err)
	}
	return resourceNode
}

// unmarshal a redelegation from a store value
func UnmarshalResourceNode(cdc *codec.Codec, value []byte) (resourceNode ResourceNode, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &resourceNode)
	return resourceNode, err
}

// AddTokensToResourceNode adds tokens to a resource node
func (v ResourceNode) AddTokensToResourceNode(amount sdk.Int) (ResourceNode, sdk.Dec) {
	// calculate the shares to issue
	var issuedShares sdk.Dec
	if v.Shares.IsZero() {
		// the first delegation to a resource node sets the exchange rate to one
		issuedShares = amount.ToDec()
	} else {
		shares, err := v.SharesFromTokens(amount)
		if err != nil {
			panic(err)
		}

		issuedShares = shares
	}

	v.Tokens = v.Tokens.Add(amount)
	v.Shares = v.Shares.Add(issuedShares)

	return v, issuedShares
}

// RemoveSharesFromResourceNode removes shares from a resource node.
// NOTE: because token fractions are left in the resource node,
//       the exchange rate of future shares of this resource node can increase.
func (v ResourceNode) RemoveSharesFromResourceNode(delShares sdk.Dec) (ResourceNode, sdk.Int) {
	remainingShares := v.Shares.Sub(delShares)
	var issuedTokens sdk.Int
	if remainingShares.IsZero() {

		// last delegation share gets any trimmings
		issuedTokens = v.Tokens
		v.Tokens = sdk.ZeroInt()
	} else {

		// leave excess tokens in the resource node
		// however fully use all the delegator shares
		issuedTokens = v.TokensFromShares(delShares).TruncateInt()
		v.Tokens = v.Tokens.Sub(issuedTokens)
		if v.Tokens.IsNegative() {
			panic("attempting to remove more tokens than available in resource node")
		}
	}

	v.Shares = remainingShares
	return v, issuedTokens
}

// RemoveTokensFromResourceNode removes tokens from a resource node
func (v ResourceNode) RemoveTokensFromResourceNode(tokens sdk.Int) ResourceNode {
	if tokens.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", tokens))
	}
	if v.Tokens.LT(tokens) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.Tokens, tokens))
	}
	v.Tokens = v.Tokens.Sub(tokens)
	return v
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the resource node has no tokens.
func (v ResourceNode) SharesFromTokens(amt sdk.Int) (sdk.Dec, error) {
	if v.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientShares
	}
	return v.GetShares().MulInt(amt).QuoInt(v.GetTokens()), nil
}

// TokensFromShares calculate the token worth of provided shares
func (v ResourceNode) TokensFromShares(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).Quo(v.Shares)
}

func (v ResourceNode) IsJailed() bool              { return v.Jailed }
func (v ResourceNode) GetMoniker() string          { return v.Description.Moniker }
func (v ResourceNode) GetStatus() sdk.BondStatus   { return v.Status }
func (v ResourceNode) GetOperator() sdk.ValAddress { return v.OperatorAddress }
func (v ResourceNode) GetPubKey() crypto.PubKey    { return v.PubKey }
func (v ResourceNode) GetAddr() sdk.ConsAddress    { return sdk.ConsAddress(v.PubKey.Address()) }
func (v ResourceNode) GetTokens() sdk.Int          { return v.Tokens }
func (v ResourceNode) GetShares() sdk.Dec          { return v.Shares }
