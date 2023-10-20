package evm

import (
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	abci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/evm/keeper"
	"github.com/stratosnet/stratos-chain/x/evm/types"
)

// BeginBlocker sets the sdk Context and EIP155 chain id to the Keeper.
func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k *keeper.Keeper) {
	baseFee := k.CalculateBaseFee(ctx)

	// return immediately if base fee is nil
	if baseFee == nil {
		return
	}

	k.SetBaseFeeParam(ctx, baseFee)

	// Store current base fee in event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFeeMarket,
			sdk.NewAttribute(types.AttributeKeyBaseFee, baseFee.String()),
		),
	})
}

// EndBlocker also retrieves the bloom filter value from the transient store and commits it to the
// KVStore. The EVM end block logic doesn't update the validator set, thus it returns
// an empty slice.
func EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock, k *keeper.Keeper) []abci.ValidatorUpdate {
	if ctx.BlockGasMeter() == nil {
		k.Logger(ctx).Error("block gas meter is nil when setting block gas used")
		panic("block gas meter is nil when setting block gas used")
	}

	gasUsed := ctx.BlockGasMeter().GasConsumedToLimit()

	k.SetBlockGasUsed(ctx, gasUsed)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"block_gas",
		sdk.NewAttribute("height", fmt.Sprintf("%d", ctx.BlockHeight())),
		sdk.NewAttribute("amount", fmt.Sprintf("%d", ctx.BlockGasMeter().GasConsumedToLimit())),
	))

	// Gas costs are handled within msg handler so costs should be ignored
	infCtx := ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

	bloom := ethtypes.BytesToBloom(k.GetBlockBloomTransient(infCtx).Bytes())
	k.EmitBlockBloomEvent(infCtx, bloom)

	return []abci.ValidatorUpdate{}
}
