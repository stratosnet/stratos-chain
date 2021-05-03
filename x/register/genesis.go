package register

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for _, resourceNode := range data.ResourceNodes {
		keeper.SetResourceNode(ctx, resourceNode)
		keeper.SetResourceNodeByPowerIndex(ctx, resourceNode)
	}

	for _, indexingNode := range data.IndexingNodes {
		keeper.SetIndexingNode(ctx, indexingNode)
		keeper.SetIndexingNodeByPowerIndex(ctx, indexingNode)
	}

	if data.Exported {
		for _, resPow := range data.LastResourceNodePowers {
			keeper.SetLastResourceNodePower(ctx, resPow.Address, resPow.Power)
		}
		keeper.SetLastResourceNodeTotalPower(ctx, data.LastResourceNodeTotalPower)
		for _, idxPow := range data.LastResourceNodePowers {
			keeper.SetLastResourceNodePower(ctx, idxPow.Address, idxPow.Power)
		}
		keeper.SetLastIndexingNodeTotalPower(ctx, data.LastIndexingNodeTotalPower)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data types.GenesisState) {
	// TODO: Define logic for exporting state
	params := keeper.GetParams(ctx)
	return types.GenesisState{
		Params: params,
	}
}
