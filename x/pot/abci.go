package pot

import (
	"bytes"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/pot/keeper"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// 	TODO: fill out if your application requires beginblock, if not you can delete this function
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) []abci.ValidatorUpdate {

	// Do not distribute rewards until the next block
	if !k.GetIsReadyToDistributeReward(ctx) && k.GetUnhandledEpoch(ctx).GT(sdk.ZeroInt()) {
		k.SetIsReadyToDistributeReward(ctx, true)
		return []abci.ValidatorUpdate{}
	}

	walletVolumes, found := k.GetUnhandledReport(ctx)
	if !found {
		return []abci.ValidatorUpdate{}
	}
	epoch := k.GetUnhandledEpoch(ctx)
	logger := k.Logger(ctx)

	//distribute POT reward
	_, err := k.DistributePotReward(ctx, walletVolumes.Volumes, epoch)
	if err != nil {
		logger.Error("An error occurred while distributing the reward. ", "ErrMsg", err.Error())
	}

	k.SetUnhandledReport(ctx, types.WalletVolumes{})
	k.SetUnhandledEpoch(ctx, sdk.ZeroInt())

	// reset total supply to 100M stos
	events := ctx.EventManager().Events()
	attrKeyAmtBytes := []byte(sdk.AttributeKeyAmount)

	totalBurnedAmount := sdk.Coins{}
	for _, event := range events {
		if event.Type == banktypes.EventTypeCoinBurn {
			attributes := event.Attributes
			for _, attr := range attributes {
				if bytes.Equal(attr.Key, attrKeyAmtBytes) {
					amount, err := sdk.ParseCoinsNormalized(string(attr.Value))
					if err != nil {
						logger.Error("An error occurred while parsing burned amount. ", "ErrMsg", err.Error())
						break
					}
					totalBurnedAmount = totalBurnedAmount.Add(amount...)
				}
			}
		}
	}

	logger.Debug("Total burned amount is:", totalBurnedAmount.String())

	// mintCoins + totalSupply should equal to 100M stos
	mintCoins := totalBurnedAmount
	totalSupplyLimit := sdk.NewInt(1e8)
	totalSupply := k.bankKeeper.GetSupply(ctx, k.BondDenom(ctx))
	newTotalSupply := totalBurnedAmount.AmountOf(k.BondDenom(ctx)).Add(totalSupply.Amount)
	if newTotalSupply.GT(totalSupplyLimit) {
		mintCoins = sdk.NewCoins(
			sdk.NewCoin(k.BondDenom(ctx), totalSupplyLimit.Sub(totalSupply.Amount)),
		)
	}

	k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins)

	senderAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	k.distrKeeper.FundCommunityPool(ctx, mintCoins, senderAddr)

	ctx.EventManager().EmitEvent(
		banktypes.NewCoinMintEvent(senderAddr, mintCoins),
	)

	return []abci.ValidatorUpdate{}
}
