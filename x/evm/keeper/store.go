package keeper

import (
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
	St := k.registerKeeper.KeeGetEffectiveTotalDeposit(kdb)
	Pt := k.KeeGetTotalUnissuedPrepay(statedb)
	Lt := k.registerKeeper.KeeGetRemainingOzoneLimit(kdb)

	purchasedNoz, remainingNoz, err := k.potKeeper.GetPrepayAmount(Lt, amount, St, Pt)
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt(), err
	}

	return purchasedNoz, remainingNoz, nil
}
