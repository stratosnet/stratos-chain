package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	k.SetParams(ctx, data.Params)
	k.SetTotalMinedTokens(ctx, data.GetTotalMinedToken())
	k.SetLastDistributedEpoch(ctx, data.LastDistributedEpoch)

	for _, immatureTotal := range data.ImmatureTotalInfo {
		walletAddr, err := sdk.AccAddressFromBech32(immatureTotal.WalletAddress)
		if err != nil {
			panic("invliad wallet address when init genesis of PoT module")
		}
		k.SetImmatureTotalReward(ctx, walletAddr, immatureTotal.Value)
	}

	for _, matureTotal := range data.MatureTotalInfo {
		walletAddr, err := sdk.AccAddressFromBech32(matureTotal.WalletAddress)
		if err != nil {
			panic("invliad wallet address when init genesis of PoT module")
		}
		k.SetMatureTotalReward(ctx, walletAddr, matureTotal.Value)
	}

	for _, individual := range data.IndividualRewardInfo {
		walletAddr, err := sdk.AccAddressFromBech32(individual.WalletAddress)
		if err != nil {
			panic("invliad wallet address when init genesis of PoT module")
		}
		k.SetIndividualReward(ctx, walletAddr, data.LastDistributedEpoch.Add(sdkmath.NewInt(data.Params.MatureEpoch)), individual)
	}

	k.SetMaturedEpoch(ctx, data.MaturedEpoch)
	// ensure total supply of bank module is LT InitialTotalSupply
	totalSupply := k.GetSupply(ctx)
	if k.GetParams(ctx).InitialTotalSupply.IsLT(totalSupply) {
		errMsg := fmt.Sprintf("current total supply[%v] is greater than total supply limit[%v]",
			totalSupply.String(), k.GetParams(ctx).InitialTotalSupply.String())
		panic(errMsg)
	}

	for _, rewardTotal := range data.RewardTotalInfo {
		epoch := rewardTotal.Epoch
		reward := rewardTotal.TotalReward
		k.SetTotalReward(ctx, epoch, reward)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func (k Keeper) ExportGenesis(ctx sdk.Context) (data *types.GenesisState) {
	params := k.GetParams(ctx)
	totalMinedToken := k.GetTotalMinedTokens(ctx)
	lastDistributedEpoch := k.GetLastDistributedEpoch(ctx)

	var individualRewardInfo []types.Reward
	var immatureTotalInfo []types.ImmatureTotal
	k.IteratorImmatureTotal(ctx, func(walletAddress sdk.AccAddress, reward sdk.Coins) (stop bool) {
		if !reward.Empty() && !reward.IsZero() {
			immatureTotal := types.NewImmatureTotal(walletAddress, reward)
			immatureTotalInfo = append(immatureTotalInfo, immatureTotal)

			miningReward := sdk.NewCoins(sdk.NewCoin(k.RewardDenom(ctx), reward.AmountOf(k.RewardDenom(ctx))))
			trafficReward := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), reward.AmountOf(k.BondDenom(ctx))))
			individualReward := types.NewReward(walletAddress, miningReward, trafficReward)
			individualRewardInfo = append(individualRewardInfo, individualReward)

		}
		return false
	})

	var matureTotalInfo []types.MatureTotal
	k.IteratorMatureTotal(ctx, func(walletAddress sdk.AccAddress, reward sdk.Coins) (stop bool) {
		if !reward.Empty() && !reward.IsZero() {
			matureTotal := types.NewMatureTotal(walletAddress, reward)
			matureTotalInfo = append(matureTotalInfo, matureTotal)
		}
		return false
	})

	maturedEpoch := k.GetMaturedEpoch(ctx)

	var rewardTotalInfo []types.RewardTotal
	k.IteratorTotalReward(ctx, func(epoch sdkmath.Int, totalReward types.TotalReward) (stop bool) {
		if epoch.GT(sdkmath.ZeroInt()) {
			info := types.NewRewardTotal(epoch, totalReward)
			rewardTotalInfo = append(rewardTotalInfo, info)
		}
		return false
	})

	return types.NewGenesisState(
		params,
		totalMinedToken,
		lastDistributedEpoch,
		immatureTotalInfo,
		matureTotalInfo,
		individualRewardInfo,
		maturedEpoch,
		rewardTotalInfo)
}
