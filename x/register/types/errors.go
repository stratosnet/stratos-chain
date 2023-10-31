package types

import (
	"cosmossdk.io/errors"
)

const (
	codeErrInvalid = uint32(iota) + 2 // NOTE: code 1 is reserved for internal errors
	codeErrInvalidNetworkAddr
	codeErrEmptyOwnerAddr
	codeErrValueNegative
	codeErrEmptyDescription
	codeErrEmptyMoniker
	codeErrEmptyResourceNodeAddr
	codeErrEmptyMetaNodeAddr
	codeErrBadDenom
	codeErrInsufficientDeposit
	codeErrResourceNodePubKeyExists
	codeErrMetaNodePubKeyExists
	codeErrNoResourceNodeFound
	codeErrNoMetaNodeFound
	codeErrNoOwnerAccountFound
	codeErrInsufficientBalance
	codeErrNodeType
	codeErrEmptyCandidateNetworkAddr
	codeErrEmptyCandidateOwnerAddr
	codeErrEmptyVoterNetworkAddr
	codeErrEmptyVoterOwnerAddr
	codeErrInvalidCandidateNetworkAddr
	codeErrInvalidCandidateOwnerAddr
	codeErrNoCandidateMetaNodeFound
	codeErrInvalidVoterNetworkAddr
	codeErrInvalidVoterOwnerAddr
	codeErrNoVoterMetaNodeFound
	codeErrSameAddr
	codeErrInvalidOwnerAddr
	codeErrInvalidBeneficiaryAddr
	codeErrInvalidVoterStatus
	codeEcoderrNoRegistrationVotePoolFound
	codeErrDuplicateVoting
	codeErrVoteExpired
	codeErrInsufficientBalanceOfNotBondedPool
	codeErrEmptyNodeNetworkAddress
	codeErrEmptyPubKey
	codeErrNoUnbondingNode
	codeErrMaxUnbondingNodeEntries
	codeErrUnbondingNode
	codeErrDepositNozRate
	codeErrRemainingNozLimit
	codeErrInvalidDepositChange
	codeErrInvalidNodeType
	codeErrUnknownAccountAddress
	codeErrUnknownPubKey
	codeErrInvalidNodeStat
	codeErrRegisterResourceNode
	codeErrRegisterMetaNode
	codeErrUnbondResourceNode
	codeErrUnbondMetaNode
	codeErrUpdateResourceNode
	codeErrUpdateMetaNode
	codeErrUpdateResourceNodeDeposit
	codeErrUpdateMetaNodeDeposit
	codeErrVoteMetaNode
	codeErrResourceNodeRegDisabled
	codeErrInvalidSuspensionStatForUnbondNode
	codeErrReporterAddress
	codeErrInvalidAmount
	codeErrReporterAddressOrOwner
	codeErrInvalidEffectiveToken
	codeErrEmitEvent
)

var (
	ErrInvalid                            = errors.Register(ModuleName, codeErrInvalid, "error invalid")
	ErrInvalidNetworkAddr                 = errors.Register(ModuleName, codeErrInvalidNetworkAddr, "invalid network address")
	ErrEmptyOwnerAddr                     = errors.Register(ModuleName, codeErrEmptyOwnerAddr, "missing owner address")
	ErrValueNegative                      = errors.Register(ModuleName, codeErrValueNegative, "value must be positive")
	ErrEmptyDescription                   = errors.Register(ModuleName, codeErrEmptyDescription, "description must be not empty")
	ErrEmptyMoniker                       = errors.Register(ModuleName, codeErrEmptyMoniker, "moniker must be not empty")
	ErrEmptyResourceNodeAddr              = errors.Register(ModuleName, codeErrEmptyResourceNodeAddr, "missing resource node address")
	ErrEmptyMetaNodeAddr                  = errors.Register(ModuleName, codeErrEmptyMetaNodeAddr, "missing Meta node address")
	ErrBadDenom                           = errors.Register(ModuleName, codeErrBadDenom, "invalid coin denomination")
	ErrInsufficientDeposit                = errors.Register(ModuleName, codeErrInsufficientDeposit, "insufficient deposit")
	ErrResourceNodePubKeyExists           = errors.Register(ModuleName, codeErrResourceNodePubKeyExists, "resource node already exist for this pubkey; must use new resource node pubkey")
	ErrMetaNodePubKeyExists               = errors.Register(ModuleName, codeErrMetaNodePubKeyExists, "meta node already exist for this pubkey; must use new meta node pubkey")
	ErrNoResourceNodeFound                = errors.Register(ModuleName, codeErrNoResourceNodeFound, "resource node does not exist")
	ErrNoMetaNodeFound                    = errors.Register(ModuleName, codeErrNoMetaNodeFound, "meta node does not exist")
	ErrNoOwnerAccountFound                = errors.Register(ModuleName, codeErrNoOwnerAccountFound, "account of owner does not exist")
	ErrInsufficientBalance                = errors.Register(ModuleName, codeErrInsufficientBalance, "insufficient balance")
	ErrNodeType                           = errors.Register(ModuleName, codeErrNodeType, "node type(s) not supported")
	ErrEmptyCandidateNetworkAddr          = errors.Register(ModuleName, codeErrEmptyCandidateNetworkAddr, "missing candidate network address")
	ErrEmptyCandidateOwnerAddr            = errors.Register(ModuleName, codeErrEmptyCandidateOwnerAddr, "missing candidate owner address")
	ErrEmptyVoterNetworkAddr              = errors.Register(ModuleName, codeErrEmptyVoterNetworkAddr, "missing voter network address")
	ErrEmptyVoterOwnerAddr                = errors.Register(ModuleName, codeErrEmptyVoterOwnerAddr, "missing voter owner address")
	ErrInvalidCandidateNetworkAddr        = errors.Register(ModuleName, codeErrInvalidCandidateNetworkAddr, "invalid candidate network address")
	ErrInvalidCandidateOwnerAddr          = errors.Register(ModuleName, codeErrInvalidCandidateOwnerAddr, "invalid candidate owner address")
	ErrNoCandidateMetaNodeFound           = errors.Register(ModuleName, codeErrNoCandidateMetaNodeFound, "candidate meta node does not exist")
	ErrInvalidVoterNetworkAddr            = errors.Register(ModuleName, codeErrInvalidVoterNetworkAddr, "invalid voter network address")
	ErrInvalidVoterOwnerAddr              = errors.Register(ModuleName, codeErrInvalidVoterOwnerAddr, "invalid voter owner address")
	ErrNoVoterMetaNodeFound               = errors.Register(ModuleName, codeErrNoVoterMetaNodeFound, "voter meta node does not exist")
	ErrSameAddr                           = errors.Register(ModuleName, codeErrSameAddr, "node address should not same as the voter address")
	ErrInvalidOwnerAddr                   = errors.Register(ModuleName, codeErrInvalidOwnerAddr, "invalid owner address")
	ErrInvalidBeneficiaryAddr             = errors.Register(ModuleName, codeErrInvalidBeneficiaryAddr, "invalid beneficiary address")
	ErrInvalidVoterStatus                 = errors.Register(ModuleName, codeErrInvalidVoterStatus, "invalid voter status")
	ErrNoRegistrationVotePoolFound        = errors.Register(ModuleName, codeEcoderrNoRegistrationVotePoolFound, "registration pool does not exist")
	ErrDuplicateVoting                    = errors.Register(ModuleName, codeErrDuplicateVoting, "duplicate voting")
	ErrVoteExpired                        = errors.Register(ModuleName, codeErrVoteExpired, "vote expired")
	ErrInsufficientBalanceOfNotBondedPool = errors.Register(ModuleName, codeErrInsufficientBalanceOfNotBondedPool, "insufficient balance of not bonded pool")
	ErrEmptyNodeNetworkAddress            = errors.Register(ModuleName, codeErrEmptyNodeNetworkAddress, "missing node network address")
	ErrEmptyPubKey                        = errors.Register(ModuleName, codeErrEmptyPubKey, "missing public key")
	ErrNoUnbondingNode                    = errors.Register(ModuleName, codeErrNoUnbondingNode, "no unbonding node found")
	ErrMaxUnbondingNodeEntries            = errors.Register(ModuleName, codeErrMaxUnbondingNodeEntries, "too many unbonding node entries for networkAddr tuple")
	ErrUnbondingNode                      = errors.Register(ModuleName, codeErrUnbondingNode, "changes cannot be made to an unbonding node")
	ErrDepositNozRate                     = errors.Register(ModuleName, codeErrDepositNozRate, "deposit noz rate must be positive")
	ErrRemainingNozLimit                  = errors.Register(ModuleName, codeErrRemainingNozLimit, "remaining Noz Limit must be non-negative")
	ErrInvalidDepositChange               = errors.Register(ModuleName, codeErrInvalidDepositChange, "invalid change for deposit")
	ErrInvalidNodeType                    = errors.Register(ModuleName, codeErrInvalidNodeType, "invalid node type")
	ErrUnknownAccountAddress              = errors.Register(ModuleName, codeErrUnknownAccountAddress, "account address does not exist")
	ErrUnknownPubKey                      = errors.Register(ModuleName, codeErrUnknownPubKey, "unknown pubKey ")
	ErrInvalidNodeStat                    = errors.Register(ModuleName, codeErrInvalidNodeStat, "invalid node status")
	ErrRegisterResourceNode               = errors.Register(ModuleName, codeErrRegisterResourceNode, "failed to register resource node")
	ErrRegisterMetaNode                   = errors.Register(ModuleName, codeErrRegisterMetaNode, "failed to register meta node")
	ErrUnbondResourceNode                 = errors.Register(ModuleName, codeErrUnbondResourceNode, "failed to unbond resource node")
	ErrUnbondMetaNode                     = errors.Register(ModuleName, codeErrUnbondMetaNode, "failed to unbond meta node")
	ErrUpdateResourceNode                 = errors.Register(ModuleName, codeErrUpdateResourceNode, "failed to update resource node")
	ErrUpdateMetaNode                     = errors.Register(ModuleName, codeErrUpdateMetaNode, "failed to update meta node")
	ErrUpdateResourceNodeDeposit          = errors.Register(ModuleName, codeErrUpdateResourceNodeDeposit, "failed to update deposit for resource node")
	ErrUpdateMetaNodeDeposit              = errors.Register(ModuleName, codeErrUpdateMetaNodeDeposit, "failed to update deposit for meta node")
	ErrVoteMetaNode                       = errors.Register(ModuleName, codeErrVoteMetaNode, "failed to vote meta node")
	ErrResourceNodeRegDisabled            = errors.Register(ModuleName, codeErrResourceNodeRegDisabled, "resource node registration is disabled")
	ErrInvalidSuspensionStatForUnbondNode = errors.Register(ModuleName, codeErrInvalidSuspensionStatForUnbondNode, "cannot unbond a suspended node")
	ErrReporterAddress                    = errors.Register(ModuleName, codeErrReporterAddress, "invalid reporter address")
	ErrInvalidAmount                      = errors.Register(ModuleName, codeErrInvalidAmount, "invalid amount")
	ErrReporterAddressOrOwner             = errors.Register(ModuleName, codeErrReporterAddressOrOwner, "invalid reporter address or owner address")
	ErrInvalidEffectiveToken              = errors.Register(ModuleName, codeErrInvalidEffectiveToken, "invalid effective token")
	ErrEmitEvent                          = errors.Register(ModuleName, codeErrEmitEvent, "failed to emit event")
)
