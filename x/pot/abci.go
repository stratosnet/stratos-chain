package pot

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/keeper"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// 	TODO: fill out if your application requires beginblock, if not you can delete this function
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) []abci.ValidatorUpdate {

	logger := k.Logger(ctx)

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic. ", "ErrMsg", r)
		}
	}()

	// Do not distribute rewards until the next block
	if !k.GetIsReadyToDistribute(ctx) && k.GetUnDistributedEpoch(ctx).GT(sdk.ZeroInt()) {
		k.SetIsReadyToDistribute(ctx, true)
	} else {
		// Start distribute reward if report found
		walletVolumes, found := k.GetUnDistributedReport(ctx)
		if found {
			epoch := k.GetUnDistributedEpoch(ctx)

			//distribute POT reward
			err := k.DistributePotReward(ctx, walletVolumes.Volumes, epoch)
			if err != nil {
				logger.Error("An error occurred while distributing the reward. ", "ErrMsg", err.Error())
			}

			// reset undistributed info after distribution
			k.SetUnDistributedReport(ctx, types.WalletVolumes{})
			k.SetUnDistributedEpoch(ctx, sdk.ZeroInt())
		}
	}

	// mature reward
	err := k.RewardMatureAndSubSlashing(ctx)
	if err != nil {
		logger.Error("An error occurred while distributing the reward. ", "ErrMsg", err.Error())
	}

	// reset total supply to 100M stos
	k.RestoreTotalSupply(ctx)
	return []abci.ValidatorUpdate{}
}
