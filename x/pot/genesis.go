package pot

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)

	foundationAccountAddr := keeper.SupplyKeeper.GetModuleAddress(types.FoundationAccount)
	err := keeper.BankKeeper.SetCoins(ctx, foundationAccountAddr, sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), data.FoundationAccountBalance)))
	if err != nil {
		return
	}

	keeper.SetInitialUOzonePrice(ctx, data.InitialUozPrice)
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data types.GenesisState) {
	params := keeper.GetParams(ctx)

	foundationAccountAddr := keeper.SupplyKeeper.GetModuleAddress(types.FoundationAccount)
	foundationAccountBalance := keeper.BankKeeper.GetCoins(ctx, foundationAccountAddr).AmountOf(keeper.BondDenom(ctx))

	initialUOzonePrice := keeper.GetInitialUOzonePrice(ctx)

	return types.NewGenesisState(params, foundationAccountBalance, initialUOzonePrice)
}
