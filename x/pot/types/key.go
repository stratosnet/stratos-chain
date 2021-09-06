package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "pot"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName
)

var (
	FoundationAccountKey   = []byte{0x01}
	InitialUOzonePriceKey  = []byte{0x02}
	TotalMinedTokensKey    = []byte{0x03}
	MinedTokensKeyPrefix   = []byte{0x04} // key: prefix_epoch
	TotalUnissuedPrepayKey = []byte{0x05}

	RewardAddressPoolKey         = []byte{0x11}
	LastMaturedEpochKey          = []byte{0x12}
	IndividualRewardKeyPrefix    = []byte{0x13} // key: prefix{address}_individual_{epoch}, the amount that is matured at {epoch}
	MatureTotalRewardKeyPrefix   = []byte{0x14} // key: prefix{address}_mature_total
	ImmatureTotalRewardKeyPrefix = []byte{0x15} // key: prefix{address}_immature_total

	EpochRewardsKeyPrefix = []byte{0x51} // key: prefix{epoch}, value: rewardDetailMap (map[node_address string]types.Reward)

	// VolumeReportStoreKeyPrefix prefix for volumeReport store
	VolumeReportStoreKeyPrefix = []byte{0x41}
)

func GetMinedTokensKey(epoch sdk.Int) []byte {
	bEpoch := []byte(epoch.String())
	return append(MinedTokensKeyPrefix, bEpoch...)
}

// VolumeReportStoreKey turns an address to key used to get it from the account store
//func VolumeReportStoreKey(reporter sdk.AccAddress) []byte {
//	return append(VolumeReportStoreKeyPrefix, reporter.Bytes()...)
//}
func VolumeReportStoreKey(epoch sdk.Int) []byte {
	return append(VolumeReportStoreKeyPrefix, epoch.String()...)
}

// GetIndividualRewardKey prefix{address}_individual_{epoch}, the amount that is matured at {epoch}
func GetIndividualRewardKey(acc sdk.AccAddress, epoch sdk.Int) []byte {
	bKeyStr := []byte("_individual_")
	bEpoch := []byte(epoch.String())

	key := append(IndividualRewardKeyPrefix, acc...)
	key = append(key, bKeyStr...)
	key = append(key, bEpoch...)
	return key
}

// GetMatureTotalRewardKey prefix{address}_mature_total
func GetMatureTotalRewardKey(acc sdk.AccAddress) []byte {
	bKeyStr := []byte("_mature_total")
	key := append(MatureTotalRewardKeyPrefix, acc.Bytes()...)
	key = append(key, bKeyStr...)
	return key
}

// GetImmatureTotalRewardKey prefix{address}_immature_total
func GetImmatureTotalRewardKey(acc sdk.AccAddress) []byte {
	bKeyStr := []byte("_immature_total")
	key := append(ImmatureTotalRewardKeyPrefix, acc.Bytes()...)
	key = append(key, bKeyStr...)
	return key
}

// GetEpochRewardsKey prefix{epoch}_rewards
func GetEpochRewardsKey(epoch sdk.Int) []byte {
	bKeyStr := []byte("_rewards")
	bEpoch := []byte(epoch.String())
	key := append(EpochRewardsKeyPrefix, bEpoch...)
	key = append(key, bKeyStr...)
	return key
}
