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
	TotalMinedTokensKey          = []byte{0x01}
	LastReportedEpochKey         = []byte{0x02}
	IndividualRewardKeyPrefix    = []byte{0x03} // key: prefix{address}_{epoch}, the amount that is matured at {epoch}
	MatureTotalRewardKeyPrefix   = []byte{0x04} // key: prefix{address}
	ImmatureTotalRewardKeyPrefix = []byte{0x05} // key: prefix{address}
	VolumeReportStoreKeyPrefix   = []byte{0x06} // VolumeReportStoreKeyPrefix prefix for volumeReport store
	MaturedEpochKeyPrefix        = []byte{0x07}
)

func VolumeReportStoreKey(epoch sdk.Int) []byte {
	return append(VolumeReportStoreKeyPrefix, epoch.String()...)
}

// GetIndividualRewardKey prefix{epoch}_{account}, the amount that is matured at {epoch}
func GetIndividualRewardKey(acc sdk.AccAddress, epoch sdk.Int) []byte {
	bKeyStr := []byte("_")
	bEpoch := []byte(epoch.String())

	key := append(IndividualRewardKeyPrefix, bEpoch...)
	key = append(key, bKeyStr...)
	key = append(key, acc...)
	return key
}

func GetIndividualRewardIteratorKey(epoch sdk.Int) []byte {
	bKeyStr := []byte("_")
	bEpoch := []byte(epoch.String())
	key := append(IndividualRewardKeyPrefix, bEpoch...)
	key = append(key, bKeyStr...)
	return key
}

// GetMatureTotalRewardKey prefix{address}
func GetMatureTotalRewardKey(acc sdk.AccAddress) []byte {
	key := append(MatureTotalRewardKeyPrefix, acc.Bytes()...)
	return key
}

// GetImmatureTotalRewardKey prefix{address}
func GetImmatureTotalRewardKey(acc sdk.AccAddress) []byte {
	key := append(ImmatureTotalRewardKeyPrefix, acc.Bytes()...)
	return key
}
