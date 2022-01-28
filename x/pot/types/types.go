package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type SingleWalletVolume struct {
	WalletAddress sdk.AccAddress `json:"wallet_address" yaml:"wallet_address"`
	Volume        sdk.Int        `json:"volume" yaml:"volume"` //uoz
}

// NewSingleWalletVolume creates a new Msg<Action> instance
func NewSingleWalletVolume(
	walletAddress sdk.AccAddress,
	volume sdk.Int,
) SingleWalletVolume {
	return SingleWalletVolume{
		WalletAddress: walletAddress,
		Volume:        volume,
	}
}

type MiningRewardParam struct {
	TotalMinedValveStart                sdk.Coin `json:"total_mined_valve_start" yaml:"total_mined_valve_start"`
	TotalMinedValveEnd                  sdk.Coin `json:"total_mined_valve_end" yaml:"total_mined_valve_end"`
	MiningReward                        sdk.Coin `json:"mining_reward" yaml:"mining_reward"`
	BlockChainPercentageInTenThousand   sdk.Int  `json:"block_chain_percentage_in_ten_thousand" yaml:"block_chain_percentage_in_ten_thousand"`
	ResourceNodePercentageInTenThousand sdk.Int  `json:"resource_node_percentage_in_ten_thousand" yaml:"resource_node_percentage_in_ten_thousand"`
	MetaNodePercentageInTenThousand     sdk.Int  `json:"meta_node_percentage_in_ten_thousand" yaml:"meta_node_percentage_in_ten_thousand"`
}

func NewMiningRewardParam(totalMinedValveStart sdk.Coin, totalMinedValveEnd sdk.Coin, miningReward sdk.Coin,
	resourceNodePercentageInTenThousand sdk.Int, metaNodePercentageInTenThousand sdk.Int, blockChainPercentageInTenThousand sdk.Int) MiningRewardParam {
	return MiningRewardParam{
		TotalMinedValveStart:                totalMinedValveStart,
		TotalMinedValveEnd:                  totalMinedValveEnd,
		MiningReward:                        miningReward,
		BlockChainPercentageInTenThousand:   blockChainPercentageInTenThousand,
		ResourceNodePercentageInTenThousand: resourceNodePercentageInTenThousand,
		MetaNodePercentageInTenThousand:     metaNodePercentageInTenThousand,
	}
}

type BLSSignatureInfo struct {
	PubKeys   [][]byte `json:"pub_keys" yaml:"pub_keys"`
	Signature []byte   `json:"signature" yaml:"signature"`
	TxData    []byte   `json:"tx_data" yaml:"tx_data"`
}

func NewBLSSignatureInfo(pubKeys [][]byte, signature []byte, txData []byte) BLSSignatureInfo {
	return BLSSignatureInfo{
		PubKeys:   pubKeys,
		Signature: signature,
		TxData:    txData,
	}
}
