package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (q Querier) ResourceNode(c context.Context, req *types.QueryResourceNodeRequest) (*types.QueryResourceNodeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.NetworkAddr == "" {
		return nil, status.Error(codes.InvalidArgument, types.ErrInvalidNetworkAddr.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)

	var params types.QueryNodesParams
	node, ok := q.GetResourceNode(ctx, params.NetworkAddr)
	if !ok {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, types.ErrNoResourceNodeFound.Error())
	}
	return &types.QueryResourceNodeResponse{Node: node}, nil
}

func (q Querier) IndexingNode(c context.Context, req *types.QueryIndexingNodeRequest) (*types.QueryIndexingNodeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.NetworkAddr == "" {
		return nil, status.Error(codes.InvalidArgument, types.ErrInvalidNetworkAddr.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)

	var params types.QueryNodesParams
	node, ok := q.GetIndexingNode(ctx, params.NetworkAddr)
	if !ok {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, types.ErrNoIndexingNodeFound.Error())
	}
	return &types.QueryIndexingNodeResponse{Node: node}, nil
}
