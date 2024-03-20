package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stratosnet/stratos-chain/x/evm/vm"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

func (k Keeper) KeeGetTotalUnissuedPrepay(statedb vm.StateDB) (value sdkmath.Int) {
	totalUnissuedPrepayAccAddr := k.accountKeeper.GetModuleAddress(regtypes.TotalUnissuedPrepay)
	if totalUnissuedPrepayAccAddr == nil {
		value = sdkmath.ZeroInt()
	} else {
		value = sdkmath.NewIntFromBigInt(statedb.GetBalance(common.BytesToAddress(totalUnissuedPrepayAccAddr)))
	}

	return
}

func (k *Keeper) KeeCalculatePrepayPurchaseAmount(statedb vm.StateDB, amount sdkmath.Int) (sdkmath.Int, sdkmath.Int, error) {
	kdb := statedb.GetKeestateDB()
	St := k.registerKeeper.KeeGetEffectiveTotalDeposit(kdb)
	Pt := k.KeeGetTotalUnissuedPrepay(statedb)
	Lt := k.registerKeeper.KeeGetRemainingOzoneLimit(kdb)

	purchase := Lt.ToLegacyDec().
		Mul(amount.ToLegacyDec()).
		Quo((St.
			Add(Pt).
			Add(amount)).ToLegacyDec()).
		TruncateInt()
	if purchase.GT(Lt) {
		return sdkmath.NewInt(0), sdkmath.NewInt(0), fmt.Errorf("not enough remaining ozone limit to complete prepay")
	}

	remaining := Lt.Sub(purchase)

	return purchase, remaining, nil
}
