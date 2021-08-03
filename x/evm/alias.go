package evm

import (
	"github.com/stratosnet/stratos-chain/x/evm/keeper"
	"github.com/stratosnet/stratos-chain/x/evm/types"
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
	ParamKeyTable   = types.ParamKeyTable
	NewGenesisState = types.NewGenesisState
)

type (
	Keeper = keeper.Keeper
)
