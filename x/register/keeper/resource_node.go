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
	ctx.Logger().Info("131")
	ownerAcc := k.accountKeeper.GetAccount(ctx, resourceNode.OwnerAddress)
	if ownerAcc == nil {
		return types.ErrNoOwnerAccountFound
	}

	coins := sdk.NewCoins(tokenToSub)

	if resourceNode.GetStatus() == sdk.Unbonded {
		ctx.Logger().Info("132")
		notBondedTokenInPool := k.GetResourceNodeNotBondedToken(ctx)
		if notBondedTokenInPool.IsLT(tokenToSub) {
			return types.ErrInsufficientBalanceOfNotBondedPool
		}
		notBondedTokenInPool = notBondedTokenInPool.Sub(tokenToSub)
		k.SetResourceNodeNotBondedToken(ctx, notBondedTokenInPool)
	} else if resourceNode.GetStatus() == sdk.Bonded {
		ctx.Logger().Info("133")
		bondedTokenInPool := k.GetResourceNodeBondedToken(ctx)
		if bondedTokenInPool.IsLT(tokenToSub) {
			return types.ErrInsufficientBalanceOfBondedPool
		}
		bondedTokenInPool = bondedTokenInPool.Sub(tokenToSub)
		if resourceNode.GetTokens().Equal(tokenToSub.Amount) {
			return types.ErrSubAllTokens
		}
		k.SetResourceNodeNotBondedToken(ctx, bondedTokenInPool)
	} else { // unbonding
		ctx.Logger().Info("1333")
		bondedTokenInPool := k.GetResourceNodeBondedToken(ctx)
		if bondedTokenInPool.IsLT(tokenToSub) {
			return types.ErrInsufficientBalanceOfBondedPool
		}
		bondedTokenInPool = bondedTokenInPool.Sub(tokenToSub)
		if resourceNode.GetTokens().LT(tokenToSub.Amount) {
			return types.ErrInsufficientBalance
		}
		k.SetResourceNodeNotBondedToken(ctx, bondedTokenInPool)
	}

	_, err := k.bankKeeper.AddCoins(ctx, resourceNode.OwnerAddress, coins)
	if err != nil {
		return err
	}

	resourceNode = resourceNode.SubToken(tokenToSub.Amount)
	newStake := resourceNode.GetTokens()

	k.SetResourceNode(ctx, resourceNode)

	ctx.Logger().Info("134")
	if newStake.IsZero() {
		ctx.Logger().Info("135")
		k.DeleteLastResourceNodeStake(ctx, resourceNode.GetNetworkAddr())
		err := k.removeResourceNode(ctx, resourceNode.GetNetworkAddr())
		if err != nil {
			return err
		}
	} else {
		ctx.Logger().Info("136")
		k.SetLastResourceNodeStake(ctx, resourceNode.GetNetworkAddr(), newStake)
	}
	ctx.Logger().Info("137")
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
