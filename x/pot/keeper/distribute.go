package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/stratosnet/stratos-chain/x/pot/types"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

var (
	foundationToFeeCollector      sdk.Coin
	unissuedPrepayToFeeCollector  sdk.Coin
	foundationToReward            sdk.Coin
	unissuedPrepayToReward        sdk.Coin
	unissuedPrepayToCommunityPool sdk.Dec
	distributeGoal                types.DistributeGoal
	rewardDetailMap               map[string]types.Reward
)

func (k Keeper) DistributePotReward(ctx sdk.Context, trafficList []*types.SingleWalletVolume, epoch sdk.Int) (
	totalConsumedNoz sdk.Dec, err error) {

	k.InitVariable(ctx)

	//1, calc traffic reward in total
	totalConsumedNoz = k.GetTotalConsumedNoz(trafficList).ToDec()
	remaining, total := k.RegisterKeeper.NozSupply(ctx)
	if totalConsumedNoz.Add(remaining.ToDec()).GT(total.ToDec()) {
		return totalConsumedNoz, errors.New("remaining+consumed Noz exceeds total Noz supply")
	}

	distributeGoal, err = k.CalcTrafficRewardInTotal(ctx, distributeGoal, totalConsumedNoz)
	if err != nil {
		return totalConsumedNoz, err
	}
	unissuedPrepayToFeeCollector = distributeGoal.StakeTrafficRewardToValidator

	//2, calc mining reward in total
	distributeGoal, err = k.CalcMiningRewardInTotal(ctx, distributeGoal)
	if err != nil && err != types.ErrOutOfIssuance {
		return totalConsumedNoz, err
	}
	foundationToFeeCollector = distributeGoal.StakeMiningRewardToValidator

	//3, calc reward for resource node, store to rewardDetailMap by wallet address(owner address)
	rewardDetailMap = k.CalcRewardForResourceNode(ctx, totalConsumedNoz, trafficList, distributeGoal, rewardDetailMap)

	//4, calc reward from meta node, store to rewardDetailMap by wallet address(owner address)
	rewardDetailMap = k.CalcRewardForMetaNode(ctx, distributeGoal, rewardDetailMap)

	//5, IMPORTANT: sort map and convert to slice to keep the order
	rewardDetailList := sortDetailMapToSlice(rewardDetailMap)

	//6, record all rewards to resource & meta nodes
	err = k.saveRewardInfo(ctx, rewardDetailList, epoch)
	if err != nil {
		return totalConsumedNoz, err
	}

	//7, mature rewards for all nodes
	totalSlashed := k.rewardMatureAndSubSlashing(ctx, epoch)

	//8, save reported epoch
	k.SetLastReportedEpoch(ctx, epoch)

	//9, update remaining ozone limit
	remainingNozLimit := k.RegisterKeeper.GetRemainingOzoneLimit(ctx)
	k.RegisterKeeper.SetRemainingOzoneLimit(ctx, remainingNozLimit.Add(totalConsumedNoz.TruncateInt()))

	//10, [TLC] transfer balance of miningReward&trafficReward pools to totalReward&totalSlashed pool, utilized for future Withdraw Tx
	err = k.transferTokens(ctx, totalSlashed)
	if err != nil {
		return totalConsumedNoz, err
	}

	return totalConsumedNoz, nil
}

func (k Keeper) CalcTrafficRewardInTotal(
	ctx sdk.Context, distributeGoal types.DistributeGoal, totalConsumedNoz sdk.Dec,
) (types.DistributeGoal, error) {

	totalTrafficReward := k.GetTrafficReward(ctx, totalConsumedNoz)
	totalMinedTokens := k.GetTotalMinedTokens(ctx)
	miningParam, err := k.GetMiningRewardParamByMinedToken(ctx, totalMinedTokens)
	if err != nil && err != types.ErrOutOfIssuance {
		return distributeGoal, err
	}
	stakeTrafficReward := totalTrafficReward.
		Mul(miningParam.BlockChainPercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	trafficRewardToResourceNodes := totalTrafficReward.
		Mul(miningParam.ResourceNodePercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	trafficRewardToMetaNodes := totalTrafficReward.
		Mul(miningParam.MetaNodePercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()

	stakeRewardToValidators, stakeRewardToResourceNodes, stakeRewardToMetaNodes := k.splitRewardByStake(ctx, stakeTrafficReward)
	distributeGoal = distributeGoal.AddStakeTrafficRewardToValidator(sdk.NewCoin(k.BondDenom(ctx), stakeRewardToValidators))
	distributeGoal = distributeGoal.AddStakeTrafficRewardToResourceNode(sdk.NewCoin(k.BondDenom(ctx), stakeRewardToResourceNodes))
	distributeGoal = distributeGoal.AddStakeTrafficRewardToMetaNode(sdk.NewCoin(k.BondDenom(ctx), stakeRewardToMetaNodes))
	distributeGoal = distributeGoal.AddTrafficRewardToResourceNode(sdk.NewCoin(k.BondDenom(ctx), trafficRewardToResourceNodes))
	distributeGoal = distributeGoal.AddTrafficRewardToMetaNode(sdk.NewCoin(k.BondDenom(ctx), trafficRewardToMetaNodes))

	return distributeGoal, nil
}

// [S] is initial genesis deposit by all resource nodes and meta nodes at t=0
// The current unissued prepay Volume Pool [pt] is the total remaining prepay wei kept by Stratos Network but not issued to Resource Node as rewards. At time t=0,  pt=0
// total consumed Ozone is [Y]
// The remaining total Ozone limit [lt] is the upper bound of total Ozone that users can purchase from Stratos blockchain.
// the total generated traffic rewards as [R]
// R = (S + Pt) * Y / (Lt + Y)
func (k Keeper) GetTrafficReward(ctx sdk.Context, totalConsumedNoz sdk.Dec) (result sdk.Dec) {
	St := k.RegisterKeeper.GetEffectiveTotalStake(ctx).ToDec()
	if St.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("effective genesis deposit by all resource nodes and meta nodes is 0")
	}
	Pt := k.RegisterKeeper.GetTotalUnissuedPrepay(ctx).Amount.ToDec()
	if Pt.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("total remaining prepay not issued is 0")
	}
	Y := totalConsumedNoz
	if Y.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("total consumed noz is 0")
	}
	Lt := k.RegisterKeeper.GetRemainingOzoneLimit(ctx).ToDec()
	if Lt.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("remaining total noz limit is 0")
	}
	R := St.Add(Pt).Mul(Y).Quo(Lt.Add(Y))
	if R.Equal(sdk.ZeroDec()) {
		ctx.Logger().Info("traffic reward to distribute is 0")
	}
	return R
}

// allocate mining reward from foundation account
func (k Keeper) CalcMiningRewardInTotal(ctx sdk.Context, distributeGoal types.DistributeGoal) (types.DistributeGoal, error) {
	totalMinedTokens := k.GetTotalMinedTokens(ctx)
	miningParam, err := k.GetMiningRewardParamByMinedToken(ctx, totalMinedTokens)

	totalMiningReward := miningParam.MiningReward
	if err != nil {
		return distributeGoal, err
	}
	stakeMiningReward := totalMiningReward.Amount.ToDec().
		Mul(miningParam.BlockChainPercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	miningRewardToResourceNodes := totalMiningReward.Amount.ToDec().
		Mul(miningParam.ResourceNodePercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	miningRewardToMetaNodes := totalMiningReward.Amount.ToDec().
		Mul(miningParam.MetaNodePercentageInTenThousand.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()

	stakeRewardToValidators, stakeRewardToResourceNodes, stakeRewardToMetaNodes := k.splitRewardByStake(ctx, stakeMiningReward)
	distributeGoal = distributeGoal.AddStakeMiningRewardToValidator(sdk.NewCoin(k.RewardDenom(ctx), stakeRewardToValidators))
	distributeGoal = distributeGoal.AddStakeMiningRewardToResourceNode(sdk.NewCoin(k.RewardDenom(ctx), stakeRewardToResourceNodes))
	distributeGoal = distributeGoal.AddStakeMiningRewardToMetaNode(sdk.NewCoin(k.RewardDenom(ctx), stakeRewardToMetaNodes))
	distributeGoal = distributeGoal.AddMiningRewardToResourceNode(sdk.NewCoin(k.RewardDenom(ctx), miningRewardToResourceNodes))
	distributeGoal = distributeGoal.AddMiningRewardToMetaNode(sdk.NewCoin(k.RewardDenom(ctx), miningRewardToMetaNodes))
	return distributeGoal, nil
}

// Iteration for each rewarded SDS node
func (k Keeper) saveRewardInfo(ctx sdk.Context, rewardDetailList []types.Reward, currentEpoch sdk.Int) (err error) {
	matureEpoch := k.getMatureEpochByCurrentEpoch(ctx, currentEpoch)

	for _, reward := range rewardDetailList {
		walletAddr, err := sdk.AccAddressFromBech32(reward.WalletAddress)
		if err != nil {
			continue
		}
		k.addNewIndividualAndUpdateImmatureTotal(ctx, walletAddr, matureEpoch, reward)
	}

	newMinedTotal := foundationToFeeCollector.Add(foundationToReward)
	oldTotalMinedToken := k.GetTotalMinedTokens(ctx)
	newTotalMinedToken := oldTotalMinedToken.Add(newMinedTotal)
	k.SetTotalMinedTokens(ctx, newTotalMinedToken)
	return nil
}

func (k Keeper) addNewIndividualAndUpdateImmatureTotal(ctx sdk.Context, account sdk.AccAddress, matureEpoch sdk.Int, newReward types.Reward) {
	newIndividualTotal := newReward.RewardFromMiningPool.Add(newReward.RewardFromTrafficPool...)
	oldImmatureTotal := k.GetImmatureTotalReward(ctx, account)
	newImmatureTotal := oldImmatureTotal.Add(newIndividualTotal...)

	k.SetIndividualReward(ctx, account, matureEpoch, newReward)
	k.SetImmatureTotalReward(ctx, account, newImmatureTotal)
}

// Iteration for mature rewards/slashing of all nodes
func (k Keeper) rewardMatureAndSubSlashing(ctx sdk.Context, currentEpoch sdk.Int) (totalSlashed sdk.Coins) {

	matureStartEpoch := k.GetLastReportedEpoch(ctx).Int64() + 1
	matureEndEpoch := currentEpoch.Int64()

	totalSlashed = sdk.Coins{}

	for i := matureStartEpoch; i <= matureEndEpoch; i++ {
		k.IteratorIndividualReward(ctx, sdk.NewInt(i), func(walletAddress sdk.AccAddress, individualReward types.Reward) (stop bool) {
			oldMatureTotal := k.GetMatureTotalReward(ctx, walletAddress)
			oldImmatureTotal := k.GetImmatureTotalReward(ctx, walletAddress)
			immatureToMature := individualReward.RewardFromMiningPool.Add(individualReward.RewardFromTrafficPool...)

			//deduct slashing amount from upcoming mature reward, don't need to deduct slashing from immatureTotal & individual
			remaining, deducted := k.RegisterKeeper.DeductSlashing(ctx, walletAddress, immatureToMature, k.RewardDenom(ctx))
			totalSlashed = totalSlashed.Add(deducted...)

			matureTotal := oldMatureTotal.Add(remaining...)
			immatureTotal := oldImmatureTotal.Sub(immatureToMature)

			k.SetMatureTotalReward(ctx, walletAddress, matureTotal)
			k.SetImmatureTotalReward(ctx, walletAddress, immatureTotal)
			return false
		})
	}
	return totalSlashed
}

// reward will mature 14 days since distribution. Each epoch interval is about 10 minutes.
func (k Keeper) getMatureEpochByCurrentEpoch(ctx sdk.Context, currentEpoch sdk.Int) (matureEpoch sdk.Int) {
	// 14 days = 20160 minutes = 2016 epochs
	paramMatureEpoch := sdk.NewInt(k.MatureEpoch(ctx))
	matureEpoch = paramMatureEpoch.Add(currentEpoch)
	return matureEpoch
}

// Iteration for calculating reward of resource nodes
func (k Keeper) CalcRewardForResourceNode(ctx sdk.Context, totalConsumedNoz sdk.Dec, trafficList []*types.SingleWalletVolume,
	distributeGoal types.DistributeGoal, rewardDetailMap map[string]types.Reward,
) map[string]types.Reward {

	// 1, calc stake reward for resource node by stake
	totalStakeOfResourceNodes := k.RegisterKeeper.GetResourceNodeBondedToken(ctx).Amount

	resourceNodeIterator := k.RegisterKeeper.GetResourceNodeIterator(ctx)
	defer resourceNodeIterator.Close()
	for ; resourceNodeIterator.Valid(); resourceNodeIterator.Next() {

		node := regtypes.MustUnmarshalResourceNode(k.cdc, resourceNodeIterator.Value())
		if node.Status != stakingtypes.Bonded {
			continue
		}

		walletAddr, err := sdk.AccAddressFromBech32(node.OwnerAddress)
		if err != nil {
			continue
		}
		tokens, ok := sdk.NewIntFromString(node.Tokens.String())
		if !ok {
			continue
		}

		shareOfToken := tokens.ToDec().Quo(totalStakeOfResourceNodes.ToDec())

		// stake reward from mining pool
		stakeMiningReward := sdk.NewCoin(k.RewardDenom(ctx),
			distributeGoal.StakeMiningRewardToResourceNode.Amount.ToDec().
				Mul(shareOfToken).
				TruncateInt())

		// stake reward from traffic pool, need to pay community tax
		stakeTrafficRewardBeforeTax := distributeGoal.StakeTrafficRewardToResourceNode.Amount.ToDec().
			Mul(shareOfToken)
		stakeTrafficRewardAfterTax, tax := k.CalcCommunityTax(ctx, stakeTrafficRewardBeforeTax)

		// update rewardDetailMap
		if _, ok := rewardDetailMap[walletAddr.String()]; !ok {
			reward := types.NewDefaultReward(walletAddr)
			rewardDetailMap[walletAddr.String()] = reward
		}
		newReward := rewardDetailMap[walletAddr.String()]
		newReward = newReward.AddRewardFromMiningPool(stakeMiningReward)
		newReward = newReward.AddRewardFromTrafficPool(stakeTrafficRewardAfterTax)
		rewardDetailMap[walletAddr.String()] = newReward

		// record value preparing for transfer
		foundationToReward = foundationToReward.
			Add(stakeMiningReward)
		unissuedPrepayToReward = unissuedPrepayToReward.
			Add(stakeTrafficRewardAfterTax)
		unissuedPrepayToCommunityPool = unissuedPrepayToCommunityPool.
			Add(tax)
	}

	// 2, calc mining & traffic reward for resource node by traffic
	for _, walletTraffic := range trafficList {
		walletAddr, err := sdk.AccAddressFromBech32(walletTraffic.WalletAddress)
		if err != nil {
			continue
		}
		trafficVolume := walletTraffic.Volume

		shareOfTraffic := trafficVolume.ToDec().Quo(totalConsumedNoz)

		// mining reward for resource node
		miningReward := sdk.NewCoin(k.RewardDenom(ctx),
			distributeGoal.MiningRewardToResourceNode.Amount.ToDec().
				Mul(shareOfTraffic).
				TruncateInt())

		// traffic reward for resource node, need to pay community tax
		trafficRewardBeforeTax := distributeGoal.TrafficRewardToResourceNode.Amount.ToDec().
			Mul(shareOfTraffic)
		trafficRewardAfterTax, tax := k.CalcCommunityTax(ctx, trafficRewardBeforeTax)

		// update rewardDetailMap
		if _, ok := rewardDetailMap[walletAddr.String()]; !ok {
			reward := types.NewDefaultReward(walletAddr)
			rewardDetailMap[walletAddr.String()] = reward
		}
		newReward := rewardDetailMap[walletAddr.String()]
		newReward = newReward.AddRewardFromMiningPool(miningReward)
		newReward = newReward.AddRewardFromTrafficPool(trafficRewardAfterTax)
		rewardDetailMap[walletAddr.String()] = newReward

		// record value preparing for transfer
		foundationToReward = foundationToReward.
			Add(miningReward)
		unissuedPrepayToReward = unissuedPrepayToReward.
			Add(trafficRewardAfterTax)
		unissuedPrepayToCommunityPool = unissuedPrepayToCommunityPool.
			Add(tax)
	}

	return rewardDetailMap
}

// Iteration for calculating reward of meta nodes
func (k Keeper) CalcRewardForMetaNode(ctx sdk.Context, distributeGoalBalance types.DistributeGoal, rewardDetailMap map[string]types.Reward,
) map[string]types.Reward {

	totalStakeOfMetaNodes := k.RegisterKeeper.GetMetaNodeBondedToken(ctx).Amount
	metaNodeCnt := k.RegisterKeeper.GetBondedMetaNodeCnt(ctx)

	mataNodeIterator := k.RegisterKeeper.GetMetaNodeIterator(ctx)
	defer mataNodeIterator.Close()

	for ; mataNodeIterator.Valid(); mataNodeIterator.Next() {
		node := regtypes.MustUnmarshalMetaNode(k.cdc, mataNodeIterator.Value())
		if node.Status != stakingtypes.Bonded {
			continue
		}

		walletAddr, err := sdk.AccAddressFromBech32(node.OwnerAddress)
		if err != nil {
			continue
		}
		tokens, ok := sdk.NewIntFromString(node.Tokens.String())
		if !ok {
			continue
		}

		shareOfToken := tokens.ToDec().Quo(totalStakeOfMetaNodes.ToDec())

		// 1, calc stake reward for meta node by stake
		stakeMiningReward := sdk.NewCoin(k.RewardDenom(ctx),
			distributeGoalBalance.StakeMiningRewardToMetaNode.Amount.ToDec().
				Mul(shareOfToken).
				TruncateInt())
		stakeTrafficReward := distributeGoalBalance.StakeTrafficRewardToMetaNode.Amount.ToDec().
			Mul(shareOfToken)

		// 2, calc mining reward for meta node (equally distributed)
		miningRewardToMetaNode := sdk.NewCoin(k.RewardDenom(ctx),
			distributeGoalBalance.MiningRewardToMetaNode.Amount.ToDec().
				Quo(metaNodeCnt.ToDec()).
				TruncateInt())

		// 3, calc traffic reward for meta node (equally distributed)
		trafficRewardToMetaNode := distributeGoalBalance.TrafficRewardToMetaNode.Amount.ToDec().
			Quo(metaNodeCnt.ToDec())

		// reward from traffic pool need to pay community tax
		rewardFromTrafficPoolBeforeTax := stakeTrafficReward.Add(trafficRewardToMetaNode)
		rewardFromTrafficPoolAfterTax, tax := k.CalcCommunityTax(ctx, rewardFromTrafficPoolBeforeTax)

		// update rewardDetailMap
		if _, ok := rewardDetailMap[walletAddr.String()]; !ok {
			reward := types.NewDefaultReward(walletAddr)
			rewardDetailMap[walletAddr.String()] = reward
		}
		newReward := rewardDetailMap[walletAddr.String()]
		newReward = newReward.AddRewardFromMiningPool(stakeMiningReward.Add(miningRewardToMetaNode))
		newReward = newReward.AddRewardFromTrafficPool(rewardFromTrafficPoolAfterTax)
		rewardDetailMap[walletAddr.String()] = newReward

		// record value preparing for transfer
		foundationToReward = foundationToReward.
			Add(stakeMiningReward).
			Add(miningRewardToMetaNode)
		unissuedPrepayToReward = unissuedPrepayToReward.
			Add(rewardFromTrafficPoolAfterTax)
		unissuedPrepayToCommunityPool = unissuedPrepayToCommunityPool.
			Add(tax)
	}

	return rewardDetailMap
}

// Iteration for getting total consumed OZ from traffic
func (k Keeper) GetTotalConsumedNoz(trafficList []*types.SingleWalletVolume) sdk.Int {
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
) (validatorReward sdk.Int, resourceNodeReward sdk.Int, metaNodeReward sdk.Int) {

	validatorBondedTokens := k.StakingKeeper.TotalBondedTokens(ctx).ToDec()
	resourceNodeBondedTokens := k.RegisterKeeper.GetResourceNodeBondedToken(ctx).Amount.ToDec()
	metaNodeBondedTokens := k.RegisterKeeper.GetMetaNodeBondedToken(ctx).Amount.ToDec()

	totalBondedTokens := validatorBondedTokens.Add(resourceNodeBondedTokens).Add(metaNodeBondedTokens)

	validatorReward = totalReward.ToDec().
		Mul(validatorBondedTokens).
		Quo(totalBondedTokens).
		TruncateInt()
	resourceNodeReward = totalReward.ToDec().
		Mul(resourceNodeBondedTokens).
		Quo(totalBondedTokens).
		TruncateInt()
	metaNodeReward = totalReward.ToDec().
		Mul(metaNodeBondedTokens).
		Quo(totalBondedTokens).
		TruncateInt()

	return
}

func (k Keeper) transferTokens(ctx sdk.Context, totalSlashed sdk.Coins) error {
	// [TLC] [FoundationAccount -> feeCollectorPool] Transfer mining reward to fee_pool for validators
	err := k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.FoundationAccount, k.feeCollectorName, sdk.NewCoins(foundationToFeeCollector))
	if err != nil {
		return err
	}
	// [TLC] [TotalUnissuedPrepay -> feeCollectorPool] Transfer traffic reward to fee_pool for validators
	err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, regtypes.TotalUnissuedPrepay, k.feeCollectorName, sdk.NewCoins(unissuedPrepayToFeeCollector))
	if err != nil {
		return err
	}

	// [TLC] [FoundationAccount -> TotalRewardPool] Transfer mining reward to TotalRewardPool for sds nodes
	err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.FoundationAccount, types.TotalRewardPool, sdk.NewCoins(foundationToReward))
	if err != nil {
		return err
	}

	// [TLC] [TotalUnissuedPrepay -> TotalRewardPool] Transfer traffic reward to TotalRewardPool for sds nodes
	err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, regtypes.TotalUnissuedPrepay, types.TotalRewardPool, sdk.NewCoins(unissuedPrepayToReward))
	if err != nil {
		return err
	}

	// [TLC] [TotalRewardPool -> Distribution] Transfer slashed reward to FeePool.CommunityPool
	totalRewardPoolAccAddr := k.AccountKeeper.GetModuleAddress(types.TotalRewardPool)
	err = k.DistrKeeper.FundCommunityPool(ctx, totalSlashed, totalRewardPoolAccAddr)
	if err != nil {
		return err
	}

	// [TLC] [TotalUnissuedPrepay -> Distribution] Transfer tax to FeePool.CommunityPool
	taxCoins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), unissuedPrepayToCommunityPool.TruncateInt()))
	prepayAccAddr := k.AccountKeeper.GetModuleAddress(regtypes.TotalUnissuedPrepay)
	err = k.DistrKeeper.FundCommunityPool(ctx, taxCoins, prepayAccAddr)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) InitVariable(ctx sdk.Context) {
	foundationToFeeCollector = sdk.NewCoin(k.RewardDenom(ctx), sdk.ZeroInt())
	unissuedPrepayToFeeCollector = sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())
	foundationToReward = sdk.NewCoin(k.RewardDenom(ctx), sdk.ZeroInt())
	unissuedPrepayToReward = sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())
	unissuedPrepayToCommunityPool = sdk.ZeroDec()
	distributeGoal = types.InitDistributeGoal()
	rewardDetailMap = make(map[string]types.Reward) //key: wallet address
}

func (k Keeper) CalcCommunityTax(ctx sdk.Context, rewardBeforeTax sdk.Dec) (reward sdk.Coin, tax sdk.Dec) {
	communityTax := k.GetCommunityTax(ctx)
	tax = rewardBeforeTax.Mul(communityTax)
	rewardAmt := rewardBeforeTax.Sub(tax).TruncateInt()
	reward = sdk.NewCoin(k.BondDenom(ctx), rewardAmt)

	return
}
