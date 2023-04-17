package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/pot/types"
)

const (
	MatureCountPerBlock = 100
)

// RewardMatureAndSubSlashing mature rewards and deduct slashing for all nodes
func (k Keeper) RewardMatureAndSubSlashing(ctx sdk.Context) error {
	lastReportedEpoch := k.GetLastReportedEpoch(ctx)
	maturedEpoch := k.GetMaturedEpoch(ctx)
	// The first batch of reward is matured from the value of mature_epoch(param) + 1
	if maturedEpoch.IsZero() {
		maturedEpoch = sdk.NewInt(k.MatureEpoch(ctx))
		k.SetMaturedEpoch(ctx, maturedEpoch)
	}

	if lastReportedEpoch.LTE(maturedEpoch) {
		return nil
	}

	maturedIndividualKeys := make([][]byte, 0)

	matureStartEpochOffset := int64(1)
	matureEndEpochOffset := lastReportedEpoch.Sub(maturedEpoch).Int64()

	processCount := 1
	for i := matureStartEpochOffset; i <= matureEndEpochOffset; i++ {
		processingEpoch := sdk.NewInt(i).Add(maturedEpoch)
		totalSlashed := sdk.Coins{}

		isBreak := true
		k.IteratorIndividualReward(ctx, processingEpoch, func(walletAddress sdk.AccAddress, individualReward types.Reward) (stop bool) {

			// Stop iteration when executed wallet reaches MatureCountPerBlock && no new volume report is received
			if processCount > MatureCountPerBlock &&
				lastReportedEpoch.Equal(processingEpoch) {
				isBreak = true
				return true
			}

			// Mature reward
			oldMatureTotal := k.GetMatureTotalReward(ctx, walletAddress)
			oldImmatureTotal := k.GetImmatureTotalReward(ctx, walletAddress)
			immatureToMature := individualReward.RewardFromMiningPool.Add(individualReward.RewardFromTrafficPool...)

			// Deduct slashing amount from upcoming mature reward, don't need to deduct slashing from immatureTotal & individual
			remaining, deducted := k.registerKeeper.DeductSlashing(ctx, walletAddress, immatureToMature, k.RewardDenom(ctx))
			totalSlashed = totalSlashed.Add(deducted...)

			matureTotal := oldMatureTotal.Add(remaining...)
			immatureTotal := oldImmatureTotal.Sub(immatureToMature)

			processCount++

			k.SetMatureTotalReward(ctx, walletAddress, matureTotal)
			k.SetImmatureTotalReward(ctx, walletAddress, immatureTotal)
			maturedIndividualKeys = append(maturedIndividualKeys, types.GetIndividualRewardKey(walletAddress, processingEpoch))
			isBreak = false
			return false
		})

		// transfer deducted slashing amount
		err := k.transferTokensForMatureReward(ctx, totalSlashed)
		if err != nil {
			return err
		}

		// when isBreak == false, means reward mature for processingEpoch is completed, update MaturedEpoch
		if !isBreak {
			k.SetMaturedEpoch(ctx, processingEpoch)
		}
	}

	// remove matured individual reward records
	for _, key := range maturedIndividualKeys {
		k.RemoveIndividualReward(ctx, key)
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
