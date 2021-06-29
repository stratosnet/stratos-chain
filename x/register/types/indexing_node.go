package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/exported"
	"github.com/tendermint/tendermint/crypto"
	"sort"
	"strings"
	"time"
)

type IndexingNode struct {
	NetworkID    string         `json:"network_id" yaml:"network_id"`       // network address of the indexing node
	PubKey       crypto.PubKey  `json:"pubkey" yaml:"pubkey"`               // the consensus public key of the indexing node; bech encoded in JSON
	Suspend      bool           `json:"suspend" yaml:"suspend"`             // has the indexing node been suspended from bonded status?
	Status       sdk.BondStatus `json:"status" yaml:"status"`               // indexing node status (bonded/unbonding/unbonded)
	Tokens       sdk.Int        `json:"tokens" yaml:"tokens"`               // delegated tokens
	OwnerAddress sdk.AccAddress `json:"owner_address" yaml:"owner_address"` // owner address of the indexing node
	Description  Description    `json:"description" yaml:"description"`     // description terms for the indexing node
}

// IndexingNodes is a collection of indexing node
type IndexingNodes []IndexingNode

func (v IndexingNodes) String() (out string) {
	for _, node := range v {
		out += node.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// ToSDKIndexingNodes -  convenience function convert []IndexingNodes to []sdk.IndexingNodes
func (v IndexingNodes) ToSDKIndexingNodes() (indexingNodes []exported.IndexingNodeI) {
	for _, node := range v {
		indexingNodes = append(indexingNodes, node)
	}
	return indexingNodes
}

// Sort IndexingNodes sorts IndexingNode array in ascending owner address order
func (v IndexingNodes) Sort() {
	sort.Sort(v)
}

// Implements sort interface
func (v IndexingNodes) Len() int {
	return len(v)
}

// Implements sort interface
func (v IndexingNodes) Less(i, j int) bool {
	return v[i].Tokens.LT(v[j].Tokens)
}

// Implements sort interface
func (v IndexingNodes) Swap(i, j int) {
	it := v[i]
	v[i] = v[j]
	v[j] = it
}

// NewIndexingNode - initialize a new indexing node
func NewIndexingNode(networkID string, pubKey crypto.PubKey, ownerAddr sdk.AccAddress, description Description) IndexingNode {
	return IndexingNode{
		NetworkID:    networkID,
		PubKey:       pubKey,
		Suspend:      false,
		Status:       sdk.Unbonded,
		Tokens:       sdk.ZeroInt(),
		OwnerAddress: ownerAddr,
		Description:  description,
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

// String returns a human readable string representation of a indexing node.
func (v IndexingNode) String() string {
	pubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, v.PubKey)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`IndexingNode:{
		Network ID:			%s
  		Pubkey:				%s
  		Suspend:			%v
  		Status:				%s
  		Tokens:				%s
		Owner Address: 		%s
  		Description:		%s
	}`, v.NetworkID, pubKey, v.Suspend, v.Status, v.Tokens, v.OwnerAddress, v.Description)
}

// AddToken adds tokens to a indexing node
func (v IndexingNode) AddToken(amount sdk.Int) IndexingNode {
	v.Tokens = v.Tokens.Add(amount)
	//if v.Status.Equal(sdk.Unbonded) {
	//	v.Status = sdk.Bonded
	//}
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

func (v IndexingNode) IsSuspended() bool              { return v.Suspend }
func (v IndexingNode) GetMoniker() string             { return v.Description.Moniker }
func (v IndexingNode) GetStatus() sdk.BondStatus      { return v.Status }
func (v IndexingNode) GetNetworkID() string           { return v.NetworkID }
func (v IndexingNode) GetPubKey() crypto.PubKey       { return v.PubKey }
func (v IndexingNode) GetNetworkAddr() sdk.AccAddress { return sdk.AccAddress(v.PubKey.Address()) }
func (v IndexingNode) GetTokens() sdk.Int             { return v.Tokens }
func (v IndexingNode) GetOwnerAddr() sdk.AccAddress   { return v.OwnerAddress }

type VoteOpinion bool

const (
	Approve            VoteOpinion = true
	Reject             VoteOpinion = false
	VoteOpinionApprove             = "Approve"
	VoteOpinionReject              = "Reject"
)

func VoteOpinionFromBool(b bool) VoteOpinion {
	if b {
		return Approve
	} else {
		return Reject
	}
}

// Equal compares two VoteOpinion instances
func (v VoteOpinion) Equal(v2 VoteOpinion) bool {
	return v == v2
}

// String implements the Stringer interface for VoteOpinion.
func (v VoteOpinion) String() string {
	if v {
		return VoteOpinionApprove
	} else {
		return VoteOpinionReject
	}
}

type IndexingNodeRegistrationVotePool struct {
	NodeAddress sdk.AccAddress   `json:"node_address" yaml:"node_address"`
	ApproveList []sdk.AccAddress `json:"approve_list" yaml:"approve_list"`
	RejectList  []sdk.AccAddress `json:"reject_list" yaml:"reject_list"`
	ExpireTime  time.Time        `json:"expire_time" yaml:"expire_time"`
}

func NewRegistrationVotePool(nodeAddress sdk.AccAddress, approveList []sdk.AccAddress, rejectList []sdk.AccAddress, expireTime time.Time) IndexingNodeRegistrationVotePool {
	return IndexingNodeRegistrationVotePool{
		NodeAddress: nodeAddress,
		ApproveList: approveList,
		RejectList:  rejectList,
		ExpireTime:  expireTime,
	}
}
