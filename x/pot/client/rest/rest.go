package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"
)

const (
	RestEpoch  = "epoch"
	RestHeight = "height"
)

// RegisterRoutes registers pot-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	registerTxRoutes(clientCtx, r)
	registerQueryRoutes(clientCtx, r)
}
