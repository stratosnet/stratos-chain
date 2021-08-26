package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"github.com/tendermint/tendermint/crypto"
	"strings"
	"time"
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

// SetResourceNode sets the main record holding resource node details
func (k Keeper) SetResourceNode(ctx sdk.Context, resourceNode types.ResourceNode) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalResourceNode(k.cdc, resourceNode)
	store.Set(types.GetResourceNodeKey(resourceNode.GetNetworkAddr()), bz)
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
func (k Keeper) AddResourceNodeStake(ctx sdk.Context, resourceNode types.ResourceNode, tokenToAdd sdk.Coin) error {
	nodeAcc := k.accountKeeper.GetAccount(ctx, resourceNode.GetNetworkAddr())
	if nodeAcc == nil {
		ctx.Logger().Info(fmt.Sprintf("create new account: %s", resourceNode.GetNetworkAddr()))
		nodeAcc = k.accountKeeper.NewAccountWithAddress(ctx, resourceNode.GetNetworkAddr())
		k.accountKeeper.SetAccount(ctx, nodeAcc)
	}

	coins := sdk.NewCoins(tokenToAdd)
	hasCoin := k.bankKeeper.HasCoins(ctx, resourceNode.GetOwnerAddr(), coins)
	if !hasCoin {
		return types.ErrInsufficientBalance
	}

	_, err := k.bankKeeper.SubtractCoins(ctx, resourceNode.GetOwnerAddr(), coins)
	if err != nil {
		return err
	}

	//TODO: Add logic for unBonding status
	if resourceNode.GetStatus() == sdk.Unbonded {
		notBondedTokenInPool := k.GetResourceNodeNotBondedToken(ctx)
		notBondedTokenInPool = notBondedTokenInPool.Add(tokenToAdd)
		k.SetResourceNodeNotBondedToken(ctx, notBondedTokenInPool)
	} else if resourceNode.GetStatus() == sdk.Bonded {
		bondedTokenInPool := k.GetResourceNodeBondedToken(ctx)
		bondedTokenInPool = bondedTokenInPool.Add(tokenToAdd)
		k.SetResourceNodeBondedToken(ctx, bondedTokenInPool)
	}

	resourceNode = resourceNode.AddToken(tokenToAdd.Amount)

	// set status from unBonded to bonded & move stake from not bonded token pool to bonded token pool
	// since resource node registration does not require voting for now
	if resourceNode.Status.Equal(sdk.Unbonded) {
		resourceNode.Status = sdk.Bonded

		tokenToBond := sdk.NewCoin(k.BondDenom(ctx), resourceNode.GetTokens())
		notBondedToken := k.GetResourceNodeNotBondedToken(ctx)
		bondedToken := k.GetResourceNodeBondedToken(ctx)

		if notBondedToken.IsLT(tokenToBond) {
			return types.ErrInsufficientBalanceOfNotBondedPool
		}
		notBondedToken = notBondedToken.Sub(tokenToBond)
		bondedToken = bondedToken.Add(tokenToBond)
		k.SetResourceNodeNotBondedToken(ctx, notBondedToken)
		k.SetResourceNodeBondedToken(ctx, bondedToken)
	}

	newStake := resourceNode.GetTokens()

	k.SetResourceNode(ctx, resourceNode)
	k.SetLastResourceNodeStake(ctx, resourceNode.GetNetworkAddr(), newStake)
	k.increaseOzoneLimitByAddStake(ctx, tokenToAdd.Amount)

	return nil
}

// SubtractResourceNodeStake Update the tokens of an existing resource node
func (k Keeper) SubtractResourceNodeStake(ctx sdk.Context, resourceNode types.ResourceNode, tokenToSub sdk.Coin) error {
	ownerAcc := k.accountKeeper.GetAccount(ctx, resourceNode.OwnerAddress)
	if ownerAcc == nil {
		return types.ErrNoOwnerAccountFound
	}

	coins := sdk.NewCoins(tokenToSub)

	//TODO: Add logic for subtracting tokens from node already bonded
	if resourceNode.GetStatus() == sdk.Unbonded {
		notBondedTokenInPool := k.GetResourceNodeNotBondedToken(ctx)
		if notBondedTokenInPool.IsLT(tokenToSub) {
			return types.ErrInsufficientBalanceOfNotBondedPool
		}
		notBondedTokenInPool = notBondedTokenInPool.Sub(tokenToSub)
		k.SetResourceNodeNotBondedToken(ctx, notBondedTokenInPool)
	} else if resourceNode.GetStatus() == sdk.Bonded {
		bondedTokenInPool := k.GetResourceNodeBondedToken(ctx)
		if bondedTokenInPool.IsLT(tokenToSub) {
			return types.ErrInsufficientBalanceOfBondedPool
		}
		bondedTokenInPool = bondedTokenInPool.Sub(tokenToSub)
		if resourceNode.GetTokens().Equal(tokenToSub.Amount) {
			return types.ErrSubAllTokens
		}
		k.SetIndexingNodeBondedToken(ctx, bondedTokenInPool)
	}

	_, err := k.bankKeeper.AddCoins(ctx, resourceNode.OwnerAddress, coins)
	if err != nil {
		return err
	}

	resourceNode = resourceNode.SubToken(tokenToSub.Amount)
	newStake := resourceNode.GetTokens()

	k.SetResourceNode(ctx, resourceNode)

	if resourceNode.GetTokens().IsZero() {
		k.DeleteLastResourceNodeStake(ctx, resourceNode.GetNetworkAddr())
		err := k.removeResourceNode(ctx, resourceNode.GetNetworkAddr())
		if err != nil {
			return err
		}
	} else {
		k.SetLastResourceNodeStake(ctx, resourceNode.GetNetworkAddr(), newStake)
	}
	k.decreaseOzoneLimitBySubtractStake(ctx, tokenToSub.Amount)

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

// GetResourceNodeList get all resource nodes by network address
func (k Keeper) GetResourceNodeList(ctx sdk.Context, networkID string) (resourceNodes []types.ResourceNode, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ResourceNodeKey)
	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalResourceNode(k.cdc, iterator.Value())
		if strings.Compare(node.NetworkID, networkID) == 0 {
			resourceNodes = append(resourceNodes, node)
		}

	}
	ctx.Logger().Info("resourceNodeList: "+networkID, types.ModuleCdc.MustMarshalJSON(resourceNodes))
	return resourceNodes, nil
}

func (k Keeper) GetResourceNodeListByMoniker(ctx sdk.Context, moniker string) (resourceNodes []types.ResourceNode, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ResourceNodeKey)
	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalResourceNode(k.cdc, iterator.Value())
		if strings.Compare(node.Description.Moniker, moniker) == 0 {
			resourceNodes = append(resourceNodes, node)
		}
	}
	ctx.Logger().Info("resourceNodeList: "+moniker, types.ModuleCdc.MustMarshalJSON(resourceNodes))
	return resourceNodes, nil
}

func (k Keeper) RegisterResourceNode(ctx sdk.Context, networkID string, pubKey crypto.PubKey, ownerAddr sdk.AccAddress,
	description types.Description, nodeType string, stake sdk.Coin) error {

	resourceNode := types.NewResourceNode(networkID, pubKey, ownerAddr, description, nodeType, time.Now())
	err := k.AddResourceNodeStake(ctx, resourceNode, stake)
	return err
}

func (k Keeper) UpdateResourceNode(ctx sdk.Context, networkID string, description types.Description, nodeType string,
	networkAddr sdk.AccAddress, ownerAddr sdk.AccAddress) error {

	node, found := k.GetResourceNode(ctx, networkAddr)
	if !found {
		return types.ErrNoResourceNodeFound
	}

	if !node.OwnerAddress.Equals(ownerAddr) {
		return types.ErrInvalidOwnerAddr
	}

	node.NetworkID = networkID
	node.Description = description
	node.NodeType = nodeType

	k.SetResourceNode(ctx, node)

	return nil
}

func (k Keeper) SetResourceNodeBondedToken(ctx sdk.Context, token sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token)
	store.Set(types.ResourceNodeBondedTokenKey, bz)
}

func (k Keeper) GetResourceNodeBondedToken(ctx sdk.Context) (token sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ResourceNodeBondedTokenKey)
	if bz == nil {
		return sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &token)
	return token
}

func (k Keeper) SetResourceNodeNotBondedToken(ctx sdk.Context, token sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token)
	store.Set(types.ResourceNodeNotBondedTokenKey, bz)
}

func (k Keeper) GetResourceNodeNotBondedToken(ctx sdk.Context) (token sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ResourceNodeNotBondedTokenKey)
	if bz == nil {
		return sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &token)
	return token
}

//
//
//// return a given amount of all the UnbondingResourceNodes
//func (k Keeper) GetURNs(ctx sdk.Context, networkAddr sdk.AccAddress,
//	maxRetrieve uint16) (unbondingResourceNodes []types.UnbondingResourceNode) {
//
//	unbondingResourceNodes = make([]types.UnbondingResourceNode, maxRetrieve)
//
//	store := ctx.KVStore(k.storeKey)
//	indexingNodePrefixKey := types.GetURNKey(networkAddr)
//	iterator := sdk.KVStorePrefixIterator(store, indexingNodePrefixKey)
//	defer iterator.Close()
//
//	i := 0
//	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
//		unbondingResourceNode := types.MustUnmarshalURN(k.cdc, iterator.Value())
//		unbondingResourceNodes[i] = unbondingResourceNode
//		i++
//	}
//	return unbondingResourceNodes[:i] // trim if the array length < maxRetrieve
//}
//
//// return a unbonding UnbondingResourceNode
//func (k Keeper) GetURN(ctx sdk.Context,
//	networkAddr sdk.AccAddress) (ubd types.UnbondingResourceNode, found bool) {
//
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetURNKey(networkAddr)
//	value := store.Get(key)
//	if value == nil {
//		return ubd, false
//	}
//
//	ubd = types.MustUnmarshalURN(k.cdc, value)
//	return ubd, true
//}
//
//// iterate through all of the unbonding indexingNodes
//func (k Keeper) IterateURNs(ctx sdk.Context, fn func(index int64, ubd types.UnbondingResourceNode) (stop bool)) {
//	store := ctx.KVStore(k.storeKey)
//	iterator := sdk.KVStorePrefixIterator(store, types.UBDResourceNodeKey)
//	defer iterator.Close()
//
//	for i := int64(0); iterator.Valid(); iterator.Next() {
//		ubd := types.MustUnmarshalURN(k.cdc, iterator.Value())
//		if stop := fn(i, ubd); stop {
//			break
//		}
//		i++
//	}
//}
//
//// HasMaxUnbondingResourceNodeEntries - check if unbonding ResourceNode has maximum number of entries
//func (k Keeper) HasMaxURNEntries(ctx sdk.Context, networkAddr sdk.AccAddress) bool {
//	ubd, found := k.GetURN(ctx, networkAddr)
//	if !found {
//		return false
//	}
//	return len(ubd.Entries) >= int(k.MaxEntries(ctx))
//}
//
//// set the unbonding ResourceNode
//func (k Keeper) SetURN(ctx sdk.Context, ubd types.UnbondingResourceNode) {
//	store := ctx.KVStore(k.storeKey)
//	bz := types.MustMarshalURN(k.cdc, ubd)
//	key := types.GetURNKey(ubd.GetNetworkAddr())
//	store.Set(key, bz)
//}
//
//// remove the unbonding ResourceNode object
//func (k Keeper) RemoveURN(ctx sdk.Context, ubd types.UnbondingResourceNode) {
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetURNKey(ubd.GetNetworkAddr())
//	store.Delete(key)
//}
//
//// SetUnbondingResourceNodeEntry adds an entry to the unbonding ResourceNode at
//// the given addresses. It creates the unbonding ResourceNode if it does not exist
//func (k Keeper) SetURNEntry(ctx sdk.Context, networkAddr sdk.AccAddress,
//	creationHeight int64, minTime time.Time, balance sdk.Int) types.UnbondingResourceNode {
//
//	ubd, found := k.GetURN(ctx, networkAddr)
//	if found {
//		ubd.AddEntry(creationHeight, minTime, balance)
//	} else {
//		ubd = types.NewUnbondingResourceNode(networkAddr, creationHeight, minTime, balance)
//	}
//	k.SetURN(ctx, ubd)
//	return ubd
//}
//
//// unbonding delegation queue timeslice operations
//
//// gets a specific unbonding queue timeslice. A timeslice is a slice of DVPairs
//// corresponding to unbonding delegations that expire at a certain time.
//func (k Keeper) GetURNQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (networkAddrs []sdk.AccAddress) {
//	store := ctx.KVStore(k.storeKey)
//	bz := store.Get(types.GetURNTimeKey(timestamp))
//	if bz == nil {
//		return []sdk.AccAddress{}
//	}
//	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &networkAddrs)
//	return networkAddrs
//}
//
//// Sets a specific unbonding queue timeslice.
//func (k Keeper) SetURNQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []sdk.AccAddress) {
//	store := ctx.KVStore(k.storeKey)
//	bz := k.cdc.MustMarshalBinaryLengthPrefixed(keys)
//	store.Set(types.GetURNTimeKey(timestamp), bz)
//}
//
//// Insert an unbonding delegation to the appropriate timeslice in the unbonding queue
//func (k Keeper) InsertURNQueue(ctx sdk.Context, ubd types.UnbondingResourceNode,
//	completionTime time.Time) {
//
//	timeSlice := k.GetURNQueueTimeSlice(ctx, completionTime)
//	networkAddr := ubd.NetworkAddr
//	if len(timeSlice) == 0 {
//		k.SetURNQueueTimeSlice(ctx, completionTime, []sdk.AccAddress{networkAddr})
//	} else {
//		timeSlice = append(timeSlice, networkAddr)
//		k.SetURNQueueTimeSlice(ctx, completionTime, timeSlice)
//	}
//}
//
//// Returns all the unbonding queue timeslices from time 0 until endTime
//func (k Keeper) URNQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
//	store := ctx.KVStore(k.storeKey)
//	return store.Iterator(types.UBDResourceNodeQueueKey,
//		sdk.InclusiveEndBytes(types.GetURNTimeKey(endTime)))
//}
//
//// Returns a concatenated list of all the timeslices inclusively previous to
//// currTime, and deletes the timeslices from the queue
//func (k Keeper) DequeueAllMatureURNQueue(ctx sdk.Context,
//	currTime time.Time) (matureUnbonds []sdk.AccAddress) {
//
//	store := ctx.KVStore(k.storeKey)
//	// gets an iterator for all timeslices from time 0 until the current Blockheader time
//	unbondingTimesliceIterator := k.URNQueueIterator(ctx, ctx.BlockHeader().Time)
//	defer unbondingTimesliceIterator.Close()
//
//	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
//		timeslice := []sdk.AccAddress{}
//		value := unbondingTimesliceIterator.Value()
//		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)
//		matureUnbonds = append(matureUnbonds, timeslice...)
//		store.Delete(unbondingTimesliceIterator.Key())
//	}
//	return matureUnbonds
//}

//
//func (k Keeper) DoRemoveResourceNode(
//	ctx sdk.Context, resourceNode register.ResourceNode, amt sdk.Int,
//) (time.Time, error) {
//
//	ownerAcc := k.accountKeeper.GetAccount(ctx, resourceNode.OwnerAddress)
//	if ownerAcc == nil {
//		return time.Time{}, types.ErrNoOwnerAccountFound
//	}
//
//	networkAddr := resourceNode.GetNetworkAddr()
//	if k.HasMaxUnbondingNodeEntries(ctx, networkAddr) {
//		return time.Time{}, types.ErrMaxUnbondingNodeEntries
//	}
//
//	returnAmount, err := k.unbond(ctx, networkAddr, false, amt)
//	if err != nil {
//		return time.Time{}, err
//	}
//
//	// transfer the node tokens to the not bonded pool
//	if resourceNode.GetStatus() == sdk.Bonded {
//		k.bondedToUnbonding(ctx, resourceNode, false)
//	}
//
//	params := k.GetParams(ctx)
//	// set the unbonding mature time and completion height appropriately
//	unbondingMatureTime := calcUnbondingMatureTime(resourceNode.CreationTime, params.UnbondingThreasholdTime, params.UnbondingCompletionTime)
//	unbondingNode := types.NewUnbondingNode(resourceNode.GetNetworkAddr(), false, ctx.BlockHeight(), unbondingMatureTime, returnAmount)
//	// Adds to unbonding node queue
//	k.InsertUnbondingNodeQueue(ctx, unbondingNode, unbondingMatureTime)
//
//	return unbondingMatureTime, nil
//}
