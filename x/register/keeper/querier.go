package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"strings"

	// this line is used by starport scaffolding # 1
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	QueryResourceNodeList         = "resource_nodes"
	QueryResourceNodesByNetworkID = "resource_nodes_network_id"
	QueryResourceNodeByMoniker    = "resource_nodes_moniker"
	QueryIndexingNodeList         = "indexing_nodes"
	QueryIndexingNodeByMoniker    = "indexing_nodes_moniker"
	QueryNodesTotalStakes         = "nodes_total_stakes"
	QueryNodeStakeByNodeAddr      = "node_stakes"
	QueryNodeStakeByOwner         = "node_stakes_by_owner"
	QueryRegisterParams           = "register_params"
	QueryDefaultLimit             = 100
	defaultDenom                  = "ustos"
)

// NewQuerier creates a new querier for register clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryResourceNodeList:
			return GetResourceNodeList(ctx, req, k)
		case QueryResourceNodesByNetworkID:
			return GetResourceNodes(ctx, req, k)
		case QueryIndexingNodeList:
			//return GetIndexingNodes(ctx, req, k)
			return GetIndexingNodeList(ctx, req, k)
		case QueryNodesTotalStakes:
			return GetNodesStakingInfo(ctx, req, k)
		case QueryNodeStakeByNodeAddr:
			return GetStakingInfoByNodeAddr(ctx, req, k)
		case QueryNodeStakeByOwner:
			return GetStakingInfoByOwnerAddr(ctx, req, k)
		//case QueryNetworkSet:
		//	return GetNetworkSet(ctx, k)
		case QueryResourceNodeByMoniker:
			return GetResourceNodesByMoniker(ctx, req, k)
		case QueryIndexingNodeByMoniker:
			return GetIndexingNodesByMoniker(ctx, req, k)
		case QueryRegisterParams:
			return GetRegisterParams(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown register query endpoint "+req.String()+string(req.Data))
		}
	}
}

func GetRegisterParams(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)
	return types.ModuleCdc.MustMarshalJSON(params), nil
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
	networks := k.GetNetworks(ctx, k)
	return []byte(strings.TrimSpace(string(networks))), nil
}

func GetResourceNodeList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params QueryNodesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	resNodes := keeper.GetResourceNodesFiltered(ctx, params)
	if resNodes == nil {
		resNodes = types.ResourceNodes{}
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, resNodes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func GetIndexingNodeList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params QueryNodesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	resNodes := keeper.GetIndexingNodesFiltered(ctx, params)
	if resNodes == nil {
		resNodes = types.IndexingNodes{}
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, resNodes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func GetNodesStakingInfo(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {

	totalBondedStakeOfResourceNodes := keeper.GetResourceNodeBondedToken(ctx).Amount
	totalBondedStakeOfIndexingNodes := keeper.GetIndexingNodeBondedToken(ctx).Amount

	totalUnbondedStakeOfResourceNodes := keeper.GetResourceNodeNotBondedToken(ctx).Amount
	totalUnbondedStakeOfIndexingNodes := keeper.GetIndexingNodeNotBondedToken(ctx).Amount

	resourceNodeList := keeper.GetAllResourceNodes(ctx)
	totalStakeOfResourceNodes := sdk.ZeroInt()
	for _, node := range resourceNodeList {
		totalStakeOfResourceNodes = totalStakeOfResourceNodes.Add(node.GetTokens())
	}

	indexingNodeList := keeper.GetAllIndexingNodes(ctx)
	totalStakeOfIndexingNodes := sdk.ZeroInt()
	for _, node := range indexingNodeList {
		totalStakeOfIndexingNodes = totalStakeOfIndexingNodes.Add(node.GetTokens())
	}

	totalBondedStake := totalBondedStakeOfResourceNodes.Add(totalBondedStakeOfIndexingNodes)
	totalUnbondedStake := totalUnbondedStakeOfResourceNodes.Add(totalUnbondedStakeOfIndexingNodes)
	totalUnbondingStake := keeper.GetAllUnbondingNodesTotalBalance(ctx)

	res := NewQueryNodesStakingInfo(
		totalStakeOfResourceNodes,
		totalStakeOfIndexingNodes,
		totalBondedStake,
		totalUnbondedStake,
		totalUnbondingStake,
	)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func GetStakingInfoByNodeAddr(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var (
		bz                 []byte
		params             QueryNodeStakingParams
		resourceNodeResult StakingInfoByResourceNodeAddr
		indexingNodeResult StakingInfoByIndexingNodeAddr
	)

	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	NodeAddr, err := sdk.AccAddressFromBech32(params.AccAddr.String())
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, err.Error())
	}

	indexingNode, found := keeper.GetIndexingNode(ctx, NodeAddr)

	if !found {
		resourceNode, ok := keeper.GetResourceNode(ctx, NodeAddr)
		// Adding resource node staking info
		if ok {
			unbondingStake, unbondedStake, bondedStake, err := getNodeStakes(
				ctx, keeper,
				resourceNode.GetStatus(),
				resourceNode.GetNetworkAddr(),
				resourceNode.GetTokens(),
			)
			if err != nil {
				return nil, err
			}
			if !resourceNode.Equal(types.ResourceNode{}) {
				resourceNodeResult = NewStakingInfoByResourceNodeAddr(
					resourceNode,
					unbondingStake,
					unbondedStake,
					bondedStake,
				)
				bzResource, err := codec.MarshalJSONIndent(keeper.cdc, resourceNodeResult)
				if err != nil {
					return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
				}
				bz = append(bz, bzResource...)
			}
		}
	} else {
		// Adding indexing node staking info

		unbondingStake, unbondedStake, bondedStake, err := getNodeStakes(
			ctx, keeper,
			indexingNode.GetStatus(),
			indexingNode.GetNetworkAddr(),
			indexingNode.GetTokens(),
		)
		if err != nil {
			return nil, err
		}
		if !indexingNode.Equal(types.IndexingNode{}) {
			indexingNodeResult = NewStakingInfoByIndexingNodeAddr(
				indexingNode,
				unbondingStake,
				unbondedStake,
				bondedStake,
			)
			bzIndexing, err := codec.MarshalJSONIndent(keeper.cdc, indexingNodeResult)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}
			bz = append(bz, bzIndexing...)
		}
	}
	return bz, nil
}

func GetStakingInfoByOwnerAddr(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var (
		params              QueryNodesParams
		resourceNodeResult  StakingInfoByResourceNodeAddr
		indexingNodeResult  StakingInfoByIndexingNodeAddr
		resourceNodeResults []StakingInfoByResourceNodeAddr
		indexingNodeResults []StakingInfoByIndexingNodeAddr
	)

	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	resNodes := keeper.GetResourceNodesFiltered(ctx, params)
	indNodes := keeper.GetIndexingNodesFiltered(ctx, params)

	for _, n := range indNodes {
		unbondingStake, unbondedStake, bondedStake, err := getNodeStakes(
			ctx, keeper,
			n.GetStatus(),
			n.GetNetworkAddr(),
			n.GetTokens(),
		)
		if err != nil {
			return nil, err
		}
		if !n.Equal(types.IndexingNode{}) {
			indexingNodeResult = NewStakingInfoByIndexingNodeAddr(
				n,
				unbondingStake,
				unbondedStake,
				bondedStake,
			)
			indexingNodeResults = append(indexingNodeResults, indexingNodeResult)
		}
	}

	for _, n := range resNodes {
		unbondingStake, unbondedStake, bondedStake, err := getNodeStakes(
			ctx, keeper,
			n.GetStatus(),
			n.GetNetworkAddr(),
			n.GetTokens(),
		)
		if err != nil {
			return nil, err
		}
		if !n.Equal(types.ResourceNode{}) {
			resourceNodeResult = NewStakingInfoByResourceNodeAddr(
				n,
				unbondingStake,
				unbondedStake,
				bondedStake,
			)
			resourceNodeResults = append(resourceNodeResults, resourceNodeResult)
		}
	}

	// pagination
	indexingResultsLen := len(indexingNodeResults)
	resourceResultsLen := len(resourceNodeResults)
	start, end := client.Paginate(indexingResultsLen+resourceResultsLen, params.Page, params.Limit, QueryDefaultLimit)
	ctx.Logger().Info("[start:end]", "start", start, "end", end)
	ctx.Logger().Info("indexingResultsLen", "indexingResultsLen", indexingResultsLen)
	ctx.Logger().Info("resourceResultsLen", "resourceResultsLen", resourceResultsLen)

	var bz []byte
	switch {
	case start < 0 || end < 0:
		return nil, nil

	case end <= indexingResultsLen:
		indexingPaginated1 := indexingNodeResults[start:end]
		bzIndexing, err := codec.MarshalJSONIndent(keeper.cdc, indexingPaginated1)

		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		if len(indexingPaginated1) != 0 {
			bz = append(bz, bzIndexing...)
		}
		return bz, nil
	case start > indexingResultsLen:
		resourcePaginated := resourceNodeResults[start-indexingResultsLen : end-indexingResultsLen]
		bzResource, err := codec.MarshalJSONIndent(keeper.cdc, resourcePaginated)

		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		if len(resourcePaginated) != 0 {
			bz = append(bz, bzResource...)
		}
		return bz, nil

	default:
		indexingPaginated2 := indexingNodeResults[start:]
		bzIndexing, err := codec.MarshalJSONIndent(keeper.cdc, indexingPaginated2)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		if len(indexingPaginated2) != 0 {
			bz = append(bz, bzIndexing...)
		} else {
			bz = append(bz, []byte("[")...)
		}

		resourcePaginated2 := resourceNodeResults[:end-indexingResultsLen]
		bzResource, err := codec.MarshalJSONIndent(keeper.cdc, resourcePaginated2)

		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		if len(resourcePaginated2) != 0 {
			// remove "["
			bzResource = bzResource[1:]
		}

		if len(bz) != 1 {
			// remove "]", add ","
			bz = bz[:len(bz)-1]
			bz = append(bz, []byte(",")...)
		}

		bz = append(bz, bzResource...)
	}

	return bz, nil
}

func (k Keeper) GetIndexingNodesFiltered(ctx sdk.Context, params QueryNodesParams) []types.IndexingNode {
	nodes := k.GetAllIndexingNodes(ctx)
	filteredNodes := make([]types.IndexingNode, 0, len(nodes))

	for _, n := range nodes {
		// match NetworkID (if supplied)
		if len(params.NetworkID) > 0 {
			if strings.Compare(n.NetworkID, params.NetworkID) != 0 {
				continue
			}
		}

		// match Moniker (if supplied)
		if len(params.Moniker) > 0 {
			if strings.Compare(n.Description.Moniker, params.Moniker) != 0 {
				continue
			}
		}

		// match OwnerAddr (if supplied)
		if params.OwnerAddr.Empty() || n.OwnerAddress.Equals(params.OwnerAddr) {
			filteredNodes = append(filteredNodes, n)
		}
	}
	return filteredNodes
}

func (k Keeper) GetResourceNodesFiltered(ctx sdk.Context, params QueryNodesParams) []types.ResourceNode {
	nodes := k.GetAllResourceNodes(ctx)
	filteredNodes := make([]types.ResourceNode, 0, len(nodes))

	for _, n := range nodes {
		// match NetworkID (if supplied)
		if len(params.NetworkID) > 0 {
			if strings.Compare(n.NetworkID, params.NetworkID) != 0 {
				continue
			}
		}

		// match Moniker (if supplied)
		if len(params.Moniker) > 0 {
			if strings.Compare(n.Description.Moniker, params.Moniker) != 0 {
				continue
			}
		}

		// match OwnerAddr (if supplied)
		if params.OwnerAddr.Empty() || n.OwnerAddress.Equals(params.OwnerAddr) {
			filteredNodes = append(filteredNodes, n)
		}
	}

	return filteredNodes
}

func (k Keeper) resourceNodesPagination(filteredNodes []types.ResourceNode, params QueryNodesParams) []types.ResourceNode {
	start, end := client.Paginate(len(filteredNodes), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		filteredNodes = nil
	} else {
		filteredNodes = filteredNodes[start:end]
	}
	return filteredNodes
}

func (k Keeper) indexingNodesPagination(filteredNodes []types.IndexingNode, params QueryNodesParams) []types.IndexingNode {
	start, end := client.Paginate(len(filteredNodes), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		filteredNodes = nil
	} else {
		filteredNodes = filteredNodes[start:end]
	}
	return filteredNodes
}

func getNodeStakes(ctx sdk.Context, keeper Keeper, bondStatus sdk.BondStatus, nodeAddress sdk.AccAddress, tokens sdk.Int) (unbondingStake, unbondedStake, bondedStake sdk.Int, err error) {
	unbondingStake = sdk.NewInt(0)
	unbondedStake = sdk.NewInt(0)
	bondedStake = sdk.NewInt(0)

	switch bondStatus {
	case sdk.Unbonding:
		unbondingStake = keeper.GetUnbondingNodeBalance(ctx, nodeAddress)
	case sdk.Unbonded:
		unbondedStake = tokens
	case sdk.Bonded:
		bondedStake = tokens
	default:
		err := fmt.Sprintf("Invalid status of node %s, expected Bonded, Unbonded, or Unbonding, got %s",
			nodeAddress.String(), bondStatus.String())
		return sdk.Int{}, sdk.Int{}, sdk.Int{}, sdkerrors.Wrap(sdkerrors.ErrPanic, err)
	}
	return unbondingStake, unbondedStake, bondedStake, nil
}
