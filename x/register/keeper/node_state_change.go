package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// Called in each EndBlock
func (k Keeper) BlockRegisteredNodesUpdates(ctx sdk.Context) {
	// Remove all mature unbonding nodes from the ubd queue.
	ctx.Logger().Debug("Enter BlockRegisteredNodesUpdates")
	matureUBDs := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, networkAddr := range matureUBDs {
		balances, isMetaNode, err := k.CompleteUnbondingWithAmount(ctx, networkAddr)
		if err != nil {
			continue
		}
		if isMetaNode {
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeCompleteUnbondingMetaNode,
					sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
					sdk.NewAttribute(types.AttributeKeyNetworkAddress, networkAddr),
				),
			)
		} else {
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeCompleteUnbondingResourceNode,
					sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
					sdk.NewAttribute(types.AttributeKeyNetworkAddress, networkAddr),
				),
			)
		}

	}

	// UpdateNode won't create UBD node
	return
}

// Node state transitions
func (k Keeper) bondedToUnbonding(ctx sdk.Context, node interface{}, isMetaNode bool, coin sdk.Coin) interface{} {
	switch isMetaNode {
	case true:
		temp := node.(types.MetaNode)
		if temp.GetStatus() != stakingtypes.Bonded {
			panic(fmt.Sprintf("bad state transition bondedToUnbonding, metaNode: %v\n", temp))
		}
		return k.beginUnbondingMetaNode(ctx, &temp, &coin)
	default:
		temp := node.(types.ResourceNode)
		if temp.GetStatus() != stakingtypes.Bonded {
			panic(fmt.Sprintf("bad state transition bondedToUnbonding, resourceNode: %v\n", temp))
		}
		return k.beginUnbondingResourceNode(ctx, &temp, &coin)
	}
}

// perform all the store operations for when a Node begins unbonding
func (k Keeper) beginUnbondingResourceNode(ctx sdk.Context, resourceNode *types.ResourceNode, coin *sdk.Coin) *types.ResourceNode {
	// set node stat to unbonding, remove token from bonded pool, add token into NotBondedPool
	err := k.RemoveTokenFromPoolWhileUnbondingResourceNode(ctx, *resourceNode, *coin)
	if err != nil {
		return &types.ResourceNode{}
	}

	networkAddr, err := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddress())
	if err != nil {
		return &types.ResourceNode{}
	}
	// trigger hook if registered
	k.AfterNodeBeginUnbonding(ctx, networkAddr, false)
	return resourceNode
}
func (k Keeper) beginUnbondingMetaNode(ctx sdk.Context, metaNode *types.MetaNode, coin *sdk.Coin) *types.MetaNode {
	// change node stat, remove token from bonded pool, add token into NotBondedPool
	err := k.RemoveTokenFromPoolWhileUnbondingMetaNode(ctx, *metaNode, *coin)
	if err != nil {
		return nil
	}
	if err != nil {
		return &types.MetaNode{}
	}
	networkAddr, err := stratos.SdsAddressFromBech32(metaNode.GetNetworkAddress())
	if err != nil {
		return &types.MetaNode{}
	}
	// trigger hook if registered
	k.AfterNodeBeginUnbonding(ctx, networkAddr, true)
	return metaNode
}

func calcUnbondingMatureTime(ctx sdk.Context, currStatus stakingtypes.BondStatus, creationTime time.Time, threasholdTime time.Duration, completionTime time.Duration) time.Time {
	switch currStatus {
	case stakingtypes.Unbonded:
		return creationTime.Add(completionTime)
	default:
		now := ctx.BlockHeader().Time
		// bonded
		if creationTime.Add(threasholdTime).After(now) {
			return creationTime.Add(threasholdTime).Add(completionTime)
		}
		return now.Add(completionTime)
	}
}

// perform all the store operations for when a validator status becomes unbonded
func (k Keeper) completeUnbondingNode(ctx sdk.Context, node interface{}, isMetaNode bool) interface{} {
	if isMetaNode {
		temp := node.(types.MetaNode)
		temp.Status = stakingtypes.Unbonded
		k.SetMetaNode(ctx, temp)
		networkAddr, _ := stratos.SdsAddressFromBech32(temp.GetNetworkAddress())
		k.RemoveMetaNodeFromBitMapIdxCache(networkAddr)
		return temp
	} else {
		temp := node.(types.ResourceNode)
		temp.Status = stakingtypes.Unbonded
		k.SetResourceNode(ctx, temp)
		return temp
	}
}
