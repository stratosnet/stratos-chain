package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// BlockRegisteredNodesUpdates Called in each EndBlock
func (k Keeper) BlockRegisteredNodesUpdates(ctx sdk.Context) {
	// Remove all mature unbonding nodes from the ubd queue.
	ctx.Logger().Debug("Enter BlockRegisteredNodesUpdates")
	matureUBDs := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, networkAddr := range matureUBDs {
		balances, isMetaNode, err := k.completeUnbondingWithAmount(ctx, networkAddr)
		if err != nil {
			continue
		}
		if isMetaNode {
			_ = ctx.EventManager().EmitTypedEvent(&types.EventCompleteUnBondingMetaNode{
				Amount:         balances.String(),
				NetworkAddress: networkAddr,
			})
		} else {
			_ = ctx.EventManager().EmitTypedEvent(&types.EventCompleteUnBondingResourceNode{
				Amount:         balances.String(),
				NetworkAddress: networkAddr,
			})
		}

	}

	// UpdateNode won't create UBD node
	return
}

// CompleteUnbondingWithAmount completes the unbonding of all mature entries in
// the retrieved unbonding delegation object and returns the total unbonding
// balance or an error upon failure.
func (k Keeper) completeUnbondingWithAmount(ctx sdk.Context, networkAddrBech32 string) (sdk.Coins, bool, error) {
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrBech32)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("NetworAddr: %s is invalid", networkAddrBech32))
		return nil, false, types.ErrInvalidNetworkAddr
	}

	ubd, found := k.GetUnbondingNode(ctx, networkAddr)
	if !found {
		ctx.Logger().Info(fmt.Sprintf("NetworAddr: %s not found while completing UnbondingWithAmount", networkAddr))
		return nil, false, types.ErrNoUnbondingNode
	}

	bondDenom := k.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time
	ctx.Logger().Debug(fmt.Sprintf("Completing UnbondingWithAmount, networAddr: %s", networkAddr))
	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(ctxTime) {
			ubd.RemoveEntry(int64(i))
			i--

			// track undelegation only when remaining or truncated shares are non-zero
			if !entry.Balance.IsZero() {
				amt := sdk.NewCoin(bondDenom, *entry.Balance)
				err = k.subtractUBDNodeDeposit(ctx, ubd, amt)
				if err != nil {
					return nil, false, err
				}

				balances = balances.Add(amt)
			}
		}
	}

	// set the unbonding node or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveUnbondingNode(ctx, networkAddr)
		err = k.unbondingToUnbonded(ctx, networkAddr, ubd.IsMetaNode)
		if err != nil {
			return balances, ubd.IsMetaNode, err
		}
	} else {
		k.SetUnbondingNode(ctx, ubd)
	}

	return balances, ubd.IsMetaNode, nil
}

// Node state transitions
func (k Keeper) bondedToUnbonding(ctx sdk.Context, node interface{}, isMetaNode bool, coin sdk.Coin) interface{} {
	switch isMetaNode {
	case true:
		temp := node.(types.MetaNode)
		if temp.GetStatus() != stakingtypes.Bonded {
			panic(fmt.Sprintf("bad state transition bondedToUnbonding, metaNode: %v\n", temp))
		}
		return k.beginUnbondingMetaNode(ctx, temp, coin)
	default:
		temp := node.(types.ResourceNode)
		if temp.GetStatus() != stakingtypes.Bonded {
			panic(fmt.Sprintf("bad state transition bondedToUnbonding, resourceNode: %v\n", temp))
		}
		return k.beginUnbondingResourceNode(ctx, temp, coin)
	}
}

// perform all the store operations for when a Node begins unbonding
func (k Keeper) beginUnbondingResourceNode(ctx sdk.Context, resourceNode types.ResourceNode, coin sdk.Coin) types.ResourceNode {
	// set node stat to unbonding, remove token from bonded pool, add token into NotBondedPool
	err := k.RemoveTokenFromPoolWhileUnbondingResourceNode(ctx, resourceNode, coin)
	if err != nil {
		return types.ResourceNode{}
	}

	networkAddr, err := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddress())
	if err != nil {
		return types.ResourceNode{}
	}
	// trigger hook if registered
	k.AfterNodeBeginUnbonding(ctx, networkAddr, false)
	return resourceNode
}
func (k Keeper) beginUnbondingMetaNode(ctx sdk.Context, metaNode types.MetaNode, coin sdk.Coin) types.MetaNode {
	// change node stat, remove token from bonded pool, add token into NotBondedPool
	err := k.RemoveTokenFromPoolWhileUnbondingMetaNode(ctx, coin)
	if err != nil {
		return types.MetaNode{}
	}
	if err != nil {
		return types.MetaNode{}
	}
	networkAddr, err := stratos.SdsAddressFromBech32(metaNode.GetNetworkAddress())
	if err != nil {
		return types.MetaNode{}
	}
	// trigger hook if registered
	k.AfterNodeBeginUnbonding(ctx, networkAddr, true)
	return metaNode
}

func (k Keeper) calcUnbondingMatureTime(ctx sdk.Context, currStatus stakingtypes.BondStatus, creationTime time.Time) time.Time {
	thresholdTime := k.UnbondingThreasholdTime(ctx)
	completionTime := k.UnbondingCompletionTime(ctx)

	switch currStatus {
	case stakingtypes.Unbonded:
		return creationTime.Add(completionTime)
	default:
		now := ctx.BlockHeader().Time
		// bonded
		if creationTime.Add(thresholdTime).After(now) {
			return creationTime.Add(thresholdTime).Add(completionTime)
		}
		return now.Add(completionTime)
	}
}

// perform all the store operations for when a validator status becomes unbonded
func (k Keeper) unbondingToUnbonded(ctx sdk.Context, networkAddr stratos.SdsAddress, isMetaNode bool) error {
	if isMetaNode {
		metaNode, found := k.GetMetaNode(ctx, networkAddr)
		if !found {
			return types.ErrNoMetaNodeFound
		}
		metaNode.Status = stakingtypes.Unbonded
		k.SetMetaNode(ctx, metaNode)
		k.RemoveMetaNodeFromBitMapIdxCache(networkAddr)
		return nil
	} else {
		resourceNode, found := k.GetResourceNode(ctx, networkAddr)
		if !found {
			return types.ErrNoResourceNodeFound
		}
		resourceNode.Status = stakingtypes.Unbonded
		k.SetResourceNode(ctx, resourceNode)

		return nil
	}
}
