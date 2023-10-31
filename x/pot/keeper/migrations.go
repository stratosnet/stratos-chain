package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v011 "github.com/stratosnet/stratos-chain/x/pot/legacy/v011"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace types.ParamsSubspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, legacySubspace types.ParamsSubspace) Migrator {
	return Migrator{keeper: keeper, legacySubspace: legacySubspace}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v011.MigrateStore(ctx, m.keeper.storeKey, m.legacySubspace, m.keeper.cdc)
}
