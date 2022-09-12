package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	FoundationAccount = "foundation_account"
	TotalRewardPool   = "total_reward_pool"
)

type DistributeGoal struct {
	StakeMiningRewardToValidator  sdk.Coin `json:"stake_mining_reward_to_validator" yaml:"stake_mining_reward_to_validator"`   // 20% mining reward * stakeOfAllValidators / totalStake
	StakeTrafficRewardToValidator sdk.Coin `json:"stake_traffic_reward_to_validator" yaml:"stake_traffic_reward_to_validator"` // 20% traffic reward * stakeOfAllValidators / totalStake

	StakeMiningRewardToMetaNode  sdk.Coin `json:"stake_mining_reward_to_meta_node" yaml:"stake_mining_reward_to_meta_node"`   // 20% mining reward * stakeOfAllMetaNodes / totalStake
	StakeTrafficRewardToMetaNode sdk.Coin `json:"stake_traffic_reward_to_meta_node" yaml:"stake_traffic_reward_to_meta_node"` // 20% traffic reward * stakeOfAllValidators / totalStake
	MiningRewardToMetaNode       sdk.Coin `json:"mining_reward_to_meta_node" yaml:"mining_reward_to_meta_node"`               // 20% of mining reward, distribute equally
	TrafficRewardToMetaNode      sdk.Coin `json:"traffic_reward_to_meta_node" yaml:"traffic_reward_to_meta_node"`             // 20% of traffic reward, distribute equally

	StakeMiningRewardToResourceNode  sdk.Coin `json:"stake_mining_reward_to_resource_node" yaml:"stake_mining_reward_to_resource_node"`   // 20% mining reward * stakeOfAllResourceNodes / totalStake
	StakeTrafficRewardToResourceNode sdk.Coin `json:"stake_traffic_reward_to_resource_node" yaml:"stake_traffic_reward_to_resource_node"` // 20% traffic reward * stakeOfAllValidators / totalStake
	MiningRewardToResourceNode       sdk.Coin `json:"mining_reward_to_resource_node" yaml:"mining_reward_to_resource_node"`               // 60% of mining reward, distribute by traffic contribution
	TrafficRewardToResourceNode      sdk.Coin `json:"traffic_reward_to_resource_node" yaml:"traffic_reward_to_resource_node"`             // 60% of traffic reward, distribute by traffic contribution
}

func NewDistributeGoal(
	stakeMiningRewardToValidator sdk.Coin, stakeMiningRewardToResourceNode sdk.Coin, stakeMiningRewardToMetaNode sdk.Coin,
	stakeTrafficRewardToValidator sdk.Coin, stakeTrafficRewardToResourceNode sdk.Coin, stakeTrafficRewardToMetaNode sdk.Coin,
	miningRewardToMetaNode sdk.Coin, trafficRewardToMetaNode sdk.Coin, miningRewardToResourceNode sdk.Coin,
	trafficRewardToResourceNode sdk.Coin) DistributeGoal {
	return DistributeGoal{
		StakeMiningRewardToValidator:     stakeMiningRewardToValidator,
		StakeMiningRewardToResourceNode:  stakeMiningRewardToResourceNode,
		StakeMiningRewardToMetaNode:      stakeMiningRewardToMetaNode,
		StakeTrafficRewardToValidator:    stakeTrafficRewardToValidator,
		StakeTrafficRewardToResourceNode: stakeTrafficRewardToResourceNode,
		StakeTrafficRewardToMetaNode:     stakeTrafficRewardToMetaNode,
		MiningRewardToMetaNode:           miningRewardToMetaNode,
		TrafficRewardToMetaNode:          trafficRewardToMetaNode,
		MiningRewardToResourceNode:       miningRewardToResourceNode,
		TrafficRewardToResourceNode:      trafficRewardToResourceNode,
	}
}

func InitDistributeGoal() DistributeGoal {
	return NewDistributeGoal(
		sdk.Coin{},
		sdk.Coin{},
		sdk.Coin{},
		sdk.Coin{},
		sdk.Coin{},
		sdk.Coin{},
		sdk.Coin{},
		sdk.Coin{},
		sdk.Coin{},
		sdk.Coin{},
	)
}

func (d DistributeGoal) AddStakeMiningRewardToValidator(reward sdk.Coin) DistributeGoal {
	if d.StakeMiningRewardToValidator.IsEqual(sdk.Coin{}) {
		d.StakeMiningRewardToValidator = reward
	} else {
		d.StakeMiningRewardToValidator = d.StakeMiningRewardToValidator.Add(reward)
	}
	return d
}

func (d DistributeGoal) AddStakeMiningRewardToResourceNode(reward sdk.Coin) DistributeGoal {
	if d.StakeMiningRewardToResourceNode.IsEqual(sdk.Coin{}) {
		d.StakeMiningRewardToResourceNode = reward
	} else {
		d.StakeMiningRewardToResourceNode = d.StakeMiningRewardToResourceNode.Add(reward)
	}
	return d
}

func (d DistributeGoal) AddStakeMiningRewardToMetaNode(reward sdk.Coin) DistributeGoal {
	if d.StakeMiningRewardToMetaNode.IsEqual(sdk.Coin{}) {
		d.StakeMiningRewardToMetaNode = reward
	} else {
		d.StakeMiningRewardToMetaNode = d.StakeMiningRewardToMetaNode.Add(reward)
	}
	return d
}
func (d DistributeGoal) AddStakeTrafficRewardToValidator(reward sdk.Coin) DistributeGoal {
	if d.StakeTrafficRewardToValidator.IsEqual(sdk.Coin{}) {
		d.StakeTrafficRewardToValidator = reward
	} else {
		d.StakeTrafficRewardToValidator = d.StakeTrafficRewardToValidator.Add(reward)
	}
	return d
}

func (d DistributeGoal) AddStakeTrafficRewardToResourceNode(reward sdk.Coin) DistributeGoal {
	if d.StakeTrafficRewardToResourceNode.IsEqual(sdk.Coin{}) {
		d.StakeTrafficRewardToResourceNode = reward
	} else {
		d.StakeTrafficRewardToResourceNode = d.StakeTrafficRewardToResourceNode.Add(reward)
	}
	return d
}

func (d DistributeGoal) AddStakeTrafficRewardToMetaNode(reward sdk.Coin) DistributeGoal {
	if d.StakeTrafficRewardToMetaNode.IsEqual(sdk.Coin{}) {
		d.StakeTrafficRewardToMetaNode = reward
	} else {
		d.StakeTrafficRewardToMetaNode = d.StakeTrafficRewardToMetaNode.Add(reward)
	}
	return d
}
func (d DistributeGoal) AddMiningRewardToMetaNode(reward sdk.Coin) DistributeGoal {
	if d.MiningRewardToMetaNode.IsEqual(sdk.Coin{}) {
		d.MiningRewardToMetaNode = reward
	} else {
		d.MiningRewardToMetaNode = d.MiningRewardToMetaNode.Add(reward)
	}
	return d
}
func (d DistributeGoal) AddTrafficRewardToMetaNode(reward sdk.Coin) DistributeGoal {
	if d.TrafficRewardToMetaNode.IsEqual(sdk.Coin{}) {
		d.TrafficRewardToMetaNode = reward
	} else {
		d.TrafficRewardToMetaNode = d.TrafficRewardToMetaNode.Add(reward)
	}
	return d
}
func (d DistributeGoal) AddMiningRewardToResourceNode(reward sdk.Coin) DistributeGoal {
	if d.MiningRewardToResourceNode.IsEqual(sdk.Coin{}) {
		d.MiningRewardToResourceNode = reward
	} else {
		d.MiningRewardToResourceNode = d.MiningRewardToResourceNode.Add(reward)
	}
	return d
}
func (d DistributeGoal) AddTrafficRewardToResourceNode(reward sdk.Coin) DistributeGoal {
	if d.TrafficRewardToResourceNode.IsEqual(sdk.Coin{}) {
		d.TrafficRewardToResourceNode = reward
	} else {
		d.TrafficRewardToResourceNode = d.TrafficRewardToResourceNode.Add(reward)
	}
	return d
}

// String returns a human readable string representation of a Reward.
func (d DistributeGoal) String() string {
	return fmt.Sprintf(`DistributeGoal:{
		StakeMiningRewardToValidator:       %s
		StakeTrafficRewardToValidator:      %s
		StakeMiningRewardToMetaNode:        %s
		StakeTrafficRewardToMetaNode:       %s
		MiningRewardToMetaNode:             %s
		TrafficRewardToMetaNode:            %s
		StakeMiningRewardToResourceNode:    %s
		StakeTrafficRewardToResourceNode:   %s
		MiningRewardToResourceNode:         %s
		TrafficRewardToResourceNode:        %s
	}`,
		d.StakeMiningRewardToValidator,
		d.StakeTrafficRewardToValidator,
		d.StakeMiningRewardToMetaNode,
		d.StakeTrafficRewardToMetaNode,
		d.MiningRewardToMetaNode,
		d.TrafficRewardToMetaNode,
		d.StakeMiningRewardToResourceNode,
		d.StakeTrafficRewardToResourceNode,
		d.MiningRewardToResourceNode,
		d.TrafficRewardToResourceNode,
	)
}

func NewReward(walletAddress sdk.AccAddress, rewardFromMiningPool sdk.Coins, rewardFromTrafficPool sdk.Coins) Reward {
	return Reward{
		WalletAddress:         walletAddress.String(),
		RewardFromMiningPool:  rewardFromMiningPool,
		RewardFromTrafficPool: rewardFromTrafficPool,
	}
}

func NewDefaultReward(walletAddress sdk.AccAddress) Reward {
	return NewReward(walletAddress, sdk.Coins{}, sdk.Coins{})
}

func (r Reward) AddRewardFromMiningPool(reward sdk.Coin) Reward {
	r.RewardFromMiningPool = r.RewardFromMiningPool.Add(reward)
	return r
}

func (r Reward) AddRewardFromTrafficPool(reward sdk.Coin) Reward {
	r.RewardFromTrafficPool = r.RewardFromTrafficPool.Add(reward)
	return r
}

// String returns a human readable string representation of a Reward.
func (r Reward) HrpString() string {
	return fmt.Sprintf(`Reward:{
		WalletAddress:			%s
  		RewardFromMiningPool:	%s
  		RewardFromTrafficPool:	%s
	}`, r.WalletAddress, r.RewardFromMiningPool.String(), r.RewardFromTrafficPool.String())
}
