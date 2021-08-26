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

// IndexingNodes is a collection of indexing node
type IndexingNodes []IndexingNode

func (v IndexingNodes) String() (out string) {
	for _, node := range v {
		out += node.String() + "\n"
	}
	return strings.TrimSpace(out)
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

func (v IndexingNodes) Validate() error {
	for _, node := range v {
		if err := node.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type IndexingNode struct {
	NetworkID    string         `json:"network_id" yaml:"network_id"`       // network id of the indexing node
	PubKey       crypto.PubKey  `json:"pubkey" yaml:"pubkey"`               // the consensus public key of the indexing node; bech encoded in JSON
	Suspend      bool           `json:"suspend" yaml:"suspend"`             // has the indexing node been suspended from bonded status?
	Status       sdk.BondStatus `json:"status" yaml:"status"`               // indexing node status (bonded/unbonding/unbonded)
	Tokens       sdk.Int        `json:"tokens" yaml:"tokens"`               // delegated tokens
	OwnerAddress sdk.AccAddress `json:"owner_address" yaml:"owner_address"` // owner address of the indexing node
	Description  Description    `json:"description" yaml:"description"`     // description terms for the indexing node
	CreationTime time.Time      `json:"creation_time" yaml:"creation_time"`
}

// NewIndexingNode - initialize a new indexing node
func NewIndexingNode(networkID string, pubKey crypto.PubKey, ownerAddr sdk.AccAddress, description Description, creationTime time.Time) IndexingNode {
	return IndexingNode{
		NetworkID:    networkID,
		PubKey:       pubKey,
		Suspend:      false,
		Status:       sdk.Unbonded,
		Tokens:       sdk.ZeroInt(),
		OwnerAddress: ownerAddr,
		Description:  description,
		CreationTime: creationTime,
	}
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
		CreationTime:		%s
	}`, v.NetworkID, pubKey, v.Suspend, v.Status, v.Tokens, v.OwnerAddress, v.Description, v.CreationTime)
}

// AddToken adds tokens to a indexing node
func (v IndexingNode) AddToken(amount sdk.Int) IndexingNode {
	v.Tokens = v.Tokens.Add(amount)
	return v
}

// SubToken removes tokens from a indexing node
func (v IndexingNode) SubToken(tokens sdk.Int) IndexingNode {
	if tokens.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", tokens))
	}
	if v.Tokens.LT(tokens) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.Tokens, tokens))
	}
	v.Tokens = v.Tokens.Sub(tokens)
	return v
}

func (v IndexingNode) Validate() error {
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

func (v IndexingNode) IsSuspended() bool              { return v.Suspend }
func (v IndexingNode) GetMoniker() string             { return v.Description.Moniker }
func (v IndexingNode) GetStatus() sdk.BondStatus      { return v.Status }
func (v IndexingNode) GetNetworkID() string           { return v.NetworkID }
func (v IndexingNode) GetPubKey() crypto.PubKey       { return v.PubKey }
func (v IndexingNode) GetNetworkAddr() sdk.AccAddress { return sdk.AccAddress(v.PubKey.Address()) }
func (v IndexingNode) GetTokens() sdk.Int             { return v.Tokens }
func (v IndexingNode) GetOwnerAddr() sdk.AccAddress   { return v.OwnerAddress }
func (v IndexingNode) GetCreationTime() time.Time     { return v.CreationTime }

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

// =======================
//
//// UnbondingIndexingNode stores all of a single delegator's unbonding bonds
//// for a single validator in an time-ordered list
//type UnbondingIndexingNode struct {
//	NetworkAddr sdk.AccAddress               `json:"network_addr" yaml:"network_addr"` // network id of the indexing node, for instance, sds://blablabla
//	Entries     []UnbondingIndexingNodeEntry `json:"entries" yaml:"entries"`           // unbonding node entries
//}
//
//// UnbondingIndexingNodeEntry - entry to an UnbondingIndexingNode
//type UnbondingIndexingNodeEntry struct {
//	CreationHeight int64     `json:"creation_height" yaml:"creation_height"` // height which the unbonding took place
//	CompletionTime time.Time `json:"completion_time" yaml:"completion_time"` // time at which the unbonding delegation will complete
//	InitialBalance sdk.Int   `json:"initial_balance" yaml:"initial_balance"` // atoms initially scheduled to receive at completion
//	Balance        sdk.Int   `json:"balance" yaml:"balance"`                 // atoms to receive at completion
//}
//
//// IsMature - is the current entry mature
//func (e UnbondingIndexingNodeEntry) IsMature(currentTime time.Time) bool {
//	return !e.CompletionTime.After(currentTime)
//}
//
//// NewUnbondingIndexingNode - create a new unbonding IndexingNode object
//func NewUnbondingIndexingNode(networkAddr sdk.AccAddress, creationHeight int64, minTime time.Time,
//	balance sdk.Int) UnbondingIndexingNode {
//
//	entry := NewUnbondingIndexingNodeEntry(creationHeight, minTime, balance)
//	return UnbondingIndexingNode{
//		NetworkAddr: networkAddr,
//		Entries:     []UnbondingIndexingNodeEntry{entry},
//	}
//}
//
//// NewUnbondingIndexingNodeEntry - create a new unbonding IndexingNode object
//func NewUnbondingIndexingNodeEntry(creationHeight int64, completionTime time.Time,
//	balance sdk.Int) UnbondingIndexingNodeEntry {
//
//	return UnbondingIndexingNodeEntry{
//		CreationHeight: creationHeight,
//		CompletionTime: completionTime,
//		InitialBalance: balance,
//		Balance:        balance,
//	}
//}
//
//// AddEntry - append entry to the unbonding IndexingNode
//func (uin *UnbondingIndexingNode) AddEntry(creationHeight int64,
//	minTime time.Time, balance sdk.Int) {
//
//	entry := NewUnbondingIndexingNodeEntry(creationHeight, minTime, balance)
//	uin.Entries = append(uin.Entries, entry)
//}
//
//// RemoveEntry - remove entry at index i to the unbonding IndexingNode
//func (uin *UnbondingIndexingNode) RemoveEntry(i int64) {
//	uin.Entries = append(uin.Entries[:i], uin.Entries[i+1:]...)
//}
//
//// return the unbonding IndexingNode
//func MustMarshalUIN(cdc *codec.Codec, uin UnbondingIndexingNode) []byte {
//	return cdc.MustMarshalBinaryLengthPrefixed(uin)
//}
//
//// unmarshal a unbonding IndexingNode from a store value
//func MustUnmarshalUIN(cdc *codec.Codec, value []byte) UnbondingIndexingNode {
//	uin, err := UnmarshalUIN(cdc, value)
//	if err != nil {
//		panic(err)
//	}
//	return uin
//}
//
//// unmarshal a unbonding IndexingNode from a store value
//func UnmarshalUIN(cdc *codec.Codec, value []byte) (uin UnbondingIndexingNode, err error) {
//	err = cdc.UnmarshalBinaryLengthPrefixed(value, &uin)
//	return uin, err
//}
//
//// nolint
//// inefficient but only used in testing
//func (uin UnbondingIndexingNode) Equal(uin2 UnbondingIndexingNode) bool {
//	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&uin)
//	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&uin2)
//	return bytes.Equal(bz1, bz2)
//}
//
//func (uin UnbondingIndexingNode) GetNetworkAddr() sdk.AccAddress {
//	return uin.NetworkAddr
//}
//
//// String returns a human readable string representation of an UnbondingIndexingNode.
//func (uin UnbondingIndexingNode) String() string {
//	out := fmt.Sprintf(`Unbonding IndexingNodes between:
//	NetworkAddr:    %s,
//	Entries:`, uin.NetworkAddr)
//	for i, entry := range uin.Entries {
//		out += fmt.Sprintf(`    Unbonding IndexingNode %d:
//      Creation Height:           %v
//      Min time to unbond (unix): %v
//      Expected balance:          %s`, i, entry.CreationHeight,
//			entry.CompletionTime, entry.Balance)
//	}
//	return out
//}
//
//// UnbondingIndexingNodes is a collection of UnbondingIndexingNode
//type UnbondingIndexingNodes []UnbondingIndexingNode
//
//func (uins UnbondingIndexingNodes) String() (out string) {
//	for _, u := range uins {
//		out += u.String() + "\n"
//	}
//	return strings.TrimSpace(out)
//}
