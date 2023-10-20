package types

import (
	"cosmossdk.io/errors"
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
	codeErrTotalSupplyCapHit
	codeErrInsufficientCommunityPool
)

var (
	ErrInvalid                        = errors.Register(ModuleName, codeErrInvalid, "error invalid")
	ErrUnknownAccountAddress          = errors.Register(ModuleName, codeErrUnknownAccountAddress, "account address does not exist")
	ErrOutOfIssuance                  = errors.Register(ModuleName, codeErrOutOfIssuance, "mining reward reaches the issuance limit")
	ErrWithdrawAmountInvalid          = errors.Register(ModuleName, codeErrWithdrawAmountInvalid, "withdraw amount is invalid")
	ErrMissingWalletAddress           = errors.Register(ModuleName, codeErrMissingWalletAddress, "missing wallet address")
	ErrMissingTargetAddress           = errors.Register(ModuleName, codeErrMissingTargetAddress, "missing target address")
	ErrInsufficientMatureTotal        = errors.Register(ModuleName, codeErrInsufficientMatureTotal, "insufficient mature total")
	ErrMatureEpoch                    = errors.Register(ModuleName, codeErrMatureEpoch, "the value of epoch must be positive and greater than its previous one")
	ErrMiningRewardParams             = errors.Register(ModuleName, codeErrMiningRewardParams, "invalid mining reward param")
	ErrCommunityTax                   = errors.Register(ModuleName, codeErrCommunityTax, "invalid community tax param")
	ErrInitialTotalSupply             = errors.Register(ModuleName, codeErrInitialTotalSupply, "invalid initial total supply param")
	ErrEmptyFromAddr                  = errors.Register(ModuleName, codeErrEmptyFromAddr, "missing from address")
	ErrEmptyReporterAddr              = errors.Register(ModuleName, codeErrEmptyReporterAddr, "missing reporter address")
	ErrEmptyWalletVolumes             = errors.Register(ModuleName, codeErrEmptyWalletVolumes, "wallet volumes list empty")
	ErrEpochNotPositive               = errors.Register(ModuleName, codeErrEpochNotPositive, "report epoch is not positive")
	ErrEmptyReportReference           = errors.Register(ModuleName, codeErrEmptyReportReference, "missing report reference")
	ErrReporterOwnerAddr              = errors.Register(ModuleName, codeErrReporterOwnerAddr, "invalid reporter owner address")
	ErrNegativeVolume                 = errors.Register(ModuleName, codeErrNegativeVolume, "report volume is negative")
	ErrFoundationDepositAmountInvalid = errors.Register(ModuleName, codeErrFoundationDepositAmountInvalid, "foundation deposit amount is invalid")
	ErrBLSSignatureInvalid            = errors.Register(ModuleName, codeErrBLSSignatureInvalid, "BLS signature is invalid")
	ErrBLSTxDataInvalid               = errors.Register(ModuleName, codeErrBLSTxDataInvalid, "BLS signature txData is invalid")
	ErrBLSPubkeysInvalid              = errors.Register(ModuleName, codeErrBLSPubkeysInvalid, "BLS signature pubkeys are invalid")
	ErrBLSVerifyFailed                = errors.Register(ModuleName, codeErrBLSVerifyFailed, "BLS signature verify failed")
	ErrBLSNotReachThreshold           = errors.Register(ModuleName, codeErrBLSNotReachThreshold, "BLS signed meta-nodes does not reach the threshold")
	ErrReporterAddress                = errors.Register(ModuleName, codeErrReporterAddress, "invalid reporter address")
	ErrInvalidAmount                  = errors.Register(ModuleName, codeErrInvalidAmount, "invalid amount")
	ErrCannotFindReport               = errors.Register(ModuleName, codeErrCannotFindReport, "Can not find report")
	ErrCannotFindReward               = errors.Register(ModuleName, codeErrCannotFindReward, "Can not find Pot rewards")
	ErrInvalidAddress                 = errors.Register(ModuleName, codeErrInvalidAddress, "invalid address")
	ErrInvalidDenom                   = errors.Register(ModuleName, codeErrInvalidDenom, "invalid denomination")
	ErrWithdrawFailure                = errors.Register(ModuleName, codeErrWithdrawFailure, "failure during withdraw")
	ErrFoundationDepositFailure       = errors.Register(ModuleName, codeErrFoundationDepositFailure, "failure during foundation deposit")
	ErrSlashingResourceNodeFailure    = errors.Register(ModuleName, codeErrSlashingResourceNodeFailure, "failure during slashing resource node")
	ErrVolumeReport                   = errors.Register(ModuleName, codeErrVolumeReport, "volume report failed")
	ErrLegacyWithdrawFailure          = errors.Register(ModuleName, codeErrLegacyWithdrawFailure, "failure during legacyWithdraw")
	ErrReporterAddressOrOwner         = errors.Register(ModuleName, codeErrReporterAddressOrOwner, "invalid reporter address or owner address")
	ErrTotalSupplyCapHit              = errors.Register(ModuleName, codeErrTotalSupplyCapHit, "minting not completed because total supply cap is hit")
	ErrInsufficientCommunityPool      = errors.Register(ModuleName, codeErrInsufficientCommunityPool, "burning not completed as a result of insufficient balance in community pool")
)
