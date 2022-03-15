package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// Keeper of the pot store
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	paramSpace       params.Subspace
	feeCollectorName string // name of the FeeCollector ModuleAccount
	BankKeeper       bank.Keeper
	SupplyKeeper     supply.Keeper
	AccountKeeper    auth.AccountKeeper
	StakingKeeper    staking.Keeper
	RegisterKeeper   register.Keeper
}

// NewKeeper creates a pot keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, feeCollectorName string,
	bankKeeper bank.Keeper, supplyKeeper supply.Keeper, accountKeeper auth.AccountKeeper, stakingKeeper staking.Keeper,
	registerKeeper register.Keeper,
) Keeper {
	keeper := Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace.WithKeyTable(types.ParamKeyTable()),
		feeCollectorName: feeCollectorName,
		BankKeeper:       bankKeeper,
		SupplyKeeper:     supplyKeeper,
		AccountKeeper:    accountKeeper,
		StakingKeeper:    stakingKeeper,
		RegisterKeeper:   registerKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) VolumeReport(ctx sdk.Context, walletVolumes []types.SingleWalletVolume, reporter stratos.SdsAddress,
	epoch sdk.Int, reportReference string, txHash string) (totalConsumedOzone sdk.Dec, err error) {
	//record volume report
	reportRecord := types.NewReportRecord(reporter, reportReference, txHash)
	k.SetVolumeReport(ctx, epoch, reportRecord)
	//distribute POT reward
	//TODO: recovery when shift to main net
	//totalConsumedOzone, err = k.DistributePotReward(ctx, walletVolumes, epoch) // Main net
	//TODO: remove when shift to main net
	totalConsumedOzone, err = k.DistributePotRewardForTestnet(ctx, walletVolumes, epoch) // Incentive test net

	return totalConsumedOzone, err
}

func (k Keeper) IsSPNode(ctx sdk.Context, p2pAddr stratos.SdsAddress) (found bool) {
	_, found = k.RegisterKeeper.GetIndexingNode(ctx, p2pAddr)
	return found
}

func (k Keeper) FoundationDeposit(ctx sdk.Context, amount sdk.Coins, from sdk.AccAddress) (err error) {
	_, err = k.BankKeeper.SubtractCoins(ctx, from, amount)
	if err != nil {
		return err
	}

	foundationAccountAddr := k.SupplyKeeper.GetModuleAddress(types.FoundationAccount)
	_, err = k.BankKeeper.AddCoins(ctx, foundationAccountAddr, amount)
	if err != nil {
		return err
	}

	return nil

}
