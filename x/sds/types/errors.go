package types

import (
	"cosmossdk.io/errors"
)

const (
	codeErrInvalid = uint32(iota) + 2 // NOTE: code 1 is reserved for internal errors
	codeErrInvalidHeight
	codeErrEmptyUploaderAddr
	codeErrEmptyReporters
	codeErrEmptyFileHash
	codeErrInvalidFileHash
	codeErrNoFileFound
	codeErrInvalidDenom
	codeErrPrepayFailure
	codeErrReporterAddressOrOwner
	codeErrInvalidSenderAddr
	codeErrInvalidBeneficiaryAddr
	codeErrOzoneLimitNotEnough
	codeErrEmitEvent
)

var (
	ErrInvalid                = errors.Register(ModuleName, codeErrInvalid, "error invalid")
	ErrInvalidHeight          = errors.Register(ModuleName, codeErrInvalidHeight, "invalid height")
	ErrEmptyUploaderAddr      = errors.Register(ModuleName, codeErrEmptyUploaderAddr, "missing uploader address")
	ErrEmptyReporters         = errors.Register(ModuleName, codeErrEmptyReporters, "missing reporters")
	ErrEmptyFileHash          = errors.Register(ModuleName, codeErrEmptyFileHash, "missing file hash")
	ErrInvalidFileHash        = errors.Register(ModuleName, codeErrInvalidFileHash, "invalid file hash")
	ErrNoFileFound            = errors.Register(ModuleName, codeErrNoFileFound, "file does not exist")
	ErrInvalidDenom           = errors.Register(ModuleName, codeErrInvalidDenom, "invalid denomination")
	ErrPrepayFailure          = errors.Register(ModuleName, codeErrPrepayFailure, "failure during prepay")
	ErrReporterAddressOrOwner = errors.Register(ModuleName, codeErrReporterAddressOrOwner, "invalid reporter address or owner address")
	ErrInvalidSenderAddr      = errors.Register(ModuleName, codeErrInvalidSenderAddr, "invalid sender address")
	ErrInvalidBeneficiaryAddr = errors.Register(ModuleName, codeErrInvalidBeneficiaryAddr, "invalid beneficiary address")
	ErrOzoneLimitNotEnough    = errors.Register(ModuleName, codeErrOzoneLimitNotEnough, "not enough remaining ozone limit to complete prepay")
	ErrEmitEvent              = errors.Register(ModuleName, codeErrEmitEvent, "failed to emit event")
)
