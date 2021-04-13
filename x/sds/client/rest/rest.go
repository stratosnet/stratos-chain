package rest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"github.com/stratosnet/stratos-chain/x/sds/types"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers sds-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/sds/file/upload", FileUploadRequestHandlerFn(cliCtx)).Methods("POST")

}

// FileUploadReq defines the properties of a send request's body.
type FileUploadReq struct {
	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req"`
	FileHash string       `json:"file_hash" yaml:"file_hash"`
}

// FileUploadRequestHandlerFn - http request handler to send coins to a address.
func FileUploadRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req FileUploadReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fileHash, err1 := sdk.AccAddressFromHex(req.FileHash)
		if err1 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgUpload(fileHash, fromAddr)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
