package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid           = sdkerrors.Register(ModuleName, 1, "error invalid")
	ErrInvalidHeight     = sdkerrors.Register(ModuleName, 2, "invalid height")
	ErrEmptyUploaderAddr = sdkerrors.Register(ModuleName, 3, "missing uploader address")
	ErrEmptyReporterAddr = sdkerrors.Register(ModuleName, 4, "missing reporter address")
	ErrEmptyFileHash     = sdkerrors.Register(ModuleName, 5, "missing file hash")
	ErrEmptySenderAddr   = sdkerrors.Register(ModuleName, 6, "missing sender address")
	ErrInvalidCoins      = sdkerrors.Register(ModuleName, 7, "invalid coins")
)
