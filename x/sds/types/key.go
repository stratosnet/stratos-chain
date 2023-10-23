package types

const (
	// ModuleName is the name of the module
	ModuleName = "sds"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName
)

var (
	FileStoreKeyPrefix = []byte{0x01} // FileStorage prefix for sds store

	ParamsKey = []byte{0x20}
)

// GetFileStoreKey turn an address to key used to get it from the account store
func GetFileStoreKey(sender []byte) []byte {
	return append(FileStoreKeyPrefix, sender...)
}
