package keeper

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"

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

	QueryDefaultLimit = 100
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

func getNodesStakingInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	totalBondedStakeOfResourceNodes := k.GetResourceNodeBondedToken(ctx).Amount
	totalBondedStakeOfMetaNodes := k.GetMetaNodeBondedToken(ctx).Amount

	totalUnbondedStakeOfResourceNodes := k.GetResourceNodeNotBondedToken(ctx).Amount
	totalUnbondedStakeOfMetaNodes := k.GetMetaNodeNotBondedToken(ctx).Amount

	resourceNodeList := k.GetAllResourceNodes(ctx)
	totalStakeOfResourceNodes := sdk.ZeroInt()
	for _, node := range resourceNodeList {
		totalStakeOfResourceNodes = totalStakeOfResourceNodes.Add(node.Tokens)
	}

	metaNodeList := k.GetAllMetaNodes(ctx)
	totalStakeOfMetaNodes := sdk.ZeroInt()
	for _, node := range metaNodeList {
		totalStakeOfMetaNodes = totalStakeOfMetaNodes.Add(node.Tokens)
	}

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
		params       types.QueryNodesParams
		stakingInfo  types.StakingInfo
		stakingInfos []types.StakingInfo
	)

	err = legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	resNodes := k.GetResourceNodesFiltered(ctx, params)
	metaNodes := k.GetMetaNodesFiltered(ctx, params)

	for _, n := range metaNodes {
		networkAddr, _ := stratos.SdsAddressFromBech32(n.GetNetworkAddress())
		unBondingStake, unBondedStake, bondedStake, err := k.getNodeStakes(
			ctx,
			n.GetStatus(),
			networkAddr,
			n.Tokens,
		)
		if err != nil {
			return nil, err
		}
		if !n.Equal(types.MetaNode{}) {
			stakingInfo = types.NewStakingInfoByMetaNodeAddr(
				n,
				unBondingStake,
				unBondedStake,
				bondedStake,
			)
			stakingInfos = append(stakingInfos, stakingInfo)
		}
	}

	for _, n := range resNodes {
		networkAddr, _ := stratos.SdsAddressFromBech32(n.GetNetworkAddress())
		unBondingStake, unBondedStake, bondedStake, err := k.getNodeStakes(
			ctx,
			n.GetStatus(),
			networkAddr,
			n.Tokens,
		)
		if err != nil {
			return nil, err
		}
		if !n.Equal(types.ResourceNode{}) {
			stakingInfo = types.NewStakingInfoByResourceNodeAddr(
				n,
				unBondingStake,
				unBondedStake,
				bondedStake,
			)
			stakingInfos = append(stakingInfos, stakingInfo)
		}
	}

	start, end := client.Paginate(len(stakingInfos), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		return nil, nil
	} else {
		stakingInfos = stakingInfos[start:end]
		result, err = codec.MarshalJSONIndent(legacyQuerierCdc, stakingInfos)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return result, nil
	}
}

func (k Keeper) resourceNodesPagination(filteredNodes []types.ResourceNode, params types.QueryNodesParams) []types.ResourceNode {
	start, end := client.Paginate(len(filteredNodes), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		filteredNodes = nil
	} else {
		filteredNodes = filteredNodes[start:end]
	}
	return filteredNodes
}

func (k Keeper) metaNodesPagination(filteredNodes []types.MetaNode, params types.QueryNodesParams) []types.MetaNode {
	start, end := client.Paginate(len(filteredNodes), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		filteredNodes = nil
	} else {
		filteredNodes = filteredNodes[start:end]
	}
	return filteredNodes
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

func (k Keeper) GetMetaNodesFiltered(ctx sdk.Context, params types.QueryNodesParams) []types.MetaNode {
	nodes := k.GetAllMetaNodes(ctx)
	filteredNodes := make([]types.MetaNode, 0, len(nodes))

	for _, n := range nodes {
		// match NetworkAddr (if supplied)
		nodeNetworkAddr, er := stratos.SdsAddressFromBech32(n.GetNetworkAddress())
		if er != nil {
			continue
		}
		if !params.NetworkAddr.Empty() {
			if nodeNetworkAddr.Equals(params.NetworkAddr) {
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
		nodeOwnerAddr, er := sdk.AccAddressFromBech32(n.GetNetworkAddress())
		if er != nil {
			continue
		}
		if params.OwnerAddr.Empty() || nodeOwnerAddr.Equals(params.OwnerAddr) {
			filteredNodes = append(filteredNodes, n)
		}
	}
	return filteredNodes
}

func (k Keeper) GetResourceNodesFiltered(ctx sdk.Context, params types.QueryNodesParams) []types.ResourceNode {
	nodes := k.GetAllResourceNodes(ctx)
	filteredNodes := make([]types.ResourceNode, 0, len(nodes))

	for _, n := range nodes {
		// match NetworkAddr (if supplied)
		nodeNetworkAddr, er := stratos.SdsAddressFromBech32(n.GetNetworkAddress())
		if er != nil {
			continue
		}
		if !params.NetworkAddr.Empty() {
			if nodeNetworkAddr.Equals(params.NetworkAddr) {
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
		nodeOwnerAddr, er := sdk.AccAddressFromBech32(n.GetNetworkAddress())
		if er != nil {
			continue
		}
		if params.OwnerAddr.Empty() || nodeOwnerAddr.Equals(params.OwnerAddr) {
			filteredNodes = append(filteredNodes, n)
		}
	}
	return filteredNodes
}
