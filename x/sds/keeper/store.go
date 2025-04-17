package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/crypto/merkle"

	"github.com/stratosnet/stratos-chain/x/sds/types"
)

// GetFileInfoByFileHash Returns the fileInfo
func (k Keeper) GetFileInfoByFileHash(ctx sdk.Context, fileHash []byte) (fileInfo types.FileInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetFileStoreKey(fileHash))
	if bz == nil {
		return fileInfo, false
	}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &fileInfo)
	return fileInfo, true
}

// Deprecated: not used for new files anymore
func (k Keeper) SetFileInfo(ctx sdk.Context, fileHash []byte, fileInfo types.FileInfo) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.GetFileStoreKey(fileHash)
	bz := k.cdc.MustMarshalLengthPrefixed(&fileInfo)
	store.Set(storeKey, bz)
}

// IterateFileInfo Iterate over all uploaded files.
// Iteration for all uploaded files
func (k Keeper) IterateFileInfo(ctx sdk.Context, handler func(string, types.FileInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.FileStoreKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		fileHash := string(iter.Key()[len(types.FileStoreKeyPrefix):])
		var fileInfo types.FileInfo
		k.cdc.MustUnmarshalLengthPrefixed(iter.Value(), &fileInfo)
		if handler(fileHash, fileInfo) {
			break
		}
	}
}

func (k Keeper) AddNewFile(ctx sdk.Context, fileHash []byte) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.GetNewFileKey(fileHash)
	store.Set(storeKey, []byte("1"))
}

// ClearNewFiles returns the list of new files uploaded in the current block, and removes them from the store
func (k Keeper) ClearNewFiles(ctx sdk.Context) []string {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.NewFileKeyPrefix)
	defer iter.Close()

	var newFiles []string
	var toDelete [][]byte
	for ; iter.Valid(); iter.Next() {
		fileHash := string(iter.Key()[len(types.NewFileKeyPrefix):])
		newFiles = append(newFiles, fileHash)
		toDelete = append(toDelete, iter.Key())
	}

	for _, key := range toDelete {
		store.Delete(key)
	}
	return newFiles
}

func (k Keeper) GetMerkleRoot(ctx sdk.Context) []byte {
	bz := ctx.KVStore(k.storeKey).Get(types.MerkleRootKeyPrefix)
	if bz == nil {
		return merkle.NullCommitment[:]
	}
	return bz
}

func (k Keeper) SetMerkleRoot(ctx sdk.Context, root []byte) {
	ctx.KVStore(k.storeKey).Set(types.MerkleRootKeyPrefix, root)
}

func (k Keeper) GetMerkleRootByHeight(ctx sdk.Context, height int64) []byte {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.GetPreviousMerkleRootKey(height)
	if store.Has(storeKey) {
		return store.Get(storeKey)
	} else {
		return merkle.NullCommitment[:]
	}
}

func (k Keeper) SetMerkleRootByHeight(ctx sdk.Context, height int64, root []byte) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.GetPreviousMerkleRootKey(height)
	store.Set(storeKey, root)
}
