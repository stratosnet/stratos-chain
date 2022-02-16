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
	NewKeeper = keeper.NewKeeper
	TxDecoder = types.TxDecoder
)

type (
	Keeper        = keeper.Keeper
	GenesisState  = types.GenesisState
	CommitStateDB = *types.CommitStateDB
)
