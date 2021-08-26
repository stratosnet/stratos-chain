package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"sort"
	"strings"
	"time"
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
	NetworkID    string         `json:"network_id" yaml:"network_id"`       // network id of the resource node
	PubKey       crypto.PubKey  `json:"pubkey" yaml:"pubkey"`               // the public key of the resource node; bech encoded in JSON
	Suspend      bool           `json:"suspend" yaml:"suspend"`             // has the resource node been suspended from bonded status?
	Status       sdk.BondStatus `json:"status" yaml:"status"`               // resource node bond status (bonded/unbonding/unbonded)
	Tokens       sdk.Int        `json:"tokens" yaml:"tokens"`               // delegated tokens
	OwnerAddress sdk.AccAddress `json:"owner_address" yaml:"owner_address"` // owner address of the resource node
	Description  Description    `json:"description" yaml:"description"`     // description terms for the resource node
	NodeType     string         `json:"node_type" yaml:"node_type"`
	CreationTime time.Time      `json:"creation_time" yaml:"creation_time"`
}

// NewResourceNode - initialize a new resource node
func NewResourceNode(networkID string, pubKey crypto.PubKey, ownerAddr sdk.AccAddress,
	description Description, nodeType string, creationTime time.Time) ResourceNode {
	return ResourceNode{
		NetworkID:    networkID,
		PubKey:       pubKey,
		Suspend:      false,
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
	pubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, v.PubKey)
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
	if v.NetworkID == "" {
		return ErrEmptyNodeId
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

func (v ResourceNode) IsSuspended() bool              { return v.Suspend }
func (v ResourceNode) GetMoniker() string             { return v.Description.Moniker }
func (v ResourceNode) GetStatus() sdk.BondStatus      { return v.Status }
func (v ResourceNode) GetNetworkID() string           { return v.NetworkID }
func (v ResourceNode) GetPubKey() crypto.PubKey       { return v.PubKey }
func (v ResourceNode) GetNetworkAddr() sdk.AccAddress { return sdk.AccAddress(v.PubKey.Address()) }
func (v ResourceNode) GetTokens() sdk.Int             { return v.Tokens }
func (v ResourceNode) GetOwnerAddr() sdk.AccAddress   { return v.OwnerAddress }
func (v ResourceNode) GetNodeType() string            { return v.NodeType }
func (v ResourceNode) GetCreationTime() time.Time     { return v.CreationTime }

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

//
//// UnbondingResourceNode stores all of a single delegator's unbonding bonds
//// for a single validator in an time-ordered list
//type UnbondingResourceNode struct {
//	NetworkAddr sdk.AccAddress               `json:"network_addr" yaml:"network_addr"` // network id of the indexing node, for instance, sds://blablabla
//	Entries      []UnbondingResourceNodeEntry `json:"entries" yaml:"entries"` // unbonding node entries
//	//NetworkID    string                       `json:"network_id" yaml:"network_id"`       // network id of the indexing node, for instance, sds://blablabla
//	//PubKey       crypto.PubKey                `json:"pubkey" yaml:"pubkey"`               // the consensus public key of the indexing node; bech encoded in JSON
//	//Suspend      bool                         `json:"suspend" yaml:"suspend"`             // has the indexing node been suspended from bonded status?
//	//Status       sdk.BondStatus               `json:"status" yaml:"status"`               // indexing node status (bonded/unbonding/unbonded)
//	//Tokens       sdk.Int                      `json:"tokens" yaml:"tokens"`               // delegated tokens
//	//OwnerAddress sdk.AccAddress               `json:"owner_address" yaml:"owner_address"` // owner address of the indexing node
//	//Description  Description                  `json:"description" yaml:"description"`     // description terms for the indexing node
//	//CreationTime time.Time                    `json:"creation_time" yaml:"creation_time"`
//}
//
//// UnbondingResourceNodeEntry - entry to an UnbondingResourceNode
//type UnbondingResourceNodeEntry struct {
//	CreationHeight int64     `json:"creation_height" yaml:"creation_height"` // height which the unbonding took place
//	CompletionTime time.Time `json:"completion_time" yaml:"completion_time"` // time at which the unbonding delegation will complete
//	InitialBalance sdk.Int   `json:"initial_balance" yaml:"initial_balance"` // atoms initially scheduled to receive at completion
//	Balance        sdk.Int   `json:"balance" yaml:"balance"`                 // atoms to receive at completion
//}
//
//// IsMature - is the current entry mature
//func (e UnbondingResourceNodeEntry) IsMature(currentTime time.Time) bool {
//	return !e.CompletionTime.After(currentTime)
//}
//
//// NewUnbondingResourceNode - create a new unbonding delegation object
//func NewUnbondingResourceNode(networkAddr sdk.AccAddress, creationHeight int64, minTime time.Time,
//	balance sdk.Int) UnbondingResourceNode {
//
//	entry := NewUnbondingResourceNodeEntry(creationHeight, minTime, balance)
//	return UnbondingResourceNode{
//		NetworkAddr: networkAddr,
//		Entries:     []UnbondingResourceNodeEntry{entry},
//		//NetworkID:    resourceNode.NetworkID,
//		//PubKey:       resourceNode.PubKey,
//		//Suspend:      resourceNode.Suspend,
//		//Status:       resourceNode.Status,
//		//Tokens:       resourceNode.Tokens,
//		//OwnerAddress: resourceNode.OwnerAddress,
//		//Description:  resourceNode.Description,
//		//CreationTime: resourceNode.CreationTime,
//	}
//}
//
//// NewUnbondingResourceNodeEntry - create a new unbonding ResourceNode object
//func NewUnbondingResourceNodeEntry(creationHeight int64, completionTime time.Time,
//	balance sdk.Int) UnbondingResourceNodeEntry {
//
//	return UnbondingResourceNodeEntry{
//		CreationHeight: creationHeight,
//		CompletionTime: completionTime,
//		InitialBalance: balance,
//		Balance:        balance,
//	}
//}
//
//// AddEntry - append entry to the unbonding ResourceNode
//func (urn *UnbondingResourceNode) AddEntry(creationHeight int64,
//	minTime time.Time, balance sdk.Int) {
//
//	entry := NewUnbondingResourceNodeEntry(creationHeight, minTime, balance)
//	urn.Entries = append(urn.Entries, entry)
//}
//
//// RemoveEntry - remove entry at index i to the unbonding ResourceNode
//func (urn *UnbondingResourceNode) RemoveEntry(i int64) {
//	urn.Entries = append(urn.Entries[:i], urn.Entries[i+1:]...)
//}
//
//// return the unbonding ResourceNode
//func MustMarshalURN(cdc *codec.Codec, uin UnbondingResourceNode) []byte {
//	return cdc.MustMarshalBinaryLengthPrefixed(uin)
//}
//
//// unmarshal a unbonding ResourceNode from a store value
//func MustUnmarshalURN(cdc *codec.Codec, value []byte) UnbondingResourceNode {
//	uin, err := UnmarshalURN(cdc, value)
//	if err != nil {
//		panic(err)
//	}
//	return uin
//}
//
//// unmarshal a unbonding ResourceNode from a store value
//func UnmarshalURN(cdc *codec.Codec, value []byte) (uin UnbondingResourceNode, err error) {
//	err = cdc.UnmarshalBinaryLengthPrefixed(value, &uin)
//	return uin, err
//}
//
//// nolint
//// inefficient but only used in testing
//func (urn UnbondingResourceNode) Equal(urn2 UnbondingResourceNode) bool {
//	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&urn)
//	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&urn2)
//	return bytes.Equal(bz1, bz2)
//}
//
//func (urn UnbondingResourceNode) GetNetworkAddr() sdk.AccAddress { return urn.NetworkAddr}
//
//// String returns a human readable string representation of an UnbondingResourceNode.
//func (urn UnbondingResourceNode) String() string {
//	out := fmt.Sprintf(`Unbonding ResourceNodes between:
//	NetworkAddr:    %s,
//	Entries:`, urn.NetworkAddr)
//	for i, entry := range urn.Entries {
//		out += fmt.Sprintf(`    Unbonding ResourceNode %d:
//      Creation Height:           %v
//      Min time to unbond (unix): %v
//      Expected balance:          %s`, i, entry.CreationHeight,
//			entry.CompletionTime, entry.Balance)
//	}
//	return out
//}
//
//// UnbondingResourceNodes is a collection of UnbondingResourceNode
//type UnbondingResourceNodes []UnbondingResourceNode
//
//func (uins UnbondingResourceNodes) String() (out string) {
//	for _, u := range uins {
//		out += u.String() + "\n"
//	}
//	return strings.TrimSpace(out)
//}
