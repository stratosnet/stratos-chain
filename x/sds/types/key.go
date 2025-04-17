package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "sds"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName
)

var (
	FileStoreKeyPrefix          = []byte{0x01} // FileStorage prefix for sds store
	NewFileKeyPrefix            = []byte{0x02} // Prefix for the newly uploaded files in the current block
	MerkleRootKeyPrefix         = []byte{0x03} // Prefix for the current Merkle root
	PreviousMerkleRootKeyPrefix = []byte{0x04} // Prefix for the Merkle root of previous blocks

	ParamsKey = []byte{0x20}
)

// GetFileStoreKey turns a file hash to a key used to get the file info from the sds store
func GetFileStoreKey(fileHash []byte) []byte {
	return append(FileStoreKeyPrefix, fileHash...)
}

// GetNewFileKey turns a file hash to a key used to get the new file uploaded in the current block from the sds store
func GetNewFileKey(fileHash []byte) []byte {
	return append(NewFileKeyPrefix, fileHash...)
}

// GetPreviousMerkleRootKey gets the previous Merkle root key associated with a given height
func GetPreviousMerkleRootKey(height int64) []byte {
	return append(PreviousMerkleRootKeyPrefix, sdk.Uint64ToBigEndian(uint64(height))...)
}
