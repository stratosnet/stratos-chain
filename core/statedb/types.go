package statedb

import "bytes"

const StorageKeyLength = 256

type StorageKey [StorageKeyLength]byte

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (sk *StorageKey) SetBytes(b []byte) {
	if len(b) > len(sk) {
		b = b[len(b)-StorageKeyLength:]
	}

	copy(sk[StorageKeyLength-len(b):], b)
}

// BytesToStorageKey sets b to StorageKey.
// If b is larger than len(h), b will be cropped from the left.
func BytesToStorageKey(b []byte) StorageKey {
	var sk StorageKey
	sk.SetBytes(b)
	return sk
}

type StorageValue struct {
	value []byte
}

func (sv *StorageValue) Result() []byte {
	return sv.value
}

func (sv *StorageValue) Eq(v []byte) bool {
	if sv.value == nil {
		return false
	}
	return bytes.Equal(sv.value, v)
}

func NewStorageValue(v []byte) *StorageValue {
	return &StorageValue{value: v}
}

// Storage represents in-memory cache/buffer of contract storage.
type Storage map[StorageKey]*StorageValue
