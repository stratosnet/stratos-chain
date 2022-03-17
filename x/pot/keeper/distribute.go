package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func (k Keeper) DistributePotReward(ctx sdk.Context, trafficList []types.SingleWalletVolume, epoch sdk.Int) (totalConsumedOzone sdk.Dec, err error) {
	distributeGoal := types.InitDistributeGoal()
	rewardDetailMap := make(map[string]types.Reward) //key: wallet address

	//1, calc traffic reward in total
	totalConsumedOzone, distributeGoal, err = k.CalcTrafficRewardInTotal(ctx, trafficList, distributeGoal)
	if err != nil {
		return totalConsumedOzone, err
	}

	//2, calc mining reward in total
	distributeGoal, err = k.CalcMiningRewardInTotal(ctx, distributeGoal)
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
	rewardDetailMap, distributeGoalBalance = k.CalcRewardForResourceNode(ctx, trafficList, distributeGoalBalance, rewardDetailMap)

	//4, calc reward from indexing node, store to rewardDetailMap by wallet address(owner address)
	rewardDetailMap, distributeGoalBalance = k.CalcRewardForIndexingNode(ctx, distributeGoalBalance, rewardDetailMap)

	//5, deduct reward from provider account (the value of parameter of distributeGoal will not change)
	err = k.deductRewardFromRewardProviderAccount(ctx, distributeGoal, epoch)
	if err != nil {
		return totalConsumedOzone, err
	}

	//6, distribute skate reward to fee pool for validators
	distributeGoalBalance, err = k.distributeValidatorRewardToFeePool(ctx, distributeGoalBalance)
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
	err = k.returnBalance(ctx, distributeGoalBalance, epoch)
	if err != nil {
		return totalConsumedOzone, err
	}

	return totalConsumedOzone, nil
}

func (k Keeper) deductRewardFromRewardProviderAccount(ctx sdk.Context, goal types.DistributeGoal, epoch sdk.Int) (err error) {
	totalRewardFromMiningPool := goal.BlockChainRewardToValidatorFromMiningPool.
		Add(goal.BlockChainRewardToIndexingNodeFromMiningPool).
		Add(goal.MetaNodeRewardToIndexingNodeFromMiningPool).
		Add(goal.BlockChainRewardToResourceNodeFromMiningPool).
		Add(goal.TrafficRewardToResourceNodeFromMiningPool)
	totalRewardFromTrafficPool := goal.BlockChainRewardToValidatorFromTrafficPool.
		Add(goal.BlockChainRewardToIndexingNodeFromTrafficPool).
		Add(goal.MetaNodeRewardToIndexingNodeFromTrafficPool).
		Add(goal.BlockChainRewardToResourceNodeFromTrafficPool).
		Add(goal.TrafficRewardToResourceNodeFromTrafficPool)

	// deduct mining reward from foundation account
	foundationAccountAddr := k.SupplyKeeper.GetModuleAddress(types.FoundationAccount)
	if foundationAccountAddr == nil {
		ctx.Logger().Error("foundation account address of distribution module does not exist.")
		return types.ErrUnknownAccountAddress
	}

	amountToDeduct := sdk.NewCoins(totalRewardFromMiningPool)
	hasCoin := k.BankKeeper.HasCoins(ctx, foundationAccountAddr, amountToDeduct)
	if !hasCoin {
		ctx.Logger().Info("balance of foundation account is 0")
		return types.ErrInsufficientFoundationAccBalance
	}
	_, err = k.BankKeeper.SubtractCoins(ctx, foundationAccountAddr, amountToDeduct)
	if err != nil {
		return err
	}

	// update mined token record by adding mining reward
	oldTotalMinedToken := k.GetTotalMinedTokens(ctx)
	newTotalMinedToken := oldTotalMinedToken.Add(totalRewardFromMiningPool)
	k.SetTotalMinedTokens(ctx, newTotalMinedToken)
	k.setMinedTokens(ctx, epoch, totalRewardFromMiningPool)

	// deduct traffic reward from prepay pool
	totalUnIssuedPrepay := k.RegisterKeeper.GetTotalUnissuedPrepay(ctx)
	newTotalUnIssuedPrePay := totalUnIssuedPrepay.Sub(totalRewardFromTrafficPool)
	if newTotalUnIssuedPrePay.IsNegative() {
		return types.ErrInsufficientUnissuedPrePayBalance
	}
	k.RegisterKeeper.SetTotalUnissuedPrepay(ctx, newTotalUnIssuedPrePay)

	return nil
}

func (k Keeper) returnBalance(ctx sdk.Context, goal types.DistributeGoal, currentEpoch sdk.Int) (err error) {
	balanceOfMiningPool := goal.BlockChainRewardToIndexingNodeFromMiningPool.
		Add(goal.MetaNodeRewardToIndexingNodeFromMiningPool).
		Add(goal.BlockChainRewardToResourceNodeFromMiningPool).
		Add(goal.TrafficRewardToResourceNodeFromMiningPool)
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
	oldMinedToken := k.GetMinedTokens(ctx, currentEpoch)
	newMinedToken := oldMinedToken.Sub(balanceOfMiningPool)
	k.SetTotalMinedTokens(ctx, newTotalMinedToken)
	k.setMinedTokens(ctx, currentEpoch, newMinedToken)

	// return balance to prepay pool
	totalUnIssuedPrepay := k.RegisterKeeper.GetTotalUnissuedPrepay(ctx)
	newTotalUnIssuedPrePay := totalUnIssuedPrepay.Add(balanceOfTrafficPool)
	k.RegisterKeeper.SetTotalUnissuedPrepay(ctx, newTotalUnIssuedPrePay)

	return nil
}

func (k Keeper) CalcTrafficRewardInTotal(
	ctx sdk.Context, trafficList []types.SingleWalletVolume, distributeGoal types.DistributeGoal,
) (sdk.Dec, types.DistributeGoal, error) {

	totalConsumedOzone, totalTrafficReward := k.getTrafficReward(ctx, trafficList)
	totalMinedTokens := k.GetTotalMinedTokens(ctx)
	miningParam, err := k.GetMiningRewardParamByMinedToken(ctx, totalMinedTokens)
	if err != nil && err != types.ErrOutOfIssuance {
		return sdk.Dec{}, distributeGoal, err
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
	distributeGoal = distributeGoal.AddBlockChainRewardToValidatorFromTrafficPool(sdk.NewCoin(k.BondDenom(ctx), stakeRewardToValidators))
	distributeGoal = distributeGoal.AddBlockChainRewardToResourceNodeFromTrafficPool(sdk.NewCoin(k.BondDenom(ctx), stakeRewardToResourceNodes))
	distributeGoal = distributeGoal.AddBlockChainRewardToIndexingNodeFromTrafficPool(sdk.NewCoin(k.BondDenom(ctx), stakeRewardToIndexingNodes))
	distributeGoal = distributeGoal.AddTrafficRewardToResourceNodeFromTrafficPool(sdk.NewCoin(k.BondDenom(ctx), trafficReward))
	distributeGoal = distributeGoal.AddMetaNodeRewardToIndexingNodeFromTrafficPool(sdk.NewCoin(k.BondDenom(ctx), indexingReward))

	return totalConsumedOzone, distributeGoal, nil
}

// [S] is initial genesis deposit by all resource nodes and meta nodes at t=0
// The current unissued prepay Volume Pool [pt] is the total remaining prepay uSTOS kept by Stratos Network but not issued to Resource Node as rewards. At time t=0,  pt=0
// total consumed Ozone is [Y]
// The remaining total Ozone limit [lt] is the upper bound of total Ozone that users can purchase from Stratos blockchain.
// the total generated traffic rewards as [R]
// R = (S + Pt) * Y / (Lt + Y)
func (k Keeper) getTrafficReward(ctx sdk.Context, trafficList []types.SingleWalletVolume) (totalConsumedOzone, result sdk.Dec) {
	S := k.RegisterKeeper.GetInitialGenesisStakeTotal(ctx).ToDec()
	if S.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("initial genesis deposit by all resource nodes and meta nodes is 0")
	}
	Pt := k.RegisterKeeper.GetTotalUnissuedPrepay(ctx).Amount.ToDec()
	if Pt.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("total remaining prepay not issued is 0")
	}
	Y := k.GetTotalConsumedUoz(trafficList).ToDec()
	if Y.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("total consumed uoz is 0")
	}
	Lt := k.RegisterKeeper.GetRemainingOzoneLimit(ctx).ToDec()
	if Lt.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("remaining total uoz limit is 0")
	}
	R := S.Add(Pt).Mul(Y).Quo(Lt.Add(Y))
	if R.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("traffic reward to distribute is 0")
	}
	return Y, R
}

// allocate mining reward from foundation account
func (k Keeper) CalcMiningRewardInTotal(ctx sdk.Context, distributeGoal types.DistributeGoal) (types.DistributeGoal, error) {
	totalMinedTokens := k.GetTotalMinedTokens(ctx)
	miningParam, err := k.GetMiningRewardParamByMinedToken(ctx, totalMinedTokens)

	totalMiningReward := miningParam.MiningReward
	if err != nil {
		return distributeGoal, err
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

	stakeRewardToValidators, stakeRewardToResourceNodes, stakeRewardToIndexingNodes := k.splitRewardByStake(ctx, stakeReward)
	distributeGoal = distributeGoal.AddBlockChainRewardToValidatorFromMiningPool(sdk.NewCoin(k.RewardDenom(ctx), stakeRewardToValidators))
	distributeGoal = distributeGoal.AddBlockChainRewardToResourceNodeFromMiningPool(sdk.NewCoin(k.RewardDenom(ctx), stakeRewardToResourceNodes))
	distributeGoal = distributeGoal.AddBlockChainRewardToIndexingNodeFromMiningPool(sdk.NewCoin(k.RewardDenom(ctx), stakeRewardToIndexingNodes))
	distributeGoal = distributeGoal.AddTrafficRewardToResourceNodeFromMiningPool(sdk.NewCoin(k.RewardDenom(ctx), trafficReward))
	distributeGoal = distributeGoal.AddMetaNodeRewardToIndexingNodeFromMiningPool(sdk.NewCoin(k.RewardDenom(ctx), indexingReward))
	return distributeGoal, nil
}

func (k Keeper) distributeRewardToSdsNodes(ctx sdk.Context, rewardDetailList []types.Reward, currentEpoch sdk.Int) (err error) {
	matureEpoch := k.getMatureEpochByCurrentEpoch(ctx, currentEpoch)

	for _, reward := range rewardDetailList {
		walletAddr := reward.WalletAddress
		k.addNewRewardAndReCalcTotal(ctx, walletAddr, currentEpoch, matureEpoch, reward)
	}
	k.SetLastReportedEpoch(ctx, currentEpoch)
	return nil
}

func (k Keeper) addNewRewardAndReCalcTotal(ctx sdk.Context, account sdk.AccAddress, currentEpoch sdk.Int, matureEpoch sdk.Int, newReward types.Reward) {
	newRewardTotal := newReward.RewardFromMiningPool.Add(newReward.RewardFromTrafficPool...)
	oldMatureTotal := k.GetMatureTotalReward(ctx, account)
	oldImmatureTotal := k.GetImmatureTotalReward(ctx, account)
	matureStartEpoch := k.GetLastReportedEpoch(ctx).Int64() + 1
	matureEndEpoch := currentEpoch.Int64()

	immatureToMature := sdk.Coins{}
	for i := matureStartEpoch; i <= matureEndEpoch; i++ {
		rewardTotal := sdk.Coins{}
		reward, found := k.GetIndividualReward(ctx, account, sdk.NewInt(i))
		if found {
			rewardTotal = reward.RewardFromMiningPool.Add(reward.RewardFromTrafficPool...)
		}
		immatureToMature = immatureToMature.Add(rewardTotal...)
	}

	//deduct slashing amount from mature total pool
	finalMatureTotal := k.RegisterKeeper.DeductSlashing(ctx, account, oldMatureTotal)
	//deduct slashing amount from upcoming mature reward, don't need to deduct slashing from immatureTotal & individual
	finalNewMature := k.RegisterKeeper.DeductSlashing(ctx, account, immatureToMature)

	matureTotal := finalMatureTotal.Add(finalNewMature...)
	immatureTotal := oldImmatureTotal.Sub(immatureToMature).Add(newRewardTotal...)

	rewardAddressPool := k.GetRewardAddressPool(ctx)
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

	k.SetMatureTotalReward(ctx, account, matureTotal)
	k.SetImmatureTotalReward(ctx, account, immatureTotal)
	k.SetIndividualReward(ctx, account, matureEpoch, newReward)
}

// reward will mature 14 days since distribution. Each epoch interval is about 10 minutes.
func (k Keeper) getMatureEpochByCurrentEpoch(ctx sdk.Context, currentEpoch sdk.Int) (matureEpoch sdk.Int) {
	// 14 days = 20160 minutes = 2016 epochs
	paramMatureEpoch := sdk.NewInt(k.MatureEpoch(ctx))
	matureEpoch = paramMatureEpoch.Add(currentEpoch)
	return matureEpoch
}

// move reward to fee pool for validator traffic reward distribution
func (k Keeper) distributeValidatorRewardToFeePool(ctx sdk.Context, distributeGoal types.DistributeGoal) (types.DistributeGoal, error) {
	rewardFromMiningPool := distributeGoal.BlockChainRewardToValidatorFromMiningPool
	rewardFromTrafficPool := distributeGoal.BlockChainRewardToValidatorFromTrafficPool
	totalRewardSendToFeePool := sdk.NewCoins(rewardFromMiningPool).Add(rewardFromTrafficPool)

	feePoolAccAddr := k.SupplyKeeper.GetModuleAddress(k.feeCollectorName)

	if feePoolAccAddr == nil {
		ctx.Logger().Error("account address of distribution module does not exist.")
		return distributeGoal, types.ErrUnknownAccountAddress
	}

	_, err := k.BankKeeper.AddCoins(ctx, feePoolAccAddr, totalRewardSendToFeePool)
	if err != nil {
		return distributeGoal, err
	}

	distributeGoal.BlockChainRewardToValidatorFromMiningPool = sdk.Coin{}
	distributeGoal.BlockChainRewardToValidatorFromTrafficPool = sdk.Coin{}

	return distributeGoal, nil
}

func (k Keeper) CalcRewardForResourceNode(ctx sdk.Context, trafficList []types.SingleWalletVolume,
	distributeGoal types.DistributeGoal, rewardDetailMap map[string]types.Reward,
) (map[string]types.Reward, types.DistributeGoal) {

	totalUsedFromMiningPool := sdk.NewCoin(k.RewardDenom(ctx), sdk.ZeroInt())
	totalUsedFromTrafficPool := sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())

	// 1, calc stake reward
	totalStakeOfResourceNodes := k.RegisterKeeper.GetResourceNodeBondedToken(ctx).Amount
	resourceNodeList := k.RegisterKeeper.GetAllResourceNodes(ctx)
	for _, node := range resourceNodeList {
		walletAddr := node.GetOwnerAddr()

		shareOfToken := node.GetTokens().ToDec().Quo(totalStakeOfResourceNodes.ToDec())
		stakeRewardFromMiningPool := sdk.NewCoin(k.RewardDenom(ctx),
			distributeGoal.BlockChainRewardToResourceNodeFromMiningPool.Amount.ToDec().Mul(shareOfToken).TruncateInt())
		stakeRewardFromTrafficPool := sdk.NewCoin(k.BondDenom(ctx),
			distributeGoal.BlockChainRewardToResourceNodeFromTrafficPool.Amount.ToDec().Mul(shareOfToken).TruncateInt())

		totalUsedFromMiningPool = totalUsedFromMiningPool.Add(stakeRewardFromMiningPool)
		totalUsedFromTrafficPool = totalUsedFromTrafficPool.Add(stakeRewardFromTrafficPool)

		if _, ok := rewardDetailMap[walletAddr.String()]; !ok {
			reward := types.NewDefaultReward(walletAddr)
			rewardDetailMap[walletAddr.String()] = reward
		}

		newReward := rewardDetailMap[walletAddr.String()]
		newReward = newReward.AddRewardFromMiningPool(stakeRewardFromMiningPool)
		newReward = newReward.AddRewardFromTrafficPool(stakeRewardFromTrafficPool)
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
		trafficRewardFromMiningPool := sdk.NewCoin(k.RewardDenom(ctx),
			distributeGoal.TrafficRewardToResourceNodeFromMiningPool.Amount.ToDec().Mul(shareOfTraffic).TruncateInt())
		trafficRewardFromTrafficPool := sdk.NewCoin(k.BondDenom(ctx),
			distributeGoal.TrafficRewardToResourceNodeFromTrafficPool.Amount.ToDec().Mul(shareOfTraffic).TruncateInt())

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

func (k Keeper) CalcRewardForIndexingNode(ctx sdk.Context, distributeGoal types.DistributeGoal, rewardDetailMap map[string]types.Reward,
) (map[string]types.Reward, types.DistributeGoal) {

	totalUsedStakeRewardFromMiningPool := sdk.NewCoin(k.RewardDenom(ctx), sdk.ZeroInt())
	totalUsedStakeRewardFromTrafficPool := sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())
	totalUsedIndexingRewardFromMiningPool := sdk.NewCoin(k.RewardDenom(ctx), sdk.ZeroInt())
	totalUsedIndexingRewardFromTrafficPool := sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())

	totalStakeOfIndexingNodes := k.RegisterKeeper.GetIndexingNodeBondedToken(ctx).Amount
	indexingNodeList := k.RegisterKeeper.GetAllIndexingNodes(ctx)
	indexingNodeCnt := sdk.NewInt(int64(len(indexingNodeList)))
	for _, node := range indexingNodeList {
		walletAddr := node.GetOwnerAddr()

		// 1, calc stake reward
		shareOfToken := node.GetTokens().ToDec().Quo(totalStakeOfIndexingNodes.ToDec())
		stakeRewardFromMiningPool := sdk.NewCoin(k.RewardDenom(ctx),
			distributeGoal.BlockChainRewardToIndexingNodeFromMiningPool.Amount.ToDec().Mul(shareOfToken).TruncateInt())
		stakeRewardFromTrafficPool := sdk.NewCoin(k.BondDenom(ctx),
			distributeGoal.BlockChainRewardToIndexingNodeFromTrafficPool.Amount.ToDec().Mul(shareOfToken).TruncateInt())

		totalUsedStakeRewardFromMiningPool = totalUsedStakeRewardFromMiningPool.Add(stakeRewardFromMiningPool)
		totalUsedStakeRewardFromTrafficPool = totalUsedStakeRewardFromTrafficPool.Add(stakeRewardFromTrafficPool)

		// 2, calc indexing reward
		indexingRewardFromMiningPool := sdk.NewCoin(k.RewardDenom(ctx),
			distributeGoal.MetaNodeRewardToIndexingNodeFromMiningPool.Amount.ToDec().Quo(indexingNodeCnt.ToDec()).TruncateInt())
		indexingRewardFromTrafficPool := sdk.NewCoin(k.BondDenom(ctx),
			distributeGoal.MetaNodeRewardToIndexingNodeFromTrafficPool.Amount.ToDec().Quo(indexingNodeCnt.ToDec()).TruncateInt())

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

func (k Keeper) GetTotalConsumedUoz(trafficList []types.SingleWalletVolume) sdk.Int {
	totalTraffic := sdk.ZeroInt()
	for _, vol := range trafficList {
		totalTraffic = totalTraffic.Add(vol.Volume)
	}
	return totalTraffic
}

func (k Keeper) splitRewardByStake(ctx sdk.Context, totalReward sdk.Int,
) (validatorReward sdk.Int, resourceNodeReward sdk.Int, indexingNodeReward sdk.Int) {

	validatorBondedTokens := k.StakingKeeper.TotalBondedTokens(ctx).ToDec()
	resourceNodeBondedTokens := k.RegisterKeeper.GetResourceNodeBondedToken(ctx).Amount.ToDec()
	indexingNodeBondedTokens := k.RegisterKeeper.GetIndexingNodeBondedToken(ctx).Amount.ToDec()

	totalBondedTokens := validatorBondedTokens.Add(resourceNodeBondedTokens).Add(indexingNodeBondedTokens)

	validatorReward = totalReward.ToDec().Mul(validatorBondedTokens).Quo(totalBondedTokens).TruncateInt()
	resourceNodeReward = totalReward.ToDec().Mul(resourceNodeBondedTokens).Quo(totalBondedTokens).TruncateInt()
	indexingNodeReward = totalReward.ToDec().Mul(indexingNodeBondedTokens).Quo(totalBondedTokens).TruncateInt()

	return
}

func (k Keeper) IteratorImmatureTotal(ctx sdk.Context, handler func(walletAddress sdk.AccAddress, immatureTotal sdk.Coins) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ImmatureTotalRewardKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.AccAddress(iter.Key()[len(types.ImmatureTotalRewardKeyPrefix) : len(types.ImmatureTotalRewardKeyPrefix)+sdk.AddrLen])
		var immatureTotal sdk.Coins
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &immatureTotal)
		if handler(addr, immatureTotal) {
			break
		}
	}
}

func (k Keeper) IteratorMatureTotal(ctx sdk.Context, handler func(walletAddress sdk.AccAddress, matureTotal sdk.Coins) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.MatureTotalRewardKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.AccAddress(iter.Key()[len(types.MatureTotalRewardKeyPrefix) : len(types.MatureTotalRewardKeyPrefix)+sdk.AddrLen])
		var matureTotal sdk.Coins
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &matureTotal)
		if handler(addr, matureTotal) {
			break
		}
	}
}
