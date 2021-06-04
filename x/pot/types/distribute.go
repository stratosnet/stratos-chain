package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DistributeGoal struct {
	BlockChainRewardToValidatorFromMiningPool  sdk.Int `json:"block_chain_reward_to_validator_from_mining_pool" yaml:"block_chain_reward_to_validator_from_mining_pool"`   // 20% mining reward * stakeOfAllValidators / totalStake
	BlockChainRewardToValidatorFromTrafficPool sdk.Int `json:"block_chain_reward_to_validator_from_traffic_pool" yaml:"block_chain_reward_to_validator_from_traffic_pool"` // 20% traffic reward * stakeOfAllValidators / totalStake

	BlockChainRewardToIndexingNodeFromMiningPool  sdk.Int `json:"block_chain_reward_to_indexing_node_from_mining_pool" yaml:"block_chain_reward_to_indexing_node_from_mining_pool"`   // 20% mining reward * stakeOfAllIndexingNodes / totalStake
	BlockChainRewardToIndexingNodeFromTrafficPool sdk.Int `json:"block_chain_reward_to_indexing_node_from_traffic_pool" yaml:"block_chain_reward_to_indexing_node_from_traffic_pool"` // 20% traffic reward * stakeOfAllValidators / totalStake
	MetaNodeRewardToIndexingNodeFromMiningPool    sdk.Int `json:"meta_node_reward_to_indexing_node_from_mining_pool" yaml:"meta_node_reward_to_indexing_node_from_mining_pool"`       // 20% of mining reward, distribute equally
	MetaNodeRewardToIndexingNodeFromTrafficPool   sdk.Int `json:"meta_node_reward_to_indexing_node_from_traffic_pool" yaml:"meta_node_reward_to_indexing_node_from_traffic_pool"`     // 20% of traffic reward, distribute equally

	BlockChainRewardToResourceNodeFromMiningPool  sdk.Int `json:"block_chain_reward_to_resource_node_from_mining_pool" yaml:"block_chain_reward_to_resource_node_from_mining_pool"`   // 20% mining reward * stakeOfAllResourceNodes / totalStake
	BlockChainRewardToResourceNodeFromTrafficPool sdk.Int `json:"block_chain_reward_to_resource_node_from_traffic_pool" yaml:"block_chain_reward_to_resource_node_from_traffic_pool"` // 20% traffic reward * stakeOfAllValidators / totalStake
	TrafficRewardToResourceNodeFromMiningPool     sdk.Int `json:"traffic_reward_to_resource_node_from_mining_pool" yaml:"traffic_reward_to_resource_node_from_mining_pool"`           // 60% of mining reward, distribute by traffic contribution
	TrafficRewardToResourceNodeFromTrafficPool    sdk.Int `json:"traffic_reward_to_resource_node_from_traffic_pool" yaml:"traffic_reward_to_resource_node_from_traffic_pool"`         // 60% of traffic reward, distribute by traffic contribution
}

func NewDistributeGoal(
	blockChainRewardToValidatorFromMiningPool sdk.Int, blockChainRewardToResourceNodeFromMiningPool sdk.Int, blockChainRewardToIndexingNodeFromMiningPool sdk.Int,
	blockChainRewardToValidatorFromTrafficPool sdk.Int, blockChainRewardToResourceNodeFromTrafficPool sdk.Int, blockChainRewardToIndexingNodeFromTrafficPool sdk.Int,
	metaNodeRewardToIndexingNodeFromMiningPool sdk.Int, metaNodeRewardToIndexingNodeFromTrafficPool sdk.Int, trafficRewardToResourceNodeFromMiningPool sdk.Int,
	trafficRewardToResourceNodeFromTrafficPool sdk.Int) DistributeGoal {
	return DistributeGoal{
		BlockChainRewardToValidatorFromMiningPool:     blockChainRewardToValidatorFromMiningPool,
		BlockChainRewardToResourceNodeFromMiningPool:  blockChainRewardToResourceNodeFromMiningPool,
		BlockChainRewardToIndexingNodeFromMiningPool:  blockChainRewardToIndexingNodeFromMiningPool,
		BlockChainRewardToValidatorFromTrafficPool:    blockChainRewardToValidatorFromTrafficPool,
		BlockChainRewardToResourceNodeFromTrafficPool: blockChainRewardToResourceNodeFromTrafficPool,
		BlockChainRewardToIndexingNodeFromTrafficPool: blockChainRewardToIndexingNodeFromTrafficPool,
		MetaNodeRewardToIndexingNodeFromMiningPool:    metaNodeRewardToIndexingNodeFromMiningPool,
		MetaNodeRewardToIndexingNodeFromTrafficPool:   metaNodeRewardToIndexingNodeFromTrafficPool,
		TrafficRewardToResourceNodeFromMiningPool:     trafficRewardToResourceNodeFromMiningPool,
		TrafficRewardToResourceNodeFromTrafficPool:    trafficRewardToResourceNodeFromTrafficPool,
	}
}

func InitDistributeGoal() DistributeGoal {
	return NewDistributeGoal(sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt())
}

func (d DistributeGoal) AddBlockChainRewardToValidatorFromMiningPool(reward sdk.Int) DistributeGoal {
	d.BlockChainRewardToValidatorFromMiningPool = d.BlockChainRewardToValidatorFromMiningPool.Add(reward)
	return d
}

func (d DistributeGoal) AddBlockChainRewardToResourceNodeFromMiningPool(reward sdk.Int) DistributeGoal {
	d.BlockChainRewardToResourceNodeFromMiningPool = d.BlockChainRewardToResourceNodeFromMiningPool.Add(reward)
	return d
}

func (d DistributeGoal) AddBlockChainRewardToIndexingNodeFromMiningPool(reward sdk.Int) DistributeGoal {
	d.BlockChainRewardToIndexingNodeFromMiningPool = d.BlockChainRewardToIndexingNodeFromMiningPool.Add(reward)
	return d
}
func (d DistributeGoal) AddBlockChainRewardToValidatorFromTrafficPool(reward sdk.Int) DistributeGoal {
	d.BlockChainRewardToValidatorFromTrafficPool = d.BlockChainRewardToValidatorFromTrafficPool.Add(reward)
	return d
}

func (d DistributeGoal) AddBlockChainRewardToResourceNodeFromTrafficPool(reward sdk.Int) DistributeGoal {
	d.BlockChainRewardToResourceNodeFromTrafficPool = d.BlockChainRewardToResourceNodeFromTrafficPool.Add(reward)
	return d
}

func (d DistributeGoal) AddBlockChainRewardToIndexingNodeFromTrafficPool(reward sdk.Int) DistributeGoal {
	d.BlockChainRewardToIndexingNodeFromTrafficPool = d.BlockChainRewardToIndexingNodeFromTrafficPool.Add(reward)
	return d
}
func (d DistributeGoal) AddMetaNodeRewardToIndexingNodeFromMiningPool(reward sdk.Int) DistributeGoal {
	d.MetaNodeRewardToIndexingNodeFromMiningPool = d.MetaNodeRewardToIndexingNodeFromMiningPool.Add(reward)
	return d
}
func (d DistributeGoal) AddMetaNodeRewardToIndexingNodeFromTrafficPool(reward sdk.Int) DistributeGoal {
	d.MetaNodeRewardToIndexingNodeFromTrafficPool = d.MetaNodeRewardToIndexingNodeFromTrafficPool.Add(reward)
	return d
}
func (d DistributeGoal) AddTrafficRewardToResourceNodeFromMiningPool(reward sdk.Int) DistributeGoal {
	d.TrafficRewardToResourceNodeFromMiningPool = d.TrafficRewardToResourceNodeFromMiningPool.Add(reward)
	return d
}
func (d DistributeGoal) AddTrafficRewardToResourceNodeFromTrafficPool(reward sdk.Int) DistributeGoal {
	d.TrafficRewardToResourceNodeFromTrafficPool = d.TrafficRewardToResourceNodeFromTrafficPool.Add(reward)
	return d
}

// String returns a human readable string representation of a Reward.
func (d DistributeGoal) String() string {
	return fmt.Sprintf(`DistributeGoal:{
		BlockChainRewardToValidatorFromMiningPool:    	%s
		BlockChainRewardToValidatorFromTrafficPool:   	%s
		BlockChainRewardToIndexingNodeFromMiningPool: 	%s
		BlockChainRewardToIndexingNodeFromTrafficPool:	%s
		MetaNodeRewardToIndexingNodeFromMiningPool:   	%s
		MetaNodeRewardToIndexingNodeFromTrafficPool:  	%s
		BlockChainRewardToResourceNodeFromMiningPool: 	%s
		BlockChainRewardToResourceNodeFromTrafficPool:	%s
		TrafficRewardToResourceNodeFromMiningPool:    	%s
		TrafficRewardToResourceNodeFromTrafficPool:   	%s
	}`, d.BlockChainRewardToValidatorFromMiningPool,
		d.BlockChainRewardToValidatorFromTrafficPool,
		d.BlockChainRewardToIndexingNodeFromMiningPool,
		d.BlockChainRewardToIndexingNodeFromTrafficPool,
		d.MetaNodeRewardToIndexingNodeFromMiningPool,
		d.MetaNodeRewardToIndexingNodeFromTrafficPool,
		d.BlockChainRewardToResourceNodeFromMiningPool,
		d.BlockChainRewardToResourceNodeFromTrafficPool,
		d.TrafficRewardToResourceNodeFromMiningPool,
		d.TrafficRewardToResourceNodeFromTrafficPool,
	)
}

type Reward struct {
	NodeAddress           sdk.AccAddress `json:"node_address" yaml:"node_address"` // account address of node
	RewardFromMiningPool  sdk.Int        `json:"reward_from_mining_pool" yaml:"reward_from_mining_pool"`
	RewardFromTrafficPool sdk.Int        `json:"reward_from_traffic_pool" yaml:"reward_from_traffic_pool"`
}

func NewReward(nodeAddress sdk.AccAddress, rewardFromMiningPool sdk.Int, rewardFromTrafficPool sdk.Int) Reward {
	return Reward{
		NodeAddress:           nodeAddress,
		RewardFromMiningPool:  rewardFromMiningPool,
		RewardFromTrafficPool: rewardFromTrafficPool,
	}
}

func NewDefaultReward(nodeAddress sdk.AccAddress) Reward {
	return NewReward(nodeAddress, sdk.ZeroInt(), sdk.ZeroInt())
}

func (r Reward) AddRewardFromMiningPool(reward sdk.Int) Reward {
	r.RewardFromMiningPool = r.RewardFromMiningPool.Add(reward)
	return r
}

func (r Reward) AddRewardFromTrafficPool(reward sdk.Int) Reward {
	r.RewardFromTrafficPool = r.RewardFromTrafficPool.Add(reward)
	return r
}

// String returns a human readable string representation of a Reward.
func (r Reward) String() string {
	return fmt.Sprintf(`Reward:{
		NodeAddress:			%s
  		RewardFromMiningPool:	%s
  		RewardFromTrafficPool:	%s
	}`, r.NodeAddress, r.RewardFromMiningPool, r.RewardFromTrafficPool)
}
