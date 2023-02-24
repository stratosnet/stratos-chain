package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/stratosnet/stratos-chain/x/sds/client/common"
	"github.com/stratosnet/stratos-chain/x/sds/types"
)

func sdsQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/sds/simulatePrepay/{amtToPrepay}", SimulatePrepayHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/sds/nozPrice", NozPriceHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/sds/nozSupply", NozSupplyHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/sds/params", sdsParamsHandlerFn(clientCtx, types.QueryParams)).Methods("GET")
}

// GET request handler to query params of POT module
func sdsParamsHandlerFn(clientCtx client.Context, queryPath string) http.HandlerFunc {

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

// SimulatePrepayHandlerFn HTTP request handler to query the simulated purchased amt of prepay
func SimulatePrepayHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		amtToPrepay, ok := checkAmtToPrepayVar(w, r)
		if !ok {
			return
		}
		resp, height, err := common.QuerySimulatePrepay(cliCtx, amtToPrepay)

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

// NozPriceHandlerFn HTTP request handler to query ongoing noz price
func NozPriceHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		resp, height, err := common.QueryCurrNozPrice(cliCtx)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var nozPrice sdk.Dec
		err = nozPrice.UnmarshalJSON(resp)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, nozPrice)
	}
}

// NozSupplyHandlerFn HTTP request handler to query noz supply details
func NozSupplyHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		resp, height, err := common.QueryNozSupply(cliCtx)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		type Supply struct {
			Remaining sdk.Int
			Total     sdk.Int
		}
		var nozSupply Supply
		err = json.Unmarshal(resp, &nozSupply)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, nozSupply)
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
