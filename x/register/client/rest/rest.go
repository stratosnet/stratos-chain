package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"
)

// REST Variable names
// nolint
const (
	RestNetworkAddr = "network"
	RestNumLimit    = "limit"
	RestMoniker     = "moniker"
	RestOwner       = "owner"
	RestQueryType   = "query_type"
)

// RegisterHandlers registers register-related REST handlers to a router
func RegisterHandlers(clientCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
}
