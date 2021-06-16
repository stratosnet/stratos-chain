package pot

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)
	keeper.SetFoundationAccount(ctx, data.FoundationAccount)
	keeper.SetInitialUOzonePrice(ctx, data.InitialUozPrice)
	keeper.SetMatureEpoch(ctx, data.MatureEpoch)
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data types.GenesisState) {
	params := keeper.GetParams(ctx)
	foundationAccount := keeper.GetFoundationAccount(ctx)
	initialUOzonePrice := keeper.GetInitialUOzonePrice(ctx)
	matureEpoch := keeper.GetMatureEpoch(ctx)

	return types.NewGenesisState(params, foundationAccount, initialUOzonePrice, matureEpoch)
}
