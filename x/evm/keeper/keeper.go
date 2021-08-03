package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/stratosnet/stratos-chain/x/evm/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the pot store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramSpace    params.Subspace
	AccountKeeper auth.AccountKeeper
}

// NewKeeper creates a pot keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, accountKeeper auth.AccountKeeper) Keeper {
	keeper := Keeper{
		cdc:           cdc,
		storeKey:      key,
		paramSpace:    paramSpace.WithKeyTable(types.ParamKeyTable()),
		AccountKeeper: accountKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
