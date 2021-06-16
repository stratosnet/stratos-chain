package sds

import (
	"github.com/stratosnet/stratos-chain/x/sds/keeper"
	"github.com/stratosnet/stratos-chain/x/sds/types"
)

const (
	DefaultParamSpace = types.DefaultParamSpace
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
)

var (
	NewKeeper = keeper.NewKeeper
)

type (
	Keeper = keeper.Keeper
)
