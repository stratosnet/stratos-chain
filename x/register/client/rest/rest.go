package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// REST Variable names
// nolint
const (
	RestNetworkID = "network"
	RestNumLimit  = "limit"
	RestMoniker   = "moniker"
	RestOwner     = "owner"
)

// RegisterRoutes registers register-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerTxRoutes(cliCtx, r)
	registerQueryRoutes(cliCtx, r)
}
