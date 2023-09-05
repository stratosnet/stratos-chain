package keeper

import (
	"errors"
	"sort"
	"strings"

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

func (k Keeper) DistributePotReward(ctx sdk.Context, trafficList []types.SingleWalletVolume, epoch sdk.Int) (err error) {

	k.InitVariable(ctx)

	//1, calc traffic reward in total
	totalConsumedNoz := k.GetTotalConsumedNoz(trafficList).ToDec()
	remaining, total := k.NozSupply(ctx)
	if totalConsumedNoz.Add(remaining.ToDec()).GT(total.ToDec()) {
		return errors.New("remaining+consumed Noz exceeds total Noz supply")
	}

	distributeGoal, err = k.CalcTrafficRewardInTotal(ctx, distributeGoal, totalConsumedNoz)
	if err != nil {
		return err
	}
	unissuedPrepayToFeeCollector = distributeGoal.StakeTrafficRewardToValidator

	//2, calc mining reward in total
	distributeGoal, err = k.CalcMiningRewardInTotal(ctx, distributeGoal)
	if err != nil && err != types.ErrOutOfIssuance {
		return err
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
		return err
	}

	//7, update remaining ozone limit
	remainingNozLimit := k.registerKeeper.GetRemainingOzoneLimit(ctx)
	k.registerKeeper.SetRemainingOzoneLimit(ctx, remainingNozLimit.Add(totalConsumedNoz.TruncateInt()))

	//8, [TLC] transfer balance of miningReward&trafficReward pools to totalReward&totalSlashed pool, utilized for future Withdraw Tx
	err = k.transferTokensForDistribution(ctx)
	if err != nil {
		return err
	}

	return nil
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
		Mul(miningParam.BlockChainPercentageInBp.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	trafficRewardToResourceNodes := totalTrafficReward.
		Mul(miningParam.ResourceNodePercentageInBp.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	trafficRewardToMetaNodes := totalTrafficReward.
		Mul(miningParam.MetaNodePercentageInBp.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()

	// all stake reward distribute to validators
	distributeGoal = distributeGoal.AddStakeTrafficRewardToValidator(sdk.NewCoin(k.BondDenom(ctx), stakeTrafficReward))
	distributeGoal = distributeGoal.AddTrafficRewardToResourceNode(sdk.NewCoin(k.BondDenom(ctx), trafficRewardToResourceNodes))
	distributeGoal = distributeGoal.AddTrafficRewardToMetaNode(sdk.NewCoin(k.BondDenom(ctx), trafficRewardToMetaNodes))

	return distributeGoal, nil
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
		Mul(miningParam.BlockChainPercentageInBp.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	miningRewardToResourceNodes := totalMiningReward.Amount.ToDec().
		Mul(miningParam.ResourceNodePercentageInBp.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()
	miningRewardToMetaNodes := totalMiningReward.Amount.ToDec().
		Mul(miningParam.MetaNodePercentageInBp.ToDec()).
		Quo(sdk.NewDec(10000)).TruncateInt()

	// all stake reward distribute to validators
	distributeGoal = distributeGoal.AddStakeMiningRewardToValidator(sdk.NewCoin(k.RewardDenom(ctx), stakeMiningReward))
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
	k.SetLastDistributedEpoch(ctx, currentEpoch)
	return nil
}

func (k Keeper) addNewIndividualAndUpdateImmatureTotal(ctx sdk.Context, account sdk.AccAddress, matureEpoch sdk.Int, newReward types.Reward) {
	newIndividualTotal := newReward.RewardFromMiningPool.Add(newReward.RewardFromTrafficPool...)
	oldImmatureTotal := k.GetImmatureTotalReward(ctx, account)
	newImmatureTotal := oldImmatureTotal.Add(newIndividualTotal...)

	k.SetIndividualReward(ctx, account, matureEpoch, newReward)
	k.SetImmatureTotalReward(ctx, account, newImmatureTotal)
}

// reward will mature 14 days since distribution. Each epoch interval is about 10 minutes.
func (k Keeper) getMatureEpochByCurrentEpoch(ctx sdk.Context, currentEpoch sdk.Int) (matureEpoch sdk.Int) {
	// 14 days = 20160 minutes = 2016 epochs
	paramMatureEpoch := sdk.NewInt(k.MatureEpoch(ctx))
	matureEpoch = paramMatureEpoch.Add(currentEpoch)
	return matureEpoch
}

// Iteration for calculating reward of resource nodes
func (k Keeper) CalcRewardForResourceNode(ctx sdk.Context, totalConsumedNoz sdk.Dec, trafficList []types.SingleWalletVolume,
	distributeGoal types.DistributeGoal, rewardDetailMap map[string]types.Reward,
) map[string]types.Reward {

	// calc mining & traffic reward for resource node by traffic
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

	metaNodeCnt := k.registerKeeper.GetBondedMetaNodeCnt(ctx)

	mataNodeIterator := k.registerKeeper.GetMetaNodeIterator(ctx)
	defer mataNodeIterator.Close()

	for ; mataNodeIterator.Valid(); mataNodeIterator.Next() {
		node := regtypes.MustUnmarshalMetaNode(k.cdc, mataNodeIterator.Value())
		if node.Status != stakingtypes.Bonded {
			continue
		}

		var walletAddrStr string
		if len(strings.TrimSpace(node.BeneficiaryAddress)) == 0 {
			walletAddrStr = node.OwnerAddress
		} else {
			walletAddrStr = node.BeneficiaryAddress
		}

		walletAddr, err := sdk.AccAddressFromBech32(walletAddrStr)
		if err != nil {
			continue
		}

		// 1, calc mining reward for meta node (equally distributed)
		miningRewardToMetaNode := sdk.NewCoin(k.RewardDenom(ctx),
			distributeGoalBalance.MiningRewardToMetaNode.Amount.ToDec().
				Quo(metaNodeCnt.ToDec()).
				TruncateInt())

		// 2, calc traffic reward for meta node (equally distributed)
		trafficRewardToMetaNode := distributeGoalBalance.TrafficRewardToMetaNode.Amount.ToDec().
			Quo(metaNodeCnt.ToDec())

		// reward from traffic pool need to pay community tax
		rewardFromTrafficPoolAfterTax, tax := k.CalcCommunityTax(ctx, trafficRewardToMetaNode)

		// update rewardDetailMap
		if _, ok := rewardDetailMap[walletAddr.String()]; !ok {
			reward := types.NewDefaultReward(walletAddr)
			rewardDetailMap[walletAddr.String()] = reward
		}
		newReward := rewardDetailMap[walletAddr.String()]
		newReward = newReward.AddRewardFromMiningPool(miningRewardToMetaNode)
		newReward = newReward.AddRewardFromTrafficPool(rewardFromTrafficPoolAfterTax)
		rewardDetailMap[walletAddr.String()] = newReward

		// record value preparing for transfer
		foundationToReward = foundationToReward.
			Add(miningRewardToMetaNode)
		unissuedPrepayToReward = unissuedPrepayToReward.
			Add(rewardFromTrafficPoolAfterTax)
		unissuedPrepayToCommunityPool = unissuedPrepayToCommunityPool.
			Add(tax)
	}

	return rewardDetailMap
}

// Iteration for getting total consumed OZ from traffic
func (k Keeper) GetTotalConsumedNoz(trafficList []types.SingleWalletVolume) sdk.Int {
	totalTraffic := sdk.ZeroInt()
	for _, vol := range trafficList {
		toAdd := vol.Volume
		totalTraffic = totalTraffic.Add(toAdd)
	}
	return totalTraffic
}

func (k Keeper) transferTokensForDistribution(ctx sdk.Context) error {

	// [TLC] [FoundationAccount -> feeCollectorPool] Transfer mining reward to fee_pool for validators
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.FoundationAccount, k.feeCollectorName, sdk.NewCoins(foundationToFeeCollector))
	if err != nil {
		return err
	}

	// [TLC] [TotalUnissuedPrepay -> feeCollectorPool] Transfer traffic reward to fee_pool for validators
	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, regtypes.TotalUnissuedPrepay, k.feeCollectorName, sdk.NewCoins(unissuedPrepayToFeeCollector))
	if err != nil {
		return err
	}

	// [TLC] [FoundationAccount -> TotalRewardPool] Transfer mining reward to TotalRewardPool for sds nodes
	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.FoundationAccount, types.TotalRewardPool, sdk.NewCoins(foundationToReward))
	if err != nil {
		return err
	}

	// [TLC] [TotalUnissuedPrepay -> TotalRewardPool] Transfer traffic reward to TotalRewardPool for sds nodes
	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, regtypes.TotalUnissuedPrepay, types.TotalRewardPool, sdk.NewCoins(unissuedPrepayToReward))
	if err != nil {
		return err
	}

	// [TLC] [TotalUnissuedPrepay -> Distribution] Transfer tax to FeePool.CommunityPool
	taxCoins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), unissuedPrepayToCommunityPool.TruncateInt()))
	prepayAccAddr := k.accountKeeper.GetModuleAddress(regtypes.TotalUnissuedPrepay)
	err = k.distrKeeper.FundCommunityPool(ctx, taxCoins, prepayAccAddr)
	if err != nil {
		return err
	}

	return nil
}

// Iteration for sorting map to slice
func sortDetailMapToSlice(rewardDetailMap map[string]types.Reward) (rewardDetailList []types.Reward) {
	keys := make([]string, 0, len(rewardDetailMap))
	for key := range rewardDetailMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		reward := rewardDetailMap[key]
		rewardDetailList = append(rewardDetailList, reward)
	}
	return rewardDetailList
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
