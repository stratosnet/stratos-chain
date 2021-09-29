package keeper

import (
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
	description types.Description, stake sdk.Coin) (ozoneLimitChange sdk.Int, err error) {

	indexingNode := types.NewIndexingNode(networkID, pubKey, ownerAddr, description, ctx.BlockHeader().Time)

	ozoneLimitChange, err = k.AddIndexingNodeStake(ctx, indexingNode, stake)
	if err != nil {
		return ozoneLimitChange, err
	}

	var approveList = make([]sdk.AccAddress, 0)
	var rejectList = make([]sdk.AccAddress, 0)
	votingValidityPeriod := votingValidityPeriodInSecond * time.Second
	expireTime := ctx.BlockHeader().Time.Add(votingValidityPeriod)

	votePool := types.NewRegistrationVotePool(indexingNode.GetNetworkAddr(), approveList, rejectList, expireTime)
	k.SetIndexingNodeRegistrationVotePool(ctx, votePool)

	return ozoneLimitChange, nil
}

// AddIndexingNodeStake Update the tokens of an existing indexing node
func (k Keeper) AddIndexingNodeStake(ctx sdk.Context, indexingNode types.IndexingNode, tokenToAdd sdk.Coin,
) (ozoneLimitChange sdk.Int, err error) {

	nodeAcc := k.accountKeeper.GetAccount(ctx, indexingNode.GetNetworkAddr())
	if nodeAcc == nil {
		nodeAcc = k.accountKeeper.NewAccountWithAddress(ctx, indexingNode.GetNetworkAddr())
		k.accountKeeper.SetAccount(ctx, nodeAcc)
	}

	coins := sdk.NewCoins(tokenToAdd)
	hasCoin := k.bankKeeper.HasCoins(ctx, indexingNode.OwnerAddress, coins)
	if !hasCoin {
		return sdk.ZeroInt(), types.ErrInsufficientBalance
	}

	_, err = k.bankKeeper.SubtractCoins(ctx, indexingNode.GetOwnerAddr(), coins)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	indexingNode = indexingNode.AddToken(tokenToAdd.Amount)

	switch indexingNode.GetStatus() {
	case sdk.Unbonded:
		notBondedTokenInPool := k.GetIndexingNodeNotBondedToken(ctx)
		notBondedTokenInPool = notBondedTokenInPool.Add(tokenToAdd)
		k.SetIndexingNodeNotBondedToken(ctx, notBondedTokenInPool)
	case sdk.Bonded:
		bondedTokenInPool := k.GetIndexingNodeBondedToken(ctx)
		bondedTokenInPool = bondedTokenInPool.Add(tokenToAdd)
		k.SetIndexingNodeBondedToken(ctx, bondedTokenInPool)
	case sdk.Unbonding:
		return sdk.ZeroInt(), types.ErrUnbondingNode
	}

	newStake := indexingNode.GetTokens()

	k.SetIndexingNode(ctx, indexingNode)
	k.SetLastIndexingNodeStake(ctx, indexingNode.GetNetworkAddr(), newStake)
	ozoneLimitChange = k.increaseOzoneLimitByAddStake(ctx, tokenToAdd.Amount)

	return ozoneLimitChange, nil
}

func (k Keeper) RemoveTokenFromPoolWhileUnbondingIndexingNode(ctx sdk.Context, indexingNode types.IndexingNode, tokenToSub sdk.Coin) error {
	// change node status to unbonding
	indexingNode.Status = sdk.Unbonding
	k.SetIndexingNode(ctx, indexingNode)
	// get pools
	bondedTokenInPool := k.GetResourceNodeBondedToken(ctx)
	notBondedTokenInPool := k.GetResourceNodeNotBondedToken(ctx)
	if bondedTokenInPool.IsLT(tokenToSub) {
		return types.ErrInsufficientBalanceOfBondedPool
	}
	// remove token from BondedPool
	bondedTokenInPool = bondedTokenInPool.Sub(tokenToSub)
	k.SetIndexingNodeBondedToken(ctx, bondedTokenInPool)
	// add token into NotBondedPool
	notBondedTokenInPool = notBondedTokenInPool.Add(tokenToSub)
	k.SetIndexingNodeNotBondedToken(ctx, notBondedTokenInPool)
	return nil
}

// SubtractIndexingNodeStake Update the tokens of an existing indexing node
func (k Keeper) SubtractIndexingNodeStake(ctx sdk.Context, indexingNode types.IndexingNode, tokenToSub sdk.Coin) error {
	ownerAcc := k.accountKeeper.GetAccount(ctx, indexingNode.OwnerAddress)
	if ownerAcc == nil {
		return types.ErrNoOwnerAccountFound
	}

	coins := sdk.NewCoins(tokenToSub)

	if indexingNode.GetStatus() == sdk.Unbonded || indexingNode.GetStatus() == sdk.Unbonding {
		notBondedTokenInPool := k.GetIndexingNodeNotBondedToken(ctx)
		if notBondedTokenInPool.IsLT(tokenToSub) {
			return types.ErrInsufficientBalanceOfNotBondedPool
		}
		notBondedTokenInPool = notBondedTokenInPool.Sub(tokenToSub)
		k.SetIndexingNodeNotBondedToken(ctx, notBondedTokenInPool)
	}
	if indexingNode.GetStatus() == sdk.Bonded {
		return types.ErrInvalidNodeStatBonded
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
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalIndexingNode(k.cdc, iterator.Value())
		if strings.Compare(node.NetworkID, networkID) == 0 {
			indexingNodes = append(indexingNodes, node)
		}
	}
	return indexingNodes, nil
}

func (k Keeper) GetIndexingNodeListByMoniker(ctx sdk.Context, moniker string) (resourceNodes []types.IndexingNode, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.IndexingNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalIndexingNode(k.cdc, iterator.Value())
		if strings.Compare(node.Description.Moniker, moniker) == 0 {
			resourceNodes = append(resourceNodes, node)
		}
	}
	return resourceNodes, nil
}

func (k Keeper) HandleVoteForIndexingNodeRegistration(ctx sdk.Context, nodeAddr sdk.AccAddress, ownerAddr sdk.AccAddress,
	opinion types.VoteOpinion, voterAddr sdk.AccAddress) (nodeStatus sdk.BondStatus, err error) {

	votePool, found := k.GetIndexingNodeRegistrationVotePool(ctx, nodeAddr)
	if !found {
		return sdk.Unbonded, types.ErrNoRegistrationVotePoolFound
	}
	if votePool.ExpireTime.Before(ctx.BlockHeader().Time) {
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

func (k Keeper) GetNodeOwnerMapFromIndexingNodes(ctx sdk.Context, nodeOwnerMap map[string]sdk.AccAddress) map[string]sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.IndexingNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalIndexingNode(k.cdc, iterator.Value())
		nodeOwnerMap[node.GetNetworkAddr().String()] = node.OwnerAddress
	}
	return nodeOwnerMap
}
