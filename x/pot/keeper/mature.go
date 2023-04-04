package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/pot/types"
)

const (
	MatureCountPerBlock = 100
)

// mature rewards and deduct slashing for all nodes
func (k Keeper) RewardMatureAndSubSlashing(ctx sdk.Context) error {
	maturedEpoch := k.GetMaturedEpoch(ctx)
	lastReportedEpoch := k.GetLastReportedEpoch(ctx)

	matureStartEpochOffset := int64(1)
	matureEndEpochOffset := lastReportedEpoch.Sub(maturedEpoch).Int64()

	for i := matureStartEpochOffset; i <= matureEndEpochOffset; i++ {
		processingEpoch := sdk.NewInt(i).Add(maturedEpoch)
		totalSlashed := sdk.Coins{}

		individualIndex := sdk.ZeroInt()
		individualStartMature := k.GetNextMatureIndividualIndex(ctx)
		count := 0
		isBreak := false
		k.IteratorIndividualReward(ctx, processingEpoch, func(walletAddress sdk.AccAddress, individualReward types.Reward) (stop bool) {
			if individualIndex.LT(individualStartMature) {
				individualIndex = individualIndex.Add(sdk.OneInt())
				return false
			}

			// stop iteration when executed wallet reaches MatureCountPerBlock
			if count > MatureCountPerBlock {
				isBreak = true
				return true
			}

			oldMatureTotal := k.GetMatureTotalReward(ctx, walletAddress)
			oldImmatureTotal := k.GetImmatureTotalReward(ctx, walletAddress)
			immatureToMature := individualReward.RewardFromMiningPool.Add(individualReward.RewardFromTrafficPool...)

			//deduct slashing amount from upcoming mature reward, don't need to deduct slashing from immatureTotal & individual
			remaining, deducted := k.registerKeeper.DeductSlashing(ctx, walletAddress, immatureToMature, k.RewardDenom(ctx))
			totalSlashed = totalSlashed.Add(deducted...)

			matureTotal := oldMatureTotal.Add(remaining...)
			immatureTotal := oldImmatureTotal.Sub(immatureToMature)

			individualIndex = individualIndex.Add(sdk.OneInt())
			count++

			k.SetMatureTotalReward(ctx, walletAddress, matureTotal)
			k.SetImmatureTotalReward(ctx, walletAddress, immatureTotal)
			k.SetNextMatureIndividualIndex(ctx, individualIndex)

			return false
		})

		err := k.transferTokensForMatureReward(ctx, totalSlashed)
		if err != nil {
			return err
		}

		if !isBreak {
			k.SetNextMatureIndividualIndex(ctx, sdk.ZeroInt())
			k.SetMaturedEpoch(ctx, processingEpoch)
		}
	}

	return nil
}

func (k Keeper) transferTokensForMatureReward(ctx sdk.Context, totalSlashed sdk.Coins) error {

	// [TLC] [TotalRewardPool -> Distribution] Transfer slashed reward to FeePool.CommunityPool
	totalRewardPoolAccAddr := k.accountKeeper.GetModuleAddress(types.TotalRewardPool)
	err := k.distrKeeper.FundCommunityPool(ctx, totalSlashed, totalRewardPoolAccAddr)
	if err != nil {
		return err
	}

	return nil
}
