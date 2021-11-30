package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/stratosnet/stratos-chain/x/pot/keeper"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/pot/report/epoch/{epoch}", getVolumeReportHandlerFn(cliCtx, keeper.QueryVolumeReport)).Methods("GET")
	r.HandleFunc("/pot/rewards/epoch/{epoch}", getPotRewardsByEpochHandlerFn(cliCtx, keeper.QueryPotRewardsByReportEpoch)).Methods("GET")
	r.HandleFunc("/pot/rewards/wallet/{walletAddress}", getPotRewardsByWalletAddrHandlerFn(cliCtx, keeper.QueryPotRewardsByWalletAddr)).Methods("GET")
}

func getPotRewardsByEpochHandlerFn(cliCtx context.CLIContext, queryPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and verify params
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		epochStr := mux.Vars(r)["epoch"]
		epoch, ok := sdk.NewIntFromString(epochStr)
		if !ok {
			return
		}

		walletAddressStr := ""
		if v := r.URL.Query().Get(RestWalletAddress); len(v) != 0 {
			walletAddressStr = v
		}
		walletAddress, err := sdk.AccAddressFromBech32(walletAddressStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.NewQueryPotRewardsByEpochParams(page, limit, epoch, walletAddress)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			cliCtx = cliCtx.WithHeight(queryHeight)
			rest.PostProcessResponse(w, cliCtx, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			cliCtx = cliCtx.WithHeight(queryHeight)
			rest.PostProcessResponse(w, cliCtx, err.Error())
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

		v := mux.Vars(r)["epoch"]
		if len(v) == 0 {
			return
		}
		//epoch, ok := validateEpoch(w, v)
		epoch, ok := sdk.NewIntFromString(v)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid epoch")
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.QueryWithData(route, []byte(epoch.String()))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// GET request handler to query potRewards info by walletAddr
func getPotRewardsByWalletAddrHandlerFn(cliCtx context.CLIContext, queryPath string) http.HandlerFunc {
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

		walletAddrStr := mux.Vars(r)["walletAddress"]
		walletAddr, err := sdk.AccAddressFromBech32(walletAddrStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var (
			queryHeight int64
		)

		if v := r.URL.Query().Get(RestHeight); len(v) != 0 {
			queryHeight, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := types.NewQueryPotRewardsByWalletAddrParams(page, limit, walletAddr, queryHeight)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(queryHeight)
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, currentHeight, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(currentHeight)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
