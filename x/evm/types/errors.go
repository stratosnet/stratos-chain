package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidPrice = sdkerrors.Register(ModuleName, 1, "error invalid price")
	// ErrInvalidState returns an error resulting from an invalid Storage State.
	ErrInvalidState = sdkerrors.Register(ModuleName, 2, "invalid storage state")
)
