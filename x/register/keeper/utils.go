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
	//TotalBondedStakeOfResourceNodes    sdk.Coin
	//TotalBondedStakeOfIndexingNodes    sdk.Coin
	//TotalUnbondedStakeOfResourceNodes  sdk.Coin
	//TotalUnbondedStakeOfIndexingNodes  sdk.Coin
	//TotalUnbondingStakeOfResourceNodes sdk.Coin
	//TotalUnbondingStakeOfIndexingNodes sdk.Coin
}

// NewQueryNodesStakingInfo creates a new instance of NodesStakingInfo
func NewQueryNodesStakingInfo(
	totalStakeOfResourceNodes,
	totalStakeOfIndexingNodes,
	totalBondedStake,
	totalUnbondedStake,
	totalUnbondingStake sdk.Int,
	//totalBondedStakeOfResourceNodes,
	//totalBondedStakeOfIndexingNodes,
	//totalUnbondedStakeOfResourceNodes,
	//totalUnbondedStakeOfIndexingNodes,
	//totalUnbondingStakeOfResourceNodes,
	//totalUnbondingStakeOfIndexingNodes sdk.Int,
) NodesStakingInfo {
	return NodesStakingInfo{
		TotalStakeOfResourceNodes: sdk.NewCoin(defaultDenom, totalStakeOfResourceNodes),
		TotalStakeOfIndexingNodes: sdk.NewCoin(defaultDenom, totalStakeOfIndexingNodes),
		TotalBondedStake:          sdk.NewCoin(defaultDenom, totalBondedStake),
		TotalUnbondedStake:        sdk.NewCoin(defaultDenom, totalUnbondedStake),
		TotalUnbondingStake:       sdk.NewCoin(defaultDenom, totalUnbondingStake),
		//TotalBondedStakeOfResourceNodes:    sdk.NewCoin(defaultDenom, totalBondedStakeOfResourceNodes),
		//TotalBondedStakeOfIndexingNodes:    sdk.NewCoin(defaultDenom, totalBondedStakeOfIndexingNodes),
		//TotalUnbondedStakeOfResourceNodes:  sdk.NewCoin(defaultDenom, totalUnbondedStakeOfResourceNodes),
		//TotalUnbondedStakeOfIndexingNodes:  sdk.NewCoin(defaultDenom, totalUnbondedStakeOfIndexingNodes),
		//TotalUnbondingStakeOfResourceNodes: sdk.NewCoin(defaultDenom, totalUnbondingStakeOfResourceNodes),
		//TotalUnbondingStakeOfIndexingNodes: sdk.NewCoin(defaultDenom, totalUnbondingStakeOfIndexingNodes),
	}
}

// StakingInfoByNodeAddr Params for query 'custom/register/staking'
type StakingInfoByNodeAddr struct {
	//NodePubKey   crypto.PublicKey
	OwnerAddress   sdk.AccAddress
	NodeAddress    sdk.AccAddress
	Moniker        string
	Tokens         sdk.Coin
	BondedStake    sdk.Coin
	UnbondedStake  sdk.Coin
	UnbondingStake sdk.Coin
}

// NewStakingInfoByNodeAddr creates a new instance of StakingInfoByNodeAddr
func NewStakingInfoByNodeAddr(
	//nodePubKey crypto.PublicKey,
	ownerAddress sdk.AccAddress,
	nodeAddress sdk.AccAddress,
	moniker string,
	tokens sdk.Int,
	bondedStake sdk.Int,
	unbondedStake sdk.Int,
	unbondingStake sdk.Int,

) StakingInfoByNodeAddr {
	return StakingInfoByNodeAddr{
		//NodePubKey:   nodePubKey,
		OwnerAddress:   ownerAddress,
		NodeAddress:    nodeAddress,
		Moniker:        moniker,
		Tokens:         sdk.NewCoin(defaultDenom, tokens),
		BondedStake:    sdk.NewCoin(defaultDenom, bondedStake),
		UnbondedStake:  sdk.NewCoin(defaultDenom, unbondedStake),
		UnbondingStake: sdk.NewCoin(defaultDenom, unbondingStake),
	}
}

//// StakingInfoByNodeWalletAddr Params for query 'custom/register/staking/owner/{nodeWalletAddress}'
//type StakingInfoByNodeWalletAddr struct {
//	StakingInfoByNodeWalletAddr []StakingInfoByNodeAddr
//}

// NewStakingInfoByNodeWalletAddr creates a new instance of StakingInfoByNodeAddr
//func NewStakingInfoByNodeWalletAddr(
//	nodesStakeByNodeWalletAddr []StakingInfoByNodeAddr,
//
//) StakingInfoByNodeWalletAddr {
//	return nodesStakeByNodeWalletAddr
//}

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
		if !params.OwnerAddr.Empty() {
			if !n.OwnerAddress.Equals(params.OwnerAddr) {
				continue
			}
		}

		filteredNodes = append(filteredNodes, n)
	}

	filteredNodes = k.resPagination(filteredNodes, params)
	return filteredNodes
}

func (k Keeper) resPagination(filteredNodes []types.ResourceNode, params QueryNodesParams) []types.ResourceNode {
	start, end := client.Paginate(len(filteredNodes), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		filteredNodes = []types.ResourceNode{}
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
		if len(params.OwnerAddr) > 0 {
			if !n.OwnerAddress.Equals(params.OwnerAddr) {
				continue
			}
		}

		filteredNodes = append(filteredNodes, n)
	}
	filteredNodes = k.indPagination(filteredNodes, params)

	filteredNodes = k.indPagination(filteredNodes, params)
	return filteredNodes
}

func (k Keeper) indPagination(filteredNodes []types.IndexingNode, params QueryNodesParams) []types.IndexingNode {
	start, end := client.Paginate(len(filteredNodes), params.Page, params.Limit, QueryDefaultLimit)
	if start < 0 || end < 0 {
		filteredNodes = []types.IndexingNode{}
	} else {
		filteredNodes = filteredNodes[start:end]
	}
	return filteredNodes
}

func CheckAccAddr(w http.ResponseWriter, r *http.Request, data string) (sdk.AccAddress, bool) {
	//NodeAddrStr := mux.Vars(r)["nodeAddress"]
	//NodeAddr, err := typesTypes.GetPubKeyFromBech32(typesTypes.Bech32PubKeyTypeSdsP2PPub, NodeAddrStr)
	AccAddr, err := sdk.AccAddressFromBech32(data)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid NodeAddress.")
		return nil, false
	}
	return AccAddr, true
}
