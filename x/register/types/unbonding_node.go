package types

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

// =======================

// UnbondingNode stores all of a single delegator's unbonding bonds
// for a single unbonding node in an time-ordered list
type UnbondingNode struct {
	NetworkAddr    stratos.SdsAddress   `json:"network_addr" yaml:"network_addr"`
	IsIndexingNode bool                 `json:"is_indexing_node yaml:"is_indexing_node`
	Entries        []UnbondingNodeEntry `json:"entries" yaml:"entries"` // unbonding node entries
}

// UnbondingNodeEntry - entry to an UnbondingNode
type UnbondingNodeEntry struct {
	CreationHeight int64     `json:"creation_height" yaml:"creation_height"` // height which the unbonding took place
	CompletionTime time.Time `json:"completion_time" yaml:"completion_time"` // time at which the unbonding delegation will complete
	InitialBalance sdk.Int   `json:"initial_balance" yaml:"initial_balance"` // ustos initially scheduled to receive at completion
	Balance        sdk.Int   `json:"balance" yaml:"balance"`                 // ustos to receive at completion
}

// IsMature - is the current entry mature
func (e UnbondingNodeEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

// NewUnbondingNode - create a new unbonding Node object
func NewUnbondingNode(networkAddr stratos.SdsAddress, isIndexingNode bool, creationHeight int64, minTime time.Time,
	balance sdk.Int) UnbondingNode {

	entry := NewUnbondingNodeEntry(creationHeight, minTime, balance)
	return UnbondingNode{
		NetworkAddr:    networkAddr,
		IsIndexingNode: isIndexingNode,
		Entries:        []UnbondingNodeEntry{entry},
	}
}

// NewUnbondingNodeEntry - create a new unbonding Node object
func NewUnbondingNodeEntry(creationHeight int64, completionTime time.Time,
	balance sdk.Int) UnbondingNodeEntry {

	return UnbondingNodeEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: balance,
		Balance:        balance,
	}
}

// AddEntry - append entry to the unbonding Node
func (un *UnbondingNode) AddEntry(creationHeight int64,
	minTime time.Time, balance sdk.Int) {

	entry := NewUnbondingNodeEntry(creationHeight, minTime, balance)
	un.Entries = append(un.Entries, entry)
}

// RemoveEntry - remove entry at index i to the unbonding Node
func (un *UnbondingNode) RemoveEntry(i int64) {
	un.Entries = append(un.Entries[:i], un.Entries[i+1:]...)
}

// return the unbonding Node
func MustMarshalUnbondingNode(cdc *codec.Codec, uin UnbondingNode) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(uin)
}

// unmarshal a unbonding Node from a store value
func MustUnmarshalUnbondingNode(cdc *codec.Codec, value []byte) UnbondingNode {
	un, err := UnmarshalUnbondingNode(cdc, value)
	if err != nil {
		panic(err)
	}
	return un
}

// unmarshal a unbonding Node from a store value
func UnmarshalUnbondingNode(cdc *codec.Codec, value []byte) (uin UnbondingNode, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &uin)
	return uin, err
}

// nolint
// inefficient but only used in testing
func (un UnbondingNode) Equal(un2 UnbondingNode) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&un)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&un2)
	return bytes.Equal(bz1, bz2)
}

func (un UnbondingNode) GetNetworkAddr() stratos.SdsAddress {
	return un.NetworkAddr
}

// String returns a human readable string representation of an UnbondingNode.
func (un UnbondingNode) String() string {
	out := fmt.Sprintf(`Unbonding Nodes between:
	NetworkAddr:    %s,
	IsIndexingNode: %t,
	Entries:`, un.NetworkAddr, un.IsIndexingNode)
	for i, entry := range un.Entries {
		out += fmt.Sprintf(`    Unbonding Node %d:
      Creation Height:           %v
      Min time to unbond (unix): %v
      Expected balance:          %s`, i, entry.CreationHeight,
			entry.CompletionTime, entry.Balance)
	}
	return out
}

// UnbondingNodes is a collection of UnbondingNode
type UnbondingNodes []UnbondingNode

func (uns UnbondingNodes) String() (out string) {
	for _, u := range uns {
		out += u.String() + "\n"
	}
	return strings.TrimSpace(out)
}
