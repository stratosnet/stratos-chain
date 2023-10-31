package types

import (
	"time"

	sdkmath "cosmossdk.io/math"

	stratos "github.com/stratosnet/stratos-chain/types"
)

// IsMature - is the current entry mature
func (e UnbondingNodeEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

// NewUnbondingNode - create a new unbonding Node object
func NewUnbondingNode(networkAddr stratos.SdsAddress, isMetaNode bool, creationHeight int64, minTime time.Time,
	balance sdkmath.Int) UnbondingNode {

	entry := NewUnbondingNodeEntry(creationHeight, minTime, balance)
	return UnbondingNode{
		NetworkAddr: networkAddr.String(),
		IsMetaNode:  isMetaNode,
		Entries:     []*UnbondingNodeEntry{&entry},
	}
}

// NewUnbondingNodeEntry - create a new unbonding Node object
func NewUnbondingNodeEntry(creationHeight int64, completionTime time.Time, balance sdkmath.Int) UnbondingNodeEntry {
	return UnbondingNodeEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: &balance,
		Balance:        &balance,
	}
}

// AddEntry - append entry to the unbonding Node
func (un *UnbondingNode) AddEntry(creationHeight int64, minTime time.Time, balance sdkmath.Int) {
	entry := NewUnbondingNodeEntry(creationHeight, minTime, balance)
	un.Entries = append(un.Entries, &entry)
}

// RemoveEntry - remove entry at index i to the unbonding Node
func (un *UnbondingNode) RemoveEntry(i int64) {
	un.Entries = append(un.Entries[:i], un.Entries[i+1:]...)
}
