package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/stratosnet/stratos-chain/x/pot/keeper"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	"net/http"
	"strconv"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/pot/rewards", getPotRewardsHandlerFn(cliCtx, keeper.QueryPotRewards)).Methods("GET")
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

		var nodeAddr sdk.AccAddress
		var epoch int64

		if v := r.URL.Query().Get(RestNodeAddress); len(v) != 0 {
			nodeAddr, err = sdk.AccAddressFromBech32(v)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else {
			nodeAddr = sdk.AccAddress{}
		}

		if v := r.URL.Query().Get(RestEpoch); len(v) != 0 {
			epoch, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "'epoch' is incorrect.")
				return
			}
		} else {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "missing query 'epoch'.")
		}

		params := keeper.NewQueryPotRewardsParams(page, limit, nodeAddr, epoch)
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

		epoch, ok := checkEpochValue(w, r)
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

func checkEpochValue(w http.ResponseWriter, r *http.Request) (sdk.Int, bool) {
	epochStr := mux.Vars(r)["epoch"]
	epoch, ok := sdk.NewIntFromString(epochStr)
	if ok != true {
		rest.WriteErrorResponse(w, http.StatusBadRequest, sdkerrors.ErrInvalidRequest.Error())
		return sdk.Int{}, false
	}
	return epoch, true
}
