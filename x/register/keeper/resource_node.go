package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

const resourceNodeCacheSize = 500

// Cache the amino decoding of validators, as it can be the case that repeated slashing calls
// cause many calls to GetValidator, which were shown to throttle the state machine in our
// simulation. Note this is quite biased though, as the simulator does more slashes than a
// live chain should, however we require the slashing to be fast as noone pays gas for it.
type cachedResourceNode struct {
	resourceNode types.ResourceNode
	marshalled   string // marshalled amino bytes for the validator object (not operator address)
}

func newCachedResourceNode(resourceNode types.ResourceNode, marshalled string) cachedResourceNode {
	return cachedResourceNode{
		resourceNode: resourceNode,
		marshalled:   marshalled,
	}
}

// get a single resource node
func (k Keeper) GetResourceNode(ctx sdk.Context, operatorAddr sdk.ValAddress) (resourceNode types.ResourceNode, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetNodeKey(types.NodeTypeResource, operatorAddr))

	if value == nil {
		return resourceNode, false
	}

	// If these amino encoded bytes are in the cache, return the cached validator
	strValue := string(value)
	if val, ok := k.resourceNodeCache[strValue]; ok {
		valToReturn := val.resourceNode
		// Doesn't mutate the cache's value
		valToReturn.OperatorAddress = operatorAddr
		return valToReturn, true
	}

	// amino bytes weren't found in cache, so amino unmarshal and add it to the cache
	resourceNode = types.MustUnmarshalResourceNode(k.cdc, value)
	cachedVal := newCachedResourceNode(resourceNode, strValue)
	k.resourceNodeCache[strValue] = newCachedResourceNode(resourceNode, strValue)
	k.resourceNodeCacheList.PushBack(cachedVal)

	// if the cache is too big, pop off the last element from it
	if k.resourceNodeCacheList.Len() > resourceNodeCacheSize {
		valToRemove := k.resourceNodeCacheList.Remove(k.resourceNodeCacheList.Front()).(cachedResourceNode)
		delete(k.resourceNodeCache, valToRemove.marshalled)
	}

	resourceNode = types.MustUnmarshalResourceNode(k.cdc, value)
	return resourceNode, true
}

// get a single resource node by node address
func (k Keeper) GetResourceNodeByAddr(ctx sdk.Context, addr sdk.ConsAddress) (resourceNode types.ResourceNode, found bool) {
	store := ctx.KVStore(k.storeKey)
	opAddr := store.Get(types.GetResourceNodeByAddrKey(addr))
	if opAddr == nil {
		return resourceNode, false
	}
	return k.GetResourceNode(ctx, opAddr)
}

// SetResourceNode set the main record holding resource node details
func (k Keeper) SetResourceNode(ctx sdk.Context, resourceNode types.ResourceNode) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalResourceNode(k.cdc, resourceNode)
	store.Set(types.GetNodeKey(types.NodeTypeResource, resourceNode.OperatorAddress), bz)
}

// validator index
func (k Keeper) SetResourceNodeByAddr(ctx sdk.Context, resourceNode types.ResourceNode) {
	store := ctx.KVStore(k.storeKey)
	addr := sdk.ConsAddress(resourceNode.PubKey.Address())
	store.Set(types.GetResourceNodeByAddrKey(addr), resourceNode.OperatorAddress)
}

// SetResourceNodeByPowerIndex ResourceNode index
func (k Keeper) SetResourceNodeByPowerIndex(ctx sdk.Context, resourceNode types.ResourceNode) {
	// jailed resource node are not kept in the power index
	if resourceNode.Jailed {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetResourceNodesByPowerIndexKey(resourceNode), resourceNode.OperatorAddress)
}

// validator index
func (k Keeper) SetNewResourceNodeByPowerIndex(ctx sdk.Context, resourceNode types.ResourceNode) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetResourceNodesByPowerIndexKey(resourceNode), resourceNode.OperatorAddress)
}

// DeleteResourceNodeByPowerIndex ResourceNode index
func (k Keeper) DeleteResourceNodeByPowerIndex(ctx sdk.Context, resourceNode types.ResourceNode) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetResourceNodesByPowerIndexKey(resourceNode))
}

// Update the tokens of an existing resource node, update the resource nodes power index key
func (k Keeper) AddResourceNodeTokensAndShares(
	ctx sdk.Context, resourceNode types.ResourceNode, tokensToAdd sdk.Int,
) (nodeOut types.ResourceNode, addedShares sdk.Dec) {

	k.DeleteResourceNodeByPowerIndex(ctx, resourceNode)
	resourceNode, addedShares = resourceNode.AddTokensToResourceNode(tokensToAdd)
	k.SetResourceNode(ctx, resourceNode)
	k.SetResourceNodeByPowerIndex(ctx, resourceNode)
	return resourceNode, addedShares
}

// Update the tokens of an existing resource node, update the resource nodes power index key
func (k Keeper) RemoveResourceNodeTokensAndShares(
	ctx sdk.Context, resourceNode types.ResourceNode, sharesToRemove sdk.Dec,
) (nodeOut types.ResourceNode, removedTokens sdk.Int) {

	k.DeleteResourceNodeByPowerIndex(ctx, resourceNode)
	resourceNode, removedTokens = resourceNode.RemoveSharesFromResourceNode(sharesToRemove)
	k.SetResourceNode(ctx, resourceNode)
	k.SetResourceNodeByPowerIndex(ctx, resourceNode)
	return resourceNode, removedTokens
}

// Update the tokens of an existing resource node, update the resource nodes power index key
func (k Keeper) RemoveResourceNodeTokens(
	ctx sdk.Context, resourceNode types.ResourceNode, tokensToRemove sdk.Int,
) types.ResourceNode {

	k.DeleteResourceNodeByPowerIndex(ctx, resourceNode)
	resourceNode = resourceNode.RemoveTokensFromResourceNode(tokensToRemove)
	k.SetResourceNode(ctx, resourceNode)
	k.SetResourceNodeByPowerIndex(ctx, resourceNode)
	return resourceNode
}
