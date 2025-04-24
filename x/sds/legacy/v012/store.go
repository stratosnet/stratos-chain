package v012

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/sds/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey) error {
	store := ctx.KVStore(storeKey)

	removeFileInfo(store)

	return nil
}

func removeFileInfo(store storetypes.KVStore) {
	iter := sdk.KVStorePrefixIterator(store, types.FileStoreKeyPrefix)
	defer iter.Close()
	var keysToDelete [][]byte
	for ; iter.Valid(); iter.Next() {
		keysToDelete = append(keysToDelete, iter.Key())
	}

	for _, key := range keysToDelete {
		store.Delete(key)
	}
}
