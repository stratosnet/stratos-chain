package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	codeErrInvalid = uint32(iota) + 2 // NOTE: code 1 is reserved for internal errors
	codeErrInvalidHeight
	codeErrEmptyUploaderAddr
	codeErrEmptyReporterAddr
	codeErrEmptyReporters
	codeErrEmptyFileHash
	codeErrInvalidFileHash
	codeErrNoFileFound
	codeErrInvalidDenom
	codeErrPrepayFailure
	codeErrInvalidSenderAddr
	codeErrInvalidBeneficiaryAddr
	codeErrReporterAddressOrOwner
)

var (
	ErrInvalid                = sdkerrors.Register(ModuleName, codeErrInvalid, "error invalid")
	ErrInvalidHeight          = sdkerrors.Register(ModuleName, codeErrInvalidHeight, "invalid height")
	ErrEmptyUploaderAddr      = sdkerrors.Register(ModuleName, codeErrEmptyUploaderAddr, "missing uploader address")
	ErrEmptyReporterAddr      = sdkerrors.Register(ModuleName, codeErrEmptyReporterAddr, "missing reporter address")
	ErrEmptyReporters         = sdkerrors.Register(ModuleName, codeErrEmptyReporters, "missing reporters")
	ErrEmptyFileHash          = sdkerrors.Register(ModuleName, codeErrEmptyFileHash, "missing file hash")
	ErrInvalidFileHash        = sdkerrors.Register(ModuleName, codeErrInvalidFileHash, "invalid file hash")
	ErrNoFileFound            = sdkerrors.Register(ModuleName, codeErrNoFileFound, "file does not exist")
	ErrInvalidDenom           = sdkerrors.Register(ModuleName, codeErrInvalidDenom, "invalid denomination")
	ErrPrepayFailure          = sdkerrors.Register(ModuleName, codeErrPrepayFailure, "failure during prepay")
	ErrInvalidSenderAddr      = sdkerrors.Register(ModuleName, codeErrInvalidSenderAddr, "invalid sender address")
	ErrInvalidBeneficiaryAddr = sdkerrors.Register(ModuleName, codeErrInvalidBeneficiaryAddr, "invalid beneficiary address")
	ErrReporterAddressOrOwner = sdkerrors.Register(ModuleName, codeErrReporterAddressOrOwner, "invalid reporter address or owner address")
)
