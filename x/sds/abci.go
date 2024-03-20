package sds

import (
	abci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/sds/keeper"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock, _ keeper.Keeper) {
	// 	TODO: fill out if your application requires beginBlock, if not you can delete this function
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(_ sdk.Context, _ abci.RequestEndBlock, _ keeper.Keeper) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
