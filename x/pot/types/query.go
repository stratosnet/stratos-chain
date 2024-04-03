package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryDefaultLimit = 100
)

// NewRewardByWallet creates a new instance of RewardByWallet
func NewRewardByWallet(
	walletAddress sdk.AccAddress,
	matureTotal sdk.Coins,
	immatureTotal sdk.Coins,
) *RewardByWallet {
	return &RewardByWallet{
		WalletAddress:       walletAddress.String(),
		MatureTotalReward:   matureTotal,
		ImmatureTotalReward: immatureTotal,
	}
}
