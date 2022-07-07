package register

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stratosnet/stratos-chain/x/register/keeper"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data *types.GenesisState) {
	keeper.SetParams(ctx, *data.Params)

	initialStakeTotal := sdk.ZeroInt()
	lenOfGenesisBondedResourceNode := int64(0)
	initResourceNodeNotBondedPoolBalances := sdk.NewCoins(sdk.Coin{
		Denom:  keeper.BondDenom(ctx),
		Amount: sdk.ZeroInt(),
	})
	initResourceNodeBondedPoolBalances := sdk.NewCoins(sdk.Coin{
		Denom:  keeper.BondDenom(ctx),
		Amount: sdk.ZeroInt(),
	})

	for _, resourceNode := range data.GetResourceNodes() {
		if resourceNode.GetStatus() == stakingtypes.Bonded {
			lenOfGenesisBondedResourceNode++
			initialStakeTotal = initialStakeTotal.Add(resourceNode.Tokens)
			initResourceNodeBondedPoolBalances = initResourceNodeBondedPoolBalances.Add(sdk.Coin{
				Denom:  keeper.BondDenom(ctx),
				Amount: resourceNode.Tokens,
			})
		} else if resourceNode.GetStatus() == stakingtypes.Unbonded {
			initResourceNodeNotBondedPoolBalances = initResourceNodeNotBondedPoolBalances.Add(sdk.Coin{
				Denom:  keeper.BondDenom(ctx),
				Amount: resourceNode.Tokens,
			})
		}
		keeper.SetResourceNode(ctx, resourceNode)
	}
	// set initial genesis number of resource nodes
	keeper.SetBondedResourceNodeCnt(ctx, sdk.NewInt(lenOfGenesisBondedResourceNode))
	err := keeper.MintResourceNodeNotBondedPoolWhenInitGenesis(ctx, initResourceNodeNotBondedPoolBalances)
	if err != nil {
		panic(err)
	}
	err = keeper.MintResourceNodeBondedPoolWhenInitGenesis(ctx, initResourceNodeBondedPoolBalances)
	if err != nil {
		panic(err)
	}

	lenOfGenesisBondedMetaNode := int64(0)
	initMetaNodeNotBondedPoolBalances := sdk.NewCoins(sdk.Coin{
		Denom:  keeper.BondDenom(ctx),
		Amount: sdk.ZeroInt(),
	})
	initMetaNodeBondedPoolBalances := sdk.NewCoins(sdk.Coin{
		Denom:  keeper.BondDenom(ctx),
		Amount: sdk.ZeroInt(),
	})
	for _, metaNode := range data.GetMetaNodes() {
		if metaNode.GetStatus() == stakingtypes.Bonded {
			lenOfGenesisBondedMetaNode++
			initialStakeTotal = initialStakeTotal.Add(metaNode.Tokens)
			initMetaNodeBondedPoolBalances = initMetaNodeBondedPoolBalances.Add(sdk.Coin{
				Denom:  keeper.BondDenom(ctx),
				Amount: metaNode.Tokens,
			})
		} else if metaNode.GetStatus() == stakingtypes.Unbonded {
			initMetaNodeNotBondedPoolBalances = initMetaNodeNotBondedPoolBalances.Add(sdk.Coin{
				Denom:  keeper.BondDenom(ctx),
				Amount: metaNode.Tokens,
			})
		}
		keeper.SetMetaNode(ctx, metaNode)
	}
	// set initial genesis number of meta nodes
	keeper.SetBondedMetaNodeCnt(ctx, sdk.NewInt(lenOfGenesisBondedMetaNode))
	err = keeper.MintMetaNodeNotBondedPoolWhenInitGenesis(ctx, initMetaNodeNotBondedPoolBalances)
	if err != nil {
		panic(err)
	}
	err = keeper.MintMetaNodeBondedPoolWhenInitGenesis(ctx, initMetaNodeBondedPoolBalances)
	if err != nil {
		panic(err)
	}

	totalUnissuedPrepay := keeper.GetTotalUnissuedPrepay(ctx).Amount
	initialUOzonePrice := sdk.ZeroDec()
	initialUOzonePrice = initialUOzonePrice.Add(data.InitialUozPrice)
	keeper.SetInitialGenesisStakeTotal(ctx, initialStakeTotal)
	keeper.SetInitialUOzonePrice(ctx, initialUOzonePrice)
	initOzoneLimit := initialStakeTotal.Add(totalUnissuedPrepay).ToDec().Quo(initialUOzonePrice).TruncateInt()
	keeper.SetRemainingOzoneLimit(ctx, initOzoneLimit)

	for _, slashing := range data.Slashing {
		walletAddress, err := sdk.AccAddressFromBech32(slashing.GetWalletAddress())
		if err != nil {
			panic(err)
		}

		keeper.SetSlashing(ctx, walletAddress, sdk.NewInt(slashing.Value))
	}
	return
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) (data *types.GenesisState) {
	params := keeper.GetParams(ctx)

	resourceNodes := keeper.GetAllResourceNodes(ctx)
	metaNodes := keeper.GetAllMetaNodes(ctx)
	initialUOzonePrice := keeper.CurrUozPrice(ctx)

	var slashingInfo []*types.Slashing
	keeper.IteratorSlashingInfo(ctx, func(walletAddress sdk.AccAddress, val sdk.Int) (stop bool) {
		if val.GT(sdk.ZeroInt()) {
			slashing := types.NewSlashing(walletAddress, val)
			slashingInfo = append(slashingInfo, slashing)
		}
		return false
	})

	return &types.GenesisState{
		Params:          &params,
		ResourceNodes:   resourceNodes,
		MetaNodes:       metaNodes,
		InitialUozPrice: initialUOzonePrice,
		Slashing:        slashingInfo,
	}
}
