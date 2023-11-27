package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/rest"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/keeper"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/register/resource-node/{nodeAddress}", nodesWithParamsFn(clientCtx, keeper.QueryResourceNodeByNetworkAddr)).Methods("GET")
	r.HandleFunc("/register/meta-node/{nodeAddress}", nodesWithParamsFn(clientCtx, keeper.QueryMetaNodeByNetworkAddr)).Methods("GET")
	r.HandleFunc("/register/deposit", nodeDepositHandlerFn(clientCtx, keeper.QueryNodesDepositTotal)).Methods("GET")
	r.HandleFunc("/register/deposit/address/{nodeAddress}", nodeDepositByNodeAddressFn(clientCtx, keeper.QueryNodeDepositByNodeAddr)).Methods("GET")
	r.HandleFunc("/register/deposit/owner/{ownerAddress}", nodeDepositByOwnerFn(clientCtx, keeper.QueryNodeDepositByOwner)).Methods("GET")
	r.HandleFunc("/register/params", registerParamsHandlerFn(clientCtx, keeper.QueryRegisterParams)).Methods("GET")
	r.HandleFunc("/register/resource-count", resourceNodesCountFn(clientCtx, keeper.QueryResourceNodesCount)).Methods("GET")
	r.HandleFunc("/register/meta-count", metaNodesCountFn(clientCtx, keeper.QueryMetaNodesCount)).Methods("GET")
	r.HandleFunc("/register/remaining-ozone-limit", remainingOzoneLimitFn(clientCtx, keeper.QueryRemainingOzoneLimit)).Methods("GET")
}

// GET request handler to query total number of bonded resource nodes
func resourceNodesCountFn(clientCtx client.Context, queryPath string) http.HandlerFunc {

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

// GET request handler to query total number of bonded resource nodes
func metaNodesCountFn(clientCtx client.Context, queryPath string) http.HandlerFunc {

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

// GET request handler to query params of Register module
func registerParamsHandlerFn(clientCtx client.Context, queryPath string) http.HandlerFunc {

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

// GET request handler to query all resource/meta nodes
func nodesWithParamsFn(clientCtx client.Context, queryPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		var (
			networkAddr stratos.SdsAddress
			moniker     string
			ownerAddr   sdk.AccAddress
		)

		networkAddrStr := mux.Vars(r)["nodeAddress"]
		networkAddr, ok = keeper.CheckSdsAddr(w, r, networkAddrStr)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		countTotal, err := strconv.ParseBool(r.FormValue("count_total"))
		if err != nil {
			countTotal = true
		}

		reverse, err := strconv.ParseBool(r.FormValue("reverse"))
		if err != nil {
			reverse = false
		}
		offset := page * limit

		NodesPageRequest := query.PageRequest{Offset: uint64(offset), Limit: uint64(limit), CountTotal: countTotal, Reverse: reverse}

		params := types.NewQueryNodesParams(networkAddr, moniker, ownerAddr, NodesPageRequest)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// GET request handler to query nodes total deposit info
func nodeDepositHandlerFn(clientCtx client.Context, queryPath string) http.HandlerFunc {

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

// GET request handler to query node deposit info
func nodeDepositByNodeAddressFn(cliCtx client.Context, queryPath string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		NodeAddrStr := mux.Vars(r)["nodeAddress"]
		nodeAddress, ok := keeper.CheckSdsAddr(w, r, NodeAddrStr)
		if !ok {
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		var (
			err       error
			queryType int64
		)
		if v := r.URL.Query().Get(RestQueryType); len(v) != 0 {
			queryType, err = strconv.ParseInt(v, 10, 64)
			if err != nil || queryType < types.QueryType_All || queryType > types.QueryType_PP {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else {
			queryType = 0
		}

		params := types.NewQueryNodeDepositParams(nodeAddress, queryType)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// GET request handler to query nodes deposit info by Node wallet address
func nodeDepositByOwnerFn(cliCtx client.Context, queryPath string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ownerAddressStr := mux.Vars(r)["ownerAddress"]
		ownerAddress, ok := keeper.CheckAccAddr(w, r, ownerAddressStr)
		if !ok {
			return
		}

		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		countTotal, err := strconv.ParseBool(r.FormValue("count_total"))
		if err != nil {
			countTotal = true
		}

		reverse, err := strconv.ParseBool(r.FormValue("reverse"))
		if err != nil {
			reverse = false
		}

		offset := (page - 1) * limit

		if limit <= 0 {
			limit = types.QueryDefaultLimit
		}

		NodesPageRequest := query.PageRequest{Offset: uint64(offset), Limit: uint64(limit), CountTotal: countTotal, Reverse: reverse}
		params := types.NewQueryNodesParams(nil, "", ownerAddress, NodesPageRequest)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, queryPath)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// GET request handler to query remaining ozone limit
func remainingOzoneLimitFn(clientCtx client.Context, queryPath string) http.HandlerFunc {
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
