package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

const indexingNodeCacheSize = 500

// Cache the amino decoding of validators, as it can be the case that repeated slashing calls
// cause many calls to GetValidator, which were shown to throttle the state machine in our
// simulation. Note this is quite biased though, as the simulator does more slashes than a
// live chain should, however we require the slashing to be fast as noone pays gas for it.
type cachedIndexingNode struct {
	indexingNode types.IndexingNode
	marshalled   string // marshalled amino bytes for the validator object (not operator address)
}

func newCachedIndexingNode(indexingNode types.IndexingNode, marshalled string) cachedIndexingNode {
	return cachedIndexingNode{
		indexingNode: indexingNode,
		marshalled:   marshalled,
	}
}

// get a single indexing node
func (k Keeper) GetIndexingNode(ctx sdk.Context, operatorAddr sdk.ValAddress) (indexingNode types.IndexingNode, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetNodeKey(types.NodeTypeIndexing, operatorAddr))

	if value == nil {
		return indexingNode, false
	}

	// If these amino encoded bytes are in the cache, return the cached validator
	strValue := string(value)
	if val, ok := k.indexingNodeCache[strValue]; ok {
		valToReturn := val.indexingNode
		// Doesn't mutate the cache's value
		valToReturn.OperatorAddress = operatorAddr
		return valToReturn, true
	}

	// amino bytes weren't found in cache, so amino unmarshal and add it to the cache
	indexingNode = types.MustUnmarshalIndexingNode(k.cdc, value)
	cachedVal := newCachedIndexingNode(indexingNode, strValue)
	k.indexingNodeCache[strValue] = newCachedIndexingNode(indexingNode, strValue)
	k.indexingNodeCacheList.PushBack(cachedVal)

	// if the cache is too big, pop off the last element from it
	if k.indexingNodeCacheList.Len() > indexingNodeCacheSize {
		valToRemove := k.indexingNodeCacheList.Remove(k.indexingNodeCacheList.Front()).(cachedIndexingNode)
		delete(k.indexingNodeCache, valToRemove.marshalled)
	}

	indexingNode = types.MustUnmarshalIndexingNode(k.cdc, value)
	return indexingNode, true
}

// get a single resource node by node address
func (k Keeper) GetIndexingNodeByAddr(ctx sdk.Context, addr sdk.ConsAddress) (indexingNode types.IndexingNode, found bool) {
	store := ctx.KVStore(k.storeKey)
	opAddr := store.Get(types.GetIndexingNodeByAddrKey(addr))
	if opAddr == nil {
		return indexingNode, false
	}
	return k.GetIndexingNode(ctx, opAddr)
}

// SetIndexingNode set the main record holding indexing node details
func (k Keeper) SetIndexingNode(ctx sdk.Context, indexingNode types.IndexingNode) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalIndexingNode(k.cdc, indexingNode)
	store.Set(types.GetNodeKey(types.NodeTypeIndexing, indexingNode.OperatorAddress), bz)
}

// validator index
func (k Keeper) SetIndexingNodeByAddr(ctx sdk.Context, indexingNode types.IndexingNode) {
	store := ctx.KVStore(k.storeKey)
	addr := sdk.ConsAddress(indexingNode.PubKey.Address())
	store.Set(types.GetIndexingNodeByAddrKey(addr), indexingNode.OperatorAddress)
}

// SetIndexingNodeByPowerIndex IndexingNode index
func (k Keeper) SetIndexingNodeByPowerIndex(ctx sdk.Context, indexingNode types.IndexingNode) {
	// jailed indexing node are not kept in the power index
	if indexingNode.Jailed {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetIndexingNodesByPowerIndexKey(indexingNode), indexingNode.OperatorAddress)
}

// validator index
func (k Keeper) SetNewIndexingNodeByPowerIndex(ctx sdk.Context, indexingNode types.IndexingNode) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetIndexingNodesByPowerIndexKey(indexingNode), indexingNode.OperatorAddress)
}

// DeleteIndexingNodeByPowerIndex IndexingNode index
func (k Keeper) DeleteIndexingNodeByPowerIndex(ctx sdk.Context, indexingNode types.IndexingNode) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetIndexingNodesByPowerIndexKey(indexingNode))
}

// Update the tokens of an existing validator, update the validators power index key
func (k Keeper) AddIndexingNodeTokensAndShares(
	ctx sdk.Context, indexingNode types.IndexingNode, tokensToAdd sdk.Int,
) (nodeOut types.IndexingNode, addedShares sdk.Dec) {

	k.DeleteIndexingNodeByPowerIndex(ctx, indexingNode)
	indexingNode, addedShares = indexingNode.AddTokensToIndexingNode(tokensToAdd)
	k.SetIndexingNode(ctx, indexingNode)
	k.SetIndexingNodeByPowerIndex(ctx, indexingNode)
	return indexingNode, addedShares
}

// Update the tokens of an existing validator, update the validators power index key
func (k Keeper) RemoveIndexingNodeTokensAndShares(
	ctx sdk.Context, indexingNode types.IndexingNode, sharesToRemove sdk.Dec,
) (nodeOut types.IndexingNode, removedTokens sdk.Int) {

	k.DeleteIndexingNodeByPowerIndex(ctx, indexingNode)
	indexingNode, removedTokens = indexingNode.RemoveSharesFromIndexingNode(sharesToRemove)
	k.SetIndexingNode(ctx, indexingNode)
	k.SetIndexingNodeByPowerIndex(ctx, indexingNode)
	return indexingNode, removedTokens
}

// Update the tokens of an existing validator, update the validators power index key
func (k Keeper) RemoveIndexingNodeTokens(
	ctx sdk.Context, indexingNode types.IndexingNode, tokensToRemove sdk.Int,
) types.IndexingNode {

	k.DeleteIndexingNodeByPowerIndex(ctx, indexingNode)
	indexingNode = indexingNode.RemoveTokensFromIndexingNode(tokensToRemove)
	k.SetIndexingNode(ctx, indexingNode)
	k.SetIndexingNodeByPowerIndex(ctx, indexingNode)
	return indexingNode
}
