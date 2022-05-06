package types

import (
	"bytes"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

// =======================

// IsMature - is the current entry mature
func (e UnbondingNodeEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

// NewUnbondingNode - create a new unbonding Node object
func NewUnbondingNode(networkAddr stratos.SdsAddress, isIndexingNode bool, creationHeight int64, minTime time.Time,
	balance sdk.Int) UnbondingNode {

	entry := NewUnbondingNodeEntry(creationHeight, minTime, balance)
	return UnbondingNode{
		NetworkAddr:    networkAddr.String(),
		IsIndexingNode: isIndexingNode,
		Entries:        []*UnbondingNodeEntry{&entry},
	}
}

// NewUnbondingNodeEntry - create a new unbonding Node object
func NewUnbondingNodeEntry(creationHeight int64, completionTime time.Time,
	balance sdk.Int) UnbondingNodeEntry {

	return UnbondingNodeEntry{
		CreationHeight: creationHeight,
		CompletionTime: &completionTime,
		InitialBalance: &balance,
		Balance:        &balance,
	}
}

// AddEntry - append entry to the unbonding Node
func (un *UnbondingNode) AddEntry(creationHeight int64,
	minTime time.Time, balance sdk.Int) {

	entry := NewUnbondingNodeEntry(creationHeight, minTime, balance)
	un.Entries = append(un.Entries, &entry)
}

// RemoveEntry - remove entry at index i to the unbonding Node
func (un *UnbondingNode) RemoveEntry(i int64) {
	un.Entries = append(un.Entries[:i], un.Entries[i+1:]...)
}

// MustMarshalUnbondingNode return the unbonding Node
func MustMarshalUnbondingNode(cdc codec.BinaryCodec, uin UnbondingNode) []byte {
	return cdc.MustMarshalLengthPrefixed(&uin)
}

// MustUnmarshalUnbondingNode unmarshal a unbonding Node from a store value
func MustUnmarshalUnbondingNode(cdc codec.BinaryCodec, value []byte) UnbondingNode {
	un, err := UnmarshalUnbondingNode(cdc, value)
	if err != nil {
		panic(err)
	}
	return un
}

// UnmarshalUnbondingNode unmarshal a unbonding Node from a store value
func UnmarshalUnbondingNode(cdc codec.BinaryCodec, value []byte) (uin UnbondingNode, err error) {
	err = cdc.UnmarshalLengthPrefixed(value, &uin)
	return uin, err
}

// Equal nolint
// inefficient but only used in testing
func (un UnbondingNode) Equal(un2 UnbondingNode) bool {
	var cdc codec.BinaryCodec
	bz1 := cdc.MustMarshalLengthPrefixed(&un)
	bz2 := cdc.MustMarshalLengthPrefixed(&un2)
	return bytes.Equal(bz1, bz2)
}

// UnbondingNodes is a collection of UnbondingNode
//type UnbondingNodes []UnbondingNode

//func (uns UnbondingNodes) String() (out string) {
//	for _, u := range uns {
//		out += u.String() + "\n"
//	}
//	return strings.TrimSpace(out)
//}
