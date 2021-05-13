package register

import (
	"github.com/stratosnet/stratos-chain/x/register/keeper"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

const (
	DefaultParamSpace = types.DefaultParamSpace
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
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
)

type (
	Keeper = keeper.Keeper
)
