package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	//stratos "github.com/stratosnet/stratos-chain/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, totalMinedToken sdk.Coin, lastReportedEpoch sdk.Int,
	immatureTotalInfo []ImmatureTotal, matureTotalInfo []MatureTotal, individualRewardInfo []Reward,
	maturedEpoch sdk.Int,
) *GenesisState {

	return &GenesisState{
		Params:               params,
		TotalMinedToken:      totalMinedToken,
		LastReportedEpoch:    lastReportedEpoch,
		ImmatureTotalInfo:    immatureTotalInfo,
		MatureTotalInfo:      matureTotalInfo,
		IndividualRewardInfo: individualRewardInfo,
		MaturedEpoch:         maturedEpoch,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	params := DefaultParams()
	coin := sdk.NewCoin(DefaultRewardDenom, sdk.ZeroInt())
	return &GenesisState{
		Params:               params,
		TotalMinedToken:      coin,
		LastReportedEpoch:    sdk.ZeroInt(),
		ImmatureTotalInfo:    make([]ImmatureTotal, 0),
		MatureTotalInfo:      make([]MatureTotal, 0),
		IndividualRewardInfo: make([]Reward, 0),
		MaturedEpoch:         sdk.ZeroInt(),
	}
}

// ValidateGenesis validates the pot genesis parameters
func ValidateGenesis(data *GenesisState) error {
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
