package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (q Querier) ResourceNode(c context.Context, req *types.QueryResourceNodeRequest) (*types.QueryResourceNodeResponse, error) {
	if req == nil {
		return &types.QueryResourceNodeResponse{}, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GetNetworkAddr() == "" {
		return &types.QueryResourceNodeResponse{}, status.Error(codes.InvalidArgument, " Network address cannot be empty")
	}

	networkAddr, err := stratos.SdsAddressFromBech32(req.GetNetworkAddr())
	if err != nil {
		return &types.QueryResourceNodeResponse{}, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	node, found := q.GetResourceNode(ctx, networkAddr)
	if !found {
		return &types.QueryResourceNodeResponse{}, status.Errorf(codes.NotFound, "network address %s not found", req.NetworkAddr)
	}

	return &types.QueryResourceNodeResponse{Node: &node}, nil
}

func (q Querier) MetaNode(c context.Context, req *types.QueryMetaNodeRequest) (*types.QueryMetaNodeResponse, error) {
	if req == nil {
		return &types.QueryMetaNodeResponse{}, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GetNetworkAddr() == "" {
		return &types.QueryMetaNodeResponse{}, status.Error(codes.InvalidArgument, " network address cannot be empty")
	}

	networkAddr, err := stratos.SdsAddressFromBech32(req.GetNetworkAddr())
	if err != nil {
		return &types.QueryMetaNodeResponse{}, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	node, found := q.GetMetaNode(ctx, networkAddr)
	if !found {
		return &types.QueryMetaNodeResponse{}, status.Errorf(codes.NotFound, "network address %s not found", req.NetworkAddr)
	}

	return &types.QueryMetaNodeResponse{Node: &node}, nil
}

func (q Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.GetParams(ctx)

	return &types.QueryParamsResponse{Params: &params}, nil
}

func (q Querier) BondedResourceNodeCount(c context.Context, _ *types.QueryBondedResourceNodeCountRequest) (*types.QueryBondedResourceNodeCountResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	number := q.GetBondedResourceNodeCnt(ctx).Int64()

	return &types.QueryBondedResourceNodeCountResponse{Number: uint64(number)}, nil
}

func (q Querier) BondedMetaNodeCount(c context.Context, _ *types.QueryBondedMetaNodeCountRequest) (*types.QueryBondedMetaNodeCountResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	number := q.GetBondedMetaNodeCnt(ctx).Int64()

	return &types.QueryBondedMetaNodeCountResponse{Number: uint64(number)}, nil
}

func (q Querier) DepositByNode(c context.Context, req *types.QueryDepositByNodeRequest) (*types.QueryDepositByNodeResponse, error) {
	if req == nil {
		return &types.QueryDepositByNodeResponse{}, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.GetNetworkAddr() == "" {
		return &types.QueryDepositByNodeResponse{}, status.Error(codes.InvalidArgument, "node network address cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	queryType := req.GetQueryType()
	networkAddr, err := stratos.SdsAddressFromBech32(req.GetNetworkAddr())
	if err != nil {
		return &types.QueryDepositByNodeResponse{}, err
	}
	depositInfo := types.DepositInfo{}

	if queryType == types.QueryType_All || queryType == types.QueryType_SP {
		metaNode, found := q.GetMetaNode(ctx, networkAddr)
		if found {
			// Adding meta node deposit info
			networkAddr, _ := stratos.SdsAddressFromBech32(metaNode.GetNetworkAddress())
			unBondingDeposit, unBondedDeposit, bondedDeposit, err := q.getNodeDeposit(
				ctx,
				metaNode.GetStatus(),
				networkAddr,
				metaNode.Tokens,
			)
			if err != nil {
				return &types.QueryDepositByNodeResponse{}, err
			}
			if !metaNode.Equal(types.MetaNode{}) {
				depositInfo = types.NewDepositInfoByMetaNodeAddr(
					q.BondDenom(ctx),
					metaNode,
					unBondingDeposit,
					unBondedDeposit,
					bondedDeposit,
				)
			}
		}
	}

	if queryType == types.QueryType_All || queryType == types.QueryType_PP {
		networkAddr, err := stratos.SdsAddressFromBech32(req.GetNetworkAddr())
		if err != nil {
			return &types.QueryDepositByNodeResponse{}, err
		}
		resourceNode, found := q.GetResourceNode(ctx, networkAddr)
		if found {
			// Adding resource node deposit info
			networkAddr, _ := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddress())
			unBondingDeposit, unBondedDeposit, bondedDeposit, err := q.getNodeDeposit(
				ctx,
				resourceNode.GetStatus(),
				networkAddr,
				resourceNode.Tokens,
			)
			if err != nil {
				return nil, err
			}
			if !resourceNode.Equal(types.ResourceNode{}) {
				depositInfo = types.NewDepositInfoByResourceNodeAddr(
					q.BondDenom(ctx),
					resourceNode,
					unBondingDeposit,
					unBondedDeposit,
					bondedDeposit,
				)
			}
		}
	}
	return &types.QueryDepositByNodeResponse{DepositInfo: &depositInfo}, nil

}

func (q Querier) DepositByOwner(c context.Context, req *types.QueryDepositByOwnerRequest) (*types.QueryDepositByOwnerResponse, error) {
	if req == nil {
		return &types.QueryDepositByOwnerResponse{}, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.GetOwnerAddr() == "" {
		return &types.QueryDepositByOwnerResponse{}, status.Error(codes.InvalidArgument, "owner address cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	ownerAddr, er := sdk.AccAddressFromBech32(req.GetOwnerAddr())
	if er != nil {
		return &types.QueryDepositByOwnerResponse{}, er
	}

	store := ctx.KVStore(q.storeKey)
	var depositInfoResponses []*types.DepositInfo

	// get resource nodes
	var resourceNodes types.ResourceNodes
	resourceNodeStore := prefix.NewStore(store, types.ResourceNodeKey)

	resourceNodesPageRes, err := FilteredPaginate(q.cdc, resourceNodeStore, ownerAddr, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		val, err := types.UnmarshalResourceNode(q.cdc, value)
		if err != nil {
			return true, err
		}

		if accumulate {
			resourceNodes = append(resourceNodes, val)
		}

		return true, nil
	})

	if err != nil {
		return &types.QueryDepositByOwnerResponse{}, status.Error(codes.Internal, err.Error())
	}
	depositInfoResponses, err = GetDepositInfosByResourceNodes(ctx, q.Keeper, resourceNodes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Continue to get meta nodes
	if req.Pagination.Limit < resourceNodesPageRes.Total {
		resourceNodesPageRes.Total = uint64(len(depositInfoResponses))
		return &types.QueryDepositByOwnerResponse{DepositInfos: depositInfoResponses, Pagination: resourceNodesPageRes}, nil

	}

	metaNodesPageLimit := req.Pagination.Limit - resourceNodesPageRes.Total

	metaNodesPageOffset := uint64(0)
	if req.Pagination.Offset > resourceNodesPageRes.Total {
		metaNodesPageOffset = req.Pagination.Offset - resourceNodesPageRes.Total
	}
	metaNodesPageRequest := query.PageRequest{Offset: metaNodesPageOffset, Limit: metaNodesPageLimit, CountTotal: req.Pagination.CountTotal, Reverse: req.Pagination.CountTotal}

	var metaNodes types.MetaNodes
	metaNodeStore := prefix.NewStore(store, types.MetaNodeKey)

	_, err = FilteredPaginate(q.cdc, metaNodeStore, ownerAddr, &metaNodesPageRequest, func(key []byte, value []byte, accumulate bool) (bool, error) {
		val, err := types.UnmarshalMetaNode(q.cdc, value)
		if err != nil {
			return true, err
		}

		if accumulate {
			metaNodes = append(metaNodes, val)
		}

		return true, nil
	})

	if err != nil {
		return &types.QueryDepositByOwnerResponse{}, status.Error(codes.Internal, err.Error())
	}

	metaNodesDepositInfoResponses, err := GetDepositInfosByMetaNodes(ctx, q.Keeper, metaNodes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	depositInfoResponses = append(depositInfoResponses, metaNodesDepositInfoResponses...)
	PageRes := resourceNodesPageRes
	PageRes.Total = uint64(len(depositInfoResponses))
	return &types.QueryDepositByOwnerResponse{DepositInfos: depositInfoResponses, Pagination: PageRes}, nil
}

func (q Querier) DepositTotal(c context.Context, _ *types.QueryDepositTotalRequest) (*types.QueryDepositTotalResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	totalBondedDepositOfResourceNodes := q.GetResourceNodeBondedToken(ctx).Amount
	totalBondedDepositOfMetaNodes := q.GetMetaNodeBondedToken(ctx).Amount

	totalUnbondedDepositOfResourceNodes := q.GetResourceNodeNotBondedToken(ctx).Amount
	totalUnbondedDepositOfMetaNodes := q.GetMetaNodeNotBondedToken(ctx).Amount

	totalDepositOfResourceNodes := totalBondedDepositOfResourceNodes.Add(totalUnbondedDepositOfResourceNodes)
	totalDepositOfMetaNodes := totalBondedDepositOfMetaNodes.Add(totalUnbondedDepositOfMetaNodes)

	totalBondedDeposit := totalBondedDepositOfResourceNodes.Add(totalBondedDepositOfMetaNodes)
	totalUnbondedDeposit := totalUnbondedDepositOfResourceNodes.Add(totalUnbondedDepositOfMetaNodes)
	totalUnbondingDeposit := q.GetAllUnbondingNodesTotalBalance(ctx)
	totalUnbondedDeposit = totalUnbondedDeposit.Sub(totalUnbondingDeposit)
	res := types.NewQueryDepositTotalInfo(
		q.BondDenom(ctx),
		totalDepositOfResourceNodes,
		totalDepositOfMetaNodes,
		totalBondedDeposit,
		totalUnbondedDeposit,
		totalUnbondingDeposit,
	)

	return res, nil
}
