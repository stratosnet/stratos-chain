package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/evm/types"
)

// GetParams returns the total set of evm parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the evm parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// ----------------------------------------------------------------------------
// Parent Base Fee
// Required by EIP1559 base fee calculation.
// ----------------------------------------------------------------------------

// GetBaseFeeParam get's the base fee from the paramSpace
// return nil if base fee is not enabled
func (k Keeper) GetBaseFeeParam(ctx sdk.Context) *big.Int {
	params := k.GetParams(ctx)
	if params.FeeMarketParams.NoBaseFee {
		return nil
	}

	return params.FeeMarketParams.BaseFee.BigInt()
}

// SetBaseFeeParam set's the base fee in the paramSpace
func (k Keeper) SetBaseFeeParam(ctx sdk.Context, baseFee *big.Int) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyBaseFee, sdk.NewIntFromBigInt(baseFee))
}
