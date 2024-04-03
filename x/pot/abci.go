package pot

import (
	abci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/pot/keeper"
)

const CHECK_TOTAL_SUPPLY_INTERVAL = 10000

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// 	TODO: fill out if your application requires beginblock, if not you can delete this function
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) []abci.ValidatorUpdate {

	logger := k.Logger(ctx)

	// mature reward
	err := k.RewardMatureAndSubSlashing(ctx)
	if err != nil {
		logger.Error("An error occurred while distributing the reward. ", "ErrMsg", err.Error())
	}

	// reset total supply to 100M stos every 10k blocks
	if ctx.BlockHeight()%CHECK_TOTAL_SUPPLY_INTERVAL == 1 {
		logger.Info("start RestoreTotalSupply")
		k.RestoreTotalSupply(ctx)
	}
	return []abci.ValidatorUpdate{}
}
