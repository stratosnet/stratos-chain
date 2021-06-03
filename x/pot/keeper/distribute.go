package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func (k Keeper) DistributePotReward(ctx sdk.Context, trafficList []types.SingleNodeVolume, epoch sdk.Int) (err error) {
	distributeGoal := types.InitDistributeGoal()
	rewardDetailMap := make(map[string]types.Reward) //key: node address

	//1, calc traffic reward in total
	err = k.calcTrafficRewardInTotal(ctx, trafficList, distributeGoal)
	if err != nil {
		return err
	}
	//2, calc mining reward in total
	err = k.calcMiningRewardInTotal(ctx, distributeGoal)
	if err != nil {
		return err
	}
	//3, calc reward for resource node
	rewardDetailMap, distributeGoal = k.calcRewardForResourceNode(ctx, trafficList, distributeGoal, rewardDetailMap)
	//4, calc reward from indexing node
	rewardDetailMap, distributeGoal = k.calcRewardForIndexingNode(ctx, distributeGoal, rewardDetailMap)

	//5, deduct reward from provider account
	err = k.deductRewardFromRewardProviderAccount(ctx, distributeGoal)
	if err != nil {
		return err
	}

	//6, distribute skate reward to fee pool for validators
	_, err = k.distributeRewardToFeePool(ctx, distributeGoal)
	if err != nil {
		return err
	}

	//7, distribute all rewards to resource nodes & indexing nodes
	err = k.distributeRewardToSdsNodes(ctx, rewardDetailMap, epoch.Int64())
	if err != nil {
		return err
	}

	//8, return balance to traffic pool & mining pool
	err = k.returnBalance(ctx, distributeGoal)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) deductRewardFromRewardProviderAccount(ctx sdk.Context, goal types.DistributeGoal) (err error) {
	totalRewardFromMiningPool := goal.BlockChainRewardToIndexingNodeFromMiningPool.Add(goal.MetaNodeRewardToIndexingNodeFromMiningPool).Add(goal.BlockChainRewardToResourceNodeFromMiningPool).Add(goal.TrafficRewardToResourceNodeFromMiningPool)
	totalRewardFromTrafficPool := goal.BlockChainRewardToIndexingNodeFromTrafficPool.Add(goal.MetaNodeRewardToIndexingNodeFromTrafficPool).Add(goal.BlockChainRewardToResourceNodeFromTrafficPool).Add(goal.TrafficRewardToResourceNodeFromTrafficPool)

	// deduct mining reward from foundation account
	foundationAccountAddr := k.GetFoundationAccount(ctx)
	amountToDeduct := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), totalRewardFromMiningPool))
	hasCoin := k.bankKeeper.HasCoins(ctx, foundationAccountAddr, amountToDeduct)
	if !hasCoin {
		return types.ErrInsufficientBalance
	}
	_, err = k.bankKeeper.SubtractCoins(ctx, foundationAccountAddr, amountToDeduct)
	if err != nil {
		return err
	}

	// deduct traffic reward from prepay pool
	totalUnIssuedPrepay := k.getTotalUnissuedPrepay(ctx)
	newTotalUnIssuedPrePay := totalUnIssuedPrepay.Sub(totalRewardFromTrafficPool)
	if newTotalUnIssuedPrePay.IsNegative() {
		return types.ErrInsufficientBalance
	}
	k.setTotalUnissuedPrepay(ctx, newTotalUnIssuedPrePay)

	return nil
}

func (k Keeper) returnBalance(ctx sdk.Context, goal types.DistributeGoal) (err error) {
	balanceOfMiningPool := goal.BlockChainRewardToIndexingNodeFromMiningPool.Add(goal.MetaNodeRewardToIndexingNodeFromMiningPool).Add(goal.BlockChainRewardToResourceNodeFromMiningPool).Add(goal.TrafficRewardToResourceNodeFromMiningPool)
	balanceOfTrafficPool := goal.BlockChainRewardToIndexingNodeFromTrafficPool.Add(goal.MetaNodeRewardToIndexingNodeFromTrafficPool).Add(goal.BlockChainRewardToResourceNodeFromTrafficPool).Add(goal.TrafficRewardToResourceNodeFromTrafficPool)

	// return balance to foundation account
	foundationAccountAddr := k.GetFoundationAccount(ctx)
	amountToAdd := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), balanceOfMiningPool))
	_, err = k.bankKeeper.AddCoins(ctx, foundationAccountAddr, amountToAdd)
	if err != nil {
		return err
	}

	// return balance to prepay pool
	totalUnIssuedPrepay := k.getTotalUnissuedPrepay(ctx)
	newTotalUnIssuedPrePay := totalUnIssuedPrepay.Add(balanceOfTrafficPool)
	k.setTotalUnissuedPrepay(ctx, newTotalUnIssuedPrePay)

	return nil
}

func (k Keeper) calcTrafficRewardInTotal(ctx sdk.Context, trafficList []types.SingleNodeVolume, distributeGoal types.DistributeGoal) error {
	totalTrafficReward := k.getTrafficReward(ctx, trafficList)
	minedTokens := k.getMinedTokens(ctx)
	miningParam, err := k.GetMiningRewardParamByMinedToken(ctx, minedTokens)
	if err != nil {
		return err
	}
	stakeReward := totalTrafficReward.Mul(miningParam.BlockChainPercentage).TruncateInt()
	trafficReward := totalTrafficReward.Mul(miningParam.ResourceNodePercentage).TruncateInt()
	indexingReward := totalTrafficReward.Mul(miningParam.MetaNodePercentage).TruncateInt()

	stakeRewardToValidators, stakeRewardToResourceNodes, stakeRewardToIndexingNodes := k.splitRewardByStake(ctx, stakeReward)
	distributeGoal.AddBlockChainRewardToValidatorFromTrafficPool(stakeRewardToValidators)
	distributeGoal.AddBlockChainRewardToResourceNodeFromTrafficPool(stakeRewardToResourceNodes)
	distributeGoal.AddBlockChainRewardToIndexingNodeFromTrafficPool(stakeRewardToIndexingNodes)
	distributeGoal.AddTrafficRewardToResourceNodeFromTrafficPool(trafficReward)
	distributeGoal.AddMetaNodeRewardToIndexingNodeFromTrafficPool(indexingReward)
	return nil
}

// [S] is initial genesis deposit by all resource nodes and meta nodes at t=0
// The current unissued prepay Volume Pool [pt] is the total remaining prepay STOS kept by Stratos Network but not issued to Resource Node as rewards. At time t=0,  pt=0
// total consumed Ozone is [Y]
// The remaining total Ozone limit [lt] is the upper bound of total Ozone that users can purchase from Stratos blockchain.
// the total generated traffic rewards as [R]
// R = (S + Pt) * Y / (Lt + Y)
func (k Keeper) getTrafficReward(ctx sdk.Context, trafficList []types.SingleNodeVolume) sdk.Dec {
	S := k.registerKeeper.GetInitialGenesisStakeTotal(ctx).ToDec()
	Pt := k.getTotalUnissuedPrepay(ctx).ToDec()
	Y := k.getTotalConsumedOzone(trafficList).ToDec()
	Lt := k.registerKeeper.GetUpperBoundOfTotalOzone(ctx).ToDec()
	R := S.Add(Pt).Mul(Y).Quo(Lt.Add(Y))

	return R
}

// allocate mining reward from foundation account
func (k Keeper) calcMiningRewardInTotal(ctx sdk.Context, distributeGoal types.DistributeGoal) error {
	minedTokens := k.getMinedTokens(ctx)
	miningParam, err := k.GetMiningRewardParamByMinedToken(ctx, minedTokens)

	totalMiningReward := miningParam.MiningReward
	if err != nil {
		return err
	}
	stakeReward := totalMiningReward.Mul(miningParam.BlockChainPercentage).TruncateInt()
	trafficReward := totalMiningReward.Mul(miningParam.ResourceNodePercentage).TruncateInt()
	indexingReward := totalMiningReward.Mul(miningParam.MetaNodePercentage).TruncateInt()

	stakeRewardToValidators, stakeRewardToResourceNodes, stakeRewardToIndexingNodes := k.splitRewardByStake(ctx, stakeReward)
	distributeGoal.AddBlockChainRewardToValidatorFromMiningPool(stakeRewardToValidators)
	distributeGoal.AddBlockChainRewardToResourceNodeFromMiningPool(stakeRewardToResourceNodes)
	distributeGoal.AddBlockChainRewardToIndexingNodeFromMiningPool(stakeRewardToIndexingNodes)
	distributeGoal.AddTrafficRewardToResourceNodeFromMiningPool(trafficReward)
	distributeGoal.AddMetaNodeRewardToIndexingNodeFromMiningPool(indexingReward)
	return nil
}

func (k Keeper) distributeRewardToSdsNodes(ctx sdk.Context, rewardDetailMap map[string]types.Reward, currentEpoch int64) (err error) {
	for _, reward := range rewardDetailMap {
		nodeAddr := reward.NodeAddress
		totalReward := reward.RewardFromMiningPool.Add(reward.RewardFromTrafficPool)
		matureEpoch := k.getMatureEpoch(currentEpoch)
		k.addNewRewardAndReCalcTotal(ctx, nodeAddr, currentEpoch, matureEpoch, totalReward)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) addNewRewardAndReCalcTotal(ctx sdk.Context, account sdk.AccAddress, currentEpoch int64, matureEpoch int64, newReward sdk.Int) {
	oldMatureTotal := k.getMatureTotalReward(ctx, account)
	oldImmatureTotal := k.getImmatureTotalReward(ctx, account)
	matureStartEpoch := k.getLastMaturedEpoch(ctx) + 1

	immatureToMature := sdk.ZeroInt()
	for i := matureStartEpoch; i <= currentEpoch; i++ {
		reward := k.getIndividualReward(ctx, account, i)
		immatureToMature = immatureToMature.Add(reward)
	}

	matureTotal := oldMatureTotal.Add(immatureToMature)
	immatureTotal := oldImmatureTotal.Sub(immatureToMature).Add(newReward)

	k.setLastMaturedEpoch(ctx, currentEpoch)
	k.setMatureTotalReward(ctx, account, matureTotal)
	k.setImmatureTotalReward(ctx, account, immatureTotal)
	k.setIndividualReward(ctx, account, matureEpoch, newReward)
}

// reward will mature 14 days since distribution. Each epoch interval is about 10 minutes.
func (k Keeper) getMatureEpoch(currentEpoch int64) (matureEpoch int64) {
	// 14 days = 20160 minutes = 2016 epochs
	return currentEpoch + 2016
}

// move reward to fee pool for validator traffic reward distribution
func (k Keeper) distributeRewardToFeePool(ctx sdk.Context, distributeGoal types.DistributeGoal) (distributed sdk.Int, err error) {
	rewardFromMiningPool := distributeGoal.BlockChainRewardToValidatorFromMiningPool
	rewardFromTrafficPool := distributeGoal.BlockChainRewardToValidatorFromTrafficPool
	totalRewardSendToFeePool := rewardFromMiningPool.Add(rewardFromTrafficPool)

	feePoolAccAddr := k.supplyKeeper.GetModuleAddress(k.feeCollectorName)

	if feePoolAccAddr == nil {
		ctx.Logger().Error("account address of distribution module does not exist.")
		return sdk.ZeroInt(), types.ErrUnknownAccountAddress
	}

	_, err = k.bankKeeper.AddCoins(ctx, feePoolAccAddr, sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), totalRewardSendToFeePool)))
	if err != nil {
		return sdk.ZeroInt(), err
	}

	distributeGoal.BlockChainRewardToValidatorFromMiningPool = sdk.ZeroInt()
	distributeGoal.BlockChainRewardToValidatorFromTrafficPool = sdk.ZeroInt()

	return totalRewardSendToFeePool, nil
}

func (k Keeper) calcRewardForResourceNode(ctx sdk.Context, trafficList []types.SingleNodeVolume,
	distributeGoal types.DistributeGoal, rewardDetailMap map[string]types.Reward,
) (map[string]types.Reward, types.DistributeGoal) {

	var totalUsedFromMiningPool sdk.Int
	var totalUsedFromTrafficPool sdk.Int

	totalUsedFromMiningPool = sdk.ZeroInt()
	totalUsedFromTrafficPool = sdk.ZeroInt()

	// 1, calc stake reward
	totalStakeOfResourceNodes := k.registerKeeper.GetLastResourceNodeTotalStake(ctx)
	resourceNodeList := k.registerKeeper.GetAllResourceNodes(ctx)
	for _, node := range resourceNodeList {
		nodeAddr := node.GetNetworkAddr()

		shareOfToken := node.GetTokens().ToDec().Quo(totalStakeOfResourceNodes.ToDec())
		stakeRewardFromMiningPool := distributeGoal.BlockChainRewardToResourceNodeFromMiningPool.ToDec().Mul(shareOfToken).TruncateInt()
		stakeRewardFromTrafficPool := distributeGoal.BlockChainRewardToResourceNodeFromTrafficPool.ToDec().Mul(shareOfToken).TruncateInt()

		totalUsedFromMiningPool = totalUsedFromMiningPool.Add(stakeRewardFromMiningPool)
		totalUsedFromTrafficPool = totalUsedFromTrafficPool.Add(stakeRewardFromTrafficPool)

		if _, ok := rewardDetailMap[nodeAddr.String()]; !ok {
			reward := types.NewDefaultReward(nodeAddr)
			rewardDetailMap[nodeAddr.String()] = reward
		}

		newReward := rewardDetailMap[nodeAddr.String()]
		newReward.AddRewardFromMiningPool(stakeRewardFromMiningPool)
		newReward.AddRewardFromTrafficPool(stakeRewardFromTrafficPool)
		rewardDetailMap[nodeAddr.String()] = newReward

	}
	// deduct used reward from distributeGoal
	distributeGoal.BlockChainRewardToResourceNodeFromMiningPool = distributeGoal.BlockChainRewardToResourceNodeFromMiningPool.Sub(totalUsedFromMiningPool)
	distributeGoal.BlockChainRewardToResourceNodeFromTrafficPool = distributeGoal.BlockChainRewardToResourceNodeFromTrafficPool.Sub(totalUsedFromTrafficPool)

	totalUsedFromMiningPool = sdk.ZeroInt()
	totalUsedFromTrafficPool = sdk.ZeroInt()

	// calc total traffic
	totalConsumedOzone := k.getTotalConsumedOzone(trafficList)
	// 2, calc traffic reward
	for _, nodeTraffic := range trafficList {
		nodeAddr := nodeTraffic.NodeAddress
		nodeTraffic := nodeTraffic.Volume

		shareOfTraffic := nodeTraffic.ToDec().Quo(totalConsumedOzone.ToDec())
		trafficRewardFromMiningPool := distributeGoal.TrafficRewardToResourceNodeFromMiningPool.ToDec().Mul(shareOfTraffic).TruncateInt()
		trafficRewardFromTrafficPool := distributeGoal.TrafficRewardToResourceNodeFromTrafficPool.ToDec().Mul(shareOfTraffic).TruncateInt()

		totalUsedFromMiningPool = totalUsedFromMiningPool.Add(trafficRewardFromMiningPool)
		totalUsedFromTrafficPool = totalUsedFromTrafficPool.Add(trafficRewardFromTrafficPool)

		if _, ok := rewardDetailMap[nodeAddr.String()]; !ok {
			reward := types.NewDefaultReward(nodeAddr)
			rewardDetailMap[nodeAddr.String()] = reward
		}

		newReward := rewardDetailMap[nodeAddr.String()]
		newReward.AddRewardFromMiningPool(trafficRewardFromMiningPool)
		newReward.AddRewardFromTrafficPool(trafficRewardFromTrafficPool)
		rewardDetailMap[nodeAddr.String()] = newReward
	}
	// deduct used reward from distributeGoal
	distributeGoal.TrafficRewardToResourceNodeFromMiningPool = distributeGoal.TrafficRewardToResourceNodeFromMiningPool.Sub(totalUsedFromMiningPool)
	distributeGoal.TrafficRewardToResourceNodeFromTrafficPool = distributeGoal.TrafficRewardToResourceNodeFromTrafficPool.Sub(totalUsedFromTrafficPool)

	return rewardDetailMap, distributeGoal
}

func (k Keeper) calcRewardForIndexingNode(ctx sdk.Context, distributeGoal types.DistributeGoal, rewardDetailMap map[string]types.Reward,
) (map[string]types.Reward, types.DistributeGoal) {

	totalUsedStakeRewardFromMiningPool := sdk.ZeroInt()
	totalUsedStakeRewardFromTrafficPool := sdk.ZeroInt()
	totalUsedIndexingRewardFromMiningPool := sdk.ZeroInt()
	totalUsedIndexingRewardFromTrafficPool := sdk.ZeroInt()

	totalStakeOfIndexingNodes := k.registerKeeper.GetLastIndexingNodeTotalStake(ctx)
	indexingNodeList := k.registerKeeper.GetAllIndexingNodes(ctx)
	indexingNodeCnt := sdk.NewDec(int64(len(indexingNodeList)))
	for _, node := range indexingNodeList {
		nodeAddr := node.GetNetworkAddr()

		// 1, calc stake reward
		shareOfToken := node.GetTokens().ToDec().Quo(totalStakeOfIndexingNodes.ToDec())
		stakeRewardFromMiningPool := distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool.ToDec().Mul(shareOfToken).TruncateInt()
		stakeRewardFromTrafficPool := distributeGoal.BlockChainRewardToIndexingNodeFromTrafficPool.ToDec().Mul(shareOfToken).TruncateInt()

		totalUsedStakeRewardFromMiningPool = totalUsedStakeRewardFromMiningPool.Add(stakeRewardFromMiningPool)
		totalUsedStakeRewardFromTrafficPool = totalUsedStakeRewardFromTrafficPool.Add(stakeRewardFromTrafficPool)

		// 2, calc indexing reward
		indexingRewardFromMiningPool := distributeGoal.MetaNodeRewardToIndexingNodeFromMiningPool.ToDec().Quo(indexingNodeCnt).TruncateInt()
		indexingRewardFromTrafficPool := distributeGoal.MetaNodeRewardToIndexingNodeFromTrafficPool.ToDec().Quo(indexingNodeCnt).TruncateInt()

		totalUsedIndexingRewardFromMiningPool = totalUsedIndexingRewardFromMiningPool.Add(indexingRewardFromMiningPool)
		totalUsedIndexingRewardFromTrafficPool = totalUsedIndexingRewardFromTrafficPool.Add(indexingRewardFromTrafficPool)

		if _, ok := rewardDetailMap[nodeAddr.String()]; !ok {
			reward := types.NewDefaultReward(nodeAddr)
			rewardDetailMap[nodeAddr.String()] = reward
		}

		newReward := rewardDetailMap[nodeAddr.String()]
		newReward.AddRewardFromMiningPool(stakeRewardFromMiningPool)
		newReward.AddRewardFromTrafficPool(stakeRewardFromTrafficPool)
		rewardDetailMap[nodeAddr.String()] = newReward
	}
	// deduct used reward from distributeGoal
	distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool = distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool.Sub(totalUsedStakeRewardFromMiningPool)
	distributeGoal.BlockChainRewardToIndexingNodeFromTrafficPool = distributeGoal.BlockChainRewardToIndexingNodeFromTrafficPool.Sub(totalUsedStakeRewardFromTrafficPool)
	distributeGoal.MetaNodeRewardToIndexingNodeFromMiningPool = distributeGoal.MetaNodeRewardToIndexingNodeFromMiningPool.Sub(totalUsedIndexingRewardFromMiningPool)
	distributeGoal.MetaNodeRewardToIndexingNodeFromTrafficPool = distributeGoal.MetaNodeRewardToIndexingNodeFromTrafficPool.Sub(totalUsedIndexingRewardFromTrafficPool)

	return rewardDetailMap, distributeGoal
}

func (k Keeper) getTotalConsumedOzone(trafficList []types.SingleNodeVolume) sdk.Int {
	totalTraffic := sdk.ZeroInt()
	for _, vol := range trafficList {
		totalTraffic = totalTraffic.Add(vol.Volume)
	}
	return totalTraffic
}

func (k Keeper) splitRewardByStake(ctx sdk.Context, totalReward sdk.Int,
) (validatorReward sdk.Int, resourceNodeReward sdk.Int, indexingNodeReward sdk.Int) {

	validatorBondedTokens := k.stakingKeeper.TotalBondedTokens(ctx).ToDec()
	resourceNodeBondedTokens := k.registerKeeper.GetLastResourceNodeTotalStake(ctx).ToDec()
	indexingNodeBondedTokens := k.registerKeeper.GetLastIndexingNodeTotalStake(ctx).ToDec()
	totalBondedTokens := validatorBondedTokens.Add(resourceNodeBondedTokens).Add(indexingNodeBondedTokens)

	validatorReward = totalReward.ToDec().Mul(validatorBondedTokens.Quo(totalBondedTokens)).TruncateInt()
	resourceNodeReward = totalReward.ToDec().Mul(resourceNodeBondedTokens.Quo(totalBondedTokens)).TruncateInt()
	indexingNodeReward = totalReward.ToDec().Mul(indexingNodeBondedTokens.Quo(totalBondedTokens)).TruncateInt()

	return
}
