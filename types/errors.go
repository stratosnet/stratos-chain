package types

import (
	"cosmossdk.io/errors"
)

const (
	// RootCodespace is the codespace for all errors defined in this package
	RootCodespace = "stratos"
)

// NOTE: We can't use 1 since that error code is reserved for internal errors.

var (
	// ErrInvalidValue returns an error resulting from an invalid value.
	ErrInvalidValue = errors.Register(RootCodespace, 2, "invalid value")

	// ErrInvalidChainID returns an error resulting from an invalid chain ID.
	ErrInvalidChainID = errors.Register(RootCodespace, 3, "invalid chain ID")

	// ErrMarshalBigInt returns an error resulting from marshaling a big.Int to a string.
	ErrMarshalBigInt = errors.Register(RootCodespace, 5, "cannot marshal big.Int to string")

	// ErrUnmarshalBigInt returns an error resulting from unmarshaling a big.Int from a string.
	ErrUnmarshalBigInt = errors.Register(RootCodespace, 6, "cannot unmarshal big.Int from string")
)
