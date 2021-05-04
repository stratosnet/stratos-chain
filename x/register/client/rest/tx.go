package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"net/http"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/register/createResourceNode",
		postCreateResourceNodeHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/register/createIndexingNode",
		postCreateIndexingNodeHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/register/removeResourceNode",
		postRemoveResourceNodeHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/register/removeIndexingNode",
		postRemoveIndexingNodeHandlerFn(cliCtx),
	).Methods("POST")
}

type (
	CreateResourceNodeRequest struct {
		BaseReq        rest.BaseReq      `json:"base_req" yaml:"base_req"`
		NetworkAddress string            `json:"network_address" yaml:"network_address"` // in bech32
		PubKey         string            `json:"pubkey" yaml:"pubkey"`                   // in bech32
		Amount         sdk.Coin          `json:"amount" yaml:"amount"`
		Description    types.Description `json:"description" yaml:"description"`
	}

	CreateIndexingNodeRequest struct {
		BaseReq        rest.BaseReq      `json:"base_req" yaml:"base_req"`
		NetworkAddress string            `json:"network_address" yaml:"network_address"` // in bech32
		PubKey         string            `json:"pubkey" yaml:"pubkey"`                   // in bech32
		Amount         sdk.Coin          `json:"amount" yaml:"amount"`
		Description    types.Description `json:"description" yaml:"description"`
	}

	RemoveResourceNodeRequest struct {
		BaseReq             rest.BaseReq `json:"base_req" yaml:"base_req"`
		ResourceNodeAddress string       `json:"resource_node_address" yaml:"resource_node_address"` // in bech32
	}

	RemoveIndexingNodeRequest struct {
		BaseReq             rest.BaseReq `json:"base_req" yaml:"base_req"`
		IndexingNodeAddress string       `json:"indexing_node_address" yaml:"indexing_node_address"` // in bech32
	}
)

func postCreateResourceNodeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateResourceNodeRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		pubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, req.PubKey)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ownerAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCreateResourceNode(req.NetworkAddress, pubkey, req.Amount, ownerAddr, req.Description)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postCreateIndexingNodeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateIndexingNodeRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		pubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, req.PubKey)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ownerAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCreateIndexingNode(req.NetworkAddress, pubkey, req.Amount, ownerAddr, req.Description)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postRemoveResourceNodeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RemoveResourceNodeRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		nodeAddr, err := sdk.AccAddressFromBech32(req.ResourceNodeAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ownerAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgRemoveResourceNode(nodeAddr, ownerAddr)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postRemoveIndexingNodeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RemoveIndexingNodeRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		nodeAddr, err := sdk.AccAddressFromBech32(req.IndexingNodeAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ownerAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgRemoveIndexingNode(nodeAddr, ownerAddr)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
