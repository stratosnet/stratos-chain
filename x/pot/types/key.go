package types

import sdk "github.com/cosmos/cosmos-sdk/types"

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
	// VolumeReportStoreKeyPrefix prefix for volumeReport store
	VolumeReportStoreKeyPrefix     = []byte{0x41}
	SingleNodeVolumeStoreKeyPrefix = []byte{0x51}
)

// VolumeReportStoreKey turns an address to key used to get it from the account store
func VolumeReportStoreKey(reporter sdk.AccAddress) []byte {
	return append(VolumeReportStoreKeyPrefix, reporter.Bytes()...)
}

// SingleNodeVolumeStoreKey turns an address to key used to get it from the account store
func SingleNodeVolumeStoreKey(nodeAddress sdk.AccAddress) []byte {
	return append(SingleNodeVolumeStoreKeyPrefix, nodeAddress.Bytes()...)
}
