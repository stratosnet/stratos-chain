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

	//// RegisterRoutes registers sds-related REST handlers to a router
	//func RegisterRoutes(clientCtx client.Context, r *mux.Router, queryRoute string) {
	//	r.HandleFunc("/sds/file/upload", FileUploadRequestHandlerFn(cliCtx)).Methods("POST")
	//	r.HandleFunc("/sds/prepay", PrepayRequestHandlerFn(cliCtx)).Methods("POST")
	//	registerSdsQueryRoutes(cliCtx, r, queryRoute)
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

//// FileUploadRequestHandlerFn - http request handler for file uploading.
//func FileUploadRequestHandlerFn(clientCtx client.Context, queryPath string) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
//		if !ok {
//			return
//		}
//
//		var req FileUploadReq
//		if !rest.ReadRESTReq(w, r, clientCtx.Codec, &req) {
//			return
//		}
//
//		req.BaseReq = req.BaseReq.Sanitize()
//		if !req.BaseReq.ValidateBasic(w) {
//			return
//		}
//
//		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		reporter, err := stratos.SdsAddressFromBech32(req.Reporter)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		fileHash := req.FileHash
//		_, err = hex.DecodeString(fileHash)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		uploader, err := sdk.AccAddressFromBech32(req.Uploader)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		msg := types.NewMsgUpload(fileHash, fromAddr, reporter, uploader)
//		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
//	}
//}
//
//// PrepayRequestHandlerFn - http request handler for prepay.
//func PrepayRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var req PrepayReq
//		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
//			return
//		}
//
//		req.BaseReq = req.BaseReq.Sanitize()
//		if !req.BaseReq.ValidateBasic(w) {
//			return
//		}
//
//		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		msg := types.NewMsgPrepay(fromAddr, req.Amount)
//		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
//	}
//}
