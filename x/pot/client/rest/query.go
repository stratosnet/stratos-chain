package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/stratosnet/stratos-chain/x/pot/keeper"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	"net/http"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/pot/rewards/{epoch}", getPotRewardsHandlerFn(cliCtx, keeper.QueryPotRewards)).Methods("GET")
	r.HandleFunc("/pot/rewards/{epoch}/{NodeWalletAddress}", getPotRewardsByNodeWalletAddrHandlerFn(cliCtx, keeper.QueryPotRewards)).Methods("GET")
	r.HandleFunc("/pot/report/{epoch}", getVolumeReportHandlerFn(cliCtx, keeper.QueryVolumeReport)).Methods("GET")
}

// GET request handler to query potRewards info
func getPotRewardsHandlerFn(cliCtx context.CLIContext, queryPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		var epoch sdk.Int
		epoch, ok = checkEpoch(w, r)
		if !ok {
			return
		}

		params := keeper.NewQueryPotRewardsParams(page, limit, sdk.AccAddress{}, epoch)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// GET request handler to query potRewards info by nodeWalletAddr
func getPotRewardsByNodeWalletAddrHandlerFn(cliCtx context.CLIContext, queryPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		nodeWalletAddr, ok := checkNodeWalletAddr(w, r)
		if !ok {
			return
		}
		epoch, ok := checkEpoch(w, r)
		if !ok {
			return
		}

		params := keeper.NewQueryPotRewardsParams(1, 1, nodeWalletAddr, epoch)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// GET request handler to query Volume report info
func getVolumeReportHandlerFn(cliCtx context.CLIContext, queryPath string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		epoch, ok := checkEpoch(w, r)
		if !ok {
			return
		}
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.QueryWithData(route, []byte(epoch.String()))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var resp types.ReportInfo
		resp = types.NewReportInfo(epoch, string(res))

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, resp)
	}
}

func checkEpoch(w http.ResponseWriter, r *http.Request) (sdk.Int, bool) {
	epochStr := mux.Vars(r)["epoch"]
	epoch, ok := sdk.NewIntFromString(epochStr)
	if ok != true {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid 'epoch'.")
		return sdk.Int{}, false
	}
	return epoch, true
}

func checkNodeWalletAddr(w http.ResponseWriter, r *http.Request) (sdk.AccAddress, bool) {
	NodeWalletAddrStr := mux.Vars(r)["NodeWalletAddress"]
	NodeWalletAddr, err := sdk.AccAddressFromBech32(NodeWalletAddrStr)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid 'NodeWalletAddress'.")
		return sdk.AccAddress{}, false
	}
	return NodeWalletAddr, true
}
