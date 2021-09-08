package keeper

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"net/http"
	"strings"
)

// QueryNodesParams Params for query 'custom/register/resource-nodes'
type QueryNodesParams struct {
	Page      int
	Limit     int
	NetworkID string
	Moniker   string
	OwnerAddr sdk.AccAddress
}

// NewQueryNodesParams creates a new instance of QueryNodesParams
func NewQueryNodesParams(page, limit int, networkID, moniker string, ownerAddr sdk.AccAddress) QueryNodesParams {
	return QueryNodesParams{
		Page:      page,
		Limit:     limit,
		NetworkID: networkID,
		Moniker:   moniker,
		OwnerAddr: ownerAddr,
	}
}

// QuerynodeStakingParams Params for query 'custom/register/staking/owner/{NodeWalletAddress}'
type QuerynodeStakingParams struct {
	AccAddr sdk.AccAddress
}

// NewQuerynodeStakingParams creates a new instance of QueryNodesParams
func NewQuerynodeStakingParams(nodeAddr sdk.AccAddress) QuerynodeStakingParams {
	return QuerynodeStakingParams{
		AccAddr: nodeAddr,
	}
}

// NodesStakingInfo Params for query 'custom/register/staking'
type NodesStakingInfo struct {
	TotalStakeOfResourceNodes sdk.Coin
	TotalStakeOfIndexingNodes sdk.Coin
	TotalBondedStake          sdk.Coin
	TotalUnbondedStake        sdk.Coin
	TotalUnbondingStake       sdk.Coin
}

// NewQueryNodesStakingInfo creates a new instance of NodesStakingInfo
func NewQueryNodesStakingInfo(
	totalStakeOfResourceNodes,
	totalStakeOfIndexingNodes,
	totalBondedStake,
	totalUnbondedStake,
	totalUnbondingStake sdk.Int,
) NodesStakingInfo {
	return NodesStakingInfo{
		TotalStakeOfResourceNodes: sdk.NewCoin(defaultDenom, totalStakeOfResourceNodes),
		TotalStakeOfIndexingNodes: sdk.NewCoin(defaultDenom, totalStakeOfIndexingNodes),
		TotalBondedStake:          sdk.NewCoin(defaultDenom, totalBondedStake),
		TotalUnbondedStake:        sdk.NewCoin(defaultDenom, totalUnbondedStake),
		TotalUnbondingStake:       sdk.NewCoin(defaultDenom, totalUnbondingStake),
	}
}

// StakingInfoByResourceNodeAddr Params for query 'custom/register/staking'
type StakingInfoByResourceNodeAddr struct {
	types.ResourceNode
	UnbondingStake       sdk.Coin
	UnbondingNodeEntries []types.UnbondingNodeEntry
}

// NewStakingInfoByResourceNodeAddr creates a new instance of StakingInfoByNodeAddr
func NewStakingInfoByResourceNodeAddr(
	resourceNode types.ResourceNode,
	unbondingStake sdk.Int,
	unbondingNodeEntries []types.UnbondingNodeEntry,

) StakingInfoByResourceNodeAddr {
	return StakingInfoByResourceNodeAddr{
		ResourceNode:         resourceNode,
		UnbondingStake:       sdk.NewCoin(defaultDenom, unbondingStake),
		UnbondingNodeEntries: unbondingNodeEntries,
	}
}

// StakingInfoByIndexingNodeAddr Params for query 'custom/register/staking'
type StakingInfoByIndexingNodeAddr struct {
	types.IndexingNode
	UnbondingStake       sdk.Coin
	UnbondingNodeEntries []types.UnbondingNodeEntry
}

// NewStakingInfoByIndexingNodeAddr creates a new instance of StakingInfoByNodeAddr
func NewStakingInfoByIndexingNodeAddr(
	indexingNode types.IndexingNode,
	unbondingStake sdk.Int,
	unbondingNodeEntries []types.UnbondingNodeEntry,
) StakingInfoByIndexingNodeAddr {
	return StakingInfoByIndexingNodeAddr{
		IndexingNode:         indexingNode,
		UnbondingStake:       sdk.NewCoin(defaultDenom, unbondingStake),
		UnbondingNodeEntries: unbondingNodeEntries,
	}
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

	filteredNodes = k.resPagination(filteredNodes, params)
	return filteredNodes
}

func (k Keeper) resPagination(filteredNodes []types.ResourceNode, params QueryNodesParams) []types.ResourceNode {
	start, end := client.Paginate(len(filteredNodes), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		//filteredNodes = []types.ResourceNode{}
		filteredNodes = nil
	} else {
		filteredNodes = filteredNodes[start:end]
	}
	return filteredNodes
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
	filteredNodes = k.indPagination(filteredNodes, params)
	return filteredNodes
}

func (k Keeper) indPagination(filteredNodes []types.IndexingNode, params QueryNodesParams) []types.IndexingNode {
	start, end := client.Paginate(len(filteredNodes), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		//filteredNodes = []types.IndexingNode{}
		filteredNodes = nil
	} else {
		filteredNodes = filteredNodes[start:end]
	}
	return filteredNodes
}

func CheckAccAddr(w http.ResponseWriter, r *http.Request, data string) (sdk.AccAddress, bool) {
	AccAddr, err := sdk.AccAddressFromBech32(data)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid NodeAddress.")
		return nil, false
	}
	return AccAddr, true
}
