package pot

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/pot/keeper"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data *types.GenesisState) {
	keeper.SetParams(ctx, data.Params)
	keeper.SetTotalMinedTokens(ctx, sdk.NewCoin(keeper.RewardDenom(ctx), sdk.NewInt(0)))
	keeper.SetLastDistributedEpoch(ctx, data.LastDistributedEpoch)

	for _, immatureTotal := range data.ImmatureTotalInfo {
		walletAddr, err := sdk.AccAddressFromBech32(immatureTotal.WalletAddress)
		if err != nil {
			panic("invliad wallet address when init genesis of PoT module")
		}
		keeper.SetImmatureTotalReward(ctx, walletAddr, immatureTotal.Value)
	}

	for _, matureTotal := range data.MatureTotalInfo {
		walletAddr, err := sdk.AccAddressFromBech32(matureTotal.WalletAddress)
		if err != nil {
			panic("invliad wallet address when init genesis of PoT module")
		}
		keeper.SetMatureTotalReward(ctx, walletAddr, matureTotal.Value)
	}

	for _, individual := range data.IndividualRewardInfo {
		walletAddr, err := sdk.AccAddressFromBech32(individual.WalletAddress)
		if err != nil {
			panic("invliad wallet address when init genesis of PoT module")
		}
		keeper.SetIndividualReward(ctx, walletAddr, data.LastDistributedEpoch.Add(sdk.NewInt(data.Params.MatureEpoch)), individual)
	}

	keeper.SetMaturedEpoch(ctx, data.MaturedEpoch)
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) (data *types.GenesisState) {
	params := keeper.GetParams(ctx)
	totalMinedToken := keeper.GetTotalMinedTokens(ctx)
	lastDistributedEpoch := keeper.GetLastDistributedEpoch(ctx)

	var individualRewardInfo []types.Reward
	var immatureTotalInfo []types.ImmatureTotal
	keeper.IteratorImmatureTotal(ctx, func(walletAddress sdk.AccAddress, reward sdk.Coins) (stop bool) {
		if !reward.Empty() && !reward.IsZero() {
			immatureTotal := types.NewImmatureTotal(walletAddress, reward)
			immatureTotalInfo = append(immatureTotalInfo, immatureTotal)

			miningReward := sdk.NewCoins(sdk.NewCoin(keeper.RewardDenom(ctx), reward.AmountOf(keeper.RewardDenom(ctx))))
			trafficReward := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), reward.AmountOf(keeper.BondDenom(ctx))))
			individualReward := types.NewReward(walletAddress, miningReward, trafficReward)
			individualRewardInfo = append(individualRewardInfo, individualReward)

		}
		return false
	})

	var matureTotalInfo []types.MatureTotal
	keeper.IteratorMatureTotal(ctx, func(walletAddress sdk.AccAddress, reward sdk.Coins) (stop bool) {
		if !reward.Empty() && !reward.IsZero() {
			matureTotal := types.NewMatureTotal(walletAddress, reward)
			matureTotalInfo = append(matureTotalInfo, matureTotal)
		}
		return false
	})

	maturedEpoch := keeper.GetMaturedEpoch(ctx)

	return types.NewGenesisState(
		params,
		totalMinedToken,
		lastDistributedEpoch,
		immatureTotalInfo,
		matureTotalInfo,
		individualRewardInfo,
		maturedEpoch)
}
