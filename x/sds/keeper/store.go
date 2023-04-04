package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/sds/types"
)

// GetFileInfoByFileHash Returns the fileInfo
func (k Keeper) GetFileInfoByFileHash(ctx sdk.Context, fileHash []byte) (fileInfo types.FileInfo, found bool) {
	store := ctx.KVStore(k.key)
	bz := store.Get(types.FileStoreKey(fileHash))
	if bz == nil {
		return fileInfo, false
	}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &fileInfo)
	return fileInfo, true
}

func (k Keeper) SetFileInfo(ctx sdk.Context, fileHash []byte, fileInfo types.FileInfo) {
	store := ctx.KVStore(k.key)
	storeKey := types.FileStoreKey(fileHash)
	bz := k.cdc.MustMarshalLengthPrefixed(&fileInfo)
	store.Set(storeKey, bz)
}

// IterateFileInfo Iterate over all uploaded files.
// Iteration for all uploaded files
func (k Keeper) IterateFileInfo(ctx sdk.Context, handler func(string, types.FileInfo) (stop bool)) {
	store := ctx.KVStore(k.key)
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
