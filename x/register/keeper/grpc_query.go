package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
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

func (q Querier) StakeByNode(c context.Context, req *types.QueryStakeByNodeRequest) (*types.QueryStakeByNodeResponse, error) {
	if req == nil {
		return &types.QueryStakeByNodeResponse{}, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.GetAccAddr() == "" {
		return &types.QueryStakeByNodeResponse{}, status.Error(codes.InvalidArgument, "node network address cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	queryType := req.QueryType
	accAddr, err := stratos.SdsAddressFromBech32(req.AccAddr)
	if err != nil {
		return &types.QueryStakeByNodeResponse{}, err
	}
	stakingInfo := types.StakingInfo{}

	if queryType == types.QueryType_All || queryType == types.QueryType_SP {
		metaNode, found := q.GetMetaNode(ctx, accAddr)
		if found {
			// Adding meta node staking info
			networkAddr, _ := stratos.SdsAddressFromBech32(metaNode.GetNetworkAddress())
			unBondingStake, unBondedStake, bondedStake, err := q.getNodeStakes(
				ctx,
				metaNode.GetStatus(),
				networkAddr,
				metaNode.Tokens,
			)
			if err != nil {
				return &types.QueryStakeByNodeResponse{}, err
			}
			if !metaNode.Equal(types.MetaNode{}) {
				stakingInfo = types.NewStakingInfoByMetaNodeAddr(
					metaNode,
					unBondingStake,
					unBondedStake,
					bondedStake,
				)
			}
		}
	}

	if queryType == types.QueryType_All || queryType == types.QueryType_PP {
		accAddr, err := stratos.SdsAddressFromBech32(req.GetAccAddr())
		if err != nil {
			return &types.QueryStakeByNodeResponse{}, err
		}
		resourceNode, found := q.GetResourceNode(ctx, accAddr)
		if found {
			// Adding resource node staking info
			networkAddr, _ := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddress())
			unBondingStake, unBondedStake, bondedStake, err := q.getNodeStakes(
				ctx,
				resourceNode.GetStatus(),
				networkAddr,
				resourceNode.Tokens,
			)
			if err != nil {
				return nil, err
			}
			if !resourceNode.Equal(types.ResourceNode{}) {
				stakingInfo = types.NewStakingInfoByResourceNodeAddr(
					resourceNode,
					unBondingStake,
					unBondedStake,
					bondedStake,
				)
			}
		}
	}
	return &types.QueryStakeByNodeResponse{StakingInfo: &stakingInfo}, nil

}

func (q Querier) StakeByOwner(c context.Context, req *types.QueryStakeByOwnerRequest) (*types.QueryStakeByOwnerResponse, error) {
	if req == nil {
		return &types.QueryStakeByOwnerResponse{}, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.GetOwnerAddr() == "" {
		return &types.QueryStakeByOwnerResponse{}, status.Error(codes.InvalidArgument, "owner address cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	var metaNodes types.MetaNodes

	ownerAddr, er := sdk.AccAddressFromBech32(req.GetOwnerAddr())
	if er != nil {
		return &types.QueryStakeByOwnerResponse{}, er
	}

	store := ctx.KVStore(q.storeKey)
	resourceNodeStore := prefix.NewStore(store, types.MetaNodeKey)

	pageRes, err := FilteredPaginate(q.cdc, resourceNodeStore, ownerAddr, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		val, err := types.UnmarshalMetaNode(q.cdc, value)
		if err != nil {
			return false, err
		}

		if accumulate {
			metaNodes = append(metaNodes, val)
		}

		return true, nil
	})

	if err != nil {
		return &types.QueryStakeByOwnerResponse{}, status.Error(codes.Internal, err.Error())
	}
	fmt.Println("MetaNodes: ", metaNodes)

	stakingInfoResponses, err := StakingInfosToStakingInfoResponses(ctx, q.Keeper, metaNodes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryStakeByOwnerResponse{StakingInfos: stakingInfoResponses, Pagination: pageRes}, nil

	//page := req.GetPage()
	//if page == 0 {
	//	page = QueryDefaultPage
	//}
	//
	//limit := req.GetLimit()
	//if limit == 0 {
	//	limit = QueryDefaultLimit
	//}
	//
	//params = types.NewQueryNodesParams(int(page), int(limit), nil, "", ownerAddr)
	//
	//resNodes := q.GetResourceNodesFiltered(ctx, params)
	//metaNodes := q.GetMetaNodesFiltered(ctx, params)
	//
	//for i, _ := range metaNodes {
	//	networkAddr, _ := stratos.SdsAddressFromBech32(metaNodes[i].GetNetworkAddress())
	//	unBondingStake, unBondedStake, bondedStake, er := q.getNodeStakes(
	//		ctx,
	//		metaNodes[i].GetStatus(),
	//		networkAddr,
	//		metaNodes[i].Tokens,
	//	)
	//	if er != nil {
	//		return nil, er
	//	}
	//	if !metaNodes[i].Equal(types.MetaNode{}) {
	//		stakingInfo := types.NewStakingInfoByMetaNodeAddr(
	//			metaNodes[i],
	//			unBondingStake,
	//			unBondedStake,
	//			bondedStake,
	//		)
	//		stakingInfos = append(stakingInfos, stakingInfo)
	//	}
	//}
	//
	//for i, _ := range resNodes {
	//	networkAddr, _ := stratos.SdsAddressFromBech32(resNodes[i].GetNetworkAddress())
	//	unBondingStake, unBondedStake, bondedStake, er := q.getNodeStakes(
	//		ctx,
	//		resNodes[i].GetStatus(),
	//		networkAddr,
	//		resNodes[i].Tokens,
	//	)
	//	if er != nil {
	//		return nil, er
	//	}
	//	if !resNodes[i].Equal(types.ResourceNode{}) {
	//		stakingInfo := types.NewStakingInfoByResourceNodeAddr(
	//			resNodes[i],
	//			unBondingStake,
	//			unBondedStake,
	//			bondedStake,
	//		)
	//		stakingInfos = append(stakingInfos, stakingInfo)
	//	}
	//}
	//
	//start, end := client.Paginate(len(stakingInfos), params.Page, params.Limit, QueryDefaultLimit)
	//if start < 0 || end < 0 {
	//	return &types.QueryStakeByOwnerResponse{}, nil
	//} else {
	//	stakingInfos = stakingInfos[start:end]
	//	return &types.QueryStakeByOwnerResponse{StakingInfos: stakingInfos}, nil
	//}

}

func (q Querier) StakeTotal(c context.Context, _ *types.QueryTotalStakeRequest) (*types.QueryTotalStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	totalBondedStakeOfResourceNodes := q.GetResourceNodeBondedToken(ctx).Amount
	totalBondedStakeOfMetaNodes := q.GetMetaNodeBondedToken(ctx).Amount

	totalUnbondedStakeOfResourceNodes := q.GetResourceNodeNotBondedToken(ctx).Amount
	totalUnbondedStakeOfMetaNodes := q.GetMetaNodeNotBondedToken(ctx).Amount

	resourceNodeList := q.GetAllResourceNodes(ctx)
	totalStakeOfResourceNodes := sdk.ZeroInt()
	for _, node := range resourceNodeList {
		totalStakeOfResourceNodes = totalStakeOfResourceNodes.Add(node.Tokens)
	}

	metaNodeList := q.GetAllMetaNodes(ctx)
	totalStakeOfMetaNodes := sdk.ZeroInt()
	for _, node := range metaNodeList {
		totalStakeOfMetaNodes = totalStakeOfMetaNodes.Add(node.Tokens)
	}

	totalBondedStake := totalBondedStakeOfResourceNodes.Add(totalBondedStakeOfMetaNodes)
	totalUnbondedStake := totalUnbondedStakeOfResourceNodes.Add(totalUnbondedStakeOfMetaNodes)
	totalUnbondingStake := q.GetAllUnbondingNodesTotalBalance(ctx)
	totalUnbondedStake = totalUnbondedStake.Sub(totalUnbondingStake)
	res := types.NewQueryNodesStakingInfo(
		totalStakeOfResourceNodes,
		totalStakeOfMetaNodes,
		totalBondedStake,
		totalUnbondedStake,
		totalUnbondingStake,
	)

	return &types.QueryTotalStakeResponse{TotalStakes: res}, nil
}
