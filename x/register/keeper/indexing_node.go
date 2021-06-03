package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

const indexingNodeCacheSize = 500

// Cache the amino decoding of indexing nodes, as it can be the case that repeated slashing calls
// cause many calls to GetIndexingNode, which were shown to throttle the state machine in our
// simulation. Note this is quite biased though, as the simulator does more slashes than a
// live chain should, however we require the slashing to be fast as no one pays gas for it.
type cachedIndexingNode struct {
	indexingNode types.IndexingNode
	marshalled   string // marshalled amino bytes for the IndexingNode object (not address)
}

func newCachedIndexingNode(indexingNode types.IndexingNode, marshalled string) cachedIndexingNode {
	return cachedIndexingNode{
		indexingNode: indexingNode,
		marshalled:   marshalled,
	}
}

// GetIndexingNode get a single indexing node
func (k Keeper) GetIndexingNode(ctx sdk.Context, addr sdk.AccAddress) (indexingNode types.IndexingNode, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetIndexingNodeKey(addr))

	if value == nil {
		return indexingNode, false
	}

	// If these amino encoded bytes are in the cache, return the cached indexing node
	strValue := string(value)
	if val, ok := k.indexingNodeCache[strValue]; ok {
		valToReturn := val.indexingNode
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

// set the main record holding indexing node details
func (k Keeper) SetIndexingNode(ctx sdk.Context, indexingNode types.IndexingNode) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalIndexingNode(k.cdc, indexingNode)
	store.Set(types.GetIndexingNodeKey(indexingNode.GetAddr()), bz)
}

// GetLastIndexingNodeStake Load the last indexing node stake.
// Returns zero if the node was not a indexing node last block.
func (k Keeper) GetLastIndexingNodeStake(ctx sdk.Context, nodeAddr sdk.AccAddress) (stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastIndexingNodeStakeKey(nodeAddr))
	if bz == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &stake)
	return
}

// SetLastIndexingNodeStake Set the last indexing node stake.
func (k Keeper) SetLastIndexingNodeStake(ctx sdk.Context, nodeAddr sdk.AccAddress, stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(stake)
	store.Set(types.GetLastIndexingNodeStakeKey(nodeAddr), bz)
}

// DeleteLastIndexingNodeStake Delete the last indexing node stake.
func (k Keeper) DeleteLastIndexingNodeStake(ctx sdk.Context, nodeAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLastIndexingNodeStakeKey(nodeAddr))
}

// GetAllIndexingNodes get the set of all indexing nodes with no limits, used during genesis dump
func (k Keeper) GetAllIndexingNodes(ctx sdk.Context) (indexingNodes []types.IndexingNode) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.IndexingNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalIndexingNode(k.cdc, iterator.Value())
		indexingNodes = append(indexingNodes, node)
	}
	return indexingNodes
}

// IterateLastIndexingNodeStakes Iterate over last indexing node stakes.
func (k Keeper) IterateLastIndexingNodeStakes(ctx sdk.Context, handler func(nodeAddr sdk.AccAddress, stake sdk.Int) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.LastIndexingNodeStakeKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.AccAddress(iter.Key()[len(types.LastIndexingNodeStakeKey):])
		var stake sdk.Int
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &stake)
		if handler(addr, stake) {
			break
		}
	}
}

// AddIndexingNodeStake Update the tokens of an existing indexing node
func (k Keeper) AddIndexingNodeStake(ctx sdk.Context, indexingNode types.IndexingNode, coinToAdd sdk.Coin) error {
	nodeAcc := k.accountKeeper.GetAccount(ctx, indexingNode.GetAddr())
	if nodeAcc == nil {
		ctx.Logger().Info(fmt.Sprintf("create new account: %s", indexingNode.GetAddr()))
		nodeAcc = k.accountKeeper.NewAccountWithAddress(ctx, indexingNode.GetAddr())
		k.accountKeeper.SetAccount(ctx, nodeAcc)
	}

	coins := sdk.NewCoins(coinToAdd)
	hasCoin := k.bankKeeper.HasCoins(ctx, indexingNode.OwnerAddress, coins)
	if !hasCoin {
		return types.ErrInsufficientBalance
	}

	err := k.bankKeeper.SendCoins(ctx, indexingNode.GetOwnerAddr(), indexingNode.GetAddr(), coins)
	if err != nil {
		return err
	}

	oldStake := k.GetLastIndexingNodeStake(ctx, indexingNode.GetAddr())
	oldTotalStake := k.GetLastIndexingNodeTotalStake(ctx)

	indexingNode = indexingNode.AddToken(coinToAdd.Amount)
	newStake := indexingNode.GetTokens()
	newTotalStake := oldTotalStake.Sub(oldStake).Add(newStake)

	k.SetIndexingNode(ctx, indexingNode)
	k.SetLastIndexingNodeStake(ctx, indexingNode.GetAddr(), newStake)
	k.SetLastIndexingNodeTotalStake(ctx, newTotalStake)
	k.increaseOzoneLimitByAddStake(ctx, coinToAdd.Amount)

	return nil
}

// SubtractIndexingNodeStake Update the tokens of an existing indexing node
func (k Keeper) SubtractIndexingNodeStake(ctx sdk.Context, indexingNode types.IndexingNode, tokensToRemove sdk.Int) error {
	ownerAcc := k.accountKeeper.GetAccount(ctx, indexingNode.OwnerAddress)
	if ownerAcc == nil {
		return types.ErrNoOwnerAccountFound
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), tokensToRemove))
	hasCoin := k.bankKeeper.HasCoins(ctx, indexingNode.GetAddr(), coins)
	if !hasCoin {
		return types.ErrInsufficientBalance
	}
	_, err := k.bankKeeper.SubtractCoins(ctx, indexingNode.GetAddr(), coins)
	if err != nil {
		return err
	}
	_, err = k.bankKeeper.AddCoins(ctx, indexingNode.OwnerAddress, coins)
	if err != nil {
		return err
	}

	oldStake := k.GetLastIndexingNodeStake(ctx, indexingNode.GetAddr())
	oldTotalStake := k.GetLastIndexingNodeTotalStake(ctx)

	indexingNode = indexingNode.RemoveToken(tokensToRemove)
	newStake := indexingNode.GetTokens()
	newTotalStake := oldTotalStake.Sub(oldStake).Add(newStake)

	k.SetIndexingNode(ctx, indexingNode)

	if indexingNode.GetTokens().IsZero() {
		k.DeleteLastIndexingNodeStake(ctx, indexingNode.GetAddr())
		err := k.removeIndexingNode(ctx, indexingNode.GetAddr())
		if err != nil {
			return err
		}
	} else {
		k.SetLastIndexingNodeStake(ctx, indexingNode.GetAddr(), newStake)
	}
	k.SetLastIndexingNodeTotalStake(ctx, newTotalStake)
	k.decreaseOzoneLimitBySubtractStake(ctx, tokensToRemove)
	return nil
}

// remove the indexing node record and associated indexes
func (k Keeper) removeIndexingNode(ctx sdk.Context, addr sdk.AccAddress) error {
	// first retrieve the old indexing node record
	indexingNode, found := k.GetIndexingNode(ctx, addr)
	if !found {
		return types.ErrNoIndexingNodeFound
	}

	if indexingNode.Tokens.IsPositive() {
		panic("attempting to remove a indexing node which still contains tokens")
	}

	// delete the old indexing node record
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetIndexingNodeKey(addr))
	return nil
}
