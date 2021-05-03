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

	for _, resPow := range data.LastResourceNodePowers {
		keeper.SetLastResourceNodePower(ctx, resPow.Address, resPow.Power)
	}
	keeper.SetLastResourceNodeTotalPower(ctx, data.LastResourceNodeTotalPower)
	for _, idxPow := range data.LastIndexingNodePowers {
		keeper.SetLastIndexingNodePower(ctx, idxPow.Address, idxPow.Power)
	}
	keeper.SetLastIndexingNodeTotalPower(ctx, data.LastIndexingNodeTotalPower)
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data types.GenesisState) {
	params := keeper.GetParams(ctx)

	lastResourceNodeTotalPower := keeper.GetLastResourceNodeTotalPower(ctx)
	lastIndexingNodeTotalPower := keeper.GetLastIndexingNodeTotalPower(ctx)

	var lastResourceNodePowers []types.LastResourceNodePower
	keeper.IterateLastResourceNodePowers(ctx, func(addr sdk.AccAddress, power int64) (stop bool) {
		lastResourceNodePowers = append(lastResourceNodePowers, types.LastResourceNodePower{Address: addr, Power: power})
		return false
	})

	var lastIndexingNodePowers []types.LastIndexingNodePower
	keeper.IterateLastIndexingNodePowers(ctx, func(addr sdk.AccAddress, power int64) (stop bool) {
		lastIndexingNodePowers = append(lastIndexingNodePowers, types.LastIndexingNodePower{Address: addr, Power: power})
		return false
	})

	resourceNodes := keeper.GetAllResourceNodes(ctx)
	indexingNodex := keeper.GetAllIndexingNodes(ctx)

	return types.GenesisState{
		Params:                     params,
		LastResourceNodeTotalPower: lastResourceNodeTotalPower,
		LastResourceNodePowers:     lastResourceNodePowers,
		ResourceNodes:              resourceNodes,
		LastIndexingNodeTotalPower: lastIndexingNodeTotalPower,
		LastIndexingNodePowers:     lastIndexingNodePowers,
		IndexingNodes:              indexingNodex,
	}
}
