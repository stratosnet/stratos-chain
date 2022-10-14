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
	codeErrInvalidCoins
	codeErrInvalidFileHash
	codeErrInvalidDenom
	codeErrPrepayFailure
)

var (
	ErrInvalid           = sdkerrors.Register(ModuleName, codeErrInvalid, "error invalid")
	ErrInvalidHeight     = sdkerrors.Register(ModuleName, codeErrInvalidHeight, "invalid height")
	ErrEmptyUploaderAddr = sdkerrors.Register(ModuleName, codeErrEmptyUploaderAddr, "missing uploader address")
	ErrEmptyReporterAddr = sdkerrors.Register(ModuleName, codeErrEmptyReporterAddr, "missing reporter address")
	ErrEmptyFileHash     = sdkerrors.Register(ModuleName, codeErrEmptyFileHash, "missing file hash")
	ErrInvalidCoins      = sdkerrors.Register(ModuleName, codeErrInvalidCoins, "invalid coins")
	ErrInvalidFileHash   = sdkerrors.Register(ModuleName, codeErrInvalidFileHash, "invalid file hash")
	ErrInvalidDenom      = sdkerrors.Register(ModuleName, codeErrInvalidDenom, "invalid denomination")
	ErrPrepayFailure     = sdkerrors.Register(ModuleName, codeErrPrepayFailure, "failure during prepay")
)
