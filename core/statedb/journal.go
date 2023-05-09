package statedb

import (
	"bytes"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type dirtyObj struct {
	storeKey sdk.StoreKey
	key      StorageKey
}

func (do *dirtyObj) ObjKey() []byte {
	return append([]byte(do.storeKey.Name()), do.key[:]...)
}

// journalEntry is a modification entry in the state change journal that can be
// reverted on demand.
type journalEntry interface {
	// revert undoes the changes introduced by this journal entry.
	revert(*KeestateDB)

	// dirtied returns the key modified by this journal entry.
	dirtied() *dirtyObj
}

// journal contains the list of state modifications applied since the last state
// commit. These are tracked to be able to be reverted in the case of an execution
// exception or request for reversal.
type journal struct {
	entries []journalEntry                      // Current changes tracked by the journal
	dirties map[sdk.StoreKey]map[StorageKey]int // Dirty accounts and the number of changes
}

// newJournal creates a new initialized journal.
func newJournal() *journal {
	return &journal{
		dirties: make(map[sdk.StoreKey]map[StorageKey]int),
	}
}

// sortedDirties sort the dirty addresses for deterministic iteration
func (j *journal) sortedDirties() []*dirtyObj {
	keys := make([]*dirtyObj, 0)
	t := 0
	for i := range j.dirties {
		for k := range j.dirties[i] {
			keys[t] = &dirtyObj{
				storeKey: i,
				key:      k,
			}
			t++
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i].ObjKey(), keys[j].ObjKey()) < 0
	})
	return keys
}

// append inserts a new modification entry to the end of the change journal.
func (j *journal) append(entry journalEntry) {
	j.entries = append(j.entries, entry)
	if dirty := entry.dirtied(); dirty != nil {
		if _, ok := j.dirties[dirty.storeKey]; !ok {
			j.dirties[dirty.storeKey] = make(map[StorageKey]int)
		}
		j.dirties[dirty.storeKey][dirty.key]++
	}
}

// revert undoes a batch of journalled modifications along with any reverted
// dirty handling too.
func (j *journal) revert(statedb *KeestateDB, snapshot int) {
	for i := len(j.entries) - 1; i >= snapshot; i-- {
		// Undo the changes made by the operation
		j.entries[i].revert(statedb)

		// Drop any dirty tracking induced by the change
		if dirty := j.entries[i].dirtied(); dirty != nil {
			if j.dirties[dirty.storeKey][dirty.key]--; j.dirties[dirty.storeKey][dirty.key] == 0 {
				delete(j.dirties[dirty.storeKey], dirty.key)
			}
		}
	}
	j.entries = j.entries[:snapshot]
}

// length returns the current number of entries in the journal.
func (j *journal) length() int {
	return len(j.entries)
}

type (
	storageChange struct {
		storeKey sdk.StoreKey
		key      StorageKey
		prevalue StorageValue
	}
	createObjectChange struct {
		storeKey sdk.StoreKey
		key      StorageKey
	}
	resetObjectChange struct {
		prev *stateObject
	}
)

func (ch storageChange) revert(s *KeestateDB) {
	s.getStateObject(ch.storeKey, ch.key[:]).SetState(ch.prevalue.Result())
}

func (ch storageChange) dirtied() *dirtyObj {
	return &dirtyObj{storeKey: ch.storeKey, key: ch.key}
}

func (ch createObjectChange) revert(s *KeestateDB) {
	delete(s.stateObjects[ch.storeKey], ch.key)
}

func (ch createObjectChange) dirtied() *dirtyObj {
	return &dirtyObj{storeKey: ch.storeKey, key: ch.key}
}

func (ch resetObjectChange) revert(s *KeestateDB) {
	s.setStateObject(ch.prev)
}

func (ch resetObjectChange) dirtied() *dirtyObj {
	return nil
}
