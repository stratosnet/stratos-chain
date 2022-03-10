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
	resNodeBondedToken := sdk.ZeroInt()
	resNodeNotBondedToken := sdk.ZeroInt()
	for _, resourceNode := range data.ResourceNodes {
		if resourceNode.GetStatus() == sdk.Bonded {
			initialStakeTotal = initialStakeTotal.Add(resourceNode.GetTokens())
			resNodeBondedToken = resNodeBondedToken.Add(resourceNode.GetTokens())
		} else if resourceNode.GetStatus() == sdk.Unbonded {
			resNodeNotBondedToken = resNodeNotBondedToken.Add(resourceNode.GetTokens())
		}
		keeper.SetResourceNode(ctx, resourceNode)
	}
	keeper.SetResourceNodeBondedToken(ctx, sdk.NewCoin(keeper.BondDenom(ctx), resNodeBondedToken))
	keeper.SetResourceNodeNotBondedToken(ctx, sdk.NewCoin(keeper.BondDenom(ctx), resNodeNotBondedToken))

	idxNodeBondedToken := sdk.ZeroInt()
	idxNodeNotBondedToken := sdk.ZeroInt()
	for _, indexingNode := range data.IndexingNodes {
		if indexingNode.GetStatus() == sdk.Bonded {
			initialStakeTotal = initialStakeTotal.Add(indexingNode.GetTokens())
			idxNodeBondedToken = idxNodeBondedToken.Add(indexingNode.GetTokens())
		} else if indexingNode.GetStatus() == sdk.Unbonded {
			idxNodeNotBondedToken = idxNodeNotBondedToken.Add(indexingNode.GetTokens())
		}
		keeper.SetIndexingNode(ctx, indexingNode)
	}
	keeper.SetIndexingNodeBondedToken(ctx, sdk.NewCoin(keeper.BondDenom(ctx), idxNodeBondedToken))
	keeper.SetIndexingNodeNotBondedToken(ctx, sdk.NewCoin(keeper.BondDenom(ctx), idxNodeNotBondedToken))

	for _, resStake := range data.LastResourceNodeStakes {
		keeper.SetLastResourceNodeStake(ctx, resStake.Address, resStake.Stake)
	}

	for _, idxStake := range data.LastIndexingNodeStakes {
		keeper.SetLastIndexingNodeStake(ctx, idxStake.Address, idxStake.Stake)
	}

	totalUnissuedPrepay := data.TotalUnissuedPrepay
	initialUOzonePrice := sdk.ZeroDec()
	initialUOzonePrice = initialUOzonePrice.Add(data.InitialUozPrice)
	keeper.SetInitialGenesisStakeTotal(ctx, initialStakeTotal)
	keeper.SetInitialUOzonePrice(ctx, initialUOzonePrice)
	initOzoneLimit := initialStakeTotal.Add(totalUnissuedPrepay).ToDec().Quo(initialUOzonePrice).TruncateInt()
	keeper.SetRemainingOzoneLimit(ctx, initOzoneLimit)
	keeper.SetTotalUnissuedPrepay(ctx, sdk.Coin{
		Denom:  data.Params.BondDenom,
		Amount: totalUnissuedPrepay,
	})
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data types.GenesisState) {
	params := keeper.GetParams(ctx)

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
	totalUnissuedPrepay := keeper.GetTotalUnissuedPrepay(ctx).Amount
	initialUOzonePrice := keeper.CurrUozPrice(ctx)

	return types.GenesisState{
		Params:                 params,
		LastResourceNodeStakes: lastResourceNodeStakes,
		ResourceNodes:          resourceNodes,
		LastIndexingNodeStakes: lastIndexingNodeStakes,
		IndexingNodes:          indexingNodes,
		InitialUozPrice:        initialUOzonePrice,
		TotalUnissuedPrepay:    totalUnissuedPrepay,
	}
}
