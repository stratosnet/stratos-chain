package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryVolumeReport                   = "query_volume_report"
	QueryIndividualRewardsByReportEpoch = "query_pot_individual_rewards_by_report_epoch"
	QueryRewardsByWalletAddr            = "query_pot_rewards_by_wallet_address"
	QuerySlashingByWalletAddr           = "query_pot_slashing_by_wallet_address"
	QueryPotParams                      = "query_pot_params"
	QueryTotalMinedToken                = "query_total_mined_token"
	QueryCirculationSupply              = "query_circulation_supply"
	QueryTotalRewardByEpoch             = "total-reward"
	QueryMetrics                        = "metrics"
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
	Epoch sdkmath.Int
}

// NewQueryIndividualRewardsByEpochParams creates a new instance of QueryIndividualRewardsByReportEpochParams
func NewQueryIndividualRewardsByEpochParams(page, limit int, epoch sdkmath.Int) QueryIndividualRewardsByReportEpochParams {
	return QueryIndividualRewardsByReportEpochParams{
		Page:  page,
		Limit: limit,
		Epoch: epoch,
	}
}

type QueryRewardsByWalletAddrParams struct {
	WalletAddr sdk.AccAddress
	Height     int64
	Epoch      sdkmath.Int
}

func NewQueryRewardsByWalletAddrParams(walletAddr sdk.AccAddress, height int64, epoch sdkmath.Int) QueryRewardsByWalletAddrParams {
	return QueryRewardsByWalletAddrParams{
		WalletAddr: walletAddr,
		Height:     height,
		Epoch:      epoch,
	}
}

type QueryTotalRewardByEpochParams struct {
	Epoch sdkmath.Int
}

func NewQueryTotalRewardByEpochParams(epoch sdkmath.Int) QueryTotalRewardByEpochParams {
	return QueryTotalRewardByEpochParams{
		Epoch: epoch,
	}
}
