package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type SingleNodeVolume struct {
	NodeAddress sdk.AccAddress `json:"node_address" yaml:"node_address"`
	Volume      sdk.Int        `json:"node_volume" yaml:"node_volume"`
}

// NewSingleNodeVolume creates a new Msg<Action> instance
func NewSingleNodeVolume(
	nodeAddress sdk.AccAddress,
	volume sdk.Int,
) SingleNodeVolume {
	return SingleNodeVolume{
		NodeAddress: nodeAddress,
		Volume:      volume,
	}
}

type MiningRewardParam struct {
	TotalMinedValveStart   sdk.Int `json:"total_mined_valve_start" yaml:"total_mined_valve_start"`
	TotalMinedValveEnd     sdk.Int `json:"total_mined_valve_end" yaml:"total_mined_valve_end"`
	MiningReward           sdk.Dec `json:"mining_reward" yaml:"mining_reward"`
	BlockChainPercentage   sdk.Dec `json:"block_chain_percentage" yaml:"block_chain_percentage"`
	ResourceNodePercentage sdk.Dec `json:"resource_node_percentage" yaml:"resource_node_percentage"`
	MetaNodePercentage     sdk.Dec `json:"meta_node_percentage" yaml:"meta_node_percentage"`
}

func NewMiningRewardParam(totalMinedValveStart sdk.Int, totalMinedValveEnd sdk.Int, miningReward sdk.Dec,
	resourceNodePercentage sdk.Dec, metaNodePercentage sdk.Dec, blockChainPercentage sdk.Dec) MiningRewardParam {
	return MiningRewardParam{
		TotalMinedValveStart:   totalMinedValveStart,
		TotalMinedValveEnd:     totalMinedValveEnd,
		MiningReward:           miningReward,
		BlockChainPercentage:   blockChainPercentage,
		ResourceNodePercentage: resourceNodePercentage,
		MetaNodePercentage:     metaNodePercentage,
	}
}
