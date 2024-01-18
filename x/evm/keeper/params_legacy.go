package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stratosnet/stratos-chain/x/evm/types"
)

// GetParams returns the total set of evm parameters.
func (k Keeper) GetParamsV011(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// ----------------------------------------------------------------------------
// Parent Base Fee
// Required by EIP1559 base fee calculation.
// ----------------------------------------------------------------------------

// GetBaseFeeParam get's the base fee from the paramSpace
// return nil if base fee is not enabled
func (k Keeper) GetBaseFeeParamV011(ctx sdk.Context) *big.Int {
	params := k.GetParams(ctx)
	if params.FeeMarketParams.NoBaseFee {
		return nil
	}

	return params.FeeMarketParams.BaseFee.BigInt()
}

func (k Keeper) GetBaseFeeV011(ctx sdk.Context, ethCfg *params.ChainConfig) *big.Int {
	if !types.IsLondon(ethCfg, ctx.BlockHeight()) {
		return nil
	}
	baseFee := k.GetBaseFeeParamV011(ctx)
	if baseFee == nil {
		// return 0 if feemarket not enabled.
		baseFee = big.NewInt(0)
	}
	return baseFee
}
