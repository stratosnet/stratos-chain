package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	//stratos "github.com/stratosnet/stratos-chain/types"
)

//type GenesisState struct {
//	Params               Params          `json:"params" yaml:"params"`
//	TotalMinedToken      sdk.Coin        `json:"total_mined_token" yaml:"total_mined_token"`
//	LastReportedEpoch    int64           `json:"last_reported_epoch" yaml:"last_reported_epoch"`
//	ImmatureTotalInfo    []ImmatureTotal `json:"immature_total_info" yaml:"immature_total_info"`
//	MatureTotalInfo      []MatureTotal   `json:"mature_total_info" yaml:"mature_total_info"`
//	IndividualRewardInfo []Reward        `json:"individual_reward_info" yaml:"individual_reward_info"`
//}
//
// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, totalMinedToken sdk.Coin, lastReportedEpoch int64,
	immatureTotalInfo []*ImmatureTotal, matureTotalInfo []*MatureTotal, individualRewardInfo []*Reward,
) GenesisState {

	return GenesisState{
		Params:               &params,
		TotalMinedToken:      &totalMinedToken,
		LastReportedEpoch:    lastReportedEpoch,
		ImmatureTotalInfo:    immatureTotalInfo,
		MatureTotalInfo:      matureTotalInfo,
		IndividualRewardInfo: individualRewardInfo,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	params := DefaultParams()
	coin := sdk.NewCoin(DefaultRewardDenom, sdk.ZeroInt())
	return GenesisState{
		Params:               &params,
		TotalMinedToken:      &coin,
		LastReportedEpoch:    0,
		ImmatureTotalInfo:    make([]*ImmatureTotal, 0),
		MatureTotalInfo:      make([]*MatureTotal, 0),
		IndividualRewardInfo: make([]*Reward, 0),
	}
}

// ValidateGenesis validates the pot genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}

//type ImmatureTotal struct {
//	WalletAddress sdk.AccAddress `json:"wallet_address" yaml:"wallet_address"`
//	Value         sdk.Coins      `json:"value" yaml:"value"`
//}
//
func NewImmatureTotal(walletAddress sdk.AccAddress, value sdk.Coins) ImmatureTotal {
	return ImmatureTotal{
		WalletAddress: walletAddress.String(),
		Value:         value,
	}
}

//type MatureTotal struct {
//	WalletAddress sdk.AccAddress `json:"wallet_address" yaml:"wallet_address"`
//	Value         sdk.Coins      `json:"value" yaml:"value"`
//}
//
func NewMatureTotal(walletAddress sdk.AccAddress, value sdk.Coins) MatureTotal {
	return MatureTotal{
		WalletAddress: walletAddress.String(),
		Value:         value,
	}
}
