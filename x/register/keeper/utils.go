package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"net/http"
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

type QueryNodeStakingParams struct {
	AccAddr sdk.AccAddress
}

// NewQuerynodeStakingParams creates a new instance of QueryNodesParams
func NewQuerynodeStakingParams(nodeAddr sdk.AccAddress) QueryNodeStakingParams {
	return QueryNodeStakingParams{
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

type StakingInfoByResourceNodeAddr struct {
	types.ResourceNode
	BondedStake    sdk.Coin
	UnbondingStake sdk.Coin
	UnbondedStake  sdk.Coin
}

// NewStakingInfoByResourceNodeAddr creates a new instance of StakingInfoByNodeAddr
func NewStakingInfoByResourceNodeAddr(
	resourceNode types.ResourceNode,
	unbondingStake sdk.Int,
	unbondedStake sdk.Int,
	bondedStake sdk.Int,

) StakingInfoByResourceNodeAddr {
	return StakingInfoByResourceNodeAddr{
		ResourceNode:   resourceNode,
		UnbondingStake: sdk.NewCoin(defaultDenom, unbondingStake),
		UnbondedStake:  sdk.NewCoin(defaultDenom, unbondedStake),
		BondedStake:    sdk.NewCoin(defaultDenom, bondedStake),
	}
}

type StakingInfoByIndexingNodeAddr struct {
	types.IndexingNode
	BondedStake    sdk.Coin
	UnbondingStake sdk.Coin
	UnbondedStake  sdk.Coin
}

// NewStakingInfoByIndexingNodeAddr creates a new instance of StakingInfoByNodeAddr
func NewStakingInfoByIndexingNodeAddr(
	indexingNode types.IndexingNode,
	unbondingStake sdk.Int,
	unbondedStake sdk.Int,
	bondedStake sdk.Int,
) StakingInfoByIndexingNodeAddr {
	return StakingInfoByIndexingNodeAddr{
		IndexingNode:   indexingNode,
		UnbondingStake: sdk.NewCoin(defaultDenom, unbondingStake),
		UnbondedStake:  sdk.NewCoin(defaultDenom, unbondedStake),
		BondedStake:    sdk.NewCoin(defaultDenom, bondedStake),
	}
}

func CheckAccAddr(w http.ResponseWriter, r *http.Request, data string) (sdk.AccAddress, bool) {
	AccAddr, err := sdk.AccAddressFromBech32(data)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid NodeAddress.")
		return nil, false
	}
	return AccAddr, true
}
