package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktrest "github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

// RegisterHandlers registers register-related REST handlers to a router
func RegisterHandlers(clientCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	sdsQueryRoutes(clientCtx, r)
	sdsTxHandlers(clientCtx, r)
}

// FileUploadReq defines the properties of a file upload request's body.
type FileUploadReq struct {
	BaseReq  sdktrest.BaseReq `json:"base_req" yaml:"base_req"`
	Reporter string           `json:"reporter" yaml:"reporter"`
	FileHash string           `json:"file_hash" yaml:"file_hash"`
	Uploader string           `json:"uploader" yaml:"uploader"`
}

// PrepayReq defines the properties of a prepay request's body.
type PrepayReq struct {
	BaseReq sdktrest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coins        `json:"amount" yaml:"amount"`
}
