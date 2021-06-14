package register

import (
	"github.com/stratosnet/stratos-chain/x/register/keeper"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

const (
	DefaultParamSpace   = types.DefaultParamSpace
	ModuleName          = types.ModuleName
	StoreKey            = types.StoreKey
	NodeTypeComputation = types.COMPUTATION
	NodeTypeDataBase    = types.DATABASE
	NodeTypeStorage     = types.STORAGE
)

var (
	NewKeeper = keeper.NewKeeper

	ErrInvalid                  = types.ErrInvalid
	ErrEmptyNetworkAddr         = types.ErrEmptyNetworkAddr
	ErrEmptyOwnerAddr           = types.ErrEmptyOwnerAddr
	ErrValueNegative            = types.ErrValueNegative
	ErrEmptyDescription         = types.ErrEmptyDescription
	ErrEmptyResourceNodeAddr    = types.ErrEmptyResourceNodeAddr
	ErrEmptyIndexingNodeAddr    = types.ErrEmptyIndexingNodeAddr
	ErrBadDenom                 = types.ErrBadDenom
	ErrResourceNodePubKeyExists = types.ErrResourceNodePubKeyExists
	ErrIndexingNodePubKeyExists = types.ErrIndexingNodePubKeyExists
	ErrNoResourceNodeFound      = types.ErrNoResourceNodeFound
	ErrNoIndexingNodeFound      = types.ErrNoIndexingNodeFound
	ErrInvalidOwnerAddr         = types.ErrInvalidOwnerAddr
	ErrInvalidApproverAddr      = types.ErrInvalidApproverAddr
	ErrInvalidApproverStatus    = types.ErrInvalidApproverStatus

	DefaultParams            = types.DefaultParams
	NewDescription           = types.NewDescription
	NewMsgCreateResourceNode = types.NewMsgCreateResourceNode
	NewMsgCreateIndexingNode = types.NewMsgCreateIndexingNode
)

type (
	Keeper                = keeper.Keeper
	Description           = types.Description
	MsgCreateResourceNode = types.MsgCreateResourceNode
	MsgCreateIndexingNode = types.MsgCreateIndexingNode

	VoteOpinion = types.VoteOpinion
)
