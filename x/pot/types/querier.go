package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryVolumeReport                   = "query_volume_report"
	QueryIndividualRewardsByReportEpoch = "query_pot_individual_rewards_by_report_epoch"
	QueryRewardsByWalletAddr            = "query_pot_rewards_by_wallet_address"
	QuerySlashingByWalletAddr           = "query_pot_slashing_by_wallet_address"
	QueryPotParams                      = "query_pot_params"
	QueryTotalMinedToken                = "query_total_mined_token"
	QueryDefaultLimit                   = 100
)

// NewRewardInfo creates a new instance of PotRewardInfo
func NewRewardInfo(
	walletAddress sdk.AccAddress,
	matureTotal sdk.Coins,
	immatureTotal sdk.Coins,
) RewardByOwner {
	return RewardByOwner{
		WalletAddress:       walletAddress.String(),
		MatureTotalReward:   matureTotal,
		ImmatureTotalReward: immatureTotal,
	}
}

type QueryIndividualRewardsByReportEpochParams struct {
	Page  int
	Limit int
	Epoch sdk.Int
}

// NewQueryIndividualRewardsByEpochParams creates a new instance of QueryIndividualRewardsByReportEpochParams
func NewQueryIndividualRewardsByEpochParams(page, limit int, epoch sdk.Int) QueryIndividualRewardsByReportEpochParams {
	return QueryIndividualRewardsByReportEpochParams{
		Page:  page,
		Limit: limit,
		Epoch: epoch,
	}
}

type QueryRewardsByWalletAddrParams struct {
	WalletAddr sdk.AccAddress
	Height     int64
	Epoch      sdk.Int
}

func NewQueryRewardsByWalletAddrParams(walletAddr sdk.AccAddress, height int64, epoch sdk.Int) QueryRewardsByWalletAddrParams {
	return QueryRewardsByWalletAddrParams{
		WalletAddr: walletAddr,
		Height:     height,
		Epoch:      epoch,
	}
}
