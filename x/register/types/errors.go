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
	codeErrInvalidVoterNetworkAddr
	codeErrInvalidVoterOwnerAddr
	codeErrSameAddr
	codeErrInvalidOwnerAddr
	codeErrInvalidVoterAddr
	codeErrInvalidVoterStatus
	codeEcoderrNoRegistrationVotePoolFound
	codeErrDuplicateVoting
	codeErrVoteExpired
	codeErrInsufficientBalanceOfBondedPool
	codeErrInsufficientBalanceOfNotBondedPool
	codeErrSubAllTokens
	codeErrED25519InvalidPubKey
	codeErrEmptyNodeNetworkAddress
	codeErrEmptyPubKey
	codeErrInvalidGenesisToken
	codeErrNoUnbondingNode
	codeErrMaxUnbondingNodeEntries
	codeErrNoNodeForAddress
	codeErrUnbondingNode
	codeErrInvalidNodeStatBonded
	codeErrInitialUOzonePrice
	codeErrInvalidStakeChange
	codeErrTotalUnissuedPrepay
	codeErrInvalidNodeType
	codeErrUnknownAccountAddress
	codeErrUnknownPubKey
	codeErrNoNodeFound
	codeErrInitialBalanceNotZero
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
	ErrInvalidVoterNetworkAddr            = sdkerrors.Register(ModuleName, codeErrInvalidVoterNetworkAddr, "invalid voter network address")
	ErrInvalidVoterOwnerAddr              = sdkerrors.Register(ModuleName, codeErrInvalidVoterOwnerAddr, "invalid voter owner address")
	ErrSameAddr                           = sdkerrors.Register(ModuleName, codeErrSameAddr, "node address should not same as the voter address")
	ErrInvalidOwnerAddr                   = sdkerrors.Register(ModuleName, codeErrInvalidOwnerAddr, "invalid owner address")
	ErrInvalidVoterAddr                   = sdkerrors.Register(ModuleName, codeErrInvalidVoterAddr, "invalid voter address")
	ErrInvalidVoterStatus                 = sdkerrors.Register(ModuleName, codeErrInvalidVoterStatus, "invalid voter status")
	ErrNoRegistrationVotePoolFound        = sdkerrors.Register(ModuleName, codeEcoderrNoRegistrationVotePoolFound, "registration pool does not exist")
	ErrDuplicateVoting                    = sdkerrors.Register(ModuleName, codeErrDuplicateVoting, "duplicate voting")
	ErrVoteExpired                        = sdkerrors.Register(ModuleName, codeErrVoteExpired, "vote expired")
	ErrInsufficientBalanceOfBondedPool    = sdkerrors.Register(ModuleName, codeErrInsufficientBalanceOfBondedPool, "insufficient balance of bonded pool")
	ErrInsufficientBalanceOfNotBondedPool = sdkerrors.Register(ModuleName, codeErrInsufficientBalanceOfNotBondedPool, "insufficient balance of not bonded pool")
	ErrSubAllTokens                       = sdkerrors.Register(ModuleName, codeErrSubAllTokens, "can not sub all tokens since the node is still bonded")
	ErrED25519InvalidPubKey               = sdkerrors.Register(ModuleName, codeErrED25519InvalidPubKey, "ED25519 public keys are unsupported")
	ErrEmptyNodeNetworkAddress            = sdkerrors.Register(ModuleName, codeErrEmptyNodeNetworkAddress, "missing node network address")
	ErrEmptyPubKey                        = sdkerrors.Register(ModuleName, codeErrEmptyPubKey, "missing public key")
	ErrInvalidGenesisToken                = sdkerrors.Register(ModuleName, codeErrInvalidGenesisToken, "invalid genesis token")
	ErrNoUnbondingNode                    = sdkerrors.Register(ModuleName, codeErrNoUnbondingNode, "no unbonding node found")
	ErrMaxUnbondingNodeEntries            = sdkerrors.Register(ModuleName, codeErrMaxUnbondingNodeEntries, "too many unbonding node entries for networkAddr tuple")
	ErrNoNodeForAddress                   = sdkerrors.Register(ModuleName, codeErrNoNodeForAddress, "registered node does not contain address")
	ErrUnbondingNode                      = sdkerrors.Register(ModuleName, codeErrUnbondingNode, "changes cannot be made to an unbonding node")
	ErrInvalidNodeStatBonded              = sdkerrors.Register(ModuleName, codeErrInvalidNodeStatBonded, "invalid node status: bonded")
	ErrInitialUOzonePrice                 = sdkerrors.Register(ModuleName, codeErrInitialUOzonePrice, "initial uOzone price must be positive")
	ErrInvalidStakeChange                 = sdkerrors.Register(ModuleName, codeErrInvalidStakeChange, "invalid change for stake")
	ErrTotalUnissuedPrepay                = sdkerrors.Register(ModuleName, codeErrTotalUnissuedPrepay, "total unissued prepay must be non-negative")
	ErrInvalidNodeType                    = sdkerrors.Register(ModuleName, codeErrInvalidNodeType, "invalid node type")
	ErrUnknownAccountAddress              = sdkerrors.Register(ModuleName, codeErrUnknownAccountAddress, "account address does not exist")
	ErrUnknownPubKey                      = sdkerrors.Register(ModuleName, codeErrUnknownPubKey, "unknown pubKey ")
	ErrNoNodeFound                        = sdkerrors.Register(ModuleName, codeErrNoNodeFound, "node does not exist ")
	ErrInitialBalanceNotZero              = sdkerrors.Register(ModuleName, codeErrInitialBalanceNotZero, "initial balance isn't zero ")
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
)
