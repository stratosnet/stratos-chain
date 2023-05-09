package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/core/statedb"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

func (k Keeper) KeeGetEffectiveTotalStake(kstatedb *statedb.KeestateDB) (stake sdk.Int) {
	bz := kstatedb.GetState(k.storeKey, types.EffectiveGenesisStakeTotalKey)
	if bz == nil {
		return sdk.ZeroInt()
	}
	intValue := stratos.Int{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intValue)
	stake = *intValue.Value
	return
}

func (k Keeper) KeeGetTotalUnissuedPrepay(kstatedb *statedb.KeestateDB) (value sdk.Int) {
	// TODO
	value = sdk.ZeroInt()
	return
}

func (k Keeper) KeeGetRemainingOzoneLimit(kstatedb *statedb.KeestateDB) (value sdk.Int) {
	bz := kstatedb.GetState(k.storeKey, types.UpperBoundOfTotalOzoneKey)
	if bz == nil {
		return sdk.ZeroInt()
	}
	intVal := stratos.Int{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intVal)
	value = *intVal.Value
	return
}

func (k *Keeper) KeeCalculatePurchaseAmount(kstatedb *statedb.KeestateDB, amount sdk.Int) (sdk.Int, sdk.Int, error) {
	St := k.KeeGetEffectiveTotalStake(kstatedb)
	Pt := k.KeeGetTotalUnissuedPrepay(kstatedb)
	Lt := k.KeeGetRemainingOzoneLimit(kstatedb)

	fmt.Println("KeeCalculatePurchaseAmount St", St)
	fmt.Println("KeeCalculatePurchaseAmount Pt", Pt)
	fmt.Println("KeeCalculatePurchaseAmount Lt", Lt)

	purchase := Lt.ToDec().
		Mul(amount.ToDec()).
		Quo((St.
			Add(Pt).
			Add(amount)).ToDec()).
		TruncateInt()
	fmt.Println("KeeCalculatePurchaseAmount purchase", purchase)
	if purchase.GT(Lt) {
		return sdk.NewInt(0), sdk.NewInt(0), fmt.Errorf("not enough remaining ozone limit to complete prepay")
	}
	remaining := Lt.Sub(purchase)
	fmt.Println("KeeCalculatePurchaseAmount remaining", remaining)

	return purchase, remaining, nil
}

func (k Keeper) KeeSetRemainingOzoneLimit(kstatedb *statedb.KeestateDB, value sdk.Int) {
	bz := k.cdc.MustMarshalLengthPrefixed(&stratos.Int{Value: &value})
	kstatedb.SetState(k.storeKey, types.UpperBoundOfTotalOzoneKey, bz)
}
