package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

func (k Keeper) DistributePotRewardForTestnet(ctx sdk.Context, trafficList []types.SingleWalletVolume, epoch sdk.Int) (totalConsumedOzone sdk.Dec, err error) {
	distributeGoal := types.InitDistributeGoal()
	rewardDetailMap := make(map[string]types.Reward) //key: wallet address

	//1, calc traffic reward in total
	totalConsumedOzone, distributeGoal, err = k.CalcTrafficRewardInTotal(ctx, trafficList, distributeGoal)
	if err != nil {
		return totalConsumedOzone, err
	}

	//2, calc mining reward in total,
	//*************** For incentive testnet, the mining reward to blockchain (stakeReward) are evenly distributed to all the nodes *****************
	distributeGoal, idxNodes, resNodes, err := k.CalcMiningRewardInTotalForTestnet(ctx, distributeGoal)
	if err != nil && err != types.ErrOutOfIssuance {
		return totalConsumedOzone, err
	}

	/**
	distributeGoalBalance is used for keeping the balance to return to the reward provider's account

	After calculation by step 3,4,
	distributeGoalBalance retains FULL rewards that will be distributed to validators
	& the REMAINING rewards that will be returned to the reward provider’s account
	*/
	distributeGoalBalance := distributeGoal

	//3, calc reward for resource node, store to rewardDetailMap by wallet address(owner address)
	rewardDetailMap, distributeGoalBalance = k.CalcRewardForResourceNodeForTestnet(ctx, trafficList, distributeGoalBalance, rewardDetailMap, resNodes)

	//4, calc reward from indexing node, store to rewardDetailMap by wallet address(owner address)
	rewardDetailMap, distributeGoalBalance = k.CalcRewardForIndexingNodeForTestnet(ctx, distributeGoalBalance, rewardDetailMap, idxNodes)

	//5, deduct reward from provider account (the value of parameter of distributeGoal will not change)
	err = k.deductRewardFromRewardProviderAccount(ctx, distributeGoal, epoch)
	if err != nil {
		return totalConsumedOzone, err
	}

	//6, distribute skate reward to fee pool for validators
	distributeGoalBalance, err = k.distributeValidatorRewardForTestnet(ctx, distributeGoalBalance)
	if err != nil {
		return totalConsumedOzone, err
	}

	//7, IMPORTANT: sort map and convert to slice to keep the order
	rewardDetailList := sortDetailMapToSlice(rewardDetailMap)

	//8, distribute all rewards to resource nodes & indexing nodes
	err = k.distributeRewardToSdsNodes(ctx, rewardDetailList, epoch)
	if err != nil {
		return totalConsumedOzone, err
	}

	//9, return balance to traffic pool & mining pool
	err = k.returnBalanceForTestnet(ctx, distributeGoalBalance, epoch)
	if err != nil {
		return totalConsumedOzone, err
	}

	return totalConsumedOzone, nil
}

// allocate mining reward from foundation account
func (k Keeper) CalcMiningRewardInTotalForTestnet(ctx sdk.Context, distributeGoal types.DistributeGoal) (
	types.DistributeGoal, []regtypes.IndexingNode, []regtypes.ResourceNode, error) {
	totalMinedTokens := k.GetTotalMinedTokens(ctx)
	miningParam, err := k.GetMiningRewardParamByMinedToken(ctx, totalMinedTokens)

	totalMiningReward := miningParam.MiningReward
	if err != nil {
		return distributeGoal, nil, nil, err
	}
	stakeReward := totalMiningReward.Amount.ToDec().
		Mul(miningParam.BlockChainPercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	trafficReward := totalMiningReward.Amount.ToDec().
		Mul(miningParam.ResourceNodePercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	indexingReward := totalMiningReward.Amount.ToDec().
		Mul(miningParam.MetaNodePercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()

	stakeRewardToValidators, stakeRewardToResourceNodes, stakeRewardToIndexingNodes, idxNodes, resNodes := k.splitRewardEvenly(ctx, stakeReward)
	distributeGoal = distributeGoal.AddBlockChainRewardToValidatorFromMiningPool(sdk.NewCoin(k.RewardDenom(ctx), stakeRewardToValidators))
	distributeGoal = distributeGoal.AddBlockChainRewardToResourceNodeFromMiningPool(sdk.NewCoin(k.RewardDenom(ctx), stakeRewardToResourceNodes))
	distributeGoal = distributeGoal.AddBlockChainRewardToIndexingNodeFromMiningPool(sdk.NewCoin(k.RewardDenom(ctx), stakeRewardToIndexingNodes))
	distributeGoal = distributeGoal.AddTrafficRewardToResourceNodeFromMiningPool(sdk.NewCoin(k.RewardDenom(ctx), trafficReward))
	distributeGoal = distributeGoal.AddMetaNodeRewardToIndexingNodeFromMiningPool(sdk.NewCoin(k.RewardDenom(ctx), indexingReward))
	return distributeGoal, idxNodes, resNodes, nil
}

func (k Keeper) splitRewardEvenly(ctx sdk.Context, totalReward sdk.Int,
) (validatorReward sdk.Int, resourceNodeReward sdk.Int, indexingNodeReward sdk.Int,
	indNodes []regtypes.IndexingNode, resNodes []regtypes.ResourceNode) {
	validatorCnt := sdk.ZeroDec()
	indexingNodeCnt := sdk.ZeroDec()
	resourceNodeCnt := sdk.ZeroDec()
	indNodes = make([]regtypes.IndexingNode, 0)
	resNodes = make([]regtypes.ResourceNode, 0)
	validatorList := k.StakingKeeper.GetAllValidators(ctx)
	for _, validator := range validatorList {
		if validator.IsBonded() && !validator.IsJailed() {
			validatorCnt = validatorCnt.Add(sdk.OneDec())
		}
	}

	indexingNodeList := k.RegisterKeeper.GetAllIndexingNodes(ctx)
	for _, indexingNode := range indexingNodeList {
		if indexingNode.IsBonded() && !indexingNode.IsSuspended() {
			indexingNodeCnt = indexingNodeCnt.Add(sdk.OneDec())
			indNodes = append(indNodes, indexingNode)
		}
	}

	resourceNodeList := k.RegisterKeeper.GetAllResourceNodes(ctx)
	for _, resourceNode := range resourceNodeList {
		if resourceNode.IsBonded() && !resourceNode.IsSuspended() {
			resourceNodeCnt = resourceNodeCnt.Add(sdk.OneDec())
			resNodes = append(resNodes, resourceNode)
		}
	}

	totalNodes := validatorCnt.Add(indexingNodeCnt).Add(resourceNodeCnt)

	validatorReward = totalReward.ToDec().Mul(validatorCnt).Quo(totalNodes).TruncateInt()
	indexingNodeReward = totalReward.ToDec().Mul(indexingNodeCnt).Quo(totalNodes).TruncateInt()
	resourceNodeReward = totalReward.ToDec().Mul(resourceNodeCnt).Quo(totalNodes).TruncateInt()
	return
}

func (k Keeper) CalcRewardForResourceNodeForTestnet(ctx sdk.Context, trafficList []types.SingleWalletVolume,
	distributeGoal types.DistributeGoal, rewardDetailMap map[string]types.Reward, resourceNodes []regtypes.ResourceNode,
) (map[string]types.Reward, types.DistributeGoal) {

	totalUsedFromMiningPool := sdk.NewCoin(k.RewardDenom(ctx), sdk.ZeroInt())
	totalUsedFromTrafficPool := sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())

	// 1, calc stake reward
	totalStakeOfResourceNodes := k.RegisterKeeper.GetResourceNodeBondedToken(ctx).Amount
	resourceNodeCnt := sdk.NewDec(int64(len(resourceNodes)))
	for _, node := range resourceNodes {
		walletAddr := node.GetOwnerAddr()

		shareOfToken := node.GetTokens().ToDec().Quo(totalStakeOfResourceNodes.ToDec())
		stakeRewardFromMiningPool := distributeGoal.BlockChainRewardToResourceNodeFromMiningPool.Amount.ToDec().Quo(resourceNodeCnt).TruncateInt()
		stakeRewardFromTrafficPool := distributeGoal.BlockChainRewardToResourceNodeFromTrafficPool.Amount.ToDec().Mul(shareOfToken).TruncateInt()

		totalUsedFromMiningPool = totalUsedFromMiningPool.Add(sdk.NewCoin(k.RewardDenom(ctx), stakeRewardFromMiningPool))
		totalUsedFromTrafficPool = totalUsedFromTrafficPool.Add(sdk.NewCoin(k.BondDenom(ctx), stakeRewardFromTrafficPool))

		if _, ok := rewardDetailMap[walletAddr.String()]; !ok {
			reward := types.NewDefaultReward(walletAddr)
			rewardDetailMap[walletAddr.String()] = reward
		}

		newReward := rewardDetailMap[walletAddr.String()]
		newReward = newReward.AddRewardFromMiningPool(sdk.NewCoin(k.RewardDenom(ctx), stakeRewardFromMiningPool))
		newReward = newReward.AddRewardFromTrafficPool(sdk.NewCoin(k.BondDenom(ctx), stakeRewardFromTrafficPool))
		rewardDetailMap[walletAddr.String()] = newReward

	}
	// deduct used reward from distributeGoal
	distributeGoal.BlockChainRewardToResourceNodeFromMiningPool =
		distributeGoal.BlockChainRewardToResourceNodeFromMiningPool.Sub(totalUsedFromMiningPool)
	distributeGoal.BlockChainRewardToResourceNodeFromTrafficPool =
		distributeGoal.BlockChainRewardToResourceNodeFromTrafficPool.Sub(totalUsedFromTrafficPool)

	totalUsedFromMiningPool = sdk.NewCoin(k.RewardDenom(ctx), sdk.ZeroInt())
	totalUsedFromTrafficPool = sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())

	// calc total traffic
	totalConsumedOzone := k.GetTotalConsumedUoz(trafficList)
	// 2, calc traffic reward
	for _, walletTraffic := range trafficList {
		walletAddr := walletTraffic.WalletAddress
		trafficVolume := walletTraffic.Volume

		shareOfTraffic := trafficVolume.ToDec().Quo(totalConsumedOzone.ToDec())
		trafficRewardFromMiningPool :=
			sdk.NewCoin(k.RewardDenom(ctx), distributeGoal.TrafficRewardToResourceNodeFromMiningPool.Amount.ToDec().Mul(shareOfTraffic).TruncateInt())
		trafficRewardFromTrafficPool :=
			sdk.NewCoin(k.BondDenom(ctx), distributeGoal.TrafficRewardToResourceNodeFromTrafficPool.Amount.ToDec().Mul(shareOfTraffic).TruncateInt())

		totalUsedFromMiningPool = totalUsedFromMiningPool.Add(trafficRewardFromMiningPool)
		totalUsedFromTrafficPool = totalUsedFromTrafficPool.Add(trafficRewardFromTrafficPool)

		if _, ok := rewardDetailMap[walletAddr.String()]; !ok {
			reward := types.NewDefaultReward(walletAddr)
			rewardDetailMap[walletAddr.String()] = reward
		}

		newReward := rewardDetailMap[walletAddr.String()]
		newReward = newReward.AddRewardFromMiningPool(trafficRewardFromMiningPool)
		newReward = newReward.AddRewardFromTrafficPool(trafficRewardFromTrafficPool)
		rewardDetailMap[walletAddr.String()] = newReward
	}
	// deduct used reward from distributeGoal
	distributeGoal.TrafficRewardToResourceNodeFromMiningPool =
		distributeGoal.TrafficRewardToResourceNodeFromMiningPool.Sub(totalUsedFromMiningPool)
	distributeGoal.TrafficRewardToResourceNodeFromTrafficPool =
		distributeGoal.TrafficRewardToResourceNodeFromTrafficPool.Sub(totalUsedFromTrafficPool)

	return rewardDetailMap, distributeGoal
}

func (k Keeper) CalcRewardForIndexingNodeForTestnet(ctx sdk.Context, distributeGoal types.DistributeGoal,
	rewardDetailMap map[string]types.Reward, indexNodes []regtypes.IndexingNode,
) (map[string]types.Reward, types.DistributeGoal) {

	totalUsedStakeRewardFromMiningPool := sdk.NewCoin(k.RewardDenom(ctx), sdk.ZeroInt())
	totalUsedStakeRewardFromTrafficPool := sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())
	totalUsedIndexingRewardFromMiningPool := sdk.NewCoin(k.RewardDenom(ctx), sdk.ZeroInt())
	totalUsedIndexingRewardFromTrafficPool := sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())

	totalStakeOfIndexingNodes := k.RegisterKeeper.GetIndexingNodeBondedToken(ctx).Amount
	indexingNodeCnt := sdk.NewDec(int64(len(indexNodes)))
	for _, node := range indexNodes {
		walletAddr := node.GetOwnerAddr()

		// 1, calc stake reward
		shareOfToken := node.GetTokens().ToDec().Quo(totalStakeOfIndexingNodes.ToDec())
		stakeRewardFromMiningPool :=
			sdk.NewCoin(k.RewardDenom(ctx), distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool.Amount.ToDec().Quo(indexingNodeCnt).TruncateInt())
		stakeRewardFromTrafficPool :=
			sdk.NewCoin(k.BondDenom(ctx), distributeGoal.BlockChainRewardToIndexingNodeFromTrafficPool.Amount.ToDec().Mul(shareOfToken).TruncateInt())

		totalUsedStakeRewardFromMiningPool = totalUsedStakeRewardFromMiningPool.Add(stakeRewardFromMiningPool)
		totalUsedStakeRewardFromTrafficPool = totalUsedStakeRewardFromTrafficPool.Add(stakeRewardFromTrafficPool)

		// 2, calc indexing reward
		indexingRewardFromMiningPool :=
			sdk.NewCoin(k.RewardDenom(ctx), distributeGoal.MetaNodeRewardToIndexingNodeFromMiningPool.Amount.ToDec().Quo(indexingNodeCnt).TruncateInt())
		indexingRewardFromTrafficPool :=
			sdk.NewCoin(k.BondDenom(ctx), distributeGoal.MetaNodeRewardToIndexingNodeFromTrafficPool.Amount.ToDec().Quo(indexingNodeCnt).TruncateInt())

		totalUsedIndexingRewardFromMiningPool = totalUsedIndexingRewardFromMiningPool.Add(indexingRewardFromMiningPool)
		totalUsedIndexingRewardFromTrafficPool = totalUsedIndexingRewardFromTrafficPool.Add(indexingRewardFromTrafficPool)

		if _, ok := rewardDetailMap[walletAddr.String()]; !ok {
			reward := types.NewDefaultReward(walletAddr)
			rewardDetailMap[walletAddr.String()] = reward
		}

		newReward := rewardDetailMap[walletAddr.String()]
		newReward = newReward.AddRewardFromMiningPool(stakeRewardFromMiningPool.Add(indexingRewardFromMiningPool))
		newReward = newReward.AddRewardFromTrafficPool(stakeRewardFromTrafficPool.Add(indexingRewardFromTrafficPool))
		rewardDetailMap[walletAddr.String()] = newReward
	}
	// deduct used reward from distributeGoal
	distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool =
		distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool.Sub(totalUsedStakeRewardFromMiningPool)
	distributeGoal.BlockChainRewardToIndexingNodeFromTrafficPool =
		distributeGoal.BlockChainRewardToIndexingNodeFromTrafficPool.Sub(totalUsedStakeRewardFromTrafficPool)
	distributeGoal.MetaNodeRewardToIndexingNodeFromMiningPool =
		distributeGoal.MetaNodeRewardToIndexingNodeFromMiningPool.Sub(totalUsedIndexingRewardFromMiningPool)
	distributeGoal.MetaNodeRewardToIndexingNodeFromTrafficPool =
		distributeGoal.MetaNodeRewardToIndexingNodeFromTrafficPool.Sub(totalUsedIndexingRewardFromTrafficPool)

	return rewardDetailMap, distributeGoal
}

// move reward to fee pool for validator traffic reward distribution
func (k Keeper) distributeValidatorRewardForTestnet(ctx sdk.Context, distributeGoal types.DistributeGoal) (types.DistributeGoal, error) {
	rewardFromMiningPool := distributeGoal.BlockChainRewardToValidatorFromMiningPool
	rewardFromTrafficPool := distributeGoal.BlockChainRewardToValidatorFromTrafficPool
	usedRewardFromMiningPool := sdk.NewCoin(k.RewardDenom(ctx), sdk.ZeroInt())

	// distribute reward from mining pool to wallet directly
	validatorWalletList := make([]sdk.AccAddress, 0)
	validators := k.StakingKeeper.GetAllValidators(ctx)
	for _, validator := range validators {
		if validator.IsBonded() && !validator.IsJailed() {
			validatorWalletList = append(validatorWalletList, sdk.AccAddress(validator.GetOperator()))
		}
	}
	rewardPerValidator := sdk.NewCoin(k.RewardDenom(ctx), rewardFromMiningPool.Amount.ToDec().Quo(sdk.NewDec(int64(len(validatorWalletList)))).TruncateInt())
	for _, validatorWallet := range validatorWalletList {
		_, err := k.BankKeeper.AddCoins(ctx, validatorWallet, sdk.NewCoins(rewardPerValidator))
		if err != nil {
			return distributeGoal, err
		}
		usedRewardFromMiningPool = usedRewardFromMiningPool.Add(rewardPerValidator)
	}

	// distribute rewards from traffic pool to fee_pool
	feePoolAccAddr := k.SupplyKeeper.GetModuleAddress(k.feeCollectorName)

	if feePoolAccAddr == nil {
		ctx.Logger().Error("account address of distribution module does not exist.")
		return distributeGoal, types.ErrUnknownAccountAddress
	}

	_, err := k.BankKeeper.AddCoins(ctx, feePoolAccAddr, sdk.NewCoins(rewardFromTrafficPool))
	if err != nil {
		return distributeGoal, err
	}

	distributeGoal.BlockChainRewardToValidatorFromMiningPool = rewardFromMiningPool.Sub(usedRewardFromMiningPool)
	distributeGoal.BlockChainRewardToValidatorFromTrafficPool = sdk.Coin{}

	return distributeGoal, nil
}

func (k Keeper) returnBalanceForTestnet(ctx sdk.Context, goal types.DistributeGoal, epoch sdk.Int) (err error) {
	balanceOfMiningPool := goal.BlockChainRewardToIndexingNodeFromMiningPool.
		Add(goal.MetaNodeRewardToIndexingNodeFromMiningPool).
		Add(goal.BlockChainRewardToResourceNodeFromMiningPool).
		Add(goal.TrafficRewardToResourceNodeFromMiningPool).
		Add(goal.BlockChainRewardToValidatorFromMiningPool)
	balanceOfTrafficPool := goal.BlockChainRewardToIndexingNodeFromTrafficPool.
		Add(goal.MetaNodeRewardToIndexingNodeFromTrafficPool).
		Add(goal.BlockChainRewardToResourceNodeFromTrafficPool).
		Add(goal.TrafficRewardToResourceNodeFromTrafficPool)

	// return balance to foundation account
	foundationAccountAddr := k.SupplyKeeper.GetModuleAddress(types.FoundationAccount)
	if foundationAccountAddr == nil {
		ctx.Logger().Error("foundation account address of distribution module does not exist.")
		return types.ErrUnknownAccountAddress
	}
	amountToAdd := sdk.NewCoins(balanceOfMiningPool)
	_, err = k.BankKeeper.AddCoins(ctx, foundationAccountAddr, amountToAdd)
	if err != nil {
		return err
	}

	//return balance to minedToken record
	oldTotalMinedToken := k.GetTotalMinedTokens(ctx)
	newTotalMinedToken := oldTotalMinedToken.Sub(balanceOfMiningPool)
	oldMinedToken := k.GetMinedTokens(ctx, epoch)
	newMinedToken := oldMinedToken.Sub(balanceOfMiningPool)
	k.SetTotalMinedTokens(ctx, newTotalMinedToken)
	k.setMinedTokens(ctx, epoch, newMinedToken)

	// return balance to prepay pool
	totalUnIssuedPrepay := k.RegisterKeeper.GetTotalUnissuedPrepay(ctx)
	newTotalUnIssuedPrePay := totalUnIssuedPrepay.Add(balanceOfTrafficPool)
	k.RegisterKeeper.SetTotalUnissuedPrepay(ctx, newTotalUnIssuedPrePay)

	return nil
}
