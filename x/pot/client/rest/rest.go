package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	//"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	"net/http"
)

// RegisterRoutes registers pot-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/pot/volume/report", VolumeReportRequestHandlerFn(cliCtx)).Methods("POST")

}

// VolumeReportReq defines the properties of a send request's body.
type VolumeReportReq struct {
	BaseReq         rest.BaseReq             `json:"base_req" yaml:"base_req"`
	NodesVolume     []types.SingleNodeVolume `json:"nodes_volume" yaml:"nodes_volume"`               // volume report
	Reporter        string                   `json:"volume_reporter" yaml:"volume_reporter"`         // volume reporter
	Epoch           int                      `json:"volume_report_epoch" yaml:"volume_report_epoch"` // volume report epoch
	ReportReference string                   `json:"volume_report_hash" yaml:"volume_report_hash"`   // volume report reference
}

// VolumeReportRequestHandlerFn - http request handler to send coins to a address.
func VolumeReportRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req VolumeReportReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		reporter, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		reportReference := req.ReportReference
		epoch := sdk.NewInt(int64(req.Epoch))

		var nodesVolume []types.SingleNodeVolume
		for _, v := range req.NodesVolume {
			singleNodeVolume := types.NewSingleNodeVolume(v.NodeAddress, v.Volume)
			nodesVolume = append(nodesVolume, singleNodeVolume)
		}

		msg := types.NewMsgVolumeReport(nodesVolume, reporter, epoch, reportReference)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
