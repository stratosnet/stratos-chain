package sds

import (
	"github.com/stratosnet/stratos-chain/x/sds/keeper"
	"github.com/stratosnet/stratos-chain/x/sds/types"
)

const (
	DefaultParamSpace = types.DefaultParamSpace
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
)

var (
	NewKeeper       = keeper.NewKeeper
	RegisterCodec   = types.RegisterCodec
	NewGenesisState = types.NewGenesisState
)

type (
	Keeper = keeper.Keeper
)
