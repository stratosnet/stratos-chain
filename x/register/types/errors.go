package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	codeErrStakeNozRate
	codeErrRemainingNozLimit
	codeErrInvalidStakeChange
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
	codeErrUpdateResourceNodeStake
	codeErrUpdateMetaNodeStake
	codeErrVoteMetaNode
	codeErrResourceNodeRegDisabled
	codeErrInvalidSuspensionStatForUnbondNode
	codeErrReporterAddress
	codeErrInvalidAmount
	codeErrReporterAddressOrOwner
	codeErrReporterNotReachThreshold
)

var (
	ErrInvalid                            = sdkerrors.Register(ModuleName, codeErrInvalid, "error invalid")
	ErrInvalidNetworkAddr                 = sdkerrors.Register(ModuleName, codeErrInvalidNetworkAddr, "invalid network address")
	ErrEmptyOwnerAddr                     = sdkerrors.Register(ModuleName, codeErrEmptyOwnerAddr, "missing owner address")
	ErrValueNegative                      = sdkerrors.Register(ModuleName, codeErrValueNegative, "value must be positive")
	ErrEmptyDescription                   = sdkerrors.Register(ModuleName, codeErrEmptyDescription, "description must be not empty")
	ErrEmptyMoniker                       = sdkerrors.Register(ModuleName, codeErrEmptyMoniker, "moniker must be not empty")
	ErrEmptyResourceNodeAddr              = sdkerrors.Register(ModuleName, codeErrEmptyResourceNodeAddr, "missing resource node address")
	ErrEmptyMetaNodeAddr                  = sdkerrors.Register(ModuleName, codeErrEmptyMetaNodeAddr, "missing Meta node address")
	ErrBadDenom                           = sdkerrors.Register(ModuleName, codeErrBadDenom, "invalid coin denomination")
	ErrResourceNodePubKeyExists           = sdkerrors.Register(ModuleName, codeErrResourceNodePubKeyExists, "resource node already exist for this pubkey; must use new resource node pubkey")
	ErrMetaNodePubKeyExists               = sdkerrors.Register(ModuleName, codeErrMetaNodePubKeyExists, "meta node already exist for this pubkey; must use new meta node pubkey")
	ErrNoResourceNodeFound                = sdkerrors.Register(ModuleName, codeErrNoResourceNodeFound, "resource node does not exist")
	ErrNoMetaNodeFound                    = sdkerrors.Register(ModuleName, codeErrNoMetaNodeFound, "meta node does not exist")
	ErrNoOwnerAccountFound                = sdkerrors.Register(ModuleName, codeErrNoOwnerAccountFound, "account of owner does not exist")
	ErrInsufficientBalance                = sdkerrors.Register(ModuleName, codeErrInsufficientBalance, "insufficient balance")
	ErrNodeType                           = sdkerrors.Register(ModuleName, codeErrNodeType, "node type(s) not supported")
	ErrEmptyCandidateNetworkAddr          = sdkerrors.Register(ModuleName, codeErrEmptyCandidateNetworkAddr, "missing candidate network address")
	ErrEmptyCandidateOwnerAddr            = sdkerrors.Register(ModuleName, codeErrEmptyCandidateOwnerAddr, "missing candidate owner address")
	ErrEmptyVoterNetworkAddr              = sdkerrors.Register(ModuleName, codeErrEmptyVoterNetworkAddr, "missing voter network address")
	ErrEmptyVoterOwnerAddr                = sdkerrors.Register(ModuleName, codeErrEmptyVoterOwnerAddr, "missing voter owner address")
	ErrInvalidCandidateNetworkAddr        = sdkerrors.Register(ModuleName, codeErrInvalidCandidateNetworkAddr, "invalid candidate network address")
	ErrInvalidCandidateOwnerAddr          = sdkerrors.Register(ModuleName, codeErrInvalidCandidateOwnerAddr, "invalid candidate owner address")
	ErrNoCandidateMetaNodeFound           = sdkerrors.Register(ModuleName, codeErrNoCandidateMetaNodeFound, "candidate meta node does not exist")
	ErrInvalidVoterNetworkAddr            = sdkerrors.Register(ModuleName, codeErrInvalidVoterNetworkAddr, "invalid voter network address")
	ErrInvalidVoterOwnerAddr              = sdkerrors.Register(ModuleName, codeErrInvalidVoterOwnerAddr, "invalid voter owner address")
	ErrNoVoterMetaNodeFound               = sdkerrors.Register(ModuleName, codeErrNoVoterMetaNodeFound, "voter meta node does not exist")
	ErrSameAddr                           = sdkerrors.Register(ModuleName, codeErrSameAddr, "node address should not same as the voter address")
	ErrInvalidOwnerAddr                   = sdkerrors.Register(ModuleName, codeErrInvalidOwnerAddr, "invalid owner address")
	ErrInvalidVoterStatus                 = sdkerrors.Register(ModuleName, codeErrInvalidVoterStatus, "invalid voter status")
	ErrNoRegistrationVotePoolFound        = sdkerrors.Register(ModuleName, codeEcoderrNoRegistrationVotePoolFound, "registration pool does not exist")
	ErrDuplicateVoting                    = sdkerrors.Register(ModuleName, codeErrDuplicateVoting, "duplicate voting")
	ErrVoteExpired                        = sdkerrors.Register(ModuleName, codeErrVoteExpired, "vote expired")
	ErrInsufficientBalanceOfNotBondedPool = sdkerrors.Register(ModuleName, codeErrInsufficientBalanceOfNotBondedPool, "insufficient balance of not bonded pool")
	ErrEmptyNodeNetworkAddress            = sdkerrors.Register(ModuleName, codeErrEmptyNodeNetworkAddress, "missing node network address")
	ErrEmptyPubKey                        = sdkerrors.Register(ModuleName, codeErrEmptyPubKey, "missing public key")
	ErrNoUnbondingNode                    = sdkerrors.Register(ModuleName, codeErrNoUnbondingNode, "no unbonding node found")
	ErrMaxUnbondingNodeEntries            = sdkerrors.Register(ModuleName, codeErrMaxUnbondingNodeEntries, "too many unbonding node entries for networkAddr tuple")
	ErrUnbondingNode                      = sdkerrors.Register(ModuleName, codeErrUnbondingNode, "changes cannot be made to an unbonding node")
	ErrStakeNozRate                       = sdkerrors.Register(ModuleName, codeErrStakeNozRate, "stake noz rate must be positive")
	ErrRemainingNozLimit                  = sdkerrors.Register(ModuleName, codeErrRemainingNozLimit, "remaining Noz Limit must be non-negative")
	ErrInvalidStakeChange                 = sdkerrors.Register(ModuleName, codeErrInvalidStakeChange, "invalid change for stake")
	ErrInvalidNodeType                    = sdkerrors.Register(ModuleName, codeErrInvalidNodeType, "invalid node type")
	ErrUnknownAccountAddress              = sdkerrors.Register(ModuleName, codeErrUnknownAccountAddress, "account address does not exist")
	ErrUnknownPubKey                      = sdkerrors.Register(ModuleName, codeErrUnknownPubKey, "unknown pubKey ")
	ErrInvalidNodeStat                    = sdkerrors.Register(ModuleName, codeErrInvalidNodeStat, "invalid node status")
	ErrRegisterResourceNode               = sdkerrors.Register(ModuleName, codeErrRegisterResourceNode, "failed to register resource node")
	ErrRegisterMetaNode                   = sdkerrors.Register(ModuleName, codeErrRegisterMetaNode, "failed to register meta node")
	ErrUnbondResourceNode                 = sdkerrors.Register(ModuleName, codeErrUnbondResourceNode, "failed to unbond resource node")
	ErrUnbondMetaNode                     = sdkerrors.Register(ModuleName, codeErrUnbondMetaNode, "failed to unbond meta node")
	ErrUpdateResourceNode                 = sdkerrors.Register(ModuleName, codeErrUpdateResourceNode, "failed to update resource node")
	ErrUpdateMetaNode                     = sdkerrors.Register(ModuleName, codeErrUpdateMetaNode, "failed to update meta node")
	ErrUpdateResourceNodeStake            = sdkerrors.Register(ModuleName, codeErrUpdateResourceNodeStake, "failed to update stake for resource node")
	ErrUpdateMetaNodeStake                = sdkerrors.Register(ModuleName, codeErrUpdateMetaNodeStake, "failed to update stake for meta node")
	ErrVoteMetaNode                       = sdkerrors.Register(ModuleName, codeErrVoteMetaNode, "failed to vote meta node")
	ErrResourceNodeRegDisabled            = sdkerrors.Register(ModuleName, codeErrResourceNodeRegDisabled, "resource node registration is disabled")
	ErrInvalidSuspensionStatForUnbondNode = sdkerrors.Register(ModuleName, codeErrInvalidSuspensionStatForUnbondNode, "cannot unbond a suspended node")
	ErrReporterAddress                    = sdkerrors.Register(ModuleName, codeErrReporterAddress, "invalid reporter address")
	ErrInvalidAmount                      = sdkerrors.Register(ModuleName, codeErrInvalidAmount, "invalid amount")
	ErrReporterAddressOrOwner             = sdkerrors.Register(ModuleName, codeErrReporterAddressOrOwner, "invalid reporter address or owner address")
	ErrReporterNotReachThreshold          = sdkerrors.Register(ModuleName, codeErrReporterNotReachThreshold, "reporter meta-nodes does not reach the threshold")
)
