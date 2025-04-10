package testutil

import (
	"errors"
)

type MockIterator struct {
	keys   [][]byte
	values [][]byte
	index  int
	closed bool
}

func NewMockIterator(data map[string][]byte) *MockIterator {
	keys := make([][]byte, 0, len(data))
	values := make([][]byte, 0, len(data))
	for k, v := range data {
		keys = append(keys, []byte(k))
		values = append(values, v)
	}
	return &MockIterator{
		keys:   keys,
		values: values,
		index:  0,
		closed: false,
	}
}

func (it *MockIterator) Domain() (start []byte, end []byte) {
	if len(it.keys) == 0 {
		return nil, nil
	}
	return it.keys[0], it.keys[len(it.keys)-1]
}

func (it *MockIterator) Valid() bool {
	return !it.closed && it.index >= 0 && it.index < len(it.keys)
}

func (it *MockIterator) Next() {
	if !it.Valid() {
		panic("Next() called on an invalid iterator")
	}
	it.index++
}

func (it *MockIterator) Key() []byte {
	if !it.Valid() {
		panic("Key() called on an invalid iterator")
	}
	return it.keys[it.index]
}

func (it *MockIterator) Value() []byte {
	if !it.Valid() {
		panic("Value() called on an invalid iterator")
	}
	return it.values[it.index]
}

func (it *MockIterator) Error() error {
	return nil
}

func (it *MockIterator) Close() error {
	if it.closed {
		return errors.New("iterator already closed")
	}
	it.closed = true
	return nil
}
