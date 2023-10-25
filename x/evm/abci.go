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
	err := ctx.EventManager().EmitTypedEvent(&types.EventFeeMarket{
		BaseFee: baseFee.String(),
	})
	if err != nil {
		k.Logger(ctx).Error(err.Error())
	}

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

	err := ctx.EventManager().EmitTypedEvent(&types.EventBlockGas{
		Height: fmt.Sprintf("%d", ctx.BlockHeight()),
		Amount: fmt.Sprintf("%d", ctx.BlockGasMeter().GasConsumedToLimit()),
	})
	if err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	// Gas costs are handled within msg handler so costs should be ignored
	infCtx := ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

	bloom := ethtypes.BytesToBloom(k.GetBlockBloomTransient(infCtx).Bytes())
	err = k.EmitBlockBloomEvent(infCtx, bloom)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	return []abci.ValidatorUpdate{}
}
