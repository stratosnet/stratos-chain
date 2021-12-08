package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	FoundationAccount = "foundation_account"
)

type DistributeGoal struct {
	BlockChainRewardToValidatorFromMiningPool  sdk.Coin `json:"block_chain_reward_to_validator_from_mining_pool" yaml:"block_chain_reward_to_validator_from_mining_pool"`   // 20% mining reward * stakeOfAllValidators / totalStake
	BlockChainRewardToValidatorFromTrafficPool sdk.Coin `json:"block_chain_reward_to_validator_from_traffic_pool" yaml:"block_chain_reward_to_validator_from_traffic_pool"` // 20% traffic reward * stakeOfAllValidators / totalStake

	BlockChainRewardToIndexingNodeFromMiningPool  sdk.Coin `json:"block_chain_reward_to_indexing_node_from_mining_pool" yaml:"block_chain_reward_to_indexing_node_from_mining_pool"`   // 20% mining reward * stakeOfAllIndexingNodes / totalStake
	BlockChainRewardToIndexingNodeFromTrafficPool sdk.Coin `json:"block_chain_reward_to_indexing_node_from_traffic_pool" yaml:"block_chain_reward_to_indexing_node_from_traffic_pool"` // 20% traffic reward * stakeOfAllValidators / totalStake
	MetaNodeRewardToIndexingNodeFromMiningPool    sdk.Coin `json:"meta_node_reward_to_indexing_node_from_mining_pool" yaml:"meta_node_reward_to_indexing_node_from_mining_pool"`       // 20% of mining reward, distribute equally
	MetaNodeRewardToIndexingNodeFromTrafficPool   sdk.Coin `json:"meta_node_reward_to_indexing_node_from_traffic_pool" yaml:"meta_node_reward_to_indexing_node_from_traffic_pool"`     // 20% of traffic reward, distribute equally

	BlockChainRewardToResourceNodeFromMiningPool  sdk.Coin `json:"block_chain_reward_to_resource_node_from_mining_pool" yaml:"block_chain_reward_to_resource_node_from_mining_pool"`   // 20% mining reward * stakeOfAllResourceNodes / totalStake
	BlockChainRewardToResourceNodeFromTrafficPool sdk.Coin `json:"block_chain_reward_to_resource_node_from_traffic_pool" yaml:"block_chain_reward_to_resource_node_from_traffic_pool"` // 20% traffic reward * stakeOfAllValidators / totalStake
	TrafficRewardToResourceNodeFromMiningPool     sdk.Coin `json:"traffic_reward_to_resource_node_from_mining_pool" yaml:"traffic_reward_to_resource_node_from_mining_pool"`           // 60% of mining reward, distribute by traffic contribution
	TrafficRewardToResourceNodeFromTrafficPool    sdk.Coin `json:"traffic_reward_to_resource_node_from_traffic_pool" yaml:"traffic_reward_to_resource_node_from_traffic_pool"`         // 60% of traffic reward, distribute by traffic contribution
}

func NewDistributeGoal(
	blockChainRewardToValidatorFromMiningPool sdk.Coin, blockChainRewardToResourceNodeFromMiningPool sdk.Coin, blockChainRewardToIndexingNodeFromMiningPool sdk.Coin,
	blockChainRewardToValidatorFromTrafficPool sdk.Coin, blockChainRewardToResourceNodeFromTrafficPool sdk.Coin, blockChainRewardToIndexingNodeFromTrafficPool sdk.Coin,
	metaNodeRewardToIndexingNodeFromMiningPool sdk.Coin, metaNodeRewardToIndexingNodeFromTrafficPool sdk.Coin, trafficRewardToResourceNodeFromMiningPool sdk.Coin,
	trafficRewardToResourceNodeFromTrafficPool sdk.Coin) DistributeGoal {
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

func (d DistributeGoal) AddBlockChainRewardToValidatorFromMiningPool(reward sdk.Coin) DistributeGoal {
	if d.BlockChainRewardToValidatorFromMiningPool.IsEqual(sdk.Coin{}) {
		d.BlockChainRewardToValidatorFromMiningPool = reward
	} else {
		d.BlockChainRewardToValidatorFromMiningPool = d.BlockChainRewardToValidatorFromMiningPool.Add(reward)
	}
	return d
}

func (d DistributeGoal) AddBlockChainRewardToResourceNodeFromMiningPool(reward sdk.Coin) DistributeGoal {
	if d.BlockChainRewardToResourceNodeFromMiningPool.IsEqual(sdk.Coin{}) {
		d.BlockChainRewardToResourceNodeFromMiningPool = reward
	} else {
		d.BlockChainRewardToResourceNodeFromMiningPool = d.BlockChainRewardToResourceNodeFromMiningPool.Add(reward)
	}
	return d
}

func (d DistributeGoal) AddBlockChainRewardToIndexingNodeFromMiningPool(reward sdk.Coin) DistributeGoal {
	if d.BlockChainRewardToIndexingNodeFromMiningPool.IsEqual(sdk.Coin{}) {
		d.BlockChainRewardToIndexingNodeFromMiningPool = reward
	} else {
		d.BlockChainRewardToIndexingNodeFromMiningPool = d.BlockChainRewardToIndexingNodeFromMiningPool.Add(reward)
	}
	return d
}
func (d DistributeGoal) AddBlockChainRewardToValidatorFromTrafficPool(reward sdk.Coin) DistributeGoal {
	if d.BlockChainRewardToValidatorFromTrafficPool.IsEqual(sdk.Coin{}) {
		d.BlockChainRewardToValidatorFromTrafficPool = reward
	} else {
		d.BlockChainRewardToValidatorFromTrafficPool = d.BlockChainRewardToValidatorFromTrafficPool.Add(reward)
	}
	return d
}

func (d DistributeGoal) AddBlockChainRewardToResourceNodeFromTrafficPool(reward sdk.Coin) DistributeGoal {
	if d.BlockChainRewardToResourceNodeFromTrafficPool.IsEqual(sdk.Coin{}) {
		d.BlockChainRewardToResourceNodeFromTrafficPool = reward
	} else {
		d.BlockChainRewardToResourceNodeFromTrafficPool = d.BlockChainRewardToResourceNodeFromTrafficPool.Add(reward)
	}
	return d
}

func (d DistributeGoal) AddBlockChainRewardToIndexingNodeFromTrafficPool(reward sdk.Coin) DistributeGoal {
	if d.BlockChainRewardToIndexingNodeFromTrafficPool.IsEqual(sdk.Coin{}) {
		d.BlockChainRewardToIndexingNodeFromTrafficPool = reward
	} else {
		d.BlockChainRewardToIndexingNodeFromTrafficPool = d.BlockChainRewardToIndexingNodeFromTrafficPool.Add(reward)
	}
	return d
}
func (d DistributeGoal) AddMetaNodeRewardToIndexingNodeFromMiningPool(reward sdk.Coin) DistributeGoal {
	if d.MetaNodeRewardToIndexingNodeFromMiningPool.IsEqual(sdk.Coin{}) {
		d.MetaNodeRewardToIndexingNodeFromMiningPool = reward
	} else {
		d.MetaNodeRewardToIndexingNodeFromMiningPool = d.MetaNodeRewardToIndexingNodeFromMiningPool.Add(reward)
	}
	return d
}
func (d DistributeGoal) AddMetaNodeRewardToIndexingNodeFromTrafficPool(reward sdk.Coin) DistributeGoal {
	if d.MetaNodeRewardToIndexingNodeFromTrafficPool.IsEqual(sdk.Coin{}) {
		d.MetaNodeRewardToIndexingNodeFromTrafficPool = reward
	} else {
		d.MetaNodeRewardToIndexingNodeFromTrafficPool = d.MetaNodeRewardToIndexingNodeFromTrafficPool.Add(reward)
	}
	return d
}
func (d DistributeGoal) AddTrafficRewardToResourceNodeFromMiningPool(reward sdk.Coin) DistributeGoal {
	if d.TrafficRewardToResourceNodeFromMiningPool.IsEqual(sdk.Coin{}) {
		d.TrafficRewardToResourceNodeFromMiningPool = reward
	} else {
		d.TrafficRewardToResourceNodeFromMiningPool = d.TrafficRewardToResourceNodeFromMiningPool.Add(reward)
	}
	return d
}
func (d DistributeGoal) AddTrafficRewardToResourceNodeFromTrafficPool(reward sdk.Coin) DistributeGoal {
	if d.TrafficRewardToResourceNodeFromTrafficPool.IsEqual(sdk.Coin{}) {
		d.TrafficRewardToResourceNodeFromTrafficPool = reward
	} else {
		d.TrafficRewardToResourceNodeFromTrafficPool = d.TrafficRewardToResourceNodeFromTrafficPool.Add(reward)
	}
	return d
}

// String returns a human readable string representation of a Reward.
func (d DistributeGoal) String() string {
	return fmt.Sprintf(`DistributeGoal:{
		BlockChainRewardToValidatorFromMiningPool:    	%s
		BlockChainRewardToIndexingNodeFromMiningPool: 	%s
		BlockChainRewardToResourceNodeFromMiningPool: 	%s
		MetaNodeRewardToIndexingNodeFromMiningPool:   	%s
		TrafficRewardToResourceNodeFromMiningPool:    	%s
		BlockChainRewardToValidatorFromTrafficPool:   	%s
		BlockChainRewardToIndexingNodeFromTrafficPool:	%s
		BlockChainRewardToResourceNodeFromTrafficPool:	%s
		MetaNodeRewardToIndexingNodeFromTrafficPool:  	%s
		TrafficRewardToResourceNodeFromTrafficPool:   	%s
	}`,
		d.BlockChainRewardToValidatorFromMiningPool,
		d.BlockChainRewardToIndexingNodeFromMiningPool,
		d.BlockChainRewardToResourceNodeFromMiningPool,
		d.MetaNodeRewardToIndexingNodeFromMiningPool,
		d.TrafficRewardToResourceNodeFromMiningPool,
		d.BlockChainRewardToValidatorFromTrafficPool,
		d.BlockChainRewardToIndexingNodeFromTrafficPool,
		d.BlockChainRewardToResourceNodeFromTrafficPool,
		d.MetaNodeRewardToIndexingNodeFromTrafficPool,
		d.TrafficRewardToResourceNodeFromTrafficPool,
	)
}

type Reward struct {
	WalletAddress         sdk.AccAddress `json:"wallet_address" yaml:"wallet_address"` // account address of node
	RewardFromMiningPool  sdk.Coins      `json:"reward_from_mining_pool" yaml:"reward_from_mining_pool"`
	RewardFromTrafficPool sdk.Coins      `json:"reward_from_traffic_pool" yaml:"reward_from_traffic_pool"`
}

func NewReward(walletAddress sdk.AccAddress, rewardFromMiningPool sdk.Coins, rewardFromTrafficPool sdk.Coins) Reward {
	return Reward{
		WalletAddress:         walletAddress,
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
func (r Reward) String() string {
	return fmt.Sprintf(`Reward:{
		WalletAddress:			%s
  		RewardFromMiningPool:	%s
  		RewardFromTrafficPool:	%s
	}`, r.WalletAddress, r.RewardFromMiningPool.String(), r.RewardFromTrafficPool.String())
}
