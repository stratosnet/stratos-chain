package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/keeper"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/register/resource-nodes", nodesWithParamsFn(clientCtx, keeper.QueryResourceNodeByNetworkAddr)).Methods("GET")
	r.HandleFunc("/register/meta-nodes", nodesWithParamsFn(clientCtx, keeper.QueryMetaNodeByNetworkAddr)).Methods("GET")
	r.HandleFunc("/register/staking", nodeStakingHandlerFn(clientCtx, keeper.QueryNodesTotalStakes)).Methods("GET")
	r.HandleFunc("/register/staking/address/{nodeAddress}", nodeStakingByNodeAddressFn(clientCtx, keeper.QueryNodeStakeByNodeAddr)).Methods("GET")
	r.HandleFunc("/register/staking/owner/{ownerAddress}", nodeStakingByOwnerFn(clientCtx, keeper.QueryNodeStakeByOwner)).Methods("GET")
	r.HandleFunc("/register/params", registerParamsHandlerFn(clientCtx, keeper.QueryRegisterParams)).Methods("GET")
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

		moniker = r.URL.Query().Get(RestMoniker)

		if v := r.URL.Query().Get(RestOwner); len(v) != 0 {
			ownerAddr, err = sdk.AccAddressFromBech32(v)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		if v := r.URL.Query().Get(RestNetworkAddr); len(v) != 0 {
			networkAddr, err = stratos.SdsAddressFromBech32(v)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
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

// GET request handler to query nodes total staking info
func nodeStakingHandlerFn(clientCtx client.Context, queryPath string) http.HandlerFunc {

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

// GET request handler to query node staking info
func nodeStakingByNodeAddressFn(cliCtx client.Context, queryPath string) http.HandlerFunc {

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

		params := types.NewQueryNodeStakingParams(nodeAddress, queryType)
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

// GET request handler to query nodes staking info by Node wallet address
func nodeStakingByOwnerFn(cliCtx client.Context, queryPath string) http.HandlerFunc {

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
