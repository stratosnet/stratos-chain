package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/stratosnet/stratos-chain/x/register/keeper"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"net/http"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/register/resource-nodes", nodesWithParamsFn(cliCtx, keeper.QueryResourceNodeList)).Methods("GET")
	//r.HandleFunc(fmt.Sprintf("/register/resource-nodes/{%s}", RestNetworkID), resourceNodeListByNetworkIDHandlerFn(cliCtx, queryRoute)).Methods("GET")
	r.HandleFunc("/register/indexing-nodes", nodesWithParamsFn(cliCtx, keeper.QueryIndexingNodeList)).Methods("GET")
	r.HandleFunc("/register/staking", nodeStakingHandlerFn(cliCtx, keeper.QueryNodesTotalStakes)).Methods("GET")
	r.HandleFunc("/register/staking/address/{nodeAddress}", nodeStakingByNodeAddressFn(cliCtx, keeper.QueryNodeStakeByNodeAddress)).Methods("GET")
}

// GET request handler to query all resource/indexing nodes
func nodesWithParamsFn(cliCtx context.CLIContext, queryPath string) http.HandlerFunc {
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

		var (
			networkID string
			moniker   string
			ownerAddr sdk.AccAddress
		)

		networkID = r.URL.Query().Get(RestNetworkID)
		moniker = r.URL.Query().Get(RestMoniker)

		if v := r.URL.Query().Get(RestOwner); len(v) != 0 {
			ownerAddr, err = sdk.AccAddressFromBech32(v)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := keeper.NewQueryNodesParams(page, limit, networkID, moniker, ownerAddr)
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

// GET request handler to query nodes total staking info
func nodeStakingHandlerFn(cliCtx context.CLIContext, queryPath string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
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
func nodeStakingByNodeAddressFn(cliCtx context.CLIContext, queryPath string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		nodeAddress, ok := keeper.CheckNodeAddr(w, r)
		if !ok {
			return
		}

		params := keeper.NewQuerynodeStakingByNodeAddressParams(nodeAddress)
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

//func checkNodeAddr(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
//	NodeAddrStr := mux.Vars(r)["nodeAddress"]
//	//NodeAddr, err := typesTypes.GetPubKeyFromBech32(typesTypes.Bech32PubKeyTypeSdsP2PPub, NodeAddrStr)
//	NodeAddr, err := sdk.AccAddressFromBech32(NodeAddrStr)
//	if err != nil {
//		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid 'NodeAddress'.")
//		return nil, false
//	}
//	return NodeAddr.Bytes(), true
//}
