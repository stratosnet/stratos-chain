package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type IndexingNode struct {
	OperatorAddress sdk.ValAddress `json:"operator_address" yaml:"operator_address"` // address of the indexing node's operator; bech encoded in JSON
	PubKey          crypto.PubKey  `json:"pubkey" yaml:"pubkey"`                     // the consensus public key of the indexing node; bech encoded in JSON
	Jailed          bool           `json:"jailed" yaml:"jailed"`                     // has the indexing node been jailed from bonded status?
	Status          sdk.BondStatus `json:"status" yaml:"status"`                     // indexing node status (bonded/unbonding/unbonded)
	Tokens          sdk.Int        `json:"tokens" yaml:"tokens"`                     // delegated tokens (incl. self-delegation)
	Shares          sdk.Dec        `json:"shares" yaml:"shares"`                     // total shares issued to a indexing node's delegators
	Description     Description    `json:"description" yaml:"description"`           // description terms for the indexing node
}

// NewValidator - initialize a new validator
func NewIndexingNode(operator sdk.ValAddress, pubKey crypto.PubKey, description Description) IndexingNode {
	return IndexingNode{
		OperatorAddress: operator,
		PubKey:          pubKey,
		Jailed:          false,
		Status:          sdk.Unbonded,
		Tokens:          sdk.ZeroInt(),
		Description:     description,
	}
}

// MustMarshalIndexingNode returns the indexingNode bytes. Panics if fails
func MustMarshalIndexingNode(cdc *codec.Codec, indexingNode IndexingNode) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(indexingNode)
}

// unmarshal a redelegation from a store value
func MustUnmarshalIndexingNode(cdc *codec.Codec, value []byte) IndexingNode {
	indexingNode, err := UnmarshalIndexingNode(cdc, value)
	if err != nil {
		panic(err)
	}
	return indexingNode
}

// unmarshal a redelegation from a store value
func UnmarshalIndexingNode(cdc *codec.Codec, value []byte) (indexingNode IndexingNode, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &indexingNode)
	return indexingNode, err
}

// AddTokensToIndexingNode adds tokens to a indexing node
func (v IndexingNode) AddTokensToIndexingNode(amount sdk.Int) (IndexingNode, sdk.Dec) {
	// calculate the shares to issue
	var issuedShares sdk.Dec
	if v.Shares.IsZero() {
		// the first delegation to a indexing node sets the exchange rate to one
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

// RemoveSharesFromIndexingNode removes shares from a indexing node.
// NOTE: because token fractions are left in the indexing node,
//       the exchange rate of future shares of this indexing node can increase.
func (v IndexingNode) RemoveSharesFromIndexingNode(delShares sdk.Dec) (IndexingNode, sdk.Int) {
	remainingShares := v.Shares.Sub(delShares)
	var issuedTokens sdk.Int
	if remainingShares.IsZero() {

		// last delegation share gets any trimmings
		issuedTokens = v.Tokens
		v.Tokens = sdk.ZeroInt()
	} else {

		// leave excess tokens in the indexing node
		// however fully use all the delegator shares
		issuedTokens = v.TokensFromShares(delShares).TruncateInt()
		v.Tokens = v.Tokens.Sub(issuedTokens)
		if v.Tokens.IsNegative() {
			panic("attempting to remove more tokens than available in indexing node")
		}
	}

	v.Shares = remainingShares
	return v, issuedTokens
}

// RemoveTokensFromIndexingNode removes tokens from a indexing node
func (v IndexingNode) RemoveTokensFromIndexingNode(tokens sdk.Int) IndexingNode {
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
// returns an error if the indexing node has no tokens.
func (v IndexingNode) SharesFromTokens(amt sdk.Int) (sdk.Dec, error) {
	if v.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientShares
	}
	return v.GetShares().MulInt(amt).QuoInt(v.GetTokens()), nil
}

// TokensFromShares calculate the token worth of provided shares
func (v IndexingNode) TokensFromShares(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).Quo(v.Shares)
}

func (v IndexingNode) IsJailed() bool              { return v.Jailed }
func (v IndexingNode) GetMoniker() string          { return v.Description.Moniker }
func (v IndexingNode) GetStatus() sdk.BondStatus   { return v.Status }
func (v IndexingNode) GetOperator() sdk.ValAddress { return v.OperatorAddress }
func (v IndexingNode) GetPubKey() crypto.PubKey    { return v.PubKey }
func (v IndexingNode) GetAddr() sdk.ConsAddress    { return sdk.ConsAddress(v.PubKey.Address()) }
func (v IndexingNode) GetTokens() sdk.Int          { return v.Tokens }
func (v IndexingNode) GetShares() sdk.Dec          { return v.Shares }
