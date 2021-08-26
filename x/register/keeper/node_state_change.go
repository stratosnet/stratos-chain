package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Calculate the ValidatorUpdates for the current block
// Called in each EndBlock
func (k Keeper) BlockRegisteredNodesUpdates(ctx sdk.Context) []abci.ValidatorUpdate {
	// Remove all mature unbonding nodes from the ubd queue.
	ctx.Logger().Info("Enter BlockRegisteredNodesUpdates")
	matureUBDs := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, networkAddr := range matureUBDs {
		balances, err := k.CompleteUnbondingWithAmount(ctx, networkAddr)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbondingNode,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyNetworkAddr, networkAddr.String()),
			),
		)
	}

	// UpdateNode won't create UBD node
	return []abci.ValidatorUpdate{}
}

// Validator state transitions

func (k Keeper) bondedToUnbonding(ctx sdk.Context, node interface{}, isIndexingNode bool) interface{} {
	if isIndexingNode {
		temp := node.(types.IndexingNode)
		if temp.GetStatus() != sdk.Bonded {
			panic(fmt.Sprintf("bad state transition bondedToUnbonding, indexingNode: %v\n", temp))
		}
		return k.beginUnbondingIndexingNode(ctx, temp)
	} else {
		temp := node.(types.ResourceNode)
		if temp.GetStatus() != sdk.Bonded {
			panic(fmt.Sprintf("bad state transition bondedToUnbonding, resourceNode: %v\n", temp))
		}
		return k.beginUnbondingResourceNode(ctx, temp)
	}
}

//func (k Keeper) unbondingToBonded(ctx sdk.Context, resourceNode types.ResourceNode) types.ResourceNode {
//	if resourceNode.GetStatus() != sdk.Unbonding {
//		panic(fmt.Sprintf("bad state transition unbondingToBonded, resourceNode: %v\n", resourceNode))
//	}
//	return k.bondNode(ctx, resourceNode)
//}
//
//func (k Keeper) unbondedToBonded(ctx sdk.Context, resourceNode types.ResourceNode) types.ResourceNode {
//	if resourceNode.GetStatus() != sdk.Unbonding {
//		panic(fmt.Sprintf("bad state transition unbondedToBonded, resourceNode: %v\n", resourceNode))
//	}
//	return k.bondNode(ctx, resourceNode)
//}

// switches a validator from unbonding state to unbonded state
func (k Keeper) unbondingToUnbonded(ctx sdk.Context, node interface{}, isIndexingNode bool) interface{} {
	if isIndexingNode {
		temp := node.(types.IndexingNode)
		if temp.GetStatus() != sdk.Unbonding {
			panic(fmt.Sprintf("bad state transition unbondingToBonded, indexingNode: %v\n", temp))
		}
		return k.completeUnbondingNode(ctx, temp, isIndexingNode)
	} else {
		temp := node.(types.ResourceNode)
		if temp.GetStatus() != sdk.Unbonding {
			panic(fmt.Sprintf("bad state transition unbondingToBonded, resourceNode: %v\n", temp))
		}
		return k.completeUnbondingNode(ctx, temp, isIndexingNode)
	}
}

//func (k Keeper) unbondingToUnbonded(ctx sdk.Context, node types.ResourceNode) types.ResourceNode {
//	if node.GetStatus() != sdk.Unbonding {
//		panic(fmt.Sprintf("bad state transition unbondingToBonded, resourceNode: %v\n", node))
//	}
//	return k.completeUnbondingResourceNode(ctx, node)
//}

//// perform all the store operations for when a validator status becomes bonded
//func (k Keeper) bondValidator(ctx sdk.Context, validator types.Validator) types.Validator {
//
//	// delete the validator by power index, as the key will change
//	k.DeleteValidatorByPowerIndex(ctx, validator)
//
//	validator = validator.UpdateStatus(sdk.Bonded)
//
//	// save the now bonded validator record to the two referenced stores
//	k.SetValidator(ctx, validator)
//	k.SetValidatorByPowerIndex(ctx, validator)
//
//	// delete from queue if present
//	k.DeleteValidatorQueue(ctx, validator)
//
//	// trigger hook
//	k.AfterValidatorBonded(ctx, validator.ConsAddress(), validator.OperatorAddress)
//
//	return validator
//}

// perform all the store operations for when a validator begins unbonding
func (k Keeper) beginUnbondingResourceNode(ctx sdk.Context, resourceNode types.ResourceNode) types.ResourceNode {
	resourceNode.Status = sdk.Unbonding
	// save the now unbonded node record and power index
	k.SetResourceNode(ctx, resourceNode)
	// trigger hook if registered
	k.AfterNodeBeginUnbonding(ctx, resourceNode.GetNetworkAddr(), false)
	return resourceNode
}
func (k Keeper) beginUnbondingIndexingNode(ctx sdk.Context, indexingNode types.IndexingNode) types.IndexingNode {
	indexingNode.Status = sdk.Unbonding
	// save the now unbonded node record and power index
	k.SetIndexingNode(ctx, indexingNode)
	// trigger hook if registered
	k.AfterNodeBeginUnbonding(ctx, indexingNode.GetNetworkAddr(), false)
	return indexingNode
}

func calcUnbondingMatureTime(creationTime time.Time, threasholdTime time.Duration, completionTime time.Duration) time.Time {
	threasholdTime = 18 * time.Second
	completionTime = 14 * time.Second
	if creationTime.Add(threasholdTime).After(time.Now()) {
		return creationTime.Add(threasholdTime).Add(completionTime)
	}
	return time.Now().Add(completionTime)
}

// perform all the store operations for when a validator status becomes unbonded
func (k Keeper) completeUnbondingNode(ctx sdk.Context, node interface{}, isIndexingNode bool) interface{} {
	if isIndexingNode {
		temp := node.(types.IndexingNode)
		temp.Status = sdk.Unbonded
		k.SetIndexingNode(ctx, temp)
		return temp
	} else {
		temp := node.(types.ResourceNode)
		temp.Status = sdk.Unbonded
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
func (k Keeper) GetAllMatureUBDNodeQueue(ctx sdk.Context, currTime time.Time) (matureValsAddrs []sdk.ValAddress) {
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	validatorTimesliceIterator := k.UBDNodeQueueIterator(ctx, ctx.BlockHeader().Time)
	defer validatorTimesliceIterator.Close()

	for ; validatorTimesliceIterator.Valid(); validatorTimesliceIterator.Next() {
		timeslice := []sdk.ValAddress{}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(validatorTimesliceIterator.Value(), &timeslice)
		matureValsAddrs = append(matureValsAddrs, timeslice...)
	}

	return matureValsAddrs
}

// Unbonds all the unbonding validators that have finished their unbonding period
func (k Keeper) UnbondAllMatureUBDNodeQueue(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	nodeTimesliceIterator := k.UBDNodeQueueIterator(ctx, ctx.BlockHeader().Time)
	defer nodeTimesliceIterator.Close()

	for ; nodeTimesliceIterator.Valid(); nodeTimesliceIterator.Next() {
		timeslice := []sdk.AccAddress{}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(nodeTimesliceIterator.Value(), &timeslice)

		for _, networkAddr := range timeslice {
			ubd, found := k.GetUnbondingNode(ctx, networkAddr)
			if !found {
				panic("node in the unbonding queue was not found")
			}

			if ubd.IsIndexingNode {
				node, found := k.GetIndexingNode(ctx, ubd.NetworkAddr)
				if !found {
					panic("cannot find indexing node " + ubd.NetworkAddr.String())
				}
				if node.GetStatus() != sdk.Unbonding {
					panic("unexpected node in unbonding queue; status was not unbonding")
				}
				k.unbondingToUnbonded(ctx, node, ubd.IsIndexingNode)
				k.removeIndexingNode(ctx, ubd.NetworkAddr)
			} else {
				node, found := k.GetResourceNode(ctx, ubd.NetworkAddr)
				if !found {
					panic("cannot find resource node " + ubd.NetworkAddr.String())
				}
				if node.GetStatus() != sdk.Unbonding {
					panic("unexpected node in unbonding queue; status was not unbonding")
				}
				k.unbondingToUnbonded(ctx, node, ubd.IsIndexingNode)
				k.removeResourceNode(ctx, ubd.NetworkAddr)
			}
		}
		store.Delete(nodeTimesliceIterator.Key())
	}
}
