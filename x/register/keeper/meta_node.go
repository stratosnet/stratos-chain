package keeper

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// getMetaNode get a single meta node
func (k Keeper) GetMetaNode(ctx sdk.Context, p2pAddress stratos.SdsAddress) (metaNode types.MetaNode, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetMetaNodeKey(p2pAddress))
	if value == nil {
		return metaNode, false
	}
	metaNode = types.MustUnmarshalMetaNode(k.cdc, value)
	return metaNode, true
}

// SetMetaNode sets the main record holding meta node details
func (k Keeper) SetMetaNode(ctx sdk.Context, metaNode types.MetaNode) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalMetaNode(k.cdc, metaNode)
	networkAddr, _ := stratos.SdsAddressFromBech32(metaNode.GetNetworkAddress())
	store.Set(types.GetMetaNodeKey(networkAddr), bz)
}

// GetAllMetaNodes get the set of all meta nodes with no limits, used during genesis dump
// Iteration for all meta nodes
func (k Keeper) GetAllMetaNodes(ctx sdk.Context) (metaNodes types.MetaNodes) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.MetaNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalMetaNode(k.cdc, iterator.Value())
		metaNodes = append(metaNodes, node)
	}
	return metaNodes
}

// GetAllValidMetaNodes get the set of all bonded & not suspended meta nodes
func (k Keeper) GetAllValidMetaNodes(ctx sdk.Context) (metaNodes []types.MetaNode) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.MetaNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalMetaNode(k.cdc, iterator.Value())
		if !node.GetSuspend() && node.GetStatus() == stakingtypes.Bonded {
			metaNodes = append(metaNodes, node)
		}
	}
	return metaNodes
}

func (k Keeper) RegisterMetaNode(ctx sdk.Context, networkAddr stratos.SdsAddress, pubKey cryptotypes.PubKey, ownerAddr sdk.AccAddress,
	description types.Description, deposit sdk.Coin) (ozoneLimitChange sdk.Int, err error) {

	if _, found := k.GetMetaNode(ctx, networkAddr); found {
		ctx.Logger().Error("Meta node already exist")
		return ozoneLimitChange, types.ErrMetaNodePubKeyExists
	}
	if _, found := k.GetResourceNode(ctx, networkAddr); found {
		ctx.Logger().Error("Resource node with same network address already exist")
		return ozoneLimitChange, types.ErrResourceNodePubKeyExists
	}

	if deposit.GetDenom() != k.BondDenom(ctx) {
		return ozoneLimitChange, types.ErrBadDenom
	}

	metaNode, err := types.NewMetaNode(networkAddr, pubKey, ownerAddr, description, ctx.BlockHeader().Time)
	if err != nil {
		return ozoneLimitChange, err
	}
	ozoneLimitChange, err = k.AddMetaNodeDeposit(ctx, metaNode, deposit)
	if err != nil {
		return ozoneLimitChange, err
	}

	var approveList = make([]stratos.SdsAddress, 0)
	var rejectList = make([]stratos.SdsAddress, 0)
	votingValidityPeriod := k.VotingPeriod(ctx)
	expireTime := ctx.BlockHeader().Time.Add(votingValidityPeriod)

	votePool := types.NewRegistrationVotePool(networkAddr, approveList, rejectList, expireTime)
	k.SetMetaNodeRegistrationVotePool(ctx, votePool)

	return ozoneLimitChange, nil
}

// AddMetaNodeDeposit Update the tokens of an existing meta node
func (k Keeper) AddMetaNodeDeposit(ctx sdk.Context, metaNode types.MetaNode, tokenToAdd sdk.Coin,
) (ozoneLimitChange sdk.Int, err error) {

	coins := sdk.NewCoins(tokenToAdd)

	ownerAddr, err := sdk.AccAddressFromBech32(metaNode.GetOwnerAddress())
	if err != nil {
		return sdk.ZeroInt(), types.ErrInvalidOwnerAddr
	}
	// sub coins from owner's wallet
	hasCoin := k.bankKeeper.HasBalance(ctx, ownerAddr, tokenToAdd)
	if !hasCoin {
		return sdk.ZeroInt(), types.ErrInsufficientBalance
	}
	targetModuleAccName := ""

	switch metaNode.GetStatus() {
	case stakingtypes.Unbonded:
		targetModuleAccName = types.MetaNodeNotBondedPool
	case stakingtypes.Bonded:
		targetModuleAccName = types.MetaNodeBondedPool
	case stakingtypes.Unbonding:
		return sdk.ZeroInt(), types.ErrUnbondingNode
	}

	if len(targetModuleAccName) > 0 {
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, ownerAddr, targetModuleAccName, coins)
		if err != nil {
			return sdk.ZeroInt(), err
		}
	}

	metaNode = metaNode.AddToken(tokenToAdd.Amount)
	k.SetMetaNode(ctx, metaNode)

	if !metaNode.Suspend {
		ozoneLimitChange = k.IncreaseOzoneLimitByAddDeposit(ctx, tokenToAdd.Amount)
	} else {
		// if node is currently suspended, ozone limit will be increased upon unsuspension instead of NOW
		ozoneLimitChange = sdk.ZeroInt()
	}

	return ozoneLimitChange, nil
}

// TODO: Unused parameter: metaNode
func (k Keeper) RemoveTokenFromPoolWhileUnbondingMetaNode(ctx sdk.Context, metaNode types.MetaNode, tokenToSub sdk.Coin) error {
	bondedMetaAccountAddr := k.accountKeeper.GetModuleAddress(types.MetaNodeBondedPool)
	if bondedMetaAccountAddr == nil {
		ctx.Logger().Error("bonded pool account address for meta nodes does not exist.")
		return types.ErrUnknownAccountAddress
	}

	hasCoin := k.bankKeeper.HasBalance(ctx, bondedMetaAccountAddr, tokenToSub)
	if !hasCoin {
		return types.ErrInsufficientBalance
	}

	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.MetaNodeBondedPool, types.MetaNodeNotBondedPool, sdk.NewCoins(tokenToSub))
	if err != nil {
		return types.ErrInsufficientBalance
	}
	return nil
}

// SubtractMetaNodeDeposit Update the tokens of an existing meta node
func (k Keeper) SubtractMetaNodeDeposit(ctx sdk.Context, metaNode types.MetaNode, tokenToSub sdk.Coin) error {
	networkAddr, err := stratos.SdsAddressFromBech32(metaNode.GetNetworkAddress())
	if err != nil {
		return types.ErrInvalidNetworkAddr
	}
	ownerAddr, err := sdk.AccAddressFromBech32(metaNode.GetOwnerAddress())
	if err != nil {
		return types.ErrInvalidOwnerAddr
	}
	ownerAcc := k.accountKeeper.GetAccount(ctx, ownerAddr)
	if ownerAcc == nil {
		return types.ErrNoOwnerAccountFound
	}

	coins := sdk.NewCoins(tokenToSub)

	if metaNode.Tokens.LT(tokenToSub.Amount) {
		return types.ErrInsufficientBalance
	}

	// deduct tokens from NotBondedPool
	nBondedMetaAccountAddr := k.accountKeeper.GetModuleAddress(types.MetaNodeNotBondedPool)
	if nBondedMetaAccountAddr == nil {
		ctx.Logger().Error("not bonded account address for meta nodes does not exist.")
		return types.ErrUnknownAccountAddress
	}

	hasCoin := k.bankKeeper.HasBalance(ctx, nBondedMetaAccountAddr, tokenToSub)
	if !hasCoin {
		return types.ErrInsufficientBalanceOfNotBondedPool
	}

	// deduct slashing amount first, slashed amt goes into TotalSlashedPool
	remaining, slashed := k.DeductSlashing(ctx, ownerAddr, coins, k.BondDenom(ctx))
	if !remaining.IsZero() {
		// add remaining tokens to owner acc
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.MetaNodeNotBondedPool, ownerAddr, remaining)
		if err != nil {
			return err
		}
	}
	if !slashed.IsZero() {
		// slashed token send to community_pool
		metaNodeNotBondedPoolAddr := k.accountKeeper.GetModuleAddress(types.MetaNodeNotBondedPool)
		err = k.distrKeeper.FundCommunityPool(ctx, slashed, metaNodeNotBondedPoolAddr)
		if err != nil {
			return err
		}
	}

	metaNode = metaNode.SubToken(tokenToSub.Amount)
	newDeposit := metaNode.Tokens

	k.SetMetaNode(ctx, metaNode)

	if newDeposit.IsZero() {
		err = k.removeMetaNode(ctx, networkAddr)
		if err != nil {
			return err
		}
	}
	return nil
}

// remove the meta node record and associated indexes
func (k Keeper) removeMetaNode(ctx sdk.Context, addr stratos.SdsAddress) error {
	// first retrieve the old meta node record
	metaNode, found := k.GetMetaNode(ctx, addr)
	if !found {
		return types.ErrNoMetaNodeFound
	}

	if metaNode.Tokens.IsPositive() {
		panic("attempting to remove a meta node which still contains tokens")
	}

	// delete the old meta node record
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetMetaNodeKey(addr))
	return nil
}

// getMetaNodeList get all meta nodes by networkAddr
func (k Keeper) GetMetaNodeList(ctx sdk.Context, networkAddr stratos.SdsAddress) (metaNodes []types.MetaNode, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.MetaNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalMetaNode(k.cdc, iterator.Value())
		networkAddrNode, err := stratos.SdsAddressFromBech32(node.GetNetworkAddress())
		if err != nil {
			continue
		}
		if networkAddrNode.Equals(networkAddr) {
			metaNodes = append(metaNodes, node)
		}
	}
	return metaNodes, nil
}

func (k Keeper) GetMetaNodeListByMoniker(ctx sdk.Context, moniker string) (resourceNodes []types.MetaNode, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.MetaNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalMetaNode(k.cdc, iterator.Value())
		if strings.Compare(node.Description.Moniker, moniker) == 0 {
			resourceNodes = append(resourceNodes, node)
		}
	}
	return resourceNodes, nil
}

func (k Keeper) WithdrawMetaNodeRegistrationDeposit(ctx sdk.Context, networkAddr stratos.SdsAddress, ownerAddr sdk.AccAddress) (
	unbondingMatureTime time.Time, err error) {

	node, found := k.GetMetaNode(ctx, networkAddr)
	if !found {
		return time.Time{}, types.ErrNoMetaNodeFound
	}
	if node.GetOwnerAddress() != ownerAddr.String() {
		return time.Time{}, types.ErrInvalidOwnerAddr
	}
	votePool, exist := k.GetMetaNodeRegistrationVotePool(ctx, networkAddr)
	if !exist {
		return time.Time{}, types.ErrNoRegistrationVotePoolFound
	}
	// to be qualified to withdraw, meta node must be unbonded && suspended && of non-passed vote
	if node.Status != stakingtypes.Unbonded || !node.Suspend || votePool.IsVotePassed {
		return time.Time{}, types.ErrInvalidNodeStat
	}
	// check available_deposit (node_token - unbonding_token > amt_to_withdraw)
	unbondingDeposit := k.GetUnbondingNodeBalance(ctx, networkAddr)
	availableDeposit := node.Tokens.Sub(unbondingDeposit)
	if availableDeposit.LTE(sdk.ZeroInt()) {
		return time.Time{}, types.ErrInsufficientBalance
	}
	if k.HasMaxUnbondingNodeEntries(ctx, networkAddr) {
		return time.Time{}, types.ErrMaxUnbondingNodeEntries
	}
	unbondingMatureTime = calcUnbondingMatureTime(ctx, node.Status, node.CreationTime, k.UnbondingThreasholdTime(ctx), k.UnbondingCompletionTime(ctx))

	// Set the unbonding mature time and completion height appropriately
	unbondingNode := k.SetUnbondingNodeEntry(ctx, networkAddr, true, ctx.BlockHeight(), unbondingMatureTime, availableDeposit)
	// Add to unbonding node queue
	k.InsertUnbondingNodeQueue(ctx, unbondingNode, unbondingMatureTime)
	ctx.Logger().Info("Unbonding meta node " + unbondingNode.String() + "\n after mature time" + unbondingMatureTime.String())

	// all deposit is being unbonded, update status to Unbonding
	node.Status = stakingtypes.Unbonding
	k.SetMetaNode(ctx, node)

	// remove from cache just in case
	k.RemoveMetaNodeFromBitMapIdxCache(networkAddr)
	// delete vote pool after withdraw is done
	k.DeleteMetaNodeRegistrationVotePool(ctx, networkAddr)

	return unbondingMatureTime, nil
}

func (k Keeper) HandleVoteForMetaNodeRegistration(ctx sdk.Context, candidateNetworkAddr stratos.SdsAddress, candidateOwnerAddr sdk.AccAddress,
	opinion types.VoteOpinion, voterNetworkAddr stratos.SdsAddress, voterOwnerAddr sdk.AccAddress) (nodeStatus stakingtypes.BondStatus, err error) {

	// voter validation
	voterNode, found := k.GetMetaNode(ctx, voterNetworkAddr)
	if !found {
		return stakingtypes.Unbonded, types.ErrNoVoterMetaNodeFound
	}
	if voterNode.GetOwnerAddress() != voterOwnerAddr.String() {
		return stakingtypes.Unbonded, types.ErrInvalidVoterOwnerAddr
	}
	if voterNode.Status != stakingtypes.Bonded || voterNode.Suspend {
		return stakingtypes.Unbonded, types.ErrInvalidVoterStatus
	}

	// candidate validation
	candidateNode, found := k.GetMetaNode(ctx, candidateNetworkAddr)
	if !found {
		return stakingtypes.Unbonded, types.ErrNoCandidateMetaNodeFound
	}
	if candidateNode.GetOwnerAddress() != candidateOwnerAddr.String() {
		return candidateNode.Status, types.ErrInvalidCandidateOwnerAddr
	}

	// vote validation and handle voting
	votePool, found := k.GetMetaNodeRegistrationVotePool(ctx, candidateNetworkAddr)
	if !found {
		return stakingtypes.Unbonded, types.ErrNoRegistrationVotePoolFound
	}
	if votePool.ExpireTime.Before(ctx.BlockHeader().Time) {
		return stakingtypes.Unbonded, types.ErrVoteExpired
	}
	if hasStringValue(votePool.ApproveList, voterNetworkAddr.String()) || hasStringValue(votePool.RejectList, voterNetworkAddr.String()) {
		return stakingtypes.Unbonded, types.ErrDuplicateVoting
	}

	if opinion.Equal(types.Approve) {
		votePool.ApproveList = append(votePool.ApproveList, voterNetworkAddr.String())
	} else {
		votePool.RejectList = append(votePool.RejectList, voterNetworkAddr.String())
	}
	k.SetMetaNodeRegistrationVotePool(ctx, votePool)

	// if vote had already passed before
	if votePool.IsVotePassed {
		return candidateNode.Status, nil
	}

	//if vote is yet to pass
	totalSpCount := len(k.GetAllValidMetaNodes(ctx))
	voteCountRequiredToPass := totalSpCount*2/3 + 1
	//unbounded to bounded
	if len(votePool.ApproveList) >= voteCountRequiredToPass {
		candidateNode.Status = stakingtypes.Bonded
		candidateNode.Suspend = false
		k.SetMetaNode(ctx, candidateNode)
		// add new available meta node to cache
		networkAddr, _ := stratos.SdsAddressFromBech32(candidateNode.GetNetworkAddress())
		k.AddMetaNodeToBitMapIdxCache(networkAddr)
		// increase ozone limit after vote is approved
		_ = k.IncreaseOzoneLimitByAddDeposit(ctx, candidateNode.Tokens)
		// increase mata node count
		v := k.GetBondedMetaNodeCnt(ctx)
		count := v.Add(sdk.NewInt(1))
		k.SetBondedMetaNodeCnt(ctx, count)
		// move deposit from not bonded pool to bonded pool
		tokenToBond := sdk.NewCoin(k.BondDenom(ctx), candidateNode.Tokens)
		// sub coins from not bonded pool
		nBondedMetaAccountAddr := k.accountKeeper.GetModuleAddress(types.MetaNodeNotBondedPool)
		if nBondedMetaAccountAddr == nil {
			ctx.Logger().Error("not bonded account address for meta nodes does not exist.")
			return candidateNode.Status, types.ErrUnknownAccountAddress
		}

		hasCoin := k.bankKeeper.HasBalance(ctx, nBondedMetaAccountAddr, tokenToBond)
		if !hasCoin {
			return candidateNode.Status, types.ErrInsufficientBalance
		}

		err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.MetaNodeNotBondedPool, types.MetaNodeBondedPool, sdk.NewCoins(tokenToBond))
		if err != nil {
			return candidateNode.Status, err
		}

		votePool.IsVotePassed = true
		k.SetMetaNodeRegistrationVotePool(ctx, votePool)
	}

	return candidateNode.Status, nil
}

func (k Keeper) UpdateMetaNode(ctx sdk.Context, description types.Description,
	networkAddr stratos.SdsAddress, ownerAddr sdk.AccAddress) error {

	node, found := k.GetMetaNode(ctx, networkAddr)
	if !found {
		return types.ErrNoMetaNodeFound
	}

	ownerAddrNode, _ := sdk.AccAddressFromBech32(node.GetOwnerAddress())
	if !ownerAddrNode.Equals(ownerAddr) {
		return types.ErrInvalidOwnerAddr
	}

	node.Description = description

	k.SetMetaNode(ctx, node)

	return nil
}

func (k Keeper) UpdateMetaNodeDeposit(ctx sdk.Context, networkAddr stratos.SdsAddress, ownerAddr sdk.AccAddress, depositDelta sdk.Coin, incrDeposit bool) (
	ozoneLimitChange sdk.Int, unbondingMatureTime time.Time, err error) {

	if depositDelta.GetDenom() != k.BondDenom(ctx) {
		return sdk.ZeroInt(), time.Time{}, types.ErrBadDenom
	}

	node, found := k.GetMetaNode(ctx, networkAddr)
	if !found {
		return sdk.ZeroInt(), time.Time{}, types.ErrNoMetaNodeFound
	}

	ownerAddrNode, _ := sdk.AccAddressFromBech32(node.GetOwnerAddress())
	if !ownerAddrNode.Equals(ownerAddr) {
		return sdk.ZeroInt(), time.Time{}, types.ErrInvalidOwnerAddr
	}

	if incrDeposit {
		blockTime := ctx.BlockHeader().Time
		ozoneLimitChange, err = k.AddMetaNodeDeposit(ctx, node, depositDelta)
		if err != nil {
			return sdk.ZeroInt(), time.Time{}, err
		}
		return ozoneLimitChange, blockTime, nil
	} else {
		ozoneLimitChange, completionTime, err := k.UnbondMetaNode(ctx, node, depositDelta.Amount)
		if err != nil {
			return sdk.ZeroInt(), time.Time{}, err
		}
		return ozoneLimitChange, completionTime, nil
	}
}

func (k Keeper) GetMetaNodeBondedToken(ctx sdk.Context) (token sdk.Coin) {
	metaNodeBondedAccAddr := k.accountKeeper.GetModuleAddress(types.MetaNodeBondedPool)
	if metaNodeBondedAccAddr == nil {
		ctx.Logger().Error("account address for meta node bonded pool does not exist.")
		return sdk.Coin{
			Denom:  k.BondDenom(ctx),
			Amount: sdk.ZeroInt(),
		}
	}
	return k.bankKeeper.GetBalance(ctx, metaNodeBondedAccAddr, k.BondDenom(ctx))
}

func (k Keeper) GetMetaNodeNotBondedToken(ctx sdk.Context) (token sdk.Coin) {
	metaNodeNotBondedAccAddr := k.accountKeeper.GetModuleAddress(types.MetaNodeNotBondedPool)
	if metaNodeNotBondedAccAddr == nil {
		ctx.Logger().Error("account address for meta node Not bonded pool does not exist.")
		return sdk.Coin{
			Denom:  k.BondDenom(ctx),
			Amount: sdk.ZeroInt(),
		}
	}
	return k.bankKeeper.GetBalance(ctx, metaNodeNotBondedAccAddr, k.BondDenom(ctx))
}

func (k Keeper) GetMetaNodeIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.MetaNodeKey)
	return iterator
}

func (k Keeper) OwnMetaNode(ctx sdk.Context, ownerAddr sdk.AccAddress, p2pAddr stratos.SdsAddress) bool {
	metaNode, found := k.GetMetaNode(ctx, p2pAddr)
	if !found {
		return false
	}

	if metaNode.OwnerAddress != ownerAddr.String() {
		return false
	}
	return true
}

func (k Keeper) GetMetaNodeBitMapIndex(ctx sdk.Context, networkAddr stratos.SdsAddress) (index int, err error) {
	k.UpdateMetaNodeBitMapIdxCache(ctx)

	k.cacheMutex.RLock()
	defer k.cacheMutex.RUnlock()

	index, ok := k.metaNodeBitMapIndexCache[networkAddr.String()]
	if !ok {
		return index, errors.New(fmt.Sprintf("Can not find meta-node %v from cache", networkAddr.String()))
	}
	if index < 0 {
		return index, errors.New(fmt.Sprintf("Can not find correct index of meta-node %v from cache", networkAddr.String()))
	}

	return index, nil
}

func (k Keeper) AddMetaNodeToBitMapIdxCache(networkAddr stratos.SdsAddress) {
	k.cacheMutex.Lock()
	defer k.cacheMutex.Unlock()

	k.metaNodeBitMapIndexCache[networkAddr.String()] = -1
	k.metaNodeBitMapIndexCacheStatus = types.CACHE_DIRTY
}

func (k Keeper) RemoveMetaNodeFromBitMapIdxCache(networkAddr stratos.SdsAddress) {
	k.cacheMutex.Lock()
	defer k.cacheMutex.Unlock()

	delete(k.metaNodeBitMapIndexCache, networkAddr.String())
	k.metaNodeBitMapIndexCacheStatus = types.CACHE_DIRTY
}

func (k Keeper) UpdateMetaNodeBitMapIdxCache(ctx sdk.Context) {
	k.cacheMutex.Lock()

	if k.metaNodeBitMapIndexCacheStatus == types.CACHE_NOT_DIRTY {
		k.cacheMutex.Unlock()
		return
	}
	if len(k.metaNodeBitMapIndexCache) == 0 {
		k.cacheMutex.Unlock()
		k.ReloadMetaNodeBitMapIdxCache(ctx)
		return
	}

	keys := make([]string, 0)
	for key, _ := range k.metaNodeBitMapIndexCache {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	for index, key := range keys {
		k.metaNodeBitMapIndexCache[key] = index
	}
	k.metaNodeBitMapIndexCacheStatus = types.CACHE_NOT_DIRTY
	k.cacheMutex.Unlock()
}

func (k Keeper) ReloadMetaNodeBitMapIdxCache(ctx sdk.Context) {
	k.cacheMutex.Lock()
	defer k.cacheMutex.Unlock()

	if k.metaNodeBitMapIndexCacheStatus == types.CACHE_NOT_DIRTY {
		return
	}
	keys := make([]string, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.MetaNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalMetaNode(k.cdc, iterator.Value())
		if node.GetSuspend() || node.GetStatus() == stakingtypes.Unbonded {
			continue
		}
		keys = append(keys, node.GetNetworkAddress())
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	for index, key := range keys {
		k.metaNodeBitMapIndexCache[key] = index
	}
	k.metaNodeBitMapIndexCacheStatus = types.CACHE_NOT_DIRTY
}
