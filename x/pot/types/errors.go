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
	codeErrMatureEpoch
	codeErrMiningRewardParams
	codeErrCommunityTax
	codeErrInitialTotalSupply
	codeErrEmptyFromAddr
	codeErrEmptyReporterAddr
	codeErrEmptyWalletVolumes
	codeErrEpochNotPositive
	codeErrEmptyReportReference
	codeErrReporterOwnerAddr
	codeErrNegativeVolume
	codeErrFoundationDepositAmountInvalid
	codeErrBLSSignatureInvalid
	codeErrBLSTxDataInvalid
	codeErrBLSPubkeysInvalid
	codeErrBLSVerifyFailed
	codeErrBLSNotReachThreshold
	codeErrReporterAddress
	codeErrInvalidAmount
	codeErrCannotFindReport
	codeErrCannotFindReward
	codeErrInvalidAddress
	codeErrInvalidDenom
	codeErrWithdrawFailure
	codeErrFoundationDepositFailure
	codeErrSlashingResourceNodeFailure
	codeErrVolumeReport
	codeErrLegacyWithdrawFailure
	codeErrReporterAddressOrOwner
)

var (
	ErrInvalid                        = sdkerrors.Register(ModuleName, codeErrInvalid, "error invalid")
	ErrUnknownAccountAddress          = sdkerrors.Register(ModuleName, codeErrUnknownAccountAddress, "account address does not exist")
	ErrOutOfIssuance                  = sdkerrors.Register(ModuleName, codeErrOutOfIssuance, "mining reward reaches the issuance limit")
	ErrWithdrawAmountInvalid          = sdkerrors.Register(ModuleName, codeErrWithdrawAmountInvalid, "withdraw amount is invalid")
	ErrMissingWalletAddress           = sdkerrors.Register(ModuleName, codeErrMissingWalletAddress, "missing wallet address")
	ErrMissingTargetAddress           = sdkerrors.Register(ModuleName, codeErrMissingTargetAddress, "missing target address")
	ErrInsufficientMatureTotal        = sdkerrors.Register(ModuleName, codeErrInsufficientMatureTotal, "insufficient mature total")
	ErrMatureEpoch                    = sdkerrors.Register(ModuleName, codeErrMatureEpoch, "the value of epoch must be positive and greater than its previous one")
	ErrMiningRewardParams             = sdkerrors.Register(ModuleName, codeErrMiningRewardParams, "invalid mining reward param")
	ErrCommunityTax                   = sdkerrors.Register(ModuleName, codeErrCommunityTax, "invalid community tax param")
	ErrInitialTotalSupply             = sdkerrors.Register(ModuleName, codeErrInitialTotalSupply, "invalid initial total supply param")
	ErrEmptyFromAddr                  = sdkerrors.Register(ModuleName, codeErrEmptyFromAddr, "missing from address")
	ErrEmptyReporterAddr              = sdkerrors.Register(ModuleName, codeErrEmptyReporterAddr, "missing reporter address")
	ErrEmptyWalletVolumes             = sdkerrors.Register(ModuleName, codeErrEmptyWalletVolumes, "wallet volumes list empty")
	ErrEpochNotPositive               = sdkerrors.Register(ModuleName, codeErrEpochNotPositive, "report epoch is not positive")
	ErrEmptyReportReference           = sdkerrors.Register(ModuleName, codeErrEmptyReportReference, "missing report reference")
	ErrReporterOwnerAddr              = sdkerrors.Register(ModuleName, codeErrReporterOwnerAddr, "invalid reporter owner address")
	ErrNegativeVolume                 = sdkerrors.Register(ModuleName, codeErrNegativeVolume, "report volume is negative")
	ErrFoundationDepositAmountInvalid = sdkerrors.Register(ModuleName, codeErrFoundationDepositAmountInvalid, "foundation deposit amount is invalid")
	ErrBLSSignatureInvalid            = sdkerrors.Register(ModuleName, codeErrBLSSignatureInvalid, "BLS signature is invalid")
	ErrBLSTxDataInvalid               = sdkerrors.Register(ModuleName, codeErrBLSTxDataInvalid, "BLS signature txData is invalid")
	ErrBLSPubkeysInvalid              = sdkerrors.Register(ModuleName, codeErrBLSPubkeysInvalid, "BLS signature pubkeys are invalid")
	ErrBLSVerifyFailed                = sdkerrors.Register(ModuleName, codeErrBLSVerifyFailed, "BLS signature verify failed")
	ErrBLSNotReachThreshold           = sdkerrors.Register(ModuleName, codeErrBLSNotReachThreshold, "BLS signed meta-nodes does not reach the threshold")
	ErrReporterAddress                = sdkerrors.Register(ModuleName, codeErrReporterAddress, "invalid reporter address")
	ErrInvalidAmount                  = sdkerrors.Register(ModuleName, codeErrInvalidAmount, "invalid amount")
	ErrCannotFindReport               = sdkerrors.Register(ModuleName, codeErrCannotFindReport, "Can not find report")
	ErrCannotFindReward               = sdkerrors.Register(ModuleName, codeErrCannotFindReward, "Can not find Pot rewards")
	ErrInvalidAddress                 = sdkerrors.Register(ModuleName, codeErrInvalidAddress, "invalid address")
	ErrInvalidDenom                   = sdkerrors.Register(ModuleName, codeErrInvalidDenom, "invalid denomination")
	ErrWithdrawFailure                = sdkerrors.Register(ModuleName, codeErrWithdrawFailure, "failure during withdraw")
	ErrFoundationDepositFailure       = sdkerrors.Register(ModuleName, codeErrFoundationDepositFailure, "failure during foundation deposit")
	ErrSlashingResourceNodeFailure    = sdkerrors.Register(ModuleName, codeErrSlashingResourceNodeFailure, "failure during slashing resource node")
	ErrVolumeReport                   = sdkerrors.Register(ModuleName, codeErrVolumeReport, "volume report failed")
	ErrLegacyWithdrawFailure          = sdkerrors.Register(ModuleName, codeErrLegacyWithdrawFailure, "failure during legacyWithdraw")
	ErrReporterAddressOrOwner         = sdkerrors.Register(ModuleName, codeErrReporterAddressOrOwner, "invalid reporter address or owner address")
)
