package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"strings"
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
	store.Set(types.GetIndexingNodeKey(indexingNode.GetNetworkAddr()), bz)
}

// IndexingNode index
func (k Keeper) SetIndexingNodeByPowerIndex(ctx sdk.Context, indexingNode types.IndexingNode) {
	// suspended indexing node are not kept in the power index
	if indexingNode.IsSuspended() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetIndexingNodesByPowerIndexKey(indexingNode), indexingNode.GetNetworkAddr())
}

// IndexingNode index
func (k Keeper) deleteIndexingNodeByPowerIndex(ctx sdk.Context, indexingNode types.IndexingNode) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetIndexingNodesByPowerIndexKey(indexingNode))
}

// GetLastIndexingNodePower Load the last indexing node power.
// Returns zero if the node was not a indexing node last block.
func (k Keeper) GetLastIndexingNodePower(ctx sdk.Context, nodeAddr sdk.AccAddress) (power int64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastIndexingNodePowerKey(nodeAddr))
	if bz == nil {
		return 0
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &power)
	return
}

// SetLastIndexingNodePower Set the last indexing node power.
func (k Keeper) SetLastIndexingNodePower(ctx sdk.Context, nodeAddr sdk.AccAddress, power int64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(power)
	store.Set(types.GetLastIndexingNodePowerKey(nodeAddr), bz)
}

// DeleteLastIndexingNodePower Delete the last indexing node power.
func (k Keeper) DeleteLastIndexingNodePower(ctx sdk.Context, nodeAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLastIndexingNodePowerKey(nodeAddr))
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

// IterateLastIndexingNodePowers Iterate over last indexing node powers.
func (k Keeper) IterateLastIndexingNodePowers(ctx sdk.Context, handler func(nodeAddr sdk.AccAddress, power int64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.LastIndexingNodePowerKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.AccAddress(iter.Key()[len(types.LastIndexingNodePowerKey):])
		var power int64
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &power)
		if handler(addr, power) {
			break
		}
	}
}

// AddIndexingNodeTokens Update the tokens of an existing indexing node, update the indexing nodes power index key
func (k Keeper) AddIndexingNodeTokens(ctx sdk.Context, indexingNode types.IndexingNode, coinToAdd sdk.Coin) error {
	nodeAcc := k.accountKeeper.GetAccount(ctx, indexingNode.GetNetworkAddr())
	if nodeAcc == nil {
		ctx.Logger().Info(fmt.Sprintf("create new account: %s", indexingNode.GetNetworkAddr()))
		nodeAcc = k.accountKeeper.NewAccountWithAddress(ctx, indexingNode.GetNetworkAddr())
		k.accountKeeper.SetAccount(ctx, nodeAcc)
	}

	coins := sdk.NewCoins(coinToAdd)
	hasCoin := k.bankKeeper.HasCoins(ctx, indexingNode.OwnerAddress, coins)
	if !hasCoin {
		return types.ErrInsufficientBalance
	}

	err := k.bankKeeper.SendCoins(ctx, indexingNode.GetOwnerAddr(), indexingNode.GetNetworkAddr(), coins)
	if err != nil {
		return err
	}

	oldPow := k.GetLastIndexingNodePower(ctx, indexingNode.GetNetworkAddr())
	oldTotalPow := k.GetLastIndexingNodeTotalPower(ctx)

	k.deleteIndexingNodeByPowerIndex(ctx, indexingNode)
	indexingNode = indexingNode.AddToken(coinToAdd.Amount)
	newPow := indexingNode.GetPower()
	newTotalPow := oldTotalPow.Sub(sdk.NewInt(oldPow)).Add(sdk.NewInt(newPow))
	k.SetIndexingNode(ctx, indexingNode)
	k.SetIndexingNodeByPowerIndex(ctx, indexingNode)
	k.SetLastIndexingNodePower(ctx, indexingNode.GetNetworkAddr(), newPow)
	k.SetLastIndexingNodeTotalPower(ctx, newTotalPow)
	return nil
}

// SubtractIndexingNodeTokens Update the tokens of an existing indexing node, update the indexing nodes power index key
func (k Keeper) SubtractIndexingNodeTokens(ctx sdk.Context, indexingNode types.IndexingNode, tokensToRemove sdk.Int) error {
	ownerAcc := k.accountKeeper.GetAccount(ctx, indexingNode.OwnerAddress)
	if ownerAcc == nil {
		return types.ErrNoOwnerAccountFound
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), tokensToRemove))
	hasCoin := k.bankKeeper.HasCoins(ctx, indexingNode.GetNetworkAddr(), coins)
	if !hasCoin {
		return types.ErrInsufficientBalance
	}
	_, err := k.bankKeeper.SubtractCoins(ctx, indexingNode.GetNetworkAddr(), coins)
	if err != nil {
		return err
	}
	_, err = k.bankKeeper.AddCoins(ctx, indexingNode.OwnerAddress, coins)
	if err != nil {
		return err
	}

	oldPow := k.GetLastIndexingNodePower(ctx, indexingNode.GetNetworkAddr())
	oldTotalPow := k.GetLastIndexingNodeTotalPower(ctx)

	k.deleteIndexingNodeByPowerIndex(ctx, indexingNode)
	indexingNode = indexingNode.RemoveToken(tokensToRemove)
	newPow := indexingNode.GetPower()
	newTotalPow := oldTotalPow.Sub(sdk.NewInt(oldPow)).Add(sdk.NewInt(newPow))
	k.SetIndexingNode(ctx, indexingNode)
	k.SetIndexingNodeByPowerIndex(ctx, indexingNode)

	if indexingNode.GetTokens().IsZero() {
		k.DeleteLastIndexingNodePower(ctx, indexingNode.GetNetworkAddr())
		err := k.removeIndexingNode(ctx, indexingNode.GetNetworkAddr())
		if err != nil {
			return err
		}
	} else {
		k.SetLastIndexingNodePower(ctx, indexingNode.GetNetworkAddr(), newPow)
	}
	k.SetLastIndexingNodeTotalPower(ctx, newTotalPow)

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
	store.Delete(types.GetIndexingNodesByPowerIndexKey(indexingNode))
	return nil
}

// GetIndexingNodeList get all indexing nodes by network ID
func (k Keeper) GetIndexingNodeList(ctx sdk.Context, networkID string) (indexingNodes []types.IndexingNode, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.IndexingNodeKey)
	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalIndexingNode(k.cdc, iterator.Value())
		if strings.Compare(node.NetworkID, networkID) == 0 {
			indexingNodes = append(indexingNodes, node)
		}

	}
	ctx.Logger().Info("IndexingNodeList: "+networkID, types.ModuleCdc.MustMarshalJSON(indexingNodes))
	return indexingNodes, nil
}

func (k Keeper) GetIndexingNodeListByMoniker(ctx sdk.Context, moniker string) (resourceNodes []types.IndexingNode, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.IndexingNodeKey)
	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalIndexingNode(k.cdc, iterator.Value())
		if strings.Compare(node.Description.Moniker, moniker) == 0 {
			resourceNodes = append(resourceNodes, node)
		}
	}
	ctx.Logger().Info("resourceNodeList: "+moniker, types.ModuleCdc.MustMarshalJSON(resourceNodes))
	return resourceNodes, nil
}
