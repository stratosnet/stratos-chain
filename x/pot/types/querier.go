package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// querier keys
const (
	QueryParams           = "params"
	QueryVolumeReportHash = "volume_report"
)

type PotRewardInfo struct {
	WalletAddress       sdk.AccAddress
	MatureTotalReward   sdk.Coin
	ImmatureTotalReward sdk.Coin
}

// NewPotRewardInfo creates a new instance of PotRewardInfo
func NewPotRewardInfo(
	walletAddress sdk.AccAddress,
	matureTotal sdk.Coin,
	immatureTotal sdk.Coin,
) PotRewardInfo {
	return PotRewardInfo{
		WalletAddress:       walletAddress,
		MatureTotalReward:   matureTotal,
		ImmatureTotalReward: immatureTotal,
	}
}

type QueryPotRewardsByEpochParams struct {
	Page  int
	Limit int
	Epoch sdk.Int
}

// NewQueryPotRewardsByEpochParams creates a new instance of QueryPotRewardsParams
func NewQueryPotRewardsByEpochParams(page, limit int, epoch sdk.Int) QueryPotRewardsByEpochParams {
	return QueryPotRewardsByEpochParams{
		Page:  page,
		Limit: limit,
		Epoch: epoch,
	}
}

type QueryPotRewardsByWalletAddrParams struct {
	Page       int
	Limit      int
	WalletAddr sdk.AccAddress
	Height     int64
}

func NewQueryPotRewardsByWalletAddrParams(page, limit int, walletAddr sdk.AccAddress, height int64) QueryPotRewardsByWalletAddrParams {
	return QueryPotRewardsByWalletAddrParams{
		Page:       page,
		Limit:      limit,
		WalletAddr: walletAddr,
		Height:     height,
	}
}

type ReportInfo struct {
	Epoch     sdk.Int
	Reference string
}

func NewReportInfo(epoch sdk.Int, reference string) ReportInfo {
	return ReportInfo{
		Epoch:     epoch,
		Reference: reference,
	}
}
