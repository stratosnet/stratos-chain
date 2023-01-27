package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v09register "github.com/stratosnet/stratos-chain/x/register/legacy/v09"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper     Keeper
	aminoCodec *codec.LegacyAmino
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, aminoCodec *codec.LegacyAmino) Migrator {
	return Migrator{keeper: keeper, aminoCodec: aminoCodec}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v09register.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc, m.aminoCodec)
}
