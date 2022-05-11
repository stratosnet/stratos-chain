package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	codeErrInvalid = uint32(iota) + 2 // NOTE: code 1 is reserved for internal errors
	codeErrInvalidHeight
	codeErrEmptyUploaderAddr
	codeErrEmptyReporterAddr
	codeErrEmptyFileHash
	codeErrEmptySenderAddr
	codeErrInvalidCoins
)

var (
	ErrInvalid           = sdkerrors.Register(ModuleName, codeErrInvalid, "error invalid")
	ErrInvalidHeight     = sdkerrors.Register(ModuleName, codeErrInvalidHeight, "invalid height")
	ErrEmptyUploaderAddr = sdkerrors.Register(ModuleName, codeErrEmptyUploaderAddr, "missing uploader address")
	ErrEmptyReporterAddr = sdkerrors.Register(ModuleName, codeErrEmptyReporterAddr, "missing reporter address")
	ErrEmptyFileHash     = sdkerrors.Register(ModuleName, codeErrEmptyFileHash, "missing file hash")
	ErrEmptySenderAddr   = sdkerrors.Register(ModuleName, codeErrEmptySenderAddr, "missing sender address")
	ErrInvalidCoins      = sdkerrors.Register(ModuleName, codeErrInvalidCoins, "invalid coins")
)
