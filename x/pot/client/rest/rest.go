package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/gorilla/mux"
)

const (
	RestWalletAddress = "wallet_address"
	RestEpoch         = "epoch"
	RestHeight        = "height"
)

// RegisterRoutes registers pot-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	registerTxRoutes(clientCtx, r)
	registerQueryRoutes(clientCtx, r)
}
