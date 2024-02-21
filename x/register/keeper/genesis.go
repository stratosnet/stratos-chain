package keeper

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/stratosnet/stratos-chain/x/register/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	k.SetParams(ctx, data.Params)

	freshStart := k.GetResourceNodeNotBondedToken(ctx).IsZero() &&
		k.GetResourceNodeBondedToken(ctx).IsZero() &&
		k.GetMetaNodeNotBondedToken(ctx).IsZero() &&
		k.GetMetaNodeBondedToken(ctx).IsZero()

	initialDepositTotal := sdkmath.ZeroInt()
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
				initialDepositTotal = initialDepositTotal.Add(resourceNode.EffectiveTokens)
			}
			if freshStart {
				amount := sdk.NewCoin(k.BondDenom(ctx), resourceNode.Tokens)
				err = k.GetBankKeeper().SendCoinsFromAccountToModule(ctx, ownerAddr, types.ResourceNodeBondedPool, sdk.NewCoins(amount))
				if err != nil {
					panic(err)
				}
			}
		case stakingtypes.Unbonded:
			if freshStart {
				amount := sdk.NewCoin(k.BondDenom(ctx), resourceNode.Tokens)
				err = k.GetBankKeeper().SendCoinsFromAccountToModule(ctx, ownerAddr, types.ResourceNodeNotBondedPool, sdk.NewCoins(amount))
				if err != nil {
					panic(err)
				}
			}
		default:
			panic(types.ErrInvalidNodeStat)
		}

		if len(strings.TrimSpace(resourceNode.BeneficiaryAddress)) == 0 {
			resourceNode.BeneficiaryAddress = resourceNode.OwnerAddress
		}
		k.SetResourceNode(ctx, resourceNode)
	}
	// set initial genesis number of resource nodes
	k.SetBondedResourceNodeCnt(ctx, sdkmath.NewInt(lenOfGenesisBondedResourceNode))

	lenOfGenesisBondedMetaNode := int64(0)
	for _, metaNode := range data.GetMetaNodes() {
		ownerAddr, err := sdk.AccAddressFromBech32(metaNode.OwnerAddress)
		if err != nil {
			panic(err)
		}
		switch metaNode.GetStatus() {
		case stakingtypes.Bonded:
			lenOfGenesisBondedMetaNode++
			if !metaNode.Suspend {
				initialDepositTotal = initialDepositTotal.Add(metaNode.Tokens)
			}
			if freshStart {
				amount := sdk.NewCoin(k.BondDenom(ctx), metaNode.Tokens)
				err = k.GetBankKeeper().SendCoinsFromAccountToModule(ctx, ownerAddr, types.MetaNodeBondedPool, sdk.NewCoins(amount))
				if err != nil {
					panic(err)
				}
			}
		case stakingtypes.Unbonded:
			if freshStart {
				amount := sdk.NewCoin(k.BondDenom(ctx), metaNode.Tokens)
				err = k.GetBankKeeper().SendCoinsFromAccountToModule(ctx, ownerAddr, types.MetaNodeNotBondedPool, sdk.NewCoins(amount))
				if err != nil {
					panic(err)
				}
			}
		default:
			panic(types.ErrInvalidNodeStat)
		}

		if len(strings.TrimSpace(metaNode.BeneficiaryAddress)) == 0 {
			metaNode.BeneficiaryAddress = metaNode.OwnerAddress
		}
		k.SetMetaNode(ctx, metaNode)
	}
	// set initial genesis number of meta nodes
	k.SetBondedMetaNodeCnt(ctx, sdkmath.NewInt(lenOfGenesisBondedMetaNode))

	totalUnissuedPrepay := k.GetTotalUnissuedPrepay(ctx).Amount
	k.SetInitialGenesisDepositTotal(ctx, initialDepositTotal)
	k.SetEffectiveTotalDeposit(ctx, initialDepositTotal)
	depositNozRate := sdkmath.LegacyZeroDec()
	depositNozRate = depositNozRate.Add(data.DepositNozRate)
	k.SetDepositNozRate(ctx, depositNozRate)

	// calc total noz supply with EffectiveGenesisDepositTotal and depositNozRate
	totalNozSupply := initialDepositTotal.ToLegacyDec().Quo(depositNozRate).TruncateInt()
	initOzoneLimit := sdkmath.ZeroInt()
	if freshStart && totalUnissuedPrepay.Equal(sdkmath.ZeroInt()) {
		// fresh start
		initOzoneLimit = initOzoneLimit.Add(totalNozSupply)
	} else {
		// not fresh start
		initOzoneLimit = initOzoneLimit.Add(data.RemainingNozLimit)
	}
	k.SetRemainingOzoneLimit(ctx, initOzoneLimit)

	for _, slashing := range data.Slashing {
		walletAddress, err := sdk.AccAddressFromBech32(slashing.GetWalletAddress())
		if err != nil {
			panic(err)
		}

		k.SetSlashing(ctx, walletAddress, sdkmath.NewInt(slashing.Value))
	}

	k.ReloadMetaNodeBitMapIdxCache(ctx)

	return
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func (k Keeper) ExportGenesis(ctx sdk.Context) (data *types.GenesisState) {
	params := k.GetParams(ctx)

	resourceNodes := k.GetAllResourceNodes(ctx)
	metaNodes := k.GetAllMetaNodes(ctx)
	remainingNozLimit := k.GetRemainingOzoneLimit(ctx)
	depositNozRate := k.GetDepositNozRate(ctx)

	var slashingInfo []types.Slashing
	k.IteratorSlashingInfo(ctx, func(walletAddress sdk.AccAddress, val sdkmath.Int) (stop bool) {
		if val.GT(sdkmath.ZeroInt()) {
			slashing := types.NewSlashing(walletAddress, val)
			slashingInfo = append(slashingInfo, slashing)
		}
		return false
	})

	return types.NewGenesisState(params, resourceNodes, metaNodes, remainingNozLimit, slashingInfo, depositNozRate)
}
