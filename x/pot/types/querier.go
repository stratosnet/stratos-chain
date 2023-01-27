package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewPotRewardInfo creates a new instance of PotRewardInfo
func NewPotRewardInfo(
	walletAddress sdk.AccAddress,
	matureTotal sdk.Coins,
	immatureTotal sdk.Coins,
) PotRewardByOwner {
	return PotRewardByOwner{
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
