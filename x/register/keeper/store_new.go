package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/stratosnet/stratos-chain/core/statedb"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

func (k *Keeper) KeeGetEffectiveTotalDeposit(kstatedb *statedb.KeestateDB) (stake sdkmath.Int) {
	bz := kstatedb.GetState(k.storeKey, types.EffectiveGenesisDepositTotalKey)
	if bz == nil {
		return sdkmath.ZeroInt()
	}
	intValue := stratos.Int{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intValue)
	stake = *intValue.Value
	return
}

func (k *Keeper) KeeGetRemainingOzoneLimit(kstatedb *statedb.KeestateDB) (value sdkmath.Int) {
	bz := kstatedb.GetState(k.storeKey, types.UpperBoundOfTotalOzoneKey)
	if bz == nil {
		return sdkmath.ZeroInt()
	}
	intVal := stratos.Int{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intVal)
	value = *intVal.Value
	return
}

func (k *Keeper) KeeSetRemainingOzoneLimit(kstatedb *statedb.KeestateDB, value sdkmath.Int) {
	bz := k.cdc.MustMarshalLengthPrefixed(&stratos.Int{Value: &value})
	kstatedb.SetState(k.storeKey, types.UpperBoundOfTotalOzoneKey, bz)
}
