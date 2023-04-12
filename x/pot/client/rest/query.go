package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/pot/report/epoch/{epoch}", getVolumeReportHandlerFn(clientCtx, types.QueryVolumeReport)).Methods("GET")
	r.HandleFunc("/pot/rewards/epoch/{epoch}", getIndividualRewardsByEpochHandlerFn(clientCtx, types.QueryIndividualRewardsByReportEpoch)).Methods("GET")
	r.HandleFunc("/pot/rewards/wallet/{walletAddress}", getRewardsByWalletAddrHandlerFn(clientCtx, types.QueryRewardsByWalletAddr)).Methods("GET")
	r.HandleFunc("/pot/slashing/{walletAddress}", getSlashingByWalletAddressHandlerFn(clientCtx, types.QuerySlashingByWalletAddr)).Methods("GET")
	r.HandleFunc("/pot/params", potParamsHandlerFn(clientCtx, types.QueryPotParams)).Methods("GET")
	r.HandleFunc("/pot/total-mined-token", getTotalMinedTokenHandlerFn(clientCtx, types.QueryTotalMinedToken)).Methods("GET")
}

// GET request handler to query params of POT module
func potParamsHandlerFn(clientCtx client.Context, queryPath string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.Query(route)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getIndividualRewardsByEpochHandlerFn(clientCtx client.Context, queryPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and verify params
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		epochStr := mux.Vars(r)["epoch"]
		epoch, ok := sdk.NewIntFromString(epochStr)
		if !ok {
			return
		}

		params := types.NewQueryIndividualRewardsByEpochParams(page, limit, epoch)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.PostProcessResponse(w, cliCtx, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.PostProcessResponse(w, cliCtx, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// GET request handler to query Volume report info
func getVolumeReportHandlerFn(clientCtx client.Context, queryPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
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
func getRewardsByWalletAddrHandlerFn(clientCtx client.Context, queryPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		walletAddrStr := mux.Vars(r)["walletAddress"]
		walletAddr, err := sdk.AccAddressFromBech32(walletAddrStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		queryEpoch := sdk.ZeroInt()
		queryHeight := cliCtx.Height

		if v := r.URL.Query().Get(RestEpoch); len(v) != 0 {
			queryEpoch, ok = sdk.NewIntFromString(v)
			if !ok {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid epoch")
				return
			}
		}

		if v := r.URL.Query().Get(RestHeight); len(v) != 0 {
			queryHeight, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := types.NewQueryRewardsByWalletAddrParams(walletAddr, queryHeight, queryEpoch)

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(queryHeight)
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

func getSlashingByWalletAddressHandlerFn(clientCtx client.Context, queryPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		v := mux.Vars(r)["walletAddress"]
		if len(v) == 0 {
			return
		}
		walletAddr, err := sdk.AccAddressFromBech32(v)

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.QueryWithData(route, []byte(walletAddr.String()))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getTotalMinedTokenHandlerFn(clientCtx client.Context, queryPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.Query(route)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
