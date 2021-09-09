package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

const (
	RestNodeAddress  = "nodeAddress"
	RestOwnerAddress = "owner"
	RestHeight       = "height"
	RestEpoch        = "epoch"
)

// RegisterRoutes registers pot-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerTxRoutes(cliCtx, r)
	registerQueryRoutes(cliCtx, r)
}
