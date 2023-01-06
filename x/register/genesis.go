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

	freshStart := keeper.GetResourceNodeNotBondedToken(ctx).IsZero() &&
		keeper.GetResourceNodeBondedToken(ctx).IsZero() &&
		keeper.GetMetaNodeNotBondedToken(ctx).IsZero() &&
		keeper.GetMetaNodeBondedToken(ctx).IsZero()

	initialStakeTotal := sdk.ZeroInt()
	lenOfGenesisBondedResourceNode := int64(0)

	for _, resourceNode := range data.GetResourceNodes() {
		ownerAddr, err := sdk.AccAddressFromBech32(resourceNode.OwnerAddress)
		if err != nil {
			panic(err)
		}
		switch resourceNode.GetStatus() {
		case stakingtypes.Bonded:
			lenOfGenesisBondedResourceNode++
			if !resourceNode.Suspend {
				initialStakeTotal = initialStakeTotal.Add(resourceNode.Tokens)
			}
			if freshStart {
				err = keeper.SendCoinsFromAccountToResNodeBondedPool(ctx, ownerAddr, sdk.NewCoin(keeper.BondDenom(ctx), resourceNode.Tokens))
				if err != nil {
					panic(err)
				}
			}
		case stakingtypes.Unbonded:
			if freshStart {
				err = keeper.SendCoinsFromAccountToResNodeNotBondedPool(ctx, ownerAddr, sdk.NewCoin(keeper.BondDenom(ctx), resourceNode.Tokens))
				if err != nil {
					panic(err)
				}
			}
		default:
			panic(types.ErrInvalidNodeStat)
		}
		keeper.SetResourceNode(ctx, resourceNode)
	}
	// set initial genesis number of resource nodes
	keeper.SetBondedResourceNodeCnt(ctx, sdk.NewInt(lenOfGenesisBondedResourceNode))

	lenOfGenesisBondedMetaNode := int64(0)
	for _, metaNode := range data.GetMetaNodes() {
		ownerAddr, err := sdk.AccAddressFromBech32(metaNode.OwnerAddress)
		if err != nil {
			panic(err)
		}
		switch metaNode.GetStatus() {
		case stakingtypes.Bonded:
			lenOfGenesisBondedMetaNode++
			initialStakeTotal = initialStakeTotal.Add(metaNode.Tokens)
			if freshStart {
				err = keeper.SendCoinsFromAccountToMetaNodeBondedPool(ctx, ownerAddr, sdk.NewCoin(keeper.BondDenom(ctx), metaNode.Tokens))
				if err != nil {
					panic(err)
				}
			}
		case stakingtypes.Unbonded:
			if freshStart {
				err = keeper.SendCoinsFromAccountToMetaNodeNotBondedPool(ctx, ownerAddr, sdk.NewCoin(keeper.BondDenom(ctx), metaNode.Tokens))
				if err != nil {
					panic(err)
				}
			}
		default:
			panic(types.ErrInvalidNodeStat)
		}
		keeper.SetMetaNode(ctx, metaNode)
	}
	// set initial genesis number of meta nodes
	keeper.SetBondedMetaNodeCnt(ctx, sdk.NewInt(lenOfGenesisBondedMetaNode))

	//totalUnissuedPrepay := keeper.GetTotalUnissuedPrepay(ctx).Amount
	//effectiveNOzonePrice := sdk.ZeroDec()
	//effectiveNOzonePrice = effectiveNOzonePrice.Add(data.NozPrice)
	//keeper.SetNOzonePrice(ctx, effectiveNOzonePrice)
	keeper.SetInitialGenesisStakeTotal(ctx, initialStakeTotal)
	keeper.SetEffectiveGenesisStakeTotal(ctx, initialStakeTotal)
	stakeNozRate := sdk.ZeroDec()
	stakeNozRate = stakeNozRate.Add(data.StakeNozRate)
	keeper.SetStakeNozRate(ctx, stakeNozRate)

	// calc total noz supply with EffectiveGenesisStakeTotal and stakeNozRate
	totalNozSupply := initialStakeTotal.ToDec().Quo(stakeNozRate).TruncateInt()

	initOzoneLimit := sdk.ZeroInt()
	if data.RemainingNozLimit.Equal(sdk.ZeroInt()) {
		// fresh start
		initOzoneLimit = initOzoneLimit.Add(totalNozSupply)
	} else {
		// not fresh start
		initOzoneLimit = initOzoneLimit.Add(data.RemainingNozLimit)
	}
	keeper.SetRemainingOzoneLimit(ctx, initOzoneLimit)
	//initOzoneLimit := initialStakeTotal.Add(totalUnissuedPrepay).ToDec().Quo(effectiveNOzonePrice).TruncateInt()
	//keeper.SetRemainingOzoneLimit(ctx, initOzoneLimit)

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
	//effectiveNOzonePrice := keeper.CurrNozPrice(ctx)
	remainingNozLimit := keeper.GetRemainingOzoneLimit(ctx)
	stakeNozRate := keeper.GetStakeNozRate(ctx)

	var slashingInfo []*types.Slashing
	keeper.IteratorSlashingInfo(ctx, func(walletAddress sdk.AccAddress, val sdk.Int) (stop bool) {
		if val.GT(sdk.ZeroInt()) {
			slashing := types.NewSlashing(walletAddress, val)
			slashingInfo = append(slashingInfo, slashing)
		}
		return false
	})

	return &types.GenesisState{
		Params:            &params,
		ResourceNodes:     resourceNodes,
		MetaNodes:         metaNodes,
		RemainingNozLimit: remainingNozLimit,
		Slashing:          slashingInfo,
		StakeNozRate:      stakeNozRate,
	}
}
