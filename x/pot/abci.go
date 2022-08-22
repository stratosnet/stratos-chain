package pot

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
	// abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// 	TODO: fill out if your application requires beginblock, if not you can delete this function
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) []abci.ValidatorUpdate {

	walletVolumes, found := k.GetUnhandledReport(ctx)
	if !found {
		return []abci.ValidatorUpdate{}
	}
	epoch := k.GetUnhandledEpoch(ctx)
	logger := k.Logger(ctx)

	//distribute POT reward
	_, err := k.DistributePotReward(ctx, walletVolumes, epoch)
	if err != nil {
		logger.Error("An error occurred while distributing the reward. ", err)
	}

	k.SetUnhandledReport(ctx, nil)
	k.SetUnhandledEpoch(ctx, sdk.Int{})

	return []abci.ValidatorUpdate{}
}
