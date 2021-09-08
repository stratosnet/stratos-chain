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
	"strconv"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/pot/rewards/epoch/{epoch}", getPotRewardsByEpochHandlerFn(cliCtx, keeper.QueryPotRewardsByEpoch)).Methods("GET")
	//r.HandleFunc("/pot/rewards/owner/{ownerAddress}", getPotRewardsByOwnerHandlerFn(cliCtx, keeper.QueryPotRewardsByOwner)).Methods("GET")
	r.HandleFunc("/pot/report/epoch/{epoch}", getVolumeReportHandlerFn(cliCtx, keeper.QueryVolumeReport)).Methods("GET")
}

func getPotRewardsByEpochHandlerFn(cliCtx context.CLIContext, queryPath string) http.HandlerFunc {
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

		epochStr := mux.Vars(r)["epoch"]
		epoch, ok := checkEpoch(epochStr)
		if !ok {
			return
		}

		ownerAddressStr := ""
		if v := r.URL.Query().Get(RestOwnerAddress); len(v) != 0 {
			ownerAddressStr = v
		}
		ownerAddress, err := sdk.AccAddressFromBech32(ownerAddressStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := keeper.NewQueryPotRewardsByepochParams(page, limit, ownerAddress, epoch)
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
func getPotRewardsByOwnerHandlerFn(cliCtx context.CLIContext, queryPath string) http.HandlerFunc {
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

		ownerAddrStr := mux.Vars(r)["ownerAddress"]
		ownerAddr, err := sdk.AccAddressFromBech32(ownerAddrStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var height int64
		if v := r.URL.Query().Get(RestHeight); len(v) != 0 {
			height, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			}
		}

		params := keeper.NewQueryPotRewardsByOwnerParams(page, limit, ownerAddr, height)
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

		v := mux.Vars(r)["epoch"]
		epoch, ok := checkEpoch(v)
		if len(v) == 0 || !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid epoch.")
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.QueryWithData(route, []byte(epoch.String()))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkEpoch(epochStr string) (sdk.Int, bool) {
	//epochStr := mux.Vars(r)["epoch"]
	epoch, ok := sdk.NewIntFromString(epochStr)
	if ok != true {
		//rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid 'epoch'.")
		return sdk.Int{}, false
	}
	return epoch, true
}
