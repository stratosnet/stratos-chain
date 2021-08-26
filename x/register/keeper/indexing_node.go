package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"github.com/tendermint/tendermint/crypto"
	"strings"
	"time"
)

const (
	indexingNodeCacheSize        = 500
	votingValidityPeriodInSecond = 7 * 24 * 60 * 60 // 7 days
)

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

// GetAllValidIndexingNodes get the set of all bonded & not suspended indexing nodes
func (k Keeper) GetAllValidIndexingNodes(ctx sdk.Context) (indexingNodes []types.IndexingNode) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.IndexingNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalIndexingNode(k.cdc, iterator.Value())
		if !node.IsSuspended() && node.GetStatus().Equal(sdk.Bonded) {
			indexingNodes = append(indexingNodes, node)
		}
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

func (k Keeper) RegisterIndexingNode(ctx sdk.Context, networkID string, pubKey crypto.PubKey, ownerAddr sdk.AccAddress,
	description types.Description, stake sdk.Coin) error {

	indexingNode := types.NewIndexingNode(networkID, pubKey, ownerAddr, description, time.Now())

	err := k.AddIndexingNodeStake(ctx, indexingNode, stake)
	if err != nil {
		return err
	}

	var approveList = make([]sdk.AccAddress, 0)
	var rejectList = make([]sdk.AccAddress, 0)
	votingValidityPeriod := votingValidityPeriodInSecond * time.Second
	expireTime := time.Now().Add(votingValidityPeriod)

	votePool := types.NewRegistrationVotePool(indexingNode.GetNetworkAddr(), approveList, rejectList, expireTime)
	k.SetIndexingNodeRegistrationVotePool(ctx, votePool)

	return nil
}

// AddIndexingNodeStake Update the tokens of an existing indexing node
func (k Keeper) AddIndexingNodeStake(ctx sdk.Context, indexingNode types.IndexingNode, tokenToAdd sdk.Coin) error {
	nodeAcc := k.accountKeeper.GetAccount(ctx, indexingNode.GetNetworkAddr())
	if nodeAcc == nil {
		ctx.Logger().Info(fmt.Sprintf("create new account: %s", indexingNode.GetNetworkAddr()))
		nodeAcc = k.accountKeeper.NewAccountWithAddress(ctx, indexingNode.GetNetworkAddr())
		k.accountKeeper.SetAccount(ctx, nodeAcc)
	}

	coins := sdk.NewCoins(tokenToAdd)
	hasCoin := k.bankKeeper.HasCoins(ctx, indexingNode.OwnerAddress, coins)
	if !hasCoin {
		return types.ErrInsufficientBalance
	}

	_, err := k.bankKeeper.SubtractCoins(ctx, indexingNode.GetOwnerAddr(), coins)
	if err != nil {
		return err
	}

	indexingNode = indexingNode.AddToken(tokenToAdd.Amount)

	//TODO: Add logic for unBonding status
	if indexingNode.GetStatus() == sdk.Unbonded {
		notBondedTokenInPool := k.GetIndexingNodeNotBondedToken(ctx)
		notBondedTokenInPool = notBondedTokenInPool.Add(tokenToAdd)
		k.SetIndexingNodeNotBondedToken(ctx, notBondedTokenInPool)
	} else if indexingNode.GetStatus() == sdk.Bonded {
		bondedTokenInPool := k.GetIndexingNodeBondedToken(ctx)
		bondedTokenInPool = bondedTokenInPool.Add(tokenToAdd)
		k.SetIndexingNodeBondedToken(ctx, bondedTokenInPool)
	}

	newStake := indexingNode.GetTokens()

	k.SetIndexingNode(ctx, indexingNode)
	k.SetLastIndexingNodeStake(ctx, indexingNode.GetNetworkAddr(), newStake)
	k.increaseOzoneLimitByAddStake(ctx, tokenToAdd.Amount)

	return nil
}

// SubtractIndexingNodeStake Update the tokens of an existing indexing node
func (k Keeper) SubtractIndexingNodeStake(ctx sdk.Context, indexingNode types.IndexingNode, tokenToSub sdk.Coin) error {
	ownerAcc := k.accountKeeper.GetAccount(ctx, indexingNode.OwnerAddress)
	if ownerAcc == nil {
		return types.ErrNoOwnerAccountFound
	}

	coins := sdk.NewCoins(tokenToSub)

	//TODO: Add logic for unBonding status
	if indexingNode.GetStatus() == sdk.Unbonded {
		notBondedTokenInPool := k.GetIndexingNodeNotBondedToken(ctx)
		if notBondedTokenInPool.IsLT(tokenToSub) {
			return types.ErrInsufficientBalanceOfNotBondedPool
		}
		notBondedTokenInPool = notBondedTokenInPool.Sub(tokenToSub)
		k.SetIndexingNodeNotBondedToken(ctx, notBondedTokenInPool)
	} else if indexingNode.GetStatus() == sdk.Bonded {
		bondedTokenInPool := k.GetIndexingNodeBondedToken(ctx)
		if bondedTokenInPool.IsLT(tokenToSub) {
			return types.ErrInsufficientBalanceOfBondedPool
		}
		bondedTokenInPool = bondedTokenInPool.Sub(tokenToSub)
		if indexingNode.GetTokens().Equal(tokenToSub.Amount) {
			return types.ErrSubAllTokens
		}
		k.SetIndexingNodeBondedToken(ctx, bondedTokenInPool)
	}

	_, err := k.bankKeeper.AddCoins(ctx, indexingNode.OwnerAddress, coins)
	if err != nil {
		return err
	}

	indexingNode = indexingNode.SubToken(tokenToSub.Amount)
	newStake := indexingNode.GetTokens()

	k.SetIndexingNode(ctx, indexingNode)

	if indexingNode.GetTokens().IsZero() {
		k.DeleteLastIndexingNodeStake(ctx, indexingNode.GetNetworkAddr())
		err := k.removeIndexingNode(ctx, indexingNode.GetNetworkAddr())
		if err != nil {
			return err
		}
	} else {
		k.SetLastIndexingNodeStake(ctx, indexingNode.GetNetworkAddr(), newStake)
	}

	k.decreaseOzoneLimitBySubtractStake(ctx, tokenToSub.Amount)
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
	ctx.Logger().Info("indexingNodeList: "+moniker, types.ModuleCdc.MustMarshalJSON(resourceNodes))
	return resourceNodes, nil
}

func (k Keeper) HandleVoteForIndexingNodeRegistration(ctx sdk.Context, nodeAddr sdk.AccAddress, ownerAddr sdk.AccAddress,
	opinion types.VoteOpinion, voterAddr sdk.AccAddress) (nodeStatus sdk.BondStatus, err error) {

	votePool, found := k.GetIndexingNodeRegistrationVotePool(ctx, nodeAddr)
	if !found {
		return sdk.Unbonded, types.ErrNoRegistrationVotePoolFound
	}
	if votePool.ExpireTime.Before(time.Now()) {
		return sdk.Unbonded, types.ErrVoteExpired
	}
	if k.hasValue(votePool.ApproveList, voterAddr) || k.hasValue(votePool.RejectList, voterAddr) {
		return sdk.Unbonded, types.ErrDuplicateVoting
	}

	node, found := k.GetIndexingNode(ctx, nodeAddr)
	if !found {
		return sdk.Unbonded, types.ErrNoIndexingNodeFound
	}
	if !node.OwnerAddress.Equals(ownerAddr) {
		return node.Status, types.ErrInvalidOwnerAddr
	}

	if opinion.Equal(types.Approve) {
		votePool.ApproveList = append(votePool.ApproveList, voterAddr)
	} else {
		votePool.RejectList = append(votePool.RejectList, voterAddr)
	}
	k.SetIndexingNodeRegistrationVotePool(ctx, votePool)

	if node.Status == sdk.Bonded {
		return node.Status, nil
	}

	totalSpCount := len(k.GetAllValidIndexingNodes(ctx))
	voteCountRequiredToPass := totalSpCount*2/3 + 1
	//unbounded to bounded
	if len(votePool.ApproveList) >= voteCountRequiredToPass {
		node.Status = sdk.Bonded
		k.SetIndexingNode(ctx, node)

		// move stake from not bonded pool to bonded pool
		tokenToBond := sdk.NewCoin(k.BondDenom(ctx), node.GetTokens())
		notBondedToken := k.GetIndexingNodeNotBondedToken(ctx)
		bondedToken := k.GetIndexingNodeBondedToken(ctx)

		if notBondedToken.IsLT(tokenToBond) {
			return node.Status, types.ErrInsufficientBalance
		}
		notBondedToken = notBondedToken.Sub(tokenToBond)
		bondedToken = bondedToken.Add(tokenToBond)
		k.SetIndexingNodeNotBondedToken(ctx, notBondedToken)
		k.SetIndexingNodeBondedToken(ctx, bondedToken)
	}

	return node.Status, nil
}

func (k Keeper) hasValue(items []sdk.AccAddress, item sdk.AccAddress) bool {
	for _, eachItem := range items {
		if eachItem.Equals(item) {
			return true
		}
	}
	return false
}

func (k Keeper) GetIndexingNodeRegistrationVotePool(ctx sdk.Context, nodeAddr sdk.AccAddress) (votePool types.IndexingNodeRegistrationVotePool, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetIndexingNodeRegistrationVotesKey(nodeAddr))
	if bz == nil {
		return votePool, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &votePool)
	return votePool, true
}

func (k Keeper) SetIndexingNodeRegistrationVotePool(ctx sdk.Context, votePool types.IndexingNodeRegistrationVotePool) {
	nodeAddr := votePool.NodeAddress
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(votePool)
	store.Set(types.GetIndexingNodeRegistrationVotesKey(nodeAddr), bz)
}

func (k Keeper) UpdateIndexingNode(ctx sdk.Context, networkID string, description types.Description,
	networkAddr sdk.AccAddress, ownerAddr sdk.AccAddress) error {

	node, found := k.GetIndexingNode(ctx, networkAddr)
	if !found {
		return types.ErrNoIndexingNodeFound
	}

	if !node.OwnerAddress.Equals(ownerAddr) {
		return types.ErrInvalidOwnerAddr
	}

	node.NetworkID = networkID
	node.Description = description

	k.SetIndexingNode(ctx, node)

	return nil
}

func (k Keeper) SetIndexingNodeBondedToken(ctx sdk.Context, token sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token)
	store.Set(types.IndexingNodeBondedTokenKey, bz)
}

func (k Keeper) GetIndexingNodeBondedToken(ctx sdk.Context) (token sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.IndexingNodeBondedTokenKey)
	if bz == nil {
		return sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &token)
	return token
}

func (k Keeper) SetIndexingNodeNotBondedToken(ctx sdk.Context, token sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token)
	store.Set(types.IndexingNodeNotBondedTokenKey, bz)
}

func (k Keeper) GetIndexingNodeNotBondedToken(ctx sdk.Context) (token sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.IndexingNodeNotBondedTokenKey)
	if bz == nil {
		return sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &token)
	return token
}

//
//
//// return a given amount of all the UnbondingIndexingNodes
//func (k Keeper) GetUINs(ctx sdk.Context, networkAddr sdk.AccAddress,
//	maxRetrieve uint16) (unbondingIndexingNodes []types.UnbondingIndexingNode) {
//
//	unbondingIndexingNodes = make([]types.UnbondingIndexingNode, maxRetrieve)
//
//	store := ctx.KVStore(k.storeKey)
//	indexingNodePrefixKey := types.GetUINKey(networkAddr)
//	iterator := sdk.KVStorePrefixIterator(store, indexingNodePrefixKey)
//	defer iterator.Close()
//
//	i := 0
//	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
//		unbondingIndexingNode := types.MustUnmarshalUIN(k.cdc, iterator.Value())
//		unbondingIndexingNodes[i] = unbondingIndexingNode
//		i++
//	}
//	return unbondingIndexingNodes[:i] // trim if the array length < maxRetrieve
//}
//
//// return a unbonding UnbondingIndexingNode
//func (k Keeper) GetUIN(ctx sdk.Context,
//	networkAddr sdk.AccAddress) (ubd types.UnbondingIndexingNode, found bool) {
//
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetUINKey(networkAddr)
//	value := store.Get(key)
//	if value == nil {
//		return ubd, false
//	}
//
//	ubd = types.MustUnmarshalUIN(k.cdc, value)
//	return ubd, true
//}
//
//// iterate through all of the unbonding indexingNodes
//func (k Keeper) IterateUINs(ctx sdk.Context, fn func(index int64, ubd types.UnbondingIndexingNode) (stop bool)) {
//	store := ctx.KVStore(k.storeKey)
//	iterator := sdk.KVStorePrefixIterator(store, types.UBDIndexingNodeKey)
//	defer iterator.Close()
//
//	for i := int64(0); iterator.Valid(); iterator.Next() {
//		ubd := types.MustUnmarshalUIN(k.cdc, iterator.Value())
//		if stop := fn(i, ubd); stop {
//			break
//		}
//		i++
//	}
//}
//
//// HasMaxUnbondingIndexingNodeEntries - check if unbonding IndexingNode has maximum number of entries
//func (k Keeper) HasMaxUINEntries(ctx sdk.Context, networkAddr sdk.AccAddress) bool {
//	ubd, found := k.GetUIN(ctx, networkAddr)
//	if !found {
//		return false
//	}
//	return len(ubd.Entries) >= int(k.MaxEntries(ctx))
//}
//
//// set the unbonding IndexingNode
//func (k Keeper) SetUIN(ctx sdk.Context, ubd types.UnbondingIndexingNode) {
//	store := ctx.KVStore(k.storeKey)
//	bz := types.MustMarshalUIN(k.cdc, ubd)
//	key := types.GetUINKey(ubd.GetNetworkAddr())
//	store.Set(key, bz)
//}
//
//// remove the unbonding IndexingNode object
//func (k Keeper) RemoveUIN(ctx sdk.Context, ubd types.UnbondingIndexingNode) {
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetUINKey(ubd.GetNetworkAddr())
//	store.Delete(key)
//}
//
//// SetUnbondingIndexingNodeEntry adds an entry to the unbonding IndexingNode at
//// the given addresses. It creates the unbonding IndexingNode if it does not exist
//func (k Keeper) SetUINEntry(ctx sdk.Context, networkAddr sdk.AccAddress,
//	creationHeight int64, minTime time.Time, balance sdk.Int) types.UnbondingIndexingNode {
//
//	ubd, found := k.GetUIN(ctx, networkAddr)
//	if found {
//		ubd.AddEntry(creationHeight, minTime, balance)
//	} else {
//		ubd = types.NewUnbondingIndexingNode(networkAddr, creationHeight, minTime, balance)
//	}
//	k.SetUIN(ctx, ubd)
//	return ubd
//}
//
//// unbonding delegation queue timeslice operations
//
//// gets a specific unbonding queue timeslice. A timeslice is a slice of DVPairs
//// corresponding to unbonding delegations that expire at a certain time.
//func (k Keeper) GetUINQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (networkAddrs []sdk.AccAddress) {
//	store := ctx.KVStore(k.storeKey)
//	bz := store.Get(types.GetUINTimeKey(timestamp))
//	if bz == nil {
//		return []sdk.AccAddress{}
//	}
//	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &networkAddrs)
//	return networkAddrs
//}
//
//// Sets a specific unbonding queue timeslice.
//func (k Keeper) SetUINQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []sdk.AccAddress) {
//	store := ctx.KVStore(k.storeKey)
//	bz := k.cdc.MustMarshalBinaryLengthPrefixed(keys)
//	store.Set(types.GetUINTimeKey(timestamp), bz)
//}
//
//// Insert an unbonding delegation to the appropriate timeslice in the unbonding queue
//func (k Keeper) InsertUINQueue(ctx sdk.Context, ubd types.UnbondingIndexingNode,
//	completionTime time.Time) {
//
//	timeSlice := k.GetUINQueueTimeSlice(ctx, completionTime)
//	networkAddr := ubd.NetworkAddr
//	if len(timeSlice) == 0 {
//		k.SetUINQueueTimeSlice(ctx, completionTime, []sdk.AccAddress{networkAddr})
//	} else {
//		timeSlice = append(timeSlice, networkAddr)
//		k.SetUINQueueTimeSlice(ctx, completionTime, timeSlice)
//	}
//}
//
//// Returns all the unbonding queue timeslices from time 0 until endTime
//func (k Keeper) UINQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
//	store := ctx.KVStore(k.storeKey)
//	return store.Iterator(types.UBDIndexingNodeQueueKey,
//		sdk.InclusiveEndBytes(types.GetUINTimeKey(endTime)))
//}
//
//// Returns a concatenated list of all the timeslices inclusively previous to
//// currTime, and deletes the timeslices from the queue
//func (k Keeper) DequeueAllMatureUINQueue(ctx sdk.Context,
//	currTime time.Time) (matureUnbonds []sdk.AccAddress) {
//
//	store := ctx.KVStore(k.storeKey)
//	// gets an iterator for all timeslices from time 0 until the current Blockheader time
//	unbondingTimesliceIterator := k.UINQueueIterator(ctx, ctx.BlockHeader().Time)
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

//func (k Keeper) DoRemoveIndexingNode(
//	ctx sdk.Context, indexingNode register.IndexingNode, amt sdk.Int,
//) (time.Time, error) {
//
//	ownerAcc := k.accountKeeper.GetAccount(ctx, indexingNode.OwnerAddress)
//	if ownerAcc == nil {
//		return time.Time{}, types.ErrNoOwnerAccountFound
//	}
//
//	networkAddr := indexingNode.GetNetworkAddr()
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
//	if indexingNode.GetStatus() == sdk.Bonded {
//		k.bondedToUnbonding(ctx, indexingNode, false)
//	}
//
//	params := k.GetParams(ctx)
//	// set the unbonding mature time and completion height appropriately
//	unbondingMatureTime := calcUnbondingMatureTime(indexingNode.CreationTime, params.UnbondingThreasholdTime, params.UnbondingCompletionTime)
//	unbondingNode := types.NewUnbondingNode(indexingNode.GetNetworkAddr(), false, ctx.BlockHeight(), unbondingMatureTime, returnAmount)
//	// Adds to unbonding node queue
//	k.InsertUnbondingNodeQueue(ctx, unbondingNode, unbondingMatureTime)
//
//	return unbondingMatureTime, nil
//}
