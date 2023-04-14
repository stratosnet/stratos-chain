package pot

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/pot/keeper"
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

	err := k.RewardMatureAndSubSlashing(ctx)
	if err != nil {
		logger.Error("An error occurred while distributing the reward. ", "ErrMsg", err.Error())
	}

	return []abci.ValidatorUpdate{}
}
