package keeper

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
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

func (q Querier) IndexingNode(c context.Context, req *types.QueryIndexingNodeRequest) (*types.QueryIndexingNodeResponse, error) {
	if req == nil {
		return &types.QueryIndexingNodeResponse{}, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GetNetworkAddr() == "" {
		return &types.QueryIndexingNodeResponse{}, status.Error(codes.InvalidArgument, " network address cannot be empty")
	}

	networkAddr, err := stratos.SdsAddressFromBech32(req.GetNetworkAddr())
	if err != nil {
		return &types.QueryIndexingNodeResponse{}, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	node, found := q.GetIndexingNode(ctx, networkAddr)
	if !found {
		return &types.QueryIndexingNodeResponse{}, status.Errorf(codes.NotFound, "network address %s not found", req.NetworkAddr)
	}

	return &types.QueryIndexingNodeResponse{Node: &node}, nil
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
		indexingNode, found := q.GetIndexingNode(ctx, accAddr)
		if found {
			// Adding indexing node staking info
			networkAddr, _ := stratos.SdsAddressFromBech32(indexingNode.GetNetworkAddr())
			unBondingStake, unBondedStake, bondedStake, err := q.getNodeStakes(
				ctx,
				indexingNode.GetStatus(),
				networkAddr,
				indexingNode.Tokens,
			)
			if err != nil {
				return &types.QueryStakeByNodeResponse{}, err
			}
			if !indexingNode.Equal(types.IndexingNode{}) {
				stakingInfo = types.NewStakingInfoByIndexingNodeAddr(
					indexingNode,
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
			networkAddr, _ := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddr())
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

	var (
		params       types.QueryNodesParams
		stakingInfo  types.StakingInfo
		stakingInfos []*types.StakingInfo
	)

	networkAddr, er := stratos.SdsAddressFromBech32(req.GetNetworkAddr())
	if er != nil {
		return &types.QueryStakeByOwnerResponse{}, er
	}

	ownerAddr, er := sdk.AccAddressFromBech32(req.GetOwnerAddr())
	if er != nil {
		return &types.QueryStakeByOwnerResponse{}, er
	}

	page, er := strconv.Atoi(string(req.Pagination.Key))
	if er != nil {
		return &types.QueryStakeByOwnerResponse{}, er
	}

	params = types.NewQueryNodesParams(page, int(req.Pagination.Limit), networkAddr, req.GetMoniker(), ownerAddr)

	resNodes := q.GetResourceNodesFiltered(ctx, params)
	indNodes := q.GetIndexingNodesFiltered(ctx, params)

	for _, n := range indNodes {
		networkAddr, _ := stratos.SdsAddressFromBech32(n.GetNetworkAddr())
		unBondingStake, unBondedStake, bondedStake, er := q.getNodeStakes(
			ctx,
			n.GetStatus(),
			networkAddr,
			n.Tokens,
		)
		if er != nil {
			return nil, er
		}
		if !n.Equal(types.IndexingNode{}) {
			stakingInfo = types.NewStakingInfoByIndexingNodeAddr(
				n,
				unBondingStake,
				unBondedStake,
				bondedStake,
			)
			stakingInfos = append(stakingInfos, &stakingInfo)
		}
	}

	for _, n := range resNodes {
		networkAddr, _ := stratos.SdsAddressFromBech32(n.GetNetworkAddr())
		unBondingStake, unBondedStake, bondedStake, er := q.getNodeStakes(
			ctx,
			n.GetStatus(),
			networkAddr,
			n.Tokens,
		)
		if er != nil {
			return nil, er
		}
		if !n.Equal(types.ResourceNode{}) {
			stakingInfo = types.NewStakingInfoByResourceNodeAddr(
				n,
				unBondingStake,
				unBondedStake,
				bondedStake,
			)
			stakingInfos = append(stakingInfos, &stakingInfo)
		}
	}

	start, end := client.Paginate(len(stakingInfos), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		return &types.QueryStakeByOwnerResponse{}, nil
	} else {
		stakingInfos = stakingInfos[start:end]
		return &types.QueryStakeByOwnerResponse{StakingInfos: stakingInfos}, nil
	}

}

func (q Querier) StakeTotal(c context.Context, _ *types.QueryTotalStakeRequest) (*types.QueryTotalStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	totalBondedStakeOfResourceNodes := q.GetResourceNodeBondedToken(ctx).Amount
	totalBondedStakeOfIndexingNodes := q.GetIndexingNodeBondedToken(ctx).Amount

	totalUnbondedStakeOfResourceNodes := q.GetResourceNodeNotBondedToken(ctx).Amount
	totalUnbondedStakeOfIndexingNodes := q.GetIndexingNodeNotBondedToken(ctx).Amount

	resourceNodeList := q.GetAllResourceNodes(ctx)
	totalStakeOfResourceNodes := sdk.ZeroInt()
	for _, node := range resourceNodeList.GetResourceNodes() {
		totalStakeOfResourceNodes = totalStakeOfResourceNodes.Add(node.Tokens)
	}

	indexingNodeList := q.GetAllIndexingNodes(ctx)
	totalStakeOfIndexingNodes := sdk.ZeroInt()
	for _, node := range indexingNodeList.GetIndexingNodes() {
		totalStakeOfIndexingNodes = totalStakeOfIndexingNodes.Add(node.Tokens)
	}

	totalBondedStake := totalBondedStakeOfResourceNodes.Add(totalBondedStakeOfIndexingNodes)
	totalUnbondedStake := totalUnbondedStakeOfResourceNodes.Add(totalUnbondedStakeOfIndexingNodes)
	totalUnbondingStake := q.GetAllUnbondingNodesTotalBalance(ctx)
	totalUnbondedStake = totalUnbondedStake.Sub(totalUnbondingStake)
	res := types.NewQueryNodesStakingInfo(
		totalStakeOfResourceNodes,
		totalStakeOfIndexingNodes,
		totalBondedStake,
		totalUnbondedStake,
		totalUnbondingStake,
	)

	return &types.QueryTotalStakeResponse{TotalStakes: res}, nil
}