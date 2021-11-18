package rest

import (
	"encoding/json"
	"net/http"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stratosnet/stratos-chain/x/sds/client/common"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerSdsQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		"/sds/simulatePrepay/{amtToPrepay}",
		SimulatePrepayHandlerFn(cliCtx, queryRoute),
	).Methods("GET")
	r.HandleFunc(
		"/sds/uozPrice",
		UozPriceHandlerFn(cliCtx, queryRoute),
	).Methods("GET")
	r.HandleFunc(
		"/sds/uozSupply",
		UozSupplyHandlerFn(cliCtx, queryRoute),
	).Methods("GET")
}

// HTTP request handler to query the simulated purchased amt of prepay
func SimulatePrepayHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
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

// HTTP request handler to query ongoing uoz price
func UozPriceHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		resp, height, err := common.QueryCurrUozPrice(cliCtx, queryRoute)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, resp)
	}
}

// HTTP request handler to query uoz supply details
func UozSupplyHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		resp, height, err := common.QueryUozSupply(cliCtx, queryRoute)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		type Supply struct {
			Remaining sdk.Int
			Total     sdk.Int
		}
		var uozSupply Supply
		err = json.Unmarshal(resp, &uozSupply)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, uozSupply)
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
