package pot

import (
	"github.com/stratosnet/stratos-chain/x/pot/keeper"
	"github.com/stratosnet/stratos-chain/x/pot/types"
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
