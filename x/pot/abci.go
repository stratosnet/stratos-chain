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

	// Do not distribute rewards until the next block
	if !k.GetIsReadyToDistributeReward(ctx) && k.GetUnhandledEpoch(ctx).GT(sdk.ZeroInt()) {
		k.SetIsReadyToDistributeReward(ctx, true)
		return []abci.ValidatorUpdate{}
	}

	walletVolumes, found := k.GetUnhandledReport(ctx)
	if !found {
		return []abci.ValidatorUpdate{}
	}
	epoch := k.GetUnhandledEpoch(ctx)
	logger := k.Logger(ctx)

	//distribute POT reward
	_, err := k.DistributePotReward(ctx, walletVolumes.Volumes, epoch)
	if err != nil {
		logger.Error("An error occurred while distributing the reward. ", "ErrMsg", err.Error())
	}

	k.SetUnhandledReport(ctx, types.WalletVolumes{})
	k.SetUnhandledEpoch(ctx, sdk.ZeroInt())

	return []abci.ValidatorUpdate{}
}
