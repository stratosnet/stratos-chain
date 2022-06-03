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
	resNodeBondedToken := sdk.ZeroInt()
	resNodeNotBondedToken := sdk.ZeroInt()
	for _, resourceNode := range data.GetResourceNodes() {
		if resourceNode.GetStatus() == stakingtypes.Bonded {
			initialStakeTotal = initialStakeTotal.Add(resourceNode.Tokens)
			resNodeBondedToken = resNodeBondedToken.Add(resourceNode.Tokens)
		} else if resourceNode.GetStatus() == stakingtypes.Unbonded {
			resNodeNotBondedToken = resNodeNotBondedToken.Add(resourceNode.Tokens)
		}
		keeper.SetResourceNode(ctx, resourceNode)
	}
	err := keeper.MintResourceNodeBondedTokenPool(ctx, sdk.NewCoin(keeper.BondDenom(ctx), resNodeBondedToken))
	if err != nil {
		panic(err)
	}
	err = keeper.MintResourceNodeNotBondedTokenPool(ctx, sdk.NewCoin(keeper.BondDenom(ctx), resNodeNotBondedToken))
	if err != nil {
		panic(err)
	}
	idxNodeBondedToken := sdk.ZeroInt()
	idxNodeNotBondedToken := sdk.ZeroInt()
	for _, metaNode := range data.GetMetaNodes() {
		if metaNode.GetStatus() == stakingtypes.Bonded {
			initialStakeTotal = initialStakeTotal.Add(metaNode.Tokens)
			idxNodeBondedToken = idxNodeBondedToken.Add(metaNode.Tokens)
		} else if metaNode.GetStatus() == stakingtypes.Unbonded {
			idxNodeNotBondedToken = idxNodeNotBondedToken.Add(metaNode.Tokens)
		}
		keeper.SetMetaNode(ctx, metaNode)
	}

	totalUnissuedPrepay := data.TotalUnissuedPrepay
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
	totalUnissuedPrepay := keeper.GetTotalUnissuedPrepay(ctx).Amount
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
		Params:              &params,
		ResourceNodes:       resourceNodes,
		MetaNodes:           metaNodes,
		InitialUozPrice:     initialUOzonePrice,
		TotalUnissuedPrepay: totalUnissuedPrepay,
		Slashing:            slashingInfo,
	}
}
