package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/crypto/merkle"

	"github.com/stratosnet/stratos-chain/x/sds/types"
)

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

// IterateMerkleRoots Iterate over the merkle roots of previous heights.
func (k Keeper) IterateMerkleRoots(ctx sdk.Context, handler func(int64, []byte) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.PreviousMerkleRootKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		height := int64(sdk.BigEndianToUint64(iter.Key()[len(types.PreviousMerkleRootKeyPrefix):]))
		root := iter.Value()
		if handler(height, root) {
			break
		}
	}
}
