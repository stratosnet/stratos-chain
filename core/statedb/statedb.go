package statedb

import (
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

	stateObjects map[StorageKey]*stateObject
}

// New creates a new state from a given trie.
func New(ctx sdk.Context) *KeestateDB {
	return &KeestateDB{
		ctx:          ctx,
		journal:      newJournal(),
		stateObjects: make(map[StorageKey]*stateObject),
	}
}

func (ks *KeestateDB) getStateObject(storeKey sdk.StoreKey, key []byte) *stateObject {
	skey := BytesToStorageKey(key)
	// Prefer live objects if any is available
	if obj := ks.stateObjects[skey]; obj != nil {
		return obj
	}

	// Insert into the live set
	obj := newObject(ks, storeKey, skey)
	if obj.value == nil {
		// if not, means issue with db, return nil obj
		return nil
	}
	ks.setStateObject(obj)
	return obj
}

func (ks *KeestateDB) getOrNewStateObject(storeKey sdk.StoreKey, key []byte) *stateObject {
	stateObject := ks.getStateObject(storeKey, key)
	if stateObject == nil {
		stateObject, _ = ks.createObject(storeKey, key)
	}
	return stateObject
}

func (ks *KeestateDB) setStateObject(object *stateObject) {
	ks.stateObjects[object.key] = object
}

func (ks *KeestateDB) createObject(storeKey sdk.StoreKey, key []byte) (newobj, prev *stateObject) {
	prev = ks.getStateObject(storeKey, key)

	skey := BytesToStorageKey(key)
	newobj = newObject(ks, storeKey, skey)
	if prev == nil {
		ks.journal.append(createObjectChange{storeKey, BytesToStorageKey(key)})
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
func (ks *KeestateDB) SetState(storeKey sdk.StoreKey, key, value []byte) {
	stateObject := ks.getOrNewStateObject(storeKey, key)
	if stateObject != nil {
		stateObject.SetState(value)
	}
}

// GetState retrieves a value from the given key's storage trie.
func (ks *KeestateDB) GetState(storeKey sdk.StoreKey, key []byte) []byte {
	stateObject := ks.getStateObject(storeKey, key)
	if stateObject != nil {
		stateValue := stateObject.GetState()
		if stateValue != nil {
			return stateValue.Result()
		}
	}
	return nil
}

// GetCommittedState retrieves a value from the given key's committed storage trie.
func (ks *KeestateDB) GetCommittedState(storeKey sdk.StoreKey, key []byte) []byte {
	stateObject := ks.getStateObject(storeKey, key)
	if stateObject != nil {
		stateValue := stateObject.GetCommittedState()
		if stateValue != nil {
			return stateValue.Result()
		}
	}
	return nil
}
