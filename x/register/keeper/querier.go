package keeper

import (
	"github.com/stratosnet/stratos-chain/x/register/types"
	"strings"

	// this line is used by starport scaffolding # 1
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	QueryResourceNodeList      = "resource_nodes"
	QueryResourceNodeByMoniker = "resource_nodes_moniker"
	QueryIndexingNodeList      = "indexing_nodes"
	QueryIndexingNodeByMoniker = "indexing_nodes_moniker"
	QueryNetworkSet            = "networks"
)

// NewQuerier creates a new querier for register clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryResourceNodeList:
			return GetResourceNodes(ctx, req, k)
		case QueryIndexingNodeList:
			return GetIndexingNodes(ctx, req, k)
		case QueryNetworkSet:
			return GetNetworkSet(ctx, k)
		case QueryResourceNodeByMoniker:
			return GetResourceNodesByMoniker(ctx, req, k)
		case QueryIndexingNodeByMoniker:
			return GetIndexingNodesByMoniker(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown register query endpoint "+req.String()+string(req.Data))
		}
	}
}

func GetResourceNodesByMoniker(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	nodeList, err := k.GetResourceNodeListByMoniker(ctx, string(req.Data))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return types.ModuleCdc.MustMarshalJSON(nodeList), nil
}

func GetIndexingNodesByMoniker(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	nodeList, err := k.GetIndexingNodeListByMoniker(ctx, string(req.Data))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return types.ModuleCdc.MustMarshalJSON(nodeList), nil
}

// GetResourceNodes fetches all resource nodes by network address.
func GetResourceNodes(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	nodeList, err := k.GetResourceNodeList(ctx, string(req.Data))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return types.ModuleCdc.MustMarshalJSON(nodeList), nil
}

// GetIndexingNodes fetches all indexing nodes by network address.
func GetIndexingNodes(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	nodeList, err := k.GetIndexingNodeList(ctx, string(req.Data))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return types.ModuleCdc.MustMarshalJSON(nodeList), nil
}

// GetNetworkSet fetches all network addresses.
func GetNetworkSet(ctx sdk.Context, k Keeper) ([]byte, error) {
	networks, err := k.GetNetworks(ctx, k)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return []byte(strings.TrimSpace(string(networks))), nil
}
