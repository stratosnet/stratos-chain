package v011

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, legacySubspace types.ParamsSubspace, cdc codec.Codec) error {
	store := ctx.KVStore(storeKey)

	// migrate params
	if err := migrateParams(ctx, store, cdc, legacySubspace); err != nil {
		return err
	}

	// add 1M stos which should have been initialized from v0.11.0 genesis file
	if err := fixTotalMinedToken(ctx, store, cdc); err != nil {
		return err
	}

	return nil
}

// migrateParams will set the params to store from legacySubspace
func migrateParams(ctx sdk.Context, store storetypes.KVStore, cdc codec.Codec, legacySubspace types.ParamsSubspace) error {
	var legacyParams types.Params
	legacySubspace.GetParamSet(ctx, &legacyParams)

	if err := legacyParams.Validate(); err != nil {
		return err
	}

	bz := cdc.MustMarshal(&legacyParams)
	store.Set(types.ParamsKey, bz)
	return nil
}

func fixTotalMinedToken(_ sdk.Context, store storetypes.KVStore, cdc codec.Codec) error {
	var oldTotalMinedToken sdk.Coin

	oldBz := store.Get(types.TotalRewardKeyPrefix)
	if oldBz == nil {
		return fmt.Errorf("total mined token should not be empty")
	}
	cdc.MustUnmarshalLengthPrefixed(oldBz, &oldTotalMinedToken)

	initialTotalMinedAmount := sdkmath.NewInt(1e6).MulRaw(stratos.StosToWei)
	newTotalMinedToken := oldTotalMinedToken.AddAmount(initialTotalMinedAmount)
	newBz := cdc.MustMarshalLengthPrefixed(&newTotalMinedToken)

	store.Set(types.TotalRewardKeyPrefix, newBz)
	return nil
}
