package statedb

import (
	"fmt"
	"sort"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// revision is the identifier of a version of state.
// it consists of an auto-increment id and a journal index.
// it's safer to use than using journal index alone.
type revision struct {
	id           int
	journalIndex int
}

type KeestateDB struct {
	ctx sdk.Context

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	journal        *journal
	validRevisions []revision
	nextRevisionID int

	stateObjects map[storetypes.StoreKey]map[StorageKey]*stateObject
}

// New creates a new state from a given trie.
func New(ctx sdk.Context) *KeestateDB {
	return &KeestateDB{
		ctx:          ctx,
		journal:      newJournal(),
		stateObjects: make(map[storetypes.StoreKey]map[StorageKey]*stateObject),
	}
}

func (ks *KeestateDB) getStateObject(storeKey storetypes.StoreKey, key []byte) *stateObject {
	skey := StorageKey(key)
	// Prefer live objects if any is available
	if obj := ks.stateObjects[storeKey][skey]; obj != nil {
		return obj
	}

	// Insert into the live set
	obj := newObject(ks, storeKey, skey)
	if obj.value.IsNil() {
		// if not, means issue with db, return nil obj
		return nil
	}
	ks.setStateObject(obj)
	return obj
}

func (ks *KeestateDB) getOrNewStateObject(storeKey storetypes.StoreKey, key []byte) *stateObject {
	stateObject := ks.getStateObject(storeKey, key)
	if stateObject == nil {
		stateObject, _ = ks.createObject(storeKey, key)
	}
	return stateObject
}

func (ks *KeestateDB) setStateObject(object *stateObject) {
	if _, ok := ks.stateObjects[object.storeKey]; !ok {
		ks.stateObjects[object.storeKey] = make(map[StorageKey]*stateObject)
	}
	ks.stateObjects[object.storeKey][object.key] = object
}

func (ks *KeestateDB) createObject(storeKey storetypes.StoreKey, key []byte) (newobj, prev *stateObject) {
	prev = ks.getStateObject(storeKey, key)

	skey := StorageKey(key)
	newobj = newObject(ks, storeKey, skey)
	if prev == nil {
		ks.journal.append(createObjectChange{storeKey, skey})
	} else {
		ks.journal.append(resetObjectChange{prev: prev})
	}
	ks.setStateObject(newobj)
	if prev != nil {
		return newobj, prev
	}
	return newobj, nil
}

// SetState sets the keeper state.
func (ks *KeestateDB) SetState(storeKey storetypes.StoreKey, key, value []byte) {
	stateObject := ks.getOrNewStateObject(storeKey, key)
	if stateObject != nil {
		stateObject.SetState(value)
	}
}

// GetState retrieves a value from the given key's storage trie.
func (ks *KeestateDB) GetState(storeKey storetypes.StoreKey, key []byte) []byte {
	stateObject := ks.getStateObject(storeKey, key)
	if stateObject != nil {
		stateValue := stateObject.GetState()
		if !stateValue.IsNil() {
			return stateValue.Result()
		}
	}
	return nil
}

// GetCommittedState retrieves a value from the given key's committed storage trie.
func (ks *KeestateDB) GetCommittedState(storeKey storetypes.StoreKey, key []byte) []byte {
	stateObject := ks.getStateObject(storeKey, key)
	if stateObject != nil {
		stateValue := stateObject.GetCommittedState()
		if !stateValue.IsNil() {
			return stateValue.Result()
		}
	}
	return nil
}

// Commit all changes to a storage trie
func (ks *KeestateDB) Commit() error {
	for _, dirtyObj := range ks.journal.sortedDirties() {
		obj := ks.stateObjects[dirtyObj.storeKey][dirtyObj.key]
		for _, key := range obj.dirtyStorage.SortedKeys() {
			value := obj.dirtyStorage[key]
			origin := obj.originStorage[key]
			// Skip noop changes, persist actual changes
			if value.Eq(origin.Result()) {
				continue
			}
			obj.store(value)
		}
	}
	// no need to clean up as it will be always on fresh ctx
	return nil
}

// Snapshot returns an identifier for the current revision of the state.
func (ks *KeestateDB) Snapshot() int {
	id := ks.nextRevisionID
	ks.nextRevisionID++
	ks.validRevisions = append(ks.validRevisions, revision{id, ks.journal.length()})
	return id
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (ks *KeestateDB) RevertToSnapshot(revid int) {
	// Find the snapshot in the stack of valid snapshots.
	idx := sort.Search(len(ks.validRevisions), func(i int) bool {
		return ks.validRevisions[i].id >= revid
	})
	if idx == len(ks.validRevisions) || ks.validRevisions[idx].id != revid {
		panic(fmt.Errorf("revision id %v cannot be reverted", revid))
	}
	snapshot := ks.validRevisions[idx].journalIndex

	// Replay the journal to undo changes and remove invalidated snapshots
	ks.journal.revert(ks, snapshot)
	ks.validRevisions = ks.validRevisions[:idx]
}

// GetSdkCtx returns current cosmos sdk context
func (ks *KeestateDB) GetSdkCtx() sdk.Context {
	return ks.ctx
}
