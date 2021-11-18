package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

//registerTxRoutes registers pot-related REST Tx handlers to a router
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/pot/volume/report", volumeReportRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/pot/address/{nodeAddr}/rewards", withdrawPotRewardsHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/pot/foundation_deposit", foundationDepositHandlerFn(cliCtx)).Methods("POST")
}

type (
	foundationDepositReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
		Amount  string       `json:"amount" yaml:"amount"`
	}

	withdrawRewardsReq struct {
		BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
		Amount     string       `json:"amount" yaml:"amount"`
		TargetAddr string       `json:"target_addr" yaml:"target_addr"`
	}

	volumeReportReq struct {
		BaseReq         rest.BaseReq             `json:"base_req" yaml:"base_req"`
		NodesVolume     []types.SingleNodeVolume `json:"nodes_volume" yaml:"nodes_volume"`         // volume report
		Reporter        string                   `json:"reporter" yaml:"reporter"`                 // volume reporter
		Epoch           int64                    `json:"report_epoch" yaml:"report_epoch"`         // volume report epoch
		ReportReference string                   `json:"report_reference" yaml:"report_reference"` // volume report reference
	}
)

// volumeReportRequestHandlerFn rest API handler to create a volume report tx.
func volumeReportRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req volumeReportReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		reporterStr := req.Reporter
		reporter, err := sdk.AccAddressFromBech32(reporterStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		reportReference := req.ReportReference
		epoch := sdk.NewInt(req.Epoch)

		var nodesVolume []types.SingleNodeVolume
		for _, v := range req.NodesVolume {
			singleNodeVolume := types.NewSingleNodeVolume(v.NodeAddress, v.Volume)
			nodesVolume = append(nodesVolume, singleNodeVolume)
		}

		reporterOwner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgVolumeReport(nodesVolume, reporter, epoch, reportReference, reporterOwner)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// rest API handler Withdraw pot rewards
func withdrawPotRewardsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawRewardsReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// read and validate URL's variables
		amountStr := req.Amount
		amount, ok := checkAmountVar(w, r, amountStr)
		if !ok {
			return
		}
		nodeAddrStr := mux.Vars(r)["nodeAddr"]
		nodeAddr, ok := checkAccountAddressVar(w, r, nodeAddrStr)
		if !ok {
			return
		}

		targetAddrStr := req.TargetAddr
		targetAddr, ok := checkAccountAddressVar(w, r, targetAddrStr)
		if !ok {
			return
		}

		//TODO: Add targetAddr after NewMsgWithdraw updates
		fmt.Println("targetAddr", targetAddr)

		ownerAddrStr := req.BaseReq.From
		ownerAddr, ok := checkAccountAddressVar(w, r, ownerAddrStr)
		if !ok {
			return
		}

		msg := types.NewMsgWithdraw(sdk.NewCoin(types.DefaultRewardDenom, amount), nodeAddr, ownerAddr)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func foundationDepositHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req foundationDepositReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// read and validate URL's variables
		amountStr := req.Amount
		amount, err := sdk.ParseCoin(amountStr)
		if err != nil {
			return
		}

		fromStr := req.BaseReq.From
		fromAddr, ok := checkAccountAddressVar(w, r, fromStr)
		if !ok {
			return
		}

		msg := types.NewMsgFoundationDeposit(amount, fromAddr)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func checkAccountAddressVar(w http.ResponseWriter, r *http.Request, accountAddrStr string) (sdk.AccAddress, bool) {
	addr, err := sdk.AccAddressFromBech32(accountAddrStr)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return nil, false
	}

	return addr, true
}

func checkAmountVar(w http.ResponseWriter, r *http.Request, amountStr string) (sdk.Int, bool) {
	amount, ok := sdk.NewIntFromString(amountStr)
	if !ok {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid withdraw amount")
		return sdk.NewInt(0), false
	}
	return amount, true
}
