package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/core/statedb"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

func (k *Keeper) KeeGetEffectiveTotalDeposit(kstatedb *statedb.KeestateDB) (stake sdk.Int) {
	bz := kstatedb.GetState(k.storeKey, types.EffectiveGenesisDepositTotalKey)
	if bz == nil {
		return sdk.ZeroInt()
	}
	intValue := stratos.Int{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intValue)
	stake = *intValue.Value
	return
}

func (k *Keeper) KeeGetRemainingOzoneLimit(kstatedb *statedb.KeestateDB) (value sdk.Int) {
	bz := kstatedb.GetState(k.storeKey, types.UpperBoundOfTotalOzoneKey)
	if bz == nil {
		return sdk.ZeroInt()
	}
	intVal := stratos.Int{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intVal)
	value = *intVal.Value
	return
}

func (k *Keeper) KeeSetRemainingOzoneLimit(kstatedb *statedb.KeestateDB, value sdk.Int) {
	bz := k.cdc.MustMarshalLengthPrefixed(&stratos.Int{Value: &value})
	kstatedb.SetState(k.storeKey, types.UpperBoundOfTotalOzoneKey, bz)
}
