package statedb

import (
	"bytes"
	"sort"
	"strings"
)

type StorageKey string

func (sk StorageKey) Len() int {
	return len(sk)
}

type StorageValue struct {
	value []byte
}

func (sv *StorageValue) Result() []byte {
	return sv.value
}

func (sv *StorageValue) Eq(v []byte) bool {
	if sv.IsNil() && v == nil {
		return true
	}
	if sv.IsNil() && v != nil {
		return false
	}
	if !sv.IsNil() && v == nil {
		return false
	}
	return bytes.Equal(sv.value, v)
}

func (sv *StorageValue) IsNil() bool {
	return sv.value == nil
}

func NewStorageValue(v []byte) StorageValue {
	return StorageValue{value: v}
}

// Storage represents in-memory cache/buffer of contract storage.
type Storage map[StorageKey]StorageValue

// SortedKeys sort the keys for deterministic iteration
func (s Storage) SortedKeys() []StorageKey {
	keys := make([]StorageKey, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(string(keys[i]), string(keys[j])) < 0
	})
	return keys
}
