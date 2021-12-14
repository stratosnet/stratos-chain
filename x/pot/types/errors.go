package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid                           = sdkerrors.Register(ModuleName, 1, "error invalid")
	ErrUnknownAccountAddress             = sdkerrors.Register(ModuleName, 2, "account address does not exist")
	ErrOutOfIssuance                     = sdkerrors.Register(ModuleName, 3, "mining reward reaches the issuance limit")
	ErrWithdrawAmountInvalid             = sdkerrors.Register(ModuleName, 4, "withdraw amount is invalid")
	ErrMissingWalletAddress              = sdkerrors.Register(ModuleName, 5, "missing wallet address")
	ErrMissingTargetAddress              = sdkerrors.Register(ModuleName, 6, "missing target address")
	ErrInsufficientMatureTotal           = sdkerrors.Register(ModuleName, 7, "insufficient mature total")
	ErrInsufficientFoundationAccBalance  = sdkerrors.Register(ModuleName, 8, "insufficient foundation account balance")
	ErrInsufficientUnissuedPrePayBalance = sdkerrors.Register(ModuleName, 9, "insufficient unissued prepay balance")
	ErrNotTheOwner                       = sdkerrors.Register(ModuleName, 10, "not the owner of the node")
	ErrMatureEpoch                       = sdkerrors.Register(ModuleName, 11, "the value of epoch must be positive and greater than its previous one")
	ErrEmptyFromAddr                     = sdkerrors.Register(ModuleName, 12, "missing from address")
	ErrEmptyReporterAddr                 = sdkerrors.Register(ModuleName, 13, "missing reporter address")
	ErrEmptyWalletVolumes                = sdkerrors.Register(ModuleName, 14, "wallet volumes list empty")
	ErrEpochNotPositive                  = sdkerrors.Register(ModuleName, 15, "report epoch is not positive")
	ErrEmptyReportReference              = sdkerrors.Register(ModuleName, 16, "missing report reference")
	ErrEmptyReporterOwnerAddr            = sdkerrors.Register(ModuleName, 17, "missing reporter owner address")
	ErrNegativeVolume                    = sdkerrors.Register(ModuleName, 18, "report volume is negative")
	ErrFoundationDepositAmountInvalid    = sdkerrors.Register(ModuleName, 19, "foundation deposit amount is invalid")
)
