package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Called in each EndBlock
func (k Keeper) BlockRegisteredNodesUpdates(ctx sdk.Context) []abci.ValidatorUpdate {
	// Remove all mature unbonding nodes from the ubd queue.
	ctx.Logger().Debug("Enter BlockRegisteredNodesUpdates")
	matureUBDs := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, networkAddr := range matureUBDs {
		balances, isIndexingNode, err := k.CompleteUnbondingWithAmount(ctx, networkAddr)
		if err != nil {
			continue
		}
		if isIndexingNode {
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeCompleteUnbondingIndexingNode,
					sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
					sdk.NewAttribute(types.AttributeKeyNetworkAddress, networkAddr.String()),
				),
			)
		} else {
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeCompleteUnbondingResourceNode,
					sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
					sdk.NewAttribute(types.AttributeKeyNetworkAddress, networkAddr.String()),
				),
			)
		}

	}

	// UpdateNode won't create UBD node
	return []abci.ValidatorUpdate{}
}

// Node state transitions
func (k Keeper) bondedToUnbonding(ctx sdk.Context, node interface{}, isIndexingNode bool, coin sdk.Coin) interface{} {
	switch isIndexingNode {
	case true:
		temp := node.(types.IndexingNode)
		if temp.GetStatus() != stakingtypes.Bonded {
			panic(fmt.Sprintf("bad state transition bondedToUnbonding, indexingNode: %v\n", temp))
		}
		return k.beginUnbondingIndexingNode(ctx, &temp, &coin)
	default:
		temp := node.(types.ResourceNode)
		if temp.GetStatus() != stakingtypes.Bonded {
			panic(fmt.Sprintf("bad state transition bondedToUnbonding, resourceNode: %v\n", temp))
		}
		return k.beginUnbondingResourceNode(ctx, &temp, &coin)
	}
}

// switches a Node from unbonding state to unbonded state
func (k Keeper) unbondingToUnbonded(ctx sdk.Context, node interface{}, isIndexingNode bool) interface{} {
	switch isIndexingNode {
	case true:
		temp := node.(types.IndexingNode)
		if temp.GetStatus() != stakingtypes.Unbonding {
			panic(fmt.Sprintf("bad state transition unbondingToBonded, indexingNode: %v\n", temp))
		}
		return k.completeUnbondingNode(ctx, temp, isIndexingNode)
	default:
		temp := node.(types.ResourceNode)
		if temp.GetStatus() != stakingtypes.Unbonding {
			panic(fmt.Sprintf("bad state transition unbondingToBonded, resourceNode: %v\n", temp))
		}
		return k.completeUnbondingNode(ctx, temp, isIndexingNode)
	}
}

// perform all the store operations for when a Node begins unbonding
func (k Keeper) beginUnbondingResourceNode(ctx sdk.Context, resourceNode *types.ResourceNode, coin *sdk.Coin) *types.ResourceNode {
	// set node stat to unbonding, remove token from bonded pool, add token into NotBondedPool
	err := k.RemoveTokenFromPoolWhileUnbondingResourceNode(ctx, *resourceNode, *coin)
	if err != nil {
		return &types.ResourceNode{}
	}

	networkAddr, err := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddr())
	if err != nil {
		return &types.ResourceNode{}
	}
	// trigger hook if registered
	k.AfterNodeBeginUnbonding(ctx, networkAddr, false)
	return resourceNode
}
func (k Keeper) beginUnbondingIndexingNode(ctx sdk.Context, indexingNode *types.IndexingNode, coin *sdk.Coin) *types.IndexingNode {
	// change node stat, remove token from bonded pool, add token into NotBondedPool
	err := k.RemoveTokenFromPoolWhileUnbondingIndexingNode(ctx, *indexingNode, *coin)
	if err != nil {
		return nil
	}
	if err != nil {
		return &types.IndexingNode{}
	}
	networkAddr, err := stratos.SdsAddressFromBech32(indexingNode.GetNetworkAddr())
	if err != nil {
		return &types.IndexingNode{}
	}
	// trigger hook if registered
	k.AfterNodeBeginUnbonding(ctx, networkAddr, true)
	return indexingNode
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
func (k Keeper) completeUnbondingNode(ctx sdk.Context, node interface{}, isIndexingNode bool) interface{} {
	if isIndexingNode {
		temp := node.(types.IndexingNode)
		temp.Status = stakingtypes.Unbonded
		k.SetIndexingNode(ctx, temp)
		return temp
	} else {
		temp := node.(types.ResourceNode)
		temp.Status = stakingtypes.Unbonded
		k.SetResourceNode(ctx, temp)
		return temp
	}
}

// Returns all the validator queue timeslices from time 0 until endTime
func (k Keeper) UBDNodeQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UBDNodeQueueKey, sdk.InclusiveEndBytes(types.GetUBDTimeKey(endTime)))
}

// Returns a concatenated list of all the timeslices before currTime, and deletes the timeslices from the queue
func (k Keeper) GetAllMatureUBDNodeQueue(ctx sdk.Context, currTime time.Time) (matureNetworkAddrs []sdk.AccAddress) {
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	ubdTimesliceIterator := k.UBDNodeQueueIterator(ctx, ctx.BlockHeader().Time)
	defer ubdTimesliceIterator.Close()

	for ; ubdTimesliceIterator.Valid(); ubdTimesliceIterator.Next() {
		timeslice := []sdk.AccAddress{}
		amino.MustUnmarshalBinaryLengthPrefixed(ubdTimesliceIterator.Value(), &timeslice)
		matureNetworkAddrs = append(matureNetworkAddrs, timeslice...)
	}

	return matureNetworkAddrs
}

// Unbonds all the unbonding validators that have finished their unbonding period
func (k Keeper) UnbondAllMatureUBDNodeQueue(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	nodeTimesliceIterator := k.UBDNodeQueueIterator(ctx, ctx.BlockHeader().Time)
	defer nodeTimesliceIterator.Close()

	for ; nodeTimesliceIterator.Valid(); nodeTimesliceIterator.Next() {
		timeslice := []stratos.SdsAddress{}
		amino.MustUnmarshalBinaryLengthPrefixed(nodeTimesliceIterator.Value(), &timeslice)

		for _, networkAddr := range timeslice {
			ubd, found := k.GetUnbondingNode(ctx, networkAddr)
			ubdNetworkAddr, _ := stratos.SdsAddressFromBech32(ubd.NetworkAddr)
			if !found {
				panic("node in the unbonding queue was not found")
			}

			if ubd.IsIndexingNode {

				node, found := k.GetIndexingNode(ctx, ubdNetworkAddr)
				if !found {
					panic("cannot find indexing node " + ubd.NetworkAddr)
				}
				if node.GetStatus() != stakingtypes.Unbonding {
					panic("unexpected node in unbonding queue; status was not unbonding")
				}
				k.unbondingToUnbonded(ctx, node, ubd.IsIndexingNode)
				k.removeIndexingNode(ctx, ubdNetworkAddr)
				_, found1 := k.GetIndexingNode(ctx, ubdNetworkAddr)
				if found1 {
					ctx.Logger().Info("Removed indexing node with addr " + ubd.NetworkAddr)
				}
			} else {
				node, found := k.GetResourceNode(ctx, ubdNetworkAddr)
				if !found {
					panic("cannot find resource node " + ubd.NetworkAddr)
				}
				if node.GetStatus() != stakingtypes.Unbonding {
					panic("unexpected node in unbonding queue; status was not unbonding")
				}
				k.unbondingToUnbonded(ctx, node, ubd.IsIndexingNode)
				k.removeResourceNode(ctx, ubdNetworkAddr)
				_, found1 := k.GetResourceNode(ctx, ubdNetworkAddr)
				if found1 {
					ctx.Logger().Info("Removed resource node with addr " + ubd.NetworkAddr)
				}

			}
		}
		store.Delete(nodeTimesliceIterator.Key())
	}
}
