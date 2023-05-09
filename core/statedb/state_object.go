package statedb

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// stateObject is the state of an acount
type stateObject struct {
	ctx sdk.Context
	db  *KeestateDB

	storeKey sdk.StoreKey
	key      StorageKey
	value    *StorageValue

	// state storage
	originStorage Storage
	dirtyStorage  Storage
}

// newObject creates a state object.
func newObject(db *KeestateDB, storeKey sdk.StoreKey, key StorageKey) *stateObject {
	so := &stateObject{
		ctx:           db.ctx,
		db:            db,
		storeKey:      storeKey,
		key:           key,
		originStorage: make(Storage),
		dirtyStorage:  make(Storage),
	}
	if value := so.read(key); value != nil {
		so.value = value
	}
	return so
}

func (s *stateObject) read(key StorageKey) *StorageValue {
	value := s.ctx.KVStore(s.storeKey).Get(key[:])
	if value != nil {
		return NewStorageValue(value)
	}
	return nil
}

// GetCommittedState query the committed state
func (s *stateObject) GetCommittedState() *StorageValue {
	if value, cached := s.originStorage[s.key]; cached {
		return value
	}
	// If no live objects are available, load it from indexdb
	value := s.read(s.key)
	if value == nil {
		return nil
	}
	s.originStorage[s.key] = value
	return value
}

// GetState query the current state (including dirty state)
func (s *stateObject) GetState() *StorageValue {
	if value, dirty := s.dirtyStorage[s.key]; dirty {
		return value
	}
	return s.GetCommittedState()
}

// SetState sets the contract state
func (s *stateObject) SetState(value []byte) {
	// If the new value is the same as old, don't set
	prev := s.GetState()
	if prev.Eq(value) {
		return
	}
	// New value is different, update and journal the change
	s.db.journal.append(storageChange{
		storeKey: s.storeKey,
		key:      s.key,
		prevalue: prev.Result(),
	})
	s.setState(value)
}

func (s *stateObject) setState(value []byte) {
	s.dirtyStorage[s.key] = NewStorageValue(value)
}
