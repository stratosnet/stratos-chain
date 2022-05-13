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
	MatureTotalReward   sdk.Coins
	ImmatureTotalReward sdk.Coins
}

// NewPotRewardInfo creates a new instance of PotRewardInfo
func NewPotRewardInfo(
	walletAddress sdk.AccAddress,
	matureTotal sdk.Coins,
	immatureTotal sdk.Coins,
) PotRewardInfo {
	return PotRewardInfo{
		WalletAddress:       walletAddress,
		MatureTotalReward:   matureTotal,
		ImmatureTotalReward: immatureTotal,
	}
}

type QueryPotRewardsByReportEpochParams struct {
	Page          int
	Limit         int
	Epoch         sdk.Int
	WalletAddress sdk.AccAddress
}

// NewQueryPotRewardsByEpochParams creates a new instance of QueryPotRewardsParams
func NewQueryPotRewardsByEpochParams(page, limit int, epoch sdk.Int, walletAddress sdk.AccAddress) QueryPotRewardsByReportEpochParams {
	return QueryPotRewardsByReportEpochParams{
		Page:          page,
		Limit:         limit,
		Epoch:         epoch,
		WalletAddress: walletAddress,
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

//type ReportInfo struct {
//	Epoch     sdk.Int
//	Reference string
//}

func NewReportInfo(epoch sdk.Int, reference string) ReportInfo {
	return ReportInfo{
		Epoch:     epoch,
		Reference: reference,
	}
}
