package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

// NewSingleWalletVolume creates a new Msg<Action> instance
func NewSingleWalletVolume(
	walletAddress sdk.AccAddress,
	volume sdk.Int,
) *SingleWalletVolume {
	return &SingleWalletVolume{
		WalletAddress: walletAddress.String(),
		Volume:        &volume,
	}
}

func NewMiningRewardParam(totalMinedValveStart sdk.Coin, totalMinedValveEnd sdk.Coin, miningReward sdk.Coin,
	resourceNodePercentageInTenThousand sdk.Int, metaNodePercentageInTenThousand sdk.Int, blockChainPercentageInTenThousand sdk.Int) *MiningRewardParam {
	return &MiningRewardParam{
		TotalMinedValveStart:                &totalMinedValveStart,
		TotalMinedValveEnd:                  &totalMinedValveEnd,
		MiningReward:                        &miningReward,
		BlockChainPercentageInTenThousand:   &blockChainPercentageInTenThousand,
		ResourceNodePercentageInTenThousand: &resourceNodePercentageInTenThousand,
		MetaNodePercentageInTenThousand:     &metaNodePercentageInTenThousand,
	}
}

func NewReportRecord(reporter stratos.SdsAddress, reportReference string, txHash string) VolumeReportRecord {
	return VolumeReportRecord{
		Reporter:        reporter.String(),
		ReportReference: reportReference,
		TxHash:          txHash,
	}
}

func NewBLSSignatureInfo(pubKeys [][]byte, signature []byte, txData []byte) BLSSignatureInfo {
	return BLSSignatureInfo{
		PubKeys:   pubKeys,
		Signature: signature,
		TxData:    txData,
	}
}

type BaseBLSSignatureInfo struct {
	PubKeys   []string `json:"pub_keys" yaml:"pub_keys"`
	Signature string   `json:"signature" yaml:"signature"`
	TxData    string   `json:"tx_data" yaml:"tx_data"`
}
