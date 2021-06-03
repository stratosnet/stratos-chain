package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

const resourceNodeCacheSize = 500

// Cache the amino decoding of resource nodes, as it can be the case that repeated slashing calls
// cause many calls to GetResourceNode, which were shown to throttle the state machine in our
// simulation. Note this is quite biased though, as the simulator does more slashes than a
// live chain should, however we require the slashing to be fast as no one pays gas for it.
type cachedResourceNode struct {
	resourceNode types.ResourceNode
	marshalled   string // marshalled amino bytes for the ResourceNode object (not address)
}

func newCachedResourceNode(resourceNode types.ResourceNode, marshalled string) cachedResourceNode {
	return cachedResourceNode{
		resourceNode: resourceNode,
		marshalled:   marshalled,
	}
}

// GetResourceNode get a single resource node
func (k Keeper) GetResourceNode(ctx sdk.Context, addr sdk.AccAddress) (resourceNode types.ResourceNode, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetResourceNodeKey(addr))

	if value == nil {
		return resourceNode, false
	}

	// If these amino encoded bytes are in the cache, return the cached resource node
	strValue := string(value)
	if val, ok := k.resourceNodeCache[strValue]; ok {
		valToReturn := val.resourceNode
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

// set the main record holding resource node details
func (k Keeper) SetResourceNode(ctx sdk.Context, resourceNode types.ResourceNode) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalResourceNode(k.cdc, resourceNode)
	store.Set(types.GetResourceNodeKey(resourceNode.GetAddr()), bz)
}

// GetLastResourceNodeStake Load the last resource node stake.
// Returns zero if the node was not a resource node last block.
func (k Keeper) GetLastResourceNodeStake(ctx sdk.Context, nodeAddr sdk.AccAddress) (stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastResourceNodeStakeKey(nodeAddr))
	if bz == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &stake)
	return
}

// SetLastResourceNodeStake Set the last resource node stake.
func (k Keeper) SetLastResourceNodeStake(ctx sdk.Context, nodeAddr sdk.AccAddress, stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(stake)
	store.Set(types.GetLastResourceNodeStakeKey(nodeAddr), bz)
}

// DeleteLastResourceNodeStake Delete the last resource node stake.
func (k Keeper) DeleteLastResourceNodeStake(ctx sdk.Context, nodeAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLastResourceNodeStakeKey(nodeAddr))
}

// GetAllResourceNodes get the set of all resource nodes with no limits, used during genesis dump
func (k Keeper) GetAllResourceNodes(ctx sdk.Context) (resourceNodes []types.ResourceNode) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ResourceNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalResourceNode(k.cdc, iterator.Value())
		resourceNodes = append(resourceNodes, node)
	}
	return resourceNodes
}

// IterateLastResourceNodeStakes Iterate over last resource node stakes.
func (k Keeper) IterateLastResourceNodeStakes(ctx sdk.Context, handler func(nodeAddr sdk.AccAddress, stake sdk.Int) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.LastResourceNodeStakeKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.AccAddress(iter.Key()[len(types.LastResourceNodeStakeKey):])
		var stake sdk.Int
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &stake)
		if handler(addr, stake) {
			break
		}
	}
}

// AddResourceNodeStake Update the tokens of an existing resource node
func (k Keeper) AddResourceNodeStake(ctx sdk.Context, resourceNode types.ResourceNode, coinToAdd sdk.Coin) error {
	nodeAcc := k.accountKeeper.GetAccount(ctx, resourceNode.GetAddr())
	if nodeAcc == nil {
		k.accountKeeper.NewAccountWithAddress(ctx, resourceNode.GetAddr())
	}

	coins := sdk.NewCoins(coinToAdd)
	hasCoin := k.bankKeeper.HasCoins(ctx, resourceNode.OwnerAddress, coins)
	if !hasCoin {
		return types.ErrInsufficientBalance
	}

	err := k.bankKeeper.SendCoins(ctx, resourceNode.GetOwnerAddr(), resourceNode.GetAddr(), coins)
	if err != nil {
		return err
	}

	oldStake := k.GetLastResourceNodeStake(ctx, resourceNode.GetAddr())
	oldTotalStake := k.GetLastResourceNodeTotalStake(ctx)

	resourceNode = resourceNode.AddToken(coinToAdd.Amount)
	newStake := resourceNode.GetTokens()
	newTotalStake := oldTotalStake.Sub(oldStake).Add(newStake)

	k.SetResourceNode(ctx, resourceNode)
	k.SetLastResourceNodeStake(ctx, resourceNode.GetAddr(), newStake)
	k.SetLastResourceNodeTotalStake(ctx, newTotalStake)
	k.increaseOzoneLimitByAddStake(ctx, coinToAdd.Amount)

	return nil
}

// SubtractResourceNodeStake Update the tokens of an existing resource node
func (k Keeper) SubtractResourceNodeStake(ctx sdk.Context, resourceNode types.ResourceNode, tokensToRemove sdk.Int) error {
	ownerAcc := k.accountKeeper.GetAccount(ctx, resourceNode.OwnerAddress)
	if ownerAcc == nil {
		return types.ErrNoOwnerAccountFound
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), tokensToRemove))
	hasCoin := k.bankKeeper.HasCoins(ctx, resourceNode.GetAddr(), coins)
	if !hasCoin {
		return types.ErrInsufficientBalance
	}
	_, err := k.bankKeeper.SubtractCoins(ctx, resourceNode.GetAddr(), coins)
	if err != nil {
		return err
	}
	_, err = k.bankKeeper.AddCoins(ctx, resourceNode.OwnerAddress, coins)
	if err != nil {
		return err
	}

	oldStake := k.GetLastResourceNodeStake(ctx, resourceNode.GetAddr())
	oldTotalStake := k.GetLastResourceNodeTotalStake(ctx)

	resourceNode = resourceNode.RemoveToken(tokensToRemove)
	newStake := resourceNode.GetTokens()
	newTotalStake := oldTotalStake.Sub(oldStake).Add(newStake)

	k.SetResourceNode(ctx, resourceNode)

	if resourceNode.GetTokens().IsZero() {
		k.DeleteLastResourceNodeStake(ctx, resourceNode.GetAddr())
		err := k.removeResourceNode(ctx, resourceNode.GetAddr())
		if err != nil {
			return err
		}
	} else {
		k.SetLastResourceNodeStake(ctx, resourceNode.GetAddr(), newStake)
	}
	k.SetLastResourceNodeTotalStake(ctx, newTotalStake)
	k.decreaseOzoneLimitBySubtractStake(ctx, tokensToRemove)

	return nil
}

// remove the resource node record and associated indexes
func (k Keeper) removeResourceNode(ctx sdk.Context, addr sdk.AccAddress) error {
	// first retrieve the old resource node record
	resourceNode, found := k.GetResourceNode(ctx, addr)
	if !found {
		return types.ErrNoResourceNodeFound
	}

	if resourceNode.Tokens.IsPositive() {
		panic("attempting to remove a resource node which still contains tokens")
	}

	// delete the old resource node record
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetResourceNodeKey(addr))
	return nil
}
