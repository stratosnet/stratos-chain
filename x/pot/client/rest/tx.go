package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	//"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// registerTxRoutes registers pot-related REST Tx handlers to a router
func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/pot/volume_report", volumeReportRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/pot/withdraw", withdrawPotRewardsHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/pot/foundation_deposit", foundationDepositHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/pot/slashing", slashingResourceNodeHandlerFn(cliCtx)).Methods("POST")
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
		BLSSignature    types.BaseBLSSignatureInfo `json:"bls_signature" yaml:"bls_signature"`       // bls signature
	}

	slashingResourceNodeReq struct {
		BaseReq        rest.BaseReq         `json:"base_req" yaml:"base_req"`
		Reporters      []stratos.SdsAddress `json:"reporters" yaml:"reporters"`             // reporter(sp node) p2p address
		ReporterOwner  []sdk.AccAddress     `json:"reporter_owner" yaml:"reporter_owner"`   // report(sp node) wallet address
		NetworkAddress stratos.SdsAddress   `json:"network_address" yaml:"network_address"` // p2p address of the pp node
		WalletAddress  sdk.AccAddress       `json:"wallet_address" yaml:"wallet_address"`   // wallet address of the pp node
		Slashing       int64                `json:"slashing" yaml:"slashing"`
		Suspend        bool                 `json:"suspend" yaml:"suspend"`
	}
)

// volumeReportRequestHandlerFn rest API handler to create a volume report tx.
func volumeReportRequestHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req volumeReportReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		reporterStr := req.Reporter
		reporter, err := stratos.SdsAddressFromBech32(reporterStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		reportReference := req.ReportReference
		epoch := sdk.NewInt(req.Epoch)

		var walletVolumes []types.SingleWalletVolume
		for _, v := range req.WalletVolumes {
			walletAddr, err := sdk.AccAddressFromBech32(v.WalletAddress)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			volumeStr := v.Volume.String()
			volume, ok := sdk.NewIntFromString(volumeStr)
			if !ok {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid volume")
				return
			}
			singleWalletVolume := types.NewSingleWalletVolume(walletAddr, volume)
			walletVolumes = append(walletVolumes, singleWalletVolume)
		}

		reporterOwner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		sig := req.BLSSignature

		pubKeys := make([][]byte, len(sig.PubKeys))
		for i, v := range sig.PubKeys {
			pubKeys[i] = []byte(v)
		}
		blsSignature := types.NewBLSSignatureInfo(pubKeys, []byte(sig.Signature), []byte(sig.TxData))

		msg := types.NewMsgVolumeReport(walletVolumes, reporter, epoch, reportReference, reporterOwner, blsSignature)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
		//utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// rest API handler Withdraw pot rewards
func withdrawPotRewardsHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawRewardsReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
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
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
		//utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func foundationDepositHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req foundationDepositReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
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

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
		//utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
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
	amount, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid withdraw amount")
		return sdk.Coins{}, false
	}
	return amount, true
}

func slashingResourceNodeHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req slashingResourceNodeReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		slashing := sdk.NewInt(req.Slashing)

		msg := types.NewMsgSlashingResourceNode(req.Reporters, req.ReporterOwner, req.NetworkAddress, req.WalletAddress, slashing, req.Suspend)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
		//utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
