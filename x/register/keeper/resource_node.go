package keeper

import (
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// GetResourceNode get a single resource node
func (k Keeper) GetResourceNode(ctx sdk.Context, p2pAddress stratos.SdsAddress) (resourceNode types.ResourceNode, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetResourceNodeKey(p2pAddress))
	if value == nil {
		return resourceNode, false
	}
	resourceNode = types.MustUnmarshalResourceNode(k.cdc, value)
	return resourceNode, true
}

// SetResourceNode sets the main record holding resource node details
func (k Keeper) SetResourceNode(ctx sdk.Context, resourceNode types.ResourceNode) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalResourceNode(k.cdc, resourceNode)
	networkAddr, _ := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddress())
	store.Set(types.GetResourceNodeKey(networkAddr), bz)
}

// GetAllResourceNodes get the set of all resource nodes with no limits, used during genesis dump
// Iteration for all resource nodes
func (k Keeper) GetAllResourceNodes(ctx sdk.Context) (resourceNodes types.ResourceNodes) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ResourceNodeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		node := types.MustUnmarshalResourceNode(k.cdc, iterator.Value())
		resourceNodes = append(resourceNodes, node)
	}
	return
}

func (k Keeper) GetResourceNodeIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ResourceNodeKey)
	return iterator
}

// AddResourceNodeDeposit Update the tokens of an existing resource node
func (k Keeper) AddResourceNodeDeposit(ctx sdk.Context, resourceNode types.ResourceNode, tokenToAdd sdk.Coin,
) (ozoneLimitChange, availableTokenAmtBefore, availableTokenAmtAfter sdk.Int, err error) {

	needAddCount := true
	networkAddr, err := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddress())
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), types.ErrInvalidNetworkAddr
	}
	nodeStored, found := k.GetResourceNode(ctx, networkAddr)
	if found && nodeStored.IsBonded() {
		needAddCount = false
	}
	unbondingDeposit := k.GetUnbondingNodeBalance(ctx, networkAddr)
	availableTokenAmtBefore = resourceNode.Tokens.Sub(unbondingDeposit)

	coins := sdk.NewCoins(tokenToAdd)

	ownerAddr, err := sdk.AccAddressFromBech32(resourceNode.GetOwnerAddress())
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), types.ErrInvalidOwnerAddr
	}

	// sub coins from owner's wallet
	hasCoin := k.bankKeeper.HasBalance(ctx, ownerAddr, tokenToAdd)
	if !hasCoin {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), types.ErrInsufficientBalance
	}

	targetModuleAccName := ""

	switch resourceNode.GetStatus() {
	case stakingtypes.Unbonded:
		targetModuleAccName = types.ResourceNodeNotBondedPool
	case stakingtypes.Bonded:
		targetModuleAccName = types.ResourceNodeBondedPool
	case stakingtypes.Unbonding:
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), types.ErrUnbondingNode
	}

	if len(targetModuleAccName) > 0 {
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, ownerAddr, targetModuleAccName, coins)
		if err != nil {
			return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), err
		}
	}

	resourceNode = resourceNode.AddToken(tokenToAdd.Amount)
	//resourceNode.Suspend = false

	// set status from unBonded to bonded & move deposit from not bonded token pool to bonded token pool
	// since resource node registration does not require voting for now
	if resourceNode.Status == stakingtypes.Unbonded {
		resourceNode.Status = stakingtypes.Bonded

		tokenToTrasfer := sdk.NewCoin(k.BondDenom(ctx), resourceNode.Tokens)
		nBondedResourceAccountAddr := k.accountKeeper.GetModuleAddress(types.ResourceNodeNotBondedPool)
		if nBondedResourceAccountAddr == nil {
			ctx.Logger().Error("not bonded account address for resource nodes does not exist.")
			return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), types.ErrUnknownAccountAddress
		}

		hasCoin := k.bankKeeper.HasBalance(ctx, nBondedResourceAccountAddr, tokenToTrasfer)
		if !hasCoin {
			return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), types.ErrInsufficientBalance
		}

		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ResourceNodeNotBondedPool, types.ResourceNodeBondedPool, sdk.NewCoins(tokenToTrasfer))
		if err != nil {
			return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), types.ErrInsufficientBalance
		}
	}

	k.SetResourceNode(ctx, resourceNode)

	if needAddCount {
		// increase resource node count
		v := k.GetBondedResourceNodeCnt(ctx)
		count := v.Add(sdk.NewInt(1))
		k.SetBondedResourceNodeCnt(ctx, count)
	}

	availableTokenAmtAfter = availableTokenAmtBefore.Add(tokenToAdd.Amount)
	return sdk.ZeroInt(), availableTokenAmtBefore, availableTokenAmtAfter, nil
}

func (k Keeper) RemoveTokenFromPoolWhileUnbondingResourceNode(ctx sdk.Context, resourceNode types.ResourceNode, tokenToSub sdk.Coin) error {
	bondedResourceAccountAddr := k.accountKeeper.GetModuleAddress(types.ResourceNodeBondedPool)
	if bondedResourceAccountAddr == nil {
		ctx.Logger().Error("bonded pool account address for resource nodes does not exist.")
		return types.ErrUnknownAccountAddress
	}

	hasCoin := k.bankKeeper.HasBalance(ctx, bondedResourceAccountAddr, tokenToSub)
	if !hasCoin {
		return types.ErrInsufficientBalance
	}

	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ResourceNodeBondedPool, types.ResourceNodeNotBondedPool, sdk.NewCoins(tokenToSub))
	if err != nil {
		return types.ErrInsufficientBalance
	}
	return nil
}

// SubtractResourceNodeDeposit Update the tokens of an existing resource node
func (k Keeper) SubtractResourceNodeDeposit(ctx sdk.Context, resourceNode types.ResourceNode, tokenToSub sdk.Coin) error {
	networkAddr, err := stratos.SdsAddressFromBech32(resourceNode.GetNetworkAddress())
	if err != nil {
		return types.ErrInvalidNetworkAddr
	}
	ownerAddr, err := sdk.AccAddressFromBech32(resourceNode.GetOwnerAddress())
	if err != nil {
		return types.ErrInvalidOwnerAddr
	}

	ownerAcc := k.accountKeeper.GetAccount(ctx, ownerAddr)
	if ownerAcc == nil {
		return types.ErrNoOwnerAccountFound
	}

	coins := sdk.NewCoins(tokenToSub)

	if resourceNode.Tokens.LT(tokenToSub.Amount) {
		return types.ErrInsufficientBalance
	}

	// deduct tokens from NotBondedPool
	nBondedResourceAccountAddr := k.accountKeeper.GetModuleAddress(types.ResourceNodeNotBondedPool)
	if nBondedResourceAccountAddr == nil {
		ctx.Logger().Error("not bonded account address for resource nodes does not exist.")
		return types.ErrUnknownAccountAddress
	}

	hasCoin := k.bankKeeper.HasBalance(ctx, nBondedResourceAccountAddr, tokenToSub)
	if !hasCoin {
		return types.ErrInsufficientBalanceOfNotBondedPool
	}

	// deduct slashing amount first, slashed amt goes into TotalSlashedPool
	remaining, slashed := k.DeductSlashing(ctx, ownerAddr, coins, k.BondDenom(ctx))
	if !remaining.IsZero() {
		// add remaining tokens to owner acc
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ResourceNodeNotBondedPool, ownerAddr, remaining)
		if err != nil {
			return err
		}
	}
	if !slashed.IsZero() {
		// slashed token send to community_pool
		resNodeNotBondedPoolAddr := k.accountKeeper.GetModuleAddress(types.ResourceNodeNotBondedPool)
		err = k.distrKeeper.FundCommunityPool(ctx, slashed, resNodeNotBondedPoolAddr)
		if err != nil {
			return err
		}
	}

	resourceNode = resourceNode.SubToken(tokenToSub.Amount)
	newDeposit := resourceNode.Tokens

	k.SetResourceNode(ctx, resourceNode)

	if newDeposit.IsZero() {
		err = k.removeResourceNode(ctx, networkAddr)
		if err != nil {
			return err
		}
	}
	return nil
}

// remove the resource node record and associated indexes
func (k Keeper) removeResourceNode(ctx sdk.Context, addr stratos.SdsAddress) error {
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

func (k Keeper) RegisterResourceNode(ctx sdk.Context, networkAddr stratos.SdsAddress, pubKey cryptotypes.PubKey, ownerAddr sdk.AccAddress,
	description types.Description, nodeType types.NodeType, deposit sdk.Coin) (ozoneLimitChange sdk.Int, err error) {

	if _, found := k.GetResourceNode(ctx, networkAddr); found {
		ctx.Logger().Error("Resource node already exist")
		return ozoneLimitChange, types.ErrResourceNodePubKeyExists
	}
	if _, found := k.GetMetaNode(ctx, networkAddr); found {
		ctx.Logger().Error("Meta node with same network address already exist")
		return ozoneLimitChange, types.ErrMetaNodePubKeyExists
	}

	if deposit.GetDenom() != k.BondDenom(ctx) {
		return ozoneLimitChange, types.ErrBadDenom
	}
	if deposit.IsLT(k.ResourceNodeMinDeposit(ctx)) {
		return ozoneLimitChange, types.ErrInsufficientDeposit
	}

	resourceNode, err := types.NewResourceNode(networkAddr, pubKey, ownerAddr, description, nodeType, ctx.BlockHeader().Time)
	if err != nil {
		return ozoneLimitChange, err
	}
	ozoneLimitChange, _, _, err = k.AddResourceNodeDeposit(ctx, resourceNode, deposit)
	return ozoneLimitChange, err
}

func (k Keeper) UpdateResourceNode(ctx sdk.Context, description types.Description, nodeType types.NodeType,
	networkAddr stratos.SdsAddress, ownerAddr sdk.AccAddress) error {

	node, found := k.GetResourceNode(ctx, networkAddr)
	if !found {
		return types.ErrNoResourceNodeFound
	}

	ownerAddrNode, _ := sdk.AccAddressFromBech32(node.GetOwnerAddress())
	if !ownerAddrNode.Equals(ownerAddr) {
		return types.ErrInvalidOwnerAddr
	}

	node.Description = description
	if nodeType != 0 {
		node.NodeType = uint32(nodeType)
	}

	k.SetResourceNode(ctx, node)

	return nil
}

// Add deposit only
func (k Keeper) UpdateResourceNodeDeposit(ctx sdk.Context, networkAddr stratos.SdsAddress, ownerAddr sdk.AccAddress, depositDelta sdk.Coin) (
	ozoneLimitChange, availableTokenAmtBefore, availableTokenAmtAfter sdk.Int, completionTime time.Time, resourcenode types.ResourceNode, err error) {

	if depositDelta.GetDenom() != k.BondDenom(ctx) {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), time.Time{}, types.ResourceNode{}, types.ErrBadDenom
	}

	node, found := k.GetResourceNode(ctx, networkAddr)
	if !found {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), time.Time{}, types.ResourceNode{}, types.ErrNoResourceNodeFound
	}

	ownerAddrNode, _ := sdk.AccAddressFromBech32(node.GetOwnerAddress())
	if !ownerAddrNode.Equals(ownerAddr) {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), time.Time{}, types.ResourceNode{}, types.ErrInvalidOwnerAddr
	}

	completionTime = ctx.BlockHeader().Time
	ozoneLimitChange, availableTokenAmtBefore, availableTokenAmtAfter, err = k.AddResourceNodeDeposit(ctx, node, depositDelta)
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), time.Time{}, types.ResourceNode{}, err
	}
	return ozoneLimitChange, availableTokenAmtBefore, availableTokenAmtAfter, completionTime, node, nil
}

func (k Keeper) UpdateEffectiveDeposit(ctx sdk.Context, networkAddr stratos.SdsAddress, effectiveDepositAfter sdk.Int) (
	ozoneLimitChange, effectiveDepositChange sdk.Int, isUnsuspendedDuringUpdate bool, err error) {

	node, found := k.GetResourceNode(ctx, networkAddr)
	if !found {
		return sdk.ZeroInt(), sdk.ZeroInt(), false, types.ErrNoResourceNodeFound
	}

	// before calc ozone limit change, get unbonding deposit and calc effective deposit to trigger ozLimit change
	unbondingDeposit := k.GetUnbondingNodeBalance(ctx, networkAddr)
	// no effective deposit after subtracting unbonding deposit
	if node.Tokens.LTE(unbondingDeposit) {
		return sdk.ZeroInt(), sdk.ZeroInt(), false, types.ErrInsufficientBalance
	}
	availableDeposit := node.Tokens.Sub(unbondingDeposit)
	if availableDeposit.LT(effectiveDepositAfter) {
		return sdk.ZeroInt(), sdk.ZeroInt(), false, types.ErrInsufficientBalance
	}

	isUnsuspendedDuringUpdate = node.Suspend == true && node.EffectiveTokens.Equal(sdk.ZeroInt()) && effectiveDepositAfter.GT(sdk.ZeroInt())

	effectiveDepositBefore := sdk.NewInt(0).Add(node.EffectiveTokens)
	effectiveDepositChange = effectiveDepositAfter.Sub(effectiveDepositBefore)

	node.EffectiveTokens = effectiveDepositAfter
	// effectiveDepositAfter > 0 means node.Suspend = false
	node.Suspend = false
	k.SetResourceNode(ctx, node)

	if effectiveDepositChange.IsNegative() && k.IsUnbondable(ctx, effectiveDepositChange.Abs()) {
		ozoneLimitChange = k.DecreaseOzoneLimitBySubtractDeposit(ctx, effectiveDepositChange.Abs())
	}
	if effectiveDepositChange.IsPositive() {
		ozoneLimitChange = k.IncreaseOzoneLimitByAddDeposit(ctx, effectiveDepositChange)
	}
	return ozoneLimitChange, effectiveDepositChange, isUnsuspendedDuringUpdate, nil
}

func (k Keeper) GetResourceNodeBondedToken(ctx sdk.Context) (token sdk.Coin) {
	resourceNodeBondedAccAddr := k.accountKeeper.GetModuleAddress(types.ResourceNodeBondedPool)
	if resourceNodeBondedAccAddr == nil {
		ctx.Logger().Error("account address for resource node bonded pool does not exist.")
		return sdk.Coin{
			Denom:  k.BondDenom(ctx),
			Amount: sdk.ZeroInt(),
		}
	}
	return k.bankKeeper.GetBalance(ctx, resourceNodeBondedAccAddr, k.BondDenom(ctx))
}

func (k Keeper) GetResourceNodeNotBondedToken(ctx sdk.Context) (token sdk.Coin) {
	resourceNodeNotBondedAccAddr := k.accountKeeper.GetModuleAddress(types.ResourceNodeNotBondedPool)
	if resourceNodeNotBondedAccAddr == nil {
		ctx.Logger().Error("account address for resource node Not bonded pool does not exist.")
		return sdk.Coin{
			Denom:  k.BondDenom(ctx),
			Amount: sdk.ZeroInt(),
		}
	}
	return k.bankKeeper.GetBalance(ctx, resourceNodeNotBondedAccAddr, k.BondDenom(ctx))
}
