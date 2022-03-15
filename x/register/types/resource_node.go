package types

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/tendermint/tendermint/crypto"
)

type NodeType uint8

const (
	STORAGE     NodeType = 4
	DATABASE    NodeType = 2
	COMPUTATION NodeType = 1
)

func (n NodeType) Type() string {
	switch n {
	case 7:
		return "storage/database/computation"
	case 6:
		return "database/storage"
	case 5:
		return "computation/storage"
	case 4:
		return "storage"
	case 3:
		return "computation/database"
	case 2:
		return "database"
	case 1:
		return "computation"
	}
	return "UNKNOWN"
}

func (n NodeType) String() string {
	return n.Type()
}

// ResourceNodes is a collection of resource node
type ResourceNodes []ResourceNode

func (v ResourceNodes) String() (out string) {
	for _, node := range v {
		out += node.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// Sort ResourceNodes sorts ResourceNode array in ascending owner address order
func (v ResourceNodes) Sort() {
	sort.Sort(v)
}

// Len implements sort interface
func (v ResourceNodes) Len() int {
	return len(v)
}

// Less implements sort interface
func (v ResourceNodes) Less(i, j int) bool {
	return v[i].Tokens.LT(v[j].Tokens)
}

// Swap implements sort interface
func (v ResourceNodes) Swap(i, j int) {
	it := v[i]
	v[i] = v[j]
	v[j] = it
}

func (v ResourceNodes) Validate() error {
	for _, node := range v {
		if err := node.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type ResourceNode struct {
	NetworkID    stratos.SdsAddress `json:"network_id" yaml:"network_id"`       // network id of the resource node, sds://...
	PubKey       crypto.PubKey      `json:"pubkey" yaml:"pubkey"`               // the public key of the resource node; bech encoded in JSON
	Suspend      bool               `json:"suspend" yaml:"suspend"`             // has the resource node been suspended from bonded status?
	Status       sdk.BondStatus     `json:"status" yaml:"status"`               // resource node bond status (bonded/unbonding/unbonded)
	Tokens       sdk.Int            `json:"tokens" yaml:"tokens"`               // delegated tokens
	OwnerAddress sdk.AccAddress     `json:"owner_address" yaml:"owner_address"` // owner address of the resource node
	Description  Description        `json:"description" yaml:"description"`     // description terms for the resource node
	NodeType     NodeType           `json:"node_type" yaml:"node_type"`
	CreationTime time.Time          `json:"creation_time" yaml:"creation_time"`
}

// NewResourceNode - initialize a new resource node
func NewResourceNode(networkID stratos.SdsAddress, pubKey crypto.PubKey, ownerAddr sdk.AccAddress,
	description Description, nodeType NodeType, creationTime time.Time) ResourceNode {
	return ResourceNode{
		NetworkID:    networkID,
		PubKey:       pubKey,
		Suspend:      true,
		Status:       sdk.Unbonded,
		Tokens:       sdk.ZeroInt(),
		OwnerAddress: ownerAddr,
		Description:  description,
		NodeType:     nodeType,
		CreationTime: creationTime,
	}
}

// String returns a human readable string representation of a resource node.
func (v ResourceNode) String() string {
	pubKey, err := stratos.Bech32ifyPubKey(stratos.Bech32PubKeyTypeAccPub, v.PubKey)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`ResourceNode:{
		Network Id:	        %s
  		Pubkey:				%s
  		Suspend:			%v
  		Status:				%s
  		Tokens:				%s
		Owner Address: 		%s
  		Description:		%s
  		CreationTime:		%s
	}`, v.NetworkID, pubKey, v.Suspend, v.Status, v.Tokens, v.OwnerAddress, v.Description, v.CreationTime)
}

// AddToken adds tokens to a resource node
func (v ResourceNode) AddToken(amount sdk.Int) ResourceNode {
	v.Tokens = v.Tokens.Add(amount)
	return v
}

// SubToken removes tokens from a resource node
func (v ResourceNode) SubToken(tokens sdk.Int) ResourceNode {
	if tokens.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", tokens))
	}
	if v.Tokens.LT(tokens) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.Tokens, tokens))
	}
	v.Tokens = v.Tokens.Sub(tokens)
	return v
}

func (v ResourceNode) Validate() error {
	if v.NetworkID.Empty() {
		return ErrEmptyNodeId
	}
	if v.NetworkID.Equals(stratos.SdsAddress(v.PubKey.Address())) {
		return ErrInvalidNetworkAddr
	}
	if len(v.PubKey.Bytes()) == 0 {
		return ErrEmptyPubKey
	}
	if v.OwnerAddress.Empty() {
		return ErrEmptyOwnerAddr
	}
	if v.Tokens.LT(sdk.ZeroInt()) {
		return ErrValueNegative
	}
	if v.Description.Moniker == "" {
		return ErrEmptyMoniker
	}
	return nil
}

// IsBonded checks if the node status equals Bonded
func (v ResourceNode) IsBonded() bool {
	return v.GetStatus().Equal(sdk.Bonded)
}

// IsUnBonded checks if the node status equals Unbonded
func (v ResourceNode) IsUnBonded() bool {
	return v.GetStatus().Equal(sdk.Unbonded)
}

// IsUnBonding checks if the node status equals Unbonding
func (v ResourceNode) IsUnBonding() bool {
	return v.GetStatus().Equal(sdk.Unbonding)
}

func (v ResourceNode) IsSuspended() bool                { return v.Suspend }
func (v ResourceNode) GetMoniker() string               { return v.Description.Moniker }
func (v ResourceNode) GetStatus() sdk.BondStatus        { return v.Status }
func (v ResourceNode) GetNetworkID() stratos.SdsAddress { return v.NetworkID }
func (v ResourceNode) GetPubKey() crypto.PubKey         { return v.PubKey }
func (v ResourceNode) GetNetworkAddr() stratos.SdsAddress {
	return stratos.SdsAddress(v.PubKey.Address())
}
func (v ResourceNode) GetTokens() sdk.Int           { return v.Tokens }
func (v ResourceNode) GetOwnerAddr() sdk.AccAddress { return v.OwnerAddress }
func (v ResourceNode) GetNodeType() string          { return v.NodeType.String() }
func (v ResourceNode) GetCreationTime() time.Time   { return v.CreationTime }

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

func (resourceNode ResourceNode) Equal(resourceNode2 ResourceNode) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&resourceNode)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&resourceNode2)
	return bytes.Equal(bz1, bz2)
}
