package keeper

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	abci "github.com/tendermint/tendermint/abci/types"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	pagiquery "github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

const (
	QueryResourceNodeByNetworkAddr = "resource-node"
	QueryMetaNodeByNetworkAddr     = "meta_node"
	QueryNodesDepositTotal         = "nodes_deposit_total"
	QueryNodeDepositByNodeAddr     = "node_deposit_by_addr"
	QueryNodeDepositByOwner        = "node_deposit_by_owner"
	QueryRegisterParams            = "register_params"
	QueryResourceNodesCount        = "resource_nodes_count"
	QueryMetaNodesCount            = "meta_nodes_count"
)

// NewQuerier creates a new querier for register clients.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryResourceNodeByNetworkAddr:
			return getResourceNodeByNetworkAddr(ctx, req, k, legacyQuerierCdc)
		case QueryMetaNodeByNetworkAddr:
			return getMetaNodeNetworkAddr(ctx, req, k, legacyQuerierCdc)
		case QueryNodesDepositTotal:
			return getNodesDepositTotalInfo(ctx, req, k, legacyQuerierCdc)
		case QueryNodeDepositByNodeAddr:
			return getDepositInfoByNodeAddr(ctx, req, k, legacyQuerierCdc)
		case QueryNodeDepositByOwner:
			return getDepositInfoByOwnerAddr(ctx, req, k, legacyQuerierCdc)
		case QueryRegisterParams:
			return getRegisterParams(ctx, req, k, legacyQuerierCdc)
		case QueryResourceNodesCount:
			return getResourceNodeCnt(ctx, req, k, legacyQuerierCdc)
		case QueryMetaNodesCount:
			return getMetaNodeCnt(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown register query endpoint "+req.String()+string(req.Data))
		}
	}
}

func getResourceNodeCnt(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	number := k.GetBondedResourceNodeCnt(ctx).Int64()
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, number)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func getMetaNodeCnt(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	number := k.GetBondedMetaNodeCnt(ctx).Int64()
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, number)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
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
		node   types.ResourceNode
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

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, node)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func getMetaNodeNetworkAddr(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	var (
		params types.QueryNodesParams
		node   types.MetaNode
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
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, node)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// Query total deposit of all resource/meta nodes
func getNodesDepositTotalInfo(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	totalBondedDepositOfResourceNodes := k.GetResourceNodeBondedToken(ctx).Amount
	totalBondedDepositOfMetaNodes := k.GetMetaNodeBondedToken(ctx).Amount

	totalUnbondedDepositOfResourceNodes := k.GetResourceNodeNotBondedToken(ctx).Amount
	totalUnbondedDepositOfMetaNodes := k.GetMetaNodeNotBondedToken(ctx).Amount

	totalDepositOfResourceNodes := totalBondedDepositOfResourceNodes.Add(totalUnbondedDepositOfResourceNodes)
	totalDepositOfMetaNodes := totalBondedDepositOfMetaNodes.Add(totalUnbondedDepositOfMetaNodes)

	totalBondedDeposit := totalBondedDepositOfResourceNodes.Add(totalBondedDepositOfMetaNodes)
	totalUnbondedDeposit := totalUnbondedDepositOfResourceNodes.Add(totalUnbondedDepositOfMetaNodes)
	totalUnbondingDeposit := k.GetAllUnbondingNodesTotalBalance(ctx)
	totalUnbondedDeposit = totalUnbondedDeposit.Sub(totalUnbondingDeposit)
	res := types.NewQueryDepositTotalInfo(
		k.BondDenom(ctx),
		totalDepositOfResourceNodes,
		totalDepositOfMetaNodes,
		totalBondedDeposit,
		totalUnbondedDeposit,
		totalUnbondingDeposit,
	)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func getDepositInfoByNodeAddr(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var (
		bz          []byte
		params      types.QueryNodeDepositParams
		depositInfo types.DepositInfo
	)

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	queryType := params.QueryType

	if queryType == types.QueryType_All || queryType == types.QueryType_SP {
		metaNode, found := k.GetMetaNode(ctx, params.AccAddr)
		if found {
			// Adding meta node deposit info
			networkAddr, _ := stratos.SdsAddressFromBech32(metaNode.GetNetworkAddress())
			unBondingDeposit, unBondedDeposit, bondedDeposit, err := k.getNodeDeposit(
				ctx,
				metaNode.GetStatus(),
				networkAddr,
				metaNode.Tokens,
			)
			if err != nil {
				return nil, err
			}
			if !metaNode.Equal(types.MetaNode{}) {
				depositInfo = types.NewDepositInfoByMetaNodeAddr(
					metaNode,
					unBondingDeposit,
					unBondedDeposit,
					bondedDeposit,
				)
				bzMeta, err := codec.MarshalJSONIndent(legacyQuerierCdc, depositInfo)
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
			// Adding resource node deposit info
			networkAddr, _ := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddress())
			unBondingDeposit, unBondedDeposit, bondedDeposit, err := k.getNodeDeposit(
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
					resourceNode,
					unBondingDeposit,
					unBondedDeposit,
					bondedDeposit,
				)
				bzResource, err := codec.MarshalJSONIndent(legacyQuerierCdc, depositInfo)
				if err != nil {
					return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
				}
				bz = append(bz, bzResource...)
			}
		}
	}

	return bz, nil
}

func getDepositInfoByOwnerAddr(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (
	result []byte, err error) {

	var (
		params               types.QueryNodesParams
		depositInfoResponses types.DepositInfos
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
	depositInfoResponses, err = DepositInfosResourceNodes(ctx, k, resourceNodes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Continue to get meta nodes
	if PageRequest.Limit < resourceNodesPageRes.Total {
		resourceNodesPageRes.Total = uint64(len(depositInfoResponses))
		result, err = codec.MarshalJSONIndent(legacyQuerierCdc, depositInfoResponses)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return result, nil
	}

	metaNodesPageLimit := limit - uint64(len(depositInfoResponses))

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

	metaNodesDepositInfoResponses, err := DepositInfosMetaNodes(ctx, k, metaNodes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	depositInfoResponses = append(depositInfoResponses, metaNodesDepositInfoResponses...)
	PageRes := resourceNodesPageRes
	PageRes.Total = uint64(len(depositInfoResponses))
	result, err = codec.MarshalJSONIndent(legacyQuerierCdc, depositInfoResponses)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return result, nil
}

func (k Keeper) getNodeDeposit(ctx sdk.Context, bondStatus stakingtypes.BondStatus, nodeAddress stratos.SdsAddress, tokens sdk.Int) (
	unbondingDeposit, unbondedDeposit, bondedDeposit sdk.Int, err error) {

	unbondingDeposit = sdk.NewInt(0)
	unbondedDeposit = sdk.NewInt(0)
	bondedDeposit = sdk.NewInt(0)

	switch bondStatus {
	case stakingtypes.Unbonding:
		unbondingDeposit = k.GetUnbondingNodeBalance(ctx, nodeAddress)
	case stakingtypes.Unbonded:
		unbondedDeposit = tokens
	case stakingtypes.Bonded:
		bondedDeposit = tokens
	default:
		err := fmt.Sprintf("Invalid status of node %s, expected Bonded, Unbonded, or Unbonding, got %s",
			nodeAddress.String(), bondStatus)
		return sdk.Int{}, sdk.Int{}, sdk.Int{}, sdkerrors.Wrap(sdkerrors.ErrPanic, err)
	}
	return unbondingDeposit, unbondedDeposit, bondedDeposit, nil
}

func GetIterator(prefixStore storetypes.KVStore, start []byte, reverse bool) db.Iterator {
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
		iterator := GetIterator(prefixStore, key, reverse)
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

	iterator := GetIterator(prefixStore, nil, reverse)
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

// DepositInfosResourceNodes Iteration for querying DepositInfos of resource nodes by owner(cmd and rest)
func DepositInfosResourceNodes(
	ctx sdk.Context, k Keeper, resourceNodes types.ResourceNodes,
) (types.DepositInfos, error) {
	res := types.DepositInfos{}
	resp := make([]*types.DepositInfo, len(resourceNodes))

	for i, resourceNode := range resourceNodes {
		depositInfo, err := GetDepositInfoByResourceNode(ctx, k, resourceNode)
		if err != nil {
			return nil, err
		}

		resp[i] = &depositInfo
		res = append(res, *resp[i])
	}

	return res, nil
}

// DepositInfosMetaNodes Iteration for querying DepositInfos of meta nodes by owner(cmd and rest)
func DepositInfosMetaNodes(
	ctx sdk.Context, k Keeper, metaNodes types.MetaNodes,
) (types.DepositInfos, error) {
	res := types.DepositInfos{}
	resp := make([]*types.DepositInfo, len(metaNodes))

	for i, metaNode := range metaNodes {
		depositInfo, err := GetDepositInfoByMetaNode(ctx, k, metaNode)
		if err != nil {
			return nil, err
		}

		resp[i] = &depositInfo
		res = append(res, *resp[i])
	}

	return res, nil
}

// GetDepositInfosByResourceNodes Iteration for querying DepositInfos of resource nodes by owner(grpc)
func GetDepositInfosByResourceNodes(
	ctx sdk.Context, k Keeper, resourceNodes types.ResourceNodes,
) ([]*types.DepositInfo, error) {
	resp := make([]*types.DepositInfo, len(resourceNodes))

	for i, resourceNode := range resourceNodes {
		depositInfo, err := GetDepositInfoByResourceNode(ctx, k, resourceNode)
		if err != nil {
			return nil, err
		}

		resp[i] = &depositInfo
	}

	return resp, nil
}

// GetDepositInfosByMetaNodes Iteration for querying DepositInfos of meta nodes by owner(grpc)
func GetDepositInfosByMetaNodes(
	ctx sdk.Context, k Keeper, metaNodes types.MetaNodes,
) ([]*types.DepositInfo, error) {

	resp := make([]*types.DepositInfo, len(metaNodes))

	for i, metaNode := range metaNodes {
		depositInfo, err := GetDepositInfoByMetaNode(ctx, k, metaNode)
		if err != nil {
			return nil, err
		}

		resp[i] = &depositInfo

	}

	return resp, nil
}

func GetDepositInfoByResourceNode(ctx sdk.Context, k Keeper, node types.ResourceNode) (types.DepositInfo, error) {
	networkAddr, _ := stratos.SdsAddressFromBech32(node.GetNetworkAddress())
	depositInfo := types.DepositInfo{}
	unBondingDeposit, unBondedDeposit, bondedDeposit, er := k.getNodeDeposit(
		ctx,
		node.GetStatus(),
		networkAddr,
		node.Tokens,
	)
	if er != nil {
		return depositInfo, er
	}

	if !node.Equal(types.ResourceNode{}) {
		depositInfo = types.NewDepositInfoByResourceNodeAddr(
			node,
			unBondingDeposit,
			unBondedDeposit,
			bondedDeposit,
		)
	}
	return depositInfo, nil
}

func GetDepositInfoByMetaNode(ctx sdk.Context, k Keeper, node types.MetaNode) (types.DepositInfo, error) {
	networkAddr, _ := stratos.SdsAddressFromBech32(node.GetNetworkAddress())
	depositInfo := types.DepositInfo{}
	unBondingDeposit, unBondedDeposit, bondedDeposit, er := k.getNodeDeposit(
		ctx,
		node.GetStatus(),
		networkAddr,
		node.Tokens,
	)
	if er != nil {
		return depositInfo, er
	}

	if !node.Equal(types.MetaNode{}) {
		depositInfo = types.NewDepositInfoByMetaNodeAddr(
			node,
			unBondingDeposit,
			unBondedDeposit,
			bondedDeposit,
		)
	}
	return depositInfo, nil
}
