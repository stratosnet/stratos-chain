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

// SetResourceNodeByPowerIndex resource node index
func (k Keeper) SetResourceNodeByPowerIndex(ctx sdk.Context, resourceNode types.ResourceNode) {
	// suspended resource node are not kept in the power index
	if resourceNode.IsSuspended() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetResourceNodesByPowerIndexKey(resourceNode), resourceNode.GetAddr())
}

// ResourceNode index
func (k Keeper) deleteResourceNodeByPowerIndex(ctx sdk.Context, resourceNode types.ResourceNode) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetResourceNodesByPowerIndexKey(resourceNode))
}

// GetLastResourceNodePower Load the last resource node power.
// Returns zero if the node was not a resource node last block.
func (k Keeper) GetLastResourceNodePower(ctx sdk.Context, nodeAddr sdk.AccAddress) (power int64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastResourceNodePowerKey(nodeAddr))
	if bz == nil {
		return 0
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &power)
	return
}

// SetLastResourceNodePower Set the last resource node power.
func (k Keeper) SetLastResourceNodePower(ctx sdk.Context, nodeAddr sdk.AccAddress, power int64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(power)
	store.Set(types.GetLastResourceNodePowerKey(nodeAddr), bz)
}

// DeleteLastResourceNodePower Delete the last resource node power.
func (k Keeper) DeleteLastResourceNodePower(ctx sdk.Context, nodeAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLastResourceNodePowerKey(nodeAddr))
}

// AddResourceNodeTokens Update the tokens of an existing resource node, update the resource nodes power index key
func (k Keeper) AddResourceNodeTokens(ctx sdk.Context, resourceNode types.ResourceNode, coinToAdd sdk.Coin) error {
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

	oldPow := k.GetLastResourceNodePower(ctx, resourceNode.GetAddr())
	oldTotalPow := k.GetLastResourceNodeTotalPower(ctx)

	k.deleteResourceNodeByPowerIndex(ctx, resourceNode)
	resourceNode = resourceNode.AddToken(coinToAdd.Amount)
	newPow := resourceNode.GetPower()
	newTotalPow := oldTotalPow.Sub(sdk.NewInt(oldPow)).Add(sdk.NewInt(newPow))
	k.SetResourceNode(ctx, resourceNode)
	k.SetResourceNodeByPowerIndex(ctx, resourceNode)
	k.SetLastResourceNodePower(ctx, resourceNode.GetAddr(), newPow)
	k.SetLastResourceNodeTotalPower(ctx, newTotalPow)
	return nil
}

// SubtractResourceNodeTokens Update the tokens of an existing resource node, update the resource nodes power index key
func (k Keeper) SubtractResourceNodeTokens(ctx sdk.Context, resourceNode types.ResourceNode, tokensToRemove sdk.Int) error {
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

	oldPow := k.GetLastResourceNodePower(ctx, resourceNode.GetAddr())
	oldTotalPow := k.GetLastResourceNodeTotalPower(ctx)

	k.deleteResourceNodeByPowerIndex(ctx, resourceNode)
	resourceNode = resourceNode.RemoveToken(tokensToRemove)
	newPow := resourceNode.GetPower()
	newTotalPow := oldTotalPow.Sub(sdk.NewInt(oldPow)).Add(sdk.NewInt(newPow))
	k.SetResourceNode(ctx, resourceNode)
	k.SetResourceNodeByPowerIndex(ctx, resourceNode)

	if resourceNode.GetTokens().IsZero() {
		k.DeleteLastResourceNodePower(ctx, resourceNode.GetAddr())
		err := k.removeResourceNode(ctx, resourceNode.GetAddr())
		if err != nil {
			return err
		}
	} else {
		k.SetLastResourceNodePower(ctx, resourceNode.GetAddr(), newPow)
	}
	k.SetLastResourceNodeTotalPower(ctx, newTotalPow)

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
	store.Delete(types.GetResourceNodesByPowerIndexKey(resourceNode))
	return nil
}
