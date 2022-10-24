package keeper

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/rpc/core"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/stratosnet/stratos-chain/x/evm/types"
)

var (
	overriddenTxHashes map[string]struct{}
)

// BeginBlock sets the sdk Context and EIP155 chain id to the Keeper.
func (k *Keeper) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	baseFee := k.CalculateBaseFee(ctx)

	// return immediately if base fee is nil
	if baseFee == nil {
		return
	}

	k.SetBaseFeeParam(ctx, baseFee)

	if !ctx.IsCheckTx() && len(ctx.HeaderHash()) > 0 {
		overriddenTxHashes = make(map[string]struct{}, 0)
		k.Logger(ctx).Info(fmt.Sprintf("evm_keeper_begin_block: len(ctx.HeaderHash().Bytes()) = %v", len(ctx.HeaderHash().Bytes())))
		k.Logger(ctx).Info(fmt.Sprintf("evm_keeper_begin_block: ctx.HeaderHash().Bytes() = = %v", common.Bytes2Hex(ctx.HeaderHash().Bytes())))
		block, err := core.BlockByHash(&rpctypes.Context{}, ctx.HeaderHash().Bytes())
		rawTxs := block.Block.Txs
		if err != nil {
			panic(err)
		}

		for _, rawTx := range rawTxs {
			tx, err := k.txDecoder(rawTx)
			if err != nil {
				panic(err)
			}

			for _, msg := range tx.GetMsgs() {
				msgEthTx, ok := msg.(*types.MsgEthereumTx)
				if !ok {
					continue
				}

				txData, err := types.UnpackTxData(msgEthTx.Data)
				if err != nil {
					panic(err)
				}

				from := msgEthTx.GetFrom()
				nonce := txData.GetNonce()
				gasPrice := txData.GetGasPrice()

				for _, rawTx2 := range rawTxs {
					tx2, err := k.txDecoder(rawTx2)
					if err != nil {
						continue
					}

					for _, msg2 := range tx2.GetMsgs() {
						msgEthTx2, ok := msg2.(*types.MsgEthereumTx)
						if !ok {
							continue
						}

						txData2, err := types.UnpackTxData(msgEthTx2.Data)
						if err != nil {
							panic(err)
						}

						if from.Equals(msgEthTx2.GetFrom()) && nonce == txData2.GetNonce() && gasPrice.Cmp(txData2.GetGasPrice()) < 0 {
							// find tx has same nonce with higher gas price from same sender,
							// record tx hash to overriddenTxHashes, let it fail in anteHandler
							// overriddenTxHashes need to be cleared in the EndBlock function
							overriddenTxHashes[msgEthTx.Hash] = struct{}{}
							break
						}
					}
				}

			}
		}
	}

	// Store current base fee in event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFeeMarket,
			sdk.NewAttribute(types.AttributeKeyBaseFee, baseFee.String()),
		),
	})
}

// EndBlock also retrieves the bloom filter value from the transient store and commits it to the
// KVStore. The EVM end block logic doesn't update the validator set, thus it returns
// an empty slice.
func (k *Keeper) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	if ctx.BlockGasMeter() == nil {
		k.Logger(ctx).Error("block gas meter is nil when setting block gas used")
		panic("block gas meter is nil when setting block gas used")
	}
	overriddenTxHashes = nil
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

func (k *Keeper) GetOverriddenTxHashMap(ctx sdk.Context) map[string]struct{} {
	if ctx.IsCheckTx() {
		return nil
	}
	return overriddenTxHashes
}
