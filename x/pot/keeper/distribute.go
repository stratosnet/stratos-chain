package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func (k Keeper) DistributePotReward(ctx sdk.Context, trafficList []types.SingleNodeVolume, epoch sdk.Int) (err error) {
	distributeGoal := types.InitDistributeGoal()
	rewardDetailMap := make(map[string]types.Reward) //key: node address

	//1, calc traffic reward in total
	distributeGoal, err = k.calcTrafficRewardInTotal(ctx, trafficList, distributeGoal)
	if err != nil {
		return err
	}

	//2, calc mining reward in total
	distributeGoal, err = k.calcMiningRewardInTotal(ctx, distributeGoal)
	if err != nil && err != types.ErrOutOfIssuance {
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
	distributeGoal, err = k.distributeValidatorRewardToFeePool(ctx, distributeGoal)
	if err != nil {
		return err
	}

	//7, distribute all rewards to resource nodes & indexing nodes
	err = k.distributeRewardToSdsNodes(ctx, rewardDetailMap, epoch)
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
	totalRewardFromMiningPool := goal.BlockChainRewardToIndexingNodeFromMiningPool.
		Add(goal.MetaNodeRewardToIndexingNodeFromMiningPool).
		Add(goal.BlockChainRewardToResourceNodeFromMiningPool).
		Add(goal.TrafficRewardToResourceNodeFromMiningPool)
	totalRewardFromTrafficPool := goal.BlockChainRewardToIndexingNodeFromTrafficPool.
		Add(goal.MetaNodeRewardToIndexingNodeFromTrafficPool).
		Add(goal.BlockChainRewardToResourceNodeFromTrafficPool).
		Add(goal.TrafficRewardToResourceNodeFromTrafficPool)

	// deduct mining reward from foundation account
	foundationAccountAddr := k.GetFoundationAccount(ctx)
	amountToDeduct := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), totalRewardFromMiningPool))
	hasCoin := k.BankKeeper.HasCoins(ctx, foundationAccountAddr, amountToDeduct)
	if !hasCoin {
		return types.ErrInsufficientBalance
	}
	_, err = k.BankKeeper.SubtractCoins(ctx, foundationAccountAddr, amountToDeduct)
	if err != nil {
		return err
	}

	// update mined token record by adding mining reward
	oldMinedToken := k.getMinedTokens(ctx)
	newMinedToken := oldMinedToken.Add(totalRewardFromMiningPool)
	k.setMinedTokens(ctx, newMinedToken)

	// deduct traffic reward from prepay pool
	totalUnIssuedPrepay := k.GetTotalUnissuedPrepay(ctx)
	newTotalUnIssuedPrePay := totalUnIssuedPrepay.Sub(totalRewardFromTrafficPool)
	if newTotalUnIssuedPrePay.IsNegative() {
		return types.ErrInsufficientBalance
	}
	k.SetTotalUnissuedPrepay(ctx, newTotalUnIssuedPrePay)

	return nil
}

func (k Keeper) returnBalance(ctx sdk.Context, goal types.DistributeGoal) (err error) {
	balanceOfMiningPool := goal.BlockChainRewardToIndexingNodeFromMiningPool.
		Add(goal.MetaNodeRewardToIndexingNodeFromMiningPool).
		Add(goal.BlockChainRewardToResourceNodeFromMiningPool).
		Add(goal.TrafficRewardToResourceNodeFromMiningPool)
	balanceOfTrafficPool := goal.BlockChainRewardToIndexingNodeFromTrafficPool.
		Add(goal.MetaNodeRewardToIndexingNodeFromTrafficPool).
		Add(goal.BlockChainRewardToResourceNodeFromTrafficPool).
		Add(goal.TrafficRewardToResourceNodeFromTrafficPool)

	// return balance to foundation account
	foundationAccountAddr := k.GetFoundationAccount(ctx)
	amountToAdd := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), balanceOfMiningPool))
	_, err = k.BankKeeper.AddCoins(ctx, foundationAccountAddr, amountToAdd)
	if err != nil {
		return err
	}

	//return balance to minedToken record
	oldMinedToken := k.getMinedTokens(ctx)
	newMinedToken := oldMinedToken.Sub(balanceOfMiningPool)
	k.setMinedTokens(ctx, newMinedToken)

	// return balance to prepay pool
	totalUnIssuedPrepay := k.GetTotalUnissuedPrepay(ctx)
	newTotalUnIssuedPrePay := totalUnIssuedPrepay.Add(balanceOfTrafficPool)
	k.SetTotalUnissuedPrepay(ctx, newTotalUnIssuedPrePay)

	return nil
}

func (k Keeper) calcTrafficRewardInTotal(
	ctx sdk.Context, trafficList []types.SingleNodeVolume, distributeGoal types.DistributeGoal,
) (types.DistributeGoal, error) {

	totalTrafficReward := k.getTrafficReward(ctx, trafficList)
	minedTokens := k.getMinedTokens(ctx)
	miningParam, err := k.GetMiningRewardParamByMinedToken(ctx, minedTokens)
	if err != nil && err != types.ErrOutOfIssuance {
		return distributeGoal, err
	}
	stakeReward := totalTrafficReward.
		Mul(miningParam.BlockChainPercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	trafficReward := totalTrafficReward.
		Mul(miningParam.ResourceNodePercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	indexingReward := totalTrafficReward.
		Mul(miningParam.MetaNodePercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()

	stakeRewardToValidators, stakeRewardToResourceNodes, stakeRewardToIndexingNodes := k.splitRewardByStake(ctx, stakeReward)
	distributeGoal = distributeGoal.AddBlockChainRewardToValidatorFromTrafficPool(stakeRewardToValidators)
	distributeGoal = distributeGoal.AddBlockChainRewardToResourceNodeFromTrafficPool(stakeRewardToResourceNodes)
	distributeGoal = distributeGoal.AddBlockChainRewardToIndexingNodeFromTrafficPool(stakeRewardToIndexingNodes)
	distributeGoal = distributeGoal.AddTrafficRewardToResourceNodeFromTrafficPool(trafficReward)
	distributeGoal = distributeGoal.AddMetaNodeRewardToIndexingNodeFromTrafficPool(indexingReward)

	return distributeGoal, nil
}

// [S] is initial genesis deposit by all resource nodes and meta nodes at t=0
// The current unissued prepay Volume Pool [pt] is the total remaining prepay uSTOS kept by Stratos Network but not issued to Resource Node as rewards. At time t=0,  pt=0
// total consumed Ozone is [Y]
// The remaining total Ozone limit [lt] is the upper bound of total Ozone that users can purchase from Stratos blockchain.
// the total generated traffic rewards as [R]
// R = (S + Pt) * Y / (Lt + Y)
func (k Keeper) getTrafficReward(ctx sdk.Context, trafficList []types.SingleNodeVolume) sdk.Dec {
	S := k.RegisterKeeper.GetInitialGenesisStakeTotal(ctx).ToDec()
	Pt := k.GetTotalUnissuedPrepay(ctx).ToDec()
	Y := k.GetTotalConsumedOzone(trafficList).ToDec()
	Lt := k.RegisterKeeper.GetRemainingOzoneLimit(ctx).ToDec()
	R := S.Add(Pt).Mul(Y).Quo(Lt.Add(Y))

	return R
}

// allocate mining reward from foundation account
func (k Keeper) calcMiningRewardInTotal(ctx sdk.Context, distributeGoal types.DistributeGoal) (types.DistributeGoal, error) {
	minedTokens := k.getMinedTokens(ctx)
	miningParam, err := k.GetMiningRewardParamByMinedToken(ctx, minedTokens)

	totalMiningReward := miningParam.MiningReward
	if err != nil {
		return distributeGoal, err
	}
	stakeReward := totalMiningReward.ToDec().
		Mul(miningParam.BlockChainPercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	trafficReward := totalMiningReward.ToDec().
		Mul(miningParam.ResourceNodePercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	indexingReward := totalMiningReward.ToDec().
		Mul(miningParam.MetaNodePercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()

	stakeRewardToValidators, stakeRewardToResourceNodes, stakeRewardToIndexingNodes := k.splitRewardByStake(ctx, stakeReward)
	distributeGoal = distributeGoal.AddBlockChainRewardToValidatorFromMiningPool(stakeRewardToValidators)
	distributeGoal = distributeGoal.AddBlockChainRewardToResourceNodeFromMiningPool(stakeRewardToResourceNodes)
	distributeGoal = distributeGoal.AddBlockChainRewardToIndexingNodeFromMiningPool(stakeRewardToIndexingNodes)
	distributeGoal = distributeGoal.AddTrafficRewardToResourceNodeFromMiningPool(trafficReward)
	distributeGoal = distributeGoal.AddMetaNodeRewardToIndexingNodeFromMiningPool(indexingReward)
	return distributeGoal, nil
}

func (k Keeper) distributeRewardToSdsNodes(ctx sdk.Context, rewardDetailMap map[string]types.Reward, currentEpoch sdk.Int) (err error) {
	matureEpoch := k.getMatureEpochByCurrentEpoch(ctx, currentEpoch)
	for _, reward := range rewardDetailMap {
		nodeAddr := reward.NodeAddress
		totalReward := reward.RewardFromMiningPool.Add(reward.RewardFromTrafficPool)
		k.addNewRewardAndReCalcTotal(ctx, nodeAddr, currentEpoch, matureEpoch, totalReward)
		if err != nil {
			return err
		}
	}
	k.setLastMaturedEpoch(ctx, currentEpoch)
	return nil
}

func (k Keeper) addNewRewardAndReCalcTotal(ctx sdk.Context, account sdk.AccAddress, currentEpoch sdk.Int, matureEpoch sdk.Int, newReward sdk.Int) {
	oldMatureTotal := k.getMatureTotalReward(ctx, account)
	oldImmatureTotal := k.getImmatureTotalReward(ctx, account)
	matureStartEpoch := k.getLastMaturedEpoch(ctx).Int64() + 1
	matureEndEpoch := currentEpoch.Int64()

	immatureToMature := sdk.ZeroInt()
	for i := matureStartEpoch; i <= matureEndEpoch; i++ {
		reward := k.getIndividualReward(ctx, account, sdk.NewInt(i))
		immatureToMature = immatureToMature.Add(reward)
	}

	matureTotal := oldMatureTotal.Add(immatureToMature)
	immatureTotal := oldImmatureTotal.Sub(immatureToMature).Add(newReward)

	rewardAddressPool := k.getRewardAddressPool(ctx)
	addrExist := false
	for i := 0; i < len(rewardAddressPool); i++ {
		if rewardAddressPool[i].Equals(account) {
			addrExist = true
			break
		}
	}
	if addrExist == false {
		rewardAddressPool = append(rewardAddressPool, account)
		k.setRewardAddressPool(ctx, rewardAddressPool)
	}

	k.setMatureTotalReward(ctx, account, matureTotal)
	k.setImmatureTotalReward(ctx, account, immatureTotal)
	k.setIndividualReward(ctx, account, matureEpoch, newReward)
}

// reward will mature 14 days since distribution. Each epoch interval is about 10 minutes.
func (k Keeper) getMatureEpochByCurrentEpoch(ctx sdk.Context, currentEpoch sdk.Int) (matureEpoch sdk.Int) {
	// 14 days = 20160 minutes = 2016 epochs
	matureEpoch = currentEpoch.Add(k.GetMatureEpoch(ctx))
	return matureEpoch
}

// move reward to fee pool for validator traffic reward distribution
func (k Keeper) distributeValidatorRewardToFeePool(ctx sdk.Context, distributeGoal types.DistributeGoal) (types.DistributeGoal, error) {
	rewardFromMiningPool := distributeGoal.BlockChainRewardToValidatorFromMiningPool
	rewardFromTrafficPool := distributeGoal.BlockChainRewardToValidatorFromTrafficPool
	totalRewardSendToFeePool := rewardFromMiningPool.Add(rewardFromTrafficPool)

	feePoolAccAddr := k.SupplyKeeper.GetModuleAddress(k.FeeCollectorName)

	if feePoolAccAddr == nil {
		ctx.Logger().Error("account address of distribution module does not exist.")
		return distributeGoal, types.ErrUnknownAccountAddress
	}

	_, err := k.BankKeeper.AddCoins(ctx, feePoolAccAddr, sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), totalRewardSendToFeePool)))
	if err != nil {
		return distributeGoal, err
	}

	distributeGoal.BlockChainRewardToValidatorFromMiningPool = sdk.ZeroInt()
	distributeGoal.BlockChainRewardToValidatorFromTrafficPool = sdk.ZeroInt()

	return distributeGoal, nil
}

func (k Keeper) calcRewardForResourceNode(ctx sdk.Context, trafficList []types.SingleNodeVolume,
	distributeGoal types.DistributeGoal, rewardDetailMap map[string]types.Reward,
) (map[string]types.Reward, types.DistributeGoal) {

	var totalUsedFromMiningPool sdk.Int
	var totalUsedFromTrafficPool sdk.Int

	totalUsedFromMiningPool = sdk.ZeroInt()
	totalUsedFromTrafficPool = sdk.ZeroInt()

	// 1, calc stake reward
	totalStakeOfResourceNodes := k.RegisterKeeper.GetLastResourceNodeTotalStake(ctx)
	resourceNodeList := k.RegisterKeeper.GetAllResourceNodes(ctx)
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
		newReward = newReward.AddRewardFromMiningPool(stakeRewardFromMiningPool)
		newReward = newReward.AddRewardFromTrafficPool(stakeRewardFromTrafficPool)
		rewardDetailMap[nodeAddr.String()] = newReward

	}
	// deduct used reward from distributeGoal
	distributeGoal.BlockChainRewardToResourceNodeFromMiningPool =
		distributeGoal.BlockChainRewardToResourceNodeFromMiningPool.Sub(totalUsedFromMiningPool)
	distributeGoal.BlockChainRewardToResourceNodeFromTrafficPool =
		distributeGoal.BlockChainRewardToResourceNodeFromTrafficPool.Sub(totalUsedFromTrafficPool)

	totalUsedFromMiningPool = sdk.ZeroInt()
	totalUsedFromTrafficPool = sdk.ZeroInt()

	// calc total traffic
	totalConsumedOzone := k.GetTotalConsumedOzone(trafficList)
	// 2, calc traffic reward
	for _, nodeTraffic := range trafficList {
		nodeAddr := nodeTraffic.NodeAddress
		nodeTraffic := nodeTraffic.Volume

		shareOfTraffic := nodeTraffic.ToDec().Quo(totalConsumedOzone.ToDec())
		trafficRewardFromMiningPool :=
			distributeGoal.TrafficRewardToResourceNodeFromMiningPool.ToDec().Mul(shareOfTraffic).TruncateInt()
		trafficRewardFromTrafficPool :=
			distributeGoal.TrafficRewardToResourceNodeFromTrafficPool.ToDec().Mul(shareOfTraffic).TruncateInt()

		totalUsedFromMiningPool = totalUsedFromMiningPool.Add(trafficRewardFromMiningPool)
		totalUsedFromTrafficPool = totalUsedFromTrafficPool.Add(trafficRewardFromTrafficPool)

		if _, ok := rewardDetailMap[nodeAddr.String()]; !ok {
			reward := types.NewDefaultReward(nodeAddr)
			rewardDetailMap[nodeAddr.String()] = reward
		}

		newReward := rewardDetailMap[nodeAddr.String()]
		newReward = newReward.AddRewardFromMiningPool(trafficRewardFromMiningPool)
		newReward = newReward.AddRewardFromTrafficPool(trafficRewardFromTrafficPool)
		rewardDetailMap[nodeAddr.String()] = newReward
	}
	// deduct used reward from distributeGoal
	distributeGoal.TrafficRewardToResourceNodeFromMiningPool =
		distributeGoal.TrafficRewardToResourceNodeFromMiningPool.Sub(totalUsedFromMiningPool)
	distributeGoal.TrafficRewardToResourceNodeFromTrafficPool =
		distributeGoal.TrafficRewardToResourceNodeFromTrafficPool.Sub(totalUsedFromTrafficPool)

	return rewardDetailMap, distributeGoal
}

func (k Keeper) calcRewardForIndexingNode(ctx sdk.Context, distributeGoal types.DistributeGoal, rewardDetailMap map[string]types.Reward,
) (map[string]types.Reward, types.DistributeGoal) {

	totalUsedStakeRewardFromMiningPool := sdk.ZeroInt()
	totalUsedStakeRewardFromTrafficPool := sdk.ZeroInt()
	totalUsedIndexingRewardFromMiningPool := sdk.ZeroInt()
	totalUsedIndexingRewardFromTrafficPool := sdk.ZeroInt()

	totalStakeOfIndexingNodes := k.RegisterKeeper.GetLastIndexingNodeTotalStake(ctx)
	indexingNodeList := k.RegisterKeeper.GetAllIndexingNodes(ctx)
	indexingNodeCnt := sdk.NewInt(int64(len(indexingNodeList)))
	for _, node := range indexingNodeList {
		nodeAddr := node.GetNetworkAddr()

		// 1, calc stake reward
		shareOfToken := node.GetTokens().ToDec().Quo(totalStakeOfIndexingNodes.ToDec())
		stakeRewardFromMiningPool :=
			distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool.ToDec().Mul(shareOfToken).TruncateInt()
		stakeRewardFromTrafficPool :=
			distributeGoal.BlockChainRewardToIndexingNodeFromTrafficPool.ToDec().Mul(shareOfToken).TruncateInt()

		totalUsedStakeRewardFromMiningPool = totalUsedStakeRewardFromMiningPool.Add(stakeRewardFromMiningPool)
		totalUsedStakeRewardFromTrafficPool = totalUsedStakeRewardFromTrafficPool.Add(stakeRewardFromTrafficPool)

		// 2, calc indexing reward
		indexingRewardFromMiningPool :=
			distributeGoal.MetaNodeRewardToIndexingNodeFromMiningPool.ToDec().Quo(indexingNodeCnt.ToDec()).TruncateInt()
		indexingRewardFromTrafficPool :=
			distributeGoal.MetaNodeRewardToIndexingNodeFromTrafficPool.ToDec().Quo(indexingNodeCnt.ToDec()).TruncateInt()

		totalUsedIndexingRewardFromMiningPool = totalUsedIndexingRewardFromMiningPool.Add(indexingRewardFromMiningPool)
		totalUsedIndexingRewardFromTrafficPool = totalUsedIndexingRewardFromTrafficPool.Add(indexingRewardFromTrafficPool)

		if _, ok := rewardDetailMap[nodeAddr.String()]; !ok {
			reward := types.NewDefaultReward(nodeAddr)
			rewardDetailMap[nodeAddr.String()] = reward
		}

		newReward := rewardDetailMap[nodeAddr.String()]
		newReward = newReward.AddRewardFromMiningPool(stakeRewardFromMiningPool.Add(indexingRewardFromMiningPool))
		newReward = newReward.AddRewardFromTrafficPool(stakeRewardFromTrafficPool.Add(indexingRewardFromTrafficPool))
		rewardDetailMap[nodeAddr.String()] = newReward
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

func (k Keeper) GetTotalConsumedOzone(trafficList []types.SingleNodeVolume) sdk.Int {
	totalTraffic := sdk.ZeroInt()
	for _, vol := range trafficList {
		totalTraffic = totalTraffic.Add(vol.Volume)
	}
	return totalTraffic
}

func (k Keeper) splitRewardByStake(ctx sdk.Context, totalReward sdk.Int,
) (validatorReward sdk.Int, resourceNodeReward sdk.Int, indexingNodeReward sdk.Int) {

	validatorBondedTokens := k.StakingKeeper.TotalBondedTokens(ctx).ToDec()
	resourceNodeBondedTokens := k.RegisterKeeper.GetLastResourceNodeTotalStake(ctx).ToDec()
	indexingNodeBondedTokens := k.RegisterKeeper.GetLastIndexingNodeTotalStake(ctx).ToDec()
	totalBondedTokens := validatorBondedTokens.Add(resourceNodeBondedTokens).Add(indexingNodeBondedTokens)

	validatorReward = totalReward.ToDec().Mul(validatorBondedTokens).Quo(totalBondedTokens).TruncateInt()
	resourceNodeReward = totalReward.ToDec().Mul(resourceNodeBondedTokens).Quo(totalBondedTokens).TruncateInt()
	indexingNodeReward = totalReward.ToDec().Mul(indexingNodeBondedTokens).Quo(totalBondedTokens).TruncateInt()

	return
}
