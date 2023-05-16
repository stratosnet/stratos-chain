package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stratosnet/stratos-chain/x/evm/vm"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

func (k Keeper) KeeGetTotalUnissuedPrepay(statedb vm.StateDB) (value sdk.Int) {
	totalUnissuedPrepayAccAddr := k.accountKeeper.GetModuleAddress(regtypes.TotalUnissuedPrepay)
	if totalUnissuedPrepayAccAddr == nil {
		value = sdk.ZeroInt()
	} else {
		value = sdk.NewIntFromBigInt(statedb.GetBalance(common.BytesToAddress(totalUnissuedPrepayAccAddr)))
	}

	return
}

func (k *Keeper) KeeCalculatePrepayPurchaseAmount(statedb vm.StateDB, amount sdk.Int) (sdk.Int, sdk.Int, error) {
	kdb := statedb.GetKeestateDB()
	St := k.registerKeeper.KeeGetEffectiveTotalStake(kdb)
	Pt := k.KeeGetTotalUnissuedPrepay(statedb)
	Lt := k.registerKeeper.KeeGetRemainingOzoneLimit(kdb)

	purchase := Lt.ToDec().
		Mul(amount.ToDec()).
		Quo((St.
			Add(Pt).
			Add(amount)).ToDec()).
		TruncateInt()
	if purchase.GT(Lt) {
		return sdk.NewInt(0), sdk.NewInt(0), fmt.Errorf("not enough remaining ozone limit to complete prepay")
	}

	remaining := Lt.Sub(purchase)

	return purchase, remaining, nil
}
