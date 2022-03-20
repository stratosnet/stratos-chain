package register

import (
	"github.com/stratosnet/stratos-chain/x/register/keeper"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

const (
	DefaultParamSpace   = types.DefaultParamSpace
	ModuleName          = types.ModuleName
	StoreKey            = types.StoreKey
	RouterKey           = types.RouterKey
	NodeTypeComputation = types.COMPUTATION
	NodeTypeDataBase    = types.DATABASE
	NodeTypeStorage     = types.STORAGE
)

var (
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec

	ErrInvalid                  = types.ErrInvalid
	ErrInvalidNetworkAddr       = types.ErrInvalidNetworkAddr
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
	ErrInvalidApproverAddr      = types.ErrInvalidVoterAddr
	ErrInvalidApproverStatus    = types.ErrInvalidVoterStatus

	DefaultParams            = types.DefaultParams
	DefaultGenesisState      = types.DefaultGenesisState
	NewGenesisState          = types.NewGenesisState
	NewResourceNode          = types.NewResourceNode
	NewIndexingNode          = types.NewIndexingNode
	NewDescription           = types.NewDescription
	NewMsgCreateResourceNode = types.NewMsgCreateResourceNode
	NewMsgCreateIndexingNode = types.NewMsgCreateIndexingNode

	GetGenesisStateFromAppState = types.GetGenesisStateFromAppState

	NewMultiRegisterHooks = types.NewMultiRegisterHooks
)

type (
	Keeper                = keeper.Keeper
	ResourceNode          = types.ResourceNode
	IndexingNode          = types.IndexingNode
	Description           = types.Description
	GenesisIndexingNode   = types.GenesisIndexingNode
	Slashing              = types.Slashing
	MsgCreateResourceNode = types.MsgCreateResourceNode
	MsgCreateIndexingNode = types.MsgCreateIndexingNode
	VoteOpinion           = types.VoteOpinion
)
