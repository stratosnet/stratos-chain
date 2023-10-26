package keeper

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/evm/types"
)

// SetParams sets the params on the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.KeyPrefixParams, bz)

	return nil
}

// GetParams returns the params from the store
func (k Keeper) GetParams(ctx sdk.Context) (p types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixParams)
	if bz == nil {
		return p
	}
	k.cdc.MustUnmarshal(bz, &p)
	return p
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
	params := k.GetParams(ctx)
	params.FeeMarketParams.BaseFee = sdkmath.NewIntFromBigInt(baseFee)
	k.SetParams(ctx, params)
}
