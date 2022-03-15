package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
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
	TotalMinedTokensKey  = []byte{0x03}
	MinedTokensKeyPrefix = []byte{0x04} // key: prefix_epoch

	RewardAddressPoolKey         = []byte{0x11}
	LastReportedEpochKey         = []byte{0x12}
	IndividualRewardKeyPrefix    = []byte{0x13} // key: prefix{address}_individual_{epoch}, the amount that is matured at {epoch}
	MatureTotalRewardKeyPrefix   = []byte{0x14} // key: prefix{address}_mature_total
	ImmatureTotalRewardKeyPrefix = []byte{0x15} // key: prefix{address}_immature_total
	SlashingPrefix               = []byte{0x16} // key: prefix{address}_{epoch}
	// VolumeReportStoreKeyPrefix prefix for volumeReport store
	VolumeReportStoreKeyPrefix = []byte{0x41}
)

func GetMinedTokensKey(epoch sdk.Int) []byte {
	bEpoch := []byte(epoch.String())
	return append(MinedTokensKeyPrefix, bEpoch...)
}

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

func GetSlashingKey(p2pAddress stratos.SdsAddress) []byte {
	key := append(SlashingPrefix, p2pAddress...)
	return key
}
