package register

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"strconv"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for _, resourceNode := range data.ResourceNodes {
		ctx.Logger().Info("Init resource node: " + resourceNode.String())
		keeper.SetResourceNode(ctx, resourceNode)
		keeper.SetResourceNodeByPowerIndex(ctx, resourceNode)
	}

	for _, indexingNode := range data.IndexingNodes {
		ctx.Logger().Info("Init indexing node: " + indexingNode.String())
		keeper.SetIndexingNode(ctx, indexingNode)
		keeper.SetIndexingNodeByPowerIndex(ctx, indexingNode)
	}

	for _, resPow := range data.LastResourceNodePowers {
		ctx.Logger().Info("Init LastResourceNodePowers: address = " + resPow.Address.String() + ", power = " + strconv.FormatInt(resPow.Power, 10))
		keeper.SetLastResourceNodePower(ctx, resPow.Address, resPow.Power)
	}

	for _, idxPow := range data.LastIndexingNodePowers {
		ctx.Logger().Info("Init LastIndexingNodePowers: address = " + idxPow.Address.String() + ", power = " + strconv.FormatInt(idxPow.Power, 10))
		keeper.SetLastIndexingNodePower(ctx, idxPow.Address, idxPow.Power)
	}

	ctx.Logger().Info("Init LastResourceNodeTotalPower: " + data.LastResourceNodeTotalPower.String())
	keeper.SetLastResourceNodeTotalPower(ctx, data.LastResourceNodeTotalPower)
	ctx.Logger().Info("Init LastIndexingNodeTotalPower: " + data.LastIndexingNodeTotalPower.String())
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
	indexingNodes := keeper.GetAllIndexingNodes(ctx)

	return types.GenesisState{
		Params:                     params,
		LastResourceNodeTotalPower: lastResourceNodeTotalPower,
		LastResourceNodePowers:     lastResourceNodePowers,
		ResourceNodes:              resourceNodes,
		LastIndexingNodeTotalPower: lastIndexingNodeTotalPower,
		LastIndexingNodePowers:     lastIndexingNodePowers,
		IndexingNodes:              indexingNodes,
	}
}
