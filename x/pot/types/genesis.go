package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, totalMinedToken sdk.Coin, lastDistributedEpoch sdkmath.Int,
	immatureTotalInfo []ImmatureTotal, matureTotalInfo []MatureTotal, individualRewardInfo []Reward,
	maturedEpoch sdkmath.Int, rewardTotalInfo []RewardTotal,
) *GenesisState {

	return &GenesisState{
		Params:               params,
		TotalMinedToken:      totalMinedToken,
		LastDistributedEpoch: lastDistributedEpoch,
		ImmatureTotalInfo:    immatureTotalInfo,
		MatureTotalInfo:      matureTotalInfo,
		IndividualRewardInfo: individualRewardInfo,
		MaturedEpoch:         maturedEpoch,
		RewardTotalInfo:      rewardTotalInfo,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	params := DefaultParams()
	coin := sdk.NewCoin(DefaultRewardDenom, sdkmath.ZeroInt())
	return &GenesisState{
		Params:               params,
		TotalMinedToken:      coin,
		LastDistributedEpoch: sdkmath.ZeroInt(),
		ImmatureTotalInfo:    make([]ImmatureTotal, 0),
		MatureTotalInfo:      make([]MatureTotal, 0),
		IndividualRewardInfo: make([]Reward, 0),
		MaturedEpoch:         sdkmath.ZeroInt(),
		RewardTotalInfo:      make([]RewardTotal, 0),
	}
}

// ValidateGenesis validates the pot genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}

func NewImmatureTotal(walletAddress sdk.AccAddress, value sdk.Coins) ImmatureTotal {
	return ImmatureTotal{
		WalletAddress: walletAddress.String(),
		Value:         value,
	}
}

func NewMatureTotal(walletAddress sdk.AccAddress, value sdk.Coins) MatureTotal {
	return MatureTotal{
		WalletAddress: walletAddress.String(),
		Value:         value,
	}
}

func NewRewardTotal(epoch sdkmath.Int, reward TotalReward) RewardTotal {
	return RewardTotal{
		Epoch:       epoch,
		TotalReward: reward,
	}
}
