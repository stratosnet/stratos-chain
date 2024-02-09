package keeper

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	sdkmath "cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

var (
	metaNodeBitMapIndexCache       = make(map[string]int)
	metaNodeBitMapIndexCacheStatus = &types.CacheStatus{Status: types.CACHE_DIRTY}
	cacheMutex                     = &sync.Mutex{}
)

// GetMetaNode get a single meta node
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

// GetAllActiveMetaNodes get the set of all bonded & not suspended meta nodes
func (k Keeper) GetAllActiveMetaNodes(ctx sdk.Context) (metaNodes []types.MetaNode) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.MetaNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalMetaNode(k.cdc, iterator.Value())
		if node.IsActivate() {
			metaNodes = append(metaNodes, node)
		}
	}
	return metaNodes
}

func (k Keeper) RegisterMetaNode(ctx sdk.Context, networkAddr stratos.SdsAddress, pubKey cryptotypes.PubKey, ownerAddr sdk.AccAddress,
	beneficiaryAddress sdk.AccAddress, description types.Description, deposit sdk.Coin) error {

	if _, found := k.GetMetaNode(ctx, networkAddr); found {
		ctx.Logger().Error("Meta node already exist")
		return types.ErrMetaNodePubKeyExists
	}
	if _, found := k.GetResourceNode(ctx, networkAddr); found {
		ctx.Logger().Error("Resource node with same network address already exist")
		return types.ErrResourceNodePubKeyExists
	}

	if deposit.GetDenom() != k.BondDenom(ctx) {
		return types.ErrBadDenom
	}

	metaNode, err := types.NewMetaNode(networkAddr, pubKey, ownerAddr, beneficiaryAddress, description, ctx.BlockHeader().Time)
	if err != nil {
		return err
	}
	_, _, _, err = k.AddMetaNodeDeposit(ctx, metaNode, deposit)
	if err != nil {
		return err
	}

	votingValidityPeriod := k.VotingPeriod(ctx)
	expireTime := ctx.BlockHeader().Time.Add(votingValidityPeriod)

	votePool := types.NewRegistrationVotePool(networkAddr, expireTime)
	k.SetMetaNodeRegistrationVotePool(ctx, votePool)

	return nil
}

// AddMetaNodeDeposit Update the tokens of an existing meta node
func (k Keeper) AddMetaNodeDeposit(ctx sdk.Context, metaNode types.MetaNode, tokenToAdd sdk.Coin,
) (ozoneLimitChange, availableTokenAmtBefore, availableTokenAmtAfter sdkmath.Int, err error) {

	coins := sdk.NewCoins(tokenToAdd)

	ownerAddr, err := sdk.AccAddressFromBech32(metaNode.GetOwnerAddress())
	if err != nil {
		return sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), types.ErrInvalidOwnerAddr
	}
	networkAddr, err := stratos.SdsAddressFromBech32(metaNode.GetNetworkAddress())
	if err != nil {
		return sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), types.ErrInvalidNetworkAddr
	}
	// sub coins from owner's wallet
	hasCoin := k.bankKeeper.HasBalance(ctx, ownerAddr, tokenToAdd)
	if !hasCoin {
		return sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), types.ErrInsufficientBalance
	}
	unbondingDeposit := k.GetUnbondingNodeBalance(ctx, networkAddr)
	availableTokenAmtBefore = metaNode.Tokens.Sub(unbondingDeposit)

	targetModuleAccName := ""

	switch metaNode.GetStatus() {
	case stakingtypes.Unbonded:
		targetModuleAccName = types.MetaNodeNotBondedPool
	case stakingtypes.Bonded:
		targetModuleAccName = types.MetaNodeBondedPool
	case stakingtypes.Unbonding:
		return sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), types.ErrUnbondingNode
	}

	if len(targetModuleAccName) > 0 {
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, ownerAddr, targetModuleAccName, coins)
		if err != nil {
			return sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), err
		}
	}

	metaNode = metaNode.AddToken(tokenToAdd.Amount)
	k.SetMetaNode(ctx, metaNode)

	if !metaNode.Suspend {
		ozoneLimitChange = k.IncreaseOzoneLimitByAddDeposit(ctx, tokenToAdd.Amount)
	} else {
		// if node is currently suspended, ozone limit will be increased upon unsuspension instead of NOW
		ozoneLimitChange = sdkmath.ZeroInt()
	}
	availableTokenAmtAfter = availableTokenAmtBefore.Add(tokenToAdd.Amount)
	return ozoneLimitChange, availableTokenAmtBefore, availableTokenAmtAfter, nil
}

func (k Keeper) RemoveTokenFromPoolWhileUnbondingMetaNode(ctx sdk.Context, tokenToSub sdk.Coin) error {
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

func (k Keeper) HandleVoteForMetaNodeRegistration(ctx sdk.Context, candidateNetworkAddr stratos.SdsAddress, candidateOwnerAddr sdk.AccAddress,
	opinion types.VoteOpinion, voterNetworkAddr stratos.SdsAddress, voterOwnerAddr sdk.AccAddress) (
	ozoneLimitChange sdkmath.Int, nodeStatus stakingtypes.BondStatus, err error) {

	// voter validation
	if !k.OwnActivateMetaNode(ctx, voterOwnerAddr, voterNetworkAddr) {
		err = types.ErrNoActiveVoterMetaNodeFound
	}

	// candidate validation
	candidateNode, found := k.GetMetaNode(ctx, candidateNetworkAddr)
	if !found {
		err = types.ErrNoCandidateMetaNodeFound
		return
	}
	if candidateNode.GetOwnerAddress() != candidateOwnerAddr.String() {
		err = types.ErrInvalidCandidateOwnerAddr
		return
	}

	// vote validation and handle voting
	votePool, found := k.GetMetaNodeRegistrationVotePool(ctx, candidateNetworkAddr)
	if !found {
		err = types.ErrNoRegistrationVotePoolFound
		return
	}
	if votePool.ExpireTime.Before(ctx.BlockHeader().Time) {
		err = types.ErrVoteExpired
		return
	}
	if hasStringValue(votePool.ApproveList, voterNetworkAddr.String()) || hasStringValue(votePool.RejectList, voterNetworkAddr.String()) {
		err = types.ErrDuplicateVoting
		return
	}

	if opinion.Equal(types.Approve) {
		votePool.ApproveList = append(votePool.ApproveList, voterNetworkAddr.String())
	} else {
		votePool.RejectList = append(votePool.RejectList, voterNetworkAddr.String())
	}
	k.SetMetaNodeRegistrationVotePool(ctx, votePool)

	ozoneLimitChange = sdkmath.ZeroInt()
	// if vote had already passed before, ozoneLimitChange should be zero.
	if votePool.IsVotePassed {
		return ozoneLimitChange, candidateNode.Status, nil
	}

	//if vote is yet to pass
	activeMetaNodeCount := len(k.GetAllActiveMetaNodes(ctx))
	voteCountRequiredToPass := activeMetaNodeCount*2/3 + 1
	//unbounded to bounded
	if len(votePool.ApproveList) >= voteCountRequiredToPass {
		candidateNode.Status = stakingtypes.Bonded
		candidateNode.Suspend = false
		k.SetMetaNode(ctx, candidateNode)
		// add new available meta node to cache
		networkAddr, _ := stratos.SdsAddressFromBech32(candidateNode.GetNetworkAddress())
		k.AddMetaNodeToBitMapIdxCache(networkAddr)
		// increase ozone limit after vote is approved
		ozoneLimitChange = k.IncreaseOzoneLimitByAddDeposit(ctx, candidateNode.Tokens)
		// increase mata node count
		newBondedMetaNodeCount := k.GetBondedMetaNodeCnt(ctx).Add(sdkmath.OneInt())
		k.SetBondedMetaNodeCnt(ctx, newBondedMetaNodeCount)
		// move deposit from not bonded pool to bonded pool
		tokenToBond := sdk.NewCoin(k.BondDenom(ctx), candidateNode.Tokens)
		// sub coins from not bonded pool
		err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.MetaNodeNotBondedPool, types.MetaNodeBondedPool, sdk.NewCoins(tokenToBond))
		if err != nil {
			return
		}

		votePool.IsVotePassed = true
		k.SetMetaNodeRegistrationVotePool(ctx, votePool)
	}

	return ozoneLimitChange, candidateNode.Status, nil
}

func (k Keeper) UpdateMetaNode(ctx sdk.Context, description types.Description,
	networkAddr stratos.SdsAddress, ownerAddr sdk.AccAddress, beneficiaryAddr sdk.AccAddress) error {

	node, found := k.GetMetaNode(ctx, networkAddr)
	if !found {
		return types.ErrNoMetaNodeFound
	}

	ownerAddrNode, _ := sdk.AccAddressFromBech32(node.GetOwnerAddress())
	if !ownerAddrNode.Equals(ownerAddr) {
		return types.ErrInvalidOwnerAddr
	}

	if len(beneficiaryAddr) > 0 {
		node.BeneficiaryAddress = beneficiaryAddr.String()
	}

	node.Description = description

	k.SetMetaNode(ctx, node)

	return nil
}

func (k Keeper) UpdateMetaNodeDeposit(ctx sdk.Context, networkAddr stratos.SdsAddress, ownerAddr sdk.AccAddress, depositDelta sdk.Coin) (
	ozoneLimitChange, availableTokenAmtBefore, availableTokenAmtAfter sdkmath.Int, unbondingMatureTime time.Time, metaNode types.MetaNode, err error) {

	if depositDelta.GetDenom() != k.BondDenom(ctx) {
		return sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), time.Time{}, types.MetaNode{}, types.ErrBadDenom
	}

	node, found := k.GetMetaNode(ctx, networkAddr)
	if !found {
		return sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), time.Time{}, types.MetaNode{}, types.ErrNoMetaNodeFound
	}

	ownerAddrNode, _ := sdk.AccAddressFromBech32(node.GetOwnerAddress())
	if !ownerAddrNode.Equals(ownerAddr) {
		return sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), time.Time{}, types.MetaNode{}, types.ErrInvalidOwnerAddr
	}

	// not allow to decrease deposit
	if depositDelta.Amount.IsPositive() {
		blockTime := ctx.BlockHeader().Time
		ozoneLimitChange, availableTokenAmtBefore, availableTokenAmtAfter, err = k.AddMetaNodeDeposit(ctx, node, depositDelta)
		if err != nil {
			return sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), time.Time{}, types.MetaNode{}, err
		}
		return ozoneLimitChange, availableTokenAmtBefore, availableTokenAmtAfter, blockTime, node, nil
	}
	return sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), time.Time{}, types.MetaNode{}, err
}

// Remove all deposit from meta-node, start unbonding process
func (k Keeper) UnbondMetaNode(ctx sdk.Context, p2pAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
) (ozoneLimitChange sdkmath.Int, depositToRemove sdk.Coin, unbondingMatureTime time.Time, err error) {

	metaNode, availableDeposit, err := k.getMetaNodeAvailableDeposit(ctx, p2pAddress)
	if err != nil {
		return
	}
	if availableDeposit.LTE(sdkmath.ZeroInt()) {
		err = types.ErrInsufficientBalance
		return
	}
	if ownerAddress.String() != metaNode.GetOwnerAddress() {
		err = types.ErrInvalidOwnerAddr
		return
	}
	if k.HasMaxUnbondingNodeEntries(ctx, p2pAddress) {
		err = types.ErrMaxUnbondingNodeEntries
		return
	}

	switch metaNode.GetStatus() {

	case stakingtypes.Bonded:
		// suspended node cannot be unbonded (avoid dup deposit decrease with node suspension)
		if metaNode.GetSuspend() {
			err = types.ErrInvalidSuspensionStatForUnbondNode
			return
		}
		// to prevent remainingOzoneLimit from being negative value
		if !k.IsUnbondable(ctx, availableDeposit) {
			err = types.ErrInsufficientBalance
			return
		}
		depositToRemove = sdk.NewCoin(k.BondDenom(ctx), availableDeposit)
		// transfer the node tokens to the not bonded pool
		k.bondedToUnbonding(ctx, metaNode, true, depositToRemove)
		// adjust ozone limit
		ozoneLimitChange = k.DecreaseOzoneLimitBySubtractDeposit(ctx, availableDeposit)
		// decrease bonded meta node count
		newBondedMetaNodeCnt := k.GetBondedMetaNodeCnt(ctx).Sub(sdkmath.OneInt())
		k.SetBondedMetaNodeCnt(ctx, newBondedMetaNodeCnt)

	case stakingtypes.Unbonded:
		// unbond the meta node which is not passed the registration voting
		votePool, exist := k.GetMetaNodeRegistrationVotePool(ctx, p2pAddress)
		if !exist {
			err = types.ErrNoRegistrationVotePoolFound
			return
		}
		// to be qualified to withdraw, meta node must be unbonded && suspended && of non-passed vote
		if !metaNode.GetSuspend() || votePool.IsVotePassed {
			err = types.ErrInvalidNodeStat
			return
		}

	default:
		err = types.ErrInvalidNodeStat
		return

	}

	unbondingMatureTime = k.calcUnbondingMatureTime(ctx, metaNode.GetStatus(), metaNode.GetCreationTime())
	// Set the unbonding mature time and completion height appropriately
	unbondingNode := k.SetUnbondingNodeEntry(ctx, p2pAddress, true, ctx.BlockHeight(), unbondingMatureTime, availableDeposit)
	// Add to unbonding node queue
	k.InsertUnbondingNodeQueue(ctx, unbondingNode, unbondingMatureTime)
	ctx.Logger().Info("Unbonding meta node " + unbondingNode.String() + "\n after mature time" + unbondingMatureTime.String())

	// all deposit is being unbonded, update status to Unbonding
	metaNode.Status = stakingtypes.Unbonding
	k.SetMetaNode(ctx, metaNode)

	// remove from cache just in case
	k.RemoveMetaNodeFromBitMapIdxCache(p2pAddress)
	// remove record from vote pool
	if _, found := k.GetMetaNodeRegistrationVotePool(ctx, p2pAddress); found {
		ctx.Logger().Info("DeleteMetaNodeRegistrationVotePool of meta node " + p2pAddress.String())
		k.DeleteMetaNodeRegistrationVotePool(ctx, p2pAddress)
	}

	return
}

// Active meta node kick other meta node
func (k Keeper) HandleVoteForKickMetaNode(ctx sdk.Context, targetNetworkAddr stratos.SdsAddress,
	opinion types.VoteOpinion, voterNetworkAddr stratos.SdsAddress, voterOwnerAddr sdk.AccAddress) (
	nodeStatus stakingtypes.BondStatus, ozoneLimitChange sdkmath.Int, depositToRemove sdk.Coin, unbondingMatureTime time.Time, err error) {

	// voter validation
	if !k.OwnActivateMetaNode(ctx, voterOwnerAddr, voterNetworkAddr) {
		err = types.ErrNoActiveVoterMetaNodeFound
	}
	// target meta node validation
	targetMetaNode, found := k.GetMetaNode(ctx, targetNetworkAddr)
	if !found {
		err = types.ErrNoActiveTargetMetaNodeFound
		return
	}
	// target meta node should be active
	if !targetMetaNode.IsActivate() {
		err = types.ErrNoActiveTargetMetaNodeFound
		return
	}
	// vote validation and handle voting
	votePool, found := k.GetKickMetaNodeVotePool(ctx, targetNetworkAddr)
	if !found {
		// create voting pool
		votingValidityPeriod := k.VotingPeriod(ctx)
		expireTime := ctx.BlockHeader().Time.Add(votingValidityPeriod)

		votePool = types.NewKickMetaNodeVotePool(targetNetworkAddr, expireTime)
		k.SetKickMetaNodeVotePool(ctx, votePool)
	}
	if votePool.ExpireTime.Before(ctx.BlockHeader().Time) {
		err = types.ErrVoteExpired
		return
	}
	if hasStringValue(votePool.ApproveList, voterNetworkAddr.String()) || hasStringValue(votePool.RejectList, voterNetworkAddr.String()) {
		err = types.ErrDuplicateVoting
		return
	}

	if opinion.Equal(types.Approve) {
		votePool.ApproveList = append(votePool.ApproveList, voterNetworkAddr.String())
	} else {
		votePool.RejectList = append(votePool.RejectList, voterNetworkAddr.String())
	}
	k.SetKickMetaNodeVotePool(ctx, votePool)

	ozoneLimitChange = sdkmath.ZeroInt()
	// if vote had already passed before, don't trigger unBonding meta node
	if votePool.IsVotePassed {
		nodeStatus = targetMetaNode.Status
		return
	}

	//if vote is yet to pass
	activeMetaNodeCount := len(k.GetAllActiveMetaNodes(ctx))
	voteCountRequiredToPass := activeMetaNodeCount*2/3 + 1
	//bounded to unbonding
	if len(votePool.ApproveList) >= voteCountRequiredToPass {
		targetOwnerAddr, _ := sdk.AccAddressFromBech32(targetMetaNode.GetOwnerAddress())

		ozoneLimitChange, depositToRemove, unbondingMatureTime, err = k.UnbondMetaNode(ctx, targetNetworkAddr, targetOwnerAddr)
		if err != nil {
			return
		}

		votePool.IsVotePassed = true
		k.SetKickMetaNodeVotePool(ctx, votePool)
	}
	nodeStatus = targetMetaNode.Status

	return
}

func (k Keeper) GetMetaNodeBondedToken(ctx sdk.Context) (token sdk.Coin) {
	metaNodeBondedAccAddr := k.accountKeeper.GetModuleAddress(types.MetaNodeBondedPool)
	if metaNodeBondedAccAddr == nil {
		ctx.Logger().Error("account address for meta node bonded pool does not exist.")
		return sdk.Coin{
			Denom:  k.BondDenom(ctx),
			Amount: sdkmath.ZeroInt(),
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
			Amount: sdkmath.ZeroInt(),
		}
	}
	return k.bankKeeper.GetBalance(ctx, metaNodeNotBondedAccAddr, k.BondDenom(ctx))
}

func (k Keeper) GetMetaNodeIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.MetaNodeKey)
	return iterator
}

func (k Keeper) OwnActivateMetaNode(ctx sdk.Context, ownerAddr sdk.AccAddress, p2pAddr stratos.SdsAddress) bool {
	metaNode, found := k.GetMetaNode(ctx, p2pAddr)
	if !found {
		return false
	}
	if !metaNode.IsActivate() {
		return false
	}
	if metaNode.OwnerAddress != ownerAddr.String() {
		return false
	}
	return true
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

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	index, ok := metaNodeBitMapIndexCache[networkAddr.String()]
	if !ok {
		return index, errors.New(fmt.Sprintf("Can not find meta-node %v from cache", networkAddr.String()))
	}
	if index < 0 {
		return index, errors.New(fmt.Sprintf("Can not find correct index of meta-node %v from cache", networkAddr.String()))
	}

	return index, nil
}

func (k Keeper) AddMetaNodeToBitMapIdxCache(networkAddr stratos.SdsAddress) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	metaNodeBitMapIndexCache[networkAddr.String()] = -1
	metaNodeBitMapIndexCacheStatus.Status = types.CACHE_DIRTY
}

func (k Keeper) RemoveMetaNodeFromBitMapIdxCache(networkAddr stratos.SdsAddress) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	delete(metaNodeBitMapIndexCache, networkAddr.String())
	metaNodeBitMapIndexCacheStatus.Status = types.CACHE_DIRTY
}

func (k Keeper) UpdateMetaNodeBitMapIdxCache(ctx sdk.Context) {
	cacheMutex.Lock()

	if metaNodeBitMapIndexCacheStatus.Status == types.CACHE_NOT_DIRTY {
		cacheMutex.Unlock()
		return
	}
	if len(metaNodeBitMapIndexCache) == 0 {
		cacheMutex.Unlock()
		k.ReloadMetaNodeBitMapIdxCache(ctx)
		return
	}

	keys := make([]string, 0)
	for key, _ := range metaNodeBitMapIndexCache {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	for index, key := range keys {
		metaNodeBitMapIndexCache[key] = index
	}
	metaNodeBitMapIndexCacheStatus.Status = types.CACHE_NOT_DIRTY
	cacheMutex.Unlock()
}

func (k Keeper) ReloadMetaNodeBitMapIdxCache(ctx sdk.Context) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if metaNodeBitMapIndexCacheStatus.Status == types.CACHE_NOT_DIRTY {
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
		metaNodeBitMapIndexCache[key] = index
	}
	metaNodeBitMapIndexCacheStatus.Status = types.CACHE_NOT_DIRTY
}

func (k Keeper) getMetaNodeAvailableDeposit(ctx sdk.Context, p2pAddress stratos.SdsAddress) (types.MetaNode, sdkmath.Int, error) {
	metaNode, found := k.GetMetaNode(ctx, p2pAddress)
	if !found {
		return types.MetaNode{}, sdkmath.ZeroInt(), types.ErrNoMetaNodeFound
	}
	unBondingDeposit := k.GetUnbondingNodeBalance(ctx, p2pAddress)
	availableDeposit := metaNode.Tokens.Sub(unBondingDeposit)

	return metaNode, availableDeposit, nil
}
