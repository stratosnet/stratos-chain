package register

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)

	initialStakeTotal := sdk.ZeroInt()

	for _, resourceNode := range data.ResourceNodes {
		initialStakeTotal = initialStakeTotal.Add(resourceNode.GetTokens())
		keeper.SetResourceNode(ctx, resourceNode)
	}

	for _, indexingNode := range data.IndexingNodes {
		initialStakeTotal = initialStakeTotal.Add(indexingNode.GetTokens())
		keeper.SetIndexingNode(ctx, indexingNode)
	}

	for _, resStake := range data.LastResourceNodeStakes {
		keeper.SetLastResourceNodeStake(ctx, resStake.Address, resStake.Stake)
	}

	for _, idxStake := range data.LastIndexingNodeStakes {
		keeper.SetLastIndexingNodeStake(ctx, idxStake.Address, idxStake.Stake)
	}

	keeper.SetLastResourceNodeTotalStake(ctx, data.LastResourceNodeTotalStake)
	keeper.SetLastIndexingNodeTotalStake(ctx, data.LastIndexingNodeTotalStake)
	keeper.SetInitialGenesisStakeTotal(ctx, initialStakeTotal)
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data types.GenesisState) {
	params := keeper.GetParams(ctx)

	lastResourceNodeTotalStake := keeper.GetLastResourceNodeTotalStake(ctx)
	lastIndexingNodeTotalStake := keeper.GetLastIndexingNodeTotalStake(ctx)

	var lastResourceNodeStakes []types.LastResourceNodeStake
	keeper.IterateLastResourceNodeStakes(ctx, func(addr sdk.AccAddress, stake sdk.Int) (stop bool) {
		lastResourceNodeStakes = append(lastResourceNodeStakes, types.LastResourceNodeStake{Address: addr, Stake: stake})
		return false
	})

	var lastIndexingNodeStakes []types.LastIndexingNodeStake
	keeper.IterateLastIndexingNodeStakes(ctx, func(addr sdk.AccAddress, stake sdk.Int) (stop bool) {
		lastIndexingNodeStakes = append(lastIndexingNodeStakes, types.LastIndexingNodeStake{Address: addr, Stake: stake})
		return false
	})

	resourceNodes := keeper.GetAllResourceNodes(ctx)
	indexingNodes := keeper.GetAllIndexingNodes(ctx)

	return types.GenesisState{
		Params:                     params,
		LastResourceNodeTotalStake: lastResourceNodeTotalStake,
		LastResourceNodeStakes:     lastResourceNodeStakes,
		ResourceNodes:              resourceNodes,
		LastIndexingNodeTotalStake: lastIndexingNodeTotalStake,
		LastIndexingNodeStakes:     lastIndexingNodeStakes,
		IndexingNodes:              indexingNodes,
	}
}
