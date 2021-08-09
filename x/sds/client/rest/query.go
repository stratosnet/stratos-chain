package rest

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stratosnet/stratos-chain/x/sds/client/common"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerSdsQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		"/sds/simulatePrepay/{amtToPrepay}",
		SimulatePrelayHandlerFn(cliCtx, queryRoute),
	).Methods("GET")
}

// HTTP request handler to query the total rewards balance from all delegations
func SimulatePrelayHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		amtToPrepay, ok := checkAmtToPrepayVar(w, r)
		if !ok {
			return
		}
		resp, height, err := common.QuerySimulatePrepay(cliCtx, queryRoute, amtToPrepay)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var simulatePrepayOut sdk.Int
		err = simulatePrepayOut.UnmarshalJSON(resp)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, simulatePrepayOut)
	}
}

func checkAmtToPrepayVar(w http.ResponseWriter, r *http.Request) (sdk.Int, bool) {
	prepayAmtStr := mux.Vars(r)["amtToPrepay"]
	amtToPrepay, ok := sdk.NewIntFromString(prepayAmtStr)
	if ok != true {
		rest.WriteErrorResponse(w, http.StatusBadRequest, sdkerrors.ErrInvalidRequest.Error())
		return sdk.Int{}, false
	}
	return amtToPrepay, true
}
