package rest

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	regTypes "github.com/stratosnet/stratos-chain/x/register/types"
	"github.com/stratosnet/stratos-chain/x/sds/types"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers sds-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/sds/file/upload", FileUploadRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/sds/prepay", PrepayRequestHandlerFn(cliCtx)).Methods("POST")

}

// FileUploadReq defines the properties of a file upload request's body.
type FileUploadReq struct {
	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req"`
	FileHash string       `json:"file_hash" yaml:"file_hash"`
	Uploader string       `json:"uploader" yaml:"uploader"`
}

// PrepayReq defines the properties of a prepay request's body.
type PrepayReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Int      `json:"amount" yaml:"amount"`
}

// FileUploadRequestHandlerFn - http request handler for file uploading.
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

		fileHash, err1 := hex.DecodeString(req.FileHash)
		if err1 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		uploader, err := sdk.AccAddressFromBech32(req.Uploader)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgUpload(fileHash, fromAddr, uploader)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// PrepayRequestHandlerFn - http request handler for prepay.
func PrepayRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PrepayReq
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
		prepayCoin := sdk.Coin{Denom: regTypes.DefaultBondDenom, Amount: req.Amount}
		coins := sdk.Coins{prepayCoin}
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgPrepay(fromAddr, coins)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
