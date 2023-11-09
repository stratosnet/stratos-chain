package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryDefaultLimit = 100
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
