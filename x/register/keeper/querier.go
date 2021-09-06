package keeper

import (
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
	QueryResourceNodeList          = "resource_nodes"
	QueryResourceNodeByMoniker     = "resource_nodes_moniker"
	QueryIndexingNodeList          = "indexing_nodes"
	QueryIndexingNodeByMoniker     = "indexing_nodes_moniker"
	QueryNodesTotalStakes          = "nodes_total_stakes"
	QueryNodeStakeByNodeAddr       = "node_stakes"
	QueryNodeStakeByNodeWalletAddr = "node_stakes_by_owner"
	QueryDefaultLimit              = 20
	defaultDenom                   = "ustos"
)

// NewQuerier creates a new querier for register clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryResourceNodeList:
			//return GetResourceNodes(ctx, req, k)
			return GetResourceNodeList(ctx, req, k)
		case QueryIndexingNodeList:
			//return GetIndexingNodes(ctx, req, k)
			return GetIndexingNodeList(ctx, req, k)
		case QueryNodesTotalStakes:
			return GetNodesStakingInfo(ctx, req, k)
		case QueryNodeStakeByNodeAddr:
			return GetStakingInfoByNodeAddr(ctx, req, k)
		case QueryNodeStakeByNodeWalletAddr:
			return GetStakingInfoByNodeWalletAddr(ctx, req, k)
		//case QueryNetworkSet:
		//	return GetNetworkSet(ctx, k)
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

	totalUnbondingStakeOfResourceNodes := totalStakeOfResourceNodes.Sub(totalBondedStakeOfResourceNodes).Sub(totalUnbondedStakeOfResourceNodes)
	totalUnbondingStakeOfIndexingNodes := totalStakeOfIndexingNodes.Sub(totalBondedStakeOfIndexingNodes).Sub(totalUnbondedStakeOfIndexingNodes)

	totalBondedStake := totalBondedStakeOfResourceNodes.Add(totalBondedStakeOfIndexingNodes)
	totalUnbondedStake := totalUnbondedStakeOfResourceNodes.Add(totalUnbondedStakeOfIndexingNodes)
	totalUnbondingStake := totalUnbondingStakeOfResourceNodes.Add(totalUnbondingStakeOfIndexingNodes)

	ctx.Logger().Info("Info:", "totalStakeOfResourceNodes", totalStakeOfResourceNodes,
		"totalStakeOfIndexingNodes", totalStakeOfIndexingNodes,
		"totalBondedStakeOfResourceNodes", totalBondedStakeOfResourceNodes,
		"totalBondedStakeOfIndexingNodes", totalBondedStakeOfIndexingNodes,
		"totalUnbondedStakeOfResourceNodes", totalUnbondedStakeOfResourceNodes,
		"totalUnbondedStakeOfIndexingNodes", totalUnbondedStakeOfIndexingNodes,
		"totalUnbondingStakeOfResourceNodes", totalUnbondingStakeOfResourceNodes,
		"totalUnbondingStakeOfIndexingNodes", totalUnbondingStakeOfIndexingNodes,
		"totalBondedStake", totalBondedStake,
		"totalUnbondedStake", totalUnbondedStake,
		"totalUnbondingStake", totalUnbondingStake,
	)

	res := NewQueryNodesStakingInfo(
		totalStakeOfResourceNodes,
		totalStakeOfIndexingNodes,
		totalBondedStake,
		//totalBondedStakeOfResourceNodes,
		//totalBondedStakeOfIndexingNodes,
		totalUnbondedStake,
		//totalUnbondedStakeOfResourceNodes,
		//totalUnbondedStakeOfIndexingNodes,
		totalUnbondingStake,
		//totalUnbondingStakeOfResourceNodes,
		//totalUnbondingStakeOfIndexingNodes,
	)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func GetStakingInfoByNodeAddr(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	ctx.Logger().Info("NodeAddr", "NodeAddr", string(req.Data))
	var params QuerynodeStakingParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	ctx.Logger().Info("params", "params.NodeAddr", params.AccAddr)

	//NodeAddr, err := sdk.AccAddressFromBech32("st1v0r46n9vr62q3xac80xmtsf5sct3qazp7azfya")
	NodeAddr, err := sdk.AccAddressFromBech32(params.AccAddr.String())
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, err.Error())
	}
	ctx.Logger().Info("NodeAddr after converting", "NodeAddr", NodeAddr)

	res := StakingInfoByNodeAddr{}
	indexingNode, found := keeper.GetIndexingNode(ctx, NodeAddr)
	//ctx.Logger().Info("indexingNode from GetIndexingNode", "indexingNode", indexingNode)

	if !found {
		ctx.Logger().Info("indexingNode not found")
		resourceNode, ok := keeper.GetResourceNode(ctx, NodeAddr)
		ctx.Logger().Info("GetResourceNode", "ok", ok)
		if ok {
			ctx.Logger().Info("resourceNode found", "resourceNode", resourceNode)
			res = NewStakingInfoByNodeAddr(
				//resourceNode.PubKey,
				resourceNode.OwnerAddress,
				resourceNode.GetNetworkAddr(),
				resourceNode.Description.Moniker,

				resourceNode.GetTokens(),
				keeper.GetLastResourceNodeStake(ctx, NodeAddr),
				sdk.NewInt(0),
				sdk.NewInt(0),
			)
		}
	} else {
		ctx.Logger().Info("indexingNode found", "indexingNode", indexingNode)
		res = NewStakingInfoByNodeAddr(
			//indexingNode.PubKey,
			indexingNode.OwnerAddress,
			indexingNode.GetNetworkAddr(),
			indexingNode.Description.Moniker,
			indexingNode.GetTokens(),
			keeper.GetLastIndexingNodeStake(ctx, NodeAddr),
			sdk.NewInt(0),
			sdk.NewInt(0),
		)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func GetStakingInfoByNodeWalletAddr(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params QueryNodesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	ctx.Logger().Info("params", "params", params)

	NodeWalletAddr, err := sdk.AccAddressFromBech32(params.OwnerAddr.String())
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, err.Error())
	}
	ctx.Logger().Info("OwnerAddr after converting", "OwnerAddr", NodeWalletAddr)

	params2 := params
	params2.Page = 1
	params2.Limit = QueryDefaultLimit
	if params.Limit > 0 {
		params2.Limit = params.Limit
	}

	ctx.Logger().Info("params2", "params2", params2)
	resNodes := keeper.GetResourceNodesFiltered(ctx, params2)
	ctx.Logger().Info("resNodes", "resNodes", resNodes)
	indNodes := keeper.GetIndexingNodesFiltered(ctx, params2)
	ctx.Logger().Info("indNodes", "indNodes", indNodes)
	var res []StakingInfoByNodeAddr
	if len(indNodes) > 0 {
		for _, n := range indNodes {
			totalStake := keeper.GetLastIndexingNodeStake(ctx, sdk.AccAddress(n.PubKey.Address()))
			bondedStake := sdk.NewInt(0)
			unBondedStake := sdk.NewInt(0)
			//unBondingStake := sdk.NewInt(0)

			switch n.GetStatus() {
			case sdk.Bonded:
				bondedStake = totalStake

			case sdk.Unbonded:
				unBondedStake = totalStake
			}

			StakingOfIndexingNodes := NewStakingInfoByNodeAddr(
				n.OwnerAddress,
				sdk.AccAddress(n.PubKey.Address()),
				n.Description.Moniker,
				//n.GetTokens(),
				totalStake,
				bondedStake,
				unBondedStake,
				totalStake.Sub(bondedStake).Sub(unBondedStake),
			)
			res = append(res, StakingOfIndexingNodes)
		}
	}

	if len(resNodes) > 0 {
		for _, n := range resNodes {
			totalStake := keeper.GetLastIndexingNodeStake(ctx, sdk.AccAddress(n.PubKey.Address()))
			bondedStake := sdk.NewInt(0)
			unBondedStake := sdk.NewInt(0)
			//unBondingStake := sdk.NewInt(0)

			switch n.GetStatus() {
			case sdk.Bonded:
				bondedStake = totalStake

			case sdk.Unbonded:
				unBondedStake = totalStake
			}

			StakingOfResourceNodes := NewStakingInfoByNodeAddr(
				n.OwnerAddress,
				sdk.AccAddress(n.PubKey.Address()),
				n.Description.Moniker,
				totalStake,
				bondedStake,
				unBondedStake,
				totalStake.Sub(bondedStake).Sub(unBondedStake),
			)
			res = append(res, StakingOfResourceNodes)
		}
	}

	ctx.Logger().Info("res", "res", res)
	start, end := client.Paginate(len(res), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		res = []StakingInfoByNodeAddr{}
	} else {
		res = res[start:end]
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
