package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid                           = sdkerrors.Register(ModuleName, 1, "error invalid")
	ErrUnknownAccountAddress             = sdkerrors.Register(ModuleName, 2, "account address does not exist")
	ErrOutOfIssuance                     = sdkerrors.Register(ModuleName, 3, "mining reward reaches the issuance limit")
	ErrWithdrawAmountNotPositive         = sdkerrors.Register(ModuleName, 4, "withdraw amount is not positive")
	ErrMissingNodeAddress                = sdkerrors.Register(ModuleName, 5, "missing node address")
	ErrMissingOwnerAddress               = sdkerrors.Register(ModuleName, 6, "missing owner address")
	ErrInsufficientMatureTotal           = sdkerrors.Register(ModuleName, 7, "insufficient mature total")
	ErrInsufficientFoundationAccBalance  = sdkerrors.Register(ModuleName, 8, "insufficient foundation account balance")
	ErrInsufficientUnissuedPrePayBalance = sdkerrors.Register(ModuleName, 9, "insufficient unissued prepay balance")
	ErrNotTheOwner                       = sdkerrors.Register(ModuleName, 10, "not the owner of the node")
	ErrMatureEpoch                       = sdkerrors.Register(ModuleName, 11, "the value of epoch must be positive and greater than its previous one")
	ErrEmptyFromAddr                     = sdkerrors.Register(ModuleName, 12, "missing from address")
	ErrEmptyReporterAddr                 = sdkerrors.Register(ModuleName, 13, "missing reporter address")
	ErrEmptyNodesVolume                  = sdkerrors.Register(ModuleName, 14, "nodes volume list empty")
	ErrEpochNotPositive                  = sdkerrors.Register(ModuleName, 15, "report epoch is not positive")
	ErrEmptyReportReference              = sdkerrors.Register(ModuleName, 16, "missing report reference")
	ErrEmptyReporterOwnerAddr            = sdkerrors.Register(ModuleName, 17, "missing reporter owner address")
	ErrNegativeVolume                    = sdkerrors.Register(ModuleName, 18, "report volume is negative")
)
