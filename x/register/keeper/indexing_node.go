package keeper

import (
	"bytes"
	"strings"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

const (
	indexingNodeCacheSize        = 500
	votingValidityPeriodInSecond = 7 * 24 * 60 * 60 // 7 days
)

// Cache the amino decoding of indexing nodes, as it can be the case that repeated slashing calls
// cause many calls to getIndexingNode, which were shown to throttle the state machine in our
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

// getIndexingNode get a single indexing node
func (k Keeper) GetIndexingNode(ctx sdk.Context, p2pAddress stratos.SdsAddress) (indexingNode types.IndexingNode, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetIndexingNodeKey(p2pAddress))
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
	networkAddr, _ := stratos.SdsAddressFromBech32(indexingNode.GetNetworkAddr())
	store.Set(types.GetIndexingNodeKey(networkAddr), bz)
}

// GetAllIndexingNodes get the set of all indexing nodes with no limits, used during genesis dump
func (k Keeper) GetAllIndexingNodes(ctx sdk.Context) (indexingNodes *types.IndexingNodes) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.IndexingNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalIndexingNode(k.cdc, iterator.Value())
		indexingNodes.IndexingNodes = append(indexingNodes.IndexingNodes, &node)
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
		if !node.IsSuspended() && node.GetStatus().Equal(stakingtypes.Bonded) {
			indexingNodes = append(indexingNodes, node)
		}
	}
	return indexingNodes
}

func (k Keeper) RegisterIndexingNode(ctx sdk.Context, networkAddr stratos.SdsAddress, pubKey cryptotypes.PubKey, ownerAddr sdk.AccAddress,
	description types.Description, stake sdk.Coin) (ozoneLimitChange sdk.Int, err error) {

	indexingNode, err := types.NewIndexingNode(networkAddr, pubKey, ownerAddr, &description, ctx.BlockHeader().Time)
	if err != nil {
		return ozoneLimitChange, err
	}
	ozoneLimitChange, err = k.AddIndexingNodeStake(ctx, indexingNode, stake)
	if err != nil {
		return ozoneLimitChange, err
	}

	var approveList = make([]stratos.SdsAddress, 0)
	var rejectList = make([]stratos.SdsAddress, 0)
	votingValidityPeriod := votingValidityPeriodInSecond * time.Second
	expireTime := ctx.BlockHeader().Time.Add(votingValidityPeriod)

	votePool := types.NewRegistrationVotePool(networkAddr, approveList, rejectList, expireTime)
	k.SetIndexingNodeRegistrationVotePool(ctx, votePool)

	return ozoneLimitChange, nil
}

// AddIndexingNodeStake Update the tokens of an existing indexing node
func (k Keeper) AddIndexingNodeStake(ctx sdk.Context, indexingNode types.IndexingNode, tokenToAdd sdk.Coin,
) (ozoneLimitChange sdk.Int, err error) {

	coins := sdk.NewCoins(tokenToAdd)

	ownerAddr, err := sdk.AccAddressFromBech32(indexingNode.GetOwnerAddress())
	if err != nil {
		return sdk.ZeroInt(), types.ErrInvalidOwnerAddr
	}
	// sub coins from owner's wallet
	hasCoin := k.bankKeeper..HasCoins(ctx, ownerAddr, coins)
	if !hasCoin {
		return sdk.ZeroInt(), types.ErrInsufficientBalance
	}
	_, err = k.bankKeeper.SubtractCoins(ctx, ownerAddr, coins)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	indexingNode = indexingNode.AddToken(tokenToAdd.Amount)

	switch indexingNode.GetStatus() {
	case stakingtypes.Unbonded:
		notBondedTokenInPool := k.GetIndexingNodeNotBondedToken(ctx)
		notBondedTokenInPool = notBondedTokenInPool.Add(tokenToAdd)
		k.SetIndexingNodeNotBondedToken(ctx, notBondedTokenInPool)
	case stakingtypes.Bonded:
		bondedTokenInPool := k.GetIndexingNodeBondedToken(ctx)
		bondedTokenInPool = bondedTokenInPool.Add(tokenToAdd)
		k.SetIndexingNodeBondedToken(ctx, bondedTokenInPool)
	case stakingtypes.Unbonding:
		return sdk.ZeroInt(), types.ErrUnbondingNode
	}

	k.SetIndexingNode(ctx, indexingNode)
	ozoneLimitChange = k.increaseOzoneLimitByAddStake(ctx, tokenToAdd.Amount)

	return ozoneLimitChange, nil
}

func (k Keeper) RemoveTokenFromPoolWhileUnbondingIndexingNode(ctx sdk.Context, indexingNode types.IndexingNode, tokenToSub sdk.Coin) error {
	// get pools
	bondedTokenInPool := k.GetIndexingNodeBondedToken(ctx)
	notBondedTokenInPool := k.GetIndexingNodeNotBondedToken(ctx)
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
	networkAddr, err := stratos.SdsAddressFromBech32(indexingNode.GetNetworkAddr())
	if err != nil {
		return types.ErrInvalidNetworkAddr
	}
	ownerAddr, err := sdk.AccAddressFromBech32(indexingNode.GetOwnerAddress())
	if err != nil {
		return types.ErrInvalidOwnerAddr
	}
	ownerAcc := k.accountKeeper.GetAccount(ctx, ownerAddr)
	if ownerAcc == nil {
		return types.ErrNoOwnerAccountFound
	}

	coins := sdk.NewCoins(tokenToSub)

	if indexingNode.Tokens.LT(tokenToSub.Amount) {
		return types.ErrInsufficientBalance
	}

	// deduct tokens from NotBondedPool
	notBondedTokenInPool := k.GetIndexingNodeNotBondedToken(ctx)
	if notBondedTokenInPool.IsLT(tokenToSub) {
		return types.ErrInsufficientBalanceOfNotBondedPool
	}
	notBondedTokenInPool = notBondedTokenInPool.Sub(tokenToSub)
	k.SetIndexingNodeNotBondedToken(ctx, notBondedTokenInPool)

	// deduct slashing amount first
	coins = k.DeductSlashing(ctx, ownerAddr, coins)
	// add tokens to owner acc
	_, err := k.bankKeeper.AddCoins(ctx, indexingNode.OwnerAddress, coins)
	if err != nil {
		return err
	}

	indexingNode = indexingNode.SubToken(tokenToSub.Amount)
	newStake := indexingNode.Tokens

	k.SetIndexingNode(ctx, indexingNode)

	if newStake.IsZero() {
		err = k.removeIndexingNode(ctx, networkAddr)
		if err != nil {
			return err
		}
	}
	return nil
}

// remove the indexing node record and associated indexes
func (k Keeper) removeIndexingNode(ctx sdk.Context, addr stratos.SdsAddress) error {
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

// getIndexingNodeList get all indexing nodes by networkAddr
func (k Keeper) GetIndexingNodeList(ctx sdk.Context, networkAddr stratos.SdsAddress) (indexingNodes []types.IndexingNode, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.IndexingNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalIndexingNode(k.cdc, iterator.Value())
		networkAddrNode, err := stratos.SdsAddressFromBech32(node.GetNetworkAddr())
		if err != nil {
			continue
		}
		if bytes.Equal(networkAddrNode, networkAddr) {
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

func (k Keeper) HandleVoteForIndexingNodeRegistration(ctx sdk.Context, nodeAddr stratos.SdsAddress, ownerAddr sdk.AccAddress,
	opinion types.VoteOpinion, voterAddr stratos.SdsAddress) (nodeStatus stakingtypes.BondStatus, err error) {

	votePool, found := k.GetIndexingNodeRegistrationVotePool(ctx, nodeAddr)
	if !found {
		return stakingtypes.Unbonded, types.ErrNoRegistrationVotePoolFound
	}
	if votePool.ExpireTime.Before(ctx.BlockHeader().Time) {
		return stakingtypes.Unbonded, types.ErrVoteExpired
	}
	if hasStringValue(votePool.ApproveList, voterAddr.String()) || hasStringValue(votePool.RejectList, voterAddr.String()) {
		return stakingtypes.Unbonded, types.ErrDuplicateVoting
	}

	node, found := k.GetIndexingNode(ctx, nodeAddr)
	if !found {
		return stakingtypes.Unbonded, types.ErrNoIndexingNodeFound
	}
	ownerAddrNode, err := sdk.AccAddressFromBech32(node.GetOwnerAddress())
	if err != nil {
		return stakingtypes.Unbonded, types.ErrInvalidOwnerAddr
	}
	if !bytes.Equal(ownerAddrNode, ownerAddr) {
		return node.Status, types.ErrInvalidOwnerAddr
	}

	if opinion.Equal(types.Approve) {
		votePool.ApproveList = append(votePool.ApproveList, voterAddr.String())
	} else {
		votePool.RejectList = append(votePool.RejectList, voterAddr.String())
	}
	k.SetIndexingNodeRegistrationVotePool(ctx, votePool)

	if node.Status == stakingtypes.Bonded {
		return node.Status, nil
	}

	totalSpCount := len(k.GetAllValidIndexingNodes(ctx))
	voteCountRequiredToPass := totalSpCount*2/3 + 1
	//unbounded to bounded
	if len(votePool.ApproveList) >= voteCountRequiredToPass {
		node.Status = stakingtypes.Bonded
		node.Suspend = false
		k.SetIndexingNode(ctx, node)

		// move stake from not bonded pool to bonded pool
		tokenToBond := sdk.NewCoin(k.BondDenom(ctx), node.Tokens)
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

func (k Keeper) GetIndexingNodeRegistrationVotePool(ctx sdk.Context, nodeAddr stratos.SdsAddress) (votePool types.IndexingNodeRegistrationVotePool, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetIndexingNodeRegistrationVotesKey(nodeAddr))
	if bz == nil {
		return votePool, false
	}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &votePool)
	return votePool, true
}

func (k Keeper) SetIndexingNodeRegistrationVotePool(ctx sdk.Context, votePool types.IndexingNodeRegistrationVotePool) {
	nodeAddr := votePool.NodeAddress
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&votePool)
	node, _ := stratos.SdsAddressFromBech32(nodeAddr)
	store.Set(types.GetIndexingNodeRegistrationVotesKey(node), bz)
}

func (k Keeper) UpdateIndexingNode(ctx sdk.Context, description types.Description,
	networkAddr stratos.SdsAddress, ownerAddr sdk.AccAddress) error {

	node, found := k.GetIndexingNode(ctx, networkAddr)
	if !found {
		return types.ErrNoIndexingNodeFound
	}

	ownerAddrNode, _ := sdk.AccAddressFromBech32(node.GetOwnerAddress())
	if !bytes.Equal(ownerAddrNode, ownerAddr) {
		return types.ErrInvalidOwnerAddr
	}

	node.Description = &description

	k.SetIndexingNode(ctx, node)

	return nil
}

func (k Keeper) UpdateIndexingNodeStake(ctx sdk.Context, networkAddr stratos.SdsAddress, ownerAddr sdk.AccAddress,
	stakeDelta sdk.Coin, incrStake bool) (ozoneLimitChange sdk.Int, unbondingMatureTime time.Time, err error) {

	blockTime := ctx.BlockHeader().Time
	node, found := k.GetIndexingNode(ctx, networkAddr)
	if !found {
		return sdk.ZeroInt(), blockTime, types.ErrNoIndexingNodeFound
	}

	ownerAddrNode, _ := sdk.AccAddressFromBech32(node.GetOwnerAddress())
	if !bytes.Equal(ownerAddrNode, ownerAddr) {
		return sdk.ZeroInt(), blockTime, types.ErrInvalidOwnerAddr
	}

	if incrStake {
		ozoneLimitChange, err = k.AddIndexingNodeStake(ctx, node, stakeDelta)
		if err != nil {
			return sdk.ZeroInt(), blockTime, err
		}
		return ozoneLimitChange, blockTime, nil
	} else {
		// if !incrStake
		if node.GetStatus() == stakingtypes.Unbonding {
			return sdk.ZeroInt(), blockTime, types.ErrUnbondingNode
		}
		ozoneLimitChange, completionTime, err := k.UnbondIndexingNode(ctx, node, stakeDelta.Amount)
		if err != nil {
			return sdk.ZeroInt(), blockTime, err
		}
		return ozoneLimitChange, completionTime, nil
	}
}

func (k Keeper) SetIndexingNodeBondedToken(ctx sdk.Context, token sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&token)
	store.Set(types.IndexingNodeBondedTokenKey, bz)
}

func (k Keeper) GetIndexingNodeBondedToken(ctx sdk.Context) (token sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.IndexingNodeBondedTokenKey)
	if bz == nil {
		return sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())
	}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &token)
	return token
}

func (k Keeper) SetIndexingNodeNotBondedToken(ctx sdk.Context, token sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&token)
	store.Set(types.IndexingNodeNotBondedTokenKey, bz)
}

func (k Keeper) GetIndexingNodeNotBondedToken(ctx sdk.Context) (token sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.IndexingNodeNotBondedTokenKey)
	if bz == nil {
		return sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())
	}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &token)
	return token
}
