package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

func (k Keeper) DistributePotReward(ctx sdk.Context, trafficList []*types.SingleWalletVolume, epoch sdk.Int) (totalConsumedOzone sdk.Dec, err error) {
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

	//5, [TLC] deduct reward from provider account (the value of parameter of distributeGoal will not change)
	err = k.deductRewardFromRewardProviderAccount(ctx, distributeGoal, epoch)
	if err != nil {
		return totalConsumedOzone, err
	}

	//6, [TLC] distribute staking reward to fee pool for validators
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

	//9, [TLC] return balance to traffic pool & mining pool
	err = k.returnBalance(ctx, distributeGoalBalance, epoch)
	if err != nil {
		return totalConsumedOzone, err
	}

	//10, mature rewards for all nodes
	k.rewardMatureAndSubSlashing(ctx, epoch)

	//11, save reported epoch
	k.SetLastReportedEpoch(ctx, epoch)

	//12, [TLC] transfer balance of miningReward&trafficReward pools to totalReward pool, utilized for future Withdraw Tx
	err = k.TransferMiningTrafficRewardsToTotalRewards(ctx)
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

	// [TLC][Foundation -> MiningRewardPool]: deduct mining reward from foundation account
	foundationAccountAddr := k.AccountKeeper.GetModuleAddress(types.FoundationAccount)
	if foundationAccountAddr == nil {
		ctx.Logger().Error("foundation account address of distribution module does not exist.")
		return types.ErrUnknownAccountAddress
	}

	hasCoin := k.BankKeeper.HasBalance(ctx, foundationAccountAddr, totalRewardFromMiningPool)
	if !hasCoin {
		ctx.Logger().Info("balance of foundation account is 0")
		return types.ErrInsufficientFoundationAccBalance
	}
	amountToDeduct := sdk.NewCoins(totalRewardFromMiningPool)
	err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, foundationAccountAddr, types.MiningRewardPool, amountToDeduct)
	if err != nil {
		return err
	}

	// [Non-TLC] update mined token record by adding mining reward
	oldTotalMinedToken := k.GetTotalMinedTokens(ctx)
	newTotalMinedToken := oldTotalMinedToken.Add(totalRewardFromMiningPool)
	k.SetTotalMinedTokens(ctx, newTotalMinedToken)
	k.setMinedTokens(ctx, epoch, totalRewardFromMiningPool)

	// [TLC][TotalUnIssuedPrepay -> TrafficRewardPool]: deduct traffic reward from prepay pool
	totalUnissuedPrepayAddr := k.AccountKeeper.GetModuleAddress(regtypes.TotalUnissuedPrepayName)
	if totalUnissuedPrepayAddr == nil {
		ctx.Logger().Error("TotalUnIssuedPrepay account address of register module does not exist.")
		return types.ErrUnknownAccountAddress
	}

	hasCoinInUnissuedPrepay := k.BankKeeper.HasBalance(ctx, totalUnissuedPrepayAddr, totalRewardFromTrafficPool)
	if !hasCoinInUnissuedPrepay {
		ctx.Logger().Info("Insufficient balance of TotalUnIssuedPrepay module account")
		return types.ErrInsufficientUnissuedPrePayBalance
	}
	err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, regtypes.TotalUnissuedPrepayName, types.TrafficRewardPool, sdk.NewCoins(totalRewardFromTrafficPool))
	//totalUnIssuedPrepay := k.RegisterKeeper.GetTotalUnissuedPrepay(ctx)
	//newTotalUnIssuedPrePay := totalUnIssuedPrepay.Sub(totalRewardFromTrafficPool)
	//if newTotalUnIssuedPrePay.IsNegative() {
	//	return types.ErrInsufficientUnissuedPrePayBalance
	//}
	//k.RegisterKeeper.SetTotalUnissuedPrepay(ctx, newTotalUnIssuedPrePay)

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
	//foundationAccountAddr := k.AccountKeeper.GetModuleAddress(types.FoundationAccount)
	//if foundationAccountAddr == nil {
	//	ctx.Logger().Error("foundation account address of distribution module does not exist.")
	//	return types.ErrUnknownAccountAddress
	//}
	// [TLC] [MiningRewardPool -> FoundationAccount]
	amountToAdd := sdk.NewCoins(balanceOfMiningPool)
	err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.MiningRewardPool, types.FoundationAccount, amountToAdd)
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
	// [TLC][TrafficRewardPool -> TotalUnIssuedPrepay]
	err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.TrafficRewardPool, regtypes.TotalUnissuedPrepayName, sdk.NewCoins(balanceOfTrafficPool))
	if err != nil {
		return err
	}
	//totalUnIssuedPrepay := k.RegisterKeeper.GetTotalUnissuedPrepay(ctx)
	//newTotalUnIssuedPrePay := totalUnIssuedPrepay.Add(balanceOfTrafficPool)
	//k.RegisterKeeper.SetTotalUnissuedPrepay(ctx, newTotalUnIssuedPrePay)

	return nil
}

func (k Keeper) CalcTrafficRewardInTotal(
	ctx sdk.Context, trafficList []*types.SingleWalletVolume, distributeGoal types.DistributeGoal,
) (sdk.Dec, types.DistributeGoal, error) {

	totalConsumedOzone, totalTrafficReward := k.GetTrafficReward(ctx, trafficList)
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
func (k Keeper) GetTrafficReward(ctx sdk.Context, trafficList []*types.SingleWalletVolume) (totalConsumedOzone, result sdk.Dec) {
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
		walletAddr, err := sdk.AccAddressFromBech32(reward.WalletAddress)
		if err != nil {
			continue
		}
		k.addNewIndividualAndUpdateImmatureTotal(ctx, walletAddr, matureEpoch, reward)
	}
	return nil
}

func (k Keeper) addNewIndividualAndUpdateImmatureTotal(ctx sdk.Context, account sdk.AccAddress, matureEpoch sdk.Int, newReward types.Reward) {
	newIndividualTotal := newReward.RewardFromMiningPool.Add(newReward.RewardFromTrafficPool...)
	oldImmatureTotal := k.GetImmatureTotalReward(ctx, account)
	newImmatureTotal := oldImmatureTotal.Add(newIndividualTotal...)

	k.SetIndividualReward(ctx, account, matureEpoch, newReward)
	k.SetImmatureTotalReward(ctx, account, newImmatureTotal)
}

func (k Keeper) rewardMatureAndSubSlashing(ctx sdk.Context, currentEpoch sdk.Int) {

	matureStartEpoch := k.GetLastReportedEpoch(ctx).Int64() + 1
	matureEndEpoch := currentEpoch.Int64()

	for i := matureStartEpoch; i <= matureEndEpoch; i++ {
		k.IteratorIndividualReward(ctx, sdk.NewInt(i), func(walletAddress sdk.AccAddress, individualReward types.Reward) (stop bool) {
			oldMatureTotal := k.GetMatureTotalReward(ctx, walletAddress)
			oldImmatureTotal := k.GetImmatureTotalReward(ctx, walletAddress)
			immatureToMature := individualReward.RewardFromMiningPool.Add(individualReward.RewardFromTrafficPool...)

			//deduct slashing amount from mature total pool
			oldMatureTotalSubSlashing := k.RegisterKeeper.DeductSlashing(ctx, walletAddress, oldMatureTotal)
			//deduct slashing amount from upcoming mature reward, don't need to deduct slashing from immatureTotal & individual
			immatureToMatureSubSlashing := k.RegisterKeeper.DeductSlashing(ctx, walletAddress, immatureToMature)

			matureTotal := oldMatureTotalSubSlashing.Add(immatureToMatureSubSlashing...)
			immatureTotal := oldImmatureTotal.Sub(immatureToMature)

			k.SetMatureTotalReward(ctx, walletAddress, matureTotal)
			k.SetImmatureTotalReward(ctx, walletAddress, immatureTotal)
			return false
		})
	}
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
	//totalRewardSendToFeePool := sdk.NewCoins(rewardFromMiningPool).Add(rewardFromTrafficPool)

	feePoolAccAddr := k.AccountKeeper.GetModuleAddress(k.feeCollectorName)

	if feePoolAccAddr == nil {
		ctx.Logger().Error("account address of distribution module does not exist.")
		return distributeGoal, types.ErrUnknownAccountAddress
	}

	// separately sending totalRerewardFromMiningPool and rewardFromTrafficPool instead of sending totalRewardSendToFeePool to feeCollector module acc
	// [TLC] [MiningRewardPool -> feeCollectorPool]
	err := k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.MiningRewardPool, feePoolAccAddr, sdk.NewCoins(rewardFromMiningPool))
	if err != nil {
		return distributeGoal, err
	}
	// [TLC] [TrafficRewardPool -> feeCollectorPool]
	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.TrafficRewardPool, feePoolAccAddr, sdk.NewCoins(rewardFromTrafficPool))
	if err != nil {
		return distributeGoal, err
	}

	distributeGoal.BlockChainRewardToValidatorFromMiningPool = sdk.Coin{}
	distributeGoal.BlockChainRewardToValidatorFromTrafficPool = sdk.Coin{}

	return distributeGoal, nil
}

func (k Keeper) CalcRewardForResourceNode(ctx sdk.Context, trafficList []*types.SingleWalletVolume,
	distributeGoal types.DistributeGoal, rewardDetailMap map[string]types.Reward,
) (map[string]types.Reward, types.DistributeGoal) {

	totalUsedFromMiningPool := sdk.NewCoin(k.RewardDenom(ctx), sdk.ZeroInt())
	totalUsedFromTrafficPool := sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())

	// 1, calc stake reward
	totalStakeOfResourceNodes := k.RegisterKeeper.GetResourceNodeBondedToken(ctx).Amount
	resourceNodeList := k.RegisterKeeper.GetAllResourceNodes(ctx)
	for _, node := range resourceNodeList.ResourceNodes {
		walletAddr, err := sdk.AccAddressFromBech32(node.OwnerAddress)
		if err != nil {
			continue
		}
		tokens, ok := sdk.NewIntFromString(node.Tokens.String())
		if !ok {
			continue
		}
		shareOfToken := tokens.ToDec().Quo(totalStakeOfResourceNodes.ToDec())
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
		walletAddr, err := sdk.AccAddressFromBech32(walletTraffic.WalletAddress)
		if err != nil {
			continue
		}
		//walletAddr := walletTraffic.WalletAddress
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
	indexingNodeCnt := sdk.NewInt(int64(len(indexingNodeList.IndexingNodes)))
	for _, node := range indexingNodeList.IndexingNodes {
		walletAddr, err := sdk.AccAddressFromBech32(node.OwnerAddress)
		if err != nil {
			continue
		}
		tokens, ok := sdk.NewIntFromString(node.Tokens.String())
		if !ok {
			continue
		}

		// 1, calc stake reward
		shareOfToken := tokens.ToDec().Quo(totalStakeOfIndexingNodes.ToDec())
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

func (k Keeper) GetTotalConsumedUoz(trafficList []*types.SingleWalletVolume) sdk.Int {
	totalTraffic := sdk.ZeroInt()
	for _, vol := range trafficList {
		toAdd, ok := sdk.NewIntFromString(vol.Volume.String())
		if !ok {
			continue
		}
		totalTraffic = totalTraffic.Add(toAdd)
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

func (k Keeper) IteratorIndividualReward(ctx sdk.Context, epoch sdk.Int, handler func(walletAddress sdk.AccAddress, individualReward types.Reward) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetIndividualRewardIteratorKey(epoch))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.AccAddress(iter.Key()[len(types.GetIndividualRewardIteratorKey(epoch)):])

		var individualReward types.Reward
		types.ModuleCdc.MustUnmarshalLengthPrefixed(iter.Value(), &individualReward)
		if handler(addr, individualReward) {
			break
		}
	}
}

func (k Keeper) IteratorImmatureTotal(ctx sdk.Context, handler func(walletAddress sdk.AccAddress, immatureTotal sdk.Coins) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ImmatureTotalRewardKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.AccAddress(iter.Key()[len(types.ImmatureTotalRewardKeyPrefix):])
		var immatureTotal sdk.Coins
		types.ModuleCdc.MustUnmarshalLengthPrefixed(iter.Value(), &immatureTotal)
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
		addr := sdk.AccAddress(iter.Key()[len(types.MatureTotalRewardKeyPrefix):])
		var matureTotal sdk.Coins
		types.ModuleCdc.MustUnmarshalLengthPrefixed(iter.Value(), &matureTotal)
		if handler(addr, matureTotal) {
			break
		}
	}
}

func (k Keeper) TransferMiningTrafficRewardsToTotalRewards(ctx sdk.Context) error {
	miningRewardAccountAddr := k.AccountKeeper.GetModuleAddress(types.MiningRewardPool)
	if miningRewardAccountAddr == nil {
		ctx.Logger().Error("mining reward account address of distribution module does not exist.")
		return types.ErrUnknownAccountAddress
	}
	miningRewardPoolBalances := k.BankKeeper.GetAllBalances(ctx, miningRewardAccountAddr)

	trafficRewardAccountAddr := k.AccountKeeper.GetModuleAddress(types.TrafficRewardPool)
	if trafficRewardAccountAddr == nil {
		ctx.Logger().Error("traffic reward account address of distribution module does not exist.")
		return types.ErrUnknownAccountAddress
	}
	trafficRewardPoolBalances := k.BankKeeper.GetAllBalances(ctx, trafficRewardAccountAddr)

	err := k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.MiningRewardPool, types.TotalRewardPool, miningRewardPoolBalances)
	if err != nil {
		return err
	}
	err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.TrafficRewardPool, types.TotalRewardPool, trafficRewardPoolBalances)
	if err != nil {
		return err
	}
	return nil
}
