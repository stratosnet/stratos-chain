package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	pagiquery "github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
	db "github.com/tendermint/tm-db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// this line is used by starport scaffolding # 1
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	QueryResourceNodeByNetworkAddr = "resource-nodes"
	QueryMetaNodeByNetworkAddr     = "meta_nodes"
	QueryNodesTotalStakes          = "nodes_total_stakes"
	QueryNodeStakeByNodeAddr       = "node_stakes"
	QueryNodeStakeByOwner          = "node_stakes_by_owner"
	QueryRegisterParams            = "register_params"
)

// NewQuerier creates a new querier for register clients.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryResourceNodeByNetworkAddr:
			return getResourceNodeByNetworkAddr(ctx, req, k, legacyQuerierCdc)
		case QueryMetaNodeByNetworkAddr:
			return getMetaNodesStakingInfo(ctx, req, k, legacyQuerierCdc)
		case QueryNodesTotalStakes:
			return getNodesStakingInfo(ctx, req, k, legacyQuerierCdc)
		case QueryNodeStakeByNodeAddr:
			return getStakingInfoByNodeAddr(ctx, req, k, legacyQuerierCdc)
		case QueryNodeStakeByOwner:
			return getStakingInfoByOwnerAddr(ctx, req, k, legacyQuerierCdc)
		case QueryRegisterParams:
			return getRegisterParams(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown register query endpoint "+req.String()+string(req.Data))
		}
	}
}

func getRegisterParams(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func getResourceNodeByNetworkAddr(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var (
		params types.QueryNodesParams
		nodes  []types.ResourceNode
	)

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if params.NetworkAddr.Empty() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, types.ErrInvalidNetworkAddr.Error())
	}
	node, found := k.GetResourceNode(ctx, params.NetworkAddr)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, types.ErrNoResourceNodeFound.Error())
	}
	nodes = append(nodes, node)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.NewResourceNodes(nodes...))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func getMetaNodesStakingInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	var (
		params types.QueryNodesParams
		nodes  []types.MetaNode
	)
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if params.NetworkAddr.Empty() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, types.ErrInvalidNetworkAddr.Error())
	}
	node, found := k.GetMetaNode(ctx, params.NetworkAddr)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, types.ErrNoMetaNodeFound.Error())
	}
	nodes = append(nodes, node)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.NewMetaNodes(nodes...))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// Iteration for querying total stakes of resource/meta nodes
func getNodesStakingInfo(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	totalBondedStakeOfResourceNodes := k.GetResourceNodeBondedToken(ctx).Amount
	totalBondedStakeOfMetaNodes := k.GetMetaNodeBondedToken(ctx).Amount

	totalUnbondedStakeOfResourceNodes := k.GetResourceNodeNotBondedToken(ctx).Amount
	totalUnbondedStakeOfMetaNodes := k.GetMetaNodeNotBondedToken(ctx).Amount

	//resourceNodeList := k.GetAllResourceNodes(ctx)
	//totalStakeOfResourceNodes := sdk.ZeroInt()
	//for _, node := range resourceNodeList {
	//	totalStakeOfResourceNodes = totalStakeOfResourceNodes.Add(node.Tokens)
	//}
	totalStakeOfResourceNodes := totalBondedStakeOfResourceNodes.Add(totalUnbondedStakeOfResourceNodes)

	//metaNodeList := k.GetAllMetaNodes(ctx)
	//totalStakeOfMetaNodes := sdk.ZeroInt()
	//for _, node := range metaNodeList {
	//	totalStakeOfMetaNodes = totalStakeOfMetaNodes.Add(node.Tokens)
	//}
	totalStakeOfMetaNodes := totalBondedStakeOfMetaNodes.Add(totalUnbondedStakeOfMetaNodes)

	totalBondedStake := totalBondedStakeOfResourceNodes.Add(totalBondedStakeOfMetaNodes)
	totalUnbondedStake := totalUnbondedStakeOfResourceNodes.Add(totalUnbondedStakeOfMetaNodes)
	totalUnbondingStake := k.GetAllUnbondingNodesTotalBalance(ctx)
	totalUnbondedStake = totalUnbondedStake.Sub(totalUnbondingStake)
	res := types.NewQueryNodesStakingInfo(
		totalStakeOfResourceNodes,
		totalStakeOfMetaNodes,
		totalBondedStake,
		totalUnbondedStake,
		totalUnbondingStake,
	)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

//
func getStakingInfoByNodeAddr(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var (
		bz          []byte
		params      types.QueryNodeStakingParams
		stakingInfo types.StakingInfo
	)

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	queryType := params.QueryType

	if queryType == types.QueryType_All || queryType == types.QueryType_SP {
		metaNode, found := k.GetMetaNode(ctx, params.AccAddr)
		if found {
			// Adding meta node staking info
			networkAddr, _ := stratos.SdsAddressFromBech32(metaNode.GetNetworkAddress())
			unBondingStake, unBondedStake, bondedStake, err := k.getNodeStakes(
				ctx,
				metaNode.GetStatus(),
				networkAddr,
				metaNode.Tokens,
			)
			if err != nil {
				return nil, err
			}
			if !metaNode.Equal(types.MetaNode{}) {
				stakingInfo = types.NewStakingInfoByMetaNodeAddr(
					metaNode,
					unBondingStake,
					unBondedStake,
					bondedStake,
				)
				bzMeta, err := codec.MarshalJSONIndent(legacyQuerierCdc, stakingInfo)
				if err != nil {
					return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
				}
				bz = append(bz, bzMeta...)
			}
		}
	}

	if queryType == types.QueryType_All || queryType == types.QueryType_PP {
		resourceNode, found := k.GetResourceNode(ctx, params.AccAddr)
		if found {
			// Adding resource node staking info
			networkAddr, _ := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddress())
			unBondingStake, unBondedStake, bondedStake, err := k.getNodeStakes(
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
				bzResource, err := codec.MarshalJSONIndent(legacyQuerierCdc, stakingInfo)
				if err != nil {
					return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
				}
				bz = append(bz, bzResource...)
			}
		}
	}

	return bz, nil
}

func getStakingInfoByOwnerAddr(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (result []byte, err error) {
	var (
		params types.QueryNodesParams
		//stakingInfo  types.StakingInfo
		stakingInfoResponses types.StakingInfos
	)

	if req.Data == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	err = legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if params.OwnerAddr.String() == "" {
		return nil, status.Error(codes.InvalidArgument, "owner address cannot be empty")
	}

	store := ctx.KVStore(k.storeKey)

	// get resource nodes
	var resourceNodes types.ResourceNodes
	resourceNodeStore := prefix.NewStore(store, types.ResourceNodeKey)

	limit := params.PageQuery.Limit
	if limit == 0 {
		limit = types.QueryDefaultLimit
	}
	offset := params.PageQuery.Offset
	countTotal := params.PageQuery.CountTotal
	reverse := params.PageQuery.Reverse
	PageRequest := pagiquery.PageRequest{Offset: offset, Limit: limit, CountTotal: countTotal, Reverse: reverse}

	resourceNodesPageRes, err := FilteredPaginate(k.cdc, resourceNodeStore, params.OwnerAddr, &PageRequest, func(key []byte, value []byte, accumulate bool) (bool, error) {
		val, err := types.UnmarshalResourceNode(k.cdc, value)
		if err != nil {
			return true, err
		}

		if accumulate {
			resourceNodes = append(resourceNodes, val)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	stakingInfoResponses, err = StakingInfosResourceNodes(ctx, k, resourceNodes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Continue to get meta nodes
	if PageRequest.Limit < resourceNodesPageRes.Total {
		resourceNodesPageRes.Total = uint64(len(stakingInfoResponses))
		result, err = codec.MarshalJSONIndent(legacyQuerierCdc, stakingInfoResponses)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return result, nil
	}

	metaNodesPageLimit := limit - uint64(len(stakingInfoResponses))

	metaNodesPageOffset := uint64(0)
	if offset > resourceNodesPageRes.Total {
		metaNodesPageOffset = offset - resourceNodesPageRes.Total
	}
	metaNodesPageRequest := pagiquery.PageRequest{Offset: metaNodesPageOffset, Limit: metaNodesPageLimit, CountTotal: countTotal, Reverse: reverse}
	var metaNodes types.MetaNodes
	metaNodeStore := prefix.NewStore(store, types.MetaNodeKey)

	_, err = FilteredPaginate(k.cdc, metaNodeStore, params.OwnerAddr, &metaNodesPageRequest, func(key []byte, value []byte, accumulate bool) (bool, error) {
		val, err := types.UnmarshalMetaNode(k.cdc, value)
		if err != nil {
			return true, err
		}

		if accumulate {
			metaNodes = append(metaNodes, val)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	metaNodesStakingInfoResponses, err := StakingInfosMetaNodes(ctx, k, metaNodes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	stakingInfoResponses = append(stakingInfoResponses, metaNodesStakingInfoResponses...)
	PageRes := resourceNodesPageRes
	PageRes.Total = uint64(len(stakingInfoResponses))
	result, err = codec.MarshalJSONIndent(legacyQuerierCdc, stakingInfoResponses)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return result, nil
}

func (k Keeper) getNodeStakes(ctx sdk.Context, bondStatus stakingtypes.BondStatus, nodeAddress stratos.SdsAddress, tokens sdk.Int) (unbondingStake, unbondedStake, bondedStake sdk.Int, err error) {
	unbondingStake = sdk.NewInt(0)
	unbondedStake = sdk.NewInt(0)
	bondedStake = sdk.NewInt(0)

	switch bondStatus {
	case stakingtypes.Unbonding:
		unbondingStake = k.GetUnbondingNodeBalance(ctx, nodeAddress)
	case stakingtypes.Unbonded:
		unbondedStake = tokens
	case stakingtypes.Bonded:
		bondedStake = tokens
	default:
		err := fmt.Sprintf("Invalid status of node %s, expected Bonded, Unbonded, or Unbonding, got %s",
			nodeAddress.String(), bondStatus)
		return sdk.Int{}, sdk.Int{}, sdk.Int{}, sdkerrors.Wrap(sdkerrors.ErrPanic, err)
	}
	return unbondingStake, unbondedStake, bondedStake, nil
}

func getIterator(prefixStore storetypes.KVStore, start []byte, reverse bool) db.Iterator {
	if reverse {
		var end []byte
		if start != nil {
			itr := prefixStore.Iterator(start, nil)
			defer itr.Close()
			if itr.Valid() {
				itr.Next()
				end = itr.Key()
			}
		}
		return prefixStore.ReverseIterator(nil, end)
	}
	return prefixStore.Iterator(start, nil)
}

// Iteration for querying total stakes of resource/meta nodes
func FilteredPaginate(cdc codec.Codec,
	prefixStore storetypes.KVStore,
	queryOwnerAddr sdk.AccAddress,
	pageRequest *pagiquery.PageRequest,
	onResult func(key []byte, value []byte, accumulate bool) (bool, error),
) (*pagiquery.PageResponse, error) {

	// if the PageRequest is nil, use default PageRequest
	if pageRequest == nil {
		pageRequest = &pagiquery.PageRequest{}
	}

	offset := pageRequest.Offset
	key := pageRequest.Key
	limit := pageRequest.Limit
	countTotal := pageRequest.CountTotal
	reverse := pageRequest.Reverse

	if offset > 0 && key != nil {
		return nil, fmt.Errorf("invalid request, either offset or key is expected, got both")
	}

	if limit == 0 {
		limit = types.QueryDefaultLimit

		// count total results when the limit is zero/not supplied
		countTotal = pageRequest.CountTotal
	}

	if len(key) != 0 {
		iterator := getIterator(prefixStore, key, reverse)
		defer iterator.Close()

		var numHits uint64
		var nextKey []byte
		var ownerAddr sdk.AccAddress

		for ; iterator.Valid(); iterator.Next() {
			if numHits == limit {
				nextKey = iterator.Key()
				break
			}

			if iterator.Error() != nil {
				return nil, iterator.Error()
			}

			if prefixStore.Has(types.MetaNodeKey) {
				metaNode, err := types.UnmarshalMetaNode(cdc, iterator.Value())
				if err != nil {
					continue
				}

				ownerAddr, err = sdk.AccAddressFromBech32(metaNode.GetOwnerAddress())
				if err != nil {
					continue
				}
			} else {
				resourceNode, err := types.UnmarshalResourceNode(cdc, iterator.Value())
				if err != nil {
					continue
				}

				ownerAddr, err = sdk.AccAddressFromBech32(resourceNode.GetOwnerAddress())
				if err != nil {
					continue
				}
			}

			if queryOwnerAddr.String() != ownerAddr.String() {
				continue
			}

			hit, err := onResult(iterator.Key(), iterator.Value(), true)
			if err != nil {
				return nil, err
			}

			if hit {
				numHits++
			}
		}

		return &pagiquery.PageResponse{
			NextKey: nextKey,
		}, nil
	}

	iterator := getIterator(prefixStore, nil, reverse)
	defer iterator.Close()

	end := offset + limit

	var numHits uint64
	var nextKey []byte
	var ownerAddr sdk.AccAddress

	for ; iterator.Valid(); iterator.Next() {
		if iterator.Error() != nil {
			return nil, iterator.Error()
		}

		if prefixStore.Has(types.MetaNodeKey) {
			metaNode, err := types.UnmarshalMetaNode(cdc, iterator.Value())
			if err != nil {
				continue
			}

			ownerAddr, err = sdk.AccAddressFromBech32(metaNode.GetOwnerAddress())
			if err != nil {
				continue
			}
		} else {
			resourceNode, err := types.UnmarshalResourceNode(cdc, iterator.Value())
			if err != nil {
				continue
			}

			ownerAddr, err = sdk.AccAddressFromBech32(resourceNode.GetOwnerAddress())
			if err != nil {
				continue
			}
		}

		if queryOwnerAddr.String() != ownerAddr.String() {
			continue
		}
		accumulate := numHits >= offset && numHits < end
		hit, err := onResult(iterator.Key(), iterator.Value(), accumulate)
		if err != nil {
			return nil, err
		}

		if hit {
			numHits++
		}

		if numHits == end+1 {
			nextKey = iterator.Key()

			if !countTotal {
				break
			}
		}
	}

	res := &pagiquery.PageResponse{NextKey: nextKey}
	if countTotal {
		res.Total = numHits
	}

	return res, nil
}

// StakingInfosResourceNodes Iteration for querying StakingInfos of resource nodes by owner(cmd and rest)
func StakingInfosResourceNodes(
	ctx sdk.Context, k Keeper, resourceNodes types.ResourceNodes,
) (types.StakingInfos, error) {
	res := types.StakingInfos{}
	resp := make([]*types.StakingInfo, len(resourceNodes))

	for i, resourceNode := range resourceNodes {
		stakingInfoResp, err := StakingInfoToStakingInfoResourceNode(ctx, k, resourceNode)
		if err != nil {
			return nil, err
		}

		resp[i] = &stakingInfoResp
		res = append(res, *resp[i])
	}

	return res, nil
}

// StakingInfosMetaNodes Iteration for querying StakingInfos of meta nodes by owner(cmd and rest)
func StakingInfosMetaNodes(
	ctx sdk.Context, k Keeper, metaNodes types.MetaNodes,
) (types.StakingInfos, error) {
	res := types.StakingInfos{}
	resp := make([]*types.StakingInfo, len(metaNodes))

	for i, metaNode := range metaNodes {
		stakingInfoResp, err := StakingInfoToStakingInfoMetaNode(ctx, k, metaNode)
		if err != nil {
			return nil, err
		}

		resp[i] = &stakingInfoResp
		res = append(res, *resp[i])
	}

	return res, nil
}

// StakingInfosToStakingResourceNodes Iteration for querying StakingInfos of resource nodes by owner(grpc)
func StakingInfosToStakingResourceNodes(
	ctx sdk.Context, k Keeper, resourceNodes types.ResourceNodes,
) ([]*types.StakingInfo, error) {
	resp := make([]*types.StakingInfo, len(resourceNodes))

	for i, resourceNode := range resourceNodes {
		stakingInfoResp, err := StakingInfoToStakingInfoResourceNode(ctx, k, resourceNode)
		if err != nil {
			return nil, err
		}

		resp[i] = &stakingInfoResp
	}

	return resp, nil
}

// StakingInfosToStakingMetaNodes Iteration for querying StakingInfos of meta nodes by owner(grpc)
func StakingInfosToStakingMetaNodes(
	ctx sdk.Context, k Keeper, metaNodes types.MetaNodes,
) ([]*types.StakingInfo, error) {

	resp := make([]*types.StakingInfo, len(metaNodes))

	for i, metaNode := range metaNodes {
		stakingInfoResp, err := StakingInfoToStakingInfoMetaNode(ctx, k, metaNode)
		if err != nil {
			return nil, err
		}

		resp[i] = &stakingInfoResp

	}

	return resp, nil
}

func StakingInfoToStakingInfoResourceNode(ctx sdk.Context, k Keeper, node types.ResourceNode) (types.StakingInfo, error) {
	networkAddr, _ := stratos.SdsAddressFromBech32(node.GetNetworkAddress())
	stakingInfo := types.StakingInfo{}
	unBondingStake, unBondedStake, bondedStake, er := k.getNodeStakes(
		ctx,
		node.GetStatus(),
		networkAddr,
		node.Tokens,
	)
	if er != nil {
		return stakingInfo, er
	}

	if !node.Equal(types.ResourceNode{}) {
		stakingInfo = types.NewStakingInfoByResourceNodeAddr(
			node,
			unBondingStake,
			unBondedStake,
			bondedStake,
		)
	}
	return stakingInfo, nil
}

func StakingInfoToStakingInfoMetaNode(ctx sdk.Context, k Keeper, node types.MetaNode) (types.StakingInfo, error) {
	networkAddr, _ := stratos.SdsAddressFromBech32(node.GetNetworkAddress())
	stakingInfo := types.StakingInfo{}
	unBondingStake, unBondedStake, bondedStake, er := k.getNodeStakes(
		ctx,
		node.GetStatus(),
		networkAddr,
		node.Tokens,
	)
	if er != nil {
		return stakingInfo, er
	}

	if !node.Equal(types.MetaNode{}) {
		stakingInfo = types.NewStakingInfoByMetaNodeAddr(
			node,
			unBondingStake,
			unBondedStake,
			bondedStake,
		)
	}
	return stakingInfo, nil
}
