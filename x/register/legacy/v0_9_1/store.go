package v0_9_1

import (
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/register/types"
)

// before: slashing 1utros -> 1ustos(handled as 1wei)
// after : slashing 1utros -> 1gwei = 1e9wei
func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.Codec, aminoCodec *codec.LegacyAmino) error {
	store := ctx.KVStore(storeKey)

	if err := migrateSlashingAmt(store, cdc, aminoCodec); err != nil {
		return err
	}

	return nil
}

func migrateSlashingAmt(store sdk.KVStore, cdc codec.Codec, aminoCodec *codec.LegacyAmino) error {
	oldSlashingStore := prefix.NewStore(store, SlashingPrefix)
	iterator := oldSlashingStore.Iterator(nil, nil)
	defer iterator.Close()

	var oldSlashingAmt sdk.Int

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		oldBz := iterator.Value()

		if oldBz == nil {
			oldSlashingStore.Delete(iterator.Key())
			continue
		}

		//old data was encoded by amino codec
		aminoCodec.MustUnmarshalLengthPrefixed(oldBz, &oldSlashingAmt)
		newSlashingAmt := oldSlashingAmt.Mul(sdk.NewInt(1e9))

		//use proto codec instead
		newBz := cdc.MustMarshalLengthPrefixed(&gogotypes.StringValue{Value: newSlashingAmt.String()})
		storeKey := types.GetSlashingKey(key)

		// slashing amount updated with new value
		oldSlashingStore.Delete(iterator.Key())
		store.Set(storeKey, newBz)
	}

	return nil
}
