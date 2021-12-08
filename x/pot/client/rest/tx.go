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
	r.HandleFunc("/pot/volume_report", volumeReportRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/pot/withdraw", withdrawPotRewardsHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/pot/foundation_deposit", foundationDepositHandlerFn(cliCtx)).Methods("POST")
}

type (
	foundationDepositReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
		Amount  string       `json:"amount" yaml:"amount"`
	}

	withdrawRewardsReq struct {
		BaseReq       rest.BaseReq `json:"base_req" yaml:"base_req"`
		Amount        string       `json:"amount" yaml:"amount"`
		TargetAddress string       `json:"target_address" yaml:"target_address"`
	}

	volumeReportReq struct {
		BaseReq         rest.BaseReq               `json:"base_req" yaml:"base_req"`
		WalletVolumes   []types.SingleWalletVolume `json:"wallet_volumes" yaml:"wallet_volumes"`     // volume report
		Reporter        string                     `json:"reporter" yaml:"reporter"`                 // volume reporter
		Epoch           int64                      `json:"epoch" yaml:"epoch"`                       // volume report epoch
		ReportReference string                     `json:"report_reference" yaml:"report_reference"` // volume report reference
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

		var walletVolumes []types.SingleWalletVolume
		for _, v := range req.WalletVolumes {
			singleWalletVolume := types.NewSingleWalletVolume(v.WalletAddress, v.Volume)
			walletVolumes = append(walletVolumes, singleWalletVolume)
		}

		reporterOwner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgVolumeReport(walletVolumes, reporter, epoch, reportReference, reporterOwner)
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

		targetAddrStr := req.TargetAddress
		targetAddr, ok := checkAccountAddressVar(w, r, targetAddrStr)
		if !ok {
			return
		}

		//TODO: Add targetAddr after NewMsgWithdraw updates
		fmt.Println("targetAddr", targetAddr)

		walletAddrStr := req.BaseReq.From
		walletAddr, ok := checkAccountAddressVar(w, r, walletAddrStr)
		if !ok {
			return
		}

		msg := types.NewMsgWithdraw(amount, walletAddr, targetAddr)
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
		amount, ok := checkAmountVar(w, r, amountStr)
		if !ok {
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

func checkAmountVar(w http.ResponseWriter, r *http.Request, amountStr string) (sdk.Coins, bool) {
	amount, err := sdk.ParseCoins(amountStr)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid withdraw amount")
		return sdk.Coins{}, false
	}
	return amount, true
}
