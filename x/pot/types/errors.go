package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	codeErrInvalid = uint32(iota) + 2 // NOTE: code 1 is reserved for internal errors
	codeErrUnknownAccountAddress
	codeErrOutOfIssuance
	codeErrWithdrawAmountInvalid
	codeErrMissingWalletAddress
	codeErrMissingTargetAddress
	codeErrInsufficientMatureTotal
	codeErrInsufficientFoundationAccBalance
	codeErrInsufficientUnissuedPrePayBalance
	codeErrNotTheOwner
	codeErrMatureEpoch
	codeErrEmptyFromAddr
	codeErrEmptyReporterAddr
	codeErrEmptyWalletVolumes
	codeErrEpochNotPositive
	codeErrEmptyReportReference
	codeErrEmptyReporterOwnerAddr
	codeErrNegativeVolume
	codeErrFoundationDepositAmountInvalid
	codeErrBLSSignatureInvalid
	codeErrBLSTxDataInvalid
	codeErrBLSPubkeysInvalid
	codeErrReporterAddress
	codeErrNodeStatusSuspend
	codeErrInvalidAmount
	codeErrCannotFindReport
	codeErrCannotFindReward
	codeErrInvalidAddress
)

var (
	ErrInvalid                           = sdkerrors.Register(ModuleName, codeErrInvalid, "error invalid")
	ErrUnknownAccountAddress             = sdkerrors.Register(ModuleName, codeErrUnknownAccountAddress, "account address does not exist")
	ErrOutOfIssuance                     = sdkerrors.Register(ModuleName, codeErrOutOfIssuance, "mining reward reaches the issuance limit")
	ErrWithdrawAmountInvalid             = sdkerrors.Register(ModuleName, codeErrWithdrawAmountInvalid, "withdraw amount is invalid")
	ErrMissingWalletAddress              = sdkerrors.Register(ModuleName, codeErrMissingWalletAddress, "missing wallet address")
	ErrMissingTargetAddress              = sdkerrors.Register(ModuleName, codeErrMissingTargetAddress, "missing target address")
	ErrInsufficientMatureTotal           = sdkerrors.Register(ModuleName, codeErrInsufficientMatureTotal, "insufficient mature total")
	ErrInsufficientFoundationAccBalance  = sdkerrors.Register(ModuleName, codeErrInsufficientFoundationAccBalance, "insufficient foundation account balance")
	ErrInsufficientUnissuedPrePayBalance = sdkerrors.Register(ModuleName, codeErrInsufficientUnissuedPrePayBalance, "insufficient unissued prepay balance")
	ErrNotTheOwner                       = sdkerrors.Register(ModuleName, codeErrNotTheOwner, "not the owner of the node")
	ErrMatureEpoch                       = sdkerrors.Register(ModuleName, codeErrMatureEpoch, "the value of epoch must be positive and greater than its previous one")
	ErrEmptyFromAddr                     = sdkerrors.Register(ModuleName, codeErrEmptyFromAddr, "missing from address")
	ErrEmptyReporterAddr                 = sdkerrors.Register(ModuleName, codeErrEmptyReporterAddr, "missing reporter address")
	ErrEmptyWalletVolumes                = sdkerrors.Register(ModuleName, codeErrEmptyWalletVolumes, "wallet volumes list empty")
	ErrEpochNotPositive                  = sdkerrors.Register(ModuleName, codeErrEpochNotPositive, "report epoch is not positive")
	ErrEmptyReportReference              = sdkerrors.Register(ModuleName, codeErrEmptyReportReference, "missing report reference")
	ErrEmptyReporterOwnerAddr            = sdkerrors.Register(ModuleName, codeErrEmptyReporterOwnerAddr, "missing reporter owner address")
	ErrNegativeVolume                    = sdkerrors.Register(ModuleName, codeErrNegativeVolume, "report volume is negative")
	ErrFoundationDepositAmountInvalid    = sdkerrors.Register(ModuleName, codeErrFoundationDepositAmountInvalid, "foundation deposit amount is invalid")
	ErrBLSSignatureInvalid               = sdkerrors.Register(ModuleName, codeErrBLSSignatureInvalid, "BLS signature is invalid")
	ErrBLSTxDataInvalid                  = sdkerrors.Register(ModuleName, codeErrBLSTxDataInvalid, "BLS signature txData is invalid")
	ErrBLSPubkeysInvalid                 = sdkerrors.Register(ModuleName, codeErrBLSPubkeysInvalid, "BLS signature pubkeys are invalid")
	ErrReporterAddress                   = sdkerrors.Register(ModuleName, codeErrReporterAddress, "invalid reporter address")
	ErrNodeStatusSuspend                 = sdkerrors.Register(ModuleName, codeErrNodeStatusSuspend, "node already in status expected")
	ErrInvalidAmount                     = sdkerrors.Register(ModuleName, codeErrInvalidAmount, "invalid amount")
	ErrCannotFindReport                  = sdkerrors.Register(ModuleName, codeErrCannotFindReport, "Can not find report")
	ErrCannotFindReward                  = sdkerrors.Register(ModuleName, codeErrCannotFindReward, "Can not find Pot rewards")
	ErrInvalidAddress                    = sdkerrors.Register(ModuleName, codeErrInvalidAddress, "invalid address")
)
