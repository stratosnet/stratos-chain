package keeper

import (
	"fmt"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	//"github.com/cosmos/cosmos-sdk/x/supply"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// Keeper of the pot store
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	paramSpace       paramstypes.Subspace
	feeCollectorName string // name of the FeeCollector ModuleAccount
	BankKeeper       types.BankKeeper
	//SupplyKeeper     supply.Keeper
	AccountKeeper  types.AccountKeeper
	StakingKeeper  types.StakingKeeper
	RegisterKeeper types.RegisterKeeper
}

// NewKeeper creates a pot keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace paramstypes.Subspace, feeCollectorName string,
	bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper, stakingKeeper types.StakingKeeper,
	registerKeeper types.RegisterKeeper,
) Keeper {
	keeper := Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace.WithKeyTable(types.ParamKeyTable()),
		feeCollectorName: feeCollectorName,
		BankKeeper:       bankKeeper,
		//SupplyKeeper:     supplyKeeper,
		AccountKeeper:  accountKeeper,
		StakingKeeper:  stakingKeeper,
		RegisterKeeper: registerKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) VolumeReport(ctx sdk.Context, walletVolumes []*types.SingleWalletVolume, reporter stratos.SdsAddress,
	epoch sdk.Int, reportReference string, txHash string) (totalConsumedOzone sdk.Dec, err error) {
	//record volume report
	reportRecord := types.NewReportRecord(reporter, reportReference, txHash)
	k.SetVolumeReport(ctx, epoch, reportRecord)
	//distribute POT reward
	//TODO: recovery when shift to main net
	totalConsumedOzone, err = k.DistributePotReward(ctx, walletVolumes, epoch) // Main net
	//TODO: remove when shift to main net
	//totalConsumedOzone, err = k.DistributePotRewardForTestnet(ctx, walletVolumes, epoch) // Incentive test net

	return totalConsumedOzone, err
}

func (k Keeper) IsSPNode(ctx sdk.Context, p2pAddr stratos.SdsAddress) (found bool) {
	_, found = k.RegisterKeeper.GetIndexingNode(ctx, p2pAddr)
	return found
}

func (k Keeper) FoundationDeposit(ctx sdk.Context, amount sdk.Coins, from sdk.AccAddress) (err error) {
	err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, amount)
	if err != nil {
		return err
	}

	//TODO
	foundationAccountAddr := k.SupplyKeeper.GetModuleAddress(types.FoundationAccount)
	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, foundationAccountAddr, amount)
	if err != nil {
		return err
	}

	return nil

}
